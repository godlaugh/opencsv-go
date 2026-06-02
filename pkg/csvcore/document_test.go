package csvcore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenDocumentDetectsMetadata(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.tsv")
	data := []byte("\ufeffid\tname\r\n1\tAda\r\n2\tLinus\r\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 2, IndexRows: 2})
	if err != nil {
		t.Fatal(err)
	}
	if document.Metadata.Encoding.Name != EncodingUTF8BOM {
		t.Fatalf("encoding = %s, want utf-8-bom", document.Metadata.Encoding.Name)
	}
	if document.Metadata.LineEnding.Name != LineEndingCRLF {
		t.Fatalf("line ending = %s, want crlf", document.Metadata.LineEnding.Name)
	}
	if document.Metadata.Dialect.Delimiter != '\t' {
		t.Fatalf("delimiter = %q, want tab", document.Metadata.Dialect.Delimiter)
	}
	if !document.Metadata.Dialect.HasHeader {
		t.Fatal("expected header")
	}
	if document.RowsScanned != 2 {
		t.Fatalf("rows scanned = %d, want 2", document.RowsScanned)
	}
}

func TestNormalizeLineEndingName(t *testing.T) {
	tests := []struct {
		raw  string
		want LineEndingName
	}{
		{raw: "lf", want: LineEndingLF},
		{raw: "unix", want: LineEndingLF},
		{raw: "\\n", want: LineEndingLF},
		{raw: "crlf", want: LineEndingCRLF},
		{raw: "windows", want: LineEndingCRLF},
		{raw: "\\r\\n", want: LineEndingCRLF},
		{raw: "cr", want: LineEndingCR},
		{raw: "classic-mac", want: LineEndingCR},
		{raw: "\\r", want: LineEndingCR},
	}
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			got, err := NormalizeLineEndingName(tt.raw)
			if err != nil {
				t.Fatal(err)
			}
			if got.Name != tt.want || got.Sequence == "" {
				t.Fatalf("NormalizeLineEndingName(%q) = %#v, want %s", tt.raw, got, tt.want)
			}
		})
	}
	if _, err := NormalizeLineEndingName("vertical-tab"); err == nil {
		t.Fatal("expected unsupported line ending error")
	}
}
