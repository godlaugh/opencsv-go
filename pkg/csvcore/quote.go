package csvcore

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	QuoteModeMinimal = "minimal"
	QuoteModeAlways  = "always"
	QuoteModeNever   = "never"
)

func NormalizeQuoteMode(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", QuoteModeMinimal:
		return QuoteModeMinimal, nil
	case QuoteModeAlways:
		return QuoteModeAlways, nil
	case QuoteModeNever, "none":
		return QuoteModeNever, nil
	default:
		return "", fmt.Errorf("quote mode must be minimal, always, or never")
	}
}

func NormalizeQuoteChar(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "double", "quote":
		return `"`, nil
	case "single":
		return "'", nil
	case "none":
		return "", nil
	}
	runes := []rune(value)
	if len(runes) != 1 {
		return "", fmt.Errorf("quote character must be double, single, none, or one character")
	}
	return value, nil
}

func DialectQuoteMode(dialect Dialect) string {
	mode, err := NormalizeQuoteMode(dialect.QuoteMode)
	if err != nil {
		return QuoteModeMinimal
	}
	return mode
}

func DialectQuoteChar(dialect Dialect) string {
	quoteChar, err := NormalizeQuoteChar(dialect.QuoteChar)
	if err != nil {
		return `"`
	}
	return quoteChar
}

func ShouldQuote(value string, dialect Dialect) bool {
	switch DialectQuoteMode(dialect) {
	case QuoteModeAlways:
		return true
	case QuoteModeNever:
		return false
	default:
		return RequiresQuote(value, dialect)
	}
}

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
	if strings.ContainsAny(value, "\r\n") {
		return true
	}
	if quote := DialectQuoteChar(dialect); quote != "" && strings.Contains(value, quote) {
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
