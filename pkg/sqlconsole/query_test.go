package sqlconsole

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"
)

func TestExecuteCurrentSelectWhereOrderLimit(t *testing.T) {
	source := Source{
		Name:    "current",
		Path:    "/tmp/current.csv",
		Headers: []string{"id", "name", "score"},
		Records: [][]string{
			{"id", "name", "score"},
			{"1", "Ada", "9"},
			{"2", "Linus", "7"},
			{"3", "Grace", "10"},
		},
	}
	result, err := Execute(`SELECT name, score FROM current WHERE CAST(score AS REAL) >= 8 ORDER BY CAST(score AS REAL) DESC LIMIT 1`, fixedResolver(source), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRows != 2 || !result.Truncated {
		t.Fatalf("total/truncated = %d/%v, want 2/true", result.TotalRows, result.Truncated)
	}
	if len(result.Rows) != 1 || result.Rows[0][0] != "Grace" || result.Rows[0][1] != "10" {
		t.Fatalf("unexpected rows: %#v", result.Rows)
	}
}

func TestExecuteStreamsSourceIntoSQLite(t *testing.T) {
	visited := 0
	source := Source{
		Name:             "current",
		Path:             "/tmp/current.csv",
		Headers:          []string{"id", "name", "score"},
		EstimatedColumns: 3,
		HasHeader:        true,
		Stream: func(ctx context.Context, handle func(row int, record []string) error) (int, error) {
			records := [][]string{
				{"id", "name", "score"},
				{"1", "Ada", "9"},
				{"2", "Linus", "7"},
				{"3", "Grace", "10"},
			}
			for row, record := range records {
				if err := ctx.Err(); err != nil {
					return row, err
				}
				visited++
				if err := handle(row, record); err != nil {
					return row, err
				}
			}
			return len(records), nil
		},
	}

	result, err := ExecuteContext(context.Background(), `SELECT name FROM current WHERE CAST(score AS REAL) >= 9 ORDER BY id`, fixedResolver(source), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Streamed {
		t.Fatal("result streamed = false, want true")
	}
	if visited != 4 {
		t.Fatalf("visited rows = %d, want 4 including header", visited)
	}
	if result.TotalRows != 2 || len(result.Rows) != 2 || result.Rows[0][0] != "Ada" || result.Rows[1][0] != "Grace" {
		t.Fatalf("streamed rows = %#v total=%d, want Ada/Grace", result.Rows, result.TotalRows)
	}
}

func TestExecuteReportsSourceLoadProgress(t *testing.T) {
	source := Source{
		Name:             "current",
		Path:             "/tmp/current.csv",
		Headers:          []string{"id", "name"},
		EstimatedColumns: 2,
		HasHeader:        true,
		Stream: func(ctx context.Context, handle func(row int, record []string) error) (int, error) {
			records := [][]string{
				{"id", "name"},
				{"1", "Ada"},
				{"2", "Linus"},
				{"3", "Grace"},
			}
			for row, record := range records {
				if err := ctx.Err(); err != nil {
					return row, err
				}
				if err := handle(row, record); err != nil {
					return row, err
				}
			}
			return len(records), nil
		},
	}
	var events []ProgressEvent
	result, err := ExecuteContext(context.Background(), `SELECT COUNT(*) AS rows FROM current`, fixedResolver(source), Options{
		ProgressInterval: 2,
		ReportProgress: func(event ProgressEvent) {
			events = append(events, event)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Progress == nil {
		t.Fatalf("progress = nil, want source loading summary")
	}
	if result.Progress.LoadedRows != 3 || result.Progress.Sources != 1 {
		t.Fatalf("progress = %#v, want three loaded rows from one source", result.Progress)
	}
	if len(events) != 2 {
		t.Fatalf("events = %#v, want checkpoint and done", events)
	}
	if events[0].Phase != "load" || events[0].Source != "/tmp/current.csv" || events[0].Rows != 2 || events[0].Done {
		t.Fatalf("first event = %#v, want load checkpoint at two rows", events[0])
	}
	if events[1].Phase != "load" || events[1].Rows != 3 || !events[1].Done {
		t.Fatalf("second event = %#v, want load done at three rows", events[1])
	}
}

func TestExecuteSupportsCColumnsAndFileRef(t *testing.T) {
	source := Source{
		Path:    "/tmp/file.csv",
		Headers: []string{"c1", "c2"},
		Records: [][]string{
			{"1", "Ada"},
			{"2", "Linus"},
		},
	}
	result, err := Execute(`SELECT c2 FROM "@file:/tmp/file.csv" WHERE c1 = 2`, fixedResolver(source), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Source != "/tmp/file.csv" {
		t.Fatalf("source = %q", result.Source)
	}
	if len(result.Rows) != 1 || result.Rows[0][0] != "Linus" {
		t.Fatalf("unexpected rows: %#v", result.Rows)
	}
}

func TestExecuteSupportsFileURLRef(t *testing.T) {
	path := "/tmp/file url.csv"
	source := Source{
		Path:    path,
		Headers: []string{"id", "name"},
		Records: [][]string{
			{"id", "name"},
			{"1", "Ada"},
			{"2", "Grace"},
		},
	}
	fileURL := (&url.URL{Scheme: "file", Path: path}).String()
	result, err := Execute(`SELECT name FROM "`+fileURL+`" WHERE id = 2`, func(ref string) (Source, error) {
		if ref != path {
			t.Fatalf("normalized file URL ref = %q, want %q", ref, path)
		}
		return source, nil
	}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Source != path {
		t.Fatalf("source = %q, want %q", result.Source, path)
	}
	if len(result.Rows) != 1 || result.Rows[0][0] != "Grace" {
		t.Fatalf("unexpected file URL rows: %#v", result.Rows)
	}
}

func TestExecuteAggregates(t *testing.T) {
	source := Source{
		Path:    "/tmp/current.csv",
		Headers: []string{"id", "region", "amount"},
		Records: [][]string{
			{"id", "region", "amount"},
			{"1", "APAC", "10.5"},
			{"2", "EU", "7"},
			{"3", "APAC", "2.5"},
		},
	}
	result, err := Execute(`SELECT COUNT(*), SUM(amount), AVG(amount), MIN(CAST(amount AS REAL)), MAX(CAST(amount AS REAL)) FROM current`, fixedResolver(source), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Rows) != 1 {
		t.Fatalf("rows = %#v, want one aggregate row", result.Rows)
	}
	want := []string{"3", "20", "6.666666666666667", "2.5", "10.5"}
	for i, value := range want {
		if result.Rows[0][i] != value {
			t.Fatalf("aggregate[%d] = %q, want %q; row=%#v", i, result.Rows[0][i], value, result.Rows[0])
		}
	}
}

func TestExecuteGroupByAggregates(t *testing.T) {
	source := Source{
		Path:    "/tmp/current.csv",
		Headers: []string{"id", "region", "amount"},
		Records: [][]string{
			{"id", "region", "amount"},
			{"1", "APAC", "10.5"},
			{"2", "EU", "7"},
			{"3", "APAC", "2.5"},
		},
	}
	result, err := Execute(`SELECT region, COUNT(*) AS rows, SUM(amount) AS total FROM current GROUP BY region ORDER BY total DESC`, fixedResolver(source), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRows != 2 || len(result.Rows) != 2 {
		t.Fatalf("rows = %#v total=%d, want two groups", result.Rows, result.TotalRows)
	}
	if result.Rows[0][0] != "APAC" || result.Rows[0][1] != "2" || result.Rows[0][2] != "13" {
		t.Fatalf("first group = %#v, want APAC/2/13", result.Rows[0])
	}
	if result.Rows[1][0] != "EU" || result.Rows[1][2] != "7" {
		t.Fatalf("second group = %#v, want EU/7", result.Rows[1])
	}
}

func TestExecuteSupportsSQLiteExpressionsJoinAndHaving(t *testing.T) {
	source := Source{
		Path:    "/tmp/orders.csv",
		Headers: []string{"id", "region", "amount"},
		Records: [][]string{
			{"id", "region", "amount"},
			{"1", "APAC", "10"},
			{"2", "EU", "7"},
			{"3", "APAC", "3"},
		},
	}
	regions := Source{
		Path:    "/tmp/regions.csv",
		Headers: []string{"region", "label"},
		Records: [][]string{
			{"region", "label"},
			{"APAC", "Asia Pacific"},
			{"EU", "Europe"},
		},
	}
	result, err := Execute(
		`SELECT r.label, SUM(CAST(o.amount AS REAL)) * 1.1 AS gross
		 FROM current AS o
		 JOIN "/tmp/regions.csv" AS r ON r.region = o.region
		 GROUP BY r.label
		 HAVING gross > 10
		 ORDER BY gross DESC`,
		multiResolver(map[string]Source{
			"current":          source,
			"/tmp/regions.csv": regions,
		}),
		Options{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRows != 1 || len(result.Rows) != 1 {
		t.Fatalf("rows = %#v total=%d, want one group", result.Rows, result.TotalRows)
	}
	if result.Columns[0] != "label" || result.Columns[1] != "gross" {
		t.Fatalf("columns = %#v", result.Columns)
	}
	if result.Rows[0][0] != "Asia Pacific" || result.Rows[0][1] != "14.3" {
		t.Fatalf("joined aggregate row = %#v, want Asia Pacific/14.3", result.Rows[0])
	}
}

func TestExecuteRejectsUnsupportedSQL(t *testing.T) {
	if _, err := Execute("DELETE FROM current", fixedResolver(Source{}), Options{}); err == nil {
		t.Fatal("expected unsupported query error")
	}
	if _, err := Execute("SELECT *", fixedResolver(Source{}), Options{}); err == nil {
		t.Fatal("expected missing FROM error")
	}
}

func TestExecuteReturnsSQLErrorLocation(t *testing.T) {
	source := Source{
		Name:    "current",
		Path:    "/tmp/current.csv",
		Headers: []string{"id", "name"},
		Records: [][]string{
			{"id", "name"},
			{"1", "Ada"},
		},
	}
	_, err := Execute("SELECT name,\nscroe FROM current", fixedResolver(source), Options{})
	if err == nil {
		t.Fatal("expected SQL error")
	}
	var queryErr *QueryError
	if !errors.As(err, &queryErr) {
		t.Fatalf("error type = %T, want QueryError", err)
	}
	if queryErr.Detail.Line != 2 || queryErr.Detail.Column != 1 {
		t.Fatalf("location = %d:%d, want 2:1", queryErr.Detail.Line, queryErr.Detail.Column)
	}
	if queryErr.Detail.Token != "scroe" {
		t.Fatalf("token = %q, want scroe", queryErr.Detail.Token)
	}
	if queryErr.Detail.Snippet != "scroe FROM current" || queryErr.Detail.Caret != "^" {
		t.Fatalf("snippet/caret = %q / %q", queryErr.Detail.Snippet, queryErr.Detail.Caret)
	}
}

func TestExecuteSQLErrorLocationSkipsStringLiteralToken(t *testing.T) {
	source := Source{
		Name:    "current",
		Path:    "/tmp/current.csv",
		Headers: []string{"id", "name"},
		Records: [][]string{
			{"id", "name"},
			{"1", "Ada"},
		},
	}
	_, err := Execute("SELECT 'scroe' AS literal,\nscroe FROM current", fixedResolver(source), Options{})
	if err == nil {
		t.Fatal("expected SQL error")
	}
	var queryErr *QueryError
	if !errors.As(err, &queryErr) {
		t.Fatalf("error type = %T, want QueryError", err)
	}
	if queryErr.Detail.Line != 2 || queryErr.Detail.Column != 1 {
		t.Fatalf("location = %d:%d, want 2:1", queryErr.Detail.Line, queryErr.Detail.Column)
	}
	if queryErr.Detail.Token != "scroe" {
		t.Fatalf("token = %q, want scroe", queryErr.Detail.Token)
	}
	if queryErr.Detail.Snippet != "scroe FROM current" || queryErr.Detail.Caret != "^" {
		t.Fatalf("snippet/caret = %q / %q", queryErr.Detail.Snippet, queryErr.Detail.Caret)
	}
}

func TestExecuteSQLErrorLocationPrefersRepeatedNearTokenContext(t *testing.T) {
	source := Source{
		Name:    "current",
		Path:    "/tmp/current.csv",
		Headers: []string{"id", "name"},
		Records: [][]string{
			{"id", "name"},
			{"1", "Ada"},
		},
	}
	_, err := Execute("SELECT name\nFROM current\nWHERE name = 'FROM'\nFROM current", fixedResolver(source), Options{})
	if err == nil {
		t.Fatal("expected SQL error")
	}
	var queryErr *QueryError
	if !errors.As(err, &queryErr) {
		t.Fatalf("error type = %T, want QueryError", err)
	}
	if queryErr.Detail.Line != 4 || queryErr.Detail.Column != 1 {
		t.Fatalf("location = %d:%d, want 4:1", queryErr.Detail.Line, queryErr.Detail.Column)
	}
	if queryErr.Detail.Token != "FROM" {
		t.Fatalf("token = %q, want FROM", queryErr.Detail.Token)
	}
	if queryErr.Detail.Snippet != "FROM current" || queryErr.Detail.Caret != "^" {
		t.Fatalf("snippet/caret = %q / %q", queryErr.Detail.Snippet, queryErr.Detail.Caret)
	}
}

func TestExecuteContextCancelsRecursiveSQLiteQuery(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(5*time.Millisecond, cancel)
	defer timer.Stop()
	defer cancel()

	_, err := ExecuteContext(
		ctx,
		`WITH RECURSIVE cnt(x) AS (VALUES(0) UNION ALL SELECT x+1 FROM cnt WHERE x < 100000000) SELECT max(x) FROM cnt`,
		fixedResolver(Source{}),
		Options{DefaultLimit: 1},
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
}

func fixedResolver(source Source) Resolver {
	return func(ref string) (Source, error) {
		return source, nil
	}
}

func multiResolver(sources map[string]Source) Resolver {
	return func(ref string) (Source, error) {
		if source, ok := sources[ref]; ok {
			return source, nil
		}
		return Source{}, errMissingSource(ref)
	}
}

type errMissingSource string

func (err errMissingSource) Error() string {
	return "missing source " + string(err)
}
