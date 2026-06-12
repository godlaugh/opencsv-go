package appserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/opencsv/opencsv/pkg/xlsx"
)

func TestServerOpenEditSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	original := []byte("\xef\xbb\xbfid,name\r\n1,Ada\r\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()

	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID, ok := openPayload["documentId"].(string)
	if !ok || documentID == "" {
		t.Fatalf("missing documentId in %#v", openPayload)
	}

	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        1,
		"value":      "Grace",
	})

	savePayload := postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     true,
	})
	backupPath, ok := savePayload["backupPath"].(string)
	if !ok || backupPath == "" {
		t.Fatalf("missing backupPath in %#v", savePayload)
	}

	backup, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(backup, original) {
		t.Fatal("backup bytes differ from original")
	}

	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(saved, []byte{0xEF, 0xBB, 0xBF}) {
		t.Fatal("saved file lost UTF-8 BOM")
	}
	if !bytes.Contains(saved, []byte("1,Grace\r\n")) {
		t.Fatalf("saved file did not include edit with CRLF: %q", saved)
	}
}

func TestServerOpenCanMarkStdinSource(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stdin-copy.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
		"sourceKind":  "stdin",
		"displayName": "stdin",
	})
	if payload["sourceKind"] != "stdin" || payload["displayName"] != "stdin" {
		t.Fatalf("stdin open response = %#v", payload)
	}
	document := payload["document"].(map[string]any)
	metadata := document["metadata"].(map[string]any)
	if metadata["path"] != path {
		t.Fatalf("stdin document path = %#v, want %q", metadata["path"], path)
	}
}

func TestServerSessionSnapshot(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.csv")
	second := filepath.Join(dir, "second.csv")
	if err := os.WriteFile(first, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("id,name\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	firstPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        first,
		"previewRows": 3,
		"indexRows":   -1,
	})
	firstID := firstPayload["documentId"].(string)
	_ = postJSON(t, handler, "/api/open", map[string]any{
		"path":        second,
		"previewRows": 3,
		"indexRows":   -1,
	})
	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": firstID,
		"row":        1,
		"col":        1,
		"value":      "Grace",
	})

	sessionPayload := getJSON(t, handler, "/api/session?id="+firstID)
	if sessionPayload["documentId"] != firstID || sessionPayload["dirty"] != true {
		t.Fatalf("session payload = %#v", sessionPayload)
	}
	session := sessionPayload["session"].(map[string]any)
	if session["id"] != firstID || session["dirty"] != true {
		t.Fatalf("session snapshot = %#v", session)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+firstID+"&start=1&count=1")
	windowRows := windowPayload["rows"].([]any)
	windowRow := windowRows[0].([]any)
	if windowRow[1] != "Grace" {
		t.Fatalf("session window row = %#v, want edited Grace", windowRow)
	}
}

func TestServerCloseProtectsDirtySession(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.csv")
	second := filepath.Join(dir, "second.csv")
	if err := os.WriteFile(first, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("id,name\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	firstPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        first,
		"previewRows": 3,
		"indexRows":   -1,
	})
	firstID := firstPayload["documentId"].(string)
	secondPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        second,
		"previewRows": 3,
		"indexRows":   -1,
	})
	secondID := secondPayload["documentId"].(string)
	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": firstID,
		"row":        1,
		"col":        1,
		"value":      "Grace",
	})

	blockedResponse := postJSONResponse(t, handler, "/api/close", map[string]any{
		"documentId": firstID,
	})
	if blockedResponse.Code != http.StatusBadRequest {
		t.Fatalf("dirty close status = %d, want 400", blockedResponse.Code)
	}
	var blockedPayload map[string]any
	if err := json.NewDecoder(blockedResponse.Body).Decode(&blockedPayload); err != nil {
		t.Fatal(err)
	}
	if blockedPayload["error"] != "document has unsaved changes" {
		t.Fatalf("dirty close payload = %#v", blockedPayload)
	}
	firstSession := getJSON(t, handler, "/api/session?id="+firstID)
	if firstSession["dirty"] != true {
		t.Fatalf("first session after blocked close = %#v", firstSession)
	}

	forcePayload := postJSON(t, handler, "/api/close", map[string]any{
		"documentId": firstID,
		"force":      true,
	})
	if forcePayload["closed"] != true || forcePayload["documentId"] != firstID || forcePayload["dirty"] != true {
		t.Fatalf("force close payload = %#v", forcePayload)
	}
	missingRequest := httptest.NewRequest(http.MethodGet, "/api/session?id="+firstID, nil)
	missingResponse := httptest.NewRecorder()
	handler.ServeHTTP(missingResponse, missingRequest)
	if missingResponse.Code != http.StatusBadRequest {
		t.Fatalf("closed session status = %d, want 400", missingResponse.Code)
	}

	secondSession := getJSON(t, handler, "/api/session?id="+secondID)
	if secondSession["documentId"] != secondID || secondSession["dirty"] != false {
		t.Fatalf("second session after closing first = %#v", secondSession)
	}
	cleanPayload := postJSON(t, handler, "/api/close", map[string]any{
		"documentId": secondID,
	})
	if cleanPayload["closed"] != true || cleanPayload["documentId"] != secondID || cleanPayload["dirty"] != false {
		t.Fatalf("clean close payload = %#v", cleanPayload)
	}
}

func TestServerOpenEditSavePreservesLateDetectedQuoteCell(t *testing.T) {
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

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	document := openPayload["document"].(map[string]any)
	dialect := document["metadata"].(map[string]any)["dialect"].(map[string]any)
	quoteCells := dialect["quoteCells"].([]any)
	if len(quoteCells) != 1 {
		t.Fatalf("quote cells = %#v, want late cell only", quoteCells)
	}
	lateCell := quoteCells[0].(map[string]any)
	if lateCell["row"] != float64(91) || lateCell["col"] != float64(1) {
		t.Fatalf("late quote cell = %#v, want row 91 col 1", lateCell)
	}

	documentID := openPayload["documentId"].(string)
	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        1,
		"value":      "edited",
	})
	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(saved, []byte("row1,edited\n")) || !bytes.Contains(saved, []byte("row91,\"late quote\"\n")) {
		t.Fatalf("saved file did not preserve edit and late quote: %q", saved)
	}
}

func TestServerListAndRestoreBackup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	original := []byte("id,name\n1,Ada\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        1,
		"value":      "Grace",
	})
	savePayload := postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     true,
	})
	backupPath := savePayload["backupPath"].(string)

	listPayload := getJSON(t, handler, "/api/backups?id="+documentID)
	backups := listPayload["backups"].([]any)
	if len(backups) != 1 {
		t.Fatalf("backups = %#v, want one", backups)
	}
	backup := backups[0].(map[string]any)
	if backup["path"] != backupPath || backup["name"] == "" {
		t.Fatalf("unexpected backup info: %#v", backup)
	}

	restorePayload := postJSON(t, handler, "/api/restore-backup", map[string]any{
		"documentId": documentID,
		"path":       backupPath,
	})
	if restorePayload["dirty"] != true {
		t.Fatalf("restore should leave document dirty: %#v", restorePayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada" {
		t.Fatalf("restored rows = %#v", rows)
	}

	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(saved, original) {
		t.Fatalf("saved restore = %q, want %q", saved, original)
	}

	otherBackup := filepath.Join(dir, "other.20260602T000000Z.csv.bak")
	if err := os.WriteFile(otherBackup, []byte("id,name\n9,Edsger\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	response := postJSONResponse(t, handler, "/api/restore-backup", map[string]any{
		"documentId": documentID,
		"path":       otherBackup,
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("restore of unrelated backup status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestServerSaveAsWritesNewPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	targetPath := filepath.Join(dir, "sample-copy.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        1,
		"value":      "Grace",
	})

	savePayload := postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"path":       targetPath,
		"backup":     true,
	})
	if savePayload["path"] != targetPath {
		t.Fatalf("save as path = %#v, want %s", savePayload["path"], targetPath)
	}
	if _, hasBackup := savePayload["backupPath"].(string); hasBackup {
		t.Fatalf("new save-as target should not create backup: %#v", savePayload)
	}
	saved, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(saved) != "id,name\n1,Grace\n" {
		t.Fatalf("save-as content = %q", string(saved))
	}
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(original) != "id,name\n1,Ada\n" {
		t.Fatalf("original changed after save-as: %q", string(original))
	}
}

func TestServerNewDocumentSaveAs(t *testing.T) {
	dir := t.TempDir()
	targetPath := filepath.Join(dir, "created.csv")

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/new", map[string]any{
		"name":    "../Untitled",
		"rows":    1,
		"columns": 2,
	})
	documentID, ok := openPayload["documentId"].(string)
	if !ok || documentID == "" {
		t.Fatalf("missing documentId in %#v", openPayload)
	}
	document := openPayload["document"].(map[string]any)
	metadata := document["metadata"].(map[string]any)
	if filepath.Base(metadata["path"].(string)) != "Untitled.csv" {
		t.Fatalf("new document path = %#v", metadata["path"])
	}
	dialect := metadata["dialect"].(map[string]any)
	if dialect["hasHeader"] != true {
		t.Fatalf("new document should default to header row: %#v", dialect)
	}
	preview := document["preview"].([]any)
	if len(preview) < 2 {
		t.Fatalf("new document preview = %#v", preview)
	}
	header := preview[0].([]any)
	if header[0] != "Column 1" || header[1] != "Column 2" {
		t.Fatalf("new document header = %#v", header)
	}

	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        0,
		"value":      "Ada",
	})
	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        1,
		"value":      "90",
	})

	savePayload := postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"path":       targetPath,
		"backup":     true,
	})
	if savePayload["path"] != targetPath {
		t.Fatalf("new document save-as path = %#v, want %s", savePayload["path"], targetPath)
	}
	if _, hasBackup := savePayload["backupPath"].(string); hasBackup {
		t.Fatalf("new document save-as target should not create backup: %#v", savePayload)
	}
	saved, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(saved) != "Column 1,Column 2\nAda,90\n" {
		t.Fatalf("new document save-as content = %q", string(saved))
	}
}

func TestServerOpenSQLResultDocument(t *testing.T) {
	handler := New(t.TempDir()).Handler()
	payload := postJSON(t, handler, "/api/sql-result/open", map[string]any{
		"name":    "../scores.sql-result.csv",
		"columns": []string{"name", "score", ""},
		"rows": [][]string{
			{"Grace", "10", "top", "ignored"},
			{"Ada"},
		},
	})

	if payload["sourceKind"] != "sql-result" || payload["displayName"] != "scores.sql-result.csv" {
		t.Fatalf("sql result source/display = %#v", payload)
	}
	documentID, ok := payload["documentId"].(string)
	if !ok || documentID == "" {
		t.Fatalf("missing documentId in %#v", payload)
	}
	document := payload["document"].(map[string]any)
	metadata := document["metadata"].(map[string]any)
	if filepath.Base(metadata["path"].(string)) != "scores.sql-result.csv" {
		t.Fatalf("sql result path = %#v", metadata["path"])
	}
	dialect := metadata["dialect"].(map[string]any)
	if dialect["hasHeader"] != true || dialect["estimatedColumns"] != float64(3) {
		t.Fatalf("sql result dialect = %#v", dialect)
	}
	preview := document["preview"].([]any)
	if len(preview) < 3 {
		t.Fatalf("sql result preview = %#v", preview)
	}
	header := preview[0].([]any)
	if header[0] != "name" || header[1] != "score" || header[2] != "Column 3" {
		t.Fatalf("sql result header = %#v", header)
	}
	first := preview[1].([]any)
	if first[0] != "Grace" || first[1] != "10" || first[2] != "top" || len(first) != 3 {
		t.Fatalf("sql result first row = %#v", first)
	}
	second := preview[2].([]any)
	if second[0] != "Ada" || second[1] != "" || second[2] != "" {
		t.Fatalf("sql result padded row = %#v", second)
	}

	sessionPayload := getJSON(t, handler, "/api/session?id="+documentID)
	if sessionPayload["sourceKind"] != "sql-result" || sessionPayload["displayName"] != "scores.sql-result.csv" {
		t.Fatalf("sql result session = %#v", sessionPayload)
	}
}

func TestServerOpenSQLResultDocumentRequiresColumns(t *testing.T) {
	response := postJSONResponse(t, New(t.TempDir()).Handler(), "/api/sql-result/open", map[string]any{
		"name": "empty.csv",
		"rows": [][]string{{"Ada"}},
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("empty sql result status = %d, want %d", response.Code, http.StatusBadRequest)
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), "columns are required") {
		t.Fatalf("empty sql result payload = %#v", payload)
	}
}

func TestServerOpenXLSXPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.xlsx")
	data, err := xlsx.EncodeWorkbook([][]string{
		{"id", "name"},
		{"1", "Ada"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	document := openPayload["document"].(map[string]any)
	metadata := document["metadata"].(map[string]any)
	if metadata["path"] != path {
		t.Fatalf("metadata path = %#v, want %q", metadata["path"], path)
	}
	preview := document["preview"].([]any)
	if got := preview[1].([]any)[1]; got != "Ada" {
		t.Fatalf("xlsx preview value = %#v, want Ada", got)
	}
}

func TestServerOpenXLSXSelectedSheet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "multi-sheet.xlsx")
	data, err := xlsx.EncodeWorkbookSheets([]xlsx.WorkbookSheet{
		{Name: "Summary", Records: [][]string{{"kind", "value"}, {"rows", "2"}}},
		{Name: "Data", Records: [][]string{{"id", "name"}, {"1", "Ada"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
		"sheetName":   "Data",
	})
	document := openPayload["document"].(map[string]any)
	metadata := document["metadata"].(map[string]any)
	workbook := metadata["workbook"].(map[string]any)
	if workbook["selectedSheetName"] != "Data" || workbook["selectedSheetIndex"] != float64(2) {
		t.Fatalf("workbook metadata = %#v", workbook)
	}
	sheets := workbook["sheets"].([]any)
	if len(sheets) != 2 || sheets[1].(map[string]any)["name"] != "Data" {
		t.Fatalf("sheets = %#v", sheets)
	}
	preview := document["preview"].([]any)
	if got := preview[1].([]any)[1]; got != "Ada" {
		t.Fatalf("selected sheet preview value = %#v, want Ada", got)
	}
}

func TestServerOpenWithFormatOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pipe.txt")
	if err := os.WriteFile(path, []byte("id|name\n1|Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":          path,
		"previewRows":   3,
		"indexRows":     -1,
		"delimiter":     "|",
		"delimiterName": "pipe",
		"hasHeader":     false,
		"encodingName":  "utf-8",
		"quoteChar":     "single",
		"quoteMode":     "always",
		"headerColumns": 1,
	})
	document := openPayload["document"].(map[string]any)
	metadata := document["metadata"].(map[string]any)
	dialect := metadata["dialect"].(map[string]any)
	if dialect["delimiterName"] != "pipe" {
		t.Fatalf("delimiterName = %#v, want pipe", dialect["delimiterName"])
	}
	if dialect["hasHeader"] != false {
		t.Fatalf("hasHeader = %#v, want false", dialect["hasHeader"])
	}
	if dialect["quoteChar"] != "'" || dialect["quoteMode"] != "always" {
		t.Fatalf("quote override = %#v", dialect)
	}
	if dialect["headerColumns"] != float64(1) {
		t.Fatalf("headerColumns = %#v, want 1", dialect["headerColumns"])
	}
	preview := document["preview"].([]any)
	if got := preview[0].([]any)[0]; got != "id" {
		t.Fatalf("preview first field = %#v, want id", got)
	}
}

func TestServerFileSample(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id|name\n1|Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	payload := postJSON(t, handler, "/api/file-sample", map[string]any{
		"path":     path,
		"maxBytes": 7,
	})
	if payload["sample"] != "id|name" {
		t.Fatalf("sample = %#v, want id|name", payload["sample"])
	}
}

func TestServerBatchEditUndo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	editPayload := postJSON(t, handler, "/api/edit-batch", map[string]any{
		"documentId": documentID,
		"edits": []map[string]any{
			{"row": 1, "col": 1, "value": "X"},
			{"row": 2, "col": 1, "value": "X"},
		},
	})
	if editPayload["count"].(float64) != 2 {
		t.Fatalf("batch edit count = %#v, want 2", editPayload["count"])
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "X" || rows[2].([]any)[1] != "X" {
		t.Fatalf("rows after batch edit = %#v", rows)
	}

	_ = postJSON(t, handler, "/api/undo", map[string]any{"documentId": documentID})
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows = windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada" || rows[2].([]any)[1] != "Linus" {
		t.Fatalf("rows after undo = %#v", rows)
	}
}

func TestServerHeaderMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "header.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	payload := postJSON(t, handler, "/api/header-mode", map[string]any{
		"documentId": documentID,
		"hasHeader":  false,
	})
	if payload["dirty"].(bool) {
		t.Fatalf("header mode payload should not be dirty: %#v", payload)
	}
	session := payload["session"].(map[string]any)
	document := session["document"].(map[string]any)
	metadata := document["metadata"].(map[string]any)
	dialect := metadata["dialect"].(map[string]any)
	if dialect["hasHeader"].(bool) {
		t.Fatalf("hasHeader = true after disabling: %#v", dialect)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if rows[0].([]any)[0] != "id" || windowPayload["totalRows"].(float64) != 3 {
		t.Fatalf("window after header mode = %#v", windowPayload)
	}
}

func TestServerQuotePolicySave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "quote.csv")
	if err := os.WriteFile(path, []byte("id,note\n1,plain\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	policyPayload := postJSON(t, handler, "/api/quote-policy", map[string]any{
		"documentId": documentID,
		"quoteChar":  "double",
		"quoteMode":  "always",
	})
	if policyPayload["dirty"] != true {
		t.Fatalf("dirty after quote policy = %#v, want true", policyPayload["dirty"])
	}
	document := policyPayload["session"].(map[string]any)["document"].(map[string]any)
	dialect := document["metadata"].(map[string]any)["dialect"].(map[string]any)
	if dialect["quoteMode"] != "always" {
		t.Fatalf("quoteMode = %#v, want always", dialect["quoteMode"])
	}

	backup := false
	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     backup,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "\"id\",\"note\"\n\"1\",\"plain\"\n"
	if string(saved) != want {
		t.Fatalf("saved quote policy = %q, want %q", string(saved), want)
	}
}

func TestServerQuotePolicyBackslashEscapeSaveAndReadback(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "quote-backslash.csv")
	if err := os.WriteFile(path, []byte("id,note\n1,\"has \"\"quote\"\"\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	policyPayload := postJSON(t, handler, "/api/quote-policy", map[string]any{
		"documentId":  documentID,
		"quoteChar":   "double",
		"quoteMode":   "minimal",
		"quoteEscape": "backslash",
	})
	document := policyPayload["session"].(map[string]any)["document"].(map[string]any)
	dialect := document["metadata"].(map[string]any)["dialect"].(map[string]any)
	if dialect["quoteEscape"] != "backslash" {
		t.Fatalf("quoteEscape = %#v, want backslash", dialect["quoteEscape"])
	}

	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "id,note\n1,\"has \\\"quote\\\"\"\n"
	if string(saved) != want {
		t.Fatalf("saved quote escape = %q, want %q", string(saved), want)
	}

	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=2")
	rows := windowPayload["rows"].([]any)
	record := rows[1].([]any)
	if record[1] != `has "quote"` {
		t.Fatalf("readback cell = %#v, want quote-preserved value", record[1])
	}

	reopenPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 2,
		"indexRows":   -1,
		"quoteEscape": "backslash",
	})
	reopenID := reopenPayload["documentId"].(string)
	reopenWindow := getJSON(t, handler, "/api/window?id="+reopenID+"&start=0&count=2")
	reopenRows := reopenWindow["rows"].([]any)
	reopenRecord := reopenRows[1].([]any)
	if reopenRecord[1] != `has "quote"` {
		t.Fatalf("reopen readback cell = %#v, want quote-preserved value", reopenRecord[1])
	}
}

func TestServerOpenLineEndingOverrideSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "line-ending-open.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
		"lineEnding":  "crlf",
	})
	documentID := openPayload["documentId"].(string)
	metadata := openPayload["document"].(map[string]any)["metadata"].(map[string]any)
	lineEnding := metadata["lineEnding"].(map[string]any)
	if lineEnding["name"] != "crlf" {
		t.Fatalf("line ending = %#v, want crlf", lineEnding)
	}

	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(saved, []byte("id,name\r\n1,Ada\r\n")) {
		t.Fatalf("saved line ending override = %q", saved)
	}
}

func TestServerLineEndingPolicySave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "line-ending-policy.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	policyPayload := postJSON(t, handler, "/api/line-ending", map[string]any{
		"documentId": documentID,
		"lineEnding": "crlf",
	})
	if policyPayload["dirty"] != true {
		t.Fatalf("dirty after line ending policy = %#v, want true", policyPayload["dirty"])
	}
	document := policyPayload["session"].(map[string]any)["document"].(map[string]any)
	metadata := document["metadata"].(map[string]any)
	lineEnding := metadata["lineEnding"].(map[string]any)
	if lineEnding["name"] != "crlf" {
		t.Fatalf("line ending = %#v, want crlf", lineEnding)
	}

	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(saved, []byte("id,name\r\n1,Ada\r\n")) {
		t.Fatalf("saved line ending policy = %q", saved)
	}
}

func TestServerDelimiterPolicySave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "delimiter-policy.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	delimiterPayload := postJSON(t, handler, "/api/delimiter", map[string]any{
		"documentId":    documentID,
		"delimiter":     "|",
		"delimiterName": "pipe",
	})
	if delimiterPayload["dirty"] != true {
		t.Fatalf("dirty after delimiter policy = %#v, want true", delimiterPayload["dirty"])
	}
	document := delimiterPayload["session"].(map[string]any)["document"].(map[string]any)
	dialect := document["metadata"].(map[string]any)["dialect"].(map[string]any)
	if dialect["delimiter"].(float64) != float64('|') || dialect["delimiterName"] != "pipe" {
		t.Fatalf("delimiter dialect = %#v, want pipe", dialect)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada" {
		t.Fatalf("rows after delimiter materialize = %#v", rows)
	}

	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(saved, []byte("id|name\n1|Ada\n2|Grace\n")) {
		t.Fatalf("saved delimiter policy = %q", saved)
	}
}

func TestServerEncodingPolicySave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "encoding-policy.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	encodingPayload := postJSON(t, handler, "/api/encoding", map[string]any{
		"documentId":   documentID,
		"encodingName": "utf-8-bom",
	})
	if encodingPayload["dirty"] != true {
		t.Fatalf("dirty after encoding policy = %#v, want true", encodingPayload["dirty"])
	}
	document := encodingPayload["session"].(map[string]any)["document"].(map[string]any)
	encoding := document["metadata"].(map[string]any)["encoding"].(map[string]any)
	if encoding["name"] != "utf-8-bom" || encoding["hasBom"] != true {
		t.Fatalf("encoding metadata = %#v, want utf-8-bom with BOM", encoding)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if rows[2].([]any)[1] != "Grace" {
		t.Fatalf("rows after encoding materialize = %#v", rows)
	}

	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(saved, []byte{0xEF, 0xBB, 0xBF}) {
		t.Fatalf("saved encoding policy missing UTF-8 BOM: % x", saved[:min(len(saved), 8)])
	}
	if !bytes.Equal(saved[3:], []byte("id,name\n1,Ada\n2,Grace\n")) {
		t.Fatalf("saved encoding payload = %q", saved)
	}
}

func TestServerManualCellQuoteSaveAndUndo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manual-quote.csv")
	if err := os.WriteFile(path, []byte("id,note\n1,plain\n2,free\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	_ = postJSON(t, handler, "/api/quote-policy", map[string]any{
		"documentId": documentID,
		"quoteChar":  "double",
		"quoteMode":  "manual",
	})
	cellPayload := postJSON(t, handler, "/api/cell-quote", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        1,
		"quoted":     true,
	})
	document := cellPayload["session"].(map[string]any)["document"].(map[string]any)
	dialect := document["metadata"].(map[string]any)["dialect"].(map[string]any)
	if dialect["quoteMode"] != "manual" || len(dialect["quoteCells"].([]any)) != 1 {
		t.Fatalf("manual cell quote dialect = %#v", dialect)
	}

	backup := false
	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     backup,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "id,note\n1,\"plain\"\n2,free\n"
	if string(saved) != want {
		t.Fatalf("manual quote save = %q, want %q", string(saved), want)
	}

	undoPayload := postJSON(t, handler, "/api/undo", map[string]any{"documentId": documentID})
	undoDocument := undoPayload["session"].(map[string]any)["document"].(map[string]any)
	undoDialect := undoDocument["metadata"].(map[string]any)["dialect"].(map[string]any)
	if _, ok := undoDialect["quoteCells"]; ok {
		t.Fatalf("undo should remove cell quote override: %#v", undoDialect)
	}
}

func TestServerQuoteRulesSaveAndUndo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "quote-rules.csv")
	if err := os.WriteFile(path, []byte("id,note\n1,plain\n2,free\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	_ = postJSON(t, handler, "/api/quote-policy", map[string]any{
		"documentId": documentID,
		"quoteChar":  "double",
		"quoteMode":  "manual",
	})
	_ = postJSON(t, handler, "/api/quote-rule", map[string]any{
		"documentId": documentID,
		"target":     "header",
		"quoted":     true,
	})
	_ = postJSON(t, handler, "/api/quote-rule", map[string]any{
		"documentId": documentID,
		"target":     "column",
		"col":        1,
		"quoted":     true,
	})
	bodyPayload := postJSON(t, handler, "/api/quote-rule", map[string]any{
		"documentId": documentID,
		"target":     "body",
		"quoted":     true,
	})
	document := bodyPayload["session"].(map[string]any)["document"].(map[string]any)
	dialect := document["metadata"].(map[string]any)["dialect"].(map[string]any)
	if dialect["quoteHeader"] != true || dialect["quoteBody"] != true || len(dialect["quoteColumns"].([]any)) != 1 {
		t.Fatalf("quote rules dialect = %#v", dialect)
	}

	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "\"id\",\"note\"\n\"1\",\"plain\"\n\"2\",\"free\"\n"
	if string(saved) != want {
		t.Fatalf("quote rules save = %q, want %q", string(saved), want)
	}

	undoPayload := postJSON(t, handler, "/api/undo", map[string]any{"documentId": documentID})
	undoDocument := undoPayload["session"].(map[string]any)["document"].(map[string]any)
	undoDialect := undoDocument["metadata"].(map[string]any)["dialect"].(map[string]any)
	if undoDialect["quoteBody"] == true {
		t.Fatalf("undo should remove body quote rule: %#v", undoDialect)
	}
}

func TestServerQuoteSettingsImportSaveAndUndo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "quote-settings.csv")
	if err := os.WriteFile(path, []byte("id,note\n1,plain\n2,free\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	settingsPayload := postJSON(t, handler, "/api/quote-settings", map[string]any{
		"documentId":   documentID,
		"quoteChar":    "double",
		"quoteMode":    "manual",
		"quoteHeader":  true,
		"quoteBody":    false,
		"quoteColumns": []int{1, 1, -1},
		"quoteCells": []map[string]int{
			{"row": 1, "col": 0},
		},
		"noQuoteCells": []map[string]int{
			{"row": 2, "col": 1},
		},
	})
	document := settingsPayload["session"].(map[string]any)["document"].(map[string]any)
	dialect := document["metadata"].(map[string]any)["dialect"].(map[string]any)
	if dialect["quoteHeader"] != true || dialect["quoteBody"] == true || len(dialect["quoteColumns"].([]any)) != 1 {
		t.Fatalf("quote settings dialect = %#v", dialect)
	}
	if len(dialect["quoteCells"].([]any)) != 1 || len(dialect["noQuoteCells"].([]any)) != 1 {
		t.Fatalf("quote settings cell maps = %#v", dialect)
	}

	_ = postJSON(t, handler, "/api/save", map[string]any{
		"documentId": documentID,
		"backup":     false,
	})
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "\"id\",\"note\"\n\"1\",\"plain\"\n2,free\n"
	if string(saved) != want {
		t.Fatalf("quote settings save = %q, want %q", string(saved), want)
	}

	undoPayload := postJSON(t, handler, "/api/undo", map[string]any{"documentId": documentID})
	undoDocument := undoPayload["session"].(map[string]any)["document"].(map[string]any)
	undoDialect := undoDocument["metadata"].(map[string]any)["dialect"].(map[string]any)
	if undoDialect["quoteHeader"] == true {
		t.Fatalf("undo should restore quote settings: %#v", undoDialect)
	}
}

func TestServerReadOnlyAndRefresh(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	readOnlyPayload := postJSON(t, handler, "/api/read-only", map[string]any{
		"documentId": documentID,
		"readOnly":   true,
	})
	if readOnlyPayload["readOnly"] != true {
		t.Fatalf("readOnly payload = %#v", readOnlyPayload)
	}

	response := postJSONResponse(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        1,
		"value":      "Grace",
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("read-only edit status = %d, want 400", response.Code)
	}
	response = postJSONResponse(t, handler, "/api/header-mode", map[string]any{
		"documentId": documentID,
		"hasHeader":  false,
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("read-only header-mode status = %d, want 400", response.Code)
	}

	_ = postJSON(t, handler, "/api/read-only", map[string]any{
		"documentId": documentID,
		"readOnly":   false,
	})
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	statusPayload := getJSON(t, handler, "/api/file-status?id="+documentID)
	if statusPayload["changedExternally"] != true {
		t.Fatalf("status payload = %#v, want changedExternally", statusPayload)
	}

	refreshPayload := postJSON(t, handler, "/api/refresh", map[string]any{"documentId": documentID})
	session := refreshPayload["session"].(map[string]any)
	document := session["document"].(map[string]any)
	rowIndex := document["rowIndex"].(map[string]any)
	if rowIndex["rows"].(float64) != 3 {
		t.Fatalf("refresh rowIndex = %#v, want 3 rows", rowIndex)
	}

	_ = postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        1,
		"col":        1,
		"value":      "Local",
	})
	if err := os.WriteFile(path, []byte("id,name\n1,Disk\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	response = postJSONResponse(t, handler, "/api/refresh", map[string]any{"documentId": documentID})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("dirty refresh status = %d, want 400", response.Code)
	}
	refreshPayload = postJSON(t, handler, "/api/refresh", map[string]any{
		"documentId":   documentID,
		"discardDirty": true,
	})
	session = refreshPayload["session"].(map[string]any)
	if session["dirty"] == true {
		t.Fatalf("discard refresh session = %#v, want clean", session)
	}
	document = session["document"].(map[string]any)
	preview := document["preview"].([]any)
	row := preview[1].([]any)
	if row[1] != "Disk" {
		t.Fatalf("discard refresh preview row = %#v, want disk value", row)
	}
}

func TestServerUndoRedoAndStructuralMutations(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	insertPayload := postJSON(t, handler, "/api/row/insert", map[string]any{
		"documentId": documentID,
		"row":        2,
	})
	if insertPayload["totalRows"].(float64) != 4 {
		t.Fatalf("insert totalRows = %#v, want 4", insertPayload["totalRows"])
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 4 {
		t.Fatalf("rows after insert = %#v, want 4 rows", rows)
	}

	undoPayload := postJSON(t, handler, "/api/undo", map[string]any{"documentId": documentID})
	if undoPayload["totalRows"].(float64) != 3 {
		t.Fatalf("undo totalRows = %#v, want 3", undoPayload["totalRows"])
	}
	redoPayload := postJSON(t, handler, "/api/redo", map[string]any{"documentId": documentID})
	if redoPayload["totalRows"].(float64) != 4 {
		t.Fatalf("redo totalRows = %#v, want 4", redoPayload["totalRows"])
	}

	columnPayload := postJSON(t, handler, "/api/column/insert", map[string]any{
		"documentId": documentID,
		"col":        1,
	})
	if columnPayload["totalRows"].(float64) != 4 {
		t.Fatalf("insert column totalRows = %#v, want 4", columnPayload["totalRows"])
	}
	deleteColumnPayload := postJSON(t, handler, "/api/column/delete", map[string]any{
		"documentId": documentID,
		"col":        1,
	})
	if deleteColumnPayload["totalRows"].(float64) != 4 {
		t.Fatalf("delete column totalRows = %#v, want 4", deleteColumnPayload["totalRows"])
	}

	transposePayload := postJSON(t, handler, "/api/transpose", map[string]any{
		"documentId": documentID,
		"maxCells":   100,
	})
	if transposePayload["totalRows"].(float64) != 2 {
		t.Fatalf("transpose totalRows = %#v, want 2", transposePayload["totalRows"])
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows = windowPayload["rows"].([]any)
	if len(rows) != 2 || rows[0].([]any)[0] != "id" || rows[1].([]any)[0] != "name" {
		t.Fatalf("rows after transpose = %#v", rows)
	}
}

func TestServerMoveRowsAndColumns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "move.csv")
	if err := os.WriteFile(path, []byte("id,name,score\n1,Ada,10\n2,Linus,20\n3,Grace,30\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	rowPayload := postJSON(t, handler, "/api/row/move", map[string]any{
		"documentId": documentID,
		"row":        2,
		"delta":      -1,
	})
	if rowPayload["changed"].(float64) != 1 {
		t.Fatalf("move row payload = %#v", rowPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Linus" || rows[2].([]any)[1] != "Ada" {
		t.Fatalf("rows after move row = %#v", rows)
	}

	columnPayload := postJSON(t, handler, "/api/column/move", map[string]any{
		"documentId": documentID,
		"col":        2,
		"delta":      -1,
	})
	if columnPayload["changed"].(float64) != 1 {
		t.Fatalf("move column payload = %#v", columnPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows = windowPayload["rows"].([]any)
	if rows[0].([]any)[1] != "score" || rows[1].([]any)[1] != "20" {
		t.Fatalf("rows after move column = %#v", rows)
	}
}

func TestServerPasteRowsAndColumns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "paste.csv")
	if err := os.WriteFile(path, []byte("id,name,score\n1,Ada,10\n2,Linus,20\n3,Grace,30\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	rowCopyPayload := postJSON(t, handler, "/api/row/paste", map[string]any{
		"documentId": documentID,
		"sourceRow":  2,
		"targetRow":  1,
		"cut":        false,
	})
	if rowCopyPayload["changed"].(float64) != 1 || rowCopyPayload["totalRows"].(float64) != 5 {
		t.Fatalf("row copy paste payload = %#v", rowCopyPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=5")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Linus" || rows[2].([]any)[1] != "Ada" || rows[3].([]any)[1] != "Linus" {
		t.Fatalf("rows after copy paste = %#v", rows)
	}

	rowCutPayload := postJSON(t, handler, "/api/row/paste", map[string]any{
		"documentId": documentID,
		"sourceRow":  4,
		"targetRow":  2,
		"cut":        true,
	})
	if rowCutPayload["changed"].(float64) != 1 || rowCutPayload["totalRows"].(float64) != 5 {
		t.Fatalf("row cut paste payload = %#v", rowCutPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=5")
	rows = windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Linus" || rows[2].([]any)[1] != "Grace" || rows[3].([]any)[1] != "Ada" {
		t.Fatalf("rows after cut paste = %#v", rows)
	}

	columnCopyPayload := postJSON(t, handler, "/api/column/paste", map[string]any{
		"documentId": documentID,
		"sourceCol":  2,
		"targetCol":  1,
		"cut":        false,
	})
	if columnCopyPayload["changed"].(float64) != 1 {
		t.Fatalf("column copy paste payload = %#v", columnCopyPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=5")
	rows = windowPayload["rows"].([]any)
	if rows[0].([]any)[1] != "score" || rows[0].([]any)[2] != "name" || rows[1].([]any)[1] != "20" || rows[1].([]any)[2] != "Linus" {
		t.Fatalf("rows after column copy paste = %#v", rows)
	}

	columnCutPayload := postJSON(t, handler, "/api/column/paste", map[string]any{
		"documentId": documentID,
		"sourceCol":  2,
		"targetCol":  1,
		"cut":        true,
	})
	if columnCutPayload["changed"].(float64) != 1 {
		t.Fatalf("column cut paste payload = %#v", columnCutPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=5")
	rows = windowPayload["rows"].([]any)
	if rows[0].([]any)[1] != "name" || rows[0].([]any)[2] != "score" || rows[1].([]any)[1] != "Linus" || rows[1].([]any)[2] != "20" {
		t.Fatalf("rows after column cut paste = %#v", rows)
	}
}

func TestServerPasteMultipleRowsAndColumns(t *testing.T) {
	dir := t.TempDir()
	rowPath := filepath.Join(dir, "paste-rows.csv")
	if err := os.WriteFile(rowPath, []byte("id,name,score\n1,Ada,10\n2,Linus,20\n3,Grace,30\n4,Ken,40\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        rowPath,
		"previewRows": 8,
		"indexRows":   -1,
	})
	rowDocumentID := openPayload["documentId"].(string)

	rowCopyPayload := postJSON(t, handler, "/api/row/paste", map[string]any{
		"documentId": rowDocumentID,
		"sourceRows": []int{2, 3},
		"targetRow":  1,
		"cut":        false,
	})
	if rowCopyPayload["changed"].(float64) != 2 || rowCopyPayload["totalRows"].(float64) != 7 {
		t.Fatalf("multi-row copy paste payload = %#v", rowCopyPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+rowDocumentID+"&start=0&count=7")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Linus" || rows[2].([]any)[1] != "Grace" || rows[3].([]any)[1] != "Ada" {
		t.Fatalf("rows after multi-row copy paste = %#v", rows)
	}

	rowCutPayload := postJSON(t, handler, "/api/row/paste", map[string]any{
		"documentId": rowDocumentID,
		"sourceRows": []int{5, 6},
		"targetRow":  3,
		"cut":        true,
	})
	if rowCutPayload["changed"].(float64) != 2 || rowCutPayload["totalRows"].(float64) != 7 {
		t.Fatalf("multi-row cut paste payload = %#v", rowCutPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+rowDocumentID+"&start=0&count=7")
	rows = windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Linus" || rows[2].([]any)[1] != "Grace" || rows[3].([]any)[1] != "Grace" || rows[4].([]any)[1] != "Ken" {
		t.Fatalf("rows after multi-row cut paste = %#v", rows)
	}

	columnPath := filepath.Join(dir, "paste-columns.csv")
	if err := os.WriteFile(columnPath, []byte("id,name,score,city\n1,Ada,10,UK\n2,Linus,20,FI\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	openPayload = postJSON(t, handler, "/api/open", map[string]any{
		"path":        columnPath,
		"previewRows": 4,
		"indexRows":   -1,
	})
	columnDocumentID := openPayload["documentId"].(string)

	columnCutPayload := postJSON(t, handler, "/api/column/paste", map[string]any{
		"documentId": columnDocumentID,
		"sourceCols": []int{2, 3},
		"targetCol":  1,
		"cut":        true,
	})
	if columnCutPayload["changed"].(float64) != 2 {
		t.Fatalf("multi-column cut paste payload = %#v", columnCutPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+columnDocumentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if rows[0].([]any)[1] != "score" || rows[0].([]any)[2] != "city" || rows[0].([]any)[3] != "name" || rows[1].([]any)[1] != "10" || rows[1].([]any)[2] != "UK" {
		t.Fatalf("rows after multi-column cut paste = %#v", rows)
	}

	columnCopyPayload := postJSON(t, handler, "/api/column/paste", map[string]any{
		"documentId": columnDocumentID,
		"sourceCols": []int{1, 2},
		"targetCol":  4,
		"cut":        false,
	})
	if columnCopyPayload["changed"].(float64) != 2 {
		t.Fatalf("multi-column copy paste payload = %#v", columnCopyPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+columnDocumentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if len(rows[0].([]any)) != 6 || rows[0].([]any)[4] != "score" || rows[0].([]any)[5] != "city" || rows[1].([]any)[4] != "10" || rows[1].([]any)[5] != "UK" {
		t.Fatalf("rows after multi-column copy paste = %#v", rows)
	}
}

func TestServerPasteRowsBelowAndColumnsRightAtSelectionEnd(t *testing.T) {
	dir := t.TempDir()
	rowPath := filepath.Join(dir, "paste-below.csv")
	if err := os.WriteFile(rowPath, []byte("id,name\n1,Ada\n2,Linus\n3,Grace\n4,Ken\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        rowPath,
		"previewRows": 8,
		"indexRows":   -1,
	})
	rowDocumentID := openPayload["documentId"].(string)

	rowCopyPayload := postJSON(t, handler, "/api/row/paste", map[string]any{
		"documentId": rowDocumentID,
		"sourceRows": []int{1, 2},
		"targetRow":  5,
		"cut":        false,
	})
	if rowCopyPayload["changed"].(float64) != 2 || rowCopyPayload["totalRows"].(float64) != 7 {
		t.Fatalf("paste rows below payload = %#v", rowCopyPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+rowDocumentID+"&start=0&count=7")
	rows := windowPayload["rows"].([]any)
	if rows[4].([]any)[1] != "Ken" || rows[5].([]any)[1] != "Ada" || rows[6].([]any)[1] != "Linus" {
		t.Fatalf("rows after paste below = %#v", rows)
	}

	rowNoopPayload := postJSON(t, handler, "/api/row/paste", map[string]any{
		"documentId": rowDocumentID,
		"sourceRows": []int{5, 6},
		"targetRow":  7,
		"cut":        true,
	})
	rowNoopChanged, _ := rowNoopPayload["changed"].(float64)
	if rowNoopChanged != 0 || rowNoopPayload["totalRows"].(float64) != 7 {
		t.Fatalf("self paste rows below payload = %#v", rowNoopPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+rowDocumentID+"&start=0&count=7")
	rows = windowPayload["rows"].([]any)
	if len(rows) != 7 || rows[5].([]any)[1] != "Ada" || rows[6].([]any)[1] != "Linus" {
		t.Fatalf("rows after self paste below = %#v", rows)
	}

	columnPath := filepath.Join(dir, "paste-right.csv")
	if err := os.WriteFile(columnPath, []byte("id,name,score,city\n1,Ada,10,UK\n2,Linus,20,FI\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	openPayload = postJSON(t, handler, "/api/open", map[string]any{
		"path":        columnPath,
		"previewRows": 4,
		"indexRows":   -1,
	})
	columnDocumentID := openPayload["documentId"].(string)

	columnCopyPayload := postJSON(t, handler, "/api/column/paste", map[string]any{
		"documentId": columnDocumentID,
		"sourceCols": []int{1, 2},
		"targetCol":  4,
		"cut":        false,
	})
	if columnCopyPayload["changed"].(float64) != 2 {
		t.Fatalf("paste columns right payload = %#v", columnCopyPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+columnDocumentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if len(rows[0].([]any)) != 6 || rows[0].([]any)[4] != "name" || rows[0].([]any)[5] != "score" || rows[1].([]any)[4] != "Ada" || rows[1].([]any)[5] != "10" {
		t.Fatalf("rows after paste columns right = %#v", rows)
	}

	columnNoopPayload := postJSON(t, handler, "/api/column/paste", map[string]any{
		"documentId": columnDocumentID,
		"sourceCols": []int{4, 5},
		"targetCol":  6,
		"cut":        true,
	})
	columnNoopChanged, _ := columnNoopPayload["changed"].(float64)
	if columnNoopChanged != 0 {
		t.Fatalf("self paste columns right payload = %#v", columnNoopPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+columnDocumentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if len(rows[0].([]any)) != 6 || rows[0].([]any)[4] != "name" || rows[0].([]any)[5] != "score" {
		t.Fatalf("rows after self paste columns right = %#v", rows)
	}
}

func TestServerPasteBelowRightPreservesSelectionsIncludingLastItem(t *testing.T) {
	dir := t.TempDir()
	rowPath := filepath.Join(dir, "paste-last-row.csv")
	if err := os.WriteFile(rowPath, []byte("id,name\n1,Ada\n2,Linus\n3,Grace\n4,Ken\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        rowPath,
		"previewRows": 8,
		"indexRows":   -1,
	})
	rowDocumentID := openPayload["documentId"].(string)

	rowCopyPayload := postJSON(t, handler, "/api/row/paste", map[string]any{
		"documentId": rowDocumentID,
		"sourceRows": []int{3, 4},
		"targetRow":  5,
		"cut":        false,
	})
	if rowCopyPayload["changed"].(float64) != 2 || rowCopyPayload["totalRows"].(float64) != 7 {
		t.Fatalf("paste rows including last payload = %#v", rowCopyPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+rowDocumentID+"&start=0&count=7")
	rows := windowPayload["rows"].([]any)
	gotNames := make([]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		gotNames = append(gotNames, row.([]any)[1].(string))
	}
	wantNames := []string{"Ada", "Linus", "Grace", "Ken", "Grace", "Ken"}
	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("rows after paste including last = %#v, want %#v", gotNames, wantNames)
	}

	columnPath := filepath.Join(dir, "paste-last-column.csv")
	if err := os.WriteFile(columnPath, []byte("id,name,score,city\n1,Ada,10,UK\n2,Linus,20,FI\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	openPayload = postJSON(t, handler, "/api/open", map[string]any{
		"path":        columnPath,
		"previewRows": 4,
		"indexRows":   -1,
	})
	columnDocumentID := openPayload["documentId"].(string)

	columnCopyPayload := postJSON(t, handler, "/api/column/paste", map[string]any{
		"documentId": columnDocumentID,
		"sourceCols": []int{2, 3},
		"targetCol":  4,
		"cut":        false,
	})
	if columnCopyPayload["changed"].(float64) != 2 {
		t.Fatalf("paste columns including last payload = %#v", columnCopyPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+columnDocumentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	gotHeader := rows[0].([]any)
	gotFirst := rows[1].([]any)
	if !reflect.DeepEqual(gotHeader, []any{"id", "name", "score", "city", "score", "city"}) ||
		!reflect.DeepEqual(gotFirst, []any{"1", "Ada", "10", "UK", "10", "UK"}) {
		t.Fatalf("columns after paste including last header=%#v first=%#v", gotHeader, gotFirst)
	}
}

func TestServerHideRowAndClearView(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hide.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n3,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	hidePayload := postJSON(t, handler, "/api/row/hide", map[string]any{
		"documentId": documentID,
		"row":        2,
	})
	if hidePayload["changed"].(float64) != 1 || hidePayload["totalRows"].(float64) != 2 {
		t.Fatalf("hide row payload = %#v", hidePayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 2 || rows[0].([]any)[1] != "Ada" || rows[1].([]any)[1] != "Grace" {
		t.Fatalf("rows after hide = %#v", rows)
	}
	clearPayload := postJSON(t, handler, "/api/view-clear", map[string]any{"documentId": documentID})
	if clearPayload["totalRows"].(float64) != 4 {
		t.Fatalf("clear payload = %#v", clearPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows = windowPayload["rows"].([]any)
	if len(rows) != 4 || rows[2].([]any)[1] != "Linus" {
		t.Fatalf("rows after clear hide = %#v", rows)
	}
}

func TestServerDedupeAndTextTransform(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tools.csv")
	if err := os.WriteFile(path, []byte("id,name,name\n1,  Ada ,Ada\n1,  Ada ,Ada\n2,Linus,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	dedupePayload := postJSON(t, handler, "/api/dedupe-rows", map[string]any{"documentId": documentID})
	if dedupePayload["changed"].(float64) != 1 || dedupePayload["totalRows"].(float64) != 3 {
		t.Fatalf("dedupe payload = %#v", dedupePayload)
	}
	transformPayload := postJSON(t, handler, "/api/text-transform", map[string]any{
		"documentId": documentID,
		"mode":       "trim",
		"cells": []map[string]any{
			{"row": 1, "col": 1},
			{"row": 2, "col": 1},
		},
	})
	if transformPayload["changed"].(float64) != 1 {
		t.Fatalf("transform payload = %#v", transformPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada" || rows[2].([]any)[1] != "Linus" {
		t.Fatalf("rows after transform = %#v", rows)
	}
	transformPayload = postJSON(t, handler, "/api/text-transform", map[string]any{
		"documentId": documentID,
		"mode":       "title",
		"cells": []map[string]any{
			{"row": 0, "col": 0},
		},
	})
	if transformPayload["changed"].(float64) != 1 {
		t.Fatalf("title transform payload = %#v", transformPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if rows[0].([]any)[0] != "Id" {
		t.Fatalf("rows after title transform = %#v", rows)
	}
	columnPayload := postJSON(t, handler, "/api/dedupe-columns", map[string]any{"documentId": documentID})
	if columnPayload["changed"].(float64) != 1 {
		t.Fatalf("dedupe column payload = %#v", columnPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if len(rows[0].([]any)) != 2 || rows[1].([]any)[1] != "Ada" {
		t.Fatalf("rows after column dedupe = %#v", rows)
	}
}

func TestServerTextTransformColumn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "column-transform.csv")
	if err := os.WriteFile(path, []byte("id,name,note\n1,Ada, keep \n2,Grace\n3,Linus, later \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	columnPayload := postJSON(t, handler, "/api/text-transform", map[string]any{
		"documentId": documentID,
		"mode":       "upper",
		"column":     1,
	})
	if columnPayload["changed"].(float64) != 3 {
		t.Fatalf("column transform payload = %#v", columnPayload)
	}
	if columnPayload["streamed"] != true {
		t.Fatalf("column transform streamed = %#v, want true", columnPayload["streamed"])
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows := windowPayload["rows"].([]any)
	if rows[0].([]any)[1] != "name" || rows[1].([]any)[1] != "ADA" || rows[2].([]any)[1] != "GRACE" || rows[3].([]any)[1] != "LINUS" {
		t.Fatalf("rows after column transform = %#v", rows)
	}
	if len(rows[2].([]any)) != 2 {
		t.Fatalf("column transform should not materialize missing ragged cells: %#v", rows[2])
	}
}

func TestServerWidthTextTransforms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "width.csv")
	if err := os.WriteFile(path, []byte("id,value\n1,ABC 123 ｶﾀｶﾅ\n2,ＡＢＣ　１２３　カタカナ\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	widePayload := postJSON(t, handler, "/api/text-transform", map[string]any{
		"documentId": documentID,
		"mode":       "wide",
		"cells": []map[string]any{
			{"row": 1, "col": 1},
		},
	})
	if widePayload["changed"].(float64) != 1 {
		t.Fatalf("wide payload = %#v", widePayload)
	}
	narrowPayload := postJSON(t, handler, "/api/text-transform", map[string]any{
		"documentId": documentID,
		"mode":       "narrow",
		"cells": []map[string]any{
			{"row": 2, "col": 1},
		},
	})
	if narrowPayload["changed"].(float64) != 1 {
		t.Fatalf("narrow payload = %#v", narrowPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "ＡＢＣ　１２３　カタカナ" || rows[2].([]any)[1] != "ABC 123 ｶﾀｶﾅ" {
		t.Fatalf("rows after width transforms = %#v", rows)
	}
}

func TestServerJoinAndSplitCells(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "join-split.csv")
	if err := os.WriteFile(path, []byte("id,first,last\n1,Ada,Lovelace\n2,Grace,Hopper\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	joinPayload := postJSON(t, handler, "/api/cells/join", map[string]any{
		"documentId": documentID,
		"separator":  " ",
		"cells": []map[string]any{
			{"row": 1, "col": 1},
			{"row": 1, "col": 2},
		},
	})
	if joinPayload["changed"].(float64) != 2 {
		t.Fatalf("join payload = %#v", joinPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada Lovelace" || rows[1].([]any)[2] != "" {
		t.Fatalf("rows after join = %#v", rows)
	}
	splitPayload := postJSON(t, handler, "/api/cells/split", map[string]any{
		"documentId": documentID,
		"separator":  " ",
		"cells": []map[string]any{
			{"row": 1, "col": 1},
		},
	})
	if splitPayload["changed"].(float64) != 2 {
		t.Fatalf("split payload = %#v", splitPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada" || rows[1].([]any)[2] != "Lovelace" {
		t.Fatalf("rows after split = %#v", rows)
	}

	swapPayload := postJSON(t, handler, "/api/cells/swap", map[string]any{
		"documentId": documentID,
		"cells": []map[string]any{
			{"row": 1, "col": 1},
			{"row": 1, "col": 2},
		},
	})
	if swapPayload["changed"].(float64) != 2 {
		t.Fatalf("swap payload = %#v", swapPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Lovelace" || rows[1].([]any)[2] != "Ada" {
		t.Fatalf("rows after swap = %#v", rows)
	}
}

func TestServerSplitCellsRejectsOverwriteUnlessAllowed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "split-conflict.csv")
	if err := os.WriteFile(path, []byte("id,name,target\n1,Ada Lovelace,occupied\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	response := postJSONResponse(t, handler, "/api/cells/split", map[string]any{
		"documentId": documentID,
		"separator":  " ",
		"cells": []map[string]any{
			{"row": 1, "col": 1},
		},
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("split conflict status = %d, body = %s", response.Code, response.Body.String())
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), `split would overwrite row 2, column 3 ("occupied")`) {
		t.Fatalf("split conflict payload = %#v", payload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada Lovelace" || rows[1].([]any)[2] != "occupied" {
		t.Fatalf("rows after rejected split = %#v", rows)
	}

	splitPayload := postJSON(t, handler, "/api/cells/split", map[string]any{
		"documentId":     documentID,
		"separator":      " ",
		"allowOverwrite": true,
		"cells": []map[string]any{
			{"row": 1, "col": 1},
		},
	})
	if splitPayload["changed"].(float64) != 2 {
		t.Fatalf("split overwrite payload = %#v", splitPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada" || rows[1].([]any)[2] != "Lovelace" {
		t.Fatalf("rows after allowed split = %#v", rows)
	}
}

func TestServerVerticalJoinAndSplitCells(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "join-split-vertical.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	joinPayload := postJSON(t, handler, "/api/cells/join", map[string]any{
		"documentId": documentID,
		"separator":  " ",
		"direction":  "vertical",
		"cells": []map[string]any{
			{"row": 1, "col": 1},
			{"row": 2, "col": 1},
		},
	})
	if joinPayload["changed"].(float64) != 2 {
		t.Fatalf("vertical join payload = %#v", joinPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada Grace" || rows[2].([]any)[1] != "" {
		t.Fatalf("rows after vertical join = %#v", rows)
	}

	splitPayload := postJSON(t, handler, "/api/cells/split", map[string]any{
		"documentId": documentID,
		"separator":  " ",
		"direction":  "vertical",
		"cells": []map[string]any{
			{"row": 1, "col": 1},
		},
	})
	if splitPayload["changed"].(float64) != 2 {
		t.Fatalf("vertical split payload = %#v", splitPayload)
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows = windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada" || rows[2].([]any)[1] != "Grace" {
		t.Fatalf("rows after vertical split = %#v", rows)
	}
}

func TestServerFindFilterSort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name,score\n1,Ada,9\n2,Linus,7\n3,Grace,10\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	findPayload := postJSON(t, handler, "/api/find", map[string]any{
		"documentId": documentID,
		"query":      "Grace",
		"limit":      10,
	})
	matches, ok := findPayload["matches"].([]any)
	if !ok || len(matches) != 1 {
		t.Fatalf("matches = %#v, want one match", findPayload["matches"])
	}
	match := matches[0].(map[string]any)
	textRange, ok := match["range"].(map[string]any)
	if !ok || textRange["start"].(float64) != 0 || textRange["end"].(float64) != 5 {
		t.Fatalf("match range = %#v, want 0-5", match["range"])
	}
	textRanges, ok := match["ranges"].([]any)
	if !ok || len(textRanges) != 1 {
		t.Fatalf("match ranges = %#v, want one range", match["ranges"])
	}
	firstRange := textRanges[0].(map[string]any)
	if firstRange["start"].(float64) != 0 || firstRange["end"].(float64) != 5 {
		t.Fatalf("match ranges[0] = %#v, want 0-5", firstRange)
	}
	if findPayload["total"].(float64) != 1 || findPayload["truncated"].(bool) {
		t.Fatalf("find total/truncated = %#v/%#v, want 1/false", findPayload["total"], findPayload["truncated"])
	}
	if findPayload["streamed"] != true {
		t.Fatalf("find streamed = %#v, want true for clean file", findPayload["streamed"])
	}

	findPayload = postJSON(t, handler, "/api/find", map[string]any{
		"documentId": documentID,
		"query":      "i",
		"limit":      1,
		"offset":     1,
	})
	matches, ok = findPayload["matches"].([]any)
	if !ok || len(matches) != 1 {
		t.Fatalf("offset matches = %#v, want one match", findPayload["matches"])
	}
	match = matches[0].(map[string]any)
	if findPayload["offset"].(float64) != 1 || match["row"].(float64) != 2 {
		t.Fatalf("offset find payload = %#v, want second matching row", findPayload)
	}

	previewPayload := postJSON(t, handler, "/api/replace", map[string]any{
		"documentId":  documentID,
		"query":       "Linus",
		"replacement": "Ken",
		"all":         true,
		"preview":     true,
	})
	if previewPayload["preview"] != true || previewPayload["count"].(float64) != 1 {
		t.Fatalf("replace preview payload = %#v, want one-change preview", previewPayload)
	}
	findPayload = postJSON(t, handler, "/api/find", map[string]any{
		"documentId": documentID,
		"query":      "Ken",
		"wholeCell":  true,
		"limit":      10,
	})
	if findPayload["total"].(float64) != 0 {
		t.Fatalf("preview mutated document: %#v", findPayload)
	}

	replacePayload := postJSON(t, handler, "/api/replace", map[string]any{
		"documentId":  documentID,
		"query":       "Linus",
		"replacement": "Ken",
		"all":         true,
	})
	if replacePayload["count"].(float64) != 1 {
		t.Fatalf("replace count = %#v, want 1", replacePayload["count"])
	}
	findPayload = postJSON(t, handler, "/api/find", map[string]any{
		"documentId": documentID,
		"query":      "Ken",
		"wholeCell":  true,
		"limit":      10,
	})
	matches, ok = findPayload["matches"].([]any)
	if !ok || len(matches) != 1 {
		t.Fatalf("post-replace matches = %#v, want one match", findPayload["matches"])
	}

	filterPayload := postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"query":      "Ada",
	})
	if filterPayload["totalRows"].(float64) != 1 {
		t.Fatalf("filter totalRows = %#v, want 1", filterPayload["totalRows"])
	}

	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=2")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 1 {
		t.Fatalf("filtered window rows = %#v, want one row", rows)
	}

	sqlFilterPayload := postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"mode":       "sql",
		"expression": "score >= 9 AND LOWER(name) LIKE 'g%'",
	})
	if sqlFilterPayload["totalRows"].(float64) != 1 {
		t.Fatalf("SQL filter totalRows = %#v, want 1", sqlFilterPayload["totalRows"])
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=2")
	rows = windowPayload["rows"].([]any)
	if len(rows) != 1 {
		t.Fatalf("SQL filtered window rows = %#v, want one row", rows)
	}

	valuesPayload := getJSON(t, handler, "/api/column-values?id="+documentID+"&column=1&limit=10")
	values := valuesPayload["values"].([]any)
	if len(values) != 3 {
		t.Fatalf("column values = %#v, want 3", valuesPayload)
	}
	valueFilterPayload := postJSON(t, handler, "/api/filter-values", map[string]any{
		"documentId": documentID,
		"column":     1,
		"values":     []string{"Ada", "Grace"},
	})
	if valueFilterPayload["totalRows"].(float64) != 2 {
		t.Fatalf("value filter totalRows = %#v, want 2", valueFilterPayload["totalRows"])
	}

	sortPayload := postJSON(t, handler, "/api/sort", map[string]any{
		"documentId": documentID,
		"column":     2,
		"numeric":    true,
		"descending": true,
		"hasHeader":  true,
		"stable":     true,
	})
	if sortPayload["totalRows"].(float64) != 2 {
		t.Fatalf("filtered sort totalRows = %#v, want 2", sortPayload["totalRows"])
	}
	if sortPayload["stable"] != true {
		t.Fatalf("filtered sort stable = %#v, want true", sortPayload["stable"])
	}
	unstableSortPayload := postJSON(t, handler, "/api/sort", map[string]any{
		"documentId": documentID,
		"column":     2,
		"numeric":    true,
		"descending": false,
		"hasHeader":  true,
		"stable":     false,
	})
	if unstableSortPayload["stable"] != false {
		t.Fatalf("unstable sort stable = %#v, want false", unstableSortPayload["stable"])
	}

	clearPayload := postJSON(t, handler, "/api/view-clear", map[string]any{
		"documentId": documentID,
	})
	if clearPayload["totalRows"].(float64) != 4 {
		t.Fatalf("clear view totalRows = %#v, want 4", clearPayload["totalRows"])
	}
	sortPayload = postJSON(t, handler, "/api/sort", map[string]any{
		"documentId": documentID,
		"column":     2,
		"numeric":    true,
		"descending": true,
		"hasHeader":  true,
	})
	if sortPayload["totalRows"].(float64) != 3 {
		t.Fatalf("full sort totalRows = %#v, want 3", sortPayload["totalRows"])
	}

	sqlPayload := postJSON(t, handler, "/api/sql", map[string]any{
		"documentId": documentID,
		"query":      "SELECT name, score FROM current WHERE CAST(score AS REAL) >= 9 ORDER BY CAST(score AS REAL) DESC LIMIT 1",
	})
	if sqlPayload["totalRows"].(float64) != 2 {
		t.Fatalf("sql totalRows = %#v, want 2", sqlPayload["totalRows"])
	}
	sqlRows := sqlPayload["rows"].([]any)
	if len(sqlRows) != 1 {
		t.Fatalf("sql rows = %#v, want one row", sqlRows)
	}
	firstSQLRow := sqlRows[0].([]any)
	if firstSQLRow[0] != "Grace" || firstSQLRow[1] != "10" {
		t.Fatalf("unexpected SQL row: %#v", firstSQLRow)
	}
	aggregatePayload := postJSON(t, handler, "/api/sql", map[string]any{
		"documentId": documentID,
		"query":      "SELECT COUNT(*) AS rows, SUM(score) AS total FROM current",
	})
	aggregateRows := aggregatePayload["rows"].([]any)
	aggregateRow := aggregateRows[0].([]any)
	if aggregateRow[0] != "3" || aggregateRow[1] != "26" {
		t.Fatalf("unexpected aggregate row: %#v", aggregateRow)
	}
	sqlErrorResponse := postJSONResponse(t, handler, "/api/sql", map[string]any{
		"documentId": documentID,
		"query":      "SELECT FROM current",
	})
	if sqlErrorResponse.Code != http.StatusBadRequest {
		t.Fatalf("sql error status = %d, want %d", sqlErrorResponse.Code, http.StatusBadRequest)
	}
	var sqlErrorPayload map[string]any
	if err := json.NewDecoder(sqlErrorResponse.Body).Decode(&sqlErrorPayload); err != nil {
		t.Fatal(err)
	}
	sqlErrorDetail, ok := sqlErrorPayload["sqlError"].(map[string]any)
	if !ok {
		t.Fatalf("missing sqlError detail in %#v", sqlErrorPayload)
	}
	if sqlErrorDetail["line"] != float64(1) || sqlErrorDetail["column"] != float64(8) {
		t.Fatalf("sql error location = %#v:%#v, want 1:8", sqlErrorDetail["line"], sqlErrorDetail["column"])
	}
	if sqlErrorDetail["token"] != "FROM" || sqlErrorDetail["snippet"] != "SELECT FROM current" {
		t.Fatalf("unexpected sql error detail: %#v", sqlErrorDetail)
	}
	if sqlErrorDetail["caret"] != "       ^" {
		t.Fatalf("sql error caret = %q, want marker under FROM", sqlErrorDetail["caret"])
	}
	repeatedTokenErrorResponse := postJSONResponse(t, handler, "/api/sql", map[string]any{
		"documentId": documentID,
		"query":      "SELECT name\nFROM current\nWHERE name = 'FROM'\nFROM current",
	})
	if repeatedTokenErrorResponse.Code != http.StatusBadRequest {
		t.Fatalf("repeated-token sql error status = %d, want %d", repeatedTokenErrorResponse.Code, http.StatusBadRequest)
	}
	var repeatedTokenPayload map[string]any
	if err := json.NewDecoder(repeatedTokenErrorResponse.Body).Decode(&repeatedTokenPayload); err != nil {
		t.Fatal(err)
	}
	repeatedTokenDetail, ok := repeatedTokenPayload["sqlError"].(map[string]any)
	if !ok {
		t.Fatalf("missing repeated-token sqlError detail in %#v", repeatedTokenPayload)
	}
	if repeatedTokenDetail["line"] != float64(4) || repeatedTokenDetail["column"] != float64(1) {
		t.Fatalf("repeated-token sql error location = %#v:%#v, want 4:1", repeatedTokenDetail["line"], repeatedTokenDetail["column"])
	}
	if repeatedTokenDetail["token"] != "FROM" || repeatedTokenDetail["snippet"] != "FROM current" {
		t.Fatalf("unexpected repeated-token sql error detail: %#v", repeatedTokenDetail)
	}
	if repeatedTokenDetail["caret"] != "^" {
		t.Fatalf("repeated-token sql error caret = %q, want marker under FROM", repeatedTokenDetail["caret"])
	}

	profilePayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
	})
	columns := profilePayload["columns"].([]any)
	if len(columns) != 3 {
		t.Fatalf("profile columns = %#v, want 3", columns)
	}
	score := columns[2].(map[string]any)
	if score["type"] != "integer" {
		t.Fatalf("score type = %#v, want integer", score["type"])
	}
}

func TestServerFindUsesRequestContextCancellation(t *testing.T) {
	body, err := json.Marshal(map[string]any{
		"query": "Ada",
		"limit": 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodPost, "/api/find", bytes.NewReader(body)).WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	New(t.TempDir()).Handler().ServeHTTP(response, request)

	if response.Code != statusClientClosedRequest {
		t.Fatalf("canceled find status = %d, want %d", response.Code, statusClientClosedRequest)
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), "context canceled") {
		t.Fatalf("canceled find payload = %#v, want context canceled", payload)
	}
}

func TestServerFilterContainsStreamsCleanFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "filter-stream.csv")
	if err := os.WriteFile(path, []byte("id,region\n1,APAC\n2,EU\n3,APAC\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 2,
		"indexRows":   0,
	})
	documentID := openPayload["documentId"].(string)

	filterPayload := postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"query":      "APAC",
	})
	if filterPayload["totalRows"] != float64(2) || filterPayload["streamed"] != true {
		t.Fatalf("filter payload = %#v, want streamed two-row match", filterPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 2 || rows[0].([]any)[1] != "APAC" || rows[1].([]any)[0] != "3" {
		t.Fatalf("streamed filter window = %#v", rows)
	}
}

func TestServerFilterSQLStreamsCleanFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sql-filter-stream.csv")
	if err := os.WriteFile(path, []byte("id,region,score\n1,APAC,8\n2,EU,9\n3,APAC,10\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 2,
		"indexRows":   0,
	})
	documentID := openPayload["documentId"].(string)

	filterPayload := postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"mode":       "sql",
		"expression": "region = 'APAC' AND CAST(score AS number) >= 9",
	})
	if filterPayload["totalRows"] != float64(1) || filterPayload["streamed"] != true {
		t.Fatalf("SQL filter payload = %#v, want one streamed match", filterPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 1 || rows[0].([]any)[0] != "3" || rows[0].([]any)[1] != "APAC" {
		t.Fatalf("streamed SQL filter window = %#v", rows)
	}
}

func TestServerSQLConsoleStreamsCleanCurrentFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sql-console-stream.csv")
	if err := os.WriteFile(path, []byte("id,region,score\n1,APAC,8\n2,EU,9\n3,APAC,10\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 2,
		"indexRows":   0,
	})
	documentID := openPayload["documentId"].(string)

	sqlPayload := postJSON(t, handler, "/api/sql", map[string]any{
		"documentId": documentID,
		"query":      "SELECT region, COUNT(*) AS rows FROM current GROUP BY region ORDER BY rows DESC",
		"limit":      10,
	})
	if sqlPayload["streamed"] != true || sqlPayload["totalRows"] != float64(2) {
		t.Fatalf("SQL payload = %#v, want streamed two groups", sqlPayload)
	}
	rows := sqlPayload["rows"].([]any)
	if len(rows) != 2 || rows[0].([]any)[0] != "APAC" || rows[0].([]any)[1] != "2" {
		t.Fatalf("streamed SQL rows = %#v, want APAC count first", rows)
	}
}

func TestServerSQLStreamsRowProgressEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sql-progress.csv")
	if err := os.WriteFile(path, []byte("id,region\n1,APAC\n2,EU\n3,APAC\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 1,
		"indexRows":   0,
	})
	response := postJSONResponse(t, handler, "/api/sql", map[string]any{
		"documentId":       openPayload["documentId"].(string),
		"query":            "SELECT region, COUNT(*) AS rows FROM current GROUP BY region ORDER BY rows DESC",
		"limit":            10,
		"streamProgress":   true,
		"progressInterval": 2,
	})
	if response.Code != http.StatusOK {
		t.Fatalf("SQL stream status = %d body=%s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/x-ndjson") {
		t.Fatalf("SQL stream content type = %q", got)
	}
	var events []map[string]any
	decoder := json.NewDecoder(response.Body)
	for {
		var event map[string]any
		if err := decoder.Decode(&event); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}
		events = append(events, event)
	}
	if len(events) != 3 {
		t.Fatalf("events = %#v, want two progress events plus final result", events)
	}
	if events[0]["type"] != "progress" || events[0]["phase"] != "load" || events[0]["rows"] != float64(2) || events[0]["done"] == true {
		t.Fatalf("first event = %#v, want loading checkpoint at two rows", events[0])
	}
	if events[1]["type"] != "progress" || events[1]["phase"] != "load" || events[1]["rows"] != float64(3) || events[1]["done"] != true {
		t.Fatalf("second event = %#v, want loading done at three rows", events[1])
	}
	final := events[len(events)-1]
	result, ok := final["result"].(map[string]any)
	if final["type"] != "result" || !ok {
		t.Fatalf("final event = %#v, want result", final)
	}
	progress := result["progress"].(map[string]any)
	if result["streamed"] != true || progress["loadedRows"] != float64(3) || progress["sources"] != float64(1) {
		t.Fatalf("result progress = %#v, want streamed three-row SQL source load", result)
	}
}

func TestServerValueChecklistStreamsCleanFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "value-filter-stream.csv")
	if err := os.WriteFile(path, []byte("id,region\n1,APAC\n2,EU\n3,\n4,APAC\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 2,
		"indexRows":   0,
	})
	documentID := openPayload["documentId"].(string)

	valuesPayload := getJSON(t, handler, "/api/column-values?id="+documentID+"&column=1&limit=10")
	if valuesPayload["streamed"] != true || valuesPayload["totalDistinct"] != float64(3) {
		t.Fatalf("column values payload = %#v, want streamed three distinct values", valuesPayload)
	}
	values := valuesPayload["values"].([]any)
	if len(values) != 3 || values[0].(map[string]any)["value"] != "APAC" || values[2].(map[string]any)["blank"] != true {
		t.Fatalf("column values = %#v, want APAC first and blank last", values)
	}

	filterPayload := postJSON(t, handler, "/api/filter-values", map[string]any{
		"documentId":   documentID,
		"column":       1,
		"values":       []string{"APAC"},
		"includeBlank": true,
	})
	if filterPayload["totalRows"] != float64(3) || filterPayload["streamed"] != true {
		t.Fatalf("value filter payload = %#v, want streamed APAC plus blank rows", filterPayload)
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 3 || rows[0].([]any)[1] != "APAC" || rows[1].([]any)[1] != "" || rows[2].([]any)[0] != "4" {
		t.Fatalf("streamed value-filter window = %#v", rows)
	}
}

func TestServerFilterUsesRequestContextCancellation(t *testing.T) {
	body, err := json.Marshal(map[string]any{
		"query": "APAC",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodPost, "/api/filter", bytes.NewReader(body)).WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	New(t.TempDir()).Handler().ServeHTTP(response, request)

	if response.Code != statusClientClosedRequest {
		t.Fatalf("canceled filter status = %d, want %d", response.Code, statusClientClosedRequest)
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), "context canceled") {
		t.Fatalf("canceled filter payload = %#v, want context canceled", payload)
	}
}

func TestServerSQLFilterUsesRequestContextCancellation(t *testing.T) {
	body, err := json.Marshal(map[string]any{
		"mode":       "sql",
		"expression": "region = 'APAC'",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodPost, "/api/filter", bytes.NewReader(body)).WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	New(t.TempDir()).Handler().ServeHTTP(response, request)

	if response.Code != statusClientClosedRequest {
		t.Fatalf("canceled SQL filter status = %d, want %d", response.Code, statusClientClosedRequest)
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), "context canceled") {
		t.Fatalf("canceled SQL filter payload = %#v, want context canceled", payload)
	}
}

func TestServerValueFilterUsesRequestContextCancellation(t *testing.T) {
	body, err := json.Marshal(map[string]any{
		"column": 1,
		"values": []string{"APAC"},
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodPost, "/api/filter-values", bytes.NewReader(body)).WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	New(t.TempDir()).Handler().ServeHTTP(response, request)

	if response.Code != statusClientClosedRequest {
		t.Fatalf("canceled value-filter status = %d, want %d", response.Code, statusClientClosedRequest)
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), "context canceled") {
		t.Fatalf("canceled value-filter payload = %#v, want context canceled", payload)
	}
}

func TestServerColumnValuesUsesRequestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodGet, "/api/column-values?id=missing&column=1", nil).WithContext(ctx)
	response := httptest.NewRecorder()

	New(t.TempDir()).Handler().ServeHTTP(response, request)

	if response.Code != statusClientClosedRequest {
		t.Fatalf("canceled column-values status = %d, want %d", response.Code, statusClientClosedRequest)
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), "context canceled") {
		t.Fatalf("canceled column-values payload = %#v, want context canceled", payload)
	}
}

func TestServerReplaceSingleScopedCell(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "replace-current.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Ada\n3,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	replacePayload := postJSON(t, handler, "/api/replace", map[string]any{
		"documentId":  documentID,
		"query":       "Ada",
		"replacement": "Grace",
		"all":         false,
		"cells": []map[string]any{
			{"row": 2, "col": 1},
		},
	})
	if replacePayload["count"].(float64) != 1 {
		t.Fatalf("replace count = %#v, want 1", replacePayload["count"])
	}
	changes := replacePayload["changes"].([]any)
	if len(changes) != 1 {
		t.Fatalf("replace changes = %#v, want one change", changes)
	}
	change := changes[0].(map[string]any)
	if change["row"].(float64) != 2 || change["col"].(float64) != 1 {
		t.Fatalf("replace change = %#v, want row 2 col 1", change)
	}

	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows := windowPayload["rows"].([]any)
	if rows[1].([]any)[1] != "Ada" || rows[2].([]any)[1] != "Grace" || rows[3].([]any)[1] != "Ada" {
		t.Fatalf("rows after scoped replace = %#v", rows)
	}
}

func TestServerProfileAcceptsSchemaRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.csv")
	if err := os.WriteFile(path, []byte("id,status,score\n1,new,10\n1,bad,105\n,done,ok\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	profilePayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
		"schema": map[string]any{
			"columns": []map[string]any{
				{"name": "id", "required": true, "type": "integer", "unique": true},
				{"name": "status", "enum": []string{"new", "done"}},
				{"name": "score", "type": "number", "min": 0, "max": 100},
			},
		},
	})
	if profilePayload["schemaRules"].(float64) != 3 {
		t.Fatalf("schemaRules = %#v, want 3", profilePayload["schemaRules"])
	}
	issues := profilePayload["issues"].([]any)
	for _, typ := range []string{"schema-unique", "schema-required", "schema-enum", "schema-range", "schema-type"} {
		if !jsonIssuesContainType(issues, typ) {
			t.Fatalf("expected %s issue in %#v", typ, issues)
		}
	}
	rangeIssue := jsonIssueByType(issues, "schema-range")
	if rangeIssue == nil {
		t.Fatalf("expected schema-range issue in %#v", issues)
	}
	quickFixes, ok := rangeIssue["quickFixes"].([]any)
	if !ok || len(quickFixes) != 1 {
		t.Fatalf("schema-range quick fixes = %#v, want one fix", rangeIssue["quickFixes"])
	}
	quickFix, ok := quickFixes[0].(map[string]any)
	if !ok || quickFix["id"] != "schema.range.clampMax" || quickFix["value"] != "100" {
		t.Fatalf("schema-range quick fix = %#v, want clampMax to 100", quickFixes[0])
	}
}

func TestServerProfileSingleEnumQuickFix(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "single-enum.csv")
	if err := os.WriteFile(path, []byte("id,status\n1,pending\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	profilePayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
		"schema": map[string]any{
			"columns": []map[string]any{
				{"name": "status", "enum": []string{"active"}},
			},
		},
	})
	issue := jsonIssueByType(profilePayload["issues"].([]any), "schema-enum")
	if issue == nil {
		t.Fatalf("expected schema-enum issue in %#v", profilePayload["issues"])
	}
	quickFixes, ok := issue["quickFixes"].([]any)
	if !ok || len(quickFixes) != 1 {
		t.Fatalf("schema-enum quick fixes = %#v, want one fix", issue["quickFixes"])
	}
	quickFix, ok := quickFixes[0].(map[string]any)
	if !ok || quickFix["id"] != "schema.enum.useValue" || quickFix["value"] != "active" {
		t.Fatalf("schema-enum quick fix = %#v, want useValue active", quickFixes[0])
	}
}

func TestServerProfileRequiredDefaultQuickFix(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "required-default.csv")
	if err := os.WriteFile(path, []byte("id,status\n1,\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	profilePayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
		"schema": map[string]any{
			"fields": []map[string]any{
				{"name": "status", "constraints": map[string]any{"required": true, "default": "active"}},
			},
		},
	})
	issue := jsonIssueByType(profilePayload["issues"].([]any), "schema-required")
	if issue == nil {
		t.Fatalf("expected schema-required issue in %#v", profilePayload["issues"])
	}
	quickFixes, ok := issue["quickFixes"].([]any)
	if !ok || len(quickFixes) != 1 {
		t.Fatalf("schema-required quick fixes = %#v, want one fix", issue["quickFixes"])
	}
	quickFix, ok := quickFixes[0].(map[string]any)
	if !ok || quickFix["id"] != "schema.default.use" || quickFix["value"] != "active" {
		t.Fatalf("schema-required quick fix = %#v, want default active", quickFixes[0])
	}
}

func TestServerProfileSemanticDefaultQuickFix(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "semantic-default.csv")
	if err := os.WriteFile(path, []byte("id,age,code\n1,n/a,bad\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	profilePayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
		"schema": map[string]any{
			"columns": []map[string]any{
				{"name": "age", "type": "integer", "default": "42"},
				{"name": "code", "pattern": "^A-", "default": "A-1"},
			},
		},
	})
	issues := profilePayload["issues"].([]any)
	for _, check := range []struct {
		typ   string
		col   float64
		value string
	}{
		{typ: "schema-type", col: 2, value: "42"},
		{typ: "schema-pattern", col: 3, value: "A-1"},
	} {
		issue := jsonIssueByTypeCol(issues, check.typ, check.col)
		if issue == nil {
			t.Fatalf("expected %s issue at col %.0f in %#v", check.typ, check.col, issues)
		}
		quickFixes, ok := issue["quickFixes"].([]any)
		if !ok || len(quickFixes) != 1 {
			t.Fatalf("%s quick fixes = %#v, want one default fix", check.typ, issue["quickFixes"])
		}
		quickFix, ok := quickFixes[0].(map[string]any)
		if !ok || quickFix["id"] != "schema.default.use" || quickFix["value"] != check.value {
			t.Fatalf("%s quick fix = %#v, want default %s", check.typ, quickFixes[0], check.value)
		}
	}
}

func TestServerProfileStreamsCleanFileOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile-stream.csv")
	if err := os.WriteFile(path, []byte("id,region,score\n1,APAC,8\n2,EU,9\n3,APAC,10\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 2,
		"indexRows":   0,
	})
	documentID := openPayload["documentId"].(string)

	profilePayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
	})
	if profilePayload["streamed"] != true || profilePayload["rows"] != float64(3) {
		t.Fatalf("profile payload = %#v, want streamed three-row profile", profilePayload)
	}
	columns := profilePayload["columns"].([]any)
	if len(columns) != 3 || columns[1].(map[string]any)["distinct"] != float64(2) {
		t.Fatalf("profile columns = %#v, want region with two distinct values", columns)
	}

	if editPayload := postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        2,
		"col":        1,
		"value":      "LATAM",
	}); editPayload["dirty"] != true {
		t.Fatalf("edit payload = %#v, want dirty edit", editPayload)
	}
	editedPayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
	})
	if editedPayload["streamed"] == true {
		t.Fatalf("edited profile streamed = true, want materialized fallback")
	}
	editedColumns := editedPayload["columns"].([]any)
	editedRegion := editedColumns[1].(map[string]any)
	if editedRegion["distinct"] != float64(2) || !jsonTopValuesContain(editedRegion["topValues"].([]any), "LATAM") {
		t.Fatalf("edited profile columns = %#v, want materialized region edit included", editedColumns)
	}
}

func TestServerProfileStreamsRowProgressEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile-progress.csv")
	if err := os.WriteFile(path, []byte("id,region\n1,APAC\n2,EU\n3,APAC\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 1,
		"indexRows":   0,
	})
	response := postJSONResponse(t, handler, "/api/profile", map[string]any{
		"documentId":       openPayload["documentId"].(string),
		"streamProgress":   true,
		"progressInterval": 2,
	})
	if response.Code != http.StatusOK {
		t.Fatalf("profile stream status = %d body=%s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/x-ndjson") {
		t.Fatalf("profile stream content type = %q", got)
	}
	var events []map[string]any
	decoder := json.NewDecoder(response.Body)
	for {
		var event map[string]any
		if err := decoder.Decode(&event); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}
		events = append(events, event)
	}
	if len(events) != 5 {
		t.Fatalf("events = %#v, want four progress events plus final result", events)
	}
	for index, want := range []struct {
		phase string
		rows  float64
		done  bool
	}{
		{"profile", 2, false},
		{"profile", 3, true},
		{"validation", 2, false},
		{"validation", 3, true},
	} {
		event := events[index]
		if event["type"] != "progress" || event["phase"] != want.phase || event["rows"] != want.rows || (event["done"] == true) != want.done {
			t.Fatalf("event %d = %#v, want %s rows %.0f done %v", index, event, want.phase, want.rows, want.done)
		}
	}
	final := events[len(events)-1]
	result, ok := final["result"].(map[string]any)
	if final["type"] != "result" || !ok {
		t.Fatalf("final event = %#v, want result", final)
	}
	progress := result["progress"].(map[string]any)
	if result["streamed"] != true || progress["profileRows"] != float64(3) || progress["validationRows"] != float64(3) {
		t.Fatalf("result progress = %#v, want streamed three-row progress", result)
	}
}

func TestServerProfileUsesRequestContextCancellation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile-cancel.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 1,
		"indexRows":   0,
	})
	body, err := json.Marshal(map[string]any{
		"documentId": openPayload["documentId"].(string),
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodPost, "/api/profile", bytes.NewReader(body)).WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != statusClientClosedRequest {
		t.Fatalf("canceled profile status = %d, want %d", response.Code, statusClientClosedRequest)
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), "context canceled") {
		t.Fatalf("canceled profile payload = %#v, want context canceled", payload)
	}
}

func TestServerProfileAcceptsCSVLintSchema(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "csvlint-schema.csv")
	if err := os.WriteFile(path, []byte("id,status,score,code\n1,new,10,A\n1,bad,105,OK\n,done,ok,ABCDE\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	profilePayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
		"schema": map[string]any{
			"fields": []map[string]any{
				{"name": "id", "datatype": "integer", "constraints": map[string]any{"required": true, "unique": true}},
				{"name": "status", "constraints": map[string]any{"enum": []string{"new", "done"}}},
				{"name": "score", "datatype": map[string]any{"base": "number", "minimum": 0, "maximum": 100}},
				{"name": "code", "datatype": "string", "constraints": map[string]any{"minLength": 2, "max_length": 4}},
			},
			"primaryKey": "id",
		},
	})
	if profilePayload["schemaRules"].(float64) != 4 {
		t.Fatalf("schemaRules = %#v, want 4", profilePayload["schemaRules"])
	}
	issues := profilePayload["issues"].([]any)
	for _, typ := range []string{"schema-unique", "schema-required", "schema-enum", "schema-range", "schema-type", "schema-min-length", "schema-max-length"} {
		if !jsonIssuesContainType(issues, typ) {
			t.Fatalf("expected %s issue in %#v", typ, issues)
		}
	}
}

func TestServerProfileAcceptsFrictionlessResourcesSchema(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "resources-schema.csv")
	if err := os.WriteFile(path, []byte("status,score\npending,105\npending,50\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	profilePayload := postJSON(t, handler, "/api/profile", map[string]any{
		"documentId": documentID,
		"schema": map[string]any{
			"resources": []map[string]any{
				{
					"name": "scores",
					"schema": map[string]any{
						"fields": []map[string]any{
							{"name": "status", "constraints": map[string]any{"enum": []string{"active"}}},
							{"name": "score", "type": "number", "constraints": map[string]any{"minimum": 0, "maximum": 100}},
						},
						"primaryKey": "status",
					},
				},
			},
		},
	})
	if profilePayload["schemaRules"].(float64) != 2 {
		t.Fatalf("schemaRules = %#v, want 2", profilePayload["schemaRules"])
	}
	issues := profilePayload["issues"].([]any)
	for _, typ := range []string{"schema-enum", "schema-unique", "schema-range"} {
		if !jsonIssuesContainType(issues, typ) {
			t.Fatalf("expected %s issue in %#v", typ, issues)
		}
	}
}

func TestServerSQLFilterNestedConditionGroups(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested-filter.csv")
	if err := os.WriteFile(path, []byte("id,status,region,score\n1,active,APAC,8\n2,active,EU,9\n3,paused,APAC,9\n4,active,APAC,10\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	sqlFilterPayload := postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"mode":       "sql",
		"expression": `"status" = 'active' AND ("region" = 'APAC' OR "score" >= 9)`,
	})
	if sqlFilterPayload["totalRows"].(float64) != 3 {
		t.Fatalf("SQL nested filter totalRows = %#v, want 3", sqlFilterPayload["totalRows"])
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=5")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 3 {
		t.Fatalf("SQL nested filtered rows = %#v, want 3 rows", rows)
	}
	gotIDs := []string{rows[0].([]any)[0].(string), rows[1].([]any)[0].(string), rows[2].([]any)[0].(string)}
	if strings.Join(gotIDs, ",") != "1,2,4" {
		t.Fatalf("SQL nested filtered IDs = %#v, want 1,2,4", gotIDs)
	}
}

func TestServerSQLFilterTemporalAndIntervalTypes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "temporal.csv")
	if err := os.WriteFile(path, []byte("id,created_at,seen_at,duration,precise_duration\n1,2026-05-27,15:30:00.123,1 02:30 day to minute,1 02:03:04.123 day(2) to second(3)\n2,2026-05-28,15:29:00,10 second,1 second(3)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	sqlFilterPayload := postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"mode":       "sql",
		"expression": "CAST(created_at AS date) = DATE '2026-05-27' AND CAST(seen_at AS time) = TIME '15:30:00.123' AND CAST(duration AS interval) = INTERVAL '1 02:30' day to minute",
	})
	if sqlFilterPayload["totalRows"].(float64) != 1 {
		t.Fatalf("SQL temporal filter totalRows = %#v, want 1", sqlFilterPayload["totalRows"])
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=2")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 1 || rows[0].([]any)[0] != "1" {
		t.Fatalf("SQL temporal filtered window rows = %#v, want row 1", rows)
	}

	sqlFilterPayload = postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"mode":       "sql",
		"expression": "CAST(precise_duration AS interval) = INTERVAL '1 02:03:04.123' day(2) to second(3)",
	})
	if sqlFilterPayload["totalRows"].(float64) != 1 {
		t.Fatalf("SQL interval precision filter totalRows = %#v, want 1", sqlFilterPayload["totalRows"])
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=2")
	rows = windowPayload["rows"].([]any)
	if len(rows) != 1 || rows[0].([]any)[0] != "1" {
		t.Fatalf("SQL interval precision filtered window rows = %#v, want row 1", rows)
	}
}

func TestServerSQLFilterTemporalArithmetic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "temporal-arithmetic.csv")
	if err := os.WriteFile(path, []byte("id,created_at,seen_at,stamp\n1,2026-05-27,15:30:00.123,2026-05-27 14:30:00\n2,2026-05-28,15:29:00,2026-05-27 13:00:00\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	sqlFilterPayload := postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"mode":       "sql",
		"expression": "CAST(created_at AS date) + INTERVAL '2' day = DATE '2026-05-29' AND CAST(seen_at AS time) + INTERVAL '30' minute = TIME '16:00:00.123' AND CAST(stamp AS timestamp) + INTERVAL '1:30' hour to minute = TIMESTAMP '2026-05-27 16:00:00'",
	})
	if sqlFilterPayload["totalRows"].(float64) != 1 {
		t.Fatalf("SQL temporal arithmetic filter totalRows = %#v, want 1", sqlFilterPayload["totalRows"])
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=2")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 1 || rows[0].([]any)[0] != "1" {
		t.Fatalf("SQL temporal arithmetic filtered window rows = %#v, want row 1", rows)
	}

	sqlFilterPayload = postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"mode":       "sql",
		"expression": "CAST(stamp AS timestamp) + INTERVAL '-1:30' hour to minute = TIMESTAMP '2026-05-27 13:00:00'",
	})
	if sqlFilterPayload["totalRows"].(float64) != 1 {
		t.Fatalf("SQL signed interval arithmetic filter totalRows = %#v, want 1", sqlFilterPayload["totalRows"])
	}
	windowPayload = getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=2")
	rows = windowPayload["rows"].([]any)
	if len(rows) != 1 || rows[0].([]any)[0] != "1" {
		t.Fatalf("SQL signed interval filtered window rows = %#v, want row 1", rows)
	}
}

func TestServerSQLFilterBigIntPrecision(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bigint.csv")
	if err := os.WriteFile(path, []byte("id,name\n9007199254740992,Ada\n9007199254740993,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)

	sqlFilterPayload := postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"mode":       "sql",
		"expression": "CAST(id AS bigint) = 9007199254740993 AND id > 9007199254740992",
	})
	if sqlFilterPayload["totalRows"].(float64) != 1 {
		t.Fatalf("SQL bigint filter totalRows = %#v, want 1", sqlFilterPayload["totalRows"])
	}
	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=2")
	rows := windowPayload["rows"].([]any)
	if len(rows) != 1 || rows[0].([]any)[0] != "9007199254740993" || rows[0].([]any)[1] != "Grace" {
		t.Fatalf("SQL bigint filtered window rows = %#v, want Grace row", rows)
	}
}

func TestServerExportView(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "export.csv")
	if err := os.WriteFile(path, []byte("id,name,score\n1,Ada,9\n2,Linus,7\n3,Grace,10\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	postJSON(t, handler, "/api/filter", map[string]any{
		"documentId": documentID,
		"query":      "Ada",
	})

	response := postJSONResponse(t, handler, "/api/export-view", map[string]any{
		"documentId": documentID,
		"columns":    []int{0, 2},
	})
	if response.Code != http.StatusOK {
		t.Fatalf("export status = %d body=%s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("X-Opencsv-Rows"); got != "1" {
		t.Fatalf("export rows header = %q, want 1", got)
	}
	if got := response.Header().Get("X-Opencsv-Columns"); got != "2" {
		t.Fatalf("export columns header = %q, want 2", got)
	}
	if got := response.Header().Get("X-Opencsv-Filename"); got != "export.view.csv" {
		t.Fatalf("export filename = %q, want export.view.csv", got)
	}
	if got := response.Body.String(); got != "id,score\n1,9\n" {
		t.Fatalf("export body = %q", got)
	}
	fullResponse := postJSONResponse(t, handler, "/api/export-view", map[string]any{
		"documentId":   documentID,
		"fullDocument": true,
	})
	if fullResponse.Code != http.StatusOK {
		t.Fatalf("full export status = %d body=%s", fullResponse.Code, fullResponse.Body.String())
	}
	if got := fullResponse.Header().Get("X-Opencsv-Rows"); got != "3" {
		t.Fatalf("full export rows header = %q, want 3", got)
	}
	if got := fullResponse.Body.String(); got != "id,name,score\n1,Ada,9\n2,Linus,7\n3,Grace,10\n" {
		t.Fatalf("full export body = %q", got)
	}

	includeHeaders := false
	rowResponse := postJSONResponse(t, handler, "/api/export-view", map[string]any{
		"documentId":     documentID,
		"columns":        []int{1, 2},
		"rows":           []int{0},
		"format":         "tsv",
		"includeHeaders": includeHeaders,
	})
	if rowResponse.Code != http.StatusOK {
		t.Fatalf("row export status = %d body=%s", rowResponse.Code, rowResponse.Body.String())
	}
	if got := rowResponse.Header().Get("Content-Type"); got != "text/tab-separated-values; charset=utf-8" {
		t.Fatalf("row export content type = %q", got)
	}
	if got := rowResponse.Body.String(); got != "Ada\t9\n" {
		t.Fatalf("row export body = %q", got)
	}

	markdownResponse := postJSONResponse(t, handler, "/api/export-view", map[string]any{
		"documentId": documentID,
		"columns":    []int{0, 1},
		"format":     "markdown",
	})
	if markdownResponse.Code != http.StatusOK {
		t.Fatalf("markdown export status = %d body=%s", markdownResponse.Code, markdownResponse.Body.String())
	}
	if got := markdownResponse.Header().Get("Content-Type"); got != "text/markdown; charset=utf-8" {
		t.Fatalf("markdown export content type = %q", got)
	}
	if got := markdownResponse.Header().Get("X-Opencsv-Filename"); got != "export.view.md" {
		t.Fatalf("markdown export filename = %q, want export.view.md", got)
	}
	if got := markdownResponse.Body.String(); got != "| id | name |\n| --- | --- |\n| 1 | Ada |" {
		t.Fatalf("markdown export body = %q", got)
	}

	xlsxResponse := postJSONResponse(t, handler, "/api/export-view", map[string]any{
		"documentId": documentID,
		"columns":    []int{0, 2},
		"format":     "xlsx",
	})
	if xlsxResponse.Code != http.StatusOK {
		t.Fatalf("xlsx export status = %d body=%s", xlsxResponse.Code, xlsxResponse.Body.String())
	}
	if got := xlsxResponse.Header().Get("Content-Type"); got != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Fatalf("xlsx content type = %q", got)
	}
	if got := xlsxResponse.Header().Get("X-Opencsv-Filename"); got != "export.view.xlsx" {
		t.Fatalf("xlsx filename = %q", got)
	}
	if !bytes.HasPrefix(xlsxResponse.Body.Bytes(), []byte("PK")) {
		t.Fatalf("xlsx export is not a zip: %q", xlsxResponse.Body.Bytes()[:2])
	}
}

func TestServerExportFullDocumentUsesDirtyEditBase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "conflict.csv")
	if err := os.WriteFile(path, []byte("id,name,score\n1,Ada,9\n2,Linus,7\n3,Grace,10\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	postJSON(t, handler, "/api/edit", map[string]any{
		"documentId": documentID,
		"row":        2,
		"col":        1,
		"value":      "LOCAL",
	})
	if err := os.WriteFile(path, []byte("id,name,score\n1,DISK,9\n2,Disk,7\n4,Nia,11\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	response := postJSONResponse(t, handler, "/api/export-view", map[string]any{
		"documentId":   documentID,
		"fullDocument": true,
	})
	if response.Code != http.StatusOK {
		t.Fatalf("export status = %d body=%s", response.Code, response.Body.String())
	}
	if got := response.Body.String(); got != "id,name,score\n1,Ada,9\n2,LOCAL,7\n3,Grace,10\n" {
		t.Fatalf("full dirty export = %q", got)
	}
}

func TestServerFindSelectionScope(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "find-selection.csv")
	if err := os.WriteFile(path, []byte("id,a,b\n1,Ada,Ada\n2,Ada,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 3,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	findPayload := postJSON(t, handler, "/api/find", map[string]any{
		"documentId": documentID,
		"query":      "Ada",
		"cells": []map[string]any{
			{"row": 1, "col": 1},
			{"row": 2, "col": 1},
		},
	})
	if findPayload["total"].(float64) != 2 {
		t.Fatalf("find total = %#v, want 2", findPayload["total"])
	}
}

func TestServerSortMultipleColumns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "multi-sort.csv")
	if err := os.WriteFile(path, []byte("id,group,score\n1,A,2\n2,B,1\n3,A,1\n4,B,2\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 5,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	sortPayload := postJSON(t, handler, "/api/sort", map[string]any{
		"documentId": documentID,
		"hasHeader":  true,
		"keys": []map[string]any{
			{"column": 1},
			{"column": 2, "descending": true, "numeric": true},
		},
	})
	if sortPayload["totalRows"].(float64) != 4 {
		t.Fatalf("sort totalRows = %#v, want 4", sortPayload["totalRows"])
	}

	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=4")
	rows := windowPayload["rows"].([]any)
	var got []string
	for _, item := range rows {
		row := item.([]any)
		got = append(got, row[0].(string))
	}
	want := []string{"1", "3", "4", "2"}
	if len(got) != len(want) {
		t.Fatalf("sorted ids = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("sorted ids = %#v, want %#v", got, want)
		}
	}
}

func TestServerSortMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sort-mode.csv")
	if err := os.WriteFile(path, []byte("id,label\n1,aa\n2,bbbb\n3,c\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	handler := New(dir).Handler()
	openPayload := postJSON(t, handler, "/api/open", map[string]any{
		"path":        path,
		"previewRows": 4,
		"indexRows":   -1,
	})
	documentID := openPayload["documentId"].(string)
	postJSON(t, handler, "/api/sort", map[string]any{
		"documentId": documentID,
		"hasHeader":  true,
		"column":     1,
		"mode":       "length",
		"descending": true,
	})

	windowPayload := getJSON(t, handler, "/api/window?id="+documentID+"&start=0&count=3")
	rows := windowPayload["rows"].([]any)
	got := rows[0].([]any)[0].(string)
	if got != "2" {
		t.Fatalf("first length-sorted id = %q, want 2", got)
	}
}

func TestServerUploadOpen(t *testing.T) {
	handler := New(t.TempDir()).Handler()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "upload.csv")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("id,name\n1,Ada\n")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/upload-open", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusOK {
		t.Fatalf("upload-open failed: status=%d payload=%#v", response.Code, payload)
	}
	if payload["documentId"] == "" {
		t.Fatalf("missing documentId in %#v", payload)
	}
}

func TestServerCommands(t *testing.T) {
	payload := getJSON(t, New(t.TempDir()).Handler(), "/api/commands")
	commands, ok := payload["commands"].([]any)
	if !ok || len(commands) == 0 {
		t.Fatalf("commands payload = %#v", payload)
	}
	foundOpen := false
	foundOfficialFilterShortcut := false
	for _, item := range commands {
		command := item.(map[string]any)
		if command["id"] == "file.open" {
			foundOpen = true
		}
		if command["id"] == "view.filter" && command["shortcut"] == "Cmd+Shift+F" {
			foundOfficialFilterShortcut = true
		}
	}
	if !foundOpen {
		t.Fatalf("file.open not found in %#v", commands)
	}
	if !foundOfficialFilterShortcut {
		t.Fatalf("view.filter Cmd+Shift+F shortcut not found in %#v", commands)
	}
}

func TestServerUpdateCheck(t *testing.T) {
	payload := getJSON(t, New(t.TempDir()).Handler(), "/api/update-check?source=startup")
	if payload["source"] != "startup" {
		t.Fatalf("source = %#v, payload %#v", payload["source"], payload)
	}
	if payload["status"] != "offline" {
		t.Fatalf("status = %#v, payload %#v", payload["status"], payload)
	}
	if payload["checkedAt"] == "" {
		t.Fatalf("missing checkedAt in %#v", payload)
	}
}

func TestServerDiff(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n1,Ada\n2,Ken\n3,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/diff", map[string]any{
		"leftPath":  left,
		"rightPath": right,
	})
	summary, ok := payload["summary"].(map[string]any)
	if !ok {
		t.Fatalf("missing summary in %#v", payload)
	}
	if summary["changed"].(float64) != 1 || summary["added"].(float64) != 1 {
		t.Fatalf("unexpected diff summary: %#v", summary)
	}
	rows, ok := payload["rows"].([]any)
	if !ok || len(rows) == 0 {
		t.Fatalf("missing rows in %#v", payload)
	}
}

func TestServerDiffOptions(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name,updated\n1,Ada Lovelace,2026-01-01\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("updated,id,name\n2026-02-01,1,ada   lovelace\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/diff", map[string]any{
		"leftPath":           left,
		"rightPath":          right,
		"ignoreCase":         true,
		"ignoreWhitespace":   true,
		"matchColumnsByName": true,
		"ignoredColumnNames": []string{"updated"},
	})
	summary := payload["summary"].(map[string]any)
	if summary["changed"].(float64) != 0 {
		t.Fatalf("unexpected diff summary with options: %#v", summary)
	}
}

func TestServerDiffKeyColumn(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n3,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n3,Grace\n1,Ada\n2,Linux\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/diff", map[string]any{
		"leftPath":  left,
		"rightPath": right,
		"keyColumn": "id",
	})
	summary := payload["summary"].(map[string]any)
	if summary["changed"].(float64) != 1 || summary["added"].(float64) != 0 || summary["deleted"].(float64) != 0 {
		t.Fatalf("unexpected key-column diff summary: %#v", summary)
	}
}

func TestServerDiffRightText(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/diff", map[string]any{
		"leftPath":   left,
		"rightText":  "id,name\n1,Ada\n2,Linux\n3,Grace\n",
		"rightLabel": "clipboard",
	})
	summary := payload["summary"].(map[string]any)
	if summary["changed"].(float64) != 1 || summary["added"].(float64) != 1 {
		t.Fatalf("unexpected text diff summary: %#v", summary)
	}
	if payload["rightPath"] != "clipboard" {
		t.Fatalf("rightPath = %#v", payload["rightPath"])
	}
}

func TestServerDiffMerge(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n1,Linux\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/diff-merge", map[string]any{
		"leftPath":    left,
		"rightPath":   right,
		"resultIndex": 1,
		"mode":        "cell",
		"col":         2,
	})
	if payload["applied"] != "cell" {
		t.Fatalf("unexpected merge payload: %#v", payload)
	}
	bytes, err := os.ReadFile(left)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "id,name\n1,Linux\n" {
		t.Fatalf("merged file = %q", string(bytes))
	}
}

func TestServerDiffMergeLeftToRight(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Linux\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n1,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/diff-merge", map[string]any{
		"leftPath":    left,
		"rightPath":   right,
		"resultIndex": 1,
		"mode":        "cell",
		"col":         2,
		"direction":   "left-to-right",
	})
	if payload["applied"] != "cell" {
		t.Fatalf("unexpected merge payload: %#v", payload)
	}
	bytes, err := os.ReadFile(right)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "id,name\n1,Linux\n" {
		t.Fatalf("merged right file = %q", string(bytes))
	}
}

func TestServerDiffEditCellCreatesBackup(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "right.csv")
	original := []byte("id,name\r\n1,Linus\r\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/diff-edit-cell", map[string]any{
		"path":  path,
		"row":   2,
		"col":   2,
		"value": "Linux",
	})
	if payload["row"].(float64) != 2 || payload["col"].(float64) != 2 || payload["value"] != "Linux" {
		t.Fatalf("unexpected edit payload: %#v", payload)
	}
	save := payload["save"].(map[string]any)
	backupPath, ok := save["backupPath"].(string)
	if !ok || backupPath == "" {
		t.Fatalf("missing backup path in %#v", payload)
	}
	backup, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(backup, original) {
		t.Fatalf("backup = %q, want %q", backup, original)
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(saved) != "id,name\r\n1,Linux\r\n" {
		t.Fatalf("edited file = %q", saved)
	}

	response := postJSONResponse(t, New(dir).Handler(), "/api/diff-edit-cell", map[string]any{
		"path":  path,
		"row":   20,
		"col":   2,
		"value": "Out of range",
	})
	if response.Code != http.StatusBadRequest {
		t.Fatalf("out-of-range edit status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestServerDiffMergeRightText(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	payload := postJSON(t, New(dir).Handler(), "/api/diff-merge", map[string]any{
		"leftPath":    left,
		"rightText":   "id,name\n1,Ada\n2,Linux\n",
		"rightLabel":  "clipboard",
		"resultIndex": 2,
		"mode":        "cell",
		"col":         2,
	})
	if payload["applied"] != "cell" {
		t.Fatalf("unexpected text merge payload: %#v", payload)
	}
	bytes, err := os.ReadFile(left)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "id,name\n1,Ada\n2,Linux\n" {
		t.Fatalf("merged text file = %q", string(bytes))
	}
}

func TestServerSQLFileReference(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	teamsPath := filepath.Join(dir, "teams.csv")
	if err := os.WriteFile(teamsPath, []byte("id,team\n1,Compiler\n2,Kernel\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	regionsPath := filepath.Join(dir, "regions.csv")
	if err := os.WriteFile(regionsPath, []byte("region,label\nAPAC,Asia Pacific\nEU,Europe\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	payload := postJSON(t, New(dir).Handler(), "/api/sql", map[string]any{
		"query": `SELECT name FROM "@file:` + path + `" WHERE id = 2`,
	})
	rows := payload["rows"].([]any)
	if len(rows) != 1 || rows[0].([]any)[0] != "Linus" {
		t.Fatalf("unexpected SQL file rows: %#v", rows)
	}

	joinPayload := postJSON(t, New(dir).Handler(), "/api/sql", map[string]any{
		"query": `SELECT s.name, t.team FROM "@file:` + path + `" AS s JOIN "@file:` + teamsPath + `" AS t ON t.id = s.id WHERE CAST(s.id AS INTEGER) = 2`,
	})
	joinRows := joinPayload["rows"].([]any)
	if len(joinRows) != 1 || joinRows[0].([]any)[0] != "Linus" || joinRows[0].([]any)[1] != "Kernel" {
		t.Fatalf("unexpected SQL join rows: %#v", joinRows)
	}

	textHeaderPayload := postJSON(t, New(dir).Handler(), "/api/sql", map[string]any{
		"query": `SELECT label FROM "@file:` + regionsPath + `" WHERE region = 'APAC'`,
	})
	textHeaderRows := textHeaderPayload["rows"].([]any)
	if len(textHeaderRows) != 1 || textHeaderRows[0].([]any)[0] != "Asia Pacific" {
		t.Fatalf("unexpected SQL text-header rows: %#v", textHeaderRows)
	}
}

func TestServerSQLUsesRequestContextCancellation(t *testing.T) {
	body, err := json.Marshal(map[string]any{
		"query": `WITH RECURSIVE cnt(x) AS (VALUES(0) UNION ALL SELECT x+1 FROM cnt WHERE x < 100000000) SELECT max(x) FROM cnt`,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodPost, "/api/sql", bytes.NewReader(body)).WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	New(t.TempDir()).Handler().ServeHTTP(response, request)

	if response.Code != statusClientClosedRequest {
		t.Fatalf("canceled SQL status = %d, want %d", response.Code, statusClientClosedRequest)
	}
	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(payload["error"].(string), "context canceled") {
		t.Fatalf("canceled SQL payload = %#v, want context canceled", payload)
	}
}

func postJSONResponse(t *testing.T, handler http.Handler, path string, value any) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func postJSON(t *testing.T, handler http.Handler, path string, value any) map[string]any {
	t.Helper()
	response := postJSONResponse(t, handler, path, value)

	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if response.Code < 200 || response.Code >= 300 {
		t.Fatalf("POST %s failed: status=%d payload=%#v", path, response.Code, payload)
	}
	return payload
}

func getJSON(t *testing.T, handler http.Handler, path string) map[string]any {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if response.Code < 200 || response.Code >= 300 {
		t.Fatalf("GET %s failed: status=%d payload=%#v", path, response.Code, payload)
	}
	return payload
}

func jsonIssuesContainType(issues []any, typ string) bool {
	return jsonIssueByType(issues, typ) != nil
}

func jsonIssueByType(issues []any, typ string) map[string]any {
	for _, raw := range issues {
		issue, ok := raw.(map[string]any)
		if ok && issue["type"] == typ {
			return issue
		}
	}
	return nil
}

func jsonIssueByTypeCol(issues []any, typ string, col float64) map[string]any {
	for _, raw := range issues {
		issue, ok := raw.(map[string]any)
		if ok && issue["type"] == typ && issue["col"] == col {
			return issue
		}
	}
	return nil
}

func jsonTopValuesContain(values []any, want string) bool {
	for _, raw := range values {
		value, ok := raw.(map[string]any)
		if ok && value["value"] == want {
			return true
		}
	}
	return false
}
