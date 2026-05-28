package appserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opencsv/opencsv/pkg/commands"
	"github.com/opencsv/opencsv/pkg/csvcore"
	"github.com/opencsv/opencsv/pkg/csvdiff"
	"github.com/opencsv/opencsv/pkg/dataprofile"
	docstore "github.com/opencsv/opencsv/pkg/document"
	"github.com/opencsv/opencsv/pkg/sqlconsole"
)

type Server struct {
	mux       *http.ServeMux
	webDir    string
	documents *docstore.Manager
}

func New(webDir string) *Server {
	server := &Server{
		mux:       http.NewServeMux(),
		webDir:    webDir,
		documents: docstore.NewManager(),
	}
	server.routes()
	return server
}

func (server *Server) Handler() http.Handler {
	return server.mux
}

func (server *Server) routes() {
	server.mux.HandleFunc("GET /api/health", server.handleHealth)
	server.mux.HandleFunc("GET /api/commands", server.handleCommands)
	server.mux.HandleFunc("POST /api/open", server.handleOpen)
	server.mux.HandleFunc("POST /api/file-sample", server.handleFileSample)
	server.mux.HandleFunc("POST /api/upload-open", server.handleUploadOpen)
	server.mux.HandleFunc("GET /api/window", server.handleWindow)
	server.mux.HandleFunc("POST /api/edit", server.handleEdit)
	server.mux.HandleFunc("POST /api/edit-batch", server.handleEditBatch)
	server.mux.HandleFunc("POST /api/undo", server.handleUndo)
	server.mux.HandleFunc("POST /api/redo", server.handleRedo)
	server.mux.HandleFunc("POST /api/header-mode", server.handleHeaderMode)
	server.mux.HandleFunc("POST /api/quote-policy", server.handleQuotePolicy)
	server.mux.HandleFunc("POST /api/read-only", server.handleReadOnly)
	server.mux.HandleFunc("GET /api/file-status", server.handleFileStatus)
	server.mux.HandleFunc("POST /api/refresh", server.handleRefresh)
	server.mux.HandleFunc("POST /api/row/insert", server.handleInsertRow)
	server.mux.HandleFunc("POST /api/row/delete", server.handleDeleteRow)
	server.mux.HandleFunc("POST /api/row/move", server.handleMoveRow)
	server.mux.HandleFunc("POST /api/row/hide", server.handleHideRow)
	server.mux.HandleFunc("POST /api/column/insert", server.handleInsertColumn)
	server.mux.HandleFunc("POST /api/column/delete", server.handleDeleteColumn)
	server.mux.HandleFunc("POST /api/column/move", server.handleMoveColumn)
	server.mux.HandleFunc("POST /api/transpose", server.handleTranspose)
	server.mux.HandleFunc("POST /api/dedupe-rows", server.handleDedupeRows)
	server.mux.HandleFunc("POST /api/dedupe-columns", server.handleDedupeColumns)
	server.mux.HandleFunc("POST /api/text-transform", server.handleTextTransform)
	server.mux.HandleFunc("POST /api/cells/swap", server.handleSwapCells)
	server.mux.HandleFunc("POST /api/cells/join", server.handleJoinCells)
	server.mux.HandleFunc("POST /api/cells/split", server.handleSplitCells)
	server.mux.HandleFunc("POST /api/save", server.handleSave)
	server.mux.HandleFunc("POST /api/export-view", server.handleExportView)
	server.mux.HandleFunc("POST /api/find", server.handleFind)
	server.mux.HandleFunc("POST /api/replace", server.handleReplace)
	server.mux.HandleFunc("POST /api/filter", server.handleFilter)
	server.mux.HandleFunc("GET /api/column-values", server.handleColumnValues)
	server.mux.HandleFunc("POST /api/filter-values", server.handleFilterValues)
	server.mux.HandleFunc("POST /api/sort", server.handleSort)
	server.mux.HandleFunc("POST /api/view-clear", server.handleViewClear)
	server.mux.HandleFunc("POST /api/diff", server.handleDiff)
	server.mux.HandleFunc("POST /api/diff-merge", server.handleDiffMerge)
	server.mux.HandleFunc("POST /api/sql", server.handleSQL)
	server.mux.HandleFunc("POST /api/profile", server.handleProfile)
	server.mux.Handle("/", http.FileServer(http.Dir(server.webDir)))
}

func (server *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (server *Server) handleCommands(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"commands": commands.Registry()})
}

type openRequest struct {
	Path          string `json:"path"`
	PreviewRows   int    `json:"previewRows"`
	IndexRows     *int   `json:"indexRows"`
	Delimiter     string `json:"delimiter,omitempty"`
	DelimiterName string `json:"delimiterName,omitempty"`
	QuoteChar     string `json:"quoteChar,omitempty"`
	QuoteMode     string `json:"quoteMode,omitempty"`
	HasHeader     *bool  `json:"hasHeader,omitempty"`
	EncodingName  string `json:"encodingName,omitempty"`
}

type openResponse struct {
	DocumentID string           `json:"documentId"`
	Document   csvcore.Document `json:"document"`
}

type fileSampleRequest struct {
	Path     string `json:"path"`
	MaxBytes int    `json:"maxBytes"`
}

type fileSampleResponse struct {
	Sample string `json:"sample"`
}

func (server *Server) handleOpen(w http.ResponseWriter, r *http.Request) {
	var request openRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.Path == "" {
		writeError(w, http.StatusBadRequest, errors.New("path is required"))
		return
	}
	if request.PreviewRows <= 0 {
		request.PreviewRows = 80
	}
	indexRows := request.PreviewRows
	if request.IndexRows != nil {
		indexRows = *request.IndexRows
	}

	options := csvcore.OpenOptions{
		PreviewRows: request.PreviewRows,
		IndexRows:   indexRows,
	}
	if request.Delimiter != "" {
		delimiter, err := delimiterRune(request.Delimiter)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		options.Delimiter = &delimiter
		options.DelimiterName = request.DelimiterName
	}
	if request.HasHeader != nil {
		options.HasHeader = request.HasHeader
	}
	if strings.TrimSpace(request.QuoteChar) != "" {
		quoteChar, err := csvcore.NormalizeQuoteChar(request.QuoteChar)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		options.QuoteChar = &quoteChar
	}
	if strings.TrimSpace(request.QuoteMode) != "" {
		quoteMode, err := csvcore.NormalizeQuoteMode(request.QuoteMode)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		options.QuoteMode = quoteMode
	}
	if strings.TrimSpace(request.EncodingName) != "" {
		encoding, err := csvcore.ParseEncodingName(request.EncodingName)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		options.Encoding = &encoding
	}

	session, err := server.documents.OpenPath(request.Path, docstore.SourceFile, filepath.Base(request.Path), options)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, openResponse{DocumentID: session.ID, Document: session.Document})
}

func delimiterRune(value string) (rune, error) {
	if value == "\\t" || strings.EqualFold(value, "tab") {
		return '\t', nil
	}
	runes := []rune(value)
	if len(runes) != 1 {
		return 0, fmt.Errorf("delimiter must be one character")
	}
	return runes[0], nil
}

func (server *Server) handleFileSample(w http.ResponseWriter, r *http.Request) {
	var request fileSampleRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.Path == "" {
		writeError(w, http.StatusBadRequest, errors.New("path is required"))
		return
	}
	if request.MaxBytes <= 0 {
		request.MaxBytes = 8192
	}
	if request.MaxBytes > 65536 {
		request.MaxBytes = 65536
	}

	file, err := os.Open(request.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	buffer := make([]byte, request.MaxBytes)
	n, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, fileSampleResponse{Sample: string(buffer[:n])})
}

func (server *Server) handleUploadOpen(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var source io.Reader
	var sourceCloser io.Closer
	filename := "upload.csv"
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if part.FormName() == "file" {
			source = part
			sourceCloser = part
			if part.FileName() != "" {
				filename = part.FileName()
			}
			break
		}
		_ = part.Close()
	}
	if source == nil {
		writeError(w, http.StatusBadRequest, errors.New("file is required"))
		return
	}
	defer sourceCloser.Close()

	uploadDir, err := os.MkdirTemp("", "opencsv-upload-*")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	targetPath := filepath.Join(uploadDir, filepath.Base(filename))
	target, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if _, err := io.Copy(target, source); err != nil {
		_ = target.Close()
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := target.Close(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	session, err := server.documents.OpenPath(targetPath, docstore.SourceUpload, filepath.Base(filename), csvcore.OpenOptions{PreviewRows: 160, IndexRows: -1})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, openResponse{DocumentID: session.ID, Document: session.Document})
}

func (server *Server) handleWindow(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	start := queryInt(r, "start", 0)
	count := queryInt(r, "count", 80)

	window, err := server.documents.Window(id, start, count)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, window)
}

type editRequest struct {
	DocumentID string `json:"documentId"`
	Row        int    `json:"row"`
	Col        int    `json:"col"`
	Value      string `json:"value"`
}

func (server *Server) handleEdit(w http.ResponseWriter, r *http.Request) {
	var request editRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.Row < 0 || request.Col < 0 {
		writeError(w, http.StatusBadRequest, errors.New("row and col must be non-negative"))
		return
	}

	if err := server.documents.Edit(request.DocumentID, request.Row, request.Col, request.Value); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"dirty": true})
}

type editBatchRequest struct {
	DocumentID string              `json:"documentId"`
	Edits      []docstore.CellEdit `json:"edits"`
}

func (server *Server) handleEditBatch(w http.ResponseWriter, r *http.Request) {
	var request editBatchRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if len(request.Edits) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{"dirty": false, "count": 0})
		return
	}
	if err := server.documents.EditCells(request.DocumentID, request.Edits); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"dirty": true, "count": len(request.Edits)})
}

type documentRequest struct {
	DocumentID string `json:"documentId"`
}

type exportViewRequest struct {
	DocumentID string `json:"documentId"`
	Columns    []int  `json:"columns,omitempty"`
	Format     string `json:"format,omitempty"`
}

type rowMutationRequest struct {
	DocumentID string `json:"documentId"`
	Row        int    `json:"row"`
	Delta      int    `json:"delta"`
}

type columnMutationRequest struct {
	DocumentID string `json:"documentId"`
	Col        int    `json:"col"`
	Delta      int    `json:"delta"`
}

type transposeRequest struct {
	DocumentID string `json:"documentId"`
	MaxCells   int    `json:"maxCells"`
}

func (server *Server) handleUndo(w http.ResponseWriter, r *http.Request) {
	var request documentRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.Undo(request.DocumentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleRedo(w http.ResponseWriter, r *http.Request) {
	var request documentRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.Redo(request.DocumentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type readOnlyRequest struct {
	DocumentID string `json:"documentId"`
	ReadOnly   bool   `json:"readOnly"`
}

type headerModeRequest struct {
	DocumentID string `json:"documentId"`
	HasHeader  bool   `json:"hasHeader"`
}

type quotePolicyRequest struct {
	DocumentID string `json:"documentId"`
	QuoteChar  string `json:"quoteChar"`
	QuoteMode  string `json:"quoteMode"`
}

func (server *Server) handleHeaderMode(w http.ResponseWriter, r *http.Request) {
	var request headerModeRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SetHeaderMode(request.DocumentID, request.HasHeader)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleQuotePolicy(w http.ResponseWriter, r *http.Request) {
	var request quotePolicyRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SetQuotePolicy(request.DocumentID, request.QuoteChar, request.QuoteMode)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleReadOnly(w http.ResponseWriter, r *http.Request) {
	var request readOnlyRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	session, err := server.documents.SetReadOnly(request.DocumentID, request.ReadOnly)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"session": session, "readOnly": session.ReadOnly})
}

func (server *Server) handleFileStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	status, err := server.documents.FileStatus(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (server *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var request documentRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	session, err := server.documents.Refresh(request.DocumentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"session": session})
}

func (server *Server) handleInsertRow(w http.ResponseWriter, r *http.Request) {
	var request rowMutationRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.InsertRow(request.DocumentID, request.Row)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleDeleteRow(w http.ResponseWriter, r *http.Request) {
	var request rowMutationRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.DeleteRow(request.DocumentID, request.Row)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleMoveRow(w http.ResponseWriter, r *http.Request) {
	var request rowMutationRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.MoveRow(request.DocumentID, request.Row, request.Delta)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleHideRow(w http.ResponseWriter, r *http.Request) {
	var request rowMutationRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.HideRow(request.DocumentID, request.Row)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleInsertColumn(w http.ResponseWriter, r *http.Request) {
	var request columnMutationRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.InsertColumn(request.DocumentID, request.Col)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleDeleteColumn(w http.ResponseWriter, r *http.Request) {
	var request columnMutationRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.DeleteColumn(request.DocumentID, request.Col)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleMoveColumn(w http.ResponseWriter, r *http.Request) {
	var request columnMutationRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.MoveColumn(request.DocumentID, request.Col, request.Delta)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleTranspose(w http.ResponseWriter, r *http.Request) {
	var request transposeRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.Transpose(request.DocumentID, request.MaxCells)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleDedupeRows(w http.ResponseWriter, r *http.Request) {
	var request documentRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.DeduplicateRows(request.DocumentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleDedupeColumns(w http.ResponseWriter, r *http.Request) {
	var request documentRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.DeduplicateColumns(request.DocumentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type textTransformRequest struct {
	DocumentID string          `json:"documentId"`
	Cells      []docstore.Cell `json:"cells"`
	Mode       string          `json:"mode"`
}

type cellsToolRequest struct {
	DocumentID string          `json:"documentId"`
	Cells      []docstore.Cell `json:"cells"`
	Separator  string          `json:"separator"`
	Direction  string          `json:"direction,omitempty"`
}

func (server *Server) handleTextTransform(w http.ResponseWriter, r *http.Request) {
	var request textTransformRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.TransformCells(
		request.DocumentID,
		request.Cells,
		docstore.TextTransformMode(request.Mode),
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleSwapCells(w http.ResponseWriter, r *http.Request) {
	var request cellsToolRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SwapCells(request.DocumentID, request.Cells)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleJoinCells(w http.ResponseWriter, r *http.Request) {
	var request cellsToolRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var result docstore.MutationResult
	var err error
	if strings.EqualFold(strings.TrimSpace(request.Direction), "vertical") {
		result, err = server.documents.JoinCellsVertical(request.DocumentID, request.Cells, request.Separator)
	} else {
		result, err = server.documents.JoinCells(request.DocumentID, request.Cells, request.Separator)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleSplitCells(w http.ResponseWriter, r *http.Request) {
	var request cellsToolRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var result docstore.MutationResult
	var err error
	if strings.EqualFold(strings.TrimSpace(request.Direction), "vertical") {
		result, err = server.documents.SplitCellsVertical(request.DocumentID, request.Cells, request.Separator)
	} else {
		result, err = server.documents.SplitCells(request.DocumentID, request.Cells, request.Separator)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type saveRequest struct {
	DocumentID string `json:"documentId"`
	Path       string `json:"path"`
	Backup     *bool  `json:"backup"`
}

func (server *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	var request saveRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	createBackup := request.Backup == nil || *request.Backup
	result, err := server.documents.Save(request.DocumentID, request.Path, createBackup)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleExportView(w http.ResponseWriter, r *http.Request) {
	var request exportViewRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.DocumentID == "" {
		writeError(w, http.StatusBadRequest, errors.New("documentId is required"))
		return
	}
	result, err := server.documents.ExportView(request.DocumentID, docstore.ExportViewOptions{Columns: request.Columns, Format: request.Format})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	filename := strings.ReplaceAll(result.Filename, `"`, "")
	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("X-Opencsv-Filename", filename)
	w.Header().Set("X-Opencsv-Rows", strconv.Itoa(result.Rows))
	w.Header().Set("X-Opencsv-Columns", strconv.Itoa(result.Columns))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(result.Data)
}

type findRequest struct {
	DocumentID    string          `json:"documentId"`
	Query         string          `json:"query"`
	Column        *int            `json:"column,omitempty"`
	Cells         []docstore.Cell `json:"cells,omitempty"`
	CaseSensitive bool            `json:"caseSensitive"`
	WholeCell     bool            `json:"wholeCell"`
	Regex         bool            `json:"regex"`
	Limit         int             `json:"limit"`
}

func (server *Server) handleFind(w http.ResponseWriter, r *http.Request) {
	var request findRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.FindDetailed(request.DocumentID, docstore.FindOptions{
		Query:         request.Query,
		Column:        request.Column,
		Cells:         request.Cells,
		CaseSensitive: request.CaseSensitive,
		WholeCell:     request.WholeCell,
		Regex:         request.Regex,
		Limit:         request.Limit,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type replaceRequest struct {
	DocumentID    string          `json:"documentId"`
	Query         string          `json:"query"`
	Replacement   string          `json:"replacement"`
	Column        *int            `json:"column,omitempty"`
	Cells         []docstore.Cell `json:"cells,omitempty"`
	CaseSensitive bool            `json:"caseSensitive"`
	WholeCell     bool            `json:"wholeCell"`
	Regex         bool            `json:"regex"`
	All           bool            `json:"all"`
}

func (server *Server) handleReplace(w http.ResponseWriter, r *http.Request) {
	var request replaceRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.Replace(request.DocumentID, docstore.ReplaceOptions{
		FindOptions: docstore.FindOptions{
			Query:         request.Query,
			Column:        request.Column,
			Cells:         request.Cells,
			CaseSensitive: request.CaseSensitive,
			WholeCell:     request.WholeCell,
			Regex:         request.Regex,
			Limit:         1,
		},
		Replacement: request.Replacement,
		All:         request.All,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type filterRequest struct {
	DocumentID    string `json:"documentId"`
	Query         string `json:"query"`
	Expression    string `json:"expression"`
	Mode          string `json:"mode"`
	Column        *int   `json:"column,omitempty"`
	CaseSensitive bool   `json:"caseSensitive"`
}

type filterValuesRequest struct {
	DocumentID   string   `json:"documentId"`
	Column       int      `json:"column"`
	Values       []string `json:"values"`
	IncludeBlank bool     `json:"includeBlank"`
}

func (server *Server) handleFilter(w http.ResponseWriter, r *http.Request) {
	var request filterRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var count int
	var err error
	if strings.EqualFold(request.Mode, "sql") || strings.TrimSpace(request.Expression) != "" {
		expression := request.Expression
		if expression == "" {
			expression = request.Query
		}
		count, err = server.documents.FilterSQL(request.DocumentID, expression)
	} else {
		count, err = server.documents.FilterContains(request.DocumentID, docstore.FilterOptions{
			Query:         request.Query,
			Column:        request.Column,
			CaseSensitive: request.CaseSensitive,
		})
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"totalRows": count})
}

func (server *Server) handleColumnValues(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	column := queryInt(r, "column", -1)
	limit := queryInt(r, "limit", 50)
	values, totalDistinct, err := server.documents.ColumnValues(id, column, limit)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"values":        values,
		"totalDistinct": totalDistinct,
	})
}

func (server *Server) handleFilterValues(w http.ResponseWriter, r *http.Request) {
	var request filterValuesRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	count, err := server.documents.FilterValues(request.DocumentID, docstore.ValueFilterOptions{
		Column:       request.Column,
		Values:       request.Values,
		IncludeBlank: request.IncludeBlank,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"totalRows": count})
}

type sortRequest struct {
	DocumentID string             `json:"documentId"`
	Column     int                `json:"column"`
	Descending bool               `json:"descending"`
	Numeric    bool               `json:"numeric"`
	Mode       docstore.SortMode  `json:"mode,omitempty"`
	HasHeader  bool               `json:"hasHeader"`
	Keys       []docstore.SortKey `json:"keys"`
}

func (server *Server) handleSort(w http.ResponseWriter, r *http.Request) {
	var request sortRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	count, err := server.documents.Sort(request.DocumentID, docstore.SortOptions{
		Column:     request.Column,
		Descending: request.Descending,
		Numeric:    request.Numeric,
		Mode:       request.Mode,
		HasHeader:  request.HasHeader,
		Keys:       request.Keys,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"totalRows": count})
}

type viewClearRequest struct {
	DocumentID string `json:"documentId"`
}

func (server *Server) handleViewClear(w http.ResponseWriter, r *http.Request) {
	var request viewClearRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	count, err := server.documents.ClearView(request.DocumentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"totalRows": count})
}

type diffRequest struct {
	LeftPath           string `json:"leftPath"`
	RightPath          string `json:"rightPath"`
	RightText          string `json:"rightText"`
	RightLabel         string `json:"rightLabel"`
	MaxRows            int    `json:"maxRows"`
	IgnoreWhitespace   bool   `json:"ignoreWhitespace"`
	IgnoreCase         bool   `json:"ignoreCase"`
	IgnoredColumns     []int  `json:"ignoredColumns"`
	MatchColumnsByName bool   `json:"matchColumnsByName"`
	KeyColumn          string `json:"keyColumn"`
}

type diffMergeRequest struct {
	diffRequest
	ResultIndex int    `json:"resultIndex"`
	Mode        string `json:"mode"`
	Col         int    `json:"col"`
	Direction   string `json:"direction"`
}

func (server *Server) handleDiff(w http.ResponseWriter, r *http.Request) {
	var request diffRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.LeftPath == "" || (request.RightPath == "" && request.RightText == "") {
		writeError(w, http.StatusBadRequest, errors.New("leftPath and rightPath or rightText are required"))
		return
	}
	if request.LeftPath == "-" || request.RightPath == "-" {
		writeError(w, http.StatusBadRequest, errors.New("diff does not support stdin"))
		return
	}
	if request.MaxRows <= 0 {
		request.MaxRows = 500
	}
	options := csvdiff.Options{
		MaxRows:            request.MaxRows,
		IgnoreWhitespace:   request.IgnoreWhitespace,
		IgnoreCase:         request.IgnoreCase,
		IgnoredColumns:     request.IgnoredColumns,
		MatchColumnsByName: request.MatchColumnsByName,
		KeyColumn:          request.KeyColumn,
	}
	var result csvdiff.Result
	var err error
	if request.RightText != "" {
		result, err = csvdiff.CompareFileToText(request.LeftPath, request.RightLabel, request.RightText, options)
	} else {
		result, err = csvdiff.CompareFiles(request.LeftPath, request.RightPath, options)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleDiffMerge(w http.ResponseWriter, r *http.Request) {
	var request diffMergeRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.LeftPath == "" || (request.RightPath == "" && request.RightText == "") {
		writeError(w, http.StatusBadRequest, errors.New("leftPath and rightPath or rightText are required"))
		return
	}
	if request.LeftPath == "-" || request.RightPath == "-" {
		writeError(w, http.StatusBadRequest, errors.New("diff merge does not support stdin"))
		return
	}
	options := csvdiff.MergeOptions{
		DiffOptions: csvdiff.Options{
			IgnoreWhitespace:   request.IgnoreWhitespace,
			IgnoreCase:         request.IgnoreCase,
			IgnoredColumns:     request.IgnoredColumns,
			MatchColumnsByName: request.MatchColumnsByName,
			KeyColumn:          request.KeyColumn,
		},
		ResultIndex: request.ResultIndex,
		Mode:        csvdiff.MergeMode(request.Mode),
		Col:         request.Col,
	}
	var result csvdiff.MergeResult
	var err error
	if request.Direction == "left-to-right" {
		if request.RightText != "" {
			writeError(w, http.StatusBadRequest, errors.New("left-to-right merge requires a rightPath"))
			return
		}
		result, err = csvdiff.MergeLeftToRight(request.LeftPath, request.RightPath, options)
	} else if request.RightText != "" {
		result, err = csvdiff.MergeRightTextToLeft(request.LeftPath, request.RightLabel, request.RightText, options)
	} else {
		result, err = csvdiff.MergeRightToLeft(request.LeftPath, request.RightPath, options)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type sqlRequest struct {
	DocumentID string `json:"documentId"`
	Query      string `json:"query"`
	Limit      int    `json:"limit"`
}

func (server *Server) handleSQL(w http.ResponseWriter, r *http.Request) {
	var request sqlRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := sqlconsole.Execute(request.Query, server.sqlResolver(request.DocumentID), sqlconsole.Options{
		DefaultLimit: max(1, request.Limit),
		MaxLimit:     5000,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type profileRequest struct {
	DocumentID string `json:"documentId"`
}

func (server *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	var request profileRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	session, records, err := server.documents.Records(request.DocumentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result := dataprofile.Analyze(records, session.Document.Metadata.Dialect.HasHeader)
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) sqlResolver(documentID string) sqlconsole.Resolver {
	return func(ref string) (sqlconsole.Source, error) {
		if ref == "" || strings.EqualFold(ref, "current") || strings.EqualFold(ref, "this") {
			if documentID == "" {
				return sqlconsole.Source{}, errors.New("current document is required")
			}
			return server.sqlSourceFromDocument(documentID)
		}
		if documentID != "" {
			if session, records, err := server.documents.Records(documentID); err == nil {
				if ref == session.OriginalPath || ref == session.WorkingPath || ref == session.DisplayName {
					return sqlSourceFromSession(session, records), nil
				}
			}
		}
		doc, err := csvcore.OpenDocument(ref, csvcore.OpenOptions{PreviewRows: 50, IndexRows: -1})
		if err != nil {
			return sqlconsole.Source{}, err
		}
		records, err := csvcore.ReadAllRecords(ref, doc.Metadata)
		if err != nil {
			return sqlconsole.Source{}, err
		}
		return sqlSourceFromDocument(ref, doc, records), nil
	}
}

func (server *Server) sqlSourceFromDocument(documentID string) (sqlconsole.Source, error) {
	session, records, err := server.documents.Records(documentID)
	if err != nil {
		return sqlconsole.Source{}, err
	}
	return sqlSourceFromSession(session, records), nil
}

func sqlSourceFromSession(session *docstore.Session, records [][]string) sqlconsole.Source {
	return sqlSourceFromDocument(session.WorkingPath, session.Document, records)
}

func sqlSourceFromDocument(path string, doc csvcore.Document, records [][]string) sqlconsole.Source {
	headers := []string(nil)
	if doc.Metadata.Dialect.HasHeader && len(records) > 0 {
		headers = records[0]
	} else if doc.Metadata.Dialect.EstimatedColumns > 0 {
		headers = make([]string, doc.Metadata.Dialect.EstimatedColumns)
		for i := range headers {
			headers[i] = fmt.Sprintf("c%d", i+1)
		}
	}
	return sqlconsole.Source{
		Name:    filepath.Base(path),
		Path:    path,
		Headers: headers,
		Records: records,
	}
}

func readJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func queryInt(r *http.Request, name string, fallback int) int {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func ResolveWebDir(value string) string {
	if value != "" {
		return value
	}
	return filepath.Join(".", "web")
}

func ListenAndServe(addr string, webDir string) error {
	if addr == "" {
		addr = "127.0.0.1:43255"
	}
	server := New(ResolveWebDir(webDir))
	fmt.Printf("opencsv listening on http://%s\n", addr)
	return http.ListenAndServe(addr, server.Handler())
}
