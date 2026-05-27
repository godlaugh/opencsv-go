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
	jsonText, err := ToJSON(matrix)
	if err != nil {
		t.Fatal(err)
	}
	wantJSON := "[\n  [\n    \"name\",\n    \"note\"\n  ],\n  [\n    \"Ada\",\n    \"uses | pipes\"\n  ],\n  [\n    \"Linus\",\n    \"line\\nbreak\"\n  ]\n]"
	if jsonText != wantJSON {
		t.Fatalf("json = %q", jsonText)
	}
	html := ToHTML(matrix)
	wantHTML := "<table>\n  <tbody>\n    <tr><td>name</td><td>note</td></tr>\n    <tr><td>Ada</td><td>uses | pipes</td></tr>\n    <tr><td>Linus</td><td>line\nbreak</td></tr>\n  </tbody>\n</table>"
	if html != wantHTML {
		t.Fatalf("html = %q", html)
	}
	sql := ToSQLValues([][]string{{"name", "quote"}, {"Ada", "O'Reilly"}})
	wantSQL := "VALUES\n  ('name', 'quote'),\n  ('Ada', 'O''Reilly');"
	if sql != wantSQL {
		t.Fatalf("sql = %q", sql)
	}
	latex := ToLaTeX([][]string{{"name", "note"}, {"Ada_1", "uses | pipes"}})
	wantLaTeX := "\\begin{tabular}{ll}\nname & note \\\\\nAda\\_1 & uses \\textbar{} pipes\n\\end{tabular}"
	if latex != wantLaTeX {
		t.Fatalf("latex = %q", latex)
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
