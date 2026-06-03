package csvcore

import (
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

	_ = dialect

	index := RowIndex{
		Path:      path,
		SizeBytes: info.Size(),
		Offsets:   make([]int64, 0, 4096),
	}

	index.Offsets = append(index.Offsets, 0)
	scanner := newRowBoundaryScanner(dialect)
	buffer := make([]byte, 256*1024)
	fileOffset := int64(0)
	lastBoundary := int64(-1)
	for maxRows <= 0 || index.Rows < maxRows {
		n, err := file.Read(buffer)
		if n > 0 {
			boundaries := scanner.scan(buffer[:n], fileOffset)
			for _, offset := range boundaries {
				lastBoundary = offset
				index.Rows++
				if maxRows > 0 && index.Rows >= maxRows {
					return index, nil
				}
				if offset < info.Size() {
					index.Offsets = append(index.Offsets, offset)
				}
			}
			fileOffset += int64(n)
		}
		if err == io.EOF {
			index.Complete = true
			break
		}
		if err != nil {
			return index, err
		}
	}
	if info.Size() > 0 && lastBoundary != info.Size() {
		index.Rows++
	}

	return index, nil
}

type rowBoundaryScanner struct {
	inQuotes     bool
	pendingQuote bool
	prevCR       bool
	escaped      bool
	quote        byte
	backslash    bool
}

func newRowBoundaryScanner(dialect Dialect) rowBoundaryScanner {
	quoteChar := DialectQuoteChar(dialect)
	scanner := rowBoundaryScanner{}
	if quoteChar == "" {
		scanner.quote = 0
	} else if len(quoteChar) == 1 {
		scanner.quote = quoteChar[0]
	} else {
		scanner.quote = '"'
	}
	scanner.backslash = DialectQuoteEscape(dialect) == QuoteEscapeBackslash
	return scanner
}

func (scanner *rowBoundaryScanner) scan(input []byte, baseOffset int64) []int64 {
	var boundaries []int64
	for i, b := range input {
		offset := baseOffset + int64(i)

		if scanner.escaped {
			scanner.escaped = false
			scanner.prevCR = false
			continue
		}

		if scanner.pendingQuote {
			if b == scanner.quote {
				scanner.pendingQuote = false
				continue
			}
			scanner.inQuotes = false
			scanner.pendingQuote = false
		}

		if scanner.inQuotes && scanner.backslash && b == '\\' {
			scanner.escaped = true
			scanner.prevCR = false
			continue
		}

		if scanner.quote != 0 && b == scanner.quote {
			if scanner.inQuotes {
				scanner.pendingQuote = true
			} else {
				scanner.inQuotes = true
			}
			scanner.prevCR = false
			continue
		}

		if scanner.inQuotes {
			scanner.prevCR = false
			continue
		}

		if b == '\n' {
			boundaries = append(boundaries, offset+1)
			scanner.prevCR = false
			continue
		}
		if scanner.prevCR {
			boundaries = append(boundaries, offset)
			scanner.prevCR = false
		}
		if b == '\r' {
			scanner.prevCR = true
		}
	}
	return boundaries
}
