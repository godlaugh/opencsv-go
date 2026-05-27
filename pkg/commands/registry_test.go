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
		"settings.keyboardShortcuts",
		"file.save",
		"file.exportView",
		"file.toggleReadOnly",
		"workbench.toggleTheme",
		"edit.copyWithHeaders",
		"selection.toggleHeaders",
		"edit.copyJson",
		"edit.copyHtml",
		"edit.copySql",
		"edit.copyLatex",
		"edit.find",
		"edit.replace",
		"edit.replaceAll",
		"row.hide",
		"row.moveUp",
		"row.moveDown",
		"column.hide",
		"column.moveLeft",
		"column.moveRight",
		"table.transpose",
		"table.dedupeRows",
		"table.dedupeColumns",
		"text.trim",
		"text.uppercase",
		"text.lowercase",
		"cell.swap",
		"cell.joinHorizontal",
		"cell.splitHorizontal",
		"view.filter",
		"view.filterBuilder",
		"view.sort",
		"view.sql",
		"diff.open",
	} {
		if _, ok := ByID(required); !ok {
			t.Fatalf("missing required command %s", required)
		}
	}
}
