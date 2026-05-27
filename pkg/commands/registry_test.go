package commands

import "testing"

func TestRegistryHasUniqueIDs(t *testing.T) {
	seen := make(map[string]bool)
	for _, command := range Registry() {
		if command.ID == "" || command.Title == "" {
			t.Fatalf("command missing id/title: %#v", command)
		}
		if command.Description == "" {
			t.Fatalf("command missing description: %#v", command)
		}
		if seen[command.ID] {
			t.Fatalf("duplicate command id %q", command.ID)
		}
		seen[command.ID] = true
	}
	for _, required := range []string{
		"file.open",
		"file.save",
		"file.toggleReadOnly",
		"workbench.toggleTheme",
		"edit.copyWithHeaders",
		"edit.copyJson",
		"edit.copyHtml",
		"edit.copySql",
		"edit.copyLatex",
		"edit.find",
		"edit.replace",
		"edit.replaceAll",
		"table.transpose",
		"view.filter",
		"view.sort",
		"view.sql",
		"diff.open",
	} {
		if _, ok := ByID(required); !ok {
			t.Fatalf("missing required command %s", required)
		}
	}
}
