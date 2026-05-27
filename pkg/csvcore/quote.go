package csvcore

import (
	"strings"
	"unicode"
)

func RequiresQuote(value string, dialect Dialect) bool {
	if value == "" {
		return false
	}
	delimiter := dialect.Delimiter
	if delimiter == 0 {
		delimiter = ','
	}
	if strings.ContainsRune(value, delimiter) {
		return true
	}
	if strings.ContainsAny(value, "\"\r\n") {
		return true
	}
	first, _ := firstRune(value)
	return unicode.IsSpace(first)
}

func firstRune(value string) (rune, bool) {
	for _, r := range value {
		return r, true
	}
	return 0, false
}
