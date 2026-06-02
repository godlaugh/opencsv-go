package csvcore

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func NormalizeDelimiter(value string, name string) (rune, string, error) {
	raw := strings.TrimSpace(value)
	label := strings.TrimSpace(name)
	if raw == "" && label != "" {
		raw = label
	}
	if value == "\t" {
		return '\t', "tab", nil
	}
	if value == " " {
		return ' ', "space", nil
	}
	switch strings.ToLower(raw) {
	case "comma", ",":
		return ',', "comma", nil
	case "tab", "tsv", "\\t":
		return '\t', "tab", nil
	case "semicolon", ";":
		return ';', "semicolon", nil
	case "pipe", "|":
		return '|', "pipe", nil
	case "space":
		return ' ', "space", nil
	}
	r, size := utf8.DecodeRuneInString(value)
	if r == utf8.RuneError || size != len(value) {
		return 0, "", fmt.Errorf("delimiter must be a known alias or single character")
	}
	if r == '\r' || r == '\n' || r == '"' || r == '\'' {
		return 0, "", fmt.Errorf("delimiter must be a non-quote, non-line-break character")
	}
	if label == "" {
		label = value
	}
	return r, label, nil
}
