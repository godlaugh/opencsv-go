package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opencsv/opencsv/pkg/csvcore"
)

func TestManagerEditFindFilterSortSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	original := []byte("\xef\xbb\xbfid,name,score\r\n1,Ada,9\r\n2,Linus,7\r\n3,Grace,10\r\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "sample.csv", csvcore.OpenOptions{PreviewRows: 4, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.Edit(session.ID, 1, 1, "Ada Lovelace"); err != nil {
		t.Fatal(err)
	}

	matches, err := manager.Find(session.ID, FindOptions{Query: "lovelace"})
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 || matches[0].Row != 1 || matches[0].Col != 1 {
		t.Fatalf("unexpected matches: %#v", matches)
	}
	regexMatches, err := manager.Find(session.ID, FindOptions{Query: `A.+lace`, Regex: true, WholeCell: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(regexMatches) != 1 {
		t.Fatalf("regex matches = %#v, want one match", regexMatches)
	}
	replaceResult, err := manager.Replace(session.ID, ReplaceOptions{
		FindOptions: FindOptions{Query: "Lovelace", CaseSensitive: true},
		Replacement: "Byron",
		All:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if replaceResult.Count != 1 {
		t.Fatalf("replace count = %d, want 1", replaceResult.Count)
	}
	matches, err = manager.Find(session.ID, FindOptions{Query: "Ada Byron"})
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 || matches[0].Row != 1 || matches[0].Col != 1 {
		t.Fatalf("unexpected post-replace matches: %#v", matches)
	}

	count, err := manager.FilterContains(session.ID, FilterOptions{Query: "Grace"})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("filter count = %d, want 1", count)
	}
	window, err := manager.Window(session.ID, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if window.TotalRows != 1 || window.Rows[0][1] != "Grace" {
		t.Fatalf("unexpected filtered window: %#v", window)
	}

	count, err = manager.FilterSQL(session.ID, "WHERE score >= 9 AND LOWER(name) LIKE 'ada%'")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("SQL filter count = %d, want 1", count)
	}
	window, err = manager.Window(session.ID, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if window.TotalRows != 1 || window.Rows[0][1] != "Ada Byron" {
		t.Fatalf("unexpected SQL filtered window: %#v", window)
	}

	values, totalDistinct, err := manager.ColumnValues(session.ID, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if totalDistinct != 3 || len(values) != 3 {
		t.Fatalf("column values = %#v total=%d, want three distinct values", values, totalDistinct)
	}
	count, err = manager.FilterValues(session.ID, ValueFilterOptions{Column: 1, Values: []string{"Grace", "Ada Byron"}})
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("value filter count = %d, want 2", count)
	}

	count, err = manager.Sort(session.ID, SortOptions{Column: 2, Descending: true, Numeric: true, HasHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("sort count = %d, want 3", count)
	}
	window, err = manager.Window(session.ID, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	if window.Rows[0][1] != "Grace" {
		t.Fatalf("first sorted row = %#v, want Grace", window.Rows[0])
	}

	result, err := manager.Save(session.ID, "", true)
	if err != nil {
		t.Fatal(err)
	}
	if result.BackupPath == "" {
		t.Fatal("expected backup path")
	}
	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(saved[:3]) != "\xef\xbb\xbf" {
		t.Fatal("saved file lost UTF-8 BOM")
	}
}

func TestManagerStructuralMutationsUndoRedo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	original := []byte("id,name\n1,Ada\n2,Linus\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "sample.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := manager.InsertRow(session.ID, 2); err != nil {
		t.Fatal(err)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 4 || records[2][0] != "" || records[3][1] != "Linus" {
		t.Fatalf("records after insert row = %#v", records)
	}
	onDisk, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(onDisk) != string(original) {
		t.Fatalf("structural edit modified original before save: %q", onDisk)
	}

	if _, err := manager.Undo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 || records[2][1] != "Linus" {
		t.Fatalf("records after undo insert row = %#v", records)
	}
	if _, err := manager.Redo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 4 || records[2][0] != "" {
		t.Fatalf("records after redo insert row = %#v", records)
	}

	if _, err := manager.InsertColumn(session.ID, 1); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records[0]) != 3 || records[0][1] != "" || records[0][2] != "name" {
		t.Fatalf("records after insert column = %#v", records)
	}
	if _, err := manager.DeleteColumn(session.ID, 1); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records[0]) != 2 || records[0][1] != "name" {
		t.Fatalf("records after delete column = %#v", records)
	}
}

func TestManagerTransposeUndoAndLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "sample.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Transpose(session.ID, 5); err == nil {
		t.Fatal("transpose over limit succeeded")
	}
	result, err := manager.Transpose(session.ID, 100)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRows != 2 {
		t.Fatalf("transpose totalRows = %d, want 2", result.TotalRows)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{{"id", "1", "2"}, {"name", "Ada", "Linus"}}
	if !recordsEqual(records, want) {
		t.Fatalf("transposed records = %#v, want %#v", records, want)
	}
	if result.Session.Document.Metadata.Dialect.HasHeader {
		t.Fatal("transpose should clear header-row interpretation")
	}
	if _, err := manager.Undo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !recordsEqual(records, [][]string{{"id", "name"}, {"1", "Ada"}, {"2", "Linus"}}) {
		t.Fatalf("records after undo transpose = %#v", records)
	}
}

func TestManagerBatchEditSingleUndo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "sample.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.EditCells(session.ID, []CellEdit{
		{Row: 1, Col: 1, Value: "X"},
		{Row: 2, Col: 1, Value: "X"},
	}); err != nil {
		t.Fatal(err)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "X" || records[2][1] != "X" {
		t.Fatalf("records after batch edit = %#v", records)
	}
	if _, err := manager.Undo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "Ada" || records[2][1] != "Linus" {
		t.Fatalf("records after single undo = %#v", records)
	}
}

func TestManagerReadOnlyAndFileStatusRefresh(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "sample.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := manager.SetReadOnly(session.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := manager.Edit(session.ID, 1, 1, "Grace"); err == nil {
		t.Fatal("read-only edit succeeded")
	}
	if _, err := manager.Save(session.ID, "", true); err == nil {
		t.Fatal("read-only save succeeded")
	}
	if _, err := manager.InsertRow(session.ID, 1); err == nil {
		t.Fatal("read-only structural mutation succeeded")
	}
	if _, err := manager.Replace(session.ID, ReplaceOptions{
		FindOptions: FindOptions{Query: "Ada"},
		Replacement: "Grace",
		All:         true,
	}); err == nil {
		t.Fatal("read-only replace succeeded")
	}

	if _, err := manager.SetReadOnly(session.ID, false); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status, err := manager.FileStatus(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !status.ChangedExternally {
		t.Fatalf("status = %#v, want changedExternally", status)
	}
	refreshed, err := manager.Refresh(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.Document.RowIndex.Rows != 3 {
		t.Fatalf("rows after refresh = %d, want 3", refreshed.Document.RowIndex.Rows)
	}
	status, err = manager.FileStatus(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if status.ChangedExternally {
		t.Fatalf("status after refresh = %#v, want no external change", status)
	}
}

func recordsEqual(left [][]string, right [][]string) bool {
	if len(left) != len(right) {
		return false
	}
	for row := range left {
		if len(left[row]) != len(right[row]) {
			return false
		}
		for col := range left[row] {
			if left[row][col] != right[row][col] {
				return false
			}
		}
	}
	return true
}
