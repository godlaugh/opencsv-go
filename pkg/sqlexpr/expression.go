package sqlexpr

import (
	"errors"
	"fmt"
	"math"
	"math/big"
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
	valueBigInt
	valueBool
	valueDate
	valueTime
	valueTimestamp
	valueInterval
)

type value struct {
	kind            valueKind
	text            string
	num             float64
	bigint          *big.Int
	bool            bool
	intervalMonths  int64
	intervalSeconds float64
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

func bigIntValue(num *big.Int) value {
	if num == nil {
		return nullValue()
	}
	return value{kind: valueBigInt, bigint: new(big.Int).Set(num)}
}

func boolValue(b bool) value {
	return value{kind: valueBool, bool: b}
}

func dateValue(text string) value {
	return value{kind: valueDate, text: text}
}

func timeValue(text string) value {
	return value{kind: valueTime, text: text}
}

func timestampValue(text string) value {
	return value{kind: valueTimestamp, text: text}
}

func intervalValue(months int64, seconds float64) value {
	return value{kind: valueInterval, intervalMonths: months, intervalSeconds: seconds}
}

func (v value) string() string {
	switch v.kind {
	case valueNull:
		return ""
	case valueString:
		return v.text
	case valueNumber:
		return strconv.FormatFloat(v.num, 'f', -1, 64)
	case valueBigInt:
		if v.bigint == nil {
			return ""
		}
		return v.bigint.String()
	case valueBool:
		if v.bool {
			return "true"
		}
		return "false"
	case valueDate, valueTime, valueTimestamp:
		return v.text
	case valueInterval:
		return fmt.Sprintf("interval:%d:%g", v.intervalMonths, v.intervalSeconds)
	default:
		return ""
	}
}

func (v value) number() (float64, bool) {
	switch v.kind {
	case valueNumber:
		return v.num, true
	case valueBigInt:
		if v.bigint == nil {
			return 0, false
		}
		result, _ := new(big.Float).SetInt(v.bigint).Float64()
		return result, true
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
	case valueBigInt:
		return v.bigint != nil && v.bigint.Sign() != 0
	case valueString:
		if strings.TrimSpace(v.text) == "" {
			return false
		}
		if num, ok := v.number(); ok {
			return num != 0
		}
		return true
	case valueDate, valueTime, valueTimestamp:
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
		if right.kind == valueBigInt {
			return bigIntValue(new(big.Int).Neg(right.bigint)), nil
		}
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
		if result, ok, err := temporalArithmetic(left, right, expr.op); ok || err != nil {
			return result, err
		}
		if result, ok, err := intervalArithmetic(left, right, expr.op); ok || err != nil {
			return result, err
		}
		if result, ok, err := bigIntArithmetic(left, right, expr.op); ok || err != nil {
			return result, err
		}
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
	if left.kind == valueInterval || right.kind == valueInterval {
		result = compareIntervals(left, right)
	} else if temporalResult, ok := compareTemporals(left, right); ok {
		result = temporalResult
	} else if isTemporal(left) || isTemporal(right) {
		result = 1
	} else if left.kind == valueNull || right.kind == valueNull {
		result = strings.Compare(left.string(), right.string())
	} else if bigintResult, ok := compareBigInts(left, right); ok {
		result = bigintResult
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

func bigIntArithmetic(left value, right value, op string) (value, bool, error) {
	leftInt, leftOK := exactBigInt(left)
	rightInt, rightOK := exactBigInt(right)
	if !leftOK || !rightOK {
		return nullValue(), false, nil
	}
	result := new(big.Int)
	switch op {
	case "+":
		result.Add(leftInt, rightInt)
	case "-":
		result.Sub(leftInt, rightInt)
	case "*":
		result.Mul(leftInt, rightInt)
	default:
		return nullValue(), false, nil
	}
	return bigIntValue(result), true, nil
}

func compareBigInts(left value, right value) (int, bool) {
	if left.kind != valueBigInt && right.kind != valueBigInt {
		return 0, false
	}
	leftInt, leftOK := exactBigInt(left)
	rightInt, rightOK := exactBigInt(right)
	if !leftOK || !rightOK {
		return 0, false
	}
	return leftInt.Cmp(rightInt), true
}

func exactBigInt(input value) (*big.Int, bool) {
	switch input.kind {
	case valueBigInt:
		if input.bigint == nil {
			return nil, false
		}
		return input.bigint, true
	case valueString:
		return parseBigIntText(input.string())
	case valueNumber:
		if math.IsNaN(input.num) || math.IsInf(input.num, 0) || math.Trunc(input.num) != input.num {
			return nil, false
		}
		const maxSafeInteger = 9007199254740991
		if input.num < -maxSafeInteger || input.num > maxSafeInteger {
			return nil, false
		}
		return big.NewInt(int64(input.num)), true
	default:
		return nil, false
	}
}

func parseBigIntText(raw string) (*big.Int, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, false
	}
	value, ok := new(big.Int).SetString(trimmed, 10)
	return value, ok
}

func isTemporal(input value) bool {
	switch input.kind {
	case valueDate, valueTime, valueTimestamp:
		return true
	default:
		return false
	}
}

func compareTemporals(left value, right value) (int, bool) {
	if !isTemporal(left) || left.kind != right.kind {
		return 0, false
	}
	return strings.Compare(left.string(), right.string()), true
}

func intervalArithmetic(left value, right value, op string) (value, bool, error) {
	if left.kind != valueInterval || right.kind != valueInterval {
		return nullValue(), false, nil
	}
	switch op {
	case "+":
		return intervalValue(left.intervalMonths+right.intervalMonths, left.intervalSeconds+right.intervalSeconds), true, nil
	case "-":
		return intervalValue(left.intervalMonths-right.intervalMonths, left.intervalSeconds-right.intervalSeconds), true, nil
	default:
		return nullValue(), false, nil
	}
}

func temporalArithmetic(left value, right value, op string) (value, bool, error) {
	if op == "+" && left.kind == valueInterval && isTemporal(right) {
		return addIntervalToTemporal(right, left, 1)
	}
	if isTemporal(left) && right.kind == valueInterval {
		switch op {
		case "+":
			return addIntervalToTemporal(left, right, 1)
		case "-":
			return addIntervalToTemporal(left, right, -1)
		}
	}
	if op == "-" && isTemporal(left) && left.kind == right.kind {
		return subtractTemporals(left, right)
	}
	return nullValue(), false, nil
}

func addIntervalToTemporal(temporal value, interval value, sign int) (value, bool, error) {
	months := int(interval.intervalMonths) * sign
	seconds := interval.intervalSeconds * float64(sign)
	switch temporal.kind {
	case valueDate:
		parsed, err := parseDate(temporal.text)
		if err != nil {
			return nullValue(), true, err
		}
		parsed = parsed.AddDate(0, months, 0)
		if whole, ok := wholeDays(seconds); ok {
			return dateValue(parsed.AddDate(0, 0, whole).Format("2006-01-02")), true, nil
		}
		return timestampValue(parsed.Add(durationFromSeconds(seconds)).Format("2006-01-02 15:04:05.999999999")), true, nil
	case valueTimestamp:
		parsed, err := parseTimestamp(temporal.text)
		if err != nil {
			return nullValue(), true, err
		}
		parsed = parsed.AddDate(0, months, 0).Add(durationFromSeconds(seconds))
		return timestampValue(parsed.Format("2006-01-02 15:04:05.999999999")), true, nil
	case valueTime:
		if months != 0 {
			return nullValue(), true, errors.New("cannot add month interval to time")
		}
		parsed, err := parseClockTime(temporal.text)
		if err != nil {
			return nullValue(), true, err
		}
		return timeValue(parsed.Add(durationFromSeconds(seconds)).Format("15:04:05.999999999")), true, nil
	default:
		return nullValue(), false, nil
	}
}

func subtractTemporals(left value, right value) (value, bool, error) {
	leftTime, err := parseTemporalTime(left)
	if err != nil {
		return nullValue(), true, err
	}
	rightTime, err := parseTemporalTime(right)
	if err != nil {
		return nullValue(), true, err
	}
	return intervalValue(0, leftTime.Sub(rightTime).Seconds()), true, nil
}

func parseTemporalTime(input value) (time.Time, error) {
	switch input.kind {
	case valueDate:
		return parseDate(input.text)
	case valueTime:
		return parseClockTime(input.text)
	case valueTimestamp:
		return parseTimestamp(input.text)
	default:
		return time.Time{}, fmt.Errorf("unsupported temporal value %s", input.string())
	}
}

func wholeDays(seconds float64) (int, bool) {
	days := seconds / 86400
	rounded := math.Round(days)
	if math.Abs(days-rounded) > 1e-9 {
		return 0, false
	}
	return int(rounded), true
}

func durationFromSeconds(seconds float64) time.Duration {
	return time.Duration(math.Round(seconds * float64(time.Second)))
}

func compareIntervals(left value, right value) int {
	if left.kind != valueInterval || right.kind != valueInterval {
		return strings.Compare(left.string(), right.string())
	}
	switch {
	case left.intervalMonths < right.intervalMonths:
		return -1
	case left.intervalMonths > right.intervalMonths:
		return 1
	case left.intervalSeconds < right.intervalSeconds:
		return -1
	case left.intervalSeconds > right.intervalSeconds:
		return 1
	default:
		return 0
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
		num, ok := parseBigIntText(input.string())
		if !ok {
			return nullValue(), fmt.Errorf("cannot cast %q as bigint", input.string())
		}
		return bigIntValue(num), nil
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
		raw, err := normalizeDate(input.string())
		if err != nil {
			return nullValue(), err
		}
		return dateValue(raw), nil
	case "time":
		raw, err := normalizeClockTime(input.string())
		if err != nil {
			return nullValue(), err
		}
		return timeValue(raw), nil
	case "timestamp":
		raw, err := normalizeTimestamp(input.string())
		if err != nil {
			return nullValue(), err
		}
		return timestampValue(raw), nil
	case "interval":
		return parseFlexibleInterval(input.string())
	default:
		return nullValue(), fmt.Errorf("unsupported cast type %s", rawType)
	}
}

func normalizeDate(raw string) (string, error) {
	parsed, err := parseDate(raw)
	if err != nil {
		return "", err
	}
	return parsed.Format("2006-01-02"), nil
}

func normalizeClockTime(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	parsed, err := parseClockTime(trimmed)
	if err == nil {
		return parsed.Format("15:04:05.999999999"), nil
	}
	return "", fmt.Errorf("cannot parse %q as time", raw)
}

func normalizeTimestamp(raw string) (string, error) {
	parsed, err := parseTimestamp(raw)
	if err != nil {
		return "", err
	}
	return parsed.Format("2006-01-02 15:04:05.999999999"), nil
}

func parseDate(raw string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(raw))
}

func parseClockTime(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	for _, layout := range []string{"15:04:05.999999999", "15:04:05"} {
		parsed, err := time.Parse(layout, trimmed)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as time", raw)
}

func parseTimestamp(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	for _, layout := range []string{"2006-01-02 15:04:05.999999999", "2006-01-02 15:04:05"} {
		parsed, err := time.Parse(layout, trimmed)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as timestamp", raw)
}

func normalizeTemporalLiteral(kind string, raw string) (string, error) {
	switch kind {
	case "date":
		return normalizeDate(raw)
	case "time":
		return normalizeClockTime(raw)
	case "timestamp":
		return normalizeTimestamp(raw)
	default:
		return "", fmt.Errorf("unsupported temporal literal %s", kind)
	}
}

var intervalFieldQualifierPattern = regexp.MustCompile(`\b(year|month|day|hour|minute|second)\s*\(\s*[0-9]+\s*\)`)

func parseFlexibleInterval(raw string) (value, error) {
	trimmed := strings.TrimSpace(stripIntervalFieldQualifiers(strings.ToLower(raw)))
	for _, suffix := range []string{
		"year to month",
		"day to second",
		"day to minute",
		"day to hour",
		"hour to second",
		"hour to minute",
		"minute to second",
		"year",
		"month",
		"day",
		"hour",
		"minute",
		"second",
	} {
		marker := " " + suffix
		if strings.HasSuffix(trimmed, marker) {
			quantity := strings.TrimSpace(strings.TrimSuffix(trimmed, marker))
			return parseIntervalLiteral(quantity, strings.Fields(suffix))
		}
	}
	return nullValue(), fmt.Errorf("cannot cast %q as interval", raw)
}

func stripIntervalFieldQualifiers(raw string) string {
	return intervalFieldQualifierPattern.ReplaceAllString(raw, "$1")
}

func parseIntervalLiteral(quantity string, fields []string) (value, error) {
	field := strings.ToLower(strings.Join(fields, " "))
	switch field {
	case "year":
		years, err := parseIntegerQuantity(quantity)
		return intervalValue(years*12, 0), err
	case "month":
		months, err := parseIntegerQuantity(quantity)
		return intervalValue(months, 0), err
	case "year to month":
		sign, unsigned := splitIntervalSign(quantity)
		parts := strings.Split(unsigned, "-")
		if len(parts) != 2 {
			return nullValue(), fmt.Errorf("interval %q does not match year to month", quantity)
		}
		years, err := parseIntegerQuantity(parts[0])
		if err != nil {
			return nullValue(), err
		}
		months, err := parseIntegerQuantity(parts[1])
		if err != nil {
			return nullValue(), err
		}
		return intervalValue(int64(sign)*(years*12+months), 0), nil
	case "day":
		days, err := parseFloatQuantity(quantity)
		return intervalValue(0, days*86400), err
	case "hour":
		hours, err := parseFloatQuantity(quantity)
		return intervalValue(0, hours*3600), err
	case "minute":
		minutes, err := parseFloatQuantity(quantity)
		return intervalValue(0, minutes*60), err
	case "second":
		seconds, err := parseFloatQuantity(quantity)
		return intervalValue(0, seconds), err
	case "day to hour", "day to minute", "day to second":
		return parseDayTimeInterval(quantity, field)
	case "hour to minute", "hour to second":
		return parseClockInterval(quantity, field, 0)
	case "minute to second":
		return parseMinuteSecondInterval(quantity)
	default:
		return nullValue(), fmt.Errorf("unsupported interval field %q", field)
	}
}

func parseIntegerQuantity(raw string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid interval quantity %q", raw)
	}
	return value, nil
}

func parseFloatQuantity(raw string) (float64, error) {
	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid interval quantity %q", raw)
	}
	return value, nil
}

func splitIntervalSign(raw string) (int, string) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "-") {
		return -1, strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
	}
	if strings.HasPrefix(trimmed, "+") {
		return 1, strings.TrimSpace(strings.TrimPrefix(trimmed, "+"))
	}
	return 1, trimmed
}

func parseDayTimeInterval(quantity string, field string) (value, error) {
	sign, unsigned := splitIntervalSign(quantity)
	parts := strings.Fields(unsigned)
	if len(parts) == 0 || len(parts) > 2 {
		return nullValue(), fmt.Errorf("interval %q does not match %s", quantity, field)
	}
	days, err := parseFloatQuantity(parts[0])
	if err != nil {
		return nullValue(), err
	}
	seconds := days * 86400
	if len(parts) == 1 {
		return intervalValue(0, float64(sign)*seconds), nil
	}
	clock, err := parseDayClockInterval(parts[1], field)
	if err != nil {
		return nullValue(), err
	}
	return intervalValue(0, float64(sign)*(seconds+clock.intervalSeconds)), nil
}

func parseDayClockInterval(quantity string, field string) (value, error) {
	parts := strings.Split(strings.TrimSpace(quantity), ":")
	switch field {
	case "day to hour":
		if len(parts) != 1 {
			return nullValue(), fmt.Errorf("interval %q does not match %s", quantity, field)
		}
		hours, err := parseFloatQuantity(parts[0])
		if err != nil {
			return nullValue(), err
		}
		return intervalValue(0, hours*3600), nil
	case "day to minute":
		if len(parts) != 2 {
			return nullValue(), fmt.Errorf("interval %q does not match %s", quantity, field)
		}
	case "day to second":
		if len(parts) != 3 {
			return nullValue(), fmt.Errorf("interval %q does not match %s", quantity, field)
		}
	default:
		return nullValue(), fmt.Errorf("unsupported interval field %q", field)
	}
	hours, err := parseFloatQuantity(parts[0])
	if err != nil {
		return nullValue(), err
	}
	minutes, err := parseFloatQuantity(parts[1])
	if err != nil {
		return nullValue(), err
	}
	seconds := hours*3600 + minutes*60
	if len(parts) == 3 {
		extra, err := parseFloatQuantity(parts[2])
		if err != nil {
			return nullValue(), err
		}
		seconds += extra
	}
	return intervalValue(0, seconds), nil
}

func parseClockInterval(quantity string, field string, baseHours float64) (value, error) {
	sign, unsigned := splitIntervalSign(quantity)
	parts := strings.Split(unsigned, ":")
	var seconds float64
	switch field {
	case "hour to minute":
		if len(parts) != 2 {
			return nullValue(), fmt.Errorf("interval %q does not match %s", quantity, field)
		}
		hours, err := parseFloatQuantity(parts[0])
		if err != nil {
			return nullValue(), err
		}
		minutes, err := parseFloatQuantity(parts[1])
		if err != nil {
			return nullValue(), err
		}
		seconds = (baseHours+hours)*3600 + minutes*60
	case "hour to second":
		if len(parts) != 3 {
			return nullValue(), fmt.Errorf("interval %q does not match %s", quantity, field)
		}
		hours, err := parseFloatQuantity(parts[0])
		if err != nil {
			return nullValue(), err
		}
		minutes, err := parseFloatQuantity(parts[1])
		if err != nil {
			return nullValue(), err
		}
		seconds = (baseHours+hours)*3600 + minutes*60
		extra, err := parseFloatQuantity(parts[2])
		if err != nil {
			return nullValue(), err
		}
		seconds += extra
	case "hour":
		hours, err := parseFloatQuantity(unsigned)
		if err != nil {
			return nullValue(), err
		}
		seconds = (baseHours + hours) * 3600
	case "minute":
		minutes, err := parseFloatQuantity(unsigned)
		if err != nil {
			return nullValue(), err
		}
		seconds = baseHours*3600 + minutes*60
	case "second":
		extra, err := parseFloatQuantity(unsigned)
		if err != nil {
			return nullValue(), err
		}
		seconds = baseHours*3600 + extra
	default:
		return nullValue(), fmt.Errorf("unsupported interval clock field %q", field)
	}
	return intervalValue(0, float64(sign)*seconds), nil
}

func parseMinuteSecondInterval(quantity string) (value, error) {
	sign, unsigned := splitIntervalSign(quantity)
	parts := strings.Split(unsigned, ":")
	if len(parts) != 2 {
		return nullValue(), fmt.Errorf("interval %q does not match minute to second", quantity)
	}
	minutes, err := parseFloatQuantity(parts[0])
	if err != nil {
		return nullValue(), err
	}
	seconds, err := parseFloatQuantity(parts[1])
	if err != nil {
		return nullValue(), err
	}
	return intervalValue(0, float64(sign)*(minutes*60+seconds)), nil
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

func parseNumberLiteral(raw string) (value, error) {
	if !strings.ContainsAny(raw, ".eE") {
		num, ok := parseBigIntText(raw)
		if ok {
			return bigIntValue(num), nil
		}
	}
	num, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nullValue(), fmt.Errorf("invalid number %q", raw)
	}
	return numberValue(num), nil
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
		value, err := parseNumberLiteral(current.text)
		if err != nil {
			return nil, err
		}
		return literalExpr{value: value}, nil
	case tokenIdentifier:
		lower := strings.ToLower(current.text)
		if lower == "null" {
			return literalExpr{value: nullValue()}, nil
		}
		if lower == "true" || lower == "false" {
			return literalExpr{value: boolValue(lower == "true")}, nil
		}
		if lower == "date" || lower == "time" || lower == "timestamp" {
			if parser.peek().kind != tokenString {
				return nil, fmt.Errorf("%s literal expects a quoted value", lower)
			}
			raw := parser.advance().text
			normalized, err := normalizeTemporalLiteral(lower, raw)
			if err != nil {
				return nil, err
			}
			switch lower {
			case "date":
				return literalExpr{value: dateValue(normalized)}, nil
			case "time":
				return literalExpr{value: timeValue(normalized)}, nil
			case "timestamp":
				return literalExpr{value: timestampValue(normalized)}, nil
			}
		}
		if lower == "interval" {
			if parser.peek().kind != tokenString {
				return nil, errors.New("interval literal expects a quoted value")
			}
			quantity := parser.advance().text
			fields, err := parser.parseIntervalFields()
			if err != nil {
				return nil, err
			}
			interval, err := parseIntervalLiteral(quantity, fields)
			if err != nil {
				return nil, err
			}
			return literalExpr{value: interval}, nil
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

func (parser *parser) parseIntervalFields() ([]string, error) {
	if parser.peek().kind != tokenIdentifier {
		return nil, errors.New("interval literal expects a field")
	}
	first := strings.ToLower(parser.advance().text)
	if !isIntervalField(first) {
		return nil, fmt.Errorf("unsupported interval field %q", first)
	}
	if err := parser.parseOptionalIntervalFieldPrecision(first); err != nil {
		return nil, err
	}
	fields := []string{first}
	if parser.matchKeyword("to") {
		if parser.peek().kind != tokenIdentifier {
			return nil, errors.New("interval literal expects a field after to")
		}
		second := strings.ToLower(parser.advance().text)
		if !isIntervalField(second) {
			return nil, fmt.Errorf("unsupported interval field %q", second)
		}
		if err := parser.parseOptionalIntervalFieldPrecision(second); err != nil {
			return nil, err
		}
		fields = append(fields, "to", second)
	}
	return fields, nil
}

func (parser *parser) parseOptionalIntervalFieldPrecision(field string) error {
	if !parser.match(tokenLeftParen) {
		return nil
	}
	if parser.peek().kind != tokenNumber {
		return fmt.Errorf("interval %s precision expects an integer", field)
	}
	raw := parser.advance().text
	if !isUnsignedIntegerToken(raw) {
		return fmt.Errorf("interval %s precision expects an integer", field)
	}
	if !parser.match(tokenRightParen) {
		return fmt.Errorf("interval %s precision expects closing parenthesis", field)
	}
	return nil
}

func isUnsignedIntegerToken(raw string) bool {
	if raw == "" {
		return false
	}
	for _, r := range raw {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isIntervalField(field string) bool {
	switch field {
	case "year", "month", "day", "hour", "minute", "second":
		return true
	default:
		return false
	}
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
