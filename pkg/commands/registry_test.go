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
		"settings.open",
		"settings.keyboardShortcuts",
		"settings.exportShortcuts",
		"settings.importShortcuts",
		"file.save",
		"file.saveAs",
		"file.exportView",
		"file.exportXlsx",
		"file.toggleReadOnly",
		"format.savePreset",
		"format.applyPreset",
		"format.saveRule",
		"format.exportPresets",
		"format.importPresets",
		"format.applyQuotePolicy",
		"format.exportQuoteRules",
		"format.importQuoteRules",
		"format.toggleHeaderQuote",
		"format.toggleBodyQuote",
		"format.toggleColumnQuote",
		"format.toggleCellQuote",
		"format.quoteSelection",
		"format.noQuoteSelection",
		"format.clearSelectionQuote",
		"workbench.toggleTheme",
		"edit.copyWithHeaders",
		"edit.copyRow",
		"edit.copyRowWithHeaders",
		"edit.copyColumn",
		"edit.copyColumnWithHeader",
		"selection.toggleHeaders",
		"selection.addRange",
		"selection.subtractRange",
		"selection.invert",
		"selection.clearRanges",
		"table.toggleHeaderRow",
		"selection.cycleAggregate",
		"selection.copyAggregate",
		"selection.insertAggregateRow",
		"edit.copySelectedFormat",
		"edit.copyJson",
		"edit.copyHtml",
		"edit.copySql",
		"edit.copyLatex",
		"cell.fillDown",
		"cell.fillRight",
		"cell.fillUp",
		"cell.fillLeft",
		"edit.find",
		"edit.replacePreview",
		"edit.replace",
		"edit.replaceAll",
		"row.cut",
		"row.paste",
		"row.hide",
		"row.moveUp",
		"row.moveDown",
		"column.cut",
		"column.paste",
		"column.hide",
		"column.moveLeft",
		"column.moveRight",
		"table.transpose",
		"table.dedupeRows",
		"table.dedupeColumns",
		"text.trim",
		"text.uppercase",
		"text.lowercase",
		"text.snakeCase",
		"text.kebabCase",
		"text.camelCase",
		"text.pascalCase",
		"text.titleCase",
		"text.fullWidth",
		"text.halfWidth",
		"cell.swap",
		"cell.joinHorizontal",
		"cell.splitHorizontal",
		"cell.joinVertical",
		"cell.splitVertical",
		"view.filter",
		"view.filterBuilder",
		"view.saveFilter",
		"view.sort",
		"view.saveSort",
		"view.exportSavedSets",
		"view.importSavedSets",
		"view.sql",
		"sql.copyResultCsv",
		"sql.exportResultCsv",
		"diff.open",
	} {
		if _, ok := ByID(required); !ok {
			t.Fatalf("missing required command %s", required)
		}
	}
}
