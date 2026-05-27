package clipboard

import (
	"bytes"
	"encoding/csv"
	"strings"
)

func ToTSV(matrix [][]string) string {
	return joinMatrix(matrix, "\t")
}

func ToRaw(matrix [][]string) string {
	return joinMatrix(matrix, "\t")
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

func ParseTSV(input string) [][]string {
	input = strings.TrimSuffix(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	if input == "" {
		return nil
	}
	lines := strings.Split(strings.ReplaceAll(input, "\r", "\n"), "\n")
	matrix := make([][]string, len(lines))
	for i, line := range lines {
		matrix[i] = strings.Split(line, "\t")
	}
	return matrix
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
