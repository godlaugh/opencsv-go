package csvcore

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

type RowIndex struct {
	Path      string  `json:"path"`
	SizeBytes int64   `json:"sizeBytes"`
	Offsets   []int64 `json:"offsets"`
	Rows      int     `json:"rows"`
	Complete  bool    `json:"complete"`
}

func BuildRowIndex(path string, dialect Dialect, maxRows int) (RowIndex, error) {
	file, err := os.Open(path)
	if err != nil {
		return RowIndex{}, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return RowIndex{}, err
	}

	reader := csv.NewReader(file)
	reader.Comma = dialect.Delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	index := RowIndex{
		Path:      path,
		SizeBytes: info.Size(),
		Offsets:   make([]int64, 0, 4096),
	}

	nextOffset := int64(0)
	for maxRows <= 0 || index.Rows < maxRows {
		index.Offsets = append(index.Offsets, nextOffset)
		_, err := reader.Read()
		if errors.Is(err, io.EOF) {
			index.Offsets = index.Offsets[:len(index.Offsets)-1]
			index.Complete = true
			break
		}
		if err != nil {
			index.Offsets = index.Offsets[:len(index.Offsets)-1]
			return index, err
		}
		index.Rows++
		nextOffset = reader.InputOffset()
	}

	return index, nil
}
