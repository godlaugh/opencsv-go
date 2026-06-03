package csvcore

import (
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
