# M1 Verification

M1 goal:

- Open and save CSV.
- Detect encoding, delimiter, header, and line endings.
- Build row offset indexes.
- Show the first screen quickly.
- Create automatic backups on save.
- Preserve original line endings and encoding on save.
- Acceptance target: 100 MB CSV reaches first screen in under 2 seconds.

## Implemented Core Paths

- `OpenDocument` detects file metadata and returns preview rows.
- `BuildRowIndex` scans row boundaries without forcing full record materialization.
- `SaveRecords` writes CSV with the detected delimiter, encoding, and line ending.
- `SaveRecords` can create a byte-for-byte backup before in-place save.
- CLI commands: `open`, `inspect`, `generate`, `save`, and `bench`.
- Local web commands through `opencsv serve`: `/api/open`, `/api/window`, `/api/edit`, and `/api/save`.
- Product workbench commands through local API: `/api/find`, `/api/filter`, `/api/sort`, and `/api/view-clear`.

## Verification Commands

```sh
make test
GOCACHE=$(pwd)/.cache/go-build go run ./cmd/opencsv bench -rows 720000 -cols 12 -out /tmp/opencsv-m1-100mb.csv
GOCACHE=$(pwd)/.cache/go-build go run ./cmd/opencsv serve -addr 127.0.0.1:43255 -web web
```

Observed on 2026-05-26:

```json
{
  "sizeBytes": 112857510,
  "openDurationMs": 0,
  "indexDurationMs": 318,
  "indexedRows": 720001,
  "encoding": {
    "name": "utf-8",
    "hasBom": false,
    "confidence": 0.98
  },
  "lineEnding": {
    "name": "lf",
    "sequence": "\n",
    "confidence": 1
  }
}
```

## Remaining Product Work

The Go core, CLI, local API, and Canvas prototype satisfy the M1 data path. The browser prototype supports normal file picking by uploading a local temporary copy, plus direct absolute-path opening through `opencsv serve`. The desktop product still needs a native shell-level file picker so it can open and save the original path without the browser upload boundary.
