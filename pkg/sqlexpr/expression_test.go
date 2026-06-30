package sqlexpr

import (
	"strings"
	"testing"
)

func TestMatchHeaderNamesAndCColumns(t *testing.T) {
	headers := []string{"id", "First Name", "amount", "status"}
	record := []string{"42", "John Smith", "125.50", "Active"}

	tests := []struct {
		name string
		expr string
		want bool
	}{
		{name: "header and numeric comparison", expr: "amount > 100 AND status = 'Active'", want: true},
		{name: "quoted header and like", expr: `"First Name" LIKE 'John%'`, want: true},
		{name: "c column alias", expr: "c1 = 42", want: true},
		{name: "case function", expr: "LOWER(status) = 'active'", want: true},
		{name: "negative branch", expr: "amount < 100 OR status = 'Archived'", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Match(tt.expr, headers, record)
			if err != nil {
				t.Fatalf("Match returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("Match(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestLikeWildcardSemantics(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		pattern string
		want    bool
	}{
		{name: "percent matches any run", text: "John Smith", pattern: "John%", want: true},
		{name: "underscore matches one rune", text: "Jo", pattern: "J_", want: true},
		{name: "underscore does not match zero runes", text: "J", pattern: "J_", want: false},
		{name: "regex dot is literal", text: "aXb", pattern: "a.b", want: false},
		{name: "regex brackets are literal", text: "ab", pattern: "a[bc]", want: false},
		{name: "unicode underscore", text: "a界b", pattern: "a_b", want: true},
		{name: "collapsed percents", text: "abcdef", pattern: "a%%%f", want: true},
		{name: "not enough suffix", text: "aaaaaac", pattern: "%a%a%a%ab", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := like(tt.text, tt.pattern); got != tt.want {
				t.Fatalf("like(%q, %q) = %v, want %v", tt.text, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestLikeEscapeSemantics(t *testing.T) {
	headers := []string{"code"}
	tests := []struct {
		expr   string
		record []string
		want   bool
	}{
		{expr: "code LIKE 'APAC!_%' ESCAPE '!'", record: []string{"APAC_10"}, want: true},
		{expr: "code LIKE 'APAC!_%' ESCAPE '!'", record: []string{"APACX10"}, want: false},
		{expr: "code NOT LIKE 'APAC!_%' ESCAPE '!'", record: []string{"APACX10"}, want: true},
		{expr: `code LIKE '100\%%' ESCAPE '\'`, record: []string{"100% done"}, want: true},
	}
	for _, tt := range tests {
		got, err := Match(tt.expr, headers, tt.record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", tt.expr, err)
		}
		if got != tt.want {
			t.Fatalf("Match(%q, %#v) = %v, want %v", tt.expr, tt.record, got, tt.want)
		}
	}
}

func TestLikeRepeatedWildcardPatternIsBounded(t *testing.T) {
	text := strings.Repeat("a", 512) + "c"
	pattern := strings.Repeat("%a", 80) + "%b"
	for range 1000 {
		if like(text, pattern) {
			t.Fatal("repeated LIKE wildcard pattern matched unexpected suffix")
		}
	}
}

func TestMatchCastTryCastAndNulls(t *testing.T) {
	headers := []string{"id", "amount", "bad", "created_at", "seen_at", "duration", "age"}
	record := []string{"7", "19.75", "nope", "2026-05-27", "15:30:00.123", "1 02:30 day to minute", "1-2 year to month"}

	tests := []string{
		"CAST(amount AS number) >= 19.75",
		"CAST(id AS bigint) = 7",
		"TRY_CAST(bad AS number) IS NULL",
		"TRY_CAST(bad AS bigint) IS NULL",
		"CAST(created_at AS date) >= DATE '2026-01-01'",
		"CAST(seen_at AS time) = TIME '15:30:00.123'",
		"TIMESTAMP '2026-05-27 15:30:00.123' > TIMESTAMP '2026-05-27 15:29:59.999'",
		"CAST(duration AS interval) = INTERVAL '1 02:30' day to minute",
		"CAST(age AS interval) = INTERVAL '1-2' year to month",
		"INTERVAL '30' second > INTERVAL '10' second",
	}

	for _, source := range tests {
		got, err := Match(source, headers, record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", source, err)
		}
		if !got {
			t.Fatalf("Match(%q) = false, want true", source)
		}
	}
}

func TestMatchBooleanIsPredicates(t *testing.T) {
	headers := []string{"flag", "disabled", "blank", "amount"}
	record := []string{"true", "0", "", "12"}

	tests := []struct {
		expr string
		want bool
	}{
		{expr: "flag IS TRUE", want: true},
		{expr: "flag IS NOT FALSE", want: true},
		{expr: "disabled IS FALSE", want: true},
		{expr: "disabled IS NOT TRUE", want: true},
		{expr: "blank IS FALSE", want: true},
		{expr: "amount IS TRUE", want: true},
		{expr: "missing IS UNKNOWN", want: true},
		{expr: "missing IS NOT TRUE", want: true},
		{expr: "missing IS FALSE", want: false},
		{expr: "TRY_CAST(flag AS boolean) IS TRUE", want: true},
		{expr: "TRY_CAST('not-bool' AS boolean) IS UNKNOWN", want: true},
	}

	for _, tt := range tests {
		got, err := Match(tt.expr, headers, record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", tt.expr, err)
		}
		if got != tt.want {
			t.Fatalf("Match(%q) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestMatchBigIntPrecision(t *testing.T) {
	headers := []string{"id", "label"}
	record := []string{"9007199254740993", "Grace"}

	tests := []struct {
		expr string
		want bool
	}{
		{expr: "CAST(id AS bigint) = 9007199254740993", want: true},
		{expr: "CAST(id AS bigint) = 9007199254740992", want: false},
		{expr: "id > 9007199254740992", want: true},
		{expr: "9007199254740993 = id", want: true},
		{expr: "CAST(id AS bigint) + 1 = 9007199254740994", want: true},
		{expr: "-9007199254740993 < 0", want: true},
	}

	for _, tt := range tests {
		got, err := Match(tt.expr, headers, record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", tt.expr, err)
		}
		if got != tt.want {
			t.Fatalf("Match(%q) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestMatchInAndNotIn(t *testing.T) {
	headers := []string{"id", "name", "region", "score"}
	record := []string{"9007199254740993", "Grace", "APAC", "10"}

	tests := []struct {
		expr string
		want bool
	}{
		{expr: "region IN ('APAC', 'EMEA')", want: true},
		{expr: "region NOT IN ('AMER', 'EMEA')", want: true},
		{expr: "region NOT IN ('APAC', 'EMEA')", want: false},
		{expr: "score IN (8, 10)", want: true},
		{expr: "(CAST(score AS number) + 1) IN (10, 11)", want: true},
		{expr: "CAST(id AS bigint) IN (9007199254740992, 9007199254740993)", want: true},
		{expr: "CAST(id AS bigint) NOT IN (9007199254740993, 9007199254740994)", want: false},
	}

	for _, tt := range tests {
		got, err := Match(tt.expr, headers, record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", tt.expr, err)
		}
		if got != tt.want {
			t.Fatalf("Match(%q) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestMatchBetweenAndNotBetween(t *testing.T) {
	headers := []string{"id", "name", "region", "score", "created_at"}
	record := []string{"9007199254740993", "Grace", "APAC", "10", "2026-05-27"}

	tests := []struct {
		expr string
		want bool
	}{
		{expr: "score BETWEEN 8 AND 10", want: true},
		{expr: "score NOT BETWEEN 8 AND 10", want: false},
		{expr: "score NOT BETWEEN 11 AND 20", want: true},
		{expr: "CAST(score AS number) + 1 BETWEEN 10 AND 12", want: true},
		{expr: "score BETWEEN 8 AND 10 AND region = 'APAC'", want: true},
		{expr: "score BETWEEN 8 AND 10 AND region = 'EMEA'", want: false},
		{expr: "CAST(created_at AS date) BETWEEN DATE '2026-01-01' AND DATE '2026-12-31'", want: true},
		{expr: "CAST(id AS bigint) BETWEEN 9007199254740992 AND 9007199254740994", want: true},
	}

	for _, tt := range tests {
		got, err := Match(tt.expr, headers, record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", tt.expr, err)
		}
		if got != tt.want {
			t.Fatalf("Match(%q) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestMatchTemporalArithmetic(t *testing.T) {
	headers := []string{"created_at", "seen_at", "stamp"}
	record := []string{"2026-05-27", "15:30:00.123", "2026-05-27 14:30:00"}

	tests := []struct {
		expr string
		want bool
	}{
		{expr: "CAST(created_at AS date) + INTERVAL '2' day = DATE '2026-05-29'", want: true},
		{expr: "DATE '2026-05-29' - INTERVAL '2' day = CAST(created_at AS date)", want: true},
		{expr: "DATE '2026-05-29' - CAST(created_at AS date) = INTERVAL '2' day", want: true},
		{expr: "CAST(created_at AS date) + INTERVAL '1' month = DATE '2026-06-27'", want: true},
		{expr: "CAST(created_at AS date) - INTERVAL '2' month = DATE '2026-03-27'", want: true},
		{expr: "CAST(seen_at AS time) + INTERVAL '30' minute = TIME '16:00:00.123'", want: true},
		{expr: "TIME '16:00:00.123' - CAST(seen_at AS time) = INTERVAL '30' minute", want: true},
		{expr: "CAST(stamp AS timestamp) + INTERVAL '1:30' hour to minute = TIMESTAMP '2026-05-27 16:00:00'", want: true},
		{expr: "CAST(stamp AS timestamp) + INTERVAL '1-2' year to month = TIMESTAMP '2027-07-27 14:30:00'", want: true},
		{expr: "TIMESTAMP '2026-05-27 16:00:00' - CAST(stamp AS timestamp) = INTERVAL '1:30' hour to minute", want: true},
		{expr: "CAST(stamp AS timestamp) + INTERVAL '-1:30' hour to minute = TIMESTAMP '2026-05-27 13:00:00'", want: true},
		{expr: "CAST(created_at AS date) + INTERVAL '-1 02:30' day to minute = TIMESTAMP '2026-05-25 21:30:00'", want: true},
		{expr: "CAST(stamp AS timestamp) + INTERVAL '+1:30' hour to minute = TIMESTAMP '2026-05-27 16:00:00'", want: true},
		{expr: "INTERVAL '-1-2' year to month = INTERVAL '-14' month", want: true},
		{expr: "INTERVAL '+1-2' year to month = INTERVAL '14' month", want: true},
		{expr: "INTERVAL '-1:30' minute to second = INTERVAL '-90' second", want: true},
		{expr: "DATE '2026-05-27' = '2026-05-27'", want: false},
	}

	for _, tt := range tests {
		got, err := Match(tt.expr, headers, record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", tt.expr, err)
		}
		if got != tt.want {
			t.Fatalf("Match(%q) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestMatchIntervalPrecisionQualifiers(t *testing.T) {
	headers := []string{"duration", "age"}
	record := []string{"1 02:03:04.123 day(2) to second(3)", "1-2 year(4) to month"}

	tests := []string{
		"INTERVAL '1.234' second(3) = INTERVAL '1.234' second",
		"INTERVAL '1 02:03:04.123' day(2) to second(3) = INTERVAL '1 02:03:04.123' day to second",
		"INTERVAL '1:30' hour(2) to minute = INTERVAL '90' minute",
		"INTERVAL '1-2' year(4) to month = INTERVAL '14' month",
		"CAST(duration AS interval) = INTERVAL '1 02:03:04.123' day to second",
		"CAST(age AS interval) = INTERVAL '14' month",
	}

	for _, source := range tests {
		got, err := Match(source, headers, record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", source, err)
		}
		if !got {
			t.Fatalf("Match(%q) = false, want true", source)
		}
	}
}

func TestMatchFunctionsAndArithmetic(t *testing.T) {
	headers := []string{"name", "amount"}
	record := []string{"  Alice  ", "9"}

	tests := []string{
		"TRIM(name) = 'Alice'",
		"LENGTH(TRIM(name)) = 5",
		"SUBSTRING(TRIM(name), 1, 3) = 'Ali'",
		"CONCAT(TRIM(name), '-ok') = 'Alice-ok'",
		"REPLACE(TRIM(name), 'Alice', 'Ada') = 'Ada'",
		"ROUND(POWER(amount, 2) / 10, 1) = 8.1",
		"MOD(amount, 2) = 1",
		"COALESCE(NULLIF(status, ''), 'fallback') = 'fallback'",
	}

	for _, source := range tests {
		got, err := Match(source, headers, record)
		if err != nil {
			t.Fatalf("Match(%q) returned error: %v", source, err)
		}
		if !got {
			t.Fatalf("Match(%q) = false, want true", source)
		}
	}
}

func TestMatchReportsInvalidSyntax(t *testing.T) {
	if _, err := Match("amount >", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected syntax error")
	}
	if _, err := Match("CAST(amount number)", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected cast syntax error")
	}
	if _, err := Match("INTERVAL '1' banana = INTERVAL '1' day", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected interval field syntax error")
	}
	if _, err := Match("INTERVAL '1:02:03' hour to minute = INTERVAL '62' minute", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected hour-to-minute interval syntax error")
	}
	if _, err := Match("INTERVAL '1' day(foo) = INTERVAL '1' day", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected interval precision syntax error")
	}
	if _, err := Match("INTERVAL '1' second(3, 2) = INTERVAL '1' second", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected interval precision scale syntax error")
	}
	if _, err := Match("code LIKE 'APAC!_' ESCAPE '!!'", []string{"code"}, []string{"APAC_"}); err == nil {
		t.Fatal("expected like escape width error")
	}
	if _, err := Match("code LIKE 'APAC!' ESCAPE '!'", []string{"code"}, []string{"APAC!"}); err == nil {
		t.Fatal("expected dangling like escape error")
	}
	if _, err := Match("amount IN ()", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected in list syntax error")
	}
	if _, err := Match("amount BETWEEN 1", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected between syntax error")
	}
	if _, err := Match("amount IS MAYBE", []string{"amount"}, []string{"10"}); err == nil {
		t.Fatal("expected is-predicate syntax error")
	}
}
