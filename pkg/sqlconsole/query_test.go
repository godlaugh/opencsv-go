package sqlconsole

import "testing"

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
	result, err := Execute(`SELECT name, score FROM current WHERE score >= 8 ORDER BY score DESC LIMIT 1`, fixedResolver(source), Options{})
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
	result, err := Execute(`SELECT COUNT(*), SUM(amount), AVG(amount), MIN(amount), MAX(amount) FROM current`, fixedResolver(source), Options{})
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

func TestExecuteRejectsUngroupedColumnWithAggregate(t *testing.T) {
	source := Source{
		Path:    "/tmp/current.csv",
		Headers: []string{"region", "amount"},
		Records: [][]string{
			{"region", "amount"},
			{"APAC", "10"},
		},
	}
	_, err := Execute(`SELECT region, SUM(amount) FROM current GROUP BY amount`, fixedResolver(source), Options{})
	if err == nil {
		t.Fatal("expected ungrouped column error")
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

func fixedResolver(source Source) Resolver {
	return func(ref string) (Source, error) {
		return source, nil
	}
}
