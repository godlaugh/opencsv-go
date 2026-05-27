package xlsx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"
)

func DecodeFile(path string) ([][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DecodeWorkbook(data)
}

func DecodeWorkbook(data []byte) ([][]string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	entries, err := zipEntries(reader)
	if err != nil {
		return nil, err
	}
	sheetPath, err := firstWorksheetPath(entries)
	if err != nil {
		return nil, err
	}
	sharedStrings, err := parseSharedStrings(entries["xl/sharedStrings.xml"])
	if err != nil {
		return nil, err
	}
	sheet, ok := entries[sheetPath]
	if !ok {
		return nil, fmt.Errorf("worksheet %q not found", sheetPath)
	}
	return parseSheet(sheet, sharedStrings)
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
	workbook, ok := entries["xl/workbook.xml"]
	if !ok {
		if _, hasSheet := entries["xl/worksheets/sheet1.xml"]; hasSheet {
			return "xl/worksheets/sheet1.xml", nil
		}
		return "", errors.New("xlsx workbook.xml not found")
	}
	firstRelID, err := firstSheetRelID(workbook)
	if err != nil {
		return "", err
	}
	rels, err := workbookRelationships(entries["xl/_rels/workbook.xml.rels"])
	if err != nil {
		return "", err
	}
	target, ok := rels[firstRelID]
	if !ok || target == "" {
		return "", fmt.Errorf("worksheet relationship %q not found", firstRelID)
	}
	if strings.HasPrefix(target, "/") {
		return strings.TrimPrefix(path.Clean(target), "/"), nil
	}
	return path.Clean(path.Join("xl", target)), nil
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

func parseSheet(data []byte, sharedStrings []string) ([][]string, error) {
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
			value, err := parseCell(decoder, start, sharedStrings)
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

func parseCell(decoder *xml.Decoder, start xml.StartElement, sharedStrings []string) (string, error) {
	cellType := ""
	for _, attr := range start.Attr {
		if attr.Name.Local == "t" {
			cellType = attr.Value
			break
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
				return cellValue(cellType, raw.String(), inline.String(), sharedStrings), nil
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

func cellValue(cellType string, raw string, inline string, sharedStrings []string) string {
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
		return raw
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
