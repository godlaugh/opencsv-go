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
)

func NormalizeQuoteMode(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", QuoteModeMinimal:
		return QuoteModeMinimal, nil
	case QuoteModeAlways:
		return QuoteModeAlways, nil
	case QuoteModeNever, "none":
		return QuoteModeNever, nil
	default:
		return "", fmt.Errorf("quote mode must be minimal, always, or never")
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
	switch DialectQuoteMode(dialect) {
	case QuoteModeAlways:
		return true
	case QuoteModeNever:
		return false
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
	for _, quoteCol := range dialect.QuoteColumns {
		if quoteCol == col {
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
