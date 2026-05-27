package csvdiff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompareRecordsMarksChangedAddedDeleted(t *testing.T) {
	left := [][]string{
		{"id", "name", "score"},
		{"1", "Ada", "9"},
		{"2", "Linus", "7"},
		{"3", "Grace", "10"},
		{"4", "Alan", "8"},
	}
	right := [][]string{
		{"id", "name", "score"},
		{"1", "Ada", "9"},
		{"2", "Ken", "7"},
		{"3", "Grace", "10"},
		{"5", "Barbara", "11"},
	}

	result := CompareRecords("left.csv", "right.csv", left, right, Options{})
	if result.Summary.Unchanged != 3 || result.Summary.Changed != 1 || result.Summary.Added != 1 || result.Summary.Deleted != 1 {
		t.Fatalf("unexpected summary: %#v", result.Summary)
	}
	var changed RowDiff
	for _, row := range result.Rows {
		if row.Status == RowChanged {
			changed = row
			break
		}
	}
	if len(changed.ChangedCells) != 1 || changed.ChangedCells[0].Col != 2 || changed.ChangedCells[0].Left != "Linus" || changed.ChangedCells[0].Right != "Ken" {
		t.Fatalf("unexpected changed cells: %#v", changed.ChangedCells)
	}
}

func TestCompareRecordsDetectsInsertedRow(t *testing.T) {
	left := [][]string{
		{"id", "name"},
		{"1", "Ada"},
		{"2", "Linus"},
	}
	right := [][]string{
		{"id", "name"},
		{"1", "Ada"},
		{"9", "New"},
		{"2", "Linus"},
	}

	result := CompareRecords("left.csv", "right.csv", left, right, Options{})
	if result.Summary.Added != 1 || result.Summary.Changed != 0 || result.Summary.Deleted != 0 {
		t.Fatalf("unexpected summary: %#v", result.Summary)
	}
	if result.Rows[2].Status != RowAdded || result.Rows[2].RightRecord[1] != "New" {
		t.Fatalf("unexpected inserted row diff: %#v", result.Rows)
	}
}

func TestCompareFilesReadsCSV(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n1,Ada\n2,Ken\n3,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := CompareFiles(left, right, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if result.LeftPath != left || result.RightPath != right {
		t.Fatalf("paths not set: %#v", result)
	}
	if result.Summary.Changed != 1 || result.Summary.Added != 1 {
		t.Fatalf("unexpected summary: %#v", result.Summary)
	}
}

func TestCompareFileToText(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := CompareFileToText(left, "clipboard", "id,name\n1,Ada\n2,Linux\n3,Grace\n", Options{})
	if err != nil {
		t.Fatal(err)
	}
	if result.RightPath != "clipboard" {
		t.Fatalf("right path = %q", result.RightPath)
	}
	if result.Summary.Changed != 1 || result.Summary.Added != 1 {
		t.Fatalf("unexpected summary: %#v", result.Summary)
	}
}

func TestMergeRightTextToLeft(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := MergeRightTextToLeft(left, "clipboard", "id,name\n1,Ada\n2,Linux\n", MergeOptions{
		ResultIndex: 2,
		Mode:        MergeCell,
		Col:         2,
	}); err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(left)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "id,name\n1,Ada\n2,Linux\n" {
		t.Fatalf("merged text result = %q", string(bytes))
	}
}

func TestCompareRecordsCanLimitRows(t *testing.T) {
	result := CompareRecords("", "", [][]string{{"a"}, {"b"}}, [][]string{{"a"}, {"c"}}, Options{MaxRows: 1})
	if !result.Truncated {
		t.Fatal("expected truncated result")
	}
	if len(result.Rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(result.Rows))
	}
}

func TestCompareRecordsCanIgnoreCaseWhitespaceAndColumns(t *testing.T) {
	left := [][]string{
		{"id", "name", "updated"},
		{"1", "Ada Lovelace", "2026-01-01"},
	}
	right := [][]string{
		{"id", "name", "updated"},
		{"1", "ada   lovelace", "2026-02-01"},
	}

	result := CompareRecords("left.csv", "right.csv", left, right, Options{
		IgnoreCase:       true,
		IgnoreWhitespace: true,
		IgnoredColumns:   []int{3},
	})
	if result.Summary.Changed != 0 || result.Summary.Unchanged != 2 {
		t.Fatalf("unexpected summary: %#v rows=%#v", result.Summary, result.Rows)
	}
}

func TestCompareRecordsCanMatchColumnsByHeaderName(t *testing.T) {
	left := [][]string{
		{"id", "name", "score"},
		{"1", "Ada", "9"},
	}
	right := [][]string{
		{"score", "id", "name"},
		{"9", "1", "Ada"},
	}

	result := CompareRecords("left.csv", "right.csv", left, right, Options{MatchColumnsByName: true})
	if result.Summary.Changed != 0 || result.Summary.Unchanged != 2 {
		t.Fatalf("unexpected summary with header matching: %#v rows=%#v", result.Summary, result.Rows)
	}
	if got := result.Rows[1].RightRecord; got[0] != "1" || got[1] != "Ada" || got[2] != "9" {
		t.Fatalf("right record was not aligned to left header order: %#v", got)
	}
}

func TestCompareRecordsCanMatchRowsByKeyColumn(t *testing.T) {
	left := [][]string{
		{"id", "name", "score"},
		{"1", "Ada", "9"},
		{"2", "Linus", "7"},
		{"3", "Grace", "10"},
	}
	right := [][]string{
		{"id", "name", "score"},
		{"3", "Grace", "10"},
		{"1", "Ada", "9"},
		{"2", "Linux", "7"},
	}

	result := CompareRecords("left.csv", "right.csv", left, right, Options{KeyColumn: "id"})
	if result.Summary.Changed != 1 || result.Summary.Added != 0 || result.Summary.Deleted != 0 || result.Summary.Unchanged != 3 {
		t.Fatalf("unexpected key-column summary: %#v rows=%#v", result.Summary, result.Rows)
	}
	var changed RowDiff
	for _, row := range result.Rows {
		if row.Status == RowChanged {
			changed = row
			break
		}
	}
	if changed.LeftRow == nil || *changed.LeftRow != 3 || changed.RightRow == nil || *changed.RightRow != 4 {
		t.Fatalf("changed row was not aligned by key: %#v", changed)
	}
	if len(changed.ChangedCells) != 1 || changed.ChangedCells[0].Col != 2 || changed.ChangedCells[0].Left != "Linus" || changed.ChangedCells[0].Right != "Linux" {
		t.Fatalf("unexpected key-column changed cells: %#v", changed.ChangedCells)
	}
}

func TestCompareRecordsCanUseNumericKeyColumn(t *testing.T) {
	left := [][]string{
		{"a", "old"},
		{"b", "same"},
	}
	right := [][]string{
		{"b", "same"},
		{"a", "new"},
	}

	result := CompareRecords("left.csv", "right.csv", left, right, Options{KeyColumn: "1"})
	if result.Summary.Changed != 1 || result.Summary.Added != 0 || result.Summary.Deleted != 0 || result.Summary.Unchanged != 1 {
		t.Fatalf("unexpected numeric key summary: %#v rows=%#v", result.Summary, result.Rows)
	}
}

func TestCompareRecordsIncludesInlineSpans(t *testing.T) {
	left := [][]string{{"id", "name"}, {"1", "Linus"}}
	right := [][]string{{"id", "name"}, {"1", "Linux"}}

	result := CompareRecords("left.csv", "right.csv", left, right, Options{})
	var changed CellDiff
	for _, row := range result.Rows {
		if row.Status == RowChanged && len(row.ChangedCells) > 0 {
			changed = row.ChangedCells[0]
			break
		}
	}
	if changed.Col != 2 || len(changed.LeftSpans) == 0 || len(changed.RightSpans) == 0 {
		t.Fatalf("missing inline spans: %#v", changed)
	}
	foundChanged := false
	for _, span := range changed.RightSpans {
		if span.Changed && span.Text == "x" {
			foundChanged = true
		}
	}
	if !foundChanged {
		t.Fatalf("right spans did not highlight changed x: %#v", changed.RightSpans)
	}
}

func TestCompareRecordsInlineSpansRespectIgnoreCase(t *testing.T) {
	left := [][]string{{"id", "name"}, {"1", "Linus"}}
	right := [][]string{{"id", "name"}, {"1", "linux"}}

	result := CompareRecords("left.csv", "right.csv", left, right, Options{IgnoreCase: true})
	var changed CellDiff
	for _, row := range result.Rows {
		if row.Status == RowChanged && len(row.ChangedCells) > 0 {
			changed = row.ChangedCells[0]
			break
		}
	}
	var highlighted []string
	for _, span := range changed.LeftSpans {
		if span.Changed {
			highlighted = append(highlighted, span.Text)
		}
	}
	for _, span := range changed.RightSpans {
		if span.Changed {
			highlighted = append(highlighted, span.Text)
		}
	}
	if strings.Join(highlighted, "|") != "s|x" {
		t.Fatalf("highlighted = %#v, want only s/x", highlighted)
	}
}

func TestMergeRightToLeftCellCreatesBackup(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n1,Linux\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := MergeRightToLeft(left, right, MergeOptions{
		ResultIndex: 1,
		Mode:        MergeCell,
		Col:         2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Applied != "cell" || result.Save.BackupPath == "" {
		t.Fatalf("unexpected merge result: %#v", result)
	}
	bytes, err := os.ReadFile(left)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "id,name\n1,Linux\n" {
		t.Fatalf("merged left file = %q", string(bytes))
	}
	backup, err := os.ReadFile(result.Save.BackupPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(backup) != "id,name\n1,Linus\n" {
		t.Fatalf("backup = %q", string(backup))
	}
}

func TestMergeRightToLeftRowWithHeaderMatching(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name,score\n1,Ada,9\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("score,id,name\n10,1,Ada Byron\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := MergeRightToLeft(left, right, MergeOptions{
		DiffOptions: Options{MatchColumnsByName: true},
		ResultIndex: 1,
		Mode:        MergeRow,
	})
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(left)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "id,name,score\n1,Ada Byron,10\n" {
		t.Fatalf("merged left file = %q", string(bytes))
	}
}

func TestMergeRightToLeftAddedAndDeletedRows(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n3,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	addResult := CompareFilesForTest(t, left, right)
	addIndex := -1
	deleteIndex := -1
	for index, row := range addResult.Rows {
		if row.Status == RowAdded {
			addIndex = index
		}
		if row.Status == RowDeleted {
			deleteIndex = index
		}
	}
	if addIndex < 0 || deleteIndex < 0 {
		t.Fatalf("missing added/deleted rows: %#v", addResult.Rows)
	}
	if _, err := MergeRightToLeft(left, right, MergeOptions{ResultIndex: addIndex, Mode: MergeRow}); err != nil {
		t.Fatal(err)
	}
	// Recompute after inserting the added row; Grace is still deleted relative to the right side.
	deleteResult := CompareFilesForTest(t, left, right)
	deleteIndex = -1
	for index, row := range deleteResult.Rows {
		if row.Status == RowDeleted {
			deleteIndex = index
			break
		}
	}
	if deleteIndex < 0 {
		t.Fatalf("missing deleted row after add merge: %#v", deleteResult.Rows)
	}
	if _, err := MergeRightToLeft(left, right, MergeOptions{ResultIndex: deleteIndex, Mode: MergeRow}); err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(left)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "id,name\n1,Ada\n2,Linus\n" {
		t.Fatalf("merged left file = %q", string(bytes))
	}
}

func CompareFilesForTest(t *testing.T, left string, right string) Result {
	t.Helper()
	result, err := CompareFiles(left, right, Options{})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
