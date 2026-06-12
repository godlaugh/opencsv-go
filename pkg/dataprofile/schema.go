package dataprofile

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// UnmarshalJSON accepts opencsv schemas plus common CSVLint/JSON Table Schema shapes.
func (schema *Schema) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*schema = Schema{}
		return nil
	}
	normalized, err := decodeSchemaJSON(data)
	if err != nil {
		return err
	}
	*schema = normalized
	return nil
}

func decodeSchemaJSON(data []byte) (Schema, error) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return Schema{}, fmt.Errorf("schema must be an object: %w", err)
	}
	if root == nil {
		return Schema{}, nil
	}
	return schemaFromObject(root)
}

func schemaFromObject(root map[string]json.RawMessage) (Schema, error) {
	source, err := schemaSourceObject(root)
	if err != nil {
		return Schema{}, err
	}
	columnRaw, err := schemaSourceColumns(source)
	if err != nil {
		return Schema{}, err
	}
	if len(columnRaw) == 0 {
		return Schema{}, nil
	}
	requiredNames := schemaRawNameSet(source["required"])
	primaryKeyNames := schemaRawNameSet(source["primaryKey"])
	if len(primaryKeyNames) == 0 {
		primaryKeyNames = schemaRawNameSet(source["primaryKeys"])
	}
	columns := make([]ColumnRule, 0, len(columnRaw))
	for _, raw := range columnRaw {
		rule, ok, err := normalizeSchemaColumnJSON(raw, requiredNames, primaryKeyNames)
		if err != nil {
			return Schema{}, err
		}
		if ok {
			columns = append(columns, rule)
		}
	}
	return Schema{Columns: columns}, nil
}

func schemaSourceObject(root map[string]json.RawMessage) (map[string]json.RawMessage, error) {
	if schemaRawIsArray(root["columns"]) || schemaRawIsArray(root["fields"]) {
		return root, nil
	}
	for _, key := range []string{"schema", "tableSchema"} {
		child, ok, err := schemaRawObject(root[key])
		if err != nil {
			return nil, fmt.Errorf("%s must be an object: %w", key, err)
		}
		if ok {
			return schemaSourceObject(child)
		}
	}
	for _, collectionKey := range []string{"tables", "resources"} {
		raw, ok := root[collectionKey]
		if !ok || len(raw) == 0 {
			continue
		}
		var tables []map[string]json.RawMessage
		if err := json.Unmarshal(raw, &tables); err != nil {
			return nil, fmt.Errorf("%s must be an array: %w", collectionKey, err)
		}
		for _, table := range tables {
			for _, key := range []string{"schema", "tableSchema"} {
				child, ok, err := schemaRawObject(table[key])
				if err != nil {
					return nil, fmt.Errorf("%s.%s must be an object: %w", collectionKey, key, err)
				}
				if ok {
					return schemaSourceObject(child)
				}
			}
		}
	}
	return root, nil
}

func schemaSourceColumns(source map[string]json.RawMessage) ([]json.RawMessage, error) {
	for _, key := range []string{"columns", "fields"} {
		raw, ok := source[key]
		if !ok || len(raw) == 0 {
			continue
		}
		var columns []json.RawMessage
		if err := json.Unmarshal(raw, &columns); err != nil {
			return nil, fmt.Errorf("%s must be an array: %w", key, err)
		}
		return columns, nil
	}
	raw, ok := source["properties"]
	if !ok || len(raw) == 0 || strings.TrimSpace(string(raw)) == "null" {
		return nil, nil
	}
	var properties map[string]json.RawMessage
	if err := json.Unmarshal(raw, &properties); err != nil {
		return nil, fmt.Errorf("properties must be an object: %w", err)
	}
	names := make([]string, 0, len(properties))
	for name := range properties {
		names = append(names, name)
	}
	sort.Strings(names)
	columns := make([]json.RawMessage, 0, len(names))
	for _, name := range names {
		column, ok, err := schemaPropertyColumnRaw(name, properties[name])
		if err != nil {
			return nil, err
		}
		if ok {
			columns = append(columns, column)
		}
	}
	return columns, nil
}

func schemaPropertyColumnRaw(name string, raw json.RawMessage) (json.RawMessage, bool, error) {
	property, ok, err := schemaRawObject(raw)
	if err != nil {
		return nil, false, fmt.Errorf("properties.%s must be an object: %w", name, err)
	}
	if !ok {
		return nil, false, nil
	}
	if firstSchemaString(property, "name", "title", "label") == "" {
		nameBytes, err := json.Marshal(name)
		if err != nil {
			return nil, false, err
		}
		property["name"] = nameBytes
	}
	column, err := json.Marshal(property)
	if err != nil {
		return nil, false, err
	}
	return column, true, nil
}

func normalizeSchemaColumnJSON(raw json.RawMessage, requiredNames map[string]struct{}, primaryKeyNames map[string]struct{}) (ColumnRule, bool, error) {
	var column map[string]json.RawMessage
	if err := json.Unmarshal(raw, &column); err != nil {
		return ColumnRule{}, false, fmt.Errorf("schema column must be an object: %w", err)
	}
	constraints, ok, err := schemaRawObject(column["constraints"])
	if err != nil {
		return ColumnRule{}, false, fmt.Errorf("constraints must be an object: %w", err)
	}
	if !ok {
		constraints = map[string]json.RawMessage{}
	}
	datatypeRaw := firstSchemaRaw(column, "datatype", "dataType", "type")
	datatypeObject, _, err := schemaRawOptionalObject(datatypeRaw)
	if err != nil {
		return ColumnRule{}, false, fmt.Errorf("datatype must be a string or object: %w", err)
	}
	name := firstSchemaString(column, "name", "title", "label")
	rule := ColumnRule{Name: name}
	if col := schemaPositiveInteger(firstSchemaRaw(column, "column", "index", "number")); col > 0 {
		rule.Column = col
	}
	if typ := normalizeSchemaType(firstSchemaString(datatypeObject, "base", "type")); typ != "" {
		rule.Type = typ
	} else if typ := normalizeSchemaType(schemaRawString(datatypeRaw)); typ != "" {
		rule.Type = typ
	} else if typ := normalizeSchemaType(firstSchemaString(constraints, "type")); typ != "" {
		rule.Type = typ
	}
	if schemaRawBool(firstSchemaRaw(column, "required")) || schemaRawBool(firstSchemaRaw(constraints, "required")) || schemaNameSetHas(requiredNames, name) || schemaNameSetHas(primaryKeyNames, name) {
		rule.Required = true
	}
	if schemaRawBool(firstSchemaRaw(column, "unique")) || schemaRawBool(firstSchemaRaw(constraints, "unique")) || schemaNameSetHas(primaryKeyNames, name) {
		rule.Unique = true
	}
	rule.Enum = firstSchemaArrayValues(column, constraints, []string{"enum", "oneOf", "values", "in"})
	rule.Pattern = firstSchemaString(column, "pattern")
	if rule.Pattern == "" {
		rule.Pattern = firstSchemaString(constraints, "pattern", "regex")
	}
	if min, ok := firstSchemaNumber([]map[string]json.RawMessage{column, constraints, datatypeObject}, "min", "minimum"); ok {
		rule.Min = &min
	}
	if max, ok := firstSchemaNumber([]map[string]json.RawMessage{column, constraints, datatypeObject}, "max", "maximum"); ok {
		rule.Max = &max
	}
	if minLength, ok := firstSchemaInteger([]map[string]json.RawMessage{column, constraints, datatypeObject}, "minLength", "minimumLength", "min_length", "minimum_length"); ok {
		rule.MinLength = &minLength
	}
	if maxLength, ok := firstSchemaInteger([]map[string]json.RawMessage{column, constraints, datatypeObject}, "maxLength", "maximumLength", "max_length", "maximum_length"); ok {
		rule.MaxLength = &maxLength
	}
	rule.Default = firstSchemaString(column, "default", "defaultValue", "default_value")
	if rule.Default == "" {
		rule.Default = firstSchemaString(constraints, "default", "defaultValue", "default_value")
	}
	return rule, ruleHasChecks(rule), nil
}

func schemaRawIsArray(raw json.RawMessage) bool {
	var values []json.RawMessage
	return len(raw) > 0 && json.Unmarshal(raw, &values) == nil
}

func schemaRawObject(raw json.RawMessage) (map[string]json.RawMessage, bool, error) {
	if len(raw) == 0 || strings.TrimSpace(string(raw)) == "null" {
		return nil, false, nil
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		return nil, false, err
	}
	if object == nil {
		return nil, false, nil
	}
	return object, true, nil
}

func schemaRawOptionalObject(raw json.RawMessage) (map[string]json.RawMessage, bool, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" || trimmed == "null" || !strings.HasPrefix(trimmed, "{") {
		return nil, false, nil
	}
	return schemaRawObject(raw)
}

func schemaRawNameSet(raw json.RawMessage) map[string]struct{} {
	names := map[string]struct{}{}
	if value := schemaRawString(raw); value != "" {
		names[value] = struct{}{}
		return names
	}
	if len(raw) == 0 {
		return names
	}
	var values []json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return names
	}
	for _, item := range values {
		if value := schemaRawString(item); value != "" {
			names[value] = struct{}{}
		}
	}
	return names
}

func schemaNameSetHas(names map[string]struct{}, value string) bool {
	if value == "" {
		return false
	}
	_, ok := names[value]
	return ok
}

func firstSchemaRaw(source map[string]json.RawMessage, keys ...string) json.RawMessage {
	if source == nil {
		return nil
	}
	for _, key := range keys {
		if raw, ok := source[key]; ok {
			return raw
		}
	}
	return nil
}

func firstSchemaString(source map[string]json.RawMessage, keys ...string) string {
	return schemaRawString(firstSchemaRaw(source, keys...))
}

func schemaRawString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		return strings.TrimSpace(number.String())
	}
	var flag bool
	if err := json.Unmarshal(raw, &flag); err == nil {
		if flag {
			return "true"
		}
		return "false"
	}
	return ""
}

func schemaRawBool(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var flag bool
	if err := json.Unmarshal(raw, &flag); err == nil {
		return flag
	}
	value := strings.ToLower(schemaRawString(raw))
	return value == "true" || value == "1" || value == "yes"
}

func schemaPositiveInteger(raw json.RawMessage) int {
	value := schemaRawString(raw)
	if value == "" {
		return 0
	}
	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return 0
	}
	return number
}

func firstSchemaArrayValues(column map[string]json.RawMessage, constraints map[string]json.RawMessage, keys []string) []string {
	for _, source := range []map[string]json.RawMessage{column, constraints} {
		for _, key := range keys {
			values := schemaArrayValues(firstSchemaRaw(source, key))
			if len(values) > 0 {
				return values
			}
		}
	}
	return nil
}

func schemaArrayValues(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	if text := schemaRawString(raw); text != "" {
		parts := strings.Split(text, ",")
		values := make([]string, 0, len(parts))
		for _, part := range parts {
			if value := strings.TrimSpace(part); value != "" {
				values = append(values, value)
			}
		}
		return values
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil
	}
	values := make([]string, 0, len(items))
	for _, item := range items {
		if value := schemaRawString(item); value != "" {
			values = append(values, value)
		}
	}
	return values
}

func firstSchemaNumber(sources []map[string]json.RawMessage, keys ...string) (float64, bool) {
	for _, source := range sources {
		if source == nil {
			continue
		}
		for _, key := range keys {
			if value, ok := schemaRawNumber(firstSchemaRaw(source, key)); ok {
				return value, true
			}
		}
	}
	return 0, false
}

func firstSchemaInteger(sources []map[string]json.RawMessage, keys ...string) (int, bool) {
	for _, source := range sources {
		if source == nil {
			continue
		}
		for _, key := range keys {
			if value, ok := schemaRawNonNegativeInteger(firstSchemaRaw(source, key)); ok {
				return value, true
			}
		}
	}
	return 0, false
}

func schemaRawNonNegativeInteger(raw json.RawMessage) (int, bool) {
	value := schemaRawString(raw)
	if value == "" {
		return 0, false
	}
	number, err := strconv.Atoi(value)
	if err != nil || number < 0 {
		return 0, false
	}
	return number, true
}

func schemaRawNumber(raw json.RawMessage) (float64, bool) {
	if len(raw) == 0 || strings.TrimSpace(string(raw)) == "null" {
		return 0, false
	}
	var number float64
	if err := json.Unmarshal(raw, &number); err == nil {
		return number, true
	}
	if value := schemaRawString(raw); value != "" {
		number, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return number, true
		}
	}
	return 0, false
}

func normalizeSchemaType(value string) string {
	normalized := normalizeSchemaTypeName(value)
	switch normalized {
	case "":
		return ""
	case "integer", "int", "long", "short", "byte", "nonnegativeinteger", "nonpositiveinteger", "positiveinteger", "negativeinteger", "unsignedlong", "unsignedint", "unsignedshort", "unsignedbyte":
		return "integer"
	case "number", "numeric", "float", "double":
		return "number"
	case "decimal":
		return "decimal"
	case "boolean", "bool":
		return "boolean"
	case "date":
		return "date"
	case "datetime", "date-time", "timestamp":
		return "datetime"
	case "time":
		return "time"
	case "year", "gyear":
		return "year"
	case "yearmonth", "year-month", "gyearmonth":
		return "yearmonth"
	case "string", "text", "normalizedstring", "token", "uri", "url", "anyuri", "email":
		return "string"
	default:
		return normalized
	}
}

func normalizeSchemaTypeName(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return ""
	}
	if index := strings.LastIndex(normalized, "#"); index >= 0 {
		normalized = normalized[index+1:]
	} else if index := strings.LastIndex(normalized, "/"); index >= 0 {
		normalized = normalized[index+1:]
	}
	normalized = strings.TrimPrefix(normalized, "xsd:")
	normalized = strings.TrimPrefix(normalized, "xmlschema:")
	return normalized
}
