package sqlconsole

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"unicode"

	_ "modernc.org/sqlite"
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
	sqlText := trimSQL(query)
	if !isReadQuery(sqlText) {
		return Result{}, errors.New("only SELECT queries are supported")
	}

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return Result{}, err
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err := db.Exec("PRAGMA temp_store = MEMORY"); err != nil {
		return Result{}, err
	}
	if _, err := db.Exec("PRAGMA journal_mode = MEMORY"); err != nil {
		return Result{}, err
	}

	refs := findSourceRefs(sqlText)
	sourceNames := make([]string, 0, len(refs))
	for _, ref := range refs {
		source, err := resolver(normalizeSourceRef(ref))
		if err != nil {
			return Result{}, err
		}
		if err := loadSource(db, ref, source); err != nil {
			return Result{}, err
		}
		sourceNames = append(sourceNames, source.Path)
	}

	countQuery := sqlText
	if base, ok := withoutTopLevelLimit(sqlText); ok {
		countQuery = base
	}
	total, countErr := countRows(db, countQuery)
	limit := options.DefaultLimit
	if hasTopLevelLimit(sqlText) {
		limit = options.MaxLimit
	}
	limitedQuery := fmt.Sprintf("SELECT * FROM (%s) LIMIT %d", sqlText, limit+1)
	rows, columns, err := queryRows(db, limitedQuery, limit+1)
	if err != nil {
		return Result{}, err
	}
	if countErr != nil {
		total = len(rows)
	}
	truncated := total > len(rows)
	if len(rows) > limit {
		rows = rows[:limit]
		truncated = true
	}
	if total < len(rows) {
		total = len(rows)
	}
	return Result{
		Columns:     columns,
		Rows:        rows,
		TotalRows:   total,
		Source:      strings.Join(sourceNames, ", "),
		Truncated:   truncated,
		ExecutedSQL: strings.TrimSpace(query),
	}, nil
}

func trimSQL(query string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(query), ";"))
}

func isReadQuery(query string) bool {
	lower := strings.ToLower(strings.TrimSpace(query))
	return strings.HasPrefix(lower, "select ") || strings.HasPrefix(lower, "with ")
}

func loadSource(db *sql.DB, tableName string, source Source) error {
	headers := sourceHeaders(source)
	if len(headers) == 0 {
		return nil
	}
	columns := uniqueHeaders(headers)
	defs := make([]string, len(columns))
	for index, column := range columns {
		defs[index] = quoteIdent(column) + " TEXT"
	}
	if _, err := db.Exec("CREATE TABLE " + quoteIdent(tableName) + " (" + strings.Join(defs, ", ") + ")"); err != nil {
		return err
	}

	placeholders := make([]string, len(columns))
	for index := range placeholders {
		placeholders[index] = "?"
	}
	statement, err := db.Prepare("INSERT INTO " + quoteIdent(tableName) + " VALUES (" + strings.Join(placeholders, ", ") + ")")
	if err != nil {
		return err
	}
	defer statement.Close()

	dataStart := 0
	if len(source.Records) > 0 && sameRecord(source.Records[0], headers) {
		dataStart = 1
	}
	for _, record := range source.Records[dataStart:] {
		values := make([]any, len(columns))
		for index := range columns {
			values[index] = cell(record, index)
		}
		if _, err := statement.Exec(values...); err != nil {
			return err
		}
	}
	return nil
}

func sourceHeaders(source Source) []string {
	width := len(source.Headers)
	for _, record := range source.Records {
		if len(record) > width {
			width = len(record)
		}
	}
	if width == 0 {
		return nil
	}
	headers := append([]string(nil), source.Headers...)
	if len(headers) == 0 {
		headers = defaultHeaders(width)
	}
	for len(headers) < width {
		headers = append(headers, fmt.Sprintf("c%d", len(headers)+1))
	}
	return headers
}

func uniqueHeaders(headers []string) []string {
	seen := make(map[string]int, len(headers))
	columns := make([]string, len(headers))
	for index, header := range headers {
		name := strings.TrimSpace(header)
		if name == "" {
			name = fmt.Sprintf("c%d", index+1)
		}
		key := strings.ToLower(name)
		seen[key]++
		if seen[key] > 1 {
			name = fmt.Sprintf("%s_%d", name, seen[key])
			key = strings.ToLower(name)
			for seen[key] > 0 {
				seen[strings.ToLower(name)]++
				name = fmt.Sprintf("%s_%d", strings.TrimSpace(header), seen[key]+1)
				key = strings.ToLower(name)
			}
			seen[key] = 1
		}
		columns[index] = name
	}
	return columns
}

func countRows(db *sql.DB, query string) (int, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM (" + trimSQL(query) + ")")
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func queryRows(db *sql.DB, query string, limit int) ([][]string, []string, error) {
	result, err := db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer result.Close()

	columns, err := result.Columns()
	if err != nil {
		return nil, nil, err
	}
	rows := make([][]string, 0)
	for result.Next() {
		values := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for index := range values {
			dest[index] = &values[index]
		}
		if err := result.Scan(dest...); err != nil {
			return nil, nil, err
		}
		row := make([]string, len(columns))
		for index, value := range values {
			if value.Valid {
				row[index] = value.String
			}
		}
		rows = append(rows, row)
		if len(rows) >= limit {
			break
		}
	}
	if err := result.Err(); err != nil {
		return nil, nil, err
	}
	return rows, columns, nil
}

func quoteIdent(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func hasTopLevelLimit(query string) bool {
	_, ok := withoutTopLevelLimit(query)
	return ok
}

func withoutTopLevelLimit(query string) (string, bool) {
	limitPos := findTopLevelKeyword(query, "limit")
	if limitPos >= 0 {
		return strings.TrimSpace(query[:limitPos]), true
	}
	return query, false
}

func findTopLevelKeyword(query string, keyword string) int {
	tokens := scanSQLTokens(query)
	for index := len(tokens) - 1; index >= 0; index-- {
		token := tokens[index]
		if token.depth == 0 && !token.quoted && strings.EqualFold(token.text, keyword) {
			return token.start
		}
	}
	return -1
}

func findSourceRefs(query string) []string {
	tokens := scanSQLTokens(query)
	seen := map[string]bool{}
	refs := make([]string, 0)
	for index, token := range tokens {
		if token.quoted || !isTableKeyword(token.text) {
			continue
		}
		ref, ok := nextSourceRef(tokens, index+1)
		if !ok || !isLoadableRef(ref) || seen[ref] {
			continue
		}
		seen[ref] = true
		refs = append(refs, ref)
	}
	return refs
}

func isTableKeyword(text string) bool {
	return strings.EqualFold(text, "from") || strings.EqualFold(text, "join")
}

func nextSourceRef(tokens []sqlToken, start int) (string, bool) {
	for index := start; index < len(tokens); index++ {
		token := tokens[index]
		if token.text == "(" || token.text == "," {
			return "", false
		}
		if strings.EqualFold(token.text, "select") {
			return "", false
		}
		return token.text, true
	}
	return "", false
}

func isLoadableRef(ref string) bool {
	lower := strings.ToLower(strings.TrimSpace(ref))
	return lower == "current" ||
		lower == "this" ||
		strings.HasPrefix(lower, "@file:") ||
		strings.HasPrefix(lower, "file://") ||
		strings.Contains(lower, "/") ||
		strings.Contains(lower, `\`) ||
		strings.Contains(lower, ".csv") ||
		strings.Contains(lower, ".tsv") ||
		strings.Contains(lower, ".xlsx")
}

type sqlToken struct {
	text   string
	start  int
	depth  int
	quoted bool
}

func scanSQLTokens(query string) []sqlToken {
	tokens := make([]sqlToken, 0)
	depth := 0
	for index := 0; index < len(query); {
		r := rune(query[index])
		if unicode.IsSpace(r) {
			index++
			continue
		}
		if r == '(' {
			tokens = append(tokens, sqlToken{text: "(", start: index, depth: depth})
			depth++
			index++
			continue
		}
		if r == ')' {
			if depth > 0 {
				depth--
			}
			tokens = append(tokens, sqlToken{text: ")", start: index, depth: depth})
			index++
			continue
		}
		if r == ',' {
			tokens = append(tokens, sqlToken{text: ",", start: index, depth: depth})
			index++
			continue
		}
		if r == '\'' || r == '"' || r == '`' {
			quote := byte(query[index])
			start := index
			index++
			var builder strings.Builder
			for index < len(query) {
				if query[index] == quote {
					if index+1 < len(query) && query[index+1] == quote {
						builder.WriteByte(quote)
						index += 2
						continue
					}
					index++
					break
				}
				builder.WriteByte(query[index])
				index++
			}
			tokens = append(tokens, sqlToken{text: builder.String(), start: start, depth: depth, quoted: true})
			continue
		}
		start := index
		for index < len(query) {
			r = rune(query[index])
			if unicode.IsSpace(r) || strings.ContainsRune("(),", r) {
				break
			}
			index++
		}
		tokens = append(tokens, sqlToken{text: query[start:index], start: start, depth: depth})
	}
	return tokens
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
