package clipboard

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"html"
	"strings"
)

func ToTSV(matrix [][]string) string {
	return joinMatrix(matrix, "\t")
}

func ToRaw(matrix [][]string) string {
	return joinMatrix(matrix, "\t")
}

func WithHeaders(matrix [][]string, headers []string) [][]string {
	result := make([][]string, 0, len(matrix)+1)
	result = append(result, append([]string(nil), headers...))
	for _, row := range matrix {
		result = append(result, append([]string(nil), row...))
	}
	return result
}

func ToCSV(matrix [][]string) (string, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	for _, row := range matrix {
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return strings.TrimSuffix(buffer.String(), "\n"), nil
}

func ToMarkdown(matrix [][]string) string {
	if len(matrix) == 0 {
		return ""
	}
	var builder strings.Builder
	writeMarkdownRow(&builder, matrix[0])
	builder.WriteByte('\n')
	builder.WriteByte('|')
	for range matrix[0] {
		builder.WriteString(" --- |")
	}
	for _, row := range matrix[1:] {
		builder.WriteByte('\n')
		writeMarkdownRow(&builder, row)
	}
	return builder.String()
}

func ToJSON(matrix [][]string) (string, error) {
	payload, err := json.MarshalIndent(matrix, "", "  ")
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func ToHTML(matrix [][]string) string {
	if len(matrix) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString("<table>\n  <tbody>")
	for _, row := range matrix {
		builder.WriteString("\n    <tr>")
		for _, value := range row {
			builder.WriteString("<td>")
			builder.WriteString(html.EscapeString(value))
			builder.WriteString("</td>")
		}
		builder.WriteString("</tr>")
	}
	builder.WriteString("\n  </tbody>\n</table>")
	return builder.String()
}

func ToSQLValues(matrix [][]string) string {
	if len(matrix) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString("VALUES")
	for rowIndex, row := range matrix {
		if rowIndex == 0 {
			builder.WriteByte('\n')
		} else {
			builder.WriteString(",\n")
		}
		builder.WriteString("  (")
		for colIndex, value := range row {
			if colIndex > 0 {
				builder.WriteString(", ")
			}
			builder.WriteByte('\'')
			builder.WriteString(strings.ReplaceAll(value, "'", "''"))
			builder.WriteByte('\'')
		}
		builder.WriteByte(')')
	}
	builder.WriteByte(';')
	return builder.String()
}

func ToLaTeX(matrix [][]string) string {
	if len(matrix) == 0 {
		return ""
	}
	cols := maxColumns(matrix)
	var builder strings.Builder
	builder.WriteString("\\begin{tabular}{")
	builder.WriteString(strings.Repeat("l", cols))
	builder.WriteString("}\n")
	for index, row := range matrix {
		if index > 0 {
			builder.WriteString(" \\\\\n")
		}
		cells := make([]string, cols)
		for col := 0; col < cols; col++ {
			if col < len(row) {
				cells[col] = escapeLaTeXCell(row[col])
			}
		}
		builder.WriteString(strings.Join(cells, " & "))
	}
	builder.WriteString("\n\\end{tabular}")
	return builder.String()
}

func ParseTSV(input string) [][]string {
	if input == "" {
		return nil
	}
	if strings.Trim(input, "\r\n") == "" {
		return nil
	}
	var matrix [][]string
	var row []string
	var field strings.Builder
	quoted := false
	fieldStart := true
	rowBreak := false
	for index := 0; index < len(input); index++ {
		ch := input[index]
		if quoted {
			if ch == '"' {
				if index+1 < len(input) && input[index+1] == '"' {
					field.WriteByte('"')
					index++
				} else {
					quoted = false
				}
				rowBreak = false
				continue
			}
			if ch == '\r' {
				if index+1 < len(input) && input[index+1] == '\n' {
					index++
				}
				field.WriteByte('\n')
			} else {
				field.WriteByte(ch)
			}
			rowBreak = false
			continue
		}
		switch ch {
		case '"':
			if fieldStart {
				quoted = true
				fieldStart = false
			} else {
				field.WriteByte(ch)
			}
			rowBreak = false
		case '\t':
			row = append(row, field.String())
			field.Reset()
			fieldStart = true
			rowBreak = false
		case '\r', '\n':
			if ch == '\r' && index+1 < len(input) && input[index+1] == '\n' {
				index++
			}
			row = append(row, field.String())
			field.Reset()
			matrix = append(matrix, row)
			row = nil
			fieldStart = true
			rowBreak = true
		default:
			field.WriteByte(ch)
			fieldStart = false
			rowBreak = false
		}
	}
	if !rowBreak || len(row) > 0 || field.Len() > 0 {
		row = append(row, field.String())
		matrix = append(matrix, row)
	}
	return matrix
}

func maxColumns(matrix [][]string) int {
	cols := 0
	for _, row := range matrix {
		if len(row) > cols {
			cols = len(row)
		}
	}
	if cols == 0 {
		return 1
	}
	return cols
}

func joinMatrix(matrix [][]string, separator string) string {
	lines := make([]string, len(matrix))
	for i, row := range matrix {
		lines[i] = strings.Join(row, separator)
	}
	return strings.Join(lines, "\n")
}

func writeMarkdownRow(builder *strings.Builder, row []string) {
	builder.WriteByte('|')
	for _, value := range row {
		builder.WriteByte(' ')
		builder.WriteString(escapeMarkdownCell(value))
		builder.WriteString(" |")
	}
}

func escapeMarkdownCell(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "|", `\|`)
	value = strings.ReplaceAll(value, "\n", "<br>")
	return value
}

func escapeLaTeXCell(value string) string {
	replacer := strings.NewReplacer(
		`\`, `\textbackslash{}`,
		"&", `\&`,
		"%", `\%`,
		"$", `\$`,
		"#", `\#`,
		"_", `\_`,
		"{", `\{`,
		"}", `\}`,
		"~", `\textasciitilde{}`,
		"^", `\textasciicircum{}`,
		"|", `\textbar{}`,
		"\n", `\\ `,
	)
	return replacer.Replace(value)
}
