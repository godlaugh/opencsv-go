package sqlexpr

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Program struct {
	root    expr
	headers map[string]int
}

func Compile(source string, headers []string) (*Program, error) {
	tokens, err := lex(source)
	if err != nil {
		return nil, err
	}
	parser := parser{tokens: tokens}
	root, err := parser.parse()
	if err != nil {
		return nil, err
	}
	headerMap := make(map[string]int, len(headers))
	for index, name := range headers {
		headerMap[strings.ToLower(name)] = index
	}
	return &Program{root: root, headers: headerMap}, nil
}

func (program *Program) Eval(record []string) (bool, error) {
	value, err := program.root.eval(evalContext{record: record, headers: program.headers})
	if err != nil {
		return false, err
	}
	return truthy(value), nil
}

func Match(source string, headers []string, record []string) (bool, error) {
	program, err := Compile(source, headers)
	if err != nil {
		return false, err
	}
	return program.Eval(record)
}

type valueKind int

const (
	valueNull valueKind = iota
	valueString
	valueNumber
	valueBool
)

type value struct {
	kind valueKind
	text string
	num  float64
	bool bool
}

func nullValue() value {
	return value{kind: valueNull}
}

func stringValue(text string) value {
	return value{kind: valueString, text: text}
}

func numberValue(num float64) value {
	return value{kind: valueNumber, num: num}
}

func boolValue(b bool) value {
	return value{kind: valueBool, bool: b}
}

func (v value) string() string {
	switch v.kind {
	case valueNull:
		return ""
	case valueString:
		return v.text
	case valueNumber:
		return strconv.FormatFloat(v.num, 'f', -1, 64)
	case valueBool:
		if v.bool {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func (v value) number() (float64, bool) {
	switch v.kind {
	case valueNumber:
		return v.num, true
	case valueBool:
		if v.bool {
			return 1, true
		}
		return 0, true
	case valueString:
		if strings.TrimSpace(v.text) == "" {
			return 0, true
		}
		num, err := strconv.ParseFloat(strings.TrimSpace(v.text), 64)
		return num, err == nil
	default:
		return 0, false
	}
}

func truthy(v value) bool {
	switch v.kind {
	case valueNull:
		return false
	case valueBool:
		return v.bool
	case valueNumber:
		return v.num != 0 && !math.IsNaN(v.num)
	case valueString:
		if strings.TrimSpace(v.text) == "" {
			return false
		}
		if num, ok := v.number(); ok {
			return num != 0
		}
		return true
	default:
		return false
	}
}

type evalContext struct {
	record  []string
	headers map[string]int
}

func (ctx evalContext) resolve(name string) value {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "c") && len(lower) > 1 {
		if index, err := strconv.Atoi(lower[1:]); err == nil && index > 0 && index <= len(ctx.record) {
			return stringValue(ctx.record[index-1])
		}
	}
	if index, ok := ctx.headers[lower]; ok && index >= 0 && index < len(ctx.record) {
		return stringValue(ctx.record[index])
	}
	return nullValue()
}

type expr interface {
	eval(ctx evalContext) (value, error)
}

type literalExpr struct {
	value value
}

func (expr literalExpr) eval(evalContext) (value, error) {
	return expr.value, nil
}

type identExpr struct {
	name string
}

func (expr identExpr) eval(ctx evalContext) (value, error) {
	return ctx.resolve(expr.name), nil
}

type unaryExpr struct {
	op    string
	right expr
}

func (expr unaryExpr) eval(ctx evalContext) (value, error) {
	right, err := expr.right.eval(ctx)
	if err != nil {
		return nullValue(), err
	}
	switch expr.op {
	case "not":
		return boolValue(!truthy(right)), nil
	case "-":
		num, ok := right.number()
		if !ok {
			return nullValue(), fmt.Errorf("cannot negate %q", right.string())
		}
		return numberValue(-num), nil
	default:
		return nullValue(), fmt.Errorf("unsupported unary operator %s", expr.op)
	}
}

type binaryExpr struct {
	left  expr
	op    string
	right expr
}

func (expr binaryExpr) eval(ctx evalContext) (value, error) {
	switch expr.op {
	case "and":
		left, err := expr.left.eval(ctx)
		if err != nil || !truthy(left) {
			return boolValue(false), err
		}
		right, err := expr.right.eval(ctx)
		return boolValue(truthy(right)), err
	case "or":
		left, err := expr.left.eval(ctx)
		if err != nil {
			return boolValue(false), err
		}
		if truthy(left) {
			return boolValue(true), nil
		}
		right, err := expr.right.eval(ctx)
		return boolValue(truthy(right)), err
	}

	left, err := expr.left.eval(ctx)
	if err != nil {
		return nullValue(), err
	}
	right, err := expr.right.eval(ctx)
	if err != nil {
		return nullValue(), err
	}

	switch expr.op {
	case "+", "-", "*", "/":
		leftNum, leftOK := left.number()
		rightNum, rightOK := right.number()
		if !leftOK || !rightOK {
			return nullValue(), fmt.Errorf("operator %s requires numbers", expr.op)
		}
		switch expr.op {
		case "+":
			return numberValue(leftNum + rightNum), nil
		case "-":
			return numberValue(leftNum - rightNum), nil
		case "*":
			return numberValue(leftNum * rightNum), nil
		case "/":
			if rightNum == 0 {
				return nullValue(), errors.New("division by zero")
			}
			return numberValue(leftNum / rightNum), nil
		}
	case "=", "==", "!=", "<>", "<", "<=", ">", ">=":
		return boolValue(compare(left, right, expr.op)), nil
	case "like", "not like":
		matched := like(left.string(), right.string())
		if expr.op == "not like" {
			matched = !matched
		}
		return boolValue(matched), nil
	}
	return nullValue(), fmt.Errorf("unsupported binary operator %s", expr.op)
}

type isNullExpr struct {
	target expr
	not    bool
}

func (expr isNullExpr) eval(ctx evalContext) (value, error) {
	target, err := expr.target.eval(ctx)
	if err != nil {
		return nullValue(), err
	}
	isNull := target.kind == valueNull
	if expr.not {
		isNull = !isNull
	}
	return boolValue(isNull), nil
}

type callExpr struct {
	name string
	args []expr
}

func (expr callExpr) eval(ctx evalContext) (value, error) {
	name := strings.ToLower(expr.name)
	values := make([]value, len(expr.args))
	for i, arg := range expr.args {
		value, err := arg.eval(ctx)
		if err != nil {
			return nullValue(), err
		}
		values[i] = value
	}
	return callFunction(name, values)
}

type castExpr struct {
	input expr
	typ   string
	try   bool
}

func (expr castExpr) eval(ctx evalContext) (value, error) {
	input, err := expr.input.eval(ctx)
	if err != nil {
		return nullValue(), err
	}
	casted, err := castValue(input, expr.typ)
	if err != nil && expr.try {
		return nullValue(), nil
	}
	return casted, err
}

func compare(left value, right value, op string) bool {
	var result int
	if left.kind == valueNull || right.kind == valueNull {
		result = strings.Compare(left.string(), right.string())
	} else if leftNum, leftOK := left.number(); leftOK {
		if rightNum, rightOK := right.number(); rightOK {
			switch {
			case leftNum < rightNum:
				result = -1
			case leftNum > rightNum:
				result = 1
			default:
				result = 0
			}
		} else {
			result = strings.Compare(left.string(), right.string())
		}
	} else {
		result = strings.Compare(left.string(), right.string())
	}

	switch op {
	case "=", "==":
		return result == 0
	case "!=", "<>":
		return result != 0
	case "<":
		return result < 0
	case "<=":
		return result <= 0
	case ">":
		return result > 0
	case ">=":
		return result >= 0
	default:
		return false
	}
}

func like(text string, pattern string) bool {
	var builder strings.Builder
	builder.WriteString("^")
	for _, r := range pattern {
		switch r {
		case '%':
			builder.WriteString(".*")
		case '_':
			builder.WriteString(".")
		default:
			builder.WriteString(regexp.QuoteMeta(string(r)))
		}
	}
	builder.WriteString("$")
	matched, err := regexp.MatchString(builder.String(), text)
	return err == nil && matched
}

func callFunction(name string, args []value) (value, error) {
	switch name {
	case "length":
		if len(args) != 1 {
			return nullValue(), fmt.Errorf("%s expects 1 argument", name)
		}
		return numberValue(float64(len([]rune(args[0].string())))), nil
	case "lower":
		return unaryString(name, args, strings.ToLower)
	case "upper":
		return unaryString(name, args, strings.ToUpper)
	case "trim":
		return unaryString(name, args, strings.TrimSpace)
	case "substring":
		if len(args) < 2 || len(args) > 3 {
			return nullValue(), errors.New("substring expects 2 or 3 arguments")
		}
		runes := []rune(args[0].string())
		start, ok := args[1].number()
		if !ok {
			return nullValue(), errors.New("substring start must be numeric")
		}
		index := max(0, int(start)-1)
		end := len(runes)
		if len(args) == 3 {
			length, ok := args[2].number()
			if !ok {
				return nullValue(), errors.New("substring length must be numeric")
			}
			end = min(len(runes), index+max(0, int(length)))
		}
		if index > len(runes) {
			return stringValue(""), nil
		}
		return stringValue(string(runes[index:end])), nil
	case "concat":
		var builder strings.Builder
		for _, arg := range args {
			builder.WriteString(arg.string())
		}
		return stringValue(builder.String()), nil
	case "replace":
		if len(args) != 3 {
			return nullValue(), errors.New("replace expects 3 arguments")
		}
		return stringValue(strings.ReplaceAll(args[0].string(), args[1].string(), args[2].string())), nil
	case "abs":
		return unaryNumber(name, args, math.Abs)
	case "ceil":
		return unaryNumber(name, args, math.Ceil)
	case "floor":
		return unaryNumber(name, args, math.Floor)
	case "round":
		if len(args) == 0 || len(args) > 2 {
			return nullValue(), errors.New("round expects 1 or 2 arguments")
		}
		num, ok := args[0].number()
		if !ok {
			return nullValue(), errors.New("round requires a number")
		}
		precision := 0
		if len(args) == 2 {
			precisionNum, ok := args[1].number()
			if !ok {
				return nullValue(), errors.New("round precision must be numeric")
			}
			precision = int(precisionNum)
		}
		factor := math.Pow10(precision)
		return numberValue(math.Round(num*factor) / factor), nil
	case "power":
		if len(args) != 2 {
			return nullValue(), errors.New("power expects 2 arguments")
		}
		left, leftOK := args[0].number()
		right, rightOK := args[1].number()
		if !leftOK || !rightOK {
			return nullValue(), errors.New("power requires numbers")
		}
		return numberValue(math.Pow(left, right)), nil
	case "sqrt":
		return unaryNumber(name, args, math.Sqrt)
	case "mod":
		if len(args) != 2 {
			return nullValue(), errors.New("mod expects 2 arguments")
		}
		left, leftOK := args[0].number()
		right, rightOK := args[1].number()
		if !leftOK || !rightOK {
			return nullValue(), errors.New("mod requires numbers")
		}
		return numberValue(math.Mod(left, right)), nil
	case "trunc":
		if len(args) == 0 || len(args) > 2 {
			return nullValue(), errors.New("trunc expects 1 or 2 arguments")
		}
		num, ok := args[0].number()
		if !ok {
			return nullValue(), errors.New("trunc requires a number")
		}
		precision := 0
		if len(args) == 2 {
			precisionNum, ok := args[1].number()
			if !ok {
				return nullValue(), errors.New("trunc precision must be numeric")
			}
			precision = int(precisionNum)
		}
		factor := math.Pow10(precision)
		return numberValue(math.Trunc(num*factor) / factor), nil
	case "coalesce":
		for _, arg := range args {
			if arg.kind != valueNull {
				return arg, nil
			}
		}
		return nullValue(), nil
	case "nullif":
		if len(args) != 2 {
			return nullValue(), errors.New("nullif expects 2 arguments")
		}
		if compare(args[0], args[1], "=") {
			return nullValue(), nil
		}
		return args[0], nil
	default:
		return nullValue(), fmt.Errorf("unsupported function %s", name)
	}
}

func unaryString(name string, args []value, fn func(string) string) (value, error) {
	if len(args) != 1 {
		return nullValue(), fmt.Errorf("%s expects 1 argument", name)
	}
	return stringValue(fn(args[0].string())), nil
}

func unaryNumber(name string, args []value, fn func(float64) float64) (value, error) {
	if len(args) != 1 {
		return nullValue(), fmt.Errorf("%s expects 1 argument", name)
	}
	num, ok := args[0].number()
	if !ok {
		return nullValue(), fmt.Errorf("%s requires a number", name)
	}
	return numberValue(fn(num)), nil
}

func castValue(input value, rawType string) (value, error) {
	typ := strings.ToLower(strings.Join(strings.Fields(rawType), " "))
	switch typ {
	case "string", "text", "varchar":
		return stringValue(input.string()), nil
	case "number", "double precision", "float":
		num, ok := input.number()
		if !ok {
			return nullValue(), fmt.Errorf("cannot cast %q as number", input.string())
		}
		return numberValue(num), nil
	case "bigint":
		num, err := strconv.ParseInt(strings.TrimSpace(input.string()), 10, 64)
		if err != nil {
			return nullValue(), fmt.Errorf("cannot cast %q as bigint", input.string())
		}
		return numberValue(float64(num)), nil
	case "boolean":
		raw := strings.ToLower(strings.TrimSpace(input.string()))
		if raw == "true" || raw == "1" {
			return boolValue(true), nil
		}
		if raw == "false" || raw == "0" || raw == "" {
			return boolValue(false), nil
		}
		return nullValue(), fmt.Errorf("cannot cast %q as boolean", input.string())
	case "date":
		if _, err := time.Parse("2006-01-02", strings.TrimSpace(input.string())); err != nil {
			return nullValue(), err
		}
		return stringValue(strings.TrimSpace(input.string())), nil
	case "time":
		if _, err := time.Parse("15:04:05", strings.TrimSpace(input.string())); err != nil {
			return nullValue(), err
		}
		return stringValue(strings.TrimSpace(input.string())), nil
	case "timestamp":
		raw := strings.TrimSpace(input.string())
		if _, err := time.Parse("2006-01-02 15:04:05", raw); err != nil {
			return nullValue(), err
		}
		return stringValue(raw), nil
	default:
		return nullValue(), fmt.Errorf("unsupported cast type %s", rawType)
	}
}

type tokenKind int

const (
	tokenEOF tokenKind = iota
	tokenIdentifier
	tokenString
	tokenNumber
	tokenOperator
	tokenLeftParen
	tokenRightParen
	tokenComma
)

type token struct {
	kind tokenKind
	text string
}

func lex(source string) ([]token, error) {
	var tokens []token
	for index := 0; index < len(source); {
		r := rune(source[index])
		if unicode.IsSpace(r) {
			index++
			continue
		}
		if isIdentifierStart(r) {
			start := index
			index++
			for index < len(source) && isIdentifierPart(rune(source[index])) {
				index++
			}
			tokens = append(tokens, token{kind: tokenIdentifier, text: source[start:index]})
			continue
		}
		if r == '"' {
			text, next, err := readQuoted(source, index, '"')
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token{kind: tokenIdentifier, text: text})
			index = next
			continue
		}
		if r == '\'' {
			text, next, err := readQuoted(source, index, '\'')
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token{kind: tokenString, text: text})
			index = next
			continue
		}
		if unicode.IsDigit(r) || (r == '.' && index+1 < len(source) && unicode.IsDigit(rune(source[index+1]))) {
			start := index
			index = scanNumber(source, index)
			tokens = append(tokens, token{kind: tokenNumber, text: source[start:index]})
			continue
		}
		switch r {
		case '(':
			tokens = append(tokens, token{kind: tokenLeftParen, text: "("})
			index++
		case ')':
			tokens = append(tokens, token{kind: tokenRightParen, text: ")"})
			index++
		case ',':
			tokens = append(tokens, token{kind: tokenComma, text: ","})
			index++
		case '=', '<', '>', '!', '+', '-', '*', '/':
			start := index
			index++
			if index < len(source) {
				next := source[index]
				if next == '=' || (source[start] == '<' && next == '>') {
					index++
				}
			}
			tokens = append(tokens, token{kind: tokenOperator, text: source[start:index]})
		default:
			return nil, fmt.Errorf("unexpected character %q", r)
		}
	}
	tokens = append(tokens, token{kind: tokenEOF})
	return tokens, nil
}

func scanNumber(source string, index int) int {
	for index < len(source) && unicode.IsDigit(rune(source[index])) {
		index++
	}
	if index < len(source) && source[index] == '.' {
		index++
		for index < len(source) && unicode.IsDigit(rune(source[index])) {
			index++
		}
	}
	if index < len(source) && (source[index] == 'e' || source[index] == 'E') {
		exponent := index + 1
		if exponent < len(source) && (source[exponent] == '+' || source[exponent] == '-') {
			exponent++
		}
		digits := exponent
		for exponent < len(source) && unicode.IsDigit(rune(source[exponent])) {
			exponent++
		}
		if exponent > digits {
			index = exponent
		}
	}
	return index
}

func readQuoted(source string, start int, quote byte) (string, int, error) {
	var builder strings.Builder
	for index := start + 1; index < len(source); index++ {
		ch := source[index]
		if ch == quote {
			if index+1 < len(source) && source[index+1] == quote {
				builder.WriteByte(quote)
				index++
				continue
			}
			return builder.String(), index + 1, nil
		}
		builder.WriteByte(ch)
	}
	return "", start, errors.New("unterminated quoted value")
}

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentifierPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

type parser struct {
	tokens []token
	pos    int
}

func (parser *parser) parse() (expr, error) {
	result, err := parser.parseOr()
	if err != nil {
		return nil, err
	}
	if parser.peek().kind != tokenEOF {
		return nil, fmt.Errorf("unexpected token %q", parser.peek().text)
	}
	return result, nil
}

func (parser *parser) parseOr() (expr, error) {
	left, err := parser.parseAnd()
	if err != nil {
		return nil, err
	}
	for parser.matchKeyword("or") {
		right, err := parser.parseAnd()
		if err != nil {
			return nil, err
		}
		left = binaryExpr{left: left, op: "or", right: right}
	}
	return left, nil
}

func (parser *parser) parseAnd() (expr, error) {
	left, err := parser.parseNot()
	if err != nil {
		return nil, err
	}
	for parser.matchKeyword("and") {
		right, err := parser.parseNot()
		if err != nil {
			return nil, err
		}
		left = binaryExpr{left: left, op: "and", right: right}
	}
	return left, nil
}

func (parser *parser) parseNot() (expr, error) {
	if parser.matchKeyword("not") {
		right, err := parser.parseNot()
		if err != nil {
			return nil, err
		}
		return unaryExpr{op: "not", right: right}, nil
	}
	return parser.parseComparison()
}

func (parser *parser) parseComparison() (expr, error) {
	left, err := parser.parseAdditive()
	if err != nil {
		return nil, err
	}

	if parser.matchKeyword("is") {
		not := parser.matchKeyword("not")
		if !parser.matchKeyword("null") {
			return nil, errors.New("expected null after is")
		}
		return isNullExpr{target: left, not: not}, nil
	}

	if parser.matchKeyword("not") {
		if !parser.matchKeyword("like") {
			return nil, errors.New("expected like after not")
		}
		right, err := parser.parseAdditive()
		if err != nil {
			return nil, err
		}
		return binaryExpr{left: left, op: "not like", right: right}, nil
	}

	if parser.matchKeyword("like") {
		right, err := parser.parseAdditive()
		if err != nil {
			return nil, err
		}
		return binaryExpr{left: left, op: "like", right: right}, nil
	}

	if parser.peek().kind == tokenOperator && isComparisonOperator(parser.peek().text) {
		op := parser.advance().text
		right, err := parser.parseAdditive()
		if err != nil {
			return nil, err
		}
		return binaryExpr{left: left, op: strings.ToLower(op), right: right}, nil
	}
	return left, nil
}

func (parser *parser) parseAdditive() (expr, error) {
	left, err := parser.parseMultiplicative()
	if err != nil {
		return nil, err
	}
	for parser.peek().kind == tokenOperator && (parser.peek().text == "+" || parser.peek().text == "-") {
		op := parser.advance().text
		right, err := parser.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		left = binaryExpr{left: left, op: op, right: right}
	}
	return left, nil
}

func (parser *parser) parseMultiplicative() (expr, error) {
	left, err := parser.parseUnary()
	if err != nil {
		return nil, err
	}
	for parser.peek().kind == tokenOperator && (parser.peek().text == "*" || parser.peek().text == "/") {
		op := parser.advance().text
		right, err := parser.parseUnary()
		if err != nil {
			return nil, err
		}
		left = binaryExpr{left: left, op: op, right: right}
	}
	return left, nil
}

func (parser *parser) parseUnary() (expr, error) {
	if parser.peek().kind == tokenOperator && parser.peek().text == "-" {
		parser.advance()
		right, err := parser.parseUnary()
		if err != nil {
			return nil, err
		}
		return unaryExpr{op: "-", right: right}, nil
	}
	return parser.parsePrimary()
}

func (parser *parser) parsePrimary() (expr, error) {
	current := parser.advance()
	switch current.kind {
	case tokenString:
		return literalExpr{value: stringValue(current.text)}, nil
	case tokenNumber:
		num, err := strconv.ParseFloat(current.text, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q", current.text)
		}
		return literalExpr{value: numberValue(num)}, nil
	case tokenIdentifier:
		lower := strings.ToLower(current.text)
		if lower == "null" {
			return literalExpr{value: nullValue()}, nil
		}
		if lower == "true" || lower == "false" {
			return literalExpr{value: boolValue(lower == "true")}, nil
		}
		if lower == "date" || lower == "time" || lower == "timestamp" || lower == "interval" {
			if parser.peek().kind != tokenString {
				return nil, fmt.Errorf("%s literal expects a quoted value", lower)
			}
			return literalExpr{value: stringValue(parser.advance().text)}, nil
		}
		if lower == "cast" || lower == "try_cast" {
			return parser.parseCast(lower == "try_cast")
		}
		if parser.match(tokenLeftParen) {
			args, err := parser.parseArguments()
			if err != nil {
				return nil, err
			}
			return callExpr{name: current.text, args: args}, nil
		}
		return identExpr{name: current.text}, nil
	case tokenLeftParen:
		nested, err := parser.parseOr()
		if err != nil {
			return nil, err
		}
		if !parser.match(tokenRightParen) {
			return nil, errors.New("expected closing parenthesis")
		}
		return nested, nil
	default:
		return nil, fmt.Errorf("unexpected token %q", current.text)
	}
}

func (parser *parser) parseCast(try bool) (expr, error) {
	if !parser.match(tokenLeftParen) {
		return nil, errors.New("cast expects opening parenthesis")
	}
	input, err := parser.parseOr()
	if err != nil {
		return nil, err
	}
	if !parser.matchKeyword("as") {
		return nil, errors.New("cast expects as")
	}
	var typ []string
	for parser.peek().kind == tokenIdentifier {
		typ = append(typ, parser.advance().text)
	}
	if len(typ) == 0 {
		return nil, errors.New("cast expects a target type")
	}
	if !parser.match(tokenRightParen) {
		return nil, errors.New("cast expects closing parenthesis")
	}
	return castExpr{input: input, typ: strings.Join(typ, " "), try: try}, nil
}

func (parser *parser) parseArguments() ([]expr, error) {
	if parser.match(tokenRightParen) {
		return nil, nil
	}
	var args []expr
	for {
		arg, err := parser.parseOr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if parser.match(tokenRightParen) {
			return args, nil
		}
		if !parser.match(tokenComma) {
			return nil, errors.New("expected comma or closing parenthesis")
		}
	}
}

func (parser *parser) peek() token {
	return parser.tokens[parser.pos]
}

func (parser *parser) advance() token {
	current := parser.tokens[parser.pos]
	if parser.pos < len(parser.tokens)-1 {
		parser.pos++
	}
	return current
}

func (parser *parser) match(kind tokenKind) bool {
	if parser.peek().kind != kind {
		return false
	}
	parser.advance()
	return true
}

func (parser *parser) matchKeyword(keyword string) bool {
	if parser.peek().kind != tokenIdentifier || strings.ToLower(parser.peek().text) != keyword {
		return false
	}
	parser.advance()
	return true
}

func isComparisonOperator(op string) bool {
	switch op {
	case "=", "==", "!=", "<>", "<", "<=", ">", ">=":
		return true
	default:
		return false
	}
}
