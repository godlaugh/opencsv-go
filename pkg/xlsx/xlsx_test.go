package xlsx

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestEncodeWorkbookWritesInlineStringSheet(t *testing.T) {
	data, err := EncodeWorkbook([][]string{
		{"id", "name"},
		{"1", "Ada & Grace"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("PK")) {
		t.Fatalf("xlsx does not start with zip magic: %q", data[:2])
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	entries := map[string]string{}
	for _, file := range reader.File {
		handle, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(handle)
		_ = handle.Close()
		if err != nil {
			t.Fatal(err)
		}
		entries[file.Name] = string(body)
	}
	sheet := entries["xl/worksheets/sheet1.xml"]
	if !strings.Contains(sheet, `<c r="A1" t="inlineStr"><is><t xml:space="preserve">id</t></is></c>`) {
		t.Fatalf("missing A1 id cell in %s", sheet)
	}
	if !strings.Contains(sheet, `Ada &amp; Grace`) {
		t.Fatalf("cell text was not escaped in %s", sheet)
	}
	if entries["[Content_Types].xml"] == "" || entries["xl/workbook.xml"] == "" {
		t.Fatalf("missing workbook metadata entries: %#v", entries)
	}
}

func TestDecodeWorkbookReadsInlineStringsAndSparseCells(t *testing.T) {
	data, err := EncodeWorkbook([][]string{
		{"id", "name", "score"},
		{"1", "Ada", "9"},
		{"2", "Grace & Hopper", ""},
	})
	if err != nil {
		t.Fatal(err)
	}
	records, err := DecodeWorkbook(data)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"id", "name", "score"},
		{"1", "Ada", "9"},
		{"2", "Grace & Hopper", ""},
	}
	if len(records) != len(want) {
		t.Fatalf("rows = %#v, want %#v", records, want)
	}
	for row := range want {
		if strings.Join(records[row], "\x00") != strings.Join(want[row], "\x00") {
			t.Fatalf("row %d = %#v, want %#v", row, records[row], want[row])
		}
	}
}

func TestDecodeWorkbookReadsSharedStringsBooleansAndGaps(t *testing.T) {
	data := zippedWorkbook(t, map[string]string{
		"xl/workbook.xml":            `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Data" sheetId="1" r:id="rId42"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId42" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet7.xml"/></Relationships>`,
		"xl/sharedStrings.xml":       `<sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><si><t>name</t></si><si><r><t>Ada</t></r><r><t> Lovelace</t></r></si></sst>`,
		"xl/worksheets/sheet7.xml":   `<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData><row r="1"><c r="B1" t="s"><v>0</v></c><c r="D1"><v>42</v></c></row><row r="3"><c r="A3" t="b"><v>1</v></c><c r="B3" t="s"><v>1</v></c></row></sheetData></worksheet>`,
	})
	records, err := DecodeWorkbook(data)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"", "name", "", "42"},
		nil,
		{"TRUE", "Ada Lovelace"},
	}
	if len(records) != len(want) {
		t.Fatalf("rows = %#v, want %#v", records, want)
	}
	for row := range want {
		if strings.Join(records[row], "\x00") != strings.Join(want[row], "\x00") {
			t.Fatalf("row %d = %#v, want %#v", row, records[row], want[row])
		}
	}
}

func TestDecodeWorkbookCanSelectWorksheet(t *testing.T) {
	data, err := EncodeWorkbookSheets([]WorkbookSheet{
		{Name: "Summary", Records: [][]string{{"kind", "value"}, {"rows", "2"}}},
		{Name: "Data", Records: [][]string{{"id", "name"}, {"1", "Ada"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	sheets, err := WorkbookSheets(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(sheets) != 2 || sheets[0].Name != "Summary" || sheets[1].Name != "Data" {
		t.Fatalf("sheets = %#v", sheets)
	}

	records, sheets, selected, err := DecodeWorkbookWithOptions(data, DecodeOptions{SheetName: "Data"})
	if err != nil {
		t.Fatal(err)
	}
	if selected.Index != 2 || selected.Name != "Data" || selected.Path != "xl/worksheets/sheet2.xml" {
		t.Fatalf("selected = %#v", selected)
	}
	if len(sheets) != 2 {
		t.Fatalf("sheet list = %#v", sheets)
	}
	if got := strings.Join(records[1], ","); got != "1,Ada" {
		t.Fatalf("selected records = %#v", records)
	}

	records, _, selected, err = DecodeWorkbookWithOptions(data, DecodeOptions{SheetIndex: 1})
	if err != nil {
		t.Fatal(err)
	}
	if selected.Name != "Summary" || strings.Join(records[1], ",") != "rows,2" {
		t.Fatalf("index-selected records = %#v selected=%#v", records, selected)
	}
}

func TestDecodeWorkbookFormatsDatesAndFormulaCachedValues(t *testing.T) {
	data := zippedWorkbook(t, map[string]string{
		"xl/workbook.xml":            `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Data" sheetId="1" r:id="rId1"/></sheets></workbook>`,
		"xl/_rels/workbook.xml.rels": `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/styles.xml":              `<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><numFmts count="1"><numFmt numFmtId="165" formatCode="yyyy-mm-dd hh:mm:ss"/></numFmts><cellXfs count="4"><xf numFmtId="0"/><xf numFmtId="14"/><xf numFmtId="22"/><xf numFmtId="165"/></cellXfs></styleSheet>`,
		"xl/worksheets/sheet1.xml":   `<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData><row r="1"><c r="A1" t="inlineStr"><is><t>date</t></is></c><c r="B1" t="inlineStr"><is><t>formula</t></is></c><c r="C1" t="inlineStr"><is><t>datetime</t></is></c><c r="D1" t="inlineStr"><is><t>custom</t></is></c></row><row r="2"><c r="A2" s="1"><v>45000</v></c><c r="B2"><f>2+3</f><v>5</v></c><c r="C2" s="2"><v>45000.5</v></c><c r="D2" s="3"><v>45000.75</v></c></row></sheetData></worksheet>`,
	})
	records, err := DecodeWorkbook(data)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"date", "formula", "datetime", "custom"},
		{"2023-03-15", "5", "2023-03-15 12:00:00", "2023-03-15 18:00:00"},
	}
	if len(records) != len(want) {
		t.Fatalf("rows = %#v, want %#v", records, want)
	}
	for row := range want {
		if strings.Join(records[row], "\x00") != strings.Join(want[row], "\x00") {
			t.Fatalf("row %d = %#v, want %#v", row, records[row], want[row])
		}
	}
}

func TestUpdateWorkbookSheetPreservesOtherSheets(t *testing.T) {
	data, err := EncodeWorkbookSheets([]WorkbookSheet{
		{Name: "Summary", Records: [][]string{{"kind", "value"}, {"rows", "2"}}},
		{Name: "Data", Records: [][]string{{"id", "name"}, {"1", "Ada"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := UpdateWorkbookSheet(data, "xl/worksheets/sheet2.xml", [][]string{{"id", "name"}, {"2", "Grace"}})
	if err != nil {
		t.Fatal(err)
	}
	summary, _, selected, err := DecodeWorkbookWithOptions(updated, DecodeOptions{SheetName: "Summary"})
	if err != nil {
		t.Fatal(err)
	}
	if selected.Name != "Summary" || strings.Join(summary[1], ",") != "rows,2" {
		t.Fatalf("summary changed: %#v selected=%#v", summary, selected)
	}
	dataRecords, _, selected, err := DecodeWorkbookWithOptions(updated, DecodeOptions{SheetName: "Data"})
	if err != nil {
		t.Fatal(err)
	}
	if selected.Name != "Data" || strings.Join(dataRecords[1], ",") != "2,Grace" {
		t.Fatalf("data sheet not updated: %#v selected=%#v", dataRecords, selected)
	}
}

func zippedWorkbook(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	for name, body := range entries {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
