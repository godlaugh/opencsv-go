package csvcore

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type FileMetadata struct {
	Path       string         `json:"path"`
	SizeBytes  int64          `json:"sizeBytes"`
	Encoding   EncodingInfo   `json:"encoding"`
	LineEnding LineEndingInfo `json:"lineEnding"`
	Dialect    Dialect        `json:"dialect"`
	Workbook   *WorkbookInfo  `json:"workbook,omitempty"`
}

type WorkbookInfo struct {
	Sheets             []WorkbookSheetInfo `json:"sheets"`
	SelectedSheetIndex int                 `json:"selectedSheetIndex"`
	SelectedSheetName  string              `json:"selectedSheetName"`
	SelectedSheetPath  string              `json:"selectedSheetPath,omitempty"`
}

type WorkbookSheetInfo struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	Path  string `json:"path,omitempty"`
}

type Document struct {
	Metadata    FileMetadata `json:"metadata"`
	Preview     [][]string   `json:"preview"`
	RowsScanned int          `json:"rowsScanned"`
	RowIndex    RowIndex     `json:"rowIndex"`
}

var ErrStreamingUnsupported = errors.New("streaming records unsupported for this file metadata")

type OpenOptions struct {
	PreviewRows int `json:"previewRows"`
	// IndexRows controls synchronous indexing during open. Use 0 for no row
	// indexing, a positive value for a small first-page index, and -1 for full
	// indexing.
	IndexRows     int `json:"indexRows"`
	SampleBytes   int `json:"sampleBytes"`
	Delimiter     *rune
	DelimiterName string
	QuoteChar     *string
	QuoteMode     string
	QuoteEscape   string
	HasHeader     *bool
	HeaderColumns *int
	Encoding      *EncodingInfo
	LineEnding    *LineEndingInfo
	SheetName     string
	SheetIndex    int
}

func OpenDocument(path string, options OpenOptions) (Document, error) {
	if options.PreviewRows <= 0 {
		options.PreviewRows = 50
	}
	if options.SampleBytes <= 0 {
		options.SampleBytes = DefaultSampleBytes
	}

	file, err := os.Open(path)
	if err != nil {
		return Document{}, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return Document{}, err
	}

	sample := make([]byte, min(options.SampleBytes, int(info.Size())))
	n, readErr := io.ReadFull(file, sample)
	if readErr != nil && !errors.Is(readErr, io.ErrUnexpectedEOF) && !errors.Is(readErr, io.EOF) {
		return Document{}, readErr
	}
	sample = sample[:n]

	encoding := DetectEncoding(sample)
	if options.Encoding != nil {
		encoding = *options.Encoding
	}
	utf8Sample, err := DecodeToUTF8(sample, encoding)
	if err != nil {
		return Document{}, err
	}
	dialect := Sniff(utf8Sample)
	if options.Delimiter != nil {
		dialect.Delimiter = *options.Delimiter
		dialect.DelimiterName = options.DelimiterName
		if dialect.DelimiterName == "" {
			dialect.DelimiterName = string(*options.Delimiter)
		}
	}
	if options.HasHeader != nil {
		dialect.HasHeader = *options.HasHeader
	}
	if options.HeaderColumns != nil {
		dialect.HeaderColumns = max(0, *options.HeaderColumns)
	}
	if options.QuoteChar != nil {
		quoteChar, err := NormalizeQuoteChar(*options.QuoteChar)
		if err != nil {
			return Document{}, err
		}
		dialect.QuoteChar = quoteChar
	}
	if strings.TrimSpace(options.QuoteMode) != "" {
		quoteMode, err := NormalizeQuoteMode(options.QuoteMode)
		if err != nil {
			return Document{}, err
		}
		dialect.QuoteMode = quoteMode
	}
	if strings.TrimSpace(options.QuoteEscape) != "" {
		quoteEscape, err := NormalizeQuoteEscape(options.QuoteEscape)
		if err != nil {
			return Document{}, err
		}
		dialect.QuoteEscape = quoteEscape
	}
	if options.IndexRows < 0 {
		detectedDialect, err := DetectQuoteRulesForFile(path, dialect, encoding)
		if err != nil {
			return Document{}, err
		}
		dialect = detectedDialect
	} else {
		dialect = DetectQuoteRules(utf8Sample, dialect)
	}

	metadata := FileMetadata{
		Path:       path,
		SizeBytes:  info.Size(),
		Encoding:   encoding,
		LineEnding: DetectLineEnding(utf8Sample),
		Dialect:    dialect,
	}
	if options.LineEnding != nil {
		metadata.LineEnding = *options.LineEnding
	}
	metadata.Dialect.UTF8BOM = encoding.Name == EncodingUTF8BOM

	preview, err := previewRecordsFromSample(utf8Sample, metadata.Dialect, options.PreviewRows)
	if err != nil {
		return Document{}, err
	}

	index := RowIndex{Path: path, SizeBytes: info.Size()}
	if options.IndexRows != 0 {
		indexRows := options.IndexRows
		if indexRows < 0 {
			indexRows = 0
		}
		index, err = BuildRowIndex(path, metadata.Dialect, indexRows)
		if err != nil {
			return Document{}, err
		}
	}

	return Document{
		Metadata:    metadata,
		Preview:     preview,
		RowsScanned: len(preview),
		RowIndex:    index,
	}, nil
}

func ReadAllRecords(path string, metadata FileMetadata) ([][]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	utf8Bytes, err := DecodeToUTF8(raw, metadata.Encoding)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(normalizeQuoteEscapeInput(utf8Bytes, metadata.Dialect)))
	reader.Comma = metadata.Dialect.Delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	var records [][]string
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(records) == 0 && len(record) > 0 && metadata.Encoding.Name == EncodingUTF8BOM {
			record[0] = strings.TrimPrefix(record[0], "\ufeff")
		}
		records = append(records, record)
	}
	return records, nil
}

func CanStreamRecords(metadata FileMetadata) bool {
	switch metadata.Encoding.Name {
	case "", EncodingUTF8, EncodingUTF8BOM:
	default:
		return false
	}
	if DialectQuoteChar(metadata.Dialect) != `"` || DialectQuoteEscape(metadata.Dialect) == QuoteEscapeBackslash {
		return false
	}
	return true
}

func ForEachRecordContext(ctx context.Context, path string, metadata FileMetadata, handle func(row int, record []string) error) (int, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !CanStreamRecords(metadata) {
		return 0, ErrStreamingUnsupported
	}
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = metadata.Dialect.Delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	rowIndex := 0
	for {
		if err := ctx.Err(); err != nil {
			return rowIndex, err
		}
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return rowIndex, err
		}
		if rowIndex == 0 && len(record) > 0 && metadata.Encoding.Name == EncodingUTF8BOM {
			record[0] = strings.TrimPrefix(record[0], "\ufeff")
		}
		if err := handle(rowIndex, record); err != nil {
			return rowIndex, err
		}
		rowIndex++
	}
	return rowIndex, nil
}

func ParseRecordsFromUTF8(input []byte, dialect Dialect) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(normalizeQuoteEscapeInput(input, dialect)))
	reader.Comma = dialect.Delimiter
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	var records [][]string
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(records) == 0 && len(record) > 0 && dialect.UTF8BOM {
			record[0] = strings.TrimPrefix(record[0], "\ufeff")
		}
		records = append(records, record)
	}
	return records, nil
}

func previewRecordsFromSample(sample []byte, dialect Dialect, maxRows int) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(normalizeQuoteEscapeInput(sample, dialect)))
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
			if len(records) > 0 {
				break
			}
			return nil, err
		}
		if dialect.UTF8BOM && len(records) == 0 && len(record) > 0 {
			record[0] = strings.TrimPrefix(record[0], "\ufeff")
		}
		records = append(records, record)
	}
	return records, nil
}

func normalizeQuoteEscapeInput(input []byte, dialect Dialect) []byte {
	quote := []byte(DialectQuoteChar(dialect))
	if len(quote) != 1 {
		return input
	}
	quoteByte := quote[0]
	normalizeQuoteChar := quoteByte != '"'
	normalizeBackslashEscape := DialectQuoteEscape(dialect) == QuoteEscapeBackslash
	if !normalizeQuoteChar && !normalizeBackslashEscape {
		return input
	}

	output := make([]byte, 0, len(input))
	delimiterBytes := []byte(",")
	if dialect.Delimiter > 0 {
		delimiterBytes = []byte(string(dialect.Delimiter))
	}
	inQuotes := false
	atFieldStart := true
	for index := 0; index < len(input); index++ {
		char := input[index]
		if inQuotes {
			if normalizeBackslashEscape && char == '\\' && index+1 < len(input) {
				next := input[index+1]
				if next == quoteByte {
					appendEscapedQuoteForReader(&output, next, normalizeQuoteChar)
					index++
					continue
				}
				if next == '\\' {
					output = append(output, '\\')
					index++
					continue
				}
			}
			if char == quoteByte {
				if !normalizeBackslashEscape && index+1 < len(input) && input[index+1] == quoteByte {
					appendNormalizedQuotedByte(&output, char, normalizeQuoteChar)
					index++
					continue
				}
				inQuotes = false
				if normalizeQuoteChar {
					output = append(output, '"')
				} else {
					output = append(output, char)
				}
				continue
			}
			appendNormalizedQuotedByte(&output, char, normalizeQuoteChar)
			continue
		}
		if atFieldStart && char == quoteByte {
			inQuotes = true
			atFieldStart = false
			if normalizeQuoteChar {
				output = append(output, '"')
			} else {
				output = append(output, char)
			}
			continue
		}
		if len(delimiterBytes) > 0 && bytes.HasPrefix(input[index:], delimiterBytes) {
			output = append(output, delimiterBytes...)
			index += len(delimiterBytes) - 1
			atFieldStart = true
			continue
		}
		output = append(output, char)
		atFieldStart = char == '\n' || char == '\r'
	}
	return output
}

func appendNormalizedQuotedByte(output *[]byte, char byte, normalizeQuoteChar bool) {
	if normalizeQuoteChar && char == '"' {
		*output = append(*output, '"', '"')
		return
	}
	*output = append(*output, char)
}

func appendEscapedQuoteForReader(output *[]byte, quote byte, normalizeQuoteChar bool) {
	if normalizeQuoteChar {
		*output = append(*output, quote)
		return
	}
	*output = append(*output, quote, quote)
}
