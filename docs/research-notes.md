# SmoothCSV Research Notes

These notes summarize public information and product signals, not private source code.

## Public Positioning

SmoothCSV presents itself as a powerful CSV editor for Windows, macOS, and Linux. Its public pitch focuses on spreadsheet-like editing, fast large-file opening, find/replace, sort/filter, SQL, multi-cell editing, and keyboard-friendly command workflows.

The `kohii/smoothcsv3` GitHub repository is used for issues and releases rather than source distribution. Its README describes the aspiration as a "VS Code of tabular editors" and lists future direction such as browser support, extensions, schema/validation, more file formats, and AI.

## Technical Signals

Public engineering posts describe:

- Tauri for desktop packaging.
- Rust backend.
- TypeScript, React, Vite, and MobX on the frontend.
- A package split around CSV, grid editor, SQL, commands, settings, i18n, and shortcuts.
- Custom virtual scrolling to avoid browser scroll-height limits on very large tables.

## UX Signals

Interesting user-requested areas:

- Excel-style value checkbox filtering.
- More reliable SQL type behavior, especially numeric comparisons.
- Large-file performance and first-screen responsiveness.
- A stronger editor mental model: command palette plus direct spreadsheet gestures.

## opencsv Strategy

opencsv should not be a clone. The better opportunity is a local-first data workbench:

- Preserve raw CSV fidelity by default.
- Make typed views and SQL explicit instead of silently coercing data.
- Put filtering, profiling, and cleaning directly on column headers.
- Treat large files as the normal case, not a stress case.
- Keep the core open and scriptable.
