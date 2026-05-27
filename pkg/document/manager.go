package document

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/opencsv/opencsv/pkg/csvcore"
	"github.com/opencsv/opencsv/pkg/sqlexpr"
)

type SourceKind string

const (
	SourceFile   SourceKind = "file"
	SourceUpload SourceKind = "upload"
)

type Cell struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type Session struct {
	ID                       string            `json:"id"`
	SourceKind               SourceKind        `json:"sourceKind"`
	DisplayName              string            `json:"displayName"`
	OriginalPath             string            `json:"originalPath"`
	WorkingPath              string            `json:"workingPath"`
	Document                 csvcore.Document  `json:"document"`
	Edits                    map[Cell]string   `json:"-"`
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
	document, err := csvcore.OpenDocument(path, options)
	if err != nil {
		return nil, err
	}
	if displayName == "" {
		displayName = path
	}
	return manager.Add(document, kind, displayName, path, path), nil
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

	if !session.ViewActive {
		window, err := csvcore.ReadRecordWindow(session.WorkingPath, session.Document.Metadata, start, count)
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

	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return csvcore.RecordWindow{}, err
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return nil, nil, err
	}
	applyEdits(records, 0, session.Edits)
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
	viewRows := append([]int(nil), session.ViewRows...)
	if !session.ViewActive {
		viewRows = make([]int, 0, max(0, len(records)-start))
		for row := start; row < len(records); row++ {
			viewRows = append(viewRows, row)
		}
	}

	exportRecords := make([][]string, 0, len(viewRows)+1)
	if hasHeader {
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

	data, err := csvcore.EncodeRecords(exportRecords, session.Document.Metadata)
	if err != nil {
		return ExportViewResult{}, err
	}
	return ExportViewResult{
		Filename:  exportViewFilename(session.DisplayName, session.Document.Metadata.Dialect.Delimiter),
		Data:      data,
		Rows:      exportedRows,
		Columns:   exportColumnCount(exportRecords, options.Columns),
		Delimiter: session.Document.Metadata.Dialect.Delimiter,
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
	for _, edit := range edits {
		session.Edits[Cell{Row: edit.Row, Col: edit.Col}] = edit.Value
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

type ExportViewOptions struct {
	Columns []int `json:"columns"`
}

type ExportViewResult struct {
	Filename  string
	Data      []byte
	Rows      int
	Columns   int
	Delimiter rune
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)
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

func (manager *Manager) HideRow(id string, row int) (MutationResult, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)
	width := maxRecordWidth(records)
	cellCount := len(records) * width
	if cellCount > maxCells {
		return MutationResult{}, fmt.Errorf("transpose would touch %d cells; limit is %d", cellCount, maxCells)
	}
	transposed := transposeRecords(records, width)
	metadata := session.Document.Metadata
	metadata.Dialect.HasHeader = false
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
	TransformTrim  TextTransformMode = "trim"
	TransformUpper TextTransformMode = "upper"
	TransformLower TextTransformMode = "lower"
)

func (manager *Manager) DeduplicateRows(id string) (MutationResult, error) {
	if err := manager.ensureWritable(id); err != nil {
		return MutationResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return MutationResult{}, errors.New("document not found")
	}
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)

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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)

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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)

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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)

	changed := joinRecords(records, cells, separator)
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)

	changed := splitRecords(records, cells, separator)
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return csvcore.SaveResult{}, err
	}
	applyEdits(records, 0, session.Edits)

	outputPath := path
	if outputPath == "" {
		outputPath = session.OriginalPath
	}
	if outputPath == "" {
		outputPath = session.WorkingPath
	}
	if outputPath != session.OriginalPath {
		createBackup = false
	}
	result, err := csvcore.SaveRecords(outputPath, session.Document.Metadata, records, csvcore.SaveOptions{CreateBackup: createBackup})
	if err != nil {
		return csvcore.SaveResult{}, err
	}

	refreshed, openErr := csvcore.OpenDocument(outputPath, csvcore.OpenOptions{
		PreviewRows: max(1, len(session.Document.Preview)),
		IndexRows:   max(1, len(session.Document.Preview)),
	})
	manager.mu.Lock()
	defer manager.mu.Unlock()
	current, ok := manager.sessions[id]
	if ok {
		current.Edits = make(map[Cell]string)
		current.Dirty = false
		current.ViewRows = nil
		current.ViewActive = false
		if openErr == nil {
			current.Document = refreshed
			current.WorkingPath = outputPath
			current.OriginalPath = outputPath
		}
		current.recordKnownFileState()
	}
	return result, nil
}

type FindOptions struct {
	Query         string `json:"query"`
	Column        *int   `json:"column,omitempty"`
	Cells         []Cell `json:"cells,omitempty"`
	CaseSensitive bool   `json:"caseSensitive"`
	WholeCell     bool   `json:"wholeCell"`
	Regex         bool   `json:"regex"`
	Limit         int    `json:"limit"`
}

type Match struct {
	Row   int    `json:"row"`
	Col   int    `json:"col"`
	Value string `json:"value"`
}

type FindResult struct {
	Matches   []Match `json:"matches"`
	Total     int     `json:"total"`
	Truncated bool    `json:"truncated"`
}

func (manager *Manager) Find(id string, options FindOptions) ([]Match, error) {
	result, err := manager.FindDetailed(id, options)
	if err != nil {
		return nil, err
	}
	return result.Matches, nil
}

func (manager *Manager) FindDetailed(id string, options FindOptions) (FindResult, error) {
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
	matcher, err := newMatcher(options)
	if err != nil {
		return FindResult{}, err
	}

	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return FindResult{}, err
	}
	applyEdits(records, 0, session.Edits)

	var matches []Match
	total := 0
	scope := cellScope(options.Cells)
	for row, record := range records {
		for col, value := range record {
			if options.Column != nil && col != *options.Column {
				continue
			}
			if !cellInScope(scope, row, col) {
				continue
			}
			if matcher.matches(value) {
				total++
				if len(matches) < options.Limit {
					matches = append(matches, Match{Row: row, Col: col, Value: value})
				}
			}
		}
	}
	return FindResult{Matches: matches, Total: total, Truncated: total > len(matches)}, nil
}

type ReplaceOptions struct {
	FindOptions
	Replacement string `json:"replacement"`
	All         bool   `json:"all"`
}

type ReplaceResult struct {
	Changes []Match `json:"changes"`
	Count   int     `json:"count"`
	Dirty   bool    `json:"dirty"`
}

func (manager *Manager) Replace(id string, options ReplaceOptions) (ReplaceResult, error) {
	if options.Query == "" {
		return ReplaceResult{}, nil
	}
	if err := manager.ensureWritable(id); err != nil {
		return ReplaceResult{}, err
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return ReplaceResult{}, errors.New("document not found")
	}
	matcher, err := newMatcher(options.FindOptions)
	if err != nil {
		return ReplaceResult{}, err
	}

	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return ReplaceResult{}, err
	}
	applyEdits(records, 0, session.Edits)

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
				return manager.applyReplaceEdits(id, edits, changes)
			}
		}
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

type ColumnValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
	Blank bool   `json:"blank"`
}

type ValueFilterOptions struct {
	Column       int      `json:"column"`
	Values       []string `json:"values"`
	IncludeBlank bool     `json:"includeBlank"`
}

func (manager *Manager) FilterContains(id string, options FilterOptions) (int, error) {
	if options.Query == "" {
		return manager.ClearView(id)
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return 0, errors.New("document not found")
	}
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return 0, err
	}
	applyEdits(records, 0, session.Edits)

	query := options.Query
	if !options.CaseSensitive {
		query = strings.ToLower(query)
	}
	viewRows := make([]int, 0)
	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	for row := start; row < len(records); row++ {
		record := records[row]
		for col, value := range record {
			if options.Column != nil && col != *options.Column {
				continue
			}
			haystack := value
			if !options.CaseSensitive {
				haystack = strings.ToLower(haystack)
			}
			if strings.Contains(haystack, query) {
				viewRows = append(viewRows, row)
				break
			}
		}
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

func (manager *Manager) ColumnValues(id string, column int, limit int) ([]ColumnValue, int, error) {
	if column < 0 {
		return nil, 0, errors.New("column must be non-negative")
	}
	if limit <= 0 {
		limit = 50
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return nil, 0, errors.New("document not found")
	}
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return nil, 0, err
	}
	applyEdits(records, 0, session.Edits)

	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	counts := make(map[string]int)
	for row := start; row < len(records); row++ {
		counts[cellValue(records[row], column)]++
	}
	values := make([]ColumnValue, 0, len(counts))
	for value, count := range counts {
		values = append(values, ColumnValue{Value: value, Count: count, Blank: value == ""})
	}
	sort.SliceStable(values, func(i, j int) bool {
		if values[i].Count != values[j].Count {
			return values[i].Count > values[j].Count
		}
		if values[i].Blank != values[j].Blank {
			return !values[i].Blank
		}
		return values[i].Value < values[j].Value
	})
	totalDistinct := len(values)
	if len(values) > limit {
		values = values[:limit]
	}
	return values, totalDistinct, nil
}

func (manager *Manager) FilterValues(id string, options ValueFilterOptions) (int, error) {
	if options.Column < 0 {
		return 0, errors.New("column must be non-negative")
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return 0, errors.New("document not found")
	}
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return 0, err
	}
	applyEdits(records, 0, session.Edits)

	selected := make(map[string]bool, len(options.Values))
	for _, value := range options.Values {
		if value == "" {
			options.IncludeBlank = true
			continue
		}
		selected[value] = true
	}
	start := 0
	if session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	viewRows := make([]int, 0)
	for row := start; row < len(records); row++ {
		value := cellValue(records[row], options.Column)
		if (value == "" && options.IncludeBlank) || selected[value] {
			viewRows = append(viewRows, row)
		}
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

func (manager *Manager) FilterSQL(id string, expression string) (int, error) {
	if strings.TrimSpace(expression) == "" {
		return manager.ClearView(id)
	}
	session, ok := manager.snapshot(id)
	if !ok {
		return 0, errors.New("document not found")
	}
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return 0, err
	}
	applyEdits(records, 0, session.Edits)

	start := 0
	headers := []string(nil)
	if session.Document.Metadata.Dialect.HasHeader && len(records) > 0 {
		start = 1
		headers = records[0]
	}
	program, err := sqlexpr.Compile(stripWhere(expression), headers)
	if err != nil {
		return 0, err
	}

	viewRows := make([]int, 0)
	for row := start; row < len(records); row++ {
		matched, err := program.Eval(records[row])
		if err != nil {
			return 0, fmt.Errorf("row %d: %w", row+1, err)
		}
		if matched {
			viewRows = append(viewRows, row)
		}
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

type SortOptions struct {
	Column     int       `json:"column"`
	Descending bool      `json:"descending"`
	Numeric    bool      `json:"numeric"`
	Mode       SortMode  `json:"mode,omitempty"`
	HasHeader  bool      `json:"hasHeader"`
	Keys       []SortKey `json:"keys"`
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return 0, err
	}
	applyEdits(records, 0, session.Edits)

	start := 0
	if options.HasHeader || session.Document.Metadata.Dialect.HasHeader {
		start = 1
	}
	viewRows := make([]int, 0, max(0, len(records)-start))
	for row := start; row < len(records); row++ {
		viewRows = append(viewRows, row)
	}
	sort.SliceStable(viewRows, func(i, j int) bool {
		return compareSortRecords(records[viewRows[i]], records[viewRows[j]], keys) < 0
	})

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
	default:
		return nil, fmt.Errorf("unsupported text transform %q", mode)
	}
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

func (manager *Manager) Refresh(id string) (*Session, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return nil, errors.New("document not found")
	}
	if session.Dirty {
		return nil, errors.New("document has unsaved edits")
	}
	path := session.trackedPath()
	if path == "" {
		return nil, errors.New("document has no file path")
	}
	delimiter := session.Document.Metadata.Dialect.Delimiter
	hasHeader := session.Document.Metadata.Dialect.HasHeader
	refreshed, err := csvcore.OpenDocument(path, csvcore.OpenOptions{
		PreviewRows: max(160, len(session.Document.Preview)),
		IndexRows:   -1,
		Delimiter:   &delimiter,
		HasHeader:   &hasHeader,
		Encoding:    &session.Document.Metadata.Encoding,
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
	current.WorkingPath = path
	current.Document = refreshed
	current.Edits = make(map[Cell]string)
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return MutationResult{}, err
	}
	applyEdits(records, 0, session.Edits)
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
	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return historySnapshot{}, err
	}
	applyEdits(records, 0, session.Edits)
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
	hasHeader := metadata.Dialect.HasHeader
	refreshed, err := csvcore.OpenDocument(tempPath, csvcore.OpenOptions{
		PreviewRows: max(160, min(len(records), 1000)),
		IndexRows:   -1,
		Delimiter:   &delimiter,
		HasHeader:   &hasHeader,
		Encoding:    &metadata.Encoding,
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
	copy.ViewRows = append([]int(nil), session.ViewRows...)
	copy.UndoStack = cloneHistory(session.UndoStack)
	copy.RedoStack = cloneHistory(session.RedoStack)
	return &copy
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

func visibleRowCount(session *Session, recordCount int, dataStart int) int {
	if session.ViewActive {
		return len(session.ViewRows)
	}
	return max(0, recordCount-dataStart)
}

func joinRecords(records [][]string, cells []Cell, separator string) int {
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

type splitOperation struct {
	row   int
	col   int
	parts []string
}

func splitRecords(records [][]string, cells []Cell, separator string) int {
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

func exportViewFilename(displayName string, delimiter rune) string {
	base := filepath.Base(displayName)
	if base == "." || base == string(filepath.Separator) || base == "" {
		base = "opencsv"
	}
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	if stem == "" {
		stem = "opencsv"
	}
	targetExt := ".csv"
	if delimiter == '\t' {
		targetExt = ".tsv"
	}
	return stem + ".view" + targetExt
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
