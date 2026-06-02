package dataprofile

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Rows        int             `json:"rows"`
	Columns     []ColumnProfile `json:"columns"`
	Issues      []Issue         `json:"issues"`
	SchemaRules int             `json:"schemaRules,omitempty"`
}

type Schema struct {
	Columns []ColumnRule `json:"columns,omitempty"`
}

type ColumnRule struct {
	Name     string   `json:"name,omitempty"`
	Column   int      `json:"column,omitempty"`
	Required bool     `json:"required,omitempty"`
	Type     string   `json:"type,omitempty"`
	Unique   bool     `json:"unique,omitempty"`
	Enum     []string `json:"enum,omitempty"`
	Pattern  string   `json:"pattern,omitempty"`
	Min      *float64 `json:"min,omitempty"`
	Max      *float64 `json:"max,omitempty"`
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
}

func Analyze(records [][]string, hasHeader bool) Result {
	return AnalyzeWithSchema(records, hasHeader, Schema{})
}

func AnalyzeWithSchema(records [][]string, hasHeader bool, schema Schema) Result {
	result := Result{Rows: dataRowCount(records, hasHeader)}
	if len(records) == 0 {
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
	return result
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

func ruleHasChecks(rule ColumnRule) bool {
	return strings.TrimSpace(rule.Name) != "" ||
		rule.Column > 0 ||
		rule.Required ||
		strings.TrimSpace(rule.Type) != "" ||
		rule.Unique ||
		len(rule.Enum) > 0 ||
		strings.TrimSpace(rule.Pattern) != "" ||
		rule.Min != nil ||
		rule.Max != nil
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
					Type:     "schema-required",
					Severity: "error",
					Row:      row + 1,
					Col:      col + 1,
					Rule:     ruleName,
					Message:  fmt.Sprintf("%q is required.", ruleName),
				})
			}
			continue
		}
		if expectedType != "" && !schemaTypeMatches(raw, expectedType) {
			issues = append(issues, Issue{
				Type:     "schema-type",
				Severity: "error",
				Row:      row + 1,
				Col:      col + 1,
				Rule:     ruleName,
				Message:  fmt.Sprintf("Value %q does not match schema type %s.", raw, expectedType),
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
					Type:     "schema-enum",
					Severity: "error",
					Row:      row + 1,
					Col:      col + 1,
					Rule:     ruleName,
					Message:  fmt.Sprintf("Value %q is not allowed by schema enum.", raw),
				})
			}
		}
		if pattern != nil && !pattern.MatchString(raw) {
			issues = append(issues, Issue{
				Type:     "schema-pattern",
				Severity: "error",
				Row:      row + 1,
				Col:      col + 1,
				Rule:     ruleName,
				Message:  fmt.Sprintf("Value %q does not match schema pattern.", raw),
			})
		}
		issues = append(issues, numericRangeIssues(raw, row, col, rule, ruleName)...)
	}
	return issues
}

func schemaTypeMatches(raw string, expectedType string) bool {
	kind := detectType(raw)
	switch expectedType {
	case "string", "text":
		return true
	case "number", "numeric":
		return kind == "integer" || kind == "decimal"
	case "integer", "decimal", "boolean", "date":
		return kind == expectedType || (expectedType == "decimal" && kind == "integer")
	default:
		return true
	}
}

func numericRangeIssues(raw string, row int, col int, rule ColumnRule, ruleName string) []Issue {
	if rule.Min == nil && rule.Max == nil {
		return nil
	}
	number, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return []Issue{{
			Type:     "schema-range",
			Severity: "error",
			Row:      row + 1,
			Col:      col + 1,
			Rule:     ruleName,
			Message:  fmt.Sprintf("Value %q is not numeric for schema range check.", raw),
		}}
	}
	if rule.Min != nil && number < *rule.Min {
		return []Issue{{
			Type:     "schema-range",
			Severity: "error",
			Row:      row + 1,
			Col:      col + 1,
			Rule:     ruleName,
			Message:  fmt.Sprintf("Value %q is below schema minimum %s.", raw, formatRuleNumber(*rule.Min)),
		}}
	}
	if rule.Max != nil && number > *rule.Max {
		return []Issue{{
			Type:     "schema-range",
			Severity: "error",
			Row:      row + 1,
			Col:      col + 1,
			Rule:     ruleName,
			Message:  fmt.Sprintf("Value %q is above schema maximum %s.", raw, formatRuleNumber(*rule.Max)),
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
	items := make([]ValueFrequency, 0, len(values))
	for value, count := range values {
		items = append(items, ValueFrequency{Value: value, Count: count})
	}
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].Count > items[i].Count || (items[j].Count == items[i].Count && items[j].Value < items[i].Value) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	if len(items) > limit {
		items = items[:limit]
	}
	return items
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
