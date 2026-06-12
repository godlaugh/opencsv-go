package dataprofile

import (
	"context"
	"encoding/json"
	"strings"
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

func TestAnalyzeStreamProfilesColumnsAndSchema(t *testing.T) {
	minScore := 0.0
	maxScore := 100.0
	records := [][]string{
		{"id", "", "score", "score"},
		{"1", "Ada", "9", "9"},
		{"1", "bad", "105", "10"},
		{"", "Grace", "bad", "11"},
		{"4", "Alan"},
	}
	result, err := AnalyzeStreamContext(context.Background(), true, Schema{Columns: []ColumnRule{
		{Name: "id", Required: true, Type: "integer", Unique: true},
		{Name: "score", Type: "number", Min: &minScore, Max: &maxScore},
		{Name: "missing", Required: true},
	}}, streamRecords(records))
	if err != nil {
		t.Fatal(err)
	}
	if !result.Streamed {
		t.Fatal("streamed = false, want true")
	}
	if result.Rows != 4 || len(result.Columns) != 4 {
		t.Fatalf("shape = rows %d columns %d, want 4/4", result.Rows, len(result.Columns))
	}
	if result.Columns[0].Type != "integer" || result.Columns[0].Distinct != 2 {
		t.Fatalf("id profile = %#v, want integer with two distinct non-empty values", result.Columns[0])
	}
	if result.Columns[2].Type != "mixed" {
		t.Fatalf("score type = %s, want mixed", result.Columns[2].Type)
	}
	for _, typ := range []string{"empty-header", "duplicate-header", "ragged-row", "schema-unique", "schema-required", "schema-range", "schema-type", "schema-missing-column"} {
		if !hasIssue(result.Issues, typ) {
			t.Fatalf("expected %s issue in %#v", typ, result.Issues)
		}
	}
}

func TestAnalyzeStreamReportsRowProgress(t *testing.T) {
	records := [][]string{
		{"id", "name"},
		{"1", "Ada"},
		{"2", "Linus"},
		{"3", "Grace"},
	}
	var events []ProgressEvent
	result, err := AnalyzeStreamContextWithOptions(context.Background(), true, Schema{}, streamRecords(records), StreamOptions{
		ProgressInterval: 2,
		ReportProgress: func(event ProgressEvent) {
			events = append(events, event)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Progress == nil {
		t.Fatalf("progress = nil, want summary")
	}
	if result.Progress.ProfileRows != 3 || result.Progress.ValidationRows != 3 || result.Progress.Passes != 2 {
		t.Fatalf("progress = %#v, want 3 profile rows and 3 validation rows", result.Progress)
	}
	if len(events) != 4 {
		t.Fatalf("events = %#v, want checkpoint and done for both passes", events)
	}
	if events[0].Phase != "profile" || events[0].Rows != 2 || events[0].Done {
		t.Fatalf("first event = %#v, want profile checkpoint at 2 rows", events[0])
	}
	if events[1].Phase != "profile" || events[1].Rows != 3 || !events[1].Done {
		t.Fatalf("second event = %#v, want profile done at 3 rows", events[1])
	}
	if events[2].Phase != "validation" || events[2].Rows != 2 || events[2].Done {
		t.Fatalf("third event = %#v, want validation checkpoint at 2 rows", events[2])
	}
	if events[3].Phase != "validation" || events[3].Rows != 3 || !events[3].Done {
		t.Fatalf("fourth event = %#v, want validation done at 3 rows", events[3])
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

func TestTopValuesKeepsHighestCountsAndLexicalTies(t *testing.T) {
	values := topValues(map[string]int{
		"zeta":  2,
		"alpha": 2,
		"beta":  3,
		"gamma": 1,
	}, 3)
	want := []ValueFrequency{
		{Value: "beta", Count: 3},
		{Value: "alpha", Count: 2},
		{Value: "zeta", Count: 2},
	}
	if len(values) != len(want) {
		t.Fatalf("top values = %#v, want %#v", values, want)
	}
	for i := range want {
		if values[i] != want[i] {
			t.Fatalf("top values = %#v, want %#v", values, want)
		}
	}
}

func streamRecords(records [][]string) RecordStream {
	return func(ctx context.Context, handle func(row int, record []string) error) (int, error) {
		for row, record := range records {
			if err := ctx.Err(); err != nil {
				return row, err
			}
			if err := handle(row, record); err != nil {
				return row, err
			}
		}
		return len(records), nil
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
	duplicates := Analyze([][]string{
		{"id", "score", "score"},
		{"1", "9", "10"},
	}, true)
	duplicateHeader := issueByType(duplicates.Issues, "duplicate-header")
	if duplicateHeader == nil || len(duplicateHeader.QuickFixes) != 1 || duplicateHeader.QuickFixes[0].ID != "header.unique" {
		t.Fatalf("duplicate header quick fix = %#v", duplicateHeader)
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
	rangeIssue := issueByType(result.Issues, "schema-range")
	if rangeIssue == nil || len(rangeIssue.QuickFixes) != 1 {
		t.Fatalf("schema range quick fix = %#v, want one clamp fix", rangeIssue)
	}
	if rangeIssue.QuickFixes[0].ID != "schema.range.clampMax" || rangeIssue.QuickFixes[0].Value != "100" {
		t.Fatalf("schema range quick fix = %#v, want clampMax to 100", rangeIssue.QuickFixes)
	}
}

func TestAnalyzeWithSchemaRangeQuickFixes(t *testing.T) {
	minScore := 0.0
	maxScore := 10.0
	result := AnalyzeWithSchema([][]string{
		{"score"},
		{"-1"},
		{"11"},
		{"bad"},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "score", Type: "number", Min: &minScore, Max: &maxScore},
	}})
	var clampMin, clampMax, nonNumeric *Issue
	for i := range result.Issues {
		issue := &result.Issues[i]
		if issue.Type != "schema-range" {
			continue
		}
		switch issue.Row {
		case 2:
			clampMin = issue
		case 3:
			clampMax = issue
		case 4:
			nonNumeric = issue
		}
	}
	if clampMin == nil || len(clampMin.QuickFixes) != 1 || clampMin.QuickFixes[0].ID != "schema.range.clampMin" || clampMin.QuickFixes[0].Value != "0" {
		t.Fatalf("min range issue = %#v, want clampMin to 0", clampMin)
	}
	if clampMax == nil || len(clampMax.QuickFixes) != 1 || clampMax.QuickFixes[0].ID != "schema.range.clampMax" || clampMax.QuickFixes[0].Value != "10" {
		t.Fatalf("max range issue = %#v, want clampMax to 10", clampMax)
	}
	if nonNumeric == nil || len(nonNumeric.QuickFixes) != 0 {
		t.Fatalf("non-numeric range issue = %#v, want no semantic quick fix", nonNumeric)
	}
}

func TestAnalyzeWithSchemaExclusiveRange(t *testing.T) {
	minScore := 0.0
	maxScore := 10.0
	result := AnalyzeWithSchema([][]string{
		{"score"},
		{"0"},
		{"10"},
		{"5"},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "score", Type: "number", Min: &minScore, Max: &maxScore, ExclusiveMin: true, ExclusiveMax: true},
	}})
	var minIssue, maxIssue *Issue
	for i := range result.Issues {
		issue := &result.Issues[i]
		if issue.Type != "schema-range" {
			continue
		}
		switch issue.Row {
		case 2:
			minIssue = issue
		case 3:
			maxIssue = issue
		}
	}
	if minIssue == nil || len(minIssue.QuickFixes) != 0 || !strings.Contains(minIssue.Message, "greater than schema minimum 0") {
		t.Fatalf("exclusive min issue = %#v, want strict lower-bound issue without clamp", minIssue)
	}
	if maxIssue == nil || len(maxIssue.QuickFixes) != 0 || !strings.Contains(maxIssue.Message, "less than schema maximum 10") {
		t.Fatalf("exclusive max issue = %#v, want strict upper-bound issue without clamp", maxIssue)
	}
}

func TestAnalyzeWithSchemaSingleEnumQuickFix(t *testing.T) {
	result := AnalyzeWithSchema([][]string{
		{"status", "category"},
		{"pending", "old"},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "status", Enum: []string{"active"}},
		{Name: "category", Enum: []string{"new", "done"}},
	}})
	var singleEnum, multiEnum *Issue
	for i := range result.Issues {
		issue := &result.Issues[i]
		if issue.Type != "schema-enum" {
			continue
		}
		switch issue.Col {
		case 1:
			singleEnum = issue
		case 2:
			multiEnum = issue
		}
	}
	if singleEnum == nil || len(singleEnum.QuickFixes) != 1 || singleEnum.QuickFixes[0].ID != "schema.enum.useValue" || singleEnum.QuickFixes[0].Value != "active" {
		t.Fatalf("single enum issue = %#v, want useValue active quick fix", singleEnum)
	}
	if multiEnum == nil || len(multiEnum.QuickFixes) != 0 {
		t.Fatalf("multi enum issue = %#v, want no guessed quick fix", multiEnum)
	}
}

func TestAnalyzeWithSchemaRequiredDefaultQuickFix(t *testing.T) {
	result := AnalyzeWithSchema([][]string{
		{"status", "owner"},
		{"", ""},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "status", Required: true, Default: "active"},
		{Name: "owner", Required: true},
	}})
	var defaulted, noDefault *Issue
	for i := range result.Issues {
		issue := &result.Issues[i]
		if issue.Type != "schema-required" {
			continue
		}
		switch issue.Col {
		case 1:
			defaulted = issue
		case 2:
			noDefault = issue
		}
	}
	if defaulted == nil || len(defaulted.QuickFixes) != 1 || defaulted.QuickFixes[0].ID != "schema.default.use" || defaulted.QuickFixes[0].Value != "active" {
		t.Fatalf("defaulted required issue = %#v, want default active quick fix", defaulted)
	}
	if noDefault == nil || len(noDefault.QuickFixes) != 0 {
		t.Fatalf("required issue without default = %#v, want no quick fix", noDefault)
	}
}

func TestAnalyzeWithSchemaMissingColumnQuickFix(t *testing.T) {
	result := AnalyzeWithSchema([][]string{
		{"id", "owner"},
		{"1", "Ada"},
		{"2", "Linus"},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "status", Required: true, Enum: []string{"active"}, Default: "active"},
		{Name: "score", Type: "integer", Default: "bad"},
	}})
	statusIssue := issueByRule(result.Issues, "schema-missing-column", "status")
	if statusIssue == nil || len(statusIssue.QuickFixes) != 1 {
		t.Fatalf("status missing-column issue = %#v, want add-column quick fix", statusIssue)
	}
	statusFix := statusIssue.QuickFixes[0]
	if statusFix.ID != "schema.column.add" || statusFix.Value != "active" {
		t.Fatalf("status missing-column fix = %#v, want add column with default active", statusFix)
	}
	scoreIssue := issueByRule(result.Issues, "schema-missing-column", "score")
	if scoreIssue == nil || len(scoreIssue.QuickFixes) != 1 {
		t.Fatalf("score missing-column issue = %#v, want add-column quick fix", scoreIssue)
	}
	scoreFix := scoreIssue.QuickFixes[0]
	if scoreFix.ID != "schema.column.add" || scoreFix.Value != "" {
		t.Fatalf("score missing-column fix = %#v, want add column without invalid default", scoreFix)
	}
	noHeader := AnalyzeWithSchema([][]string{{"1", "Ada"}}, false, Schema{Columns: []ColumnRule{
		{Name: "status", Required: true, Default: "active"},
	}})
	noHeaderIssue := issueByRule(noHeader.Issues, "schema-missing-column", "status")
	if noHeaderIssue == nil || len(noHeaderIssue.QuickFixes) != 0 {
		t.Fatalf("no-header missing-column issue = %#v, want no structural quick fix", noHeaderIssue)
	}
}

func TestAnalyzeWithSchemaSemanticDefaultQuickFixes(t *testing.T) {
	minLength := 2
	maxLength := 4
	min := 0.0
	max := 10.0
	result := AnalyzeWithSchema([][]string{
		{"age", "code", "short", "long", "score", "status", "unique"},
		{"n/a", "bad", "x", "too-long", "oops", "pending", "bad"},
		{"2", "A-2", "ok", "fine", "5", "active", "fallback"},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "age", Type: "integer", Default: "42"},
		{Name: "code", Pattern: "^A-", Default: "A-1"},
		{Name: "short", MinLength: &minLength, Default: "ok"},
		{Name: "long", MaxLength: &maxLength, Default: "fine"},
		{Name: "score", Min: &min, Max: &max, Default: "5"},
		{Name: "status", Enum: []string{"active", "paused"}, Default: "active"},
		{Name: "unique", Type: "integer", Unique: true, Default: "1"},
	}})
	for _, typ := range []string{"schema-type", "schema-pattern", "schema-min-length", "schema-max-length", "schema-range", "schema-enum"} {
		issue := issueByType(result.Issues, typ)
		if issue == nil || !hasQuickFix(issue, "schema.default.use") {
			t.Fatalf("%s issue = %#v, want schema.default.use quick fix in %#v", typ, issue, result.Issues)
		}
	}
	uniqueIssue := issueByTypeCol(result.Issues, "schema-type", 7)
	if uniqueIssue == nil {
		t.Fatalf("expected unique column type issue in %#v", result.Issues)
	}
	if hasQuickFix(uniqueIssue, "schema.default.use") {
		t.Fatalf("unique column issue = %#v, should not offer default quick fix", uniqueIssue)
	}
}

func TestAnalyzeWithSchemaStringLengthRules(t *testing.T) {
	minLength := 2
	maxLength := 3
	result := AnalyzeWithSchema([][]string{
		{"code"},
		{"A"},
		{"猫狗"},
		{"猫狗鸟鱼"},
	}, true, Schema{Columns: []ColumnRule{
		{Name: "code", MinLength: &minLength, MaxLength: &maxLength},
	}})
	if result.SchemaRules != 1 {
		t.Fatalf("schema rules = %d, want 1", result.SchemaRules)
	}
	if !hasIssue(result.Issues, "schema-min-length") {
		t.Fatalf("expected schema-min-length issue in %#v", result.Issues)
	}
	if !hasIssue(result.Issues, "schema-max-length") {
		t.Fatalf("expected schema-max-length issue in %#v", result.Issues)
	}
	if lengthIssueCount := countIssues(result.Issues, "schema-min-length") + countIssues(result.Issues, "schema-max-length"); lengthIssueCount != 2 {
		t.Fatalf("length issues = %d, want 2: %#v", lengthIssueCount, result.Issues)
	}
}

func TestSchemaUnmarshalAcceptsCSVLintFields(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"fields": [
			{"name": "id", "type": "integer", "constraints": {"required": true, "unique": true}},
			{"name": "status", "constraints": {"enum": ["new", "done"]}},
			{"name": "score", "datatype": {"base": "number", "minimum": "0", "maximum": 100}},
			{"name": "code", "datatype": "string", "constraints": {"minLength": "2", "max_length": 4}}
		],
		"primaryKey": "id"
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 4 {
		t.Fatalf("columns = %#v, want 4", schema.Columns)
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
	code := schema.Columns[3]
	if code.Type != "string" || code.MinLength == nil || *code.MinLength != 2 || code.MaxLength == nil || *code.MaxLength != 4 {
		t.Fatalf("code rule = %#v", code)
	}
}

func TestSchemaUnmarshalAcceptsCSVLintTypeURLs(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"fields": [
			{"name": "id", "constraints": {"type": "http://www.w3.org/2001/XMLSchema#integer"}},
			{"name": "seen_at", "constraints": {"type": "http://www.w3.org/2001/XMLSchema#time"}},
			{"name": "period", "constraints": {"type": "http://www.w3.org/2001/XMLSchema#gYearMonth"}},
			{"name": "year", "constraints": {"type": "http://www.w3.org/2001/XMLSchema#gYear"}},
			{"name": "stamp", "constraints": {"type": "http://www.w3.org/2001/XMLSchema#dateTime"}},
			{"name": "amount", "constraints": {"type": "http://www.w3.org/2001/XMLSchema#decimal"}},
			{"name": "flag", "constraints": {"type": "http://www.w3.org/2001/XMLSchema#boolean"}},
			{"name": "count", "constraints": {"type": "xsd:unsignedInt"}},
			{"name": "link", "constraints": {"type": "http://www.w3.org/2001/XMLSchema#anyURI"}}
		]
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 9 {
		t.Fatalf("columns = %#v, want 9", schema.Columns)
	}
	wantTypes := []string{"integer", "time", "yearmonth", "year", "datetime", "decimal", "boolean", "integer", "string"}
	for index, want := range wantTypes {
		if schema.Columns[index].Type != want {
			t.Fatalf("column %d type = %q, want %q in %#v", index, schema.Columns[index].Type, want, schema.Columns)
		}
	}
	result := AnalyzeWithSchema([][]string{
		{"id", "seen_at", "period", "year", "stamp", "amount", "flag", "count", "link"},
		{"1", "15:30:00", "2026-06", "2026", "2026-06-12T15:30:00Z", "12.5", "true", "3", "https://example.com"},
		{"bad", "25:61", "2026-13", "20X6", "2026-06-12", "oops", "maybe", "many", "not checked"},
	}, true, schema)
	if result.SchemaRules != 9 {
		t.Fatalf("schema rules = %d, want 9", result.SchemaRules)
	}
	if typeIssues := countIssues(result.Issues, "schema-type"); typeIssues != 8 {
		t.Fatalf("schema-type issues = %d, want 8: %#v", typeIssues, result.Issues)
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

func TestSchemaUnmarshalAcceptsFrictionlessResources(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"resources": [
			{
				"name": "scores",
				"schema": {
					"fields": [
						{"name": "status", "constraints": {"enum": ["active"], "default": "active"}},
						{"name": "score", "type": "number", "constraints": {"minimum": 0, "maximum": 100}}
					],
					"primaryKey": "status"
				}
			}
		]
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 2 {
		t.Fatalf("columns = %#v, want 2", schema.Columns)
	}
	status := schema.Columns[0]
	if status.Name != "status" || !status.Required || !status.Unique || len(status.Enum) != 1 || status.Enum[0] != "active" || status.Default != "active" {
		t.Fatalf("status rule = %#v", status)
	}
	score := schema.Columns[1]
	if score.Name != "score" || score.Type != "number" || score.Min == nil || *score.Min != 0 || score.Max == nil || *score.Max != 100 {
		t.Fatalf("score rule = %#v", score)
	}
}

func TestSchemaUnmarshalAcceptsJSONSchemaProperties(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"type": "object",
		"required": ["id", "status"],
		"properties": {
			"status": {"type": "string", "enum": ["active"], "default": "active"},
			"score": {"type": "number", "minimum": 0, "maximum": 100},
			"id": {"type": "integer"}
		}
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 3 {
		t.Fatalf("columns = %#v, want 3", schema.Columns)
	}
	id := schema.Columns[0]
	if id.Name != "id" || id.Type != "integer" || !id.Required {
		t.Fatalf("id rule = %#v", id)
	}
	score := schema.Columns[1]
	if score.Name != "score" || score.Type != "number" || score.Min == nil || *score.Min != 0 || score.Max == nil || *score.Max != 100 {
		t.Fatalf("score rule = %#v", score)
	}
	status := schema.Columns[2]
	if status.Name != "status" || status.Type != "string" || !status.Required || len(status.Enum) != 1 || status.Enum[0] != "active" || status.Default != "active" {
		t.Fatalf("status rule = %#v", status)
	}
	result := AnalyzeWithSchema([][]string{
		{"id", "score", "status"},
		{"1", "95", "active"},
		{"bad", "105", ""},
	}, true, schema)
	for _, typ := range []string{"schema-type", "schema-range", "schema-required"} {
		if !hasIssue(result.Issues, typ) {
			t.Fatalf("expected %s issue in %#v", typ, result.Issues)
		}
	}
}

func TestSchemaUnmarshalAcceptsJSONSchemaExclusiveRanges(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"type": "object",
		"properties": {
			"legacy": {"type": "number", "minimum": 0, "exclusiveMinimum": true, "maximum": 10, "exclusiveMaximum": true},
			"score": {"type": "number", "exclusiveMinimum": 0, "exclusiveMaximum": 10}
		}
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 2 {
		t.Fatalf("columns = %#v, want 2", schema.Columns)
	}
	for _, rule := range schema.Columns {
		if rule.Type != "number" || rule.Min == nil || *rule.Min != 0 || rule.Max == nil || *rule.Max != 10 || !rule.ExclusiveMin || !rule.ExclusiveMax {
			t.Fatalf("exclusive range rule = %#v", rule)
		}
	}
}

func TestSchemaUnmarshalAcceptsJSONSchemaFormats(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"type": "object",
		"properties": {
			"created_at": {"type": "string", "format": "date-time"},
			"day": {"type": "string", "format": "date"},
			"clock": {"type": "string", "format": "time"}
		}
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 3 {
		t.Fatalf("columns = %#v, want 3", schema.Columns)
	}
	want := []struct {
		name string
		typ  string
	}{
		{"clock", "time"},
		{"created_at", "datetime"},
		{"day", "date"},
	}
	for index, item := range want {
		if schema.Columns[index].Name != item.name || schema.Columns[index].Type != item.typ {
			t.Fatalf("column %d = %#v, want %s/%s", index, schema.Columns[index], item.name, item.typ)
		}
	}
	result := AnalyzeWithSchema([][]string{
		{"clock", "created_at", "day"},
		{"15:30:00", "2026-06-12T15:30:00Z", "2026-06-12"},
		{"25:61", "2026-06-12", "Friday"},
	}, true, schema)
	if typeIssues := countIssues(result.Issues, "schema-type"); typeIssues != 3 {
		t.Fatalf("schema-type issues = %d, want 3: %#v", typeIssues, result.Issues)
	}
}

func TestSchemaUnmarshalAcceptsJSONSchemaConstValues(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"type": "object",
		"properties": {
			"status": {"type": "string", "const": "active"},
			"approved": {"type": "boolean", "oneOf": [{"const": true}]}
		}
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 2 {
		t.Fatalf("columns = %#v, want 2", schema.Columns)
	}
	approved := schema.Columns[0]
	if approved.Name != "approved" || approved.Type != "boolean" || len(approved.Enum) != 1 || approved.Enum[0] != "true" {
		t.Fatalf("approved rule = %#v", approved)
	}
	status := schema.Columns[1]
	if status.Name != "status" || status.Type != "string" || len(status.Enum) != 1 || status.Enum[0] != "active" {
		t.Fatalf("status rule = %#v", status)
	}
	result := AnalyzeWithSchema([][]string{
		{"approved", "status"},
		{"false", "pending"},
	}, true, schema)
	approvedIssue := issueByTypeCol(result.Issues, "schema-enum", 1)
	if approvedIssue == nil || !hasQuickFix(approvedIssue, "schema.enum.useValue") || approvedIssue.QuickFixes[0].Value != "true" {
		t.Fatalf("approved issue = %#v, want use allowed true quick fix", approvedIssue)
	}
	statusIssue := issueByTypeCol(result.Issues, "schema-enum", 2)
	if statusIssue == nil || !hasQuickFix(statusIssue, "schema.enum.useValue") || statusIssue.QuickFixes[0].Value != "active" {
		t.Fatalf("status issue = %#v, want use allowed active quick fix", statusIssue)
	}
}

func TestSchemaUnmarshalAcceptsJSONSchemaNullableTypes(t *testing.T) {
	var schema Schema
	err := json.Unmarshal([]byte(`{
		"type": "object",
		"properties": {
			"id": {"type": ["integer", "null"]},
			"score": {"type": ["null", "number"]},
			"stamp": {"type": ["null", "string"], "format": "date-time"}
		}
	}`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Columns) != 3 {
		t.Fatalf("columns = %#v, want 3", schema.Columns)
	}
	want := []struct {
		name string
		typ  string
	}{
		{"id", "integer"},
		{"score", "number"},
		{"stamp", "datetime"},
	}
	for index, item := range want {
		if schema.Columns[index].Name != item.name || schema.Columns[index].Type != item.typ {
			t.Fatalf("column %d = %#v, want %s/%s", index, schema.Columns[index], item.name, item.typ)
		}
	}
	result := AnalyzeWithSchema([][]string{
		{"id", "score", "stamp"},
		{"1", "95", "2026-06-12T15:30:00Z"},
		{"bad", "oops", "2026-06-12"},
	}, true, schema)
	if result.SchemaRules != 3 {
		t.Fatalf("schema rules = %d, want 3", result.SchemaRules)
	}
	if typeIssues := countIssues(result.Issues, "schema-type"); typeIssues != 3 {
		t.Fatalf("schema-type issues = %d, want 3: %#v", typeIssues, result.Issues)
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

func countIssues(issues []Issue, typ string) int {
	count := 0
	for _, issue := range issues {
		if issue.Type == typ {
			count++
		}
	}
	return count
}

func issueByType(issues []Issue, typ string) *Issue {
	for i := range issues {
		if issues[i].Type == typ {
			return &issues[i]
		}
	}
	return nil
}

func issueByTypeCol(issues []Issue, typ string, col int) *Issue {
	for i := range issues {
		if issues[i].Type == typ && issues[i].Col == col {
			return &issues[i]
		}
	}
	return nil
}

func issueByRule(issues []Issue, typ string, rule string) *Issue {
	for i := range issues {
		if issues[i].Type == typ && issues[i].Rule == rule {
			return &issues[i]
		}
	}
	return nil
}

func hasQuickFix(issue *Issue, id string) bool {
	if issue == nil {
		return false
	}
	for _, fix := range issue.QuickFixes {
		if fix.ID == id {
			return true
		}
	}
	return false
}
