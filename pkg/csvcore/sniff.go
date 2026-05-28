package csvcore

import (
	"bytes"
	"encoding/csv"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Dialect struct {
	Delimiter        rune   `json:"delimiter"`
	DelimiterName    string `json:"delimiterName"`
	QuoteChar        string `json:"quoteChar"`
	QuoteMode        string `json:"quoteMode"`
	QuoteHeader      bool   `json:"quoteHeader,omitempty"`
	QuoteColumns     []int  `json:"quoteColumns,omitempty"`
	QuoteCells       []Cell `json:"quoteCells,omitempty"`
	HasHeader        bool   `json:"hasHeader"`
	UTF8BOM          bool   `json:"utf8Bom"`
	EstimatedColumns int    `json:"estimatedColumns"`
}

type Cell struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func Sniff(sample []byte) Dialect {
	dialect := Dialect{
		Delimiter:     ',',
		DelimiterName: "comma",
		QuoteChar:     `"`,
		QuoteMode:     QuoteModeMinimal,
	}

	if bytes.HasPrefix(sample, []byte{0xEF, 0xBB, 0xBF}) {
		dialect.UTF8BOM = true
		sample = sample[3:]
	}

	if !utf8.Valid(sample) {
		// M0 only records the signal. Encoding conversion belongs in the file I/O
		// layer once we add Shift-JIS, GBK, and Windows-1252 support.
		sample = bytes.ToValidUTF8(sample, []byte("?"))
	}

	best := dialectCandidate{delimiter: ',', name: "comma", score: math.Inf(-1)}
	for _, candidate := range []dialectCandidate{
		{delimiter: ',', name: "comma"},
		{delimiter: '\t', name: "tab"},
		{delimiter: ';', name: "semicolon"},
		{delimiter: '|', name: "pipe"},
	} {
		score, columns, rows := scoreDelimiter(sample, candidate.delimiter)
		candidate.score = score
		candidate.columns = columns
		candidate.rows = rows
		if candidate.betterThan(best) {
			best = candidate
		}
	}

	dialect.Delimiter = best.delimiter
	dialect.DelimiterName = best.name
	dialect.EstimatedColumns = best.columns
	dialect.HasHeader = inferHeader(sample, best.delimiter)
	return dialect
}

type dialectCandidate struct {
	delimiter rune
	name      string
	score     float64
	columns   int
	rows      int
}

func (candidate dialectCandidate) betterThan(other dialectCandidate) bool {
	if candidate.score != other.score {
		return candidate.score > other.score
	}
	if candidate.columns != other.columns {
		return candidate.columns > other.columns
	}
	return candidate.delimiter == ','
}

func scoreDelimiter(sample []byte, delimiter rune) (float64, int, int) {
	reader := csv.NewReader(bytes.NewReader(sample))
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	var widths []int
	for len(widths) < 80 {
		record, err := reader.Read()
		if err != nil {
			break
		}
		if len(record) == 1 && strings.TrimSpace(record[0]) == "" {
			continue
		}
		widths = append(widths, len(record))
	}

	if len(widths) == 0 {
		return math.Inf(-1), 0, 0
	}

	sum := 0
	minWidth := widths[0]
	for _, width := range widths {
		sum += width
		if width < minWidth {
			minWidth = width
		}
	}
	mean := float64(sum) / float64(len(widths))

	var variance float64
	for _, width := range widths {
		delta := float64(width) - mean
		variance += delta * delta
	}
	variance /= float64(len(widths))

	if mean <= 1 {
		return float64(len(widths)) - variance - 1000, int(mean), len(widths)
	}

	score := float64(len(widths))*mean - variance*8 + float64(minWidth)
	return score, int(math.Round(mean)), len(widths)
}

func inferHeader(sample []byte, delimiter rune) bool {
	reader := csv.NewReader(bytes.NewReader(sample))
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records := make([][]string, 0, 8)
	for len(records) < cap(records) {
		record, err := reader.Read()
		if err != nil {
			break
		}
		if len(record) == 1 && strings.TrimSpace(record[0]) == "" {
			continue
		}
		records = append(records, record)
	}
	if len(records) < 2 || len(records[0]) == 0 || len(records[1]) == 0 {
		return false
	}

	first := records[0]
	columns := len(first)
	if columns <= 1 {
		return false
	}
	for _, record := range records[1:] {
		if len(record) < columns {
			columns = len(record)
		}
	}
	if columns <= 1 {
		return false
	}

	nameLikeFirst := 0
	semanticNameFirst := 0
	typeChanges := 0
	uniqueHeaderNames := map[string]bool{}
	for col := 0; col < columns; col++ {
		header := strings.TrimSpace(first[col])
		if header == "" {
			return false
		}
		lowerHeader := strings.ToLower(header)
		if uniqueHeaderNames[lowerHeader] {
			return false
		}
		uniqueHeaderNames[lowerHeader] = true

		if looksHeaderName(header) {
			nameLikeFirst++
		}
		if looksSemanticHeaderName(header) {
			semanticNameFirst++
		}
		for _, record := range records[1:] {
			if col >= len(record) {
				continue
			}
			if fieldKind(header) != fieldKind(record[col]) {
				typeChanges++
				break
			}
		}
	}

	strongNames := nameLikeFirst >= max(1, (columns+1)/2)
	enoughTypeSignal := typeChanges >= max(1, columns/4)
	enoughSemanticSignal := semanticNameFirst >= max(2, (columns+1)/2)
	return strongNames && (enoughTypeSignal || enoughSemanticSignal)
}

func looksNumeric(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func looksHeaderName(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || looksNumeric(value) {
		return false
	}
	hasLetter := false
	for _, r := range value {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r), r == '_', r == '-', r == ' ', r == '.', r == '/', r == '#':
		default:
			return false
		}
	}
	return hasLetter
}

func fieldKind(value string) string {
	value = strings.TrimSpace(value)
	switch {
	case value == "":
		return "blank"
	case looksNumeric(value):
		return "number"
	case looksBool(value):
		return "bool"
	case looksISODate(value):
		return "date"
	default:
		return "text"
	}
}

func looksBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "false", "yes", "no":
		return true
	default:
		return false
	}
}

func looksSemanticHeaderName(value string) bool {
	normalized := normalizeHeaderName(value)
	if _, ok := commonHeaderNames[normalized]; ok {
		return true
	}
	return hasNumericHeaderSuffix(normalized, "field") || hasNumericHeaderSuffix(normalized, "column") || hasNumericHeaderSuffix(normalized, "col")
}

func normalizeHeaderName(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func hasNumericHeaderSuffix(value, prefix string) bool {
	if !strings.HasPrefix(value, prefix) || len(value) == len(prefix) {
		return false
	}
	for _, r := range value[len(prefix):] {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func looksISODate(value string) bool {
	value = strings.TrimSpace(value)
	if len(value) < len("2006-01-02") {
		return false
	}
	return isDigitAt(value, 0) &&
		isDigitAt(value, 1) &&
		isDigitAt(value, 2) &&
		isDigitAt(value, 3) &&
		value[4] == '-' &&
		isDigitAt(value, 5) &&
		isDigitAt(value, 6) &&
		value[7] == '-' &&
		isDigitAt(value, 8) &&
		isDigitAt(value, 9)
}

func isDigitAt(value string, index int) bool {
	return index < len(value) && value[index] >= '0' && value[index] <= '9'
}

var commonHeaderNames = map[string]struct{}{
	"address":     {},
	"amount":      {},
	"category":    {},
	"city":        {},
	"code":        {},
	"count":       {},
	"country":     {},
	"createdat":   {},
	"date":        {},
	"description": {},
	"email":       {},
	"firstname":   {},
	"fullname":    {},
	"id":          {},
	"key":         {},
	"label":       {},
	"lastname":    {},
	"name":        {},
	"note":        {},
	"notes":       {},
	"phone":       {},
	"price":       {},
	"qty":         {},
	"quantity":    {},
	"region":      {},
	"score":       {},
	"segment":     {},
	"state":       {},
	"status":      {},
	"team":        {},
	"time":        {},
	"timestamp":   {},
	"total":       {},
	"type":        {},
	"updatedat":   {},
	"valid":       {},
	"value":       {},
}
