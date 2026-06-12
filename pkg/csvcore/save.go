package csvcore

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SaveOptions struct {
	CreateBackup bool   `json:"createBackup"`
	BackupPath   string `json:"backupPath"`
}

type SaveResult struct {
	Path       string `json:"path"`
	BackupPath string `json:"backupPath,omitempty"`
	SizeBytes  int64  `json:"sizeBytes"`
}

func SaveRecords(path string, metadata FileMetadata, records [][]string, options SaveOptions) (SaveResult, error) {
	var backupPath string
	if options.CreateBackup {
		var err error
		backupPath, err = CreateBackup(path, options.BackupPath)
		if err != nil {
			return SaveResult{}, err
		}
	}

	utf8Bytes, err := EncodeRecords(records, metadata)
	if err != nil {
		return SaveResult{}, err
	}
	output, err := EncodeFromUTF8(utf8Bytes, metadata.Encoding)
	if err != nil {
		return SaveResult{}, err
	}

	if err := AtomicWriteFile(path, output); err != nil {
		return SaveResult{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return SaveResult{}, err
	}
	return SaveResult{Path: path, BackupPath: backupPath, SizeBytes: info.Size()}, nil
}

func AtomicWriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	temp, err := os.CreateTemp(dir, "."+base+".*.tmp")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	keepTemp := false
	defer func() {
		if !keepTemp {
			_ = os.Remove(tempPath)
		}
	}()

	if err := temp.Chmod(mode); err != nil {
		_ = temp.Close()
		return err
	}
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempPath, path); err != nil {
		return err
	}
	keepTemp = true
	syncParentDir(path)
	return nil
}

func syncParentDir(path string) {
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return
	}
	defer dir.Close()
	_ = dir.Sync()
}

// EncodeRecords formats records as UTF-8 text using the document dialect and
// line-ending policy. It does not apply the original file encoding or BOM.
func EncodeRecords(records [][]string, metadata FileMetadata) ([]byte, error) {
	var buffer bytes.Buffer
	for rowIndex, record := range records {
		if err := WriteRecordUTF8(&buffer, metadata, rowIndex, record); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func WriteRecordUTF8(writer io.Writer, metadata FileMetadata, rowIndex int, record []string) error {
	encoded, err := EncodeRecord(record, metadata, rowIndex)
	if err != nil {
		return err
	}
	_, err = writer.Write(encoded)
	return err
}

func EncodeRecord(record []string, metadata FileMetadata, rowIndex int) ([]byte, error) {
	var buffer bytes.Buffer
	delimiter := metadata.Dialect.Delimiter
	if delimiter == 0 {
		delimiter = ','
	}
	quoteChar := DialectQuoteChar(metadata.Dialect)
	for index, field := range record {
		if index > 0 {
			buffer.WriteRune(delimiter)
		}
		if ShouldQuoteField(field, metadata.Dialect, rowIndex, index) && quoteChar != "" {
			buffer.WriteString(quoteChar)
			buffer.WriteString(EscapeQuotedField(field, metadata.Dialect))
			buffer.WriteString(quoteChar)
		} else {
			buffer.WriteString(field)
		}
	}
	buffer.WriteByte('\n')

	if metadata.LineEnding.Name == LineEndingCRLF {
		return bytes.ReplaceAll(buffer.Bytes(), []byte("\n"), []byte("\r\n")), nil
	}
	if metadata.LineEnding.Name == LineEndingCR {
		return bytes.ReplaceAll(buffer.Bytes(), []byte("\n"), []byte("\r")), nil
	}
	return buffer.Bytes(), nil
}

func EscapeQuotedField(value string, dialect Dialect) string {
	quoteChar := DialectQuoteChar(dialect)
	if quoteChar == "" {
		return value
	}
	if DialectQuoteEscape(dialect) == QuoteEscapeBackslash {
		escaped := strings.ReplaceAll(value, `\`, `\\`)
		if quoteChar != `\` {
			escaped = strings.ReplaceAll(escaped, quoteChar, `\`+quoteChar)
		}
		return escaped
	}
	return strings.ReplaceAll(value, quoteChar, quoteChar+quoteChar)
}

func CreateBackup(path string, requestedPath string) (string, error) {
	backupPath := requestedPath
	if backupPath == "" {
		backupPath = defaultBackupPath(path)
	}

	source, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer source.Close()

	target, err := os.OpenFile(backupPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(target, source); err != nil {
		_ = target.Close()
		_ = os.Remove(backupPath)
		return "", err
	}
	if err := target.Sync(); err != nil {
		_ = target.Close()
		_ = os.Remove(backupPath)
		return "", err
	}
	if err := target.Close(); err != nil {
		_ = os.Remove(backupPath)
		return "", err
	}
	syncParentDir(backupPath)
	return backupPath, nil
}

func defaultBackupPath(path string) string {
	ext := filepath.Ext(path)
	stem := strings.TrimSuffix(path, ext)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	base := fmt.Sprintf("%s.%s%s.bak", stem, timestamp, ext)
	if _, err := os.Stat(base); errors.Is(err, os.ErrNotExist) {
		return base
	}
	for index := 2; ; index++ {
		candidate := fmt.Sprintf("%s.%s.%d%s.bak", stem, timestamp, index, ext)
		if _, err := os.Stat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate
		}
	}
}
