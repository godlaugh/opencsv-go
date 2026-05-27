package csvcore

const DefaultSampleBytes = 128 * 1024

type PreviewResult struct {
	Path        string     `json:"path"`
	SizeBytes   int64      `json:"sizeBytes"`
	Dialect     Dialect    `json:"dialect"`
	Records     [][]string `json:"records"`
	RowsScanned int        `json:"rowsScanned"`
}

func OpenPreview(path string, maxRows int) (PreviewResult, error) {
	if maxRows <= 0 {
		maxRows = 50
	}
	document, err := OpenDocument(path, OpenOptions{PreviewRows: maxRows, IndexRows: maxRows})
	if err != nil {
		return PreviewResult{}, err
	}

	return PreviewResult{
		Path:        path,
		SizeBytes:   document.Metadata.SizeBytes,
		Dialect:     document.Metadata.Dialect,
		Records:     document.Preview,
		RowsScanned: document.RowsScanned,
	}, nil
}
