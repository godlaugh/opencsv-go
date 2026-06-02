package dataprofile

import (
	"encoding/json"
	"testing"
)

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

func TestAnalyzeAddsQuickFixesForSafeStructureIssues(t *testing.T) {
	result := Analyze([][]string{
		{"id", ""},
		{"1"},
		{"2", "Ada"},
	}, true)
	emptyHeader := issueByType(result.Issues, "empty-header")
	if emptyHeader == nil || len(emptyHeader.QuickFixes) != 1 || emptyHeader.QuickFixes[0].ID != "header.fill" {
		t.Fatalf("empty header quick fix = %#v", emptyHeader)
	}
	ragged := issueByType(result.Issues, "ragged-row")
	if ragged == nil || ragged.Actual != 1 || ragged.Expected != 2 {
		t.Fatalf("ragged issue = %#v, want actual/expected", ragged)
	}
	if len(ragged.QuickFixes) != 1 || ragged.QuickFixes[0].ID != "row.pad" {
		t.Fatalf("ragged quick fix = %#v", ragged.QuickFixes)
	}
}

func TestAnalyzeWithSchemaRules(t *testing.T) {
	minScore := 0.0
	maxScore := 100.0
	result := AnalyzeWithSchema([][]string{
		{"id", "status", "score"},
		{"1", "new", "10"},
		{"1", "bad", "105"},
		{"", "done", "ok"},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "id", Required: true, Type: "integer", Unique: true},
		{Name: "status", Enum: []string{"new", "done"}, Pattern: "^[a-z]+$"},
		{Name: "score", Type: "number", Min: &minScore, Max: &maxScore},
		{Name: "missing", Required: true},
	}})
	if result.SchemaRules != 4 {
		t.Fatalf("schema rules = %d, want 4", result.SchemaRules)
	}
	for _, typ := range []string{"schema-unique", "schema-required", "schema-enum", "schema-range", "schema-type", "schema-missing-column"} {
		if !hasIssue(result.Issues, typ) {
			t.Fatalf("expected %s issue in %#v", typ, result.Issues)
		}
	}
}

func TestSchemaUnmarshalAcceptsCSVLintFields(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"fields": [
			{"name": "id", "type": "integer", "constraints": {"required": true, "unique": true}},
			{"name": "status", "constraints": {"enum": ["new", "done"]}},
			{"name": "score", "datatype": {"base": "number", "minimum": "0", "maximum": 100}}
		],
		"primaryKey": "id"
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 3 {
		t.Fatalf("columns = %#v, want 3", schema.Columns)
	}
	id := schema.Columns[0]
	if id.Name != "id" || id.Type != "integer" || !id.Required || !id.Unique {
		t.Fatalf("id rule = %#v", id)
	}
	status := schema.Columns[1]
	if len(status.Enum) != 2 || status.Enum[0] != "new" || status.Enum[1] != "done" {
		t.Fatalf("status enum = %#v", status.Enum)
	}
	score := schema.Columns[2]
	if score.Type != "number" || score.Min == nil || *score.Min != 0 || score.Max == nil || *score.Max != 100 {
		t.Fatalf("score rule = %#v", score)
	}
}

func TestSchemaUnmarshalAcceptsNestedTableSchema(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"tables": [{
			"schema": {
				"fields": [
					{"title": "code", "datatype": "string", "constraints": {"pattern": "^A"}},
					{"name": "amount", "dataType": "float", "constraints": {"minimum": 1}}
				],
				"required": ["amount"]
			}
		}]
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 2 {
		t.Fatalf("columns = %#v, want 2", schema.Columns)
	}
	if schema.Columns[0].Name != "code" || schema.Columns[0].Type != "string" || schema.Columns[0].Pattern != "^A" {
		t.Fatalf("code rule = %#v", schema.Columns[0])
	}
	if schema.Columns[1].Name != "amount" || schema.Columns[1].Type != "number" || !schema.Columns[1].Required || schema.Columns[1].Min == nil || *schema.Columns[1].Min != 1 {
		t.Fatalf("amount rule = %#v", schema.Columns[1])
	}
}

func TestAnalyzeWithSchemaInvalidPattern(t *testing.T) {
	result := AnalyzeWithSchema([][]string{
		{"code"},
		{"A1"},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "code", Pattern: "["},
	}})
	issue := issueByType(result.Issues, "schema-pattern")
	if issue == nil || issue.Row != 1 || issue.Col != 1 {
		t.Fatalf("pattern issue = %#v", issue)
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

func issueByType(issues []Issue, typ string) *Issue {
	for i := range issues {
		if issues[i].Type == typ {
			return &issues[i]
		}
	}
	return nil
}
