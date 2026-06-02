package xlsx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)

const contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

func ContentType() string {
	return contentType
}

type WorkbookSheet struct {
	Name    string
	Records [][]string
}

func EncodeWorkbook(records [][]string) ([]byte, error) {
	return EncodeWorkbookSheets([]WorkbookSheet{{Name: "Sheet1", Records: records}})
}

func EncodeWorkbookSheets(sheets []WorkbookSheet) ([]byte, error) {
	if len(sheets) == 0 {
		sheets = []WorkbookSheet{{Name: "Sheet1"}}
	}
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	entries := map[string]string{
		"[Content_Types].xml":        contentTypesXMLForSheets(len(sheets)),
		"_rels/.rels":                rootRelsXML,
		"docProps/app.xml":           appXML,
		"docProps/core.xml":          coreXML(time.Now().UTC()),
		"xl/workbook.xml":            workbookXMLForSheets(sheets),
		"xl/_rels/workbook.xml.rels": workbookRelsXMLForSheets(len(sheets)),
	}
	order := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"docProps/app.xml",
		"docProps/core.xml",
		"xl/workbook.xml",
		"xl/_rels/workbook.xml.rels",
	}
	for index, sheet := range sheets {
		name := fmt.Sprintf("xl/worksheets/sheet%d.xml", index+1)
		entries[name] = sheetXML(sheet.Records)
		order = append(order, name)
	}
	for _, name := range order {
		if err := writeEntry(writer, name, entries[name]); err != nil {
			_ = writer.Close()
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func UpdateWorkbookSheet(data []byte, sheetPath string, records [][]string) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	cleanSheetPath := strings.TrimPrefix(sheetPath, "/")
	replaced := false
	for _, file := range reader.File {
		entry, err := writer.Create(file.Name)
		if err != nil {
			_ = writer.Close()
			return nil, err
		}
		if file.Name == cleanSheetPath {
			if _, err := entry.Write([]byte(sheetXML(records))); err != nil {
				_ = writer.Close()
				return nil, err
			}
			replaced = true
			continue
		}
		handle, err := file.Open()
		if err != nil {
			_ = writer.Close()
			return nil, err
		}
		body, err := io.ReadAll(handle)
		_ = handle.Close()
		if err != nil {
			_ = writer.Close()
			return nil, err
		}
		if _, err := entry.Write(body); err != nil {
			_ = writer.Close()
			return nil, err
		}
	}
	if !replaced {
		_ = writer.Close()
		return nil, fmt.Errorf("worksheet %q not found", sheetPath)
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func writeEntry(writer *zip.Writer, name string, content string) error {
	entry, err := writer.Create(name)
	if err != nil {
		return err
	}
	_, err = entry.Write([]byte(content))
	return err
}

func sheetXML(records [][]string) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for rowIndex, record := range records {
		rowNumber := rowIndex + 1
		builder.WriteString(fmt.Sprintf(`<row r="%d">`, rowNumber))
		for columnIndex, value := range record {
			cell := cellRef(rowNumber, columnIndex+1)
			builder.WriteString(`<c r="`)
			builder.WriteString(cell)
			builder.WriteString(`" t="inlineStr"><is><t xml:space="preserve">`)
			builder.WriteString(escapeText(value))
			builder.WriteString(`</t></is></c>`)
		}
		builder.WriteString(`</row>`)
	}
	builder.WriteString(`</sheetData></worksheet>`)
	return builder.String()
}

func cellRef(row int, column int) string {
	var letters []byte
	for column > 0 {
		column--
		letters = append([]byte{byte('A' + column%26)}, letters...)
		column /= 26
	}
	return fmt.Sprintf("%s%d", string(letters), row)
}

func escapeText(value string) string {
	var cleaned strings.Builder
	for _, r := range value {
		if r == '\t' || r == '\n' || r == '\r' || (r >= 0x20 && r <= 0xD7FF) || (r >= 0xE000 && r <= 0xFFFD) || (r >= 0x10000 && r <= 0x10FFFF) {
			cleaned.WriteRune(r)
		}
	}
	var escaped bytes.Buffer
	_ = xml.EscapeText(&escaped, []byte(cleaned.String()))
	return escaped.String()
}

func coreXML(now time.Time) string {
	timestamp := now.Format(time.RFC3339)
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">` +
		`<dc:creator>opencsv</dc:creator><cp:lastModifiedBy>opencsv</cp:lastModifiedBy>` +
		`<dcterms:created xsi:type="dcterms:W3CDTF">` + timestamp + `</dcterms:created>` +
		`<dcterms:modified xsi:type="dcterms:W3CDTF">` + timestamp + `</dcterms:modified>` +
		`</cp:coreProperties>`
}

func contentTypesXMLForSheets(count int) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`)
	builder.WriteString(`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`)
	builder.WriteString(`<Default Extension="xml" ContentType="application/xml"/>`)
	builder.WriteString(`<Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>`)
	builder.WriteString(`<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>`)
	builder.WriteString(`<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>`)
	for index := 1; index <= count; index++ {
		builder.WriteString(fmt.Sprintf(`<Override PartName="/xl/worksheets/sheet%d.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>`, index))
	}
	builder.WriteString(`</Types>`)
	return builder.String()
}

func workbookXMLForSheets(sheets []WorkbookSheet) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString(`<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets>`)
	for index, sheet := range sheets {
		name := strings.TrimSpace(sheet.Name)
		if name == "" {
			name = fmt.Sprintf("Sheet%d", index+1)
		}
		builder.WriteString(fmt.Sprintf(`<sheet name="%s" sheetId="%d" r:id="rId%d"/>`, escapeAttr(name), index+1, index+1))
	}
	builder.WriteString(`</sheets></workbook>`)
	return builder.String()
}

func workbookRelsXMLForSheets(count int) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for index := 1; index <= count; index++ {
		builder.WriteString(fmt.Sprintf(`<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet%d.xml"/>`, index, index))
	}
	builder.WriteString(`</Relationships>`)
	return builder.String()
}

func escapeAttr(value string) string {
	var escaped bytes.Buffer
	_ = xml.EscapeText(&escaped, []byte(value))
	return escaped.String()
}

const contentTypesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>
  <Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
</Types>`

const rootRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
  <Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>`

const appXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">
  <Application>opencsv</Application>
</Properties>`

const workbookXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheets>
    <sheet name="Sheet1" sheetId="1" r:id="rId1"/>
  </sheets>
</workbook>`

const workbookRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
</Relationships>`
