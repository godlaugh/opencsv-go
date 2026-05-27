# Milestones

## M0: Foundation

Goal: prove the Go-first core and UI performance path.

- Initialize Go module and repo structure.
- Implement CSV dialect sniffing.
- Implement preview loading without whole-file reads.
- Implement row offset indexing.
- Add synthetic data generator and benchmark CLI.
- Add static Canvas grid prototype.
- Document architecture and next milestones.

Status: complete in `v0.1`.

## M1: Real File Opening

Goal: make opencsv useful for read-only inspection of real CSV and TSV files.

- Core and CLI open path.
- Encoding detection for UTF-8, UTF-8 BOM, UTF-16 LE/BE with BOM, heuristic UTF-16, and legacy single-byte fallback.
- Line ending detection for LF, CRLF, and CR.
- Delimiter and header detection.
- Fast preview loading without waiting on full-file indexing.
- Row offset indexing.
- Basic save, automatic backup, and preservation of detected encoding and line endings.
- Local API and Canvas prototype integration through `opencsv serve`.
- File picker through the chosen desktop shell.
- Background row indexing with progress and cancellation.
- Cell window API: fetch rows by offset range.
- Encoding conversion support for Shift-JIS, GBK, and Windows-1252.
- Delimiter and header override UI.
- Error row reporting for malformed records.

Acceptance target: 100 MB CSV opens to first screen in under 2 seconds on a modern laptop.

Core verification: a 112,857,510 byte generated CSV opened to first preview in 0 ms and completed full row indexing in 318 ms on 2026-05-26.

## M2: Editing MVP

Goal: spreadsheet-like editing that does not corrupt raw data.

- Single-cell edit, fill, paste, insert/delete rows and columns.
- Undo/redo operation log.
- Dirty state and atomic save.
- Backup and recovery file.
- Preserve delimiter, line endings, quoting style, and BOM by default.

Acceptance target: common CSV editing no longer requires Excel.

## M3: Grid UX

Goal: beat the feel of existing CSV editors.

- Frozen row and column headers.
- Column resize, reorder, hide, pin.
- Range selection and copy/paste.
- Keyboard shortcuts and command palette.
- Column menu with sort, filter, profile, rename, clean.
- Multi-column sort composition from column headers.
- Responsive inspector for selected column and cell.

Acceptance target: million-row scrolling stays fluid and selection feels native.

## M4: Filtering and Profiling

Goal: make data quality visible and filters effortless.

- Excel-style value checklist filters.
- Text, number, date, null, and invalid filters.
- Type inference with confidence and sample evidence.
- Column stats: nulls, distinct, min/max, common values.
- Filter chips that can switch between UI builder and expression mode.

Acceptance target: filtering is faster than writing SQL for everyday cases.

## M5: SQL Workbench

Goal: turn opencsv from editor into data workbench.

- Query current file and multiple opened files.
- Typed schema override.
- Query history and saved snippets.
- Result grid and export.
- Friendly errors with column names and sample bad values.

Acceptance target: numeric and date comparisons behave intuitively without hidden source mutation.

## M6: Reliability and Scale

Goal: make large-file work boring in the best way.

- 1 GB CSV first-screen benchmark.
- Sort/filter/query progress with cancellation.
- External file change detection.
- Crash recovery.
- Memory budget tests.
- Cross-platform packaging.

Acceptance target: opencsv is stable enough for daily use on macOS, Windows, and Linux.

## M7: Extensions and AI

Goal: make repeatable data work programmable.

- Command API.
- User scripts or WASM transforms.
- Plugin manifest and permission model.
- AI-assisted filters, SQL, transforms, and schema suggestions.
- Preview-first execution model for AI actions.

Acceptance target: AI speeds up data work without silently changing files.
