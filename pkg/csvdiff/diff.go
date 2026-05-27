package csvdiff

import (
	"errors"
	"strings"

	"github.com/opencsv/opencsv/pkg/csvcore"
)

type RowStatus string

const (
	RowUnchanged RowStatus = "unchanged"
	RowChanged   RowStatus = "changed"
	RowAdded     RowStatus = "added"
	RowDeleted   RowStatus = "deleted"
)

type Options struct {
	PreviewRows        int
	MaxRows            int
	IgnoreWhitespace   bool
	IgnoreCase         bool
	IgnoredColumns     []int
	MatchColumnsByName bool
}

type Result struct {
	LeftPath  string    `json:"leftPath"`
	RightPath string    `json:"rightPath"`
	Rows      []RowDiff `json:"rows"`
	Summary   Summary   `json:"summary"`
	Truncated bool      `json:"truncated"`
}

type Summary struct {
	Unchanged int `json:"unchanged"`
	Changed   int `json:"changed"`
	Added     int `json:"added"`
	Deleted   int `json:"deleted"`
}

type RowDiff struct {
	Status       RowStatus  `json:"status"`
	LeftRow      *int       `json:"leftRow,omitempty"`
	RightRow     *int       `json:"rightRow,omitempty"`
	LeftRecord   []string   `json:"leftRecord,omitempty"`
	RightRecord  []string   `json:"rightRecord,omitempty"`
	ChangedCells []CellDiff `json:"changedCells,omitempty"`
}

type CellDiff struct {
	Col        int        `json:"col"`
	Left       string     `json:"left"`
	Right      string     `json:"right"`
	LeftSpans  []TextSpan `json:"leftSpans,omitempty"`
	RightSpans []TextSpan `json:"rightSpans,omitempty"`
}

type TextSpan struct {
	Text    string `json:"text"`
	Changed bool   `json:"changed"`
}

type MergeMode string

const (
	MergeRow  MergeMode = "row"
	MergeCell MergeMode = "cell"
)

type MergeOptions struct {
	DiffOptions Options
	ResultIndex int
	Mode        MergeMode
	Col         int
}

type MergeResult struct {
	Save    csvcore.SaveResult `json:"save"`
	Applied string             `json:"applied"`
}

func CompareFiles(leftPath string, rightPath string, options Options) (Result, error) {
	if options.PreviewRows <= 0 {
		options.PreviewRows = 50
	}
	leftDoc, err := csvcore.OpenDocument(leftPath, csvcore.OpenOptions{PreviewRows: options.PreviewRows, IndexRows: -1})
	if err != nil {
		return Result{}, err
	}
	rightDoc, err := csvcore.OpenDocument(rightPath, csvcore.OpenOptions{PreviewRows: options.PreviewRows, IndexRows: -1})
	if err != nil {
		return Result{}, err
	}
	leftRecords, err := csvcore.ReadAllRecords(leftPath, leftDoc.Metadata)
	if err != nil {
		return Result{}, err
	}
	rightRecords, err := csvcore.ReadAllRecords(rightPath, rightDoc.Metadata)
	if err != nil {
		return Result{}, err
	}
	return CompareRecords(leftPath, rightPath, leftRecords, rightRecords, options), nil
}

func CompareRecords(leftPath string, rightPath string, leftRecords [][]string, rightRecords [][]string, options Options) Result {
	displayLeft, displayRight := displayRecords(leftRecords, rightRecords, options)
	compareLeft := comparableRecords(displayLeft, options)
	compareRight := comparableRecords(displayRight, options)
	ops := diffOps(compareLeft, compareRight)
	rows := materializeRows(displayLeft, displayRight, compareLeft, compareRight, ops, options)
	result := Result{LeftPath: leftPath, RightPath: rightPath, Rows: rows}
	if options.MaxRows > 0 && len(result.Rows) > options.MaxRows {
		result.Rows = result.Rows[:options.MaxRows]
		result.Truncated = true
	}
	for _, row := range result.Rows {
		switch row.Status {
		case RowUnchanged:
			result.Summary.Unchanged++
		case RowChanged:
			result.Summary.Changed++
		case RowAdded:
			result.Summary.Added++
		case RowDeleted:
			result.Summary.Deleted++
		}
	}
	return result
}

func MergeRightToLeft(leftPath string, rightPath string, options MergeOptions) (MergeResult, error) {
	leftDoc, err := csvcore.OpenDocument(leftPath, csvcore.OpenOptions{PreviewRows: 50, IndexRows: -1})
	if err != nil {
		return MergeResult{}, err
	}
	rightDoc, err := csvcore.OpenDocument(rightPath, csvcore.OpenOptions{PreviewRows: 50, IndexRows: -1})
	if err != nil {
		return MergeResult{}, err
	}
	leftRecords, err := csvcore.ReadAllRecords(leftPath, leftDoc.Metadata)
	if err != nil {
		return MergeResult{}, err
	}
	rightRecords, err := csvcore.ReadAllRecords(rightPath, rightDoc.Metadata)
	if err != nil {
		return MergeResult{}, err
	}
	result := CompareRecords(leftPath, rightPath, leftRecords, rightRecords, options.DiffOptions)
	if options.ResultIndex < 0 || options.ResultIndex >= len(result.Rows) {
		return MergeResult{}, errors.New("diff result index is out of range")
	}
	row := result.Rows[options.ResultIndex]
	projection := buildProjection(leftRecords, rightRecords, options.DiffOptions)

	applied, err := applyMerge(&leftRecords, rightRecords, row, projection, options)
	if err != nil {
		return MergeResult{}, err
	}
	save, err := csvcore.SaveRecords(leftPath, leftDoc.Metadata, leftRecords, csvcore.SaveOptions{CreateBackup: true})
	if err != nil {
		return MergeResult{}, err
	}
	return MergeResult{Save: save, Applied: applied}, nil
}

type opKind int

const (
	opEqual opKind = iota
	opDelete
	opAdd
)

type op struct {
	kind  opKind
	left  int
	right int
}

func diffOps(left [][]string, right [][]string) []op {
	lcs := make([][]int, len(left)+1)
	for i := range lcs {
		lcs[i] = make([]int, len(right)+1)
	}
	for i := len(left) - 1; i >= 0; i-- {
		for j := len(right) - 1; j >= 0; j-- {
			if recordsEqual(left[i], right[j]) {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}

	var ops []op
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if recordsEqual(left[i], right[j]) {
			ops = append(ops, op{kind: opEqual, left: i, right: j})
			i++
			j++
		} else if lcs[i+1][j] >= lcs[i][j+1] {
			ops = append(ops, op{kind: opDelete, left: i})
			i++
		} else {
			ops = append(ops, op{kind: opAdd, right: j})
			j++
		}
	}
	for i < len(left) {
		ops = append(ops, op{kind: opDelete, left: i})
		i++
	}
	for j < len(right) {
		ops = append(ops, op{kind: opAdd, right: j})
		j++
	}
	return ops
}

func materializeRows(left [][]string, right [][]string, compareLeft [][]string, compareRight [][]string, ops []op, options Options) []RowDiff {
	var rows []RowDiff
	for index := 0; index < len(ops); {
		if ops[index].kind == opEqual {
			leftRow := displayRow(ops[index].left)
			rightRow := displayRow(ops[index].right)
			rows = append(rows, RowDiff{
				Status:      RowUnchanged,
				LeftRow:     &leftRow,
				RightRow:    &rightRow,
				LeftRecord:  left[ops[index].left],
				RightRecord: right[ops[index].right],
			})
			index++
			continue
		}

		start := index
		for index < len(ops) && ops[index].kind != opEqual {
			index++
		}
		rows = append(rows, materializeChangeBlock(left, right, compareLeft, compareRight, ops[start:index], options)...)
	}
	return rows
}

func materializeChangeBlock(left [][]string, right [][]string, compareLeft [][]string, compareRight [][]string, ops []op, options Options) []RowDiff {
	var deleted []int
	var added []int
	for _, op := range ops {
		if op.kind == opDelete {
			deleted = append(deleted, op.left)
		}
		if op.kind == opAdd {
			added = append(added, op.right)
		}
	}

	var rows []RowDiff
	leftPos, rightPos := 0, 0
	for leftPos < len(deleted) && rightPos < len(added) {
		leftIndex := deleted[leftPos]
		rightIndex := added[rightPos]
		if likelySameLogicalRow(compareLeft[leftIndex], compareRight[rightIndex]) {
			leftRow := displayRow(leftIndex)
			rightRow := displayRow(rightIndex)
			rows = append(rows, RowDiff{
				Status:       RowChanged,
				LeftRow:      &leftRow,
				RightRow:     &rightRow,
				LeftRecord:   left[leftIndex],
				RightRecord:  right[rightIndex],
				ChangedCells: changedCells(left[leftIndex], right[rightIndex], compareLeft[leftIndex], compareRight[rightIndex], options),
			})
			leftPos++
			rightPos++
			continue
		}
		leftRow := displayRow(leftIndex)
		rows = append(rows, RowDiff{
			Status:     RowDeleted,
			LeftRow:    &leftRow,
			LeftRecord: left[leftIndex],
		})
		leftPos++
	}
	for _, leftIndex := range deleted[leftPos:] {
		leftRow := displayRow(leftIndex)
		rows = append(rows, RowDiff{
			Status:     RowDeleted,
			LeftRow:    &leftRow,
			LeftRecord: left[leftIndex],
		})
	}
	for _, rightIndex := range added[rightPos:] {
		rightRow := displayRow(rightIndex)
		rows = append(rows, RowDiff{
			Status:      RowAdded,
			RightRow:    &rightRow,
			RightRecord: right[rightIndex],
		})
	}
	return rows
}

func likelySameLogicalRow(left []string, right []string) bool {
	if len(left) == 0 || len(right) == 0 {
		return len(left) == len(right)
	}
	return left[0] == right[0]
}

func changedCells(left []string, right []string, compareLeft []string, compareRight []string, options Options) []CellDiff {
	width := max(len(left), len(right))
	cells := make([]CellDiff, 0)
	for col := 0; col < width; col++ {
		leftValue := cell(left, col)
		rightValue := cell(right, col)
		if cell(compareLeft, col) != cell(compareRight, col) {
			leftSpans, rightSpans := inlineSpans(leftValue, rightValue, options)
			cells = append(cells, CellDiff{
				Col:        col + 1,
				Left:       leftValue,
				Right:      rightValue,
				LeftSpans:  leftSpans,
				RightSpans: rightSpans,
			})
		}
	}
	return cells
}

func recordsEqual(left []string, right []string) bool {
	return strings.Join(left, "\x00") == strings.Join(right, "\x00")
}

func displayRecords(left [][]string, right [][]string, options Options) ([][]string, [][]string) {
	projection := buildProjection(left, right, options)
	return projectRecords(left, projection.Left), projectRecords(right, projection.Right)
}

type columnProjection struct {
	Left  []int
	Right []int
}

func buildProjection(left [][]string, right [][]string, options Options) columnProjection {
	projection := defaultProjection(left, right)
	if options.MatchColumnsByName {
		projection = headerProjection(left, right)
	}
	if len(options.IgnoredColumns) > 0 {
		projection = removeIgnoredProjectionColumns(projection, ignoredColumnSet(options.IgnoredColumns))
	}
	return projection
}

func defaultProjection(left [][]string, right [][]string) columnProjection {
	width := max(maxRecordWidth(left), maxRecordWidth(right))
	projection := columnProjection{
		Left:  make([]int, width),
		Right: make([]int, width),
	}
	for index := 0; index < width; index++ {
		projection.Left[index] = index
		projection.Right[index] = index
	}
	return projection
}

func maxRecordWidth(records [][]string) int {
	width := 0
	for _, record := range records {
		width = max(width, len(record))
	}
	return width
}

func headerProjection(left [][]string, right [][]string) columnProjection {
	if len(left) == 0 || len(right) == 0 {
		return defaultProjection(left, right)
	}
	leftHeader := left[0]
	rightHeader := right[0]
	rightByHeader := make(map[string]int, len(rightHeader))
	for index, header := range rightHeader {
		if _, exists := rightByHeader[header]; !exists {
			rightByHeader[header] = index
		}
	}

	usedRight := make(map[int]bool, len(rightHeader))
	leftColumns := make([]int, 0, len(leftHeader)+len(rightHeader))
	rightColumns := make([]int, 0, len(leftHeader)+len(rightHeader))
	for leftIndex, header := range leftHeader {
		leftColumns = append(leftColumns, leftIndex)
		if rightIndex, ok := rightByHeader[header]; ok {
			rightColumns = append(rightColumns, rightIndex)
			usedRight[rightIndex] = true
		} else {
			rightColumns = append(rightColumns, -1)
		}
	}
	for rightIndex := range rightHeader {
		if usedRight[rightIndex] {
			continue
		}
		leftColumns = append(leftColumns, -1)
		rightColumns = append(rightColumns, rightIndex)
	}
	return columnProjection{Left: leftColumns, Right: rightColumns}
}

func ignoredColumnSet(columns []int) map[int]bool {
	ignored := make(map[int]bool, len(columns))
	for _, column := range columns {
		if column > 0 {
			ignored[column-1] = true
		}
	}
	return ignored
}

func removeIgnoredProjectionColumns(projection columnProjection, ignored map[int]bool) columnProjection {
	next := columnProjection{}
	for index := range projection.Left {
		if ignored[index] {
			continue
		}
		next.Left = append(next.Left, projection.Left[index])
		next.Right = append(next.Right, projection.Right[index])
	}
	return next
}

func projectRecords(records [][]string, columns []int) [][]string {
	result := make([][]string, 0, len(records))
	for _, record := range records {
		next := make([]string, len(columns))
		for index, source := range columns {
			if source >= 0 && source < len(record) {
				next[index] = record[source]
			}
		}
		result = append(result, next)
	}
	return result
}

func cloneRecords(records [][]string) [][]string {
	result := make([][]string, len(records))
	for index, record := range records {
		result[index] = append([]string(nil), record...)
	}
	return result
}

func comparableRecords(records [][]string, options Options) [][]string {
	result := make([][]string, len(records))
	for rowIndex, record := range records {
		next := make([]string, len(record))
		for colIndex, value := range record {
			next[colIndex] = comparableCell(value, options)
		}
		result[rowIndex] = next
	}
	return result
}

func comparableCell(value string, options Options) string {
	if options.IgnoreWhitespace {
		value = strings.Join(strings.Fields(value), "")
	}
	if options.IgnoreCase {
		value = strings.ToLower(value)
	}
	return value
}

func inlineSpans(left string, right string, options Options) ([]TextSpan, []TextSpan) {
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	lcs := make([][]int, len(leftRunes)+1)
	for index := range lcs {
		lcs[index] = make([]int, len(rightRunes)+1)
	}
	for i := len(leftRunes) - 1; i >= 0; i-- {
		for j := len(rightRunes) - 1; j >= 0; j-- {
			if sameInlineRune(leftRunes[i], rightRunes[j], options) {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}

	var leftSpans []TextSpan
	var rightSpans []TextSpan
	i, j := 0, 0
	for i < len(leftRunes) && j < len(rightRunes) {
		if sameInlineRune(leftRunes[i], rightRunes[j], options) {
			leftSpans = appendSpan(leftSpans, string(leftRunes[i]), false)
			rightSpans = appendSpan(rightSpans, string(rightRunes[j]), false)
			i++
			j++
		} else if lcs[i+1][j] >= lcs[i][j+1] {
			leftSpans = appendSpan(leftSpans, string(leftRunes[i]), true)
			i++
		} else {
			rightSpans = appendSpan(rightSpans, string(rightRunes[j]), true)
			j++
		}
	}
	for i < len(leftRunes) {
		leftSpans = appendSpan(leftSpans, string(leftRunes[i]), true)
		i++
	}
	for j < len(rightRunes) {
		rightSpans = appendSpan(rightSpans, string(rightRunes[j]), true)
		j++
	}
	return leftSpans, rightSpans
}

func sameInlineRune(left rune, right rune, options Options) bool {
	if left == right {
		return true
	}
	if options.IgnoreCase {
		return strings.EqualFold(string(left), string(right))
	}
	return false
}

func appendSpan(spans []TextSpan, text string, changed bool) []TextSpan {
	if text == "" {
		return spans
	}
	if len(spans) > 0 && spans[len(spans)-1].Changed == changed {
		spans[len(spans)-1].Text += text
		return spans
	}
	return append(spans, TextSpan{Text: text, Changed: changed})
}

func applyMerge(leftRecords *[][]string, rightRecords [][]string, row RowDiff, projection columnProjection, options MergeOptions) (string, error) {
	switch options.Mode {
	case MergeCell:
		return applyCellMerge(leftRecords, rightRecords, row, projection, options.Col)
	case MergeRow:
		return applyRowMerge(leftRecords, rightRecords, row, projection)
	default:
		return "", errors.New("unsupported merge mode")
	}
}

func applyCellMerge(leftRecords *[][]string, rightRecords [][]string, row RowDiff, projection columnProjection, col int) (string, error) {
	if row.Status != RowChanged {
		return "", errors.New("cell merge requires a changed row")
	}
	if row.LeftRow == nil || row.RightRow == nil {
		return "", errors.New("changed row is missing row numbers")
	}
	leftCol, rightCol, err := projectedColumns(projection, col)
	if err != nil {
		return "", err
	}
	leftIndex := *row.LeftRow - 1
	rightIndex := *row.RightRow - 1
	if leftIndex < 0 || leftIndex >= len(*leftRecords) || rightIndex < 0 || rightIndex >= len(rightRecords) {
		return "", errors.New("merge row is out of range")
	}
	ensureWidth(&(*leftRecords)[leftIndex], leftCol+1)
	(*leftRecords)[leftIndex][leftCol] = cell(rightRecords[rightIndex], rightCol)
	return "cell", nil
}

func applyRowMerge(leftRecords *[][]string, rightRecords [][]string, row RowDiff, projection columnProjection) (string, error) {
	switch row.Status {
	case RowChanged:
		if row.LeftRow == nil || row.RightRow == nil {
			return "", errors.New("changed row is missing row numbers")
		}
		leftIndex := *row.LeftRow - 1
		rightIndex := *row.RightRow - 1
		if leftIndex < 0 || leftIndex >= len(*leftRecords) || rightIndex < 0 || rightIndex >= len(rightRecords) {
			return "", errors.New("merge row is out of range")
		}
		merged, err := projectedRightRecordForLeft((*leftRecords)[leftIndex], rightRecords[rightIndex], projection)
		if err != nil {
			return "", err
		}
		(*leftRecords)[leftIndex] = merged
		return "row", nil
	case RowAdded:
		if row.RightRow == nil {
			return "", errors.New("added row is missing a right row number")
		}
		rightIndex := *row.RightRow - 1
		if rightIndex < 0 || rightIndex >= len(rightRecords) {
			return "", errors.New("right row is out of range")
		}
		emptyLeft := make([]string, maxRecordWidth(*leftRecords))
		merged, err := projectedRightRecordForLeft(emptyLeft, rightRecords[rightIndex], projection)
		if err != nil {
			return "", err
		}
		insertIndex := min(max(0, rightIndex), len(*leftRecords))
		*leftRecords = append((*leftRecords)[:insertIndex], append([][]string{merged}, (*leftRecords)[insertIndex:]...)...)
		return "added row", nil
	case RowDeleted:
		if row.LeftRow == nil {
			return "", errors.New("deleted row is missing a left row number")
		}
		leftIndex := *row.LeftRow - 1
		if leftIndex < 0 || leftIndex >= len(*leftRecords) {
			return "", errors.New("left row is out of range")
		}
		*leftRecords = append((*leftRecords)[:leftIndex], (*leftRecords)[leftIndex+1:]...)
		return "deleted row", nil
	default:
		return "", errors.New("unchanged rows cannot be merged")
	}
}

func projectedColumns(projection columnProjection, col int) (int, int, error) {
	index := col - 1
	if index < 0 || index >= len(projection.Left) || index >= len(projection.Right) {
		return 0, 0, errors.New("merge column is out of range")
	}
	leftCol := projection.Left[index]
	rightCol := projection.Right[index]
	if leftCol < 0 || rightCol < 0 {
		return 0, 0, errors.New("cannot merge a one-sided column")
	}
	return leftCol, rightCol, nil
}

func projectedRightRecordForLeft(left []string, right []string, projection columnProjection) ([]string, error) {
	merged := append([]string(nil), left...)
	for index := range projection.Left {
		leftCol := projection.Left[index]
		rightCol := projection.Right[index]
		if leftCol < 0 {
			continue
		}
		if rightCol < 0 {
			continue
		}
		ensureWidth(&merged, leftCol+1)
		merged[leftCol] = cell(right, rightCol)
	}
	return merged, nil
}

func ensureWidth(record *[]string, width int) {
	for len(*record) < width {
		*record = append(*record, "")
	}
}

func cell(record []string, col int) string {
	if col < 0 || col >= len(record) {
		return ""
	}
	return record[col]
}

func displayRow(index int) int {
	return index + 1
}
