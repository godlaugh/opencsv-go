package dataprofile

import "testing"

func TestAnalyzeProfilesColumnsAndIssues(t *testing.T) {
	records := [][]string{
		{"id", "", "score", "score"},
		{"1", "Ada", "9", "9"},
		{"2", "Linus", "7", "7"},
		{"3", "Grace", "bad", "10"},
		{"4", "Alan"},
	}

	result := Analyze(records, true)
	if result.Rows != 4 {
		t.Fatalf("rows = %d, want 4", result.Rows)
	}
	if len(result.Columns) != 4 {
		t.Fatalf("columns = %d, want 4", len(result.Columns))
	}
	if result.Columns[0].Type != "integer" || result.Columns[0].Distinct != 4 {
		t.Fatalf("unexpected id profile: %#v", result.Columns[0])
	}
	if result.Columns[2].Type != "mixed" {
		t.Fatalf("score type = %s, want mixed", result.Columns[2].Type)
	}
	if !hasIssue(result.Issues, "empty-header") {
		t.Fatalf("expected empty-header issue: %#v", result.Issues)
	}
	if !hasIssue(result.Issues, "duplicate-header") {
		t.Fatalf("expected duplicate-header issue: %#v", result.Issues)
	}
	if !hasIssue(result.Issues, "ragged-row") {
		t.Fatalf("expected ragged-row issue: %#v", result.Issues)
	}
}

func TestAnalyzeNumericSummary(t *testing.T) {
	result := Analyze([][]string{
		{"amount"},
		{"1.5"},
		{"2.5"},
		{""},
	}, true)
	amount := result.Columns[0]
	if amount.Type != "decimal" {
		t.Fatalf("type = %s, want decimal", amount.Type)
	}
	if amount.Nulls != 1 {
		t.Fatalf("nulls = %d, want 1", amount.Nulls)
	}
	if amount.Sum == nil || *amount.Sum != 4 {
		t.Fatalf("sum = %#v, want 4", amount.Sum)
	}
	if amount.Average == nil || *amount.Average != 2 {
		t.Fatalf("average = %#v, want 2", amount.Average)
	}
	if amount.Min != "1.5" || amount.Max != "2.5" {
		t.Fatalf("min/max = %q/%q", amount.Min, amount.Max)
	}
}

func hasIssue(issues []Issue, typ string) bool {
	for _, issue := range issues {
		if issue.Type == typ {
			return true
		}
	}
	return false
}
