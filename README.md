# opencsv

opencsv is a local-first CSV and tabular data workbench. The goal is to keep the fast, familiar feel of a spreadsheet while handling large files, typed queries, and repeatable data-cleaning workflows more gracefully than traditional CSV editors.

This repository starts with a Go-first M0 milestone:

- Go CSV core for dialect sniffing, preview, synthetic data generation, and row offset indexing.
- CLI for inspection and benchmark runs.
- Static Canvas grid prototype for UI and virtual scrolling experiments.
- Architecture and milestone notes for the desktop app that will follow.
- Product PRD and architecture v2 notes for parity development.

## Why Go

SmoothCSV's public technical notes describe a Rust/Tauri architecture. opencsv is intentionally exploring a Go core instead:

- Go is productive for file I/O, concurrency, background indexing, and cross-platform services.
- A Go backend can power desktop, local web, and future server modes from the same core.
- The core can expose a narrow RPC surface to any shell: Wails, a Tauri sidecar, Electron, or a native WebView.

The current direction is `Go core + browser-grade Canvas UI`, with the desktop shell still open. Wails is the most natural first candidate if we want a pure Go desktop packaging story.

## Run

Format and test:

```sh
make fmt
make test
```

Generate and inspect a CSV:

```sh
go run ./cmd/opencsv generate -rows 1000 -cols 12 /tmp/opencsv-sample.csv
go run ./cmd/opencsv open -rows 5 -index-rows 5 /tmp/opencsv-sample.csv
go run ./cmd/opencsv inspect -rows 5 /tmp/opencsv-sample.csv
```

Round-trip save a CSV while preserving detected encoding and line endings:

```sh
go run ./cmd/opencsv save -out /tmp/opencsv-sample-copy.csv /tmp/opencsv-sample.csv
go run ./cmd/opencsv save /tmp/opencsv-sample.csv
```

Run the M0 benchmark:

```sh
make bench
```

Run the grid prototype:

```sh
make serve
```

Then open `http://localhost:43255`.

## Current Shape

```text
cmd/opencsv       CLI entrypoint
pkg/appserver     Local browser prototype API
pkg/csvcore       CSV dialect, preview, indexing, sample generation
web/              static Canvas grid prototype
docs/             planning and architecture notes
```

Key planning documents:

- [Product PRD](docs/product-prd.md)
- [Architecture v2](docs/architecture-v2.md)
- [Product Parity Progress](docs/parity-progress.md)
- [M1 Verification](docs/m1-verification.md)

## M0 Acceptance Criteria

- A clean Go module that builds without third-party dependencies.
- CSV preview returns a first screen of records without loading the whole file.
- Row indexing records byte offsets for future random access.
- A local UI prototype demonstrates large-grid scrolling and cell editing.
- Benchmarks give us a baseline before adding DuckDB/DataFusion-style query capabilities.

## M1 Core Progress

The core and CLI now support:

- Encoding detection for UTF-8, UTF-8 BOM, UTF-16 LE/BE with BOM, heuristic UTF-16, and a conservative legacy single-byte fallback.
- Line ending detection for LF, CRLF, and CR.
- Delimiter and header detection.
- Fast document open with preview rows that does not wait for full-file indexing.
- Row offset indexing as a separate operation.
- CSV save with automatic backup and preservation of detected encoding and line endings.
- Local web API and Canvas grid integration through `opencsv serve`.

Current local benchmark:

```text
112,857,510 byte CSV
openDurationMs: 0
indexDurationMs: 318
indexedRows: 720001
```

Open a real file in the prototype:

```sh
make serve
```

Then click Open with an empty path field to use the browser file picker, or enter an absolute CSV path in the toolbar. Browser-picked files are uploaded into a local temporary copy before parsing; absolute paths are opened directly by the local server. Edit a cell and press Save. In-place saves create a timestamped `.bak` file first.
