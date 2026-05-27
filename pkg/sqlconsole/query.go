package sqlconsole

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/opencsv/opencsv/pkg/sqlexpr"
)

type Source struct {
	Name    string
	Path    string
	Headers []string
	Records [][]string
}

type Resolver func(ref string) (Source, error)

type Result struct {
	Columns     []string   `json:"columns"`
	Rows        [][]string `json:"rows"`
	TotalRows   int        `json:"totalRows"`
	Source      string     `json:"source"`
	Truncated   bool       `json:"truncated"`
	ExecutedSQL string     `json:"executedSql"`
}

type Options struct {
	DefaultLimit int
	MaxLimit     int
}

type queryParts struct {
	selectClause string
	fromRef      string
	whereClause  string
	groupBy      []string
	orderBy      string
	descending   bool
	limit        int
}

type selectExpr struct {
	Raw      string
	Alias    string
	Kind     string
	Function string
	Column   int
	Star     bool
	Output   string
}

func Execute(query string, resolver Resolver, options Options) (Result, error) {
	if strings.TrimSpace(query) == "" {
		return Result{}, errors.New("query is required")
	}
	if options.DefaultLimit <= 0 {
		options.DefaultLimit = 100
	}
	if options.MaxLimit <= 0 {
		options.MaxLimit = 5000
	}
	parts, err := parseQuery(query)
	if err != nil {
		return Result{}, err
	}
	if parts.limit <= 0 {
		parts.limit = options.DefaultLimit
	}
	if parts.limit > options.MaxLimit {
		parts.limit = options.MaxLimit
	}

	source, err := resolver(normalizeSourceRef(parts.fromRef))
	if err != nil {
		return Result{}, err
	}
	if len(source.Headers) == 0 && len(source.Records) > 0 {
		source.Headers = defaultHeaders(len(source.Records[0]))
	}

	dataStart := 0
	headers := source.Headers
	if len(source.Records) > 0 && sameRecord(source.Records[0], headers) {
		dataStart = 1
	}
	sourceRows := make([][]string, 0)
	if parts.whereClause != "" {
		program, err := sqlexpr.Compile(parts.whereClause, headers)
		if err != nil {
			return Result{}, err
		}
		for row := dataStart; row < len(source.Records); row++ {
			matched, err := program.Eval(source.Records[row])
			if err != nil {
				return Result{}, fmt.Errorf("row %d: %w", row+1, err)
			}
			if matched {
				sourceRows = append(sourceRows, source.Records[row])
			}
		}
	} else {
		sourceRows = append(sourceRows, source.Records[dataStart:]...)
	}

	if hasAggregation(parts.selectClause) || len(parts.groupBy) > 0 {
		expressions, err := parseSelectExpressions(parts.selectClause, headers)
		if err != nil {
			return Result{}, err
		}
		return executeAggregate(query, source, headers, sourceRows, expressions, parts, options)
	}

	selected, columns, err := parseSelectColumns(parts.selectClause, headers)
	if err != nil {
		return Result{}, err
	}

	if parts.orderBy != "" {
		orderIndex, err := resolveColumn(parts.orderBy, headers)
		if err != nil {
			return Result{}, err
		}
		sort.SliceStable(sourceRows, func(i, j int) bool {
			left := cell(sourceRows[i], orderIndex)
			right := cell(sourceRows[j], orderIndex)
			less := compareValues(left, right) < 0
			if parts.descending {
				return !less && left != right
			}
			return less
		})
	}

	total := len(sourceRows)
	if total > parts.limit {
		sourceRows = sourceRows[:parts.limit]
	}
	rows := make([][]string, 0, len(sourceRows))
	for _, record := range sourceRows {
		row := make([]string, len(selected))
		for i, col := range selected {
			row[i] = cell(record, col)
		}
		rows = append(rows, row)
	}
	return Result{
		Columns:     columns,
		Rows:        rows,
		TotalRows:   total,
		Source:      source.Path,
		Truncated:   total > len(rows),
		ExecutedSQL: strings.TrimSpace(query),
	}, nil
}

func parseQuery(query string) (queryParts, error) {
	text := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(query), ";"))
	lower := strings.ToLower(text)
	if !strings.HasPrefix(lower, "select ") {
		return queryParts{}, errors.New("only SELECT queries are supported")
	}
	fromPos := findKeyword(text, " from ", len("select "))
	if fromPos < 0 {
		return queryParts{}, errors.New("SELECT query requires FROM")
	}
	parts := queryParts{
		selectClause: strings.TrimSpace(text[len("select "):fromPos]),
	}
	afterFromStart := fromPos + len(" from ")
	wherePos := findKeyword(text, " where ", afterFromStart)
	groupPos := findKeyword(text, " group by ", afterFromStart)
	orderPos := findKeyword(text, " order by ", afterFromStart)
	limitPos := findKeyword(text, " limit ", afterFromStart)
	endFrom := firstPositive(len(text), wherePos, groupPos, orderPos, limitPos)
	parts.fromRef = strings.TrimSpace(text[afterFromStart:endFrom])
	if parts.fromRef == "" {
		return queryParts{}, errors.New("FROM source is required")
	}
	if wherePos >= 0 {
		whereStart := wherePos + len(" where ")
		whereEnd := firstPositive(len(text), groupPos, orderPos, limitPos)
		parts.whereClause = strings.TrimSpace(text[whereStart:whereEnd])
	}
	if groupPos >= 0 {
		groupStart := groupPos + len(" group by ")
		groupEnd := firstPositive(len(text), orderPos, limitPos)
		for _, token := range splitComma(strings.TrimSpace(text[groupStart:groupEnd])) {
			column := strings.Trim(strings.TrimSpace(token), `"`)
			if column != "" {
				parts.groupBy = append(parts.groupBy, column)
			}
		}
		if len(parts.groupBy) == 0 {
			return queryParts{}, errors.New("GROUP BY requires at least one column")
		}
	}
	if orderPos >= 0 {
		orderStart := orderPos + len(" order by ")
		orderEnd := firstPositive(len(text), limitPos)
		orderClause := strings.Fields(strings.TrimSpace(text[orderStart:orderEnd]))
		if len(orderClause) == 0 {
			return queryParts{}, errors.New("ORDER BY requires a column")
		}
		parts.orderBy = strings.Trim(orderClause[0], `"`)
		if len(orderClause) > 1 {
			parts.descending = strings.EqualFold(orderClause[1], "desc")
		}
	}
	if limitPos >= 0 {
		limitRaw := strings.TrimSpace(text[limitPos+len(" limit "):])
		limitFields := strings.Fields(limitRaw)
		if len(limitFields) == 0 {
			return queryParts{}, errors.New("LIMIT requires a number")
		}
		limit, err := strconv.Atoi(limitFields[0])
		if err != nil || limit < 0 {
			return queryParts{}, errors.New("LIMIT must be a non-negative integer")
		}
		parts.limit = limit
	}
	return parts, nil
}

func parseSelectColumns(selectClause string, headers []string) ([]int, []string, error) {
	if strings.TrimSpace(selectClause) == "*" {
		indexes := make([]int, len(headers))
		for i := range indexes {
			indexes[i] = i
		}
		return indexes, append([]string(nil), headers...), nil
	}
	tokens := splitComma(selectClause)
	indexes := make([]int, 0, len(tokens))
	names := make([]string, 0, len(tokens))
	for _, token := range tokens {
		name := strings.Trim(strings.TrimSpace(token), `"`)
		index, err := resolveColumn(name, headers)
		if err != nil {
			return nil, nil, err
		}
		indexes = append(indexes, index)
		names = append(names, headerName(headers, index))
	}
	return indexes, names, nil
}

func parseSelectExpressions(selectClause string, headers []string) ([]selectExpr, error) {
	tokens := splitComma(selectClause)
	expressions := make([]selectExpr, 0, len(tokens))
	for _, token := range tokens {
		raw, alias := splitAlias(strings.TrimSpace(token))
		expr := selectExpr{Raw: raw, Alias: alias, Kind: "column"}
		lower := strings.ToLower(raw)
		if open := strings.Index(lower, "("); open > 0 && strings.HasSuffix(strings.TrimSpace(raw), ")") {
			function := strings.ToLower(strings.TrimSpace(raw[:open]))
			arg := strings.TrimSpace(raw[open+1 : len(raw)-1])
			if isAggregateFunction(function) {
				expr.Kind = "aggregate"
				expr.Function = function
				if arg == "*" {
					expr.Star = true
				} else {
					index, err := resolveColumn(strings.Trim(arg, `"`), headers)
					if err != nil {
						return nil, err
					}
					expr.Column = index
				}
			}
		}
		if expr.Kind == "column" {
			index, err := resolveColumn(strings.Trim(raw, `"`), headers)
			if err != nil {
				return nil, err
			}
			expr.Column = index
		}
		if expr.Alias != "" {
			expr.Output = expr.Alias
		} else if expr.Kind == "aggregate" {
			if expr.Star {
				expr.Output = strings.ToUpper(expr.Function) + "(*)"
			} else {
				expr.Output = strings.ToUpper(expr.Function) + "(" + headerName(headers, expr.Column) + ")"
			}
		} else {
			expr.Output = headerName(headers, expr.Column)
		}
		expressions = append(expressions, expr)
	}
	return expressions, nil
}

func executeAggregate(query string, source Source, headers []string, sourceRows [][]string, expressions []selectExpr, parts queryParts, options Options) (Result, error) {
	groupColumns := make([]int, 0, len(parts.groupBy))
	for _, column := range parts.groupBy {
		index, err := resolveColumn(column, headers)
		if err != nil {
			return Result{}, err
		}
		groupColumns = append(groupColumns, index)
	}
	if len(groupColumns) == 0 {
		for _, expr := range expressions {
			if expr.Kind == "column" {
				groupColumns = append(groupColumns, expr.Column)
			}
		}
	}
	for _, expr := range expressions {
		if expr.Kind == "column" && !containsInt(groupColumns, expr.Column) {
			return Result{}, fmt.Errorf("column %q must appear in GROUP BY or be aggregated", headerName(headers, expr.Column))
		}
	}

	type groupState struct {
		key    string
		values []string
		rows   [][]string
	}
	groups := make(map[string]*groupState)
	order := make([]string, 0)
	for _, record := range sourceRows {
		values := make([]string, len(groupColumns))
		for i, col := range groupColumns {
			values[i] = cell(record, col)
		}
		key := strings.Join(values, "\x00")
		group := groups[key]
		if group == nil {
			group = &groupState{key: key, values: values}
			groups[key] = group
			order = append(order, key)
		}
		group.rows = append(group.rows, record)
	}
	if len(groupColumns) == 0 {
		key := ""
		groups[key] = &groupState{key: key, rows: sourceRows}
		order = []string{key}
	}

	columns := make([]string, len(expressions))
	for i, expr := range expressions {
		columns[i] = expr.Output
	}
	rows := make([][]string, 0, len(order))
	for _, key := range order {
		group := groups[key]
		row := make([]string, len(expressions))
		for i, expr := range expressions {
			if expr.Kind == "column" {
				row[i] = cell(group.rows[0], expr.Column)
				continue
			}
			row[i] = aggregateValue(expr, group.rows)
		}
		rows = append(rows, row)
	}

	if parts.orderBy != "" {
		orderIndex, err := resolveOutputColumn(parts.orderBy, columns)
		if err != nil {
			return Result{}, err
		}
		sort.SliceStable(rows, func(i, j int) bool {
			less := compareValues(rows[i][orderIndex], rows[j][orderIndex]) < 0
			if parts.descending {
				return !less && rows[i][orderIndex] != rows[j][orderIndex]
			}
			return less
		})
	}

	total := len(rows)
	if total > parts.limit {
		rows = rows[:parts.limit]
	}
	return Result{
		Columns:     columns,
		Rows:        rows,
		TotalRows:   total,
		Source:      source.Path,
		Truncated:   total > len(rows),
		ExecutedSQL: strings.TrimSpace(query),
	}, nil
}

func aggregateValue(expr selectExpr, rows [][]string) string {
	switch expr.Function {
	case "count":
		if expr.Star {
			return strconv.Itoa(len(rows))
		}
		count := 0
		for _, row := range rows {
			if strings.TrimSpace(cell(row, expr.Column)) != "" {
				count++
			}
		}
		return strconv.Itoa(count)
	case "sum", "avg":
		sum := 0.0
		count := 0
		for _, row := range rows {
			value, err := strconv.ParseFloat(strings.TrimSpace(cell(row, expr.Column)), 64)
			if err == nil {
				sum += value
				count++
			}
		}
		if expr.Function == "avg" {
			if count == 0 {
				return ""
			}
			return formatFloat(sum / float64(count))
		}
		return formatFloat(sum)
	case "min", "max":
		result := ""
		for _, row := range rows {
			value := cell(row, expr.Column)
			if strings.TrimSpace(value) == "" {
				continue
			}
			if result == "" || (expr.Function == "min" && compareValues(value, result) < 0) || (expr.Function == "max" && compareValues(value, result) > 0) {
				result = value
			}
		}
		return result
	default:
		return ""
	}
}

func splitAlias(token string) (string, string) {
	if pos := findKeyword(token, " as ", 0); pos >= 0 {
		return strings.TrimSpace(token[:pos]), strings.Trim(strings.TrimSpace(token[pos+len(" as "):]), `"`)
	}
	return token, ""
}

func hasAggregation(selectClause string) bool {
	for _, token := range splitComma(selectClause) {
		lower := strings.ToLower(strings.TrimSpace(token))
		for _, function := range []string{"count(", "sum(", "avg(", "min(", "max("} {
			if strings.HasPrefix(lower, function) {
				return true
			}
		}
	}
	return false
}

func isAggregateFunction(name string) bool {
	switch name {
	case "count", "sum", "avg", "min", "max":
		return true
	default:
		return false
	}
}

func resolveOutputColumn(name string, columns []string) (int, error) {
	trimmed := strings.Trim(strings.TrimSpace(name), `"`)
	if strings.HasPrefix(strings.ToLower(trimmed), "c") {
		if index, err := strconv.Atoi(trimmed[1:]); err == nil && index > 0 && index <= len(columns) {
			return index - 1, nil
		}
	}
	for index, column := range columns {
		if strings.EqualFold(column, trimmed) {
			return index, nil
		}
	}
	return 0, fmt.Errorf("unknown result column %q", name)
}

func containsInt(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func resolveColumn(name string, headers []string) (int, error) {
	lower := strings.ToLower(strings.TrimSpace(strings.Trim(name, `"`)))
	if strings.HasPrefix(lower, "c") && len(lower) > 1 {
		if index, err := strconv.Atoi(lower[1:]); err == nil && index > 0 {
			return index - 1, nil
		}
	}
	for index, header := range headers {
		if strings.EqualFold(header, name) {
			return index, nil
		}
	}
	return 0, fmt.Errorf("unknown column %q", name)
}

func normalizeSourceRef(ref string) string {
	ref = strings.TrimSpace(strings.Trim(ref, `"`))
	if strings.HasPrefix(ref, "@file:") {
		return strings.TrimPrefix(ref, "@file:")
	}
	if strings.HasPrefix(ref, "file://") {
		parsed, err := url.Parse(ref)
		if err == nil {
			return parsed.Path
		}
	}
	return ref
}

func defaultHeaders(width int) []string {
	headers := make([]string, width)
	for i := range headers {
		headers[i] = fmt.Sprintf("c%d", i+1)
	}
	return headers
}

func sameRecord(record []string, headers []string) bool {
	if len(record) != len(headers) {
		return false
	}
	for i := range record {
		if record[i] != headers[i] {
			return false
		}
	}
	return true
}

func cell(record []string, col int) string {
	if col < 0 || col >= len(record) {
		return ""
	}
	return record[col]
}

func headerName(headers []string, col int) string {
	if col >= 0 && col < len(headers) && headers[col] != "" {
		return headers[col]
	}
	return fmt.Sprintf("c%d", col+1)
}

func compareValues(left string, right string) int {
	leftNum, leftErr := strconv.ParseFloat(strings.TrimSpace(left), 64)
	rightNum, rightErr := strconv.ParseFloat(strings.TrimSpace(right), 64)
	if leftErr == nil && rightErr == nil {
		if leftNum < rightNum {
			return -1
		}
		if leftNum > rightNum {
			return 1
		}
		return 0
	}
	return strings.Compare(left, right)
}

func findKeyword(text string, keyword string, start int) int {
	lower := strings.ToLower(text)
	inQuote := rune(0)
	for i := start; i <= len(text)-len(keyword); i++ {
		r := rune(text[i])
		if (r == '\'' || r == '"') && (i == 0 || text[i-1] != '\\') {
			if inQuote == 0 {
				inQuote = r
			} else if inQuote == r {
				inQuote = 0
			}
		}
		if inQuote == 0 && strings.HasPrefix(lower[i:], keyword) {
			return i
		}
	}
	return -1
}

func firstPositive(fallback int, values ...int) int {
	result := fallback
	for _, value := range values {
		if value >= 0 && value < result {
			result = value
		}
	}
	return result
}

func splitComma(value string) []string {
	var parts []string
	start := 0
	inQuote := rune(0)
	for index, r := range value {
		if r == '\'' || r == '"' {
			if inQuote == 0 {
				inQuote = r
			} else if inQuote == r {
				inQuote = 0
			}
		}
		if r == ',' && inQuote == 0 {
			parts = append(parts, value[start:index])
			start = index + 1
		}
	}
	parts = append(parts, value[start:])
	return parts
}
