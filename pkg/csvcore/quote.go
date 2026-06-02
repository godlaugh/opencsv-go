package csvcore

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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

const maxDetectedQuoteCells = 100_000

func DetectQuoteRules(sample []byte, dialect Dialect) Dialect {
	flags := quotedFieldFlags(sample, dialect, 80)
	return applyDetectedQuoteRules(quoteRuleStatsFromFlags(flags, dialect), dialect)
}

func DetectQuoteRulesForFile(path string, dialect Dialect, encoding EncodingInfo) (Dialect, error) {
	if _, _, ok := rawQuoteScanBytes(dialect); !ok {
		return dialect, nil
	}
	if encoding.Name == EncodingUTF16LE || encoding.Name == EncodingUTF16BE {
		raw, err := os.ReadFile(path)
		if err != nil {
			return Dialect{}, err
		}
		utf8Bytes, err := DecodeToUTF8(raw, encoding)
		if err != nil {
			return Dialect{}, err
		}
		return DetectQuoteRulesAll(utf8Bytes, dialect), nil
	}
	file, err := os.Open(path)
	if err != nil {
		return Dialect{}, err
	}
	defer file.Close()
	stats, err := scanQuoteRuleStats(file, dialect, encoding.Name == EncodingUTF8BOM)
	if err != nil {
		return Dialect{}, err
	}
	return applyDetectedQuoteRules(stats, dialect), nil
}

func DetectQuoteRulesAll(input []byte, dialect Dialect) Dialect {
	flags := quotedFieldFlags(input, dialect, 0)
	return applyDetectedQuoteRules(quoteRuleStatsFromFlags(flags, dialect), dialect)
}

type quoteRuleStats struct {
	header      []bool
	seen        map[int]int
	quoted      map[int]int
	quotedCells []Cell
}

func quoteRuleStatsFromFlags(flags [][]bool, dialect Dialect) quoteRuleStats {
	stats := newQuoteRuleStats()
	for row, record := range flags {
		stats.addRecord(row, record, dialect)
	}
	return stats
}

func newQuoteRuleStats() quoteRuleStats {
	return quoteRuleStats{
		seen:   map[int]int{},
		quoted: map[int]int{},
	}
}

func (stats *quoteRuleStats) addRecord(row int, record []bool, dialect Dialect) {
	if dialect.HasHeader && row == 0 {
		stats.header = append([]bool(nil), record...)
	} else {
		for col, wasQuoted := range record {
			stats.seen[col]++
			if wasQuoted {
				stats.quoted[col]++
			}
		}
	}
	for col, wasQuoted := range record {
		if wasQuoted && len(stats.quotedCells) < maxDetectedQuoteCells {
			stats.quotedCells = append(stats.quotedCells, Cell{Row: row, Col: col})
		}
	}
}

func applyDetectedQuoteRules(stats quoteRuleStats, dialect Dialect) Dialect {
	if len(stats.header) == 0 && len(stats.seen) == 0 && len(stats.quotedCells) == 0 {
		return dialect
	}
	coveredHeader := false
	if dialect.HasHeader && len(stats.header) > 0 && allTrue(stats.header) {
		dialect.QuoteHeader = true
		coveredHeader = true
	}

	quoteColumns := map[int]struct{}{}
	for col, total := range stats.seen {
		if total > 0 && stats.quoted[col] == total {
			dialect.QuoteColumns = append(dialect.QuoteColumns, col)
			quoteColumns[col] = struct{}{}
		}
	}
	sort.Ints(dialect.QuoteColumns)

	for _, cell := range stats.quotedCells {
		if coveredHeader && cell.Row == 0 {
			continue
		}
		if _, ok := quoteColumns[cell.Col]; ok && (!dialect.HasHeader || cell.Row > 0) {
			continue
		}
		dialect.QuoteCells = append(dialect.QuoteCells, cell)
	}
	sortCells(dialect.QuoteCells)
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

func rawQuoteScanBytes(dialect Dialect) (byte, byte, bool) {
	delimiter := dialect.Delimiter
	if delimiter == 0 {
		delimiter = ','
	}
	if delimiter > 127 {
		return 0, 0, false
	}
	quote := []byte(DialectQuoteChar(dialect))
	if len(quote) != 1 || quote[0] > 127 {
		return 0, 0, false
	}
	return byte(delimiter), quote[0], true
}

func scanQuoteRuleStats(reader io.Reader, dialect Dialect, trimUTF8BOM bool) (quoteRuleStats, error) {
	delimiter, quote, ok := rawQuoteScanBytes(dialect)
	if !ok {
		return newQuoteRuleStats(), nil
	}
	buffered := bufio.NewReaderSize(reader, 64*1024)
	if trimUTF8BOM {
		if prefix, err := buffered.Peek(3); err == nil && bytes.Equal(prefix, []byte{0xEF, 0xBB, 0xBF}) {
			_, _ = buffered.Discard(3)
		}
	}
	stats := newQuoteRuleStats()
	record := []bool{}
	row := 0
	quoted := false
	atStart := true
	inQuotes := false
	for {
		char, err := buffered.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return quoteRuleStats{}, err
		}
		if inQuotes {
			if char == quote {
				next, err := buffered.Peek(1)
				if err == nil && len(next) == 1 && next[0] == quote {
					_, _ = buffered.Discard(1)
				} else {
					inQuotes = false
				}
			}
			continue
		}
		if atStart && char == quote {
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
			stats.addRecord(row, record, dialect)
			row++
			record = []bool{}
			quoted = false
			atStart = true
			if char == '\r' {
				if next, err := buffered.Peek(1); err == nil && len(next) == 1 && next[0] == '\n' {
					_, _ = buffered.Discard(1)
				}
			}
			continue
		}
		atStart = false
	}
	if len(record) > 0 || !atStart || quoted {
		record = append(record, quoted)
		stats.addRecord(row, record, dialect)
	}
	return stats, nil
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
	for index := 0; index < len(input) && (maxRows <= 0 || len(records) < maxRows); index++ {
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
	if (maxRows <= 0 || len(records) < maxRows) && (len(record) > 0 || !atStart || quoted) {
		record = append(record, quoted)
		records = append(records, record)
	}
	return records
}
