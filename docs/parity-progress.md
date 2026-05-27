# Product Parity Progress

Last updated: 2026-05-27

This document tracks implementation status against `docs/product-prd.md`.

## Current Verified Scope

| PRD ID | Capability | Status | Evidence |
| --- | --- | --- | --- |
| FR-001 | Open CSV/TSV from path and browser picker/upload | In progress | API tests, Browser DOM verification |
| FR-002 | Detect delimiter, encoding, line ending, header | In progress | `pkg/csvcore` tests |
| FR-002/004 | Minimal quote-policy detection and quote indicators | Started | `pkg/csvcore` tests, Browser Inspector/status verification |
| FR-003 | Canvas grid with stable header/body clipping | In progress | Browser DOM workflow plus JS check |
| FR-003 | Cell value panel and hyperlink opening foundation | Started | Browser Inspector verification |
| FR-004 | Edit cell, dirty state, save with backup | In progress | `pkg/document` and `pkg/appserver` tests |
| FR-004 | Undo/redo and row/column insert/delete foundation | Started | `pkg/document` and `pkg/appserver` tests, Browser command-palette verification |
| FR-004/016 | Multi-cell selection, batch edit, and selection aggregation foundation | Started | `/api/edit-batch`, `pkg/document` and `pkg/appserver` tests, Browser verification |
| FR-013 | Clipboard copy/paste foundation | Started | `pkg/clipboard` tests, command registry, Browser command-palette verification |
| FR-005 | Find text | Started | `/api/find`, toolbar query, Browser verification |
| FR-005 | Replace current/all foundation | Started | `/api/replace`, `pkg/document` tests, Browser verification |
| FR-006 | Basic contains filter | Started | `/api/filter`, toolbar query, Browser verification |
| FR-006 | SQL WHERE row filter subset | Started | `pkg/sqlexpr` tests, `/api/filter mode=sql`, Browser verification |
| FR-006 | Excel-like value checklist filter foundation | Started | `pkg/document`/`pkg/appserver` tests, Browser column menu verification |
| FR-007 | Basic selected-column sort | Started | `/api/sort`, Browser verification |
| FR-006/007 | Column quick filter/sort menu foundation | Started | Browser column-header menu verification |
| FR-009 | CSV diff/compare foundation and visual view | Started | `pkg/csvdiff` tests, `/api/diff`, CLI diff smoke, Browser Diff View verification |
| FR-008 | SQL Console foundation plus aggregate query support | Started | `pkg/sqlconsole` tests, `/api/sql`, Browser verification |
| FR-014/016 | Column profile and CSVLint-style validation foundation | Started | `pkg/dataprofile` tests, `/api/profile`, Browser Profile View verification |
| FR-010/011 | CLI goto, format options, diff argument validation, and SmoothCSV URL parsing foundation | Started | `cmd/opencsv` tests, `pkg/protocol` tests, CLI smoke |
| FR-012 | Command registry and Command Palette foundation | Started | `pkg/commands` tests, `/api/commands`, Browser command palette verification |
| FR-017 | Transpose rows/columns foundation | Started | `pkg/document`/`pkg/appserver`/`pkg/commands` tests, Browser command-palette verification |
| FR-018 | Read-only mode and external file refresh foundation | Started | `pkg/document`/`pkg/appserver` tests, Browser verification |
| NFR-001 | 100 MB first screen under 2s | Passing core benchmark | `opencsv bench` |

## Browser Verification Notes

Verified in the in-app browser on `http://127.0.0.1:43255/`:

- Initial UI shows no fake files.
- Path input can open `/tmp/opencsv-m1-100mb.csv`.
- Workspace shows `opencsv-m1-100mb.csv`.
- Status shows `720,000 rows, 12 columns`.
- Find query `segment_7` reports a match at row 8, column 2.
- Filter query `segment_7` reports `60,000 rows`.
- Sort selected `id` column reports descending sort.
- SQL WHERE filter `where field_01 = 'segment_7' and cast(id as number) > 700000` reports `1,667` rows.
- Replace UI opens a 720,000-row file and replaces the first `segment_7` match with `segment_x`, reporting `Replaced 1 cell`.
- Command registry disables document-scoped commands before open and enables implemented commands after opening a CSV, including SQL and Profile.
- `/api/diff` compares `/tmp/opencsv-diff-left.csv` and `/tmp/opencsv-diff-right.csv`, reporting `changed: 1`, `added: 1`, `deleted: 0`.
- Diff View opens from Compare, accepts left/right paths, runs the same files, shows unchanged/changed/added rows, and displays changed cell `c2: Linus -> Ken`.
- SQL Console opens after loading `/tmp/opencsv-diff-right.csv`; query `SELECT name FROM current WHERE id >= 2 ORDER BY id DESC LIMIT 2` returns `Grace` and `Ken`, status `SQL returned 2 of 2 rows`.
- SQL Console aggregate query on `/tmp/opencsv-sql-agg.csv` returns grouped results for `SELECT region, COUNT(*) AS rows, SUM(amount) AS total FROM current GROUP BY region ORDER BY total DESC`: `APAC / 2 / 13` and `EU / 2 / 12`.
- Profile View opens after loading `/tmp/opencsv-profile.csv`, reports `4 rows`, `4 columns`, `3 issues`, and displays column type/null/distinct/min/max/average/top-value statistics.
- Profile validation detects empty header, duplicate header, and ragged row issues in `/tmp/opencsv-profile.csv`.
- Selecting a Profile column row updates the right-side inspector to that column's type, null count, and distinct count.
- Command Palette opens from the `Cmd` toolbar button, filters commands by text, and executing `Profile` opens the Profile View for `/tmp/opencsv-profile.csv`.
- Command Palette executes `Insert Row`, `Undo`, `Redo`, `Insert Column`, and `Delete Column` on `/tmp/opencsv-structure.csv`; row count changes `2 -> 3 -> 2 -> 3`, column count changes `2 -> 3 -> 2`, and the original file remains unchanged until Save.
- Shift-arrow range selection on `/tmp/opencsv-multicell.csv` selects a 2x2 numeric range and status reports `selected 4 cells, sum 10, avg 2.5, min 1, max 4`.
- Editing that 2x2 range to `9` updates all selected cells through `/api/edit-batch` and status reports `sum 36, avg 9, min 9, max 9`; `Cmd+Z` restores the previous values in one undo step.
- Command Palette exposes `Copy`, `Copy as CSV`, `Copy as Markdown`, and `Copy Raw`.
- Browser plugin clipboard virtualization is unavailable in the current runtime, so system clipboard read/write is not directly verified through Browser; clipboard formatters, paste matrix parsing, and `/api/edit-batch` are covered by automated tests.
- Quote indicator verification on `/tmp/opencsv-quotes.csv`: selecting `comma, here` reports `will quote on save`, selecting `plain` reports `no quote needed`, selecting a leading-space value reports `will quote on save`, and the status/Inspector format summary reports `UTF-8`, `Comma`, `LF`, `Quote minimal`, and `Header row`.
- Header regression verification after wheel scroll: `#frozenGrid` keeps `clientHeight`, CSS height, and bounding height at `36px`.
- Hidden panel regression verification: `#profilePanel[hidden]` computes to `display: none`, so hidden SQL/Diff/Profile layers no longer intercept grid/header clicks.
- Column menu verification on `/tmp/opencsv-column-menu.csv`: clicking the `amount` header opens a menu, `Sort descending` reports `Sorted amount descending`, the numeric sampling path places the row with `group = B` first, clicking the `group` header then offers `Filter: B`, filtering reports `Filtered group matched 1 rows`, clearing reports `Filter/sort cleared`, and `Show blanks` reports `Blank group matched 1 rows` with query `where c2 = ''`.
- Read-only/refresh verification on `/tmp/opencsv-readonly-refresh.csv`: toggling `Read Only` sets status `Read-only`, disables Save, sets `aria-pressed=true`, and pressing Enter on the grid reports `Edit blocked: document is read-only` without opening the editor. Toggling back restores editing controls. After an external file rewrite from 1 data row to 2 data rows, file-status polling auto-refreshes the document and status returns to `2 rows, 2 columns` with `Synced`.
- Transpose verification on `/tmp/opencsv-transpose.csv`: Command Palette search `transpose` shows `Transpose`; executing it changes `2 rows, 2 columns` into `2 rows, 3 columns`, clears header-row interpretation to generic columns, marks the document edited, and `Cmd+Z` restores `2 rows, 2 columns`.
- Cell value panel verification on `/tmp/opencsv-cell-panel.csv`: selecting a URL cell sets Cell Value to `https://example.com/demo`, enables `Open Link`, and exposes the same href. Selecting a non-URL `note` cell disables `Open Link`; editing the Cell Value panel to `long note with spaces updated` enables `Apply`, commits through the normal edit path, leaves the updated value visible, and marks the document edited. In read-only mode, the Cell Value textarea and Apply button are disabled.
- Value checklist filter verification on `/tmp/opencsv-value-filter.csv`: clicking the `segment` column menu loads values `A (2)`, `B (2)`, `C (1)`, and `(blank) (1)`; checking `A` and `B` enables `Apply selected values`; applying reports `Value filter segment matched 4 rows` and keeps the selected column on `segment`.

## Implementation Notes

- `pkg/document` now owns document sessions, edit overlay, window reads, save, find, contains-filter, sort, and view mapping.
- `pkg/document` now includes full-record snapshot history for first-pass undo/redo and materializes row/column structural edits into temporary working copies so the original file is not modified before Save.
- `pkg/document.EditCells` and `/api/edit-batch` apply multi-cell edits as one undoable operation.
- `pkg/appserver` delegates document behavior to `pkg/document` instead of owning session logic directly.
- `csvcore.RecordWindow` now includes `sourceRows` so UI edits can target the original source row even when a filtered/sorted view is active.
- `pkg/sqlexpr` provides the first SQL WHERE expression engine for row filters: identifiers/header names, `c1` columns, comparisons, `AND`/`OR`/`NOT`, `LIKE`, `IS NULL`, arithmetic, `CAST`/`TRY_CAST`, and core string/numeric/conditional functions.
- `pkg/protocol` parses official-style goto locations and `smoothcsv://file/...` URLs, including row/cell/range locations and shared format query parameters.
- Root CLI syntax now accepts official-style `opencsv [options] [path ...]`, `--goto`, `--diff`, `--delimiter`, `--encoding`, header flags, window flags, `--wait`, `--version`, and `--help` while retaining developer subcommands.
- CLI open applies delimiter, encoding, and header-row overrides to `csvcore.OpenDocument`; quote and quote-mode are parsed and preserved in invocation metadata but not yet applied by the CSV reader/writer.
- Find now supports regex/case/whole-cell options in the API layer. Replace current/all stores edits through the document edit overlay and marks the session dirty.
- Toolbar wraps into two rows on narrower workspaces so Find/Replace buttons stay clickable instead of sliding under the inspector.
- `pkg/csvdiff` computes row-level and cell-level compare results using LCS anchors plus changed-row pairing by leading key.
- `/api/diff` exposes diff results for two local paths, and root CLI `--diff` now includes a diff result payload in prototype mode.
- Frontend Diff View provides a side-by-side-ish result table with status, left row/record, right row/record, and changed cell details.
- `pkg/commands` exposes command IDs, labels, descriptions, shortcuts, scope, implementation status, and document requirements. The frontend fetches `/api/commands` and applies disabled states/titles to toolbar buttons.
- Frontend grid supports Shift-arrow/Shift-click range selection, range highlighting, batch edit submission, and numeric selection aggregation in the status bar.
- `pkg/clipboard` supports TSV, CSV, Markdown table, and raw value serialization plus TSV matrix parsing.
- `pkg/csvcore.RequiresQuote` mirrors the current Go writer's minimal quote behavior for delimiter, quote, CR/LF, and leading whitespace detection.
- Frontend copy/paste handlers copy selected ranges as TSV and paste tab-separated matrices through `/api/edit-batch`; Command Palette exposes CSV/Markdown/raw copy variants.
- Frontend Command Palette searches `/api/commands` metadata, respects document requirements, and dispatches implemented UI handlers for open/save/undo/redo/find/replace/row/column/filter/sort/clear/sql/profile/diff.
- Frontend grid renders a small quote indicator on cells that will be quoted on save. Status and Inspector surfaces expose encoding, delimiter, line ending, header mode, quote policy, and selected-cell quote state.
- Frontend hidden overlay panels now explicitly honor `[hidden]` to avoid blocking grid interactions.
- Frontend column headers expose a quick menu for ascending/descending sort, filtering by the selected cell value, blank/non-blank filters, and clearing the active view. Quick sort samples loaded cells to treat obvious numeric columns numerically even before full profiling is available.
- `pkg/document` stores `ReadOnly` session state plus last-known file size/mod time, rejects write operations while read-only, exposes file-status checks, and can refresh a clean document from disk.
- `/api/read-only`, `/api/file-status`, and `/api/refresh` expose read-only and external-change behavior to the frontend.
- Frontend status bar shows editable/read-only and synced/changed-on-disk states. Clean, editable documents auto-refresh when the backing file changes; dirty or read-only documents surface the external-change state instead of overwriting local context.
- `pkg/document.Transpose` materializes a real undoable mutation, fills ragged missing cells with blanks, clears header-row interpretation after transpose, marks the document dirty, and rejects operations above the configured cell-count limit.
- `/api/transpose` exposes transpose with a max-cell guard. `table.transpose` is registered in the command registry and dispatched through the frontend Command Palette.
- Frontend Inspector now includes a Cell Value panel for the active cell. It applies edits through the same `/api/edit` path as grid editing, respects read-only mode, keeps local edits visible after async window refreshes, and detects HTTP(S), `mailto:`, `file://`, and bare-domain links for `Open Link`.
- `pkg/document.ColumnValues` returns top distinct values with counts for a column, and `pkg/document.FilterValues` maps rows by selected values plus optional blank inclusion.
- `/api/column-values` and `/api/filter-values` expose value checklist data and filtering. The column header menu renders a bounded checkbox list and applies selected values without mutating source data.
- `pkg/sqlconsole` provides the first SQL Console execution path: `SELECT`, explicit columns or `*`, `FROM current`, quoted path refs, `@file:` refs, `file://` refs, `WHERE` via `pkg/sqlexpr`, `ORDER BY`, `LIMIT`, `COUNT`, `SUM`, `AVG`, `MIN`, `MAX`, and `GROUP BY`.
- `/api/sql` executes SQL against the current document or a local CSV path. The frontend SQL tab now shows an editor, run button, and result table.
- `pkg/dataprofile` provides first-pass profiling and validation: inferred dominant type, null ratio, distinct count, min/max, numeric sum/average, top values, empty/duplicate headers, ragged rows, and type anomalies.
- `/api/profile` profiles the current document. The frontend Profile tab renders summary metrics, a column profile table, issue table, and inspector integration.

## Known Gaps

- Find/Filter/Sort are first-pass implementations. Column quick actions and value checklist filtering exist, but saved filter chips, multi-column sort, and richer type-aware menus are still missing.
- Replace supports current/all and regex basics, but the UI still needs explicit toggles for case sensitivity, whole-cell, and regex mode.
- Undo/redo and structural row/column editing are correct first-pass implementations, but history is full-record snapshots and must become operation-based or chunk-based before it is suitable for very large files.
- Row/column operations currently run through Command Palette and whole-document materialization; direct grid context menus, multi-row/multi-column operations, and virtualized structural mutation are still missing.
- Multi-cell selection/editing is a keyboard-first foundation; mouse drag selection, fill handle, Excel-specific paste behavior, and richer copy/export surfaces are still missing.
- Cell Value panel and hyperlink opening exist as a foundation, but link auditing, custom protocol handling, context menu integration, and richer large-text editing are still missing.
- Clipboard support is implemented as a foundation, but Browser-level system clipboard verification is blocked by the current Browser tool runtime; Excel-specific paste quirks, copy-to-file export, and explicit copy-format UI buttons remain incomplete.
- Quote indicator support currently reflects minimal save behavior only; explicit quote modes, quote character changes, escape rule controls, and a full format preferences UI are still missing.
- SQL WHERE row filter is implemented as a useful subset, but not yet complete for every documented type/function edge case.
- SQL console is implemented as a useful subset, including basic aggregation, but it is not yet fully SQLite-compatible and does not yet support joins, subqueries, expression projection, HAVING, or query history.
- Diff engine/API/CLI payload and a basic visual Diff View exist, but navigation between diffs and a fully virtualized side-by-side editor remain incomplete.
- Profile/validation exists as a useful foundation, but CSVLint-compatible rule sets, saved schemas, quick fixes, and incremental profiling for very large files remain incomplete.
- Command Palette is implemented as a first-pass command runner, but fuzzy scoring, command history, keybinding customization, and full action coverage remain incomplete.
- Full CLI/URL parity is not implemented yet: root CLI parsing/opening exists, but native desktop launch, wait-for-close, protocol handler registration, and actual diff view wiring remain.
- CLI/URL parsing exists as a reusable package, but it is not wired into a desktop protocol handler yet.
- Read-only and external file refresh foundations exist in local app mode, but OS-level file watching, conflict resolution UI for dirty documents, and user-configurable auto-refresh policy are still missing.
- Transpose is implemented as a whole-document materialized operation with a cell-count guard; selection-only transpose, large-file streaming transpose, and a richer confirmation UI are still missing.
- Native desktop shell is not implemented yet.
- Browser screenshot capture is timing out in the current tool runtime, so visual verification is currently DOM/API based.
