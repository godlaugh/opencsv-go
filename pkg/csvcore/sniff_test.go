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

func TestSniffDetectsSingleQuoteWithEmbeddedDelimiter(t *testing.T) {
	dialect := Sniff([]byte("id,note\n1,'hello, world'\n2,'plain'\n"))
	if dialect.Delimiter != ',' {
		t.Fatalf("delimiter = %q, want comma", dialect.Delimiter)
	}
	if dialect.QuoteChar != "'" {
		t.Fatalf("quote char = %q, want single quote", dialect.QuoteChar)
	}
	if dialect.EstimatedColumns != 2 {
		t.Fatalf("columns = %d, want 2", dialect.EstimatedColumns)
	}
}

func TestSniffDetectsSingleQuoteWithoutEmbeddedDelimiter(t *testing.T) {
	dialect := Sniff([]byte("id,note\n'1','Ada'\n'2','Grace'\n"))
	if dialect.QuoteChar != "'" {
		t.Fatalf("quote char = %q, want single quote", dialect.QuoteChar)
	}
	if dialect.EstimatedColumns != 2 {
		t.Fatalf("columns = %d, want 2", dialect.EstimatedColumns)
	}
}

func TestSniffDetectsTextDimensionHeader(t *testing.T) {
	dialect := Sniff([]byte("region,label\nAPAC,Asia Pacific\nEU,Europe\n"))
	if !dialect.HasHeader {
		t.Fatal("expected text dimension table header")
	}
}

func TestSniffAvoidsWeakAllTextHeader(t *testing.T) {
	dialect := Sniff([]byte("A,B\nC,D\nE,F\n"))
	if dialect.HasHeader {
		t.Fatal("did not expect weak all-text rows to be promoted to header")
	}
}

func TestSniffAvoidsPersonNameRowsAsHeader(t *testing.T) {
	dialect := Sniff([]byte("Ada,Lovelace\nGrace,Hopper\nLinus,Torvalds\n"))
	if dialect.HasHeader {
		t.Fatal("did not expect person-name data rows to be promoted to header")
	}
}

func TestSniffRequiresMoreThanOneSemanticTextHeader(t *testing.T) {
	dialect := Sniff([]byte("Ada,score\nGrace,Hopper\nLinus,Torvalds\n"))
	if dialect.HasHeader {
		t.Fatal("did not expect one common-looking text value to promote a text table header")
	}
}

func TestSniffDetectsGeneratedFieldHeaders(t *testing.T) {
	dialect := Sniff([]byte("field_01,field_02\nsegment_0,2026-05-01\nsegment_1,2026-05-02\n"))
	if !dialect.HasHeader {
		t.Fatal("expected generated field headers")
	}
}

func TestSniffRejectsDuplicateHeaderNames(t *testing.T) {
	dialect := Sniff([]byte("name,name\nAda,Lovelace\nGrace,Hopper\n"))
	if dialect.HasHeader {
		t.Fatal("did not expect duplicate header names to be promoted")
	}
}
