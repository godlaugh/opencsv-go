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

func TestManagerSortMultipleColumns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "multi-sort.csv")
	if err := os.WriteFile(path, []byte("id,group,score\n1,A,2\n2,B,1\n3,A,1\n4,B,2\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "multi-sort.csv", csvcore.OpenOptions{PreviewRows: 5, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	count, err := manager.Sort(session.ID, SortOptions{
		HasHeader: true,
		Keys: []SortKey{
			{Column: 1},
			{Column: 2, Descending: true, Numeric: true},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Fatalf("sort count = %d, want 4", count)
	}

	window, err := manager.Window(session.ID, 0, 4)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, row := range window.Rows {
		got = append(got, row[0])
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

func TestManagerSmartSortModes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "smart-sort.csv")
	if err := os.WriteFile(path, []byte("id,date,label,amount\n1,2026-05-02,aa,１０\n2,2025-01-01,bbbb,2\n3,2026-04-01,c,９\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "smart-sort.csv", csvcore.OpenOptions{PreviewRows: 4, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := manager.Sort(session.ID, SortOptions{Column: 3, Mode: SortAuto, HasHeader: true}); err != nil {
		t.Fatal(err)
	}
	assertSortedIDs(t, manager, session.ID, []string{"2", "3", "1"})

	if _, err := manager.Sort(session.ID, SortOptions{Column: 1, Mode: SortDate, HasHeader: true}); err != nil {
		t.Fatal(err)
	}
	assertSortedIDs(t, manager, session.ID, []string{"2", "3", "1"})

	if _, err := manager.Sort(session.ID, SortOptions{Column: 2, Descending: true, Mode: SortLength, HasHeader: true}); err != nil {
		t.Fatal(err)
	}
	assertSortedIDs(t, manager, session.ID, []string{"2", "1", "3"})
}

func assertSortedIDs(t *testing.T, manager *Manager, sessionID string, want []string) {
	t.Helper()
	window, err := manager.Window(sessionID, 0, len(want))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, row := range window.Rows {
		got = append(got, row[0])
	}
	if len(got) != len(want) {
		t.Fatalf("sorted ids = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("sorted ids = %#v, want %#v", got, want)
		}
	}
}

func TestManagerFindDetailedCountsBeyondLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "find.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,ada\n3,Ada Byron\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "find.csv", csvcore.OpenOptions{PreviewRows: 4, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	result, err := manager.FindDetailed(session.ID, FindOptions{Query: "ada", Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 3 || len(result.Matches) != 1 || !result.Truncated {
		t.Fatalf("find result = %#v, want total 3, one stored match, truncated", result)
	}
}

func TestManagerFindAndReplaceSelectionScope(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "find-selection.csv")
	if err := os.WriteFile(path, []byte("id,a,b\n1,Ada,Ada\n2,Ada,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "find-selection.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	matches, err := manager.Find(session.ID, FindOptions{
		Query: "Ada",
		Cells: []Cell{{Row: 1, Col: 1}, {Row: 2, Col: 1}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 2 {
		t.Fatalf("selection matches = %#v, want 2", matches)
	}
	for _, match := range matches {
		if match.Col != 1 {
			t.Fatalf("selection match outside column 1: %#v", matches)
		}
	}

	result, err := manager.Replace(session.ID, ReplaceOptions{
		FindOptions: FindOptions{
			Query: "Ada",
			Cells: []Cell{{Row: 1, Col: 2}},
		},
		Replacement: "Grace",
		All:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 1 {
		t.Fatalf("replace count = %d, want 1", result.Count)
	}
	window, err := manager.Window(session.ID, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if window.Rows[0][1] != "Ada" || window.Rows[0][2] != "Grace" || window.Rows[1][2] != "Ada" {
		t.Fatalf("unexpected scoped replace window: %#v", window.Rows)
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

func TestManagerMoveRowsAndColumnsUndo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "move.csv")
	if err := os.WriteFile(path, []byte("id,name,score\n1,Ada,10\n2,Linus,20\n3,Grace,30\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "move.csv", csvcore.OpenOptions{PreviewRows: 5, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	result, err := manager.MoveRow(session.ID, 2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 1 {
		t.Fatalf("move row changed = %d, want 1", result.Changed)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "Linus" || records[2][1] != "Ada" {
		t.Fatalf("records after move row = %#v", records)
	}
	if _, err := manager.Undo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "Ada" || records[2][1] != "Linus" {
		t.Fatalf("records after undo move row = %#v", records)
	}

	result, err = manager.MoveColumn(session.ID, 2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 1 {
		t.Fatalf("move column changed = %d, want 1", result.Changed)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !recordsEqual(records, [][]string{
		{"id", "score", "name"},
		{"1", "10", "Ada"},
		{"2", "20", "Linus"},
		{"3", "30", "Grace"},
	}) {
		t.Fatalf("records after move column = %#v", records)
	}
}

func TestManagerMoveBoundariesAreNoops(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "move-noop.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "move-noop.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	result, err := manager.MoveRow(session.ID, 1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 0 || result.Dirty {
		t.Fatalf("top data row move result = %#v, want clean no-op", result)
	}
	result, err = manager.MoveColumn(session.ID, 0, -1)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 0 || result.Dirty {
		t.Fatalf("first column move result = %#v, want clean no-op", result)
	}
}

func TestManagerHideRowAndClearView(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hide.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n3,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "hide.csv", csvcore.OpenOptions{PreviewRows: 4, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	result, err := manager.HideRow(session.ID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 1 || result.TotalRows != 2 || result.Dirty {
		t.Fatalf("hide result = %#v, want one clean hidden row", result)
	}
	window, err := manager.Window(session.ID, 0, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(window.Rows) != 2 || window.Rows[0][1] != "Ada" || window.Rows[1][1] != "Grace" {
		t.Fatalf("window after hide row = %#v", window.Rows)
	}
	count, err := manager.ClearView(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Fatalf("clear view count = %d, want 4 source records", count)
	}
	window, err = manager.Window(session.ID, 0, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(window.Rows) != 4 || window.Rows[2][1] != "Linus" {
		t.Fatalf("window after clear hide = %#v", window.Rows)
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

func TestManagerDeduplicateRowsAndUndo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dupes.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n1,Ada\n2,Linus\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "dupes.csv", csvcore.OpenOptions{PreviewRows: 5, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	result, err := manager.DeduplicateRows(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 2 || result.TotalRows != 3 {
		t.Fatalf("dedupe result = %#v, want 2 removed and 3 rows", result)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !recordsEqual(records, [][]string{{"id", "name"}, {"1", "Ada"}, {"2", "Linus"}}) {
		t.Fatalf("deduped records = %#v", records)
	}
	if _, err := manager.Undo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 5 {
		t.Fatalf("records after undo dedupe = %#v", records)
	}
}

func TestDeduplicateRowsKeepsDistinctRaggedKeys(t *testing.T) {
	records := [][]string{
		{"value"},
		{"a\x00b"},
		{"a", "b"},
	}
	deduped, removed := deduplicateRecords(records, true)
	if removed != 0 || len(deduped) != len(records) {
		t.Fatalf("deduped = %#v, removed = %d", deduped, removed)
	}
}

func TestManagerDeduplicateColumnsAndUndo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dupe-cols.csv")
	content := "id,name,name,score,score\n1,Ada,Ada,10,10\n2,Linus,Linus,20,20\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "dupe-cols.csv", csvcore.OpenOptions{PreviewRows: 5, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	result, err := manager.DeduplicateColumns(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 2 || result.TotalRows != 3 {
		t.Fatalf("dedupe columns result = %#v, want 2 removed and 3 rows", result)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	expected := [][]string{
		{"id", "name", "score"},
		{"1", "Ada", "10"},
		{"2", "Linus", "20"},
	}
	if !recordsEqual(records, expected) {
		t.Fatalf("deduped column records = %#v", records)
	}
	if _, err := manager.Undo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(records[0]) != 5 {
		t.Fatalf("records after undo column dedupe = %#v", records)
	}
}

func TestDeduplicateColumnsKeepsDistinctHeaders(t *testing.T) {
	records := [][]string{
		{"name", "name_copy"},
		{"Ada", "Ada"},
		{"Linus", "Linus"},
	}
	deduped, removed := deduplicateColumns(records)
	if removed != 0 || !recordsEqual(deduped, records) {
		t.Fatalf("deduped = %#v, removed = %d", deduped, removed)
	}
}

func TestManagerTextTransforms(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "text.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,  Ada \n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "text.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	result, err := manager.TransformCells(session.ID, []Cell{{Row: 1, Col: 1}, {Row: 2, Col: 1}}, TransformTrim)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 1 {
		t.Fatalf("trim changed = %d, want 1", result.Changed)
	}
	result, err = manager.TransformCells(session.ID, []Cell{{Row: 1, Col: 1}, {Row: 2, Col: 1}}, TransformUpper)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 2 {
		t.Fatalf("upper changed = %d, want 2", result.Changed)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "ADA" || records[2][1] != "LINUS" {
		t.Fatalf("transformed records = %#v", records)
	}
}

func TestManagerJoinAndSplitCells(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "join-split.csv")
	if err := os.WriteFile(path, []byte("id,first,last\n1,Ada,Lovelace\n2,Grace,Hopper\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "join-split.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	result, err := manager.JoinCells(session.ID, []Cell{{Row: 1, Col: 1}, {Row: 1, Col: 2}}, " ")
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 2 {
		t.Fatalf("join changed = %d, want 2", result.Changed)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "Ada Lovelace" || records[1][2] != "" {
		t.Fatalf("records after join = %#v", records)
	}
	if _, err := manager.Undo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "Ada" || records[1][2] != "Lovelace" {
		t.Fatalf("records after undo join = %#v", records)
	}

	result, err = manager.SplitCells(session.ID, []Cell{{Row: 1, Col: 1}}, "d")
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 2 {
		t.Fatalf("split changed = %d, want 2", result.Changed)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "A" || records[1][2] != "a" {
		t.Fatalf("records after split = %#v", records)
	}
}

func TestManagerSwapCells(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "swap.csv")
	if err := os.WriteFile(path, []byte("id,a,b\n1,Ada,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manager := NewManager()
	session, err := manager.OpenPath(path, SourceFile, "swap.csv", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}

	result, err := manager.SwapCells(session.ID, []Cell{{Row: 1, Col: 1}, {Row: 1, Col: 2}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 2 || !result.Dirty {
		t.Fatalf("swap result = %#v", result)
	}
	_, records, err := manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "Linus" || records[1][2] != "Ada" {
		t.Fatalf("records after swap = %#v", records)
	}
	if _, err := manager.Undo(session.ID); err != nil {
		t.Fatal(err)
	}
	_, records, err = manager.Records(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if records[1][1] != "Ada" || records[1][2] != "Linus" {
		t.Fatalf("records after undo = %#v", records)
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
