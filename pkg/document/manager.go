package document

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

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
	UndoStack                []historySnapshot `json:"-"`
	RedoStack                []historySnapshot `json:"-"`
	Dirty                    bool              `json:"dirty"`
	ReadOnly                 bool              `json:"readOnly"`
	LastKnownPath            string            `json:"-"`
	LastKnownSizeBytes       int64             `json:"-"`
	LastKnownModTimeUnixNano int64             `json:"-"`
}

type historySnapshot struct {
	Records  [][]string
	Metadata csvcore.FileMetadata
	Dirty    bool
	ViewRows []int
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

	if len(session.ViewRows) == 0 {
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

func (manager *Manager) Find(id string, options FindOptions) ([]Match, error) {
	session, ok := manager.snapshot(id)
	if !ok {
		return nil, errors.New("document not found")
	}
	if options.Query == "" {
		return nil, nil
	}
	if options.Limit <= 0 {
		options.Limit = 100
	}
	matcher, err := newMatcher(options)
	if err != nil {
		return nil, err
	}

	records, err := csvcore.ReadAllRecords(session.WorkingPath, session.Document.Metadata)
	if err != nil {
		return nil, err
	}
	applyEdits(records, 0, session.Edits)

	var matches []Match
	for row, record := range records {
		for col, value := range record {
			if options.Column != nil && col != *options.Column {
				continue
			}
			if matcher.matches(value) {
				matches = append(matches, Match{Row: row, Col: col, Value: value})
				if len(matches) >= options.Limit {
					return matches, nil
				}
			}
		}
	}
	return matches, nil
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
	for row, record := range records {
		for col, value := range record {
			if options.Column != nil && col != *options.Column {
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
	return len(viewRows), nil
}

type SortOptions struct {
	Column     int       `json:"column"`
	Descending bool      `json:"descending"`
	Numeric    bool      `json:"numeric"`
	HasHeader  bool      `json:"hasHeader"`
	Keys       []SortKey `json:"keys"`
}

type SortKey struct {
	Column     int  `json:"column"`
	Descending bool `json:"descending"`
	Numeric    bool `json:"numeric"`
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
	return len(viewRows), nil
}

func normalizeSortKeys(options SortOptions) ([]SortKey, error) {
	keys := append([]SortKey(nil), options.Keys...)
	if len(keys) == 0 {
		keys = []SortKey{{
			Column:     options.Column,
			Descending: options.Descending,
			Numeric:    options.Numeric,
		}}
	}
	for _, key := range keys {
		if key.Column < 0 {
			return nil, errors.New("column must be non-negative")
		}
	}
	return keys, nil
}

func compareSortRecords(leftRecord []string, rightRecord []string, keys []SortKey) int {
	for _, key := range keys {
		left := cellValue(leftRecord, key.Column)
		right := cellValue(rightRecord, key.Column)
		result := compareSortValues(left, right, key.Numeric)
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

func compareSortValues(left string, right string, numeric bool) int {
	if numeric {
		leftNum, leftErr := strconv.ParseFloat(strings.TrimSpace(left), 64)
		rightNum, rightErr := strconv.ParseFloat(strings.TrimSpace(right), 64)
		if leftErr == nil && rightErr == nil {
			switch {
			case leftNum < rightNum:
				return -1
			case leftNum > rightNum:
				return 1
			default:
				return 0
			}
		}
	}
	return strings.Compare(left, right)
}

func (manager *Manager) ClearView(id string) (int, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	session, ok := manager.sessions[id]
	if !ok {
		return 0, errors.New("document not found")
	}
	session.ViewRows = nil
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
		Records:  cloneRecords(records),
		Metadata: session.Document.Metadata,
		Dirty:    session.Dirty,
		ViewRows: append([]int(nil), session.ViewRows...),
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
			Records:  cloneRecords(snapshot.Records),
			Metadata: snapshot.Metadata,
			Dirty:    snapshot.Dirty,
			ViewRows: append([]int(nil), snapshot.ViewRows...),
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

func maxRecordWidth(records [][]string) int {
	width := 0
	for _, record := range records {
		width = max(width, len(record))
	}
	return width
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
