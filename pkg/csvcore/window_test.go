package csvcore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadRecordWindow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Linus\n3,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 2})
	if err != nil {
		t.Fatal(err)
	}

	window, err := ReadRecordWindow(path, document.Metadata, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(window.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(window.Rows))
	}
	if window.Rows[0][1] != "Linus" || window.Rows[1][1] != "Grace" {
		t.Fatalf("unexpected window rows: %#v", window.Rows)
	}
}
