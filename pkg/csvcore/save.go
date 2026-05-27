package csvcore

import (
	"bytes"
	"encoding/csv"
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
		backupPath, err = createBackup(path, options.BackupPath)
		if err != nil {
			return SaveResult{}, err
		}
	}

	utf8Bytes, err := encodeCSVRecords(records, metadata)
	if err != nil {
		return SaveResult{}, err
	}
	output, err := EncodeFromUTF8(utf8Bytes, metadata.Encoding)
	if err != nil {
		return SaveResult{}, err
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, output, 0o644); err != nil {
		return SaveResult{}, err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return SaveResult{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return SaveResult{}, err
	}
	return SaveResult{Path: path, BackupPath: backupPath, SizeBytes: info.Size()}, nil
}

func encodeCSVRecords(records [][]string, metadata FileMetadata) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	writer.Comma = metadata.Dialect.Delimiter
	writer.UseCRLF = metadata.LineEnding.Name == LineEndingCRLF
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	if metadata.LineEnding.Name == LineEndingCR {
		return bytes.ReplaceAll(buffer.Bytes(), []byte("\n"), []byte("\r")), nil
	}
	return buffer.Bytes(), nil
}

func createBackup(path string, requestedPath string) (string, error) {
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
	defer target.Close()

	if _, err := io.Copy(target, source); err != nil {
		return "", err
	}
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
