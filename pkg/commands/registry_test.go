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
		"macro.startRecording",
		"macro.stopRecording",
		"macro.playLast",
		"macro.clearLast",
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
		"view.editFilterBuilder",
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

func TestCommandMacroRecordingIsWired(t *testing.T) {
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
		`const commandMacroStorageKey = "opencsv-command-macro-last"`,
		`const macroControlCommandIDs = new Set(["macro.startRecording", "macro.stopRecording", "macro.playLast", "macro.clearLast", "command.palette"])`,
		"let macroRecording = false",
		"let macroPlaybackActive = false",
		"let macroDraftCommands = []",
		"let lastCommandMacro = loadLastCommandMacro()",
		"function renderMacroManager()",
		"function startMacroRecording()",
		"function stopMacroRecording()",
		"function clearLastCommandMacro()",
		"async function playLastCommandMacro()",
		"function recordMacroCommand(commandID)",
		"if (options.recordMacro !== false) recordMacroCommand(commandID)",
		`await executeCommand(commandID, { recordMacro: false })`,
		`case "macro.startRecording"`,
		`case "macro.stopRecording"`,
		`case "macro.playLast"`,
		`case "macro.clearLast"`,
		`playMacroButton?.addEventListener("click", () => executeCommand("macro.playLast"))`,
		`clearMacroButton?.addEventListener("click", () => executeCommand("macro.clearLast"))`,
		"document.body.dataset.macroRecording",
		"document.body.dataset.macroPlayback",
		"document.body.dataset.macroSummary",
		`t("macro.status.started")`,
		`t("macro.status.inactive")`,
		`tf("macro.status.saved"`,
		`t("macro.status.savedEmpty")`,
		`t("macro.status.cleared")`,
		`t("macro.status.emptyPlay")`,
		`tf("macro.status.played"`,
		`t("macro.none")`,
		`localizedCommandLabel(commandID, commandID)`,
		`"command.macro.playLast.title": "播放最近宏"`,
		"localStorage.setItem(commandMacroStorageKey",
		"localStorage.removeItem(commandMacroStorageKey)",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing command macro wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`class="shortcut-macro-bar"`,
		`id="macroSummary"`,
		`data-i18n="macro.title"`,
		`id="playMacroButton"`,
		`data-command="macro.playLast"`,
		`data-i18n="macro.play"`,
		`id="clearMacroButton"`,
		`data-command="macro.clearLast"`,
		`data-i18n="macro.clear"`,
		`v=command-i18n-complete-1`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing command macro html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		".shortcut-macro-bar",
		".shortcut-macro-actions",
		"grid-template-rows: auto auto auto auto auto minmax(0, 1fr)",
	} {
		if !strings.Contains(css, snippet) {
			t.Fatalf("missing command macro style wiring %q", snippet)
		}
	}
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing macro script cache-bust version")
	}
}

func TestCommandPaletteShellLocalizationIsWired(t *testing.T) {
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
		`id="commandPalette"`,
		`data-i18n-aria-label="command.paletteAria"`,
		`id="paletteInput"`,
		`data-i18n-aria-label="command.searchAria"`,
		`data-i18n-placeholder="command.searchPlaceholder"`,
		`id="paletteList"`,
		`data-i18n-aria-label="command.listAria"`,
		`v=command-i18n-complete-1`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing command palette localized html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"command.paletteAria": "Command palette"`,
		`"command.paletteAria": "命令面板"`,
		`"command.searchPlaceholder": "搜索命令"`,
		`"command.noMatches": "没有匹配命令"`,
		`"command.status.loadFailed": "命令加载失败：{{message}}"`,
		`"command.status.unavailable": "{{command}}不可用。"`,
		`"command.status.genericUnavailable": "命令不可用。"`,
		`"command.status.historySaveFailed": "命令历史无法保存。"`,
		`"command.status.noHandler": "命令 {{command}} 没有 UI 处理器。"`,
		`function localizedCommandStatus(key, values = {})`,
		`localizedCommandStatus("loadFailed"`,
		`localizedCommandStatus("unavailable"`,
		`localizedCommandStatus("genericUnavailable")`,
		`localizedCommandStatus("historySaveFailed")`,
		`localizedCommandStatus("noHandler"`,
		`emptyPaletteNode(t("command.noMatches"))`,
		`paletteInput.addEventListener("input"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing command palette localized frontend wiring %q", snippet)
		}
	}
}

func TestInspectorShellLocalizationIsWired(t *testing.T) {
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
		`data-i18n-aria-label="inspector.aria"`,
		`data-i18n="inspector.cell"`,
		`data-i18n-aria-label="inspector.cellValueAria"`,
		`data-i18n="inspector.popOut"`,
		`data-i18n="inspector.hiddenColumns"`,
		`data-i18n="inspector.cellQuote"`,
		`data-i18n="inspector.toggleQuote"`,
		`data-i18n-aria-label="inspector.dialogAria"`,
		`data-i18n="inspector.cellValue"`,
		`v=command-i18n-complete-1`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing inspector localized html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"inspector.cell": "单元格"`,
		`"inspector.toggleQuote": "切换"`,
		`"inspector.quote.willQuote": "保存时将加引号"`,
		`tf("inspector.cellStats"`,
		`syncHiddenColumnControls();`,
		`t("inspector.unknown")`,
		`t("inspector.expand")`,
		`t("inspector.collapse")`,
		`t("inspector.status.openCsvFirst")`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing inspector localized frontend wiring %q", snippet)
		}
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

func TestPrintableTypingAndIMEEditingAreWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function shouldStartCellTextEdit(event)",
		`event.isComposing || event.key === "Process" || event.key === "Dead"`,
		`event.metaKey || event.ctrlKey || event.altKey`,
		`openEditor(event.key, { source: "keydown" });`,
		`window.addEventListener("compositionstart"`,
		`openEditor("", { source: "composition", composition: true });`,
		`document.body.dataset.imeEditStart`,
		`if (event.isComposing || editorCompositionActive) return;`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing printable/IME edit wiring %q", snippet)
		}
	}
}

func TestInlineEditorLineBreakShortcutIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function insertEditorLineBreak()",
		"editor.selectionStart",
		"`" + "${editor.value.slice(0, start)}\\n${editor.value.slice(end)}" + "`",
		"document.body.dataset.inlineEditorLineBreak",
		"if (event.shiftKey || event.altKey)",
		"insertEditorLineBreak();",
		"commitEditor();",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing inline editor line-break wiring %q", snippet)
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
		`status.textContent = localizedCopyStatus("copiedSlice", { label, shape, format: copyFormatLabel(format) })`,
		`if (!copied && !includeHeaders && format === "tsv") clearStructureClipboard()`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing row/column copy-as wiring %q", snippet)
		}
	}
}

func TestCopyExportLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing copy/export script cache-bust version")
	}
	for _, snippet := range []string{
		`button.textContent = t("file.export")`,
		`button.textContent = t("edit.copy")`,
		`"command.edit.copySelectedFormat.description": "使用工具栏复制格式复制选中单元格。"`,
		`"copy.status.copiedSelection": "已复制 {{cells}}{{headers}}为 {{format}}"`,
		`"copy.status.sliceFailed": "复制{{label}}失败：{{message}}"`,
		`"copy.slice.rowWithHeaders": "带表头的行"`,
		`"export.status.exported": "已导出 {{rows}}，{{columns}}到 {{filename}}"`,
		`function localizedCopyStatus(key, values = {})`,
		`function localizedCopyCellUnit(count)`,
		`function localizedCopyShape(rows, cols, cellCount)`,
		`function localizedCopySliceLabel(kind, count, includeHeaders = false)`,
		`function localizedExportStatus(key, values = {})`,
		`localizedCopyStatus("noCells")`,
		`localizedCopyStatus("copiedSelection"`,
		`localizedCopyStatus("failed"`,
		`localizedCopySliceLabel("row", rowCount, includeHeaders)`,
		`localizedCopySliceLabel("column", colCount, includeHeaders)`,
		`localizedCopyStatus("sliceFailed"`,
		`localizedCopyStatus("noSliceData"`,
		`localizedCopyShape(rows, cols, cellCount)`,
		`localizedExportStatus("failed"`,
		`localizedExportStatus("exported"`,
	} {
		if !strings.Contains(source+html, snippet) {
			t.Fatalf("missing copy/export localization wiring %q", snippet)
		}
	}
}

func TestTabularCopyWritesHTMLFlavorWhenAvailable(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function tabularClipboardHTML",
		`["csv", "tsv", "raw"].includes(normalizeCopyFormat(format)) ? matrixToHTML(matrix) : ""`,
		"window.__opencsvLastCopiedHTML = html",
		"document.body.dataset.copyHTMLPreview = html.slice(0, 500)",
		`await writeClipboardText(text, { html })`,
		"new ClipboardItem({",
		`"text/plain": new Blob([text], { type: "text/plain" })`,
		`"text/html": new Blob([html], { type: "text/html" })`,
		`document.body.dataset.clipboardFlavors = "text/html,text/plain"`,
		`document.body.dataset.clipboardFlavors = "text/plain"`,
		`document.body.dataset.clipboardFlavors = "text/plain:fallback"`,
		`document.body.dataset.clipboardFlavors = "unavailable"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing tabular copy HTML flavor wiring %q", snippet)
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

func TestSpreadsheetPasteCanExtendGrid(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function sourceColForPasteTarget",
		"const previousRows = sourceRowCount",
		"const previousColumns = columns.length",
		"if (viewMapped && row >= sourceRowCount) break",
		"const sourceRow = sourceRowForDisplayRow(row, \"Paste\")",
		"const sourceCol = sourceColForPasteTarget(col)",
		"ensureSourceRowsForEdits(edits)",
		"ensureSourceColumnsForEdits(edits)",
		"document.body.dataset.pasteExtent",
		"document.body.dataset.pasteExpanded",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing expanding spreadsheet paste wiring %q", snippet)
		}
	}
}

func TestPlainTextCSVPasteIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"function parsePlainTextMatrix",
		"function hasUnquotedDelimiter",
		"function parseDelimitedText",
		"function parseCSV",
		`if (hasUnquotedDelimiter(input, "\t"))`,
		`if (hasUnquotedDelimiter(input, ","))`,
		`char === "," || char === "\t"`,
		`return { matrix: parseTSV(text), source: "tsv" }`,
		`return { matrix: parseCSV(text), source: "csv" }`,
		`const parsed = parsePlainTextMatrix(text)`,
		`pasteParsedMatrix(matrix, parsed.source)`,
		`window.__opencsvParseCSV = parseCSV`,
		`window.__opencsvParsePlainTextMatrix = parsePlainTextMatrix`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing CSV paste wiring %q", snippet)
		}
	}
}

func TestClipboardPasteStatusLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing paste-status script cache-bust version")
	}
	for _, snippet := range []string{
		`"paste.status.empty": "剪贴板为空。"`,
		`"paste.status.failed": "粘贴失败：{{message}}"`,
		`"paste.status.pasted": "已粘贴 {{count}} 个单元格"`,
		`function localizedPasteStatus(key, values = {})`,
		`localizedPasteStatus("failed"`,
		`localizedPasteStatus("empty")`,
		`localizedPasteStatus("pasted"`,
	} {
		if !strings.Contains(source+html, snippet) {
			t.Fatalf("missing clipboard paste status localization wiring %q", snippet)
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

func TestSelectionAggregateLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatalf("missing aggregate localization cache bust")
	}
	if !strings.Contains(html, `data-i18n="aggregate.button.empty"`) {
		t.Fatalf("missing aggregate button localized html wiring")
	}
	for _, snippet := range []string{
		`"aggregate.metric.sum": "求和"`,
		`"aggregate.summary.numeric": " - 已选 {{count}} 个单元格，求和 {{sum}}，平均 {{avg}}，最小 {{min}}，最大 {{max}}"`,
		`"aggregate.status.selectMultiple": "请选择多个单元格进行聚合。"`,
		`"aggregate.status.copied": "已复制 {{label}} {{value}}"`,
		`"aggregate.status.addedRow": "已添加{{label}}聚合行"`,
		`function localizedAggregateStatus(key, values = {})`,
		`function localizedAggregateMetric(metric)`,
		`function localizedAggregateValue(aggregate)`,
		`t("aggregate.button.empty")`,
		`t("aggregate.button.title")`,
		`tf("aggregate.summary.numeric"`,
		`localizedAggregateStatus("selectMultiple")`,
		`localizedAggregateStatus("selectMultipleCopy")`,
		`localizedAggregateStatus("copyFailed"`,
		`localizedAggregateStatus("selectDataRows")`,
		`localizedAggregateStatus("selectCellsCount")`,
		`localizedAggregateStatus("selectNumericCells")`,
		`localizedAggregateStatus("addedRow"`,
		`localizedAggregateStatus("noNumeric")`,
		`localizedAggregateStatus("removedRow")`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing aggregate localization wiring %q", snippet)
		}
	}
}

func TestGridInteractionLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing grid interaction script cache-bust version")
	}
	for _, snippet := range []string{
		`"grid.canvas.resizeColumn": "调整 {{column}} 宽度"`,
		`"grid.canvas.dragMoveColumn": "拖动以移动所选列"`,
		`"grid.status.columnWidth": "{{column}} 宽度 {{width}}px"`,
		`"grid.status.columnAutoFit": "{{column}} 已自动适配为 {{width}}px"`,
		`"grid.status.columnMoveCanceled": "列移动已取消。"`,
		`function localizedGridStatus(key, values = {})`,
		`function localizedGridCanvas(key, values = {})`,
		`function localizedGridColumnLabel(column, index)`,
		`localizedGridCanvas("moveRow")`,
		`localizedGridCanvas("moveColumn")`,
		`localizedGridCanvas("fillSelection")`,
		`localizedGridCanvas("resizeColumn"`,
		`localizedGridCanvas("dragMoveColumn")`,
		`localizedGridCanvas("dragMoveRow")`,
		`localizedGridStatus("movingColumn"`,
		`localizedGridStatus("movingRow"`,
		`localizedGridStatus("columnMoveCanceled")`,
		`localizedGridStatus("rowMoveCanceled")`,
		`localizedGridStatus("columnWidth"`,
		`localizedGridStatus("columnAutoFit"`,
	} {
		if !strings.Contains(source+html, snippet) {
			t.Fatalf("missing grid interaction localization wiring %q", snippet)
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
		"let openTabVisibilitySequence = 0",
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
		"document.body.dataset.openTabVisibilityCheck",
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

func TestSQLConsolePerTabStateIsWired(t *testing.T) {
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
	for _, snippet := range []string{
		`const defaultSQLQuery = "SELECT * FROM current LIMIT 25"`,
		"const defaultSQLPaneSize = 320",
		"sqlState: currentSQLConsoleState()",
		"function currentSQLConsoleState()",
		"paneSize: currentSQLPaneSize()",
		"function applySQLPaneSize",
		"document.body.dataset.sqlPaneSize",
		"function cloneSQLResult(result)",
		"function restoreSQLConsoleState(state)",
		"restoreSQLConsoleState(state.sqlState)",
		"applySQLPaneSize(state?.paneSize",
		"sqlInput.value = defaultSQLQuery",
		`sqlInput.addEventListener("input", syncActiveSQLConsoleTabState)`,
		`sqlPaneSizeInput?.addEventListener("input"`,
		"syncActiveSQLConsoleTabState();",
		"document.body.dataset.sqlTabState",
		"renderSQLResult(lastSQLResult)",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing SQL per-tab state wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`id="sqlPaneSizeInput"`,
		`type="range"`,
		`aria-label="SQL result pane height"`,
		`v=command-i18n-complete-1`,
	} {
		if !strings.Contains(string(indexHTML), snippet) {
			t.Fatalf("missing SQL pane-size html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		"--sql-result-height: 320px",
		".sql-pane-size-input",
		"height: var(--sql-result-height)",
	} {
		if !strings.Contains(string(styles), snippet) {
			t.Fatalf("missing SQL pane-size style wiring %q", snippet)
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
		`id="resetQuotePolicyButton"`,
		`Reset Quotes`,
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
		"function resetQuotePolicySettings",
		"function quotePolicySettingsModified",
		"document.body.dataset.settingsModifiedSections",
		"document.body.dataset.settingsModifiedSectionCount",
		"document.body.dataset.savedViewsReset",
		"document.body.dataset.quotePolicyReset",
		`badge.className = "settings-section-state"`,
		`badge.dataset.state = modified ? "modified" : "default"`,
		`"Saved filters and sorts": savedViewsSettingsModified({ filters, recent, sorts })`,
		`"Editor appearance": editorAppearanceModified()`,
		`"File format presets": presets.length > 0`,
		`resetSavedViewsButton.disabled = !sectionState.state["Saved filters and sorts"]`,
		`resetSavedViewsButton.addEventListener("click", resetSavedViewsSettings)`,
		`function quoteOnlyPolicyModified`,
		`resetQuotePolicyButton.disabled = !documentId || readOnly || !quoteOnlyPolicyModified()`,
		`resetQuotePolicyButton.addEventListener("click", resetQuotePolicySettings)`,
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

func TestSettingsShellLocalizationIsWired(t *testing.T) {
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
		`id="settingsPanel"`,
		`data-i18n-aria-label="settings.title"`,
		`data-i18n-placeholder="settings.searchPlaceholder"`,
		`data-i18n-aria-label="settings.themeModeAria"`,
		`data-i18n="settings.theme.system"`,
		`data-i18n="settings.autoRefresh.clean"`,
		`data-i18n="settings.update.manual"`,
		`data-i18n="settings.checkNow"`,
		`data-i18n="settings.noUpdateCheck"`,
		`data-i18n="settings.savedViews.export"`,
		`data-i18n="settings.savedViews.import"`,
		`data-i18n="settings.savedViews.reset"`,
		`data-i18n-aria-label="settings.savedViews.recentHistoryAria"`,
		`data-i18n="settings.savedViews.clearRecent"`,
		`data-i18n-aria-label="settings.savedViews.importJsonAria"`,
		`data-i18n="settings.savedViews.applyImport"`,
		`data-i18n="settings.linkProtocols.allowed"`,
		`data-i18n-aria-label="settings.linkProtocols.allowedAria"`,
		`data-i18n="settings.linkProtocols.apply"`,
		`data-i18n="settings.linkProtocols.reset"`,
		`data-i18n="settings.quotePolicy"`,
		`data-i18n="settings.quotes.export"`,
		`data-i18n="settings.quotes.import"`,
		`data-i18n="settings.quotes.reset"`,
		`data-i18n="settings.policy.delimiter"`,
		`data-i18n-aria-label="settings.policy.delimiterAria"`,
		`data-i18n="settings.policy.customDelimiter"`,
		`data-i18n="settings.policy.applyDelimiter"`,
		`data-i18n="settings.policy.encoding"`,
		`data-i18n="settings.policy.applyEncoding"`,
		`data-i18n="settings.policy.lineEnding"`,
		`data-i18n="settings.policy.applyLineEnding"`,
		`data-i18n="settings.quotes.mode"`,
		`data-i18n-aria-label="settings.quotes.modeAria"`,
		`data-i18n="settings.quotes.mode.manual"`,
		`data-i18n="settings.quotes.char.single"`,
		`data-i18n="settings.quotes.escape.backslash"`,
		`data-i18n="settings.quotes.apply"`,
		`data-i18n-aria-label="settings.quotes.importJsonAria"`,
		`data-i18n="settings.quotes.applyImport"`,
		`data-i18n="settings.quotes.cancelImport"`,
		`data-i18n="settings.appearance.gridFont"`,
		`data-i18n-aria-label="settings.appearance.gridFontAria"`,
		`data-i18n="settings.appearance.row.default"`,
		`data-i18n="settings.appearance.reset"`,
		`data-i18n="settings.formats.export"`,
		`data-i18n="settings.formats.import"`,
		`data-i18n="settings.formats.reset"`,
		`data-i18n-aria-label="settings.formats.importJsonAria"`,
		`data-i18n="settings.formats.applyImport"`,
		`data-i18n="settings.formats.cancelImport"`,
		`data-i18n-aria-label="settings.formats.presetAria"`,
		`data-i18n="settings.formats.presetEmpty"`,
		`data-i18n="settings.formats.savePreset"`,
		`data-i18n="settings.formats.applyPreset"`,
		`data-i18n-placeholder="settings.formats.rulePatternPlaceholder"`,
		`data-i18n-aria-label="settings.formats.rulePatternAria"`,
		`data-i18n="settings.formats.saveRule"`,
		`data-i18n="settings.noMatches"`,
		`v=command-i18n-complete-1`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing settings localized html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"settings.theme.system": "跟随系统"`,
		`"settings.autoRefresh.notify": "仅通知"`,
		`"settings.update.manual": "仅手动"`,
		`"settings.noUpdateCheck": "尚未运行更新检查。"`,
		`"settings.update.resultLocalPrototype": "本地原型尚未配置远程更新源。"`,
		`"settings.update.statusCompleted": "更新检查已完成。"`,
		`"settings.update.statusFailed": "更新检查失败：{{message}}"`,
		`"settings.update.statusResultSaveFailed": "无法保存更新检查结果。"`,
		`"window.status.geometrySaveFailed": "无法保存窗口布局。"`,
		`"window.status.loadFailed": "窗口加载失败：{{message}}"`,
		`"settings.noMatchesFor": "没有匹配“{{query}}”的设置。"`,
		`"settings.savedViews.summary": "{{filters}} 个筛选，{{recent}} 条最近记录，{{sorts}} 个排序 · 历史 {{history}}"`,
		`"settings.savedViews.statusReset": "保存的筛选、最近筛选和保存的排序已重置。"`,
		`"settings.linkProtocols.allowed": "允许的自定义协议"`,
		`"settings.linkProtocols.statusAllowed": "允许的链接协议：{{protocols}}"`,
		`"settings.linkProtocols.errorUnsafe": "{{protocol}} 不安全，不能允许。"`,
		`"settings.quotePolicy": "引号策略"`,
		`"settings.quotes.mode.manual": "手动"`,
		`"settings.quotes.char.single": "单引号"`,
		`"settings.quotes.escape.backslash": "反斜杠"`,
		`"settings.quotes.statusSet": "引号策略已设为 {{summary}}。"`,
		`"settings.quotes.statusReset": "引号策略已重置。"`,
		`"settings.policy.applyDelimiter": "应用分隔符"`,
		`"settings.policy.statusDelimiterSet": "分隔符已设为 {{delimiter}}。"`,
		`"settings.policy.statusEncodingSet": "编码已设为 {{encoding}}。"`,
		`"settings.policy.statusLineEndingSet": "换行已设为 {{lineEnding}}。"`,
		`"format.status.synced": "已同步"`,
		`"format.status.readOnlyEnabled": "只读模式已启用"`,
		`"format.quote.summary": "引号 {{mode}} {{character}} {{escape}} 转义"`,
		`"format.header.rowPlusColumns": "表头行 + {{columns}}"`,
		`"settings.appearance.gridFont": "表格字体"`,
		`"settings.appearance.row.default": "默认 · 28px"`,
		`"settings.formats.summary": "{{presets}} 个预设，{{rules}} 条规则，{{quoteRules}} 条引号规则"`,
		`"settings.formats.statusImported": "已导入 {{count}} 个格式预设。"`,
		`"settings.formats.errorInvalidJson": "JSON 无效"`,
		`"settings.formats.rowMeta": "{{preview}} · {{rules}} 条规则 · {{quoteRules}} 条引号规则"`,
		`"settings.formats.field.delimiter": "分隔符"`,
		`"settings.formats.quoteMode.manual": "手动"`,
		`"settings.formats.presetAria": "格式预设"`,
		`"settings.formats.presetEmpty": "预设"`,
		`"settings.formats.savePreset": "保存预设"`,
		`"settings.formats.applyPreset": "应用"`,
		`"settings.formats.rulePatternAria": "格式规则模式"`,
		`"settings.formats.rulePatternPlaceholder": "*.csv"`,
		`"settings.formats.saveRule": "保存规则"`,
		`"settings.formats.statusOpenFileBeforeSave": "保存格式预设前请先打开文件。"`,
		`"settings.formats.statusRuleAdded": "已添加格式规则：{{rule}}"`,
		`"settings.formats.statusSavedPreset": "已保存格式预设：{{preset}}"`,
		`tf("settings.theme.summarySystem"`,
		`t("settings.autoRefresh.clean")`,
		`tf("settings.update.summaryLast"`,
		`t("settings.update.resultLocalPrototype")`,
		`t("settings.update.statusCompleted")`,
		`tf("settings.update.statusFailed"`,
		`t("settings.update.statusResultSaveFailed")`,
		`t("window.status.geometrySaveFailed")`,
		`tf("window.status.loadFailed"`,
		`tf("settings.savedViews.summary"`,
		`t("settings.savedViews.noMatches")`,
		`t("settings.savedViews.apply")`,
		`t("settings.savedViews.delete")`,
		`t("settings.linkProtocols.none")`,
		`tf("settings.linkProtocols.statusAllowed"`,
		`tf("settings.linkProtocols.errorUnsafe"`,
		`t("settings.quotes.apply")`,
		`t("settings.policy.applyDelimiter")`,
		`t("settings.policy.applyEncoding")`,
		`t("settings.policy.applyLineEnding")`,
		`tf("settings.policy.statusDelimiterSet"`,
		`tf("settings.policy.statusEncodingSet"`,
		`tf("settings.policy.statusLineEndingSet"`,
		`t("format.status.encoding")`,
		`t("file.readOnly")`,
		`t("format.status.readOnlyEnabled")`,
		`updateFormatStatus();`,
		`quotePolicyDisplayLabel(metadata.dialect)`,
		`headerModeDisplayLabel(metadata.dialect)`,
		`documentFormatDisplaySummary()`,
		`tf("settings.quotes.summary"`,
		`tf("settings.quotes.statusSet"`,
		`t("settings.quotes.statusReset")`,
		`tf("settings.appearance.summary"`,
		`t("settings.appearance.defaultSummary")`,
		`t("settings.appearance.statusUpdated")`,
		`t("settings.appearance.statusReset")`,
		`tf("settings.formats.summary"`,
		`t("settings.formats.empty")`,
		`tf("settings.formats.rowMeta"`,
		`t("settings.formats.presetEmpty")`,
		`t("settings.formats.savePreset")`,
		`t("settings.formats.applyPreset")`,
		`t("settings.formats.saveRule")`,
		`t("settings.formats.addRule")`,
		`t("settings.formats.noAutomaticRules")`,
		`tf("settings.formats.movePresetUpAria"`,
		`formatPresetFieldAria("settings.formats.field.delimiter"`,
		`t("settings.formats.quoteMode.manual")`,
		`t("settings.formats.saveEdit")`,
		`tf("settings.formats.statusExported"`,
		`tf("settings.formats.statusImported"`,
		`localizedFormatPresetImportError(error.message)`,
		`t("settings.formats.statusReset")`,
		`tf("settings.formats.statusUpdatedPreset"`,
		`t("settings.formats.statusCustomDelimiterRequired")`,
		`t("settings.formats.statusOpenFileBeforeSave")`,
		`tf("settings.formats.statusRuleRemoved"`,
		`tf("settings.formats.statusRuleMoved"`,
		`t("settings.formats.statusRuleRequired")`,
		`tf("settings.formats.statusRuleAdded"`,
		`tf("settings.formats.statusSavedPreset"`,
		`t("settings.formats.statusChoosePreset")`,
		`t("settings.checkNow")`,
		`t("settings.savedViews.export")`,
		`t("settings.savedViews.import")`,
		`t("edit.all")`,
		`t("nav.command")`,
		`t("diff.compare")`,
		`t("file.restore")`,
		`t("edit.preview")`,
		`t("view.builder")`,
		`t("inspector.toggleQuote")`,
		`${t("sql.copy")} CSV`,
		`${t("sql.export")} CSV`,
		`"command.file.restoreBackup.description": "从带时间戳的备份中恢复当前文档。"`,
		`"command.format.toggleCellQuote.description": "切换当前单元格的显式引号保留。"`,
		`"command.sql.copyResultCsv.description": "将最近 SQL 结果的表头和行复制为 CSV 文本。"`,
		`"command.sql.exportResultCsv.description": "将最近 SQL 结果的表头和行下载为 CSV 文件。"`,
		`"command.settings.checkForUpdates.description": "在启用更新检查时检查配置的更新来源。"`,
		`"command.workbench.toggleTheme.description": "在浅色、深色、高对比度和跟随系统主题之间切换。"`,
		`"command.workbench.splitEditor.description": "为当前工作区打开第二个编辑器组。"`,
		`"command.workbench.toggleSplitOrientation.description": "在左右并排和上下堆叠布局之间切换编辑器组。"`,
		`"command.workbench.focusNextGroup.description": "在编辑器组之间移动焦点。"`,
		`"command.workbench.moveTabToOtherGroup.description": "将当前文档标签移动到另一个编辑器组。"`,
		`"command.workbench.closeSplitEditor.description": "关闭第二个编辑器组。"`,
		`"command.file.closeTab.description": "关闭当前打开的文件标签页，并保护未保存编辑。"`,
		`"command.format.savePreset.description": "将当前分隔符、表头和编码检测结果保存为可复用的格式预设。"`,
		`"command.format.applyPreset.description": "使用选中的格式预设重新打开当前路径。"`,
		`"command.format.saveRule.description": "把文件名、路径或内容规则关联到当前格式预设，用于自动应用。"`,
		`"command.format.exportPresets.description": "将保存的文件格式预设和自动规则下载为 JSON。"`,
		`"command.format.importPresets.description": "从 JSON 导入文件格式预设和自动规则。"`,
		`"command.format.applyDelimiter.description": "更新当前文档的保存分隔符。"`,
		`"command.format.applyEncoding.description": "更新当前文档的保存编码。"`,
		`"command.format.applyLineEnding.description": "更新当前文档的保存换行符。"`,
		`"command.format.applyQuotePolicy.description": "更新当前文档保存格式的引号字符和引号模式。"`,
		`"command.format.switchWorksheet.description": "聚焦当前工作簿的 XLSX 工作表选择器。"`,
		`"command.format.exportQuoteRules.description": "将当前文档的引号策略和规则映射下载为 JSON。"`,
		`"command.format.importQuoteRules.description": "将引号策略和规则映射 JSON 导入当前文档。"`,
		`"command.format.toggleHeaderQuote.description": "切换表头行的显式引号保留。"`,
		`"command.format.toggleBodyQuote.description": "切换正文行的显式引号保留。"`,
		`"command.format.toggleColumnQuote.description": "切换当前列的显式引号保留。"`,
		`"command.format.quoteSelection.description": "将每个选中单元格标记为保存时加引号。"`,
		`"command.format.noQuoteSelection.description": "为每个选中单元格设置显式不加引号覆盖。"`,
		`"command.format.clearSelectionQuote.description": "移除选中单元格的显式加引号和不加引号覆盖。"`,
		`"command.edit.copyCsv.description": "将选中单元格复制为 CSV 文本。"`,
		`"command.edit.copyMarkdown.description": "将选中单元格复制为 Markdown 表格。"`,
		`"command.edit.copyRaw.description": "不使用 CSV 引号和转义，复制选中单元格的原始值。"`,
		`"command.edit.copyRow.description": "将当前行复制为制表符分隔文本。"`,
		`"command.edit.copyRowWithHeaders.description": "复制当前行并包含表头行。"`,
		`"command.edit.copyRowCsv.description": "将当前行复制为 CSV 文本。"`,
		`"command.edit.copyRowMarkdown.description": "将当前行复制为 Markdown 表格。"`,
		`"command.edit.copyRowRaw.description": "不使用 CSV 引号，复制当前行的原始值。"`,
		`"command.edit.copyColumn.description": "将当前列复制为制表符分隔文本。"`,
		`"command.edit.copyColumnWithHeader.description": "复制当前列并包含它的表头值。"`,
		`"command.edit.copyColumnCsv.description": "将当前列复制为 CSV 文本。"`,
		`"command.edit.copyColumnMarkdown.description": "将当前列复制为 Markdown 表格。"`,
		`"command.edit.copyColumnRaw.description": "不使用 CSV 引号，复制当前列的原始值。"`,
		`t("settings.status.updateOff")`,
		`tf("settings.noMatchesFor"`,
		`t("settings.noMatches")`,
		`syncThemeControls();`,
		`syncAutoRefreshControls();`,
		`syncUpdateCheckControls();`,
		`syncQuotePolicyControls();`,
		`syncEditorAppearanceControls();`,
		`renderFormatPresetSelect();`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing settings localized frontend wiring %q", snippet)
		}
	}
}

func TestFormatPresetsPersistQuoteRules(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		"const quoteSettings = currentQuoteRuleSettings()",
		"...quoteSettings",
		"...formatPresetQuoteSettings(preset)",
		"function formatPresetQuoteSettings(preset)",
		"function formatPresetQuoteRuleCount(preset)",
		"quoteHeader: normalizeQuoteRuleBoolean(raw.quoteHeader)",
		"quoteColumns: normalizeIntegerList(raw.quoteColumns)",
		"quoteCells: normalizeQuoteRuleCells(raw.quoteCells)",
		"noQuoteCells: normalizeQuoteRuleCells(raw.noQuoteCells)",
		"const presetQuoteRulesApplied = preset ? await applyPresetQuoteSettings(preset) : false",
		`fetch("/api/quote-settings"`,
		`tf("settings.quotes.statusAppliedPreset"`,
		"document.body.dataset.formatPresetQuoteRules = quoteRulesSignature(settings)",
		"|qr:${formatPresetQuoteRuleCount(preset)}",
		"quoteRulesSummary(formatPresetQuoteSettings(preset))",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing format preset quote rule wiring %q", snippet)
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
		`t("recovery.section")`,
		`t("recovery.draftName")`,
		`t("recovery.detail")`,
		`tf("recovery.discardAria"`,
		`t("recovery.status.discarded")`,
		`tf("recovery.status.restored"`,
		`"recovery.section": "恢复草稿"`,
		`t("file.openTabs")`,
		`tf("file.closeTabAria"`,
		`tf(entry.rowCount === 1 ? "file.row" : "file.rows"`,
		`t("file.flag.edited")`,
		`t("file.flag.readOnly")`,
		`"file.rows": "{{count}} 行"`,
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
		`data-i18n-aria-label="externalChange.dialogAria"`,
		`data-i18n="externalChange.keepEditing"`,
		`data-i18n="externalChange.compareDisk"`,
		`data-i18n="externalChange.overwriteDisk"`,
		`data-i18n="externalChange.reloadFromDisk"`,
	} {
		if !strings.Contains(markup, snippet) {
			t.Fatalf("missing external overwrite markup %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"externalChange.title.conflict": "未保存编辑冲突"`,
		`"externalChange.status.changedDirty": "磁盘文件已更改，而本地编辑尚未保存。"`,
		`"externalChange.status.refreshedManual": "已刷新外部更改"`,
		`t("externalChange.title.conflict")`,
		`t("externalChange.status.preparingCompare")`,
		`tf("externalChange.status.compareFailed"`,
		`t("externalChange.status.reloadingDiscard")`,
		`t("externalChange.status.saveBlocked")`,
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

func TestSaveRestoreLocalizationIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	htmlBytes, err := os.ReadFile(filepath.Join("..", "..", "web", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	html := string(htmlBytes)
	for _, snippet := range []string{
		`data-i18n-aria-label="saveAs.dialogAria"`,
		`data-i18n="saveAs.title"`,
		`data-i18n="saveAs.path"`,
		`data-i18n-aria-label="saveAs.pathAria"`,
		`data-i18n-aria-label="restore.dialogAria"`,
		`data-i18n="restore.title"`,
		`data-i18n-aria-label="restore.backupFilesAria"`,
		`data-i18n="common.refresh"`,
		`v=command-i18n-complete-1`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing save/restore localized html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"saveAs.title": "另存为"`,
		`"saveAs.status.savedBackup": "已保存 {{path}}；备份 {{backup}}"`,
		`"restore.title": "从备份恢复"`,
		`"restore.status.noBackups": "没有找到此文件的备份。"`,
		`t("saveAs.status.openCsv")`,
		`tf("saveAs.status.savingAs"`,
		`tf("saveAs.status.failed"`,
		`tf("saveAs.status.savedBackup"`,
		`t("restore.status.loading")`,
		`tf("restore.status.failedBackups"`,
		`t("restore.status.noBackups")`,
		`tf("restore.status.restoring"`,
		`tf("restore.status.restored"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing save/restore localized frontend wiring %q", snippet)
		}
	}
}

func TestSplitWorkbenchLocalizationIsWired(t *testing.T) {
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
		`data-command="workbench.splitEditor"`,
		`data-command="workbench.toggleSplitOrientation"`,
		`data-command="workbench.focusNextGroup"`,
		`data-command="workbench.moveTabToOtherGroup"`,
		`data-command="workbench.closeSplitEditor"`,
		`v=command-i18n-complete-1`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing split workbench html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"workbench.status.splitOpened": "已打开拆分编辑器组"`,
		`"workbench.status.movedToGroup": "已将 {{name}} 移到组 {{group}}"`,
		`"workbench.status.gridSummary": "{{rows}}{{aggregate}}，{{columns}}{{path}}{{dirty}}{{mode}}{{external}} - {{selection}}{{stats}}"`,
		`"workbench.moveToGroup": "移到组 {{number}}"`,
		`editorGroupTitle("primary"`,
		`tf("workbench.rowsColumns"`,
		`tf("workbench.status.gridSummary"`,
		`t("workbench.status.splitOpened")`,
		`tf("workbench.status.layout"`,
		`tf("workbench.status.noTabInGroup"`,
		`tf("workbench.status.movedToGroup"`,
		`t("workbench.editorTab")`,
		`tf("workbench.moveToGroup"`,
		`tf("workbench.statusAria"`,
		`renderEditorGroups();`,
		`syncStatusbarSurface();`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing split workbench localized frontend wiring %q", snippet)
		}
	}
}

func TestXLSXWorksheetLocalizationIsWired(t *testing.T) {
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
		`id="workbookSheetSelect"`,
		`data-i18n-aria-label="workbook.worksheetAria"`,
		`data-i18n-title="workbook.worksheetTitle"`,
		`v=command-i18n-complete-1`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing XLSX worksheet localized html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"workbook.worksheetAria": "XLSX 工作表"`,
		`"workbook.status.chooseWorksheet": "请选择 {{count}} 个工作表中的一个。"`,
		`"workbook.status.opened": "已打开工作表 {{sheet}}。"`,
		`const renderSignature = `,
		`workbookSheetSelect.dataset.signature = renderSignature`,
		`t("workbook.sheetFallback")`,
		`tf("workbook.sheetNumber"`,
		`t("workbook.worksheetTitle")`,
		`t("workbook.noAlternateWorksheet")`,
		`t("workbook.status.openWorkbook")`,
		`t("workbook.status.onlyOneWorksheet")`,
		`t("workbook.status.saveOrReload")`,
		`tf("workbook.status.chooseWorksheet"`,
		`t("workbook.status.needsLocalPath")`,
		`tf("workbook.status.opening"`,
		`showOpeningStatus: false`,
		`tf("workbook.status.opened"`,
		`tf("workbook.status.failed"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing XLSX worksheet localized frontend wiring %q", snippet)
		}
	}
}

func TestFileOpenLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing file-open script cache-bust version")
	}
	for _, snippet := range []string{
		`"file.status.opening": "正在打开 {{path}}..."`,
		`"file.status.openFailed": "打开失败：{{message}}"`,
		`"file.status.openedMultiple": "已打开 {{count}} 个文件；当前 {{name}}。"`,
		`"file.status.appliedUrlFormat": "已应用 URL 格式：{{summary}}"`,
		`"file.status.appliedFormatPreset": "已应用格式预设：{{preset}}"`,
		`"file.status.switchingTab": "正在切换到 {{name}}..."`,
		`"file.status.openTabUnavailable": "打开的标签不可用：{{message}}"`,
		`"file.status.switchedTab": "已切换到 {{name}}。"`,
		`"file.status.recentSaveFailed": "无法保存最近文件。"`,
		`"file.status.recentUpdateFailed": "无法更新最近文件。"`,
		`"file.status.removedMissingRecent": "已移除缺失的最近文件：{{message}}"`,
		`"settings.quotes.statusAppliedPreset": "已应用预设引号规则：{{summary}}。"`,
		`tf("file.status.opening"`,
		`tf("file.status.openFailed"`,
		`tf("file.status.openedMultiple"`,
		`tf("file.status.appliedUrlFormat"`,
		`tf("file.status.appliedPresetQuoteRules"`,
		`tf("settings.quotes.statusAppliedPreset"`,
		`tf("file.status.appliedUrlLocation"`,
		`tf("file.status.appliedFormatRule"`,
		`tf("file.status.appliedFormatPreset"`,
		`tf("file.status.switchingTab"`,
		`tf("file.status.openTabUnavailable"`,
		`tf("file.status.switchedTab"`,
		`t("file.status.recentSaveFailed")`,
		`t("file.status.recentUpdateFailed")`,
		`tf("file.status.removedMissingRecent"`,
		`t("file.status.stdin")`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing file-open localized frontend wiring %q", snippet)
		}
	}
}

func TestStructuralEditLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing structural-edit script cache-bust version")
	}
	for _, snippet := range []string{
		`"edit.status.aggregateRowsViewOnly": "聚合行仅可查看。"`,
		`"edit.status.failed": "{{action}}失败：{{message}}"`,
		`"edit.status.cutForPaste": "已剪切 {{unit}}，可用于粘贴。"`,
		`"edit.status.restoredHiddenRows": "已恢复隐藏的 {{unit}}"`,
		`"edit.status.columnHidden": "已隐藏列：{{column}}"`,
		`"edit.status.movedRowTo": "已将第 {{from}} 行移到第 {{to}} 行"`,
		`"edit.status.movedColumnTo": "已将第 {{from}} 列移到第 {{to}} 列"`,
		`"edit.status.noDuplicateRows": "未找到重复行"`,
		`"command.edit.undo.title": "撤销"`,
		`"command.row.hide.title": "隐藏行"`,
		`"command.column.pasteRight.title": "向右粘贴列"`,
		`function localizedEditStatus(key, values = {})`,
		`function localizedEditAction(commandID, fallback)`,
		`function localizedEditUnit(kind, count)`,
		`function localizedEditKind(kind)`,
		`localizedEditFailed(action, payload.error ?? "unknown error")`,
		`localizedEditStatus("cutForPaste"`,
		`t("edit.status.copyOrCutRowFirst")`,
		`t("edit.status.copyOrCutColumnFirst")`,
		`t("edit.status.rowHidden")`,
		`t("edit.status.allRowsAlreadyShown")`,
		`localizedEditStatus("movedRowTo"`,
		`localizedEditStatus("movedColumnTo"`,
		`localizedEditStatus("columnHidden"`,
		`localizedEditStatus("restoredHiddenColumns"`,
		`localizedEditStatus("removedDuplicateRows"`,
		`localizedEditStatus("removedDuplicateColumns"`,
		`localizedEditStatus("appliedRows"`,
		`t("edit.status.aggregateRowsViewOnly")`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing structural edit localized frontend wiring %q", snippet)
		}
	}
}

func TestCellToolLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing cell-tool script cache-bust version")
	}
	for _, snippet := range []string{
		`data-i18n-aria-label="cell.split.dialogAria"`,
		`data-i18n="cell.split.title"`,
		`data-i18n="cell.split.confirm"`,
		`"cell.status.changed": "{{action}}：已更改 {{count}} 个单元格"`,
		`"cell.status.edited": "单元格已编辑"`,
		`"cell.status.editedCells": "已编辑 {{count}} 个单元格"`,
		`"cell.status.editReverted": "单元格编辑已还原"`,
		`"cell.status.editFailed": "编辑失败：{{message}}"`,
		`"cell.status.selectTwoToJoin": "请选择至少两个单元格进行合并。"`,
		`"cell.split.summaryHorizontal": "横向拆分覆盖"`,
		`"cell.split.statusBlockedCount": "拆分被 {{count}} 个已有单元格阻止。"`,
		`"command.cell.fillDown.title": "向下填充"`,
		`"command.cell.swapRight.title": "向右交换"`,
		`function localizedCellStatus(key, values = {})`,
		`function localizedCellAction(commandID, fallback)`,
		`function localizedCellChanged(action, count)`,
		`t("cell.status.selectTwoToSwap")`,
		`localizedCellStatus("noAdjacent"`,
		`t(vertical ? "cell.status.selectTwoToJoinVertical" : "cell.status.selectTwoToJoin")`,
		`t(request.direction === "vertical" ? "cell.split.summaryVertical" : "cell.split.summaryHorizontal")`,
		`tf("cell.split.detail"`,
		`tf("cell.split.statusBlockedCount"`,
		`localizedCellStatus("autoFillSelectTwo"`,
		`localizedCellStatus("autoFillLimit"`,
		`localizedCellStatus("loadFilteredRows"`,
		`t("cell.status.editReverted")`,
		`localizedCellStatus("editedCells"`,
		`t("cell.status.edited")`,
		`localizedCellStatus("editFailed"`,
		`status.textContent = localizedCellRunning(options.label)`,
		`status.textContent = localizedCellFailed(options.label, error)`,
		`const message = changed > 0 ? localizedCellChanged(options.label, changed) : localizedCellNoChange(options.label)`,
	} {
		if !strings.Contains(source+html, snippet) {
			t.Fatalf("missing cell-tool localized wiring %q", snippet)
		}
	}
}

func TestSmartSelectionStatusLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing selection-status script cache-bust version")
	}
	for _, snippet := range []string{
		`"selection.status.selectAtLeastOne": "请选择至少一个单元格。"`,
		`"selection.status.transformLimit": "选区有 {{count}} 个单元格；转换上限为 {{limit}}。"`,
		`"selection.status.addedRange": "已添加范围到选区：已选 {{count}} 个单元格"`,
		`"selection.status.subtractedRange": "已减去范围：已选 {{count}} 个单元格"`,
		`"selection.status.inverted": "已反转选区：已选 {{count}} 个单元格"`,
		`"selection.status.extraCleared": "已清除额外选区范围"`,
		`function localizedSelectionStatus(key, values = {})`,
		`localizedSelectionStatus("loadFilteredRowsStructure")`,
		`localizedSelectionStatus("selectAtLeastOne")`,
		`localizedSelectionStatus("transformLimit"`,
		`localizedSelectionStatus("headerIncluded")`,
		`localizedSelectionStatus("addedGesture"`,
		`localizedSelectionStatus("addedRange"`,
		`localizedSelectionStatus("subtractedRange"`,
		`localizedSelectionStatus("inverted"`,
		`localizedSelectionStatus("extraCleared")`,
	} {
		if !strings.Contains(source+html, snippet) {
			t.Fatalf("missing smart selection status localization wiring %q", snippet)
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

func TestImplementedCommandsHaveSimplifiedChineseLocalization(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	var missing []string
	for _, command := range Registry() {
		if !command.Implemented {
			continue
		}
		for _, suffix := range []string{"title", "description"} {
			key := `"` + "command." + command.ID + "." + suffix + `"`
			if !strings.Contains(source, key) {
				missing = append(missing, key)
			}
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("implemented commands missing zh-CN localization: %s", strings.Join(missing, ", "))
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
	html := string(indexHTML)
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing SQL result command localization script cache-bust version")
	}
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
	for _, id := range []string{
		"sql.openResult",
		"sql.copyCell",
		"sql.copyRow",
		"sql.copyRowWithHeaders",
		"sql.copyColumn",
		"sql.copyColumnWithHeader",
		"sql.clearSelection",
	} {
		for _, suffix := range []string{"title", "description"} {
			key := `"` + "command." + id + "." + suffix + `"`
			if !strings.Contains(source, key) {
				missing = append(missing, key+" zh-CN")
			}
		}
	}
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
			for _, keySuffix := range []string{"title", "description"} {
				key := `"` + "command." + id + "." + keySuffix + `"`
				if !strings.Contains(source, key) {
					missing = append(missing, key+" zh-CN")
				}
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
		`tf("diff.status.goto"`,
		"localizedDiffSide(targetSide)",
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
		`t("diff.status.applying")`,
		`tf("diff.status.mergeFailed"`,
		`tf("diff.status.merged"`,
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
		`"diff.status.chooseInputs": "请选择左侧 CSV 和右侧 CSV 路径或文本进行比较。"`,
		`"diff.status.comparing": "正在比较文件..."`,
		`"diff.status.failed": "差异比较失败：{{message}}"`,
		`"diff.status.clipboardUnavailable": "剪贴板读取不可用；请将 CSV 文本粘贴到右侧。"`,
		`"diff.status.clipboardLoaded": "已加载 {{count}} 个剪贴板字符用于比较。"`,
		`"diff.status.goto": "差异跳转到{{side}}第 {{row}} 行。"`,
		`"diff.status.merged": "已合并 {{applied}}{{backup}}"`,
		"applyDiffLocaleStrings",
		"setSelectOptionLabels(diffStatusFilter",
		"localizedDiffStatus(row.status)",
		"appendDiffOpenSourceActions",
		"openDiffSourceRow(\"left\"",
		"document.body.dataset.diffOpenSource",
		`t("diff.status.chooseInputs")`,
		`t("diff.status.comparing")`,
		`tf("diff.status.failed"`,
		`t("diff.status.clipboardUnavailable")`,
		`tf("diff.status.clipboardLoaded"`,
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
		"document.body.dataset.columnValuesOrder",
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
		`id="conditionPreviewText"`,
		`data-i18n="condition.previewTitle"`,
		`data-i18n="condition.previewEmpty"`,
		`id="previewConditionButton"`,
		`data-i18n="condition.previewButton"`,
		`id="conditionTreeActionSelect"`,
		`data-i18n="condition.treeAction"`,
		`data-i18n="condition.treeAction.choose"`,
		`id="runConditionTreeActionButton"`,
		`data-i18n="condition.runTreeAction"`,
		`id="groupAllConditionsButton"`,
		`data-i18n="condition.groupAll"`,
		`id="conditionAllAndButton"`,
		`id="conditionAllOrButton"`,
		`data-i18n="condition.allAnd"`,
		`data-i18n="condition.allOr"`,
		`id="balanceConditionGroupsButton"`,
		`data-i18n="condition.balanceGroups"`,
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
		`"condition.previewTitle": "SQL 预览"`,
		`"condition.previewButton": "预览 SQL"`,
		`"condition.treeAction": "树操作"`,
		`"condition.runTreeAction": "运行"`,
		`"condition.rowAction": "行操作"`,
		`"condition.runRowAction": "运行"`,
		`"condition.groupAll": "全部分组"`,
		`"condition.allAnd": "全部且"`,
		`"condition.allOr": "全部或"`,
		`"condition.ungroupPrevious": "解组"`,
		`"condition.balanceGroups": "平衡分组"`,
		`"condition.clearGroups": "清除分组"`,
		`"condition.status.ready": "快速条件已就绪`,
		`"condition.status.duplicated": "已复制条件`,
		`"condition.status.groupedPrevious": "已将条件`,
		`"condition.status.ungroupedPrevious": "已取消条件`,
		`"condition.status.groupedAll": "已将全部`,
		`"condition.status.allAnd": "已将所有条件连接设为且。`,
		`"condition.status.allOr": "已将所有条件连接设为或。`,
		`"condition.status.chooseTreeAction": "请选择条件树操作。`,
		`"condition.status.chooseRowAction": "请选择条件`,
		`"condition.status.previewReady": "SQL 预览已就绪。`,
		`"view.chip.edit": "Edit"`,
		`"view.chip.edit": "编辑"`,
		`"view.chip.editBuilder": "Edit this filter in Builder"`,
		`"view.status.filterToBuilder": "Loaded current filter in Filter Builder"`,
		`"command.view.editFilterBuilder.title": "在筛选构建器中编辑"`,
		`"condition.status.balancedGroups": "已平衡条件分组`,
		`"condition.status.clearedGroups": "已清除`,
		"applyConditionLocaleStrings",
		"setSelectOptionLabels(conditionOperator",
		"setSelectOptionLabels(conditionTemplateSelect",
		"setSelectOptionLabels(conditionTreeActionSelect",
		"conditionTreeActionLabelKeys",
		"renderConditionPreview",
		"conditionBuilderExpressionFromDOM",
		"previewConditionBuilder",
		"openActiveFilterInBuilder",
		"openSavedFilterInBuilder",
		"replayFilterRowsInBuilder",
		`case "view.editFilterBuilder":`,
		`const needsBuilderFilter = command.id === "view.editFilterBuilder"`,
		"setConditionBuilderRows(builderRows)",
		"handleConditionBuilderEnter",
		"handleConditionBuilderRowShortcut",
		`event.key.toLowerCase() === "g"`,
		`event.key.toLowerCase() === "u"`,
		`event.key === "Backspace" || event.key === "Delete"`,
		`tf("condition.status.ready"`,
		"moveConditionRow",
		"duplicateConditionRow",
		"groupConditionWithPrevious",
		"groupAllConditionRows",
		"setAllConditionJoins",
		"balanceConditionGroups",
		"runConditionTreeAction",
		"conditionRowActionSelect",
		"runConditionRowMenuAction",
		"clearConditionGroups",
		`button.dataset.conditionAction === "row-menu"`,
		`button.dataset.conditionAction === "move-up"`,
		`button.dataset.conditionAction === "duplicate"`,
		`button.dataset.conditionAction === "group-previous"`,
		`button.dataset.conditionAction === "ungroup-previous"`,
		"document.body.dataset.conditionReorder",
		"document.body.dataset.conditionDuplicate",
		"document.body.dataset.conditionGroupAction",
		"document.body.dataset.conditionGroupAll",
		"document.body.dataset.conditionJoinAction",
		"document.body.dataset.conditionTreeMenuAction",
		"document.body.dataset.conditionRowMenuAction",
		"document.body.dataset.conditionBuilderPreview",
		"document.body.dataset.conditionBuilderPreviewState",
		"document.body.dataset.conditionBuilderReopened",
		"document.body.dataset.savedFilterBuilderReopened",
		"document.body.dataset.conditionBuilderPreviewShortcut",
		"document.body.dataset.conditionRowKeyboardAction",
		"document.body.dataset.conditionGroupBalance",
		"document.body.dataset.conditionGroupClear",
		`runConditionTreeActionButton.addEventListener("click", runConditionTreeAction)`,
		`previewConditionButton.addEventListener("click", previewConditionBuilder)`,
		`editButton.addEventListener("click", openActiveFilterInBuilder)`,
		`editButton.addEventListener("click", () => openSavedFilterInBuilder(filter.id))`,
		`normalizeConditionBuilderRows(filter.builderRows).length`,
		`balanceConditionGroupsButton.addEventListener("click", balanceConditionGroups)`,
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
		`document.body.dataset.sortViewUndo`,
		`document.body.dataset.sortViewRedo`,
		`sortViewUndoStack`,
		`pushSortViewUndoSnapshot`,
		`undoSortView`,
		`redoSortView`,
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
		`data-i18n-aria-label="about.dialogAria"`,
		`id="aboutUpdateSummary"`,
		`id="aboutCheckUpdatesButton"`,
		`data-i18n="about.licenseValue"`,
		`data-i18n="settings.noUpdateCheck"`,
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
		`"about.dialogAria": "关于与支持"`,
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
		`data-i18n-aria-label="settings.keyboardShortcuts"`,
		`data-i18n="settings.keyboardShortcuts"`,
		`data-i18n="shortcut.export"`,
		`data-i18n="shortcut.import"`,
		`data-i18n="shortcut.preset.analysis"`,
		`data-i18n-placeholder="shortcut.searchPlaceholder"`,
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
		`root.querySelectorAll("[data-i18n-aria-label]")`,
		`"shortcut.preset.analysis": "分析"`,
		`"command.settings.exportShortcuts.title": "导出键盘快捷键"`,
		`t("shortcut.pressShortcut")`,
		`t("shortcut.record")`,
		`t("shortcut.reset")`,
		`t("shortcut.replace")`,
		`tf("shortcut.status.exported"`,
		`tf("shortcut.status.imported"`,
		`t("shortcut.export")`,
		`t("shortcut.import")`,
		`t("shortcut.error.invalidJson")`,
		`tf("shortcut.status.assigned"`,
		`shortcutsTabButton.addEventListener("click", () => executeCommand("settings.keyboardShortcuts"))`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing keyboard shortcuts shell frontend wiring %q", snippet)
		}
	}
}

func TestKeyboardShortcutResetAllIsWired(t *testing.T) {
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
		`id="resetAllShortcutsButton"`,
		`aria-label="Reset All Keyboard Shortcuts"`,
		`data-i18n="shortcut.resetAll"`,
		`data-i18n-aria-label="shortcut.resetAllAria"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing reset-all shortcuts html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`const resetAllShortcutsButton = document.querySelector("#resetAllShortcutsButton")`,
		"function syncResetAllShortcutsButton()",
		"resetAllShortcutsButton.disabled = customShortcuts.size === 0",
		"function resetAllShortcuts()",
		"customShortcuts = new Map()",
		"document.body.dataset.shortcutResetAll = String(resetCount)",
		`t("shortcut.status.defaults")`,
		`tf("shortcut.status.resetTitle"`,
		`tf("shortcut.status.resetCount"`,
		`t("shortcut.status.alreadyDefaults")`,
		`resetAllShortcutsButton.addEventListener("click", resetAllShortcuts)`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing reset-all shortcuts frontend wiring %q", snippet)
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
		`data-i18n-aria-label="dirtyClose.dialogAria"`,
		`data-i18n="dirtyClose.title"`,
		`data-i18n="dirtyClose.message.generic"`,
		`data-i18n="dirtyClose.warning"`,
		`id="cancelDirtyCloseButton" type="button" data-dialog-role="cancel" data-i18n="dirtyClose.keepEditing"`,
		`id="discardDirtyCloseButton" type="button" data-dialog-role="confirm" data-i18n="dirtyClose.discardTab"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing dirty-close html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"dirtyClose.title": "关闭未保存标签"`,
		`"dirtyClose.keepEditing": "继续编辑"`,
		`"dirtyClose.status.closed": "已关闭 {{name}}。"`,
		`const dirtyClosePanel = document.querySelector("#dirtyClosePanel")`,
		"let pendingDirtyCloseTab = null",
		`message === "document has unsaved changes" && !options.force`,
		"openDirtyClosePanel(target, message)",
		"function updateDirtyClosePanel()",
		`tf("dirtyClose.message.detail"`,
		`tf("dirtyClose.status.saveBeforeClosing"`,
		`tf("dirtyClose.status.closed"`,
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
		`id="confirmTransposeButton" type="button" data-i18n="transpose.confirm" data-dialog-role="confirm"`,
		`id="cancelTransposeConfirmButton" type="button" data-i18n="transpose.cancel" data-dialog-role="cancel"`,
		`id="confirmOpenLinkButton" type="button" data-i18n="link.prompt.open" data-dialog-role="confirm"`,
		`id="cancelOpenLinkButton" type="button" data-i18n="link.prompt.cancel" data-dialog-role="cancel"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing platform dialog button order html %q", snippet)
		}
	}
}

func TestTransposeLocalizationIsWired(t *testing.T) {
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
	if !strings.Contains(html, `v=command-i18n-complete-1`) {
		t.Fatal("missing transpose localization script cache-bust version")
	}
	for _, snippet := range []string{
		`data-i18n-aria-label="transpose.dialogAria"`,
		`data-i18n="transpose.title"`,
		`data-i18n="transpose.close"`,
		`data-i18n="transpose.summaryDefault"`,
		`data-i18n="transpose.detailDefault"`,
		`data-i18n="transpose.warningDefault"`,
		`data-i18n="transpose.cancel"`,
		`data-i18n="transpose.confirm"`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing transpose localized html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"transpose.title": "转置表格"`,
		`"transpose.summary": "转置 {{rows}} 行 x {{columns}} 列"`,
		`"transpose.status.confirm": "请确认表格转置。"`,
		`"transpose.status.selected": "已转置 {{rows}}x{{columns}} 选区"`,
		`"transpose.status.multiSelected": "已转置 {{ranges}} 个范围（{{cells}} 个单元格更改）"`,
		`"edit.status.loadFilteredRowsAction": "使用{{action}}前请先加载选中的筛选行。"`,
		`tf("transpose.summary"`,
		`t("transpose.detail")`,
		`tf("transpose.warningBlocked"`,
		`t(plan.blocked ? "transpose.status.tooLarge" : "transpose.status.confirm")`,
		`t("transpose.status.noSelectionChange")`,
		`tf("transpose.status.selectionLimit"`,
		`t("transpose.status.rangeTooSmall")`,
		`t("transpose.status.overlap")`,
		`tf("transpose.status.selected"`,
		`tf("transpose.status.multiSelected"`,
		`t("transpose.status.filteredRows")`,
		`tf("edit.status.loadFilteredRowsAction"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing transpose localized frontend wiring %q", snippet)
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
		`tf("file.status.openedMultiple"`,
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

func TestQuoteRuleListEditingIsWired(t *testing.T) {
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
		`id="quoteRuleList"`,
		`class="quote-rule-list"`,
		`aria-label="Active quote rules"`,
		`id="addQuoteCellRuleButton"`,
		`id="addNoQuoteCellRuleButton"`,
		`id="clearCellQuoteRuleButton"`,
		`id="quoteSelectionRuleButton"`,
		`id="noQuoteSelectionRuleButton"`,
		`id="clearSelectionQuoteRuleButton"`,
		`id="quoteRuleColumnInput"`,
		`aria-label="Quote column rule source"`,
		`id="quoteRuleColumnActionSelect"`,
		`id="applyQuoteColumnRuleButton"`,
		`id="quoteRuleColumnPatternInput"`,
		`id="quoteRuleColumnPatternModeSelect"`,
		`id="quoteRuleColumnPatternActionSelect"`,
		`id="applyQuoteColumnPatternRuleButton"`,
		`id="quoteRuleValuePatternInput"`,
		`id="quoteRuleValuePatternModeSelect"`,
		`id="quoteRuleValuePatternActionSelect"`,
		`id="previewQuoteValuePatternRuleButton"`,
		`id="applyQuoteValuePatternRuleButton"`,
		`id="quoteRuleRowStartInput"`,
		`id="quoteRuleRowEndInput"`,
		`id="quoteRuleRowActionSelect"`,
		`id="applyQuoteRowRangeRuleButton"`,
		`id="quoteRuleRowInput"`,
		`id="quoteRuleColInput"`,
		`id="quoteRuleActionSelect"`,
		`id="applyQuoteCoordinateRuleButton"`,
		`Quote Cell`,
		`No Quote Cell`,
		`Clear Cell`,
		`Quote Selection`,
		`No Quote Selection`,
		`Clear Selection`,
		`data-i18n="settings.quotes.quoteCell"`,
		`data-i18n="settings.quotes.noQuoteSelection"`,
		`data-i18n="settings.quotes.sourceColumn"`,
		`data-i18n-aria-label="settings.quotes.sourceColumnAria"`,
		`data-i18n="settings.quotes.columnAction"`,
		`data-i18n="settings.quotes.action.quoteColumn"`,
		`data-i18n="settings.quotes.columnPattern"`,
		`data-i18n="settings.quotes.match.regex"`,
		`data-i18n="settings.quotes.valuePattern"`,
		`data-i18n="settings.quotes.action.noQuoteValues"`,
		`data-i18n="settings.quotes.rowStart"`,
		`data-i18n="settings.quotes.action.noQuoteRows"`,
		`data-i18n="settings.quotes.ruleAction"`,
		`data-i18n="settings.quotes.applyRule"`,
		`data-i18n-aria-label="settings.quotes.activeRulesAria"`,
		`Apply Column`,
		`Apply Pattern`,
		`Apply Values`,
		`Preview Values`,
		`Apply Rows`,
		`Apply Rule`,
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing quote rule list html wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`const quoteRuleList = document.querySelector("#quoteRuleList")`,
		"const quoteRuleListRenderLimit = 500",
		`const addQuoteCellRuleButton = document.querySelector("#addQuoteCellRuleButton")`,
		`const quoteSelectionRuleButton = document.querySelector("#quoteSelectionRuleButton")`,
		`const quoteRuleColumnInput = document.querySelector("#quoteRuleColumnInput")`,
		`const applyQuoteColumnRuleButton = document.querySelector("#applyQuoteColumnRuleButton")`,
		`const quoteRuleColumnPatternInput = document.querySelector("#quoteRuleColumnPatternInput")`,
		`const applyQuoteColumnPatternRuleButton = document.querySelector("#applyQuoteColumnPatternRuleButton")`,
		`const quoteRuleValuePatternInput = document.querySelector("#quoteRuleValuePatternInput")`,
		`const previewQuoteValuePatternRuleButton = document.querySelector("#previewQuoteValuePatternRuleButton")`,
		`const applyQuoteValuePatternRuleButton = document.querySelector("#applyQuoteValuePatternRuleButton")`,
		`const quoteRuleRowStartInput = document.querySelector("#quoteRuleRowStartInput")`,
		`const applyQuoteRowRangeRuleButton = document.querySelector("#applyQuoteRowRangeRuleButton")`,
		`const quoteRuleRowInput = document.querySelector("#quoteRuleRowInput")`,
		`const applyQuoteCoordinateRuleButton = document.querySelector("#applyQuoteCoordinateRuleButton")`,
		"function renderQuoteRuleList()",
		"function quoteMapLimitLabel(dialect = metadata?.dialect)",
		"document.body.dataset.quoteMapTruncated",
		"document.body.dataset.quoteMapDetected",
		"document.body.dataset.quoteMapLimit",
		"function quoteRuleItemCount(settings)",
		"function activeQuoteRuleItems(settings, limit = Number.POSITIVE_INFINITY)",
		"activeQuoteRuleItems(settings, quoteRuleListRenderLimit)",
		"document.body.dataset.quoteRuleListRendered = String(items.length)",
		`"settings.quotes.quoteSelection": "选区加引号"`,
		`"settings.quotes.action.noQuoteValues": "值不加引号"`,
		`"settings.quotes.item.headerTitle": "表头行加引号"`,
		`"settings.quotes.item.noCellTitle": "单元格 {{cell}} 不加引号"`,
		`"settings.quotes.statusCellApplied": "{{action}}：{{cell}}。"`,
		`"settings.quotes.statusRemoved": "已移除引号规则：{{rule}}。"`,
		`"settings.quotes.statusRulesOpenFile": "更改引号规则前请先打开文件。"`,
		`"settings.quotes.statusCellStateEnabled": "单元格引号已启用。"`,
		`"settings.quotes.statusRuleToggled": "{{target}}引号已{{state}}。"`,
		`"settings.quotes.statusSelectionApplied": "{{action}}：已更新 {{count}} 个单元格。"`,
		`"settings.quotes.statusActionReadOnly": "{{action}}被阻止：文档为只读"`,
		`"settings.quotes.statusValuePatternMatched": "{{action}}：全文件 {{scanned}} 个单元格中匹配 {{matched}} 个{{limit}}。"`,
		`"settings.quotes.statusExported": "已导出引号规则：{{summary}}。"`,
		`"settings.quotes.statusImported": "已导入引号规则：{{summary}}。"`,
		`t("settings.quotes.empty.noRules")`,
		`t("settings.quotes.removeRule")`,
		`t("settings.quotes.statusRulesOpenFile")`,
		`settings.quotes.statusCellStateEnabled`,
		`t(quoted ? "settings.quotes.statusCellStateEnabled"`,
		`tf("settings.quotes.statusRuleToggled"`,
		`tf("settings.quotes.statusSelectionApplied"`,
		`tf("settings.quotes.statusCellApplied"`,
		`tf("settings.quotes.statusRemoved"`,
		`tf("settings.quotes.statusExported"`,
		`tf("settings.quotes.statusImported"`,
		`tf("settings.quotes.statusActionReadOnly"`,
		`tf("settings.quotes.statusValuePatternMatched"`,
		`tf("settings.quotes.empty.showingFirst"`,
		`tf("settings.quotes.item.columnTitle"`,
		`tf("settings.quotes.item.noCellTitle"`,
		"function applyActiveCellQuoteRule(mode)",
		"function resolveQuoteRuleColumn(value)",
		"function applyColumnQuoteRule()",
		"function matchingQuoteRuleColumns(pattern, matchMode)",
		"function applyColumnPatternQuoteRule()",
		"function validateQuoteValuePattern(pattern, matchMode)",
		"function applyValuePatternQuoteRule(preview = false)",
		"function quoteCellsForSourceRowRange(startValue, endValue, maxCells = 50_000)",
		"function applyRowRangeQuoteRule()",
		"function applyCoordinateQuoteRule()",
		"function removeQuoteRule(item)",
		`document.body.dataset.quoteRuleList = String(totalItems)`,
		"document.body.dataset.quoteCellRule = `${mode}:${row}:${col}:${signature}`",
		"document.body.dataset.quoteColumnRule = `${mode}:${col}:${signature}`",
		"document.body.dataset.quoteColumnPatternRule = `${matchMode}:${action}:${matched.columns.join(\",\")}:${signature}`",
		`fetch("/api/quote-value-rule"`,
		"document.body.dataset.quoteValuePatternRule = response.ok",
		`${preview ? "preview" : matchMode}:${action}:${payload.matched}:${payload.scannedCells}:${payload.truncated ? "truncated" : "full"}`,
		"document.body.dataset.quoteValuePatternPreview = document.body.dataset.quoteValuePatternRule",
		`tf("settings.quotes.statusValuePatternMatched"`,
		"document.body.dataset.quoteRowRangeRule = `${action}:${range.startRow}-${range.endRow}:${range.cells.length}:${signature}`",
		"document.body.dataset.quoteCoordinateRule = `${mode}:${row}:${col}:${signature}`",
		"quoteRuleColumnInput.dataset.syncedValue = defaultColumn",
		"quoteRuleRowInput.dataset.syncedValue = defaultRow",
		"document.body.dataset.quoteRuleRemove = `${item.kind}:${item.key}`",
		`fetch("/api/quote-settings"`,
		`addQuoteCellRuleButton.addEventListener("click", () => applyActiveCellQuoteRule("quote"))`,
		`addNoQuoteCellRuleButton.addEventListener("click", () => applyActiveCellQuoteRule("no-quote"))`,
		`clearCellQuoteRuleButton.addEventListener("click", () => applyActiveCellQuoteRule("clear"))`,
		`quoteSelectionRuleButton.addEventListener("click", () => setSelectedCellQuoteOverrides("quote"))`,
		`noQuoteSelectionRuleButton.addEventListener("click", () => setSelectedCellQuoteOverrides("no-quote"))`,
		`clearSelectionQuoteRuleButton.addEventListener("click", () => setSelectedCellQuoteOverrides("clear"))`,
		`applyQuoteColumnRuleButton.addEventListener("click", applyColumnQuoteRule)`,
		`applyQuoteColumnPatternRuleButton.addEventListener("click", applyColumnPatternQuoteRule)`,
		`previewQuoteValuePatternRuleButton.addEventListener("click", () => applyValuePatternQuoteRule(true))`,
		`applyQuoteValuePatternRuleButton.addEventListener("click", () => applyValuePatternQuoteRule(false))`,
		`applyQuoteRowRangeRuleButton.addEventListener("click", applyRowRangeQuoteRule)`,
		`applyQuoteCoordinateRuleButton.addEventListener("click", applyCoordinateQuoteRule)`,
		`remove.addEventListener("click", () => removeQuoteRule(item))`,
		"renderQuoteRuleList()",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing quote rule list frontend wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		".quote-rule-list",
		".quote-rule-item",
		".quote-rule-empty",
		".quote-rule-grid button",
		`.quote-rule-grid input[type="number"]`,
		`.quote-rule-grid input[type="text"]`,
	} {
		if !strings.Contains(css, snippet) {
			t.Fatalf("missing quote rule list styles %q", snippet)
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
		`"link.prompt.dialogAria": "打开自定义链接"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing custom link prompt wiring %q", snippet)
		}
	}
	for _, snippet := range []string{
		`id="linkConfirmDialog"`,
		`id="linkConfirmRemember"`,
		`id="confirmOpenLinkButton"`,
		`data-i18n-aria-label="link.prompt.dialogAria"`,
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
