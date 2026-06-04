# opencsv Product Requirements Document

Status: Active product specification for parity development

Last updated: 2026-06-02

## 1. Source Baseline

This PRD is derived from public SmoothCSV product materials and official project documentation:

- SmoothCSV homepage: https://smoothcsv.com/
- SmoothCSV feature catalog: https://smoothcsv.com/features
- SmoothCSV blog index: https://smoothcsv.com/blog/
- SmoothCSV GitHub README/release repository: https://github.com/kohii/smoothcsv3
- SmoothCSV CLI specification: https://github.com/kohii/smoothcsv3/blob/main/docs/cli.md
- SmoothCSV URL protocol specification: https://github.com/kohii/smoothcsv3/blob/main/docs/smoothcsv-url.md
- SmoothCSV SQL row filter documentation: https://github.com/kohii/smoothcsv3/blob/main/docs/sql-in-row-filter.md
- SmoothCSV SQL console documentation: https://github.com/kohii/smoothcsv3/blob/main/docs/sql-in-sql-console.md
- SmoothCSV technical articles describing Tauri/Rust/React architecture and virtual scrolling.

Latest public source check: 2026-06-02. The homepage still centers spreadsheet-like editing, find/replace, sort/filter, SQL queries, multi-cell editing, command palette access, large-file speed, and cross-platform desktop availability. The current feature catalog adds more explicit product-surface requirements: smart selection, auto-fill, swap cells, aggregate rows, quick value filters, multi-column sort with type/direction per column, CSV compare against files or clipboard, quoting control including manual cell-level preservation, file format presets/rules, copy/export formats, CLI/deep links, Windows portable and AUR distribution, auto-backup, and auto-refresh. The GitHub README currently lists v3 feature areas beyond the first homepage pass, including CSV Compare, quote indicators, automatic backup, presets/rules, quick column filters/sorting, shortcut customization, aggregation, read-only mode, transpose, auto-refresh, opening hyperlinks from cells, row/column dedupe, text transforms, join/split, Excel import/export, dark/light themes, and localization. The GitHub release API currently shows v3.21.0 as latest, published 2026-06-02, with recent release notes adding split editor groups, drag/drop tab movement between groups, group focus/move commands, saved and recent filter query history, macOS stay-open behavior after closing all windows, active-tab visibility fixes, URL-cell selection behavior, Windows preset rule matching fixes, and tab overflow behavior. Recent release cadence remains active: v3.20.0 shipped 2026-05-14, v3.20.1 shipped 2026-05-16, v3.20.2 shipped 2026-05-22, and v3.21.0 shipped 2026-06-02.

Additional official source details captured from the current public pages:

- The homepage positions SmoothCSV as an intuitive spreadsheet-like CSV editor with find/replace, sort/filter, SQL queries, multi-cell editing, command palette access to all features, large-file speed, and Windows/macOS/Linux availability.
- The feature catalog positions reusable column workflows as first-class product behavior: quick value filters, visual conditions, SQL filters, smart sort, multi-column sort, and type-aware sorting must remain discoverable from column/header surfaces and command workflows. It also defines Smart Selection as explicit range add/subtract/invert behavior, not just ordinary drag selection.
- The download page publishes macOS Universal, Windows user/MSI/portable builds for x64 and ARM64, Linux `.deb`, `.rpm`, AppImage packages, and AUR distribution. Product parity must therefore keep packaging and CLI assumptions cross-platform.
- The pricing/license pages state that every feature remains free, licenses are one-time support purchases, work offline, and have no per-person device limit. opencsv does not need to clone SmoothCSV's commercial flow, but it must preserve the same local/offline feature availability principle.
- The "Best CSV Editors in 2026" blog emphasizes dedicated CSV editing, data preservation, large-file handling, quote preservation, frequent updates, thoughtful UX, and SQL console support as the competitive baseline.
- The SmoothCSV v1 archive calls out long-running requirements that still matter: multiple encodings, many CSV/TSV styles, Excel-like operation, ragged-row editing, embedded newline cells, SQL `SELECT`, fast large-file open, tabbed UI, macros/key macros, and CSV grep search.
- The v3 technical article describes a command-centric, VS Code-like architecture with command metadata, keybindings, context-aware enablement, menu/context-menu/button dispatch, lazy-loaded system extensions, and local-first desktop operation with no external communication except update checks.
- The virtual scrolling article identifies browser element-size limits around million-row datasets and requires a custom scrolling strategy or canvas rendering that is not capped by DOM spacer height.
- The official SQL Console documentation requires SQLite syntax, direct path/`file://`/`@file:` table references, current opened-file references, text-backed CSV values, and explicit `CAST` for reliable numeric comparisons.
- The current public roadmap lists SmoothCSV in Browser, an extension system, CSV schema/validation, additional formats such as JSON/Avro/Parquet, and AI integration. These are tracked as strategic post-parity directions unless an official shipping claim moves them into the product baseline.

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
- Split the editor into multiple tab groups arranged side-by-side, stacked, or nested.
- Move tabs between editor groups by drag/drop and command palette actions.
- Sort rows.
- Filter rows.
- Recent filter query history with configurable retention and a shared dropdown.
- Visual condition-builder filtering.
- Search and replace.
- Multi-cell selection/editing.
- Undo/redo.
- Body-first selection where row/column label selection initially excludes headers and can toggle/include headers.
- Swap selected cells with adjacent cells in up/down/left/right directions.
- Command palette.
- Keyboard shortcuts.
- Keyboard shortcut customization.
- Context-aware command/menu/keybinding enablement.
- Theme support, including light/dark/system auto-detect.
- Settings UI with modified indicators and reset-to-default affordances.
- Custom grid/text-editor fonts.
- File format presets.
- File format rules for automatic preset application.
- Manual quote mode and per-cell quote toggling for irregular quoting.
- CSVLint-like validation.
- Query with SQL.
- Compare/diff two CSV files.
- Compare inline character-level changes, merge changes, ignore whitespace/case/columns, align columns by header name, align rows by key column, and compare with clipboard.
- Edit, sort, filter, and search within compare view while keeping the diff live.
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
- Restore from timestamped backup history.
- Quick filters and sorting from column menus.
- Aggregation summaries and aggregate rows.
- Smart sort by text, number, date, and length, including full-width numeral support.
- Read-only mode.
- Transpose rows and columns.
- Auto-refresh for external file modifications.
- Open hyperlinks from cells.
- Start page with quick actions and recent files.
- Persistent window size, position, and maximized state.
- Auto-update checks with manual mode.
- Internationalization for English, Japanese, French, Spanish, Italian, Simplified Chinese, Portuguese, Russian, German, and Korean.

### 5.3 Official Docs Feature Areas

CLI:

- `smoothcsv [options] [path ...]`
- Open single/multiple files.
- Open stdin.
- CLI `--diff old.csv new.csv` launches Compare view.
- `--goto` row/cell/range selection.
- `--diff` side-by-side compare.
- Delimiter, quote, quote mode, encoding, line ending, header rows, header columns overrides.
- New window/reuse window.
- Wait for close.
- Version/help.
- Exit codes.

URL:

- `smoothcsv://file/path/to/file.csv[:location][?params]`
- `opencsv://file/path/to/file.csv[:location][?params]`
- Row/cell/range locations.
- Delimiter, quote, quote mode, encoding, line ending query params.

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
- Restore from backup history.
- Quote indicators.
- Copy/export formats.
- Read-only mode.
- Auto-refresh.
- Aggregation status.
- Automated tests and performance benchmarks.

### 6.2 Can Ship After Product Parity v1

- SmoothCSV-in-browser parity beyond a local prototype.
- Extension marketplace/system.
- AI-assisted transforms.
- CSV schema authoring beyond validation/profile basics.
- Additional structured formats beyond CSV/TSV/XLSX import/export, such as JSON, Avro, and Parquet.
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
- Quote mode: minimal, always, never, manual, unknown.
- Quote override rules by row type (header/body), column, and individual cell when the source file contains irregular quoting.
- Encoding: UTF-8, UTF-8 BOM, UTF-16 LE/BE, Shift-JIS, GBK, Windows-1252, legacy single-byte fallback.
- Line ending: LF, CRLF, CR.
- Header rows.
- Header columns.
- Sparse/metadata-leading rows without over-promoting weak header signals.

Acceptance:

- Detection result is visible in status/format panel.
- User can override detection.
- User can change delimiter, quote policy, line ending, header interpretation, and other save-format policies without silent data loss.
- Saving preserves detected values unless overridden.
- Files with mixed quote patterns round-trip without silently normalizing quote state.
- Per-cell quote toggling is visible and undoable.

### FR-003 Grid View

The app shall render a spreadsheet-like grid with:

- Stable column header row.
- Stable row number gutter.
- Smooth vertical and horizontal scrolling.
- Custom scroll mechanics or canvas rendering that are not limited by browser maximum element dimensions.
- Cell selection.
- Range selection.
- Multi-cell selection.
- Additive, subtractive, and inverted selection through commands and native-feeling mouse gestures.
- Body-first row/column/corner selection that initially excludes headers and can include/toggle headers via Alt/Option or command.
- Stable sticky header behavior.
- Short/ragged rows where missing cells are visually distinct from real empty cells.
- Selected cell/range visibility in light and dark themes.
- Column resize.
- Column reorder.
- Column hide/show with per-file hidden-column persistence.
- Quote indicators for cells that will be quoted on save.
- Scroll position persistence per file.
- Open hyperlink command for URL-like cell values.
- Link audit for recognized URL values, including trusted standard schemes, safe custom protocols, and blocked unsafe protocols.
- Configurable allow-list for user-approved custom link protocols.
- Cell Value panel support for multiline and large cell values, including edit statistics, dirty draft recovery, expandable editing space, and a detached large-cell editor.

Acceptance:

- Header must not be covered by body rows during scroll.
- Header height must remain visually stable.
- Missing cells in rows shorter than the estimated table width must not look like editable empty cells.
- Unsafe protocols such as `javascript:` and `data:` must not be opened from cells.
- Unknown custom protocols must stay blocked until explicitly allowed in Settings.
- Multiline values edited from the Cell Value panel must keep embedded line breaks after apply and save.
- 1M rows must scroll without DOM row explosion.

### FR-004 Editing

The app shall support:

- Edit single cell.
- Start editing by typing printable characters, including IME-safe CJK input behavior.
- Paste tabular clipboard data.
- Paste spreadsheet clipboard data with quoted cells, embedded tabs, and embedded line breaks.
- Fill selection.
- Auto-fill with intelligent pattern detection.
- Multi-cell edit.
- Insert/delete rows.
- Insert/delete columns.
- Move rows/columns through commands and native-feeling drag reordering where available.
- Swap selected cells with adjacent cells in any direction.
- Copy/cut/paste rows and columns.
- Remove duplicate rows and columns.
- Transform selected text cells: trim, uppercase, lowercase, naming/case conversions, full-width/half-width conversions, and join/split horizontally and vertically with a configurable separator.
- Undo/redo.
- Dirty state.
- Read-only mode that disables editing operations.
- Atomic save.
- Save as.
- Auto-backup before destructive in-place save.
- Restore from backup history without overwriting the live file until the user saves.
- Auto-fit column width on edit when enabled.
- Insert line breaks inside cell editor with platform-consistent shortcut.

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
- Exact match highlighting inside cells.
- Case sensitivity toggle.
- Whole cell toggle.
- Regex toggle.
- Search within selection.
- Seed find/replace from selected text while editing.
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
- Natural sort of value checklist entries.
- Visual condition-builder filter with operators such as equals, contains, starts with, blank/non-blank, and between/not-between.
- Selected-cell quick condition templates for common filter predicates.
- Empty/non-empty filters.
- Numeric, text, date, boolean filters.
- Filter chips.
- Recent filter query dropdown with configurable retention or disabled history.
- Saved filter dropdown shared across open editors.
- Clear all filters.
- Open current filter in SQL console.
- Preserve filters across clean external reloads.

Acceptance:

- Header names work when header rows are enabled.
- `c1`, `c2`, ... work when no header row exists.
- `CAST` and `TRY_CAST` work for numeric, date, time, timestamp, and interval cases.
- Filtered row count is visible.
- Reapplying a recent or saved filter restores the query without retyping.

### FR-007 Sort

The app shall support:

- Sort ascending/descending by column.
- Multi-column sort.
- Smart sort modes for text, number, date, and length.
- Full-width numeral support.
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
- SQLCipher-equivalent protection for temporary databases if a temp database is used.
- Query current opened files.
- Query file paths.
- Query unsaved/temporary documents.
- Direct path, `file://`, and `@file:` table references.
- Results grid.
- Copy selected ranges from SQL result grid.
- Export/save query result.
- Query history.
- Per-tab console state and split-pane size persistence.
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
- Row alignment by key column when row order differs.
- Highlight changed cells.
- Inline character-level diff inside modified cells.
- Highlight added/deleted rows.
- Navigate next/previous diff.
- Merge row/cell changes in either direction.
- Ignore whitespace, case, selected columns, or configured ignored columns.
- Align columns by header name when order differs.
- Compare active file with clipboard.
- Edit, sort, filter, and search inside compare view.
- Respect per-file format detection/override.
- Recover gracefully when one side is missing or replaced externally.

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
- `--line-ending`, `--eol`.
- `--header-rows`.
- `--header-columns`.
- `--header`, `--no-header`.
- `--new-window`, `-n`.
- `--reuse-window`, `-r`.
- `--wait`, `-w`.
- `--version`, `-v`.
- `--help`, `-h`.
- `opencsv grep` for cell search from scripts.

Acceptance:

- Invalid CLI usage exits with code 2.
- General runtime errors exit with code 1.
- Success exits with code 0.

### FR-011 URL Protocol

The app shall parse and handle:

- Official SmoothCSV-compatible URLs:
  - `smoothcsv://file/path/to/file.csv`
  - `smoothcsv://file/path/to/file.csv:5`
  - `smoothcsv://file/path/to/file.csv:1:2`
- `opencsv://file/path/to/file.csv`
- `opencsv://file/path/to/file.csv:5`
- `opencsv://file/path/to/file.csv:1:2`
- `opencsv://file/path/to/file.csv:1-5`
- `opencsv://file/path/to/file.csv:1:2-10:20`
- Query params for delimiter, quote, quote mode, encoding, line ending.
- Official URL query aliases such as `quotemode`, `lineending`, `headerrows`, and `headercolumns`, with case-insensitive parameter names.

Acceptance:

- Coordinates are 1-indexed.
- Windows-style paths are accepted.
- URL open can jump to target row/cell/range.
- URL format parameters accept official spelling plus existing dashed/camel/snake aliases.

### FR-012 Command Palette and Shortcuts

The app shall support:

- Command palette.
- Search commands by name.
- Keyboard shortcuts.
- Shortcut customization.
- Optional shortcut context conditions such as document, grid, SQL, profile, settings, diff, and typing contexts.
- Context-aware command enablement using declarative conditions.
- Menus, toolbar buttons, context menus, and shortcuts all dispatch through the same command IDs.
- Grid context menus expose common copy, text transform, quote, cell operation, table, row, column, filter, and sort commands without requiring the command palette.
- Commands for open, save, find, filter, SQL console, format options, diff, sort.

Acceptance:

- Every menu action has a command ID.
- Commands can be invoked by palette or shortcut.
- Customized shortcuts preserve command enablement/context conditions and update menus.

### FR-013 File Format Presets and Rules

The app shall support named presets and automatic rules:

- Delimiter.
- Quote character.
- Quote mode, including minimal, always, never, and manual.
- Header/body/column/cell-specific quote rules.
- Encoding.
- Header rows.
- Header columns.
- Line ending.
- Rule matching by filename/path/content sample.
- Explicit ordering for presets and rules so overlapping automatic rules are deterministic.
- Direct editing for preset names and saved format metadata.

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
- User-supplied schema rules for required columns, primitive type checks, uniqueness, allowed values, regular-expression patterns, and numeric ranges.
- CSVLint/JSON Table Schema style field definitions with common `fields`, `constraints`, `datatype`, `tableSchema`, and `schema` shapes.
- Validation schemas supplied through the Profile UI, URL boot parameters, saved-schema import, or profile API requests.
- Saved schema/rule sets with JSON import and export.
- Visual rule editing for common schema checks without hand-writing JSON.

Acceptance:

- Problems panel lists issue type, row, column, message.
- Clicking issue navigates to location.
- Safe quick fixes can repair straightforward structural issues such as empty generated headers and short ragged rows without guessing semantic data changes.

### FR-015 Clipboard and Export

The app shall support:

- Copy cells as TSV.
- Copy cells as CSV.
- Copy cells as raw values without CSV quoting/escaping.
- Copy selections as Markdown, HTML, JSON, SQL, and LaTeX.
- Copy row/column.
- Paste spreadsheet ranges.
- Parse quoted spreadsheet TSV paste payloads, including embedded tabs, escaped quotes, and embedded line breaks.
- Export filtered rows.
- Export SQL result.
- Open and save `.xlsx` files at least as import/export conversion.

Acceptance:

- Copy/paste round-trips with spreadsheet apps for simple ranges.

### FR-016 Aggregation and Profiling

The app shall support:

- Selection summary in status bar: sum, average, minimum, maximum, count.
- Switch/copy the active status-bar aggregate value.
- Add aggregate rows to the grid.
- Column profile basics: type, null ratio, distinct count, min/max for numeric/date-like values.

Acceptance:

- Selecting numeric cells updates aggregate metrics without blocking scroll.
- Aggregate rows can be removed without changing source data.

### FR-020 App Shell, Settings, and Recent Files

The app shall support:

- Start page with New File, Open, Settings, and recent files.
- Recent files display with file name as primary label and path as secondary label.
- Missing recent files are cleaned up automatically.
- Split editor area into multiple tab groups arranged horizontally, vertically, or nested.
- Move tabs between groups by drag/drop, tab context menu, and command palette actions.
- Switch focus between editor groups by command.
- Persistent window size, position, and maximized state.
- Title/menu entry points for Settings and Keyboard Shortcuts.
- Settings search, modified indicators, and reset-to-default controls.
- Automatic update checks plus a manual update mode.
- Offline operation for all editing features.

Acceptance:

- Relaunch restores window geometry and recent file list.
- Splitting, moving tabs, and switching editor groups preserves open documents and selections.
- Settings changed from defaults are visibly marked and can be reset.
- Disabling automatic update checks prevents startup update network calls.

### FR-017 Transpose

The app shall support transposing rows and columns.

Acceptance:

- Transpose works on small and medium files.
- Transpose can apply to the full table or to one rectangular selected range when the transposed target fits in the current grid.
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
- Million-row scrolling must not depend on DOM spacer dimensions that hit Safari/Chromium element-size caps.
- Full row index for 100 MB: under 2 seconds.
- Find over 100 MB: progress visible, cancellable.
- SQL query over 100 MB: progress visible, cancellable.
- Startup should defer heavy optional surfaces and render a usable start/open surface before loading SQL, diff, settings, and extension-like modules.

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
- Timestamped backups can be listed and restored into an undoable dirty document state.
- Crash recovery for dirty documents.
- External file change detection.
- Clear warning before overwriting changed external file.

### NFR-004 Security and Privacy

Targets:

- Local file contents stay local.
- No telemetry by default.
- No external communication except explicit update checks/downloads and user-initiated links.
- SQL temporary storage must not leave plain-text copies of sensitive data when an encrypted temp database mode is available.
- Browser prototype labels temporary uploads.
- Desktop shell uses OS-native file access.

### NFR-005 Cross-Platform

Targets:

- macOS first.
- Windows and Linux supported by architecture and packaging plan.
- CLI behavior consistent across OS path conventions.
- Build artifacts eventually cover macOS universal, Windows x64/ARM64 installer and portable forms, and Linux `.deb`, `.rpm`, AppImage, and AUR-style distribution.

## 9. Information Architecture

Primary surfaces:

- Start page: quick actions, recent files, update/license/support surfaces.
- Workspace/sidebar: opened files, diff sessions, SQL consoles.
- Toolbar: open, save, find, filter, sort, command palette.
- Grid view: table editing.
- Inspector: selected column/cell profile and metadata.
- Status bar: row/column counts, encoding, delimiter, header, dirty state.
- Problems panel: validation issues.
- SQL console: query editor and results.
- Diff view: side-by-side comparison.
- Settings: app, files, grid editor, text editor, keybindings, update, theme, and file-format presets.

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
- Original quote map / quote-rule overrides when manual or mixed quote mode is active.
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
- Save/backup/restore.
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
| A-004B | Restore From Backup lists timestamped backups and restores one into the session | API/UI test |
| A-005 | Browser Open shows file picker | Browser test with file chooser |
| A-006 | Direct path open works | API/CLI test |
| A-007 | Windowed row loading works | Unit/API test |
| A-008 | Filter SQL works | Unit/API/UI test |
| A-009 | SQL console works | Unit/API/UI test |
| A-010 | CLI options match spec | CLI golden tests |
| A-011 | URL parser handles locations/params | Unit tests |
| A-012 | Diff rejects invalid arg counts | CLI tests |
| A-013 | Diff highlights changed cells | UI test |
| A-014 | Mixed quote patterns round-trip | Fixture save test with per-cell/header/body/column quote assertions |
| A-015 | Missing cells in short rows are visually distinct | Browser screenshot/pixel assertion and ragged-row fixture |
| A-016 | Command enablement is context-aware | Command registry and UI state tests |
| A-017 | Keybinding customization preserves conditions | Settings/keybinding integration tests |
| A-018 | Compare can align rows by key and merge both directions | Diff engine and Browser tests |
| A-019 | SQL Console supports SQLite path refs and secure temp storage | SQL integration tests |
| A-020 | Recent/start/settings shell behavior persists across relaunch | Desktop shell smoke tests |

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

- In-memory SQLite engine integration.
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

### M8 Product Shell and Settings Parity

- Start page with recent files and quick actions.
- Settings and keybinding customization UI.
- Persistent window geometry.
- Update settings.
- File format preset/rule management UI.

### M9 Fidelity and Advanced Grid Parity

- Manual quote mode and per-cell quote map.
- Auto-fill and up/down/left/right swap cells.
- Aggregate rows.
- Smart sort by date/length and full-width numerals.
- Compare key-column row alignment and bidirectional merge.

## 15. Open Product Decisions

- Desktop shell: Wails vs Tauri sidecar vs custom native shell.
- SQL engine: in-memory SQLite row loading is the first parity path; streaming SQLite virtual tables, DuckDB, or DataFusion-like integration remain open for 1 GB-scale query performance.
- Encoding library choice for Shift-JIS/GBK/Windows-1252.
- Whether browser prototype should support save-back-to-original through File System Access API where available.
- Exact UI shape of command palette and menu system.
- Whether opencsv should expose an extension API in v1 or first stabilize command/context APIs internally.
- How much of SmoothCSV's license/update shell should be mirrored versus replaced with an open-source-native support/about/update model.
