package csvcore

import "testing"

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
}
