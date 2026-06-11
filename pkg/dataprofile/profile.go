package dataprofile

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Result struct {
	Rows        int              `json:"rows"`
	Columns     []ColumnProfile  `json:"columns"`
	Issues      []Issue          `json:"issues"`
	SchemaRules int              `json:"schemaRules,omitempty"`
	Streamed    bool             `json:"streamed,omitempty"`
	Progress    *ProgressSummary `json:"progress,omitempty"`
}

type Schema struct {
	Columns []ColumnRule `json:"columns,omitempty"`
}

type ProgressSummary struct {
	ProfileRows    int `json:"profileRows"`
	ValidationRows int `json:"validationRows"`
	TotalRows      int `json:"totalRows"`
	Passes         int `json:"passes"`
	Events         int `json:"events,omitempty"`
}

type ProgressEvent struct {
	Phase   string `json:"phase"`
	Rows    int    `json:"rows"`
	Records int    `json:"records"`
	Done    bool   `json:"done,omitempty"`
}

type StreamOptions struct {
	ProgressInterval int
	ReportProgress   func(ProgressEvent)
}

type ColumnRule struct {
	Name      string   `json:"name,omitempty"`
	Column    int      `json:"column,omitempty"`
	Required  bool     `json:"required,omitempty"`
	Type      string   `json:"type,omitempty"`
	Unique    bool     `json:"unique,omitempty"`
	Enum      []string `json:"enum,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	MinLength *int     `json:"minLength,omitempty"`
	MaxLength *int     `json:"maxLength,omitempty"`
	Default   string   `json:"default,omitempty"`
}

type ColumnProfile struct {
	Index      int              `json:"index"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Rows       int              `json:"rows"`
	Nulls      int              `json:"nulls"`
	NullRatio  float64          `json:"nullRatio"`
	Distinct   int              `json:"distinct"`
	Min        string           `json:"min,omitempty"`
	Max        string           `json:"max,omitempty"`
	Sum        *float64         `json:"sum,omitempty"`
	Average    *float64         `json:"average,omitempty"`
	TypeCounts map[string]int   `json:"typeCounts"`
	TopValues  []ValueFrequency `json:"topValues,omitempty"`
}

type ValueFrequency struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type Issue struct {
	Type       string     `json:"type"`
	Severity   string     `json:"severity"`
	Row        int        `json:"row,omitempty"`
	Col        int        `json:"col,omitempty"`
	Message    string     `json:"message"`
	Rule       string     `json:"rule,omitempty"`
	Actual     int        `json:"actual,omitempty"`
	Expected   int        `json:"expected,omitempty"`
	QuickFixes []QuickFix `json:"quickFixes,omitempty"`
}

type QuickFix struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value,omitempty"`
}

type RecordStream func(ctx context.Context, handle func(row int, record []string) error) (int, error)

func Analyze(records [][]string, hasHeader bool) Result {
	return AnalyzeWithSchema(records, hasHeader, Schema{})
}

func AnalyzeStreamContext(ctx context.Context, hasHeader bool, schema Schema, stream RecordStream) (Result, error) {
	return AnalyzeStreamContextWithOptions(ctx, hasHeader, schema, stream, StreamOptions{})
}

func AnalyzeStreamContextWithOptions(ctx context.Context, hasHeader bool, schema Schema, stream RecordStream, options StreamOptions) (Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if stream == nil {
		return Result{}, fmt.Errorf("record stream is required")
	}
	acc := newProfileAccumulator(hasHeader)
	profileProgress := newStreamProgressTracker("profile", hasHeader, options)
	if _, err := stream(ctx, func(row int, record []string) error {
		profileProgress.observe(row)
		return acc.consume(row, record)
	}); err != nil {
		return Result{}, err
	}
	profileProgress.finish()
	result := acc.finalize()
	validationProgress := newStreamProgressTracker("validation", hasHeader, options)
	if err := streamProfileValidationPass(ctx, stream, hasHeader, acc.headers(), &result, schema, validationProgress); err != nil {
		return Result{}, err
	}
	validationProgress.finish()
	result.Streamed = true
	result.Progress = &ProgressSummary{
		ProfileRows:    profileProgress.rows,
		ValidationRows: validationProgress.rows,
		TotalRows:      result.Rows,
		Passes:         2,
		Events:         profileProgress.events + validationProgress.events,
	}
	return result, nil
}

func AnalyzeWithSchema(records [][]string, hasHeader bool, schema Schema) Result {
	result := Result{Rows: dataRowCount(records, hasHeader)}
	if len(records) == 0 {
		result.Progress = materializedProgressSummary(result.Rows)
		return result
	}

	width := maxWidth(records)
	headers := headersFor(records, hasHeader, width)
	result.Columns = make([]ColumnProfile, width)
	for col := 0; col < width; col++ {
		result.Columns[col] = ColumnProfile{
			Index:      col,
			Name:       headers[col],
			Rows:       dataRowCount(records, hasHeader),
			TypeCounts: make(map[string]int),
		}
	}

	if hasHeader {
		result.Issues = append(result.Issues, headerIssues(records[0])...)
	}
	result.Issues = append(result.Issues, rowWidthIssues(records, hasHeader, width)...)

	start := 0
	if hasHeader {
		start = 1
	}
	values := make([]map[string]int, width)
	for col := range values {
		values[col] = make(map[string]int)
	}
	for row := start; row < len(records); row++ {
		for col := 0; col < width; col++ {
			profile := &result.Columns[col]
			raw := strings.TrimSpace(cell(records[row], col))
			if raw == "" {
				profile.Nulls++
				continue
			}
			values[col][raw]++
			kind := detectType(raw)
			profile.TypeCounts[kind]++
			updateMinMax(profile, raw, kind)
			if kind == "integer" || kind == "decimal" {
				num, _ := strconv.ParseFloat(raw, 64)
				if profile.Sum == nil {
					zero := 0.0
					profile.Sum = &zero
				}
				*profile.Sum += num
			}
		}
	}

	for col := range result.Columns {
		profile := &result.Columns[col]
		profile.Distinct = len(values[col])
		if profile.Rows > 0 {
			profile.NullRatio = float64(profile.Nulls) / float64(profile.Rows)
		}
		profile.Type = dominantType(profile.TypeCounts)
		if profile.Sum != nil {
			numericRows := profile.TypeCounts["integer"] + profile.TypeCounts["decimal"]
			if numericRows > 0 {
				average := *profile.Sum / float64(numericRows)
				profile.Average = &average
			}
		}
		profile.TopValues = topValues(values[col], 5)
		result.Issues = append(result.Issues, typeIssues(records, hasHeader, col, profile.Type)...)
	}
	schemaRules, schemaIssues := schemaValidationIssues(records, hasHeader, headers, schema)
	result.SchemaRules = schemaRules
	result.Issues = append(result.Issues, schemaIssues...)
	result.Progress = materializedProgressSummary(result.Rows)
	return result
}

func materializedProgressSummary(rows int) *ProgressSummary {
	return &ProgressSummary{
		ProfileRows:    rows,
		ValidationRows: rows,
		TotalRows:      rows,
		Passes:         2,
	}
}

type streamProgressTracker struct {
	phase     string
	hasHeader bool
	interval  int
	report    func(ProgressEvent)
	rows      int
	records   int
	events    int
}

func newStreamProgressTracker(phase string, hasHeader bool, options StreamOptions) *streamProgressTracker {
	interval := options.ProgressInterval
	if interval <= 0 {
		interval = 1000
	}
	return &streamProgressTracker{
		phase:     phase,
		hasHeader: hasHeader,
		interval:  interval,
		report:    options.ReportProgress,
	}
}

func (tracker *streamProgressTracker) observe(row int) {
	tracker.records++
	if tracker.hasHeader && row == 0 {
		return
	}
	tracker.rows++
	if tracker.report != nil && tracker.interval > 0 && tracker.rows > 0 && tracker.rows%tracker.interval == 0 {
		tracker.emit(false)
	}
}

func (tracker *streamProgressTracker) finish() {
	if tracker.report != nil {
		tracker.emit(true)
	}
}

func (tracker *streamProgressTracker) emit(done bool) {
	tracker.events++
	tracker.report(ProgressEvent{
		Phase:   tracker.phase,
		Rows:    tracker.rows,
		Records: tracker.records,
		Done:    done,
	})
}

type rowWidth struct {
	row   int
	width int
}

type profileAccumulator struct {
	hasHeader bool
	header    []string
	rows      int
	widths    []rowWidth
	result    Result
	values    []map[string]int
}

func newProfileAccumulator(hasHeader bool) *profileAccumulator {
	return &profileAccumulator{hasHeader: hasHeader}
}

func (acc *profileAccumulator) consume(row int, record []string) error {
	if row == 0 && acc.hasHeader {
		acc.header = append([]string(nil), record...)
		acc.ensureWidth(len(record))
		acc.result.Issues = append(acc.result.Issues, headerIssues(record)...)
		return nil
	}
	acc.ensureWidth(len(record))
	acc.rows++
	acc.widths = append(acc.widths, rowWidth{row: row + 1, width: len(record)})
	for col := 0; col < len(acc.result.Columns); col++ {
		profile := &acc.result.Columns[col]
		raw := strings.TrimSpace(cell(record, col))
		if raw == "" {
			profile.Nulls++
			continue
		}
		acc.values[col][raw]++
		kind := detectType(raw)
		profile.TypeCounts[kind]++
		updateMinMax(profile, raw, kind)
		if kind == "integer" || kind == "decimal" {
			num, _ := strconv.ParseFloat(raw, 64)
			if profile.Sum == nil {
				zero := 0.0
				profile.Sum = &zero
			}
			*profile.Sum += num
		}
	}
	return nil
}

func (acc *profileAccumulator) ensureWidth(width int) {
	for len(acc.result.Columns) < width {
		col := len(acc.result.Columns)
		name := fmt.Sprintf("c%d", col+1)
		if acc.hasHeader && col < len(acc.header) && strings.TrimSpace(acc.header[col]) != "" {
			name = acc.header[col]
		}
		acc.result.Columns = append(acc.result.Columns, ColumnProfile{
			Index:      col,
			Name:       name,
			Nulls:      acc.rows,
			TypeCounts: make(map[string]int),
		})
		acc.values = append(acc.values, make(map[string]int))
	}
}

func (acc *profileAccumulator) headers() []string {
	headers := make([]string, len(acc.result.Columns))
	for col := range headers {
		headers[col] = acc.result.Columns[col].Name
	}
	return headers
}

func (acc *profileAccumulator) finalize() Result {
	acc.result.Rows = acc.rows
	for _, width := range acc.widths {
		if width.width != len(acc.result.Columns) {
			issue := Issue{
				Type:     "ragged-row",
				Severity: "error",
				Row:      width.row,
				Actual:   width.width,
				Expected: len(acc.result.Columns),
				Message:  fmt.Sprintf("Row has %d columns, expected %d.", width.width, len(acc.result.Columns)),
			}
			if width.width < len(acc.result.Columns) {
				issue.QuickFixes = []QuickFix{{ID: "row.pad", Label: "Pad missing cells"}}
			}
			acc.result.Issues = append(acc.result.Issues, issue)
		}
	}
	for col := range acc.result.Columns {
		profile := &acc.result.Columns[col]
		profile.Rows = acc.rows
		profile.Distinct = len(acc.values[col])
		if profile.Rows > 0 {
			profile.NullRatio = float64(profile.Nulls) / float64(profile.Rows)
		}
		profile.Type = dominantType(profile.TypeCounts)
		if profile.Sum != nil {
			numericRows := profile.TypeCounts["integer"] + profile.TypeCounts["decimal"]
			if numericRows > 0 {
				average := *profile.Sum / float64(numericRows)
				profile.Average = &average
			}
		}
		profile.TopValues = topValues(acc.values[col], 5)
	}
	return acc.result
}

func headersFor(records [][]string, hasHeader bool, width int) []string {
	headers := make([]string, width)
	for col := range headers {
		headers[col] = fmt.Sprintf("c%d", col+1)
	}
	if !hasHeader || len(records) == 0 {
		return headers
	}
	for col := 0; col < min(width, len(records[0])); col++ {
		if strings.TrimSpace(records[0][col]) != "" {
			headers[col] = records[0][col]
		}
	}
	return headers
}

func headerIssues(header []string) []Issue {
	seen := make(map[string]int)
	var issues []Issue
	for col, name := range header {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			issues = append(issues, Issue{
				Type:     "empty-header",
				Severity: "warning",
				Row:      1,
				Col:      col + 1,
				Message:  fmt.Sprintf("Header column %d is empty.", col+1),
				QuickFixes: []QuickFix{
					{ID: "header.fill", Label: "Fill header"},
				},
			})
			continue
		}
		key := strings.ToLower(trimmed)
		if first, ok := seen[key]; ok {
			issues = append(issues, Issue{
				Type:     "duplicate-header",
				Severity: "warning",
				Row:      1,
				Col:      col + 1,
				Message:  fmt.Sprintf("Header %q duplicates column %d.", trimmed, first+1),
				QuickFixes: []QuickFix{
					{ID: "header.unique", Label: "Make header unique"},
				},
			})
		}
		seen[key] = col
	}
	return issues
}

func rowWidthIssues(records [][]string, hasHeader bool, width int) []Issue {
	var issues []Issue
	start := 0
	if hasHeader {
		start = 1
	}
	for row := start; row < len(records); row++ {
		if len(records[row]) != width {
			issue := Issue{
				Type:     "ragged-row",
				Severity: "error",
				Row:      row + 1,
				Actual:   len(records[row]),
				Expected: width,
				Message:  fmt.Sprintf("Row has %d columns, expected %d.", len(records[row]), width),
			}
			if len(records[row]) < width {
				issue.QuickFixes = []QuickFix{{ID: "row.pad", Label: "Pad missing cells"}}
			}
			issues = append(issues, issue)
		}
	}
	return issues
}

func schemaValidationIssues(records [][]string, hasHeader bool, headers []string, schema Schema) (int, []Issue) {
	rules := 0
	var issues []Issue
	for _, rule := range schema.Columns {
		if !ruleHasChecks(rule) {
			continue
		}
		rules++
		col, ok := resolveRuleColumn(rule, headers)
		ruleName := columnRuleName(rule)
		if !ok {
			issues = append(issues, Issue{
				Type:     "schema-missing-column",
				Severity: "error",
				Row:      1,
				Rule:     ruleName,
				Message:  fmt.Sprintf("Schema column %q was not found.", ruleName),
			})
			continue
		}
		issues = append(issues, validateColumnRule(records, hasHeader, col, rule, ruleName)...)
	}
	return rules, issues
}

type compiledColumnRule struct {
	col          int
	rule         ColumnRule
	ruleName     string
	expectedType string
	enum         map[string]struct{}
	pattern      *regexp.Regexp
	seen         map[string]int
}

func streamProfileValidationPass(ctx context.Context, stream RecordStream, hasHeader bool, headers []string, result *Result, schema Schema, progress *streamProgressTracker) error {
	if ctx == nil {
		ctx = context.Background()
	}
	compiledRules, ruleCount, issues := compileColumnRules(headers, schema)
	result.SchemaRules = ruleCount
	result.Issues = append(result.Issues, issues...)
	if len(result.Columns) == 0 && len(compiledRules) == 0 {
		return nil
	}
	_, err := stream(ctx, func(row int, record []string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if progress != nil {
			progress.observe(row)
		}
		if row == 0 && hasHeader {
			return nil
		}
		result.Issues = append(result.Issues, streamedTypeIssues(record, row, result.Columns)...)
		for index := range compiledRules {
			result.Issues = append(result.Issues, streamedColumnRuleIssues(record, row, &compiledRules[index])...)
		}
		return nil
	})
	return err
}

func compileColumnRules(headers []string, schema Schema) ([]compiledColumnRule, int, []Issue) {
	var compiled []compiledColumnRule
	var issues []Issue
	rules := 0
	for _, rule := range schema.Columns {
		if !ruleHasChecks(rule) {
			continue
		}
		rules++
		col, ok := resolveRuleColumn(rule, headers)
		ruleName := columnRuleName(rule)
		if !ok {
			issues = append(issues, Issue{
				Type:     "schema-missing-column",
				Severity: "error",
				Row:      1,
				Rule:     ruleName,
				Message:  fmt.Sprintf("Schema column %q was not found.", ruleName),
			})
			continue
		}
		compiledRule := compiledColumnRule{
			col:          col,
			rule:         rule,
			ruleName:     ruleName,
			expectedType: strings.ToLower(strings.TrimSpace(rule.Type)),
			enum:         make(map[string]struct{}, len(rule.Enum)),
			seen:         make(map[string]int),
		}
		for _, value := range rule.Enum {
			compiledRule.enum[value] = struct{}{}
		}
		if strings.TrimSpace(rule.Pattern) != "" {
			pattern, err := regexp.Compile(rule.Pattern)
			if err != nil {
				issues = append(issues, Issue{
					Type:     "schema-pattern",
					Severity: "error",
					Row:      1,
					Col:      col + 1,
					Rule:     ruleName,
					Message:  fmt.Sprintf("Schema pattern for %q is invalid: %v.", ruleName, err),
				})
				continue
			}
			compiledRule.pattern = pattern
		}
		compiled = append(compiled, compiledRule)
	}
	return compiled, rules, issues
}

func streamedTypeIssues(record []string, row int, columns []ColumnProfile) []Issue {
	var issues []Issue
	for col, profile := range columns {
		dominant := profile.Type
		if dominant == "" || dominant == "string" || dominant == "mixed" {
			continue
		}
		raw := strings.TrimSpace(cell(record, col))
		if raw == "" {
			continue
		}
		kind := detectType(raw)
		if kind != dominant && !(dominant == "decimal" && kind == "integer") {
			issues = append(issues, Issue{
				Type:     "type-anomaly",
				Severity: "warning",
				Row:      row + 1,
				Col:      col + 1,
				Message:  fmt.Sprintf("Value %q is %s in a mostly %s column.", raw, kind, dominant),
			})
		}
	}
	return issues
}

func streamedColumnRuleIssues(record []string, row int, compiled *compiledColumnRule) []Issue {
	raw := strings.TrimSpace(cell(record, compiled.col))
	if raw == "" {
		if compiled.rule.Required {
			return []Issue{{
				Type:       "schema-required",
				Severity:   "error",
				Row:        row + 1,
				Col:        compiled.col + 1,
				Rule:       compiled.ruleName,
				Message:    fmt.Sprintf("%q is required.", compiled.ruleName),
				QuickFixes: defaultQuickFixesForRule(compiled.rule),
			}}
		}
		return nil
	}
	var issues []Issue
	if compiled.expectedType != "" && !schemaTypeMatches(raw, compiled.expectedType) {
		issues = append(issues, Issue{
			Type:       "schema-type",
			Severity:   "error",
			Row:        row + 1,
			Col:        compiled.col + 1,
			Rule:       compiled.ruleName,
			Message:    fmt.Sprintf("Value %q does not match schema type %s.", raw, compiled.expectedType),
			QuickFixes: defaultQuickFixesForRule(compiled.rule),
		})
	}
	if compiled.rule.Unique {
		if firstRow, ok := compiled.seen[raw]; ok {
			issues = append(issues, Issue{
				Type:     "schema-unique",
				Severity: "error",
				Row:      row + 1,
				Col:      compiled.col + 1,
				Rule:     compiled.ruleName,
				Message:  fmt.Sprintf("Value %q duplicates row %d.", raw, firstRow),
			})
		} else {
			compiled.seen[raw] = row + 1
		}
	}
	if len(compiled.enum) > 0 {
		if _, ok := compiled.enum[raw]; !ok {
			issues = append(issues, Issue{
				Type:       "schema-enum",
				Severity:   "error",
				Row:        row + 1,
				Col:        compiled.col + 1,
				Rule:       compiled.ruleName,
				Message:    fmt.Sprintf("Value %q is not allowed by schema enum.", raw),
				QuickFixes: mergeQuickFixes(enumQuickFixes(compiled.rule.Enum), defaultQuickFixesForRule(compiled.rule)),
			})
		}
	}
	if compiled.pattern != nil && !compiled.pattern.MatchString(raw) {
		issues = append(issues, Issue{
			Type:       "schema-pattern",
			Severity:   "error",
			Row:        row + 1,
			Col:        compiled.col + 1,
			Rule:       compiled.ruleName,
			Message:    fmt.Sprintf("Value %q does not match schema pattern.", raw),
			QuickFixes: defaultQuickFixesForRule(compiled.rule),
		})
	}
	issues = append(issues, stringLengthIssues(raw, row, compiled.col, compiled.rule, compiled.ruleName)...)
	issues = append(issues, numericRangeIssues(raw, row, compiled.col, compiled.rule, compiled.ruleName)...)
	return issues
}

func ruleHasChecks(rule ColumnRule) bool {
	return strings.TrimSpace(rule.Name) != "" ||
		rule.Column > 0 ||
		rule.Required ||
		strings.TrimSpace(rule.Type) != "" ||
		rule.Unique ||
		len(rule.Enum) > 0 ||
		strings.TrimSpace(rule.Pattern) != "" ||
		rule.Min != nil ||
		rule.Max != nil ||
		rule.MinLength != nil ||
		rule.MaxLength != nil
}

func resolveRuleColumn(rule ColumnRule, headers []string) (int, bool) {
	if rule.Column > 0 {
		col := rule.Column - 1
		return col, col >= 0 && col < len(headers)
	}
	name := strings.TrimSpace(rule.Name)
	if name == "" {
		return -1, false
	}
	for col, header := range headers {
		if header == name {
			return col, true
		}
	}
	for col, header := range headers {
		if strings.EqualFold(header, name) {
			return col, true
		}
	}
	return -1, false
}

func columnRuleName(rule ColumnRule) string {
	if strings.TrimSpace(rule.Name) != "" {
		return strings.TrimSpace(rule.Name)
	}
	if rule.Column > 0 {
		return fmt.Sprintf("column %d", rule.Column)
	}
	return "unnamed rule"
}

func validateColumnRule(records [][]string, hasHeader bool, col int, rule ColumnRule, ruleName string) []Issue {
	var issues []Issue
	start := 0
	if hasHeader {
		start = 1
	}

	expectedType := strings.ToLower(strings.TrimSpace(rule.Type))
	enum := make(map[string]struct{}, len(rule.Enum))
	for _, value := range rule.Enum {
		enum[value] = struct{}{}
	}
	var pattern *regexp.Regexp
	if strings.TrimSpace(rule.Pattern) != "" {
		compiled, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return []Issue{{
				Type:     "schema-pattern",
				Severity: "error",
				Row:      1,
				Col:      col + 1,
				Rule:     ruleName,
				Message:  fmt.Sprintf("Schema pattern for %q is invalid: %v.", ruleName, err),
			}}
		}
		pattern = compiled
	}
	seen := make(map[string]int)
	for row := start; row < len(records); row++ {
		raw := strings.TrimSpace(cell(records[row], col))
		if raw == "" {
			if rule.Required {
				issues = append(issues, Issue{
					Type:       "schema-required",
					Severity:   "error",
					Row:        row + 1,
					Col:        col + 1,
					Rule:       ruleName,
					Message:    fmt.Sprintf("%q is required.", ruleName),
					QuickFixes: defaultQuickFixesForRule(rule),
				})
			}
			continue
		}
		if expectedType != "" && !schemaTypeMatches(raw, expectedType) {
			issues = append(issues, Issue{
				Type:       "schema-type",
				Severity:   "error",
				Row:        row + 1,
				Col:        col + 1,
				Rule:       ruleName,
				Message:    fmt.Sprintf("Value %q does not match schema type %s.", raw, expectedType),
				QuickFixes: defaultQuickFixesForRule(rule),
			})
		}
		if rule.Unique {
			if firstRow, ok := seen[raw]; ok {
				issues = append(issues, Issue{
					Type:     "schema-unique",
					Severity: "error",
					Row:      row + 1,
					Col:      col + 1,
					Rule:     ruleName,
					Message:  fmt.Sprintf("Value %q duplicates row %d.", raw, firstRow),
				})
			} else {
				seen[raw] = row + 1
			}
		}
		if len(enum) > 0 {
			if _, ok := enum[raw]; !ok {
				issues = append(issues, Issue{
					Type:       "schema-enum",
					Severity:   "error",
					Row:        row + 1,
					Col:        col + 1,
					Rule:       ruleName,
					Message:    fmt.Sprintf("Value %q is not allowed by schema enum.", raw),
					QuickFixes: mergeQuickFixes(enumQuickFixes(rule.Enum), defaultQuickFixesForRule(rule)),
				})
			}
		}
		if pattern != nil && !pattern.MatchString(raw) {
			issues = append(issues, Issue{
				Type:       "schema-pattern",
				Severity:   "error",
				Row:        row + 1,
				Col:        col + 1,
				Rule:       ruleName,
				Message:    fmt.Sprintf("Value %q does not match schema pattern.", raw),
				QuickFixes: defaultQuickFixesForRule(rule),
			})
		}
		issues = append(issues, stringLengthIssues(raw, row, col, rule, ruleName)...)
		issues = append(issues, numericRangeIssues(raw, row, col, rule, ruleName)...)
	}
	return issues
}

func enumQuickFixes(values []string) []QuickFix {
	if len(values) != 1 || values[0] == "" {
		return nil
	}
	return []QuickFix{{ID: "schema.enum.useValue", Label: "Use allowed value", Value: values[0]}}
}

func defaultQuickFixes(value string) []QuickFix {
	if value == "" {
		return nil
	}
	return []QuickFix{{ID: "schema.default.use", Label: "Use default", Value: value}}
}

func defaultQuickFixesForRule(rule ColumnRule) []QuickFix {
	if !schemaDefaultMatchesRule(rule) {
		return nil
	}
	return defaultQuickFixes(rule.Default)
}

func schemaDefaultMatchesRule(rule ColumnRule) bool {
	value := strings.TrimSpace(rule.Default)
	if value == "" || rule.Unique {
		return false
	}
	expectedType := strings.ToLower(strings.TrimSpace(rule.Type))
	if expectedType != "" && !schemaTypeMatches(value, expectedType) {
		return false
	}
	if len(rule.Enum) > 0 {
		allowed := false
		for _, enumValue := range rule.Enum {
			if value == enumValue {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	if strings.TrimSpace(rule.Pattern) != "" {
		pattern, err := regexp.Compile(rule.Pattern)
		if err != nil || !pattern.MatchString(value) {
			return false
		}
	}
	if rule.Min != nil || rule.Max != nil {
		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false
		}
		if rule.Min != nil && number < *rule.Min {
			return false
		}
		if rule.Max != nil && number > *rule.Max {
			return false
		}
	}
	length := utf8.RuneCountInString(value)
	if rule.MinLength != nil && length < *rule.MinLength {
		return false
	}
	if rule.MaxLength != nil && length > *rule.MaxLength {
		return false
	}
	return true
}

func mergeQuickFixes(groups ...[]QuickFix) []QuickFix {
	var fixes []QuickFix
	seen := make(map[string]bool)
	for _, group := range groups {
		for _, fix := range group {
			key := fix.ID + "\x00" + fix.Value
			if seen[key] {
				continue
			}
			seen[key] = true
			fixes = append(fixes, fix)
		}
	}
	return fixes
}

func schemaTypeMatches(raw string, expectedType string) bool {
	if schemaTemporalTypeMatches(raw, expectedType) {
		return true
	}
	kind := detectType(raw)
	switch expectedType {
	case "string", "text":
		return true
	case "number", "numeric":
		return kind == "integer" || kind == "decimal"
	case "integer", "decimal", "boolean", "date", "time", "year", "yearmonth":
		return kind == expectedType || (expectedType == "decimal" && kind == "integer")
	default:
		return true
	}
}

func schemaTemporalTypeMatches(raw string, expectedType string) bool {
	layouts := map[string][]string{
		"date":      {"2006-01-02", "2006-01-02 15:04:05", time.RFC3339},
		"time":      {"15:04", "15:04:05", "15:04:05.999999999"},
		"year":      {"2006"},
		"yearmonth": {"2006-01"},
	}
	for _, layout := range layouts[expectedType] {
		if _, err := time.Parse(layout, raw); err == nil {
			return true
		}
	}
	return false
}

func numericRangeIssues(raw string, row int, col int, rule ColumnRule, ruleName string) []Issue {
	if rule.Min == nil && rule.Max == nil {
		return nil
	}
	number, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return []Issue{{
			Type:       "schema-range",
			Severity:   "error",
			Row:        row + 1,
			Col:        col + 1,
			Rule:       ruleName,
			Message:    fmt.Sprintf("Value %q is not numeric for schema range check.", raw),
			QuickFixes: defaultQuickFixesForRule(rule),
		}}
	}
	if rule.Min != nil && number < *rule.Min {
		target := formatRuleNumber(*rule.Min)
		return []Issue{{
			Type:     "schema-range",
			Severity: "error",
			Row:      row + 1,
			Col:      col + 1,
			Rule:     ruleName,
			Message:  fmt.Sprintf("Value %q is below schema minimum %s.", raw, target),
			QuickFixes: mergeQuickFixes([]QuickFix{
				{ID: "schema.range.clampMin", Label: "Use minimum", Value: target},
			}, defaultQuickFixesForRule(rule)),
		}}
	}
	if rule.Max != nil && number > *rule.Max {
		target := formatRuleNumber(*rule.Max)
		return []Issue{{
			Type:     "schema-range",
			Severity: "error",
			Row:      row + 1,
			Col:      col + 1,
			Rule:     ruleName,
			Message:  fmt.Sprintf("Value %q is above schema maximum %s.", raw, target),
			QuickFixes: mergeQuickFixes([]QuickFix{
				{ID: "schema.range.clampMax", Label: "Use maximum", Value: target},
			}, defaultQuickFixesForRule(rule)),
		}}
	}
	return nil
}

func stringLengthIssues(raw string, row int, col int, rule ColumnRule, ruleName string) []Issue {
	if rule.MinLength == nil && rule.MaxLength == nil {
		return nil
	}
	length := utf8.RuneCountInString(raw)
	if rule.MinLength != nil && length < *rule.MinLength {
		return []Issue{{
			Type:       "schema-min-length",
			Severity:   "error",
			Row:        row + 1,
			Col:        col + 1,
			Rule:       ruleName,
			Message:    fmt.Sprintf("Value %q is shorter than schema minimum length %d.", raw, *rule.MinLength),
			QuickFixes: defaultQuickFixesForRule(rule),
		}}
	}
	if rule.MaxLength != nil && length > *rule.MaxLength {
		return []Issue{{
			Type:       "schema-max-length",
			Severity:   "error",
			Row:        row + 1,
			Col:        col + 1,
			Rule:       ruleName,
			Message:    fmt.Sprintf("Value %q is longer than schema maximum length %d.", raw, *rule.MaxLength),
			QuickFixes: defaultQuickFixesForRule(rule),
		}}
	}
	return nil
}

func formatRuleNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func typeIssues(records [][]string, hasHeader bool, col int, dominant string) []Issue {
	if dominant == "" || dominant == "string" || dominant == "mixed" {
		return nil
	}
	start := 0
	if hasHeader {
		start = 1
	}
	var issues []Issue
	for row := start; row < len(records); row++ {
		raw := strings.TrimSpace(cell(records[row], col))
		if raw == "" {
			continue
		}
		kind := detectType(raw)
		if kind != dominant && !(dominant == "decimal" && kind == "integer") {
			issues = append(issues, Issue{
				Type:     "type-anomaly",
				Severity: "warning",
				Row:      row + 1,
				Col:      col + 1,
				Message:  fmt.Sprintf("Value %q is %s in a mostly %s column.", raw, kind, dominant),
			})
		}
	}
	return issues
}

func dominantType(counts map[string]int) string {
	total := 0
	bestType := ""
	bestCount := 0
	for kind, count := range counts {
		total += count
		if count > bestCount {
			bestType = kind
			bestCount = count
		}
	}
	if total == 0 {
		return "empty"
	}
	if counts["decimal"] > 0 && counts["integer"] > 0 && counts["decimal"]+counts["integer"] == total {
		return "decimal"
	}
	if float64(bestCount)/float64(total) < 0.8 {
		return "mixed"
	}
	return bestType
}

func detectType(value string) string {
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "integer"
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "decimal"
	}
	switch strings.ToLower(value) {
	case "true", "false":
		return "boolean"
	}
	for _, layout := range []string{"2006-01-02", "2006-01-02 15:04:05", time.RFC3339} {
		if _, err := time.Parse(layout, value); err == nil {
			return "date"
		}
	}
	for _, layout := range []string{"15:04", "15:04:05", "15:04:05.999999999"} {
		if _, err := time.Parse(layout, value); err == nil {
			return "time"
		}
	}
	if _, err := time.Parse("2006", value); err == nil {
		return "year"
	}
	if _, err := time.Parse("2006-01", value); err == nil {
		return "yearmonth"
	}
	return "string"
}

func updateMinMax(profile *ColumnProfile, raw string, kind string) {
	if profile.Min == "" || compareProfileValue(raw, profile.Min, kind) < 0 {
		profile.Min = raw
	}
	if profile.Max == "" || compareProfileValue(raw, profile.Max, kind) > 0 {
		profile.Max = raw
	}
}

func compareProfileValue(left string, right string, kind string) int {
	if kind == "integer" || kind == "decimal" {
		leftNum, leftErr := strconv.ParseFloat(left, 64)
		rightNum, rightErr := strconv.ParseFloat(right, 64)
		if leftErr == nil && rightErr == nil {
			if leftNum < rightNum {
				return -1
			}
			if leftNum > rightNum {
				return 1
			}
			return 0
		}
	}
	return strings.Compare(left, right)
}

func topValues(values map[string]int, limit int) []ValueFrequency {
	if limit <= 0 {
		return nil
	}
	items := make([]ValueFrequency, 0, min(limit, len(values)))
	for value, count := range values {
		candidate := ValueFrequency{Value: value, Count: count}
		if len(items) < limit {
			items = append(items, candidate)
			continue
		}
		worst := 0
		for i := 1; i < len(items); i++ {
			if topValueLess(items[i], items[worst]) {
				worst = i
			}
		}
		if topValueLess(items[worst], candidate) {
			items[worst] = candidate
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return topValueLess(items[j], items[i])
	})
	return items
}

func topValueLess(left ValueFrequency, right ValueFrequency) bool {
	if left.Count != right.Count {
		return left.Count < right.Count
	}
	return left.Value > right.Value
}

func maxWidth(records [][]string) int {
	width := 0
	for _, record := range records {
		width = max(width, len(record))
	}
	return width
}

func dataRowCount(records [][]string, hasHeader bool) int {
	if hasHeader && len(records) > 0 {
		return len(records) - 1
	}
	return len(records)
}

func cell(record []string, col int) string {
	if col < 0 || col >= len(record) {
		return ""
	}
	return record[col]
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
