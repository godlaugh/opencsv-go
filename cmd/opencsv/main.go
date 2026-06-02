package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/opencsv/opencsv/pkg/appserver"
	"github.com/opencsv/opencsv/pkg/csvcore"
	"github.com/opencsv/opencsv/pkg/csvdiff"
	"github.com/opencsv/opencsv/pkg/protocol"
)

const version = "0.2.0-dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "open":
		err = openDocument(os.Args[2:])
	case "inspect":
		err = inspect(os.Args[2:])
	case "generate":
		err = generate(os.Args[2:])
	case "save":
		err = save(os.Args[2:])
	case "grep":
		err = grep(os.Args[2:])
	case "serve":
		err = serve(os.Args[2:])
	case "bench":
		err = bench(os.Args[2:])
	default:
		err = rootOpen(os.Args[1:])
	}

	if err != nil {
		var usageErr usageError
		if errors.As(err, &usageErr) {
			fmt.Fprintln(os.Stderr, "opencsv:", usageErr.Error())
			usage()
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, "opencsv:", err)
		os.Exit(1)
	}
}

type usageError string

func (err usageError) Error() string {
	return string(err)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  opencsv [options] [path ...]")
	fmt.Fprintln(os.Stderr, "  opencsv --goto file.csv:10:3")
	fmt.Fprintln(os.Stderr, "  opencsv --diff old.csv new.csv")
	fmt.Fprintln(os.Stderr, "  opencsv open [-rows 50] [-index-rows 0] <file>")
	fmt.Fprintln(os.Stderr, "  opencsv inspect [-rows 20] <file>")
	fmt.Fprintln(os.Stderr, "  opencsv generate [-rows 1000] [-cols 12] <file>")
	fmt.Fprintln(os.Stderr, "  opencsv save [-out copy.csv] [-backup=true] <file>")
	fmt.Fprintln(os.Stderr, "  opencsv grep [-i] [-regex] [-col name|index] [-limit 100] <query> <file>")
	fmt.Fprintln(os.Stderr, "  opencsv serve [-addr 127.0.0.1:43255] [-web web]")
	fmt.Fprintln(os.Stderr, "  opencsv bench [-rows 100000] [-cols 16] [-out /tmp/opencsv-bench.csv]")
}

type cliInvocation struct {
	Mode        string                `json:"mode"`
	Files       []cliOpenedFile       `json:"files"`
	Diff        *csvdiff.Result       `json:"diff,omitempty"`
	Goto        *protocol.Location    `json:"goto,omitempty"`
	Format      protocol.FormatParams `json:"format,omitempty"`
	NewWindow   bool                  `json:"newWindow,omitempty"`
	ReuseWindow bool                  `json:"reuseWindow,omitempty"`
	Wait        bool                  `json:"wait,omitempty"`
}

type cliOpenedFile struct {
	Path     string            `json:"path"`
	Source   string            `json:"source"`
	Location protocol.Location `json:"location,omitempty"`
	Document csvcore.Document  `json:"document"`
}

type rootOptions struct {
	paths       []string
	gotoTarget  string
	diff        bool
	format      protocol.FormatParams
	newWindow   bool
	reuseWindow bool
	wait        bool
	version     bool
	help        bool
}

func rootOpen(args []string) error {
	options, err := parseRootOptions(args)
	if err != nil {
		return err
	}
	if options.help {
		usage()
		return nil
	}
	if options.version {
		fmt.Println(version)
		return nil
	}
	invocation := cliInvocation{
		Mode:        "open",
		Format:      options.format,
		NewWindow:   options.newWindow,
		ReuseWindow: options.reuseWindow,
		Wait:        options.wait,
	}
	if options.diff {
		invocation.Mode = "diff"
	}

	openOptions, err := openOptionsFromFormat(options.format)
	if err != nil {
		return err
	}
	openOptions.PreviewRows = 50
	openOptions.IndexRows = 0

	var gotoTarget *protocol.LocationTarget
	if options.gotoTarget != "" {
		target, err := protocol.ParseGoto(options.gotoTarget)
		if err != nil {
			return usageError(err.Error())
		}
		gotoTarget = &target
		invocation.Goto = &target.Location
		if len(options.paths) == 0 {
			options.paths = append(options.paths, target.Path)
		}
	}
	if len(options.paths) == 0 {
		return usageError("at least one path is required")
	}
	if options.diff {
		if len(options.paths) != 2 {
			return usageError("--diff requires exactly two file paths")
		}
		for _, path := range options.paths {
			if path == "-" {
				return usageError("--diff does not support stdin")
			}
		}
	}

	for _, rawPath := range options.paths {
		target := protocol.FileTarget{Path: rawPath, Format: options.format}
		if strings.HasPrefix(strings.ToLower(rawPath), "smoothcsv://") {
			parsed, err := protocol.ParseSmoothCSVURL(rawPath)
			if err != nil {
				return usageError(err.Error())
			}
			target = parsed
			openOptions, err = openOptionsFromFormat(mergeFormat(options.format, target.Format))
			if err != nil {
				return err
			}
			openOptions.PreviewRows = 50
			openOptions.IndexRows = 0
		}
		if gotoTarget != nil && samePath(gotoTarget.Path, target.Path) {
			target.Location = gotoTarget.Location
		}

		document, source, err := openCLITarget(target.Path, openOptions)
		if err != nil {
			return err
		}
		invocation.Files = append(invocation.Files, cliOpenedFile{
			Path:     target.Path,
			Source:   source,
			Location: target.Location,
			Document: document,
		})
	}

	if options.diff {
		result, err := csvdiff.CompareFiles(options.paths[0], options.paths[1], csvdiff.Options{MaxRows: 500})
		if err != nil {
			return err
		}
		invocation.Diff = &result
	}

	return printJSON(invocation)
}

func parseRootOptions(args []string) (rootOptions, error) {
	options := rootOptions{}
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--" {
			options.paths = append(options.paths, args[index+1:]...)
			return options, nil
		}
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			options.paths = append(options.paths, arg)
			continue
		}
		name, value, hasValue := strings.Cut(arg, "=")
		takeValue := func() (string, error) {
			if hasValue {
				return value, nil
			}
			if index+1 >= len(args) {
				return "", usageError(fmt.Sprintf("%s requires a value", name))
			}
			index++
			return args[index], nil
		}
		switch name {
		case "--help", "-h":
			options.help = true
		case "--version", "-v":
			options.version = true
		case "--goto", "-g":
			value, err := takeValue()
			if err != nil {
				return rootOptions{}, err
			}
			options.gotoTarget = value
		case "--diff":
			options.diff = true
		case "--delimiter", "-d":
			value, err := takeValue()
			if err != nil {
				return rootOptions{}, err
			}
			options.format.Delimiter = value
		case "--quote", "-q":
			value, err := takeValue()
			if err != nil {
				return rootOptions{}, err
			}
			options.format.Quote = value
		case "--quote-mode", "--qm":
			value, err := takeValue()
			if err != nil {
				return rootOptions{}, err
			}
			options.format.QuoteMode = value
		case "--encoding", "-e":
			value, err := takeValue()
			if err != nil {
				return rootOptions{}, err
			}
			options.format.Encoding = value
		case "--line-ending", "--eol":
			value, err := takeValue()
			if err != nil {
				return rootOptions{}, err
			}
			options.format.LineEnding = value
		case "--header-rows":
			raw, err := takeValue()
			if err != nil {
				return rootOptions{}, err
			}
			value, err := parseNonNegative(raw, "--header-rows")
			if err != nil {
				return rootOptions{}, err
			}
			options.format.HeaderRows = &value
		case "--header-columns":
			raw, err := takeValue()
			if err != nil {
				return rootOptions{}, err
			}
			value, err := parseNonNegative(raw, "--header-columns")
			if err != nil {
				return rootOptions{}, err
			}
			options.format.HeaderCols = &value
		case "--header":
			value := true
			options.format.Header = &value
		case "--no-header":
			value := true
			options.format.NoHeader = &value
		case "--new-window", "-n":
			options.newWindow = true
		case "--reuse-window", "-r":
			options.reuseWindow = true
		case "--wait", "-w":
			options.wait = true
		default:
			return rootOptions{}, usageError(fmt.Sprintf("unknown option %s", name))
		}
	}
	return options, nil
}

func openOptionsFromFormat(format protocol.FormatParams) (csvcore.OpenOptions, error) {
	options := csvcore.OpenOptions{}
	if format.Delimiter != "" {
		delimiter, name, err := parseDelimiter(format.Delimiter)
		if err != nil {
			return csvcore.OpenOptions{}, err
		}
		options.Delimiter = &delimiter
		options.DelimiterName = name
	}
	if format.Encoding != "" {
		encoding, err := csvcore.ParseEncodingName(format.Encoding)
		if err != nil {
			return csvcore.OpenOptions{}, err
		}
		options.Encoding = &encoding
	}
	if format.LineEnding != "" {
		lineEnding, err := csvcore.NormalizeLineEndingName(format.LineEnding)
		if err != nil {
			return csvcore.OpenOptions{}, err
		}
		options.LineEnding = &lineEnding
	}
	if format.Quote != "" {
		quoteChar, err := csvcore.NormalizeQuoteChar(format.Quote)
		if err != nil {
			return csvcore.OpenOptions{}, err
		}
		options.QuoteChar = &quoteChar
	}
	if format.QuoteMode != "" {
		quoteMode, err := csvcore.NormalizeQuoteMode(format.QuoteMode)
		if err != nil {
			return csvcore.OpenOptions{}, err
		}
		options.QuoteMode = quoteMode
	}
	if format.HeaderRows != nil {
		hasHeader := *format.HeaderRows > 0
		options.HasHeader = &hasHeader
	}
	if format.HeaderCols != nil {
		options.HeaderColumns = format.HeaderCols
	}
	if format.Header != nil {
		hasHeader := *format.Header
		options.HasHeader = &hasHeader
	}
	if format.NoHeader != nil && *format.NoHeader {
		hasHeader := false
		options.HasHeader = &hasHeader
	}
	return options, nil
}

func parseDelimiter(raw string) (rune, string, error) {
	switch strings.ToLower(raw) {
	case "comma":
		return ',', "comma", nil
	case "tab", "\\t":
		return '\t', "tab", nil
	case "semicolon":
		return ';', "semicolon", nil
	case "pipe":
		return '|', "pipe", nil
	case "space":
		return ' ', "space", nil
	}
	r, size := utf8.DecodeRuneInString(raw)
	if r == utf8.RuneError || size != len(raw) {
		return 0, "", fmt.Errorf("delimiter must be a known alias or single character")
	}
	return r, raw, nil
}

func openCLITarget(path string, options csvcore.OpenOptions) (csvcore.Document, string, error) {
	if path != "-" {
		document, err := csvcore.OpenDocument(path, options)
		return document, "file", err
	}
	tmp, err := os.CreateTemp("", "opencsv-stdin-*.csv")
	if err != nil {
		return csvcore.Document{}, "", err
	}
	defer tmp.Close()
	if _, err := io.Copy(tmp, os.Stdin); err != nil {
		return csvcore.Document{}, "", err
	}
	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return csvcore.Document{}, "", err
	}
	document, err := csvcore.OpenDocument(tmp.Name(), options)
	return document, "stdin", err
}

func mergeFormat(base protocol.FormatParams, override protocol.FormatParams) protocol.FormatParams {
	result := base
	if override.Delimiter != "" {
		result.Delimiter = override.Delimiter
	}
	if override.Quote != "" {
		result.Quote = override.Quote
	}
	if override.QuoteMode != "" {
		result.QuoteMode = override.QuoteMode
	}
	if override.Encoding != "" {
		result.Encoding = override.Encoding
	}
	if override.LineEnding != "" {
		result.LineEnding = override.LineEnding
	}
	if override.HeaderRows != nil {
		result.HeaderRows = override.HeaderRows
	}
	if override.HeaderCols != nil {
		result.HeaderCols = override.HeaderCols
	}
	if override.Header != nil {
		result.Header = override.Header
	}
	if override.NoHeader != nil {
		result.NoHeader = override.NoHeader
	}
	return result
}

func parseNonNegative(raw string, name string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, usageError(fmt.Sprintf("%s must be a non-negative integer", name))
	}
	return value, nil
}

func samePath(left string, right string) bool {
	if left == right {
		return true
	}
	leftAbs, leftErr := filepath.Abs(left)
	rightAbs, rightErr := filepath.Abs(right)
	return leftErr == nil && rightErr == nil && leftAbs == rightAbs
}

func openDocument(args []string) error {
	flags := flag.NewFlagSet("open", flag.ExitOnError)
	rows := flags.Int("rows", 50, "preview rows")
	indexRows := flags.Int("index-rows", 0, "rows to index synchronously; -1 indexes the whole file")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("open expects one file path")
	}

	document, err := csvcore.OpenDocument(flags.Arg(0), csvcore.OpenOptions{
		PreviewRows: *rows,
		IndexRows:   *indexRows,
	})
	if err != nil {
		return err
	}
	return printJSON(document)
}

func inspect(args []string) error {
	flags := flag.NewFlagSet("inspect", flag.ExitOnError)
	rows := flags.Int("rows", 20, "preview rows")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("inspect expects one file path")
	}

	preview, err := csvcore.OpenPreview(flags.Arg(0), *rows)
	if err != nil {
		return err
	}
	return printJSON(preview)
}

func generate(args []string) error {
	flags := flag.NewFlagSet("generate", flag.ExitOnError)
	rows := flags.Int("rows", 1000, "data rows")
	cols := flags.Int("cols", 12, "columns")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("generate expects one output path")
	}

	file, err := os.Create(flags.Arg(0))
	if err != nil {
		return err
	}
	defer file.Close()
	return csvcore.GenerateSynthetic(file, *rows, *cols)
}

func save(args []string) error {
	flags := flag.NewFlagSet("save", flag.ExitOnError)
	out := flags.String("out", "", "output path; defaults to in-place save")
	backup := flags.Bool("backup", true, "create a backup before in-place save")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("save expects one file path")
	}

	inputPath := flags.Arg(0)
	document, err := csvcore.OpenDocument(inputPath, csvcore.OpenOptions{PreviewRows: 1})
	if err != nil {
		return err
	}
	records, err := csvcore.ReadAllRecords(inputPath, document.Metadata)
	if err != nil {
		return err
	}

	outputPath := inputPath
	createBackup := *backup
	if *out != "" {
		outputPath = *out
		createBackup = false
	}
	result, err := csvcore.SaveRecords(outputPath, document.Metadata, records, csvcore.SaveOptions{CreateBackup: createBackup})
	if err != nil {
		return err
	}
	return printJSON(result)
}

type grepOptions struct {
	IgnoreCase bool
	Regex      bool
	Column     string
	Limit      int
}

type grepMatch struct {
	RowNumber    int    `json:"rowNumber"`
	ColumnNumber int    `json:"columnNumber"`
	ColumnName   string `json:"columnName,omitempty"`
	Value        string `json:"value"`
}

type grepResult struct {
	Path       string      `json:"path"`
	Query      string      `json:"query"`
	Column     string      `json:"column,omitempty"`
	IgnoreCase bool        `json:"ignoreCase"`
	Regex      bool        `json:"regex"`
	Total      int         `json:"total"`
	Truncated  bool        `json:"truncated"`
	Matches    []grepMatch `json:"matches"`
}

func grep(args []string) error {
	flags := flag.NewFlagSet("grep", flag.ExitOnError)
	ignoreCase := flags.Bool("i", false, "case-insensitive search")
	ignoreCaseLong := flags.Bool("ignore-case", false, "case-insensitive search")
	regex := flags.Bool("regex", false, "treat query as a regular expression")
	column := flags.String("col", "", "column name or 1-based column number to search")
	limit := flags.Int("limit", 100, "maximum matches to include; 0 means unlimited")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 2 {
		return fmt.Errorf("grep expects a query and one file path")
	}
	result, err := grepFile(flags.Arg(1), flags.Arg(0), grepOptions{
		IgnoreCase: *ignoreCase || *ignoreCaseLong,
		Regex:      *regex,
		Column:     *column,
		Limit:      *limit,
	})
	if err != nil {
		return err
	}
	return printJSON(result)
}

func grepFile(path string, query string, options grepOptions) (grepResult, error) {
	if path == "-" {
		return grepResult{}, usageError("grep expects a file path, not stdin")
	}
	if query == "" {
		return grepResult{}, usageError("grep query must not be empty")
	}
	if options.Limit < 0 {
		return grepResult{}, usageError("-limit must be non-negative")
	}
	document, err := csvcore.OpenDocument(path, csvcore.OpenOptions{PreviewRows: 1, IndexRows: -1})
	if err != nil {
		return grepResult{}, err
	}
	records, err := csvcore.ReadAllRecords(path, document.Metadata)
	if err != nil {
		return grepResult{}, err
	}
	columnIndex, err := resolveGrepColumn(options.Column, records, document.Metadata)
	if err != nil {
		return grepResult{}, err
	}
	matches, err := grepMatcher(query, options)
	if err != nil {
		return grepResult{}, err
	}
	result := grepResult{
		Path:       path,
		Query:      query,
		Column:     options.Column,
		IgnoreCase: options.IgnoreCase,
		Regex:      options.Regex,
		Matches:    []grepMatch{},
	}
	header := recordsHeader(records, document.Metadata)
	for rowIndex, record := range records {
		startCol := 0
		endCol := len(record)
		if columnIndex >= 0 {
			startCol = columnIndex
			endCol = columnIndex + 1
		}
		for colIndex := startCol; colIndex < endCol; colIndex++ {
			value := ""
			if colIndex < len(record) {
				value = record[colIndex]
			}
			if !matches(value) {
				continue
			}
			result.Total++
			if options.Limit == 0 || len(result.Matches) < options.Limit {
				result.Matches = append(result.Matches, grepMatch{
					RowNumber:    rowIndex + 1,
					ColumnNumber: colIndex + 1,
					ColumnName:   headerName(header, colIndex),
					Value:        value,
				})
			}
		}
	}
	result.Truncated = result.Total > len(result.Matches)
	return result, nil
}

func resolveGrepColumn(raw string, records [][]string, metadata csvcore.FileMetadata) (int, error) {
	if raw == "" {
		return -1, nil
	}
	if columnNumber, err := strconv.Atoi(raw); err == nil {
		if columnNumber <= 0 {
			return -1, usageError("-col index must be 1 or greater")
		}
		return columnNumber - 1, nil
	}
	header := recordsHeader(records, metadata)
	for index, name := range header {
		if name == raw || strings.EqualFold(name, raw) {
			return index, nil
		}
	}
	return -1, usageError(fmt.Sprintf("column %q not found", raw))
}

func grepMatcher(query string, options grepOptions) (func(string) bool, error) {
	if options.Regex {
		pattern := query
		if options.IgnoreCase {
			pattern = "(?i:" + query + ")"
		}
		expression, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		return func(value string) bool {
			return expression.MatchString(value)
		}, nil
	}
	if options.IgnoreCase {
		needle := strings.ToLower(query)
		return func(value string) bool {
			return strings.Contains(strings.ToLower(value), needle)
		}, nil
	}
	return func(value string) bool {
		return strings.Contains(value, query)
	}, nil
}

func recordsHeader(records [][]string, metadata csvcore.FileMetadata) []string {
	if !metadata.Dialect.HasHeader || len(records) == 0 {
		return nil
	}
	return records[0]
}

func headerName(header []string, index int) string {
	if index >= 0 && index < len(header) {
		return header[index]
	}
	return ""
}

func serve(args []string) error {
	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := flags.String("addr", "127.0.0.1:43255", "listen address")
	webDir := flags.String("web", "web", "static web directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	return appserver.ListenAndServe(*addr, *webDir)
}

func bench(args []string) error {
	flags := flag.NewFlagSet("bench", flag.ExitOnError)
	rows := flags.Int("rows", 100000, "data rows")
	cols := flags.Int("cols", 16, "columns")
	out := flags.String("out", "/tmp/opencsv-bench.csv", "output csv")
	if err := flags.Parse(args); err != nil {
		return err
	}

	start := time.Now()
	file, err := os.Create(*out)
	if err != nil {
		return err
	}
	if err := csvcore.GenerateSynthetic(file, *rows, *cols); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	generateDuration := time.Since(start)

	openStart := time.Now()
	document, err := csvcore.OpenDocument(*out, csvcore.OpenOptions{PreviewRows: 50})
	if err != nil {
		return err
	}
	openDuration := time.Since(openStart)

	indexStart := time.Now()
	index, err := csvcore.BuildRowIndex(*out, document.Metadata.Dialect, 0)
	if err != nil {
		return err
	}
	indexDuration := time.Since(indexStart)

	return printJSON(map[string]any{
		"path":               *out,
		"rows":               *rows,
		"cols":               *cols,
		"sizeBytes":          document.Metadata.SizeBytes,
		"encoding":           document.Metadata.Encoding,
		"lineEnding":         document.Metadata.LineEnding,
		"dialect":            document.Metadata.Dialect,
		"indexedRows":        index.Rows,
		"generateDurationMs": generateDuration.Milliseconds(),
		"openDurationMs":     openDuration.Milliseconds(),
		"indexDurationMs":    indexDuration.Milliseconds(),
	})
}

func printJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
