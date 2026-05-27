package csvdiff

import (
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
	PreviewRows int
	MaxRows     int
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
	Col   int    `json:"col"`
	Left  string `json:"left"`
	Right string `json:"right"`
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
	ops := diffOps(leftRecords, rightRecords)
	rows := materializeRows(leftRecords, rightRecords, ops)
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
			if sameRecord(left[i], right[j]) {
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
		if sameRecord(left[i], right[j]) {
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

func materializeRows(left [][]string, right [][]string, ops []op) []RowDiff {
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
		rows = append(rows, materializeChangeBlock(left, right, ops[start:index])...)
	}
	return rows
}

func materializeChangeBlock(left [][]string, right [][]string, ops []op) []RowDiff {
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
		if likelySameLogicalRow(left[leftIndex], right[rightIndex]) {
			leftRow := displayRow(leftIndex)
			rightRow := displayRow(rightIndex)
			rows = append(rows, RowDiff{
				Status:       RowChanged,
				LeftRow:      &leftRow,
				RightRow:     &rightRow,
				LeftRecord:   left[leftIndex],
				RightRecord:  right[rightIndex],
				ChangedCells: changedCells(left[leftIndex], right[rightIndex]),
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

func changedCells(left []string, right []string) []CellDiff {
	width := max(len(left), len(right))
	cells := make([]CellDiff, 0)
	for col := 0; col < width; col++ {
		leftValue := cell(left, col)
		rightValue := cell(right, col)
		if leftValue != rightValue {
			cells = append(cells, CellDiff{Col: col + 1, Left: leftValue, Right: rightValue})
		}
	}
	return cells
}

func sameRecord(left []string, right []string) bool {
	return strings.Join(left, "\x00") == strings.Join(right, "\x00")
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
