package xlsx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type WorksheetInfo struct {
	Index          int    `json:"index"`
	Name           string `json:"name"`
	Path           string `json:"path,omitempty"`
	RelationshipID string `json:"-"`
}

type DecodeOptions struct {
	SheetName  string
	SheetIndex int
}

func DecodeFile(path string) ([][]string, error) {
	records, _, _, err := DecodeFileWithOptions(path, DecodeOptions{})
	return records, err
}

func DecodeFileWithOptions(path string, options DecodeOptions) ([][]string, []WorksheetInfo, WorksheetInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, WorksheetInfo{}, err
	}
	return DecodeWorkbookWithOptions(data, options)
}

func DecodeWorkbook(data []byte) ([][]string, error) {
	records, _, _, err := DecodeWorkbookWithOptions(data, DecodeOptions{})
	return records, err
}

func DecodeWorkbookWithOptions(data []byte, options DecodeOptions) ([][]string, []WorksheetInfo, WorksheetInfo, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, nil, WorksheetInfo{}, err
	}
	entries, err := zipEntries(reader)
	if err != nil {
		return nil, nil, WorksheetInfo{}, err
	}
	sheets, err := workbookSheets(entries)
	if err != nil {
		return nil, nil, WorksheetInfo{}, err
	}
	selected, err := selectWorksheet(sheets, options)
	if err != nil {
		return nil, nil, WorksheetInfo{}, err
	}
	sharedStrings, err := parseSharedStrings(entries["xl/sharedStrings.xml"])
	if err != nil {
		return nil, nil, WorksheetInfo{}, err
	}
	styles, err := parseStyles(entries["xl/styles.xml"])
	if err != nil {
		return nil, nil, WorksheetInfo{}, err
	}
	sheet, ok := entries[selected.Path]
	if !ok {
		return nil, nil, WorksheetInfo{}, fmt.Errorf("worksheet %q not found", selected.Path)
	}
	records, err := parseSheet(sheet, sharedStrings, styles)
	if err != nil {
		return nil, nil, WorksheetInfo{}, err
	}
	return records, sheets, selected, nil
}

func WorkbookSheets(data []byte) ([]WorksheetInfo, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	entries, err := zipEntries(reader)
	if err != nil {
		return nil, err
	}
	return workbookSheets(entries)
}

func zipEntries(reader *zip.Reader) (map[string][]byte, error) {
	entries := make(map[string][]byte, len(reader.File))
	for _, file := range reader.File {
		handle, err := file.Open()
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(handle)
		_ = handle.Close()
		if err != nil {
			return nil, err
		}
		entries[file.Name] = body
	}
	return entries, nil
}

func firstWorksheetPath(entries map[string][]byte) (string, error) {
	sheets, err := workbookSheets(entries)
	if err != nil {
		return "", err
	}
	return sheets[0].Path, nil
}

func firstSheetRelID(workbook []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(workbook))
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sheet" {
			continue
		}
		for _, attr := range start.Attr {
			if attr.Name.Local == "id" {
				return attr.Value, nil
			}
		}
		return "", errors.New("first worksheet is missing relationship id")
	}
	return "", errors.New("xlsx workbook contains no worksheets")
}

func workbookSheets(entries map[string][]byte) ([]WorksheetInfo, error) {
	workbook, ok := entries["xl/workbook.xml"]
	if !ok {
		if _, hasSheet := entries["xl/worksheets/sheet1.xml"]; hasSheet {
			return []WorksheetInfo{{Index: 1, Name: "Sheet1", Path: "xl/worksheets/sheet1.xml"}}, nil
		}
		return nil, errors.New("xlsx workbook.xml not found")
	}
	rels, err := workbookRelationships(entries["xl/_rels/workbook.xml.rels"])
	if err != nil {
		return nil, err
	}
	sheets, err := parseWorkbookSheets(workbook, rels)
	if err != nil {
		return nil, err
	}
	return sheets, nil
}

func parseWorkbookSheets(workbook []byte, rels map[string]string) ([]WorksheetInfo, error) {
	decoder := xml.NewDecoder(bytes.NewReader(workbook))
	var sheets []WorksheetInfo
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sheet" {
			continue
		}
		sheet := WorksheetInfo{Index: len(sheets) + 1}
		for _, attr := range start.Attr {
			switch attr.Name.Local {
			case "name":
				sheet.Name = attr.Value
			case "id":
				sheet.RelationshipID = attr.Value
			}
		}
		if sheet.Name == "" {
			sheet.Name = fmt.Sprintf("Sheet%d", sheet.Index)
		}
		target, ok := rels[sheet.RelationshipID]
		if !ok || target == "" {
			return nil, fmt.Errorf("worksheet relationship %q not found", sheet.RelationshipID)
		}
		sheet.Path = worksheetPath(target)
		sheets = append(sheets, sheet)
	}
	if len(sheets) == 0 {
		return nil, errors.New("xlsx workbook contains no worksheets")
	}
	return sheets, nil
}

func worksheetPath(target string) string {
	if strings.HasPrefix(target, "/") {
		return strings.TrimPrefix(path.Clean(target), "/")
	}
	return path.Clean(path.Join("xl", target))
}

func selectWorksheet(sheets []WorksheetInfo, options DecodeOptions) (WorksheetInfo, error) {
	if strings.TrimSpace(options.SheetName) != "" {
		name := strings.TrimSpace(options.SheetName)
		for _, sheet := range sheets {
			if sheet.Name == name {
				return sheet, nil
			}
		}
		for _, sheet := range sheets {
			if strings.EqualFold(sheet.Name, name) {
				return sheet, nil
			}
		}
		return WorksheetInfo{}, fmt.Errorf("worksheet %q not found", options.SheetName)
	}
	if options.SheetIndex > 0 {
		for _, sheet := range sheets {
			if sheet.Index == options.SheetIndex {
				return sheet, nil
			}
		}
		return WorksheetInfo{}, fmt.Errorf("worksheet index %d not found", options.SheetIndex)
	}
	return sheets[0], nil
}

func workbookRelationships(data []byte) (map[string]string, error) {
	if len(data) == 0 {
		return nil, errors.New("workbook relationships not found")
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	relationships := map[string]string{}
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "Relationship" {
			continue
		}
		var id, target string
		for _, attr := range start.Attr {
			switch attr.Name.Local {
			case "Id":
				id = attr.Value
			case "Target":
				target = attr.Value
			}
		}
		if id != "" {
			relationships[id] = target
		}
	}
	return relationships, nil
}

func parseSharedStrings(data []byte) ([]string, error) {
	if len(data) == 0 {
		return nil, nil
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var values []string
	var builder strings.Builder
	inSI := false
	inText := false
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch value := token.(type) {
		case xml.StartElement:
			if value.Name.Local == "si" {
				inSI = true
				builder.Reset()
			} else if inSI && value.Name.Local == "t" {
				inText = true
			}
		case xml.EndElement:
			if value.Name.Local == "t" {
				inText = false
			} else if value.Name.Local == "si" {
				inSI = false
				values = append(values, builder.String())
			}
		case xml.CharData:
			if inSI && inText {
				builder.Write([]byte(value))
			}
		}
	}
	return values, nil
}

type cellStyle struct {
	NumFmtID   int
	FormatCode string
}

func parseStyles(data []byte) ([]cellStyle, error) {
	if len(data) == 0 {
		return nil, nil
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	customFormats := map[int]string{}
	var styles []cellStyle
	inCellXfs := false
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch value := token.(type) {
		case xml.StartElement:
			switch value.Name.Local {
			case "numFmt":
				var id int
				var code string
				for _, attr := range value.Attr {
					switch attr.Name.Local {
					case "numFmtId":
						id, _ = strconv.Atoi(attr.Value)
					case "formatCode":
						code = attr.Value
					}
				}
				if id >= 0 && code != "" {
					customFormats[id] = code
				}
			case "cellXfs":
				inCellXfs = true
			case "xf":
				if !inCellXfs {
					continue
				}
				var id int
				for _, attr := range value.Attr {
					if attr.Name.Local == "numFmtId" {
						id, _ = strconv.Atoi(attr.Value)
						break
					}
				}
				styles = append(styles, cellStyle{NumFmtID: id, FormatCode: customFormats[id]})
			}
		case xml.EndElement:
			if value.Name.Local == "cellXfs" {
				inCellXfs = false
			}
		}
	}
	return styles, nil
}

func parseSheet(data []byte, sharedStrings []string, styles []cellStyle) ([][]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var records [][]string
	currentRow := -1
	nextRow := 0
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		switch start.Name.Local {
		case "row":
			currentRow = rowIndexFromStart(start, nextRow)
			records = ensureRows(records, currentRow)
			nextRow = currentRow + 1
		case "c":
			if currentRow < 0 {
				currentRow = nextRow
				records = ensureRows(records, currentRow)
				nextRow = currentRow + 1
			}
			col := columnIndexFromStart(start, len(records[currentRow]))
			value, err := parseCell(decoder, start, sharedStrings, styles)
			if err != nil {
				return nil, err
			}
			records[currentRow] = setCell(records[currentRow], col, value)
		}
	}
	return trimTrailingEmptyRows(records), nil
}

func rowIndexFromStart(start xml.StartElement, fallback int) int {
	for _, attr := range start.Attr {
		if attr.Name.Local != "r" {
			continue
		}
		row, err := strconv.Atoi(attr.Value)
		if err == nil && row > 0 {
			return row - 1
		}
	}
	return fallback
}

func columnIndexFromStart(start xml.StartElement, fallback int) int {
	for _, attr := range start.Attr {
		if attr.Name.Local != "r" {
			continue
		}
		col, ok := columnIndexFromCellRef(attr.Value)
		if ok {
			return col
		}
	}
	return fallback
}

func columnIndexFromCellRef(ref string) (int, bool) {
	col := 0
	seen := false
	for _, r := range ref {
		if !unicode.IsLetter(r) {
			break
		}
		seen = true
		upper := unicode.ToUpper(r)
		if upper < 'A' || upper > 'Z' {
			return 0, false
		}
		col = col*26 + int(upper-'A'+1)
	}
	if !seen {
		return 0, false
	}
	return col - 1, true
}

func parseCell(decoder *xml.Decoder, start xml.StartElement, sharedStrings []string, styles []cellStyle) (string, error) {
	cellType := ""
	styleIndex := -1
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "t":
			cellType = attr.Value
		case "s":
			styleIndex, _ = strconv.Atoi(attr.Value)
		}
	}
	var raw strings.Builder
	var inline strings.Builder
	inValue := false
	inInlineText := false
	for {
		token, err := decoder.Token()
		if err != nil {
			return "", err
		}
		switch value := token.(type) {
		case xml.StartElement:
			if value.Name.Local == "v" {
				inValue = true
			} else if value.Name.Local == "t" {
				inInlineText = true
			}
		case xml.EndElement:
			if value.Name.Local == "v" {
				inValue = false
			} else if value.Name.Local == "t" {
				inInlineText = false
			} else if value.Name.Local == "c" {
				return cellValue(cellType, raw.String(), inline.String(), sharedStrings, styles, styleIndex), nil
			}
		case xml.CharData:
			if inValue {
				raw.Write([]byte(value))
			}
			if inInlineText {
				inline.Write([]byte(value))
			}
		}
	}
}

func cellValue(cellType string, raw string, inline string, sharedStrings []string, styles []cellStyle, styleIndex int) string {
	switch cellType {
	case "s":
		index, err := strconv.Atoi(strings.TrimSpace(raw))
		if err == nil && index >= 0 && index < len(sharedStrings) {
			return sharedStrings[index]
		}
		return ""
	case "inlineStr":
		return inline
	case "str":
		if raw != "" {
			return raw
		}
		return inline
	case "b":
		if strings.TrimSpace(raw) == "1" {
			return "TRUE"
		}
		return "FALSE"
	default:
		if dateLikeStyle(styles, styleIndex) && strings.TrimSpace(raw) != "" {
			if formatted, ok := formatExcelDateSerial(strings.TrimSpace(raw), styles[styleIndex]); ok {
				return formatted
			}
		}
		return raw
	}
}

func dateLikeStyle(styles []cellStyle, index int) bool {
	return index >= 0 && index < len(styles) && dateStyleKind(styles[index]) != ""
}

func dateStyleKind(style cellStyle) string {
	switch style.NumFmtID {
	case 14, 15, 16, 17:
		return "date"
	case 18, 19, 20, 21, 45, 46, 47:
		return "time"
	case 22:
		return "datetime"
	}
	code := normalizedFormatCode(style.FormatCode)
	if code == "" {
		return ""
	}
	hasDate := strings.ContainsAny(code, "yd")
	hasTime := strings.ContainsAny(code, "hs")
	if hasDate && hasTime {
		return "datetime"
	}
	if hasDate {
		return "date"
	}
	if hasTime {
		return "time"
	}
	return ""
}

func normalizedFormatCode(code string) string {
	var builder strings.Builder
	inQuote := false
	inBracket := false
	escaped := false
	for _, r := range strings.ToLower(code) {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '"' {
			inQuote = !inQuote
			continue
		}
		if r == '[' {
			inBracket = true
			continue
		}
		if r == ']' {
			inBracket = false
			continue
		}
		if inQuote || inBracket {
			continue
		}
		if strings.ContainsRune("ymdhs", r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func formatExcelDateSerial(raw string, style cellStyle) (string, bool) {
	serial, err := strconv.ParseFloat(raw, 64)
	if err != nil || serial < 0 {
		return "", false
	}
	wholeDays := int(math.Floor(serial))
	fraction := serial - float64(wholeDays)
	if wholeDays >= 60 {
		wholeDays--
	}
	base := time.Date(1899, 12, 31, 0, 0, 0, 0, time.UTC)
	seconds := int(math.Round(fraction * 86400))
	if seconds >= 86400 {
		wholeDays++
		seconds -= 86400
	}
	value := base.AddDate(0, 0, wholeDays).Add(time.Duration(seconds) * time.Second)
	switch dateStyleKind(style) {
	case "time":
		return value.Format("15:04:05"), true
	case "datetime":
		return value.Format("2006-01-02 15:04:05"), true
	default:
		return value.Format("2006-01-02"), true
	}
}

func ensureRows(records [][]string, row int) [][]string {
	for len(records) <= row {
		records = append(records, nil)
	}
	return records
}

func setCell(record []string, col int, value string) []string {
	for len(record) <= col {
		record = append(record, "")
	}
	record[col] = value
	return record
}

func trimTrailingEmptyRows(records [][]string) [][]string {
	end := len(records)
	for end > 0 && len(records[end-1]) == 0 {
		end--
	}
	return records[:end]
}
