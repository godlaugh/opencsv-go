package csvdiff

import (
	"os"
	"path/filepath"
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

func TestCompareRecordsCanLimitRows(t *testing.T) {
	result := CompareRecords("", "", [][]string{{"a"}, {"b"}}, [][]string{{"a"}, {"c"}}, Options{MaxRows: 1})
	if !result.Truncated {
		t.Fatal("expected truncated result")
	}
	if len(result.Rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(result.Rows))
	}
}
