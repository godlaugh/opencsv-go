package commands

type Scope string

const (
	ScopeGlobal   Scope = "global"
	ScopeDocument Scope = "document"
	ScopeGrid     Scope = "grid"
	ScopeView     Scope = "view"
)

type Command struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Shortcut    string   `json:"shortcut,omitempty"`
	Scope       Scope    `json:"scope"`
	RequiresDoc bool     `json:"requiresDocument"`
	Implemented bool     `json:"implemented"`
	Aliases     []string `json:"aliases,omitempty"`
}

func Registry() []Command {
	return []Command{
		{ID: "command.palette", Title: "Command Palette", Description: "Search and run available commands.", Shortcut: "Cmd+Shift+P", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"commands"}},
		{ID: "settings.open", Title: "Settings", Description: "Open application settings and file format preset management.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"preferences", "format presets"}},
		{ID: "settings.keyboardShortcuts", Title: "Keyboard Shortcuts", Description: "Customize command keyboard shortcuts.", Shortcut: "Cmd+,", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"keybindings", "shortcuts", "hotkeys"}},
		{ID: "settings.exportShortcuts", Title: "Export Keyboard Shortcuts", Description: "Download customized keyboard shortcuts as JSON.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"export keybindings", "backup shortcuts", "shortcut json"}},
		{ID: "settings.importShortcuts", Title: "Import Keyboard Shortcuts", Description: "Import customized keyboard shortcuts from JSON.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"import keybindings", "restore shortcuts", "shortcut json"}},
		{ID: "settings.checkForUpdates", Title: "Check for Updates", Description: "Check the configured update source when update checks are enabled.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"update check", "software update", "manual update"}},
		{ID: "workbench.toggleTheme", Title: "Toggle Theme", Description: "Cycle between light, dark, high contrast, and system themes.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"dark mode", "light mode", "high contrast", "color theme"}},
		{ID: "workbench.splitEditor", Title: "Split Editor Group", Description: "Open a second editor group for the current workspace.", Shortcut: "Cmd+\\", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"split editor", "split group", "new editor group"}},
		{ID: "workbench.toggleSplitOrientation", Title: "Toggle Split Orientation", Description: "Switch editor groups between side-by-side and stacked layouts.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"split horizontal", "split vertical", "editor group layout"}},
		{ID: "workbench.focusNextGroup", Title: "Focus Next Editor Group", Description: "Move focus between editor groups.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"next group", "focus group", "switch group"}},
		{ID: "workbench.moveTabToOtherGroup", Title: "Move Tab to Other Group", Description: "Move the current document tab into the other editor group.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"move tab", "move editor", "tab group"}},
		{ID: "workbench.closeSplitEditor", Title: "Close Split Editor Group", Description: "Close the secondary editor group.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"close split", "close group", "unsplit"}},
		{ID: "file.open", Title: "Open", Description: "Open a CSV or delimited text file.", Shortcut: "Cmd+O", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"open file"}},
		{ID: "file.save", Title: "Save", Description: "Save the current document with backup policy.", Shortcut: "Cmd+S", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "file.saveAs", Title: "Save As", Description: "Save the current document to a chosen local path.", Shortcut: "Cmd+Shift+S", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"save copy", "save to path"}},
		{ID: "file.restoreBackup", Title: "Restore From Backup", Description: "Restore the current document from one of its timestamped backups.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"backup history", "recover backup", "restore backup"}},
		{ID: "file.exportView", Title: "Export View", Description: "Download the current filtered, sorted, and visible-column view.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"export filtered rows", "download csv", "export csv"}},
		{ID: "file.exportXlsx", Title: "Export XLSX", Description: "Download the current filtered, sorted, and visible-column view as an Excel workbook.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"excel export", "xlsx", "download excel"}},
		{ID: "file.toggleReadOnly", Title: "Toggle Read Only", Description: "Prevent or allow edits to the current document.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"readonly", "lock"}},
		{ID: "format.savePreset", Title: "Save Format Preset", Description: "Save the current delimiter, header, and encoding detection as a reusable format preset.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"file format preset", "format preset"}},
		{ID: "format.applyPreset", Title: "Apply Format Preset", Description: "Reopen the current path with the selected format preset.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"apply preset", "open with preset"}},
		{ID: "format.saveRule", Title: "Save Format Rule", Description: "Attach a filename, path, or content rule to the current format preset for automatic application.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"format rule", "auto format"}},
		{ID: "format.exportPresets", Title: "Export Format Presets", Description: "Download saved file format presets and automatic rules as JSON.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"export formats", "backup format presets", "format preset json"}},
		{ID: "format.importPresets", Title: "Import Format Presets", Description: "Import file format presets and automatic rules from JSON.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"import formats", "restore format presets", "format preset json"}},
		{ID: "format.applyDelimiter", Title: "Apply Delimiter", Description: "Update the current document save delimiter.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"delimiter", "separator", "comma", "tab", "pipe"}},
		{ID: "format.applyEncoding", Title: "Apply Encoding", Description: "Update the current document save encoding.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"encoding", "charset", "bom", "utf-8", "shift-jis"}},
		{ID: "format.applyLineEnding", Title: "Apply Line Ending", Description: "Update the current document save line ending.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"line ending", "newline", "crlf", "lf"}},
		{ID: "format.applyQuotePolicy", Title: "Apply Quote Policy", Description: "Update quote character and quote mode for the current document save format.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"quote mode", "quote character"}},
		{ID: "format.switchWorksheet", Title: "Switch Worksheet", Description: "Focus the XLSX worksheet selector for the current workbook.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"worksheet", "sheet", "excel sheet", "xlsx sheet"}},
		{ID: "format.exportQuoteRules", Title: "Export Quote Rules", Description: "Download the current document quote policy and rule map as JSON.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"export quotes", "backup quote rules", "quote json"}},
		{ID: "format.importQuoteRules", Title: "Import Quote Rules", Description: "Import quote policy and rule map JSON into the current document.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"import quotes", "restore quote rules", "quote json"}},
		{ID: "format.toggleHeaderQuote", Title: "Toggle Header Quote", Description: "Toggle explicit quote preservation for the header row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"quote header", "header quote"}},
		{ID: "format.toggleBodyQuote", Title: "Toggle Body Quote", Description: "Toggle explicit quote preservation for body rows.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"quote body", "body quote"}},
		{ID: "format.toggleColumnQuote", Title: "Toggle Column Quote", Description: "Toggle explicit quote preservation for the active column.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"quote column", "column quote"}},
		{ID: "format.toggleCellQuote", Title: "Toggle Cell Quote", Description: "Toggle explicit quote preservation for the active cell.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"cell quote", "manual quote", "quote cell"}},
		{ID: "format.quoteSelection", Title: "Quote Selection", Description: "Mark every selected cell to be quoted on save.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"quote selected cells", "bulk quote", "manual quote selection"}},
		{ID: "format.noQuoteSelection", Title: "No Quote Selection", Description: "Mark every selected cell with an explicit no-quote override.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"unquote selected cells", "bulk no quote", "manual no quote"}},
		{ID: "format.clearSelectionQuote", Title: "Clear Selection Quote Overrides", Description: "Remove explicit quote and no-quote overrides from the selected cells.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"clear quote selection", "reset selected quote rules", "clear manual quotes"}},
		{ID: "edit.undo", Title: "Undo", Description: "Undo the last editing operation.", Shortcut: "Cmd+Z", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.redo", Title: "Redo", Description: "Redo the last undone editing operation.", Shortcut: "Cmd+Shift+Z", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.copy", Title: "Copy", Description: "Copy the selected cells as tab-separated text.", Shortcut: "Cmd+C", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.copyWithHeaders", Title: "Copy with Headers", Description: "Copy the selected cells as tab-separated text with a header row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy header", "header copy"}},
		{ID: "edit.copyRow", Title: "Copy Row", Description: "Copy the active row as tab-separated text.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy current row", "row copy"}},
		{ID: "edit.copyRowWithHeaders", Title: "Copy Row with Headers", Description: "Copy the active row with a header row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy row headers", "row copy headers"}},
		{ID: "edit.copyColumn", Title: "Copy Column", Description: "Copy the active column as tab-separated text.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy current column", "column copy"}},
		{ID: "edit.copyColumnWithHeader", Title: "Copy Column with Header", Description: "Copy the active column with its header value.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy column header", "column copy header"}},
		{ID: "selection.toggleHeaders", Title: "Toggle Header in Selection", Description: "Include or exclude column headers when copying the current selection.", Shortcut: "Shift+Alt+H", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"header selection", "include headers"}},
		{ID: "selection.addRange", Title: "Add Selection Range", Description: "Add the current range to a multi-range smart selection.", Shortcut: "Shift+Alt+M", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"smart selection add", "add range", "multi range selection"}},
		{ID: "selection.subtractRange", Title: "Subtract Selection Range", Description: "Remove the current range from the multi-range smart selection.", Shortcut: "Shift+Alt+-", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"smart selection subtract", "subtract range", "remove range"}},
		{ID: "selection.invert", Title: "Invert Selection", Description: "Invert the current selection across the table body.", Shortcut: "Shift+Alt+I", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"smart selection invert", "invert range"}},
		{ID: "selection.clearRanges", Title: "Clear Extra Selection Ranges", Description: "Clear multi-range smart selection state and keep the active cell selected.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"clear smart selection", "clear ranges"}},
		{ID: "table.toggleHeaderRow", Title: "Toggle Header Row", Description: "Treat the first row as column headers or as regular data.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"header row", "no header", "set headers"}},
		{ID: "selection.cycleAggregate", Title: "Cycle Selection Aggregate", Description: "Switch the active status-bar aggregate for the current selection.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"status aggregate", "switch aggregate"}},
		{ID: "selection.copyAggregate", Title: "Copy Selection Aggregate", Description: "Copy the active status-bar aggregate value.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy aggregate", "copy sum", "copy average"}},
		{ID: "selection.insertAggregateRow", Title: "Add Aggregate Row", Description: "Add a removable aggregate row for the current selection without changing file data.", Shortcut: "Shift+Alt+A", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"aggregate row", "add aggregate", "insert aggregate"}},
		{ID: "edit.copySelectedFormat", Title: "Copy Selected Format", Description: "Copy the selected cells using the toolbar copy format.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy format", "copy selected format"}},
		{ID: "edit.copyCsv", Title: "Copy as CSV", Description: "Copy the selected cells as CSV text.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"csv copy"}},
		{ID: "edit.copyMarkdown", Title: "Copy as Markdown", Description: "Copy the selected cells as a Markdown table.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"markdown copy"}},
		{ID: "edit.copyJson", Title: "Copy as JSON", Description: "Copy the selected cells as a JSON array.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"json copy"}},
		{ID: "edit.copyHtml", Title: "Copy as HTML", Description: "Copy the selected cells as an HTML table.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"html copy"}},
		{ID: "edit.copySql", Title: "Copy as SQL", Description: "Copy the selected cells as SQL VALUES.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"sql copy", "values copy"}},
		{ID: "edit.copyLatex", Title: "Copy as LaTeX", Description: "Copy the selected cells as a LaTeX tabular block.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"latex copy"}},
		{ID: "edit.copyRaw", Title: "Copy Raw", Description: "Copy selected raw values without CSV quoting.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"raw copy"}},
		{ID: "edit.paste", Title: "Paste", Description: "Paste a tab-separated cell matrix at the selected cell.", Shortcut: "Cmd+V", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "cell.fillDown", Title: "Fill Down", Description: "Auto-fill the selected range downward with detected number, date, or text patterns.", Shortcut: "Shift+Alt+F", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"auto fill", "autofill", "fill series"}},
		{ID: "cell.fillRight", Title: "Fill Right", Description: "Auto-fill the selected range to the right with detected number, date, or text patterns.", Shortcut: "Shift+Alt+G", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"auto fill right", "fill series right"}},
		{ID: "cell.fillUp", Title: "Fill Up", Description: "Auto-fill the selected range upward with detected number, date, or text patterns.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"auto fill up", "fill series up"}},
		{ID: "cell.fillLeft", Title: "Fill Left", Description: "Auto-fill the selected range to the left with detected number, date, or text patterns.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"auto fill left", "fill series left"}},
		{ID: "edit.find", Title: "Find", Description: "Find cells matching the search query.", Shortcut: "Cmd+F", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.replacePreview", Title: "Preview Replace", Description: "Preview the cells that would change before replacing text.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"replace preview", "preview replacement"}},
		{ID: "edit.replace", Title: "Replace", Description: "Replace the current match.", Shortcut: "Cmd+Alt+F", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.replaceAll", Title: "Replace All", Description: "Replace all matches in the current document.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "row.insert", Title: "Insert Row", Description: "Insert an empty row before the selected row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"add row"}},
		{ID: "row.delete", Title: "Delete Row", Description: "Delete the selected row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"remove row"}},
		{ID: "row.cut", Title: "Cut Row", Description: "Cut the active row for structural paste.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"cut row", "row clipboard"}},
		{ID: "row.paste", Title: "Paste Row", Description: "Paste the copied or cut row before the active row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"paste row", "insert copied row"}},
		{ID: "row.moveUp", Title: "Move Row Up", Description: "Move the selected row up one position.", Shortcut: "Shift+Alt+Up", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"row up", "move row"}},
		{ID: "row.moveDown", Title: "Move Row Down", Description: "Move the selected row down one position.", Shortcut: "Shift+Alt+Down", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"row down", "move row"}},
		{ID: "row.hide", Title: "Hide Row", Description: "Hide the selected row from the current view without deleting it.", Shortcut: "Shift+Alt+R", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"hide row"}},
		{ID: "column.insert", Title: "Insert Column", Description: "Insert an empty column before the selected column.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"add column"}},
		{ID: "column.delete", Title: "Delete Column", Description: "Delete the selected column.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"remove column"}},
		{ID: "column.cut", Title: "Cut Column", Description: "Cut the active column for structural paste.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"cut column", "column clipboard"}},
		{ID: "column.paste", Title: "Paste Column", Description: "Paste the copied or cut column before the active column.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"paste column", "insert copied column"}},
		{ID: "column.moveLeft", Title: "Move Column Left", Description: "Move the selected column left one position.", Shortcut: "Shift+Alt+Left", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"column left", "move column"}},
		{ID: "column.moveRight", Title: "Move Column Right", Description: "Move the selected column right one position.", Shortcut: "Shift+Alt+Right", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"column right", "move column"}},
		{ID: "column.hide", Title: "Hide Column", Description: "Hide the selected column from the current view without deleting it.", Shortcut: "Shift+Alt+X", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"hide column"}},
		{ID: "column.showAll", Title: "Show All Columns", Description: "Restore every hidden column in the current view.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"unhide columns", "restore columns", "show hidden columns"}},
		{ID: "table.transpose", Title: "Transpose", Description: "Swap rows and columns for the current document.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"swap rows columns"}},
		{ID: "table.dedupeRows", Title: "Remove Duplicate Rows", Description: "Remove duplicate data rows while preserving the header row.", Shortcut: "Shift+Alt+D", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"dedupe rows", "unique rows", "duplicate rows"}},
		{ID: "table.dedupeColumns", Title: "Remove Duplicate Columns", Description: "Remove duplicate columns whose full cell values match.", Shortcut: "Shift+Alt+C", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"dedupe columns", "unique columns", "duplicate columns"}},
		{ID: "text.trim", Title: "Trim Text", Description: "Trim leading and trailing whitespace in selected cells.", Shortcut: "Shift+Alt+T", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"trim cells", "strip whitespace"}},
		{ID: "text.uppercase", Title: "Uppercase Text", Description: "Convert selected cells to uppercase text.", Shortcut: "Shift+Alt+U", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"uppercase cells", "upper"}},
		{ID: "text.lowercase", Title: "Lowercase Text", Description: "Convert selected cells to lowercase text.", Shortcut: "Shift+Alt+L", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"lowercase cells", "lower"}},
		{ID: "text.snakeCase", Title: "Snake Case", Description: "Convert selected cells to snake_case names.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"snake_case", "snake case", "naming conversion"}},
		{ID: "text.kebabCase", Title: "Kebab Case", Description: "Convert selected cells to kebab-case names.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"kebab-case", "dash case", "naming conversion"}},
		{ID: "text.camelCase", Title: "Camel Case", Description: "Convert selected cells to camelCase names.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"camelCase", "camel case", "naming conversion"}},
		{ID: "text.pascalCase", Title: "Pascal Case", Description: "Convert selected cells to PascalCase names.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"PascalCase", "pascal case", "naming conversion"}},
		{ID: "text.titleCase", Title: "Title Case", Description: "Convert selected cells to Title Case names.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"titlecase", "start case", "naming conversion"}},
		{ID: "text.fullWidth", Title: "Full Width", Description: "Convert selected cells to full-width Unicode text.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"wide text", "zenkaku", "fullwidth"}},
		{ID: "text.halfWidth", Title: "Half Width", Description: "Convert selected cells to half-width Unicode text.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"narrow text", "hankaku", "halfwidth"}},
		{ID: "cell.swap", Title: "Swap Cells", Description: "Swap the values of exactly two selected cells.", Shortcut: "Shift+Alt+W", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"swap cell values", "exchange cells"}},
		{ID: "cell.swapUp", Title: "Swap Cell Up", Description: "Swap the active cell with the visible cell above it.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"swap up", "exchange with above", "swap cell above"}},
		{ID: "cell.swapDown", Title: "Swap Cell Down", Description: "Swap the active cell with the visible cell below it.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"swap down", "exchange with below", "swap cell below"}},
		{ID: "cell.swapLeft", Title: "Swap Cell Left", Description: "Swap the active cell with the visible cell to its left.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"swap left", "exchange with left", "swap cell left"}},
		{ID: "cell.swapRight", Title: "Swap Cell Right", Description: "Swap the active cell with the visible cell to its right.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"swap right", "exchange with right", "swap cell right"}},
		{ID: "cell.joinHorizontal", Title: "Join Cells", Description: "Join selected cells horizontally into the first selected cell.", Shortcut: "Shift+Alt+J", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"merge cells", "join horizontal"}},
		{ID: "cell.splitHorizontal", Title: "Split Cell", Description: "Split the active cell horizontally by separator text.", Shortcut: "Shift+Alt+S", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"split cells", "split horizontal"}},
		{ID: "cell.joinVertical", Title: "Join Cells Vertically", Description: "Join selected cells vertically into the top selected cell.", Shortcut: "Shift+Alt+V", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"merge cells vertically", "join vertical"}},
		{ID: "cell.splitVertical", Title: "Split Cell Vertically", Description: "Split the active cell downward by separator text.", Shortcut: "Shift+Alt+B", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"split cells vertically", "split vertical"}},
		{ID: "view.filter", Title: "Filter", Description: "Filter rows by text or SQL WHERE expression.", Shortcut: "Cmd+Shift+F", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "view.filterBuilder", Title: "Filter Builder", Description: "Build a column filter without writing SQL.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"condition filter", "visual filter", "filter condition"}},
		{ID: "view.saveFilter", Title: "Save Current Filter", Description: "Save the active filter expression for quick reuse.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"saved filter", "filter preset", "save filter"}},
		{ID: "view.sort", Title: "Sort", Description: "Sort the selected column.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "view.saveSort", Title: "Save Current Sort", Description: "Save the active sort keys for quick reuse.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"saved sort", "sort preset", "save sort"}},
		{ID: "view.clear", Title: "Clear View", Description: "Clear active filter and sort view.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "view.exportSavedSets", Title: "Export Saved Filters and Sorts", Description: "Download saved filter and sort sets as JSON.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"export filters", "export sorts", "saved view json"}},
		{ID: "view.importSavedSets", Title: "Import Saved Filters and Sorts", Description: "Import saved filter and sort sets from JSON.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"import filters", "import sorts", "saved view json"}},
		{ID: "view.sql", Title: "SQL Console", Description: "Open the SQL console for the current workspace.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"query"}},
		{ID: "sql.copyResultCsv", Title: "Copy SQL Result as CSV", Description: "Copy the latest SQL result headers and rows as CSV text.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"copy sql result", "sql result csv"}},
		{ID: "sql.copyResultTsv", Title: "Copy SQL Result as TSV", Description: "Copy the latest SQL result headers and rows as tab-separated text.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"copy sql result tsv", "sql result tab"}},
		{ID: "sql.copyResultMarkdown", Title: "Copy SQL Result as Markdown", Description: "Copy the latest SQL result headers and rows as a Markdown table.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"copy sql result markdown", "sql result md"}},
		{ID: "sql.copyResultJson", Title: "Copy SQL Result as JSON", Description: "Copy the latest SQL result headers and rows as JSON arrays.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"copy sql result json", "sql result json"}},
		{ID: "sql.exportResultCsv", Title: "Export SQL Result as CSV", Description: "Download the latest SQL result headers and rows as a CSV file.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"export sql result", "download sql result"}},
		{ID: "sql.exportResultTsv", Title: "Export SQL Result as TSV", Description: "Download the latest SQL result headers and rows as a TSV file.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"export sql result tsv", "download sql result tab"}},
		{ID: "sql.exportResultMarkdown", Title: "Export SQL Result as Markdown", Description: "Download the latest SQL result headers and rows as a Markdown table.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"export sql result markdown", "download sql result md"}},
		{ID: "sql.exportResultJson", Title: "Export SQL Result as JSON", Description: "Download the latest SQL result headers and rows as JSON arrays.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"export sql result json", "download sql result json"}},
		{ID: "view.profile", Title: "Profile", Description: "Show column profile and data quality metrics.", Scope: ScopeView, RequiresDoc: true, Implemented: true},
		{ID: "diff.open", Title: "Compare Files", Description: "Compare two CSV files side by side.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"diff", "compare"}},
	}
}

func ByID(id string) (Command, bool) {
	for _, command := range Registry() {
		if command.ID == id {
			return command, true
		}
	}
	return Command{}, false
}
