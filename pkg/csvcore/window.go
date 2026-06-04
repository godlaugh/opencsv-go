package csvcore

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type RecordWindow struct {
	Start      int        `json:"start"`
	Rows       [][]string `json:"rows"`
	TotalRows  int        `json:"totalRows"`
	SourceRows []int      `json:"sourceRows,omitempty"`
}

func ReadRecordWindowWithIndex(path string, metadata FileMetadata, index RowIndex, start int, count int) (RecordWindow, error) {
	if canReadIndexedWindow(metadata, index, start) {
		return readRecordWindowFromOffset(path, metadata, index, start, count)
	}
	return ReadRecordWindow(path, metadata, start, count)
}

func ReadRecordWindow(path string, metadata FileMetadata, start int, count int) (RecordWindow, error) {
	if start < 0 {
		start = 0
	}
	if count <= 0 {
		count = 50
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return RecordWindow{}, err
	}
	utf8Bytes, err := DecodeToUTF8(raw, metadata.Encoding)
	if err != nil {
		return RecordWindow{}, err
	}

	reader := csv.NewReader(bytes.NewReader(normalizeQuoteEscapeInput(utf8Bytes, metadata.Dialect)))
	reader.Comma = metadata.Dialect.Delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	window := RecordWindow{
		Start:      start,
		Rows:       make([][]string, 0, count),
		SourceRows: make([]int, 0, count),
	}
	rowIndex := 0
	for len(window.Rows) < count {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return window, err
		}
		if rowIndex == 0 && len(record) > 0 && metadata.Encoding.Name == EncodingUTF8BOM {
			record[0] = strings.TrimPrefix(record[0], "\ufeff")
		}
		if rowIndex >= start {
			window.Rows = append(window.Rows, record)
			window.SourceRows = append(window.SourceRows, rowIndex)
		}
		rowIndex++
	}
	window.TotalRows = rowIndex

	return window, nil
}

func canReadIndexedWindow(metadata FileMetadata, index RowIndex, start int) bool {
	if start < 0 {
		start = 0
	}
	if !index.Complete || start < 0 {
		return false
	}
	switch metadata.Encoding.Name {
	case "", EncodingUTF8, EncodingUTF8BOM:
	default:
		return false
	}
	if DialectQuoteChar(metadata.Dialect) != `"` || DialectQuoteEscape(metadata.Dialect) == QuoteEscapeBackslash {
		return false
	}
	return start >= index.Rows || start < len(index.Offsets)
}

func readRecordWindowFromOffset(path string, metadata FileMetadata, index RowIndex, start int, count int) (RecordWindow, error) {
	if start < 0 {
		start = 0
	}
	if count <= 0 {
		count = 50
	}
	window := RecordWindow{
		Start:      start,
		TotalRows:  index.Rows,
		Rows:       make([][]string, 0, count),
		SourceRows: make([]int, 0, count),
	}
	if start >= index.Rows || start >= len(index.Offsets) {
		return window, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return RecordWindow{}, err
	}
	defer file.Close()
	if _, err := file.Seek(index.Offsets[start], io.SeekStart); err != nil {
		return RecordWindow{}, err
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = metadata.Dialect.Delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	for len(window.Rows) < count {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return window, err
		}
		sourceRow := start + len(window.Rows)
		if sourceRow == 0 && len(record) > 0 && metadata.Encoding.Name == EncodingUTF8BOM {
			record[0] = strings.TrimPrefix(record[0], "\ufeff")
		}
		window.Rows = append(window.Rows, record)
		window.SourceRows = append(window.SourceRows, sourceRow)
	}
	return window, nil
}
