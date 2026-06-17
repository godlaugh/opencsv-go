package csvcore

import (
	"context"
	"errors"
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

func TestOpenDocumentParsesSingleQuoteDialect(t *testing.T) {
	path := filepath.Join(t.TempDir(), "single-quote.csv")
	data := []byte("id,note\n1,'hello, world'\n2,'line\nbreak with \"quote\"'\n3,'Bob''s note'\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 4, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if document.Metadata.Dialect.QuoteChar != "'" {
		t.Fatalf("quote char = %q, want single quote", document.Metadata.Dialect.QuoteChar)
	}
	if len(document.Preview) != 4 || document.Preview[1][1] != "hello, world" || document.Preview[2][1] != "line\nbreak with \"quote\"" || document.Preview[3][1] != "Bob's note" {
		t.Fatalf("preview = %#v", document.Preview)
	}

	records, err := ReadAllRecords(path, document.Metadata)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 4 || records[2][1] != "line\nbreak with \"quote\"" || records[3][1] != "Bob's note" {
		t.Fatalf("records = %#v", records)
	}
}

func TestOpenDocumentDetectsBackslashQuoteEscape(t *testing.T) {
	path := filepath.Join(t.TempDir(), "backslash-quote.csv")
	data := []byte("id,note\n1,\"has \\\"quote\\\" and comma, ok\"\n2,\"line \\\"one\\\"\nline two\"\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if document.Metadata.Dialect.QuoteEscape != QuoteEscapeBackslash {
		t.Fatalf("quote escape = %q, want backslash", document.Metadata.Dialect.QuoteEscape)
	}
	if len(document.Preview) != 3 || document.Preview[1][1] != `has "quote" and comma, ok` || document.Preview[2][1] != "line \"one\"\nline two" {
		t.Fatalf("preview = %#v", document.Preview)
	}

	records, err := ReadAllRecords(path, document.Metadata)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 || records[1][1] != `has "quote" and comma, ok` || records[2][1] != "line \"one\"\nline two" {
		t.Fatalf("records = %#v", records)
	}
	if !PreservesQuote(document.Metadata.Dialect, 1, 1) || !PreservesQuote(document.Metadata.Dialect, 2, 1) {
		t.Fatalf("detected quote preservation missing: %#v", document.Metadata.Dialect)
	}
}

func TestOpenDocumentDetectsShiftJISEncoding(t *testing.T) {
	path := filepath.Join(t.TempDir(), "shift-jis.csv")
	encoding := EncodingInfo{Name: EncodingShiftJIS, Confidence: 1}
	data, err := EncodeFromUTF8([]byte("名前,メモ\n太郎,かな\n"), encoding)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 2, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if document.Metadata.Encoding.Name != EncodingShiftJIS {
		t.Fatalf("encoding = %+v, want shift-jis", document.Metadata.Encoding)
	}
	if len(document.Preview) != 2 || document.Preview[0][0] != "名前" || document.Preview[1][1] != "かな" {
		t.Fatalf("preview = %#v", document.Preview)
	}
}

func TestForEachRecordContextStreamsUTF8BOM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stream.csv")
	if err := os.WriteFile(path, []byte("\ufeffid,name\n1,Ada\n2,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	document, err := OpenDocument(path, OpenOptions{PreviewRows: 2, IndexRows: 0})
	if err != nil {
		t.Fatal(err)
	}
	var rows [][]string
	count, err := ForEachRecordContext(context.Background(), path, document.Metadata, func(row int, record []string) error {
		rows = append(rows, append([]string(nil), record...))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 || len(rows) != 3 {
		t.Fatalf("streamed count/rows = %d/%#v, want 3 rows", count, rows)
	}
	if rows[0][0] != "id" || rows[2][1] != "Grace" {
		t.Fatalf("streamed rows = %#v", rows)
	}
}

func TestForEachRecordContextCancellation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stream-cancel.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	document, err := OpenDocument(path, OpenOptions{PreviewRows: 2, IndexRows: 0})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	_, err = ForEachRecordContext(ctx, path, document.Metadata, func(row int, record []string) error {
		cancel()
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("stream cancellation error = %v, want context canceled", err)
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

func TestNormalizeDelimiter(t *testing.T) {
	tests := []struct {
		raw      string
		name     string
		wantRune rune
		wantName string
	}{
		{raw: "comma", wantRune: ',', wantName: "comma"},
		{raw: "tab", wantRune: '\t', wantName: "tab"},
		{raw: "\\t", wantRune: '\t', wantName: "tab"},
		{raw: "\t", wantRune: '\t', wantName: "tab"},
		{raw: "semicolon", wantRune: ';', wantName: "semicolon"},
		{raw: "pipe", wantRune: '|', wantName: "pipe"},
		{raw: "space", wantRune: ' ', wantName: "space"},
		{raw: " ", wantRune: ' ', wantName: "space"},
		{raw: "^", wantRune: '^', wantName: "^"},
		{raw: "^", name: "caret", wantRune: '^', wantName: "caret"},
	}
	for _, tt := range tests {
		t.Run(tt.raw+tt.name, func(t *testing.T) {
			gotRune, gotName, err := NormalizeDelimiter(tt.raw, tt.name)
			if err != nil {
				t.Fatal(err)
			}
			if gotRune != tt.wantRune || gotName != tt.wantName {
				t.Fatalf("NormalizeDelimiter(%q, %q) = %q/%q, want %q/%q", tt.raw, tt.name, gotRune, gotName, tt.wantRune, tt.wantName)
			}
		})
	}
	if _, _, err := NormalizeDelimiter("\n", ""); err == nil {
		t.Fatal("expected line-break delimiter error")
	}
	if _, _, err := NormalizeDelimiter("::", ""); err == nil {
		t.Fatal("expected multi-character delimiter error")
	}
}
