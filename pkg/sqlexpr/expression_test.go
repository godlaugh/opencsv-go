package sqlexpr

import "testing"

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

func TestMatchCastTryCastAndNulls(t *testing.T) {
	headers := []string{"id", "amount", "bad", "created_at"}
	record := []string{"7", "19.75", "nope", "2026-05-27"}

	tests := []string{
		"CAST(amount AS number) >= 19.75",
		"CAST(id AS bigint) = 7",
		"TRY_CAST(bad AS number) IS NULL",
		"CAST(created_at AS date) >= DATE '2026-01-01'",
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
}
