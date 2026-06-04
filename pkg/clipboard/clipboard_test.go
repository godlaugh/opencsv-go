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

func TestWithHeaders(t *testing.T) {
	matrix := WithHeaders([][]string{{"1", "Ada"}, {"2", "Linus"}}, []string{"id", "name"})
	if len(matrix) != 3 || matrix[0][0] != "id" || matrix[1][1] != "Ada" {
		t.Fatalf("matrix with headers = %#v", matrix)
	}
	matrix[0][0] = "changed"
	if got := ToTSV(WithHeaders([][]string{{"1"}}, []string{"id"})); got != "id\n1" {
		t.Fatalf("header copy tsv = %q", got)
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
	if ParseTSV("\r\n") != nil {
		t.Fatal("newline-only input should produce nil matrix")
	}
}

func TestParseTSVQuotedSpreadsheetCells(t *testing.T) {
	matrix := ParseTSV("name\tnote\r\nAda\t\"line\r\nbreak\"\r\nLinus\t\"has\t tab and \"\"quote\"\"\"\r\n")
	if len(matrix) != 3 {
		t.Fatalf("matrix row count = %#v", matrix)
	}
	if got := matrix[1][1]; got != "line\nbreak" {
		t.Fatalf("quoted newline cell = %q", got)
	}
	if got := matrix[2][1]; got != "has\t tab and \"quote\"" {
		t.Fatalf("quoted tab/quote cell = %q", got)
	}
}
