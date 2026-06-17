package appserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

const (
	statusClientClosedRequest   = 499
	maxSQLResultDocumentRows    = 5000
	maxSQLResultDocumentColumns = 512
)

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
	server.mux.HandleFunc("GET /api/update-check", server.handleUpdateCheck)
	server.mux.HandleFunc("POST /api/open", server.handleOpen)
	server.mux.HandleFunc("POST /api/file-sample", server.handleFileSample)
	server.mux.HandleFunc("POST /api/new", server.handleNew)
	server.mux.HandleFunc("POST /api/upload-open", server.handleUploadOpen)
	server.mux.HandleFunc("GET /api/session", server.handleSession)
	server.mux.HandleFunc("POST /api/close", server.handleClose)
	server.mux.HandleFunc("GET /api/window", server.handleWindow)
	server.mux.HandleFunc("POST /api/edit", server.handleEdit)
	server.mux.HandleFunc("POST /api/edit-batch", server.handleEditBatch)
	server.mux.HandleFunc("POST /api/undo", server.handleUndo)
	server.mux.HandleFunc("POST /api/redo", server.handleRedo)
	server.mux.HandleFunc("POST /api/header-mode", server.handleHeaderMode)
	server.mux.HandleFunc("POST /api/delimiter", server.handleDelimiter)
	server.mux.HandleFunc("POST /api/encoding", server.handleEncoding)
	server.mux.HandleFunc("POST /api/line-ending", server.handleLineEnding)
	server.mux.HandleFunc("POST /api/quote-policy", server.handleQuotePolicy)
	server.mux.HandleFunc("POST /api/quote-rule", server.handleQuoteRule)
	server.mux.HandleFunc("POST /api/quote-settings", server.handleQuoteSettings)
	server.mux.HandleFunc("POST /api/quote-value-rule", server.handleQuoteValueRule)
	server.mux.HandleFunc("POST /api/cell-quote", server.handleCellQuote)
	server.mux.HandleFunc("POST /api/read-only", server.handleReadOnly)
	server.mux.HandleFunc("GET /api/file-status", server.handleFileStatus)
	server.mux.HandleFunc("POST /api/refresh", server.handleRefresh)
	server.mux.HandleFunc("POST /api/row/insert", server.handleInsertRow)
	server.mux.HandleFunc("POST /api/row/delete", server.handleDeleteRow)
	server.mux.HandleFunc("POST /api/row/move", server.handleMoveRow)
	server.mux.HandleFunc("POST /api/row/paste", server.handlePasteRow)
	server.mux.HandleFunc("POST /api/row/hide", server.handleHideRow)
	server.mux.HandleFunc("POST /api/column/insert", server.handleInsertColumn)
	server.mux.HandleFunc("POST /api/column/delete", server.handleDeleteColumn)
	server.mux.HandleFunc("POST /api/column/move", server.handleMoveColumn)
	server.mux.HandleFunc("POST /api/column/paste", server.handlePasteColumn)
	server.mux.HandleFunc("POST /api/transpose", server.handleTranspose)
	server.mux.HandleFunc("POST /api/dedupe-rows", server.handleDedupeRows)
	server.mux.HandleFunc("POST /api/dedupe-columns", server.handleDedupeColumns)
	server.mux.HandleFunc("POST /api/text-transform", server.handleTextTransform)
	server.mux.HandleFunc("POST /api/cells/swap", server.handleSwapCells)
	server.mux.HandleFunc("POST /api/cells/join", server.handleJoinCells)
	server.mux.HandleFunc("POST /api/cells/split", server.handleSplitCells)
	server.mux.HandleFunc("POST /api/save", server.handleSave)
	server.mux.HandleFunc("GET /api/backups", server.handleBackups)
	server.mux.HandleFunc("POST /api/restore-backup", server.handleRestoreBackup)
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
	server.mux.HandleFunc("POST /api/diff-edit-cell", server.handleDiffEditCell)
	server.mux.HandleFunc("POST /api/sql", server.handleSQL)
	server.mux.HandleFunc("POST /api/sql-result/open", server.handleOpenSQLResult)
	server.mux.HandleFunc("POST /api/profile", server.handleProfile)
	server.mux.Handle("/", http.FileServer(http.Dir(server.webDir)))
}

func (server *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (server *Server) handleCommands(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"commands": commands.Registry()})
}

func (server *Server) handleUpdateCheck(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")
	if source == "" {
		source = "manual"
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":         "offline",
		"source":         source,
		"currentVersion": "dev",
		"latestVersion":  "dev",
		"checkedAt":      time.Now().UTC().Format(time.RFC3339),
		"message":        "No remote update feed is configured for the local prototype.",
	})
}

type openRequest struct {
	Path          string `json:"path"`
	PreviewRows   int    `json:"previewRows"`
	IndexRows     *int   `json:"indexRows"`
	SourceKind    string `json:"sourceKind,omitempty"`
	DisplayName   string `json:"displayName,omitempty"`
	Delimiter     string `json:"delimiter,omitempty"`
	DelimiterName string `json:"delimiterName,omitempty"`
	QuoteChar     string `json:"quoteChar,omitempty"`
	QuoteMode     string `json:"quoteMode,omitempty"`
	QuoteEscape   string `json:"quoteEscape,omitempty"`
	LineEnding    string `json:"lineEnding,omitempty"`
	HasHeader     *bool  `json:"hasHeader,omitempty"`
	HeaderColumns *int   `json:"headerColumns,omitempty"`
	EncodingName  string `json:"encodingName,omitempty"`
	SheetName     string `json:"sheetName,omitempty"`
	SheetIndex    int    `json:"sheetIndex,omitempty"`
}

type openResponse struct {
	DocumentID  string              `json:"documentId"`
	Document    csvcore.Document    `json:"document"`
	SourceKind  docstore.SourceKind `json:"sourceKind,omitempty"`
	DisplayName string              `json:"displayName,omitempty"`
}

type newDocumentRequest struct {
	Name    string `json:"name"`
	Rows    int    `json:"rows"`
	Columns int    `json:"columns"`
}

type sqlResultOpenRequest struct {
	Name    string     `json:"name"`
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
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
		delimiter, delimiterName, err := csvcore.NormalizeDelimiter(request.Delimiter, request.DelimiterName)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		options.Delimiter = &delimiter
		options.DelimiterName = delimiterName
	}
	if request.HasHeader != nil {
		options.HasHeader = request.HasHeader
	}
	if request.HeaderColumns != nil {
		if *request.HeaderColumns < 0 {
			writeError(w, http.StatusBadRequest, errors.New("headerColumns must be non-negative"))
			return
		}
		options.HeaderColumns = request.HeaderColumns
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
	if strings.TrimSpace(request.QuoteEscape) != "" {
		quoteEscape, err := csvcore.NormalizeQuoteEscape(request.QuoteEscape)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		options.QuoteEscape = quoteEscape
	}
	if strings.TrimSpace(request.LineEnding) != "" {
		lineEnding, err := csvcore.NormalizeLineEndingName(request.LineEnding)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		options.LineEnding = &lineEnding
	}
	if strings.TrimSpace(request.EncodingName) != "" {
		encoding, err := csvcore.ParseEncodingName(request.EncodingName)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		options.Encoding = &encoding
	}
	if strings.TrimSpace(request.SheetName) != "" {
		options.SheetName = strings.TrimSpace(request.SheetName)
	}
	if request.SheetIndex < 0 {
		writeError(w, http.StatusBadRequest, errors.New("sheetIndex must be positive"))
		return
	}
	if request.SheetIndex > 0 {
		options.SheetIndex = request.SheetIndex
	}

	sourceKind, displayName, err := openRequestSource(request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	session, err := server.documents.OpenPath(request.Path, sourceKind, displayName, options)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, openResponseFromSession(session))
}

func openResponseFromSession(session *docstore.Session) openResponse {
	return openResponse{
		DocumentID:  session.ID,
		Document:    session.Document,
		SourceKind:  session.SourceKind,
		DisplayName: session.DisplayName,
	}
}

func openRequestSource(request openRequest) (docstore.SourceKind, string, error) {
	displayName := cleanOpenDisplayName(request.DisplayName)
	if displayName == "" {
		displayName = cleanOpenDisplayName(filepath.Base(request.Path))
	}
	sourceKind := strings.ToLower(strings.TrimSpace(request.SourceKind))
	switch sourceKind {
	case "", string(docstore.SourceFile):
		return docstore.SourceFile, displayName, nil
	case string(docstore.SourceStdin):
		if displayName == "" || displayName == "." {
			displayName = "stdin"
		}
		return docstore.SourceStdin, displayName, nil
	default:
		return "", "", fmt.Errorf("unsupported sourceKind %q", request.SourceKind)
	}
}

func (server *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("documentId is required"))
		return
	}
	session, ok := server.documents.Snapshot(id)
	if !ok {
		writeError(w, http.StatusBadRequest, errors.New("document not found"))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"session":     session,
		"documentId":  session.ID,
		"document":    session.Document,
		"sourceKind":  session.SourceKind,
		"displayName": session.DisplayName,
		"dirty":       session.Dirty,
		"readOnly":    session.ReadOnly,
	})
}

type closeRequest struct {
	DocumentID string `json:"documentId"`
	Force      bool   `json:"force,omitempty"`
}

func (server *Server) handleClose(w http.ResponseWriter, r *http.Request) {
	var request closeRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.DocumentID == "" {
		writeError(w, http.StatusBadRequest, errors.New("documentId is required"))
		return
	}
	session, err := server.documents.Close(request.DocumentID, request.Force)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"closed":      true,
		"documentId":  session.ID,
		"dirty":       session.Dirty,
		"readOnly":    session.ReadOnly,
		"sourceKind":  session.SourceKind,
		"displayName": session.DisplayName,
	})
}

func cleanOpenDisplayName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	name = filepath.Base(name)
	if name == "." || name == ".." || name == string(filepath.Separator) {
		return ""
	}
	return name
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

func (server *Server) handleNew(w http.ResponseWriter, r *http.Request) {
	var request newDocumentRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	filename := newDocumentFilename(request.Name)
	rows := clampNewDocumentDimension(request.Rows, 1, 1000)
	columns := clampNewDocumentDimension(request.Columns, 2, 128)
	records := newDocumentRecords(rows, columns)

	tempDir, err := os.MkdirTemp("", "opencsv-new-*")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	targetPath := filepath.Join(tempDir, filename)
	metadata := csvcore.FileMetadata{
		Path:       targetPath,
		Encoding:   csvcore.EncodingInfo{Name: csvcore.EncodingUTF8, Confidence: 1},
		LineEnding: csvcore.LineEndingInfo{Name: csvcore.LineEndingLF, Sequence: "\n", Confidence: 1},
		Dialect: csvcore.Dialect{
			Delimiter:        ',',
			DelimiterName:    "comma",
			HasHeader:        true,
			EstimatedColumns: columns,
		},
	}
	if _, err := csvcore.SaveRecords(targetPath, metadata, records, csvcore.SaveOptions{}); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	hasHeader := true
	session, err := server.documents.OpenPath(targetPath, docstore.SourceUpload, filename, csvcore.OpenOptions{PreviewRows: 160, IndexRows: -1, HasHeader: &hasHeader})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, openResponse{DocumentID: session.ID, Document: session.Document})
}

func (server *Server) handleOpenSQLResult(w http.ResponseWriter, r *http.Request) {
	var request sqlResultOpenRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	records, err := sqlResultDocumentRecords(request.Columns, request.Rows)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	filename := sqlResultDocumentFilename(request.Name)
	tempDir, err := os.MkdirTemp("", "opencsv-sql-result-*")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	targetPath := filepath.Join(tempDir, filename)
	metadata := csvcore.FileMetadata{
		Path:       targetPath,
		Encoding:   csvcore.EncodingInfo{Name: csvcore.EncodingUTF8, Confidence: 1},
		LineEnding: csvcore.LineEndingInfo{Name: csvcore.LineEndingLF, Sequence: "\n", Confidence: 1},
		Dialect: csvcore.Dialect{
			Delimiter:        ',',
			DelimiterName:    "comma",
			HasHeader:        true,
			EstimatedColumns: len(records[0]),
		},
	}
	if _, err := csvcore.SaveRecords(targetPath, metadata, records, csvcore.SaveOptions{}); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	hasHeader := true
	session, err := server.documents.OpenPath(targetPath, docstore.SourceSQLResult, filename, csvcore.OpenOptions{PreviewRows: 160, IndexRows: -1, HasHeader: &hasHeader})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, openResponseFromSession(session))
}

func newDocumentFilename(raw string) string {
	name := strings.TrimSpace(raw)
	if name == "" {
		name = "Untitled.csv"
	}
	name = filepath.Base(name)
	if name == "." || name == ".." || name == string(filepath.Separator) {
		name = "Untitled.csv"
	}
	if filepath.Ext(name) == "" {
		name += ".csv"
	}
	return name
}

func sqlResultDocumentFilename(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return "SQL Result.csv"
	}
	return newDocumentFilename(raw)
}

func sqlResultDocumentRecords(columns []string, rows [][]string) ([][]string, error) {
	if len(columns) == 0 {
		return nil, errors.New("columns are required")
	}
	if len(columns) > maxSQLResultDocumentColumns {
		return nil, fmt.Errorf("columns exceed limit of %d", maxSQLResultDocumentColumns)
	}
	if len(rows) > maxSQLResultDocumentRows {
		return nil, fmt.Errorf("rows exceed limit of %d", maxSQLResultDocumentRows)
	}

	width := len(columns)
	records := make([][]string, 0, len(rows)+1)
	header := make([]string, width)
	for index, column := range columns {
		header[index] = column
		if strings.TrimSpace(header[index]) == "" {
			header[index] = fmt.Sprintf("Column %d", index+1)
		}
	}
	records = append(records, header)
	for _, row := range rows {
		record := make([]string, width)
		copy(record, row)
		records = append(records, record)
	}
	return records, nil
}

func clampNewDocumentDimension(value int, fallback int, limit int) int {
	if value <= 0 {
		return fallback
	}
	if value > limit {
		return limit
	}
	return value
}

func newDocumentRecords(rows int, columns int) [][]string {
	records := make([][]string, 0, rows+1)
	header := make([]string, columns)
	for column := range header {
		header[column] = fmt.Sprintf("Column %d", column+1)
	}
	records = append(records, header)
	for range rows {
		records = append(records, make([]string, columns))
	}
	return records
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
	DocumentID   string `json:"documentId"`
	DiscardDirty bool   `json:"discardDirty,omitempty"`
}

type exportViewRequest struct {
	DocumentID     string `json:"documentId"`
	Columns        []int  `json:"columns,omitempty"`
	Rows           []int  `json:"rows,omitempty"`
	Format         string `json:"format,omitempty"`
	IncludeHeaders *bool  `json:"includeHeaders,omitempty"`
	FullDocument   bool   `json:"fullDocument,omitempty"`
}

type rowMutationRequest struct {
	DocumentID string `json:"documentId"`
	Row        int    `json:"row"`
	Delta      int    `json:"delta"`
}

type rowPasteRequest struct {
	DocumentID string `json:"documentId"`
	SourceRow  int    `json:"sourceRow"`
	SourceRows []int  `json:"sourceRows"`
	TargetRow  int    `json:"targetRow"`
	Cut        bool   `json:"cut"`
}

type columnMutationRequest struct {
	DocumentID string `json:"documentId"`
	Col        int    `json:"col"`
	Delta      int    `json:"delta"`
}

type columnPasteRequest struct {
	DocumentID string `json:"documentId"`
	SourceCol  int    `json:"sourceCol"`
	SourceCols []int  `json:"sourceCols"`
	TargetCol  int    `json:"targetCol"`
	Cut        bool   `json:"cut"`
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

type delimiterRequest struct {
	DocumentID    string `json:"documentId"`
	Delimiter     string `json:"delimiter"`
	DelimiterName string `json:"delimiterName,omitempty"`
}

type lineEndingRequest struct {
	DocumentID string `json:"documentId"`
	LineEnding string `json:"lineEnding"`
}

type encodingRequest struct {
	DocumentID   string `json:"documentId"`
	EncodingName string `json:"encodingName"`
}

type quotePolicyRequest struct {
	DocumentID  string `json:"documentId"`
	QuoteChar   string `json:"quoteChar"`
	QuoteMode   string `json:"quoteMode"`
	QuoteEscape string `json:"quoteEscape"`
}

type cellQuoteRequest struct {
	DocumentID string `json:"documentId"`
	Row        int    `json:"row"`
	Col        int    `json:"col"`
	Quoted     bool   `json:"quoted"`
}

type quoteRuleRequest struct {
	DocumentID string `json:"documentId"`
	Target     string `json:"target"`
	Col        int    `json:"col"`
	Quoted     bool   `json:"quoted"`
}

type quoteSettingsRequest struct {
	DocumentID   string         `json:"documentId"`
	QuoteChar    string         `json:"quoteChar"`
	QuoteMode    string         `json:"quoteMode"`
	QuoteEscape  string         `json:"quoteEscape"`
	QuoteHeader  bool           `json:"quoteHeader"`
	QuoteBody    bool           `json:"quoteBody"`
	QuoteColumns []int          `json:"quoteColumns"`
	QuoteCells   []csvcore.Cell `json:"quoteCells"`
	NoQuoteCells []csvcore.Cell `json:"noQuoteCells"`
}

type quoteValueRuleRequest struct {
	DocumentID string `json:"documentId"`
	Pattern    string `json:"pattern"`
	MatchMode  string `json:"matchMode"`
	Action     string `json:"action"`
	MaxCells   int    `json:"maxCells"`
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

func (server *Server) handleDelimiter(w http.ResponseWriter, r *http.Request) {
	var request delimiterRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SetDelimiter(request.DocumentID, request.Delimiter, request.DelimiterName)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleLineEnding(w http.ResponseWriter, r *http.Request) {
	var request lineEndingRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SetLineEnding(request.DocumentID, request.LineEnding)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleEncoding(w http.ResponseWriter, r *http.Request) {
	var request encodingRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SetEncoding(request.DocumentID, request.EncodingName)
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
	result, err := server.documents.SetQuotePolicy(request.DocumentID, request.QuoteChar, request.QuoteMode, request.QuoteEscape)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleQuoteRule(w http.ResponseWriter, r *http.Request) {
	var request quoteRuleRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SetQuoteRule(request.DocumentID, request.Target, request.Col, request.Quoted)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleQuoteSettings(w http.ResponseWriter, r *http.Request) {
	var request quoteSettingsRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SetQuoteSettings(request.DocumentID, docstore.QuoteSettings{
		QuoteChar:    request.QuoteChar,
		QuoteMode:    request.QuoteMode,
		QuoteEscape:  request.QuoteEscape,
		QuoteHeader:  request.QuoteHeader,
		QuoteBody:    request.QuoteBody,
		QuoteColumns: request.QuoteColumns,
		QuoteCells:   request.QuoteCells,
		NoQuoteCells: request.NoQuoteCells,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleQuoteValueRule(w http.ResponseWriter, r *http.Request) {
	var request quoteValueRuleRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.ApplyQuoteValueRule(request.DocumentID, docstore.QuoteValueRuleOptions{
		Pattern:   request.Pattern,
		MatchMode: request.MatchMode,
		Action:    request.Action,
		MaxCells:  request.MaxCells,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleCellQuote(w http.ResponseWriter, r *http.Request) {
	var request cellQuoteRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.SetCellQuote(request.DocumentID, request.Row, request.Col, request.Quoted)
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
	session, err := server.documents.Refresh(request.DocumentID, request.DiscardDirty)
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

func (server *Server) handlePasteRow(w http.ResponseWriter, r *http.Request) {
	var request rowPasteRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	sourceRows := request.SourceRows
	if len(sourceRows) == 0 {
		sourceRows = []int{request.SourceRow}
	}
	result, err := server.documents.PasteRows(request.DocumentID, sourceRows, request.TargetRow, request.Cut)
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

func (server *Server) handlePasteColumn(w http.ResponseWriter, r *http.Request) {
	var request columnPasteRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	sourceCols := request.SourceCols
	if len(sourceCols) == 0 {
		sourceCols = []int{request.SourceCol}
	}
	result, err := server.documents.PasteColumns(request.DocumentID, sourceCols, request.TargetCol, request.Cut)
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
	Column     *int            `json:"column,omitempty"`
	Mode       string          `json:"mode"`
}

type cellsToolRequest struct {
	DocumentID     string          `json:"documentId"`
	Cells          []docstore.Cell `json:"cells"`
	Separator      string          `json:"separator"`
	Direction      string          `json:"direction,omitempty"`
	AllowOverwrite bool            `json:"allowOverwrite,omitempty"`
}

func (server *Server) handleTextTransform(w http.ResponseWriter, r *http.Request) {
	var request textTransformRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.Column != nil {
		result, err := server.documents.TransformColumn(
			request.DocumentID,
			*request.Column,
			docstore.TextTransformMode(request.Mode),
		)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
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
	options := docstore.SplitOptions{AllowOverwrite: request.AllowOverwrite}
	if strings.EqualFold(strings.TrimSpace(request.Direction), "vertical") {
		result, err = server.documents.SplitCellsVerticalWithOptions(request.DocumentID, request.Cells, request.Separator, options)
	} else {
		result, err = server.documents.SplitCellsWithOptions(request.DocumentID, request.Cells, request.Separator, options)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type saveRequest struct {
	DocumentID             string `json:"documentId"`
	Path                   string `json:"path"`
	Backup                 *bool  `json:"backup"`
	AllowExternalOverwrite bool   `json:"allowExternalOverwrite"`
}

type restoreBackupRequest struct {
	DocumentID string `json:"documentId"`
	Path       string `json:"path"`
}

func (server *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	var request saveRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	createBackup := request.Backup == nil || *request.Backup
	result, err := server.documents.SaveWithOptions(request.DocumentID, request.Path, createBackup, docstore.SaveOptions{
		AllowExternalOverwrite: request.AllowExternalOverwrite,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleBackups(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	backups, err := server.documents.ListBackups(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"backups": backups})
}

func (server *Server) handleRestoreBackup(w http.ResponseWriter, r *http.Request) {
	var request restoreBackupRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.RestoreBackup(request.DocumentID, request.Path)
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
	result, err := server.documents.ExportView(request.DocumentID, docstore.ExportViewOptions{
		Columns:        request.Columns,
		Rows:           request.Rows,
		Format:         request.Format,
		IncludeHeaders: request.IncludeHeaders,
		FullDocument:   request.FullDocument,
	})
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
	Offset        int             `json:"offset"`
}

func (server *Server) handleFind(w http.ResponseWriter, r *http.Request) {
	var request findRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.FindDetailedContext(r.Context(), request.DocumentID, docstore.FindOptions{
		Query:         request.Query,
		Column:        request.Column,
		Cells:         request.Cells,
		CaseSensitive: request.CaseSensitive,
		WholeCell:     request.WholeCell,
		Regex:         request.Regex,
		Limit:         request.Limit,
		Offset:        request.Offset,
	})
	if err != nil {
		writeError(w, requestContextStatus(err, http.StatusBadRequest), err)
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
	Preview       bool            `json:"preview"`
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
		Preview:     request.Preview,
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
	var result docstore.FilterResult
	var err error
	if strings.EqualFold(request.Mode, "sql") || strings.TrimSpace(request.Expression) != "" {
		expression := request.Expression
		if expression == "" {
			expression = request.Query
		}
		result, err = server.documents.FilterSQLDetailedContext(r.Context(), request.DocumentID, expression)
	} else {
		result, err = server.documents.FilterContainsDetailedContext(r.Context(), request.DocumentID, docstore.FilterOptions{
			Query:         request.Query,
			Column:        request.Column,
			CaseSensitive: request.CaseSensitive,
		})
	}
	if err != nil {
		writeError(w, requestContextStatus(err, http.StatusBadRequest), err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleColumnValues(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	column := queryInt(r, "column", -1)
	limit := queryInt(r, "limit", 50)
	result, err := server.documents.ColumnValuesContext(r.Context(), id, column, limit)
	if err != nil {
		writeError(w, requestContextStatus(err, http.StatusBadRequest), err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleFilterValues(w http.ResponseWriter, r *http.Request) {
	var request filterValuesRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := server.documents.FilterValuesDetailedContext(r.Context(), request.DocumentID, docstore.ValueFilterOptions{
		Column:       request.Column,
		Values:       request.Values,
		IncludeBlank: request.IncludeBlank,
	})
	if err != nil {
		writeError(w, requestContextStatus(err, http.StatusBadRequest), err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type sortRequest struct {
	DocumentID string             `json:"documentId"`
	Column     int                `json:"column"`
	Descending bool               `json:"descending"`
	Numeric    bool               `json:"numeric"`
	Mode       docstore.SortMode  `json:"mode,omitempty"`
	HasHeader  bool               `json:"hasHeader"`
	Keys       []docstore.SortKey `json:"keys"`
	Stable     *bool              `json:"stable,omitempty"`
}

func (server *Server) handleSort(w http.ResponseWriter, r *http.Request) {
	var request sortRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	stable := request.Stable == nil || *request.Stable
	count, err := server.documents.Sort(request.DocumentID, docstore.SortOptions{
		Column:     request.Column,
		Descending: request.Descending,
		Numeric:    request.Numeric,
		Mode:       request.Mode,
		HasHeader:  request.HasHeader,
		Keys:       request.Keys,
		Stable:     &stable,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"totalRows": count, "stable": stable})
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
	LeftPath           string   `json:"leftPath"`
	RightPath          string   `json:"rightPath"`
	RightText          string   `json:"rightText"`
	RightLabel         string   `json:"rightLabel"`
	MaxRows            int      `json:"maxRows"`
	IgnoreWhitespace   bool     `json:"ignoreWhitespace"`
	IgnoreCase         bool     `json:"ignoreCase"`
	IgnoredColumns     []int    `json:"ignoredColumns"`
	IgnoredColumnNames []string `json:"ignoredColumnNames"`
	MatchColumnsByName bool     `json:"matchColumnsByName"`
	KeyColumn          string   `json:"keyColumn"`
}

type diffMergeRequest struct {
	diffRequest
	ResultIndex int    `json:"resultIndex"`
	Mode        string `json:"mode"`
	Col         int    `json:"col"`
	Direction   string `json:"direction"`
}

type diffEditCellRequest struct {
	Path  string `json:"path"`
	Row   int    `json:"row"`
	Col   int    `json:"col"`
	Value string `json:"value"`
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
		IgnoredColumnNames: request.IgnoredColumnNames,
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

func (server *Server) handleDiffEditCell(w http.ResponseWriter, r *http.Request) {
	var request diffEditCellRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.Path == "" {
		writeError(w, http.StatusBadRequest, errors.New("path is required"))
		return
	}
	if request.Path == "-" {
		writeError(w, http.StatusBadRequest, errors.New("diff edit does not support stdin"))
		return
	}
	if request.Row < 1 || request.Col < 1 {
		writeError(w, http.StatusBadRequest, errors.New("row and col must be 1-based positive numbers"))
		return
	}

	document, err := csvcore.OpenDocument(request.Path, csvcore.OpenOptions{PreviewRows: 50, IndexRows: -1})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	records, err := csvcore.ReadAllRecords(request.Path, document.Metadata)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	rowIndex := request.Row - 1
	colIndex := request.Col - 1
	if rowIndex >= len(records) {
		writeError(w, http.StatusBadRequest, errors.New("diff edit row is out of range"))
		return
	}
	for len(records[rowIndex]) <= colIndex {
		records[rowIndex] = append(records[rowIndex], "")
	}
	records[rowIndex][colIndex] = request.Value
	save, err := csvcore.SaveRecords(request.Path, document.Metadata, records, csvcore.SaveOptions{CreateBackup: true})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"save":  save,
		"path":  request.Path,
		"row":   request.Row,
		"col":   request.Col,
		"value": request.Value,
	})
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
			IgnoredColumnNames: request.IgnoredColumnNames,
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
	DocumentID       string `json:"documentId"`
	Query            string `json:"query"`
	Limit            int    `json:"limit"`
	StreamProgress   bool   `json:"streamProgress,omitempty"`
	ProgressInterval int    `json:"progressInterval,omitempty"`
}

type sqlStreamEvent struct {
	Type      string                  `json:"type"`
	Phase     string                  `json:"phase,omitempty"`
	Source    string                  `json:"source,omitempty"`
	Rows      int                     `json:"rows,omitempty"`
	Records   int                     `json:"records,omitempty"`
	Done      bool                    `json:"done,omitempty"`
	ElapsedMs int64                   `json:"elapsedMs,omitempty"`
	Result    *sqlconsole.Result      `json:"result,omitempty"`
	Error     string                  `json:"error,omitempty"`
	SQLError  *sqlconsole.ErrorDetail `json:"sqlError,omitempty"`
}

func (server *Server) handleSQL(w http.ResponseWriter, r *http.Request) {
	var request sqlRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.StreamProgress {
		server.handleSQLProgressStream(w, r, request)
		return
	}
	result, err := server.sqlResult(r.Context(), request, sqlconsole.Options{})
	if err != nil {
		writeError(w, requestContextStatus(err, http.StatusBadRequest), err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleSQLProgressStream(w http.ResponseWriter, r *http.Request, request sqlRequest) {
	startedAt := time.Now()
	w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	flusher, _ := w.(http.Flusher)
	writeEvent := func(event sqlStreamEvent) {
		if event.ElapsedMs == 0 {
			event.ElapsedMs = time.Since(startedAt).Milliseconds()
		}
		_ = json.NewEncoder(w).Encode(event)
		if flusher != nil {
			flusher.Flush()
		}
	}
	result, err := server.sqlResult(r.Context(), request, sqlconsole.Options{
		ProgressInterval: request.ProgressInterval,
		ReportProgress: func(progress sqlconsole.ProgressEvent) {
			writeEvent(sqlStreamEvent{
				Type:    "progress",
				Phase:   progress.Phase,
				Source:  progress.Source,
				Rows:    progress.Rows,
				Records: progress.Records,
				Done:    progress.Done,
			})
		},
	})
	if err != nil {
		event := sqlStreamEvent{Type: "error", Error: err.Error()}
		var queryErr *sqlconsole.QueryError
		if errors.As(err, &queryErr) {
			event.SQLError = &queryErr.Detail
		}
		writeEvent(event)
		return
	}
	writeEvent(sqlStreamEvent{Type: "result", Result: &result})
}

func (server *Server) sqlResult(ctx context.Context, request sqlRequest, options sqlconsole.Options) (sqlconsole.Result, error) {
	options.DefaultLimit = max(1, request.Limit)
	options.MaxLimit = 5000
	return sqlconsole.ExecuteContext(ctx, request.Query, server.sqlResolver(request.DocumentID), options)
}

type profileRequest struct {
	DocumentID       string             `json:"documentId"`
	Schema           dataprofile.Schema `json:"schema"`
	StreamProgress   bool               `json:"streamProgress,omitempty"`
	ProgressInterval int                `json:"progressInterval,omitempty"`
}

type profileStreamEvent struct {
	Type      string              `json:"type"`
	Phase     string              `json:"phase,omitempty"`
	Rows      int                 `json:"rows,omitempty"`
	Records   int                 `json:"records,omitempty"`
	Done      bool                `json:"done,omitempty"`
	ElapsedMs int64               `json:"elapsedMs,omitempty"`
	Result    *dataprofile.Result `json:"result,omitempty"`
	Error     string              `json:"error,omitempty"`
}

func (server *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	var request profileRequest
	if err := readJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if request.StreamProgress {
		server.handleProfileProgressStream(w, r, request)
		return
	}
	result, err := server.profileResult(r.Context(), request, dataprofile.StreamOptions{})
	if err != nil {
		writeError(w, requestContextStatus(err, http.StatusBadRequest), err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (server *Server) handleProfileProgressStream(w http.ResponseWriter, r *http.Request, request profileRequest) {
	startedAt := time.Now()
	w.Header().Set("Content-Type", "application/x-ndjson; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	flusher, _ := w.(http.Flusher)
	writeEvent := func(event profileStreamEvent) {
		if event.ElapsedMs == 0 {
			event.ElapsedMs = time.Since(startedAt).Milliseconds()
		}
		_ = json.NewEncoder(w).Encode(event)
		if flusher != nil {
			flusher.Flush()
		}
	}
	result, err := server.profileResult(r.Context(), request, dataprofile.StreamOptions{
		ProgressInterval: request.ProgressInterval,
		ReportProgress: func(progress dataprofile.ProgressEvent) {
			writeEvent(profileStreamEvent{
				Type:    "progress",
				Phase:   progress.Phase,
				Rows:    progress.Rows,
				Records: progress.Records,
				Done:    progress.Done,
			})
		},
	})
	if err != nil {
		writeEvent(profileStreamEvent{Type: "error", Error: err.Error()})
		return
	}
	writeEvent(profileStreamEvent{Type: "result", Result: &result})
}

func (server *Server) profileResult(ctx context.Context, request profileRequest, options dataprofile.StreamOptions) (dataprofile.Result, error) {
	if session, ok := server.documents.Snapshot(request.DocumentID); ok && canStreamDocumentSource(session) {
		return dataprofile.AnalyzeStreamContextWithOptions(
			ctx,
			session.Document.Metadata.Dialect.HasHeader,
			request.Schema,
			func(streamCtx context.Context, handle func(row int, record []string) error) (int, error) {
				return csvcore.ForEachRecordContext(streamCtx, session.WorkingPath, session.Document.Metadata, handle)
			},
			options,
		)
	}
	session, records, err := server.documents.Records(request.DocumentID)
	if err != nil {
		return dataprofile.Result{}, err
	}
	result := dataprofile.AnalyzeWithSchema(records, session.Document.Metadata.Dialect.HasHeader, request.Schema)
	return result, nil
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
			if session, ok := server.documents.Snapshot(documentID); ok &&
				(ref == session.OriginalPath || ref == session.WorkingPath || ref == session.DisplayName) {
				return server.sqlSourceFromSession(documentID, session)
			}
		}
		doc, err := csvcore.OpenDocument(ref, csvcore.OpenOptions{PreviewRows: 50, IndexRows: 0})
		if err != nil {
			return sqlconsole.Source{}, err
		}
		if csvcore.CanStreamRecords(doc.Metadata) {
			return sqlSourceFromCleanDocument(ref, doc), nil
		}
		records, err := csvcore.ReadAllRecords(ref, doc.Metadata)
		if err != nil {
			return sqlconsole.Source{}, err
		}
		return sqlSourceFromDocumentRecords(ref, doc, records), nil
	}
}

func (server *Server) sqlSourceFromDocument(documentID string) (sqlconsole.Source, error) {
	session, ok := server.documents.Snapshot(documentID)
	if !ok {
		return sqlconsole.Source{}, errors.New("document not found")
	}
	return server.sqlSourceFromSession(documentID, session)
}

func (server *Server) sqlSourceFromSession(documentID string, session *docstore.Session) (sqlconsole.Source, error) {
	if canStreamDocumentSource(session) {
		return sqlSourceFromCleanDocument(session.WorkingPath, session.Document), nil
	}
	session, records, err := server.documents.Records(documentID)
	if err != nil {
		return sqlconsole.Source{}, err
	}
	return sqlSourceFromDocumentRecords(session.WorkingPath, session.Document, records), nil
}

func canStreamDocumentSource(session *docstore.Session) bool {
	return session.EditBaseRecords == nil &&
		len(session.Edits) == 0 &&
		csvcore.CanStreamRecords(session.Document.Metadata)
}

func sqlSourceFromCleanDocument(path string, doc csvcore.Document) sqlconsole.Source {
	source := sqlSourceFromDocumentRecords(path, doc, nil)
	source.HasHeader = doc.Metadata.Dialect.HasHeader
	source.Stream = func(ctx context.Context, handle func(row int, record []string) error) (int, error) {
		return csvcore.ForEachRecordContext(ctx, path, doc.Metadata, handle)
	}
	return source
}

func sqlSourceFromDocumentRecords(path string, doc csvcore.Document, records [][]string) sqlconsole.Source {
	headers := []string(nil)
	if doc.Metadata.Dialect.HasHeader && len(records) > 0 {
		headers = records[0]
	} else if doc.Metadata.Dialect.HasHeader && len(doc.Preview) > 0 {
		headers = doc.Preview[0]
	} else if doc.Metadata.Dialect.EstimatedColumns > 0 {
		headers = make([]string, doc.Metadata.Dialect.EstimatedColumns)
		for i := range headers {
			headers[i] = fmt.Sprintf("c%d", i+1)
		}
	}
	return sqlconsole.Source{
		Name:             filepath.Base(path),
		Path:             path,
		Headers:          headers,
		Records:          records,
		EstimatedColumns: doc.Metadata.Dialect.EstimatedColumns,
		HasHeader:        doc.Metadata.Dialect.HasHeader,
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
	payload := map[string]any{"error": err.Error()}
	var queryErr *sqlconsole.QueryError
	if errors.As(err, &queryErr) {
		payload["sqlError"] = queryErr.Detail
	}
	writeJSON(w, status, payload)
}

func requestContextStatus(err error, fallback int) int {
	if errors.Is(err, context.Canceled) {
		return statusClientClosedRequest
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout
	}
	return fallback
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
