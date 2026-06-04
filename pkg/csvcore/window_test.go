package csvcore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadRecordWindow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n3,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 2})
	if err != nil {
		t.Fatal(err)
	}

	window, err := ReadRecordWindow(path, document.Metadata, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(window.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(window.Rows))
	}
	if window.Rows[0][1] != "Linus" || window.Rows[1][1] != "Grace" {
		t.Fatalf("unexpected window rows: %#v", window.Rows)
	}
}

func TestReadRecordWindowWithIndexSeeksToWindow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "indexed.csv")
	data := "id,note\n1,\"line\nbreak\"\n2,plain\n3,last\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 2, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if !document.RowIndex.Complete {
		t.Fatalf("row index should be complete: %#v", document.RowIndex)
	}

	window, err := ReadRecordWindowWithIndex(path, document.Metadata, document.RowIndex, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if window.TotalRows != 4 {
		t.Fatalf("total rows = %d, want 4", window.TotalRows)
	}
	if got := window.SourceRows; len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Fatalf("source rows = %#v, want 2,3", got)
	}
	if len(window.Rows) != 2 || window.Rows[0][1] != "plain" || window.Rows[1][1] != "last" {
		t.Fatalf("indexed window rows = %#v", window.Rows)
	}

	empty, err := ReadRecordWindowWithIndex(path, document.Metadata, document.RowIndex, 4, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty.Rows) != 0 || empty.TotalRows != 4 {
		t.Fatalf("empty tail window = %#v, want no rows and total 4", empty)
	}
}

func TestReadRecordWindowWithIndexFallsBackForNonStandardQuote(t *testing.T) {
	path := filepath.Join(t.TempDir(), "single-quote.csv")
	data := "id,note\n1,'line\nbreak'\n2,plain\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	metadata := FileMetadata{
		Path:     path,
		Encoding: EncodingInfo{Name: EncodingUTF8},
		Dialect:  Dialect{Delimiter: ',', QuoteChar: "'", QuoteMode: QuoteModeMinimal},
	}
	index, err := BuildRowIndex(path, metadata.Dialect, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !index.Complete {
		t.Fatalf("row index should be complete: %#v", index)
	}

	window, err := ReadRecordWindowWithIndex(path, metadata, index, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(window.Rows) != 2 || window.Rows[0][0] != "break'" || window.Rows[1][0] != "2" {
		t.Fatalf("fallback window rows = %#v", window.Rows)
	}
}
