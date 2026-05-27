package csvcore

import "testing"

func TestRequiresQuote(t *testing.T) {
	dialect := Dialect{Delimiter: ','}
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
