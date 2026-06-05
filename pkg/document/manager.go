package document

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/opencsv/opencsv/pkg/clipboard"
	"github.com/opencsv/opencsv/pkg/csvcore"
	"github.com/opencsv/opencsv/pkg/sqlexpr"
	"github.com/opencsv/opencsv/pkg/xlsx"
	"golang.org/x/text/width"
)

type SourceKind string

const (
	SourceFile   SourceKind = "file"
	SourceUpload SourceKind = "upload"
	SourceStdin  SourceKind = "stdin"
)

type Cell struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type QuoteSettings struct {
	QuoteChar    string
	QuoteMode    string
	QuoteEscape  string
	QuoteHeader  bool
	QuoteBody    bool
	QuoteColumns []int
	QuoteCells   []csvcore.Cell
	NoQuoteCells []csvcore.Cell
}

type Session struct {
	ID                       string            `json:"id"`
	SourceKind               SourceKind        `json:"sourceKind"`
	DisplayName              string            `json:"displayName"`
	OriginalPath             string            `json:"originalPath"`
	WorkingPath              string            `json:"workingPath"`
	Document                 csvcore.Document  `json:"document"`
	Edits                    map[Cell]string   `json:"-"`
	EditBaseRecords          [][]string        `json:"-"`
	ViewRows                 []int             `json:"-"`
	ViewActive               bool              `json:"-"`
	UndoStack                []historySnapshot `json:"-"`
	RedoStack                []historySnapshot `json:"-"`
	Dirty                    bool              `json:"dirty"`
	ReadOnly                 bool              `json:"readOnly"`
	LastKnownPath            string            `json:"-"`
	LastKnownSizeBytes       int64             `json:"-"`
	LastKnownModTimeUnixNano int64             `json:"-"`
}

type historySnapshot struct {
	Records    [][]string
	Metadata   csvcore.FileMetadata
	Dirty      bool
	ViewRows   []int
	ViewActive bool
}

type FileStatus struct {
	Path                 string     `json:"path"`
	SizeBytes            int64      `json:"sizeBytes"`
	ModTimeUnixNano      int64      `json:"modTimeUnixNano"`
	ChangedExternally    bool       `json:"changedExternally"`
	Missing              bool       `json:"missing"`
	Dirty                bool       `json:"dirty"`
	ReadOnly             bool       `json:"readOnly"`
	SourceKind           SourceKind `json:"sourceKind"`
	LastKnownSizeBytes   int64      `json:"lastKnownSizeBytes"`
	LastKnownModTimeNano int64      `json:"lastKnownModTimeNano"`
}

type BackupInfo struct {
	Path            string `json:"path"`
	Name            string `json:"name"`
	SizeBytes       int64  `json:"sizeBytes"`
	ModTimeUnixNano int64  `json:"modTimeUnixNano"`
}

type Manager struct {
	mu       sync.Mutex
	nextID   int
	sessions map[string]*Session
}

const defaultTransposeCellLimit = 1_000_000

func NewManager() *Manager {
	return &Manager{sessions: make(map[string]*Session)}
}

func (manager *Manager) OpenPath(path string, kind SourceKind, displayName string, options csvcore.OpenOptions) (*Session, error) {
	document, workingPath, err := openDocumentForPath(path, options)
	if err != nil {
		return nil, err
	}
	if displayName == "" {
		displayName = path
	}
	return manager.Add(document, kind, displayName, path, workingPath), nil
}

func openDocumentForPath(path string, options csvcore.OpenOptions) (csvcore.Document, string, error) {
	if isXLSXPath(path) {
		document, workingPath, err := openXLSXDocument(path, options)
		return document, workingPath, err
	}
	document, err := csvcore.OpenDocument(path, options)
	return document, path, err
}

func openXLSXDocument(path string, options csvcore.OpenOptions) (csvcore.Document, string, error) {
	records, sheets, selected, err := xlsx.DecodeFileWithOptions(path, xlsx.DecodeOptions{
		SheetName:  options.SheetName,
		SheetIndex: options.SheetIndex,
	})
	if err != nil {
		return csvcore.Document{}, "", err
	}
	workingPath, err := writeImportedXLSXWorkingCopy(path, records)
	if err != nil {
		return csvcore.Document{}, "", err
	}
	document, err := csvcore.OpenDocument(workingPath, options)
	if err != nil {
		return csvcore.Document{}, "", err
	}
	info, err := os.Stat(path)
	if err == nil {
		document.Metadata.SizeBytes = info.Size()
	}
	document.Metadata.Path = path
	document.Metadata.Workbook = workbookInfo(sheets, selected)
	return document, workingPath, nil
}

func workbookInfo(sheets []xlsx.WorksheetInfo, selected xlsx.WorksheetInfo) *csvcore.WorkbookInfo {
	if len(sheets) == 0 {
		return nil
	}
	info := &csvcore.WorkbookInfo{
		Sheets:             make([]csvcore.WorkbookSheetInfo, 0, len(sheets)),
		SelectedSheetIndex: selected.Index,
		SelectedSheetName:  selected.Name,
		SelectedSheetPath:  selected.Path,
	}
	for _, sheet := range sheets {
		info.Sheets = append(info.Sheets, csvcore.WorkbookSheetInfo{
			Index: sheet.Index,
			Name:  sheet.Name,
			Path:  sheet.Path,
		})
	}
	return info
}

func writeImportedXLSXWorkingCopy(sourcePath string, records [][]string) (string, error) {
	dir, err := os.MkdirTemp("", "opencsv-xlsx-*")
	if err != nil {
		return "", err
	}
	stem := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
	if stem == "" {
		stem = "workbook"
	}
	workingPath := filepath.Join(dir, stem+".csv")
	metadata := csvcore.FileMetadata{
		Path:       workingPath,
		Encoding:   csvcore.EncodingInfo{Name: csvcore.EncodingUTF8, Confidence: 1},
		LineEnding: csvcore.LineEndingInfo{Name: csvcore.LineEndingLF, Sequence: "\n", Confidence: 1},
		Dialect: csvcore.Dialect{
			Delimiter:        ',',
			DelimiterName:    "comma",
			EstimatedColumns: maxRecordWidth(records),
		},
	}
	if _, err := csvcore.SaveRecords(workingPath, metadata, records, csvcore.SaveOptions{}); err != nil {
		return "", err
	}
	return workingPath, nil
}

func (manager *Manager) Add(doc csvcore.Document, kind SourceKind, displayName string, originalPath string, workingPath string) *Session {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.nextID++
	id := strconv.Itoa(manager.nextID)
	session := &Session{
		ID:           id,
		SourceKind:   kind,
		DisplayName:  displayName,
		OriginalPath: originalPath,
		WorkingPath:  workingPath,
		Document:     doc,
		Edits:        make(map[Cell]string),
	}
	session.recordKnownFileState()
	manager.sessions[id] = session
	return cloneSession(session)
}

func (manager *Manager) Window(id string, start int, count int) (csvcore.RecordWindow, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return csvcore.RecordWindow{}, errors.New("document not found")
	}
	if start < 0 {
		start = 0
	}
	if count <= 0 {
		count = 80
	}

	if !session.ViewActive && session.EditBaseRecords == nil {
		window, err := csvcore.ReadRecordWindowWithIndex(
			session.WorkingPath,
			session.Document.Metadata,
			session.Document.RowIndex,
			start,
			count,
		)
		if err != nil {
			return csvcore.RecordWindow{}, err
		}
		applyEdits(window.Rows, start, session.Edits)
		window.TotalRows = session.Document.RowIndex.Rows
		if len(window.SourceRows) == 0 {
			window.SourceRows = make([]int, len(window.Rows))
			for i := range window.Rows {
				window.SourceRows[i] = start + i
			}
		}
		return window, nil
	}

	records, err := recordsForSession(session)
	if err != nil {
		return csvcore.RecordWindow{}, err
	}
	if !session.ViewActive {
		if start > len(records) {
			start = len(records)
		}
		end := min(start+count, len(records))
		rows := cloneRecords(records[start:end])
		sourceRows := make([]int, len(rows))
		for i := range rows {
			sourceRows[i] = start + i
		}
		return csvcore.RecordWindow{Start: start, Rows: rows, TotalRows: len(records), SourceRows: sourceRows}, nil
	}
	if start > len(session.ViewRows) {
		start = len(session.ViewRows)
	}
	end := min(start+count, len(session.ViewRows))
	rows := make([][]string, 0, max(0, end-start))
	for _, sourceRow := range session.ViewRows[start:end] {
		if sourceRow < 0 || sourceRow >= len(records) {
			continue
		}
		record := append([]string(nil), records[sourceRow]...)
		windowRecords := [][]string{record}
		applyEdits(windowRecords, sourceRow, session.Edits)
		rows = append(rows, windowRecords[0])
	}
	sourceRows := append([]int(nil), session.ViewRows[start:end]...)
	return csvcore.RecordWindow{Start: start, Rows: rows, TotalRows: len(session.ViewRows), SourceRows: sourceRows}, nil
}

func (manager *Manager) Records(id string) (*Session, [][]string, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return nil, nil, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return nil, nil, err
	}
	return session, records, nil
}

func (manager *Manager) ExportView(id string, options ExportViewOptions) (ExportViewResult, error) {
	session, records, err := manager.Records(id)
	if err != nil {
		return ExportViewResult{}, err
	}
	for _, column := range options.Columns {
		if column < 0 {
			return ExportViewResult{}, errors.New("column must be non-negative")
		}
	}

	hasHeader := session.Document.Metadata.Dialect.HasHeader && len(records) > 0
	start := 0
	if hasHeader {
		start = 1
	}
	includeHeaders := hasHeader
	if options.IncludeHeaders != nil {
		includeHeaders = hasHeader && *options.IncludeHeaders
	}
	viewRows := append([]int(nil), session.ViewRows...)
	if options.FullDocument || !session.ViewActive {
		viewRows = make([]int, 0, max(0, len(records)-start))
		for row := start; row < len(records); row++ {
			viewRows = append(viewRows, row)
		}
	}
	if len(options.Rows) > 0 {
		slicedRows := make([]int, 0, len(options.Rows))
		for _, displayRow := range options.Rows {
			if displayRow < 0 {
				return ExportViewResult{}, errors.New("row must be non-negative")
			}
			if displayRow >= len(viewRows) {
				return ExportViewResult{}, fmt.Errorf("row %d is outside the current view", displayRow)
			}
			slicedRows = append(slicedRows, viewRows[displayRow])
		}
		viewRows = slicedRows
	}

	exportRecords := make([][]string, 0, len(viewRows)+1)
	if includeHeaders {
		exportRecords = append(exportRecords, projectRecord(records[0], options.Columns))
	}
	exportedRows := 0
	for _, row := range viewRows {
		if row < start || row >= len(records) {
			continue
		}
		exportRecords = append(exportRecords, projectRecord(records[row], options.Columns))
		exportedRows++
	}

	format := normalizeExportViewFormat(options.Format)
	metadata := session.Document.Metadata
	if format == "tsv" {
		metadata.Dialect.Delimiter = '\t'
		metadata.Dialect.DelimiterName = "tab"
	}
	contentType := exportViewContentType(metadata.Dialect.Delimiter)
	filename := exportViewFilename(session.DisplayName, metadata.Dialect.Delimiter)
	var data []byte
	if format == "xlsx" {
		data, err = xlsx.EncodeWorkbook(exportRecords)
		contentType = xlsx.ContentType()
		filename = exportViewXLSXFilename(session.DisplayName)
	} else if isRichExportViewFormat(format) {
		data, err = encodeRichExportView(exportRecords, format)
		contentType = exportViewRichContentType(format)
		filename = exportViewFilenameWithExtension(session.DisplayName, exportViewRichExtension(format))
	} else {
		data, err = csvcore.EncodeRecords(exportRecords, metadata)
	}
	if err != nil {
		return ExportViewResult{}, err
	}
	return ExportViewResult{
		Filename:    filename,
		Data:        data,
		Rows:        exportedRows,
		Columns:     exportColumnCount(exportRecords, options.Columns),
		Delimiter:   session.Document.Metadata.Dialect.Delimiter,
		ContentType: contentType,
	}, nil
}

func (manager *Manager) Edit(id string, row int, col int, value string) error {
	return manager.EditCells(id, []CellEdit{{Row: row, Col: col, Value: value}})
}

type CellEdit struct {
	Row   int    `json:"row"`
	Col   int    `json:"col"`
	Value string `json:"value"`
}

func (manager *Manager) EditCells(id string, edits []CellEdit) error {
	if len(edits) == 0 {
		return nil
	}
	for _, edit := range edits {
		if edit.Row < 0 || edit.Col < 0 {
			return errors.New("row and col must be non-negative")
		}
	}
	if err := manager.ensureWritable(id); err != nil {
		return err
	}
	if err := manager.pushUndo(id); err != nil {
		return err
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return errors.New("document not found")
	}
	if session.EditBaseRecords == nil && len(session.Edits) == 0 && len(session.UndoStack) > 0 {
		session.EditBaseRecords = cloneRecords(session.UndoStack[len(session.UndoStack)-1].Records)
	}
	for _, edit := range edits {
		session.Edits[Cell{Row: edit.Row, Col: edit.Col}] = edit.Value
		if edit.Col+1 > session.Document.Metadata.Dialect.EstimatedColumns {
			session.Document.Metadata.Dialect.EstimatedColumns = edit.Col + 1
		}
	}
	session.Dirty = true
	return nil
}

type MutationResult struct {
	Session   *Session `json:"session"`
	TotalRows int      `json:"totalRows"`
	Dirty     bool     `json:"dirty"`
	Changed   int      `json:"changed,omitempty"`
}

func (manager *Manager) SetHeaderMode(id string, hasHeader bool) (MutationResult, error) {
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	if session.Document.Metadata.Dialect.HasHeader == hasHeader {
		return MutationResult{Session: cloneSession(session), TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	session.Document.Metadata.Dialect.HasHeader = hasHeader
	session.ViewRows = nil
	session.ViewActive = false
	return MutationResult{Session: cloneSession(session), TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty, Changed: 1}, nil
}

func (manager *Manager) SetQuotePolicy(id string, quoteChar string, quoteMode string, quoteEscape string) (MutationResult, error) {
	normalizedQuoteChar, err := csvcore.NormalizeQuoteChar(quoteChar)
	if err != nil {
		return MutationResult{}, err
	}
	normalizedQuoteMode, err := csvcore.NormalizeQuoteMode(quoteMode)
	if err != nil {
		return MutationResult{}, err
	}
	normalizedQuoteEscape, err := csvcore.NormalizeQuoteEscape(quoteEscape)
	if err != nil {
		return MutationResult{}, err
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	if session.ReadOnly {
		return MutationResult{}, errors.New("document is read-only")
	}
	if session.Document.Metadata.Dialect.QuoteChar == normalizedQuoteChar &&
		csvcore.DialectQuoteMode(session.Document.Metadata.Dialect) == normalizedQuoteMode &&
		csvcore.DialectQuoteEscape(session.Document.Metadata.Dialect) == normalizedQuoteEscape {
		return MutationResult{Session: cloneSession(session), TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	session.Document.Metadata.Dialect.QuoteChar = normalizedQuoteChar
	session.Document.Metadata.Dialect.QuoteMode = normalizedQuoteMode
	session.Document.Metadata.Dialect.QuoteEscape = normalizedQuoteEscape
	session.Dirty = true
	return MutationResult{Session: cloneSession(session), TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty, Changed: 1}, nil
}

func (manager *Manager) SetLineEnding(id string, lineEnding string) (MutationResult, error) {
	normalized, err := csvcore.NormalizeLineEndingName(lineEnding)
	if err != nil {
		return MutationResult{}, err
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	if session.ReadOnly {
		return MutationResult{}, errors.New("document is read-only")
	}
	if session.Document.Metadata.LineEnding.Name == normalized.Name {
		return MutationResult{Session: cloneSession(session), TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	session.Document.Metadata.LineEnding = normalized
	session.Dirty = true
	return MutationResult{Session: cloneSession(session), TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty, Changed: 1}, nil
}

func (manager *Manager) SetEncoding(id string, encodingName string) (MutationResult, error) {
	encoding, err := csvcore.ParseEncodingName(encodingName)
	if err != nil {
		return MutationResult{}, err
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	current := session.Document.Metadata.Encoding
	if current.Name == encoding.Name && current.HasBOM == encoding.HasBOM {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	metadata := session.Document.Metadata
	metadata.Encoding = encoding
	metadata.Dialect.UTF8BOM = encoding.Name == csvcore.EncodingUTF8BOM
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: 1}, nil
}

func (manager *Manager) SetDelimiter(id string, delimiterValue string, delimiterName string) (MutationResult, error) {
	delimiter, name, err := csvcore.NormalizeDelimiter(delimiterValue, delimiterName)
	if err != nil {
		return MutationResult{}, err
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	if session.Document.Metadata.Dialect.Delimiter == delimiter && session.Document.Metadata.Dialect.DelimiterName == name {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	metadata := session.Document.Metadata
	metadata.Dialect.Delimiter = delimiter
	metadata.Dialect.DelimiterName = name
	metadata.Dialect.EstimatedColumns = maxRecordWidth(records)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: 1}, nil
}

func (manager *Manager) SetCellQuote(id string, row int, col int, quoted bool) (MutationResult, error) {
	if row < 0 || col < 0 {
		return MutationResult{}, errors.New("cell coordinates must be non-negative")
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	if session.ReadOnly {
		return MutationResult{}, errors.New("document is read-only")
	}
	nextDialect, changed := csvcore.SetCellQuote(session.Document.Metadata.Dialect, row, col, quoted)
	if !changed {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	current.Document.Metadata.Dialect = nextDialect
	current.Dirty = true
	return MutationResult{Session: cloneSession(current), TotalRows: current.Document.RowIndex.Rows, Dirty: current.Dirty, Changed: 1}, nil
}

func (manager *Manager) SetQuoteRule(id string, target string, col int, quoted bool) (MutationResult, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	if session.ReadOnly {
		return MutationResult{}, errors.New("document is read-only")
	}
	nextDialect := session.Document.Metadata.Dialect
	var changed bool
	switch strings.ToLower(strings.TrimSpace(target)) {
	case "header":
		nextDialect, changed = csvcore.SetQuoteHeader(nextDialect, quoted)
	case "body":
		nextDialect, changed = csvcore.SetQuoteBody(nextDialect, quoted)
	case "column":
		if col < 0 {
			return MutationResult{}, errors.New("column must be non-negative")
		}
		nextDialect, changed = csvcore.SetQuoteColumn(nextDialect, col, quoted)
	default:
		return MutationResult{}, errors.New("quote rule target must be header, body, or column")
	}
	if !changed {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	current.Document.Metadata.Dialect = nextDialect
	current.Dirty = true
	return MutationResult{Session: cloneSession(current), TotalRows: current.Document.RowIndex.Rows, Dirty: current.Dirty, Changed: 1}, nil
}

func (manager *Manager) SetQuoteSettings(id string, settings QuoteSettings) (MutationResult, error) {
	normalizedQuoteChar, err := csvcore.NormalizeQuoteChar(settings.QuoteChar)
	if err != nil {
		return MutationResult{}, err
	}
	normalizedQuoteMode, err := csvcore.NormalizeQuoteMode(settings.QuoteMode)
	if err != nil {
		return MutationResult{}, err
	}
	normalizedQuoteEscape, err := csvcore.NormalizeQuoteEscape(settings.QuoteEscape)
	if err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	if session.ReadOnly {
		return MutationResult{}, errors.New("document is read-only")
	}
	nextDialect := session.Document.Metadata.Dialect
	nextDialect.QuoteChar = normalizedQuoteChar
	nextDialect.QuoteMode = normalizedQuoteMode
	nextDialect.QuoteEscape = normalizedQuoteEscape
	nextDialect.QuoteHeader = settings.QuoteHeader
	nextDialect.QuoteBody = settings.QuoteBody
	nextDialect.QuoteColumns = normalizeQuoteColumns(settings.QuoteColumns)
	nextDialect.QuoteCells = normalizeQuoteCells(settings.QuoteCells)
	nextDialect.NoQuoteCells = normalizeQuoteCells(settings.NoQuoteCells)
	if reflect.DeepEqual(session.Document.Metadata.Dialect, nextDialect) {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	current.Document.Metadata.Dialect = nextDialect
	current.Dirty = true
	return MutationResult{Session: cloneSession(current), TotalRows: current.Document.RowIndex.Rows, Dirty: current.Dirty, Changed: 1}, nil
}

func normalizeQuoteColumns(columns []int) []int {
	seen := make(map[int]bool)
	normalized := make([]int, 0, len(columns))
	for _, col := range columns {
		if col < 0 || seen[col] {
			continue
		}
		seen[col] = true
		normalized = append(normalized, col)
	}
	sort.Ints(normalized)
	return normalized
}

func normalizeQuoteCells(cells []csvcore.Cell) []csvcore.Cell {
	seen := make(map[csvcore.Cell]bool)
	normalized := make([]csvcore.Cell, 0, len(cells))
	for _, cell := range cells {
		if cell.Row < 0 || cell.Col < 0 || seen[cell] {
			continue
		}
		seen[cell] = true
		normalized = append(normalized, cell)
	}
	sort.Slice(normalized, func(i, j int) bool {
		if normalized[i].Row == normalized[j].Row {
			return normalized[i].Col < normalized[j].Col
		}
		return normalized[i].Row < normalized[j].Row
	})
	return normalized
}

type ExportViewOptions struct {
	Columns        []int `json:"columns"`
	Rows           []int `json:"rows"`
	Format         string
	IncludeHeaders *bool
	FullDocument   bool
}

type ExportViewResult struct {
	Filename    string
	Data        []byte
	Rows        int
	Columns     int
	Delimiter   rune
	ContentType string
}

func (manager *Manager) Undo(id string) (MutationResult, error) {
	manager.mu.Lock()
	session, ok := manager.sessions[id]
	if !ok {
		manager.mu.Unlock()
		return MutationResult{}, errors.New("document not found")
	}
	if session.ReadOnly {
		manager.mu.Unlock()
		return MutationResult{}, errors.New("document is read-only")
	}
	if len(session.UndoStack) == 0 {
		manager.mu.Unlock()
		return MutationResult{Session: cloneSession(session), TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	current := cloneSession(session)
	target := session.UndoStack[len(session.UndoStack)-1]
	session.UndoStack = session.UndoStack[:len(session.UndoStack)-1]
	manager.mu.Unlock()

	redo, err := snapshotFromSession(current)
	if err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.applySnapshot(id, target, func(session *Session) {
		session.RedoStack = append(session.RedoStack, redo)
	})
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty}, nil
}

func (manager *Manager) Redo(id string) (MutationResult, error) {
	manager.mu.Lock()
	session, ok := manager.sessions[id]
	if !ok {
		manager.mu.Unlock()
		return MutationResult{}, errors.New("document not found")
	}
	if session.ReadOnly {
		manager.mu.Unlock()
		return MutationResult{}, errors.New("document is read-only")
	}
	if len(session.RedoStack) == 0 {
		manager.mu.Unlock()
		return MutationResult{Session: cloneSession(session), TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	current := cloneSession(session)
	target := session.RedoStack[len(session.RedoStack)-1]
	session.RedoStack = session.RedoStack[:len(session.RedoStack)-1]
	manager.mu.Unlock()

	undo, err := snapshotFromSession(current)
	if err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.applySnapshot(id, target, func(session *Session) {
		session.UndoStack = append(session.UndoStack, undo)
	})
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty}, nil
}

func (manager *Manager) InsertRow(id string, row int) (MutationResult, error) {
	if row < 0 {
		row = 0
	}
	return manager.mutateRecords(id, func(records [][]string) ([][]string, error) {
		if row > len(records) {
			row = len(records)
		}
		width := maxRecordWidth(records)
		inserted := make([]string, width)
		records = append(records, nil)
		copy(records[row+1:], records[row:])
		records[row] = inserted
		return records, nil
	})
}

func (manager *Manager) DeleteRow(id string, row int) (MutationResult, error) {
	return manager.mutateRecords(id, func(records [][]string) ([][]string, error) {
		if row < 0 || row >= len(records) {
			return records, errors.New("row is out of range")
		}
		return append(records[:row], records[row+1:]...), nil
	})
}

func (manager *Manager) MoveRow(id string, row int, delta int) (MutationResult, error) {
	if delta == 0 {
		session, ok := manager.snapshot(id)
		if !ok {
			return MutationResult{}, errors.New("document not found")
		}
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	if row < 0 || row >= len(records) {
		return MutationResult{}, errors.New("row is out of range")
	}
	minRow := 0
	if session.Document.Metadata.Dialect.HasHeader {
		minRow = 1
	}
	target := row + delta
	if target < minRow || target >= len(records) {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	records = moveRecord(records, row, target)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, session.Document.Metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: 1}, nil
}

func (manager *Manager) PasteRow(id string, sourceRow int, targetRow int, cut bool) (MutationResult, error) {
	return manager.PasteRows(id, []int{sourceRow}, targetRow, cut)
}

func (manager *Manager) PasteRows(id string, sourceRows []int, targetRow int, cut bool) (MutationResult, error) {
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	records = cloneRecords(records)
	sourceRows = uniqueIndexesInOrder(sourceRows)
	if len(sourceRows) == 0 {
		return MutationResult{}, errors.New("source rows are required")
	}
	minRow := 0
	if session.Document.Metadata.Dialect.HasHeader {
		minRow = 1
	}
	for _, sourceRow := range sourceRows {
		if sourceRow < 0 || sourceRow >= len(records) {
			return MutationResult{}, errors.New("source row is out of range")
		}
		if sourceRow < minRow {
			return MutationResult{}, errors.New("header row cannot be pasted")
		}
	}
	if targetRow < minRow {
		targetRow = minRow
	}
	if targetRow > len(records) {
		targetRow = len(records)
	}
	before := cloneRecords(records)
	pasted := make([][]string, 0, len(sourceRows))
	for _, sourceRow := range sourceRows {
		pasted = append(pasted, append([]string(nil), records[sourceRow]...))
	}
	if cut {
		sourceSet := indexSet(sourceRows)
		targetRow -= countIndexesBefore(sourceRows, targetRow)
		filtered := make([][]string, 0, len(records)-len(sourceSet))
		for row, record := range records {
			if sourceSet[row] {
				continue
			}
			filtered = append(filtered, record)
		}
		records = filtered
	}
	if targetRow < minRow {
		targetRow = minRow
	}
	if targetRow > len(records) {
		targetRow = len(records)
	}
	records = insertRecords(records, targetRow, pasted)
	if reflect.DeepEqual(before, records) {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	metadata := session.Document.Metadata
	metadata.Dialect.EstimatedColumns = maxRecordWidth(records)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: len(sourceRows)}, nil
}

func (manager *Manager) HideRow(id string, row int) (MutationResult, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	if row < 0 || row >= len(records) {
		return MutationResult{}, errors.New("row is out of range")
	}
	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	if row < start {
		return MutationResult{Session: session, TotalRows: visibleRowCount(session, len(records), start), Dirty: session.Dirty}, nil
	}
	viewRows := append([]int(nil), session.ViewRows...)
	if !session.ViewActive {
		viewRows = make([]int, 0, max(0, len(records)-start))
		for sourceRow := start; sourceRow < len(records); sourceRow++ {
			viewRows = append(viewRows, sourceRow)
		}
	}
	filtered := viewRows[:0]
	changed := 0
	for _, sourceRow := range viewRows {
		if sourceRow == row {
			changed++
			continue
		}
		filtered = append(filtered, sourceRow)
	}
	if changed == 0 {
		return MutationResult{Session: session, TotalRows: len(viewRows), Dirty: session.Dirty}, nil
	}
	manager.mu.Lock()
	current, ok := manager.sessions[id]
	if !ok {
		manager.mu.Unlock()
		return MutationResult{}, errors.New("document not found")
	}
	current.ViewRows = append([]int(nil), filtered...)
	current.ViewActive = true
	updated := cloneSession(current)
	manager.mu.Unlock()
	return MutationResult{Session: updated, TotalRows: len(filtered), Dirty: updated.Dirty, Changed: changed}, nil
}

func (manager *Manager) InsertColumn(id string, col int) (MutationResult, error) {
	if col < 0 {
		col = 0
	}
	return manager.mutateRecords(id, func(records [][]string) ([][]string, error) {
		width := maxRecordWidth(records)
		if col > width {
			col = width
		}
		for row := range records {
			for len(records[row]) < width {
				records[row] = append(records[row], "")
			}
			records[row] = append(records[row], "")
			copy(records[row][col+1:], records[row][col:])
			records[row][col] = ""
		}
		return records, nil
	})
}

func (manager *Manager) DeleteColumn(id string, col int) (MutationResult, error) {
	return manager.mutateRecords(id, func(records [][]string) ([][]string, error) {
		if col < 0 || col >= maxRecordWidth(records) {
			return records, errors.New("column is out of range")
		}
		for row := range records {
			if col >= len(records[row]) {
				continue
			}
			records[row] = append(records[row][:col], records[row][col+1:]...)
		}
		return records, nil
	})
}

func (manager *Manager) MoveColumn(id string, col int, delta int) (MutationResult, error) {
	if delta == 0 {
		session, ok := manager.snapshot(id)
		if !ok {
			return MutationResult{}, errors.New("document not found")
		}
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	width := maxRecordWidth(records)
	if col < 0 || col >= width {
		return MutationResult{}, errors.New("column is out of range")
	}
	target := col + delta
	if target < 0 || target >= width {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	records = moveColumnRecords(records, col, target, width)
	metadata := session.Document.Metadata
	metadata.Dialect.EstimatedColumns = maxRecordWidth(records)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: 1}, nil
}

func (manager *Manager) PasteColumn(id string, sourceCol int, targetCol int, cut bool) (MutationResult, error) {
	return manager.PasteColumns(id, []int{sourceCol}, targetCol, cut)
}

func (manager *Manager) PasteColumns(id string, sourceCols []int, targetCol int, cut bool) (MutationResult, error) {
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	records = cloneRecords(records)
	width := maxRecordWidth(records)
	sourceCols = uniqueIndexesInOrder(sourceCols)
	if len(sourceCols) == 0 {
		return MutationResult{}, errors.New("source columns are required")
	}
	for _, sourceCol := range sourceCols {
		if sourceCol < 0 || sourceCol >= width {
			return MutationResult{}, errors.New("source column is out of range")
		}
	}
	if targetCol < 0 {
		targetCol = 0
	}
	if targetCol > width {
		targetCol = width
	}
	before := cloneRecords(records)
	values := make([][]string, len(records))
	for row := range records {
		ensureRecordWidth(&records[row], width)
		values[row] = make([]string, len(sourceCols))
		for index, sourceCol := range sourceCols {
			values[row][index] = records[row][sourceCol]
		}
	}
	if cut {
		sourceSet := indexSet(sourceCols)
		for row := range records {
			filtered := make([]string, 0, width-len(sourceSet))
			for col, value := range records[row] {
				if sourceSet[col] {
					continue
				}
				filtered = append(filtered, value)
			}
			records[row] = filtered
		}
		targetCol -= countIndexesBefore(sourceCols, targetCol)
		width -= len(sourceCols)
	}
	if targetCol < 0 {
		targetCol = 0
	}
	if targetCol > width {
		targetCol = width
	}
	for row := range records {
		ensureRecordWidth(&records[row], width)
		records[row] = insertRecordValues(records[row], targetCol, values[row])
	}
	if reflect.DeepEqual(before, records) {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	metadata := session.Document.Metadata
	metadata.Dialect.EstimatedColumns = maxRecordWidth(records)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: len(sourceCols)}, nil
}

func (manager *Manager) Transpose(id string, maxCells int) (MutationResult, error) {
	if maxCells <= 0 {
		maxCells = defaultTransposeCellLimit
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	width := maxRecordWidth(records)
	cellCount := len(records) * width
	if cellCount > maxCells {
		return MutationResult{}, fmt.Errorf("transpose would touch %d cells; limit is %d", cellCount, maxCells)
	}
	transposed := transposeRecords(records, width)
	metadata := session.Document.Metadata
	metadata.Dialect.HasHeader = false
	metadata.Dialect.HeaderColumns = 0
	metadata.Dialect.EstimatedColumns = len(records)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, transposed, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty}, nil
}

type TextTransformMode string

const (
	TransformTrim   TextTransformMode = "trim"
	TransformUpper  TextTransformMode = "upper"
	TransformLower  TextTransformMode = "lower"
	TransformSnake  TextTransformMode = "snake"
	TransformKebab  TextTransformMode = "kebab"
	TransformCamel  TextTransformMode = "camel"
	TransformPascal TextTransformMode = "pascal"
	TransformTitle  TextTransformMode = "title"
	TransformWide   TextTransformMode = "wide"
	TransformNarrow TextTransformMode = "narrow"
)

func (manager *Manager) DeduplicateRows(id string) (MutationResult, error) {
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}

	deduped, removed := deduplicateRecords(records, session.Document.Metadata.Dialect.HasHeader)
	if removed == 0 {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, deduped, session.Document.Metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: removed}, nil
}

func (manager *Manager) DeduplicateColumns(id string) (MutationResult, error) {
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}

	deduped, removed := deduplicateColumns(records)
	if removed == 0 {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	metadata := session.Document.Metadata
	metadata.Dialect.EstimatedColumns = maxRecordWidth(deduped)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, deduped, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: removed}, nil
}

func (manager *Manager) TransformCells(id string, cells []Cell, mode TextTransformMode) (MutationResult, error) {
	if len(cells) == 0 {
		return MutationResult{}, errors.New("no cells selected")
	}
	transform, err := textTransformFunc(mode)
	if err != nil {
		return MutationResult{}, err
	}
	for _, cell := range cells {
		if cell.Row < 0 || cell.Col < 0 {
			return MutationResult{}, errors.New("row and col must be non-negative")
		}
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}

	changed := 0
	seen := make(map[Cell]bool, len(cells))
	for _, cell := range cells {
		if seen[cell] {
			continue
		}
		seen[cell] = true
		if cell.Row >= len(records) {
			continue
		}
		ensureRecordWidth(&records[cell.Row], cell.Col+1)
		next := transform(records[cell.Row][cell.Col])
		if next != records[cell.Row][cell.Col] {
			changed++
			records[cell.Row][cell.Col] = next
		}
	}
	if changed == 0 {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, session.Document.Metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: changed}, nil
}

func (manager *Manager) TransformColumn(id string, col int, mode TextTransformMode) (MutationResult, error) {
	if col < 0 {
		return MutationResult{}, errors.New("column must be non-negative")
	}
	transform, err := textTransformFunc(mode)
	if err != nil {
		return MutationResult{}, err
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	startRow := 0
	if session.Document.Metadata.Dialect.HasHeader {
		startRow = 1
	}
	changed := 0
	for row := startRow; row < len(records); row++ {
		if col >= len(records[row]) {
			continue
		}
		next := transform(records[row][col])
		if next != records[row][col] {
			changed++
			records[row][col] = next
		}
	}
	if changed == 0 {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, session.Document.Metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: changed}, nil
}

func (manager *Manager) SwapCells(id string, cells []Cell) (MutationResult, error) {
	if len(cells) != 2 {
		return MutationResult{}, errors.New("select exactly two cells to swap")
	}
	first := cells[0]
	second := cells[1]
	if first.Row < 0 || first.Col < 0 || second.Row < 0 || second.Col < 0 {
		return MutationResult{}, errors.New("row and col must be non-negative")
	}
	if first == second {
		return MutationResult{}, errors.New("select two different cells to swap")
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	if first.Row >= len(records) || second.Row >= len(records) {
		return MutationResult{}, errors.New("swap row is out of range")
	}
	firstValue := cellValue(records[first.Row], first.Col)
	secondValue := cellValue(records[second.Row], second.Col)
	if firstValue == secondValue {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	ensureRecordWidth(&records[first.Row], first.Col+1)
	ensureRecordWidth(&records[second.Row], second.Col+1)
	records[first.Row][first.Col], records[second.Row][second.Col] = secondValue, firstValue
	metadata := session.Document.Metadata
	metadata.Dialect.EstimatedColumns = max(metadata.Dialect.EstimatedColumns, maxRecordWidth(records))
	updated, err := manager.materializeRecords(id, records, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: 2}, nil
}

func (manager *Manager) JoinCells(id string, cells []Cell, separator string) (MutationResult, error) {
	return manager.joinCells(id, cells, separator, "horizontal")
}

func (manager *Manager) JoinCellsVertical(id string, cells []Cell, separator string) (MutationResult, error) {
	return manager.joinCells(id, cells, separator, "vertical")
}

func (manager *Manager) joinCells(id string, cells []Cell, separator string, direction string) (MutationResult, error) {
	if len(cells) == 0 {
		return MutationResult{}, errors.New("no cells selected")
	}
	if separator == "" {
		return MutationResult{}, errors.New("separator is required")
	}
	for _, cell := range cells {
		if cell.Row < 0 || cell.Col < 0 {
			return MutationResult{}, errors.New("row and col must be non-negative")
		}
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}

	changed := joinRecords(records, cells, separator, direction)
	if changed == 0 {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, session.Document.Metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: changed}, nil
}

func (manager *Manager) SplitCells(id string, cells []Cell, separator string) (MutationResult, error) {
	return manager.splitCells(id, cells, separator, "horizontal")
}

func (manager *Manager) SplitCellsVertical(id string, cells []Cell, separator string) (MutationResult, error) {
	return manager.splitCells(id, cells, separator, "vertical")
}

func (manager *Manager) splitCells(id string, cells []Cell, separator string, direction string) (MutationResult, error) {
	if len(cells) == 0 {
		return MutationResult{}, errors.New("no cells selected")
	}
	if separator == "" {
		return MutationResult{}, errors.New("separator is required")
	}
	for _, cell := range cells {
		if cell.Row < 0 || cell.Col < 0 {
			return MutationResult{}, errors.New("row and col must be non-negative")
		}
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}

	changed := splitRecords(&records, cells, separator, direction)
	if changed == 0 {
		return MutationResult{Session: session, TotalRows: session.Document.RowIndex.Rows, Dirty: session.Dirty}, nil
	}
	metadata := session.Document.Metadata
	metadata.Dialect.EstimatedColumns = maxRecordWidth(records)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: changed}, nil
}

func (manager *Manager) Save(id string, path string, createBackup bool) (csvcore.SaveResult, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return csvcore.SaveResult{}, errors.New("document not found")
	}
	if session.ReadOnly {
		return csvcore.SaveResult{}, errors.New("document is read-only")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return csvcore.SaveResult{}, err
	}

	outputPath := path
	if outputPath == "" {
		outputPath = session.OriginalPath
	}
	if outputPath == "" {
		outputPath = session.WorkingPath
	}
	if outputPath != session.OriginalPath {
		if _, statErr := os.Stat(outputPath); errors.Is(statErr, os.ErrNotExist) {
			createBackup = false
		} else if statErr != nil {
			return csvcore.SaveResult{}, statErr
		}
	}
	var result csvcore.SaveResult
	if isXLSXPath(outputPath) {
		sourceWorkbookPath := ""
		if isXLSXPath(session.OriginalPath) {
			sourceWorkbookPath = session.OriginalPath
		}
		result, err = saveXLSXRecords(outputPath, records, createBackup, selectedWorkbookSheetPath(session.Document.Metadata), sourceWorkbookPath)
	} else {
		result, err = csvcore.SaveRecords(outputPath, session.Document.Metadata, records, csvcore.SaveOptions{CreateBackup: createBackup})
	}
	if err != nil {
		return csvcore.SaveResult{}, err
	}

	delimiter := session.Document.Metadata.Dialect.Delimiter
	delimiterName := session.Document.Metadata.Dialect.DelimiterName
	hasHeader := session.Document.Metadata.Dialect.HasHeader
	headerColumns := session.Document.Metadata.Dialect.HeaderColumns
	refreshed, refreshedWorkingPath, openErr := openDocumentForPath(outputPath, csvcore.OpenOptions{
		PreviewRows:   max(1, len(session.Document.Preview)),
		IndexRows:     max(1, len(session.Document.Preview)),
		Delimiter:     &delimiter,
		DelimiterName: delimiterName,
		HasHeader:     &hasHeader,
		HeaderColumns: &headerColumns,
		Encoding:      &session.Document.Metadata.Encoding,
		QuoteChar:     &session.Document.Metadata.Dialect.QuoteChar,
		QuoteMode:     session.Document.Metadata.Dialect.QuoteMode,
		QuoteEscape:   session.Document.Metadata.Dialect.QuoteEscape,
		SheetName:     selectedWorkbookSheetName(session.Document.Metadata),
		SheetIndex:    selectedWorkbookSheetIndex(session.Document.Metadata),
	})
	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if ok {
		current.Edits = make(map[Cell]string)
		current.EditBaseRecords = nil
		current.Dirty = false
		current.ViewRows = nil
		current.ViewActive = false
		if openErr == nil {
			current.Document = refreshed
			current.WorkingPath = refreshedWorkingPath
			current.OriginalPath = outputPath
		}
		current.recordKnownFileState()
	}
	return result, nil
}

func saveXLSXRecords(path string, records [][]string, createBackup bool, sheetPath string, sourceWorkbookPath string) (csvcore.SaveResult, error) {
	var backupPath string
	if createBackup {
		var err error
		backupPath, err = csvcore.CreateBackup(path, "")
		if err != nil {
			return csvcore.SaveResult{}, err
		}
	}
	var data []byte
	var err error
	if sheetPath != "" {
		var existing []byte
		existing, err = os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) && sourceWorkbookPath != "" && sourceWorkbookPath != path {
			existing, err = os.ReadFile(sourceWorkbookPath)
		}
		if err == nil {
			data, err = xlsx.UpdateWorkbookSheet(existing, sheetPath, records)
			if err != nil {
				return csvcore.SaveResult{}, err
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return csvcore.SaveResult{}, err
		}
	}
	if data == nil {
		data, err = xlsx.EncodeWorkbook(records)
		if err != nil {
			return csvcore.SaveResult{}, err
		}
	}
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return csvcore.SaveResult{}, err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return csvcore.SaveResult{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return csvcore.SaveResult{}, err
	}
	return csvcore.SaveResult{Path: path, BackupPath: backupPath, SizeBytes: info.Size()}, nil
}

func selectedWorkbookSheetPath(metadata csvcore.FileMetadata) string {
	if metadata.Workbook == nil {
		return ""
	}
	return metadata.Workbook.SelectedSheetPath
}

func selectedWorkbookSheetName(metadata csvcore.FileMetadata) string {
	if metadata.Workbook == nil {
		return ""
	}
	return metadata.Workbook.SelectedSheetName
}

func selectedWorkbookSheetIndex(metadata csvcore.FileMetadata) int {
	if metadata.Workbook == nil {
		return 0
	}
	return metadata.Workbook.SelectedSheetIndex
}

func (manager *Manager) ListBackups(id string) ([]BackupInfo, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return nil, errors.New("document not found")
	}
	path := session.trackedPath()
	if path == "" {
		return nil, errors.New("document has no file path")
	}
	dir := filepath.Dir(path)
	pattern := backupNamePattern(path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	backups := make([]BackupInfo, 0)
	for _, entry := range entries {
		if !pattern.MatchString(entry.Name()) {
			continue
		}
		info, err := entry.Info()
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		backups = append(backups, BackupInfo{
			Path:            filepath.Join(dir, entry.Name()),
			Name:            entry.Name(),
			SizeBytes:       info.Size(),
			ModTimeUnixNano: info.ModTime().UnixNano(),
		})
	}
	sort.Slice(backups, func(left, right int) bool {
		if backups[left].ModTimeUnixNano == backups[right].ModTimeUnixNano {
			return backups[left].Name > backups[right].Name
		}
		return backups[left].ModTimeUnixNano > backups[right].ModTimeUnixNano
	})
	return backups, nil
}

func backupNamePattern(path string) *regexp.Regexp {
	ext := filepath.Ext(path)
	stem := strings.TrimSuffix(filepath.Base(path), ext)
	namePattern := fmt.Sprintf(`^%s\.\d{8}T\d{6}Z(?:\.\d+)?%s\.bak$`, regexp.QuoteMeta(stem), regexp.QuoteMeta(ext))
	return regexp.MustCompile(namePattern)
}

func (manager *Manager) RestoreBackup(id string, backupPath string) (MutationResult, error) {
	if strings.TrimSpace(backupPath) == "" {
		return MutationResult{}, errors.New("backup path is required")
	}
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	cleanBackupPath := filepath.Clean(backupPath)
	backups, err := manager.ListBackups(id)
	if err != nil {
		return MutationResult{}, err
	}
	allowed := false
	for _, backup := range backups {
		if filepath.Clean(backup.Path) == cleanBackupPath {
			allowed = true
			break
		}
	}
	if !allowed {
		return MutationResult{}, errors.New("backup does not belong to this document")
	}

	var records [][]string
	if isXLSXPath(session.trackedPath()) {
		records, _, _, err = xlsx.DecodeFileWithOptions(cleanBackupPath, xlsx.DecodeOptions{
			SheetName:  selectedWorkbookSheetName(session.Document.Metadata),
			SheetIndex: selectedWorkbookSheetIndex(session.Document.Metadata),
		})
	} else {
		records, err = csvcore.ReadAllRecords(cleanBackupPath, session.Document.Metadata)
	}
	if err != nil {
		return MutationResult{}, err
	}
	metadata := session.Document.Metadata
	metadata.Dialect.EstimatedColumns = maxRecordWidth(records)
	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty, Changed: len(records)}, nil
}

type FindOptions struct {
	Query         string `json:"query"`
	Column        *int   `json:"column,omitempty"`
	Cells         []Cell `json:"cells,omitempty"`
	CaseSensitive bool   `json:"caseSensitive"`
	WholeCell     bool   `json:"wholeCell"`
	Regex         bool   `json:"regex"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
}

type Match struct {
	Row   int    `json:"row"`
	Col   int    `json:"col"`
	Value string `json:"value"`
}

type FindResult struct {
	Matches   []Match `json:"matches"`
	Total     int     `json:"total"`
	Offset    int     `json:"offset"`
	Truncated bool    `json:"truncated"`
	Streamed  bool    `json:"streamed,omitempty"`
}

func (manager *Manager) Find(id string, options FindOptions) ([]Match, error) {
	result, err := manager.FindDetailed(id, options)
	if err != nil {
		return nil, err
	}
	return result.Matches, nil
}

func (manager *Manager) FindDetailed(id string, options FindOptions) (FindResult, error) {
	return manager.FindDetailedContext(context.Background(), id, options)
}

func (manager *Manager) FindDetailedContext(ctx context.Context, id string, options FindOptions) (FindResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return FindResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return FindResult{}, errors.New("document not found")
	}
	if options.Query == "" {
		return FindResult{}, nil
	}
	if options.Limit <= 0 {
		options.Limit = 100
	}
	if options.Offset < 0 {
		options.Offset = 0
	}
	matcher, err := newMatcher(options)
	if err != nil {
		return FindResult{}, err
	}
	if canStreamFind(session) {
		return findDetailedInStream(ctx, session, options, matcher)
	}

	records, err := recordsForSession(session)
	if err != nil {
		return FindResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return FindResult{}, err
	}

	var matches []Match
	total := 0
	scope := cellScope(options.Cells)
	for row, record := range records {
		if err := ctx.Err(); err != nil {
			return FindResult{}, err
		}
		for col, value := range record {
			if options.Column != nil && col != *options.Column {
				continue
			}
			if !cellInScope(scope, row, col) {
				continue
			}
			if matcher.matches(value) {
				total++
				if total <= options.Offset {
					continue
				}
				if len(matches) < options.Limit {
					matches = append(matches, Match{Row: row, Col: col, Value: value})
				}
			}
		}
	}
	return FindResult{Matches: matches, Total: total, Offset: options.Offset, Truncated: total > options.Offset+len(matches)}, nil
}

func canStreamFind(session *Session) bool {
	return session.EditBaseRecords == nil &&
		!session.ViewActive &&
		len(session.Edits) == 0 &&
		csvcore.CanStreamRecords(session.Document.Metadata)
}

func findDetailedInStream(ctx context.Context, session *Session, options FindOptions, matcher matcher) (FindResult, error) {
	var matches []Match
	total := 0
	scope := cellScope(options.Cells)
	_, err := csvcore.ForEachRecordContext(ctx, session.WorkingPath, session.Document.Metadata, func(row int, record []string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		for col, value := range record {
			if options.Column != nil && col != *options.Column {
				continue
			}
			if !cellInScope(scope, row, col) {
				continue
			}
			if matcher.matches(value) {
				total++
				if total <= options.Offset {
					continue
				}
				if len(matches) < options.Limit {
					matches = append(matches, Match{Row: row, Col: col, Value: value})
				}
			}
		}
		return nil
	})
	if err != nil {
		return FindResult{}, err
	}
	return FindResult{Matches: matches, Total: total, Offset: options.Offset, Truncated: total > options.Offset+len(matches), Streamed: true}, nil
}

type ReplaceOptions struct {
	FindOptions
	Replacement string `json:"replacement"`
	All         bool   `json:"all"`
	Preview     bool   `json:"preview"`
}

type ReplaceResult struct {
	Changes []Match `json:"changes"`
	Count   int     `json:"count"`
	Dirty   bool    `json:"dirty"`
	Preview bool    `json:"preview"`
}

func (manager *Manager) Replace(id string, options ReplaceOptions) (ReplaceResult, error) {
	if options.Query == "" {
		return ReplaceResult{}, nil
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return ReplaceResult{}, errors.New("document not found")
	}
	if !options.Preview {
		if err := manager.ensureWritable(id); err != nil {
			return ReplaceResult{}, err
		}
	}
	matcher, err := newMatcher(options.FindOptions)
	if err != nil {
		return ReplaceResult{}, err
	}

	records, err := recordsForSession(session)
	if err != nil {
		return ReplaceResult{}, err
	}

	edits := make(map[Cell]string)
	var changes []Match
	scope := cellScope(options.Cells)
	for row, record := range records {
		for col, value := range record {
			if options.Column != nil && col != *options.Column {
				continue
			}
			if !cellInScope(scope, row, col) {
				continue
			}
			if !matcher.matches(value) {
				continue
			}
			replaced := matcher.replace(value, options.Replacement, options.All)
			if replaced == value {
				continue
			}
			edits[Cell{Row: row, Col: col}] = replaced
			changes = append(changes, Match{Row: row, Col: col, Value: replaced})
			if !options.All {
				if options.Preview {
					return ReplaceResult{Changes: changes, Count: len(changes), Dirty: session.Dirty, Preview: true}, nil
				}
				return manager.applyReplaceEdits(id, edits, changes)
			}
		}
	}
	if options.Preview {
		return ReplaceResult{Changes: changes, Count: len(changes), Dirty: session.Dirty, Preview: true}, nil
	}
	return manager.applyReplaceEdits(id, edits, changes)
}

func (manager *Manager) applyReplaceEdits(id string, edits map[Cell]string, changes []Match) (ReplaceResult, error) {
	if len(edits) > 0 {
		if err := manager.pushUndo(id); err != nil {
			return ReplaceResult{}, err
		}
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return ReplaceResult{}, errors.New("document not found")
	}
	for key, value := range edits {
		session.Edits[key] = value
	}
	if len(edits) > 0 {
		session.Dirty = true
	}
	return ReplaceResult{Changes: changes, Count: len(changes), Dirty: session.Dirty}, nil
}

func cellScope(cells []Cell) map[Cell]struct{} {
	if len(cells) == 0 {
		return nil
	}
	scope := make(map[Cell]struct{}, len(cells))
	for _, cell := range cells {
		if cell.Row < 0 || cell.Col < 0 {
			continue
		}
		scope[cell] = struct{}{}
	}
	return scope
}

func cellInScope(scope map[Cell]struct{}, row int, col int) bool {
	if scope == nil {
		return true
	}
	_, ok := scope[Cell{Row: row, Col: col}]
	return ok
}

type FilterOptions struct {
	Query         string `json:"query"`
	Column        *int   `json:"column,omitempty"`
	CaseSensitive bool   `json:"caseSensitive"`
}

type FilterResult struct {
	TotalRows int  `json:"totalRows"`
	Streamed  bool `json:"streamed,omitempty"`
}

type ColumnValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
	Blank bool   `json:"blank"`
}

type ColumnValuesResult struct {
	Values        []ColumnValue `json:"values"`
	TotalDistinct int           `json:"totalDistinct"`
	Streamed      bool          `json:"streamed,omitempty"`
}

type ValueFilterOptions struct {
	Column       int      `json:"column"`
	Values       []string `json:"values"`
	IncludeBlank bool     `json:"includeBlank"`
}

func (manager *Manager) FilterContains(id string, options FilterOptions) (int, error) {
	result, err := manager.FilterContainsDetailedContext(context.Background(), id, options)
	if err != nil {
		return 0, err
	}
	return result.TotalRows, nil
}

func (manager *Manager) FilterContainsDetailedContext(ctx context.Context, id string, options FilterOptions) (FilterResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}
	if options.Query == "" {
		count, err := manager.ClearView(id)
		return FilterResult{TotalRows: count}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return FilterResult{}, errors.New("document not found")
	}
	query := options.Query
	if !options.CaseSensitive {
		query = strings.ToLower(query)
	}
	if canStreamContainsFilter(session) {
		return manager.filterContainsInStream(ctx, id, session, options, query)
	}
	records, err := recordsForSession(session)
	if err != nil {
		return FilterResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}

	viewRows := make([]int, 0)
	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	for row := start; row < len(records); row++ {
		if err := ctx.Err(); err != nil {
			return FilterResult{}, err
		}
		record := records[row]
		if filterRecordContains(record, options, query) {
			viewRows = append(viewRows, row)
		}
	}
	return manager.applyFilterViewRows(id, viewRows, false)
}

func canStreamContainsFilter(session *Session) bool {
	return canStreamCleanSession(session)
}

func canStreamCleanSession(session *Session) bool {
	return session.EditBaseRecords == nil &&
		len(session.Edits) == 0 &&
		csvcore.CanStreamRecords(session.Document.Metadata)
}

func (manager *Manager) filterContainsInStream(ctx context.Context, id string, session *Session, options FilterOptions, query string) (FilterResult, error) {
	viewRows := make([]int, 0)
	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	_, err := csvcore.ForEachRecordContext(ctx, session.WorkingPath, session.Document.Metadata, func(row int, record []string) error {
		if row < start {
			return nil
		}
		if filterRecordContains(record, options, query) {
			viewRows = append(viewRows, row)
		}
		return nil
	})
	if err != nil {
		return FilterResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}
	return manager.applyFilterViewRows(id, viewRows, true)
}

func filterRecordContains(record []string, options FilterOptions, query string) bool {
	for col, value := range record {
		if options.Column != nil && col != *options.Column {
			continue
		}
		haystack := value
		if !options.CaseSensitive {
			haystack = strings.ToLower(haystack)
		}
		if strings.Contains(haystack, query) {
			return true
		}
	}
	return false
}

func (manager *Manager) applyFilterViewRows(id string, viewRows []int, streamed bool) (FilterResult, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if !ok {
		return FilterResult{}, errors.New("document not found")
	}
	current.ViewRows = viewRows
	current.ViewActive = true
	return FilterResult{TotalRows: len(viewRows), Streamed: streamed}, nil
}

func naturalCompare(left string, right string) int {
	leftRunes := []rune(strings.ToLower(width.Narrow.String(left)))
	rightRunes := []rune(strings.ToLower(width.Narrow.String(right)))
	leftIndex := 0
	rightIndex := 0
	for leftIndex < len(leftRunes) && rightIndex < len(rightRunes) {
		leftRune := leftRunes[leftIndex]
		rightRune := rightRunes[rightIndex]
		if isASCIIDigit(leftRune) && isASCIIDigit(rightRune) {
			compare, nextLeft, nextRight := compareNaturalNumberRun(leftRunes, leftIndex, rightRunes, rightIndex)
			if compare != 0 {
				return compare
			}
			leftIndex = nextLeft
			rightIndex = nextRight
			continue
		}
		if leftRune < rightRune {
			return -1
		}
		if leftRune > rightRune {
			return 1
		}
		leftIndex++
		rightIndex++
	}
	if leftIndex < len(leftRunes) {
		return 1
	}
	if rightIndex < len(rightRunes) {
		return -1
	}
	return strings.Compare(left, right)
}

func compareNaturalNumberRun(left []rune, leftStart int, right []rune, rightStart int) (int, int, int) {
	leftEnd := scanASCIIDigits(left, leftStart)
	rightEnd := scanASCIIDigits(right, rightStart)
	leftSignificant := trimLeadingZeroRun(left, leftStart, leftEnd)
	rightSignificant := trimLeadingZeroRun(right, rightStart, rightEnd)
	leftLength := leftEnd - leftSignificant
	rightLength := rightEnd - rightSignificant
	if leftLength != rightLength {
		if leftLength < rightLength {
			return -1, leftEnd, rightEnd
		}
		return 1, leftEnd, rightEnd
	}
	for offset := 0; offset < leftLength; offset++ {
		leftDigit := left[leftSignificant+offset]
		rightDigit := right[rightSignificant+offset]
		if leftDigit < rightDigit {
			return -1, leftEnd, rightEnd
		}
		if leftDigit > rightDigit {
			return 1, leftEnd, rightEnd
		}
	}
	if leftEnd-leftStart != rightEnd-rightStart {
		if leftEnd-leftStart < rightEnd-rightStart {
			return -1, leftEnd, rightEnd
		}
		return 1, leftEnd, rightEnd
	}
	return 0, leftEnd, rightEnd
}

func scanASCIIDigits(value []rune, start int) int {
	index := start
	for index < len(value) && isASCIIDigit(value[index]) {
		index++
	}
	return index
}

func trimLeadingZeroRun(value []rune, start int, end int) int {
	index := start
	for index < end-1 && value[index] == '0' {
		index++
	}
	return index
}

func isASCIIDigit(value rune) bool {
	return value >= '0' && value <= '9'
}

func (manager *Manager) ColumnValues(id string, column int, limit int) ([]ColumnValue, int, error) {
	result, err := manager.ColumnValuesContext(context.Background(), id, column, limit)
	if err != nil {
		return nil, 0, err
	}
	return result.Values, result.TotalDistinct, nil
}

func (manager *Manager) ColumnValuesContext(ctx context.Context, id string, column int, limit int) (ColumnValuesResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if column < 0 {
		return ColumnValuesResult{}, errors.New("column must be non-negative")
	}
	if limit <= 0 {
		limit = 50
	}
	if err := ctx.Err(); err != nil {
		return ColumnValuesResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return ColumnValuesResult{}, errors.New("document not found")
	}
	if canStreamCleanSession(session) {
		return columnValuesInStream(ctx, session, column, limit)
	}
	records, err := recordsForSession(session)
	if err != nil {
		return ColumnValuesResult{}, err
	}
	counts, err := columnValueCountsFromRecords(ctx, records, session.Document.Metadata, column)
	if err != nil {
		return ColumnValuesResult{}, err
	}
	return columnValuesResultFromCounts(counts, limit, false), nil
}

func columnValuesInStream(ctx context.Context, session *Session, column int, limit int) (ColumnValuesResult, error) {
	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	counts := make(map[string]int)
	_, err := csvcore.ForEachRecordContext(ctx, session.WorkingPath, session.Document.Metadata, func(row int, record []string) error {
		if row < start {
			return nil
		}
		counts[cellValue(record, column)]++
		return nil
	})
	if err != nil {
		return ColumnValuesResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return ColumnValuesResult{}, err
	}
	return columnValuesResultFromCounts(counts, limit, true), nil
}

func columnValueCountsFromRecords(ctx context.Context, records [][]string, metadata csvcore.FileMetadata, column int) (map[string]int, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	start := 0
	if metadata.Dialect.HasHeader {
		start = 1
	}
	counts := make(map[string]int)
	for row := start; row < len(records); row++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		counts[cellValue(records[row], column)]++
	}
	return counts, nil
}

func columnValuesResultFromCounts(counts map[string]int, limit int, streamed bool) ColumnValuesResult {
	values := make([]ColumnValue, 0, len(counts))
	for value, count := range counts {
		values = append(values, ColumnValue{Value: value, Count: count, Blank: value == ""})
	}
	sort.SliceStable(values, func(i, j int) bool {
		if values[i].Blank != values[j].Blank {
			return !values[i].Blank
		}
		return naturalCompare(values[i].Value, values[j].Value) < 0
	})
	totalDistinct := len(values)
	if len(values) > limit {
		values = values[:limit]
	}
	return ColumnValuesResult{Values: values, TotalDistinct: totalDistinct, Streamed: streamed}
}

func (manager *Manager) FilterValues(id string, options ValueFilterOptions) (int, error) {
	result, err := manager.FilterValuesDetailedContext(context.Background(), id, options)
	if err != nil {
		return 0, err
	}
	return result.TotalRows, nil
}

func (manager *Manager) FilterValuesDetailedContext(ctx context.Context, id string, options ValueFilterOptions) (FilterResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if options.Column < 0 {
		return FilterResult{}, errors.New("column must be non-negative")
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return FilterResult{}, errors.New("document not found")
	}
	selected, includeBlank := selectedValueFilterSet(options)
	if canStreamCleanSession(session) {
		return manager.filterValuesInStream(ctx, id, session, options.Column, selected, includeBlank)
	}
	records, err := recordsForSession(session)
	if err != nil {
		return FilterResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}
	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	viewRows := make([]int, 0)
	for row := start; row < len(records); row++ {
		if err := ctx.Err(); err != nil {
			return FilterResult{}, err
		}
		value := cellValue(records[row], options.Column)
		if valueFilterMatches(value, selected, includeBlank) {
			viewRows = append(viewRows, row)
		}
	}
	return manager.applyFilterViewRows(id, viewRows, false)
}

func selectedValueFilterSet(options ValueFilterOptions) (map[string]bool, bool) {
	selected := make(map[string]bool, len(options.Values))
	includeBlank := options.IncludeBlank
	for _, value := range options.Values {
		if value == "" {
			includeBlank = true
			continue
		}
		selected[value] = true
	}
	return selected, includeBlank
}

func (manager *Manager) filterValuesInStream(ctx context.Context, id string, session *Session, column int, selected map[string]bool, includeBlank bool) (FilterResult, error) {
	viewRows := make([]int, 0)
	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	_, err := csvcore.ForEachRecordContext(ctx, session.WorkingPath, session.Document.Metadata, func(row int, record []string) error {
		if row < start {
			return nil
		}
		if valueFilterMatches(cellValue(record, column), selected, includeBlank) {
			viewRows = append(viewRows, row)
		}
		return nil
	})
	if err != nil {
		return FilterResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}
	return manager.applyFilterViewRows(id, viewRows, true)
}

func valueFilterMatches(value string, selected map[string]bool, includeBlank bool) bool {
	return (value == "" && includeBlank) || selected[value]
}

func (manager *Manager) FilterSQL(id string, expression string) (int, error) {
	result, err := manager.FilterSQLDetailedContext(context.Background(), id, expression)
	if err != nil {
		return 0, err
	}
	return result.TotalRows, nil
}

func (manager *Manager) FilterSQLDetailedContext(ctx context.Context, id string, expression string) (FilterResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}
	if strings.TrimSpace(expression) == "" {
		count, err := manager.ClearView(id)
		return FilterResult{TotalRows: count}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return FilterResult{}, errors.New("document not found")
	}
	if canStreamCleanSession(session) {
		return manager.filterSQLInStream(ctx, id, session, expression)
	}
	records, err := recordsForSession(session)
	if err != nil {
		return FilterResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}

	start := 0
	headers := []string(nil)
	if session.Document.Metadata.Dialect.HasHeader && len(records) > 0 {
		start = 1
		headers = records[0]
	}
	program, err := sqlexpr.Compile(stripWhere(expression), headers)
	if err != nil {
		return FilterResult{}, err
	}

	viewRows := make([]int, 0)
	for row := start; row < len(records); row++ {
		if err := ctx.Err(); err != nil {
			return FilterResult{}, err
		}
		matched, err := program.Eval(records[row])
		if err != nil {
			return FilterResult{}, fmt.Errorf("row %d: %w", row+1, err)
		}
		if matched {
			viewRows = append(viewRows, row)
		}
	}
	return manager.applyFilterViewRows(id, viewRows, false)
}

func (manager *Manager) filterSQLInStream(ctx context.Context, id string, session *Session, expression string) (FilterResult, error) {
	viewRows := make([]int, 0)
	source := stripWhere(expression)
	var program *sqlexpr.Program
	if !session.Document.Metadata.Dialect.HasHeader {
		compiled, err := sqlexpr.Compile(source, nil)
		if err != nil {
			return FilterResult{}, err
		}
		program = compiled
	}
	_, err := csvcore.ForEachRecordContext(ctx, session.WorkingPath, session.Document.Metadata, func(row int, record []string) error {
		if row == 0 && session.Document.Metadata.Dialect.HasHeader {
			compiled, err := sqlexpr.Compile(source, record)
			if err != nil {
				return err
			}
			program = compiled
			return nil
		}
		if program == nil {
			compiled, err := sqlexpr.Compile(source, nil)
			if err != nil {
				return err
			}
			program = compiled
		}
		matched, err := program.Eval(record)
		if err != nil {
			return fmt.Errorf("row %d: %w", row+1, err)
		}
		if matched {
			viewRows = append(viewRows, row)
		}
		return nil
	})
	if err != nil {
		return FilterResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return FilterResult{}, err
	}
	return manager.applyFilterViewRows(id, viewRows, true)
}

type SortOptions struct {
	Column     int       `json:"column"`
	Descending bool      `json:"descending"`
	Numeric    bool      `json:"numeric"`
	Mode       SortMode  `json:"mode,omitempty"`
	HasHeader  bool      `json:"hasHeader"`
	Keys       []SortKey `json:"keys"`
	Stable     *bool     `json:"stable,omitempty"`
}

type SortMode string

const (
	SortAuto   SortMode = "auto"
	SortText   SortMode = "text"
	SortNumber SortMode = "number"
	SortDate   SortMode = "date"
	SortLength SortMode = "length"
)

type SortKey struct {
	Column     int      `json:"column"`
	Descending bool     `json:"descending"`
	Numeric    bool     `json:"numeric"`
	Mode       SortMode `json:"mode,omitempty"`
}

func (manager *Manager) Sort(id string, options SortOptions) (int, error) {
	keys, err := normalizeSortKeys(options)
	if err != nil {
		return 0, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return 0, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return 0, err
	}

	start := 0
	if options.HasHeader || session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	var viewRows []int
	if session.ViewActive {
		viewRows = append([]int(nil), session.ViewRows...)
	} else {
		viewRows = make([]int, 0, max(0, len(records)-start))
		for row := start; row < len(records); row++ {
			viewRows = append(viewRows, row)
		}
	}
	less := func(i, j int) bool {
		return compareSortRecords(records[viewRows[i]], records[viewRows[j]], keys) < 0
	}
	if sortStable(options.Stable) {
		sort.SliceStable(viewRows, less)
	} else {
		sort.Slice(viewRows, less)
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if !ok {
		return 0, errors.New("document not found")
	}
	current.ViewRows = viewRows
	current.ViewActive = true
	return len(viewRows), nil
}

func sortStable(stable *bool) bool {
	return stable == nil || *stable
}

func normalizeSortKeys(options SortOptions) ([]SortKey, error) {
	keys := append([]SortKey(nil), options.Keys...)
	if len(keys) == 0 {
		keys = []SortKey{{
			Column:     options.Column,
			Descending: options.Descending,
			Numeric:    options.Numeric,
			Mode:       options.Mode,
		}}
	}
	for index, key := range keys {
		if key.Column < 0 {
			return nil, errors.New("column must be non-negative")
		}
		keys[index].Mode = normalizeSortMode(key.Mode, key.Numeric)
	}
	return keys, nil
}

func compareSortRecords(leftRecord []string, rightRecord []string, keys []SortKey) int {
	for _, key := range keys {
		left := cellValue(leftRecord, key.Column)
		right := cellValue(rightRecord, key.Column)
		result := compareSortValues(left, right, key.Mode)
		if result == 0 {
			continue
		}
		if key.Descending {
			return -result
		}
		return result
	}
	return 0
}

func normalizeSortMode(mode SortMode, numeric bool) SortMode {
	if mode == "" {
		if numeric {
			return SortNumber
		}
		return SortAuto
	}
	switch mode {
	case SortAuto, SortText, SortNumber, SortDate, SortLength:
		return mode
	default:
		return SortAuto
	}
}

func compareSortValues(left string, right string, mode SortMode) int {
	switch mode {
	case SortNumber:
		if result, ok := compareNumbers(left, right); ok {
			return result
		}
	case SortDate:
		if result, ok := compareDates(left, right); ok {
			return result
		}
	case SortLength:
		return compareInts(utf8.RuneCountInString(left), utf8.RuneCountInString(right))
	case SortAuto:
		if result, ok := compareNumbers(left, right); ok {
			return result
		}
		if result, ok := compareDates(left, right); ok {
			return result
		}
	}
	return strings.Compare(left, right)
}

func compareNumbers(left string, right string) (int, bool) {
	leftNum, leftErr := parseSortNumber(left)
	rightNum, rightErr := parseSortNumber(right)
	if leftErr == nil && rightErr == nil {
		return compareFloats(leftNum, rightNum), true
	}
	return 0, false
}

func parseSortNumber(value string) (float64, error) {
	normalized := strings.TrimSpace(normalizeFullWidthNumeric(value))
	normalized = strings.ReplaceAll(normalized, ",", "")
	if normalized == "" {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseFloat(normalized, 64)
}

func compareDates(left string, right string) (int, bool) {
	leftTime, leftErr := parseSortDate(left)
	rightTime, rightErr := parseSortDate(right)
	if leftErr == nil && rightErr == nil {
		return compareTimes(leftTime, rightTime), true
	}
	return 0, false
}

func parseSortDate(value string) (time.Time, error) {
	normalized := strings.TrimSpace(normalizeFullWidthNumeric(value))
	for _, layout := range []string{
		time.RFC3339,
		"2006-01-02",
		"2006/01/02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02 15:04:05",
		"01/02/2006",
		"01/02/2006 15:04",
		"01/02/2006 15:04:05",
	} {
		if parsed, err := time.Parse(layout, normalized); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("not a supported date")
}

func normalizeFullWidthNumeric(value string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= '０' && r <= '９':
			return '0' + (r - '０')
		case r == '．':
			return '.'
		case r == '－':
			return '-'
		case r == '＋':
			return '+'
		case r == '，':
			return ','
		default:
			return r
		}
	}, value)
}

func compareFloats(left float64, right float64) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func compareTimes(left time.Time, right time.Time) int {
	switch {
	case left.Before(right):
		return -1
	case left.After(right):
		return 1
	default:
		return 0
	}
}

func compareInts(left int, right int) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func textTransformFunc(mode TextTransformMode) (func(string) string, error) {
	switch mode {
	case TransformTrim:
		return strings.TrimSpace, nil
	case TransformUpper:
		return strings.ToUpper, nil
	case TransformLower:
		return strings.ToLower, nil
	case TransformSnake:
		return snakeCase, nil
	case TransformKebab:
		return kebabCase, nil
	case TransformCamel:
		return camelCase, nil
	case TransformPascal:
		return pascalCase, nil
	case TransformTitle:
		return titleCase, nil
	case TransformWide:
		return width.Widen.String, nil
	case TransformNarrow:
		return width.Narrow.String, nil
	default:
		return nil, fmt.Errorf("unsupported text transform %q", mode)
	}
}

func snakeCase(value string) string {
	return strings.Join(lowerNamingWords(value), "_")
}

func kebabCase(value string) string {
	return strings.Join(lowerNamingWords(value), "-")
}

func camelCase(value string) string {
	words := lowerNamingWords(value)
	if len(words) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString(words[0])
	for _, word := range words[1:] {
		builder.WriteString(upperFirst(word))
	}
	return builder.String()
}

func pascalCase(value string) string {
	words := lowerNamingWords(value)
	if len(words) == 0 {
		return ""
	}
	var builder strings.Builder
	for _, word := range words {
		builder.WriteString(upperFirst(word))
	}
	return builder.String()
}

func titleCase(value string) string {
	words := lowerNamingWords(value)
	if len(words) == 0 {
		return ""
	}
	for index, word := range words {
		words[index] = upperFirst(word)
	}
	return strings.Join(words, " ")
}

func lowerNamingWords(value string) []string {
	words := splitNamingWords(value)
	for index, word := range words {
		words[index] = strings.ToLower(word)
	}
	return words
}

func splitNamingWords(value string) []string {
	runes := []rune(value)
	words := make([]string, 0, 4)
	current := make([]rune, 0, len(runes))
	flush := func() {
		if len(current) == 0 {
			return
		}
		words = append(words, string(current))
		current = current[:0]
	}

	for index, currentRune := range runes {
		if !unicode.IsLetter(currentRune) && !unicode.IsDigit(currentRune) {
			flush()
			continue
		}
		if len(current) > 0 {
			var next rune
			hasNext := false
			if index+1 < len(runes) {
				next = runes[index+1]
				hasNext = true
			}
			if startsNamingWord(current, currentRune, next, hasNext) {
				flush()
			}
		}
		current = append(current, currentRune)
	}
	flush()
	return words
}

func startsNamingWord(current []rune, next rune, afterNext rune, hasAfterNext bool) bool {
	previous := current[len(current)-1]
	if unicode.IsDigit(previous) && unicode.IsLetter(next) {
		return true
	}
	if unicode.IsLetter(previous) && unicode.IsDigit(next) {
		return true
	}
	if (unicode.IsLower(previous) || unicode.IsDigit(previous)) && unicode.IsUpper(next) {
		return true
	}
	return len(current) > 1 &&
		unicode.IsUpper(previous) &&
		unicode.IsUpper(next) &&
		hasAfterNext &&
		unicode.IsLower(afterNext)
}

func upperFirst(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (manager *Manager) ClearView(id string) (int, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return 0, errors.New("document not found")
	}
	session.ViewRows = nil
	session.ViewActive = false
	return session.Document.RowIndex.Rows, nil
}

func (manager *Manager) SetReadOnly(id string, readOnly bool) (*Session, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return nil, errors.New("document not found")
	}
	session.ReadOnly = readOnly
	return cloneSession(session), nil
}

func (manager *Manager) FileStatus(id string) (FileStatus, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return FileStatus{}, errors.New("document not found")
	}
	status := FileStatus{
		Path:                 session.trackedPath(),
		Dirty:                session.Dirty,
		ReadOnly:             session.ReadOnly,
		SourceKind:           session.SourceKind,
		LastKnownSizeBytes:   session.LastKnownSizeBytes,
		LastKnownModTimeNano: session.LastKnownModTimeUnixNano,
	}
	if status.Path == "" {
		return status, nil
	}
	info, err := os.Stat(status.Path)
	if err != nil {
		if os.IsNotExist(err) {
			status.Missing = true
			status.ChangedExternally = true
			return status, nil
		}
		return FileStatus{}, err
	}
	status.SizeBytes = info.Size()
	status.ModTimeUnixNano = info.ModTime().UnixNano()
	status.ChangedExternally =
		status.SizeBytes != session.LastKnownSizeBytes ||
			status.ModTimeUnixNano != session.LastKnownModTimeUnixNano
	return status, nil
}

func (manager *Manager) Refresh(id string, discardDirty bool) (*Session, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return nil, errors.New("document not found")
	}
	if session.Dirty && !discardDirty {
		return nil, errors.New("document has unsaved edits")
	}
	path := session.trackedPath()
	if path == "" {
		return nil, errors.New("document has no file path")
	}
	delimiter := session.Document.Metadata.Dialect.Delimiter
	delimiterName := session.Document.Metadata.Dialect.DelimiterName
	hasHeader := session.Document.Metadata.Dialect.HasHeader
	headerColumns := session.Document.Metadata.Dialect.HeaderColumns
	refreshed, refreshedWorkingPath, err := openDocumentForPath(path, csvcore.OpenOptions{
		PreviewRows:   max(160, len(session.Document.Preview)),
		IndexRows:     -1,
		Delimiter:     &delimiter,
		DelimiterName: delimiterName,
		HasHeader:     &hasHeader,
		HeaderColumns: &headerColumns,
		Encoding:      &session.Document.Metadata.Encoding,
		SheetName:     selectedWorkbookSheetName(session.Document.Metadata),
		SheetIndex:    selectedWorkbookSheetIndex(session.Document.Metadata),
	})
	if err != nil {
		return nil, err
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if !ok {
		return nil, errors.New("document not found")
	}
	current.WorkingPath = refreshedWorkingPath
	current.Document = refreshed
	current.Edits = make(map[Cell]string)
	current.EditBaseRecords = nil
	current.ViewRows = nil
	current.ViewActive = false
	current.UndoStack = nil
	current.RedoStack = nil
	current.Dirty = false
	current.recordKnownFileState()
	return cloneSession(current), nil
}

func (manager *Manager) mutateRecords(id string, mutate func(records [][]string) ([][]string, error)) (MutationResult, error) {
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := recordsForSession(session)
	if err != nil {
		return MutationResult{}, err
	}
	records = cloneRecords(records)

	if err := manager.pushUndo(id); err != nil {
		return MutationResult{}, err
	}
	records, err = mutate(records)
	if err != nil {
		return MutationResult{}, err
	}
	updated, err := manager.materializeRecords(id, records, session.Document.Metadata, true)
	if err != nil {
		return MutationResult{}, err
	}
	return MutationResult{Session: updated, TotalRows: updated.Document.RowIndex.Rows, Dirty: updated.Dirty}, nil
}

func (manager *Manager) pushUndo(id string) error {
	session, ok := manager.snapshot(id)
	if !ok {
		return errors.New("document not found")
	}
	snapshot, err := snapshotFromSession(session)
	if err != nil {
		return err
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if !ok {
		return errors.New("document not found")
	}
	current.UndoStack = append(current.UndoStack, snapshot)
	if len(current.UndoStack) > 100 {
		current.UndoStack = current.UndoStack[len(current.UndoStack)-100:]
	}
	current.RedoStack = nil
	return nil
}

func snapshotFromSession(session *Session) (historySnapshot, error) {
	records, err := recordsForSession(session)
	if err != nil {
		return historySnapshot{}, err
	}
	return historySnapshot{
		Records:    cloneRecords(records),
		Metadata:   session.Document.Metadata,
		Dirty:      session.Dirty,
		ViewRows:   append([]int(nil), session.ViewRows...),
		ViewActive: session.ViewActive,
	}, nil
}

func (manager *Manager) applySnapshot(id string, snapshot historySnapshot, adjustHistory func(session *Session)) (*Session, error) {
	if _, err := manager.materializeRecords(id, cloneRecords(snapshot.Records), snapshot.Metadata, snapshot.Dirty); err != nil {
		return nil, err
	}
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return nil, errors.New("document not found")
	}
	session.ViewRows = append([]int(nil), snapshot.ViewRows...)
	session.ViewActive = snapshot.ViewActive
	if adjustHistory != nil {
		adjustHistory(session)
	}
	return cloneSession(session), nil
}

func (manager *Manager) materializeRecords(id string, records [][]string, metadata csvcore.FileMetadata, dirty bool) (*Session, error) {
	temp, err := os.CreateTemp("", "opencsv-work-*.csv")
	if err != nil {
		return nil, err
	}
	tempPath := temp.Name()
	if err := temp.Close(); err != nil {
		_ = os.Remove(tempPath)
		return nil, err
	}
	if err := os.Remove(tempPath); err != nil {
		return nil, err
	}
	if _, err := csvcore.SaveRecords(tempPath, metadata, records, csvcore.SaveOptions{}); err != nil {
		return nil, err
	}
	delimiter := metadata.Dialect.Delimiter
	delimiterName := metadata.Dialect.DelimiterName
	hasHeader := metadata.Dialect.HasHeader
	headerColumns := metadata.Dialect.HeaderColumns
	refreshed, err := csvcore.OpenDocument(tempPath, csvcore.OpenOptions{
		PreviewRows:   max(160, min(len(records), 1000)),
		IndexRows:     -1,
		Delimiter:     &delimiter,
		DelimiterName: delimiterName,
		HasHeader:     &hasHeader,
		HeaderColumns: &headerColumns,
		Encoding:      &metadata.Encoding,
		QuoteChar:     &metadata.Dialect.QuoteChar,
		QuoteMode:     metadata.Dialect.QuoteMode,
		QuoteEscape:   metadata.Dialect.QuoteEscape,
	})
	if err != nil {
		return nil, err
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return nil, errors.New("document not found")
	}
	session.WorkingPath = tempPath
	session.Document = refreshed
	session.Edits = make(map[Cell]string)
	session.EditBaseRecords = nil
	session.ViewRows = nil
	session.ViewActive = false
	session.Dirty = dirty
	return cloneSession(session), nil
}

func (manager *Manager) snapshot(id string) (*Session, bool) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return nil, false
	}
	return cloneSession(session), true
}

func (manager *Manager) Snapshot(id string) (*Session, bool) {
	return manager.snapshot(id)
}

func (manager *Manager) Close(id string, force bool) (*Session, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return nil, errors.New("document not found")
	}
	if session.Dirty && !force {
		return nil, errors.New("document has unsaved changes")
	}
	delete(manager.sessions, id)
	return cloneSession(session), nil
}

func (manager *Manager) ensureWritable(id string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return errors.New("document not found")
	}
	if session.ReadOnly {
		return errors.New("document is read-only")
	}
	return nil
}

func (session *Session) trackedPath() string {
	if session.SourceKind == SourceFile && session.OriginalPath != "" {
		return session.OriginalPath
	}
	return session.WorkingPath
}

func (session *Session) recordKnownFileState() {
	path := session.trackedPath()
	session.LastKnownPath = path
	session.LastKnownSizeBytes = 0
	session.LastKnownModTimeUnixNano = 0
	if path == "" {
		return
	}
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	session.LastKnownSizeBytes = info.Size()
	session.LastKnownModTimeUnixNano = info.ModTime().UnixNano()
}

func cloneSession(session *Session) *Session {
	copy := *session
	copy.Edits = make(map[Cell]string, len(session.Edits))
	for key, value := range session.Edits {
		copy.Edits[key] = value
	}
	copy.EditBaseRecords = session.EditBaseRecords
	copy.ViewRows = append([]int(nil), session.ViewRows...)
	copy.UndoStack = cloneHistory(session.UndoStack)
	copy.RedoStack = cloneHistory(session.RedoStack)
	return &copy
}

func recordsForSession(session *Session) ([][]string, error) {
	if session.EditBaseRecords != nil {
		records := cloneRecords(session.EditBaseRecords)
		applyEdits(records, 0, session.Edits)
		return records, nil
	}
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return nil, err
	}
	applyEdits(records, 0, session.Edits)
	return records, nil
}

func cloneHistory(input []historySnapshot) []historySnapshot {
	if len(input) == 0 {
		return nil
	}
	output := make([]historySnapshot, len(input))
	for i, snapshot := range input {
		output[i] = historySnapshot{
			Records:    cloneRecords(snapshot.Records),
			Metadata:   snapshot.Metadata,
			Dirty:      snapshot.Dirty,
			ViewRows:   append([]int(nil), snapshot.ViewRows...),
			ViewActive: snapshot.ViewActive,
		}
	}
	return output
}

func cloneRecords(records [][]string) [][]string {
	if records == nil {
		return nil
	}
	output := make([][]string, len(records))
	for i, record := range records {
		output[i] = append([]string(nil), record...)
	}
	return output
}

func transposeRecords(records [][]string, width int) [][]string {
	if len(records) == 0 || width == 0 {
		return nil
	}
	output := make([][]string, width)
	for col := 0; col < width; col++ {
		output[col] = make([]string, len(records))
		for row := range records {
			if col < len(records[row]) {
				output[col][row] = records[row][col]
			}
		}
	}
	return output
}

func moveRecord(records [][]string, from int, to int) [][]string {
	record := records[from]
	if from < to {
		copy(records[from:to], records[from+1:to+1])
	} else {
		copy(records[to+1:from+1], records[to:from])
	}
	records[to] = record
	return records
}

func insertRecords(records [][]string, index int, inserted [][]string) [][]string {
	if len(inserted) == 0 {
		return records
	}
	if index < 0 {
		index = 0
	}
	if index > len(records) {
		index = len(records)
	}
	output := make([][]string, 0, len(records)+len(inserted))
	output = append(output, records[:index]...)
	output = append(output, inserted...)
	output = append(output, records[index:]...)
	return output
}

func moveColumnRecords(records [][]string, from int, to int, width int) [][]string {
	order := make([]int, width)
	for i := range order {
		order[i] = i
	}
	moved := order[from]
	if from < to {
		copy(order[from:to], order[from+1:to+1])
	} else {
		copy(order[to+1:from+1], order[to:from])
	}
	order[to] = moved

	output := make([][]string, len(records))
	for row, record := range records {
		padded := append([]string(nil), record...)
		ensureRecordWidth(&padded, width)
		next := make([]string, width)
		for col, sourceCol := range order {
			next[col] = padded[sourceCol]
		}
		output[row] = next
	}
	return output
}

func insertRecordValues(record []string, index int, inserted []string) []string {
	if len(inserted) == 0 {
		return record
	}
	if index < 0 {
		index = 0
	}
	if index > len(record) {
		index = len(record)
	}
	output := make([]string, 0, len(record)+len(inserted))
	output = append(output, record[:index]...)
	output = append(output, inserted...)
	output = append(output, record[index:]...)
	return output
}

func uniqueIndexesInOrder(indexes []int) []int {
	seen := make(map[int]bool, len(indexes))
	output := make([]int, 0, len(indexes))
	for _, index := range indexes {
		if seen[index] {
			continue
		}
		seen[index] = true
		output = append(output, index)
	}
	return output
}

func indexSet(indexes []int) map[int]bool {
	set := make(map[int]bool, len(indexes))
	for _, index := range indexes {
		set[index] = true
	}
	return set
}

func countIndexesBefore(indexes []int, target int) int {
	count := 0
	for _, index := range indexes {
		if index < target {
			count++
		}
	}
	return count
}

func visibleRowCount(session *Session, recordCount int, dataStart int) int {
	if session.ViewActive {
		return len(session.ViewRows)
	}
	return max(0, recordCount-dataStart)
}

func joinRecords(records [][]string, cells []Cell, separator string, direction string) int {
	if direction == "vertical" {
		return joinRecordsVertical(records, cells, separator)
	}
	return joinRecordsHorizontal(records, cells, separator)
}

func joinRecordsHorizontal(records [][]string, cells []Cell, separator string) int {
	grouped := make(map[int][]int)
	seen := make(map[Cell]bool, len(cells))
	for _, cell := range cells {
		if cell.Row < 0 || cell.Row >= len(records) || seen[cell] {
			continue
		}
		seen[cell] = true
		grouped[cell.Row] = append(grouped[cell.Row], cell.Col)
	}
	changed := 0
	for row, cols := range grouped {
		sort.Ints(cols)
		if len(cols) < 2 {
			continue
		}
		ensureRecordWidth(&records[row], cols[len(cols)-1]+1)
		values := make([]string, len(cols))
		for index, col := range cols {
			values[index] = records[row][col]
		}
		joined := strings.Join(values, separator)
		if records[row][cols[0]] != joined {
			records[row][cols[0]] = joined
			changed++
		}
		for _, col := range cols[1:] {
			if records[row][col] != "" {
				records[row][col] = ""
				changed++
			}
		}
	}
	return changed
}

func joinRecordsVertical(records [][]string, cells []Cell, separator string) int {
	grouped := make(map[int][]int)
	seen := make(map[Cell]bool, len(cells))
	for _, cell := range cells {
		if cell.Row < 0 || cell.Row >= len(records) || seen[cell] {
			continue
		}
		seen[cell] = true
		grouped[cell.Col] = append(grouped[cell.Col], cell.Row)
	}
	changed := 0
	for col, rows := range grouped {
		sort.Ints(rows)
		if len(rows) < 2 {
			continue
		}
		values := make([]string, len(rows))
		for index, row := range rows {
			values[index] = cellValue(records[row], col)
		}
		joined := strings.Join(values, separator)
		ensureRecordWidth(&records[rows[0]], col+1)
		if records[rows[0]][col] != joined {
			records[rows[0]][col] = joined
			changed++
		}
		for _, row := range rows[1:] {
			ensureRecordWidth(&records[row], col+1)
			if records[row][col] != "" {
				records[row][col] = ""
				changed++
			}
		}
	}
	return changed
}

type splitOperation struct {
	row   int
	col   int
	parts []string
}

func splitRecords(records *[][]string, cells []Cell, separator string, direction string) int {
	if direction == "vertical" {
		return splitRecordsVertical(records, cells, separator)
	}
	return splitRecordsHorizontal(*records, cells, separator)
}

func splitRecordsHorizontal(records [][]string, cells []Cell, separator string) int {
	seen := make(map[Cell]bool, len(cells))
	operations := make([]splitOperation, 0, len(cells))
	for _, cell := range cells {
		if cell.Row < 0 || cell.Row >= len(records) || seen[cell] {
			continue
		}
		seen[cell] = true
		value := cellValue(records[cell.Row], cell.Col)
		parts := strings.Split(value, separator)
		if len(parts) < 2 {
			continue
		}
		operations = append(operations, splitOperation{row: cell.Row, col: cell.Col, parts: parts})
	}
	sort.Slice(operations, func(i int, j int) bool {
		if operations[i].row == operations[j].row {
			return operations[i].col > operations[j].col
		}
		return operations[i].row < operations[j].row
	})
	changed := 0
	for _, operation := range operations {
		ensureRecordWidth(&records[operation.row], operation.col+len(operation.parts))
		for offset, part := range operation.parts {
			targetCol := operation.col + offset
			if records[operation.row][targetCol] != part {
				records[operation.row][targetCol] = part
				changed++
			}
		}
	}
	return changed
}

func splitRecordsVertical(records *[][]string, cells []Cell, separator string) int {
	seen := make(map[Cell]bool, len(cells))
	operations := make([]splitOperation, 0, len(cells))
	for _, cell := range cells {
		if cell.Row < 0 || seen[cell] {
			continue
		}
		seen[cell] = true
		value := ""
		if cell.Row < len(*records) {
			value = cellValue((*records)[cell.Row], cell.Col)
		}
		parts := strings.Split(value, separator)
		if len(parts) < 2 {
			continue
		}
		operations = append(operations, splitOperation{row: cell.Row, col: cell.Col, parts: parts})
	}
	sort.Slice(operations, func(i int, j int) bool {
		if operations[i].col == operations[j].col {
			return operations[i].row < operations[j].row
		}
		return operations[i].col < operations[j].col
	})
	changed := 0
	for _, operation := range operations {
		lastRow := operation.row + len(operation.parts) - 1
		for len(*records) <= lastRow {
			*records = append(*records, nil)
		}
		for offset, part := range operation.parts {
			targetRow := operation.row + offset
			ensureRecordWidth(&(*records)[targetRow], operation.col+1)
			if (*records)[targetRow][operation.col] != part {
				(*records)[targetRow][operation.col] = part
				changed++
			}
		}
	}
	return changed
}

func deduplicateRecords(records [][]string, hasHeader bool) ([][]string, int) {
	if len(records) == 0 {
		return records, 0
	}
	start := 0
	deduped := make([][]string, 0, len(records))
	if hasHeader {
		start = 1
		deduped = append(deduped, append([]string(nil), records[0]...))
	}
	seen := make(map[string]bool, max(0, len(records)-start))
	removed := 0
	for _, record := range records[start:] {
		key := recordKey(record)
		if seen[key] {
			removed++
			continue
		}
		seen[key] = true
		deduped = append(deduped, record)
	}
	return deduped, removed
}

func deduplicateColumns(records [][]string) ([][]string, int) {
	width := maxRecordWidth(records)
	if len(records) == 0 || width == 0 {
		return records, 0
	}
	seen := make(map[string]bool, width)
	keep := make([]bool, width)
	removed := 0
	for col := 0; col < width; col++ {
		key := columnKey(records, col)
		if seen[key] {
			removed++
			continue
		}
		seen[key] = true
		keep[col] = true
	}
	if removed == 0 {
		return records, 0
	}
	deduped := make([][]string, len(records))
	for row, record := range records {
		next := make([]string, 0, len(record))
		for col, value := range record {
			if col < len(keep) && keep[col] {
				next = append(next, value)
			}
		}
		deduped[row] = next
	}
	return deduped, removed
}

func columnKey(records [][]string, col int) string {
	var builder strings.Builder
	builder.WriteString(strconv.Itoa(len(records)))
	builder.WriteByte('|')
	for _, record := range records {
		value := ""
		if col < len(record) {
			value = record[col]
		}
		builder.WriteString(strconv.Itoa(len(value)))
		builder.WriteByte(':')
		builder.WriteString(value)
		builder.WriteByte(0)
	}
	return builder.String()
}

func recordKey(record []string) string {
	var builder strings.Builder
	for _, value := range record {
		builder.WriteString(strconv.Itoa(len(value)))
		builder.WriteByte(':')
		builder.WriteString(value)
		builder.WriteByte(0)
	}
	return builder.String()
}

func maxRecordWidth(records [][]string) int {
	width := 0
	for _, record := range records {
		width = max(width, len(record))
	}
	return width
}

func ensureRecordWidth(record *[]string, width int) {
	for len(*record) < width {
		*record = append(*record, "")
	}
}

func applyEdits(records [][]string, start int, edits map[Cell]string) {
	for key, value := range edits {
		localRow := key.Row - start
		if localRow < 0 || localRow >= len(records) {
			continue
		}
		for len(records[localRow]) <= key.Col {
			records[localRow] = append(records[localRow], "")
		}
		records[localRow][key.Col] = value
	}
}

func cellValue(record []string, col int) string {
	if col < 0 || col >= len(record) {
		return ""
	}
	return record[col]
}

func projectRecord(record []string, columns []int) []string {
	if len(columns) == 0 {
		return append([]string(nil), record...)
	}
	projected := make([]string, len(columns))
	for index, column := range columns {
		projected[index] = cellValue(record, column)
	}
	return projected
}

func exportColumnCount(records [][]string, columns []int) int {
	if len(columns) > 0 {
		return len(columns)
	}
	return maxRecordWidth(records)
}

func normalizeExportViewFormat(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "", "csv":
		return "csv"
	case "md":
		return "markdown"
	case "tex":
		return "latex"
	default:
		return format
	}
}

func isRichExportViewFormat(format string) bool {
	switch format {
	case "markdown", "json", "html", "sql", "latex", "raw":
		return true
	default:
		return false
	}
}

func encodeRichExportView(records [][]string, format string) ([]byte, error) {
	switch format {
	case "markdown":
		return []byte(clipboard.ToMarkdown(records)), nil
	case "json":
		text, err := clipboard.ToJSON(records)
		if err != nil {
			return nil, err
		}
		return []byte(text), nil
	case "html":
		return []byte(clipboard.ToHTML(records)), nil
	case "sql":
		return []byte(clipboard.ToSQLValues(records)), nil
	case "latex":
		return []byte(clipboard.ToLaTeX(records)), nil
	case "raw":
		return []byte(clipboard.ToRaw(records)), nil
	default:
		return nil, fmt.Errorf("unsupported export format %q", format)
	}
}

func exportViewRichExtension(format string) string {
	switch format {
	case "markdown":
		return ".md"
	case "json":
		return ".json"
	case "html":
		return ".html"
	case "sql":
		return ".sql"
	case "latex":
		return ".tex"
	default:
		return ".txt"
	}
}

func exportViewRichContentType(format string) string {
	switch format {
	case "markdown":
		return "text/markdown; charset=utf-8"
	case "json":
		return "application/json; charset=utf-8"
	case "html":
		return "text/html; charset=utf-8"
	case "sql":
		return "application/sql; charset=utf-8"
	case "latex":
		return "application/x-latex; charset=utf-8"
	default:
		return "text/plain; charset=utf-8"
	}
}

func exportViewFilename(displayName string, delimiter rune) string {
	targetExt := ".csv"
	if delimiter == '\t' {
		targetExt = ".tsv"
	}
	return exportViewFilenameWithExtension(displayName, targetExt)
}

func exportViewXLSXFilename(displayName string) string {
	return exportViewFilenameWithExtension(displayName, ".xlsx")
}

func exportViewFilenameWithExtension(displayName string, targetExt string) string {
	base := filepath.Base(displayName)
	if base == "." || base == string(filepath.Separator) || base == "" {
		base = "opencsv"
	}
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	if stem == "" {
		stem = "opencsv"
	}
	return stem + ".view" + targetExt
}

func exportViewContentType(delimiter rune) string {
	if delimiter == '\t' {
		return "text/tab-separated-values; charset=utf-8"
	}
	return "text/csv; charset=utf-8"
}

func isXLSXPath(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".xlsx")
}

type matcher struct {
	options FindOptions
	query   string
	regex   *regexp.Regexp
}

func newMatcher(options FindOptions) (matcher, error) {
	result := matcher{options: options, query: options.Query}
	if !options.CaseSensitive {
		result.query = strings.ToLower(result.query)
	}
	if options.Regex {
		pattern := options.Query
		if !options.CaseSensitive {
			pattern = "(?i:" + pattern + ")"
		}
		if options.WholeCell {
			pattern = "^(?:" + pattern + ")$"
		}
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return matcher{}, err
		}
		result.regex = compiled
	}
	return result, nil
}

func (matcher matcher) matches(value string) bool {
	if matcher.regex != nil {
		return matcher.regex.MatchString(value)
	}
	haystack := value
	if !matcher.options.CaseSensitive {
		haystack = strings.ToLower(haystack)
	}
	if matcher.options.WholeCell {
		return haystack == matcher.query
	}
	return strings.Contains(haystack, matcher.query)
}

func (matcher matcher) replace(value string, replacement string, all bool) string {
	if matcher.regex != nil {
		if !all {
			indexes := matcher.regex.FindStringSubmatchIndex(value)
			if indexes == nil {
				return value
			}
			return string(matcher.regex.ExpandString([]byte{}, replacement, value, indexes))
		}
		return matcher.regex.ReplaceAllString(value, replacement)
	}
	if matcher.options.WholeCell {
		return replacement
	}
	if matcher.options.CaseSensitive {
		limit := 1
		if all {
			limit = -1
		}
		return strings.Replace(value, matcher.options.Query, replacement, limit)
	}
	return replaceInsensitive(value, matcher.options.Query, replacement, all)
}

func replaceInsensitive(value string, query string, replacement string, all bool) string {
	lowerValue := strings.ToLower(value)
	lowerQuery := strings.ToLower(query)
	var builder strings.Builder
	start := 0
	for {
		index := strings.Index(lowerValue[start:], lowerQuery)
		if index < 0 {
			builder.WriteString(value[start:])
			break
		}
		index += start
		builder.WriteString(value[start:index])
		builder.WriteString(replacement)
		start = index + len(query)
		if !all {
			builder.WriteString(value[start:])
			break
		}
	}
	result := builder.String()
	if result == "" && value != "" {
		return value
	}
	return result
}

func stripWhere(expression string) string {
	trimmed := strings.TrimSpace(expression)
	if len(trimmed) >= len("where ") && strings.EqualFold(trimmed[:len("where ")], "where ") {
		return strings.TrimSpace(trimmed[len("where "):])
	}
	return trimmed
}
