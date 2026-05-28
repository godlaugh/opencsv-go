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
		{ID: "settings.keyboardShortcuts", Title: "Keyboard Shortcuts", Description: "Customize command keyboard shortcuts.", Shortcut: "Cmd+,", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"keybindings", "shortcuts", "hotkeys"}},
		{ID: "workbench.toggleTheme", Title: "Toggle Theme", Description: "Cycle between light, dark, and system themes.", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"dark mode", "light mode", "color theme"}},
		{ID: "file.open", Title: "Open", Description: "Open a CSV or delimited text file.", Shortcut: "Cmd+O", Scope: ScopeGlobal, Implemented: true, Aliases: []string{"open file"}},
		{ID: "file.save", Title: "Save", Description: "Save the current document with backup policy.", Shortcut: "Cmd+S", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "file.exportView", Title: "Export View", Description: "Download the current filtered, sorted, and visible-column view.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"export filtered rows", "download csv", "export csv"}},
		{ID: "file.exportXlsx", Title: "Export XLSX", Description: "Download the current filtered, sorted, and visible-column view as an Excel workbook.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"excel export", "xlsx", "download excel"}},
		{ID: "file.toggleReadOnly", Title: "Toggle Read Only", Description: "Prevent or allow edits to the current document.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"readonly", "lock"}},
		{ID: "edit.undo", Title: "Undo", Description: "Undo the last editing operation.", Shortcut: "Cmd+Z", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.redo", Title: "Redo", Description: "Redo the last undone editing operation.", Shortcut: "Cmd+Shift+Z", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.copy", Title: "Copy", Description: "Copy the selected cells as tab-separated text.", Shortcut: "Cmd+C", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.copyWithHeaders", Title: "Copy with Headers", Description: "Copy the selected cells as tab-separated text with a header row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy header", "header copy"}},
		{ID: "selection.toggleHeaders", Title: "Toggle Header in Selection", Description: "Include or exclude column headers when copying the current selection.", Shortcut: "Shift+Alt+H", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"header selection", "include headers"}},
		{ID: "table.toggleHeaderRow", Title: "Toggle Header Row", Description: "Treat the first row as column headers or as regular data.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"header row", "no header", "set headers"}},
		{ID: "selection.cycleAggregate", Title: "Cycle Selection Aggregate", Description: "Switch the active status-bar aggregate for the current selection.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"status aggregate", "switch aggregate"}},
		{ID: "selection.copyAggregate", Title: "Copy Selection Aggregate", Description: "Copy the active status-bar aggregate value.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"copy aggregate", "copy sum", "copy average"}},
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
		{ID: "edit.replace", Title: "Replace", Description: "Replace the current match.", Shortcut: "Cmd+Alt+F", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "edit.replaceAll", Title: "Replace All", Description: "Replace all matches in the current document.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "row.insert", Title: "Insert Row", Description: "Insert an empty row before the selected row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"add row"}},
		{ID: "row.delete", Title: "Delete Row", Description: "Delete the selected row.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"remove row"}},
		{ID: "row.moveUp", Title: "Move Row Up", Description: "Move the selected row up one position.", Shortcut: "Shift+Alt+Up", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"row up", "move row"}},
		{ID: "row.moveDown", Title: "Move Row Down", Description: "Move the selected row down one position.", Shortcut: "Shift+Alt+Down", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"row down", "move row"}},
		{ID: "row.hide", Title: "Hide Row", Description: "Hide the selected row from the current view without deleting it.", Shortcut: "Shift+Alt+R", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"hide row"}},
		{ID: "column.insert", Title: "Insert Column", Description: "Insert an empty column before the selected column.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"add column"}},
		{ID: "column.delete", Title: "Delete Column", Description: "Delete the selected column.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"remove column"}},
		{ID: "column.moveLeft", Title: "Move Column Left", Description: "Move the selected column left one position.", Shortcut: "Shift+Alt+Left", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"column left", "move column"}},
		{ID: "column.moveRight", Title: "Move Column Right", Description: "Move the selected column right one position.", Shortcut: "Shift+Alt+Right", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"column right", "move column"}},
		{ID: "column.hide", Title: "Hide Column", Description: "Hide the selected column from the current view without deleting it.", Shortcut: "Shift+Alt+X", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"hide column"}},
		{ID: "table.transpose", Title: "Transpose", Description: "Swap rows and columns for the current document.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"swap rows columns"}},
		{ID: "table.dedupeRows", Title: "Remove Duplicate Rows", Description: "Remove duplicate data rows while preserving the header row.", Shortcut: "Shift+Alt+D", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"dedupe rows", "unique rows", "duplicate rows"}},
		{ID: "table.dedupeColumns", Title: "Remove Duplicate Columns", Description: "Remove duplicate columns whose full cell values match.", Shortcut: "Shift+Alt+C", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"dedupe columns", "unique columns", "duplicate columns"}},
		{ID: "text.trim", Title: "Trim Text", Description: "Trim leading and trailing whitespace in selected cells.", Shortcut: "Shift+Alt+T", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"trim cells", "strip whitespace"}},
		{ID: "text.uppercase", Title: "Uppercase Text", Description: "Convert selected cells to uppercase text.", Shortcut: "Shift+Alt+U", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"uppercase cells", "upper"}},
		{ID: "text.lowercase", Title: "Lowercase Text", Description: "Convert selected cells to lowercase text.", Shortcut: "Shift+Alt+L", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"lowercase cells", "lower"}},
		{ID: "cell.swap", Title: "Swap Cells", Description: "Swap the values of exactly two selected cells.", Shortcut: "Shift+Alt+W", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"swap cell values", "exchange cells"}},
		{ID: "cell.joinHorizontal", Title: "Join Cells", Description: "Join selected cells horizontally into the first selected cell.", Shortcut: "Shift+Alt+J", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"merge cells", "join horizontal"}},
		{ID: "cell.splitHorizontal", Title: "Split Cell", Description: "Split the active cell horizontally by separator text.", Shortcut: "Shift+Alt+S", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"split cells", "split horizontal"}},
		{ID: "cell.joinVertical", Title: "Join Cells Vertically", Description: "Join selected cells vertically into the top selected cell.", Shortcut: "Shift+Alt+V", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"merge cells vertically", "join vertical"}},
		{ID: "cell.splitVertical", Title: "Split Cell Vertically", Description: "Split the active cell downward by separator text.", Shortcut: "Shift+Alt+B", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"split cells vertically", "split vertical"}},
		{ID: "view.filter", Title: "Filter", Description: "Filter rows by text or SQL WHERE expression.", Shortcut: "Cmd+L", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "view.filterBuilder", Title: "Filter Builder", Description: "Build a column filter without writing SQL.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true, Aliases: []string{"condition filter", "visual filter", "filter condition"}},
		{ID: "view.sort", Title: "Sort", Description: "Sort the selected column.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "view.clear", Title: "Clear View", Description: "Clear active filter and sort view.", Scope: ScopeDocument, RequiresDoc: true, Implemented: true},
		{ID: "view.sql", Title: "SQL Console", Description: "Open the SQL console for the current workspace.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"query"}},
		{ID: "sql.copyResultCsv", Title: "Copy SQL Result as CSV", Description: "Copy the latest SQL result headers and rows as CSV text.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"copy sql result", "sql result csv"}},
		{ID: "sql.exportResultCsv", Title: "Export SQL Result as CSV", Description: "Download the latest SQL result headers and rows as a CSV file.", Scope: ScopeView, RequiresDoc: true, Implemented: true, Aliases: []string{"export sql result", "download sql result"}},
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
