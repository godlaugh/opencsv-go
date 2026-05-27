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
