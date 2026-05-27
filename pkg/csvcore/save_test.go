package csvcore

import (
	"bytes"
	"os"
	"path/filepath"
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
