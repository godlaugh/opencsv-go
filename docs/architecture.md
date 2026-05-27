# Architecture

## Direction

opencsv is organized around a Go core and a high-performance grid UI.

```mermaid
flowchart LR
  File["CSV / TSV / future formats"] --> Core["Go csvcore"]
  Core --> Index["Row offset index"]
  Core --> Profile["Type and quality profile"]
  Core --> Query["Query engine"]
  Core --> Export["View export"]
  Core --> RPC["Local RPC boundary"]
  RPC --> Grid["Canvas grid UI"]
  RPC --> Shell["Desktop shell"]
```

## Core Principles

- First screen fast: parse enough to show data, then continue indexing in the background.
- Raw fidelity: preserve delimiter, quoting, BOM, encoding, and line endings unless the user asks to normalize.
- Typed views: types help filtering and SQL, but raw cell values remain available.
- Active views are explicit state: an empty filter result must remain distinct from no filter.
- Cancelable work: indexing, profiling, sorting, and querying must be cancellable.
- Narrow UI boundary: frontend asks for windows of cells, column metadata, and operations; it does not own parsing.

## Go Packages

`pkg/csvcore` currently owns:

- `Sniff`: delimiter, BOM, column count, and header heuristics.
- `OpenPreview`: first-screen records without whole-file load.
- `BuildRowIndex`: record start byte offsets using `encoding/csv.Reader.InputOffset`.
- `EncodeRecords`: UTF-8 CSV/TSV serialization using source dialect and line-ending policy.
- `GenerateSynthetic`: deterministic benchmark data.

`pkg/document` currently owns session view state, edit overlays, undo/redo snapshots, visible row windows, filtering/sorting, and current-view export. `ViewActive` intentionally separates an empty active view from the unfiltered full table so filters, windows, and exports agree on row counts.

Planned packages:

- `pkg/table`: immutable table window protocol and edit overlay.
- `pkg/profile`: type inference, nulls, distinct counts, min/max, invalid samples.
- `pkg/query`: SQL execution over current files.
- `pkg/project`: workspace, recent files, settings, and session state.
- `pkg/rpc`: JSON-RPC or Connect/gRPC bridge to the UI.

## Desktop Shell Options

Wails is the leading candidate for a Go-native desktop shell. Tauri remains possible with Go as a sidecar, but that gives up the simplicity of a mostly Go stack. The M1 decision should be based on file dialog quality, menu integration, packaging friction, and UI RPC ergonomics.
