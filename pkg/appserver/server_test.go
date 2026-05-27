package appserver

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
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
	if findPayload["total"].(float64) != 1 || findPayload["truncated"].(bool) {
		t.Fatalf("find total/truncated = %#v/%#v, want 1/false", findPayload["total"], findPayload["truncated"])
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
	})
	if sortPayload["totalRows"].(float64) != 3 {
		t.Fatalf("sort totalRows = %#v, want 3", sortPayload["totalRows"])
	}

	sqlPayload := postJSON(t, handler, "/api/sql", map[string]any{
		"documentId": documentID,
		"query":      "SELECT name, score FROM current WHERE score >= 9 ORDER BY score DESC LIMIT 1",
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
	for _, item := range commands {
		command := item.(map[string]any)
		if command["id"] == "file.open" {
			foundOpen = true
		}
	}
	if !foundOpen {
		t.Fatalf("file.open not found in %#v", commands)
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
		"ignoredColumns":     []int{3},
	})
	summary := payload["summary"].(map[string]any)
	if summary["changed"].(float64) != 0 {
		t.Fatalf("unexpected diff summary with options: %#v", summary)
	}
}

func TestServerSQLFileReference(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	payload := postJSON(t, New(dir).Handler(), "/api/sql", map[string]any{
		"query": `SELECT name FROM "@file:` + path + `" WHERE id = 2`,
	})
	rows := payload["rows"].([]any)
	if len(rows) != 1 || rows[0].([]any)[0] != "Linus" {
		t.Fatalf("unexpected SQL file rows: %#v", rows)
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
