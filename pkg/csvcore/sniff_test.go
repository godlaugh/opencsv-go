package csvcore

import "testing"

func TestSniffDetectsCommaAndHeader(t *testing.T) {
	dialect := Sniff([]byte("id,name,score\n1,Ada,98\n2,Linus,95\n"))
	if dialect.Delimiter != ',' {
		t.Fatalf("delimiter = %q, want comma", dialect.Delimiter)
	}
	if !dialect.HasHeader {
		t.Fatal("expected header")
	}
	if dialect.EstimatedColumns != 3 {
		t.Fatalf("columns = %d, want 3", dialect.EstimatedColumns)
	}
}

func TestSniffDetectsTab(t *testing.T) {
	dialect := Sniff([]byte("id\tname\tscore\n1\tAda\t98\n2\tLinus\t95\n"))
	if dialect.Delimiter != '\t' {
		t.Fatalf("delimiter = %q, want tab", dialect.Delimiter)
	}
	if dialect.DelimiterName != "tab" {
		t.Fatalf("delimiter name = %q, want tab", dialect.DelimiterName)
	}
}
