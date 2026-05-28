package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
