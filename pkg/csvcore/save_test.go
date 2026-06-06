package csvcore

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestSaveRecordsCreatesBackupAndPreservesUTF8BOMAndCRLF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	original := []byte("\xef\xbb\xbfid,name\r\n1,Ada\r\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 5})
	if err != nil {
		t.Fatal(err)
	}
	records, err := ReadAllRecords(path, document.Metadata)
	if err != nil {
		t.Fatal(err)
	}
	records = append(records, []string{"2", "Linus"})

	backupPath := filepath.Join(dir, "sample.csv.bak")
	result, err := SaveRecords(path, document.Metadata, records, SaveOptions{
		CreateBackup: true,
		BackupPath:   backupPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.BackupPath != backupPath {
		t.Fatalf("backup path = %q, want %q", result.BackupPath, backupPath)
	}

	backup, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(backup, original) {
		t.Fatal("backup does not match original bytes")
	}

	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(saved, []byte{0xEF, 0xBB, 0xBF}) {
		t.Fatal("saved file lost UTF-8 BOM")
	}
	if bytes.Contains(saved, []byte("\n")) && !bytes.Contains(saved, []byte("\r\n")) {
		t.Fatal("saved file lost CRLF line endings")
	}
	if !bytes.Contains(saved, []byte("2,Linus\r\n")) {
		t.Fatalf("saved file did not include appended record: %q", saved)
	}
}

func TestEncodeRecordsUsesDialectAndLineEnding(t *testing.T) {
	metadata := FileMetadata{
		LineEnding: LineEndingInfo{Name: LineEndingCRLF},
		Dialect:    Dialect{Delimiter: '\t'},
	}
	got, err := EncodeRecords([][]string{
		{"id", "note"},
		{"1", "has\ttab"},
	}, metadata)
	if err != nil {
		t.Fatal(err)
	}
	want := "id\tnote\r\n1\t\"has\ttab\"\r\n"
	if string(got) != want {
		t.Fatalf("encoded records = %q, want %q", string(got), want)
	}
}

func TestWriteRecordUTF8UsesRowQuoteRulesAndLineEnding(t *testing.T) {
	metadata := FileMetadata{
		LineEnding: LineEndingInfo{Name: LineEndingCRLF},
		Dialect: Dialect{
			Delimiter:   ',',
			QuoteChar:   `"`,
			QuoteMode:   QuoteModeManual,
			HasHeader:   true,
			QuoteHeader: true,
			QuoteBody:   true,
		},
	}
	var output bytes.Buffer
	if err := WriteRecordUTF8(&output, metadata, 0, []string{"id", "name"}); err != nil {
		t.Fatal(err)
	}
	if err := WriteRecordUTF8(&output, metadata, 1, []string{"1", "Ada"}); err != nil {
		t.Fatal(err)
	}
	want := "\"id\",\"name\"\r\n\"1\",\"Ada\"\r\n"
	if output.String() != want {
		t.Fatalf("streamed records = %q, want %q", output.String(), want)
	}
}

func TestEncodeRecordsQuoteModes(t *testing.T) {
	records := [][]string{
		{"id", "note"},
		{"1", "plain"},
		{"2", "has,comma"},
	}
	metadata := FileMetadata{
		Dialect: Dialect{Delimiter: ',', QuoteChar: `"`, QuoteMode: QuoteModeAlways},
	}
	got, err := EncodeRecords(records, metadata)
	if err != nil {
		t.Fatal(err)
	}
	want := "\"id\",\"note\"\n\"1\",\"plain\"\n\"2\",\"has,comma\"\n"
	if string(got) != want {
		t.Fatalf("always quote = %q, want %q", string(got), want)
	}

	metadata.Dialect.QuoteMode = QuoteModeNever
	got, err = EncodeRecords(records, metadata)
	if err != nil {
		t.Fatal(err)
	}
	want = "id,note\n1,plain\n2,has,comma\n"
	if string(got) != want {
		t.Fatalf("never quote = %q, want %q", string(got), want)
	}

	metadata.Dialect = Dialect{Delimiter: ',', QuoteChar: "'", QuoteMode: QuoteModeMinimal}
	got, err = EncodeRecords([][]string{{"id", "note"}, {"1", "O'Reilly"}}, metadata)
	if err != nil {
		t.Fatal(err)
	}
	want = "id,note\n1,'O''Reilly'\n"
	if string(got) != want {
		t.Fatalf("single quote = %q, want %q", string(got), want)
	}

	metadata.Dialect = Dialect{Delimiter: ',', QuoteChar: `"`, QuoteMode: QuoteModeMinimal, QuoteEscape: QuoteEscapeBackslash}
	got, err = EncodeRecords([][]string{{"id", "note"}, {"1", `has "quote" and \ slash`}}, metadata)
	if err != nil {
		t.Fatal(err)
	}
	want = "id,note\n1,\"has \\\"quote\\\" and \\\\ slash\"\n"
	if string(got) != want {
		t.Fatalf("backslash quote escape = %q, want %q", string(got), want)
	}
	readRecords, err := ParseRecordsFromUTF8(got, metadata.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if readRecords[1][1] != `has "quote" and \ slash` {
		t.Fatalf("backslash quote readback = %q", readRecords[1][1])
	}

	metadata.Dialect = Dialect{
		Delimiter:    ',',
		QuoteChar:    `"`,
		QuoteMode:    QuoteModeManual,
		QuoteCells:   []Cell{{Row: 1, Col: 1}},
		NoQuoteCells: []Cell{{Row: 2, Col: 1}},
	}
	got, err = EncodeRecords(records, metadata)
	if err != nil {
		t.Fatal(err)
	}
	want = "id,note\n1,\"plain\"\n2,has,comma\n"
	if string(got) != want {
		t.Fatalf("manual quote = %q, want %q", string(got), want)
	}

	metadata.Dialect = Dialect{
		Delimiter:    ',',
		QuoteChar:    `"`,
		QuoteMode:    QuoteModeManual,
		HasHeader:    true,
		QuoteHeader:  true,
		QuoteColumns: []int{1},
	}
	got, err = EncodeRecords(records, metadata)
	if err != nil {
		t.Fatal(err)
	}
	want = "\"id\",\"note\"\n1,\"plain\"\n2,\"has,comma\"\n"
	if string(got) != want {
		t.Fatalf("manual header/column quote = %q, want %q", string(got), want)
	}

	metadata.Dialect = Dialect{Delimiter: ',', QuoteChar: `"`, QuoteMode: QuoteModeManual, HasHeader: true, QuoteBody: true}
	got, err = EncodeRecords(records, metadata)
	if err != nil {
		t.Fatal(err)
	}
	want = "id,note\n\"1\",\"plain\"\n\"2\",\"has,comma\"\n"
	if string(got) != want {
		t.Fatalf("manual body quote = %q, want %q", string(got), want)
	}
}

func TestBackslashQuoteEscapeRowIndex(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "backslash-escape.csv")
	metadata := FileMetadata{
		Encoding: EncodingInfo{Name: EncodingUTF8},
		Dialect:  Dialect{Delimiter: ',', QuoteChar: `"`, QuoteMode: QuoteModeMinimal, QuoteEscape: QuoteEscapeBackslash},
	}
	records := [][]string{{"id", "note"}, {"1", "line \"one\"\nline two"}, {"2", "plain"}}
	encoded, err := EncodeRecords(records, metadata)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatal(err)
	}
	index, err := BuildRowIndex(path, metadata.Dialect, 0)
	if err != nil {
		t.Fatal(err)
	}
	if index.Rows != 3 {
		t.Fatalf("row index rows = %d, want 3 for %q", index.Rows, string(encoded))
	}
	readBack, err := ReadAllRecords(path, metadata)
	if err != nil {
		t.Fatal(err)
	}
	if readBack[1][1] != records[1][1] {
		t.Fatalf("readback multiline field = %q, want %q", readBack[1][1], records[1][1])
	}
}

func TestSaveRecordsPreservesDetectedQuoteRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mixed-quotes.csv")
	original := []byte("\"id\",name,note\n\"1\",Ada,\"plain\"\n\"2\",Linus,free\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}
	document, err := OpenDocument(path, OpenOptions{PreviewRows: 5})
	if err != nil {
		t.Fatal(err)
	}
	dialect := document.Metadata.Dialect
	if !PreservesQuote(dialect, 0, 0) || !PreservesQuote(dialect, 1, 0) || !PreservesQuote(dialect, 1, 2) {
		t.Fatalf("quote preservation rules missing: %+v", dialect)
	}
	if PreservesQuote(dialect, 0, 1) || PreservesQuote(dialect, 2, 2) {
		t.Fatalf("quote preservation rules too broad: %+v", dialect)
	}
	records, err := ReadAllRecords(path, document.Metadata)
	if err != nil {
		t.Fatal(err)
	}
	records[1][1] = "Grace"
	if _, err := SaveRecords(path, document.Metadata, records, SaveOptions{}); err != nil {
		t.Fatal(err)
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "\"id\",name,note\n\"1\",Grace,\"plain\"\n\"2\",Linus,free\n"
	if string(saved) != want {
		t.Fatalf("saved quote preservation = %q, want %q", string(saved), want)
	}
}

func TestSaveRecordsPreservesLateDetectedQuoteCells(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "late-quotes.csv")
	var input bytes.Buffer
	input.WriteString("id,note\n")
	for row := 1; row <= 90; row++ {
		input.WriteString("row")
		input.WriteString(strconv.Itoa(row))
		input.WriteString(",plain\n")
	}
	input.WriteString("row91,\"late quote\"\n")
	if err := os.WriteFile(path, input.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	sampledDocument, err := OpenDocument(path, OpenOptions{PreviewRows: 5, IndexRows: 0})
	if err != nil {
		t.Fatal(err)
	}
	if PreservesQuote(sampledDocument.Metadata.Dialect, 91, 1) {
		t.Fatalf("sampled quote detection should not see the late cell: %+v", sampledDocument.Metadata.Dialect)
	}
	document, err := OpenDocument(path, OpenOptions{PreviewRows: 5, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if !PreservesQuote(document.Metadata.Dialect, 91, 1) {
		t.Fatalf("late quote cell was not preserved: %+v", document.Metadata.Dialect)
	}
	records, err := ReadAllRecords(path, document.Metadata)
	if err != nil {
		t.Fatal(err)
	}
	records[1][1] = "edited"
	if _, err := SaveRecords(path, document.Metadata, records, SaveOptions{}); err != nil {
		t.Fatal(err)
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(saved, []byte("row91,\"late quote\"\n")) {
		t.Fatalf("late quote cell was not preserved on save: %q", saved)
	}
	if !bytes.Contains(saved, []byte("row1,edited\n")) {
		t.Fatalf("unquoted sampled cell did not save edit unquoted: %q", saved)
	}
	if bytes.Contains(saved, []byte("row1,\"edited\"\n")) {
		t.Fatalf("unquoted sampled cell should remain unquoted: %q", saved)
	}
}

func TestSaveRecordsPreservesUTF16LE(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "utf16.csv")
	encoding := EncodingInfo{Name: EncodingUTF16LE, HasBOM: true}
	encoded, err := EncodeFromUTF8([]byte("id,name\r\n1,Ada\r\n"), encoding)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 5})
	if err != nil {
		t.Fatal(err)
	}
	records, err := ReadAllRecords(path, document.Metadata)
	if err != nil {
		t.Fatal(err)
	}
	records[1][1] = "Grace"

	if _, err := SaveRecords(path, document.Metadata, records, SaveOptions{}); err != nil {
		t.Fatal(err)
	}

	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	detected := DetectEncoding(saved)
	if detected.Name != EncodingUTF16LE || !detected.HasBOM {
		t.Fatalf("encoding = %+v, want utf-16le with BOM", detected)
	}
	decoded, err := DecodeToUTF8(saved, detected)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(decoded, []byte("1,Grace\r\n")) {
		t.Fatalf("saved UTF-16 data did not preserve content/line ending: %q", decoded)
	}
}

func TestSaveRecordsDefaultBackupAvoidsNameCollision(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	document, err := OpenDocument(path, OpenOptions{PreviewRows: 5})
	if err != nil {
		t.Fatal(err)
	}

	first, err := SaveRecords(path, document.Metadata, [][]string{{"id", "name"}, {"1", "Linus"}}, SaveOptions{CreateBackup: true})
	if err != nil {
		t.Fatal(err)
	}
	second, err := SaveRecords(path, document.Metadata, [][]string{{"id", "name"}, {"1", "Grace"}}, SaveOptions{CreateBackup: true})
	if err != nil {
		t.Fatal(err)
	}
	if first.BackupPath == "" || second.BackupPath == "" || first.BackupPath == second.BackupPath {
		t.Fatalf("backup paths were not unique: first=%q second=%q", first.BackupPath, second.BackupPath)
	}
}
