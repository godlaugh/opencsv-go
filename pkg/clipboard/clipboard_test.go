package clipboard

import "testing"

func TestCopyFormats(t *testing.T) {
	matrix := [][]string{
		{"name", "note"},
		{"Ada", "uses | pipes"},
		{"Linus", "line\nbreak"},
	}
	if got := ToTSV(matrix); got != "name\tnote\nAda\tuses | pipes\nLinus\tline\nbreak" {
		t.Fatalf("tsv = %q", got)
	}
	csvText, err := ToCSV(matrix)
	if err != nil {
		t.Fatal(err)
	}
	if csvText != "name,note\nAda,uses | pipes\nLinus,\"line\nbreak\"" {
		t.Fatalf("csv = %q", csvText)
	}
	markdown := ToMarkdown(matrix)
	wantMarkdown := "| name | note |\n| --- | --- |\n| Ada | uses \\| pipes |\n| Linus | line<br>break |"
	if markdown != wantMarkdown {
		t.Fatalf("markdown = %q", markdown)
	}
}

func TestParseTSV(t *testing.T) {
	matrix := ParseTSV("a\tb\r\n1\t2\n3\t4\n")
	if len(matrix) != 3 || matrix[1][1] != "2" || matrix[2][0] != "3" {
		t.Fatalf("matrix = %#v", matrix)
	}
	if ParseTSV("") != nil {
		t.Fatal("empty input should produce nil matrix")
	}
}
