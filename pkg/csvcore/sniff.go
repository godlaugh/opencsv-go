package csvcore

import (
	"bytes"
	"encoding/csv"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Dialect struct {
	Delimiter        rune   `json:"delimiter"`
	DelimiterName    string `json:"delimiterName"`
	HasHeader        bool   `json:"hasHeader"`
	UTF8BOM          bool   `json:"utf8Bom"`
	EstimatedColumns int    `json:"estimatedColumns"`
}

func Sniff(sample []byte) Dialect {
	dialect := Dialect{
		Delimiter:     ',',
		DelimiterName: "comma",
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

	first, err := reader.Read()
	if err != nil || len(first) == 0 {
		return false
	}

	second, err := reader.Read()
	if err != nil || len(second) == 0 {
		return false
	}

	textyFirst := 0
	numericSecond := 0
	columns := min(len(first), len(second))
	for i := 0; i < columns; i++ {
		if !looksNumeric(first[i]) && strings.TrimSpace(first[i]) != "" {
			textyFirst++
		}
		if looksNumeric(second[i]) {
			numericSecond++
		}
	}

	return textyFirst >= max(1, columns/2) && numericSecond > 0
}

func looksNumeric(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}
