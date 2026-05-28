package main

import (
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
	if options.format.Delimiter != "tab" || options.format.Encoding != "utf-8" {
		t.Fatalf("format not parsed: %#v", options.format)
	}
	if options.gotoTarget != "new.tsv:10:3" {
		t.Fatalf("goto target = %q", options.gotoTarget)
	}
	if len(options.paths) != 2 || options.paths[0] != "old.tsv" || options.paths[1] != "new.tsv" {
		t.Fatalf("paths = %#v", options.paths)
	}
}

func TestOpenOptionsFromFormat(t *testing.T) {
	headerRows := 1
	format := protocol.FormatParams{Delimiter: "tab", Encoding: "utf8", Quote: "single", QuoteMode: "always", HeaderRows: &headerRows}
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
	if options.QuoteChar == nil || *options.QuoteChar != "'" || options.QuoteMode != "always" {
		t.Fatalf("quote override = %#v %q", options.QuoteChar, options.QuoteMode)
	}
	if options.HasHeader == nil || !*options.HasHeader {
		t.Fatalf("header override = %#v", options.HasHeader)
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
