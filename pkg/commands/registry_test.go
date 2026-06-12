package commands

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

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
		"file.new",
		"file.open",
		"settings.open",
		"settings.keyboardShortcuts",
		"settings.exportShortcuts",
		"settings.importShortcuts",
		"settings.checkForUpdates",
		"help.aboutSupport",
		"workbench.splitEditor",
		"workbench.toggleSplitOrientation",
		"workbench.focusNextGroup",
		"workbench.moveTabToOtherGroup",
		"workbench.closeSplitEditor",
		"file.save",
		"file.saveAs",
		"file.restoreBackup",
		"file.exportView",
		"file.exportXlsx",
		"file.toggleReadOnly",
		"format.savePreset",
		"format.applyPreset",
		"format.saveRule",
		"format.exportPresets",
		"format.importPresets",
		"format.applyDelimiter",
		"format.applyEncoding",
		"format.applyLineEnding",
		"format.applyQuotePolicy",
		"format.switchWorksheet",
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
		"edit.copyRowCsv",
		"edit.copyRowMarkdown",
		"edit.copyRowJson",
		"edit.copyRowHtml",
		"edit.copyRowSql",
		"edit.copyRowLatex",
		"edit.copyRowRaw",
		"edit.copyColumn",
		"edit.copyColumnWithHeader",
		"edit.copyColumnCsv",
		"edit.copyColumnMarkdown",
		"edit.copyColumnJson",
		"edit.copyColumnHtml",
		"edit.copyColumnSql",
		"edit.copyColumnLatex",
		"edit.copyColumnRaw",
		"selection.toggleHeaders",
		"selection.addRange",
		"selection.subtractRange",
		"selection.invert",
		"selection.clearRanges",
		"table.toggleHeaderRow",
		"selection.cycleAggregate",
		"selection.copyAggregate",
		"selection.insertAggregateRow",
		"selection.insertSumRow",
		"selection.insertAverageRow",
		"selection.insertMinRow",
		"selection.insertMaxRow",
		"selection.insertCountRow",
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
		"row.pasteBelow",
		"row.hide",
		"row.showAll",
		"row.moveUp",
		"row.moveDown",
		"column.cut",
		"column.paste",
		"column.pasteRight",
		"column.hide",
		"column.showAll",
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
		"text.normalizeNFC",
		"text.normalizeNFD",
		"text.normalizeNFKC",
		"text.normalizeNFKD",
		"text.column.trim",
		"text.column.uppercase",
		"text.column.lowercase",
		"text.column.snakeCase",
		"text.column.kebabCase",
		"text.column.camelCase",
		"text.column.pascalCase",
		"text.column.titleCase",
		"text.column.fullWidth",
		"text.column.halfWidth",
		"text.column.normalizeNFC",
		"text.column.normalizeNFD",
		"text.column.normalizeNFKC",
		"text.column.normalizeNFKD",
		"cell.swap",
		"cell.swapUp",
		"cell.swapDown",
		"cell.swapLeft",
		"cell.swapRight",
		"cell.openLink",
		"cell.joinHorizontal",
		"cell.splitHorizontal",
		"cell.joinVertical",
		"cell.splitVertical",
		"view.filter",
		"view.filterBuilder",
		"view.saveFilter",
		"view.sort",
		"column.sortAscending",
		"column.sortDescending",
		"column.sortLengthAscending",
		"column.sortLengthDescending",
		"column.addSortAscending",
		"column.addSortDescending",
		"column.filterBySelectedValue",
		"column.filterBlanks",
		"column.filterNonBlanks",
		"view.saveSort",
		"view.removeLastSortKey",
		"view.exportSavedSets",
		"view.importSavedSets",
		"view.sql",
		"sql.openResult",
		"sql.copyResultCsv",
		"sql.copyResultTsv",
		"sql.copyResultMarkdown",
		"sql.copyResultJson",
		"sql.copyResultHtml",
		"sql.copyResultSql",
		"sql.copyResultLatex",
		"sql.copyResultRaw",
		"sql.copyCell",
		"sql.copyRow",
		"sql.copyRowWithHeaders",
		"sql.copyColumn",
		"sql.copyColumnWithHeader",
		"sql.copySelectionCsv",
		"sql.copySelectionTsv",
		"sql.copySelectionMarkdown",
		"sql.copySelectionJson",
		"sql.copySelectionHtml",
		"sql.copySelectionSql",
		"sql.copySelectionLatex",
		"sql.copySelectionRaw",
		"sql.copyFullResultCsv",
		"sql.copyFullResultTsv",
		"sql.copyFullResultMarkdown",
		"sql.copyFullResultJson",
		"sql.copyFullResultHtml",
		"sql.copyFullResultSql",
		"sql.copyFullResultLatex",
		"sql.copyFullResultRaw",
		"sql.exportResultCsv",
		"sql.exportResultTsv",
		"sql.exportResultMarkdown",
		"sql.exportResultJson",
		"sql.exportResultHtml",
		"sql.exportResultSql",
		"sql.exportResultLatex",
		"sql.exportResultRaw",
		"sql.clearSelection",
		"diff.open",
	} {
		if _, ok := ByID(required); !ok {
			t.Fatalf("missing required command %s", required)
		}
	}
}

func TestOfficialFilterShortcutDefault(t *testing.T) {
	command, ok := ByID("view.filter")
	if !ok {
		t.Fatal("missing view.filter command")
	}
	if command.Shortcut != "Cmd+Shift+F" {
		t.Fatalf("view.filter shortcut = %q, want Cmd+Shift+F", command.Shortcut)
	}
}

func TestFileCloseTabShortcutDefault(t *testing.T) {
	command, ok := ByID("file.closeTab")
	if !ok {
		t.Fatal("missing file.closeTab command")
	}
	if command.Shortcut != "Cmd+W" {
		t.Fatalf("file.closeTab shortcut = %q, want Cmd+W", command.Shortcut)
	}
}

func TestCloseTabShortcutDispatchWiring(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`if (key === "w" && !event.shiftKey && defaultShortcutActive("file.closeTab"))`,
		`if (event.key === "F4" && event.ctrlKey && !event.metaKey && defaultShortcutActive("file.closeTab"))`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing close-tab shortcut wiring %q", snippet)
		}
	}
}

func TestURLLocationReversedRangeNormalizationIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function orderedURLLocation(",
		`if (match) return orderedURLLocation(match[1], match[2], match[3], match[4]);`,
		`if (match) return orderedURLLocation(match[1], 1, match[2], "last");`,
		"Math.min(rowStart, rowEnd)",
		"Math.max(rowStart, rowEnd)",
		"Math.min(colStart, colEnd)",
		"Math.max(colStart, colEnd)",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing URL location reversed-range wiring %q", snippet)
		}
	}
}

func TestPasteBelowAndRightCommandsAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`case "row.pasteBelow":`,
		`await pasteRowFromClipboard({ placement: "after-selection" });`,
		`case "column.pasteRight":`,
		`await pasteColumnFromClipboard({ placement: "after-selection" });`,
		`function selectedRowPasteAfterTarget()`,
		`function selectedColumnPasteRightTarget()`,
		`onCopyReady: () => {`,
		`setStructureClipboard("row", rowSelection.sourceRows, false)`,
		`setStructureClipboard("column", sourceCols, false)`,
		`addCommand("row.pasteBelow");`,
		`addCommand("column.pasteRight");`,
		`document.body.dataset.structurePaste = ` + "`row:${clipboard.cut ? \"cut\" : \"copy\"}:${placement}:${sources.join(\",\")}->${targetRow}`",
		`document.body.dataset.structurePaste = ` + "`column:${clipboard.cut ? \"cut\" : \"copy\"}:${placement}:${sources.join(\",\")}->${targetCol}`",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing paste below/right wiring %q", snippet)
		}
	}
}

func TestRowColumnCopyAsFormatsAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	formats := map[string]string{
		"Csv":      "csv",
		"Markdown": "markdown",
		"Json":     "json",
		"Html":     "html",
		"Sql":      "sql",
		"Latex":    "latex",
		"Raw":      "raw",
	}
	for suffix, format := range formats {
		for _, target := range []string{"Row", "Column"} {
			id := "edit.copy" + target + suffix
			if _, ok := ByID(id); !ok {
				t.Fatalf("missing copy-as command %s", id)
			}
			if !strings.Contains(source, `["`+id+`", "`+format+`"]`) {
				t.Fatalf("missing frontend format map for %s", id)
			}
		}
	}
	for _, snippet := range []string{
		"function normalizeCopyFormat",
		"function copyFormatLabel",
		"copyActiveRow({ format: copyRowFormatCommands.get(commandID) })",
		"copyActiveColumn({ format: copyColumnFormatCommands.get(commandID) })",
		"for (const commandID of copyRowFormatCommands.keys()) addCommand(commandID);",
		"for (const commandID of copyColumnFormatCommands.keys()) addCommand(commandID);",
		"status.textContent = `Copied ${label}: ${shape} as ${copyFormatLabel(format)}`",
		`if (!copied && !includeHeaders && format === "tsv") clearStructureClipboard()`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing row/column copy-as wiring %q", snippet)
		}
	}
}

func TestSpreadsheetHTMLPasteIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function parseHTMLTable",
		"function htmlTableCellText",
		"function pasteHTMLTable",
		`clipboardItemsText(items, "text/html")`,
		`event.clipboardData.getData("text/html")`,
		`document.body.dataset.pasteSource = source`,
		`window.__opencsvParseHTMLTable = parseHTMLTable`,
		`br.replaceWith("\n")`,
		`cell.getAttribute("colspan")`,
		`cell.getAttribute("rowspan")`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing spreadsheet HTML paste wiring %q", snippet)
		}
	}
}

func TestEditorGroupLayoutPersistenceIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`const editorGroupLayoutStorageKey = "opencsv-editor-group-layout"`,
		"let editorGroupLayout = loadEditorGroupLayout()",
		"let editorGroupActiveTabs",
		"function loadEditorGroupLayout()",
		"function saveEditorGroupLayout()",
		"function editorGroupLayoutSignature(",
		"function editorGroupStackSignature()",
		"function tabForEditorGroup(",
		"function moveEditorGroupTabToGroup(",
		"document.body.dataset.editorGroupLayoutRestored",
		"document.body.dataset.editorGroupLayoutSaved",
		"document.body.dataset.editorGroupStacks",
		"saveEditorGroupLayout();",
		"document.body.dataset.editorGroupLayout = editorGroupLayoutSignature()",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing editor group layout persistence wiring %q", snippet)
		}
	}
}

func TestSelectionAggregateStatsAreCachedAndLimited(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"const aggregateScanCellLimit = 5000",
		"let aggregateStatsCacheKey = \"\"",
		"function currentSelectionAggregateStats()",
		"aggregateStatsCache = selectionAggregateStats()",
		"if (count > aggregateScanCellLimit)",
		"document.body.dataset.aggregateLimited = String(Boolean(stats.limited))",
		"document.body.dataset.aggregateScannedCells = String(stats.scannedCells)",
		"document.body.dataset.aggregateElapsedMs = String(stats.elapsedMs)",
		"currentSelectionAggregateStats().count <= 1",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing aggregate cache/limit wiring %q", snippet)
		}
	}
}

func TestOpenFileTabsAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"let openFileTabs = []",
		"function renderOpenFileTabs()",
		"function closeOpenFileTab",
		"async function closeActiveOpenFileTab",
		"function ensureActiveOpenFileTabVisible",
		"function scrollOpenFileTabIntoView",
		"function syncOpenFileTabOverflowState",
		"function openDirtyClosePanel",
		"function closeDirtyClosePanel",
		"async function confirmDirtyCloseDiscard",
		"function upsertOpenFileTab",
		"function captureActiveOpenFileTabState",
		"async function switchOpenFileTab",
		"function pathInputValueForOpenSession",
		"pathInput.value = pathInputValueForOpenSession(session, detail, tab.name)",
		"updateFormatRulePlaceholder(detail)",
		"switchOpenFileTab(tab.documentId, { group: tabGroup(tab) })",
		"document.body.dataset.editorGroupSwitch",
		"row.dataset.editorGroup = tabGroup(tab)",
		"/api/session?id=",
		"/api/close",
		"file-tab-close",
		"file-tab-row active",
		"document.body.dataset.closeTabBlocked",
		"document.body.dataset.closeTabConfirm",
		"document.body.dataset.closeTabDiscarded",
		"closeOpenFileTab(target.documentId, { force: true })",
		"document.body.dataset.openTabOverflow",
		"document.body.dataset.activeOpenTabVisible",
		"document.body.dataset.openTabCount",
		"document.body.dataset.openTabs",
		"document.body.dataset.activeOpenTab",
		"button.dataset.documentId",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing open file tab wiring %q", snippet)
		}
	}
}

func TestSettingsSectionModifiedBadgesAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	styles, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "styles.css"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	css := string(styles)
	for _, snippet := range []string{
		`id="resetSavedViewsButton"`,
		`Reset Views`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing settings reset html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		"function settingsSectionModifiedState",
		"function syncSettingsSectionBadges",
		"function savedViewsSettingsModified",
		"function resetSavedViewsSettings",
		"function quotePolicySettingsModified",
		"document.body.dataset.settingsModifiedSections",
		"document.body.dataset.settingsModifiedSectionCount",
		"document.body.dataset.savedViewsReset",
		`badge.className = "settings-section-state"`,
		`badge.dataset.state = modified ? "modified" : "default"`,
		`"Saved filters and sorts": savedViewsSettingsModified({ filters, recent, sorts })`,
		`"Editor appearance": editorAppearanceModified()`,
		`"File format presets": presets.length > 0`,
		`resetSavedViewsButton.disabled = !sectionState.state["Saved filters and sorts"]`,
		`resetSavedViewsButton.addEventListener("click", resetSavedViewsSettings)`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing settings section modified wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		".settings-section-state",
		`.settings-section-state[data-state="modified"]`,
	} {
		if !strings.Contains(css, snippet) {
			t.Fatalf("missing settings section modified style %q", snippet)
		}
	}
}

func TestDirtyDraftRecoveryIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	styles, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "styles.css"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	css := string(styles)
	for _, snippet := range []string{
		`const dirtyDraftsStorageKey = "opencsv-dirty-drafts"`,
		"function scheduleDirtyDraftSave()",
		"function restoreDirtyDraft(id)",
		`body: JSON.stringify({ documentId, row, col, value })`,
		`label.textContent = "Recovery Drafts"`,
		"document.body.dataset.recoveryDraftRestored",
		"clearDirtyDraft(draftIDBeforeSave)",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing dirty draft recovery wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		".recovery-draft-row",
		".recovery-draft-row .file-button",
	} {
		if !strings.Contains(css, snippet) {
			t.Fatalf("missing dirty draft recovery style %q", snippet)
		}
	}
}

func TestExternalOverwriteGuardIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	html, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	markup := string(html)
	for _, snippet := range []string{
		`id="overwriteExternalChangeButton"`,
		`data-dialog-role="confirm">Overwrite Disk`,
	} {
		if !strings.Contains(markup, snippet) {
			t.Fatalf("missing external overwrite markup %q", snippet)
		}
	}
	for _, snippet := range []string{
		`allowExternalOverwrite: Boolean(options.allowExternalOverwrite)`,
		`document.body.dataset.externalSaveBlocked = "changed"`,
		"overwriteExternalChangeButton.hidden = !conflict",
		"overwriteExternalChangeOnDisk",
		`document.body.dataset.externalOverwrite = "saved"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing external overwrite wiring %q", snippet)
		}
	}
}

func TestImplementedCommandsHaveFrontendDispatch(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	handlers := make(map[string]bool)
	for _, match := range regexp.MustCompile(`case "([^"]+)"\s*:`).FindAllStringSubmatch(source, -1) {
		handlers[match[1]] = true
	}
	for _, match := range regexp.MustCompile(`"(sql\.[^"]+)"\s*:\s*\{\s*action:`).FindAllStringSubmatch(source, -1) {
		handlers[match[1]] = true
	}
	for _, match := range regexp.MustCompile(`\["(edit\.copy(?:Row|Column)[^"]+)",\s*"`).FindAllStringSubmatch(source, -1) {
		handlers[match[1]] = true
	}

	var missing []string
	for _, command := range Registry() {
		if command.Implemented && !handlers[command.ID] {
			missing = append(missing, command.ID)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("implemented commands missing frontend dispatch: %s", strings.Join(missing, ", "))
	}
}

func TestSQLResultFormatsHaveCommandsAndSelectOptions(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}

	source := string(appJS)
	formatMatch := regexp.MustCompile(`const sqlResultFormats = \[([^\]]+)\];`).FindStringSubmatch(source)
	if len(formatMatch) != 2 {
		t.Fatal("missing sqlResultFormats declaration")
	}
	formats := regexp.MustCompile(`"([^"]+)"`).FindAllStringSubmatch(formatMatch[1], -1)
	if len(formats) == 0 {
		t.Fatal("sqlResultFormats is empty")
	}

	commandMap := regexp.MustCompile(`"(sql\.[^"]+)"\s*:\s*\{\s*action:\s*"([^"]+)",\s*format:\s*"([^"]+)"`).FindAllStringSubmatch(source, -1)
	frontendCommands := make(map[string]bool)
	frontendActions := make(map[string]string)
	for _, match := range commandMap {
		frontendCommands[match[1]] = true
		frontendActions[match[1]] = match[2]
	}
	sqlSelect := regexp.MustCompile(`(?s)<select id="sqlResultFormatSelect"[^>]*>(.*?)</select>`).FindStringSubmatch(string(indexHTML))
	if len(sqlSelect) != 2 {
		t.Fatal("missing SQL result format select")
	}
	selectOptions := make(map[string]bool)
	for _, match := range regexp.MustCompile(`<option value="([^"]+)">`).FindAllStringSubmatch(sqlSelect[1], -1) {
		selectOptions[match[1]] = true
	}

	var missing []string
	for _, match := range formats {
		format := match[1]
		suffix := sqlResultCommandSuffixForTest(format)
		for _, prefix := range []string{"sql.copyResult", "sql.copySelection", "sql.copyFullResult", "sql.exportResult"} {
			id := prefix + suffix
			if _, ok := ByID(id); !ok {
				missing = append(missing, id+" registry")
			}
			if !frontendCommands[id] {
				missing = append(missing, id+" frontend")
			}
		}
		if !selectOptions[format] {
			missing = append(missing, format+" select option")
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf("SQL result format coverage missing: %s", strings.Join(missing, ", "))
	}
	if frontendActions["sql.copyFullResultSql"] != "copyFullResult" {
		t.Fatalf("sql.copyFullResultSql action = %q", frontendActions["sql.copyFullResultSql"])
	}
}

func TestDiffResultFiltersAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, id := range []string{
		"diffStatusFilter",
		"diffSearchInput",
		"diffSortSelect",
		"diffSortDirectionButton",
		"diffClearFilterButton",
	} {
		if !strings.Contains(html, `id="`+id+`"`) {
			t.Fatalf("missing diff filter control %s", id)
		}
		if !strings.Contains(source, `#`+id) {
			t.Fatalf("missing diff filter selector %s", id)
		}
	}
	for _, snippet := range []string{
		"filteredDiffRows",
		"sortedDiffRows",
		"compareDiffRows",
		"compareDiffSortMissing",
		"diffRowMatchesStatus",
		"diffRowSearchText",
		"diffStatusFilter.addEventListener(\"change\", applyDiffResultFilters)",
		"diffSearchInput.addEventListener(\"input\", applyDiffResultFilters)",
		"diffSortSelect.addEventListener(\"change\", applyDiffResultFilters)",
		"diffSortDirectionButton.addEventListener(\"click\", toggleDiffSortDirection)",
		"diffClearFilterButton.addEventListener(\"click\", clearDiffResultFilters)",
		"document.body.dataset.diffFilteredRows",
		"document.body.dataset.diffSort",
		"document.body.dataset.diffSortDirection",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing diff filter wiring %q", snippet)
		}
	}
}

func TestDiffURLGotoIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`params.get("diffGoto")`,
		`params.get("diffGotoSide")`,
		"function applyDiffURLGoto",
		"document.body.dataset.diffGotoSide",
		"document.body.dataset.diffGotoRow",
		"document.body.dataset.diffGotoActive",
		"row.dataset.resultIndex",
		"Diff goto ${targetSide} row ${targetRow}.",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing diff URL goto wiring %q", snippet)
		}
	}
}

func TestDiffEditableCellsAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	stylesCSS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "styles.css"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	styles := string(stylesCSS)
	for _, snippet := range []string{
		"function editDiffCell",
		`fetch("/api/diff-edit-cell"`,
		"appendDiffCellEditor",
		"input.dataset.diffEditSide",
		"document.body.dataset.diffEditCell",
		"document.body.dataset.diffEditBackup",
		`"diff.action.editCell"`,
		`"diff.status.editedCell"`,
		"await runDiffCompare()",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing diff editable-cell wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		".diff-cell-editors",
		".diff-cell-editor input",
		"white-space: normal",
	} {
		if !strings.Contains(styles, snippet) {
			t.Fatalf("missing diff editable-cell styles %q", snippet)
		}
	}
}

func TestDiffViewI18nIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, id := range []string{
		"diffIgnoreCaseLabel",
		"diffIgnoreWhitespaceLabel",
		"diffMatchHeadersLabel",
		"diffRightTextModeLabel",
	} {
		if !strings.Contains(html, `id="`+id+`"`) {
			t.Fatalf("missing diff i18n label anchor %s", id)
		}
		if !strings.Contains(source, `#`+id) {
			t.Fatalf("missing diff i18n selector %s", id)
		}
	}
	for _, snippet := range []string{
		`"diff.status.summary"`,
		`"diff.table.changedCells"`,
		`"diff.action.applyRightCellToLeft"`,
		`"diff.action.openLeftRow"`,
		`"diff.status.openedSource": "已从比较打开`,
		`"diff.status.summary": "差异`,
		"applyDiffLocaleStrings",
		"setSelectOptionLabels(diffStatusFilter",
		"localizedDiffStatus(row.status)",
		"appendDiffOpenSourceActions",
		"openDiffSourceRow(\"left\"",
		"document.body.dataset.diffOpenSource",
		"tf(\"diff.status.active\"",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing diff i18n wiring %q", snippet)
		}
	}
}

func TestSQLConsoleI18nIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`"sql.status.returned"`,
		`"sql.status.returnedTimed"`,
		`"sql.status.returnedTimedRows"`,
		`"sql.status.runningRows"`,
		`"sql.context.copySelectionAs"`,
		`"sql.status.returned": "SQL 返回`,
		`"sql.status.returnedTimed": "SQL 返回`,
		`"sql.status.runningRows": "正在加载`,
		"applySQLLocaleStrings",
		"renderSQLHistory()",
		"renderSQLResult(lastSQLResult)",
		`const statusKey = payload.progress?.loadedRows ? "sql.status.returnedTimedRows" : "sql.status.returnedTimed"`,
		"tf(statusKey",
		"readSQLResponse(response)",
		"applySQLProgressSummary(payload.progress)",
		"tf(\"sql.status.selected\"",
		"tf(\"sql.context.copyResultAs\"",
		"t(\"sql.empty.rows\")",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing SQL i18n wiring %q", snippet)
		}
	}
}

func TestFindReplaceI18nIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, id := range []string{
		"commandInput",
		"replaceInput",
		"findCaseButton",
		"findSelectionButton",
		"replacePreviewButton",
	} {
		if !strings.Contains(html, `id="`+id+`"`) {
			t.Fatalf("missing find/replace control %s", id)
		}
	}
	for _, snippet := range []string{
		`"find.placeholder"`,
		`"find.placeholder": "查找/筛选`,
		`"replace.status.preview": "预览`,
		"applyFindLocaleStrings",
		`commandInput.placeholder = t("find.placeholder")`,
		`findSelectionButton.textContent = t("find.selectionShort")`,
		`t("find.status.enterQuery")`,
		`tf("find.status.foundTimed"`,
		`tf("replace.status.noReplacements"`,
		`tf("replace.status.preview"`,
		"normalizeServerFindRanges",
		"normalizeServerFindRange",
		"formatServerMatchRanges",
		"replaceChangePreview",
		"match.ranges",
		"match.range",
		"flatMap((match) => normalizeServerFindRanges",
		"document.body.dataset.replacePreviewRanges",
		"document.body.dataset.replaceRanges",
		"document.body.dataset.replaceChange",
		"document.body.dataset.findHighlight",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing find/replace i18n wiring %q", snippet)
		}
	}
}

func TestLongOperationCancellationIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, id := range []string{
		"cancelFindButton",
		"cancelSqlButton",
		"cancelProfileButton",
	} {
		if !strings.Contains(html, `id="`+id+`"`) {
			t.Fatalf("missing cancellation control %s", id)
		}
	}
	for _, snippet := range []string{
		"let activeFindController = null",
		"let findProgressTimer = 0",
		"let activeSQLController = null",
		"let activeProfileController = null",
		"function startFindProgress",
		"function updateFindProgressStatus",
		"function stopFindProgressTimer",
		"function cancelFindRequest",
		"function startSQLProgress",
		"function cancelSQLRequest",
		"function updateSQLProgressStatus",
		"function readSQLResponse",
		"function applySQLProgressEvent",
		"document.body.dataset.findElapsedMs",
		"document.body.dataset.findStreamed",
		"document.body.dataset.sqlElapsedMs",
		"document.body.dataset.sqlProgressRows",
		"document.body.dataset.sqlProgressEvents",
		"function startProfileProgress",
		"function cancelProfileRequest",
		"function readProfileResponse",
		"function applyProfileProgressEvent",
		"streamProgress: true",
		"document.body.dataset.profileProgressRows",
		"document.body.dataset.profileProgressEvents",
		"signal: controller.signal",
		`params.get("find")`,
		`params.get("findRun")`,
		"document.body.dataset.findUrlBoot",
		`document.body.dataset.findProgress = "canceled"`,
		`document.body.dataset.sqlProgress = "canceled"`,
		`document.body.dataset.profileProgress = "canceled"`,
		`"find.status.runningElapsed": "正在查找`,
		`"find.status.canceledElapsed": "查找已在`,
		`"find.status.noMatchesTimed": "用时`,
		`"find.status.foundTimed": "找到`,
		`"sql.status.runningElapsed": "正在运行 SQL... {{seconds}} 秒"`,
		`"sql.status.canceledElapsed": "SQL 已在 {{seconds}} 秒后取消。"`,
		`"sql.status.returnedTimed": "SQL 返回 {{rows}} / {{total}} 行，用时 {{seconds}} 秒"`,
		`"profile.status.canceledElapsed": "概况已在`,
		`"profile.status.runningRows": "正在执行`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing long-operation cancellation wiring %q", snippet)
		}
	}
}

func TestFilterURLBootAndStreamingMarkerIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`params.get("filter")`,
		`params.get("filterRun")`,
		"document.body.dataset.filterUrlBoot",
		"document.body.dataset.filterStreamed",
		"document.body.dataset.columnValuesStreamed",
		"document.body.dataset.valueFilterStreamed",
		`"/api/filter-values"`,
		`/api/column-values?id=`,
		"await filterCurrentQuery()",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing filter URL boot/streaming wiring %q", snippet)
		}
	}
}

func TestConditionBuilderI18nIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, snippet := range []string{
		`data-i18n="condition.title"`,
		`data-i18n="condition.template.selectedEquals"`,
		`data-i18n="condition.operator.contains"`,
		`data-i18n="condition.add"`,
		`id="groupAllConditionsButton"`,
		`data-i18n="condition.groupAll"`,
		`id="conditionAllAndButton"`,
		`id="conditionAllOrButton"`,
		`data-i18n="condition.allAnd"`,
		`data-i18n="condition.allOr"`,
		`id="clearConditionGroupsButton"`,
		`data-i18n="condition.clearGroups"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing condition builder i18n html anchor %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"condition.operator.contains"`,
		`"condition.operator.contains": "包含"`,
		`"condition.moveUp": "上移"`,
		`"condition.duplicate": "复制"`,
		`"condition.groupPrevious": "分组"`,
		`"condition.groupAll": "全部分组"`,
		`"condition.allAnd": "全部且"`,
		`"condition.allOr": "全部或"`,
		`"condition.ungroupPrevious": "解组"`,
		`"condition.clearGroups": "清除分组"`,
		`"condition.status.ready": "快速条件已就绪`,
		`"condition.status.duplicated": "已复制条件`,
		`"condition.status.groupedPrevious": "已将条件`,
		`"condition.status.ungroupedPrevious": "已取消条件`,
		`"condition.status.groupedAll": "已将全部`,
		`"condition.status.allAnd": "已将所有条件连接设为且。`,
		`"condition.status.allOr": "已将所有条件连接设为或。`,
		`"condition.status.clearedGroups": "已清除`,
		"applyConditionLocaleStrings",
		"setSelectOptionLabels(conditionOperator",
		"setSelectOptionLabels(conditionTemplateSelect",
		`tf("condition.status.ready"`,
		"moveConditionRow",
		"duplicateConditionRow",
		"groupConditionWithPrevious",
		"groupAllConditionRows",
		"setAllConditionJoins",
		"clearConditionGroups",
		`button.dataset.conditionAction === "move-up"`,
		`button.dataset.conditionAction === "duplicate"`,
		`button.dataset.conditionAction === "group-previous"`,
		`button.dataset.conditionAction === "ungroup-previous"`,
		"document.body.dataset.conditionReorder",
		"document.body.dataset.conditionDuplicate",
		"document.body.dataset.conditionGroupAction",
		"document.body.dataset.conditionGroupAll",
		"document.body.dataset.conditionJoinAction",
		"document.body.dataset.conditionGroupClear",
		`clearConditionGroupsButton.addEventListener("click", clearConditionGroups)`,
		`groupAllConditionsButton.addEventListener("click", groupAllConditionRows)`,
		`conditionAllAndButton.addEventListener("click", () => setAllConditionJoins("AND"))`,
		`conditionAllOrButton.addEventListener("click", () => setAllConditionJoins("OR"))`,
		`t("condition.status.enterValue")`,
		`t("condition.status.enterDate")`,
		`matched: "condition.status.matched"`,
		"statusKeys.matched",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing condition builder i18n wiring %q", snippet)
		}
	}
}

func TestStableSortOptionIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, snippet := range []string{
		`id="columnStableSortToggle"`,
		`Stable sort`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing stable sort html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`let sortStable = true`,
		`const stable = options.stable ?? sortStable`,
		`stable,`,
		`payload.stable !== false`,
		`document.body.dataset.sortStable`,
		`document.body.dataset.sortRemovedKey`,
		`document.body.dataset.sortKeyCount`,
		`removeLastSortKey`,
		`resetViewMappingForReplay`,
		`case "view.removeLastSortKey"`,
		`"view.chip.removeSortKey"`,
		`columnStableSortToggle.addEventListener("change"`,
		`stable: savedSort.stable !== false`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing stable sort frontend wiring %q", snippet)
		}
	}
}

func TestShowAllRowsCommandIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`case "row.showAll":`,
		"async function showAllRows()",
		`fetch("/api/view-clear"`,
		`document.body.dataset.rowShowAll`,
		`addCommand("row.showAll")`,
		`command.id === "row.showAll"`,
		`"command.row.showAll.title": "显示所有行"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing show-all-rows wiring %q", snippet)
		}
	}
	start := strings.Index(source, "async function showAllRows()")
	end := strings.Index(source, "async function moveSelectedRow")
	if start < 0 || end < start {
		t.Fatal("could not locate showAllRows function body")
	}
	if strings.Contains(source[start:end], "hiddenSourceColumns = new Set()") {
		t.Fatal("showAllRows must not restore hidden columns")
	}
}

func TestStartPageNewFileIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`case "file.new":`,
		"async function createNewFile()",
		`fetch("/api/new"`,
		`JSON.stringify({ name, rows: 1, columns: 2 })`,
		`document.body.dataset.newFile`,
		`t("start.newFile")`,
		`t("start.settings")`,
		`activeFileEntry?.detail === "temporary upload"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing start-page new-file wiring %q", snippet)
		}
	}
}

func TestStartPageSupportSurfaceIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, snippet := range []string{
		`id="aboutSupportPanel"`,
		`id="aboutUpdateSummary"`,
		`id="aboutCheckUpdatesButton"`,
		`data-i18n="about.licenseValue"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing about/support html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`case "help.aboutSupport":`,
		"function openAboutSupportPanel()",
		"function syncAboutSupportPanel()",
		"function runAboutSupportUpdateCheck()",
		`t("start.support")`,
		`document.body.dataset.aboutSupportPanel`,
		`document.body.dataset.aboutSupportStatus`,
		`"command.help.aboutSupport.title": "关于与支持"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing about/support frontend wiring %q", snippet)
		}
	}
}

func TestKeyboardShortcutsShellEntryIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, snippet := range []string{
		`id="shortcutsTabButton"`,
		`data-command="settings.keyboardShortcuts"`,
		`aria-label="Keyboard Shortcuts"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing keyboard shortcuts shell html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`const shortcutsTabButton = document.querySelector("#shortcutsTabButton")`,
		`"nav.shortcuts": "Keys"`,
		`[shortcutsTabButton, "nav.shortcuts"]`,
		`shortcutsTabButton.title = t("settings.keyboardShortcuts")`,
		`shortcutsTabButton.addEventListener("click", () => executeCommand("settings.keyboardShortcuts"))`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing keyboard shortcuts shell frontend wiring %q", snippet)
		}
	}
}

func TestNonGridStatusbarVisibilityIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, snippet := range []string{
		`id="statusbar"`,
		`class="statusbar"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing statusbar html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`const statusbar = document.querySelector("#statusbar")`,
		"function syncStatusbarSurface(",
		"statusbar.hidden = false",
		"statusbar.dataset.viewMode = mode",
		"statusbar.dataset.editorSurface = surface",
		"document.body.dataset.statusbarView = mode",
		"document.body.dataset.statusbarSurface = surface",
		"document.body.dataset.statusbarVisible = String(!statusbar.hidden)",
		"syncStatusbarSurface(mode)",
		`syncStatusbarSurface("grid")`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing non-grid statusbar wiring %q", snippet)
		}
	}
}

func TestAppTitleStateIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function syncAppTitle",
		"document.title = title",
		"document.body.dataset.appTitle = title",
		"document.body.dataset.appTitleDirty",
		"document.body.dataset.appTitleSurface",
		"syncAppTitle();",
		"viewMode !== \"grid\" ? viewModeLabel(viewMode) : \"\"",
		"fileLabel ? `${fileLabel}${dirtyLabel}` : \"\"",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing app title wiring %q", snippet)
		}
	}
}

func TestWindowGeometryStateIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`const windowGeometryStorageKey = "opencsv-window-geometry"`,
		"function loadWindowGeometry",
		"function currentWindowGeometry",
		"function saveWindowGeometryNow",
		"function queueWindowGeometrySave",
		"document.body.dataset.windowGeometry",
		"document.body.dataset.windowGeometryRestored",
		"document.body.dataset.windowGeometryMaximized",
		"document.body.dataset.windowGeometrySaved",
		`window.addEventListener("pagehide", saveWindowGeometryNow)`,
		"queueWindowGeometrySave();",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing window geometry wiring %q", snippet)
		}
	}
}

func TestDirtyCloseDiscardPanelIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, snippet := range []string{
		`id="dirtyClosePanel"`,
		`aria-label="Close unsaved tab"`,
		`id="cancelDirtyCloseButton" type="button" data-dialog-role="cancel"`,
		`id="discardDirtyCloseButton" type="button" data-dialog-role="confirm"`,
		`Discard Tab`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing dirty-close html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`const dirtyClosePanel = document.querySelector("#dirtyClosePanel")`,
		"let pendingDirtyCloseTab = null",
		`message === "document has unsaved changes" && !options.force`,
		"openDirtyClosePanel(target, message)",
		"closeOpenFileTab(target.documentId, { force: true })",
		"document.body.dataset.closeTabConfirm",
		"document.body.dataset.closeTabConfirmReason",
		"document.body.dataset.closeTabDiscarded",
		`dirtyClosePanel.addEventListener("mousedown"`,
		`discardDirtyCloseButton.addEventListener("click", confirmDirtyCloseDiscard)`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing dirty-close frontend wiring %q", snippet)
		}
	}
}

func TestPlatformDialogButtonOrderIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, snippet := range []string{
		"function syncPlatformDialogButtonOrder()",
		"function platformDialogButtonOrder()",
		"function dialogButtonOrderValue(role, confirmLeft)",
		`document.querySelectorAll("[data-dialog-actions]")`,
		`actions.dataset.dialogButtonOrder = orderName`,
		`document.body.dataset.dialogButtonPlatform = platform`,
		`document.body.dataset.dialogButtonOrder = orderName`,
		`params.get("dialogPlatform")`,
		`platform.startsWith("win") || platform.includes("windows") ? "windows" : "standard"`,
		"syncPlatformDialogButtonOrder()",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing platform dialog button order wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`data-dialog-actions`,
		`id="confirmSaveAsButton" type="button" data-dialog-role="confirm"`,
		`id="cancelSaveAsButton" type="button" data-dialog-role="cancel"`,
		`id="confirmTransposeButton" type="button" data-dialog-role="confirm"`,
		`id="cancelTransposeConfirmButton" type="button" data-dialog-role="cancel"`,
		`id="confirmOpenLinkButton" type="button" data-i18n="link.prompt.open" data-dialog-role="confirm"`,
		`id="cancelOpenLinkButton" type="button" data-i18n="link.prompt.cancel" data-dialog-role="cancel"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing platform dialog button order html %q", snippet)
		}
	}
}

func TestLaunchURLFlagsAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"applyLaunchParams(params)",
		"function applyLaunchParams(params)",
		`params.get("newWindow")`,
		`params.get("reuseWindow")`,
		`document.body.dataset.launchFlags`,
		`document.body.dataset.launchWait`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing launch URL flag wiring %q", snippet)
		}
	}
}

func TestCLIStdinURLBootIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`params.get("stdin")`,
		`source === "stdin"`,
		`displayName: stdinSource ? (openNames[index] || "stdin")`,
		`remember: !entry.stdinSource`,
		`sourceKind: entry.stdinSource ? "stdin"`,
		`body.sourceKind = options.sourceKind`,
		`body.displayName = options.displayName`,
		`document.body.dataset.stdinUrlBoot = "1"`,
		`document.body.dataset.openSourceKind = options.sourceKind`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing CLI stdin URL boot wiring %q", snippet)
		}
	}
}

func TestMultiOpenURLBootIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`params.getAll("open").filter(Boolean)`,
		`params.getAll("openSource")`,
		`params.getAll("openName")`,
		`params.getAll("openGoto")`,
		`for (const entry of openEntries)`,
		`document.body.dataset.multiOpenCount`,
		`document.body.dataset.multiOpenActive`,
		`document.body.dataset.multiOpenFiles`,
		`Opened ${openEntries.length} files`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing multi-open URL boot wiring %q", snippet)
		}
	}
}

func TestProfileViewI18nIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	for _, snippet := range []string{
		`data-i18n="profile.run"`,
		`data-i18n="profile.rule.columnName"`,
		`data-i18n="profile.rule.default"`,
		`id="profileRuleDefaultInput"`,
		`data-i18n-placeholder="profile.schemaNamePlaceholder"`,
		`<option value="time">Time</option>`,
		`<option value="datetime">Date-Time</option>`,
		`<option value="year">Year</option>`,
		`<option value="yearmonth">Year-Month</option>`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing profile i18n html anchor %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"profile.status.complete"`,
		`"profile.status.complete": "概况`,
		`"profile.quickFix.header.unique": "设为唯一"`,
		`"profile.schema.savedAs": "验证架构已保存为`,
		`"profile.summary.progress": "Scanned"`,
		`"profile.status.completeTimedRows": "Profile:`,
		`"profile.type.time": "Time"`,
		`"profile.type.datetime": "Date-Time"`,
		`"profile.type.yearmonth": "Year-Month"`,
		`"profile.rule.default": "Default"`,
		`"profile.rule.default": "默认值"`,
		`"profile.quickFix.schema.column.add": "Add Column"`,
		`"profile.quickFix.schema.column.add": "添加列"`,
		"applyProfileProgressSummary(result.progress)",
		"readProfileResponse(response)",
		"applyProfileLocaleStrings",
		"renderProfile(profileData)",
		"syncProfileSchemaControls()",
		"constraints.type",
		"isPlainSchemaObject(source.properties)",
		"Object.entries(source.properties)",
		"normalizeProfileSchemaFormat(column.format)",
		"schemaAllowedValues(column, constraints)",
		"schemaAllowedValue(source?.const ?? source?.constant)",
		"schemaNumberFromSources(rangeSources",
		"profile.rule.exclusiveMinSummary",
		`["date-time", "datetime", "timestamp"]`,
		"Array.isArray(value)",
		`type === "null"`,
		"unsignedint",
		`datetime: "profile.type.datetime"`,
		"gyearmonth",
		`summaryNode(t("profile.summary.rows")`,
		`t("profile.empty.noColumns")`,
		"profileSeverityLabel(issue.severity)",
		"profileQuickFixLabel(fix)",
		"uniqueHeaderNameForSourceCol",
		"profileRuleDefaultInput.value",
		"rule.default",
		"applyProfileMissingColumnQuickFix(issue, fix)",
		`fetch("/api/column/insert"`,
		"document.body.dataset.profileQuickFix = `schema.column.add:",
		`tf("profile.quickFix.applied"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing profile i18n wiring %q", snippet)
		}
	}
}

func TestCustomCellLinkPromptIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	stylesCSS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "styles.css"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(indexHTML)
	styles := string(stylesCSS)
	for _, snippet := range []string{
		`confirmedCustomLinkProtocolsStorageKey`,
		`function openCellLinkWithPrompt(link)`,
		`function linkRequiresCustomProtocolPrompt(link)`,
		`function openCustomLinkConfirmDialog(link)`,
		`function confirmCustomLinkOpen()`,
		`function cellLinkAtPoint(x, y)`,
		`document.body.dataset.linkConfirmDialog`,
		`document.body.dataset.linkConfirmRemembered`,
		`window.__opencsvLastOpenedLink`,
		`openCellLinkWithPrompt(link);`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing custom link prompt wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`id="linkConfirmDialog"`,
		`id="linkConfirmRemember"`,
		`id="confirmOpenLinkButton"`,
		`data-i18n="link.prompt.title"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing custom link prompt html %q", snippet)
		}
	}
	for _, snippet := range []string{
		`.link-confirm-dialog`,
		`.link-confirm-body > strong`,
		`.link-confirm-remember`,
	} {
		if !strings.Contains(styles, snippet) {
			t.Fatalf("missing custom link prompt styles %q", snippet)
		}
	}
}

func TestExternalLinkSecurityAuditIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function unsafeLinkPreflightReason(candidate)",
		"function unsafeParsedLinkReason(parsed)",
		"unsafeLinkPreflightReason(candidate)",
		"unsafeParsedLinkReason(parsed)",
		"Contains control characters",
		"Contains encoded control characters",
		"Contains unsafe path separators",
		"Contains embedded credentials",
		`/%(?:0[0-9a-f]|1[0-9a-f]|7f)/i`,
		`parsed.username || parsed.password`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing external link security audit wiring %q", snippet)
		}
	}
}

func TestAutoFillNamedSequencesAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"autoFillNamedSequences",
		"function namedSequenceAutoFill(seed, length)",
		"function namedSequenceMatch(value, sequence)",
		"function applyNamedSequenceCase(value, sample)",
		"const named = namedSequenceAutoFill(seed, values.length)",
		"window.__opencsvAutoFillSeries = autoFillSeries",
		`["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]`,
		`["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing named auto-fill wiring %q", snippet)
		}
	}
}

func sqlResultCommandSuffixForTest(format string) string {
	switch format {
	case "csv":
		return "Csv"
	case "tsv":
		return "Tsv"
	case "markdown":
		return "Markdown"
	case "json":
		return "Json"
	case "html":
		return "Html"
	case "sql":
		return "Sql"
	case "latex":
		return "Latex"
	case "raw":
		return "Raw"
	default:
		return strings.Title(format)
	}
}

func TestHTMLDataCommandsAreRegistered(t *testing.T) {
	indexHTML, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[string]bool)
	var missing []string
	for _, match := range regexp.MustCompile(`data-command="([^"]+)"`).FindAllStringSubmatch(string(indexHTML), -1) {
		id := match[1]
		if seen[id] {
			continue
		}
		seen[id] = true
		if _, ok := ByID(id); !ok {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("html data-command ids missing registry entries: %s", strings.Join(missing, ", "))
	}
}
