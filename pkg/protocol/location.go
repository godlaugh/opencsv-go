package protocol

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Location struct {
	StartRow int `json:"startRow"`
	StartCol int `json:"startCol,omitempty"`
	EndRow   int `json:"endRow,omitempty"`
	EndCol   int `json:"endCol,omitempty"`
}

func (location Location) IsZero() bool {
	return location.StartRow == 0 && location.StartCol == 0 && location.EndRow == 0 && location.EndCol == 0
}

type FormatParams struct {
	Delimiter    string     `json:"delimiter,omitempty"`
	Quote        string     `json:"quote,omitempty"`
	QuoteMode    string     `json:"quoteMode,omitempty"`
	Encoding     string     `json:"encoding,omitempty"`
	LineEnding   string     `json:"lineEnding,omitempty"`
	HeaderRows   *int       `json:"headerRows,omitempty"`
	HeaderCols   *int       `json:"headerColumns,omitempty"`
	Header       *bool      `json:"header,omitempty"`
	NoHeader     *bool      `json:"noHeader,omitempty"`
	RawParameter url.Values `json:"-"`
}

type FileTarget struct {
	Path     string       `json:"path"`
	Location Location     `json:"location,omitempty"`
	Format   FormatParams `json:"format,omitempty"`
}

type LocationTarget struct {
	Path     string
	Location Location
}

var (
	cellRangePattern = regexp.MustCompile(`^(.*):([0-9]+):([0-9]+)-([0-9]+):([0-9]+)$`)
	rowRangePattern  = regexp.MustCompile(`^(.*):([0-9]+)-([0-9]+)$`)
	cellPattern      = regexp.MustCompile(`^(.*):([0-9]+):([0-9]+)$`)
	rowPattern       = regexp.MustCompile(`^(.*):([0-9]+)$`)
)

func ParseGoto(value string) (LocationTarget, error) {
	path, location, err := splitLocation(value)
	if err != nil {
		return LocationTarget{}, err
	}
	if strings.TrimSpace(path) == "" {
		return LocationTarget{}, errors.New("goto path is required")
	}
	return LocationTarget{Path: path, Location: location}, nil
}

func ParseSmoothCSVURL(raw string) (FileTarget, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return FileTarget{}, err
	}
	if !isCSVURLScheme(parsed.Scheme) {
		return FileTarget{}, fmt.Errorf("unsupported URL scheme %q", parsed.Scheme)
	}
	if !strings.EqualFold(parsed.Host, "file") {
		return FileTarget{}, fmt.Errorf("unsupported CSV URL host %q", parsed.Host)
	}
	path, err := url.PathUnescape(parsed.EscapedPath())
	if err != nil {
		return FileTarget{}, err
	}
	path = normalizeURLPath(path)
	path, location, err := splitLocation(path)
	if err != nil {
		return FileTarget{}, err
	}
	if path == "" {
		return FileTarget{}, errors.New("file path is required")
	}

	format, err := parseFormatParams(parsed.Query())
	if err != nil {
		return FileTarget{}, err
	}
	return FileTarget{Path: path, Location: location, Format: format}, nil
}

func IsCSVURL(raw string) bool {
	parsed, err := url.Parse(raw)
	return err == nil && isCSVURLScheme(parsed.Scheme)
}

func isCSVURLScheme(scheme string) bool {
	return strings.EqualFold(scheme, "smoothcsv") || strings.EqualFold(scheme, "opencsv")
}

func splitLocation(value string) (string, Location, error) {
	if match := cellRangePattern.FindStringSubmatch(value); match != nil {
		location, err := makeLocation(match[2], match[3], match[4], match[5])
		return match[1], location, err
	}
	if match := rowRangePattern.FindStringSubmatch(value); match != nil {
		start, err := positiveInt(match[2], "start row")
		if err != nil {
			return "", Location{}, err
		}
		end, err := positiveInt(match[3], "end row")
		if err != nil {
			return "", Location{}, err
		}
		if end < start {
			return "", Location{}, errors.New("end row must be greater than or equal to start row")
		}
		return match[1], Location{StartRow: start, EndRow: end}, nil
	}
	if match := cellPattern.FindStringSubmatch(value); match != nil {
		row, err := positiveInt(match[2], "row")
		if err != nil {
			return "", Location{}, err
		}
		col, err := positiveInt(match[3], "column")
		if err != nil {
			return "", Location{}, err
		}
		return match[1], Location{StartRow: row, StartCol: col}, nil
	}
	if match := rowPattern.FindStringSubmatch(value); match != nil && !looksLikeWindowsDriveOnly(match[1], match[2]) {
		row, err := positiveInt(match[2], "row")
		if err != nil {
			return "", Location{}, err
		}
		return match[1], Location{StartRow: row}, nil
	}
	return value, Location{}, nil
}

func makeLocation(startRow string, startCol string, endRow string, endCol string) (Location, error) {
	sr, err := positiveInt(startRow, "start row")
	if err != nil {
		return Location{}, err
	}
	sc, err := positiveInt(startCol, "start column")
	if err != nil {
		return Location{}, err
	}
	er, err := positiveInt(endRow, "end row")
	if err != nil {
		return Location{}, err
	}
	ec, err := positiveInt(endCol, "end column")
	if err != nil {
		return Location{}, err
	}
	if er < sr || (er == sr && ec < sc) {
		return Location{}, errors.New("range end must be after range start")
	}
	return Location{StartRow: sr, StartCol: sc, EndRow: er, EndCol: ec}, nil
}

func parseFormatParams(values url.Values) (FormatParams, error) {
	params := FormatParams{RawParameter: values}
	params.Delimiter = first(values, "delimiter", "d")
	params.Quote = first(values, "quote", "q")
	params.QuoteMode = first(values, "quotemode", "quote-mode", "quoteMode", "quote_mode", "qm")
	params.Encoding = first(values, "encoding", "e")
	params.LineEnding = first(values, "lineending", "line-ending", "lineEnding", "line_ending", "newline", "eol", "le")
	if raw := first(values, "headerrows", "header-rows", "headerRows"); raw != "" {
		value, err := nonNegativeInt(raw, "header rows")
		if err != nil {
			return FormatParams{}, err
		}
		params.HeaderRows = &value
	}
	if raw := first(values, "headercolumns", "header-columns", "headerColumns", "header_cols"); raw != "" {
		value, err := nonNegativeInt(raw, "header columns")
		if err != nil {
			return FormatParams{}, err
		}
		params.HeaderCols = &value
	}
	if raw := first(values, "header"); raw != "" {
		value, err := parseBool(raw, "header")
		if err != nil {
			return FormatParams{}, err
		}
		params.Header = &value
	}
	if raw := first(values, "noheader", "no-header", "noHeader"); raw != "" {
		value, err := parseBool(raw, "no-header")
		if err != nil {
			return FormatParams{}, err
		}
		params.NoHeader = &value
	}
	return params, nil
}

func normalizeURLPath(path string) string {
	if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		return path[1:]
	}
	return path
}

func looksLikeWindowsDriveOnly(path string, suffix string) bool {
	return len(path) == 1 && len(suffix) == 0 && path[0] >= 'A' && path[0] <= 'z'
}

func first(values url.Values, names ...string) string {
	for _, name := range names {
		if value := values.Get(name); value != "" {
			return value
		}
	}
	for _, name := range names {
		for key, value := range values {
			if strings.EqualFold(key, name) && len(value) > 0 && value[0] != "" {
				return value[0]
			}
		}
	}
	return ""
}

func positiveInt(raw string, name string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return value, nil
}

func nonNegativeInt(raw string, name string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", name)
	}
	return value, nil
}

func parseBool(raw string, name string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("%s must be boolean", name)
	}
}
