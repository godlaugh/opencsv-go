package csvcore

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

const (
	QuoteModeMinimal = "minimal"
	QuoteModeAlways  = "always"
	QuoteModeNever   = "never"
	QuoteModeManual  = "manual"
)

func NormalizeQuoteMode(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", QuoteModeMinimal:
		return QuoteModeMinimal, nil
	case QuoteModeAlways:
		return QuoteModeAlways, nil
	case QuoteModeNever, "none":
		return QuoteModeNever, nil
	case QuoteModeManual:
		return QuoteModeManual, nil
	default:
		return "", fmt.Errorf("quote mode must be minimal, always, never, or manual")
	}
}

func NormalizeQuoteChar(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "double", "quote":
		return `"`, nil
	case "single":
		return "'", nil
	case "none":
		return "", nil
	}
	runes := []rune(value)
	if len(runes) != 1 {
		return "", fmt.Errorf("quote character must be double, single, none, or one character")
	}
	return value, nil
}

func DialectQuoteMode(dialect Dialect) string {
	mode, err := NormalizeQuoteMode(dialect.QuoteMode)
	if err != nil {
		return QuoteModeMinimal
	}
	return mode
}

func DialectQuoteChar(dialect Dialect) string {
	quoteChar, err := NormalizeQuoteChar(dialect.QuoteChar)
	if err != nil {
		return `"`
	}
	return quoteChar
}

func ShouldQuote(value string, dialect Dialect) bool {
	return ShouldQuoteField(value, dialect, -1, -1)
}

func ShouldQuoteField(value string, dialect Dialect, row int, col int) bool {
	if SuppressesQuote(dialect, row, col) {
		return false
	}
	switch DialectQuoteMode(dialect) {
	case QuoteModeAlways:
		return true
	case QuoteModeNever:
		return false
	case QuoteModeManual:
		return PreservesQuote(dialect, row, col)
	default:
		return RequiresQuote(value, dialect) || PreservesQuote(dialect, row, col)
	}
}

func RequiresQuote(value string, dialect Dialect) bool {
	if value == "" {
		return false
	}
	delimiter := dialect.Delimiter
	if delimiter == 0 {
		delimiter = ','
	}
	if strings.ContainsRune(value, delimiter) {
		return true
	}
	if strings.ContainsAny(value, "\r\n") {
		return true
	}
	if quote := DialectQuoteChar(dialect); quote != "" && strings.Contains(value, quote) {
		return true
	}
	first, _ := firstRune(value)
	return unicode.IsSpace(first)
}

func firstRune(value string) (rune, bool) {
	for _, r := range value {
		return r, true
	}
	return 0, false
}

func PreservesQuote(dialect Dialect, row int, col int) bool {
	if row < 0 || col < 0 {
		return false
	}
	if dialect.HasHeader && dialect.QuoteHeader && row == 0 {
		return true
	}
	if dialect.QuoteBody && isBodyRow(dialect, row) {
		return true
	}
	for _, quoteCol := range dialect.QuoteColumns {
		if quoteCol == col && isBodyRow(dialect, row) {
			return true
		}
	}
	for _, cell := range dialect.QuoteCells {
		if cell.Row == row && cell.Col == col {
			return true
		}
	}
	return false
}

func isBodyRow(dialect Dialect, row int) bool {
	return !dialect.HasHeader || row > 0
}

func SuppressesQuote(dialect Dialect, row int, col int) bool {
	if row < 0 || col < 0 {
		return false
	}
	for _, cell := range dialect.NoQuoteCells {
		if cell.Row == row && cell.Col == col {
			return true
		}
	}
	return false
}

func SetCellQuote(dialect Dialect, row int, col int, quoted bool) (Dialect, bool) {
	cell := Cell{Row: row, Col: col}
	beforeQuoteCells := append([]Cell(nil), dialect.QuoteCells...)
	beforeNoQuoteCells := append([]Cell(nil), dialect.NoQuoteCells...)
	if quoted {
		dialect.NoQuoteCells = removeCell(dialect.NoQuoteCells, cell)
		if !PreservesQuote(dialect, row, col) {
			dialect.QuoteCells = addCell(dialect.QuoteCells, cell)
		}
	} else {
		dialect.QuoteCells = removeCell(dialect.QuoteCells, cell)
		dialect.NoQuoteCells = addCell(dialect.NoQuoteCells, cell)
	}
	sortCells(dialect.QuoteCells)
	sortCells(dialect.NoQuoteCells)
	return dialect, !sameCells(beforeQuoteCells, dialect.QuoteCells) || !sameCells(beforeNoQuoteCells, dialect.NoQuoteCells)
}

func SetQuoteHeader(dialect Dialect, quoted bool) (Dialect, bool) {
	if dialect.QuoteHeader == quoted {
		return dialect, false
	}
	dialect.QuoteHeader = quoted
	return dialect, true
}

func SetQuoteBody(dialect Dialect, quoted bool) (Dialect, bool) {
	if dialect.QuoteBody == quoted {
		return dialect, false
	}
	dialect.QuoteBody = quoted
	return dialect, true
}

func SetQuoteColumn(dialect Dialect, col int, quoted bool) (Dialect, bool) {
	before := append([]int(nil), dialect.QuoteColumns...)
	if quoted {
		dialect.QuoteColumns = addInt(dialect.QuoteColumns, col)
	} else {
		dialect.QuoteColumns = removeInt(dialect.QuoteColumns, col)
	}
	sort.Ints(dialect.QuoteColumns)
	return dialect, !sameInts(before, dialect.QuoteColumns)
}

func DetectQuoteRules(sample []byte, dialect Dialect) Dialect {
	flags := quotedFieldFlags(sample, dialect, 80)
	if len(flags) == 0 {
		return dialect
	}

	coveredHeader := false
	if dialect.HasHeader && len(flags[0]) > 0 && allTrue(flags[0]) {
		dialect.QuoteHeader = true
		coveredHeader = true
	}

	startRow := 0
	if dialect.HasHeader {
		startRow = 1
	}
	seen := map[int]int{}
	quoted := map[int]int{}
	for row := startRow; row < len(flags); row++ {
		for col, wasQuoted := range flags[row] {
			seen[col]++
			if wasQuoted {
				quoted[col]++
			}
		}
	}
	quoteColumns := map[int]struct{}{}
	for col, total := range seen {
		if total > 0 && quoted[col] == total {
			dialect.QuoteColumns = append(dialect.QuoteColumns, col)
			quoteColumns[col] = struct{}{}
		}
	}
	sort.Ints(dialect.QuoteColumns)

	const maxQuoteCells = 1000
	for row, record := range flags {
		for col, wasQuoted := range record {
			if !wasQuoted {
				continue
			}
			if coveredHeader && row == 0 {
				continue
			}
			if _, ok := quoteColumns[col]; ok && (!dialect.HasHeader || row > 0) {
				continue
			}
			dialect.QuoteCells = append(dialect.QuoteCells, Cell{Row: row, Col: col})
			if len(dialect.QuoteCells) >= maxQuoteCells {
				return dialect
			}
		}
	}
	return dialect
}

func addCell(cells []Cell, cell Cell) []Cell {
	for _, candidate := range cells {
		if candidate == cell {
			return cells
		}
	}
	return append(cells, cell)
}

func removeCell(cells []Cell, cell Cell) []Cell {
	result := cells[:0]
	for _, candidate := range cells {
		if candidate != cell {
			result = append(result, candidate)
		}
	}
	return result
}

func sortCells(cells []Cell) {
	sort.Slice(cells, func(i, j int) bool {
		if cells[i].Row == cells[j].Row {
			return cells[i].Col < cells[j].Col
		}
		return cells[i].Row < cells[j].Row
	})
}

func sameCells(left []Cell, right []Cell) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func addInt(values []int, value int) []int {
	for _, candidate := range values {
		if candidate == value {
			return values
		}
	}
	return append(values, value)
}

func removeInt(values []int, value int) []int {
	result := values[:0]
	for _, candidate := range values {
		if candidate != value {
			result = append(result, candidate)
		}
	}
	return result
}

func sameInts(left []int, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func allTrue(values []bool) bool {
	if len(values) == 0 {
		return false
	}
	for _, value := range values {
		if !value {
			return false
		}
	}
	return true
}

func quotedFieldFlags(sample []byte, dialect Dialect, maxRows int) [][]bool {
	delimiter := dialect.Delimiter
	if delimiter == 0 {
		delimiter = ','
	}
	quote := []rune(DialectQuoteChar(dialect))
	if len(quote) != 1 {
		return nil
	}
	input := []rune(strings.TrimPrefix(string(sample), "\ufeff"))
	records := make([][]bool, 0, maxRows)
	record := []bool{}
	quoted := false
	atStart := true
	inQuotes := false
	for index := 0; index < len(input) && len(records) < maxRows; index++ {
		char := input[index]
		if inQuotes {
			if char == quote[0] {
				if index+1 < len(input) && input[index+1] == quote[0] {
					index++
				} else {
					inQuotes = false
				}
			}
			continue
		}
		if atStart && char == quote[0] {
			quoted = true
			atStart = false
			inQuotes = true
			continue
		}
		if char == delimiter {
			record = append(record, quoted)
			quoted = false
			atStart = true
			continue
		}
		if char == '\n' || char == '\r' {
			record = append(record, quoted)
			records = append(records, record)
			record = []bool{}
			quoted = false
			atStart = true
			if char == '\r' && index+1 < len(input) && input[index+1] == '\n' {
				index++
			}
			continue
		}
		atStart = false
	}
	if len(record) > 0 || !atStart || quoted {
		record = append(record, quoted)
		records = append(records, record)
	}
	return records
}
