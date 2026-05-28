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
	if findPayload["total"].(float64) != 1 || findPayload["truncated"].(bool) {
		t.Fatalf("find total/truncated = %#v/%#v, want 1/false", findPayload["total"], findPayload["truncated"])
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
	})
	if sortPayload["totalRows"].(float64) != 3 {
		t.Fatalf("sort totalRows = %#v, want 3", sortPayload["totalRows"])
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
