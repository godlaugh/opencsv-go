package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/opencsv/opencsv/pkg/csvcore"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "inspect":
		err = inspect(os.Args[2:])
	case "generate":
		err = generate(os.Args[2:])
	case "bench":
		err = bench(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "opencsv:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  opencsv inspect [-rows 20] <file>")
	fmt.Fprintln(os.Stderr, "  opencsv generate [-rows 1000] [-cols 12] <file>")
	fmt.Fprintln(os.Stderr, "  opencsv bench [-rows 100000] [-cols 16] [-out /tmp/opencsv-bench.csv]")
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

	previewStart := time.Now()
	preview, err := csvcore.OpenPreview(*out, 50)
	if err != nil {
		return err
	}
	previewDuration := time.Since(previewStart)

	indexStart := time.Now()
	index, err := csvcore.BuildRowIndex(*out, preview.Dialect, 0)
	if err != nil {
		return err
	}
	indexDuration := time.Since(indexStart)

	return printJSON(map[string]any{
		"path":               *out,
		"rows":               *rows,
		"cols":               *cols,
		"sizeBytes":          preview.SizeBytes,
		"dialect":            preview.Dialect,
		"indexedRows":        index.Rows,
		"generateDurationMs": generateDuration.Milliseconds(),
		"previewDurationMs":  previewDuration.Milliseconds(),
		"indexDurationMs":    indexDuration.Milliseconds(),
	})
}

func printJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
