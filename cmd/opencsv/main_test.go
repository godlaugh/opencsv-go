package main

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/opencsv/opencsv/pkg/csvcore"
	"github.com/opencsv/opencsv/pkg/protocol"
)

func TestParseRootOptionsInterspersedFlags(t *testing.T) {
	options, err := parseRootOptions([]string{
		"--diff",
		"old.tsv",
		"new.tsv",
		"-d",
		"tab",
		"--encoding=utf-8",
		"--line-ending",
		"crlf",
		"--goto",
		"new.tsv:10:3",
		"--wait",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !options.diff || !options.wait {
		t.Fatalf("diff/wait not parsed: %#v", options)
	}
	if options.format.Delimiter != "tab" || options.format.Encoding != "utf-8" || options.format.LineEnding != "crlf" {
		t.Fatalf("format not parsed: %#v", options.format)
	}
	if options.gotoTarget != "new.tsv:10:3" {
		t.Fatalf("goto target = %q", options.gotoTarget)
	}
	if len(options.paths) != 2 || options.paths[0] != "old.tsv" || options.paths[1] != "new.tsv" {
		t.Fatalf("paths = %#v", options.paths)
	}
}

func TestParseRootOptionsRejectsConflictingWindowFlags(t *testing.T) {
	if _, err := parseRootOptions([]string{"--new-window", "--reuse-window", "file.csv"}); err == nil {
		t.Fatal("expected conflicting window flag error")
	}
}

func TestOpenOptionsFromFormat(t *testing.T) {
	headerRows := 1
	headerColumns := 2
	format := protocol.FormatParams{Delimiter: "tab", Encoding: "utf8", LineEnding: "crlf", Quote: "single", QuoteMode: "always", HeaderRows: &headerRows, HeaderCols: &headerColumns}
	options, err := openOptionsFromFormat(format)
	if err != nil {
		t.Fatal(err)
	}
	if options.Delimiter == nil || *options.Delimiter != '\t' || options.DelimiterName != "tab" {
		t.Fatalf("delimiter override = %#v %q", options.Delimiter, options.DelimiterName)
	}
	if options.Encoding == nil || options.Encoding.Name != "utf-8" {
		t.Fatalf("encoding override = %#v", options.Encoding)
	}
	if options.LineEnding == nil || options.LineEnding.Name != "crlf" {
		t.Fatalf("line ending override = %#v", options.LineEnding)
	}
	if options.QuoteChar == nil || *options.QuoteChar != "'" || options.QuoteMode != "always" {
		t.Fatalf("quote override = %#v %q", options.QuoteChar, options.QuoteMode)
	}
	if options.HasHeader == nil || !*options.HasHeader {
		t.Fatalf("header override = %#v", options.HasHeader)
	}
	if options.HeaderColumns == nil || *options.HeaderColumns != 2 {
		t.Fatalf("header columns override = %#v", options.HeaderColumns)
	}
}

func TestRootOpenGotoWithoutSeparatePath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := rootOpen([]string{"--goto", path + ":2:1"}); err != nil {
		t.Fatal(err)
	}
}

func TestRootOpenRejectsInvalidDiffUsage(t *testing.T) {
	if err := rootOpen([]string{"--diff", "one.csv"}); err == nil {
		t.Fatal("expected --diff arity error")
	}
	if err := rootOpen([]string{"--diff", "one.csv", "-"}); err == nil {
		t.Fatal("expected --diff stdin error")
	}
}

func TestRootOpenDiffComputesResult(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n1,Ada\n2,Ken\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := rootOpen([]string{"--diff", left, right}); err != nil {
		t.Fatal(err)
	}
}

func TestRootOpenDiffIncludesWebRoute(t *testing.T) {
	dir := t.TempDir()
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	if err := os.WriteFile(left, []byte("id,name\n1,Ada\n2,Linus\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(right, []byte("id,name\n1,Ada\n2,Ken\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	output, err := captureStdout(t, func() error {
		return rootOpen([]string{"--diff", left, right})
	})
	if err != nil {
		t.Fatal(err)
	}
	var invocation cliInvocation
	if err := json.Unmarshal([]byte(output), &invocation); err != nil {
		t.Fatal(err)
	}
	if invocation.Mode != "diff" || invocation.Route == "" || invocation.WebURL == "" {
		t.Fatalf("diff invocation route = %#v", invocation)
	}
	route, err := url.Parse(invocation.Route)
	if err != nil {
		t.Fatal(err)
	}
	query := route.Query()
	if query.Get("diffRun") != "1" {
		t.Fatalf("diffRun = %q, want 1", query.Get("diffRun"))
	}
	if query.Get("diffLeft") != left || query.Get("diffRight") != right {
		t.Fatalf("diff route paths = %#v", query)
	}
	webURL, err := url.Parse(invocation.WebURL)
	if err != nil {
		t.Fatal(err)
	}
	if webURL.Query().Get("diffLeft") != left || webURL.Query().Get("diffRight") != right {
		t.Fatalf("web URL query = %#v", webURL.Query())
	}
}

func TestRootOpenIncludesOpenRouteForGotoAndFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.tsv")
	if err := os.WriteFile(path, []byte("id\tname\n1\tAda\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	output, err := captureStdout(t, func() error {
		return rootOpen([]string{
			"--goto", path + ":2:1",
			"--delimiter", "tab",
			"--encoding", "utf-8",
			"--quote", "single",
			"--quote-mode", "always",
			"--line-ending", "crlf",
			"--header-rows", "1",
			"--header-columns", "1",
			path,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	var invocation cliInvocation
	if err := json.Unmarshal([]byte(output), &invocation); err != nil {
		t.Fatal(err)
	}
	if invocation.Mode != "open" || invocation.Route == "" || invocation.WebURL == "" {
		t.Fatalf("open invocation route = %#v", invocation)
	}
	route, err := url.Parse(invocation.Route)
	if err != nil {
		t.Fatal(err)
	}
	query := route.Query()
	if query.Get("open") != path || query.Get("goto") != "2:1" {
		t.Fatalf("open route query = %#v", query)
	}
	if query.Get("delimiter") != "tab" || query.Get("encoding") != "utf-8" {
		t.Fatalf("format route query = %#v", query)
	}
	if query.Get("quote") != "single" || query.Get("quotemode") != "always" || query.Get("lineending") != "crlf" {
		t.Fatalf("official format route query = %#v", query)
	}
	if query.Get("headerrows") != "1" || query.Get("headercolumns") != "1" {
		t.Fatalf("header format route query = %#v", query)
	}
}

func TestRootOpenIncludesLaunchFlagsInRoutes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	left := filepath.Join(dir, "left.csv")
	right := filepath.Join(dir, "right.csv")
	for _, file := range []struct {
		path string
		body string
	}{
		{path: path, body: "id,name\n1,Ada\n"},
		{path: left, body: "id,name\n1,Ada\n"},
		{path: right, body: "id,name\n1,Grace\n"},
	} {
		if err := os.WriteFile(file.path, []byte(file.body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	openOutput, err := captureStdout(t, func() error {
		return rootOpen([]string{"--new-window", "--wait", path})
	})
	if err != nil {
		t.Fatal(err)
	}
	var openInvocation cliInvocation
	if err := json.Unmarshal([]byte(openOutput), &openInvocation); err != nil {
		t.Fatal(err)
	}
	openRoute, err := url.Parse(openInvocation.Route)
	if err != nil {
		t.Fatal(err)
	}
	openQuery := openRoute.Query()
	if !openInvocation.NewWindow || !openInvocation.Wait || openInvocation.ReuseWindow {
		t.Fatalf("open launch flags = %#v", openInvocation)
	}
	if openQuery.Get("newWindow") != "1" || openQuery.Get("wait") != "1" || openQuery.Get("reuseWindow") != "" {
		t.Fatalf("open route launch query = %#v", openQuery)
	}
	openWebURL, err := url.Parse(openInvocation.WebURL)
	if err != nil {
		t.Fatal(err)
	}
	if openWebURL.Query().Get("newWindow") != "1" || openWebURL.Query().Get("wait") != "1" {
		t.Fatalf("open web launch query = %#v", openWebURL.Query())
	}

	diffOutput, err := captureStdout(t, func() error {
		return rootOpen([]string{"--diff", "--reuse-window", "--wait", left, right})
	})
	if err != nil {
		t.Fatal(err)
	}
	var diffInvocation cliInvocation
	if err := json.Unmarshal([]byte(diffOutput), &diffInvocation); err != nil {
		t.Fatal(err)
	}
	diffRoute, err := url.Parse(diffInvocation.Route)
	if err != nil {
		t.Fatal(err)
	}
	diffQuery := diffRoute.Query()
	if diffInvocation.NewWindow || !diffInvocation.ReuseWindow || !diffInvocation.Wait {
		t.Fatalf("diff launch flags = %#v", diffInvocation)
	}
	if diffQuery.Get("reuseWindow") != "1" || diffQuery.Get("wait") != "1" || diffQuery.Get("newWindow") != "" {
		t.Fatalf("diff route launch query = %#v", diffQuery)
	}
}

func TestRootOpenParsesOpenCSVURLScheme(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "url.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rawURL := "opencsv://file" + path + ":2:2?Delimiter=comma&QuoteMode=always"
	output, err := captureStdout(t, func() error {
		return rootOpen([]string{rawURL})
	})
	if err != nil {
		t.Fatal(err)
	}
	var invocation cliInvocation
	if err := json.Unmarshal([]byte(output), &invocation); err != nil {
		t.Fatal(err)
	}
	if len(invocation.Files) != 1 || invocation.Files[0].Path != path {
		t.Fatalf("invocation files = %#v", invocation.Files)
	}
	if invocation.Files[0].Location != (protocol.Location{StartRow: 2, StartCol: 2}) {
		t.Fatalf("location = %#v", invocation.Files[0].Location)
	}
	if invocation.Files[0].Format.Delimiter != "comma" || invocation.Files[0].Format.QuoteMode != "always" {
		t.Fatalf("format = %#v", invocation.Files[0].Format)
	}
	webURL, err := url.Parse(invocation.WebURL)
	if err != nil {
		t.Fatal(err)
	}
	if webURL.Query().Get("open") != path || webURL.Query().Get("goto") != "2:2" || webURL.Query().Get("quotemode") != "always" {
		t.Fatalf("web URL query = %#v", webURL.Query())
	}
}

func TestOpenCLITargetReadsStdin(t *testing.T) {
	input, err := os.CreateTemp(t.TempDir(), "stdin-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := input.WriteString("id,name\n1,Ada\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := input.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	previousStdin := os.Stdin
	os.Stdin = input
	t.Cleanup(func() {
		os.Stdin = previousStdin
		_ = input.Close()
	})

	document, source, err := openCLITarget("-", csvcore.OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	if source != "stdin" {
		t.Fatalf("source = %q, want stdin", source)
	}
	if document.RowIndex.Rows != 2 {
		t.Fatalf("rows = %d, want 2", document.RowIndex.Rows)
	}
	if len(document.Preview) != 2 || document.Preview[1][1] != "Ada" {
		t.Fatalf("preview = %#v", document.Preview)
	}
}

func TestBenchReportsIndexedWindowMetrics(t *testing.T) {
	out := filepath.Join(t.TempDir(), "bench.csv")
	output, err := captureStdout(t, func() error {
		return bench([]string{
			"-rows", "20",
			"-cols", "4",
			"-preview-rows", "5",
			"-window-start", "10",
			"-window-rows", "6",
			"-out", out,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(output), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["path"] != out {
		t.Fatalf("bench path = %#v, want %q", payload["path"], out)
	}
	if payload["indexedRows"].(float64) != 21 || payload["indexComplete"] != true {
		t.Fatalf("index metrics = %#v", payload)
	}
	if payload["windowStart"].(float64) != 10 || payload["windowRowsRead"].(float64) != 6 {
		t.Fatalf("window metrics = %#v", payload)
	}
	if _, ok := payload["windowDurationMs"]; !ok {
		t.Fatalf("missing windowDurationMs in %#v", payload)
	}
	if _, ok := payload["windowDurationUs"]; !ok {
		t.Fatalf("missing windowDurationUs in %#v", payload)
	}
}

func TestGrepFileFindsCells(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grep.csv")
	if err := os.WriteFile(path, []byte("id,name,note\n1,Ada,Lovelace\n2,Grace,Ada friend\n3,Linus,kernel\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := grepFile(path, "Ada", grepOptions{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 || result.Truncated {
		t.Fatalf("grep result = %#v, want two untruncated matches", result)
	}
	if result.Matches[0].RowNumber != 2 || result.Matches[0].ColumnNumber != 2 || result.Matches[0].ColumnName != "name" {
		t.Fatalf("first match = %#v", result.Matches[0])
	}
	if result.Matches[1].RowNumber != 3 || result.Matches[1].ColumnNumber != 3 || result.Matches[1].ColumnName != "note" {
		t.Fatalf("second match = %#v", result.Matches[1])
	}
}

func TestGrepFileColumnIgnoreCaseAndLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grep-column.csv")
	if err := os.WriteFile(path, []byte("id,name,note\n1,Ada,Lovelace\n2,Grace,Ada friend\n3,Linus,kernel\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := grepFile(path, "ada", grepOptions{IgnoreCase: true, Column: "name", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 || result.Matches[0].Value != "Ada" {
		t.Fatalf("case-insensitive column grep = %#v", result)
	}

	limited, err := grepFile(path, "a", grepOptions{IgnoreCase: true, Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if limited.Total <= len(limited.Matches) || !limited.Truncated || len(limited.Matches) != 1 {
		t.Fatalf("limited grep = %#v, want one stored truncated match", limited)
	}
}

func TestGrepFileRegexAndErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grep-regex.csv")
	if err := os.WriteFile(path, []byte("id,name\n1,Ada\n2,Grace\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := grepFile(path, "^A", grepOptions{Regex: true, Column: "2", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 || result.Matches[0].Value != "Ada" {
		t.Fatalf("regex grep = %#v", result)
	}
	if _, err := grepFile(path, "Ada", grepOptions{Column: "missing"}); err == nil {
		t.Fatal("expected missing column error")
	}
	if _, err := grepFile(path, "[", grepOptions{Regex: true}); err == nil {
		t.Fatal("expected invalid regex error")
	}
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	previousStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	runErr := fn()
	closeErr := writer.Close()
	os.Stdout = previousStdout
	output, readErr := io.ReadAll(reader)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	return string(output), runErr
}
