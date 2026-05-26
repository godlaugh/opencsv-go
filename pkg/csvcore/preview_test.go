package csvcore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenPreviewAndIndex(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.csv")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := GenerateSynthetic(file, 10, 5); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	preview, err := OpenPreview(path, 4)
	if err != nil {
		t.Fatal(err)
	}
	if preview.RowsScanned != 4 {
		t.Fatalf("rows scanned = %d, want 4", preview.RowsScanned)
	}

	index, err := BuildRowIndex(path, preview.Dialect, 0)
	if err != nil {
		t.Fatal(err)
	}
	if index.Rows != 11 {
		t.Fatalf("indexed rows = %d, want 11", index.Rows)
	}
	if !index.Complete {
		t.Fatal("expected complete index")
	}
}
