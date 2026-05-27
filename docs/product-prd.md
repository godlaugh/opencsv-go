# opencsv Product Requirements Document

Status: Active product specification for parity development

Last updated: 2026-05-27

## 1. Source Baseline

This PRD is derived from public SmoothCSV product materials and official project documentation:

- SmoothCSV homepage: https://smoothcsv.com/
- SmoothCSV blog index: https://smoothcsv.com/blog/
- SmoothCSV GitHub README/release repository: https://github.com/kohii/smoothcsv3
- SmoothCSV CLI specification: https://github.com/kohii/smoothcsv3/blob/main/docs/cli.md
- SmoothCSV URL protocol specification: https://github.com/kohii/smoothcsv3/blob/main/docs/smoothcsv-url.md
- SmoothCSV SQL row filter documentation: https://github.com/kohii/smoothcsv3/blob/main/docs/sql-in-row-filter.md
- SmoothCSV SQL console documentation: https://github.com/kohii/smoothcsv3/blob/main/docs/sql-in-sql-console.md
- SmoothCSV technical articles describing Tauri/Rust/React architecture and virtual scrolling.

Latest public source check: 2026-05-27. The GitHub README currently lists v3 feature areas beyond the first homepage pass, including CSV Compare, quote indicators, automatic backup, presets/rules, quick column filters/sorting, shortcut customization, aggregation, read-only mode, transpose, auto-refresh, opening hyperlinks from cells, row/column dedupe, text transforms, join/split, Excel import/export, dark/light themes, and localization. The release page currently shows v3.20.1 as latest, with recent release notes adding CLI stdin, compare merge/actions, inline character-level diff, compare ignore options, header-name column alignment, redesigned dark mode, CLI diff launch, header-in-selection controls, and faster Canvas grid editor work.

The goal is not to clone SmoothCSV source code. The goal is to build an open, Go-first CSV workbench with comparable product coverage, stronger UX discipline, and measurable performance.

## 2. Product Summary

opencsv is a local-first desktop-grade CSV editor and tabular data workbench.

It must let users open, inspect, edit, transform, query, compare, and save CSV-like files quickly without losing file fidelity. The product should feel like a spreadsheet for direct manipulation, like VS Code for command/search workflows, and like a lightweight database console for SQL-style analysis.

## 3. Target Users

### 3.1 Data Operators

People who repeatedly inspect CSV exports, fix bad rows, normalize headers, compare files, filter records, and hand results to another tool.

Primary needs:

- Fast open, search, filter, and save.
- No silent data corruption.
- Familiar spreadsheet interactions.
- Good behavior on large files.

### 3.2 Engineers and Automation Users

People who use CSV files in scripts, ETL, tests, migrations, and debugging.

Primary needs:

- CLI open command.
- URL/deep-link support.
- Diff mode.
- SQL console.
- Reliable encoding and delimiter overrides.

### 3.3 Analysts

People who need quick ad hoc questions answered from CSV files without importing into a database.

Primary needs:

- SQL query console.
- Row filters.
- Sort/group/aggregate workflows.
- Export of query/filter results.

## 4. Product Principles

1. Preserve file fidelity by default.
   The app must preserve encoding, BOM, delimiter, quote behavior, and line endings unless the user explicitly changes them.

2. First screen before full work.
   Opening large files must show useful data before full indexing, profiling, or query preparation is complete.

3. Direct manipulation first, commands second.
   Common workflows must be visible in the UI; command palette accelerates but must not be the only path.

4. Local-first and private.
   File contents stay local unless the user explicitly exports or uses a future external integration.

5. Measurable performance.
   Large-file claims must be backed by benchmark scripts and repeatable acceptance results.

6. Browser prototype is not the final desktop shell.
   Web UI may support uploaded temporary copies. Native file picker, original-path persistence, menus, and OS integration belong in the desktop shell.

## 5. Competitive Feature Baseline

### 5.1 Public Homepage Feature Areas

Product claims to match:

- CSV editing with spreadsheet-like grid.
- Fast large-file open; public benchmark claim mentions 100 MB file much faster than Excel.
- Find and replace.
- Sort and filter.
- SQL support.
- Multi-cell edit.
- Cross-platform desktop support.
- Modern UI with smooth interaction.

### 5.2 Public README Feature Areas

Product claims to match:

- Open CSV, TSV, and other delimited files.
- Auto-detect delimiter, quote, encoding, header rows.
- Preserve and expose flexible quote rules by header, body, and column.
- Show quote indicators for cells that will be quoted on save.
- View and edit large files.
- Add, delete, move, copy, paste rows and columns.
- Hide rows/columns and dedupe rows/columns.
- Sort rows.
- Filter rows.
- Visual condition-builder filtering.
- Search and replace.
- Multi-cell selection/editing.
- Undo/redo.
- Command palette.
- Keyboard shortcuts.
- Keyboard shortcut customization.
- Theme support, including light/dark/system auto-detect.
- File format presets.
- File format rules for automatic preset application.
- CSVLint-like validation.
- Query with SQL.
- Compare/diff two CSV files.
- Compare inline character-level changes, merge changes, ignore whitespace/case/columns, align columns by header name, and compare with clipboard.
- CLI integration.
- URL protocol integration.
- Clipboard integration.
- Export/copy in multiple formats.
- Copy selections as Markdown, HTML, JSON, SQL, LaTeX, and raw values.
- Copy raw values without CSV quoting/escaping.
- Excel import/export.
- Auto-fill with intelligent pattern detection.
- Text transforms and join/split cell tools.
- Cell value panel.
- Auto-backup before overwrite.
- Quick filters and sorting from column menus.
- Aggregation summaries and aggregate rows.
- Read-only mode.
- Transpose rows and columns.
- Auto-refresh for external file modifications.
- Open hyperlinks from cells.
- Internationalization for English, Japanese, French, Spanish, Italian, Simplified Chinese, Portuguese, Russian, German, and Korean.

### 5.3 Official Docs Feature Areas

CLI:

- `smoothcsv [options] [path ...]`
- Open single/multiple files.
- Open stdin.
- CLI `--diff old.csv new.csv` launches Compare view.
- `--goto` row/cell/range selection.
- `--diff` side-by-side compare.
- Delimiter, quote, quote mode, encoding, header rows, header columns overrides.
- New window/reuse window.
- Wait for close.
- Version/help.
- Exit codes.

URL:

- `smoothcsv://file/path/to/file.csv[:location][?params]`
- Row/cell/range locations.
- Delimiter, quote, quote mode, encoding query params.

SQL row filter:

- SQL WHERE syntax.
- Header names or `c1`, `c2`, ... references.
- String, number, boolean, bigint, null, date, time, timestamp, interval types.
- `CAST` and `TRY_CAST`.
- String, numeric, and conditional functions.

SQL console:

- SQLite-compatible SQL syntax.
- Query current opened files and filesystem paths.
- File references via direct path, `file://`, and `@file:`.
- Default columns `c1`, `c2`, ...
- Explicit casts for reliable numeric comparison.

## 6. Scope

### 6.1 Must Ship for Product Parity v1

- Desktop-grade app shell or local app mode with native-feeling file open/save.
- Fast CSV/TSV opening.
- Grid editing.
- Find/replace.
- Sort/filter.
- SQL row filter.
- SQL console.
- Diff mode.
- CLI.
- URL/deep-link parser.
- File format overrides.
- Save with backup and fidelity preservation.
- Quote indicators.
- Copy/export formats.
- Read-only mode.
- Auto-refresh.
- Aggregation status.
- Automated tests and performance benchmarks.

### 6.2 Can Ship After Product Parity v1

- Extension marketplace.
- AI-assisted transforms.
- Cloud collaboration.
- Browser-hosted SaaS version.
- Full Excel workbook editing beyond CSV import/export.

## 7. Functional Requirements

### FR-001 File Open

The app shall open CSV, TSV, and delimiter-separated text files from:

- Native file picker in desktop shell.
- Absolute or relative path from CLI/API.
- Browser file picker in prototype mode, as temporary uploaded copy.
- Stdin from CLI.
- URL protocol/deep link.

Acceptance:

- Opening a local CSV displays first screen of rows.
- Opening a TSV detects tab delimiter.
- Opening a 100 MB CSV displays first screen in under 2 seconds.
- Prototype mode labels uploaded files as temporary copies.

### FR-002 Format Detection

The app shall detect:

- Delimiter: comma, tab, semicolon, pipe, space, custom single character.
- Quote character: double quote, single quote, none.
- Quote mode: minimal, always, never, unknown.
- Encoding: UTF-8, UTF-8 BOM, UTF-16 LE/BE, Shift-JIS, GBK, Windows-1252, legacy single-byte fallback.
- Line ending: LF, CRLF, CR.
- Header rows.
- Header columns.

Acceptance:

- Detection result is visible in status/format panel.
- User can override detection.
- Saving preserves detected values unless overridden.

### FR-003 Grid View

The app shall render a spreadsheet-like grid with:

- Stable column header row.
- Stable row number gutter.
- Smooth vertical and horizontal scrolling.
- Cell selection.
- Range selection.
- Multi-cell selection.
- Frozen header behavior.
- Column resize.
- Column reorder.
- Column hide/show.
- Quote indicators for cells that will be quoted on save.
- Scroll position persistence per file.
- Open hyperlink command for URL-like cell values.

Acceptance:

- Header must not be covered by body rows during scroll.
- Header height must remain visually stable.
- 1M rows must scroll without DOM row explosion.

### FR-004 Editing

The app shall support:

- Edit single cell.
- Paste tabular clipboard data.
- Fill selection.
- Multi-cell edit.
- Insert/delete rows.
- Insert/delete columns.
- Move rows/columns.
- Copy/cut/paste rows and columns.
- Undo/redo.
- Dirty state.
- Read-only mode that disables editing operations.
- Atomic save.
- Save as.
- Auto-backup before destructive in-place save.
- Auto-fit column width on edit when enabled.

Acceptance:

- Editing then saving produces expected CSV.
- Undo restores prior state.
- Backup is byte-for-byte original before save.
- Save does not silently normalize encoding or line endings.

### FR-005 Search and Replace

The app shall support:

- Find within current file.
- Replace current match.
- Replace all.
- Case sensitivity toggle.
- Whole cell toggle.
- Regex toggle.
- Search within selection.
- Match count.
- Keyboard shortcut.

Acceptance:

- Find scrolls/selects matching cell.
- Replace all reports changed cell count.

### FR-006 Row Filter

The app shall support:

- Quick filter text input.
- SQL WHERE style filter.
- Excel-like value checklist filter.
- Empty/non-empty filters.
- Numeric, text, date, boolean filters.
- Filter chips.
- Clear all filters.
- Open current filter in SQL console.

Acceptance:

- Header names work when header rows are enabled.
- `c1`, `c2`, ... work when no header row exists.
- `CAST` and `TRY_CAST` work for numeric/date cases.
- Filtered row count is visible.

### FR-007 Sort

The app shall support:

- Sort ascending/descending by column.
- Multi-column sort.
- Numeric/text/date modes.
- Stable sort option.
- Header-aware sorting.
- Quick column menu sorting.

Acceptance:

- Header rows remain at top.
- Sort can be undone.

### FR-008 SQL Console

The app shall support:

- SQL console tab.
- SQLite-compatible query syntax.
- Query current opened files.
- Query file paths.
- Query unsaved/temporary documents.
- Results grid.
- Export/save query result.
- Query history.
- Query errors with useful locations/messages.

Acceptance:

- `SELECT * FROM "/path/to/file.csv"` works.
- `SELECT * FROM "@file:/path/to/file.csv"` works.
- Header names and `c1` columns work.
- Numeric comparison docs explain cast behavior.

### FR-009 Diff/Compare

The app shall support side-by-side file compare:

- CLI `--diff old.csv new.csv`.
- GUI open diff.
- Row/column aligned view.
- Highlight changed cells.
- Inline character-level diff inside modified cells.
- Highlight added/deleted rows.
- Navigate next/previous diff.
- Merge row/cell changes from one side to the other.
- Ignore whitespace, case, selected columns, or configured ignored columns.
- Align columns by header name when order differs.
- Compare active file with clipboard.
- Edit, sort, filter, and search inside compare view.
- Respect per-file format detection/override.

Acceptance:

- Diff requires exactly two files in CLI.
- Diff rejects stdin.
- Goto can target either side.

### FR-010 CLI

The app shall provide an `opencsv` command:

- `opencsv [options] [path ...]`
- Open one or many files.
- `opencsv -` from stdin.
- `--goto`, `-g`.
- `--diff`.
- `--delimiter`, `-d`.
- `--quote`, `-q`.
- `--quote-mode`, `--qm`.
- `--encoding`, `-e`.
- `--header-rows`.
- `--header-columns`.
- `--header`, `--no-header`.
- `--new-window`, `-n`.
- `--reuse-window`, `-r`.
- `--wait`, `-w`.
- `--version`, `-v`.
- `--help`, `-h`.

Acceptance:

- Invalid CLI usage exits with code 2.
- General runtime errors exit with code 1.
- Success exits with code 0.

### FR-011 URL Protocol

The app shall parse and handle:

- `opencsv://file/path/to/file.csv`
- `opencsv://file/path/to/file.csv:5`
- `opencsv://file/path/to/file.csv:1:2`
- `opencsv://file/path/to/file.csv:1-5`
- `opencsv://file/path/to/file.csv:1:2-10:20`
- Query params for delimiter, quote, quote mode, encoding.

Acceptance:

- Coordinates are 1-indexed.
- Windows-style paths are accepted.
- URL open can jump to target row/cell/range.

### FR-012 Command Palette and Shortcuts

The app shall support:

- Command palette.
- Search commands by name.
- Keyboard shortcuts.
- Shortcut customization.
- Commands for open, save, find, filter, SQL console, format options, diff, sort.

Acceptance:

- Every menu action has a command ID.
- Commands can be invoked by palette or shortcut.

### FR-013 File Format Presets and Rules

The app shall support named presets and automatic rules:

- Delimiter.
- Quote character.
- Quote mode, including minimal, always, and never.
- Header/body/column-specific quote rules.
- Encoding.
- Header rows.
- Header columns.
- Line ending.
- Rule matching by filename/path/content sample.

Acceptance:

- User can apply preset to open file.
- User can save current detection/override as preset.
- Matching rules can automatically apply a preset.

### FR-014 Validation

The app shall validate:

- Ragged rows.
- Unclosed quotes.
- Invalid encoding sequences.
- Duplicate headers.
- Empty headers.
- Type-profile anomalies.

Acceptance:

- Problems panel lists issue type, row, column, message.
- Clicking issue navigates to location.

### FR-015 Clipboard and Export

The app shall support:

- Copy cells as TSV.
- Copy cells as CSV.
- Copy cells as raw values without CSV quoting/escaping.
- Copy selections as Markdown, HTML, JSON, SQL, and LaTeX.
- Copy row/column.
- Paste spreadsheet ranges.
- Export filtered rows.
- Export SQL result.
- Open and save `.xlsx` files at least as import/export conversion.

Acceptance:

- Copy/paste round-trips with spreadsheet apps for simple ranges.

### FR-016 Aggregation and Profiling

The app shall support:

- Selection summary in status bar: sum, average, minimum, maximum, count.
- Add aggregate rows to the grid.
- Column profile basics: type, null ratio, distinct count, min/max for numeric/date-like values.

Acceptance:

- Selecting numeric cells updates aggregate metrics without blocking scroll.
- Aggregate rows can be removed without changing source data.

### FR-017 Transpose

The app shall support transposing rows and columns.

Acceptance:

- Transpose works on small and medium files.
- The operation is undoable and warns when the result would be too large for memory.

### FR-018 Auto-Refresh and External Changes

The app shall detect external file changes and support automatic reload.

Acceptance:

- Dirty local edits are protected from silent overwrite.
- Read-only files can auto-refresh without changing dirty state.

### FR-019 Internationalization and Accessibility

The app shall support:

- English UI initially.
- Architecture for multiple locales.
- Target locales: Japanese, French, Spanish, Italian, Simplified Chinese, Portuguese, Russian, German, Korean.
- Keyboard navigation.
- Accessible names for controls.
- Light, dark, and system-following themes.
- High contrast theme support.

Acceptance:

- Main workflows can be completed by keyboard.
- UI text is not hard-coded inside rendering primitives where localization is impossible.

## 8. Non-Functional Requirements

### NFR-001 Performance

Targets:

- 100 MB CSV first screen: under 2 seconds.
- 1 GB CSV first screen: under 5 seconds.
- 1M rows scroll: responsive at 60 FPS target, no body DOM row explosion.
- Full row index for 100 MB: under 2 seconds.
- Find over 100 MB: progress visible, cancellable.
- SQL query over 100 MB: progress visible, cancellable.

### NFR-002 Memory

Targets:

- 100 MB file open should not require loading all decoded records into UI memory.
- Grid should request windows of rows.
- Edit overlay stores deltas until save.
- Large operations have explicit memory budgets.

### NFR-003 Reliability

Targets:

- Atomic save via temp file and rename.
- Byte-for-byte backup before in-place save.
- Crash recovery for dirty documents.
- External file change detection.
- Clear warning before overwriting changed external file.

### NFR-004 Security and Privacy

Targets:

- Local file contents stay local.
- No telemetry by default.
- Browser prototype labels temporary uploads.
- Desktop shell uses OS-native file access.

### NFR-005 Cross-Platform

Targets:

- macOS first.
- Windows and Linux supported by architecture.
- CLI behavior consistent across OS path conventions.

## 9. Information Architecture

Primary surfaces:

- Workspace/sidebar: opened files, diff sessions, SQL consoles.
- Toolbar: open, save, find, filter, sort, command palette.
- Grid view: table editing.
- Inspector: selected column/cell profile and metadata.
- Status bar: row/column counts, encoding, delimiter, header, dirty state.
- Problems panel: validation issues.
- SQL console: query editor and results.
- Diff view: side-by-side comparison.

## 10. Key User Flows

### UF-001 Open and Save CSV

1. User clicks Open.
2. Desktop app shows native file picker; browser prototype shows browser file picker.
3. App detects format and displays first screen.
4. User edits cell.
5. Dirty state appears.
6. User saves.
7. App creates backup and atomically writes updated file.

### UF-002 Filter Rows

1. User opens file.
2. User opens filter.
3. User chooses value checklist or SQL WHERE expression.
4. Grid shows filtered rows and count.
5. User clears filter or opens filter in SQL console.

### UF-003 SQL Analysis

1. User opens SQL console.
2. User writes query referencing current CSV.
3. App executes query.
4. Results appear in grid.
5. User exports results.

### UF-004 CLI Diff

1. User runs `opencsv --diff old.csv new.csv`.
2. App opens diff view.
3. User navigates changed cells.
4. User closes app; CLI returns.

## 11. Data Model Requirements

Document:

- ID.
- Source kind: file path, stdin, upload, untitled, SQL result.
- Original path.
- Working path.
- Format metadata.
- Row index.
- Header config.
- Edit overlay.
- Dirty state.
- Validation issues.

Format metadata:

- Encoding.
- BOM.
- Line ending.
- Delimiter.
- Quote character.
- Quote mode.
- Header rows.
- Header columns.

Grid window:

- Start row.
- Row count.
- Records.
- Applied edits.
- Filter/sort mapping if active.

## 12. Architecture Requirements

### 12.1 Go Core

The Go core owns:

- File open.
- Format detection.
- Encoding decode/encode.
- Row indexing.
- Windowed row reads.
- Edit overlay.
- Save/backup.
- Search.
- Filter.
- Sort.
- Validation.
- SQL engine integration.
- CLI.
- URL parser.

### 12.2 UI Layer

The UI owns:

- Canvas grid rendering.
- Selection.
- Keyboard/mouse interactions.
- Menus/toolbars/panels.
- Command palette.
- Status and progress display.

### 12.3 Desktop Shell

The desktop shell owns:

- Native file picker.
- Native save dialog.
- App menu.
- File association.
- URL protocol registration.
- CLI handoff to running app.
- Window management.

## 13. Acceptance Matrix

| ID | Requirement | Evidence |
| --- | --- | --- |
| A-001 | 100 MB first screen < 2s | `opencsv bench` output |
| A-002 | Header remains fixed during scroll | Browser screenshot/video and pixel/assertion test |
| A-003 | Save preserves UTF-8 BOM + CRLF | Unit/integration test |
| A-004 | Save creates byte-for-byte backup | Unit/integration test |
| A-005 | Browser Open shows file picker | Browser test with file chooser |
| A-006 | Direct path open works | API/CLI test |
| A-007 | Windowed row loading works | Unit/API test |
| A-008 | Filter SQL works | Unit/API/UI test |
| A-009 | SQL console works | Unit/API/UI test |
| A-010 | CLI options match spec | CLI golden tests |
| A-011 | URL parser handles locations/params | Unit tests |
| A-012 | Diff rejects invalid arg counts | CLI tests |
| A-013 | Diff highlights changed cells | UI test |

## 14. Release Milestones

### M1.1 Stabilize Current Open/Save Prototype

- Remove misleading placeholder UI.
- Fix grid header/body rendering.
- Implement browser file-picker upload path.
- Keep direct absolute path open for local server.
- Keep current M1 performance benchmark.

### M1.2 Product Shell Architecture

- Add command registry.
- Add document/session model.
- Add format metadata panel.
- Add progress/cancellation primitives.
- Add architecture tests.

### M2 Editing MVP

- Edit overlay.
- Save as.
- Undo/redo.
- Row/column insert/delete.
- Clipboard paste/copy.

### M3 Grid UX

- Range selection.
- Column resize/reorder/hide.
- Fixed row/column headers.
- Keyboard navigation.
- Find/replace.

### M4 Filter/Sort/Profile

- Sort engine.
- SQL row filter.
- Value checklist filter.
- Column profile.
- Problems panel.

### M5 SQL Console

- SQLite engine integration.
- Query current/open/path files.
- Results grid.
- Export result.

### M6 Diff/CLI/URL Parity

- Diff view.
- Full CLI spec.
- URL parser/handler.
- Window management hooks.

### M7 Performance and Cross-Platform Hardening

- 1 GB benchmark.
- Memory budget tests.
- Desktop packaging.
- macOS/Windows/Linux smoke tests.

## 15. Open Product Decisions

- Desktop shell: Wails vs Tauri sidecar vs custom native shell.
- SQL engine: SQLite virtual table vs DuckDB vs DataFusion-like Go integration.
- Encoding library choice for Shift-JIS/GBK/Windows-1252.
- Whether browser prototype should support save-back-to-original through File System Access API where available.
- Exact UI shape of command palette and menu system.
