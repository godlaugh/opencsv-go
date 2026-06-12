package csvcore

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestRequiresQuote(t *testing.T) {
	dialect := Dialect{Delimiter: ',', QuoteChar: `"`, QuoteMode: QuoteModeMinimal}
	cases := []struct {
		value string
		want  bool
	}{
		{"plain", false},
		{"", false},
		{"has,comma", true},
		{"has \"quote\"", true},
		{"line\nbreak", true},
		{"carriage\rreturn", true},
		{" leading", true},
		{"\tleading-tab", true},
		{"trailing ", false},
	}
	for _, tc := range cases {
		if got := RequiresQuote(tc.value, dialect); got != tc.want {
			t.Fatalf("RequiresQuote(%q) = %v, want %v", tc.value, got, tc.want)
		}
	}
	if !RequiresQuote("has\ttab", Dialect{Delimiter: '\t'}) {
		t.Fatal("tab-delimited field containing tab should require quote")
	}
	if RequiresQuote("has,comma", Dialect{Delimiter: '\t'}) {
		t.Fatal("tab-delimited field containing comma should not require quote")
	}
}

func TestNormalizeQuotePolicy(t *testing.T) {
	mode, err := NormalizeQuoteMode("always")
	if err != nil || mode != QuoteModeAlways {
		t.Fatalf("mode = %q, err=%v", mode, err)
	}
	mode, err = NormalizeQuoteMode("manual")
	if err != nil || mode != QuoteModeManual {
		t.Fatalf("manual mode = %q, err=%v", mode, err)
	}
	quoteChar, err := NormalizeQuoteChar("single")
	if err != nil || quoteChar != "'" {
		t.Fatalf("quote char = %q, err=%v", quoteChar, err)
	}
	if !ShouldQuote("plain", Dialect{Delimiter: ',', QuoteMode: QuoteModeAlways}) {
		t.Fatal("always mode should quote plain values")
	}
	if ShouldQuote("has,comma", Dialect{Delimiter: ',', QuoteMode: QuoteModeNever}) {
		t.Fatal("never mode should not quote delimiter-containing values")
	}
	dialect, changed := SetCellQuote(Dialect{Delimiter: ',', QuoteMode: QuoteModeManual}, 1, 2, true)
	if !changed || !ShouldQuoteField("plain", dialect, 1, 2) {
		t.Fatalf("manual quote toggle failed: changed=%v dialect=%+v", changed, dialect)
	}
	dialect, changed = SetCellQuote(dialect, 1, 2, false)
	if !changed || ShouldQuoteField("plain", dialect, 1, 2) {
		t.Fatalf("manual unquote toggle failed: changed=%v dialect=%+v", changed, dialect)
	}

	dialect, changed = SetQuoteHeader(Dialect{Delimiter: ',', QuoteMode: QuoteModeManual, HasHeader: true}, true)
	if !changed || !ShouldQuoteField("id", dialect, 0, 0) || ShouldQuoteField("plain", dialect, 1, 0) {
		t.Fatalf("header quote rule failed: changed=%v dialect=%+v", changed, dialect)
	}
	dialect, changed = SetQuoteBody(dialect, true)
	if !changed || !ShouldQuoteField("plain", dialect, 1, 0) {
		t.Fatalf("body quote rule failed: changed=%v dialect=%+v", changed, dialect)
	}
	dialect, changed = SetQuoteBody(dialect, false)
	if !changed || ShouldQuoteField("plain", dialect, 1, 0) {
		t.Fatalf("body quote off failed: changed=%v dialect=%+v", changed, dialect)
	}
	dialect, changed = SetQuoteHeader(dialect, false)
	if !changed || ShouldQuoteField("id", dialect, 0, 0) {
		t.Fatalf("header quote off failed: changed=%v dialect=%+v", changed, dialect)
	}
	dialect, changed = SetQuoteColumn(dialect, 1, true)
	if !changed || !ShouldQuoteField("plain", dialect, 1, 1) || ShouldQuoteField("note", dialect, 0, 1) {
		t.Fatalf("column quote rule failed: changed=%v dialect=%+v", changed, dialect)
	}
	dialect, changed = SetQuoteColumn(dialect, 1, false)
	if !changed || ShouldQuoteField("plain", dialect, 1, 1) {
		t.Fatalf("column quote off failed: changed=%v dialect=%+v", changed, dialect)
	}
}

func TestDetectQuoteRulesReportsQuoteMapLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "quote-map-limit.csv")
	var input strings.Builder
	input.WriteString("id,note\n")
	for quoted := 0; quoted <= MaxQuoteCellRules; quoted++ {
		input.WriteString("q")
		input.WriteString(strconv.Itoa(quoted))
		input.WriteString(",\"quoted\"\n")
		input.WriteString("u")
		input.WriteString(strconv.Itoa(quoted))
		input.WriteString(",plain\n")
	}
	if err := os.WriteFile(path, []byte(input.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	document, err := OpenDocument(path, OpenOptions{PreviewRows: 3, IndexRows: -1})
	if err != nil {
		t.Fatal(err)
	}
	dialect := document.Metadata.Dialect
	if dialect.QuoteMapDetected != MaxQuoteCellRules+1 || dialect.QuoteMapLimit != MaxQuoteCellRules || !dialect.QuoteMapTruncated {
		t.Fatalf("quote map metadata detected=%d limit=%d truncated=%v, want detected %d limit %d truncated",
			dialect.QuoteMapDetected, dialect.QuoteMapLimit, dialect.QuoteMapTruncated, MaxQuoteCellRules+1, MaxQuoteCellRules)
	}
	if len(dialect.QuoteCells) != MaxQuoteCellRules {
		t.Fatalf("quote cells = %d, want capped %d", len(dialect.QuoteCells), MaxQuoteCellRules)
	}
	if !PreservesQuote(dialect, 1, 1) {
		t.Fatal("first quoted cell should be preserved")
	}
	truncatedRow := 1 + (MaxQuoteCellRules * 2)
	if PreservesQuote(dialect, truncatedRow, 1) {
		t.Fatalf("truncated quote cell row %d should not be in capped quote map", truncatedRow)
	}
}
