package csvcore

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

const DefaultSampleBytes = 128 * 1024

type PreviewResult struct {
	Path        string     `json:"path"`
	SizeBytes   int64      `json:"sizeBytes"`
	Dialect     Dialect    `json:"dialect"`
	Records     [][]string `json:"records"`
	RowsScanned int        `json:"rowsScanned"`
}

func OpenPreview(path string, maxRows int) (PreviewResult, error) {
	if maxRows <= 0 {
		maxRows = 50
	}

	file, err := os.Open(path)
	if err != nil {
		return PreviewResult{}, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return PreviewResult{}, err
	}

	sample := make([]byte, min(DefaultSampleBytes, int(info.Size())))
	n, readErr := io.ReadFull(file, sample)
	if readErr != nil && !errors.Is(readErr, io.ErrUnexpectedEOF) && !errors.Is(readErr, io.EOF) {
		return PreviewResult{}, readErr
	}
	sample = sample[:n]

	dialect := Sniff(sample)
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return PreviewResult{}, err
	}

	reader := csv.NewReader(file)
	reader.Comma = dialect.Delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	records := make([][]string, 0, maxRows)
	for len(records) < maxRows {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return PreviewResult{}, err
		}
		if dialect.UTF8BOM && len(records) == 0 && len(record) > 0 {
			record[0] = strings.TrimPrefix(record[0], "\ufeff")
		}
		records = append(records, record)
	}

	return PreviewResult{
		Path:        path,
		SizeBytes:   info.Size(),
		Dialect:     dialect,
		Records:     records,
		RowsScanned: len(records),
	}, nil
}
