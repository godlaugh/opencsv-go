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
		"row.hide",
		"row.showAll",
		"row.moveUp",
		"row.moveDown",
		"column.cut",
		"column.paste",
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
		"view.exportSavedSets",
		"view.importSavedSets",
		"view.sql",
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

func TestEditorGroupLayoutPersistenceIsWired(t *testing.T) {
	appJS, err := os.ReadFile(filepath.Join("..", "..", "web", "src", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appJS)
	for _, snippet := range []string{
		`const editorGroupLayoutStorageKey = "opencsv-editor-group-layout"`,
		"let editorGroupLayout = loadEditorGroupLayout()",
		"function loadEditorGroupLayout()",
		"function saveEditorGroupLayout()",
		"function editorGroupLayoutSignature(",
		"document.body.dataset.editorGroupLayoutRestored",
		"document.body.dataset.editorGroupLayoutSaved",
		"saveEditorGroupLayout();",
		"document.body.dataset.editorGroupLayout = editorGroupLayoutSignature()",
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing editor group layout persistence wiring %q", snippet)
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
		`"sql.context.copySelectionAs"`,
		`"sql.status.returned": "SQL 返回`,
		"applySQLLocaleStrings",
		"renderSQLHistory()",
		"renderSQLResult(lastSQLResult)",
		"tf(\"sql.status.returned\"",
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
		`tf("find.status.found"`,
		`tf("replace.status.noReplacements"`,
		`tf("replace.status.preview"`,
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
	} {
		if !strings.Contains(html, `id="`+id+`"`) {
			t.Fatalf("missing cancellation control %s", id)
		}
	}
	for _, snippet := range []string{
		"let activeFindController = null",
		"let activeSQLController = null",
		"function startFindProgress",
		"function cancelFindRequest",
		"function startSQLProgress",
		"function cancelSQLRequest",
		"signal: controller.signal",
		`document.body.dataset.findProgress = "canceled"`,
		`document.body.dataset.sqlProgress = "canceled"`,
		`"find.status.running": "正在查找`,
		`"sql.status.canceled": "SQL 已取消。"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing long-operation cancellation wiring %q", snippet)
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
	} {
		if !strings.Contains(html, snippet) {
			t.Fatalf("missing condition builder i18n html anchor %q", snippet)
		}
	}
	for _, snippet := range []string{
		`"condition.operator.contains"`,
		`"condition.operator.contains": "包含"`,
		`"condition.moveUp": "上移"`,
		`"condition.status.ready": "快速条件已就绪`,
		"applyConditionLocaleStrings",
		"setSelectOptionLabels(conditionOperator",
		"setSelectOptionLabels(conditionTemplateSelect",
		`tf("condition.status.ready"`,
		"moveConditionRow",
		`button.dataset.conditionAction === "move-up"`,
		"document.body.dataset.conditionReorder",
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
		`data-i18n-placeholder="profile.schemaNamePlaceholder"`,
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
		"applyProfileLocaleStrings",
		"renderProfile(profileData)",
		"syncProfileSchemaControls()",
		`summaryNode(t("profile.summary.rows")`,
		`t("profile.empty.noColumns")`,
		"profileSeverityLabel(issue.severity)",
		"profileQuickFixLabel(fix)",
		"uniqueHeaderNameForSourceCol",
		`tf("profile.quickFix.applied"`,
	} {
		if !strings.Contains(source, snippet) {
			t.Fatalf("missing profile i18n wiring %q", snippet)
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
