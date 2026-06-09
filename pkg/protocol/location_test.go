package protocol

import (
	"fmt"
	"testing"
)

func TestParseGotoLocations(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want LocationTarget
	}{
		{
			name: "row",
			raw:  "data.csv:5",
			want: LocationTarget{Path: "data.csv", Location: Location{StartRow: 5}},
		},
		{
			name: "cell",
			raw:  "data.csv:10:3",
			want: LocationTarget{Path: "data.csv", Location: Location{StartRow: 10, StartCol: 3}},
		},
		{
			name: "row range",
			raw:  "data.csv:1-100",
			want: LocationTarget{Path: "data.csv", Location: Location{StartRow: 1, EndRow: 100}},
		},
		{
			name: "cell range",
			raw:  "data.csv:1:1-10:5",
			want: LocationTarget{Path: "data.csv", Location: Location{StartRow: 1, StartCol: 1, EndRow: 10, EndCol: 5}},
		},
		{
			name: "windows path",
			raw:  `C:\Users\me\data.csv:2:4`,
			want: LocationTarget{Path: `C:\Users\me\data.csv`, Location: Location{StartRow: 2, StartCol: 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGoto(tt.raw)
			if err != nil {
				t.Fatalf("ParseGoto returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("ParseGoto(%q) = %#v, want %#v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestParseSmoothCSVURL(t *testing.T) {
	target, err := ParseSmoothCSVURL("smoothcsv://file/tmp/data.csv:1:2-10:5?delimiter=tab&quote=double&quote-mode=minimal&encoding=utf8&line-ending=crlf&header-rows=1&header-columns=2")
	if err != nil {
		t.Fatal(err)
	}
	if target.Path != "/tmp/data.csv" {
		t.Fatalf("path = %q, want /tmp/data.csv", target.Path)
	}
	if target.Location != (Location{StartRow: 1, StartCol: 2, EndRow: 10, EndCol: 5}) {
		t.Fatalf("location = %#v", target.Location)
	}
	if target.Format.Delimiter != "tab" || target.Format.Quote != "double" || target.Format.QuoteMode != "minimal" || target.Format.Encoding != "utf8" || target.Format.LineEnding != "crlf" {
		t.Fatalf("unexpected format params: %#v", target.Format)
	}
	if target.Format.HeaderRows == nil || *target.Format.HeaderRows != 1 {
		t.Fatalf("header rows = %#v, want 1", target.Format.HeaderRows)
	}
	if target.Format.HeaderCols == nil || *target.Format.HeaderCols != 2 {
		t.Fatalf("header columns = %#v, want 2", target.Format.HeaderCols)
	}
}

func TestParseOpenCSVURLNativeScheme(t *testing.T) {
	target, err := ParseSmoothCSVURL("opencsv://file/tmp/data.csv:5?delimiter=pipe&quotemode=always")
	if err != nil {
		t.Fatal(err)
	}
	if target.Path != "/tmp/data.csv" {
		t.Fatalf("path = %q, want /tmp/data.csv", target.Path)
	}
	if target.Location != (Location{StartRow: 5}) {
		t.Fatalf("location = %#v", target.Location)
	}
	if target.Format.Delimiter != "pipe" || target.Format.QuoteMode != "always" {
		t.Fatalf("format = %#v", target.Format)
	}
	if !IsCSVURL("opencsv://file/tmp/data.csv") || !IsCSVURL("smoothcsv://file/tmp/data.csv") {
		t.Fatal("expected both opencsv and smoothcsv schemes to be recognized")
	}
	if IsCSVURL("file:///tmp/data.csv") {
		t.Fatal("file scheme should not be treated as app CSV URL")
	}
}

func TestParseSmoothCSVURLOfficialCaseInsensitiveFormatAliases(t *testing.T) {
	target, err := ParseSmoothCSVURL("smoothcsv://file/tmp/data.csv?Delimiter=pipe&Quote=single&QuoteMode=always&Encoding=UTF-8&LineEnding=crlf&HeaderRows=1&HeaderColumns=1&NoHeader=false")
	if err != nil {
		t.Fatal(err)
	}
	if target.Format.Delimiter != "pipe" || target.Format.Quote != "single" || target.Format.QuoteMode != "always" || target.Format.Encoding != "UTF-8" || target.Format.LineEnding != "crlf" {
		t.Fatalf("unexpected official format aliases: %#v", target.Format)
	}
	if target.Format.HeaderRows == nil || *target.Format.HeaderRows != 1 {
		t.Fatalf("header rows = %#v, want 1", target.Format.HeaderRows)
	}
	if target.Format.HeaderCols == nil || *target.Format.HeaderCols != 1 {
		t.Fatalf("header columns = %#v, want 1", target.Format.HeaderCols)
	}
	if target.Format.NoHeader == nil || *target.Format.NoHeader {
		t.Fatalf("no header = %#v, want false", target.Format.NoHeader)
	}
}

func TestParseSmoothCSVURLLineEndingAliases(t *testing.T) {
	for _, alias := range []string{
		"lineending",
		"line-ending",
		"lineEnding",
		"line_ending",
		"newline",
		"eol",
		"le",
	} {
		t.Run(alias, func(t *testing.T) {
			target, err := ParseSmoothCSVURL(fmt.Sprintf("smoothcsv://file/tmp/data.csv?%s=crlf", alias))
			if err != nil {
				t.Fatal(err)
			}
			if target.Format.LineEnding != "crlf" {
				t.Fatalf("line ending for %s = %q, want %q", alias, target.Format.LineEnding, "crlf")
			}
		})
	}
}

func TestParseSmoothCSVURLWindowsPath(t *testing.T) {
	target, err := ParseSmoothCSVURL("smoothcsv://file/C:/Users/me/data.csv:12?delimiter=comma")
	if err != nil {
		t.Fatal(err)
	}
	if target.Path != "C:/Users/me/data.csv" {
		t.Fatalf("path = %q, want Windows path without leading slash", target.Path)
	}
	if target.Location != (Location{StartRow: 12}) {
		t.Fatalf("location = %#v", target.Location)
	}
}

func TestParseLocationRejectsInvalidRange(t *testing.T) {
	if _, err := ParseGoto("data.csv:10-1"); err == nil {
		t.Fatal("expected invalid row range error")
	}
	if _, err := ParseGoto("data.csv:10:5-10:1"); err == nil {
		t.Fatal("expected invalid cell range error")
	}
	if _, err := ParseSmoothCSVURL("smoothcsv://file/tmp/data.csv?header-rows=-1"); err == nil {
		t.Fatal("expected invalid format parameter error")
	}
}
