const canvas = document.querySelector("#grid");
const ctx = canvas.getContext("2d");
const headerCanvas = document.querySelector("#gridHeader");
const headerCtx = headerCanvas.getContext("2d");
const editor = document.querySelector("#editor");
const status = document.querySelector("#status");
const statusbar = document.querySelector("#statusbar");
const encodingStatus = document.querySelector("#encodingStatus");
const delimiterStatus = document.querySelector("#delimiterStatus");
const quoteStatus = document.querySelector("#quoteStatus");
const headerStatus = document.querySelector("#headerStatus");
const workbookSheetSelect = document.querySelector("#workbookSheetSelect");
const readOnlyStatus = document.querySelector("#readOnlyStatus");
const externalStatus = document.querySelector("#externalStatus");
const aggregateStatusButton = document.querySelector("#aggregateStatusButton");
const formatPresetSelect = document.querySelector("#formatPresetSelect");
const saveFormatPresetButton = document.querySelector("#saveFormatPresetButton");
const applyFormatPresetButton = document.querySelector("#applyFormatPresetButton");
const formatRuleInput = document.querySelector("#formatRuleInput");
const saveFormatRuleButton = document.querySelector("#saveFormatRuleButton");
const columnName = document.querySelector("#columnName");
const hiddenColumnSummary = document.querySelector("#hiddenColumnSummary");
const showHiddenColumnsButton = document.querySelector("#showHiddenColumnsButton");
const columnType = document.querySelector("#columnType");
const columnNulls = document.querySelector("#columnNulls");
const columnDistinct = document.querySelector("#columnDistinct");
const cellQuote = document.querySelector("#cellQuote");
const formatSummary = document.querySelector("#formatSummary");
const cellValueInput = document.querySelector("#cellValueInput");
const cellValueStats = document.querySelector("#cellValueStats");
const openCellValueDialogButton = document.querySelector("#openCellValueDialogButton");
const toggleCellValueLargeButton = document.querySelector("#toggleCellValueLargeButton");
const revertCellValueButton = document.querySelector("#revertCellValueButton");
const applyCellValueButton = document.querySelector("#applyCellValueButton");
const openCellLink = document.querySelector("#openCellLink");
const cellLinkAudit = document.querySelector("#cellLinkAudit");
const cellLinkProtocol = document.querySelector("#cellLinkProtocol");
const cellLinkTarget = document.querySelector("#cellLinkTarget");
const linkConfirmDialog = document.querySelector("#linkConfirmDialog");
const linkConfirmProtocol = document.querySelector("#linkConfirmProtocol");
const linkConfirmTarget = document.querySelector("#linkConfirmTarget");
const linkConfirmRemember = document.querySelector("#linkConfirmRemember");
const cancelOpenLinkButton = document.querySelector("#cancelOpenLinkButton");
const confirmOpenLinkButton = document.querySelector("#confirmOpenLinkButton");
const cellSeparatorInput = document.querySelector("#cellSeparatorInput");
const toggleCellQuoteButton = document.querySelector("#toggleCellQuoteButton");
const fileList = document.querySelector("#fileList");
const openForm = document.querySelector("#openForm");
const openButton = document.querySelector("#openButton");
const pathInput = document.querySelector("#pathInput");
const fileInput = document.querySelector("#fileInput");
const commandInput = document.querySelector("#commandInput");
const replaceInput = document.querySelector("#replaceInput");
const copyFormatSelect = document.querySelector("#copyFormatSelect");
const copyFormatButton = document.querySelector("#copyFormatButton");
const copyHeadersToggle = document.querySelector("#copyHeadersToggle");
const filterButton = document.querySelector("#filterButton");
const recentFilterSelect = document.querySelector("#recentFilterSelect");
const filterBuilderButton = document.querySelector("#filterBuilderButton");
const sortButton = document.querySelector("#sortButton");
const findButton = document.querySelector("#findButton");
const cancelFindButton = document.querySelector("#cancelFindButton");
const findPrevButton = document.querySelector("#findPrevButton");
const findNextButton = document.querySelector("#findNextButton");
const findMatchCount = document.querySelector("#findMatchCount");
const viewChips = document.querySelector("#viewChips");
const replacePreviewButton = document.querySelector("#replacePreviewButton");
const replaceButton = document.querySelector("#replaceButton");
const replaceAllButton = document.querySelector("#replaceAllButton");
const findCaseButton = document.querySelector("#findCaseButton");
const findWholeButton = document.querySelector("#findWholeButton");
const findRegexButton = document.querySelector("#findRegexButton");
const findSelectionButton = document.querySelector("#findSelectionButton");
const diffButton = document.querySelector("#diffButton");
const exportFormatSelect = document.querySelector("#exportFormatSelect");
const exportButton = document.querySelector("#exportButton");
const exportXlsxButton = document.querySelector("#exportXlsxButton");
const saveButton = document.querySelector("#saveButton");
const saveAsButton = document.querySelector("#saveAsButton");
const restoreBackupButton = document.querySelector("#restoreBackupButton");
const readOnlyButton = document.querySelector("#readOnlyButton");
const commandButtons = Array.from(document.querySelectorAll("[data-command]"));
const gridFrame = document.querySelector(".grid-frame");
const editorGroups = document.querySelector("#editorGroups");
const primaryEditorGroup = document.querySelector("#primaryEditorGroup");
const secondaryEditorGroup = document.querySelector("#secondaryEditorGroup");
const primaryEditorGroupBody = document.querySelector("#primaryEditorGroupBody");
const secondaryEditorGroupBody = document.querySelector("#secondaryEditorGroupBody");
const primaryEditorGroupTitle = document.querySelector("#primaryEditorGroupTitle");
const primaryEditorGroupMeta = document.querySelector("#primaryEditorGroupMeta");
const secondaryEditorGroupTitle = document.querySelector("#secondaryEditorGroupTitle");
const secondaryEditorGroupMeta = document.querySelector("#secondaryEditorGroupMeta");
const primaryEditorTab = document.querySelector("#primaryEditorTab");
const secondaryEditorTab = document.querySelector("#secondaryEditorTab");
const primaryEditorGroupEmpty = document.querySelector("#primaryEditorGroupEmpty");
const secondaryEditorGroupEmpty = document.querySelector("#secondaryEditorGroupEmpty");
const secondaryEditorGroupPreview = document.querySelector("#secondaryEditorGroupPreview");
const secondaryEditorGroupFile = document.querySelector("#secondaryEditorGroupFile");
const secondaryEditorGroupView = document.querySelector("#secondaryEditorGroupView");
const secondaryEditorGroupSelection = document.querySelector("#secondaryEditorGroupSelection");
const splitEditorButton = document.querySelector("#splitEditorButton");
const splitOrientationButton = document.querySelector("#splitOrientationButton");
const focusEditorGroupButton = document.querySelector("#focusEditorGroupButton");
const moveEditorTabButton = document.querySelector("#moveEditorTabButton");
const closeSplitEditorButton = document.querySelector("#closeSplitEditorButton");
const focusSecondaryGroupButton = document.querySelector("#focusSecondaryGroupButton");
const closeSecondaryGroupButton = document.querySelector("#closeSecondaryGroupButton");
const gridTabButton = document.querySelector("#gridTabButton");
const sqlTabButton = document.querySelector("#sqlTabButton");
const profileTabButton = document.querySelector("#profileTabButton");
const settingsTabButton = document.querySelector("#settingsTabButton");
const shortcutsTabButton = document.querySelector("#shortcutsTabButton");
const paletteButton = document.querySelector("#paletteButton");
const sqlPanel = document.querySelector("#sqlPanel");
const sqlInput = document.querySelector("#sqlInput");
const runSqlButton = document.querySelector("#runSqlButton");
const cancelSqlButton = document.querySelector("#cancelSqlButton");
const sqlResultFormatSelect = document.querySelector("#sqlResultFormatSelect");
const openSqlResultButton = document.querySelector("#openSqlResultButton");
const copySqlResultButton = document.querySelector("#copySqlResultButton");
const exportSqlResultButton = document.querySelector("#exportSqlResultButton");
const sqlHistoryList = document.querySelector("#sqlHistoryList");
const clearSqlHistoryButton = document.querySelector("#clearSqlHistoryButton");
const sqlResult = document.querySelector("#sqlResult");
const diffPanel = document.querySelector("#diffPanel");
const leftDiffPath = document.querySelector("#leftDiffPath");
const rightDiffPath = document.querySelector("#rightDiffPath");
const runDiffButton = document.querySelector("#runDiffButton");
const diffIgnoreCase = document.querySelector("#diffIgnoreCase");
const diffIgnoreCaseLabel = document.querySelector("#diffIgnoreCaseLabel");
const diffIgnoreWhitespace = document.querySelector("#diffIgnoreWhitespace");
const diffIgnoreWhitespaceLabel = document.querySelector("#diffIgnoreWhitespaceLabel");
const diffMatchHeaders = document.querySelector("#diffMatchHeaders");
const diffMatchHeadersLabel = document.querySelector("#diffMatchHeadersLabel");
const diffKeyColumn = document.querySelector("#diffKeyColumn");
const diffIgnoredColumns = document.querySelector("#diffIgnoredColumns");
const diffPrevButton = document.querySelector("#diffPrevButton");
const diffNextButton = document.querySelector("#diffNextButton");
const diffStatusFilter = document.querySelector("#diffStatusFilter");
const diffSearchInput = document.querySelector("#diffSearchInput");
const diffSortSelect = document.querySelector("#diffSortSelect");
const diffSortDirectionButton = document.querySelector("#diffSortDirectionButton");
const diffClearFilterButton = document.querySelector("#diffClearFilterButton");
const diffRightTextMode = document.querySelector("#diffRightTextMode");
const diffRightTextModeLabel = document.querySelector("#diffRightTextModeLabel");
const diffRightText = document.querySelector("#diffRightText");
const diffPasteClipboardButton = document.querySelector("#diffPasteClipboardButton");
const diffSummary = document.querySelector("#diffSummary");
const diffResult = document.querySelector("#diffResult");
const profilePanel = document.querySelector("#profilePanel");
const runProfileButton = document.querySelector("#runProfileButton");
const cancelProfileButton = document.querySelector("#cancelProfileButton");
const toggleProfileSchemaButton = document.querySelector("#toggleProfileSchemaButton");
const saveProfileSchemaButton = document.querySelector("#saveProfileSchemaButton");
const resetProfileSchemaButton = document.querySelector("#resetProfileSchemaButton");
const profileSchemaStatus = document.querySelector("#profileSchemaStatus");
const profileSchemaNameInput = document.querySelector("#profileSchemaNameInput");
const savedProfileSchemaSelect = document.querySelector("#savedProfileSchemaSelect");
const applyProfileSchemaButton = document.querySelector("#applyProfileSchemaButton");
const deleteProfileSchemaButton = document.querySelector("#deleteProfileSchemaButton");
const exportProfileSchemasButton = document.querySelector("#exportProfileSchemasButton");
const showImportProfileSchemasButton = document.querySelector("#showImportProfileSchemasButton");
const profileSchemaImportPanel = document.querySelector("#profileSchemaImportPanel");
const profileSchemaImportText = document.querySelector("#profileSchemaImportText");
const applyProfileSchemaImportButton = document.querySelector("#applyProfileSchemaImportButton");
const cancelProfileSchemaImportButton = document.querySelector("#cancelProfileSchemaImportButton");
const profileSchemaInput = document.querySelector("#profileSchemaInput");
const profileRuleNameInput = document.querySelector("#profileRuleNameInput");
const profileRuleColumnInput = document.querySelector("#profileRuleColumnInput");
const profileRuleTypeSelect = document.querySelector("#profileRuleTypeSelect");
const profileRuleRequiredToggle = document.querySelector("#profileRuleRequiredToggle");
const profileRuleUniqueToggle = document.querySelector("#profileRuleUniqueToggle");
const profileRuleEnumInput = document.querySelector("#profileRuleEnumInput");
const profileRulePatternInput = document.querySelector("#profileRulePatternInput");
const profileRuleDefaultInput = document.querySelector("#profileRuleDefaultInput");
const profileRuleMinInput = document.querySelector("#profileRuleMinInput");
const profileRuleMaxInput = document.querySelector("#profileRuleMaxInput");
const profileRuleMinLengthInput = document.querySelector("#profileRuleMinLengthInput");
const profileRuleMaxLengthInput = document.querySelector("#profileRuleMaxLengthInput");
const saveProfileRuleButton = document.querySelector("#saveProfileRuleButton");
const clearProfileRuleButton = document.querySelector("#clearProfileRuleButton");
const profileRuleList = document.querySelector("#profileRuleList");
const profileSummary = document.querySelector("#profileSummary");
const profileColumns = document.querySelector("#profileColumns");
const profileIssues = document.querySelector("#profileIssues");
const settingsPanel = document.querySelector("#settingsPanel");
const settingsSearch = document.querySelector("#settingsSearch");
const settingsEmpty = document.querySelector("#settingsEmpty");
const settingsModifiedBadge = document.querySelector("#settingsModifiedBadge");
const settingsLocaleTitle = document.querySelector("#settingsLocaleTitle");
const settingsLocaleSummary = document.querySelector("#settingsLocaleSummary");
const localeSelect = document.querySelector("#localeSelect");
const localeSelectLabel = document.querySelector("#localeSelectLabel");
const resetLocaleButton = document.querySelector("#resetLocaleButton");
const settingsQuoteSummary = document.querySelector("#settingsQuoteSummary");
const delimiterSelect = document.querySelector("#delimiterSelect");
const delimiterCustomInput = document.querySelector("#delimiterCustomInput");
const applyDelimiterButton = document.querySelector("#applyDelimiterButton");
const encodingSelect = document.querySelector("#encodingSelect");
const applyEncodingButton = document.querySelector("#applyEncodingButton");
const quoteModeSelect = document.querySelector("#quoteModeSelect");
const quoteCharSelect = document.querySelector("#quoteCharSelect");
const quoteEscapeSelect = document.querySelector("#quoteEscapeSelect");
const applyQuotePolicyButton = document.querySelector("#applyQuotePolicyButton");
const lineEndingSelect = document.querySelector("#lineEndingSelect");
const applyLineEndingButton = document.querySelector("#applyLineEndingButton");
const quoteHeaderRuleToggle = document.querySelector("#quoteHeaderRuleToggle");
const quoteBodyRuleToggle = document.querySelector("#quoteBodyRuleToggle");
const quoteColumnRuleToggle = document.querySelector("#quoteColumnRuleToggle");
const addQuoteCellRuleButton = document.querySelector("#addQuoteCellRuleButton");
const addNoQuoteCellRuleButton = document.querySelector("#addNoQuoteCellRuleButton");
const clearCellQuoteRuleButton = document.querySelector("#clearCellQuoteRuleButton");
const quoteSelectionRuleButton = document.querySelector("#quoteSelectionRuleButton");
const noQuoteSelectionRuleButton = document.querySelector("#noQuoteSelectionRuleButton");
const clearSelectionQuoteRuleButton = document.querySelector("#clearSelectionQuoteRuleButton");
const quoteRuleColumnInput = document.querySelector("#quoteRuleColumnInput");
const quoteRuleColumnActionSelect = document.querySelector("#quoteRuleColumnActionSelect");
const applyQuoteColumnRuleButton = document.querySelector("#applyQuoteColumnRuleButton");
const quoteRuleColumnPatternInput = document.querySelector("#quoteRuleColumnPatternInput");
const quoteRuleColumnPatternModeSelect = document.querySelector("#quoteRuleColumnPatternModeSelect");
const quoteRuleColumnPatternActionSelect = document.querySelector("#quoteRuleColumnPatternActionSelect");
const applyQuoteColumnPatternRuleButton = document.querySelector("#applyQuoteColumnPatternRuleButton");
const quoteRuleValuePatternInput = document.querySelector("#quoteRuleValuePatternInput");
const quoteRuleValuePatternModeSelect = document.querySelector("#quoteRuleValuePatternModeSelect");
const quoteRuleValuePatternActionSelect = document.querySelector("#quoteRuleValuePatternActionSelect");
const applyQuoteValuePatternRuleButton = document.querySelector("#applyQuoteValuePatternRuleButton");
const quoteRuleRowStartInput = document.querySelector("#quoteRuleRowStartInput");
const quoteRuleRowEndInput = document.querySelector("#quoteRuleRowEndInput");
const quoteRuleRowActionSelect = document.querySelector("#quoteRuleRowActionSelect");
const applyQuoteRowRangeRuleButton = document.querySelector("#applyQuoteRowRangeRuleButton");
const quoteRuleRowInput = document.querySelector("#quoteRuleRowInput");
const quoteRuleColInput = document.querySelector("#quoteRuleColInput");
const quoteRuleActionSelect = document.querySelector("#quoteRuleActionSelect");
const applyQuoteCoordinateRuleButton = document.querySelector("#applyQuoteCoordinateRuleButton");
const exportQuoteRulesButton = document.querySelector("#exportQuoteRulesButton");
const showImportQuoteRulesButton = document.querySelector("#showImportQuoteRulesButton");
const resetQuotePolicyButton = document.querySelector("#resetQuotePolicyButton");
const quoteRuleList = document.querySelector("#quoteRuleList");
const quoteRulesImportPanel = document.querySelector("#quoteRulesImportPanel");
const quoteRulesImportText = document.querySelector("#quoteRulesImportText");
const applyQuoteRulesImportButton = document.querySelector("#applyQuoteRulesImportButton");
const cancelQuoteRulesImportButton = document.querySelector("#cancelQuoteRulesImportButton");
const settingsThemeSummary = document.querySelector("#settingsThemeSummary");
const themePreferenceSelect = document.querySelector("#themePreferenceSelect");
const resetThemePreferenceButton = document.querySelector("#resetThemePreferenceButton");
const settingsRefreshSummary = document.querySelector("#settingsRefreshSummary");
const autoRefreshModeSelect = document.querySelector("#autoRefreshModeSelect");
const resetAutoRefreshButton = document.querySelector("#resetAutoRefreshButton");
const settingsUpdateSummary = document.querySelector("#settingsUpdateSummary");
const updateCheckModeSelect = document.querySelector("#updateCheckModeSelect");
const checkForUpdatesButton = document.querySelector("#checkForUpdatesButton");
const resetUpdateChecksButton = document.querySelector("#resetUpdateChecksButton");
const updateCheckResult = document.querySelector("#updateCheckResult");
const settingsLinkSummary = document.querySelector("#settingsLinkSummary");
const customLinkProtocolsInput = document.querySelector("#customLinkProtocolsInput");
const applyCustomLinkProtocolsButton = document.querySelector("#applyCustomLinkProtocolsButton");
const resetCustomLinkProtocolsButton = document.querySelector("#resetCustomLinkProtocolsButton");
const settingsAppearanceSummary = document.querySelector("#settingsAppearanceSummary");
const gridFontSelect = document.querySelector("#gridFontSelect");
const gridFontSizeInput = document.querySelector("#gridFontSizeInput");
const gridRowHeightSelect = document.querySelector("#gridRowHeightSelect");
const textEditorFontSelect = document.querySelector("#textEditorFontSelect");
const textEditorFontSizeInput = document.querySelector("#textEditorFontSizeInput");
const resetEditorAppearanceButton = document.querySelector("#resetEditorAppearanceButton");
const settingsFormatSummary = document.querySelector("#settingsFormatSummary");
const settingsPresetList = document.querySelector("#settingsPresetList");
const settingsViewsSummary = document.querySelector("#settingsViewsSummary");
const settingsSavedViewList = document.querySelector("#settingsSavedViewList");
const filterHistoryLimitSelect = document.querySelector("#filterHistoryLimitSelect");
const clearRecentFiltersButton = document.querySelector("#clearRecentFiltersButton");
const resetSavedViewsButton = document.querySelector("#resetSavedViewsButton");
const exportSavedViewsButton = document.querySelector("#exportSavedViewsButton");
const showImportSavedViewsButton = document.querySelector("#showImportSavedViewsButton");
const savedViewImportPanel = document.querySelector("#savedViewImportPanel");
const savedViewImportText = document.querySelector("#savedViewImportText");
const applySavedViewImportButton = document.querySelector("#applySavedViewImportButton");
const cancelSavedViewImportButton = document.querySelector("#cancelSavedViewImportButton");
const resetFormatPresetsButton = document.querySelector("#resetFormatPresetsButton");
const exportFormatPresetsButton = document.querySelector("#exportFormatPresetsButton");
const showImportFormatPresetsButton = document.querySelector("#showImportFormatPresetsButton");
const formatPresetImportPanel = document.querySelector("#formatPresetImportPanel");
const formatPresetImportText = document.querySelector("#formatPresetImportText");
const applyFormatPresetImportButton = document.querySelector("#applyFormatPresetImportButton");
const cancelFormatPresetImportButton = document.querySelector("#cancelFormatPresetImportButton");
const commandPalette = document.querySelector("#commandPalette");
const paletteInput = document.querySelector("#paletteInput");
const paletteList = document.querySelector("#paletteList");
const shortcutsPanel = document.querySelector("#shortcutsPanel");
const shortcutPresetSelect = document.querySelector("#shortcutPresetSelect");
const applyShortcutPresetButton = document.querySelector("#applyShortcutPresetButton");
const shortcutsSearch = document.querySelector("#shortcutsSearch");
const shortcutsList = document.querySelector("#shortcutsList");
const closeShortcutsButton = document.querySelector("#closeShortcutsButton");
const exportShortcutsButton = document.querySelector("#exportShortcutsButton");
const showImportShortcutsButton = document.querySelector("#showImportShortcutsButton");
const resetAllShortcutsButton = document.querySelector("#resetAllShortcutsButton");
const shortcutImportPanel = document.querySelector("#shortcutImportPanel");
const shortcutImportText = document.querySelector("#shortcutImportText");
const applyShortcutImportButton = document.querySelector("#applyShortcutImportButton");
const cancelShortcutImportButton = document.querySelector("#cancelShortcutImportButton");
const aboutSupportPanel = document.querySelector("#aboutSupportPanel");
const closeAboutSupportButton = document.querySelector("#closeAboutSupportButton");
const aboutUpdateSummary = document.querySelector("#aboutUpdateSummary");
const aboutUpdateResult = document.querySelector("#aboutUpdateResult");
const aboutOpenSettingsButton = document.querySelector("#aboutOpenSettingsButton");
const aboutKeyboardShortcutsButton = document.querySelector("#aboutKeyboardShortcutsButton");
const aboutCheckUpdatesButton = document.querySelector("#aboutCheckUpdatesButton");
const saveAsPanel = document.querySelector("#saveAsPanel");
const saveAsPath = document.querySelector("#saveAsPath");
const closeSaveAsButton = document.querySelector("#closeSaveAsButton");
const cancelSaveAsButton = document.querySelector("#cancelSaveAsButton");
const confirmSaveAsButton = document.querySelector("#confirmSaveAsButton");
const restoreBackupPanel = document.querySelector("#restoreBackupPanel");
const restoreBackupList = document.querySelector("#restoreBackupList");
const closeRestoreBackupButton = document.querySelector("#closeRestoreBackupButton");
const cancelRestoreBackupButton = document.querySelector("#cancelRestoreBackupButton");
const refreshRestoreBackupsButton = document.querySelector("#refreshRestoreBackupsButton");
const dirtyClosePanel = document.querySelector("#dirtyClosePanel");
const closeDirtyCloseButton = document.querySelector("#closeDirtyCloseButton");
const cancelDirtyCloseButton = document.querySelector("#cancelDirtyCloseButton");
const discardDirtyCloseButton = document.querySelector("#discardDirtyCloseButton");
const dirtyCloseName = document.querySelector("#dirtyCloseName");
const dirtyCloseDetail = document.querySelector("#dirtyCloseDetail");
const externalChangePanel = document.querySelector("#externalChangePanel");
const externalChangeTitle = document.querySelector("#externalChangeTitle");
const externalChangeMessage = document.querySelector("#externalChangeMessage");
const closeExternalChangeButton = document.querySelector("#closeExternalChangeButton");
const keepExternalChangeButton = document.querySelector("#keepExternalChangeButton");
const compareExternalChangeButton = document.querySelector("#compareExternalChangeButton");
const saveAsExternalChangeButton = document.querySelector("#saveAsExternalChangeButton");
const overwriteExternalChangeButton = document.querySelector("#overwriteExternalChangeButton");
const reloadExternalChangeButton = document.querySelector("#reloadExternalChangeButton");
const transposeConfirmPanel = document.querySelector("#transposeConfirmPanel");
const closeTransposeConfirmButton = document.querySelector("#closeTransposeConfirmButton");
const cancelTransposeConfirmButton = document.querySelector("#cancelTransposeConfirmButton");
const confirmTransposeButton = document.querySelector("#confirmTransposeButton");
const transposeConfirmSummary = document.querySelector("#transposeConfirmSummary");
const transposeConfirmDetail = document.querySelector("#transposeConfirmDetail");
const transposeConfirmWarning = document.querySelector("#transposeConfirmWarning");
const splitConfirmPanel = document.querySelector("#splitConfirmPanel");
const closeSplitConfirmButton = document.querySelector("#closeSplitConfirmButton");
const cancelSplitConfirmButton = document.querySelector("#cancelSplitConfirmButton");
const confirmSplitConfirmButton = document.querySelector("#confirmSplitConfirmButton");
const splitConfirmSummary = document.querySelector("#splitConfirmSummary");
const splitConfirmDetail = document.querySelector("#splitConfirmDetail");
const splitConfirmWarning = document.querySelector("#splitConfirmWarning");
const cellValueDialog = document.querySelector("#cellValueDialog");
const closeCellValueDialogButton = document.querySelector("#closeCellValueDialogButton");
const cellValueDialogInput = document.querySelector("#cellValueDialogInput");
const cellValueDialogStats = document.querySelector("#cellValueDialogStats");
const revertCellValueDialogButton = document.querySelector("#revertCellValueDialogButton");
const applyCellValueDialogButton = document.querySelector("#applyCellValueDialogButton");
const gridContextMenu = document.querySelector("#gridContextMenu");
const conditionPanel = document.querySelector("#conditionPanel");
const conditionColumn = document.querySelector("#conditionColumn");
const conditionOperator = document.querySelector("#conditionOperator");
const conditionValue = document.querySelector("#conditionValue");
const conditionTemplateSelect = document.querySelector("#conditionTemplateSelect");
const useConditionTemplateButton = document.querySelector("#useConditionTemplateButton");
const clearConditionGroupsButton = document.querySelector("#clearConditionGroupsButton");
const conditionPrimaryOpenGroup = document.querySelector("#conditionPrimaryOpenGroup");
const conditionPrimaryCloseGroup = document.querySelector("#conditionPrimaryCloseGroup");
const conditionRows = document.querySelector("#conditionRows");
const closeConditionButton = document.querySelector("#closeConditionButton");
const addConditionButton = document.querySelector("#addConditionButton");
const groupAllConditionsButton = document.querySelector("#groupAllConditionsButton");
const conditionAllAndButton = document.querySelector("#conditionAllAndButton");
const conditionAllOrButton = document.querySelector("#conditionAllOrButton");
const applyConditionButton = document.querySelector("#applyConditionButton");
const clearConditionButton = document.querySelector("#clearConditionButton");
const columnMenu = document.querySelector("#columnMenu");
const columnMenuTitle = document.querySelector("#columnMenuTitle");
const columnSortAsc = document.querySelector("#columnSortAsc");
const columnSortDesc = document.querySelector("#columnSortDesc");
const columnSortLengthAsc = document.querySelector("#columnSortLengthAsc");
const columnSortLengthDesc = document.querySelector("#columnSortLengthDesc");
const columnAddSortAsc = document.querySelector("#columnAddSortAsc");
const columnAddSortDesc = document.querySelector("#columnAddSortDesc");
const columnStableSortToggle = document.querySelector("#columnStableSortToggle");
const columnFilterValue = document.querySelector("#columnFilterValue");
const columnFilterBlank = document.querySelector("#columnFilterBlank");
const columnFilterFilled = document.querySelector("#columnFilterFilled");
const columnTransformMode = document.querySelector("#columnTransformMode");
const columnTransformApply = document.querySelector("#columnTransformApply");
const columnMenuValues = document.querySelector("#columnMenuValues");
const columnApplyValues = document.querySelector("#columnApplyValues");
const columnHide = document.querySelector("#columnHide");
const columnShowAll = document.querySelector("#columnShowAll");
const columnClearView = document.querySelector("#columnClearView");

let rowHeight = 28;
const headerHeight = cssPixelValue("--grid-header-height", 48);
const rowNumberWidth = 58;
const windowSize = 160;
const minColumnWidth = 92;
const maxColumnWidth = 420;
const columnReorderDragThreshold = 8;
const rowReorderDragThreshold = 8;
const autoFitColumnWidthOnEdit = true;
const transposeConfirmThresholdCells = 50_000;
const transposeMaxCells = 1_000_000;
const recentFilesStorageKey = "opencsv-recent-files";
const recentFilesLimit = 8;
const commandHistoryStorageKey = "opencsv-command-history";
const commandHistoryLimit = 8;
const sqlResultFormats = ["csv", "tsv", "markdown", "json", "html", "sql", "latex", "raw"];
const sqlResultFormatLabels = {
  csv: "CSV",
  tsv: "TSV",
  markdown: "Markdown",
  json: "JSON",
  html: "HTML",
  sql: "SQL",
  latex: "LaTeX",
  raw: "Raw",
};
const copyRowFormatCommands = new Map([
  ["edit.copyRowCsv", "csv"],
  ["edit.copyRowMarkdown", "markdown"],
  ["edit.copyRowJson", "json"],
  ["edit.copyRowHtml", "html"],
  ["edit.copyRowSql", "sql"],
  ["edit.copyRowLatex", "latex"],
  ["edit.copyRowRaw", "raw"],
]);
const copyColumnFormatCommands = new Map([
  ["edit.copyColumnCsv", "csv"],
  ["edit.copyColumnMarkdown", "markdown"],
  ["edit.copyColumnJson", "json"],
  ["edit.copyColumnHtml", "html"],
  ["edit.copyColumnSql", "sql"],
  ["edit.copyColumnLatex", "latex"],
  ["edit.copyColumnRaw", "raw"],
]);
const sqlResultFormatExtensions = {
  csv: "csv",
  tsv: "tsv",
  markdown: "md",
  json: "json",
  html: "html",
  sql: "sql",
  latex: "tex",
  raw: "txt",
};
const sqlResultFormatContentTypes = {
  csv: "text/csv;charset=utf-8",
  tsv: "text/tab-separated-values;charset=utf-8",
  markdown: "text/markdown;charset=utf-8",
  json: "application/json;charset=utf-8",
  html: "text/html;charset=utf-8",
  sql: "application/sql;charset=utf-8",
  latex: "application/x-latex;charset=utf-8",
  raw: "text/plain;charset=utf-8",
};
const sqlResultCommandFormats = {
  "sql.copyResultCsv": { action: "copy", format: "csv" },
  "sql.copyResultTsv": { action: "copy", format: "tsv" },
  "sql.copyResultMarkdown": { action: "copy", format: "markdown" },
  "sql.copyResultJson": { action: "copy", format: "json" },
  "sql.copyResultHtml": { action: "copy", format: "html" },
  "sql.copyResultSql": { action: "copy", format: "sql" },
  "sql.copyResultLatex": { action: "copy", format: "latex" },
  "sql.copyResultRaw": { action: "copy", format: "raw" },
  "sql.copySelectionCsv": { action: "copySelection", format: "csv" },
  "sql.copySelectionTsv": { action: "copySelection", format: "tsv" },
  "sql.copySelectionMarkdown": { action: "copySelection", format: "markdown" },
  "sql.copySelectionJson": { action: "copySelection", format: "json" },
  "sql.copySelectionHtml": { action: "copySelection", format: "html" },
  "sql.copySelectionSql": { action: "copySelection", format: "sql" },
  "sql.copySelectionLatex": { action: "copySelection", format: "latex" },
  "sql.copySelectionRaw": { action: "copySelection", format: "raw" },
  "sql.copyFullResultCsv": { action: "copyFullResult", format: "csv" },
  "sql.copyFullResultTsv": { action: "copyFullResult", format: "tsv" },
  "sql.copyFullResultMarkdown": { action: "copyFullResult", format: "markdown" },
  "sql.copyFullResultJson": { action: "copyFullResult", format: "json" },
  "sql.copyFullResultHtml": { action: "copyFullResult", format: "html" },
  "sql.copyFullResultSql": { action: "copyFullResult", format: "sql" },
  "sql.copyFullResultLatex": { action: "copyFullResult", format: "latex" },
  "sql.copyFullResultRaw": { action: "copyFullResult", format: "raw" },
  "sql.exportResultCsv": { action: "export", format: "csv" },
  "sql.exportResultTsv": { action: "export", format: "tsv" },
  "sql.exportResultMarkdown": { action: "export", format: "markdown" },
  "sql.exportResultJson": { action: "export", format: "json" },
  "sql.exportResultHtml": { action: "export", format: "html" },
  "sql.exportResultSql": { action: "export", format: "sql" },
  "sql.exportResultLatex": { action: "export", format: "latex" },
  "sql.exportResultRaw": { action: "export", format: "raw" },
};
const sqlResultContextCommandIDs = new Set([
  "sql.copyCell",
  "sql.copyRow",
  "sql.copyRowWithHeaders",
  "sql.copyColumn",
  "sql.copyColumnWithHeader",
  "sql.openResult",
  "sql.clearSelection",
]);
const sqlResultTargetCommandIDs = new Set([
  "sql.copyCell",
  "sql.copyRow",
  "sql.copyRowWithHeaders",
  "sql.copyColumn",
  "sql.copyColumnWithHeader",
]);
const sqlResultSelectionCommandIDs = new Set([
  "sql.copySelectionCsv",
  "sql.copySelectionTsv",
  "sql.copySelectionMarkdown",
  "sql.copySelectionJson",
  "sql.copySelectionHtml",
  "sql.copySelectionSql",
  "sql.copySelectionLatex",
  "sql.copySelectionRaw",
  "sql.clearSelection",
]);
const sqlResultCommandIDs = new Set([...Object.keys(sqlResultCommandFormats), ...sqlResultContextCommandIDs]);
const findPageSize = 100;
const profileSchemaStorageKey = "opencsv-profile-schema";
const profileSchemasStorageKey = "opencsv-profile-schemas";
const profileSchemasExportKind = "opencsv-profile-schemas";
const profileSchemasLimit = 24;
const defaultProfileSchemaText = `{
  "columns": [
    { "name": "id", "required": true, "type": "integer", "unique": true }
  ]
}`;
const localeStorageKey = "opencsv-locale";
const defaultLocale = "en";
const supportedLocales = [
  { code: "en", label: "English" },
  { code: "ja", label: "Japanese - 日本語" },
  { code: "fr", label: "French - Français" },
  { code: "es", label: "Spanish - Español" },
  { code: "it", label: "Italian - Italiano" },
  { code: "zh-cn", label: "Simplified Chinese - 简体中文" },
  { code: "pt", label: "Portuguese - Português" },
  { code: "ru", label: "Russian - Русский" },
  { code: "de", label: "German - Deutsch" },
  { code: "ko", label: "Korean - 한국어" },
];
const localeCatalog = {
  en: {
    "nav.grid": "Grid",
    "nav.sql": "SQL",
    "nav.profile": "Profile",
    "nav.settings": "Settings",
    "nav.shortcuts": "Keys",
    "nav.command": "Cmd",
    "common.close": "Close",
    "common.cancel": "Cancel",
    "file.open": "Open",
    "edit.copy": "Copy",
    "edit.headers": "Headers",
    "view.filter": "Filter",
    "view.filters": "Filters",
    "view.builder": "Builder",
    "view.sort": "Sort",
    "view.chips.aria": "Active view filters",
    "view.filters.aria": "Recent and saved filters",
    "view.filters.recent": "Recent",
    "view.filters.saved": "Saved",
    "view.chip.sql": "SQL",
    "view.chip.openSql": "Open this filter in SQL",
    "view.chip.edit": "Edit",
    "view.chip.editBuilder": "Edit this filter in Builder",
    "view.chip.save": "Save",
    "view.chip.saveFilter": "Save this filter",
    "view.chip.saveSort": "Save this sort",
    "view.chip.removeSortKey": "Remove Key",
    "view.chip.removeLastSortKey": "Remove the last sort key",
    "view.chip.clear": "Clear",
    "view.chip.clearView": "Clear filter and sort",
    "view.chip.apply": "Apply",
    "view.chip.delete": "Delete",
    "view.chip.savedFilter": "Saved: {{name}}",
    "view.chip.savedSort": "Saved Sort: {{name}}",
    "view.chip.activeSort": "Sort: {{sort}}",
    "view.chip.applySavedFilter": "Apply saved filter {{name}}",
    "view.chip.deleteSavedFilter": "Delete saved filter {{name}}",
    "view.chip.applySavedSort": "Apply saved sort {{name}}",
    "view.chip.deleteSavedSort": "Delete saved sort {{name}}",
    "view.status.filterToSql": "Loaded current filter in SQL Console",
    "view.status.filterToBuilder": "Loaded current filter in Filter Builder",
    "view.status.saveFilterFirst": "Apply a filter before saving it.",
    "view.status.savedFilter": "Saved filter: {{name}}",
    "view.status.savedFilterNotFound": "Saved filter not found.",
    "view.status.deletedSavedFilter": "Deleted saved filter.",
    "view.status.filterMatched": "Filter matched {{count}} rows",
    "view.status.sqlFilterMatched": "SQL filter matched {{count}} rows",
    "view.status.filterFailed": "Filter failed: {{message}}",
    "view.status.sqlFilterFailed": "SQL filter failed: {{message}}",
    "view.status.columnFilterMatched": "Filtered {{column}} matched {{count}} rows",
    "view.status.blankFilterMatched": "Blank {{column}} matched {{count}} rows",
    "view.status.nonblankFilterMatched": "Non-blank {{column}} matched {{count}} rows",
    "view.status.selectedBlank": "Selected cell is blank.",
    "view.status.chooseValueFilter": "Choose at least one value to filter.",
    "view.status.valueFilterFailed": "Value filter failed: {{message}}",
    "view.status.valueFilterMatched": "Value filter {{column}} matched {{count}} rows",
    "view.status.clearFailed": "Clear failed: {{message}}",
    "view.status.cleared": "Filter/sort cleared",
    "view.status.savedFilterMatched": "Saved filter matched {{count}} rows",
    "view.status.recentFilterNotFound": "Recent filter not found.",
    "view.status.recentFilterMatched": "Recent filter matched {{count}} rows",
    "view.status.saveSortFirst": "Apply a sort before saving it.",
    "view.status.savedSort": "Saved sort: {{name}}",
    "view.status.savedSortNotFound": "Saved sort not found.",
    "view.status.savedSortNoColumns": "Saved sort has no matching columns.",
    "view.status.savedSortApplied": "Saved sort {{name}} applied: {{sort}}",
    "view.status.deletedSavedSort": "Deleted saved sort.",
    "view.status.removeSortKeyFirst": "Add a secondary sort key before removing one.",
    "view.status.removedSortKey": "Removed {{key}}; sorted {{sort}}",
    "view.status.sorted": "Sorted {{sort}}",
    "view.status.sortFailed": "Sort failed: {{message}}",
    "view.status.exportedSavedViews": "Exported {{filters}} filters, {{sorts}} sorts.",
    "view.status.importedSavedViews": "Imported {{filters}} filters, {{sorts}} sorts.",
    "view.status.savedViewImportFailed": "Saved view import failed: {{message}}",
    "view.status.filterHistoryOff": "Recent filter history is off.",
    "view.status.filterHistoryKeep": "Recent filter history will keep {{count}} queries.",
    "view.status.clearedRecentFilters": "Cleared recent filter history.",
    "sort.direction.ascending": "ascending",
    "sort.direction.descending": "descending",
    "sort.join.then": ", then ",
    "sort.mode.number": "number",
    "sort.mode.date": "date",
    "sort.mode.length": "length",
    "sort.mode.text": "text",
    "column.stableSort": "Stable sort",
    "column.textTransform": "Text transform",
    "column.applyTransform": "Apply transform to column",
    "column.values": "Column values",
    "column.valuesLoading": "Loading values...",
    "column.valuesFailed": "Values failed: {{message}}",
    "column.noValues": "No values",
    "column.valuesCount": "Values ({{count}})",
    "column.valuesShown": "Values ({{shown}} of {{total}})",
    "column.blankValue": "(blank)",
    "column.applyValues": "Apply selected values",
    "column.hide": "Hide column",
    "column.showAll": "Show all columns",
    "column.showAllShort": "Show All",
    "column.clearView": "Clear filter/sort",
    "column.filterSelected": "Filter by selected value",
    "column.filterValue": "Filter: {{value}}",
    "text.mode.trim": "Trim",
    "text.mode.upper": "Uppercase",
    "text.mode.lower": "Lowercase",
    "text.mode.snake": "Snake Case",
    "text.mode.kebab": "Kebab Case",
    "text.mode.camel": "Camel Case",
    "text.mode.pascal": "Pascal Case",
    "text.mode.title": "Title Case",
    "text.mode.wide": "Full Width",
    "text.mode.narrow": "Half Width",
    "text.mode.nfc": "Normalize NFC",
    "text.mode.nfd": "Normalize NFD",
    "text.mode.nfkc": "Normalize NFKC",
    "text.mode.nfkd": "Normalize NFKD",
    "text.status.running": "{{label}}...",
    "text.status.failed": "{{label}} failed: {{message}}",
    "text.status.changed": "{{label}}: {{count}} cells changed",
    "text.status.noChange": "{{label}}: no cells changed",
    "edit.find": "Find",
    "edit.preview": "Preview",
    "edit.replace": "Replace",
    "edit.all": "All",
    "find.placeholder": "Find/filter...",
    "find.replacePlaceholder": "Replace...",
    "find.matchCase": "Match case",
    "find.wholeCell": "Match whole cell",
    "find.regex": "Use regular expression",
    "find.selection": "Search in selection",
    "find.selectionShort": "Sel",
    "find.previous": "Previous match",
    "find.next": "Next match",
    "find.previewTitle": "Preview replacements",
    "find.replaceTitle": "Replace first match",
    "find.replaceAllTitle": "Replace all matches",
    "find.status.enterQuery": "Type text to find.",
    "find.status.running": "Finding {{query}}...",
    "find.status.runningElapsed": "Finding {{query}}... {{seconds}}s",
    "find.status.canceled": "Find canceled.",
    "find.status.canceledElapsed": "Find canceled after {{seconds}}s.",
    "find.status.noMatches": "No matches for {{query}}",
    "find.status.noMatchesTimed": "No matches for {{query}} in {{seconds}}s",
    "find.status.failed": "Find failed: {{message}}",
    "find.status.page": " (showing matches {{start}}-{{end}})",
    "find.status.found": "Found {{current}} of {{total}} for {{query}} at row {{row}}, column {{column}}{{page}}",
    "find.status.foundTimed": "Found {{current}} of {{total}} for {{query}} at row {{row}}, column {{column}}{{page}} in {{seconds}}s",
    "find.status.selectionLimit": "Selection has {{count}} cells; search-in-selection limit is {{limit}}.",
    "find.status.loadSelection": "Load the selected filtered rows before searching within selection.",
    "replace.status.enterQuery": "Type text to replace.",
    "replace.status.enterPreviewQuery": "Type text to preview replacements.",
    "replace.status.failed": "Replace failed: {{message}}",
    "replace.status.previewFailed": "Replace preview failed: {{message}}",
    "replace.status.noReplacements": "No replacements for {{query}}",
    "replace.status.replacedCells": "Replaced {{count}} cells",
    "replace.status.replacedCurrent": "Replaced current match",
    "replace.status.replacedOne": "Replaced 1 cell",
    "replace.status.previewFirst": "; first row {{row}}, column {{column}}",
    "replace.status.preview": "Preview: {{count}} cells would change{{rowLabel}}: {{before}} -> {{after}}",
    "condition.title": "Filter Builder",
    "condition.close": "Close",
    "condition.quick": "Quick Condition",
    "condition.quickConditions": "Quick filter conditions",
    "condition.additional": "Additional filter conditions",
    "condition.template.choose": "Choose template",
    "condition.template.selectedEquals": "Selected cell equals",
    "condition.template.selectedContains": "Selected cell contains",
    "condition.template.selectedNotContains": "Selected cell does not contain",
    "condition.template.selectedNotEquals": "Selected cell not equal",
    "condition.template.activeBlank": "Active column blank",
    "condition.template.activeNotBlank": "Active column not blank",
    "condition.template.activeTrue": "Active column true",
    "condition.template.activeFalse": "Active column false",
    "condition.useTemplate": "Use Template",
    "condition.legendOne": "Condition 1",
    "condition.legend": "Condition {{number}}",
    "condition.join": "Join",
    "condition.join.and": "AND",
    "condition.join.or": "OR",
    "condition.column": "Column",
    "condition.operator": "Operator",
    "condition.value": "Value",
    "condition.add": "Add Condition",
    "condition.apply": "Apply",
    "condition.clearView": "Clear View",
    "condition.groupAll": "Group All",
    "condition.allAnd": "All AND",
    "condition.allOr": "All OR",
    "condition.clearGroups": "Clear Groups",
    "condition.remove": "Remove",
    "condition.duplicate": "Duplicate",
    "condition.groupPrevious": "Group",
    "condition.ungroupPrevious": "Ungroup",
    "condition.moveUp": "Up",
    "condition.moveDown": "Down",
    "condition.operator.contains": "contains",
    "condition.operator.notContains": "does not contain",
    "condition.operator.equals": "equals",
    "condition.operator.notEquals": "not equals",
    "condition.operator.inList": "in list",
    "condition.operator.notInList": "not in list",
    "condition.operator.startsWith": "starts with",
    "condition.operator.endsWith": "ends with",
    "condition.operator.blank": "is blank",
    "condition.operator.notBlank": "is not blank",
    "condition.operator.booleanTrue": "is true",
    "condition.operator.booleanFalse": "is false",
    "condition.operator.dateEquals": "date equals",
    "condition.operator.dateBefore": "date before",
    "condition.operator.dateAfter": "date after",
    "condition.operator.dateOnOrBefore": "date on/before",
    "condition.operator.dateOnOrAfter": "date on/after",
    "condition.operator.gt": "greater than",
    "condition.operator.gte": "greater or equal",
    "condition.operator.lt": "less than",
    "condition.operator.lte": "less or equal",
    "condition.operator.between": "between",
    "condition.operator.notBetween": "not between",
    "condition.placeholder.range": "min..max",
    "condition.placeholder.list": "value1, value2",
    "condition.placeholder.date": "YYYY-MM-DD",
    "condition.group.openTitle": "Open a group before this condition",
    "condition.group.closeTitle": "Close a group after this condition",
    "condition.group.openAria": "Open group before condition {{number}}",
    "condition.group.closeAria": "Close group after condition {{number}}",
    "condition.joinAria": "Condition {{number}} join",
    "condition.columnAria": "Condition {{number}} column",
    "condition.operatorAria": "Condition {{number}} operator",
    "condition.valueAria": "Condition {{number}} value",
    "condition.removeTitle": "Remove condition {{number}}",
    "condition.duplicateTitle": "Duplicate condition {{number}}",
    "condition.groupPreviousTitle": "Group condition {{number}} with the previous condition",
    "condition.ungroupPreviousTitle": "Remove one group around condition {{number}} and the previous condition",
    "condition.moveUpTitle": "Move condition {{number}} up",
    "condition.moveDownTitle": "Move condition {{number}} down",
    "condition.status.chooseQuick": "Choose a quick filter condition.",
    "condition.status.ready": "Quick condition ready: {{label}}.",
    "condition.status.duplicated": "Duplicated condition {{number}}.",
    "condition.status.groupedPrevious": "Grouped conditions {{from}} and {{to}}.",
    "condition.status.ungroupedPrevious": "Ungrouped conditions {{from}} and {{to}}.",
    "condition.status.reordered": "Moved condition {{from}} to position {{to}}.",
    "condition.status.addConditionFirst": "Add at least two conditions before grouping all.",
    "condition.status.groupedAll": "Grouped all {{count}} conditions.",
    "condition.status.maxGroups": "The selected conditions already use the maximum group depth.",
    "condition.status.allAnd": "Set all condition joins to AND.",
    "condition.status.allOr": "Set all condition joins to OR.",
    "condition.status.clearedGroups": "Cleared {{count}} condition group markers.",
    "condition.status.selectCell": "Select a cell before using a quick condition.",
    "condition.status.unsupportedQuick": "Choose a supported quick condition.",
    "condition.status.tooManyCloses": "Condition {{number}} closes more groups than are open.",
    "condition.status.closeGroups": "Close every condition group before applying.",
    "condition.status.chooseColumn": "Choose a column to filter.",
    "condition.status.enterValue": "Enter a filter value.",
    "condition.status.chooseOperator": "Choose a filter operator.",
    "condition.status.enterRange": "Enter a range as min..max.",
    "condition.status.enterList": "Enter list values separated by commas or new lines.",
    "condition.status.enterDate": "Enter a date as YYYY-MM-DD.",
    "condition.status.matched": "Condition filter matched {{count}} rows",
    "condition.status.failed": "Condition filter failed: {{message}}",
    "condition.filterLabel": "Condition filter",
    "diff.compare": "Compare",
    "diff.run": "Run",
    "diff.ignoreCase": "Ignore case",
    "diff.ignoreWhitespace": "Ignore whitespace",
    "diff.matchHeaders": "Match headers",
    "diff.keyColumnPlaceholder": "Key column e.g. id or 1",
    "diff.ignoredColumnsPlaceholder": "Ignore columns e.g. 3, updated",
    "diff.previousDifference": "Previous difference",
    "diff.nextDifference": "Next difference",
    "diff.status.all": "All rows",
    "diff.status.differences": "Differences",
    "diff.status.changed": "Changed",
    "diff.status.added": "Added",
    "diff.status.deleted": "Deleted",
    "diff.status.unchanged": "Unchanged",
    "diff.searchPlaceholder": "Search diff rows...",
    "diff.sort.source": "Source order",
    "diff.sort.status": "Status",
    "diff.sort.leftRow": "Left row",
    "diff.sort.rightRow": "Right row",
    "diff.sort.leftText": "Left text",
    "diff.sort.rightText": "Right text",
    "diff.sort.asc": "Asc",
    "diff.sort.desc": "Desc",
    "diff.sortAscending": "Sort ascending",
    "diff.sortDescending": "Sort descending",
    "diff.clear": "Clear",
    "diff.rightTextMode": "Compare right side as text",
    "diff.pasteClipboard": "Paste Clipboard",
    "diff.rightTextPlaceholder": "Paste CSV text here",
    "diff.summary.unchanged": "Unchanged",
    "diff.summary.changed": "Changed",
    "diff.summary.added": "Added",
    "diff.summary.deleted": "Deleted",
    "diff.summary.shown": "Shown",
    "diff.table.status": "Status",
    "diff.table.leftRow": "Left Row",
    "diff.table.left": "Left",
    "diff.table.rightRow": "Right Row",
    "diff.table.right": "Right",
    "diff.table.changedCells": "Changed Cells",
    "diff.table.actions": "Actions",
    "diff.empty.none": "No diff rows",
    "diff.empty.filtered": "No diff rows match filters",
    "diff.action.insertRow": "Insert row",
    "diff.action.deleteLeft": "Delete left",
    "diff.action.applyRow": "Apply row",
    "diff.action.deleteRight": "Delete right",
    "diff.action.insertRight": "Insert right",
    "diff.action.applyLeftRow": "Apply left row",
    "diff.action.applyRightRowToLeft": "Apply right row to left file",
    "diff.action.applyLeftRowToRight": "Apply left row to right file",
    "diff.action.applyRightCellToLeft": "Apply right cell c{{col}} to left",
    "diff.action.applyLeftCellToRight": "Apply left cell c{{col}} to right",
    "diff.action.openLeftRow": "Open left row {{row}} in Grid",
    "diff.action.openRightRow": "Open right row {{row}} in Grid",
    "diff.action.openLeftShort": "Open L",
    "diff.action.openRightShort": "Open R",
    "diff.action.editCell": "Save",
    "diff.action.editLeftCell": "Edit left row {{row}} cell c{{col}}",
    "diff.action.editRightCell": "Edit right row {{row}} cell c{{col}}",
    "diff.edit.leftShort": "L",
    "diff.edit.rightShort": "R",
    "diff.side.left": "left",
    "diff.side.right": "right",
    "diff.status.openedSource": "Opened {{side}} row {{row}} from compare.",
    "diff.status.openSourceFailed": "Open compare source failed: {{message}}",
    "diff.status.editingCell": "Editing {{side}} row {{row}} c{{col}}...",
    "diff.status.editedCell": "Edited {{side}} row {{row}} c{{col}}{{backup}}.",
    "diff.status.editFailed": "Diff edit failed: {{message}}",
    "diff.status.editNoChange": "No diff edit applied; value is unchanged.",
    "diff.status.backupSuffix": "; backup {{path}}",
    "diff.status.summary": "Diff: {{changed}} changed, {{added}} added, {{deleted}} deleted",
    "diff.status.filteredSuffix": "; showing {{shown}} rows after filter",
    "diff.status.filterShown": "Diff filter shows {{shown}} of {{total}} rows.",
    "diff.status.filterCleared": "Diff filter cleared; showing {{total}} rows.",
    "diff.status.active": "Diff {{current}} of {{total}}",
    "file.readOnly": "Read Only",
    "file.new": "New",
    "file.export": "Export",
    "file.save": "Save",
    "file.saveAs": "Save As",
    "file.restore": "Restore",
    "workbench.split": "Split",
    "workbench.orientation": "H/V",
    "workbench.focus": "Focus",
    "workbench.move": "Move",
    "workbench.close": "Close",
    "settings.title": "Settings",
    "settings.modified": "Modified",
    "settings.defaults": "Defaults",
    "settings.keyboardShortcuts": "Keyboard Shortcuts",
    "settings.checkForUpdates": "Check for Updates",
    "settings.language": "Language",
    "settings.resetLanguage": "Reset Language",
    "settings.searchPlaceholder": "Search settings",
    "settings.savedViews": "Saved Filters, Recent & Sorts",
    "settings.theme": "Theme",
    "settings.externalChanges": "External Changes",
    "settings.updates": "Updates",
    "settings.linkProtocols": "Link Protocols",
    "settings.editorAppearance": "Editor Appearance",
    "settings.fileFormatPresets": "File Format Presets",
    "link.prompt.title": "Open Custom Link",
    "link.prompt.protocol": "{{protocol}} custom protocol",
    "link.prompt.target": "{{target}}",
    "link.prompt.remember": "Remember this protocol",
    "link.prompt.cancel": "Cancel",
    "link.prompt.open": "Open Link",
    "link.status.confirm": "Confirm opening {{protocol}} links.",
    "link.status.canceled": "Link open canceled.",
    "link.status.opened": "Opened {{href}}",
    "link.status.remembered": "Opened {{href}} and remembered {{protocol}} links.",
    "start.newFile": "New File",
    "start.newHelp": "Create an unsaved CSV",
    "start.openCsv": "Open CSV/TSV",
    "start.openHelp": "Type a path in the toolbar",
    "start.chooseFile": "Choose File",
    "start.chooseHelp": "Open with browser picker",
    "start.settings": "Settings",
    "start.settingsHelp": "Adjust app preferences",
    "start.support": "Support",
    "start.supportHelp": "About, updates, and help",
    "about.title": "About & Support",
    "about.product": "Product",
    "about.license": "License",
    "about.licenseValue": "Free local editor",
    "about.network": "Network",
    "about.networkValue": "Offline by default",
    "about.updates": "Updates",
    "file.recent": "Recent",
    "file.noOpen": "No file open",
    "file.status.creatingNew": "Creating new CSV...",
    "file.status.newReady": "New file ready: {{name}}",
    "file.status.newFailed": "New file failed: {{message}}",
    "status.ready": "Ready",
    "status.languageChanged": "Language updated.",
    "status.languageReset": "Language reset.",
    "sql.run": "Run",
    "sql.copy": "Copy",
    "sql.export": "Export",
    "sql.openResult": "Open Result",
    "sql.clear": "Clear",
    "sql.empty.prompt": "Run SQL to inspect results.",
    "sql.empty.history": "No query history",
    "sql.empty.rows": "No rows",
    "sql.error.title": "SQL error",
    "sql.error.location": "line {{line}}, column {{column}}",
    "sql.error.token": "token {{token}}",
    "sql.error.phase": "phase {{phase}}",
    "sql.status.enterQuery": "Enter a SQL query.",
    "sql.status.running": "Running SQL...",
    "sql.status.runningElapsed": "Running SQL... {{seconds}}s",
    "sql.status.runningRows": "Loading {{phase}}... {{rows}} source rows in {{seconds}}s",
    "sql.status.canceled": "SQL canceled.",
    "sql.status.canceledElapsed": "SQL canceled after {{seconds}}s.",
    "sql.status.returned": "SQL returned {{rows}} of {{total}} rows",
    "sql.status.returnedTimed": "SQL returned {{rows}} of {{total}} rows in {{seconds}}s",
    "sql.status.returnedTimedRows": "SQL loaded {{loadedRows}} source rows, returned {{rows}} of {{total}} rows in {{seconds}}s",
    "sql.progress.load": "source rows",
    "sql.status.historySaveFailed": "SQL query history could not be saved.",
    "sql.status.historyClearFailed": "SQL query history could not be cleared.",
    "sql.status.historyLoaded": "Loaded SQL query from history.",
    "sql.status.failed": "SQL failed: {{message}}",
    "sql.status.failedAt": "SQL failed at line {{line}}, column {{column}}: {{message}}",
    "sql.status.runBeforeCopy": "Run a SQL query before copying results.",
    "sql.status.runBeforeExport": "Run a SQL query before exporting results.",
    "sql.status.runBeforeOpen": "Run a SQL query before opening results.",
    "sql.status.openingResult": "Opening SQL result...",
    "sql.status.openedResult": "Opened SQL result as {{name}}: {{rows}} rows, {{cols}} columns",
    "sql.status.openResultFailed": "Open SQL result failed: {{message}}",
    "sql.status.selectCell": "Select a SQL result cell first.",
    "sql.status.selectRow": "Select a SQL result row first.",
    "sql.status.selectColumn": "Select a SQL result column first.",
    "sql.status.noCopy": "No {{label}} to copy.",
    "sql.status.copied": "Copied {{label}}: {{rows}} rows, {{cols}} columns as {{format}}",
    "sql.status.copyFailed": "{{label}} copy failed: {{message}}",
    "sql.status.exported": "Exported SQL result as {{format}}: {{rows}} rows, {{cols}} columns to {{filename}}",
    "sql.status.selected": "Selected SQL result: {{rows}} rows, {{cols}} columns",
    "sql.status.selectionCleared": "Cleared SQL result selection.",
    "sql.label.selection": "SQL selection",
    "sql.label.result": "SQL result",
    "sql.label.cell": "SQL cell",
    "sql.label.row": "SQL row",
    "sql.label.rowWithHeaders": "SQL row with headers",
    "sql.label.column": "SQL column",
    "sql.label.columnWithHeader": "SQL column with header",
    "sql.context.copyCell": "Copy Cell",
    "sql.context.copyRow": "Copy Row",
    "sql.context.copyRowWithHeaders": "Copy Row with Headers",
    "sql.context.copyColumn": "Copy Column",
    "sql.context.copyColumnWithHeader": "Copy Column with Header",
    "sql.context.copySelectionAs": "Copy Selection as {{format}}",
    "sql.context.copyResultAs": "Copy Result as {{format}}",
    "sql.context.exportResultAs": "Export Result as {{format}}",
    "sql.context.openResult": "Open Result as Document",
    "sql.context.clearSelection": "Clear Selection",
    "profile.run": "Run Profile",
    "profile.schema": "Schema",
    "profile.saveSchema": "Save Schema",
    "profile.resetSchema": "Reset Schema",
    "profile.apply": "Apply",
    "profile.delete": "Delete",
    "profile.export": "Export",
    "profile.import": "Import",
    "profile.cancel": "Cancel",
    "profile.schemaNamePlaceholder": "Schema name",
    "profile.savedSchemas": "Saved schemas",
    "profile.noSchema": "No schema",
    "profile.schemaRules": "{{count}} schema rules",
    "profile.invalidSchema": "Invalid schema JSON",
    "profile.rule.columnName": "Column name",
    "profile.rule.columnNumber": "Column #",
    "profile.rule.type": "Type",
    "profile.rule.required": "Required",
    "profile.rule.unique": "Unique",
    "profile.rule.enum": "Enum",
    "profile.rule.pattern": "Pattern",
    "profile.rule.default": "Default",
    "profile.rule.min": "Min",
    "profile.rule.max": "Max",
    "profile.rule.minLength": "Min Length",
    "profile.rule.maxLength": "Max Length",
    "profile.rule.defaultSummary": "default {{value}}",
    "profile.rule.save": "Save Rule",
    "profile.rule.update": "Update Rule",
    "profile.rule.clear": "Clear Rule",
    "profile.rule.none": "No schema rules.",
    "profile.rule.edit": "Edit",
    "profile.rule.remove": "Remove",
    "profile.rule.rule": "Rule {{index}}",
    "profile.rule.locationOnly": "location only",
    "profile.rule.requiredSummary": "required",
    "profile.rule.uniqueSummary": "unique",
    "profile.rule.typeSummary": "type {{type}}",
    "profile.rule.enumSummary": "enum {{values}}",
    "profile.rule.patternSummary": "pattern {{pattern}}",
    "profile.rule.minSummary": "min {{value}}",
    "profile.rule.maxSummary": "max {{value}}",
    "profile.rule.exclusiveMinSummary": "min > {{value}}",
    "profile.rule.exclusiveMaxSummary": "max < {{value}}",
    "profile.rule.minLengthSummary": "min length {{value}}",
    "profile.rule.maxLengthSummary": "max length {{value}}",
    "profile.type.any": "Any",
    "profile.type.string": "String",
    "profile.type.number": "Number",
    "profile.type.integer": "Integer",
    "profile.type.decimal": "Decimal",
    "profile.type.boolean": "Boolean",
    "profile.type.date": "Date",
    "profile.type.datetime": "Date-Time",
    "profile.type.time": "Time",
    "profile.type.year": "Year",
    "profile.type.yearmonth": "Year-Month",
    "profile.type.unknown": "Unknown",
    "profile.severity.error": "Error",
    "profile.severity.warning": "Warning",
    "profile.quickFix.header.fill": "Fill Header",
    "profile.quickFix.header.unique": "Make Unique",
    "profile.quickFix.row.pad": "Pad Row",
    "profile.quickFix.schema.column.add": "Add Column",
    "profile.quickFix.schema.default.use": "Use Default",
    "profile.quickFix.schema.enum.useValue": "Use Allowed Value",
    "profile.quickFix.schema.range.clampMin": "Use Minimum",
    "profile.quickFix.schema.range.clampMax": "Use Maximum",
    "profile.summary.status": "Status",
    "profile.summary.rows": "Rows",
    "profile.summary.columns": "Columns",
    "profile.summary.errors": "Errors",
    "profile.summary.warnings": "Warnings",
    "profile.summary.rules": "Rules",
    "profile.summary.progress": "Scanned",
    "profile.summary.progressValue": "{{passes}} passes / {{rows}} row scans",
    "profile.loading": "loading",
    "profile.columns.loading": "Profiling columns...",
    "profile.issues.loading": "Checking data quality...",
    "profile.status.running": "Profiling data...",
    "profile.status.runningElapsed": "Profiling data... {{seconds}}s",
    "profile.status.runningRows": "Profiling {{phase}}... {{rows}} rows scanned in {{seconds}}s",
    "profile.status.canceled": "Profile canceled.",
    "profile.status.canceledElapsed": "Profile canceled after {{seconds}}s.",
    "profile.status.complete": "Profile: {{rows}} rows, {{columns}} columns, {{issues}} issues",
    "profile.status.completeTimed": "Profile: {{rows}} rows, {{columns}} columns, {{issues}} issues in {{seconds}}s",
    "profile.status.completeTimedRows": "Profile: {{rows}} rows, {{columns}} columns, {{issues}} issues in {{seconds}}s (profile {{profileRows}}, validation {{validationRows}} rows)",
    "profile.status.failed": "Profile failed: {{message}}",
    "profile.progress.profile": "profile pass",
    "profile.progress.validation": "validation pass",
    "profile.status.initialColumns": "Run Profile to inspect columns.",
    "profile.status.initialIssues": "No validation issues loaded.",
    "profile.empty.noValidation": "No validation result.",
    "profile.empty.noColumns": "No columns found.",
    "profile.empty.noIssues": "No validation issues.",
    "profile.columns.column": "Column",
    "profile.columns.type": "Type",
    "profile.columns.rows": "Rows",
    "profile.columns.nulls": "Nulls",
    "profile.columns.distinct": "Distinct",
    "profile.columns.min": "Min",
    "profile.columns.max": "Max",
    "profile.columns.average": "Average",
    "profile.columns.topValues": "Top Values",
    "profile.value.empty": "(empty)",
    "profile.issues.severity": "Severity",
    "profile.issues.type": "Type",
    "profile.issues.row": "Row",
    "profile.issues.column": "Column",
    "profile.issues.message": "Message",
    "profile.issues.fix": "Fix",
    "profile.issues.defaultFix": "Fix",
    "profile.status.selectedIssue": "Selected {{type}} at row {{row}}, column {{column}}.",
    "profile.schema.statusInvalid": "Schema invalid: {{message}}",
    "profile.schema.saved": "Validation schema saved.",
    "profile.schema.savedAs": "Validation schema saved as {{name}}.",
    "profile.schema.reset": "Validation schema reset.",
    "profile.schema.notFound": "Saved schema not found.",
    "profile.schema.applied": "Applied schema {{name}}.",
    "profile.schema.deleted": "Deleted saved schema.",
    "profile.schema.exported": "Exported {{count}} profile schemas.",
    "profile.schema.imported": "Imported {{count}} profile schemas.",
    "profile.schema.importFailed": "Profile schema import failed: {{message}}",
    "profile.rule.chooseColumn": "Choose a column name or number.",
    "profile.rule.invalid": "Rule invalid: {{message}}",
    "profile.rule.added": "Added schema rule.",
    "profile.rule.updated": "Updated schema rule.",
    "profile.rule.cleared": "Cleared schema rule editor.",
    "profile.rule.notFound": "Schema rule not found.",
    "profile.rule.editing": "Editing schema rule {{index}}.",
    "profile.rule.removed": "Removed schema rule.",
    "profile.quickFix.readOnly": "Read-only mode blocks quick fixes.",
    "profile.quickFix.noFile": "Open a file before applying fixes.",
    "profile.quickFix.unavailable": "No safe quick fix is available for this issue.",
    "profile.quickFix.applied": "Applied quick fix to {{count}} cells.",
    "profile.quickFix.failed": "Quick fix failed: {{message}}",
    "command.selection.insertAverageRow.title": "Add Average Row",
    "command.selection.insertAverageRow.description": "Add a removable average row for the selected numeric cells.",
    "command.selection.insertCountRow.title": "Add Count Row",
    "command.selection.insertCountRow.description": "Add a removable count row for the selected cells.",
  },
  ja: {
    "nav.grid": "グリッド",
    "nav.profile": "プロファイル",
    "nav.settings": "設定",
    "file.open": "開く",
    "edit.copy": "コピー",
    "edit.headers": "ヘッダー",
    "view.filter": "フィルター",
    "view.filters": "フィルター",
    "view.builder": "ビルダー",
    "view.sort": "並べ替え",
    "edit.find": "検索",
    "edit.preview": "プレビュー",
    "edit.replace": "置換",
    "edit.all": "すべて",
    "diff.compare": "比較",
    "file.readOnly": "読み取り専用",
    "file.export": "エクスポート",
    "file.save": "保存",
    "file.saveAs": "別名保存",
    "file.restore": "復元",
    "workbench.split": "分割",
    "workbench.focus": "フォーカス",
    "workbench.move": "移動",
    "workbench.close": "閉じる",
    "settings.title": "設定",
    "settings.modified": "変更済み",
    "settings.defaults": "既定",
    "settings.language": "言語",
    "settings.resetLanguage": "言語をリセット",
    "settings.searchPlaceholder": "設定を検索",
    "start.openCsv": "CSV/TSV を開く",
    "start.openHelp": "ツールバーにパスを入力",
    "start.chooseFile": "ファイルを選択",
    "start.chooseHelp": "ブラウザーの選択画面で開く",
    "file.recent": "最近",
    "file.noOpen": "ファイル未選択",
    "status.ready": "準備完了",
    "status.languageChanged": "言語を更新しました。",
    "status.languageReset": "言語をリセットしました。",
  },
  fr: {
    "nav.grid": "Grille",
    "nav.profile": "Profil",
    "nav.settings": "Paramètres",
    "file.open": "Ouvrir",
    "edit.copy": "Copier",
    "view.filter": "Filtrer",
    "view.filters": "Filtres",
    "view.builder": "Assistant",
    "view.sort": "Trier",
    "edit.find": "Rechercher",
    "edit.preview": "Aperçu",
    "edit.replace": "Remplacer",
    "diff.compare": "Comparer",
    "file.readOnly": "Lecture seule",
    "file.export": "Exporter",
    "file.save": "Enregistrer",
    "file.saveAs": "Enregistrer sous",
    "file.restore": "Restaurer",
    "settings.title": "Paramètres",
    "settings.modified": "Modifié",
    "settings.defaults": "Par défaut",
    "settings.language": "Langue",
    "settings.resetLanguage": "Réinitialiser la langue",
    "settings.searchPlaceholder": "Rechercher les paramètres",
    "start.openCsv": "Ouvrir CSV/TSV",
    "start.chooseFile": "Choisir un fichier",
    "file.recent": "Récents",
    "status.ready": "Prêt",
  },
  es: {
    "nav.grid": "Cuadrícula",
    "nav.profile": "Perfil",
    "nav.settings": "Ajustes",
    "file.open": "Abrir",
    "edit.copy": "Copiar",
    "view.filter": "Filtrar",
    "view.filters": "Filtros",
    "view.sort": "Ordenar",
    "edit.find": "Buscar",
    "edit.preview": "Vista previa",
    "edit.replace": "Reemplazar",
    "diff.compare": "Comparar",
    "file.readOnly": "Solo lectura",
    "file.export": "Exportar",
    "file.save": "Guardar",
    "file.saveAs": "Guardar como",
    "settings.title": "Ajustes",
    "settings.modified": "Modificado",
    "settings.defaults": "Predeterminado",
    "settings.language": "Idioma",
    "settings.resetLanguage": "Restablecer idioma",
    "settings.searchPlaceholder": "Buscar ajustes",
    "start.openCsv": "Abrir CSV/TSV",
    "start.chooseFile": "Elegir archivo",
    "file.recent": "Recientes",
    "status.ready": "Listo",
  },
  it: {
    "nav.grid": "Griglia",
    "nav.profile": "Profilo",
    "nav.settings": "Impostazioni",
    "file.open": "Apri",
    "edit.copy": "Copia",
    "view.filter": "Filtro",
    "view.sort": "Ordina",
    "edit.find": "Trova",
    "edit.replace": "Sostituisci",
    "diff.compare": "Confronta",
    "file.export": "Esporta",
    "file.save": "Salva",
    "settings.title": "Impostazioni",
    "settings.language": "Lingua",
    "settings.resetLanguage": "Ripristina lingua",
    "settings.searchPlaceholder": "Cerca impostazioni",
    "start.openCsv": "Apri CSV/TSV",
    "start.chooseFile": "Scegli file",
    "file.recent": "Recenti",
    "status.ready": "Pronto",
  },
  "zh-cn": {
    "nav.grid": "网格",
    "nav.profile": "概况",
    "nav.settings": "设置",
    "nav.shortcuts": "快捷键",
    "common.close": "关闭",
    "common.cancel": "取消",
    "file.open": "打开",
    "edit.copy": "复制",
    "edit.headers": "表头",
    "view.filter": "筛选",
    "view.filters": "筛选",
    "view.builder": "构建器",
    "view.sort": "排序",
    "view.chips.aria": "当前视图筛选",
    "view.filters.aria": "最近和保存的筛选",
    "view.filters.recent": "最近",
    "view.filters.saved": "保存",
    "view.chip.sql": "SQL",
    "view.chip.openSql": "在 SQL 控制台打开此筛选",
    "view.chip.edit": "编辑",
    "view.chip.editBuilder": "在筛选构建器中编辑此筛选",
    "view.chip.save": "保存",
    "view.chip.saveFilter": "保存此筛选",
    "view.chip.saveSort": "保存此排序",
    "view.chip.removeSortKey": "移除键",
    "view.chip.removeLastSortKey": "移除最后一个排序键",
    "view.chip.clear": "清除",
    "view.chip.clearView": "清除筛选和排序",
    "view.chip.apply": "应用",
    "view.chip.delete": "删除",
    "view.chip.savedFilter": "保存的筛选：{{name}}",
    "view.chip.savedSort": "保存的排序：{{name}}",
    "view.chip.activeSort": "排序：{{sort}}",
    "view.chip.applySavedFilter": "应用保存的筛选 {{name}}",
    "view.chip.deleteSavedFilter": "删除保存的筛选 {{name}}",
    "view.chip.applySavedSort": "应用保存的排序 {{name}}",
    "view.chip.deleteSavedSort": "删除保存的排序 {{name}}",
    "view.status.filterToSql": "已将当前筛选载入 SQL 控制台",
    "view.status.filterToBuilder": "已将当前筛选载入筛选构建器",
    "view.status.saveFilterFirst": "请先应用筛选再保存。",
    "view.status.savedFilter": "已保存筛选：{{name}}",
    "view.status.savedFilterNotFound": "未找到保存的筛选。",
    "view.status.deletedSavedFilter": "已删除保存的筛选。",
    "view.status.filterMatched": "筛选匹配 {{count}} 行",
    "view.status.sqlFilterMatched": "SQL 筛选匹配 {{count}} 行",
    "view.status.filterFailed": "筛选失败：{{message}}",
    "view.status.sqlFilterFailed": "SQL 筛选失败：{{message}}",
    "view.status.columnFilterMatched": "已筛选 {{column}}，匹配 {{count}} 行",
    "view.status.blankFilterMatched": "{{column}} 空值匹配 {{count}} 行",
    "view.status.nonblankFilterMatched": "{{column}} 非空值匹配 {{count}} 行",
    "view.status.selectedBlank": "选中单元格为空。",
    "view.status.chooseValueFilter": "请至少选择一个筛选值。",
    "view.status.valueFilterFailed": "值筛选失败：{{message}}",
    "view.status.valueFilterMatched": "{{column}} 值筛选匹配 {{count}} 行",
    "view.status.clearFailed": "清除失败：{{message}}",
    "view.status.cleared": "已清除筛选/排序",
    "view.status.savedFilterMatched": "保存的筛选匹配 {{count}} 行",
    "view.status.recentFilterNotFound": "未找到最近筛选。",
    "view.status.recentFilterMatched": "最近筛选匹配 {{count}} 行",
    "view.status.saveSortFirst": "请先应用排序再保存。",
    "view.status.savedSort": "已保存排序：{{name}}",
    "view.status.savedSortNotFound": "未找到保存的排序。",
    "view.status.savedSortNoColumns": "保存的排序没有匹配列。",
    "view.status.savedSortApplied": "已应用保存的排序 {{name}}：{{sort}}",
    "view.status.deletedSavedSort": "已删除保存的排序。",
    "view.status.removeSortKeyFirst": "请先添加次级排序键再移除。",
    "view.status.removedSortKey": "已移除 {{key}}；已排序：{{sort}}",
    "view.status.sorted": "已排序：{{sort}}",
    "view.status.sortFailed": "排序失败：{{message}}",
    "view.status.exportedSavedViews": "已导出 {{filters}} 个筛选、{{sorts}} 个排序。",
    "view.status.importedSavedViews": "已导入 {{filters}} 个筛选、{{sorts}} 个排序。",
    "view.status.savedViewImportFailed": "保存视图导入失败：{{message}}",
    "view.status.filterHistoryOff": "最近筛选历史已关闭。",
    "view.status.filterHistoryKeep": "最近筛选历史将保留 {{count}} 个查询。",
    "view.status.clearedRecentFilters": "已清除最近筛选历史。",
    "sort.direction.ascending": "升序",
    "sort.direction.descending": "降序",
    "sort.join.then": "，然后 ",
    "sort.mode.number": "数字",
    "sort.mode.date": "日期",
    "sort.mode.length": "长度",
    "sort.mode.text": "文本",
    "column.stableSort": "稳定排序",
    "column.textTransform": "文本转换",
    "column.applyTransform": "将转换应用到列",
    "column.values": "列值",
    "column.valuesLoading": "正在加载值...",
    "column.valuesFailed": "值加载失败：{{message}}",
    "column.noValues": "没有值",
    "column.valuesCount": "值（{{count}}）",
    "column.valuesShown": "值（{{shown}} / {{total}}）",
    "column.blankValue": "（空白）",
    "column.applyValues": "应用选中的值",
    "column.hide": "隐藏列",
    "column.showAll": "显示所有列",
    "column.showAllShort": "显示全部",
    "column.clearView": "清除筛选/排序",
    "column.filterSelected": "按选中值筛选",
    "column.filterValue": "筛选：{{value}}",
    "text.mode.trim": "修剪",
    "text.mode.upper": "转大写",
    "text.mode.lower": "转小写",
    "text.mode.snake": "蛇形命名",
    "text.mode.kebab": "短横线命名",
    "text.mode.camel": "驼峰命名",
    "text.mode.pascal": "帕斯卡命名",
    "text.mode.title": "标题大小写",
    "text.mode.wide": "全角",
    "text.mode.narrow": "半角",
    "text.mode.nfc": "标准化 NFC",
    "text.mode.nfd": "标准化 NFD",
    "text.mode.nfkc": "标准化 NFKC",
    "text.mode.nfkd": "标准化 NFKD",
    "text.status.running": "正在{{label}}...",
    "text.status.failed": "{{label}}失败：{{message}}",
    "text.status.changed": "{{label}}：已更改 {{count}} 个单元格",
    "text.status.noChange": "{{label}}：没有单元格更改",
    "edit.find": "查找",
    "edit.preview": "预览",
    "edit.replace": "替换",
    "edit.all": "全部",
    "find.placeholder": "查找/筛选...",
    "find.replacePlaceholder": "替换...",
    "find.matchCase": "区分大小写",
    "find.wholeCell": "匹配整个单元格",
    "find.regex": "使用正则表达式",
    "find.selection": "在选区中搜索",
    "find.selectionShort": "选",
    "find.previous": "上一个匹配",
    "find.next": "下一个匹配",
    "find.previewTitle": "预览替换",
    "find.replaceTitle": "替换第一个匹配",
    "find.replaceAllTitle": "替换全部匹配",
    "find.status.enterQuery": "请输入要查找的文本。",
    "find.status.running": "正在查找 {{query}}...",
    "find.status.runningElapsed": "正在查找 {{query}}... {{seconds}} 秒",
    "find.status.canceled": "查找已取消。",
    "find.status.canceledElapsed": "查找已在 {{seconds}} 秒后取消。",
    "find.status.noMatches": "没有找到 {{query}} 的匹配项",
    "find.status.noMatchesTimed": "用时 {{seconds}} 秒，没有找到 {{query}} 的匹配项",
    "find.status.failed": "查找失败：{{message}}",
    "find.status.page": "（显示第 {{start}}-{{end}} 个匹配）",
    "find.status.found": "找到 {{current}} / {{total}} 个 {{query}}，第 {{row}} 行，第 {{column}} 列{{page}}",
    "find.status.foundTimed": "找到 {{current}} / {{total}} 个 {{query}}，第 {{row}} 行，第 {{column}} 列{{page}}，用时 {{seconds}} 秒",
    "find.status.selectionLimit": "选区有 {{count}} 个单元格；选区内搜索上限是 {{limit}}。",
    "find.status.loadSelection": "请先加载所选筛选行，再在选区中搜索。",
    "replace.status.enterQuery": "请输入要替换的文本。",
    "replace.status.enterPreviewQuery": "请输入要预览替换的文本。",
    "replace.status.failed": "替换失败：{{message}}",
    "replace.status.previewFailed": "替换预览失败：{{message}}",
    "replace.status.noReplacements": "没有 {{query}} 的可替换项",
    "replace.status.replacedCells": "已替换 {{count}} 个单元格",
    "replace.status.replacedCurrent": "已替换当前匹配",
    "replace.status.replacedOne": "已替换 1 个单元格",
    "replace.status.previewFirst": "；首个在第 {{row}} 行，第 {{column}} 列",
    "replace.status.preview": "预览：{{count}} 个单元格将变更{{rowLabel}}：{{before}} -> {{after}}",
    "condition.title": "筛选构建器",
    "condition.close": "关闭",
    "condition.quick": "快速条件",
    "condition.quickConditions": "快速筛选条件",
    "condition.additional": "附加筛选条件",
    "condition.template.choose": "选择模板",
    "condition.template.selectedEquals": "选中单元格等于",
    "condition.template.selectedContains": "选中单元格包含",
    "condition.template.selectedNotContains": "选中单元格不包含",
    "condition.template.selectedNotEquals": "选中单元格不等于",
    "condition.template.activeBlank": "当前列为空",
    "condition.template.activeNotBlank": "当前列非空",
    "condition.template.activeTrue": "当前列为真",
    "condition.template.activeFalse": "当前列为假",
    "condition.useTemplate": "使用模板",
    "condition.legendOne": "条件 1",
    "condition.legend": "条件 {{number}}",
    "condition.join": "连接",
    "condition.join.and": "且",
    "condition.join.or": "或",
    "condition.column": "列",
    "condition.operator": "操作符",
    "condition.value": "值",
    "condition.add": "添加条件",
    "condition.apply": "应用",
    "condition.clearView": "清除视图",
    "condition.groupAll": "全部分组",
    "condition.allAnd": "全部且",
    "condition.allOr": "全部或",
    "condition.clearGroups": "清除分组",
    "condition.remove": "移除",
    "condition.duplicate": "复制",
    "condition.groupPrevious": "分组",
    "condition.ungroupPrevious": "解组",
    "condition.moveUp": "上移",
    "condition.moveDown": "下移",
    "condition.operator.contains": "包含",
    "condition.operator.notContains": "不包含",
    "condition.operator.equals": "等于",
    "condition.operator.notEquals": "不等于",
    "condition.operator.inList": "在列表中",
    "condition.operator.notInList": "不在列表中",
    "condition.operator.startsWith": "开头为",
    "condition.operator.endsWith": "结尾为",
    "condition.operator.blank": "为空",
    "condition.operator.notBlank": "非空",
    "condition.operator.booleanTrue": "为真",
    "condition.operator.booleanFalse": "为假",
    "condition.operator.dateEquals": "日期等于",
    "condition.operator.dateBefore": "日期早于",
    "condition.operator.dateAfter": "日期晚于",
    "condition.operator.dateOnOrBefore": "日期不晚于",
    "condition.operator.dateOnOrAfter": "日期不早于",
    "condition.operator.gt": "大于",
    "condition.operator.gte": "大于或等于",
    "condition.operator.lt": "小于",
    "condition.operator.lte": "小于或等于",
    "condition.operator.between": "介于",
    "condition.operator.notBetween": "不介于",
    "condition.placeholder.range": "最小值..最大值",
    "condition.placeholder.list": "值1, 值2",
    "condition.placeholder.date": "YYYY-MM-DD",
    "condition.group.openTitle": "在此条件前打开分组",
    "condition.group.closeTitle": "在此条件后关闭分组",
    "condition.group.openAria": "在条件 {{number}} 前打开分组",
    "condition.group.closeAria": "在条件 {{number}} 后关闭分组",
    "condition.joinAria": "条件 {{number}} 连接方式",
    "condition.columnAria": "条件 {{number}} 列",
    "condition.operatorAria": "条件 {{number}} 操作符",
    "condition.valueAria": "条件 {{number}} 值",
    "condition.removeTitle": "移除条件 {{number}}",
    "condition.duplicateTitle": "复制条件 {{number}}",
    "condition.groupPreviousTitle": "将条件 {{number}} 与上一条件分组",
    "condition.ungroupPreviousTitle": "移除条件 {{number}} 与上一条件外的一层分组",
    "condition.moveUpTitle": "上移条件 {{number}}",
    "condition.moveDownTitle": "下移条件 {{number}}",
    "condition.status.chooseQuick": "请选择快速筛选条件。",
    "condition.status.ready": "快速条件已就绪：{{label}}。",
    "condition.status.duplicated": "已复制条件 {{number}}。",
    "condition.status.groupedPrevious": "已将条件 {{from}} 与 {{to}} 分组。",
    "condition.status.ungroupedPrevious": "已取消条件 {{from}} 与 {{to}} 的分组。",
    "condition.status.reordered": "已将条件 {{from}} 移到位置 {{to}}。",
    "condition.status.addConditionFirst": "请先添加至少两个条件再全部分组。",
    "condition.status.groupedAll": "已将全部 {{count}} 个条件分组。",
    "condition.status.maxGroups": "所选条件已达到最大分组深度。",
    "condition.status.allAnd": "已将所有条件连接设为且。",
    "condition.status.allOr": "已将所有条件连接设为或。",
    "condition.status.clearedGroups": "已清除 {{count}} 个条件分组标记。",
    "condition.status.selectCell": "使用快速条件前请先选择一个单元格。",
    "condition.status.unsupportedQuick": "请选择支持的快速条件。",
    "condition.status.tooManyCloses": "条件 {{number}} 关闭的分组多于已打开分组。",
    "condition.status.closeGroups": "应用前请关闭所有条件分组。",
    "condition.status.chooseColumn": "请选择要筛选的列。",
    "condition.status.enterValue": "请输入筛选值。",
    "condition.status.chooseOperator": "请选择筛选操作符。",
    "condition.status.enterRange": "请输入范围，格式为 最小值..最大值。",
    "condition.status.enterList": "请输入用逗号或换行分隔的列表值。",
    "condition.status.enterDate": "请输入 YYYY-MM-DD 格式的日期。",
    "condition.status.matched": "筛选条件匹配 {{count}} 行",
    "condition.status.failed": "筛选条件失败：{{message}}",
    "condition.filterLabel": "筛选条件",
    "diff.compare": "比较",
    "diff.run": "运行",
    "diff.ignoreCase": "忽略大小写",
    "diff.ignoreWhitespace": "忽略空白",
    "diff.matchHeaders": "按表头匹配",
    "diff.keyColumnPlaceholder": "键列，例如 id 或 1",
    "diff.ignoredColumnsPlaceholder": "忽略列，例如 3, updated",
    "diff.previousDifference": "上一个差异",
    "diff.nextDifference": "下一个差异",
    "diff.status.all": "全部行",
    "diff.status.differences": "差异",
    "diff.status.changed": "已更改",
    "diff.status.added": "已添加",
    "diff.status.deleted": "已删除",
    "diff.status.unchanged": "未更改",
    "diff.searchPlaceholder": "搜索差异行...",
    "diff.sort.source": "源顺序",
    "diff.sort.status": "状态",
    "diff.sort.leftRow": "左行",
    "diff.sort.rightRow": "右行",
    "diff.sort.leftText": "左侧文本",
    "diff.sort.rightText": "右侧文本",
    "diff.sort.asc": "升序",
    "diff.sort.desc": "降序",
    "diff.sortAscending": "升序排序",
    "diff.sortDescending": "降序排序",
    "diff.clear": "清除",
    "diff.rightTextMode": "将右侧作为文本比较",
    "diff.pasteClipboard": "粘贴剪贴板",
    "diff.rightTextPlaceholder": "在此粘贴 CSV 文本",
    "diff.summary.unchanged": "未更改",
    "diff.summary.changed": "已更改",
    "diff.summary.added": "已添加",
    "diff.summary.deleted": "已删除",
    "diff.summary.shown": "显示",
    "diff.table.status": "状态",
    "diff.table.leftRow": "左行",
    "diff.table.left": "左侧",
    "diff.table.rightRow": "右行",
    "diff.table.right": "右侧",
    "diff.table.changedCells": "更改单元格",
    "diff.table.actions": "操作",
    "diff.empty.none": "没有差异行",
    "diff.empty.filtered": "没有匹配筛选的差异行",
    "diff.action.insertRow": "插入行",
    "diff.action.deleteLeft": "删除左侧",
    "diff.action.applyRow": "应用行",
    "diff.action.deleteRight": "删除右侧",
    "diff.action.insertRight": "插入右侧",
    "diff.action.applyLeftRow": "应用左侧行",
    "diff.action.applyRightRowToLeft": "将右侧行应用到左文件",
    "diff.action.applyLeftRowToRight": "将左侧行应用到右文件",
    "diff.action.applyRightCellToLeft": "将右侧单元格 c{{col}} 应用到左侧",
    "diff.action.applyLeftCellToRight": "将左侧单元格 c{{col}} 应用到右侧",
    "diff.action.openLeftRow": "在表格中打开左侧第 {{row}} 行",
    "diff.action.openRightRow": "在表格中打开右侧第 {{row}} 行",
    "diff.action.openLeftShort": "打开左侧",
    "diff.action.openRightShort": "打开右侧",
    "diff.action.editCell": "保存",
    "diff.action.editLeftCell": "编辑左侧第 {{row}} 行 c{{col}}",
    "diff.action.editRightCell": "编辑右侧第 {{row}} 行 c{{col}}",
    "diff.edit.leftShort": "左",
    "diff.edit.rightShort": "右",
    "diff.side.left": "左侧",
    "diff.side.right": "右侧",
    "diff.status.openedSource": "已从比较打开{{side}}第 {{row}} 行。",
    "diff.status.openSourceFailed": "打开比较源失败：{{message}}",
    "diff.status.editingCell": "正在编辑{{side}}第 {{row}} 行 c{{col}}...",
    "diff.status.editedCell": "已编辑{{side}}第 {{row}} 行 c{{col}}{{backup}}。",
    "diff.status.editFailed": "差异编辑失败：{{message}}",
    "diff.status.editNoChange": "未应用差异编辑；值未改变。",
    "diff.status.backupSuffix": "；备份 {{path}}",
    "diff.status.summary": "差异：{{changed}} 个更改，{{added}} 个添加，{{deleted}} 个删除",
    "diff.status.filteredSuffix": "；筛选后显示 {{shown}} 行",
    "diff.status.filterShown": "差异筛选显示 {{shown}} / {{total}} 行。",
    "diff.status.filterCleared": "差异筛选已清除，显示 {{total}} 行。",
    "diff.status.active": "差异 {{current}} / {{total}}",
    "file.readOnly": "只读",
    "file.new": "新建",
    "file.export": "导出",
    "file.save": "保存",
    "file.saveAs": "另存为",
    "file.restore": "恢复",
    "workbench.split": "拆分",
    "workbench.orientation": "横/纵",
    "workbench.focus": "聚焦",
    "workbench.move": "移动",
    "workbench.close": "关闭",
    "settings.title": "设置",
    "settings.modified": "已修改",
    "settings.defaults": "默认",
    "settings.keyboardShortcuts": "键盘快捷键",
    "settings.checkForUpdates": "检查更新",
    "settings.language": "语言",
    "settings.resetLanguage": "重置语言",
    "settings.searchPlaceholder": "搜索设置",
    "settings.savedViews": "保存的筛选、最近记录和排序",
    "settings.theme": "主题",
    "settings.externalChanges": "外部更改",
    "settings.updates": "更新",
    "settings.linkProtocols": "链接协议",
    "settings.editorAppearance": "编辑器外观",
    "settings.fileFormatPresets": "文件格式预设",
    "link.prompt.title": "打开自定义链接",
    "link.prompt.protocol": "{{protocol}} 自定义协议",
    "link.prompt.target": "{{target}}",
    "link.prompt.remember": "记住此协议",
    "link.prompt.cancel": "取消",
    "link.prompt.open": "打开链接",
    "link.status.confirm": "确认打开 {{protocol}} 链接。",
    "link.status.canceled": "已取消打开链接。",
    "link.status.opened": "已打开 {{href}}",
    "link.status.remembered": "已打开 {{href}}，并记住 {{protocol}} 链接。",
    "start.newFile": "新建文件",
    "start.newHelp": "创建未保存的 CSV",
    "start.openCsv": "打开 CSV/TSV",
    "start.openHelp": "在工具栏中输入路径",
    "start.chooseFile": "选择文件",
    "start.chooseHelp": "使用浏览器选择器打开",
    "start.settings": "设置",
    "start.settingsHelp": "调整应用偏好",
    "start.support": "支持",
    "start.supportHelp": "关于、更新和帮助",
    "about.title": "关于与支持",
    "about.product": "产品",
    "about.license": "许可",
    "about.licenseValue": "免费的本地编辑器",
    "about.network": "网络",
    "about.networkValue": "默认离线",
    "about.updates": "更新",
    "file.recent": "最近",
    "file.noOpen": "未打开文件",
    "file.status.creatingNew": "正在创建新 CSV...",
    "file.status.newReady": "新文件已就绪：{{name}}",
    "file.status.newFailed": "新建文件失败：{{message}}",
    "status.ready": "就绪",
    "status.languageChanged": "语言已更新。",
    "status.languageReset": "语言已重置。",
    "sql.run": "运行",
    "sql.copy": "复制",
    "sql.export": "导出",
    "sql.openResult": "打开结果",
    "sql.clear": "清除",
    "sql.empty.prompt": "运行 SQL 以查看结果。",
    "sql.empty.history": "没有查询历史",
    "sql.empty.rows": "没有行",
    "sql.error.title": "SQL 错误",
    "sql.error.location": "第 {{line}} 行，第 {{column}} 列",
    "sql.error.token": "标记 {{token}}",
    "sql.error.phase": "阶段 {{phase}}",
    "sql.status.enterQuery": "请输入 SQL 查询。",
    "sql.status.running": "正在运行 SQL...",
    "sql.status.runningElapsed": "正在运行 SQL... {{seconds}} 秒",
    "sql.status.runningRows": "正在加载{{phase}}... 已加载 {{rows}} 行，用时 {{seconds}} 秒",
    "sql.status.canceled": "SQL 已取消。",
    "sql.status.canceledElapsed": "SQL 已在 {{seconds}} 秒后取消。",
    "sql.status.returned": "SQL 返回 {{rows}} / {{total}} 行",
    "sql.status.returnedTimed": "SQL 返回 {{rows}} / {{total}} 行，用时 {{seconds}} 秒",
    "sql.status.returnedTimedRows": "SQL 已加载 {{loadedRows}} 个源数据行，返回 {{rows}} / {{total}} 行，用时 {{seconds}} 秒",
    "sql.progress.load": "源数据行",
    "sql.status.historySaveFailed": "无法保存 SQL 查询历史。",
    "sql.status.historyClearFailed": "无法清除 SQL 查询历史。",
    "sql.status.historyLoaded": "已从历史加载 SQL 查询。",
    "sql.status.failed": "SQL 失败：{{message}}",
    "sql.status.failedAt": "SQL 在第 {{line}} 行、第 {{column}} 列失败：{{message}}",
    "sql.status.runBeforeCopy": "复制结果前请先运行 SQL 查询。",
    "sql.status.runBeforeExport": "导出结果前请先运行 SQL 查询。",
    "sql.status.runBeforeOpen": "打开结果前请先运行 SQL 查询。",
    "sql.status.openingResult": "正在打开 SQL 结果...",
    "sql.status.openedResult": "已将 SQL 结果打开为 {{name}}：{{rows}} 行，{{cols}} 列",
    "sql.status.openResultFailed": "打开 SQL 结果失败：{{message}}",
    "sql.status.selectCell": "请先选择一个 SQL 结果单元格。",
    "sql.status.selectRow": "请先选择一个 SQL 结果行。",
    "sql.status.selectColumn": "请先选择一个 SQL 结果列。",
    "sql.status.noCopy": "没有可复制的{{label}}。",
    "sql.status.copied": "已复制{{label}}：{{rows}} 行，{{cols}} 列，格式 {{format}}",
    "sql.status.copyFailed": "{{label}}复制失败：{{message}}",
    "sql.status.exported": "已将 SQL 结果导出为 {{format}}：{{rows}} 行，{{cols}} 列，到 {{filename}}",
    "sql.status.selected": "已选择 SQL 结果：{{rows}} 行，{{cols}} 列",
    "sql.status.selectionCleared": "已清除 SQL 结果选择。",
    "sql.label.selection": "SQL 选区",
    "sql.label.result": "SQL 结果",
    "sql.label.cell": "SQL 单元格",
    "sql.label.row": "SQL 行",
    "sql.label.rowWithHeaders": "带表头的 SQL 行",
    "sql.label.column": "SQL 列",
    "sql.label.columnWithHeader": "带表头的 SQL 列",
    "sql.context.copyCell": "复制单元格",
    "sql.context.copyRow": "复制行",
    "sql.context.copyRowWithHeaders": "复制行并带表头",
    "sql.context.copyColumn": "复制列",
    "sql.context.copyColumnWithHeader": "复制列并带表头",
    "sql.context.copySelectionAs": "复制选区为 {{format}}",
    "sql.context.copyResultAs": "复制结果为 {{format}}",
    "sql.context.exportResultAs": "导出结果为 {{format}}",
    "sql.context.openResult": "将结果打开为文档",
    "sql.context.clearSelection": "清除选择",
    "profile.run": "运行概况",
    "profile.schema": "架构",
    "profile.saveSchema": "保存架构",
    "profile.resetSchema": "重置架构",
    "profile.apply": "应用",
    "profile.delete": "删除",
    "profile.export": "导出",
    "profile.import": "导入",
    "profile.cancel": "取消",
    "profile.schemaNamePlaceholder": "架构名称",
    "profile.savedSchemas": "已保存架构",
    "profile.noSchema": "没有架构",
    "profile.schemaRules": "{{count}} 条架构规则",
    "profile.invalidSchema": "架构 JSON 无效",
    "profile.rule.columnName": "列名",
    "profile.rule.columnNumber": "列号",
    "profile.rule.type": "类型",
    "profile.rule.required": "必填",
    "profile.rule.unique": "唯一",
    "profile.rule.enum": "枚举",
    "profile.rule.pattern": "模式",
    "profile.rule.default": "默认值",
    "profile.rule.min": "最小值",
    "profile.rule.max": "最大值",
    "profile.rule.minLength": "最小长度",
    "profile.rule.maxLength": "最大长度",
    "profile.rule.defaultSummary": "默认值 {{value}}",
    "profile.rule.save": "保存规则",
    "profile.rule.update": "更新规则",
    "profile.rule.clear": "清除规则",
    "profile.rule.none": "没有架构规则。",
    "profile.rule.edit": "编辑",
    "profile.rule.remove": "移除",
    "profile.rule.rule": "规则 {{index}}",
    "profile.rule.locationOnly": "仅定位",
    "profile.rule.requiredSummary": "必填",
    "profile.rule.uniqueSummary": "唯一",
    "profile.rule.typeSummary": "类型 {{type}}",
    "profile.rule.enumSummary": "枚举 {{values}}",
    "profile.rule.patternSummary": "模式 {{pattern}}",
    "profile.rule.minSummary": "最小值 {{value}}",
    "profile.rule.maxSummary": "最大值 {{value}}",
    "profile.rule.exclusiveMinSummary": "最小值 > {{value}}",
    "profile.rule.exclusiveMaxSummary": "最大值 < {{value}}",
    "profile.rule.minLengthSummary": "最小长度 {{value}}",
    "profile.rule.maxLengthSummary": "最大长度 {{value}}",
    "profile.type.any": "任意",
    "profile.type.string": "字符串",
    "profile.type.number": "数字",
    "profile.type.integer": "整数",
    "profile.type.decimal": "小数",
    "profile.type.boolean": "布尔",
    "profile.type.date": "日期",
    "profile.type.datetime": "日期时间",
    "profile.type.time": "时间",
    "profile.type.year": "年份",
    "profile.type.yearmonth": "年月",
    "profile.type.unknown": "未知",
    "profile.severity.error": "错误",
    "profile.severity.warning": "警告",
    "profile.quickFix.header.fill": "填充表头",
    "profile.quickFix.header.unique": "设为唯一",
    "profile.quickFix.row.pad": "补齐行",
    "profile.quickFix.schema.column.add": "添加列",
    "profile.quickFix.schema.default.use": "使用默认值",
    "profile.quickFix.schema.enum.useValue": "使用允许值",
    "profile.quickFix.schema.range.clampMin": "使用最小值",
    "profile.quickFix.schema.range.clampMax": "使用最大值",
    "profile.summary.status": "状态",
    "profile.summary.rows": "行",
    "profile.summary.columns": "列",
    "profile.summary.errors": "错误",
    "profile.summary.warnings": "警告",
    "profile.summary.rules": "规则",
    "profile.summary.progress": "已扫描",
    "profile.summary.progressValue": "{{passes}} 轮 / {{rows}} 次行扫描",
    "profile.loading": "加载中",
    "profile.columns.loading": "正在分析列...",
    "profile.issues.loading": "正在检查数据质量...",
    "profile.status.running": "正在分析数据...",
    "profile.status.runningElapsed": "正在分析数据... {{seconds}} 秒",
    "profile.status.runningRows": "正在执行{{phase}}... 已扫描 {{rows}} 行，用时 {{seconds}} 秒",
    "profile.status.canceled": "概况已取消。",
    "profile.status.canceledElapsed": "概况已在 {{seconds}} 秒后取消。",
    "profile.status.complete": "概况：{{rows}} 行，{{columns}} 列，{{issues}} 个问题",
    "profile.status.completeTimed": "概况：{{rows}} 行，{{columns}} 列，{{issues}} 个问题，用时 {{seconds}} 秒",
    "profile.status.completeTimedRows": "概况：{{rows}} 行，{{columns}} 列，{{issues}} 个问题，用时 {{seconds}} 秒（分析 {{profileRows}} 行，验证 {{validationRows}} 行）",
    "profile.status.failed": "概况失败：{{message}}",
    "profile.progress.profile": "分析阶段",
    "profile.progress.validation": "验证阶段",
    "profile.status.initialColumns": "运行概况以检查列。",
    "profile.status.initialIssues": "没有加载验证问题。",
    "profile.empty.noValidation": "没有验证结果。",
    "profile.empty.noColumns": "没有找到列。",
    "profile.empty.noIssues": "没有验证问题。",
    "profile.columns.column": "列",
    "profile.columns.type": "类型",
    "profile.columns.rows": "行",
    "profile.columns.nulls": "空值",
    "profile.columns.distinct": "唯一值",
    "profile.columns.min": "最小值",
    "profile.columns.max": "最大值",
    "profile.columns.average": "平均值",
    "profile.columns.topValues": "高频值",
    "profile.value.empty": "（空）",
    "profile.issues.severity": "级别",
    "profile.issues.type": "类型",
    "profile.issues.row": "行",
    "profile.issues.column": "列",
    "profile.issues.message": "消息",
    "profile.issues.fix": "修复",
    "profile.issues.defaultFix": "修复",
    "profile.status.selectedIssue": "已选择 {{type}}：第 {{row}} 行，第 {{column}} 列。",
    "profile.schema.statusInvalid": "架构无效：{{message}}",
    "profile.schema.saved": "验证架构已保存。",
    "profile.schema.savedAs": "验证架构已保存为 {{name}}。",
    "profile.schema.reset": "验证架构已重置。",
    "profile.schema.notFound": "未找到已保存架构。",
    "profile.schema.applied": "已应用架构 {{name}}。",
    "profile.schema.deleted": "已删除保存的架构。",
    "profile.schema.exported": "已导出 {{count}} 个 Profile 架构。",
    "profile.schema.imported": "已导入 {{count}} 个 Profile 架构。",
    "profile.schema.importFailed": "Profile 架构导入失败：{{message}}",
    "profile.rule.chooseColumn": "请选择列名或列号。",
    "profile.rule.invalid": "规则无效：{{message}}",
    "profile.rule.added": "已添加架构规则。",
    "profile.rule.updated": "已更新架构规则。",
    "profile.rule.cleared": "已清除架构规则编辑器。",
    "profile.rule.notFound": "未找到架构规则。",
    "profile.rule.editing": "正在编辑架构规则 {{index}}。",
    "profile.rule.removed": "已移除架构规则。",
    "profile.quickFix.readOnly": "只读模式会阻止快速修复。",
    "profile.quickFix.noFile": "请先打开文件再应用修复。",
    "profile.quickFix.unavailable": "这个问题没有安全的快速修复。",
    "profile.quickFix.applied": "已将快速修复应用到 {{count}} 个单元格。",
    "profile.quickFix.failed": "快速修复失败：{{message}}",
    "command.command.palette.title": "命令面板",
    "command.command.palette.description": "搜索并运行可用命令。",
    "command.settings.open.title": "设置",
    "command.settings.open.description": "打开应用设置和文件格式预设管理。",
    "command.settings.keyboardShortcuts.title": "键盘快捷键",
    "command.settings.keyboardShortcuts.description": "自定义命令快捷键。",
    "command.help.aboutSupport.title": "关于与支持",
    "command.help.aboutSupport.description": "打开本地关于、支持和更新入口。",
    "command.file.new.title": "新建",
    "command.file.new.description": "创建新的未保存 CSV 文档。",
    "command.file.open.title": "打开",
    "command.file.open.description": "打开 CSV 或分隔文本文件。",
    "command.file.save.title": "保存",
    "command.file.save.description": "按备份策略保存当前文档。",
    "command.file.saveAs.title": "另存为",
    "command.file.saveAs.description": "将当前文档保存到指定本地路径。",
    "command.file.exportView.title": "导出视图",
    "command.file.exportView.description": "下载当前筛选、排序和可见列视图。",
    "command.file.exportXlsx.title": "导出 XLSX",
    "command.file.exportXlsx.description": "将当前视图下载为 Excel 工作簿。",
    "command.file.toggleReadOnly.title": "切换只读",
    "command.file.toggleReadOnly.description": "阻止或允许编辑当前文档。",
    "command.diff.open.title": "比较",
    "command.diff.open.description": "并排比较两个 CSV 文件。",
    "command.view.sql.title": "SQL 控制台",
    "command.view.sql.description": "打开当前工作区的 SQL 控制台。",
    "command.view.profile.title": "概况",
    "command.view.profile.description": "显示列概况和数据质量指标。",
    "command.profile.run.title": "运行概况",
    "command.profile.run.description": "运行当前文档的列概况和验证。",
    "command.profile.cancel.title": "取消概况",
    "command.profile.cancel.description": "取消正在进行的概况分析。",
    "command.view.filter.title": "筛选",
    "command.view.filter.description": "按文本或 SQL WHERE 表达式筛选行。",
    "command.view.filterBuilder.title": "筛选构建器",
    "command.view.filterBuilder.description": "无需编写 SQL 即可构建列筛选。",
    "command.view.sort.title": "排序",
    "command.view.sort.description": "排序选中的列。",
    "command.column.sortAscending.title": "按列升序排序",
    "command.column.sortAscending.description": "将当前列按升序排序。",
    "command.column.sortDescending.title": "按列降序排序",
    "command.column.sortDescending.description": "将当前列按降序排序。",
    "command.column.sortLengthAscending.title": "按列长度升序排序",
    "command.column.sortLengthAscending.description": "将当前列按文本长度升序排序。",
    "command.column.sortLengthDescending.title": "按列长度降序排序",
    "command.column.sortLengthDescending.description": "将当前列按文本长度降序排序。",
    "command.column.addSortAscending.title": "追加列升序排序",
    "command.column.addSortAscending.description": "将当前列追加为升序次级排序键。",
    "command.column.addSortDescending.title": "追加列降序排序",
    "command.column.addSortDescending.description": "将当前列追加为降序次级排序键。",
    "command.view.removeLastSortKey.title": "移除最后一个排序键",
    "command.view.removeLastSortKey.description": "从当前多列排序中移除最后一个排序键。",
    "command.column.filterBySelectedValue.title": "按选中值筛选",
    "command.column.filterBySelectedValue.description": "将当前列筛选为选中单元格的值。",
    "command.column.filterBlanks.title": "只显示空值",
    "command.column.filterBlanks.description": "将当前列筛选为空单元格。",
    "command.column.filterNonBlanks.title": "只显示非空值",
    "command.column.filterNonBlanks.description": "将当前列筛选为非空单元格。",
    "command.edit.copy.title": "复制",
    "command.edit.copy.description": "以制表符分隔文本复制选中单元格。",
    "command.edit.copyWithHeaders.title": "复制并带表头",
    "command.edit.copyWithHeaders.description": "复制选中单元格并包含表头行。",
    "command.edit.copyCsv.title": "复制为 CSV",
    "command.edit.copyMarkdown.title": "复制为 Markdown",
    "command.edit.copyJson.title": "复制为 JSON",
    "command.edit.copyHtml.title": "复制为 HTML",
    "command.edit.copySql.title": "复制为 SQL",
    "command.edit.copyLatex.title": "复制为 LaTeX",
    "command.edit.copyRaw.title": "复制原始值",
    "command.edit.find.title": "查找",
    "command.edit.find.description": "查找匹配查询的单元格。",
    "command.edit.replacePreview.title": "预览替换",
    "command.edit.replacePreview.description": "替换文本前预览会变更的单元格。",
    "command.edit.replace.title": "替换",
    "command.edit.replace.description": "替换当前匹配项。",
    "command.edit.replaceAll.title": "全部替换",
    "command.edit.replaceAll.description": "替换当前文档中的所有匹配项。",
    "command.selection.cycleAggregate.title": "切换选区聚合",
    "command.selection.cycleAggregate.description": "切换当前选区的状态栏聚合值。",
    "command.selection.copyAggregate.title": "复制选区聚合",
    "command.selection.copyAggregate.description": "复制当前状态栏聚合值。",
    "command.selection.insertAggregateRow.title": "添加聚合行",
    "command.selection.insertAggregateRow.description": "为当前选区添加一行可移除的聚合行，不修改文件数据。",
    "command.selection.insertSumRow.title": "添加求和行",
    "command.selection.insertSumRow.description": "为选中的数值单元格添加可移除的求和行。",
    "command.selection.insertAverageRow.title": "添加平均行",
    "command.selection.insertAverageRow.description": "为选中的数值单元格添加可移除的平均值行。",
    "command.selection.insertMinRow.title": "添加最小值行",
    "command.selection.insertMinRow.description": "为选中的数值单元格添加可移除的最小值行。",
    "command.selection.insertMaxRow.title": "添加最大值行",
    "command.selection.insertMaxRow.description": "为选中的数值单元格添加可移除的最大值行。",
    "command.selection.insertCountRow.title": "添加计数行",
    "command.selection.insertCountRow.description": "为选中的单元格添加可移除的计数行。",
    "command.table.transpose.title": "转置",
    "command.table.transpose.description": "转置当前表格或选中的一个或多个矩形区域。",
    "command.table.dedupeRows.title": "删除重复行",
    "command.table.dedupeColumns.title": "删除重复列",
    "command.text.trim.title": "修剪文本",
    "command.text.trim.description": "修剪选中单元格首尾的空白字符。",
    "command.text.uppercase.title": "转大写",
    "command.text.uppercase.description": "将选中单元格转换为大写文本。",
    "command.text.lowercase.title": "转小写",
    "command.text.lowercase.description": "将选中单元格转换为小写文本。",
    "command.text.snakeCase.title": "蛇形命名",
    "command.text.snakeCase.description": "将选中单元格转换为 snake_case 命名。",
    "command.text.kebabCase.title": "短横线命名",
    "command.text.kebabCase.description": "将选中单元格转换为 kebab-case 命名。",
    "command.text.camelCase.title": "驼峰命名",
    "command.text.camelCase.description": "将选中单元格转换为 camelCase 命名。",
    "command.text.pascalCase.title": "帕斯卡命名",
    "command.text.pascalCase.description": "将选中单元格转换为 PascalCase 命名。",
    "command.text.titleCase.title": "标题大小写",
    "command.text.titleCase.description": "将选中单元格转换为 Title Case 命名。",
    "command.text.fullWidth.title": "转全角",
    "command.text.fullWidth.description": "将选中单元格转换为全角 Unicode 文本。",
    "command.text.halfWidth.title": "转半角",
    "command.text.halfWidth.description": "将选中单元格转换为半角 Unicode 文本。",
    "command.text.normalizeNFC.title": "标准化 NFC",
    "command.text.normalizeNFC.description": "将选中单元格标准化为 Unicode NFC 组合形式。",
    "command.text.normalizeNFD.title": "标准化 NFD",
    "command.text.normalizeNFD.description": "将选中单元格标准化为 Unicode NFD 分解形式。",
    "command.text.normalizeNFKC.title": "标准化 NFKC",
    "command.text.normalizeNFKC.description": "将选中单元格标准化为 Unicode NFKC 兼容组合形式。",
    "command.text.normalizeNFKD.title": "标准化 NFKD",
    "command.text.normalizeNFKD.description": "将选中单元格标准化为 Unicode NFKD 兼容分解形式。",
    "command.text.column.trim.title": "修剪整列文本",
    "command.text.column.trim.description": "修剪当前列首尾的空白字符。",
    "command.text.column.uppercase.title": "整列转大写",
    "command.text.column.uppercase.description": "将当前列转换为大写文本。",
    "command.text.column.lowercase.title": "整列转小写",
    "command.text.column.lowercase.description": "将当前列转换为小写文本。",
    "command.text.column.snakeCase.title": "整列转蛇形命名",
    "command.text.column.snakeCase.description": "将当前列转换为 snake_case 命名。",
    "command.text.column.kebabCase.title": "整列转短横线命名",
    "command.text.column.kebabCase.description": "将当前列转换为 kebab-case 命名。",
    "command.text.column.camelCase.title": "整列转驼峰命名",
    "command.text.column.camelCase.description": "将当前列转换为 camelCase 命名。",
    "command.text.column.pascalCase.title": "整列转帕斯卡命名",
    "command.text.column.pascalCase.description": "将当前列转换为 PascalCase 命名。",
    "command.text.column.titleCase.title": "整列转标题大小写",
    "command.text.column.titleCase.description": "将当前列转换为 Title Case 命名。",
    "command.text.column.fullWidth.title": "整列转全角",
    "command.text.column.fullWidth.description": "将当前列转换为全角 Unicode 文本。",
    "command.text.column.halfWidth.title": "整列转半角",
    "command.text.column.halfWidth.description": "将当前列转换为半角 Unicode 文本。",
    "command.text.column.normalizeNFC.title": "整列标准化 NFC",
    "command.text.column.normalizeNFC.description": "将当前列标准化为 Unicode NFC 组合形式。",
    "command.text.column.normalizeNFD.title": "整列标准化 NFD",
    "command.text.column.normalizeNFD.description": "将当前列标准化为 Unicode NFD 分解形式。",
    "command.text.column.normalizeNFKC.title": "整列标准化 NFKC",
    "command.text.column.normalizeNFKC.description": "将当前列标准化为 Unicode NFKC 兼容组合形式。",
    "command.text.column.normalizeNFKD.title": "整列标准化 NFKD",
    "command.text.column.normalizeNFKD.description": "将当前列标准化为 Unicode NFKD 兼容分解形式。",
    "command.cell.openLink.title": "打开链接",
    "command.cell.openLink.description": "打开当前单元格中通过安全审计的链接。",
    "command.cell.joinHorizontal.title": "合并单元格",
    "command.cell.splitHorizontal.title": "拆分单元格",
    "command.cell.joinVertical.title": "纵向合并单元格",
    "command.cell.splitVertical.title": "纵向拆分单元格",
    "command.row.insert.title": "插入行",
    "command.row.delete.title": "删除行",
    "command.row.cut.title": "剪切行",
    "command.row.paste.title": "粘贴行",
    "command.row.showAll.title": "显示所有行",
    "command.row.showAll.description": "恢复当前行视图中隐藏的行。",
    "command.column.insert.title": "插入列",
    "command.column.delete.title": "删除列",
    "command.column.cut.title": "剪切列",
    "command.column.paste.title": "粘贴列",
    "command.view.clear.title": "清除视图",
  },
  pt: {
    "nav.grid": "Grade",
    "nav.profile": "Perfil",
    "nav.settings": "Configurações",
    "file.open": "Abrir",
    "edit.copy": "Copiar",
    "view.filter": "Filtrar",
    "view.sort": "Ordenar",
    "edit.find": "Localizar",
    "edit.replace": "Substituir",
    "diff.compare": "Comparar",
    "file.export": "Exportar",
    "file.save": "Salvar",
    "settings.title": "Configurações",
    "settings.language": "Idioma",
    "settings.resetLanguage": "Redefinir idioma",
    "settings.searchPlaceholder": "Pesquisar configurações",
    "start.openCsv": "Abrir CSV/TSV",
    "start.chooseFile": "Escolher arquivo",
    "file.recent": "Recentes",
    "status.ready": "Pronto",
  },
  ru: {
    "nav.grid": "Таблица",
    "nav.profile": "Профиль",
    "nav.settings": "Настройки",
    "file.open": "Открыть",
    "edit.copy": "Копировать",
    "view.filter": "Фильтр",
    "view.sort": "Сортировка",
    "edit.find": "Найти",
    "edit.replace": "Заменить",
    "diff.compare": "Сравнить",
    "file.export": "Экспорт",
    "file.save": "Сохранить",
    "settings.title": "Настройки",
    "settings.language": "Язык",
    "settings.resetLanguage": "Сбросить язык",
    "settings.searchPlaceholder": "Поиск настроек",
    "start.openCsv": "Открыть CSV/TSV",
    "start.chooseFile": "Выбрать файл",
    "file.recent": "Недавние",
    "status.ready": "Готово",
  },
  de: {
    "nav.grid": "Raster",
    "nav.profile": "Profil",
    "nav.settings": "Einstellungen",
    "file.open": "Öffnen",
    "edit.copy": "Kopieren",
    "view.filter": "Filtern",
    "view.sort": "Sortieren",
    "edit.find": "Suchen",
    "edit.replace": "Ersetzen",
    "diff.compare": "Vergleichen",
    "file.export": "Exportieren",
    "file.save": "Speichern",
    "settings.title": "Einstellungen",
    "settings.language": "Sprache",
    "settings.resetLanguage": "Sprache zurücksetzen",
    "settings.searchPlaceholder": "Einstellungen suchen",
    "start.openCsv": "CSV/TSV öffnen",
    "start.chooseFile": "Datei wählen",
    "file.recent": "Zuletzt",
    "status.ready": "Bereit",
  },
  ko: {
    "nav.grid": "그리드",
    "nav.profile": "프로필",
    "nav.settings": "설정",
    "file.open": "열기",
    "edit.copy": "복사",
    "view.filter": "필터",
    "view.sort": "정렬",
    "edit.find": "찾기",
    "edit.replace": "바꾸기",
    "diff.compare": "비교",
    "file.export": "내보내기",
    "file.save": "저장",
    "settings.title": "설정",
    "settings.language": "언어",
    "settings.resetLanguage": "언어 재설정",
    "settings.searchPlaceholder": "설정 검색",
    "start.openCsv": "CSV/TSV 열기",
    "start.chooseFile": "파일 선택",
    "file.recent": "최근",
    "status.ready": "준비됨",
  },
};
const customShortcutsStorageKey = "opencsv-custom-shortcuts";
const shortcutExportKind = "opencsv-shortcuts";
const shortcutContextOptions = [
  { id: "", label: "Always" },
  { id: "document", label: "Document" },
  { id: "grid", label: "Grid" },
  { id: "sql", label: "SQL" },
  { id: "profile", label: "Profile" },
  { id: "settings", label: "Settings" },
  { id: "diff", label: "Diff" },
  { id: "typing", label: "Typing" },
];
const shortcutContextAtoms = new Set(shortcutContextOptions.map((option) => option.id).filter(Boolean));
const shortcutViewContexts = new Set(["grid", "sql", "profile", "settings", "diff"]);
const shortcutPresets = [
  { id: "default", name: "Default", shortcuts: {} },
  {
    id: "analysis",
    name: "Analysis",
    shortcuts: {
      "diff.open": "Shift+Alt+O",
      "selection.copyAggregate": "Shift+Alt+K",
      "table.transpose": "Shift+Alt+E",
      "view.profile": "Shift+Alt+P",
      "view.sql": "Shift+Alt+Q",
    },
  },
  {
    id: "editing",
    name: "Editing",
    shortcuts: {
      "file.toggleReadOnly": "Shift+Alt+Y",
      "format.applyQuotePolicy": "Shift+Alt+Q",
      "format.toggleCellQuote": "Shift+Alt+P",
      "text.fullWidth": "Shift+Alt+N",
      "text.halfWidth": "Shift+Alt+Z",
    },
  },
];
const sqlHistoryStorageKey = "opencsv-sql-history";
const sqlHistoryLimit = 8;
const savedFiltersStorageKey = "opencsv-saved-filters";
const savedFiltersLimit = 12;
const recentFiltersStorageKey = "opencsv-recent-filters";
const filterHistoryLimitStorageKey = "opencsv-filter-history-limit";
const defaultFilterHistoryLimit = 8;
const filterHistoryLimitOptions = [0, 5, 8, 12, 20];
const savedSortsStorageKey = "opencsv-saved-sorts";
const savedSortsLimit = 12;
const savedViewsExportKind = "opencsv-saved-views";
const formatPresetsStorageKey = "opencsv-format-presets";
const formatPresetsLimit = 12;
const formatPresetExportKind = "opencsv-format-presets";
const quoteRulesExportKind = "opencsv-quote-rules";
const editorAppearanceStorageKey = "opencsv-editor-appearance";
const scrollPositionsStorageKey = "opencsv-scroll-positions";
const scrollPositionsLimit = 32;
const windowGeometryStorageKey = "opencsv-window-geometry";
const windowGeometrySaveDelay = 200;
const hiddenColumnsStorageKey = "opencsv-hidden-columns";
const hiddenColumnsLimit = 64;
const editorGroupLayoutStorageKey = "opencsv-editor-group-layout";
const dirtyDraftsStorageKey = "opencsv-dirty-drafts";
const dirtyDraftsLimit = 5;
const dirtyDraftMaxBytes = 1_000_000;
const cellSeparatorStorageKey = "opencsv-cell-separator";
const defaultCellSeparatorToken = "space";
const autoRefreshStorageKey = "opencsv-auto-refresh-mode";
const defaultAutoRefreshMode = "clean";
const autoRefreshModes = ["clean", "notify"];
const updateCheckStorageKey = "opencsv-update-check-mode";
const updateCheckResultStorageKey = "opencsv-update-check-result";
const defaultUpdateCheckMode = "manual";
const updateCheckModes = ["automatic", "manual", "off"];
const blockedLinkProtocols = new Set(["javascript:", "data:", "vbscript:", "blob:", "about:"]);
const trustedLinkProtocols = new Set(["http:", "https:", "mailto:", "file:"]);
const customLinkProtocolsStorageKey = "opencsv-custom-link-protocols";
const confirmedCustomLinkProtocolsStorageKey = "opencsv-confirmed-custom-link-protocols";
const defaultCustomLinkProtocols = ["obsidian", "vscode"];
const themeStorageKey = "opencsv-theme";
const defaultThemePreference = "system";
const themePreferences = ["light", "dark", "high-contrast", "system"];
const formatRuleSampleBytes = 65536;
const editorFontFamilies = {
  mono: 'ui-monospace, "SFMono-Regular", Menlo, Consolas, "Liberation Mono", monospace',
  sans: 'Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  serif: 'Georgia, "Times New Roman", Times, serif',
};
const defaultEditorAppearance = {
  gridFont: "mono",
  gridSize: 13,
  gridRowHeight: 28,
  editorFont: "mono",
  editorSize: 13,
};
const aggregateRowCommandMetrics = new Map([
  ["selection.insertSumRow", "sum"],
  ["selection.insertAverageRow", "avg"],
  ["selection.insertMinRow", "min"],
  ["selection.insertMaxRow", "max"],
  ["selection.insertCountRow", "count"],
]);
const aggregateScanCellLimit = 5000;
const columnTextTransformSpecs = [
  { id: "text.column.trim", mode: "trim", label: "Trim Column" },
  { id: "text.column.uppercase", mode: "upper", label: "Uppercase Column" },
  { id: "text.column.lowercase", mode: "lower", label: "Lowercase Column" },
  { id: "text.column.snakeCase", mode: "snake", label: "Snake Case Column" },
  { id: "text.column.kebabCase", mode: "kebab", label: "Kebab Case Column" },
  { id: "text.column.camelCase", mode: "camel", label: "Camel Case Column" },
  { id: "text.column.pascalCase", mode: "pascal", label: "Pascal Case Column" },
  { id: "text.column.titleCase", mode: "title", label: "Title Case Column" },
  { id: "text.column.fullWidth", mode: "wide", label: "Full Width Column" },
  { id: "text.column.halfWidth", mode: "narrow", label: "Half Width Column" },
  { id: "text.column.normalizeNFC", mode: "nfc", label: "Normalize Column NFC" },
  { id: "text.column.normalizeNFD", mode: "nfd", label: "Normalize Column NFD" },
  { id: "text.column.normalizeNFKC", mode: "nfkc", label: "Normalize Column NFKC" },
  { id: "text.column.normalizeNFKD", mode: "nfkd", label: "Normalize Column NFKD" },
];
const columnTextTransformCommands = new Map(columnTextTransformSpecs.map((spec) => [spec.id, spec]));
const columnTextTransformCommandByMode = new Map(columnTextTransformSpecs.map((spec) => [spec.mode, spec.id]));
const columnSortCommands = new Map([
  ["column.sortAscending", { descending: false, append: false, mode: "" }],
  ["column.sortDescending", { descending: true, append: false, mode: "" }],
  ["column.sortLengthAscending", { descending: false, append: false, mode: "length" }],
  ["column.sortLengthDescending", { descending: true, append: false, mode: "length" }],
  ["column.addSortAscending", { descending: false, append: true, mode: "" }],
  ["column.addSortDescending", { descending: true, append: true, mode: "" }],
]);
const columnFilterCommands = new Map([
  ["column.filterBySelectedValue", { kind: "value" }],
  ["column.filterBlanks", { kind: "blank" }],
  ["column.filterNonBlanks", { kind: "nonblank" }],
]);

let totalRows = 0;
let sourceRowCount = 0;
let aggregateRows = [];
let columns = [];
let allColumns = [];
let hiddenSourceColumns = new Set();

let scrollX = 0;
let scrollY = 0;
let wheelRemainderY = 0;
let scrollPositionSaveTimer = 0;
let windowGeometrySaveTimer = 0;
let restoredWindowGeometry = loadWindowGeometry();
let selected = { row: 0, col: 0 };
let selectionAnchor = { row: 0, col: 0 };
let smartSelectionRanges = [];
let draggingSelection = false;
let dragSelectionMode = "cell";
let smartSelectionGesture = "";
let smartSelectionGestureBaseRanges = [];
let pendingColumnReorder = null;
let pendingRowReorder = null;
let draggingFillHandle = false;
let hoverFillHandle = false;
let fillDragSourceRange = null;
let fillDragDirection = "";
let hoverURL = "";
let hoverResizeColumn = -1;
let resizingColumn = -1;
let resizeStartClientX = 0;
let resizeStartWidth = 0;
let editing = false;
let structureClipboard = null;
let documentId = "";
let activeFileEntry = null;
let openFileTabs = [];
let metadata = null;
let rows = new Map();
let rowSource = new Map();
let windowStart = 0;
let windowEnd = 0;
let loadingWindow = false;
let windowGeneration = 0;
let dirty = false;
let dirtyDraftSaveTimer = 0;
let dirtyDraftSaveInFlight = false;
let dataStartRow = 0;
let viewMapped = false;
let sortDescending = false;
let sortKeys = [];
let sortStable = true;
let commandsById = new Map();
let viewMode = "grid";
let editorGroupLayout = loadEditorGroupLayout();
let editorGroupsSplit = editorGroupLayout.split;
let editorGroupOrientation = editorGroupLayout.orientation;
let activeEditorGroup = editorGroupLayout.activeGroup;
let currentEditorTabGroup = editorGroupLayout.currentTabGroup;
let editorGroupActiveTabs = { primary: "", secondary: "" };
let pendingEditorTabDrag = null;
let profileData = null;
let lastSQLResult = null;
let sqlResultSelectionAnchor = null;
let sqlResultSelectionRange = null;
let sqlResultContextTarget = null;
let profileLoading = false;
let activeProfileController = null;
let profileProgressTimer = 0;
let profileProgressStartedAt = 0;
let profileProgressEventCount = 0;
let profileSchemaVisible = false;
let profileSchemaText = localStorage.getItem(profileSchemaStorageKey) || "";
let savedProfileSchemasState = loadSavedProfileSchemas();
let editingProfileRuleIndex = -1;
let paletteOpen = false;
let paletteItems = [];
let paletteIndex = 0;
let customShortcuts = new Map();
let shortcutEditorOpen = false;
let recordingShortcutID = "";
let pendingShortcutConflict = null;
let conditionBuilderOpen = false;
let conditionExtraRows = [];
let columnMenuColumn = -1;
let readOnly = false;
let externalFileChanged = false;
let externalChangePanelDismissed = false;
let fileStatusTimer = 0;
let pendingWholeTranspose = null;
let pendingSplitTool = null;
let pendingDirtyCloseTab = null;
let customLinkProtocols = loadCustomLinkProtocols();
let confirmedCustomLinkProtocols = loadConfirmedCustomLinkProtocols();
let pendingLinkConfirmation = null;
let cellValueInputDirty = false;
let cellValueInputKey = "";
let cellValueExpanded = false;
let cellValueDialogOpen = false;
let cellSeparatorToken = normalizeCellSeparatorToken(localStorage.getItem(cellSeparatorStorageKey));
let findCaseSensitive = false;
let findWholeCell = false;
let findRegex = false;
let findInSelection = false;
let findState = emptyFindState();
let activeFindController = null;
let findProgressTimer = 0;
let findProgressStartedAt = 0;
let findProgressQuery = "";
let activeSQLController = null;
let sqlProgressTimer = 0;
let sqlProgressStartedAt = 0;
let sqlProgressEventCount = 0;
let activeFilter = null;
let includeHeadersInSelection = false;
let activeAggregateMetric = "sum";
let aggregateStatsCacheKey = "";
let aggregateStatsCache = null;
let diffRows = [];
let activeDiffIndex = -1;
let lastDiffResult = null;
let quotePolicyDraftDirty = false;
let delimiterDraftDirty = false;
let encodingDraftDirty = false;
let lineEndingDraftDirty = false;
const themeMediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
let localePreference = normalizeLocalePreference(localStorage.getItem(localeStorageKey));
let themePreference = normalizeThemePreference(localStorage.getItem(themeStorageKey));
let autoRefreshMode = normalizeAutoRefreshMode(localStorage.getItem(autoRefreshStorageKey));
let updateCheckMode = normalizeUpdateCheckMode(localStorage.getItem(updateCheckStorageKey));
let updateCheckState = loadUpdateCheckState();
let filterHistoryLimit = normalizeFilterHistoryLimit(localStorage.getItem(filterHistoryLimitStorageKey));
let themeColors = {};
let editorAppearance = loadEditorAppearance();

initializeTheme();
setLocalePreference(localePreference, false);
setAutoRefreshMode(autoRefreshMode, false);
setUpdateCheckMode(updateCheckMode, false);
applyEditorAppearance(editorAppearance, false);
setCellSeparatorToken(cellSeparatorToken, false);
document.body.dataset.editorGroupLayoutRestored = editorGroupLayoutSignature(editorGroupLayout);
syncProfileSchemaControls();
syncAppTitle();
syncWindowGeometryDataset();

function resize() {
  queueWindowGeometrySave();
  draw();
}

function syncCanvasSize() {
  const rect = gridFrame.getBoundingClientRect();
  const scale = window.devicePixelRatio || 1;
  const width = Math.max(1, rect.width);
  const height = Math.max(1, rect.height);
  syncCanvasCSSSize(canvas, width, height);
  resizeCanvasBackingStore(canvas, width, height, scale);
  ctx.setTransform(scale, 0, 0, scale, 0, 0);
  syncCanvasCSSSize(headerCanvas, width, headerHeight);
  resizeCanvasBackingStore(headerCanvas, width, headerHeight, scale);
  headerCtx.setTransform(scale, 0, 0, scale, 0, 0);
  return { width, height };
}

function syncCanvasCSSSize(target, width, height) {
  const cssWidth = `${Math.round(width)}px`;
  const cssHeight = `${Math.round(height)}px`;
  if (target.style.width !== cssWidth) target.style.width = cssWidth;
  if (target.style.height !== cssHeight) target.style.height = cssHeight;
}

function resizeCanvasBackingStore(target, width, height, scale) {
  const pixelWidth = Math.max(1, Math.round(width * scale));
  const pixelHeight = Math.max(1, Math.round(height * scale));
  if (target.width !== pixelWidth) target.width = pixelWidth;
  if (target.height !== pixelHeight) target.height = pixelHeight;
}

function maxVerticalScroll() {
  if (totalRows <= 0) return 0;
  const viewportRowsHeight = Math.max(0, canvas.clientHeight - headerHeight);
  const fullVisibleRows = Math.max(1, Math.floor(viewportRowsHeight / rowHeight));
  return Math.max(0, totalRows - fullVisibleRows) * rowHeight;
}

function maxHorizontalScroll() {
  const totalWidth = columns.reduce((sum, column) => sum + column.width, 0);
  return Math.max(0, totalWidth - canvas.clientWidth + rowNumberWidth);
}

function scrollPositionKey() {
  const path = metadata?.path || activeFileEntry?.detail || "";
  return path && path !== "temporary upload" ? path : "";
}

function scrollPositionEntry() {
  const path = scrollPositionKey();
  if (!documentId || !path) return null;
  return {
    path,
    scrollX: Math.max(0, Math.round(scrollX)),
    scrollY: Math.max(0, Math.round(scrollY)),
    updatedAt: Date.now(),
  };
}

function scrollPositions() {
  try {
    const parsed = JSON.parse(localStorage.getItem(scrollPositionsStorageKey) || "[]");
    return Array.isArray(parsed) ? parsed : [];
  } catch (error) {
    return [];
  }
}

function saveScrollPositionEntry(entry) {
  if (!entry) return;
  const next = [
    entry,
    ...scrollPositions().filter((item) => item?.path && item.path !== entry.path),
  ].slice(0, scrollPositionsLimit);
  localStorage.setItem(scrollPositionsStorageKey, JSON.stringify(next));
}

function updateScrollPositionDataset() {
  if (!documentId) {
    delete document.body.dataset.scrollPosition;
    delete document.body.dataset.scrollPositionKey;
    return;
  }
  const key = scrollPositionKey();
  document.body.dataset.scrollPosition = `${Math.round(scrollX)}:${Math.round(scrollY)}`;
  document.body.dataset.scrollPositionKey = key;
}

function queueScrollPositionSave() {
  updateScrollPositionDataset();
  const entry = scrollPositionEntry();
  if (!entry) return;
  if (scrollPositionSaveTimer) window.clearTimeout(scrollPositionSaveTimer);
  scrollPositionSaveTimer = window.setTimeout(() => {
    scrollPositionSaveTimer = 0;
    saveScrollPositionEntry(entry);
  }, 200);
}

function saveCurrentScrollPosition() {
  if (scrollPositionSaveTimer) {
    window.clearTimeout(scrollPositionSaveTimer);
    scrollPositionSaveTimer = 0;
  }
  saveScrollPositionEntry(scrollPositionEntry());
}

function normalizeWindowGeometry(value) {
  if (!value || typeof value !== "object") return null;
  const width = Math.round(Number(value.width));
  const height = Math.round(Number(value.height));
  if (!Number.isFinite(width) || !Number.isFinite(height) || width <= 0 || height <= 0) return null;
  const x = Math.round(Number(value.x));
  const y = Math.round(Number(value.y));
  return {
    width,
    height,
    x: Number.isFinite(x) ? x : 0,
    y: Number.isFinite(y) ? y : 0,
    maximized: Boolean(value.maximized),
    updatedAt: Math.max(0, Math.round(Number(value.updatedAt)) || 0),
  };
}

function loadWindowGeometry() {
  try {
    return normalizeWindowGeometry(JSON.parse(localStorage.getItem(windowGeometryStorageKey) || "null"));
  } catch {
    return null;
  }
}

function currentWindowGeometry() {
  const width = Math.max(1, Math.round(window.outerWidth || window.innerWidth || document.documentElement.clientWidth || 1));
  const height = Math.max(1, Math.round(window.outerHeight || window.innerHeight || document.documentElement.clientHeight || 1));
  const x = Math.round(Number.isFinite(window.screenX) ? window.screenX : window.screenLeft || 0);
  const y = Math.round(Number.isFinite(window.screenY) ? window.screenY : window.screenTop || 0);
  const availableWidth = Math.round(window.screen?.availWidth || 0);
  const availableHeight = Math.round(window.screen?.availHeight || 0);
  const maximized = availableWidth > 0 && availableHeight > 0
    && Math.abs(width - availableWidth) <= 12
    && Math.abs(height - availableHeight) <= 96;
  return { width, height, x, y, maximized, updatedAt: Date.now() };
}

function windowGeometrySignature(entry) {
  const geometry = normalizeWindowGeometry(entry);
  if (!geometry) return "";
  return `${geometry.width}x${geometry.height}+${geometry.x}+${geometry.y}:${geometry.maximized ? "max" : "normal"}`;
}

function syncWindowGeometryDataset(current = currentWindowGeometry()) {
  document.body.dataset.windowGeometry = windowGeometrySignature(current);
  document.body.dataset.windowGeometryRestored = windowGeometrySignature(restoredWindowGeometry);
  document.body.dataset.windowGeometryMaximized = String(Boolean(current?.maximized));
  document.body.dataset.windowGeometryUpdatedAt = String(Math.max(0, Math.round(Number(current?.updatedAt)) || 0));
}

function saveWindowGeometryNow() {
  if (windowGeometrySaveTimer) {
    window.clearTimeout(windowGeometrySaveTimer);
    windowGeometrySaveTimer = 0;
  }
  const entry = currentWindowGeometry();
  try {
    localStorage.setItem(windowGeometryStorageKey, JSON.stringify(entry));
    document.body.dataset.windowGeometrySaved = windowGeometrySignature(entry);
  } catch {
    status.textContent = "Window geometry could not be saved.";
  }
  syncWindowGeometryDataset(entry);
}

function queueWindowGeometrySave() {
  syncWindowGeometryDataset();
  if (windowGeometrySaveTimer) window.clearTimeout(windowGeometrySaveTimer);
  windowGeometrySaveTimer = window.setTimeout(saveWindowGeometryNow, windowGeometrySaveDelay);
}

function restoreScrollPosition() {
  const key = scrollPositionKey();
  if (!documentId || !key) {
    delete document.body.dataset.scrollPositionRestored;
    updateScrollPositionDataset();
    return false;
  }
  const entry = scrollPositions().find((item) => item?.path === key);
  if (!entry) {
    delete document.body.dataset.scrollPositionRestored;
    updateScrollPositionDataset();
    return false;
  }
  scrollX = Math.max(0, Math.min(maxHorizontalScroll(), Number(entry.scrollX) || 0));
  scrollY = Math.max(0, Math.min(maxVerticalScroll(), Number(entry.scrollY) || 0));
  wheelRemainderY = 0;
  document.body.dataset.scrollPositionRestored = `${Math.round(scrollX)}:${Math.round(scrollY)}`;
  updateScrollPositionDataset();
  return true;
}

function hiddenColumnsKey() {
  const path = metadata?.path || activeFileEntry?.detail || "";
  return path && path !== "temporary upload" ? path : "";
}

function hiddenColumnStates() {
  try {
    const parsed = JSON.parse(localStorage.getItem(hiddenColumnsStorageKey) || "[]");
    return Array.isArray(parsed) ? parsed : [];
  } catch (error) {
    return [];
  }
}

function normalizedHiddenColumnList(columnsList) {
  const sourceColumns = new Set(allColumns.map((column) => column.sourceCol));
  const seen = new Set();
  const result = [];
  for (const value of columnsList ?? []) {
    const sourceCol = Number(value);
    if (!Number.isInteger(sourceCol) || sourceCol < 0 || !sourceColumns.has(sourceCol) || seen.has(sourceCol)) {
      continue;
    }
    seen.add(sourceCol);
    result.push(sourceCol);
  }
  result.sort((a, b) => a - b);
  return result.length >= allColumns.length ? [] : result;
}

function restoreHiddenColumnsState() {
  const key = hiddenColumnsKey();
  if (!documentId || !key) {
    delete document.body.dataset.hiddenColumnsRestored;
    return new Set();
  }
  const entry = hiddenColumnStates().find((item) => item?.path === key);
  const columnsList = normalizedHiddenColumnList(entry?.columns);
  if (!columnsList.length) {
    delete document.body.dataset.hiddenColumnsRestored;
    return new Set();
  }
  document.body.dataset.hiddenColumnsRestored = columnsList.join(",");
  return new Set(columnsList);
}

function saveHiddenColumnsState() {
  const key = hiddenColumnsKey();
  if (!documentId || !key) return;
  const columnsList = normalizedHiddenColumnList(Array.from(hiddenSourceColumns));
  const rest = hiddenColumnStates().filter((item) => item?.path && item.path !== key);
  const next = columnsList.length
    ? [{ path: key, columns: columnsList, updatedAt: Date.now() }, ...rest].slice(0, hiddenColumnsLimit)
    : rest.slice(0, hiddenColumnsLimit);
  localStorage.setItem(hiddenColumnsStorageKey, JSON.stringify(next));
  delete document.body.dataset.hiddenColumnsRestored;
  document.body.dataset.hiddenColumnsPersisted = columnsList.join(",");
  document.body.dataset.hiddenColumnsKey = key;
}

function applyVerticalWheelDelta(deltaY) {
  const maxY = maxVerticalScroll();
  if (maxY === 0) {
    scrollY = 0;
    wheelRemainderY = 0;
    return;
  }

  wheelRemainderY += deltaY;
  const rowDelta = Math.trunc(wheelRemainderY / rowHeight);
  if (rowDelta === 0) return;

  scrollY = Math.max(0, Math.min(maxY, scrollY + rowDelta * rowHeight));
  wheelRemainderY -= rowDelta * rowHeight;
  if (scrollY === 0 || scrollY === maxY) {
    wheelRemainderY = 0;
  }
}

function initializeTheme() {
  applyTheme(themePreference, false);
  themeMediaQuery.addEventListener("change", () => {
    if (themePreference === "system") applyTheme("system", false);
  });
}

function normalizeLocalePreference(value) {
  return supportedLocales.some((locale) => locale.code === value) ? value : defaultLocale;
}

function localeLanguageTag(value = localePreference) {
  return value === "zh-cn" ? "zh-CN" : value;
}

function localeLabel(value = localePreference) {
  return supportedLocales.find((locale) => locale.code === value)?.label || supportedLocales[0].label;
}

function t(key) {
  return localeCatalog[localePreference]?.[key] || localeCatalog[defaultLocale][key] || key;
}

function tf(key, values = {}) {
  return t(key).replace(/\{\{(\w+)\}\}/g, (_match, name) => values[name] ?? "");
}

function translationOrEmpty(key) {
  const value = t(key);
  return value === key ? "" : value;
}

function localizedCommandTitle(command) {
  if (!command) return "";
  return translationOrEmpty(`command.${command.id}.title`) || command.title || command.id || "";
}

function localizedCommandDescription(command) {
  if (!command) return "";
  return translationOrEmpty(`command.${command.id}.description`) || command.description || "";
}

function localizedCommandLabel(commandID, fallback = commandID) {
  const command = commandsById.get(commandID);
  return command ? localizedCommandTitle(command) : fallback;
}

function localizedTransformStatus(key, label, values = {}) {
  return tf(`text.status.${key}`, { label, ...values });
}

function setLocalePreference(locale, persist = true) {
  localePreference = normalizeLocalePreference(locale);
  if (persist) {
    localStorage.setItem(localeStorageKey, localePreference);
  }
  document.documentElement.lang = localeLanguageTag();
  document.body.dataset.locale = localePreference;
  syncLocaleControls();
  if (commandsById.size) {
    applyCommandMetadata();
  } else {
    applyLocaleStrings();
  }
  if (paletteOpen) renderCommandPalette();
  if (shortcutEditorOpen) renderKeyboardShortcuts();
  if (lastDiffResult) renderDiffResult(lastDiffResult);
  if (lastSQLResult) renderSQLResult(lastSQLResult);
  if (profileData) renderProfile(profileData);
  syncProfileSchemaControls();
  renderSQLHistory();
  renderFileList();
}

function localePreferenceModified() {
  return localePreference !== defaultLocale;
}

function syncLocaleControls() {
  if (!localeSelect) return;
  localeSelect.value = localePreference;
  settingsLocaleSummary.textContent = localeLabel();
  resetLocaleButton.disabled = !localePreferenceModified();
}

function updateLocaleFromSettings() {
  setLocalePreference(localeSelect.value);
  renderSettingsPanel();
  status.textContent = t("status.languageChanged");
}

function resetLocalePreference() {
  setLocalePreference(defaultLocale);
  renderSettingsPanel();
  status.textContent = t("status.languageReset");
}

function applyLocaleStrings() {
  const textTargets = [
    [gridTabButton, "nav.grid"],
    [sqlTabButton, "nav.sql"],
    [profileTabButton, "nav.profile"],
    [settingsTabButton, "nav.settings"],
    [shortcutsTabButton, "nav.shortcuts"],
    [paletteButton, "nav.command"],
    [openButton, "file.open"],
    [copyFormatButton, "edit.copy"],
    [filterButton, "view.filter"],
    [filterBuilderButton, "view.builder"],
    [sortButton, "view.sort"],
    [findButton, "edit.find"],
    [replacePreviewButton, "edit.preview"],
    [replaceButton, "edit.replace"],
    [replaceAllButton, "edit.all"],
    [diffButton, "diff.compare"],
    [readOnlyButton, "file.readOnly"],
    [exportButton, "file.export"],
    [saveButton, "file.save"],
    [saveAsButton, "file.saveAs"],
    [restoreBackupButton, "file.restore"],
    [splitEditorButton, "workbench.split"],
    [splitOrientationButton, "workbench.orientation"],
    [focusEditorGroupButton, "workbench.focus"],
    [moveEditorTabButton, "workbench.move"],
    [closeSplitEditorButton, "workbench.close"],
    [settingsPanel?.querySelector(".settings-header strong"), "settings.title"],
    [settingsLocaleTitle, "settings.language"],
    [localeSelectLabel, "settings.language"],
    [resetLocaleButton, "settings.resetLanguage"],
    [settingsPanel?.querySelector('[aria-label="Saved filters and sorts"] header strong'), "settings.savedViews"],
    [settingsPanel?.querySelector('[aria-label="Theme"] header strong'), "settings.theme"],
    [settingsPanel?.querySelector('[aria-label="External changes"] header strong'), "settings.externalChanges"],
    [settingsPanel?.querySelector('[aria-label="Updates"] header strong'), "settings.updates"],
    [settingsPanel?.querySelector('[aria-label="Link protocols"] header strong'), "settings.linkProtocols"],
    [settingsPanel?.querySelector('[aria-label="Editor appearance"] header strong'), "settings.editorAppearance"],
    [settingsPanel?.querySelector('[aria-label="File format presets"] header strong'), "settings.fileFormatPresets"],
  ];
  for (const [element, key] of textTargets) {
    if (element) element.textContent = t(key);
  }
  if (shortcutsTabButton) {
    shortcutsTabButton.title = t("settings.keyboardShortcuts");
    shortcutsTabButton.setAttribute("aria-label", t("settings.keyboardShortcuts"));
  }
  const copyHeaderLabel = copyHeadersToggle?.nextElementSibling;
  if (copyHeaderLabel) copyHeaderLabel.textContent = t("edit.headers");
  if (recentFilterSelect?.options?.[0]) recentFilterSelect.options[0].textContent = t("view.filters");
  if (recentFilterSelect) {
    recentFilterSelect.setAttribute("aria-label", t("view.filters.aria"));
    recentFilterSelect.title = t("view.filters.aria");
  }
  if (viewChips) viewChips.setAttribute("aria-label", t("view.chips.aria"));
  const stableSortLabel = columnStableSortToggle?.closest("label")?.querySelector("span");
  if (stableSortLabel) stableSortLabel.textContent = t("column.stableSort");
  const transformLabel = columnTransformMode?.closest("label")?.querySelector("span");
  if (transformLabel) transformLabel.textContent = t("column.textTransform");
  if (columnTransformMode) {
    columnTransformMode.setAttribute("aria-label", t("column.textTransform"));
    applyColumnTransformOptionLabels();
  }
  if (columnTransformApply) columnTransformApply.textContent = t("column.applyTransform");
  if (columnMenuValues) columnMenuValues.setAttribute("aria-label", t("column.values"));
  if (columnApplyValues) columnApplyValues.textContent = t("column.applyValues");
  if (columnHide) columnHide.textContent = t("column.hide");
  if (columnShowAll) columnShowAll.textContent = t("column.showAll");
  if (columnClearView) columnClearView.textContent = t("column.clearView");
  if (settingsSearch) settingsSearch.placeholder = t("settings.searchPlaceholder");
  applyDataI18n();
  syncAboutSupportPanel();
  applyFindLocaleStrings();
  applyConditionLocaleStrings();
  applyProfileLocaleStrings();
  applySQLLocaleStrings();
  applyDiffLocaleStrings();
  renderViewChips();
  const readyLabels = new Set(Object.values(localeCatalog).map((catalog) => catalog["status.ready"]).filter(Boolean));
  if (!documentId && status && readyLabels.has(status.textContent)) {
    status.textContent = t("status.ready");
  }
  document.body.dataset.localeLabel = localeLabel();
}

function applyColumnTransformOptionLabels() {
  if (!columnTransformMode) return;
  for (const option of columnTransformMode.options) {
    const label = translationOrEmpty(`text.mode.${option.value}`);
    if (label) option.textContent = label;
  }
}

function applyDataI18n(root = document) {
  root.querySelectorAll("[data-i18n]").forEach((element) => {
    element.textContent = t(element.dataset.i18n);
  });
  root.querySelectorAll("[data-i18n-placeholder]").forEach((element) => {
    element.placeholder = t(element.dataset.i18nPlaceholder);
  });
}

function applyFindLocaleStrings() {
  if (commandInput) commandInput.placeholder = t("find.placeholder");
  if (replaceInput) replaceInput.placeholder = t("find.replacePlaceholder");
  if (findButton) findButton.title = t("edit.find");
  if (cancelFindButton) cancelFindButton.title = t("common.cancel");
  if (findCaseButton) findCaseButton.title = t("find.matchCase");
  if (findWholeButton) findWholeButton.title = t("find.wholeCell");
  if (findRegexButton) findRegexButton.title = t("find.regex");
  if (findSelectionButton) {
    findSelectionButton.title = t("find.selection");
    findSelectionButton.textContent = t("find.selectionShort");
  }
  if (findPrevButton) findPrevButton.title = t("find.previous");
  if (findNextButton) findNextButton.title = t("find.next");
  if (replacePreviewButton) replacePreviewButton.title = t("find.previewTitle");
  if (replaceButton) replaceButton.title = t("find.replaceTitle");
  if (replaceAllButton) replaceAllButton.title = t("find.replaceAllTitle");
}

function conditionOperatorLabelKeys() {
  return {
    contains: "condition.operator.contains",
    notContains: "condition.operator.notContains",
    equals: "condition.operator.equals",
    notEquals: "condition.operator.notEquals",
    inList: "condition.operator.inList",
    notInList: "condition.operator.notInList",
    startsWith: "condition.operator.startsWith",
    endsWith: "condition.operator.endsWith",
    blank: "condition.operator.blank",
    notBlank: "condition.operator.notBlank",
    booleanTrue: "condition.operator.booleanTrue",
    booleanFalse: "condition.operator.booleanFalse",
    dateEquals: "condition.operator.dateEquals",
    dateBefore: "condition.operator.dateBefore",
    dateAfter: "condition.operator.dateAfter",
    dateOnOrBefore: "condition.operator.dateOnOrBefore",
    dateOnOrAfter: "condition.operator.dateOnOrAfter",
    gt: "condition.operator.gt",
    gte: "condition.operator.gte",
    lt: "condition.operator.lt",
    lte: "condition.operator.lte",
    between: "condition.operator.between",
    notBetween: "condition.operator.notBetween",
  };
}

function conditionTemplateLabelKeys() {
  return {
    "": "condition.template.choose",
    selectedEquals: "condition.template.selectedEquals",
    selectedContains: "condition.template.selectedContains",
    selectedNotContains: "condition.template.selectedNotContains",
    selectedNotEquals: "condition.template.selectedNotEquals",
    activeBlank: "condition.template.activeBlank",
    activeNotBlank: "condition.template.activeNotBlank",
    activeTrue: "condition.template.activeTrue",
    activeFalse: "condition.template.activeFalse",
  };
}

function applyConditionLocaleStrings() {
  const dialog = conditionPanel?.querySelector(".condition-panel");
  const templateBar = conditionPanel?.querySelector(".condition-template-bar");
  if (filterBuilderButton) filterBuilderButton.title = t("condition.title");
  if (dialog) dialog.setAttribute("aria-label", t("condition.title"));
  if (templateBar) templateBar.setAttribute("aria-label", t("condition.quickConditions"));
  if (conditionTemplateSelect) conditionTemplateSelect.setAttribute("aria-label", t("condition.quickConditions"));
  if (conditionPrimaryOpenGroup) {
    conditionPrimaryOpenGroup.setAttribute("aria-label", tf("condition.group.openAria", { number: 1 }));
    const label = conditionPrimaryOpenGroup.closest("label");
    if (label) label.title = t("condition.group.openTitle");
  }
  if (conditionPrimaryCloseGroup) {
    conditionPrimaryCloseGroup.setAttribute("aria-label", tf("condition.group.closeAria", { number: 1 }));
    const label = conditionPrimaryCloseGroup.closest("label");
    if (label) label.title = t("condition.group.closeTitle");
  }
  if (conditionColumn) conditionColumn.setAttribute("aria-label", tf("condition.columnAria", { number: 1 }));
  if (conditionOperator) conditionOperator.setAttribute("aria-label", tf("condition.operatorAria", { number: 1 }));
  if (conditionValue) conditionValue.setAttribute("aria-label", tf("condition.valueAria", { number: 1 }));
  if (conditionRows) conditionRows.setAttribute("aria-label", t("condition.additional"));
  setSelectOptionLabels(conditionTemplateSelect, conditionTemplateLabelKeys());
  setSelectOptionLabels(conditionOperator, conditionOperatorLabelKeys());
  if (conditionRows?.querySelector(".condition-row-extra")) {
    readConditionExtraRowsFromDOM();
    renderConditionExtraRows();
  } else {
    updateConditionValueState();
  }
}

function applyProfileLocaleStrings() {
  setSelectOptionLabels(profileRuleTypeSelect, {
    "": "profile.type.any",
    string: "profile.type.string",
    number: "profile.type.number",
    integer: "profile.type.integer",
    decimal: "profile.type.decimal",
    boolean: "profile.type.boolean",
    date: "profile.type.date",
    datetime: "profile.type.datetime",
    time: "profile.type.time",
    year: "profile.type.year",
    yearmonth: "profile.type.yearmonth",
  });
  if (savedProfileSchemaSelect?.options?.[0]) {
    savedProfileSchemaSelect.options[0].textContent = t("profile.savedSchemas");
  }
  if (editingProfileRuleIndex >= 0) {
    saveProfileRuleButton.textContent = t("profile.rule.update");
  }
}

function applySQLLocaleStrings() {
  const textTargets = [
    [runSqlButton, "sql.run"],
    [cancelSqlButton, "common.cancel"],
    [clearSqlHistoryButton, "sql.clear"],
  ];
  for (const [element, key] of textTargets) {
    if (element) element.textContent = t(key);
  }
  if (cancelSqlButton) cancelSqlButton.title = t("common.cancel");
  syncSQLResultActions();
}

function applyDiffLocaleStrings() {
  const textTargets = [
    [runDiffButton, "diff.run"],
    [diffIgnoreCaseLabel, "diff.ignoreCase"],
    [diffIgnoreWhitespaceLabel, "diff.ignoreWhitespace"],
    [diffMatchHeadersLabel, "diff.matchHeaders"],
    [diffClearFilterButton, "diff.clear"],
    [diffRightTextModeLabel, "diff.rightTextMode"],
    [diffPasteClipboardButton, "diff.pasteClipboard"],
  ];
  for (const [element, key] of textTargets) {
    if (element) element.textContent = t(key);
  }
  setDiffPlaceholder(diffKeyColumn, "diff.keyColumnPlaceholder");
  setDiffPlaceholder(diffIgnoredColumns, "diff.ignoredColumnsPlaceholder");
  setDiffPlaceholder(diffSearchInput, "diff.searchPlaceholder");
  setDiffPlaceholder(diffRightText, "diff.rightTextPlaceholder");
  if (diffPrevButton) diffPrevButton.title = t("diff.previousDifference");
  if (diffNextButton) diffNextButton.title = t("diff.nextDifference");
  setSelectOptionLabels(diffStatusFilter, {
    all: "diff.status.all",
    differences: "diff.status.differences",
    changed: "diff.status.changed",
    added: "diff.status.added",
    deleted: "diff.status.deleted",
    unchanged: "diff.status.unchanged",
  });
  setSelectOptionLabels(diffSortSelect, {
    source: "diff.sort.source",
    status: "diff.sort.status",
    leftRow: "diff.sort.leftRow",
    rightRow: "diff.sort.rightRow",
    leftText: "diff.sort.leftText",
    rightText: "diff.sort.rightText",
  });
  setDiffSortDirection(diffSortDirection());
}

function setDiffPlaceholder(element, key) {
  if (element) element.placeholder = t(key);
}

function setSelectOptionLabels(select, labels) {
  if (!select) return;
  for (const option of select.options) {
    const key = labels[option.value];
    if (key) option.textContent = t(key);
  }
}

function normalizeThemePreference(preference) {
  return themePreferences.includes(preference) ? preference : defaultThemePreference;
}

function themePreferenceLabel(preference = themePreference) {
  const labels = {
    light: "Light",
    dark: "Dark",
    "high-contrast": "High contrast",
    system: "System",
  };
  return labels[preference] || labels[defaultThemePreference];
}

function applyTheme(preference, announce = true) {
  preference = normalizeThemePreference(preference);
  themePreference = preference;
  localStorage.setItem(themeStorageKey, preference);
  const resolved = preference === "system"
    ? (themeMediaQuery.matches ? "dark" : "light")
    : preference;
  document.body.dataset.theme = resolved;
  document.body.dataset.themePreference = preference;
  syncThemeControls();
  refreshThemeColors();
  draw();
  if (announce && status) {
    status.textContent = `Theme: ${preference === "system" ? `System (${themePreferenceLabel(resolved)})` : themePreferenceLabel(preference)}`;
  }
}

function cycleTheme() {
  const order = themePreferences;
  const current = order.includes(themePreference) ? order.indexOf(themePreference) : order.indexOf(defaultThemePreference);
  applyTheme(order[(current + 1) % order.length]);
}

function syncThemeControls() {
  if (!themePreferenceSelect) return;
  themePreferenceSelect.value = themePreference;
  const resolved = document.body.dataset.theme || defaultThemePreference;
  settingsThemeSummary.textContent = themePreference === "system"
    ? `System (${themePreferenceLabel(resolved)})`
    : themePreferenceLabel(themePreference);
  resetThemePreferenceButton.disabled = !themePreferenceModified();
}

function themePreferenceModified() {
  return themePreference !== defaultThemePreference;
}

function updateThemeFromSettings() {
  applyTheme(themePreferenceSelect.value);
  renderSettingsPanel();
}

function resetThemePreference() {
  applyTheme(defaultThemePreference);
  renderSettingsPanel();
}

function normalizeAutoRefreshMode(value) {
  return autoRefreshModes.includes(value) ? value : defaultAutoRefreshMode;
}

function autoRefreshModeLabel(mode = autoRefreshMode) {
  return mode === "notify" ? "Notify only" : "Auto-refresh clean files";
}

function setAutoRefreshMode(mode, persist = true) {
  autoRefreshMode = normalizeAutoRefreshMode(mode);
  if (persist) {
    localStorage.setItem(autoRefreshStorageKey, autoRefreshMode);
  }
  document.body.dataset.autoRefreshMode = autoRefreshMode;
  syncAutoRefreshControls();
}

function syncAutoRefreshControls() {
  if (!autoRefreshModeSelect) return;
  autoRefreshModeSelect.value = autoRefreshMode;
  settingsRefreshSummary.textContent = autoRefreshModeLabel();
  resetAutoRefreshButton.disabled = !autoRefreshModeModified();
}

function autoRefreshModeModified() {
  return autoRefreshMode !== defaultAutoRefreshMode;
}

function updateAutoRefreshModeFromSettings() {
  setAutoRefreshMode(autoRefreshModeSelect.value);
  renderSettingsPanel();
  status.textContent = autoRefreshMode === "clean"
    ? "Clean documents will auto-refresh when files change."
    : "External file changes will be reported without auto-refresh.";
}

function resetAutoRefreshMode() {
  setAutoRefreshMode(defaultAutoRefreshMode);
  renderSettingsPanel();
  status.textContent = "External change refresh reset.";
}

function normalizeUpdateCheckMode(value) {
  return updateCheckModes.includes(value) ? value : defaultUpdateCheckMode;
}

function updateCheckModeLabel(mode = updateCheckMode) {
  const labels = {
    automatic: "Automatic checks",
    manual: "Manual only",
    off: "Off",
  };
  return labels[mode] || labels[defaultUpdateCheckMode];
}

function loadUpdateCheckState() {
  try {
    return JSON.parse(localStorage.getItem(updateCheckResultStorageKey) || "{}");
  } catch {
    return {};
  }
}

function saveUpdateCheckState(state) {
  updateCheckState = { ...state };
  try {
    localStorage.setItem(updateCheckResultStorageKey, JSON.stringify(updateCheckState));
  } catch {
    status.textContent = "Update check result could not be saved.";
  }
}

function setUpdateCheckMode(mode, persist = true) {
  updateCheckMode = normalizeUpdateCheckMode(mode);
  if (persist) {
    localStorage.setItem(updateCheckStorageKey, updateCheckMode);
  }
  document.body.dataset.updateCheckMode = updateCheckMode;
  syncUpdateCheckControls();
}

function updateCheckModeModified() {
  return updateCheckMode !== defaultUpdateCheckMode;
}

function syncUpdateCheckControls() {
  if (!updateCheckModeSelect) return;
  updateCheckModeSelect.value = updateCheckMode;
  const last = updateCheckState?.checkedAt ? ` · last ${new Date(updateCheckState.checkedAt).toLocaleString()}` : "";
  settingsUpdateSummary.textContent = `${updateCheckModeLabel()}${last}`;
  const checksDisabled = updateCheckMode === "off";
  checkForUpdatesButton.disabled = checksDisabled;
  checkForUpdatesButton.setAttribute("aria-disabled", String(checksDisabled));
  resetUpdateChecksButton.disabled = !updateCheckModeModified() && !updateCheckState?.checkedAt;
  if (updateCheckResult) {
    updateCheckResult.textContent = updateCheckState?.message || "No update check has run.";
  }
  document.body.dataset.updateCheckSummary = settingsUpdateSummary.textContent;
  document.body.dataset.updateCheckLastSource = updateCheckState?.source || "";
  document.body.dataset.updateCheckLastStatus = updateCheckState?.status || "";
  document.body.dataset.updateCheckCount = String(Number(document.body.dataset.updateCheckCount || "0"));
  syncAboutSupportPanel();
}

function syncAboutSupportPanel() {
  if (!aboutSupportPanel) return;
  if (aboutUpdateSummary) aboutUpdateSummary.textContent = settingsUpdateSummary?.textContent || updateCheckModeLabel();
  if (aboutUpdateResult) aboutUpdateResult.textContent = updateCheckState?.message || "No update check has run.";
  if (aboutCheckUpdatesButton) {
    const disabled = updateCheckMode === "off";
    aboutCheckUpdatesButton.disabled = disabled;
    aboutCheckUpdatesButton.setAttribute("aria-disabled", String(disabled));
  }
  document.body.dataset.aboutSupportStatus = `${updateCheckMode}:${updateCheckState?.status || "none"}`;
}

function openAboutSupportPanel() {
  syncAboutSupportPanel();
  aboutSupportPanel.hidden = false;
  document.body.dataset.aboutSupportPanel = "open";
  closeAboutSupportButton?.focus();
}

function closeAboutSupportPanel() {
  aboutSupportPanel.hidden = true;
  document.body.dataset.aboutSupportPanel = "closed";
  canvas.focus();
}

async function runAboutSupportUpdateCheck() {
  await checkForUpdates("manual");
  syncAboutSupportPanel();
}

function updateUpdateCheckModeFromSettings() {
  setUpdateCheckMode(updateCheckModeSelect.value);
  renderSettingsPanel();
  if (updateCheckMode === "automatic") {
    status.textContent = "Startup update checks are enabled.";
  } else if (updateCheckMode === "manual") {
    status.textContent = "Update checks run manually only.";
  } else {
    status.textContent = "Update checks are off.";
  }
}

function resetUpdateCheckMode() {
  setUpdateCheckMode(defaultUpdateCheckMode);
  updateCheckState = {};
  localStorage.removeItem(updateCheckResultStorageKey);
  syncUpdateCheckControls();
  renderSettingsPanel();
  status.textContent = "Update checks reset.";
}

async function checkForUpdates(source = "manual") {
  if (updateCheckMode === "off") {
    document.body.dataset.updateCheckSkipped = `${source}:off`;
    status.textContent = "Update checks are off.";
    syncUpdateCheckControls();
    return;
  }
  document.body.dataset.updateCheckInFlight = source;
  try {
    const response = await fetch(`/api/update-check?source=${encodeURIComponent(source)}`);
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "update check failed");
    const nextState = {
      status: payload.status || "unknown",
      source: payload.source || source,
      currentVersion: payload.currentVersion || "",
      latestVersion: payload.latestVersion || "",
      checkedAt: payload.checkedAt || new Date().toISOString(),
      message: payload.message || "Update check completed.",
    };
    saveUpdateCheckState(nextState);
    document.body.dataset.updateCheckCount = String(Number(document.body.dataset.updateCheckCount || "0") + 1);
    document.body.dataset.updateCheckSkipped = "";
    if (source !== "startup") {
      status.textContent = "Update check completed.";
    }
  } catch (error) {
    const nextState = {
      status: "error",
      source,
      checkedAt: new Date().toISOString(),
      message: `Update check failed: ${error.message}`,
    };
    saveUpdateCheckState(nextState);
    if (source !== "startup") {
      status.textContent = nextState.message;
    }
  } finally {
    document.body.dataset.updateCheckInFlight = "";
    syncUpdateCheckControls();
    if (viewMode === "settings") renderSettingsPanel();
  }
}

function maybeRunStartupUpdateCheck() {
  if (updateCheckMode !== "automatic") {
    document.body.dataset.updateStartupCheck = `skipped:${updateCheckMode}`;
    document.body.dataset.updateCheckSkipped = `startup:${updateCheckMode}`;
    syncUpdateCheckControls();
    return;
  }
  document.body.dataset.updateStartupCheck = "running";
  checkForUpdates("startup").then(() => {
    document.body.dataset.updateStartupCheck = "completed";
  });
}

function normalizeFilterHistoryLimit(value) {
  const number = Number(value);
  return filterHistoryLimitOptions.includes(number) ? number : defaultFilterHistoryLimit;
}

function filterHistoryLimitLabel(limit = filterHistoryLimit) {
  if (limit === 0) return "Off";
  return `Keep ${limit}`;
}

function setFilterHistoryLimit(limit, persist = true) {
  filterHistoryLimit = normalizeFilterHistoryLimit(limit);
  if (persist) {
    localStorage.setItem(filterHistoryLimitStorageKey, String(filterHistoryLimit));
  }
  if (filterHistoryLimit === 0) {
    saveRecentFilters([]);
  } else {
    saveRecentFilters(recentFilters());
  }
  syncFilterHistoryControls();
  refreshFilterDropdown();
}

function syncFilterHistoryControls() {
  if (!filterHistoryLimitSelect) return;
  filterHistoryLimitSelect.value = String(filterHistoryLimit);
  clearRecentFiltersButton.disabled = recentFilters().length === 0;
  document.body.dataset.filterHistoryLimit = String(filterHistoryLimit);
}

function filterHistoryLimitModified() {
  return filterHistoryLimit !== defaultFilterHistoryLimit;
}

function savedViewsSettingsModified({ filters = savedFilters(), recent = recentFilters(), sorts = savedSorts() } = {}) {
  return filters.length > 0 || recent.length > 0 || sorts.length > 0 || filterHistoryLimitModified();
}

function updateFilterHistoryLimitFromSettings() {
  setFilterHistoryLimit(filterHistoryLimitSelect.value);
  renderSettingsPanel();
  status.textContent = filterHistoryLimit === 0
    ? t("view.status.filterHistoryOff")
    : tf("view.status.filterHistoryKeep", { count: formatCount(filterHistoryLimit) });
}

function clearRecentFilters() {
  saveRecentFilters([]);
  refreshFilterDropdown();
  renderSettingsPanel();
  status.textContent = t("view.status.clearedRecentFilters");
}

function resetSavedViewsSettings() {
  const before = {
    filters: savedFilters().length,
    recent: recentFilters().length,
    sorts: savedSorts().length,
    history: filterHistoryLimit,
  };
  saveSavedFilters([]);
  saveSavedSorts([]);
  filterHistoryLimit = defaultFilterHistoryLimit;
  localStorage.setItem(filterHistoryLimitStorageKey, String(filterHistoryLimit));
  saveRecentFilters([]);
  if (savedViewImportPanel) savedViewImportPanel.hidden = true;
  if (savedViewImportText) savedViewImportText.value = "";
  document.body.dataset.savedViewImportPanel = "closed";
  document.body.dataset.savedViewsReset = `${before.filters}:${before.recent}:${before.sorts}:${before.history}->0:0:0:${filterHistoryLimit}`;
  refreshFilterDropdown([], []);
  renderViewChips();
  renderSettingsPanel();
  status.textContent = "Saved filters, recent filters, and saved sorts reset.";
}

function loadCustomLinkProtocols() {
  try {
    const stored = localStorage.getItem(customLinkProtocolsStorageKey);
    if (stored == null) return [...defaultCustomLinkProtocols];
    const parsed = JSON.parse(stored);
    const normalized = normalizeCustomLinkProtocols(parsed);
    return normalized.error ? [...defaultCustomLinkProtocols] : normalized.protocols;
  } catch {
    return [...defaultCustomLinkProtocols];
  }
}

function normalizeCustomLinkProtocols(value) {
  const rawItems = Array.isArray(value)
    ? value
    : String(value ?? "").split(/[,\s]+/);
  const protocols = [];
  const seen = new Set();
  for (const item of rawItems) {
    const protocol = normalizeCustomLinkProtocol(item);
    if (!protocol) continue;
    if (blockedLinkProtocols.has(`${protocol}:`)) {
      return { protocols: [], error: `${protocol} is unsafe and cannot be allowed.` };
    }
    if (trustedLinkProtocols.has(`${protocol}:`)) continue;
    if (!seen.has(protocol)) {
      seen.add(protocol);
      protocols.push(protocol);
    }
  }
  return { protocols: protocols.slice(0, 24), error: "" };
}

function normalizeCustomLinkProtocol(value) {
  const protocol = String(value ?? "").trim().toLowerCase().replace(/:$/, "");
  if (!protocol) return "";
  return /^[a-z][a-z0-9+.-]*$/.test(protocol) ? protocol : "";
}

function customLinkProtocolLabel() {
  return customLinkProtocols.length ? customLinkProtocols.join(", ") : "None";
}

function customLinkProtocolsModified() {
  return JSON.stringify(customLinkProtocols) !== JSON.stringify(defaultCustomLinkProtocols);
}

function syncCustomLinkProtocolControls() {
  if (!customLinkProtocolsInput) return;
  customLinkProtocolsInput.value = customLinkProtocols.join(", ");
  settingsLinkSummary.textContent = customLinkProtocolLabel();
  resetCustomLinkProtocolsButton.disabled = !customLinkProtocolsModified();
  document.body.dataset.customLinkProtocols = customLinkProtocols.join(",");
  document.body.dataset.confirmedCustomLinkProtocols = confirmedCustomLinkProtocols.join(",");
}

function applyCustomLinkProtocols(nextProtocols, announce = true) {
  customLinkProtocols = [...nextProtocols];
  localStorage.setItem(customLinkProtocolsStorageKey, JSON.stringify(customLinkProtocols));
  persistConfirmedCustomLinkProtocols();
  syncCustomLinkProtocolControls();
  syncCellLinkState(cellValueInput.value);
  if (announce) {
    status.textContent = `Allowed link protocols: ${customLinkProtocolLabel()}`;
  }
}

function updateCustomLinkProtocolsFromSettings() {
  const normalized = normalizeCustomLinkProtocols(customLinkProtocolsInput.value);
  if (normalized.error) {
    status.textContent = normalized.error;
    return;
  }
  applyCustomLinkProtocols(normalized.protocols);
  renderSettingsPanel();
}

function resetCustomLinkProtocols() {
  applyCustomLinkProtocols(defaultCustomLinkProtocols);
  renderSettingsPanel();
  status.textContent = "Link protocols reset.";
}

function loadConfirmedCustomLinkProtocols() {
  try {
    return normalizeConfirmedCustomLinkProtocols(JSON.parse(localStorage.getItem(confirmedCustomLinkProtocolsStorageKey) || "[]"));
  } catch {
    return [];
  }
}

function normalizeConfirmedCustomLinkProtocols(value) {
  const normalized = normalizeCustomLinkProtocols(value);
  if (normalized.error) return [];
  return normalized.protocols.filter((protocol) => customLinkProtocols.includes(protocol));
}

function persistConfirmedCustomLinkProtocols() {
  confirmedCustomLinkProtocols = normalizeConfirmedCustomLinkProtocols(confirmedCustomLinkProtocols);
  localStorage.setItem(confirmedCustomLinkProtocolsStorageKey, JSON.stringify(confirmedCustomLinkProtocols));
  document.body.dataset.confirmedCustomLinkProtocols = confirmedCustomLinkProtocols.join(",");
}

function rememberCustomLinkProtocol(protocol) {
  const normalized = normalizeCustomLinkProtocol(protocol);
  if (!normalized || !customLinkProtocols.includes(normalized)) return;
  confirmedCustomLinkProtocols = [
    normalized,
    ...confirmedCustomLinkProtocols.filter((item) => item !== normalized),
  ];
  persistConfirmedCustomLinkProtocols();
}

function loadEditorAppearance() {
  try {
    return normalizeEditorAppearance(JSON.parse(localStorage.getItem(editorAppearanceStorageKey) || "{}"));
  } catch {
    return { ...defaultEditorAppearance };
  }
}

function normalizeEditorAppearance(value) {
  const fontValue = (candidate, fallback) => (editorFontFamilies[candidate] ? candidate : fallback);
  return {
    gridFont: fontValue(value?.gridFont, defaultEditorAppearance.gridFont),
    gridSize: clampFontSize(value?.gridSize, defaultEditorAppearance.gridSize),
    gridRowHeight: clampGridRowHeight(value?.gridRowHeight, defaultEditorAppearance.gridRowHeight),
    editorFont: fontValue(value?.editorFont, defaultEditorAppearance.editorFont),
    editorSize: clampFontSize(value?.editorSize, defaultEditorAppearance.editorSize),
  };
}

function clampFontSize(value, fallback) {
  const parsed = Number.parseInt(String(value ?? ""), 10);
  return Number.isFinite(parsed) ? Math.max(11, Math.min(18, parsed)) : fallback;
}

function clampGridRowHeight(value, fallback) {
  const parsed = Number.parseInt(String(value ?? ""), 10);
  return Number.isFinite(parsed) ? Math.max(22, Math.min(44, parsed)) : fallback;
}

function applyEditorAppearance(next, persist = true) {
  editorAppearance = normalizeEditorAppearance(next);
  if (persist) {
    localStorage.setItem(editorAppearanceStorageKey, JSON.stringify(editorAppearance));
  }
  const root = document.documentElement;
  root.style.setProperty("--grid-font-family", editorFontFamilies[editorAppearance.gridFont]);
  root.style.setProperty("--grid-font-size", `${editorAppearance.gridSize}px`);
  root.style.setProperty("--editor-font-family", editorFontFamilies[editorAppearance.editorFont]);
  root.style.setProperty("--editor-font-size", `${editorAppearance.editorSize}px`);
  rowHeight = editorAppearance.gridRowHeight;
  scrollY = Math.max(0, Math.min(maxVerticalScroll(), scrollY));
  document.body.dataset.editorAppearance = `${editorAppearance.gridFont}:${editorAppearance.gridSize}:row${editorAppearance.gridRowHeight}/${editorAppearance.editorFont}:${editorAppearance.editorSize}`;
  syncEditorAppearanceControls();
  draw();
}

function syncEditorAppearanceControls() {
  if (!gridFontSelect) return;
  gridFontSelect.value = editorAppearance.gridFont;
  gridFontSizeInput.value = String(editorAppearance.gridSize);
  gridRowHeightSelect.value = String(editorAppearance.gridRowHeight);
  textEditorFontSelect.value = editorAppearance.editorFont;
  textEditorFontSizeInput.value = String(editorAppearance.editorSize);
  settingsAppearanceSummary.textContent = `Grid ${editorAppearance.gridSize}px · Row ${editorAppearance.gridRowHeight}px · Editor ${editorAppearance.editorSize}px`;
  resetEditorAppearanceButton.disabled = !editorAppearanceModified();
}

function editorAppearanceModified() {
  return JSON.stringify(editorAppearance) !== JSON.stringify(defaultEditorAppearance);
}

function updateEditorAppearanceFromControls() {
  applyEditorAppearance({
    gridFont: gridFontSelect.value,
    gridSize: gridFontSizeInput.value,
    gridRowHeight: gridRowHeightSelect.value,
    editorFont: textEditorFontSelect.value,
    editorSize: textEditorFontSizeInput.value,
  });
  renderSettingsPanel();
  status.textContent = "Editor appearance updated.";
}

function resetEditorAppearance() {
  applyEditorAppearance(defaultEditorAppearance);
  renderSettingsPanel();
  status.textContent = "Editor appearance reset.";
}

function gridBodyFont() {
  return `${editorAppearance.gridSize}px ${editorFontFamilies[editorAppearance.gridFont]}`;
}

function refreshThemeColors() {
  const style = getComputedStyle(document.body);
  const read = (name) => style.getPropertyValue(name).trim();
  themeColors = {
    accent: read("--accent"),
    gridBg: read("--grid-bg"),
    gridHeader: read("--grid-header"),
    gridGutter: read("--grid-gutter"),
    gridRow: read("--grid-row"),
    gridRowAlt: read("--grid-row-alt"),
    gridLine: read("--grid-line"),
    gridLineSoft: read("--grid-line-soft"),
    findActiveHighlight: read("--find-active-highlight"),
    findHighlight: read("--find-highlight"),
    muted: read("--muted"),
    panelSubtle: read("--panel-subtle"),
    quoteMarker: read("--quote-marker"),
    quoteMarkerText: read("--quote-marker-text"),
    selectionBg: read("--selection-bg"),
    text: read("--text"),
  };
}

function cssPixelValue(name, fallback) {
  const value = getComputedStyle(document.documentElement).getPropertyValue(name);
  const parsed = Number.parseFloat(value);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
}

function valueAt(row, col) {
  if (!documentId) return "";
  const aggregateRow = aggregateRowAt(row);
  if (aggregateRow) return aggregateRow.values[sourceColFor(col)] ?? "";
  const record = rows.get(row);
  if (!record) {
    maybeLoadWindow(row);
    return "";
  }
  return record[sourceColFor(col)] ?? "";
}

function valueAtSourceColumn(row, sourceCol) {
  if (!documentId) return "";
  const aggregateRow = aggregateRowAt(row);
  if (aggregateRow) return aggregateRow.values[sourceCol] ?? "";
  const record = rows.get(row);
  if (!record) {
    maybeLoadWindow(row);
    return "";
  }
  return record[sourceCol] ?? "";
}

function hasCellAt(row, col) {
  if (aggregateRowAt(row)) return true;
  const record = rows.get(row);
  if (!record) return true;
  return sourceColFor(col) < record.length;
}

function aggregateRowAt(row) {
  if (row < sourceRowCount) return null;
  return aggregateRows[row - sourceRowCount] ?? null;
}

function aggregateRowIndex(row) {
  return row >= sourceRowCount ? row - sourceRowCount : -1;
}

function syncTotalRows() {
  totalRows = sourceRowCount + aggregateRows.length;
  document.body.dataset.aggregateRows = String(aggregateRows.length);
}

function sourceColFor(col) {
  return columns[col]?.sourceCol ?? col;
}

function visibleColForSource(sourceCol) {
  return columns.findIndex((column) => column.sourceCol === sourceCol);
}

function columnNameForSource(sourceCol) {
  return allColumns.find((column) => column.sourceCol === sourceCol)?.name ?? `Column ${sourceCol + 1}`;
}

function syncHiddenColumnControls() {
  const hidden = Array.from(hiddenSourceColumns).sort((a, b) => a - b);
  const count = hidden.length;
  const label = count === 1 ? "1 hidden" : `${count} hidden`;
  hiddenColumnSummary.textContent = documentId ? label : "0 hidden";
  hiddenColumnSummary.title = hidden.map(columnNameForSource).join(", ");
  showHiddenColumnsButton.disabled = !documentId || count === 0;
  columnShowAll.disabled = !documentId || count === 0;
  document.body.dataset.hiddenColumns = hidden.join(",");
  document.body.dataset.hiddenColumnCount = documentId ? String(count) : "0";
  document.body.dataset.hiddenColumnsKey = documentId ? hiddenColumnsKey() : "";
  document.body.dataset.hiddenColumnsPersisted = hidden.join(",");
}

function setStructureClipboard(kind, sources, cut = false) {
  const normalizedSources = (Array.isArray(sources) ? sources : [sources])
    .filter((source) => Number.isInteger(source));
  structureClipboard = { documentId, kind, source: normalizedSources[0], sources: normalizedSources, cut };
  document.body.dataset.structureClipboard = `${kind}:${cut ? "cut" : "copy"}:${normalizedSources.join(",")}`;
  updateCommandState();
}

function clearStructureClipboard() {
  structureClipboard = null;
  document.body.dataset.structureClipboard = "";
  updateCommandState();
}

function structureClipboardFor(kind) {
  if (!structureClipboard || structureClipboard.documentId !== documentId || structureClipboard.kind !== kind) {
    return null;
  }
  return structureClipboard;
}

function applyHiddenColumns() {
  columns = allColumns.filter((column) => !hiddenSourceColumns.has(column.sourceCol));
  syncHiddenColumnControls();
}

function draw() {
  if (viewMode !== "grid") {
    updateCommandState();
    syncAppTitle();
    return;
  }
  const { width, height } = syncCanvasSize();
  ctx.clearRect(0, 0, width, height);
  headerCtx.clearRect(0, 0, width, headerHeight);

  ctx.fillStyle = themeColors.gridBg;
  ctx.fillRect(0, 0, width, height);

  drawRows(width, height);
  drawBodySeparators(width, height);
  drawSelection();
  drawRowReorderPreview(width, height);
  drawHeaders(headerCtx, width);

  const selectedColumn = columns[selected.col]?.name ?? "";
  const path = metadata?.path ? ` - ${metadata.path}` : "";
  const dirtyLabel = dirty ? " - edited" : "";
  const aggregateStats = currentSelectionAggregateStats();
  const aggregate = selectionAggregateSummary(aggregateStats);
  const modeLabel = readOnly ? " - read-only" : "";
  const externalLabel = externalFileChanged ? " - external change" : "";
  const aggregateLabel = aggregateRows.length ? ` + ${aggregateRows.length} aggregate` : "";
  document.body.dataset.selectionRange = documentId ? selectionSummary() : "";
  document.body.dataset.columnWidths = columns.map((column) => Math.round(column.width)).join(",");
  document.body.dataset.gridHeaderHeight = String(headerHeight);
  document.body.dataset.gridHeaderLayer = `${headerCanvas.clientHeight}:${headerCanvas.height}`;
  queueScrollPositionSave();
  copyHeadersToggle.checked = includeHeadersInSelection;
  updateAggregateStatus(aggregateStats);
  status.textContent =
    documentId
      ? `${sourceRowCount.toLocaleString()} rows${aggregateLabel}, ${columns.length} columns${path}${dirtyLabel}${modeLabel}${externalLabel} - row ${selected.row + 1}, ${selectedColumn}${aggregate}`
      : "Ready";
  renderEditorGroups();
  syncAppTitle();
  updateFormatStatus();
  updateInspector();
  updateCommandState();
}

function renderEditorGroups() {
  if (!editorGroups) return;
  if (!editorGroupsSplit) {
    activeEditorGroup = "primary";
    currentEditorTabGroup = "primary";
  }
  syncEditorGroupActiveTabs();
  mountCurrentEditorTab();
  editorGroups.classList.toggle("split", editorGroupsSplit);
  editorGroups.dataset.orientation = editorGroupOrientation;
  secondaryEditorGroup.hidden = !editorGroupsSplit;
  primaryEditorGroup.classList.toggle("active", activeEditorGroup === "primary");
  secondaryEditorGroup.classList.toggle("active", editorGroupsSplit && activeEditorGroup === "secondary");
  const primaryInfo = editorGroupInfo("primary");
  const secondaryInfo = editorGroupInfo("secondary");
  const activeInfo = activeEditorGroup === "secondary" ? secondaryInfo : primaryInfo;
  const primaryHasTab = Boolean(primaryInfo.tab);
  const secondaryHasTab = editorGroupsSplit && Boolean(secondaryInfo.tab);
  primaryEditorGroup.classList.toggle("has-tab", primaryHasTab);
  secondaryEditorGroup.classList.toggle("has-tab", secondaryHasTab);
  primaryEditorGroupTitle.textContent = `Group 1${primaryHasTab ? ` · ${primaryInfo.viewLabel}` : ""}`;
  primaryEditorGroupMeta.textContent = primaryHasTab ? primaryInfo.selectionLabel : editorGroupsSplit ? "Empty" : "Ready";
  secondaryEditorGroupTitle.textContent = `Group 2${secondaryHasTab ? ` · ${secondaryInfo.viewLabel}` : ""}`;
  secondaryEditorGroupMeta.textContent = secondaryHasTab ? secondaryInfo.selectionLabel : editorGroupsSplit ? "Empty" : primaryInfo.rowLabel;
  primaryEditorTab.hidden = !primaryHasTab;
  secondaryEditorTab.hidden = !secondaryHasTab;
  primaryEditorTab.textContent = primaryInfo.tabLabel;
  secondaryEditorTab.textContent = secondaryInfo.tabLabel;
  primaryEditorTab.title = primaryHasTab ? `${primaryInfo.fileLabel} · ${primaryInfo.selectionLabel}` : "";
  secondaryEditorTab.title = secondaryHasTab ? `${secondaryInfo.fileLabel} · ${secondaryInfo.selectionLabel}` : "";
  primaryEditorGroupEmpty.hidden = !editorGroupsSplit || primaryHasTab || !documentId;
  secondaryEditorGroupEmpty.hidden = !editorGroupsSplit || secondaryHasTab || !documentId;
  secondaryEditorGroupPreview.hidden = !editorGroupsSplit || secondaryHasTab;
  secondaryEditorGroupFile.textContent = secondaryHasTab || !editorGroupsSplit ? secondaryInfo.fileLabel : "Empty group";
  secondaryEditorGroupView.textContent = secondaryHasTab ? secondaryInfo.viewLabel : "Ready";
  secondaryEditorGroupSelection.textContent = secondaryHasTab ? secondaryInfo.selectionLabel : "Ready";
  document.body.dataset.editorGroupCount = editorGroupsSplit ? "2" : "1";
  document.body.dataset.editorGroupOrientation = editorGroupOrientation;
  document.body.dataset.activeEditorGroup = activeEditorGroup;
  document.body.dataset.currentEditorTabGroup = currentEditorTabGroup;
  document.body.dataset.editorGroupLayout = editorGroupLayoutSignature();
  const currentGroupInfo = editorGroupInfo(currentEditorTabGroup);
  document.body.dataset.editorTabLabel = documentId
    ? `${currentGroupInfo.fileLabel}|${viewModeLabel()}|${currentGroupInfo.selectionLabel}`
    : "";
  document.body.dataset.secondaryEditorGroup = editorGroupsSplit
    ? `${secondaryHasTab ? secondaryInfo.fileLabel : "Empty"}|${secondaryHasTab ? secondaryInfo.viewLabel : "Ready"}|${secondaryHasTab ? secondaryInfo.selectionLabel : "Ready"}`
    : "";
  document.body.dataset.editorGroupSwitch = activeInfo.tab
    ? `${activeEditorGroup}:${activeInfo.fileLabel}`
    : `${activeEditorGroup}:`;
  document.body.dataset.editorGroupStacks = editorGroupStackSignature();
  splitEditorButton.disabled = editorGroupsSplit;
  splitOrientationButton.disabled = !editorGroupsSplit;
  focusEditorGroupButton.disabled = !editorGroupsSplit;
  moveEditorTabButton.disabled = !editorGroupsSplit || !documentId;
  closeSplitEditorButton.disabled = !editorGroupsSplit;
}

function normalizeEditorGroup(group) {
  return editorGroupsSplit && group === "secondary" ? "secondary" : "primary";
}

function tabGroup(tab) {
  return editorGroupsSplit && tab?.group === "secondary" ? "secondary" : "primary";
}

function tabsForEditorGroup(group) {
  if (!editorGroupsSplit && group === "secondary") return [];
  const normalized = normalizeEditorGroup(group);
  return openFileTabs.filter((tab) => tabGroup(tab) === normalized);
}

function tabForEditorGroup(group) {
  const normalized = normalizeEditorGroup(group);
  const tabs = tabsForEditorGroup(normalized);
  if (!tabs.length) return null;
  const activeID = editorGroupActiveTabs[normalized];
  return tabs.find((tab) => tab.documentId && tab.documentId === activeID) || tabs[tabs.length - 1];
}

function currentOpenTargetGroup(options = {}) {
  return normalizeEditorGroup(options.group || activeEditorGroup || currentEditorTabGroup);
}

function syncEditorGroupActiveTabs() {
  if (!editorGroupsSplit) {
    currentEditorTabGroup = "primary";
    editorGroupActiveTabs.primary = documentId || editorGroupActiveTabs.primary || "";
    editorGroupActiveTabs.secondary = "";
    return;
  }
  if (documentId) {
    const currentTab = openFileTabs.find((tab) => tab.documentId === documentId);
    currentEditorTabGroup = tabGroup(currentTab || { group: currentEditorTabGroup });
    editorGroupActiveTabs[currentEditorTabGroup] = documentId;
  }
  for (const group of ["primary", "secondary"]) {
    const activeID = editorGroupActiveTabs[group];
    const activeStillValid = activeID && tabsForEditorGroup(group).some((tab) => tab.documentId === activeID);
    if (!activeStillValid) {
      editorGroupActiveTabs[group] = tabForEditorGroup(group)?.documentId || "";
    }
  }
}

function setOpenFileTabGroup(targetDocumentId, group) {
  const normalized = normalizeEditorGroup(group);
  openFileTabs = openFileTabs.map((tab) => (
    tab.documentId === targetDocumentId ? { ...tab, group: normalized } : tab
  ));
  if (targetDocumentId === documentId) currentEditorTabGroup = normalized;
  editorGroupActiveTabs[normalized] = targetDocumentId || editorGroupActiveTabs[normalized] || "";
}

function editorGroupInfo(group) {
  const normalized = normalizeEditorGroup(group);
  const tab = tabForEditorGroup(normalized);
  const isCurrent = Boolean(tab?.documentId && tab.documentId === documentId);
  const fileLabel = isCurrent
    ? tab?.name || activeFileEntry?.name || fileNameFromPath(metadata?.path || "") || "No file open"
    : tab?.name || "No file open";
  const viewLabel = isCurrent
    ? viewModeLabel()
    : viewModeLabel(tab?.viewState?.viewMode || "grid");
  const rowLabel = isCurrent
    ? `${sourceRowCount.toLocaleString()} rows, ${columns.length} columns`
    : `${formatCount(tab?.rowCount || 0)} rows`;
  const selectionLabel = isCurrent
    ? `row ${selected.row + 1}, ${columns[selected.col]?.name || "Column"}`
    : tabSelectionLabel(tab);
  const tabLabel = tab ? `${fileLabel}${(isCurrent ? dirty : tab.dirty) ? " *" : ""}` : "";
  return { tab, fileLabel, viewLabel, rowLabel, selectionLabel, tabLabel };
}

function tabSelectionLabel(tab) {
  if (!tab?.viewState?.selected) return "Ready";
  const row = Math.max(0, Number(tab.viewState.selected.row) || 0) + 1;
  const col = Math.max(0, Number(tab.viewState.selected.col) || 0) + 1;
  return `row ${row}, column ${col}`;
}

function editorGroupStackSignature() {
  const groupSignature = (group) => tabsForEditorGroup(group)
    .map((tab) => `${tab.documentId === editorGroupActiveTabs[group] ? "*" : ""}${tab.name}`)
    .join(",");
  return `primary:${groupSignature("primary")}|secondary:${groupSignature("secondary")}`;
}

function loadEditorGroupLayout() {
  try {
    const parsed = JSON.parse(localStorage.getItem(editorGroupLayoutStorageKey) || "{}");
    const split = Boolean(parsed?.split);
    const orientation = parsed?.orientation === "vertical" ? "vertical" : "horizontal";
    return {
      split,
      orientation,
      activeGroup: split && parsed?.activeGroup === "secondary" ? "secondary" : "primary",
      currentTabGroup: split && parsed?.currentTabGroup === "secondary" ? "secondary" : "primary",
    };
  } catch {
    return {
      split: false,
      orientation: "horizontal",
      activeGroup: "primary",
      currentTabGroup: "primary",
    };
  }
}

function saveEditorGroupLayout() {
  const layout = {
    split: editorGroupsSplit,
    orientation: editorGroupOrientation === "vertical" ? "vertical" : "horizontal",
    activeGroup: editorGroupsSplit && activeEditorGroup === "secondary" ? "secondary" : "primary",
    currentTabGroup: editorGroupsSplit && currentEditorTabGroup === "secondary" ? "secondary" : "primary",
    updatedAt: Date.now(),
  };
  try {
    localStorage.setItem(editorGroupLayoutStorageKey, JSON.stringify(layout));
    document.body.dataset.editorGroupLayoutSaved = editorGroupLayoutSignature(layout);
  } catch {
    status.textContent = "Editor group layout could not be saved.";
  }
}

function editorGroupLayoutSignature(layout = {
  split: editorGroupsSplit,
  orientation: editorGroupOrientation,
  activeGroup: activeEditorGroup,
  currentTabGroup: currentEditorTabGroup,
}) {
  const split = layout.split ? "split" : "single";
  const orientation = layout.orientation === "vertical" ? "vertical" : "horizontal";
  const activeGroup = layout.split && layout.activeGroup === "secondary" ? "secondary" : "primary";
  const currentTabGroup = layout.split && layout.currentTabGroup === "secondary" ? "secondary" : "primary";
  return `${split}:${orientation}:${activeGroup}:${currentTabGroup}`;
}

function editorGroupBody(group) {
  return group === "secondary" ? secondaryEditorGroupBody : primaryEditorGroupBody;
}

function mountCurrentEditorTab() {
  const targetBody = editorGroupBody(currentEditorTabGroup);
  if (!targetBody || gridFrame.parentElement === targetBody) return;
  targetBody.appendChild(gridFrame);
}

function viewModeLabel(mode = viewMode) {
  switch (mode) {
    case "sql":
      return "SQL";
    case "diff":
      return "Compare";
    case "profile":
      return "Profile";
    case "settings":
      return "Settings";
    default:
      return "Grid";
  }
}

function syncAppTitle() {
  const fileLabel = activeFileEntry?.name || fileNameFromPath(metadata?.path || "") || "";
  const dirtyLabel = dirty ? " *" : "";
  const surfaceLabel = documentId || viewMode !== "grid" ? viewModeLabel(viewMode) : "";
  const parts = [
    fileLabel ? `${fileLabel}${dirtyLabel}` : "",
    surfaceLabel,
    "opencsv",
  ].filter(Boolean);
  const title = parts.join(" - ");
  document.title = title;
  document.body.dataset.appTitle = title;
  document.body.dataset.appTitleFile = fileLabel;
  document.body.dataset.appTitleDirty = dirty ? "true" : "false";
  document.body.dataset.appTitleSurface = surfaceLabel;
}

function splitEditorGroup() {
  if (editorGroupsSplit) return;
  editorGroupsSplit = true;
  if (documentId) {
    currentEditorTabGroup = "secondary";
    setOpenFileTabGroup(documentId, "secondary");
  }
  editorGroupActiveTabs.primary = tabForEditorGroup("primary")?.documentId || "";
  editorGroupActiveTabs.secondary = documentId || editorGroupActiveTabs.secondary || "";
  activeEditorGroup = documentId ? currentEditorTabGroup : "secondary";
  renderEditorGroups();
  renderFileList();
  saveEditorGroupLayout();
  void focusEditorGroup(activeEditorGroup, { activateTab: false });
  resize();
  status.textContent = "Split editor group opened";
}

function closeSplitEditorGroup() {
  if (!editorGroupsSplit) return;
  openFileTabs = openFileTabs.map((tab) => ({ ...tab, group: "primary" }));
  editorGroupActiveTabs = { primary: documentId || editorGroupActiveTabs.primary || "", secondary: "" };
  currentEditorTabGroup = "primary";
  editorGroupsSplit = false;
  activeEditorGroup = "primary";
  renderEditorGroups();
  renderFileList();
  saveEditorGroupLayout();
  void focusEditorGroup("primary", { activateTab: false });
  resize();
  status.textContent = "Split editor group closed";
}

function toggleSplitOrientation() {
  if (!editorGroupsSplit) return;
  editorGroupOrientation = editorGroupOrientation === "horizontal" ? "vertical" : "horizontal";
  renderEditorGroups();
  saveEditorGroupLayout();
  resize();
  status.textContent = `Split layout ${editorGroupOrientation}`;
}

async function focusNextEditorGroup() {
  if (!editorGroupsSplit) return;
  await focusEditorGroup(activeEditorGroup === "primary" ? "secondary" : "primary");
}

function moveCurrentTabToOtherGroup() {
  const targetGroup = currentEditorTabGroup === "primary" ? "secondary" : "primary";
  moveCurrentTabToGroup(targetGroup, "command");
}

async function moveEditorGroupTabToGroup(sourceGroup, targetGroup, source = "command") {
  const sourceTab = tabForEditorGroup(sourceGroup);
  if (!sourceTab?.documentId) {
    status.textContent = `No tab in group ${sourceGroup === "secondary" ? "2" : "1"}.`;
    return;
  }
  if (sourceTab.documentId !== documentId) {
    await switchOpenFileTab(sourceTab.documentId, { group: sourceGroup });
  }
  moveCurrentTabToGroup(targetGroup, source);
}

function moveCurrentTabToGroup(group, source = "command") {
  if (!documentId) {
    status.textContent = "No document tab to move.";
    return;
  }
  if (group === "secondary" && !editorGroupsSplit) {
    editorGroupsSplit = true;
  }
  if (group !== "secondary") group = "primary";
  if (currentEditorTabGroup === group) {
    void focusEditorGroup(group);
    status.textContent = `Tab already in group ${group === "secondary" ? "2" : "1"}`;
    return;
  }
  const previousGroup = currentEditorTabGroup;
  setOpenFileTabGroup(documentId, group);
  editorGroupActiveTabs[group] = documentId;
  editorGroupActiveTabs[previousGroup] = tabForEditorGroup(previousGroup)?.documentId || "";
  currentEditorTabGroup = group;
  activeEditorGroup = group;
  renderEditorGroups();
  renderFileList();
  saveEditorGroupLayout();
  resize();
  void focusEditorGroup(group, { activateTab: false });
  const fileLabel = activeFileEntry?.name || fileNameFromPath(metadata?.path || "") || "Current tab";
  document.body.dataset.editorTabMove = `${previousGroup}->${group}:${selectionSummary()}`;
  document.body.dataset.editorTabMoveSource = source;
  status.textContent = `Moved ${fileLabel} to group ${group === "secondary" ? "2" : "1"}`;
}

async function focusEditorGroup(group, options = {}) {
  activeEditorGroup = editorGroupsSplit && group === "secondary" ? "secondary" : "primary";
  const targetTab = tabForEditorGroup(activeEditorGroup);
  if (options.activateTab !== false && targetTab?.documentId && targetTab.documentId !== documentId) {
    await switchOpenFileTab(targetTab.documentId, { group: activeEditorGroup });
  }
  renderEditorGroups();
  const focusedTab = tabForEditorGroup(activeEditorGroup);
  document.body.dataset.editorGroupSwitch = focusedTab?.documentId
    ? `${activeEditorGroup}:${focusedTab.name}`
    : `${activeEditorGroup}:`;
  saveEditorGroupLayout();
  const activeGroupHasTab = currentEditorTabGroup === activeEditorGroup;
  if (activeGroupHasTab && viewMode === "grid") {
    canvas.focus();
  } else if (activeGroupHasTab && viewMode === "sql") {
    sqlInput.focus();
  } else if (activeGroupHasTab && viewMode === "profile") {
    profilePanel.focus();
  } else if (activeGroupHasTab && viewMode === "settings") {
    settingsSearch.focus();
  } else if (activeEditorGroup === "secondary") {
    secondaryEditorGroup.focus();
  } else {
    primaryEditorGroup.focus();
  }
  status.textContent = activeEditorGroup === "secondary" ? "Focused editor group 2" : "Focused editor group 1";
}

function editorTabDragStart(event) {
  if (!documentId) {
    event.preventDefault();
    return;
  }
  closeGridContextMenu();
  const dragGroup = event.currentTarget === secondaryEditorTab ? "secondary" : "primary";
  event.dataTransfer.effectAllowed = "move";
  event.dataTransfer.setData("text/opencsv-editor-tab", dragGroup);
  event.dataTransfer.setData("text/plain", activeFileEntry?.name || fileNameFromPath(metadata?.path || "") || "opencsv tab");
  document.body.dataset.editorTabDrag = dragGroup;
}

function editorGroupDragOver(event) {
  if (!editorGroupsSplit || !documentId) return;
  if (!Array.from(event.dataTransfer.types).includes("text/opencsv-editor-tab")) return;
  event.preventDefault();
  event.dataTransfer.dropEffect = "move";
  event.currentTarget.classList.add("drop-target");
}

function editorGroupDragLeave(event) {
  event.currentTarget.classList.remove("drop-target");
}

function clearEditorTabDropTargets() {
  primaryEditorGroup.classList.remove("drop-target");
  secondaryEditorGroup.classList.remove("drop-target");
}

function editorGroupDrop(event) {
  if (!editorGroupsSplit || !documentId) return;
  const sourceGroup = event.dataTransfer.getData("text/opencsv-editor-tab");
  if (!sourceGroup) return;
  event.preventDefault();
  const targetGroup = event.currentTarget.dataset.groupId === "secondary" ? "secondary" : "primary";
  clearEditorTabDropTargets();
  void moveEditorGroupTabToGroup(sourceGroup, targetGroup, "drop");
}

function editorTabDragEnd() {
  clearEditorTabDropTargets();
  document.body.dataset.editorTabDrag = "";
}

function editorTabPointerDown(event) {
  if (event.button !== 0 || !documentId) return;
  pendingEditorTabDrag = {
    sourceGroup: event.currentTarget === secondaryEditorTab ? "secondary" : "primary",
    startX: event.clientX,
    startY: event.clientY,
    dragging: false,
  };
}

function editorGroupAtPoint(x, y) {
  const group = document.elementFromPoint(x, y)?.closest(".editor-group")?.dataset.groupId;
  if (group === "secondary") return "secondary";
  if (group === "primary") return "primary";
  return null;
}

function editorTabPointerMove(event) {
  if (!pendingEditorTabDrag || !editorGroupsSplit) return;
  const deltaX = event.clientX - pendingEditorTabDrag.startX;
  const deltaY = event.clientY - pendingEditorTabDrag.startY;
  if (!pendingEditorTabDrag.dragging && Math.hypot(deltaX, deltaY) < 8) return;
  pendingEditorTabDrag.dragging = true;
  document.body.dataset.editorTabDrag = pendingEditorTabDrag.sourceGroup;
  clearEditorTabDropTargets();
  const targetGroup = editorGroupAtPoint(event.clientX, event.clientY);
  if (targetGroup) (targetGroup === "secondary" ? secondaryEditorGroup : primaryEditorGroup).classList.add("drop-target");
  event.preventDefault();
}

function editorTabPointerUp(event) {
  if (!pendingEditorTabDrag) return;
  const dragState = pendingEditorTabDrag;
  pendingEditorTabDrag = null;
  clearEditorTabDropTargets();
  document.body.dataset.editorTabDrag = "";
  if (!dragState.dragging || !editorGroupsSplit || !documentId) return;
  const targetGroup = editorGroupAtPoint(event.clientX, event.clientY);
  if (!targetGroup) return;
  void moveEditorGroupTabToGroup(dragState.sourceGroup, targetGroup, "drop");
}

function openEditorTabContextMenu(event) {
  event.preventDefault();
  if (!documentId) return;
  closeColumnMenu();
  const currentGroup = event.currentTarget === secondaryEditorTab ? "secondary" : "primary";
  void focusEditorGroup(currentGroup);
  const targetGroup = currentGroup === "primary" ? "secondary" : "primary";
  const heading = document.createElement("strong");
  heading.textContent = "Editor Tab";
  gridContextMenu.replaceChildren(heading);
  const addItem = (label, action, options = {}) => {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = label;
    button.disabled = Boolean(options.disabled);
    button.addEventListener("click", () => {
      closeGridContextMenu();
      void action();
    });
    gridContextMenu.append(button);
  };
  addItem(`Move to Group ${targetGroup === "secondary" ? "2" : "1"}`, () => moveEditorGroupTabToGroup(currentGroup, targetGroup, "context"), {
    disabled: !editorGroupsSplit,
  });
  addItem("Focus Next Group", () => focusNextEditorGroup(), { disabled: !editorGroupsSplit });
  gridContextMenu.hidden = false;
  const menuRect = gridContextMenu.getBoundingClientRect();
  const left = Math.min(event.clientX, window.innerWidth - menuRect.width - 8);
  const top = Math.min(event.clientY, window.innerHeight - menuRect.height - 8);
  gridContextMenu.style.left = `${Math.max(8, left)}px`;
  gridContextMenu.style.top = `${Math.max(8, top)}px`;
  gridContextMenu.dataset.target = "editor-tab";
  document.body.dataset.gridContextTarget = "editor-tab";
}

function setViewMode(mode) {
  closeColumnMenu();
  closeGridContextMenu();
  viewMode = mode;
  gridFrame.classList.toggle("sql-mode", mode === "sql");
  gridFrame.classList.toggle("diff-mode", mode === "diff");
  gridFrame.classList.toggle("profile-mode", mode === "profile");
  gridFrame.classList.toggle("settings-mode", mode === "settings");
  sqlPanel.hidden = mode !== "sql";
  diffPanel.hidden = mode !== "diff";
  profilePanel.hidden = mode !== "profile";
  settingsPanel.hidden = mode !== "settings";
  gridTabButton.classList.toggle("active", mode === "grid");
  sqlTabButton.classList.toggle("active", mode === "sql");
  profileTabButton.classList.toggle("active", mode === "profile");
  settingsTabButton.classList.toggle("active", mode === "settings");
  syncStatusbarSurface(mode);
  renderEditorGroups();
  if (mode === "grid") {
    resize();
    canvas.focus();
  } else if (mode === "sql") {
    sqlInput.focus();
  } else if (mode === "diff") {
    leftDiffPath.focus();
  } else if (mode === "profile") {
    profilePanel.focus();
  } else {
    renderSettingsPanel();
    settingsSearch.focus();
  }
  updateCommandState();
  syncAppTitle();
}

function syncStatusbarSurface(mode = viewMode) {
  if (!statusbar) return;
  statusbar.hidden = false;
  const surface = mode === "grid" ? "grid" : "non-grid";
  statusbar.dataset.viewMode = mode;
  statusbar.dataset.editorSurface = surface;
  statusbar.setAttribute("aria-label", `${viewModeLabel(mode)} status`);
  document.body.dataset.statusbarView = mode;
  document.body.dataset.statusbarSurface = surface;
  document.body.dataset.statusbarVisible = String(!statusbar.hidden);
}

function syncPlatformDialogButtonOrder() {
  const platform = platformDialogButtonOrder();
  const confirmLeft = platform === "windows";
  const orderName = confirmLeft ? "confirm-left" : "confirm-right";
  let actionCount = 0;
  for (const actions of document.querySelectorAll("[data-dialog-actions]")) {
    actionCount++;
    actions.dataset.dialogButtonOrder = orderName;
    for (const button of actions.querySelectorAll("button")) {
      const role = button.dataset.dialogRole || "secondary";
      button.style.order = String(dialogButtonOrderValue(role, confirmLeft));
    }
  }
  document.body.dataset.dialogButtonPlatform = platform;
  document.body.dataset.dialogButtonOrder = orderName;
  document.body.dataset.dialogActionCount = String(actionCount);
}

function platformDialogButtonOrder() {
  const params = new URLSearchParams(window.location.search);
  const override = params.get("dialogPlatform") || params.get("dialog-platform") || "";
  const platform = String(
    override ||
    navigator.userAgentData?.platform ||
    navigator.platform ||
    navigator.userAgent ||
    "",
  ).toLowerCase();
  return platform.startsWith("win") || platform.includes("windows") ? "windows" : "standard";
}

function dialogButtonOrderValue(role, confirmLeft) {
  if (role === "confirm") return confirmLeft ? 0 : 2;
  if (role === "cancel") return confirmLeft ? 2 : 0;
  return 1;
}

async function loadCommandRegistry() {
  try {
    const response = await fetch("/api/commands");
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "commands failed");
    commandsById = new Map((payload.commands ?? []).map((command) => [command.id, command]));
    customShortcuts = loadCustomShortcuts();
    for (const commandID of customShortcuts.keys()) {
      if (!commandsById.has(commandID)) customShortcuts.delete(commandID);
    }
    applyCommandMetadata();
    updateCommandState();
  } catch (error) {
    status.textContent = `Commands failed: ${error.message}`;
  }
}

function applyCommandMetadata() {
  for (const button of commandButtons) {
    const command = commandsById.get(button.dataset.command);
    if (!command) continue;
    const shortcut = effectiveShortcut(command);
    button.title = shortcut
      ? `${localizedCommandDescription(command)} (${shortcut})`
      : localizedCommandDescription(command);
    if (button.dataset.command === "edit.replaceAll") {
      button.textContent = "All";
    } else if (button.dataset.command === "command.palette") {
      button.textContent = "Cmd";
    } else if (button.dataset.command === "workbench.splitEditor") {
      button.textContent = "Split";
    } else if (button.dataset.command === "workbench.toggleSplitOrientation") {
      button.textContent = "H/V";
    } else if (button.dataset.command === "workbench.focusNextGroup") {
      button.textContent = "Focus";
    } else if (button.dataset.command === "workbench.moveTabToOtherGroup") {
      button.textContent = "Move";
    } else if (button.dataset.command === "workbench.closeSplitEditor") {
      button.textContent = "Close";
    } else if (button.dataset.command === "file.toggleReadOnly") {
      button.textContent = "Read Only";
    } else if (button.dataset.command === "diff.open") {
      button.textContent = "Compare";
    } else if (button.dataset.command === "file.exportView") {
      button.textContent = "Export";
    } else if (button.dataset.command === "file.exportXlsx") {
      button.textContent = "XLSX";
    } else if (button.dataset.command === "file.restoreBackup") {
      button.textContent = "Restore";
    } else if (button.dataset.command === "edit.copySelectedFormat") {
      button.textContent = "Copy";
    } else if (button.dataset.command === "edit.replacePreview") {
      button.textContent = "Preview";
    } else if (button.dataset.command === "view.filterBuilder") {
      button.textContent = "Builder";
    } else if (button.dataset.command === "format.savePreset") {
      button.textContent = "Save Preset";
    } else if (button.dataset.command === "format.applyPreset") {
      button.textContent = "Apply";
    } else if (button.dataset.command === "format.saveRule") {
      button.textContent = "Save Rule";
    } else if (button.dataset.command === "format.applyDelimiter") {
      button.textContent = "Apply Delim";
    } else if (button.dataset.command === "format.applyEncoding") {
      button.textContent = "Apply Enc";
    } else if (button.dataset.command === "format.applyLineEnding") {
      button.textContent = "Apply EOL";
    } else if (button.dataset.command === "format.applyQuotePolicy") {
      button.textContent = "Apply Quote";
    } else if (button.dataset.command === "format.exportQuoteRules") {
      button.textContent = "Export Quotes";
    } else if (button.dataset.command === "format.importQuoteRules") {
      button.textContent = "Import Quotes";
    } else if (button.dataset.command === "format.toggleCellQuote") {
      button.textContent = "Toggle";
    } else if (button.dataset.command === "column.showAll") {
      button.textContent = t("column.showAllShort");
    } else if (button.dataset.command === "sql.copyResultCsv") {
      button.textContent = "Copy CSV";
    } else if (button.dataset.command === "sql.exportResultCsv") {
      button.textContent = "Export CSV";
    } else if (button.dataset.command === "selection.cycleAggregate") {
      updateAggregateStatus(selectionAggregateStats());
    } else if (button.dataset.command === "settings.exportShortcuts") {
      button.textContent = "Export";
    } else if (button.dataset.command === "settings.importShortcuts") {
      button.textContent = "Import";
    } else if (button.dataset.command === "settings.checkForUpdates") {
      button.textContent = "Check Now";
    } else if (button.dataset.command === "view.exportSavedSets") {
      button.textContent = "Export Views";
    } else if (button.dataset.command === "view.importSavedSets") {
      button.textContent = "Import Views";
    } else if (button.dataset.command !== "view.sql" && button.dataset.command !== "view.profile") {
      button.textContent = localizedCommandTitle(command).replace("Replace All", t("edit.all"));
    }
  }
  applyLocaleStrings();
}

function updateCommandState() {
  const hasDocument = Boolean(documentId);
  for (const button of commandButtons) {
    const command = commandsById.get(button.dataset.command);
    if (!command) continue;
    const disabled = isCommandDisabled(command);
    button.disabled = disabled;
    button.setAttribute("aria-disabled", String(disabled));
  }
  readOnlyButton.classList.toggle("active", readOnly);
  readOnlyButton.setAttribute("aria-pressed", String(readOnly));
  sqlTabButton.disabled = !hasDocument || commandsById.get("view.sql")?.implemented === false;
  profileTabButton.disabled = !hasDocument || commandsById.get("view.profile")?.implemented === false;
  if (!hasDocument) {
    findPrevButton.disabled = true;
    findNextButton.disabled = true;
  }
  copyFormatSelect.disabled = !hasDocument;
  copyHeadersToggle.disabled = !hasDocument;
  formatPresetSelect.disabled = !hasDocument || formatPresets().length === 0;
  saveFormatPresetButton.disabled = !hasDocument;
  applyFormatPresetButton.disabled = !hasDocument || readOnly || !formatPresetSelect.value;
  formatRuleInput.disabled = !hasDocument;
  saveFormatRuleButton.disabled = !hasDocument;
  delimiterSelect.disabled = !hasDocument || readOnly;
  delimiterCustomInput.disabled = !hasDocument || readOnly || delimiterSelect.value !== "custom";
  applyDelimiterButton.disabled = !hasDocument || readOnly;
  encodingSelect.disabled = !hasDocument || readOnly;
  applyEncodingButton.disabled = !hasDocument || readOnly;
  quoteModeSelect.disabled = !hasDocument;
  quoteCharSelect.disabled = !hasDocument;
  quoteEscapeSelect.disabled = !hasDocument;
  applyQuotePolicyButton.disabled = !hasDocument || readOnly;
  lineEndingSelect.disabled = !hasDocument || readOnly;
  applyLineEndingButton.disabled = !hasDocument || readOnly;
  if (workbookSheetSelect) {
    workbookSheetSelect.disabled = !canSwitchWorkbookSheet();
  }
  quoteHeaderRuleToggle.disabled = !hasDocument || readOnly;
  quoteBodyRuleToggle.disabled = !hasDocument || readOnly;
  quoteColumnRuleToggle.disabled = !hasDocument || readOnly || columns.length === 0;
  toggleCellQuoteButton.disabled = !hasDocument || readOnly || sourceRowCount <= 0;
}

function isCommandDisabled(command) {
  const hasDocument = Boolean(documentId);
  const needsSQLResult = sqlResultCommandIDs.has(command.id);
  const needsSQLTarget = sqlResultTargetCommandIDs.has(command.id);
  const needsSQLSelection = sqlResultSelectionCommandIDs.has(command.id);
  const needsActiveFilter = command.id === "view.saveFilter";
  const needsMultiSort = command.id === "view.removeLastSortKey";
  const needsCellLink = command.id === "cell.openLink";
  const needsColumnTextTransform = columnTextTransformCommands.has(command.id);
  const needsColumnAction = needsColumnTextTransform || columnSortCommands.has(command.id) || columnFilterCommands.has(command.id);
  const needsAggregate =
    command.id === "selection.copyAggregate" ||
    command.id === "selection.cycleAggregate" ||
    command.id === "selection.insertAggregateRow" ||
    aggregateRowCommandMetrics.has(command.id);
  const needsColumn = command.id === "format.toggleColumnQuote" || needsColumnAction;
  const needsBodyRow = command.id === "row.cut" || command.id === "row.paste" || command.id === "row.pasteBelow";
  const needsMappedRows = command.id === "row.showAll";
  const needsVisibleColumn = command.id === "column.cut" || command.id === "column.paste" || command.id === "column.pasteRight";
  const needsHiddenColumn = command.id === "column.showAll";
  const needsRowClipboard = command.id === "row.paste" || command.id === "row.pasteBelow";
  const needsColumnClipboard = command.id === "column.paste" || command.id === "column.pasteRight";
  const needsUpdateChecksEnabled = command.id === "settings.checkForUpdates";
  const needsWorkbookSheet = command.id === "format.switchWorksheet";
  const needsProfileIdle = command.id === "profile.run";
  const needsProfileCancel = command.id === "profile.cancel";
  const createsSplitGroup = command.id === "workbench.splitEditor";
  const needsSplitGroup =
    command.id === "workbench.closeSplitEditor" ||
    command.id === "workbench.focusNextGroup" ||
    command.id === "workbench.moveTabToOtherGroup" ||
    command.id === "workbench.toggleSplitOrientation";
  const writeBlocked =
    readOnly &&
    (
      command.id === "file.save" ||
      command.id === "file.saveAs" ||
      command.id === "file.restoreBackup" ||
      command.id === "edit.undo" ||
      command.id === "edit.redo" ||
      command.id === "edit.paste" ||
      command.id === "edit.replace" ||
      command.id === "edit.replaceAll" ||
      command.id === "row.insert" ||
      command.id === "row.delete" ||
      command.id === "row.cut" ||
      command.id === "row.paste" ||
      command.id === "row.pasteBelow" ||
      command.id === "row.hide" ||
      command.id === "row.moveUp" ||
      command.id === "row.moveDown" ||
      command.id === "column.insert" ||
      command.id === "column.delete" ||
      command.id === "column.cut" ||
      command.id === "column.paste" ||
      command.id === "column.pasteRight" ||
      command.id === "column.hide" ||
      command.id === "column.moveLeft" ||
      command.id === "column.moveRight" ||
      command.id === "table.transpose" ||
      command.id === "table.toggleHeaderRow" ||
      command.id === "table.dedupeRows" ||
      command.id === "table.dedupeColumns" ||
      command.id === "text.trim" ||
      command.id === "text.uppercase" ||
      command.id === "text.lowercase" ||
      command.id === "text.snakeCase" ||
      command.id === "text.kebabCase" ||
      command.id === "text.camelCase" ||
      command.id === "text.pascalCase" ||
      command.id === "text.titleCase" ||
      command.id === "text.fullWidth" ||
      command.id === "text.halfWidth" ||
      command.id === "text.normalizeNFC" ||
      command.id === "text.normalizeNFD" ||
      command.id === "text.normalizeNFKC" ||
      command.id === "text.normalizeNFKD" ||
      needsColumnTextTransform ||
      command.id === "cell.swap" ||
      command.id === "cell.swapUp" ||
      command.id === "cell.swapDown" ||
      command.id === "cell.swapLeft" ||
      command.id === "cell.swapRight" ||
      command.id === "cell.joinHorizontal" ||
      command.id === "cell.splitHorizontal" ||
      command.id === "cell.joinVertical" ||
      command.id === "cell.splitVertical" ||
      command.id === "cell.fillDown" ||
      command.id === "cell.fillRight" ||
      command.id === "cell.fillUp" ||
      command.id === "cell.fillLeft" ||
      command.id === "format.applyPreset" ||
      command.id === "format.applyDelimiter" ||
      command.id === "format.applyEncoding" ||
      command.id === "format.applyQuotePolicy" ||
      command.id === "format.applyLineEnding" ||
      command.id === "format.importQuoteRules" ||
      command.id === "format.toggleHeaderQuote" ||
      command.id === "format.toggleBodyQuote" ||
      command.id === "format.toggleColumnQuote" ||
      command.id === "format.toggleCellQuote" ||
      command.id === "format.quoteSelection" ||
      command.id === "format.noQuoteSelection" ||
      command.id === "format.clearSelectionQuote"
    );
  const needsSelectedCells =
    command.id === "format.quoteSelection" ||
    command.id === "format.noQuoteSelection" ||
    command.id === "format.clearSelectionQuote";
  return (
    command.implemented === false ||
    (command.requiresDocument && !hasDocument) ||
    (needsSQLResult && !sqlResultAvailable()) ||
    (needsSQLTarget && !sqlResultTargetForCommand()) ||
    (needsSQLSelection && !sqlSelectedResultMatrix(lastSQLResult)) ||
    (needsActiveFilter && !activeFilter) ||
    (needsMultiSort && sortKeys.length <= 1) ||
    (needsAggregate && currentSelectionAggregateStats().count <= 1) ||
    (needsColumn && columns.length === 0) ||
    (needsVisibleColumn && columns.length === 0) ||
    (needsMappedRows && !viewMapped) ||
    (needsHiddenColumn && hiddenSourceColumns.size === 0) ||
    (needsCellLink && !normalizeURL(valueAt(selected.row, selected.col))) ||
    (needsBodyRow && selected.row >= sourceRowCount) ||
    (needsRowClipboard && !structureClipboardFor("row")) ||
    (needsColumnClipboard && !structureClipboardFor("column")) ||
    (needsSelectedCells && selectionCellCount() === 0) ||
    (needsUpdateChecksEnabled && updateCheckMode === "off") ||
    (needsWorkbookSheet && !canSwitchWorkbookSheet()) ||
    (needsProfileIdle && profileLoading) ||
    (needsProfileCancel && !activeProfileController) ||
    (createsSplitGroup && editorGroupsSplit) ||
    (needsSplitGroup && !editorGroupsSplit) ||
    writeBlocked
  );
}

function sqlResultAvailable() {
  const columns = sqlResultColumns(lastSQLResult);
  return Boolean(lastSQLResult && (columns.length || lastSQLResult.rows?.length));
}

function updateInspector() {
  const selectedColumn = columns[selected.col]?.name ?? "";
  const profile = profileData?.columns?.[sourceColFor(selected.col)];
  columnName.textContent = selectedColumn || "none";
  columnType.textContent = profile?.type ?? columns[selected.col]?.type ?? "unknown";
  columnNulls.textContent = profile
    ? `${formatCount(profile.nulls)} (${formatPercent(profile.nullRatio)})`
    : "unknown";
  columnDistinct.textContent = profile ? formatCount(profile.distinct) : "unknown";
  cellQuote.textContent = documentId ? quoteLabel(valueAt(selected.row, selected.col), selected.row, selected.col) : "unknown";
  formatSummary.textContent = metadata ? documentFormatSummary() : "unknown";
  syncDelimiterControls();
  syncEncodingControls();
  syncQuotePolicyControls();
  syncLineEndingControls();
  syncCellValuePanel();
}

function updateFormatStatus() {
  if (!metadata) {
    encodingStatus.textContent = "Encoding";
    delimiterStatus.textContent = "Delimiter";
    quoteStatus.textContent = "Quote minimal";
    headerStatus.textContent = "Header";
    headerStatus.disabled = true;
    headerStatus.setAttribute("aria-pressed", "false");
    document.body.dataset.headerColumns = "";
    document.body.dataset.readOnly = "false";
    document.body.dataset.workbookSheet = "";
    document.body.dataset.workbookSheetCount = "";
    readOnlyStatus.textContent = "Editable";
    syncWorkbookSheetControls();
    syncExternalChangeControls(false);
    syncDelimiterControls();
    syncEncodingControls();
    syncLineEndingControls();
    return;
  }
  encodingStatus.textContent = encodingLabel();
  delimiterStatus.textContent = delimiterLabel(metadata.dialect);
  quoteStatus.textContent = quotePolicyLabel(metadata.dialect);
  headerStatus.textContent = headerModeLabel(metadata.dialect);
  headerStatus.disabled = !documentId || readOnly;
  headerStatus.setAttribute("aria-pressed", String(Boolean(metadata.dialect?.hasHeader)));
  document.body.dataset.headerColumns = String(headerColumnsValue(metadata.dialect));
  document.body.dataset.readOnly = readOnly ? "true" : "false";
  document.body.dataset.workbookSheet = metadata.workbook?.selectedSheetName || "";
  document.body.dataset.workbookSheetCount = String(metadata.workbook?.sheets?.length || "");
  readOnlyStatus.textContent = readOnly ? "Read-only" : "Editable";
  syncWorkbookSheetControls();
  syncExternalChangeControls();
  syncDelimiterControls();
  syncEncodingControls();
  syncQuotePolicyControls();
  syncLineEndingControls();
}

function syncExternalChangeControls(hasDocument = Boolean(metadata)) {
  const changed = Boolean(hasDocument && externalFileChanged);
  const conflict = Boolean(changed && dirty);
  document.body.dataset.externalFileChanged = String(changed);
  document.body.dataset.externalConflict = conflict ? "dirty" : changed ? "available" : "";
  externalStatus.textContent = changed ? "Changed on disk" : "Synced";
  externalStatus.disabled = !changed;
  externalStatus.title = changed
    ? conflict
      ? "Resolve external file change with unsaved edits"
      : "Refresh changed file from disk"
    : "Synced";
  if (!changed && externalChangePanel && !externalChangePanel.hidden) {
    closeExternalChangePanel();
  } else if (changed && externalChangePanel && !externalChangePanel.hidden) {
    updateExternalChangePanel();
  }
}

function encodingValue(info = metadata?.encoding) {
  const name = String(info?.name || "utf-8").toLowerCase();
  if (name === "utf8") return "utf-8";
  if (name === "utf8bom") return "utf-8-bom";
  if (name === "utf-16le" && info?.hasBom) return "utf-16le-bom";
  if (name === "utf-16be" && info?.hasBom) return "utf-16be-bom";
  return name;
}

function encodingLabel(info = metadata?.encoding) {
  const name = encodingValue(info);
  if (name === "utf-8-bom") return "UTF-8 BOM";
  if (name === "utf-8") return "UTF-8";
  if (name === "utf-16le-bom") return "UTF-16 LE BOM";
  if (name === "utf-16le") return info?.hasBom ? "UTF-16 LE BOM" : "UTF-16 LE";
  if (name === "utf-16be-bom") return "UTF-16 BE BOM";
  if (name === "utf-16be") return info?.hasBom ? "UTF-16 BE BOM" : "UTF-16 BE";
  if (name === "shift-jis") return "Shift-JIS";
  if (name === "windows-1252") return "Windows-1252";
  if (name === "legacy-single-byte") return "Legacy byte";
  return name.toUpperCase();
}

function delimiterLabel(dialect) {
  if (!dialect) return "Delimiter";
  const name = dialect.delimiterName || delimiterChar(dialect);
  if (name === "\t") return "Tab";
  return name.charAt(0).toUpperCase() + name.slice(1);
}

function delimiterChar(dialect = metadata?.dialect) {
  const raw = dialect?.delimiter;
  if (typeof raw === "number" && Number.isFinite(raw) && raw > 0) {
    return String.fromCodePoint(raw);
  }
  if (typeof raw === "string" && raw.length > 0) return raw;
  return ",";
}

function delimiterOptionValue(dialect = metadata?.dialect) {
  const name = String(dialect?.delimiterName || "").trim().toLowerCase();
  const char = delimiterChar(dialect);
  if (name === "comma" || char === ",") return "comma";
  if (name === "tab" || name === "tsv" || char === "\t") return "tab";
  if (name === "semicolon" || char === ";") return "semicolon";
  if (name === "pipe" || char === "|") return "pipe";
  if (name === "space" || char === " ") return "space";
  return "custom";
}

function delimiterCustomValue(dialect = metadata?.dialect) {
  return delimiterOptionValue(dialect) === "custom" ? delimiterChar(dialect) : "";
}

function delimiterControlValue() {
  const option = delimiterSelect?.value || "comma";
  if (option === "custom") {
    const custom = delimiterCustomInput.value || ",";
    return { delimiter: Array.from(custom)[0] || ",", delimiterName: Array.from(custom)[0] || "," };
  }
  return normalizeDelimiterOverride(option);
}

function lineEndingLabel(info = metadata?.lineEnding) {
  return (info?.name ?? "LF").toUpperCase();
}

function normalizeLineEndingName(value) {
  const raw = String(value || "lf").trim().toLowerCase();
  if (raw === "crlf" || raw === "windows" || raw === "\\r\\n" || raw === "\r\n") return "crlf";
  if (raw === "cr" || raw === "classic-mac" || raw === "classic mac" || raw === "\\r" || raw === "\r") return "cr";
  return "lf";
}

function lineEndingValue(info = metadata?.lineEnding) {
  return normalizeLineEndingName(info?.name);
}

function lineEndingSequence(value) {
  const normalized = normalizeLineEndingName(value);
  if (normalized === "crlf") return "\r\n";
  if (normalized === "cr") return "\r";
  return "\n";
}

function quoteModeValue(dialect = metadata?.dialect) {
  const raw = String(dialect?.quoteMode || "minimal").toLowerCase();
  return ["minimal", "always", "never", "manual"].includes(raw) ? raw : "minimal";
}

function quoteCharValue(dialect = metadata?.dialect) {
  const raw = dialect?.quoteChar;
  if (raw === "'" || raw === "single") return "single";
  if (raw === "" || raw === "none") return "none";
  return "double";
}

function quoteCharLiteral(dialect = metadata?.dialect) {
  return quoteCharValue(dialect) === "single" ? "'" : quoteCharValue(dialect) === "none" ? "" : '"';
}

function quoteCharLabel(dialect = metadata?.dialect) {
  const value = quoteCharValue(dialect);
  if (value === "single") return "single";
  if (value === "none") return "none";
  return "double";
}

function normalizeQuoteEscapeValue(value) {
  const raw = String(value || "double").toLowerCase().trim();
  return raw === "backslash" || raw === "slash" ? "backslash" : "double";
}

function quoteEscapeValue(dialect = metadata?.dialect) {
  return normalizeQuoteEscapeValue(dialect?.quoteEscape);
}

function quoteEscapeLabel(dialect = metadata?.dialect) {
  return `${quoteEscapeValue(dialect)} escape`;
}

function selectedQuotePolicySummary() {
  return `${quoteModeSelect.value} ${quoteCharSelect.value} ${quoteEscapeSelect.value} escape`;
}

function quotePolicyLabel(dialect = metadata?.dialect) {
  return `Quote ${quoteModeValue(dialect)} ${quoteCharLabel(dialect)} ${quoteEscapeLabel(dialect)}`;
}

function headerColumnsValue(dialect = metadata?.dialect) {
  return normalizeHeaderColumns(dialect?.headerColumns);
}

function normalizeHeaderColumns(value) {
  const parsed = Number.parseInt(value, 10);
  return Number.isInteger(parsed) && parsed > 0 ? parsed : 0;
}

function headerModeLabel(dialect = metadata?.dialect) {
  const headerColumns = headerColumnsValue(dialect);
  const columnLabel = headerColumns === 1 ? "1 header col" : `${headerColumns} header cols`;
  if (dialect?.hasHeader && headerColumns > 0) return `Header row + ${columnLabel}`;
  if (dialect?.hasHeader) return "Header row";
  if (headerColumns > 0) return columnLabel;
  return "No header";
}

function workbookSheetOptions() {
  const sheets = metadata?.workbook?.sheets;
  return Array.isArray(sheets) ? sheets : [];
}

function workbookSheetSignature(sheets = workbookSheetOptions()) {
  return sheets.map((sheet) => `${sheet.index}:${sheet.name}`).join("|");
}

function syncWorkbookSheetControls() {
  if (!workbookSheetSelect) return;
  const sheets = workbookSheetOptions();
  const signature = workbookSheetSignature(sheets);
  if (workbookSheetSelect.dataset.signature !== signature) {
    const fragment = document.createDocumentFragment();
    if (!sheets.length) {
      const option = document.createElement("option");
      option.value = "";
      option.textContent = "Sheet";
      fragment.append(option);
    } else {
      for (const sheet of sheets) {
        const option = document.createElement("option");
        option.value = String(sheet.index);
        option.textContent = sheet.name || `Sheet ${sheet.index}`;
        fragment.append(option);
      }
    }
    workbookSheetSelect.replaceChildren(fragment);
    workbookSheetSelect.dataset.signature = signature;
  }
  const selectedIndex = metadata?.workbook?.selectedSheetIndex;
  const selectedValue = Number.isInteger(selectedIndex) && selectedIndex > 0 ? String(selectedIndex) : "";
  if (selectedValue && [...workbookSheetSelect.options].some((option) => option.value === selectedValue)) {
    workbookSheetSelect.value = selectedValue;
  }
  workbookSheetSelect.hidden = sheets.length <= 1;
  workbookSheetSelect.disabled = !documentId || sheets.length <= 1 || dirty;
  workbookSheetSelect.title = dirty
    ? "Save or reload before switching worksheets"
    : sheets.length > 1
      ? "XLSX worksheet"
      : "No alternate worksheet";
  document.body.dataset.workbookSheetOptions = signature;
  document.body.dataset.workbookSheetSwitchBlocked = dirty && sheets.length > 1 ? "dirty" : "";
}

function canSwitchWorkbookSheet() {
  return Boolean(documentId && workbookSheetOptions().length > 1 && !dirty);
}

function focusWorkbookSheetSelector() {
  const sheets = workbookSheetOptions();
  if (!documentId || !sheets.length) {
    status.textContent = "Open an XLSX workbook before switching worksheets.";
    return;
  }
  if (sheets.length <= 1) {
    status.textContent = "The current workbook has only one worksheet.";
    return;
  }
  if (dirty) {
    status.textContent = "Save or reload before switching worksheets.";
    syncWorkbookSheetControls();
    return;
  }
  syncWorkbookSheetControls();
  workbookSheetSelect.hidden = false;
  workbookSheetSelect.focus();
  status.textContent = `Choose one of ${formatCount(sheets.length)} worksheets.`;
}

async function switchWorkbookSheet() {
  const sheets = workbookSheetOptions();
  const selectedIndex = Number.parseInt(workbookSheetSelect?.value || "", 10);
  const sheet = sheets.find((item) => item.index === selectedIndex);
  if (!sheet || !documentId) {
    syncWorkbookSheetControls();
    return;
  }
  if (sheet.index === metadata?.workbook?.selectedSheetIndex) return;
  if (dirty) {
    status.textContent = "Save or reload before switching worksheets.";
    syncWorkbookSheetControls();
    return;
  }
  const path = metadata?.path || activeFileEntry?.detail || pathInput.value.trim();
  if (!path || path === "temporary upload") {
    status.textContent = "Switch Worksheet needs a local workbook path.";
    syncWorkbookSheetControls();
    return;
  }
  const previousSheet = metadata?.workbook?.selectedSheetName || "";
  status.textContent = `Opening worksheet ${sheet.name || sheet.index}...`;
  document.body.dataset.workbookSheetSwitch = `${previousSheet}->${sheet.name || sheet.index}`;
  try {
    await openPath(path, {
      applyRules: false,
      formatOverride: { sheetIndex: sheet.index },
    });
    status.textContent = `Opened worksheet ${sheet.name || sheet.index}.`;
  } catch (error) {
    status.textContent = `Switch Worksheet failed: ${error.message}`;
    syncWorkbookSheetControls();
  }
}

function syncDelimiterControls() {
  if (!delimiterSelect || !delimiterCustomInput || !applyDelimiterButton) return;
  if (!delimiterDraftDirty) {
    delimiterSelect.value = delimiterOptionValue();
    delimiterCustomInput.value = delimiterCustomValue();
  }
  const custom = delimiterSelect.value === "custom";
  delimiterSelect.disabled = !documentId || readOnly;
  delimiterCustomInput.disabled = !documentId || readOnly || !custom;
  applyDelimiterButton.disabled = !documentId || readOnly;
  document.body.dataset.delimiter = delimiterOptionValue();
  document.body.dataset.delimiterDraft = delimiterDraftDirty
    ? `${delimiterSelect.value}:${custom ? delimiterCustomInput.value : ""}`
    : "";
}

function syncEncodingControls() {
  if (!encodingSelect || !applyEncodingButton) return;
  const value = encodingValue();
  if (!encodingDraftDirty) {
    const hasOption = Array.from(encodingSelect.options).some((option) => option.value === value);
    encodingSelect.value = hasOption ? value : "utf-8";
  }
  encodingSelect.disabled = !documentId || readOnly;
  applyEncodingButton.disabled = !documentId || readOnly;
  document.body.dataset.encoding = value;
  document.body.dataset.encodingDraft = encodingDraftDirty ? encodingSelect.value : "";
}

function syncQuotePolicyControls() {
  if (!quoteModeSelect || !quoteCharSelect || !quoteEscapeSelect) return;
  if (quotePolicyDraftDirty) {
    settingsQuoteSummary.textContent = selectedQuotePolicySummary();
  } else {
    quoteModeSelect.value = quoteModeValue();
    quoteCharSelect.value = quoteCharValue();
    quoteEscapeSelect.value = quoteEscapeValue();
    settingsQuoteSummary.textContent = metadata ? `${quoteModeValue()} ${quoteCharLabel()} ${quoteEscapeLabel()}` : "minimal double double escape";
  }
  const dialect = metadata?.dialect;
  quoteHeaderRuleToggle.checked = Boolean(dialect?.quoteHeader);
  quoteBodyRuleToggle.checked = Boolean(dialect?.quoteBody);
  quoteColumnRuleToggle.checked = Boolean(dialect && isSourceColumnQuoted(sourceColFor(selected.col)));
  const activeSourceCellAvailable = Boolean(documentId && columns.length > 0 && selected.row < sourceRowCount);
  const selectionRuleAvailable = Boolean(documentId && columns.length > 0 && sourceRowCount > 0);
  const coordinateRuleAvailable = Boolean(documentId && columns.length > 0 && sourceRowCount + dataStartRow > 0);
  quoteHeaderRuleToggle.disabled = !documentId || readOnly;
  quoteBodyRuleToggle.disabled = !documentId || readOnly;
  quoteColumnRuleToggle.disabled = !documentId || readOnly || columns.length === 0;
  addQuoteCellRuleButton.disabled = readOnly || !activeSourceCellAvailable;
  addNoQuoteCellRuleButton.disabled = readOnly || !activeSourceCellAvailable;
  clearCellQuoteRuleButton.disabled = readOnly || !activeSourceCellAvailable;
  quoteSelectionRuleButton.disabled = readOnly || !selectionRuleAvailable;
  noQuoteSelectionRuleButton.disabled = readOnly || !selectionRuleAvailable;
  clearSelectionQuoteRuleButton.disabled = readOnly || !selectionRuleAvailable;
  quoteRuleColumnInput.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleColumnActionSelect.disabled = readOnly || !coordinateRuleAvailable;
  applyQuoteColumnRuleButton.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleColumnPatternInput.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleColumnPatternModeSelect.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleColumnPatternActionSelect.disabled = readOnly || !coordinateRuleAvailable;
  applyQuoteColumnPatternRuleButton.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleValuePatternInput.disabled = readOnly || !selectionRuleAvailable;
  quoteRuleValuePatternModeSelect.disabled = readOnly || !selectionRuleAvailable;
  quoteRuleValuePatternActionSelect.disabled = readOnly || !selectionRuleAvailable;
  applyQuoteValuePatternRuleButton.disabled = readOnly || !selectionRuleAvailable;
  quoteRuleRowStartInput.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleRowEndInput.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleRowActionSelect.disabled = readOnly || !coordinateRuleAvailable;
  applyQuoteRowRangeRuleButton.disabled = readOnly || !coordinateRuleAvailable;
  if (coordinateRuleAvailable && document.activeElement !== quoteRuleRowStartInput && document.activeElement !== quoteRuleRowEndInput) {
    const defaultRow = String(selectedSourceRow() + 1);
    const previousStart = quoteRuleRowStartInput.dataset.syncedValue || "";
    const previousEnd = quoteRuleRowEndInput.dataset.syncedValue || "";
    const stillUsingDefault =
      quoteRuleRowStartInput.value === "" ||
      quoteRuleRowEndInput.value === "" ||
      (quoteRuleRowStartInput.value === previousStart && quoteRuleRowEndInput.value === previousEnd);
    if (stillUsingDefault) {
      quoteRuleRowStartInput.value = defaultRow;
      quoteRuleRowEndInput.value = defaultRow;
      quoteRuleRowStartInput.dataset.syncedValue = defaultRow;
      quoteRuleRowEndInput.dataset.syncedValue = defaultRow;
    }
  }
  if (coordinateRuleAvailable && document.activeElement !== quoteRuleColumnInput) {
    const defaultColumn = quoteRuleColumnLabel(sourceColFor(selected.col));
    const previousDefaultColumn = quoteRuleColumnInput.dataset.syncedValue || "";
    if (quoteRuleColumnInput.value === "" || quoteRuleColumnInput.value === previousDefaultColumn) {
      quoteRuleColumnInput.value = defaultColumn;
      quoteRuleColumnInput.dataset.syncedValue = defaultColumn;
    }
  }
  quoteRuleRowInput.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleColInput.disabled = readOnly || !coordinateRuleAvailable;
  quoteRuleActionSelect.disabled = readOnly || !coordinateRuleAvailable;
  applyQuoteCoordinateRuleButton.disabled = readOnly || !coordinateRuleAvailable;
  if (coordinateRuleAvailable && document.activeElement !== quoteRuleRowInput && document.activeElement !== quoteRuleColInput) {
    const defaultRow = String(selectedSourceRow() + 1);
    const defaultCol = String(sourceColFor(selected.col) + 1);
    const previousDefaultRow = quoteRuleRowInput.dataset.syncedValue || "";
    const previousDefaultCol = quoteRuleColInput.dataset.syncedValue || "";
    const stillUsingDefault =
      quoteRuleRowInput.value === "" ||
      quoteRuleColInput.value === "" ||
      (quoteRuleRowInput.value === previousDefaultRow && quoteRuleColInput.value === previousDefaultCol);
    if (stillUsingDefault) {
      quoteRuleRowInput.value = defaultRow;
      quoteRuleColInput.value = defaultCol;
      quoteRuleRowInput.dataset.syncedValue = defaultRow;
      quoteRuleColInput.dataset.syncedValue = defaultCol;
    }
  }
  if (resetQuotePolicyButton) {
    resetQuotePolicyButton.disabled = !documentId || readOnly || !quoteOnlyPolicyModified();
  }
  renderQuoteRuleList();
}

function syncLineEndingControls() {
  if (!lineEndingSelect || !applyLineEndingButton) return;
  const value = lineEndingValue();
  if (!lineEndingDraftDirty) {
    lineEndingSelect.value = value;
  }
  lineEndingSelect.disabled = !documentId || readOnly;
  applyLineEndingButton.disabled = !documentId || readOnly;
  document.body.dataset.lineEnding = value;
  document.body.dataset.lineEndingDraft = lineEndingDraftDirty ? normalizeLineEndingName(lineEndingSelect.value) : "";
}

function quotePolicySettingsModified() {
  if (delimiterDraftDirty || encodingDraftDirty || lineEndingDraftDirty || quotePolicyDraftDirty) return true;
  return quoteOnlyPolicyModified();
}

function quoteOnlyPolicyModified() {
  if (quotePolicyDraftDirty) return true;
  if (!documentId) return false;
  const settings = currentQuoteRuleSettings();
  return (
    settings.quoteMode !== "minimal" ||
    settings.quoteChar !== "double" ||
    normalizeQuoteEscapeValue(settings.quoteEscape) !== "double" ||
    settings.quoteHeader ||
    settings.quoteBody ||
    settings.quoteColumns.length > 0 ||
    settings.quoteCells.length > 0 ||
    settings.noQuoteCells.length > 0
  );
}

function documentFormatSummary() {
  return `${encodingLabel()} · ${delimiterLabel(metadata.dialect)} · ${lineEndingLabel()} · ${quotePolicyLabel(metadata.dialect)} · ${headerModeLabel(metadata.dialect)}`;
}

function currentFormatPreset() {
  if (!metadata?.dialect) return null;
  return {
    id: "",
    name: documentFormatSummary(),
    delimiter: delimiterChar(metadata.dialect),
    delimiterName: metadata.dialect.delimiterName || delimiterLabel(metadata.dialect).toLowerCase(),
    hasHeader: Boolean(metadata.dialect.hasHeader),
    headerColumns: headerColumnsValue(metadata.dialect),
    encodingName: metadata.encoding?.name ?? "utf-8",
    lineEndingName: lineEndingValue(metadata.lineEnding),
    quoteChar: quoteCharValue(metadata.dialect),
    quoteMode: quoteModeValue(metadata.dialect),
    quoteEscape: quoteEscapeValue(metadata.dialect),
    rules: [],
  };
}

function formatPresets() {
  try {
    const parsed = JSON.parse(localStorage.getItem(formatPresetsStorageKey) || "[]");
    if (!Array.isArray(parsed)) return [];
    return parsed
      .filter((preset) => preset && typeof preset === "object" && typeof preset.id === "string")
      .map((preset) => ({
        ...preset,
        headerColumns: normalizeHeaderColumns(preset.headerColumns),
        quoteChar: preset.quoteChar || "double",
        quoteMode: preset.quoteMode || "minimal",
        quoteEscape: normalizeQuoteEscapeValue(preset.quoteEscape),
        lineEndingName: normalizeLineEndingName(preset.lineEndingName),
        rules: normalizeFormatRules(preset.rules),
      }))
      .slice(0, formatPresetsLimit);
  } catch {
    return [];
  }
}

function saveFormatPresets(presets) {
  localStorage.setItem(formatPresetsStorageKey, JSON.stringify(presets.slice(0, formatPresetsLimit)));
  document.body.dataset.formatPresets = String(presets.length);
  renderFormatPresetSelect();
  renderSettingsPanel();
}

function renderFormatPresetSelect() {
  const selectedValue = formatPresetSelect.value;
  const presets = formatPresets();
  const options = [new Option("Preset", "")];
  for (const preset of presets) {
    const ruleCount = preset.rules?.length ?? 0;
    const label = ruleCount > 0 ? `${preset.name || preset.id} (${ruleCount})` : preset.name || preset.id;
    options.push(new Option(label, preset.id));
  }
  formatPresetSelect.replaceChildren(...options);
  if (presets.some((preset) => preset.id === selectedValue)) {
    formatPresetSelect.value = selectedValue;
  }
  formatPresetSelect.disabled = !documentId || presets.length === 0;
  applyFormatPresetButton.disabled = !documentId || readOnly || !formatPresetSelect.value;
  document.body.dataset.formatPresets = String(presets.length);
}

function saveCurrentFormatPreset(rulePattern = "") {
  const preset = currentFormatPreset();
  if (!preset) {
    status.textContent = "Open a file before saving a format preset.";
    return null;
  }
  const id = formatPresetID(preset);
  const existing = formatPresets().find((item) => item.id === id);
  const rules = normalizeFormatRules(existing?.rules);
  let ruleKey = "";
  if (rulePattern) {
    const normalized = normalizeFormatRule({ pattern: rulePattern });
    ruleKey = normalized ? formatRuleKey(normalized) : "";
    if (normalized && !rules.some((rule) => formatRuleKey(rule) === formatRuleKey(normalized))) {
      rules.push(normalized);
    }
  }
  const saved = { ...preset, id, name: formatPresetName(preset), rules };
  const otherPresets = formatPresets()
    .filter((item) => item.id !== id)
    .map((item) => (
      ruleKey
        ? { ...item, rules: normalizeFormatRules(item.rules).filter((rule) => formatRuleKey(rule) !== ruleKey) }
        : item
    ));
  const presets = [saved, ...otherPresets].slice(0, formatPresetsLimit);
  saveFormatPresets(presets);
  formatPresetSelect.value = id;
  document.body.dataset.formatPresetPreview = formatPresetPreview(saved);
  status.textContent = rulePattern ? `Saved format rule: ${rulePattern} -> ${saved.name}` : `Saved format preset: ${saved.name}`;
  updateCommandState();
  return saved;
}

function saveCurrentFormatRule() {
  const pattern = formatRuleInput.value.trim() || defaultFormatRulePattern();
  if (!pattern) {
    status.textContent = "Enter a rule pattern after opening a file.";
    return;
  }
  const saved = saveCurrentFormatPreset(pattern);
  if (!saved) return;
  formatRuleInput.value = pattern;
  document.body.dataset.formatRulePreview = `${formatRuleTarget(pattern)}:${formatRulePattern(pattern)}`;
}

async function applySelectedFormatPreset() {
  if (!ensureEditable("Apply Format Preset")) return;
  const preset = formatPresets().find((item) => item.id === formatPresetSelect.value);
  if (!preset) {
    status.textContent = "Choose a format preset first.";
    return;
  }
  const path = metadata?.path || pathInput.value.trim();
  if (!path) {
    status.textContent = "Open a local file before applying a format preset.";
    return;
  }
  document.body.dataset.formatPresetPreview = formatPresetPreview(preset);
  await openPath(path, { preset });
}

function exportFormatPresets() {
  const presets = formatPresets();
  const ruleCount = presets.reduce((sum, preset) => sum + normalizeFormatRules(preset.rules).length, 0);
  const payload = {
    kind: formatPresetExportKind,
    version: 1,
    exportedAt: new Date().toISOString(),
    presets: presets.map((preset) => ({
      id: preset.id,
      name: preset.name,
      delimiter: preset.delimiter,
      delimiterName: preset.delimiterName,
      hasHeader: Boolean(preset.hasHeader),
      headerColumns: normalizeHeaderColumns(preset.headerColumns),
      encodingName: preset.encodingName || "utf-8",
      lineEndingName: normalizeLineEndingName(preset.lineEndingName),
      quoteChar: preset.quoteChar || "double",
      quoteMode: preset.quoteMode || "minimal",
      quoteEscape: normalizeQuoteEscapeValue(preset.quoteEscape),
      rules: normalizeFormatRules(preset.rules),
    })),
  };
  const text = `${JSON.stringify(payload, null, 2)}\n`;
  window.__opencsvLastFormatPresetExportText = text;
  document.body.dataset.formatPresetExport = `${presets.length}:${ruleCount}`;
  document.body.dataset.formatPresetExportPreview = text.slice(0, 500);
  triggerDownload(text, "opencsv-format-presets.json", "application/json;charset=utf-8");
  status.textContent = `Exported ${formatCount(presets.length)} format presets.`;
}

function showFormatPresetImportPanel() {
  openSettingsPanel();
  formatPresetImportPanel.hidden = false;
  formatPresetImportText.value = "";
  document.body.dataset.formatPresetImportPanel = "open";
  syncFormatPresetImportControls();
  formatPresetImportText.focus();
}

function cancelFormatPresetImport() {
  formatPresetImportPanel.hidden = true;
  formatPresetImportText.value = "";
  document.body.dataset.formatPresetImportPanel = "closed";
  syncFormatPresetImportControls();
}

function syncFormatPresetImportControls() {
  if (!formatPresetImportText) return;
  applyFormatPresetImportButton.disabled = formatPresetImportText.value.trim() === "";
}

function applyFormatPresetImport() {
  try {
    const imported = parseFormatPresetImport(formatPresetImportText.value);
    const next = mergeImportedFormatPresets(formatPresets(), imported);
    saveFormatPresets(next);
    formatPresetImportPanel.hidden = true;
    formatPresetImportText.value = "";
    const importedRules = imported.reduce((sum, preset) => sum + normalizeFormatRules(preset.rules).length, 0);
    document.body.dataset.formatPresetImport = `${imported.length}:${importedRules}`;
    document.body.dataset.formatPresetImportPanel = "closed";
    status.textContent = `Imported ${formatCount(imported.length)} format presets.`;
    updateCommandState();
  } catch (error) {
    status.textContent = `Import failed: ${error.message}`;
    document.body.dataset.formatPresetImport = "error";
  }
  syncFormatPresetImportControls();
}

function parseFormatPresetImport(text) {
  let parsed;
  try {
    parsed = JSON.parse(text);
  } catch {
    throw new Error("invalid JSON");
  }
  const rawPresets = Array.isArray(parsed) ? parsed : parsed?.presets;
  if (!Array.isArray(rawPresets)) {
    throw new Error("missing presets array");
  }
  const imported = [];
  for (const raw of rawPresets) {
    const preset = normalizeImportedFormatPreset(raw);
    if (preset && !imported.some((item) => item.id === preset.id)) imported.push(preset);
  }
  if (!imported.length) throw new Error("no valid presets");
  return imported.slice(0, formatPresetsLimit);
}

function normalizeImportedFormatPreset(raw) {
  if (!raw || typeof raw !== "object") return null;
  const delimiter = normalizeImportedDelimiter(raw.delimiter, raw.delimiterName);
  const quoteChar = ["double", "single", "none"].includes(raw.quoteChar) ? raw.quoteChar : "double";
  const quoteMode = ["minimal", "always", "never", "manual"].includes(raw.quoteMode) ? raw.quoteMode : "minimal";
  const quoteEscape = normalizeQuoteEscapeValue(raw.quoteEscape);
  const preset = {
    id: "",
    name: "",
    delimiter,
    delimiterName: normalizeImportedDelimiterName(raw.delimiterName, delimiter),
    hasHeader: Boolean(raw.hasHeader),
    headerColumns: normalizeHeaderColumns(raw.headerColumns),
    encodingName: typeof raw.encodingName === "string" && raw.encodingName.trim() ? raw.encodingName.trim() : "utf-8",
    lineEndingName: normalizeLineEndingName(raw.lineEndingName),
    quoteChar,
    quoteMode,
    quoteEscape,
    rules: normalizeFormatRules(raw.rules),
  };
  preset.id = typeof raw.id === "string" && raw.id.trim() ? raw.id.trim().slice(0, 160) : formatPresetID(preset);
  preset.name = typeof raw.name === "string" && raw.name.trim() ? raw.name.trim().slice(0, 120) : formatPresetName(preset);
  return preset;
}

function normalizeImportedDelimiter(value, name) {
  if (value === "\t" || value === "\\t") return "\t";
  if (typeof value === "string" && value.length > 0) return value[0];
  const lower = typeof name === "string" ? name.trim().toLowerCase() : "";
  if (lower === "tab" || lower === "tsv") return "\t";
  if (lower === "semicolon") return ";";
  if (lower === "pipe") return "|";
  if (lower === "space") return " ";
  return ",";
}

function normalizeImportedDelimiterName(value, delimiter) {
  if (typeof value === "string" && value.trim()) return value.trim();
  if (delimiter === "\t") return "tab";
  if (delimiter === ";") return "semicolon";
  if (delimiter === "|") return "pipe";
  if (delimiter === " ") return "space";
  if (delimiter === ",") return "comma";
  return delimiter;
}

function mergeImportedFormatPresets(existing, imported) {
  const importedIDs = new Set(imported.map((preset) => preset.id));
  return [
    ...imported,
    ...existing.filter((preset) => !importedIDs.has(preset.id)),
  ].slice(0, formatPresetsLimit);
}

function openSettingsPanel() {
  setViewMode("settings");
}

function renderSettingsPanel() {
  if (!settingsPresetList) return;
  syncLocaleControls();
  syncThemeControls();
  syncAutoRefreshControls();
  syncUpdateCheckControls();
  syncFilterHistoryControls();
  syncCustomLinkProtocolControls();
  syncDelimiterControls();
  syncEncodingControls();
  syncQuotePolicyControls();
  syncEditorAppearanceControls();
  const presets = formatPresets();
  const filters = savedFilters();
  const recent = recentFilters();
  const sorts = savedSorts();
  const query = settingsSearch.value.trim().toLowerCase();
  const visiblePresets = query
    ? presets.filter((preset) => formatPresetSearchText(preset).includes(query))
    : presets;
  const sectionState = settingsSectionModifiedState({ presets, filters, recent, sorts });
  settingsModifiedBadge.textContent = sectionState.modified.length ? t("settings.modified") : t("settings.defaults");
  settingsFormatSummary.textContent = `${presets.length} presets, ${presets.reduce((sum, preset) => sum + normalizeFormatRules(preset.rules).length, 0)} rules`;
  settingsViewsSummary.textContent = `${filters.length} filters, ${recent.length} recent, ${sorts.length} sorts · history ${filterHistoryLimitLabel().toLowerCase()}`;
  resetFormatPresetsButton.disabled = presets.length === 0;
  if (resetSavedViewsButton) {
    resetSavedViewsButton.disabled = !sectionState.state["Saved filters and sorts"];
  }
  syncSettingsSectionBadges(sectionState);
  renderSavedViewSettingsList(filters, sorts, query);
  settingsPresetList.replaceChildren();
  document.body.dataset.settingsPresetCount = String(presets.length);
  document.body.dataset.settingsVisiblePresetCount = String(visiblePresets.length);
  if (visiblePresets.length === 0) {
    settingsPresetList.append(emptySettingsNode(presets.length ? "No matching presets" : "No format presets saved"));
  } else {
    for (const preset of visiblePresets) {
      settingsPresetList.append(formatPresetSettingsRow(preset));
    }
  }
  filterSettingsSections(query);
}

function settingsSectionModifiedState({ presets = formatPresets(), filters = savedFilters(), recent = recentFilters(), sorts = savedSorts() } = {}) {
  const state = {
    Language: localePreferenceModified(),
    "Saved filters and sorts": savedViewsSettingsModified({ filters, recent, sorts }),
    Theme: themePreferenceModified(),
    "External changes": autoRefreshModeModified(),
    Updates: updateCheckModeModified() || Boolean(updateCheckState?.checkedAt),
    "Link protocols": customLinkProtocolsModified(),
    "Quote policy": quotePolicySettingsModified(),
    "Editor appearance": editorAppearanceModified(),
    "File format presets": presets.length > 0,
  };
  const modified = Object.entries(state)
    .filter(([, value]) => value)
    .map(([label]) => label);
  return { state, modified };
}

function syncSettingsSectionBadges(sectionState = settingsSectionModifiedState()) {
  if (!settingsPanel) return;
  const modifiedLabels = [];
  for (const section of settingsPanel.querySelectorAll(".settings-section")) {
    const label = section.getAttribute("aria-label") || "";
    const modified = Boolean(sectionState.state[label]);
    section.dataset.modified = String(modified);
    let badge = section.querySelector(":scope > header .settings-section-state");
    if (!badge) {
      badge = document.createElement("span");
      badge.className = "settings-section-state";
      const actions = section.querySelector(":scope > header .settings-section-actions");
      section.querySelector(":scope > header")?.insertBefore(badge, actions || null);
    }
    badge.textContent = modified ? t("settings.modified") : t("settings.defaults");
    badge.dataset.state = modified ? "modified" : "default";
    badge.title = `${label}: ${badge.textContent}`;
    if (modified) modifiedLabels.push(label);
  }
  document.body.dataset.settingsModifiedSections = modifiedLabels.join("|");
  document.body.dataset.settingsModifiedSectionCount = String(modifiedLabels.length);
}

function renderSavedViewSettingsList(filters, sorts, query) {
  if (!settingsSavedViewList) return;
  const items = [
    ...filters.map((item) => ({ kind: "filter", item })),
    ...sorts.map((item) => ({ kind: "sort", item })),
  ];
  const visible = query
    ? items.filter(({ kind, item }) => savedViewSearchText(kind, item).includes(query))
    : items;
  settingsSavedViewList.replaceChildren();
  document.body.dataset.settingsSavedFilterCount = String(filters.length);
  document.body.dataset.settingsSavedSortCount = String(sorts.length);
  document.body.dataset.settingsVisibleSavedViewCount = String(visible.length);
  if (visible.length === 0) {
    settingsSavedViewList.append(emptySettingsNode(items.length ? "No matching saved views" : "No saved filters or sorts"));
    return;
  }
  for (const entry of visible) {
    settingsSavedViewList.append(savedViewSettingsRow(entry.kind, entry.item));
  }
}

function filterSettingsSections(query) {
  if (!settingsPanel) return;
  const sections = Array.from(settingsPanel.querySelectorAll(".settings-section"));
  const visibleLabels = [];
  for (const section of sections) {
    const label = section.getAttribute("aria-label") || "";
    const text = `${label} ${section.textContent}`.toLowerCase();
    const visible = !query || text.includes(query);
    section.hidden = !visible;
    if (visible) visibleLabels.push(label);
  }
  document.body.dataset.settingsVisibleSectionCount = String(visibleLabels.length);
  document.body.dataset.settingsVisibleSections = visibleLabels.join("|");
  if (settingsEmpty) {
    settingsEmpty.hidden = visibleLabels.length > 0;
    settingsEmpty.textContent = query ? `No settings match "${settingsSearch.value.trim()}".` : "No matching settings";
  }
}

function savedViewSearchText(kind, item) {
  return [
    kind,
    item.name,
    item.display,
    item.expression,
    ...(item.keys ?? []).flatMap((key) => [key.name, key.mode, key.descending ? "descending" : "ascending"]),
  ].filter(Boolean).join(" ").toLowerCase();
}

function savedViewSettingsRow(kind, item) {
  const row = document.createElement("section");
  row.className = "settings-preset-row";
  row.dataset.savedViewKind = kind;
  row.dataset.savedViewId = item.id;

  const header = document.createElement("header");
  const title = document.createElement("strong");
  title.textContent = item.name;
  const meta = document.createElement("span");
  meta.textContent = kind === "filter" ? `Filter · ${item.expression}` : `Sort · ${formatStoredSortKeys(item.keys)}`;
  const actions = document.createElement("div");
  actions.className = "settings-preset-actions";
  const applyButton = settingsActionButton("Apply");
  applyButton.disabled = !documentId;
  applyButton.addEventListener("click", () => {
    if (kind === "filter") {
      applySavedFilter(item.id);
    } else {
      applySavedSort(item.id);
    }
  });
  const deleteButton = settingsActionButton("Delete");
  deleteButton.addEventListener("click", () => {
    if (kind === "filter") {
      deleteSavedFilter(item.id);
    } else {
      deleteSavedSort(item.id);
    }
    renderSettingsPanel();
  });
  actions.append(applyButton, deleteButton);
  header.append(title, meta, actions);
  row.append(header);
  return row;
}

function formatPresetSearchText(preset) {
  return [
    preset.name,
    preset.id,
    preset.delimiterName,
    preset.encodingName,
    preset.lineEndingName,
    preset.quoteMode,
    preset.quoteChar,
    preset.quoteEscape,
    ...(preset.rules ?? []).map((rule) => rule.pattern),
  ].filter(Boolean).join(" ").toLowerCase();
}

function emptySettingsNode(message) {
  const node = document.createElement("div");
  node.className = "settings-empty";
  node.textContent = message;
  return node;
}

function formatPresetSettingsRow(preset) {
  const row = document.createElement("section");
  row.className = "settings-preset-row";
  row.dataset.presetId = preset.id;

  const header = document.createElement("header");
  const title = document.createElement("strong");
  title.textContent = preset.name || preset.id;
  const presetRules = normalizeFormatRules(preset.rules);
  const meta = document.createElement("span");
  meta.textContent = `${formatPresetPreview(preset)} · ${presetRules.length} rules`;
  const actions = document.createElement("div");
  actions.className = "settings-preset-actions";
  const presets = formatPresets();
  const presetIndex = presets.findIndex((item) => item.id === preset.id);
  const moveUpButton = settingsActionButton("Up");
  moveUpButton.disabled = presetIndex <= 0;
  moveUpButton.setAttribute("aria-label", `Move preset ${preset.name || preset.id} up`);
  moveUpButton.addEventListener("click", () => moveFormatPreset(preset.id, -1));
  const moveDownButton = settingsActionButton("Down");
  moveDownButton.disabled = presetIndex < 0 || presetIndex >= presets.length - 1;
  moveDownButton.setAttribute("aria-label", `Move preset ${preset.name || preset.id} down`);
  moveDownButton.addEventListener("click", () => moveFormatPreset(preset.id, 1));
  const applyButton = settingsActionButton("Apply");
  applyButton.disabled = !documentId;
  applyButton.addEventListener("click", async () => {
    formatPresetSelect.value = preset.id;
    await applyFormatPresetToCurrentPath(preset);
  });
  const deleteButton = settingsActionButton("Delete");
  deleteButton.addEventListener("click", () => deleteFormatPreset(preset.id));
  actions.append(moveUpButton, moveDownButton, applyButton, deleteButton);
  header.append(title, meta, actions);

  const editForm = formatPresetEditForm(preset);

  const rules = document.createElement("div");
  rules.className = "settings-rule-list";
  if (presetRules.length === 0) {
    const empty = document.createElement("span");
    empty.className = "settings-rule-empty";
    empty.textContent = "No automatic rules";
    rules.append(empty);
  } else {
    for (const [ruleIndex, rule] of presetRules.entries()) {
      const item = document.createElement("div");
      item.className = "settings-rule-row";
      item.dataset.rulePattern = rule.pattern;
      item.dataset.ruleIndex = String(ruleIndex);
      const label = document.createElement("code");
      label.textContent = rule.pattern;
      const moveRuleUp = settingsActionButton("Up");
      moveRuleUp.disabled = ruleIndex === 0;
      moveRuleUp.setAttribute("aria-label", `Move rule ${rule.pattern} up`);
      moveRuleUp.addEventListener("click", () => moveFormatRule(preset.id, rule.pattern, -1));
      const moveRuleDown = settingsActionButton("Down");
      moveRuleDown.disabled = ruleIndex >= presetRules.length - 1;
      moveRuleDown.setAttribute("aria-label", `Move rule ${rule.pattern} down`);
      moveRuleDown.addEventListener("click", () => moveFormatRule(preset.id, rule.pattern, 1));
      const remove = settingsActionButton("Remove");
      remove.addEventListener("click", () => deleteFormatRule(preset.id, rule.pattern));
      item.append(label, moveRuleUp, moveRuleDown, remove);
      rules.append(item);
    }
  }

  const addRule = document.createElement("form");
  addRule.className = "settings-add-rule";
  const input = document.createElement("input");
  input.type = "text";
  input.placeholder = "name:*.csv, path:/exports/*, content:marker";
  input.setAttribute("aria-label", `Add rule for ${preset.name || preset.id}`);
  const add = settingsActionButton("Add Rule");
  add.type = "submit";
  addRule.append(input, add);
  addRule.addEventListener("submit", (event) => {
    event.preventDefault();
    addFormatRuleToPreset(preset.id, input.value);
  });

  row.append(header, editForm, rules, addRule);
  return row;
}

function formatPresetEditForm(preset) {
  const form = document.createElement("form");
  form.className = "settings-preset-edit";
  form.dataset.presetId = preset.id;

  const presetLabel = preset.name || preset.id;
  const nameInput = document.createElement("input");
  nameInput.name = "presetName";
  nameInput.type = "text";
  nameInput.maxLength = 120;
  nameInput.value = preset.name || formatPresetName(preset);
  nameInput.setAttribute("aria-label", `Preset name for ${presetLabel}`);

  const delimiterSelect = settingsSelectControl("delimiter", [
    ["comma", "Comma"],
    ["tab", "Tab"],
    ["semicolon", "Semicolon"],
    ["pipe", "Pipe"],
    ["space", "Space"],
    ["custom", "Custom"],
  ], delimiterOptionValue(preset));
  delimiterSelect.setAttribute("aria-label", `Delimiter for ${presetLabel}`);

  const customDelimiter = document.createElement("input");
  customDelimiter.name = "customDelimiter";
  customDelimiter.type = "text";
  customDelimiter.maxLength = 1;
  customDelimiter.value = delimiterCustomValue(preset);
  customDelimiter.disabled = delimiterSelect.value !== "custom";
  customDelimiter.setAttribute("aria-label", `Custom delimiter for ${presetLabel}`);
  delimiterSelect.addEventListener("change", () => {
    const custom = delimiterSelect.value === "custom";
    customDelimiter.disabled = !custom;
    if (!custom) customDelimiter.value = "";
  });

  const headerRow = document.createElement("input");
  headerRow.name = "hasHeader";
  headerRow.type = "checkbox";
  headerRow.checked = Boolean(preset.hasHeader);
  headerRow.setAttribute("aria-label", `Header row for ${presetLabel}`);

  const headerColumns = document.createElement("input");
  headerColumns.name = "headerColumns";
  headerColumns.type = "number";
  headerColumns.min = "0";
  headerColumns.max = "64";
  headerColumns.step = "1";
  headerColumns.value = String(normalizeHeaderColumns(preset.headerColumns));
  headerColumns.setAttribute("aria-label", `Header columns for ${presetLabel}`);

  const encodingSelectEdit = settingsSelectControl("encodingName", [
    ["utf-8", "UTF-8"],
    ["utf-8-bom", "UTF-8 BOM"],
    ["utf-16le", "UTF-16 LE"],
    ["utf-16le-bom", "UTF-16 LE BOM"],
    ["utf-16be", "UTF-16 BE"],
    ["utf-16be-bom", "UTF-16 BE BOM"],
    ["shift-jis", "Shift-JIS"],
    ["gbk", "GBK"],
    ["windows-1252", "Windows-1252"],
    ["legacy-single-byte", "Legacy byte"],
  ], preset.encodingName || "utf-8");
  encodingSelectEdit.setAttribute("aria-label", `Encoding for ${presetLabel}`);

  const lineEndingSelectEdit = settingsSelectControl("lineEndingName", [
    ["lf", "LF"],
    ["crlf", "CRLF"],
    ["cr", "CR"],
  ], normalizeLineEndingName(preset.lineEndingName));
  lineEndingSelectEdit.setAttribute("aria-label", `Line ending for ${presetLabel}`);

  const quoteModeSelectEdit = settingsSelectControl("quoteMode", [
    ["minimal", "Minimal"],
    ["always", "Always"],
    ["never", "Never"],
    ["manual", "Manual"],
  ], quoteModeValue(preset));
  quoteModeSelectEdit.setAttribute("aria-label", `Quote mode for ${presetLabel}`);

  const quoteCharSelectEdit = settingsSelectControl("quoteChar", [
    ["double", "Double quote"],
    ["single", "Single quote"],
    ["none", "None"],
  ], quoteCharValue(preset));
  quoteCharSelectEdit.setAttribute("aria-label", `Quote character for ${presetLabel}`);

  const quoteEscapeSelectEdit = settingsSelectControl("quoteEscape", [
    ["double", "Double quote"],
    ["backslash", "Backslash"],
  ], quoteEscapeValue(preset));
  quoteEscapeSelectEdit.setAttribute("aria-label", `Quote escape for ${presetLabel}`);

  const save = settingsActionButton("Save Edit");
  save.type = "submit";
  save.setAttribute("aria-label", `Save preset edit for ${presetLabel}`);

  form.append(
    settingsControlLabel("Name", nameInput),
    settingsControlLabel("Delimiter", delimiterSelect),
    settingsControlLabel("Custom", customDelimiter),
    settingsCheckboxLabel("Header Row", headerRow),
    settingsControlLabel("Header Cols", headerColumns),
    settingsControlLabel("Encoding", encodingSelectEdit),
    settingsControlLabel("EOL", lineEndingSelectEdit),
    settingsControlLabel("Quote Mode", quoteModeSelectEdit),
    settingsControlLabel("Quote Char", quoteCharSelectEdit),
    settingsControlLabel("Quote Escape", quoteEscapeSelectEdit),
    save,
  );
  form.addEventListener("submit", (event) => {
    event.preventDefault();
    saveEditedFormatPreset(preset.id, form);
  });
  return form;
}

function settingsControlLabel(labelText, control) {
  const label = document.createElement("label");
  const span = document.createElement("span");
  span.textContent = labelText;
  label.append(span, control);
  return label;
}

function settingsCheckboxLabel(labelText, control) {
  const label = document.createElement("label");
  label.className = "settings-preset-checkbox";
  const span = document.createElement("span");
  span.textContent = labelText;
  label.append(control, span);
  return label;
}

function settingsSelectControl(name, options, selectedValue) {
  const select = document.createElement("select");
  select.name = name;
  for (const [value, label] of options) {
    select.append(new Option(label, value));
  }
  if (options.some(([value]) => value === selectedValue)) {
    select.value = selectedValue;
  } else if (options.length) {
    select.value = options[0][0];
  }
  return select;
}

function settingsActionButton(label) {
  const button = document.createElement("button");
  button.type = "button";
  button.textContent = label;
  return button;
}

async function applyFormatPresetToCurrentPath(preset) {
  const path = metadata?.path || pathInput.value.trim();
  if (!path) {
    status.textContent = "Open a local file before applying a format preset.";
    return;
  }
  document.body.dataset.formatPresetPreview = formatPresetPreview(preset);
  await openPath(path, { preset });
}

function deleteFormatPreset(presetID) {
  const presets = formatPresets().filter((preset) => preset.id !== presetID);
  saveFormatPresets(presets);
  if (formatPresetSelect.value === presetID) formatPresetSelect.value = "";
  document.body.dataset.formatPresetOrder = presets.map((preset) => preset.id).join("|");
  status.textContent = "Deleted format preset.";
  updateCommandState();
}

function moveFormatPreset(presetID, delta) {
  const presets = formatPresets();
  const index = presets.findIndex((preset) => preset.id === presetID);
  const nextIndex = index + delta;
  if (index < 0 || nextIndex < 0 || nextIndex >= presets.length) {
    status.textContent = "Format preset order unchanged.";
    return;
  }
  const [preset] = presets.splice(index, 1);
  presets.splice(nextIndex, 0, preset);
  saveFormatPresets(presets);
  formatPresetSelect.value = presetID;
  document.body.dataset.formatPresetOrder = presets.map((item) => item.id).join("|");
  status.textContent = `Moved format preset ${nextIndex < index ? "up" : "down"}.`;
  updateCommandState();
}

function saveEditedFormatPreset(presetID, form) {
  let updated = null;
  let found = false;
  const presets = formatPresets().map((preset) => {
    if (preset.id !== presetID) return preset;
    found = true;
    const draft = editedFormatPresetFromForm(preset, form);
    if (!draft) return preset;
    updated = draft;
    return draft;
  });
  if (!found) {
    status.textContent = "Format preset not found.";
    return;
  }
  if (!updated) return;
  saveFormatPresets(presets);
  formatPresetSelect.value = presetID;
  const signature = formatPresetPreview(updated);
  document.body.dataset.formatPresetEdit = `${presetID}:${signature}`;
  document.body.dataset.formatPresetPreview = signature;
  status.textContent = `Updated format preset: ${updated.name || formatPresetName(updated)}`;
  updateCommandState();
}

function editedFormatPresetFromForm(preset, form) {
  const delimiterOption = form.elements.delimiter?.value || "comma";
  let delimiterInfo;
  if (delimiterOption === "custom") {
    const custom = Array.from(form.elements.customDelimiter?.value || "")[0];
    if (!custom) {
      status.textContent = "Enter a custom delimiter before saving the preset.";
      return null;
    }
    delimiterInfo = { delimiter: custom, delimiterName: custom };
  } else {
    delimiterInfo = normalizeDelimiterOverride(delimiterOption);
  }
  const quoteChar = ["double", "single", "none"].includes(form.elements.quoteChar?.value)
    ? form.elements.quoteChar.value
    : "double";
  const quoteMode = ["minimal", "always", "never", "manual"].includes(form.elements.quoteMode?.value)
    ? form.elements.quoteMode.value
    : "minimal";
  const quoteEscape = normalizeQuoteEscapeValue(form.elements.quoteEscape?.value);
  const draft = {
    ...preset,
    delimiter: delimiterInfo.delimiter,
    delimiterName: delimiterInfo.delimiterName,
    hasHeader: Boolean(form.elements.hasHeader?.checked),
    headerColumns: normalizeHeaderColumns(form.elements.headerColumns?.value),
    encodingName: form.elements.encodingName?.value || "utf-8",
    lineEndingName: normalizeLineEndingName(form.elements.lineEndingName?.value),
    quoteChar,
    quoteMode,
    quoteEscape,
    rules: normalizeFormatRules(preset.rules),
  };
  const name = String(form.elements.presetName?.value || "").trim().slice(0, 120);
  draft.name = name || formatPresetName(draft);
  return draft;
}

function deleteFormatRule(presetID, pattern) {
  const key = pattern.toLowerCase();
  const presets = formatPresets().map((preset) => (
    preset.id === presetID
      ? { ...preset, rules: normalizeFormatRules(preset.rules).filter((rule) => rule.pattern.toLowerCase() !== key) }
      : preset
  ));
  saveFormatPresets(presets);
  document.body.dataset.formatRulePreview = `removed:${pattern}`;
  document.body.dataset.formatRuleOrder = formatRuleOrderSignature(presets, presetID);
  status.textContent = `Removed format rule: ${pattern}`;
}

function moveFormatRule(presetID, pattern, delta) {
  const key = pattern.toLowerCase();
  let moved = false;
  const presets = formatPresets().map((preset) => {
    const rules = normalizeFormatRules(preset.rules);
    if (preset.id !== presetID) return { ...preset, rules };
    const index = rules.findIndex((rule) => formatRuleKey(rule) === key);
    const nextIndex = index + delta;
    if (index < 0 || nextIndex < 0 || nextIndex >= rules.length) return { ...preset, rules };
    const [rule] = rules.splice(index, 1);
    rules.splice(nextIndex, 0, rule);
    moved = true;
    return { ...preset, rules };
  });
  if (!moved) {
    status.textContent = "Format rule order unchanged.";
    return;
  }
  saveFormatPresets(presets);
  document.body.dataset.formatRuleOrder = formatRuleOrderSignature(presets, presetID);
  status.textContent = `Moved format rule ${delta < 0 ? "up" : "down"}.`;
  updateCommandState();
}

function formatRuleOrderSignature(presets, presetID) {
  const preset = presets.find((item) => item.id === presetID);
  return normalizeFormatRules(preset?.rules).map((rule) => rule.pattern).join("|");
}

function addFormatRuleToPreset(presetID, pattern) {
  const normalized = normalizeFormatRule({ pattern });
  if (!normalized) {
    status.textContent = "Enter a rule pattern first.";
    return;
  }
  const key = formatRuleKey(normalized);
  const presets = formatPresets().map((preset) => {
    const rules = normalizeFormatRules(preset.rules).filter((rule) => formatRuleKey(rule) !== key);
    if (preset.id !== presetID) return { ...preset, rules };
    return { ...preset, rules: [...rules, normalized].slice(0, 12) };
  });
  saveFormatPresets(presets);
  document.body.dataset.formatRulePreview = `${formatRuleTarget(normalized.pattern)}:${formatRulePattern(normalized.pattern)}`;
  document.body.dataset.formatRuleOrder = formatRuleOrderSignature(presets, presetID);
  status.textContent = `Added format rule: ${normalized.pattern}`;
}

function resetFormatPresets() {
  localStorage.removeItem(formatPresetsStorageKey);
  formatPresetSelect.value = "";
  document.body.dataset.formatPresets = "0";
  document.body.dataset.formatRulePreview = "reset";
  document.body.dataset.formatPresetPreview = "";
  document.body.dataset.formatPresetEdit = "";
  document.body.dataset.formatPresetExport = "";
  document.body.dataset.formatPresetImport = "";
  document.body.dataset.formatPresetImportPanel = "closed";
  formatPresetImportPanel.hidden = true;
  formatPresetImportText.value = "";
  syncFormatPresetImportControls();
  renderFormatPresetSelect();
  renderSettingsPanel();
  updateCommandState();
  status.textContent = "Reset format presets.";
}

async function applyDelimiter() {
  if (!documentId) {
    status.textContent = "Open a file before changing delimiter.";
    return;
  }
  if (readOnly) {
    status.textContent = "Delimiter blocked: document is read-only";
    return;
  }
  const { delimiter, delimiterName } = delimiterControlValue();
  const response = await fetch("/api/delimiter", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, delimiter, delimiterName }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Delimiter failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  delimiterDraftDirty = false;
  applyMutationPayload(payload, `Delimiter set to ${delimiterLabel({ delimiter, delimiterName })}.`);
  syncDelimiterControls();
  updateFormatStatus();
  updateInspector();
  status.textContent = `Delimiter set to ${delimiterLabel({ delimiter, delimiterName })}.`;
  document.body.dataset.delimiterPolicy = delimiterName;
}

async function applyEncoding() {
  if (!documentId) {
    status.textContent = "Open a file before changing encoding.";
    return;
  }
  if (readOnly) {
    status.textContent = "Encoding blocked: document is read-only";
    return;
  }
  const encodingName = encodingValue({ name: encodingSelect.value });
  const response = await fetch("/api/encoding", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, encodingName }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Encoding failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  encodingDraftDirty = false;
  applyMutationPayload(payload, `Encoding set to ${encodingLabel({ name: encodingName })}.`);
  syncEncodingControls();
  updateFormatStatus();
  updateInspector();
  status.textContent = `Encoding set to ${encodingLabel(metadata?.encoding)}.`;
  document.body.dataset.encodingPolicy = encodingValue(metadata?.encoding);
}

async function applyQuotePolicy() {
  if (!documentId) {
    status.textContent = "Open a file before changing quote policy.";
    return;
  }
  if (readOnly) {
    status.textContent = "Quote policy blocked: document is read-only";
    return;
  }
  const quoteMode = quoteModeSelect.value;
  const quoteChar = quoteCharSelect.value;
  const quoteEscape = quoteEscapeSelect.value;
  const response = await fetch("/api/quote-policy", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, quoteChar, quoteMode, quoteEscape }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Quote policy failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `Quote policy set to ${quoteMode} ${quoteChar} ${quoteEscape} escape.`);
  if (metadata?.dialect) {
    metadata.dialect.quoteMode = quoteMode;
    metadata.dialect.quoteChar = quoteChar === "single" ? "'" : quoteChar === "none" ? "" : '"';
    metadata.dialect.quoteEscape = quoteEscape;
  }
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  status.textContent = `Quote policy set to ${quoteMode} ${quoteChar} ${quoteEscape} escape.`;
  document.body.dataset.quotePolicy = `${quoteMode}:${quoteChar}:${quoteEscape}`;
}

async function resetQuotePolicySettings() {
  if (!documentId) {
    status.textContent = "Open a file before resetting quote policy.";
    return;
  }
  if (readOnly) {
    status.textContent = "Quote policy reset blocked: document is read-only";
    return;
  }
  const settings = {
    quoteChar: "double",
    quoteMode: "minimal",
    quoteEscape: "double",
    quoteHeader: false,
    quoteBody: false,
    quoteColumns: [],
    quoteCells: [],
    noQuoteCells: [],
  };
  const before = quoteRulesSignature(currentQuoteRuleSettings());
  quotePolicyDraftDirty = false;
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...settings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Quote policy reset failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quotePolicyReset = "error";
    syncQuotePolicyControls();
    renderSettingsPanel();
    return;
  }
  quoteRulesImportPanel.hidden = true;
  quoteRulesImportText.value = "";
  document.body.dataset.quoteRulesImportPanel = "closed";
  document.body.dataset.quotePolicyReset = `${before}->${quoteRulesSignature(settings)}`;
  applyMutationPayload(payload, "Quote policy reset.");
  syncQuoteRulesImportControls();
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  renderSettingsPanel();
  status.textContent = "Quote policy reset.";
}

async function applyLineEnding() {
  if (!documentId) {
    status.textContent = "Open a file before changing line ending.";
    return;
  }
  if (readOnly) {
    status.textContent = "Line ending blocked: document is read-only";
    return;
  }
  const lineEnding = normalizeLineEndingName(lineEndingSelect.value);
  const response = await fetch("/api/line-ending", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, lineEnding }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Line ending failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  lineEndingDraftDirty = false;
  applyMutationPayload(payload, `Line ending set to ${lineEnding.toUpperCase()}.`);
  if (metadata) {
    metadata.lineEnding = { name: lineEnding, sequence: lineEndingSequence(lineEnding), confidence: 1 };
  }
  syncLineEndingControls();
  updateFormatStatus();
  updateInspector();
  status.textContent = `Line ending set to ${lineEnding.toUpperCase()}.`;
  document.body.dataset.lineEndingPolicy = lineEnding;
}

async function toggleCellQuote() {
  if (!documentId) {
    status.textContent = "Open a file before changing cell quote state.";
    return;
  }
  if (readOnly) {
    status.textContent = "Cell quote blocked: document is read-only";
    return;
  }
  if (selected.row >= sourceRowCount) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  const quoted = !shouldQuote(valueAt(selected.row, selected.col), selected.row, selected.col);
  const row = selectedSourceRow();
  const col = sourceColFor(selected.col);
  const response = await fetch("/api/cell-quote", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, row, col, quoted }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Cell quote failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  applyMutationPayload(payload, quoted ? "Cell quote enabled." : "Cell quote disabled.");
  document.body.dataset.cellQuoteOverride = `${row}:${col}:${quoted ? "quote" : "no-quote"}`;
}

async function toggleQuoteRule(target, quoted = null) {
  if (!documentId) {
    status.textContent = "Open a file before changing quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = "Quote rule blocked: document is read-only";
    return;
  }
  const sourceCol = sourceColFor(selected.col);
  let nextQuoted = quoted;
  if (nextQuoted == null) {
    if (target === "header") nextQuoted = !Boolean(metadata?.dialect?.quoteHeader);
    if (target === "body") nextQuoted = !Boolean(metadata?.dialect?.quoteBody);
    if (target === "column") nextQuoted = !isSourceColumnQuoted(sourceCol);
  }
  const body = { documentId, target, quoted: Boolean(nextQuoted) };
  if (target === "column") {
    if (sourceCol < 0) {
      status.textContent = "Select a column before changing column quote.";
      return;
    }
    body.col = sourceCol;
  }
  const response = await fetch("/api/quote-rule", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Quote rule failed: ${payload.error ?? "unknown error"}`;
    syncQuotePolicyControls();
    return;
  }
  applyMutationPayload(payload, `${quoteRuleLabel(target)} quote ${nextQuoted ? "enabled" : "disabled"}.`);
  document.body.dataset.quoteRuleOverride = `${target}:${target === "column" ? sourceCol : "*"}:${nextQuoted ? "quote" : "off"}`;
  syncQuotePolicyControls();
}

async function setSelectedCellQuoteOverrides(mode) {
  const labels = {
    quote: "Quote Selection",
    "no-quote": "No Quote Selection",
    clear: "Clear Selection Quote Overrides",
  };
  const label = labels[mode] ?? "Quote Selection";
  if (!documentId) {
    status.textContent = "Open a file before changing selection quote rules.";
    return;
  }
  if (!ensureEditable(label)) return;
  const cells = selectedCellsForTransform();
  if (!cells) return;
  const selectedKeys = new Set(cells.map(quoteCellKey));
  const settings = currentQuoteRuleSettings();
  const nextSettings = {
    ...settings,
    quoteCells: settings.quoteCells.filter((cell) => !selectedKeys.has(quoteCellKey(cell))),
    noQuoteCells: settings.noQuoteCells.filter((cell) => !selectedKeys.has(quoteCellKey(cell))),
  };
  if (mode === "quote") {
    nextSettings.quoteCells.push(
      ...cells.filter((cell) => !quoteSettingsPreservesCell(settings, cell)),
    );
  } else if (mode === "no-quote") {
    nextSettings.noQuoteCells.push(...cells);
  }
  nextSettings.quoteCells = normalizeQuoteRuleCells(nextSettings.quoteCells);
  nextSettings.noQuoteCells = normalizeQuoteRuleCells(nextSettings.noQuoteCells);
  const signature = quoteRulesSignature(nextSettings);
  document.body.dataset.quoteSelection = `${mode}:${cells.length}:${signature}`;
  if (signature === quoteRulesSignature(settings)) {
    status.textContent = `${label}: no quote overrides changed.`;
    return;
  }
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...nextSettings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quoteSelection = "error";
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `${label}: ${formatCount(cells.length)} cells updated.`);
  syncQuotePolicyControls();
}

async function applyActiveCellQuoteRule(mode) {
  const labels = {
    quote: "Quote Cell",
    "no-quote": "No Quote Cell",
    clear: "Clear Cell Quote",
  };
  const label = labels[mode] || "Cell Quote";
  if (!documentId) {
    status.textContent = "Open a file before changing cell quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = `${label} blocked: document is read-only`;
    return;
  }
  if (selected.row >= sourceRowCount || columns.length === 0) {
    status.textContent = "Select a source cell before changing cell quote rules.";
    return;
  }
  const row = selectedSourceRow();
  const col = sourceColFor(selected.col);
  const key = `${row}:${col}`;
  const settings = currentQuoteRuleSettings();
  const nextSettings = {
    ...settings,
    quoteCells: settings.quoteCells.filter((cell) => quoteCellKey(cell) !== key),
    noQuoteCells: settings.noQuoteCells.filter((cell) => quoteCellKey(cell) !== key),
  };
  if (mode === "quote") {
    nextSettings.quoteCells.push({ row, col });
  } else if (mode === "no-quote") {
    nextSettings.noQuoteCells.push({ row, col });
  }
  nextSettings.quoteCells = normalizeQuoteRuleCells(nextSettings.quoteCells);
  nextSettings.noQuoteCells = normalizeQuoteRuleCells(nextSettings.noQuoteCells);
  const signature = quoteRulesSignature(nextSettings);
  document.body.dataset.quoteCellRule = `${mode}:${row}:${col}:${signature}`;
  if (signature === quoteRulesSignature(settings)) {
    status.textContent = `${label}: no cell quote rules changed.`;
    return;
  }
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...nextSettings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quoteCellRule = "error";
    syncQuotePolicyControls();
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `${label}: ${quoteRuleCellLabel({ row, col })}.`);
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  renderSettingsPanel();
}

async function applyCoordinateQuoteRule() {
  const mode = quoteRuleActionSelect.value;
  const labels = {
    quote: "Quote Cell Rule",
    "no-quote": "No Quote Cell Rule",
    clear: "Clear Cell Rule",
  };
  const label = labels[mode] || "Cell Rule";
  if (!documentId) {
    status.textContent = "Open a file before editing quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = `${label} blocked: document is read-only`;
    return;
  }
  const row = Number.parseInt(quoteRuleRowInput.value, 10) - 1;
  const col = Number.parseInt(quoteRuleColInput.value, 10) - 1;
  const maxSourceRow = sourceRowCount + dataStartRow - 1;
  const maxSourceCol = Math.max(0, allColumns.length - 1);
  if (!Number.isInteger(row) || !Number.isInteger(col) || row < 0 || col < 0 || row > maxSourceRow || col > maxSourceCol) {
    status.textContent = "Enter a valid source row and column for the quote rule.";
    document.body.dataset.quoteCoordinateRule = "invalid";
    return;
  }
  const key = `${row}:${col}`;
  const settings = currentQuoteRuleSettings();
  const nextSettings = {
    ...settings,
    quoteCells: settings.quoteCells.filter((cell) => quoteCellKey(cell) !== key),
    noQuoteCells: settings.noQuoteCells.filter((cell) => quoteCellKey(cell) !== key),
  };
  if (mode === "quote") {
    nextSettings.quoteCells.push({ row, col });
  } else if (mode === "no-quote") {
    nextSettings.noQuoteCells.push({ row, col });
  }
  nextSettings.quoteCells = normalizeQuoteRuleCells(nextSettings.quoteCells);
  nextSettings.noQuoteCells = normalizeQuoteRuleCells(nextSettings.noQuoteCells);
  const signature = quoteRulesSignature(nextSettings);
  document.body.dataset.quoteCoordinateRule = `${mode}:${row}:${col}:${signature}`;
  if (signature === quoteRulesSignature(settings)) {
    status.textContent = `${label}: no quote rules changed.`;
    return;
  }
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...nextSettings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quoteCoordinateRule = "error";
    syncQuotePolicyControls();
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `${label}: ${quoteRuleCellLabel({ row, col })}.`);
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  renderSettingsPanel();
}

function resolveQuoteRuleColumn(value) {
  const raw = String(value || "").trim();
  if (!raw) return { error: "Enter a source column number or name." };
  if (/^\d+$/.test(raw)) {
    const sourceCol = Number.parseInt(raw, 10) - 1;
    if (allColumns.some((column) => column.sourceCol === sourceCol)) return { col: sourceCol };
    return { error: "Enter a valid source column number or name." };
  }
  const matches = allColumns.filter((column) => String(column.name || "").trim().toLowerCase() === raw.toLowerCase());
  if (matches.length === 1) return { col: matches[0].sourceCol };
  if (matches.length > 1) return { error: "Column name matches multiple source columns; use the source column number." };
  return { error: "Enter a valid source column number or name." };
}

async function applyColumnQuoteRule() {
  const mode = quoteRuleColumnActionSelect.value;
  const label = mode === "clear" ? "Clear Column Rule" : "Quote Column Rule";
  if (!documentId) {
    status.textContent = "Open a file before editing quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = `${label} blocked: document is read-only`;
    return;
  }
  const resolved = resolveQuoteRuleColumn(quoteRuleColumnInput.value);
  if (resolved.error) {
    status.textContent = resolved.error;
    document.body.dataset.quoteColumnRule = "invalid";
    return;
  }
  const col = resolved.col;
  const settings = currentQuoteRuleSettings();
  const nextSettings = {
    ...settings,
    quoteColumns: settings.quoteColumns.filter((sourceCol) => sourceCol !== col),
  };
  if (mode === "quote") {
    nextSettings.quoteColumns.push(col);
  }
  nextSettings.quoteColumns = normalizeIntegerList(nextSettings.quoteColumns);
  const signature = quoteRulesSignature(nextSettings);
  document.body.dataset.quoteColumnRule = `${mode}:${col}:${signature}`;
  if (signature === quoteRulesSignature(settings)) {
    status.textContent = `${label}: no quote rules changed.`;
    return;
  }
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...nextSettings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quoteColumnRule = "error";
    syncQuotePolicyControls();
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `${label}: ${quoteRuleColumnLabel(col)}.`);
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  renderSettingsPanel();
}

function regexEscape(text) {
  return String(text).replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

function globToRegex(pattern) {
  return new RegExp(`^${String(pattern).split("*").map(regexEscape).join(".*")}$`, "i");
}

function quotePatternMatcher(pattern, matchMode, emptyMessage, regexMessage) {
  const raw = String(pattern || "").trim();
  if (!raw) return { error: emptyMessage };
  if (matchMode === "regex") {
    try {
      const expression = new RegExp(raw, "i");
      return { matcher: (value) => expression.test(value) };
    } catch {
      return { error: regexMessage };
    }
  }
  if (matchMode === "glob") {
    const expression = globToRegex(raw);
    return { matcher: (value) => expression.test(value) };
  }
  if (matchMode === "exact") {
    const needle = raw.toLowerCase();
    return { matcher: (value) => value.toLowerCase() === needle };
  }
  const needle = raw.toLowerCase();
  return { matcher: (value) => value.toLowerCase().includes(needle) };
}

function matchingQuoteRuleColumns(pattern, matchMode) {
  const compiled = quotePatternMatcher(
    pattern,
    matchMode,
    "Enter a column pattern before applying quote rules.",
    "Enter a valid regular expression for the column pattern.",
  );
  if (compiled.error) return compiled;
  const columnsBySource = new Map();
  for (const column of allColumns) {
    const label = quoteRuleColumnLabel(column.sourceCol);
    if (compiled.matcher(label)) columnsBySource.set(column.sourceCol, column);
  }
  const columns = Array.from(columnsBySource.keys()).sort((left, right) => left - right);
  if (!columns.length) {
    return { error: "No source columns match the quote pattern." };
  }
  return { columns };
}

async function applyColumnPatternQuoteRule() {
  const action = quoteRuleColumnPatternActionSelect.value;
  const matchMode = quoteRuleColumnPatternModeSelect.value;
  const label = action === "clear" ? "Clear Column Pattern Rule" : "Quote Column Pattern Rule";
  if (!documentId) {
    status.textContent = "Open a file before editing quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = `${label} blocked: document is read-only`;
    return;
  }
  const matched = matchingQuoteRuleColumns(quoteRuleColumnPatternInput.value, matchMode);
  if (matched.error) {
    status.textContent = matched.error;
    document.body.dataset.quoteColumnPatternRule = "invalid";
    return;
  }
  const matchedColumns = new Set(matched.columns);
  const settings = currentQuoteRuleSettings();
  const nextSettings = {
    ...settings,
    quoteColumns: settings.quoteColumns.filter((sourceCol) => !matchedColumns.has(sourceCol)),
  };
  if (action === "quote") {
    nextSettings.quoteColumns.push(...matched.columns);
  }
  nextSettings.quoteColumns = normalizeIntegerList(nextSettings.quoteColumns);
  const signature = quoteRulesSignature(nextSettings);
  document.body.dataset.quoteColumnPatternRule = `${matchMode}:${action}:${matched.columns.join(",")}:${signature}`;
  if (signature === quoteRulesSignature(settings)) {
    status.textContent = `${label}: no quote rules changed.`;
    return;
  }
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...nextSettings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quoteColumnPatternRule = "error";
    syncQuotePolicyControls();
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `${label}: ${formatCount(matched.columns.length)} columns.`);
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  renderSettingsPanel();
}

function matchingQuoteRuleCellsByValue(pattern, matchMode, maxCells = 50_000) {
  const compiled = quotePatternMatcher(
    pattern,
    matchMode,
    "Enter a value pattern before applying quote rules.",
    "Enter a valid regular expression for the value pattern.",
  );
  if (compiled.error) return compiled;
  const scanCells = sourceRowCount * columns.length;
  if (scanCells > maxCells) {
    return { error: `Value pattern quote rules can scan up to ${formatCount(maxCells)} loaded cells.` };
  }
  const cells = [];
  for (let row = 0; row < sourceRowCount; row++) {
    if (!rows.has(row)) {
      return { error: "Load the target rows before applying value-pattern quote rules." };
    }
    const recordRow = rowSource.get(row) ?? row + dataStartRow;
    for (let col = 0; col < columns.length; col++) {
      const value = String(valueAt(row, col) ?? "");
      if (compiled.matcher(value)) {
        cells.push({ row: recordRow, col: sourceColFor(col) });
      }
    }
  }
  if (!cells.length) {
    return { error: "No loaded cells match the quote value pattern." };
  }
  return { cells };
}

async function applyValuePatternQuoteRule() {
  const action = quoteRuleValuePatternActionSelect.value;
  const matchMode = quoteRuleValuePatternModeSelect.value;
  const labels = {
    quote: "Quote Value Pattern Rule",
    "no-quote": "No Quote Value Pattern Rule",
    clear: "Clear Value Pattern Rule",
  };
  const label = labels[action] || "Value Pattern Rule";
  if (!documentId) {
    status.textContent = "Open a file before editing quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = `${label} blocked: document is read-only`;
    return;
  }
  const matched = matchingQuoteRuleCellsByValue(quoteRuleValuePatternInput.value, matchMode);
  if (matched.error) {
    status.textContent = matched.error;
    document.body.dataset.quoteValuePatternRule = "invalid";
    return;
  }
  const matchedKeys = new Set(matched.cells.map(quoteCellKey));
  const settings = currentQuoteRuleSettings();
  const nextSettings = {
    ...settings,
    quoteCells: settings.quoteCells.filter((cell) => !matchedKeys.has(quoteCellKey(cell))),
    noQuoteCells: settings.noQuoteCells.filter((cell) => !matchedKeys.has(quoteCellKey(cell))),
  };
  if (action === "quote") {
    nextSettings.quoteCells.push(
      ...matched.cells.filter((cell) => !quoteSettingsPreservesCell(settings, cell)),
    );
  } else if (action === "no-quote") {
    nextSettings.noQuoteCells.push(...matched.cells);
  }
  nextSettings.quoteCells = normalizeQuoteRuleCells(nextSettings.quoteCells);
  nextSettings.noQuoteCells = normalizeQuoteRuleCells(nextSettings.noQuoteCells);
  const signature = quoteRulesSignature(nextSettings);
  document.body.dataset.quoteValuePatternRule = `${matchMode}:${action}:${matched.cells.length}:${signature}`;
  if (signature === quoteRulesSignature(settings)) {
    status.textContent = `${label}: no quote rules changed.`;
    return;
  }
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...nextSettings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quoteValuePatternRule = "error";
    syncQuotePolicyControls();
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `${label}: ${formatCount(matched.cells.length)} cells.`);
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  renderSettingsPanel();
}

function quoteCellsForSourceRowRange(startValue, endValue, maxCells = 50_000) {
  const startRow = Number.parseInt(startValue, 10) - 1;
  const endRow = Number.parseInt(endValue, 10) - 1;
  const maxSourceRow = sourceRowCount + dataStartRow - 1;
  if (!Number.isInteger(startRow) || !Number.isInteger(endRow) || startRow < 0 || endRow < 0 || startRow > endRow || endRow > maxSourceRow) {
    return { error: "Enter a valid source row range for the quote rule." };
  }
  const sourceColumns = normalizeIntegerList(allColumns.map((column) => column.sourceCol));
  if (!sourceColumns.length) {
    return { error: "Open a file with columns before applying row quote rules." };
  }
  const scanCells = (endRow - startRow + 1) * sourceColumns.length;
  if (scanCells > maxCells) {
    return { error: `Row quote rules can update up to ${formatCount(maxCells)} cells at once.` };
  }
  const cells = [];
  for (let row = startRow; row <= endRow; row++) {
    for (const col of sourceColumns) {
      cells.push({ row, col });
    }
  }
  return { cells, startRow, endRow };
}

async function applyRowRangeQuoteRule() {
  const action = quoteRuleRowActionSelect.value;
  const labels = {
    quote: "Quote Row Range Rule",
    "no-quote": "No Quote Row Range Rule",
    clear: "Clear Row Range Rule",
  };
  const label = labels[action] || "Row Range Rule";
  if (!documentId) {
    status.textContent = "Open a file before editing quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = `${label} blocked: document is read-only`;
    return;
  }
  const range = quoteCellsForSourceRowRange(quoteRuleRowStartInput.value, quoteRuleRowEndInput.value);
  if (range.error) {
    status.textContent = range.error;
    document.body.dataset.quoteRowRangeRule = "invalid";
    return;
  }
  const matchedKeys = new Set(range.cells.map(quoteCellKey));
  const settings = currentQuoteRuleSettings();
  const nextSettings = {
    ...settings,
    quoteCells: settings.quoteCells.filter((cell) => !matchedKeys.has(quoteCellKey(cell))),
    noQuoteCells: settings.noQuoteCells.filter((cell) => !matchedKeys.has(quoteCellKey(cell))),
  };
  if (action === "quote") {
    nextSettings.quoteCells.push(
      ...range.cells.filter((cell) => !quoteSettingsPreservesCell(settings, cell)),
    );
  } else if (action === "no-quote") {
    nextSettings.noQuoteCells.push(...range.cells);
  }
  nextSettings.quoteCells = normalizeQuoteRuleCells(nextSettings.quoteCells);
  nextSettings.noQuoteCells = normalizeQuoteRuleCells(nextSettings.noQuoteCells);
  const signature = quoteRulesSignature(nextSettings);
  document.body.dataset.quoteRowRangeRule = `${action}:${range.startRow}-${range.endRow}:${range.cells.length}:${signature}`;
  if (signature === quoteRulesSignature(settings)) {
    status.textContent = `${label}: no quote rules changed.`;
    return;
  }
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...nextSettings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quoteRowRangeRule = "error";
    syncQuotePolicyControls();
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `${label}: ${formatCount(range.cells.length)} cells.`);
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  renderSettingsPanel();
}

function quoteCellKey(cell) {
  return `${cell.row}:${cell.col}`;
}

function quoteSettingsPreservesCell(settings, cell) {
  const hasHeader = Boolean(metadata?.dialect?.hasHeader);
  const bodyRow = !hasHeader || cell.row > 0;
  if (hasHeader && settings.quoteHeader && cell.row === 0) return true;
  if (settings.quoteBody && bodyRow) return true;
  if (bodyRow && settings.quoteColumns.includes(cell.col)) return true;
  return settings.quoteCells.some((quotedCell) => quotedCell.row === cell.row && quotedCell.col === cell.col);
}

function quoteRuleLabel(target) {
  if (target === "header") return "Header";
  if (target === "body") return "Body";
  return "Column";
}

function exportQuoteRules() {
  if (!metadata?.dialect) {
    status.textContent = "Open a file before exporting quote rules.";
    return;
  }
  const settings = currentQuoteRuleSettings();
  const payload = {
    kind: quoteRulesExportKind,
    version: 1,
    exportedAt: new Date().toISOString(),
    ...settings,
  };
  const text = `${JSON.stringify(payload, null, 2)}\n`;
  window.__opencsvLastQuoteRulesExportText = text;
  document.body.dataset.quoteRulesExport = quoteRulesSignature(settings);
  document.body.dataset.quoteRulesExportPreview = text.slice(0, 500);
  triggerDownload(text, "opencsv-quote-rules.json", "application/json;charset=utf-8");
  status.textContent = `Exported quote rules: ${quoteRulesSummary(settings)}.`;
}

function showQuoteRulesImportPanel() {
  openSettingsPanel();
  quoteRulesImportPanel.hidden = false;
  quoteRulesImportText.value = "";
  document.body.dataset.quoteRulesImportPanel = "open";
  syncQuoteRulesImportControls();
  quoteRulesImportText.focus();
}

function cancelQuoteRulesImport() {
  quoteRulesImportPanel.hidden = true;
  quoteRulesImportText.value = "";
  document.body.dataset.quoteRulesImportPanel = "closed";
  syncQuoteRulesImportControls();
}

function syncQuoteRulesImportControls() {
  if (!quoteRulesImportText) return;
  applyQuoteRulesImportButton.disabled = quoteRulesImportText.value.trim() === "";
}

async function applyQuoteRulesImport() {
  if (!documentId) {
    status.textContent = "Open a file before importing quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = "Quote rule import blocked: document is read-only";
    return;
  }
  try {
    const settings = parseQuoteRulesImport(quoteRulesImportText.value);
    const response = await fetch("/api/quote-settings", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ documentId, ...settings }),
    });
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "quote import failed");
    quotePolicyDraftDirty = false;
    quoteRulesImportPanel.hidden = true;
    quoteRulesImportText.value = "";
    document.body.dataset.quoteRulesImport = quoteRulesSignature(settings);
    document.body.dataset.quoteRulesImportPanel = "closed";
    applyMutationPayload(payload, `Imported quote rules: ${quoteRulesSummary(settings)}.`);
    syncQuotePolicyControls();
  } catch (error) {
    status.textContent = `Quote rule import failed: ${error.message}`;
    document.body.dataset.quoteRulesImport = "error";
  }
  syncQuoteRulesImportControls();
}

function currentQuoteRuleSettings() {
  const dialect = metadata?.dialect ?? {};
  return {
    quoteChar: quoteCharValue(dialect),
    quoteMode: quoteModeValue(dialect),
    quoteEscape: quoteEscapeValue(dialect),
    quoteHeader: Boolean(dialect.quoteHeader),
    quoteBody: Boolean(dialect.quoteBody),
    quoteColumns: normalizeIntegerList(dialect.quoteColumns),
    quoteCells: normalizeQuoteRuleCells(dialect.quoteCells),
    noQuoteCells: normalizeQuoteRuleCells(dialect.noQuoteCells),
  };
}

function parseQuoteRulesImport(text) {
  let parsed;
  try {
    parsed = JSON.parse(text);
  } catch {
    throw new Error("invalid JSON");
  }
  const source = parsed?.rules && typeof parsed.rules === "object" ? parsed.rules : parsed;
  if (!source || typeof source !== "object" || Array.isArray(source)) {
    throw new Error("missing quote rules object");
  }
  return {
    quoteChar: normalizeQuoteRuleChar(source.quoteChar),
    quoteMode: normalizeQuoteRuleMode(source.quoteMode),
    quoteEscape: normalizeQuoteRuleEscape(source.quoteEscape),
    quoteHeader: normalizeQuoteRuleBoolean(source.quoteHeader),
    quoteBody: normalizeQuoteRuleBoolean(source.quoteBody),
    quoteColumns: normalizeIntegerList(source.quoteColumns),
    quoteCells: normalizeQuoteRuleCells(source.quoteCells),
    noQuoteCells: normalizeQuoteRuleCells(source.noQuoteCells),
  };
}

function normalizeQuoteRuleMode(value) {
  const mode = String(value || "manual").toLowerCase().trim();
  return ["minimal", "always", "never", "manual"].includes(mode) ? mode : "manual";
}

function normalizeQuoteRuleChar(value) {
  const raw = String(value ?? "double").toLowerCase().trim();
  if (raw === "'" || raw === "single") return "single";
  if (raw === "" || raw === "none") return "none";
  return "double";
}

function normalizeQuoteRuleEscape(value) {
  return normalizeQuoteEscapeValue(value);
}

function normalizeQuoteRuleBoolean(value) {
  if (value === true) return true;
  if (value === false || value == null) return false;
  if (typeof value === "string") {
    return ["true", "1", "yes", "on"].includes(value.toLowerCase().trim());
  }
  return Boolean(value);
}

function normalizeIntegerList(values) {
  if (!Array.isArray(values)) return [];
  return Array.from(new Set(
    values
      .map((value) => Number.parseInt(value, 10))
      .filter((value) => Number.isInteger(value) && value >= 0),
  )).sort((left, right) => left - right);
}

function normalizeQuoteRuleCells(values) {
  if (!Array.isArray(values)) return [];
  const seen = new Set();
  const cells = [];
  for (const value of values) {
    const row = Number.parseInt(value?.row, 10);
    const col = Number.parseInt(value?.col, 10);
    if (!Number.isInteger(row) || !Number.isInteger(col) || row < 0 || col < 0) continue;
    const key = `${row}:${col}`;
    if (seen.has(key)) continue;
    seen.add(key);
    cells.push({ row, col });
  }
  return cells.sort((left, right) => left.row === right.row ? left.col - right.col : left.row - right.row);
}

function quoteRulesSignature(settings) {
  return [
    settings.quoteMode,
    settings.quoteChar,
    normalizeQuoteEscapeValue(settings.quoteEscape),
    settings.quoteHeader ? "header" : "no-header",
    settings.quoteBody ? "body" : "no-body",
    `cols:${settings.quoteColumns.join(",")}`,
    `quote:${settings.quoteCells.map((cell) => `${cell.row}:${cell.col}`).join(",")}`,
    `no:${settings.noQuoteCells.map((cell) => `${cell.row}:${cell.col}`).join(",")}`,
  ].join("|");
}

function quoteRulesSummary(settings) {
  const parts = [`${settings.quoteMode} ${settings.quoteChar} ${normalizeQuoteEscapeValue(settings.quoteEscape)} escape`];
  if (settings.quoteHeader) parts.push("header");
  if (settings.quoteBody) parts.push("body");
  if (settings.quoteColumns.length) parts.push(`${settings.quoteColumns.length} columns`);
  if (settings.quoteCells.length) parts.push(`${settings.quoteCells.length} quoted cells`);
  if (settings.noQuoteCells.length) parts.push(`${settings.noQuoteCells.length} no-quote cells`);
  return parts.join(", ");
}

function renderQuoteRuleList() {
  if (!quoteRuleList) return;
  quoteRuleList.replaceChildren();
  if (!documentId) {
    quoteRuleList.append(emptyQuoteRuleNode("Open a file to inspect quote rules."));
    document.body.dataset.quoteRuleList = "closed";
    return;
  }
  const settings = currentQuoteRuleSettings();
  const items = activeQuoteRuleItems(settings);
  document.body.dataset.quoteRuleList = String(items.length);
  if (!items.length) {
    quoteRuleList.append(emptyQuoteRuleNode("No active quote rules."));
    return;
  }
  for (const item of items) {
    const row = document.createElement("div");
    row.className = "quote-rule-item";
    row.dataset.ruleKind = item.kind;
    row.dataset.ruleKey = item.key;

    const meta = document.createElement("div");
    const title = document.createElement("strong");
    title.textContent = item.title;
    const detail = document.createElement("span");
    detail.textContent = item.detail;
    meta.append(title, detail);

    const remove = document.createElement("button");
    remove.type = "button";
    remove.textContent = "Remove";
    remove.disabled = readOnly;
    remove.addEventListener("click", () => removeQuoteRule(item));
    row.append(meta, remove);
    quoteRuleList.append(row);
  }
}

function emptyQuoteRuleNode(message) {
  const node = document.createElement("div");
  node.className = "quote-rule-empty";
  node.textContent = message;
  return node;
}

function activeQuoteRuleItems(settings) {
  const items = [];
  if (settings.quoteHeader) {
    items.push({
      kind: "header",
      key: "header",
      title: "Quote header row",
      detail: "Force all header fields to be quoted.",
    });
  }
  if (settings.quoteBody) {
    items.push({
      kind: "body",
      key: "body",
      title: "Quote body rows",
      detail: "Force all data fields to be quoted.",
    });
  }
  for (const col of settings.quoteColumns) {
    items.push({
      kind: "column",
      key: String(col),
      col,
      title: `Quote ${quoteRuleColumnLabel(col)}`,
      detail: `Source column ${col + 1}`,
    });
  }
  for (const cell of settings.quoteCells) {
    items.push({
      kind: "cell",
      key: quoteCellKey(cell),
      cell,
      title: `Quote cell ${quoteRuleCellLabel(cell)}`,
      detail: `Source row ${cell.row + 1}, column ${cell.col + 1}`,
    });
  }
  for (const cell of settings.noQuoteCells) {
    items.push({
      kind: "no-cell",
      key: quoteCellKey(cell),
      cell,
      title: `No quote cell ${quoteRuleCellLabel(cell)}`,
      detail: `Source row ${cell.row + 1}, column ${cell.col + 1}`,
    });
  }
  return items;
}

function quoteRuleColumnLabel(sourceCol) {
  const column = allColumns.find((candidate) => candidate.sourceCol === sourceCol);
  return column?.name ? `${column.name}` : `column ${sourceCol + 1}`;
}

function quoteRuleCellLabel(cell) {
  return `R${cell.row + 1}C${cell.col + 1}`;
}

async function removeQuoteRule(item) {
  if (!documentId) {
    status.textContent = "Open a file before changing quote rules.";
    return;
  }
  if (readOnly) {
    status.textContent = "Quote rule removal blocked: document is read-only";
    return;
  }
  const settings = currentQuoteRuleSettings();
  const nextSettings = {
    ...settings,
    quoteColumns: [...settings.quoteColumns],
    quoteCells: [...settings.quoteCells],
    noQuoteCells: [...settings.noQuoteCells],
  };
  if (item.kind === "header") nextSettings.quoteHeader = false;
  if (item.kind === "body") nextSettings.quoteBody = false;
  if (item.kind === "column") {
    nextSettings.quoteColumns = nextSettings.quoteColumns.filter((col) => col !== item.col);
  }
  if (item.kind === "cell") {
    nextSettings.quoteCells = nextSettings.quoteCells.filter((cell) => quoteCellKey(cell) !== item.key);
  }
  if (item.kind === "no-cell") {
    nextSettings.noQuoteCells = nextSettings.noQuoteCells.filter((cell) => quoteCellKey(cell) !== item.key);
  }
  document.body.dataset.quoteRuleRemove = `${item.kind}:${item.key}`;
  const response = await fetch("/api/quote-settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...nextSettings }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Quote rule removal failed: ${payload.error ?? "unknown error"}`;
    document.body.dataset.quoteRuleRemove = "error";
    syncQuotePolicyControls();
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `Removed quote rule: ${item.title}.`);
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  renderSettingsPanel();
}

function normalizeFormatRules(rules) {
  if (!Array.isArray(rules)) return [];
  const normalized = [];
  for (const rule of rules) {
    const item = normalizeFormatRule(rule);
    if (item && !normalized.some((existing) => formatRuleKey(existing) === formatRuleKey(item))) {
      normalized.push(item);
    }
  }
  return normalized.slice(0, 12);
}

function normalizeFormatRule(rule) {
  const pattern = typeof rule === "string" ? rule : rule?.pattern;
  if (typeof pattern !== "string" || pattern.trim() === "") return null;
  return { pattern: pattern.trim() };
}

function formatRuleKey(rule) {
  return rule.pattern.toLowerCase();
}

function formatRuleTarget(pattern) {
  const lower = pattern.trim().toLowerCase();
  if (lower.startsWith("path:")) return "path";
  if (lower.startsWith("content:")) return "content";
  return "name";
}

function formatRulePattern(pattern) {
  return pattern.replace(/^(name|path|content):/i, "").trim();
}

function defaultFormatRulePattern(path = metadata?.path || pathInput.value.trim()) {
  const name = fileNameFromPath(path);
  const dot = name.lastIndexOf(".");
  if (dot > 0 && dot < name.length - 1) return `*.${name.slice(dot + 1)}`;
  return name ? `name:${name}` : "";
}

function updateFormatRulePlaceholder(path = metadata?.path || pathInput.value.trim()) {
  const pattern = defaultFormatRulePattern(path);
  if (pattern) formatRuleInput.placeholder = pattern;
}

async function matchingFormatPresetForPath(path) {
  const pathMatch = findMatchingFormatPreset(path, "");
  if (pathMatch) return pathMatch;
  if (!formatPresets().some((preset) => preset.rules?.some((rule) => formatRuleTarget(rule.pattern) === "content"))) {
    return null;
  }
  const sample = await fileSampleForRules(path);
  if (sample === null) return null;
  return findMatchingFormatPreset(path, sample);
}

function findMatchingFormatPreset(path, sample) {
  const basename = fileNameFromPath(path).toLowerCase();
  const normalizedPath = path.replace(/\\/g, "/").toLowerCase();
  const normalizedSample = sample.toLowerCase();
  for (const preset of formatPresets()) {
    for (const rule of preset.rules ?? []) {
      const target = formatRuleTarget(rule.pattern);
      const pattern = formatRulePattern(rule.pattern);
      if (!pattern) continue;
      if (target === "content") {
        if (normalizedSample && normalizedSample.includes(pattern.toLowerCase())) {
          return { preset, rule };
        }
        continue;
      }
      const value = target === "path" ? normalizedPath : basename;
      if (matchesFormatRuleValue(value, pattern.toLowerCase())) {
        return { preset, rule };
      }
    }
  }
  return null;
}

function matchesFormatRuleValue(value, pattern) {
  if (pattern.includes("*") || pattern.includes("?")) {
    return globToRegExp(pattern).test(value);
  }
  return value === pattern || value.includes(pattern);
}

function globToRegExp(pattern) {
  const escaped = pattern
    .replace(/[.+^${}()|[\]\\]/g, "\\$&")
    .replace(/\*/g, ".*")
    .replace(/\?/g, ".");
  return new RegExp(`^${escaped}$`, "i");
}

async function fileSampleForRules(path) {
  try {
    const response = await fetch("/api/file-sample", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ path, maxBytes: formatRuleSampleBytes }),
    });
    const payload = await response.json();
    if (!response.ok) return null;
    return payload.sample ?? "";
  } catch {
    return null;
  }
}

function formatPresetID(preset) {
  const parts = [
    preset.delimiter === "\t" ? "tab" : preset.delimiter,
    preset.hasHeader ? "header" : "no-header",
    (preset.encodingName || "utf-8").toLowerCase(),
    normalizeLineEndingName(preset.lineEndingName),
    preset.quoteChar || "double",
    preset.quoteMode || "minimal",
    normalizeQuoteEscapeValue(preset.quoteEscape),
  ];
  const headerColumns = normalizeHeaderColumns(preset.headerColumns);
  if (headerColumns > 0) parts.splice(2, 0, `hcols${headerColumns}`);
  return parts.join(":");
}

function formatPresetName(preset) {
  const delimiter = preset.delimiterName || (preset.delimiter === "\t" ? "tab" : preset.delimiter);
  const lineEnding = normalizeLineEndingName(preset.lineEndingName).toUpperCase();
  return `${delimiterLabel({ delimiter: preset.delimiter, delimiterName: delimiter })} · ${formatPresetHeaderLabel(preset)} · ${(preset.encodingName || "UTF-8").toUpperCase()} · ${lineEnding} · ${quoteModeValue(preset)} ${quoteCharLabel(preset)} ${quoteEscapeLabel(preset)}`;
}

function formatPresetPreview(preset) {
  return `${preset.delimiter === "\t" ? "\\t" : preset.delimiter}|${preset.hasHeader ? "header" : "no-header"}|hcols:${normalizeHeaderColumns(preset.headerColumns)}|${preset.encodingName || "utf-8"}|${normalizeLineEndingName(preset.lineEndingName)}|${preset.quoteMode || "minimal"}|${preset.quoteChar || "double"}|${normalizeQuoteEscapeValue(preset.quoteEscape)}`;
}

function formatPresetHeaderLabel(preset) {
  const headerColumns = normalizeHeaderColumns(preset.headerColumns);
  const columnLabel = headerColumns === 1 ? "1 Header Col" : `${headerColumns} Header Cols`;
  if (preset.hasHeader && headerColumns > 0) return `Header Row + ${columnLabel}`;
  if (preset.hasHeader) return "Header";
  if (headerColumns > 0) return columnLabel;
  return "No header";
}

function quoteLabel(value, displayRow = selected.row, visibleCol = selected.col) {
  if (suppressesQuote(displayRow, visibleCol)) return "manual no quote";
  if (!shouldQuote(value, displayRow, visibleCol)) return "no quote needed";
  return preservesQuote(displayRow, visibleCol) ? "will preserve quote on save" : "will quote on save";
}

function normalizeCellSeparatorToken(value) {
  const token = String(value ?? "");
  return token.length ? token : defaultCellSeparatorToken;
}

function setCellSeparatorToken(value, persist = true) {
  cellSeparatorToken = normalizeCellSeparatorToken(value);
  if (cellSeparatorInput) cellSeparatorInput.value = cellSeparatorToken;
  if (persist) localStorage.setItem(cellSeparatorStorageKey, cellSeparatorToken);
  syncCellSeparatorState();
}

function updateCellSeparatorFromInput() {
  cellSeparatorToken = normalizeCellSeparatorToken(cellSeparatorInput.value);
  localStorage.setItem(cellSeparatorStorageKey, cellSeparatorToken);
  syncCellSeparatorState();
}

function selectedCellSeparator() {
  if (cellSeparatorInput) {
    cellSeparatorToken = normalizeCellSeparatorToken(cellSeparatorInput.value);
  }
  syncCellSeparatorState();
  return parseCellSeparatorToken(cellSeparatorToken);
}

function parseCellSeparatorToken(token) {
  const value = normalizeCellSeparatorToken(token);
  const alias = value.trim().toLowerCase();
  if (alias === "space" || value === "\\s") return " ";
  if (alias === "tab" || value === "\\t") return "\t";
  if (alias === "newline" || alias === "linefeed" || value === "\\n") return "\n";
  return value;
}

function printableCellSeparator(separator) {
  if (separator === " ") return "space";
  if (separator === "\t") return "\\t";
  if (separator === "\n") return "\\n";
  return String(separator ?? "");
}

function syncCellSeparatorState() {
  const separator = parseCellSeparatorToken(cellSeparatorToken);
  document.body.dataset.cellSeparatorToken = cellSeparatorToken;
  document.body.dataset.cellSeparator = printableCellSeparator(separator);
  document.body.dataset.cellSeparatorLength = String(separator.length);
}

function syncCellValuePanel() {
  const key = selectedCellKey();
  const value = documentId ? valueAt(selected.row, selected.col) : "";
  const aggregate = Boolean(aggregateRowAt(selected.row));
  if (key !== cellValueInputKey || !cellValueInputDirty) {
    cellValueInput.value = value;
    cellValueInputKey = key;
    cellValueInputDirty = false;
  }
  const disabled = !documentId || readOnly || aggregate;
  const panelValue = cellValueInputDirty && key === cellValueInputKey ? cellValueInput.value : value;
  cellValueInput.disabled = disabled;
  applyCellValueButton.disabled = disabled || !cellValueInputDirty;
  revertCellValueButton.disabled = disabled || !cellValueInputDirty;
  openCellValueDialogButton.disabled = !documentId;
  syncCellValueStats(panelValue);
  syncCellValueEditorSize();
  syncCellLinkState(panelValue);
  syncCellValueDialog();
}

function cellValueSummary(value) {
  const text = String(value ?? "");
  const chars = Array.from(text).length;
  const lines = text.length ? text.split(/\r\n|\r|\n/).length : 1;
  return {
    chars,
    lines,
    label: `${chars} ${chars === 1 ? "char" : "chars"} · ${lines} ${lines === 1 ? "line" : "lines"}`,
  };
}

function syncCellValueStats(value) {
  const summary = cellValueSummary(value);
  cellValueStats.textContent = summary.label;
  document.body.dataset.cellValueChars = String(summary.chars);
  document.body.dataset.cellValueLines = String(summary.lines);
  document.body.dataset.cellValueDirty = cellValueInputDirty ? "true" : "false";
}

function syncCellValueEditorSize() {
  cellValueInput.classList.toggle("expanded", cellValueExpanded);
  toggleCellValueLargeButton.textContent = cellValueExpanded ? "Collapse" : "Expand";
  toggleCellValueLargeButton.setAttribute("aria-pressed", cellValueExpanded ? "true" : "false");
  document.body.dataset.cellValueExpanded = cellValueExpanded ? "true" : "false";
  cellValueInput.style.height = "auto";
  const maxHeight = cellValueExpanded ? 360 : 220;
  const desiredHeight = Math.max(128, Math.min(maxHeight, cellValueInput.scrollHeight + 2));
  cellValueInput.style.height = `${desiredHeight}px`;
  cellValueInput.style.overflowY = cellValueInput.scrollHeight > maxHeight ? "auto" : "hidden";
}

function toggleCellValueLarge() {
  cellValueExpanded = !cellValueExpanded;
  syncCellValueEditorSize();
}

function cellValueEditorDisabled() {
  return !documentId || readOnly || Boolean(aggregateRowAt(selected.row));
}

function syncCellValueDraft(value, source = "panel") {
  cellValueInputDirty = true;
  cellValueInput.value = value;
  const disabled = cellValueEditorDisabled();
  applyCellValueButton.disabled = disabled;
  revertCellValueButton.disabled = disabled;
  syncCellValueStats(value);
  syncCellValueEditorSize();
  syncCellLinkState(value);
  if (cellValueDialogOpen) {
    if (source !== "dialog") cellValueDialogInput.value = value;
    syncCellValueDialog();
  }
}

function openCellValueDialog() {
  if (!documentId) {
    status.textContent = "Open a CSV before using the Cell Value editor.";
    return;
  }
  cellValueDialogOpen = true;
  cellValueDialog.hidden = false;
  document.body.dataset.cellValueDialog = "open";
  cellValueDialogInput.value = cellValueInput.value;
  syncCellValueDialog();
  cellValueDialogInput.focus();
  cellValueDialogInput.select();
}

function closeCellValueDialog() {
  if (!cellValueDialogOpen && cellValueDialog.hidden) return;
  cellValueDialogOpen = false;
  cellValueDialog.hidden = true;
  document.body.dataset.cellValueDialog = "closed";
  cellValueInput.focus();
}

function syncCellValueDialog() {
  document.body.dataset.cellValueDialog = cellValueDialogOpen ? "open" : "closed";
  if (!cellValueDialogOpen) return;
  if (document.activeElement !== cellValueDialogInput || !cellValueInputDirty) {
    cellValueDialogInput.value = cellValueInput.value;
  }
  const disabled = cellValueEditorDisabled();
  const summary = cellValueSummary(cellValueDialogInput.value);
  cellValueDialogStats.textContent = summary.label;
  cellValueDialogInput.disabled = disabled;
  applyCellValueDialogButton.disabled = disabled || !cellValueInputDirty;
  revertCellValueDialogButton.disabled = disabled || !cellValueInputDirty;
  document.body.dataset.cellValueDialogChars = String(summary.chars);
  document.body.dataset.cellValueDialogLines = String(summary.lines);
  document.body.dataset.cellValueDialogDirty = cellValueInputDirty ? "true" : "false";
}

function revertCellValueInput() {
  if (!cellValueInputDirty) return;
  const value = documentId ? valueAt(selected.row, selected.col) : "";
  cellValueInput.value = value;
  cellValueInputDirty = false;
  applyCellValueButton.disabled = true;
  revertCellValueButton.disabled = true;
  syncCellValueStats(value);
  syncCellValueEditorSize();
  syncCellLinkState(value);
  if (cellValueDialogOpen) {
    cellValueDialogInput.value = value;
    syncCellValueDialog();
  }
  status.textContent = "Cell edit reverted";
}

function syncCellLinkState(value) {
  const link = analyzeCellLink(value);
  syncCellLinkAudit(link);
  if (link.href) {
    openCellLink.href = link.href;
    openCellLink.title = `${link.protocol.toUpperCase()} ${link.target}`;
    openCellLink.setAttribute("aria-disabled", "false");
    openCellLink.classList.remove("disabled");
  } else {
    openCellLink.removeAttribute("href");
    openCellLink.title = link.status === "blocked" ? linkUnavailableMessage(link) : "No link";
    openCellLink.setAttribute("aria-disabled", "true");
    openCellLink.classList.add("disabled");
  }
}

function syncCellLinkAudit(link) {
  const statusValue = link.status || "none";
  document.body.dataset.cellLinkStatus = statusValue;
  document.body.dataset.cellLinkProtocol = link.protocol || "";
  document.body.dataset.cellLinkTarget = link.target || "";
  document.body.dataset.cellLinkHref = link.href || "";
  cellLinkAudit.classList.remove("link-trusted", "link-custom", "link-blocked");
  if (statusValue === "none") {
    cellLinkAudit.hidden = true;
    cellLinkProtocol.textContent = "Link";
    cellLinkTarget.textContent = "";
    return;
  }
  cellLinkAudit.hidden = false;
  cellLinkAudit.classList.add(`link-${statusValue}`);
  const protocol = link.protocol ? link.protocol.toUpperCase() : "Link";
  if (statusValue === "blocked") {
    cellLinkProtocol.textContent = `${protocol} blocked`;
  } else if (statusValue === "custom") {
    cellLinkProtocol.textContent = `${protocol} custom link`;
  } else {
    cellLinkProtocol.textContent = `${protocol} link`;
  }
  cellLinkTarget.textContent = link.target || link.href || "Link";
}

function selectedCellKey() {
  return documentId ? `${documentId}:${selected.row}:${selected.col}:${aggregateRowAt(selected.row) ? "aggregate" : "data"}` : "";
}

function normalizeURL(value) {
  return analyzeCellLink(value).href;
}

function analyzeCellLink(value) {
  const text = String(value ?? "").trim();
  if (!text) return emptyCellLink();
  const bareDomain = /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+(?:[/?#][^\s]*)?$/i;
  const hasProtocol = /^[a-z][a-z0-9+.-]*:/i.test(text);
  if (/\s/.test(text)) {
    return hasProtocol ? blockedCellLink(text, protocolNameFromText(text), "Contains whitespace") : emptyCellLink();
  }
  const candidate = bareDomain.test(text) ? `https://${text}` : text;
  if (!/^[a-z][a-z0-9+.-]*:/i.test(candidate)) return emptyCellLink();
  const preflightReason = unsafeLinkPreflightReason(candidate);
  if (preflightReason) {
    return blockedCellLink(text, protocolNameFromText(candidate), preflightReason);
  }
  let parsed;
  try {
    parsed = new URL(candidate);
  } catch {
    return blockedCellLink(text, protocolNameFromText(candidate), "Invalid link");
  }
  const protocol = parsed.protocol.toLowerCase();
  const protocolName = protocol.replace(/:$/, "");
  const parsedReason = unsafeParsedLinkReason(parsed);
  if (parsedReason) {
    return blockedCellLink(text, protocolName, parsedReason);
  }
  if (blockedLinkProtocols.has(protocol)) {
    return blockedCellLink(text, protocolName, "Unsafe protocol");
  }
  if ((protocol === "http:" || protocol === "https:") && !parsed.hostname) {
    return blockedCellLink(text, protocolName, "Missing host");
  }
  if (protocol === "mailto:" && !validMailtoLink(candidate)) {
    return blockedCellLink(text, protocolName, "Invalid email link");
  }
  if (protocol === "file:" && !parsed.hostname && !parsed.pathname) {
    return blockedCellLink(text, protocolName, "Missing file path");
  }
  if (!trustedLinkProtocols.has(protocol) && !customLinkProtocols.includes(protocolName)) {
    return blockedCellLink(text, protocolName, "Protocol not allowed");
  }
  return {
    href: candidate,
    protocol: protocolName,
    status: trustedLinkProtocols.has(protocol) ? "trusted" : "custom",
    target: linkTargetLabel(parsed, candidate),
  };
}

function unsafeLinkPreflightReason(candidate) {
  if (/[\u0000-\u001f\u007f]/.test(candidate)) return "Contains control characters";
  if (/%(?:0[0-9a-f]|1[0-9a-f]|7f)/i.test(candidate)) return "Contains encoded control characters";
  if (candidate.includes("\\")) return "Contains unsafe path separators";
  return "";
}

function unsafeParsedLinkReason(parsed) {
  const protocol = parsed.protocol.toLowerCase();
  if ((protocol === "http:" || protocol === "https:") && (parsed.username || parsed.password)) {
    return "Contains embedded credentials";
  }
  return "";
}

function emptyCellLink() {
  return { href: "", protocol: "", status: "none", target: "" };
}

function blockedCellLink(text, protocol, reason) {
  return {
    href: "",
    protocol: protocol || "",
    status: "blocked",
    target: protocol ? `${reason}: ${protocol}` : reason || text,
  };
}

function protocolNameFromText(text) {
  const match = String(text ?? "").match(/^([a-z][a-z0-9+.-]*):/i);
  return match ? match[1].toLowerCase() : "";
}

function validMailtoLink(value) {
  const address = value.replace(/^mailto:/i, "").split(/[?#]/, 1)[0];
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/i.test(address);
}

function linkTargetLabel(parsed, fallback) {
  if (parsed.protocol === "mailto:") {
    return safeDecodeURIComponent(parsed.pathname).replace(/^\/+/, "") || fallback;
  }
  if (parsed.protocol === "file:") {
    const prefix = parsed.hostname ? `//${parsed.hostname}` : "";
    return safeDecodeURIComponent(`${prefix}${parsed.pathname}`) || fallback;
  }
  const target = `${parsed.host || ""}${parsed.pathname || ""}${parsed.search || ""}${parsed.hash || ""}`;
  return safeDecodeURIComponent(target || fallback);
}

function safeDecodeURIComponent(value) {
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

function linkUnavailableMessage(link) {
  if (link.status === "blocked") {
    return link.protocol ? `Link blocked: ${link.protocol} protocol` : "Link blocked";
  }
  return "Selected cell is not a link";
}

function openSelectedCellLink() {
  const link = analyzeCellLink(valueAt(selected.row, selected.col));
  openCellLinkWithPrompt(link);
}

function openCellLinkWithPrompt(link) {
  if (!link.href) {
    status.textContent = linkUnavailableMessage(link);
    return;
  }
  if (linkRequiresCustomProtocolPrompt(link)) {
    openCustomLinkConfirmDialog(link);
    return;
  }
  openAuditedCellLink(link);
}

function linkRequiresCustomProtocolPrompt(link) {
  return link.status === "custom" && !confirmedCustomLinkProtocols.includes(link.protocol);
}

function openCustomLinkConfirmDialog(link) {
  if (!linkConfirmDialog) {
    openAuditedCellLink(link);
    return;
  }
  pendingLinkConfirmation = { ...link };
  linkConfirmDialog.hidden = false;
  if (linkConfirmProtocol) {
    linkConfirmProtocol.textContent = tf("link.prompt.protocol", { protocol: link.protocol.toUpperCase() });
  }
  if (linkConfirmTarget) {
    linkConfirmTarget.textContent = tf("link.prompt.target", { target: link.target || link.href });
  }
  if (linkConfirmRemember) {
    linkConfirmRemember.checked = false;
  }
  document.body.dataset.linkConfirmDialog = "open";
  document.body.dataset.linkConfirmProtocol = link.protocol;
  document.body.dataset.linkConfirmHref = link.href;
  document.body.dataset.linkConfirmTarget = link.target || link.href;
  document.body.dataset.linkConfirmRemembered = "false";
  status.textContent = tf("link.status.confirm", { protocol: link.protocol });
  confirmOpenLinkButton?.focus();
}

function closeCustomLinkConfirmDialog() {
  pendingLinkConfirmation = null;
  if (linkConfirmDialog) linkConfirmDialog.hidden = true;
  document.body.dataset.linkConfirmDialog = "closed";
  document.body.dataset.linkConfirmProtocol = "";
  document.body.dataset.linkConfirmHref = "";
  document.body.dataset.linkConfirmTarget = "";
}

function cancelCustomLinkConfirmDialog() {
  closeCustomLinkConfirmDialog();
  status.textContent = t("link.status.canceled");
}

function confirmCustomLinkOpen() {
  const link = pendingLinkConfirmation;
  if (!link?.href) {
    closeCustomLinkConfirmDialog();
    return;
  }
  const remember = Boolean(linkConfirmRemember?.checked);
  if (remember) {
    rememberCustomLinkProtocol(link.protocol);
  }
  closeCustomLinkConfirmDialog();
  document.body.dataset.linkConfirmRemembered = remember ? "true" : "false";
  openAuditedCellLink(link, { remembered: remember });
}

function openAuditedCellLink(link, options = {}) {
  window.__opencsvLastOpenedLink = link.href;
  document.body.dataset.lastOpenedLink = link.href;
  document.body.dataset.lastOpenedLinkProtocol = link.protocol || "";
  window.open(link.href, "_blank", "noopener");
  status.textContent = options.remembered
    ? tf("link.status.remembered", { href: link.href, protocol: link.protocol })
    : tf("link.status.opened", { href: link.href });
}

function cellAtPoint(x, y) {
  if (x < rowNumberWidth || y < headerHeight) return null;
  const row = rowAt(y);
  const col = columnAt(x);
  if (row < 0 || col < 0) return null;
  return { row, col };
}

function urlAtPoint(x, y) {
  return cellLinkAtPoint(x, y).href;
}

function cellLinkAtPoint(x, y) {
  const cell = cellAtPoint(x, y);
  if (!cell) return emptyCellLink();
  return analyzeCellLink(valueAt(cell.row, cell.col));
}

function clearURLHover() {
  if (!hoverURL) return;
  hoverURL = "";
  document.body.dataset.hoverUrl = "";
  syncCanvasHoverChrome();
}

function clearResizeHover() {
  if (hoverResizeColumn < 0) return;
  hoverResizeColumn = -1;
  document.body.dataset.resizeColumn = "";
  syncCanvasHoverChrome();
}

function clearFillHandleHover() {
  if (draggingFillHandle || !hoverFillHandle) return;
  hoverFillHandle = false;
  document.body.dataset.fillHandle = "";
  syncCanvasHoverChrome();
}

function syncCanvasHoverChrome() {
  if (pendingRowReorder?.active) {
    canvas.style.cursor = "grabbing";
    canvas.title = "Move row";
    return;
  }
  if (pendingColumnReorder?.active) {
    canvas.style.cursor = "grabbing";
    canvas.title = "Move column";
    return;
  }
  if (draggingFillHandle || hoverFillHandle) {
    canvas.style.cursor = "crosshair";
    canvas.title = "Fill selection";
    return;
  }
  const resizeColumn = resizingColumn >= 0 ? resizingColumn : hoverResizeColumn;
  if (resizeColumn >= 0 && columns[resizeColumn]) {
    canvas.style.cursor = "col-resize";
    canvas.title = `Resize ${columns[resizeColumn].name || `Column ${resizeColumn + 1}`}`;
    return;
  }
  if (hoverURL) {
    canvas.style.cursor = "pointer";
    canvas.title = hoverURL;
    return;
  }
  if (pendingColumnReorder) {
    canvas.style.cursor = "grab";
    canvas.title = "Drag to move selected column";
    return;
  }
  if (pendingRowReorder) {
    canvas.style.cursor = "grab";
    canvas.title = "Drag to move selected row";
    return;
  }
  canvas.style.cursor = "";
  canvas.removeAttribute("title");
}

function updateURLHover(event) {
  if (draggingFillHandle || draggingSelection || resizingColumn >= 0 || editing || !documentId) {
    clearURLHover();
    return;
  }
  if (!(event.metaKey || event.ctrlKey)) {
    clearURLHover();
    return;
  }
  const rect = canvas.getBoundingClientRect();
  const url = urlAtPoint(event.clientX - rect.left, event.clientY - rect.top);
  if (!url) {
    clearURLHover();
    return;
  }
  hoverURL = url;
  document.body.dataset.hoverUrl = url;
  syncCanvasHoverChrome();
}

function updateCanvasHover(event) {
  if (draggingFillHandle || draggingSelection || resizingColumn >= 0 || pendingColumnReorder?.active || pendingRowReorder?.active || editing || !documentId) {
    clearURLHover();
    clearResizeHover();
    clearFillHandleHover();
    return;
  }
  if (event.metaKey || event.ctrlKey) {
    clearResizeHover();
    clearFillHandleHover();
    updateURLHover(event);
    return;
  }
  clearURLHover();
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  if (fillHandleAt(x, y)) {
    clearResizeHover();
    hoverFillHandle = true;
    document.body.dataset.fillHandle = selectionSummary();
    syncCanvasHoverChrome();
    return;
  }
  clearFillHandleHover();
  const col = headerResizeColumnAt(x, y);
  if (col < 0) {
    if (pendingColumnReorder && !pendingColumnReorder.active) {
      pendingColumnReorder = null;
      syncCanvasHoverChrome();
    }
    clearResizeHover();
    return;
  }
  hoverResizeColumn = col;
  document.body.dataset.resizeColumn = String(col);
  if (pendingColumnReorder && !pendingColumnReorder.active) {
    pendingColumnReorder = null;
  }
  syncCanvasHoverChrome();
}

function openGridContextMenu(event) {
  if (!documentId || columns.length === 0 || totalRows === 0) return;
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  const target = gridContextTargetAt(x, y);
  event.preventDefault();
  if (!target) {
    closeGridContextMenu();
    return;
  }
  if (editing) commitEditor();
  closeColumnMenu();
  selectGridContextTarget(target);
  renderGridContextMenu(target, event.clientX, event.clientY);
  draw();
}

function gridContextTargetAt(x, y) {
  if (y < 0 || x < 0 || x > canvas.clientWidth || y > canvas.clientHeight) return null;
  if (y < headerHeight && x >= rowNumberWidth) {
    const col = columnAt(x);
    return col >= 0 ? { kind: "column", col } : null;
  }
  if (x < rowNumberWidth && y >= headerHeight) {
    const row = rowAt(y);
    return row >= 0 ? { kind: "row", row } : null;
  }
  const cell = cellAtPoint(x, y);
  return cell ? { kind: "cell", ...cell } : null;
}

function selectGridContextTarget(target) {
  if (target.kind === "column") {
    selectColumnRange(target.col);
    return;
  }
  if (target.kind === "row") {
    selectRowRange(target.row);
    return;
  }
  if (selectionCellCount() > 1 && activeSelectionRanges().some((range) => rangeContainsCell(range, target.row, target.col))) {
    return;
  }
  selected = { row: target.row, col: target.col };
  selectionAnchor = { ...selected };
}

function renderGridContextMenu(target, clientX, clientY) {
  const heading = document.createElement("strong");
  heading.textContent = gridContextHeading(target);
  gridContextMenu.replaceChildren(heading);
  const contextCommandIDs = [];

  const addSeparator = () => gridContextMenu.append(document.createElement("hr"));
  const addItem = (label, action, options = {}) => {
    const button = document.createElement("button");
    button.type = "button";
    button.setAttribute("role", "menuitem");
    button.textContent = label;
    if (options.commandID) button.dataset.command = options.commandID;
    button.disabled = Boolean(options.disabled);
    button.addEventListener("click", async () => {
      closeGridContextMenu();
      await action();
    });
    gridContextMenu.append(button);
    return button;
  };
  const addCommand = (commandID, label = "") => {
    const command = commandsById.get(commandID);
    contextCommandIDs.push(commandID);
    const fallbackLabel = label || commandTitle(command || { id: commandID });
    const localizedLabel = command && label === command.title
      ? localizedCommandTitle(command)
      : (translationOrEmpty(`context.${commandID}`) || (label ? fallbackLabel : commandTitle(command || { id: commandID })));
    addItem(localizedLabel, () => executeCommand(commandID), {
      commandID,
      disabled: !command || isCommandDisabled(command),
    });
  };

  addCommand("edit.copy", "Copy");
  addCommand("edit.copyWithHeaders", "Copy with Headers");
  addCommand("edit.copyCsv", "Copy as CSV");
  addCommand("edit.copyMarkdown", "Copy as Markdown");
  addCommand("edit.copyJson", "Copy as JSON");
  addCommand("edit.copyHtml", "Copy as HTML");
  addCommand("edit.copySql", "Copy as SQL");
  addCommand("edit.copyLatex", "Copy as LaTeX");
  addCommand("edit.copyRaw", "Copy Raw");
  if (target.kind !== "column") {
    addCommand("edit.copyRow", "Copy Row");
    addCommand("edit.copyRowWithHeaders", "Copy Row with Headers");
    for (const commandID of copyRowFormatCommands.keys()) addCommand(commandID);
  }
  if (target.kind !== "row") {
    addCommand("edit.copyColumn", "Copy Column");
    addCommand("edit.copyColumnWithHeader", "Copy Column with Header");
    for (const commandID of copyColumnFormatCommands.keys()) addCommand(commandID);
  }

  const url = target.kind === "cell" ? normalizeURL(valueAt(selected.row, selected.col)) : "";
  if (url) {
    addCommand("cell.openLink", "Open Link");
  }

  if (target.kind === "cell") {
    addSeparator();
    addCommand("format.toggleCellQuote", "Toggle Cell Quote");
    addCommand("format.quoteSelection", "Quote Selection");
    addCommand("format.noQuoteSelection", "No Quote Selection");
    addCommand("format.clearSelectionQuote", "Clear Selection Quote");
    addSeparator();
    addCommand("text.trim", "Trim");
    addCommand("text.uppercase", "Uppercase");
    addCommand("text.lowercase", "Lowercase");
    addCommand("text.snakeCase", "Snake Case");
    addCommand("text.kebabCase", "Kebab Case");
    addCommand("text.camelCase", "Camel Case");
    addCommand("text.pascalCase", "Pascal Case");
    addCommand("text.titleCase", "Title Case");
    addCommand("text.fullWidth", "Full Width");
    addCommand("text.halfWidth", "Half Width");
    addCommand("text.normalizeNFC", "Normalize NFC");
    addCommand("text.normalizeNFD", "Normalize NFD");
    addCommand("text.normalizeNFKC", "Normalize NFKC");
    addCommand("text.normalizeNFKD", "Normalize NFKD");
    addSeparator();
    addCommand("cell.fillDown", "Fill Down");
    addCommand("cell.fillRight", "Fill Right");
    addCommand("cell.fillUp", "Fill Up");
    addCommand("cell.fillLeft", "Fill Left");
    addCommand("cell.swapUp", "Swap Up");
    addCommand("cell.swapDown", "Swap Down");
    addCommand("cell.swapLeft", "Swap Left");
    addCommand("cell.swapRight", "Swap Right");
    addCommand("cell.joinHorizontal", "Join Cells");
    addCommand("cell.splitHorizontal", "Split Cell");
    addCommand("cell.joinVertical", "Join Cells Vertically");
    addCommand("cell.splitVertical", "Split Cell Vertically");
    addSeparator();
    addCommand("selection.insertAggregateRow", "Add Aggregate Row");
    addCommand("selection.insertSumRow", "Add Sum Row");
    addCommand("selection.insertAverageRow", "Add Average Row");
    addCommand("selection.insertMinRow", "Add Min Row");
    addCommand("selection.insertMaxRow", "Add Max Row");
    addCommand("selection.insertCountRow", "Add Count Row");
    addSeparator();
    addCommand("table.transpose", "Transpose");
    addCommand("table.dedupeRows", "Remove Duplicate Rows");
    addCommand("table.dedupeColumns", "Remove Duplicate Columns");
    addCommand("column.filterBySelectedValue", "Filter by Cell");
    addCommand("column.sortAscending", "Sort Column Ascending");
    addCommand("column.sortDescending", "Sort Column Descending");
  } else if (target.kind === "column") {
    addSeparator();
    addCommand("format.toggleColumnQuote", "Toggle Column Quote");
    addSeparator();
    for (const spec of columnTextTransformSpecs) {
      addCommand(spec.id, spec.label);
    }
    addSeparator();
    addCommand("column.sortAscending", "Sort Ascending");
    addCommand("column.sortDescending", "Sort Descending");
    addCommand("column.sortLengthAscending", "Sort Length Ascending");
    addCommand("column.sortLengthDescending", "Sort Length Descending");
    addCommand("column.addSortAscending", "Then Sort Ascending");
    addCommand("column.addSortDescending", "Then Sort Descending");
    addCommand("column.filterBySelectedValue", "Filter by Selected Value");
    addCommand("column.filterBlanks", "Show Blanks");
    addCommand("column.filterNonBlanks", "Show Non-Blanks");
    addCommand("view.filterBuilder", "Filter Builder");
  }

  if (target.kind !== "column") {
    addSeparator();
    addCommand("row.insert");
    addCommand("row.delete");
    addCommand("row.cut");
    addCommand("row.paste");
    addCommand("row.pasteBelow");
    addCommand("row.hide");
    addCommand("row.showAll");
  }
  if (target.kind !== "row") {
    addSeparator();
    addCommand("column.insert");
    addCommand("column.delete");
    addCommand("column.cut");
    addCommand("column.paste");
    addCommand("column.pasteRight");
    addCommand("column.hide");
    addCommand("column.showAll");
  }
  addSeparator();
  addCommand("view.clear");

  gridContextMenu.hidden = false;
  const menuRect = gridContextMenu.getBoundingClientRect();
  const left = Math.min(Math.max(8, clientX), window.innerWidth - menuRect.width - 8);
  const top = Math.min(Math.max(8, clientY), window.innerHeight - menuRect.height - 8);
  gridContextMenu.style.left = `${left}px`;
  gridContextMenu.style.top = `${top}px`;
  gridContextMenu.dataset.target = target.kind;
  document.body.dataset.gridContextTarget = target.kind;
  document.body.dataset.gridContextCommands = contextCommandIDs.join(",");
}

function gridContextHeading(target) {
  if (target.kind === "column") return `Column: ${columns[target.col]?.name ?? target.col + 1}`;
  if (target.kind === "row") return `Row ${target.row + 1}`;
  const column = columns[target.col]?.name ?? `Column ${target.col + 1}`;
  return `Cell ${target.row + 1}, ${column}`;
}

function closeGridContextMenu() {
  if (gridContextMenu.hidden) return;
  gridContextMenu.hidden = true;
  gridContextMenu.replaceChildren();
  gridContextMenu.dataset.target = "";
  document.body.dataset.gridContextTarget = "";
  document.body.dataset.gridContextCommands = "";
  document.body.dataset.sqlResultContext = "";
  document.body.dataset.sqlResultContextCommands = "";
}

function normalizedSelection() {
  const startRow = Math.max(0, Math.min(selectionAnchor.row, selected.row));
  const endRow = Math.max(0, Math.max(selectionAnchor.row, selected.row));
  const startCol = Math.max(0, Math.min(selectionAnchor.col, selected.col));
  const endCol = Math.max(0, Math.max(selectionAnchor.col, selected.col));
  return {
    startRow,
    endRow: Math.min(totalRows - 1, endRow),
    startCol,
    endCol: Math.min(columns.length - 1, endCol),
  };
}

function normalizedSmartSelectionRanges() {
  if (!smartSelectionRanges.length) return [];
  return smartSelectionRanges
    .map((range) => normalizeRange(range))
    .filter(Boolean);
}

function activeSelectionRanges() {
  const primary = normalizedSelection();
  const smart = normalizedSmartSelectionRanges();
  if (smart.some((range) => rangeContainsRange(range, primary))) return dedupeRanges(smart);
  return dedupeRanges([primary, ...smart].filter(Boolean));
}

function normalizeRange(range) {
  if (!documentId || columns.length === 0 || totalRows === 0 || !range) return null;
  const startRow = Math.max(0, range.startRow);
  const endRow = Math.min(totalRows - 1, range.endRow);
  const startCol = Math.max(0, range.startCol);
  const endCol = Math.min(columns.length - 1, range.endCol);
  if (endRow < startRow || endCol < startCol) return null;
  return { startRow, endRow, startCol, endCol };
}

function rangeKey(range) {
  return `${range.startRow}:${range.startCol}-${range.endRow}:${range.endCol}`;
}

function dedupeRanges(ranges) {
  const seen = new Set();
  const result = [];
  for (const range of ranges) {
    const normalized = normalizeRange(range);
    if (!normalized) continue;
    const key = rangeKey(normalized);
    if (seen.has(key)) continue;
    seen.add(key);
    result.push(normalized);
  }
  return result;
}

function rangeCellCount(range) {
  return (range.endRow - range.startRow + 1) * (range.endCol - range.startCol + 1);
}

function selectionCellCount() {
  if (!documentId || columns.length === 0 || totalRows === 0) return 0;
  return activeSelectionRanges().reduce((count, range) => count + rangeCellCount(range), 0);
}

function selectionSummary(range = null) {
  if (!documentId || columns.length === 0 || totalRows === 0) return "";
  const ranges = range ? [range] : activeSelectionRanges();
  const base = ranges.map(rangeKey).join(";");
  return includeHeadersInSelection ? `${base}:headers` : base;
}

function aggregateStatsKey() {
  return [
    documentId || "",
    windowGeneration,
    sourceRowCount,
    columns.length,
    selectionSummary(),
  ].join("|");
}

function currentSelectionAggregateStats() {
  const key = aggregateStatsKey();
  if (aggregateStatsCacheKey === key && aggregateStatsCache) return aggregateStatsCache;
  aggregateStatsCache = selectionAggregateStats();
  aggregateStatsCacheKey = key;
  return aggregateStatsCache;
}

function selectionAggregateStats() {
  const startedAt = Date.now();
  const count = selectionCellCount();
  const stats = {
    count,
    numericCount: 0,
    limited: false,
    scannedCells: 0,
    elapsedMs: 0,
    sum: 0,
    avg: 0,
    min: 0,
    max: 0,
  };
  if (count <= 1) return stats;
  if (count > aggregateScanCellLimit) {
    stats.limited = true;
    stats.elapsedMs = Date.now() - startedAt;
    return stats;
  }
  let minValue = Number.POSITIVE_INFINITY;
  let maxValue = Number.NEGATIVE_INFINITY;
  for (const range of activeSelectionRanges()) {
    if (range.endRow >= sourceRowCount) {
      stats.limited = true;
      stats.elapsedMs = Date.now() - startedAt;
      return stats;
    }
    for (let row = range.startRow; row <= range.endRow; row++) {
      for (let col = range.startCol; col <= range.endCol; col++) {
        stats.scannedCells++;
        const raw = valueAt(row, col).trim();
        if (raw === "") continue;
        const value = Number(raw);
        if (!Number.isFinite(value)) continue;
        stats.numericCount++;
        stats.sum += value;
        minValue = Math.min(minValue, value);
        maxValue = Math.max(maxValue, value);
      }
    }
  }
  if (stats.numericCount) {
    stats.avg = stats.sum / stats.numericCount;
    stats.min = minValue;
    stats.max = maxValue;
  }
  stats.elapsedMs = Date.now() - startedAt;
  return stats;
}

function selectionAggregateSummary(stats = currentSelectionAggregateStats()) {
  if (stats.count <= 1) return "";
  if (stats.limited) {
    const headerNote = includeHeadersInSelection ? ", headers included" : "";
    return ` - selected ${formatCount(stats.count)} cells${headerNote}`;
  }
  if (!stats.numericCount) {
    return ` - selected ${formatCount(stats.count)} cells`;
  }
  return ` - selected ${formatCount(stats.count)} cells, sum ${formatNumber(stats.sum)}, avg ${formatNumber(stats.avg)}, min ${formatNumber(stats.min)}, max ${formatNumber(stats.max)}`;
}

function aggregateMetricOptions(stats = currentSelectionAggregateStats()) {
  if (stats.count <= 1) return [];
  if (!stats.numericCount) return ["count"];
  return ["sum", "avg", "min", "max", "count"];
}

function normalizeAggregateMetric(stats = currentSelectionAggregateStats()) {
  const options = aggregateMetricOptions(stats);
  if (!options.length) return "";
  if (!options.includes(activeAggregateMetric)) activeAggregateMetric = options[0];
  return activeAggregateMetric;
}

function activeAggregate(stats = currentSelectionAggregateStats()) {
  const metric = normalizeAggregateMetric(stats);
  if (!metric) return null;
  const value = metric === "count" ? stats.count : stats[metric];
  return {
    metric,
    label: aggregateMetricLabel(metric),
    value,
    formatted: formatNumber(value),
  };
}

function aggregateMetricLabel(metric) {
  return {
    count: "Count",
    sum: "Sum",
    avg: "Avg",
    min: "Min",
    max: "Max",
  }[metric] ?? metric;
}

function updateAggregateStatus(stats = currentSelectionAggregateStats()) {
  document.body.dataset.aggregateCellCount = String(stats.count);
  document.body.dataset.aggregateNumericCount = String(stats.numericCount);
  document.body.dataset.aggregateLimited = String(Boolean(stats.limited));
  document.body.dataset.aggregateScannedCells = String(stats.scannedCells);
  document.body.dataset.aggregateElapsedMs = String(stats.elapsedMs);
  const aggregate = activeAggregate(stats);
  if (!aggregate) {
    aggregateStatusButton.hidden = true;
    aggregateStatusButton.disabled = true;
    aggregateStatusButton.textContent = "Count 0";
    aggregateStatusButton.removeAttribute("aria-label");
    delete document.body.dataset.aggregateActive;
    return;
  }
  aggregateStatusButton.hidden = false;
  aggregateStatusButton.disabled = false;
  aggregateStatusButton.textContent = `${aggregate.label} ${aggregate.formatted}`;
  aggregateStatusButton.title = "Click to switch aggregate; Command-click to copy.";
  aggregateStatusButton.setAttribute("aria-label", `${aggregate.label} ${aggregate.formatted}`);
  document.body.dataset.aggregateActive = `${aggregate.metric}:${aggregate.formatted}`;
}

function cycleAggregateMetric() {
  const stats = selectionAggregateStats();
  const options = aggregateMetricOptions(stats);
  if (!options.length) {
    status.textContent = "Select multiple cells to aggregate.";
    return;
  }
  const current = normalizeAggregateMetric(stats);
  const index = Math.max(0, options.indexOf(current));
  activeAggregateMetric = options[(index + 1) % options.length];
  const aggregate = activeAggregate(stats);
  updateAggregateStatus(stats);
  draw();
  status.textContent = `${aggregate.label} ${aggregate.formatted}`;
}

async function copyActiveAggregate() {
  const aggregate = activeAggregate();
  if (!aggregate) {
    status.textContent = "Select multiple cells to copy an aggregate.";
    return;
  }
  const text = aggregate.formatted;
  window.__opencsvLastAggregateText = text;
  document.body.dataset.aggregatePreview = `${aggregate.metric}:${text}`;
  try {
    await writeClipboardText(text);
    status.textContent = `Copied ${aggregate.label} ${text}`;
  } catch (error) {
    status.textContent = `Aggregate copy failed: ${error.message}`;
  }
}

function insertAggregateRow(metricOverride = "") {
  if (!documentId) return;
  const range = normalizedSelection();
  if (range.endRow >= sourceRowCount) {
    status.textContent = "Select data rows before adding an aggregate row.";
    return;
  }
  const stats = selectionAggregateStats();
  const options = aggregateMetricOptions(stats);
  let metric = "";
  if (metricOverride) {
    if (!options.includes(metricOverride)) {
      status.textContent = metricOverride === "count" ? "Select cells to add a count row." : "Select numeric cells to add an aggregate row.";
      return;
    }
    metric = metricOverride;
  } else {
    metric = normalizeAggregateMetric(stats);
  }
  if (!metric) {
    status.textContent = metricOverride === "count" ? "Select cells to add a count row." : "Select numeric cells to add an aggregate row.";
    return;
  }
  const row = aggregateRowForSelection(range, metric);
  if (!row) return;
  aggregateRows.push(row);
  syncTotalRows();
  selected = { row: totalRows - 1, col: range.startCol };
  selectionAnchor = { ...selected };
  document.body.dataset.aggregateRowPreview = `${metric}:${row.values.join("|")}`;
  draw();
  status.textContent = `Added ${aggregateMetricLabel(metric)} aggregate row`;
}

function aggregateRowForSelection(range, metric) {
  const values = Array.from({ length: allColumns.length }, () => "");
  const labelCol = aggregateLabelColumn(range);
  if (labelCol >= 0) values[sourceColFor(labelCol)] = aggregateMetricLabel(metric);
  let populated = 0;
  for (let col = range.startCol; col <= range.endCol; col++) {
    const value = aggregateValueForColumn(range, col, metric);
    if (value === "") continue;
    values[sourceColFor(col)] = value;
    populated++;
  }
  if (populated === 0) {
    status.textContent = "No numeric values in the selected columns.";
    return null;
  }
  return {
    metric,
    label: aggregateMetricLabel(metric),
    values,
  };
}

function aggregateLabelColumn(range) {
  if (range.startCol > 0) return 0;
  if (range.endCol + 1 < columns.length) return range.endCol + 1;
  return -1;
}

function aggregateValueForColumn(range, col, metric) {
  if (metric === "count") return String(range.endRow - range.startRow + 1);
  const values = [];
  for (let row = range.startRow; row <= range.endRow; row++) {
    const raw = valueAt(row, col).trim();
    if (raw === "") continue;
    const value = Number(raw);
    if (Number.isFinite(value)) values.push(value);
  }
  if (!values.length) return "";
  if (metric === "sum") return formatNumber(values.reduce((sum, value) => sum + value, 0));
  if (metric === "avg") return formatNumber(values.reduce((sum, value) => sum + value, 0) / values.length);
  if (metric === "min") return formatNumber(Math.min(...values));
  if (metric === "max") return formatNumber(Math.max(...values));
  return "";
}

function removeSelectedAggregateRow() {
  const index = aggregateRowIndex(selected.row);
  if (index < 0) return false;
  aggregateRows.splice(index, 1);
  syncTotalRows();
  selected = {
    row: Math.max(0, Math.min(totalRows - 1, selected.row - 1)),
    col: selected.col,
  };
  selectionAnchor = { ...selected };
  document.body.dataset.aggregateRowPreview = `removed:${aggregateRows.length}`;
  draw();
  status.textContent = "Removed aggregate row";
  return true;
}

function openCommandPalette() {
  paletteOpen = true;
  paletteIndex = 0;
  paletteInput.value = "";
  commandPalette.hidden = false;
  renderCommandPalette();
  paletteInput.focus();
}

function closeCommandPalette() {
  paletteOpen = false;
  commandPalette.hidden = true;
  paletteList.replaceChildren();
  canvas.focus();
}

function openKeyboardShortcuts() {
  shortcutEditorOpen = true;
  recordingShortcutID = "";
  pendingShortcutConflict = null;
  shortcutsSearch.value = "";
  shortcutImportPanel.hidden = true;
  shortcutImportText.value = "";
  syncShortcutImportControls();
  shortcutsPanel.hidden = false;
  renderKeyboardShortcuts();
  shortcutsSearch.focus();
}

function closeKeyboardShortcuts() {
  shortcutEditorOpen = false;
  recordingShortcutID = "";
  pendingShortcutConflict = null;
  shortcutsPanel.hidden = true;
  shortcutImportPanel.hidden = true;
  shortcutImportText.value = "";
  document.body.dataset.shortcutImportPanel = "closed";
  shortcutsList.replaceChildren();
  canvas.focus();
}

function renderKeyboardShortcuts() {
  syncResetAllShortcutsButton();
  const query = shortcutsSearch.value.trim().toLowerCase();
  const commands = Array.from(commandsById.values())
    .filter((command) => command.implemented !== false)
    .filter((command) => !query || commandFuzzyScore(command, query) > 0 || effectiveShortcut(command).toLowerCase().includes(query))
    .sort((left, right) => commandTitle(left).localeCompare(commandTitle(right)));
  if (!commands.length) {
    shortcutsList.replaceChildren(emptyPaletteNode("No matching commands"));
    return;
  }

  const fragment = document.createDocumentFragment();
  for (const command of commands) {
    const row = document.createElement("div");
    row.className = command.id === recordingShortcutID ? "shortcut-row recording" : "shortcut-row";
    row.setAttribute("role", "listitem");

    const meta = document.createElement("div");
    const title = document.createElement("strong");
    title.textContent = localizedCommandTitle(command);
    const description = document.createElement("span");
    description.textContent = localizedCommandDescription(command);
    meta.append(title, description);

    const binding = effectiveShortcutBinding(command);
    const shortcut = document.createElement("kbd");
    shortcut.textContent = command.id === recordingShortcutID
      ? "Press shortcut"
      : binding.shortcut;
    if (binding.when) shortcut.title = `When ${shortcutContextLabel(binding.when)}`;

    const contextSelect = document.createElement("select");
    contextSelect.className = "shortcut-context";
    contextSelect.setAttribute("aria-label", `${localizedCommandTitle(command)} shortcut context`);
    for (const optionValue of shortcutContextOptions) {
      const option = document.createElement("option");
      option.value = optionValue.id;
      option.textContent = optionValue.label;
      contextSelect.append(option);
    }
    if (binding.when && !shortcutContextOptions.some((option) => option.id === binding.when)) {
      const option = document.createElement("option");
      option.value = binding.when;
      option.textContent = shortcutContextLabel(binding.when);
      contextSelect.append(option);
    }
    contextSelect.value = binding.when;
    contextSelect.addEventListener("change", () => updateShortcutContext(command, contextSelect.value));

    const actions = document.createElement("div");
    actions.className = "shortcut-actions";
    const recordButton = document.createElement("button");
    recordButton.type = "button";
    recordButton.textContent = command.id === recordingShortcutID ? "Cancel" : "Record";
    recordButton.addEventListener("click", () => {
      recordingShortcutID = command.id === recordingShortcutID ? "" : command.id;
      pendingShortcutConflict = null;
      renderKeyboardShortcuts();
    });
    const resetButton = document.createElement("button");
    resetButton.type = "button";
    resetButton.textContent = "Reset";
    resetButton.disabled = !customShortcuts.has(command.id);
    resetButton.addEventListener("click", () => resetCustomShortcut(command.id));
    actions.append(recordButton, resetButton);
    if (pendingShortcutConflict?.commandID === command.id) {
      const conflictCommand = commandsById.get(pendingShortcutConflict.conflictID);
      const conflictNote = document.createElement("span");
      conflictNote.className = "shortcut-conflict";
      conflictNote.textContent = conflictCommand
        ? `${pendingShortcutConflict.shortcut} is assigned to ${localizedCommandTitle(conflictCommand)}`
        : `${pendingShortcutConflict.shortcut} is already assigned`;
      const replaceButton = document.createElement("button");
      replaceButton.type = "button";
      replaceButton.textContent = "Replace";
      replaceButton.addEventListener("click", replaceShortcutConflict);
      meta.append(conflictNote);
      actions.append(replaceButton);
    }

    row.append(meta, shortcut, contextSelect, actions);
    fragment.append(row);
  }
  shortcutsList.replaceChildren(fragment);
}

function effectiveShortcut(command) {
  return effectiveShortcutBinding(command).shortcut;
}

function effectiveShortcutBinding(command) {
  if (customShortcuts.has(command.id)) {
    return customShortcuts.get(command.id);
  }
  return { shortcut: command.shortcut || "", when: defaultShortcutWhen(command) };
}

function defaultShortcutWhen(command) {
  if (!command?.scope || command.scope === "global") return "";
  if (command.scope === "grid") return "grid";
  if (command.scope === "view") return "document";
  return "document";
}

function defaultShortcutActive(commandID) {
  return !customShortcuts.has(commandID);
}

function loadCustomShortcuts() {
  try {
    const parsed = JSON.parse(localStorage.getItem(customShortcutsStorageKey) || "{}");
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) return new Map();
    const bindings = new Map();
    for (const [id, value] of Object.entries(parsed)) {
      const binding = normalizeShortcutBinding(value);
      if (typeof id === "string" && binding) bindings.set(id, binding);
    }
    return bindings;
  } catch {
    return new Map();
  }
}

function saveCustomShortcuts() {
  const value = serializedShortcutMap(customShortcuts);
  localStorage.setItem(customShortcutsStorageKey, JSON.stringify(value));
  document.body.dataset.shortcutCount = String(customShortcuts.size);
  document.body.dataset.shortcutPreview = JSON.stringify(value);
  syncResetAllShortcutsButton();
  applyCommandMetadata();
  renderCommandPalette();
  if (shortcutEditorOpen) renderKeyboardShortcuts();
}

function syncResetAllShortcutsButton() {
  resetAllShortcutsButton.disabled = customShortcuts.size === 0;
  resetAllShortcutsButton.title = customShortcuts.size === 0
    ? "Keyboard shortcuts are using defaults"
    : `Reset ${formatCount(customShortcuts.size)} custom shortcuts`;
}

function applyShortcutPreset() {
  const preset = shortcutPresets.find((item) => item.id === shortcutPresetSelect.value) ?? shortcutPresets[0];
  const result = shortcutMapFromEntries(Object.entries(preset.shortcuts));
  customShortcuts = result.shortcuts;
  pendingShortcutConflict = null;
  saveCustomShortcuts();
  document.body.dataset.shortcutPreset = `${preset.id}:${result.applied}:${result.skipped}`;
  status.textContent = `Applied ${preset.name} shortcut preset: ${formatCount(result.applied)} shortcuts${result.skipped ? `, skipped ${formatCount(result.skipped)}` : ""}.`;
}

function shortcutMapFromEntries(entries) {
  const next = new Map();
  let applied = 0;
  let skipped = 0;
  for (const [commandID, rawBinding] of entries) {
    const command = commandsById.get(commandID);
    const binding = normalizeShortcutBinding(rawBinding);
    if (!binding || !command || command.implemented === false || !shortcutAllowed(binding.shortcut)) {
      skipped++;
      continue;
    }
    const conflict = commandForShortcutWithMap(binding.shortcut, binding.when, commandID, next);
    if (conflict) {
      skipped++;
      continue;
    }
    if (binding.shortcut !== command.shortcut || binding.when !== defaultShortcutWhen(command)) {
      next.set(commandID, binding);
    }
    applied++;
  }
  return { shortcuts: next, applied, skipped };
}

function exportShortcuts() {
  const shortcuts = serializedShortcutMap(customShortcuts);
  const payload = {
    kind: shortcutExportKind,
    version: 1,
    exportedAt: new Date().toISOString(),
    shortcuts,
  };
  const text = `${JSON.stringify(payload, null, 2)}\n`;
  window.__opencsvLastShortcutExportText = text;
  document.body.dataset.shortcutExport = String(Object.keys(shortcuts).length);
  document.body.dataset.shortcutExportPreview = text.slice(0, 500);
  triggerDownload(text, "opencsv-shortcuts.json", "application/json;charset=utf-8");
  status.textContent = `Exported ${formatCount(Object.keys(shortcuts).length)} shortcuts.`;
}

function showShortcutImportPanel() {
  openKeyboardShortcuts();
  recordingShortcutID = "";
  shortcutImportPanel.hidden = false;
  shortcutImportText.value = "";
  document.body.dataset.shortcutImportPanel = "open";
  syncShortcutImportControls();
  shortcutImportText.focus();
}

function cancelShortcutImport() {
  shortcutImportPanel.hidden = true;
  shortcutImportText.value = "";
  document.body.dataset.shortcutImportPanel = "closed";
  syncShortcutImportControls();
  shortcutsSearch.focus();
}

function syncShortcutImportControls() {
  if (!shortcutImportText) return;
  applyShortcutImportButton.disabled = shortcutImportText.value.trim() === "";
}

function applyShortcutImport() {
  try {
    const imported = parseShortcutImport(shortcutImportText.value);
    const result = mergeImportedShortcuts(imported);
    pendingShortcutConflict = null;
    saveCustomShortcuts();
    shortcutImportPanel.hidden = true;
    shortcutImportText.value = "";
    document.body.dataset.shortcutImport = `${result.imported}:${result.skipped}`;
    document.body.dataset.shortcutImportPanel = "closed";
    status.textContent = `Imported ${formatCount(result.imported)} shortcuts${result.skipped ? `, skipped ${formatCount(result.skipped)}` : ""}.`;
  } catch (error) {
    status.textContent = `Shortcut import failed: ${error.message}`;
    document.body.dataset.shortcutImport = "error";
  }
  syncShortcutImportControls();
}

function parseShortcutImport(text) {
  let parsed;
  try {
    parsed = JSON.parse(text);
  } catch {
    throw new Error("invalid JSON");
  }
  const source = parsed?.shortcuts && typeof parsed.shortcuts === "object" && !Array.isArray(parsed.shortcuts)
    ? parsed.shortcuts
    : parsed;
  if (!source || typeof source !== "object" || Array.isArray(source)) {
    throw new Error("missing shortcuts object");
  }
  return Object.entries(source)
    .filter(([commandID, value]) => typeof commandID === "string" && normalizeShortcutBinding(value))
    .map(([commandID, value]) => [commandID, normalizeShortcutBinding(value)]);
}

function normalizeShortcutText(shortcut) {
  return shortcut
    .split("+")
    .map((part) => part.trim())
    .filter(Boolean)
    .join("+");
}

function normalizeShortcutBinding(value) {
  if (typeof value === "string") {
    return { shortcut: normalizeShortcutText(value), when: "" };
  }
  if (!value || typeof value !== "object" || Array.isArray(value) || typeof value.shortcut !== "string") {
    return null;
  }
  return {
    shortcut: normalizeShortcutText(value.shortcut),
    when: normalizeShortcutWhen(value.when),
  };
}

function normalizeShortcutWhen(value) {
  const when = String(value || "").trim().toLowerCase();
  if (!when) return "";
  if (shortcutContextOptions.some((option) => option.id === when)) return when;
  const parsed = parseShortcutWhenExpression(when);
  return parsed ? serializeShortcutWhenExpression(parsed) : "";
}

function shortcutContextLabel(when) {
  const normalized = normalizeShortcutWhen(when);
  return shortcutContextOptions.find((option) => option.id === normalized)?.label
    ?? (normalized ? `Expression: ${normalized}` : "Always");
}

function serializedShortcutMap(map) {
  return Object.fromEntries(
    Array.from(map.entries()).map(([commandID, binding]) => [commandID, serializedShortcutBinding(binding)]),
  );
}

function serializedShortcutBinding(binding) {
  if (!binding?.when) return binding?.shortcut ?? "";
  return { shortcut: binding.shortcut, when: binding.when };
}

function updateShortcutContext(command, whenValue) {
  const when = normalizeShortcutWhen(whenValue);
  const binding = effectiveShortcutBinding(command);
  const next = { shortcut: binding.shortcut, when };
  if (!next.shortcut) {
    customShortcuts.delete(command.id);
    saveCustomShortcuts();
    status.textContent = `Reset shortcut context for ${localizedCommandTitle(command)}`;
    return;
  }
  const conflict = commandForShortcutWithMap(next.shortcut, next.when, command.id, customShortcuts);
  if (conflict) {
    status.textContent = `${next.shortcut} already overlaps ${conflict.title} in ${shortcutContextLabel(next.when)}.`;
    renderKeyboardShortcuts();
    return;
  }
  if (next.shortcut === command.shortcut && next.when === defaultShortcutWhen(command)) {
    customShortcuts.delete(command.id);
  } else {
    customShortcuts.set(command.id, next);
  }
  document.body.dataset.shortcutContext = `${command.id}:${next.when || "always"}`;
  saveCustomShortcuts();
  status.textContent = `${localizedCommandTitle(command)} shortcut context: ${shortcutContextLabel(next.when)}.`;
}

function mergeImportedShortcuts(imported) {
  let importedCount = 0;
  let skipped = 0;
  const next = new Map(customShortcuts);
  for (const [commandID, binding] of imported) {
    const command = commandsById.get(commandID);
    if (binding.shortcut === "") {
      if (!command || command.implemented === false) {
        skipped++;
        continue;
      }
      next.set(commandID, binding);
      importedCount++;
      continue;
    }
    if (!command || command.implemented === false || !shortcutAllowed(binding.shortcut)) {
      skipped++;
      continue;
    }
    const conflict = commandForShortcutWithMap(binding.shortcut, binding.when, commandID, next);
    if (conflict) {
      skipped++;
      continue;
    }
    if (binding.shortcut === command.shortcut && binding.when === defaultShortcutWhen(command)) {
      next.delete(commandID);
    } else {
      next.set(commandID, binding);
    }
    importedCount++;
  }
  customShortcuts = next;
  return { imported: importedCount, skipped };
}

function commandForShortcutWithMap(shortcut, when, exceptID, customMap) {
  for (const command of commandsById.values()) {
    if (command.id === exceptID) continue;
    const effective = customMap.has(command.id)
      ? customMap.get(command.id)
      : { shortcut: command.shortcut || "", when: defaultShortcutWhen(command) };
    if (effective.shortcut === shortcut && shortcutContextsOverlap(when, effective.when)) return command;
  }
  return null;
}

function resetCustomShortcut(commandID) {
  const command = commandsById.get(commandID);
  if (!command) return;
  pendingShortcutConflict = null;
  customShortcuts.delete(commandID);
  saveCustomShortcuts();
  status.textContent = `Reset shortcut for ${localizedCommandTitle(command)}`;
}

function resetAllShortcuts() {
  const resetCount = customShortcuts.size;
  recordingShortcutID = "";
  pendingShortcutConflict = null;
  customShortcuts = new Map();
  document.body.dataset.shortcutResetAll = String(resetCount);
  saveCustomShortcuts();
  status.textContent = resetCount
    ? `Reset ${formatCount(resetCount)} custom shortcuts.`
    : "Keyboard shortcuts already use defaults.";
}

function handleShortcutEditorKeydown(event) {
  if (!shortcutEditorOpen) return false;
  if (!recordingShortcutID) {
    if (event.key === "Escape") {
      event.preventDefault();
      closeKeyboardShortcuts();
      return true;
    }
    return false;
  }

  event.preventDefault();
  if (event.key === "Escape") {
    recordingShortcutID = "";
    renderKeyboardShortcuts();
    return true;
  }
  const shortcut = shortcutFromEvent(event);
  if (!shortcut || !shortcutAllowed(shortcut)) {
    status.textContent = "Use a modifier or function key.";
    return true;
  }
  const command = commandsById.get(recordingShortcutID);
  if (!command) return true;
  const when = effectiveShortcutBinding(command).when;
  const conflict = commandForShortcut(shortcut, when, command.id);
  if (conflict) {
    pendingShortcutConflict = { commandID: command.id, conflictID: conflict.id, shortcut, when };
    document.body.dataset.shortcutConflict = `${command.id}:${conflict.id}:${shortcut}:${when || "always"}`;
    renderKeyboardShortcuts();
    status.textContent = `${shortcut} is already assigned to ${conflict.title}`;
    return true;
  }
  assignShortcut(command, { shortcut, when }, customShortcuts);
  recordingShortcutID = "";
  pendingShortcutConflict = null;
  saveCustomShortcuts();
  status.textContent = `Shortcut for ${localizedCommandTitle(command)}: ${shortcut}`;
  return true;
}

function replaceShortcutConflict() {
  if (!pendingShortcutConflict) return;
  const target = commandsById.get(pendingShortcutConflict.commandID);
  const conflict = commandsById.get(pendingShortcutConflict.conflictID);
  if (!target || !conflict) return;
  clearShortcut(conflict, pendingShortcutConflict.shortcut, customShortcuts);
  assignShortcut(target, {
    shortcut: pendingShortcutConflict.shortcut,
    when: pendingShortcutConflict.when,
  }, customShortcuts);
  const message = `${pendingShortcutConflict.shortcut} moved from ${conflict.title} to ${target.title}`;
  recordingShortcutID = "";
  pendingShortcutConflict = null;
  document.body.dataset.shortcutConflict = "resolved";
  saveCustomShortcuts();
  status.textContent = message;
}

function assignShortcut(command, binding, customMap) {
  if (binding.shortcut === command.shortcut && binding.when === defaultShortcutWhen(command)) {
    customMap.delete(command.id);
  } else {
    customMap.set(command.id, binding);
  }
}

function clearShortcut(command, shortcut, customMap) {
  if ((command.shortcut || "") === shortcut) {
    customMap.set(command.id, { shortcut: "", when: "" });
  } else if (customMap.has(command.id) && customMap.get(command.id).shortcut === shortcut) {
    customMap.delete(command.id);
  }
}

function shortcutAllowed(shortcut) {
  return shortcut.includes("Cmd+") || shortcut.includes("Alt+") || /^F\d+$/.test(shortcut);
}

function shortcutContextsOverlap(left, right) {
  const normalizedLeft = normalizeShortcutWhen(left);
  const normalizedRight = normalizeShortcutWhen(right);
  if (!normalizedLeft || !normalizedRight) return true;
  const leftExpression = parseShortcutWhenExpression(normalizedLeft);
  const rightExpression = parseShortcutWhenExpression(normalizedRight);
  if (!leftExpression || !rightExpression) return true;
  return shortcutContextScenarios().some((scenario) => (
    evaluateShortcutWhenExpression(leftExpression, (atom) => Boolean(scenario[atom])) &&
    evaluateShortcutWhenExpression(rightExpression, (atom) => Boolean(scenario[atom]))
  ));
}

function shortcutContextMatches(when, event) {
  const normalized = normalizeShortcutWhen(when);
  if (!normalized) return true;
  const expression = parseShortcutWhenExpression(normalized);
  if (!expression) return true;
  return evaluateShortcutWhenExpression(expression, (atom) => shortcutContextValue(atom, event));
}

function shortcutContextValue(atom, event) {
  switch (atom) {
    case "document":
      return Boolean(documentId);
    case "grid":
      return viewMode === "grid" && !editing && !conditionBuilderOpen && !isTypingTarget(event.target);
    case "sql":
      return viewMode === "sql";
    case "profile":
      return viewMode === "profile";
    case "settings":
      return viewMode === "settings";
    case "diff":
      return viewMode === "diff";
    case "typing":
      return isTypingTarget(event.target);
    default:
      return false;
  }
}

function shortcutContextScenarios() {
  const scenarios = [];
  for (const hasDocument of [false, true]) {
    for (const view of shortcutViewContexts) {
      for (const typing of [false, true]) {
        scenarios.push({
          document: hasDocument,
          grid: view === "grid" && !typing,
          sql: view === "sql",
          profile: view === "profile",
          settings: view === "settings",
          diff: view === "diff",
          typing,
        });
      }
    }
  }
  return scenarios;
}

function parseShortcutWhenExpression(raw) {
  const tokens = tokenizeShortcutWhenExpression(raw);
  if (!tokens) return null;
  let index = 0;

  const parseOr = () => {
    let node = parseAnd();
    if (!node) return null;
    while (tokens[index] === "||") {
      index += 1;
      const right = parseAnd();
      if (!right) return null;
      node = { type: "or", left: node, right };
    }
    return node;
  };

  const parseAnd = () => {
    let node = parseUnary();
    if (!node) return null;
    while (tokens[index] === "&&") {
      index += 1;
      const right = parseUnary();
      if (!right) return null;
      node = { type: "and", left: node, right };
    }
    return node;
  };

  const parseUnary = () => {
    const token = tokens[index];
    if (token === "!") {
      index += 1;
      const operand = parseUnary();
      return operand ? { type: "not", operand } : null;
    }
    if (token === "(") {
      index += 1;
      const node = parseOr();
      if (!node || tokens[index] !== ")") return null;
      index += 1;
      return node;
    }
    if (shortcutContextAtoms.has(token)) {
      index += 1;
      return { type: "atom", value: token };
    }
    return null;
  };

  const expression = parseOr();
  return expression && index === tokens.length ? expression : null;
}

function tokenizeShortcutWhenExpression(raw) {
  const normalized = String(raw || "")
    .trim()
    .toLowerCase()
    .replace(/\band\b/g, "&&")
    .replace(/\bor\b/g, "||")
    .replace(/\bnot\b/g, "!");
  const tokens = [];
  let index = 0;
  while (index < normalized.length) {
    const char = normalized[index];
    if (/\s/.test(char)) {
      index += 1;
      continue;
    }
    if (normalized.startsWith("&&", index) || normalized.startsWith("||", index)) {
      tokens.push(normalized.slice(index, index + 2));
      index += 2;
      continue;
    }
    if (char === "!" || char === "(" || char === ")") {
      tokens.push(char);
      index += 1;
      continue;
    }
    const match = normalized.slice(index).match(/^[a-z]+/);
    if (!match) return null;
    tokens.push(match[0]);
    index += match[0].length;
  }
  return tokens;
}

function serializeShortcutWhenExpression(expression, parentPrecedence = 0) {
  const precedence = shortcutWhenExpressionPrecedence(expression);
  let text = "";
  if (expression.type === "atom") {
    text = expression.value;
  } else if (expression.type === "not") {
    text = `!${serializeShortcutWhenExpression(expression.operand, precedence)}`;
  } else {
    const operator = expression.type === "and" ? "&&" : "||";
    text = `${serializeShortcutWhenExpression(expression.left, precedence)} ${operator} ${serializeShortcutWhenExpression(expression.right, precedence + 1)}`;
  }
  return precedence < parentPrecedence ? `(${text})` : text;
}

function shortcutWhenExpressionPrecedence(expression) {
  if (expression.type === "atom") return 4;
  if (expression.type === "not") return 3;
  if (expression.type === "and") return 2;
  return 1;
}

function evaluateShortcutWhenExpression(expression, atomValue) {
  if (expression.type === "atom") return Boolean(atomValue(expression.value));
  if (expression.type === "not") return !evaluateShortcutWhenExpression(expression.operand, atomValue);
  if (expression.type === "and") {
    return evaluateShortcutWhenExpression(expression.left, atomValue) &&
      evaluateShortcutWhenExpression(expression.right, atomValue);
  }
  return evaluateShortcutWhenExpression(expression.left, atomValue) ||
    evaluateShortcutWhenExpression(expression.right, atomValue);
}

function commandForShortcut(shortcut, when = "", exceptID = "") {
  for (const command of commandsById.values()) {
    if (command.id === exceptID) continue;
    const binding = effectiveShortcutBinding(command);
    if (binding.shortcut === shortcut && shortcutContextsOverlap(when, binding.when)) return command;
  }
  return null;
}

function customCommandForShortcut(shortcut, event) {
  for (const [commandID, binding] of customShortcuts.entries()) {
    if (binding.shortcut === shortcut && shortcutContextMatches(binding.when, event)) return commandsById.get(commandID) ?? null;
  }
  return null;
}

function runCustomShortcut(event) {
  if (paletteOpen || shortcutEditorOpen) return false;
  const shortcut = shortcutFromEvent(event);
  if (!shortcut) return false;
  const command = customCommandForShortcut(shortcut, event);
  if (!command) return false;
  event.preventDefault();
  if (isCommandDisabled(command)) {
    status.textContent = `${localizedCommandTitle(command)} is unavailable.`;
    return true;
  }
  executeCommand(command.id);
  return true;
}

function shortcutFromEvent(event) {
  const key = shortcutKeyName(event.key);
  if (!key) return "";
  const parts = [];
  if (event.metaKey || event.ctrlKey) parts.push("Cmd");
  if (event.shiftKey && key !== "Shift") parts.push("Shift");
  if (event.altKey && key !== "Alt") parts.push("Alt");
  if (key !== "Cmd" && key !== "Shift" && key !== "Alt") parts.push(key);
  return parts.join("+");
}

function shortcutKeyName(key) {
  if (key === "Meta" || key === "Control") return "Cmd";
  if (key === "Shift") return "Shift";
  if (key === "Alt" || key === "Option") return "Alt";
  if (key === " ") return "Space";
  if (key === "ArrowUp") return "Up";
  if (key === "ArrowDown") return "Down";
  if (key === "ArrowLeft") return "Left";
  if (key === "ArrowRight") return "Right";
  if (/^F\d{1,2}$/i.test(key)) return key.toUpperCase();
  if (key.length === 1) return key.toUpperCase();
  if (key === "Escape") return "";
  return key;
}

function openConditionBuilder() {
  if (!documentId) return;
  conditionBuilderOpen = true;
  renderConditionColumns();
  conditionColumn.value = String(sourceColFor(selected.col));
  if (!conditionColumn.value && conditionColumn.options.length > 0) {
    conditionColumn.selectedIndex = 0;
  }
  conditionOperator.value = conditionOperator.value || "contains";
  conditionTemplateSelect.value = "";
  document.body.dataset.conditionTemplate = "";
  updateConditionValueState();
  conditionPanel.hidden = false;
  conditionValue.focus();
  conditionValue.select();
}

function closeConditionBuilder() {
  conditionBuilderOpen = false;
  conditionPanel.hidden = true;
  canvas.focus();
}

function renderConditionColumns() {
  const previous = conditionColumn.value;
  populateConditionColumnSelect(conditionColumn, previous);
  renderConditionExtraRows();
}

function populateConditionColumnSelect(select, previous = "") {
  const fragment = document.createDocumentFragment();
  for (const column of columns) {
    const option = document.createElement("option");
    option.value = String(column.sourceCol);
    option.textContent = column.name;
    fragment.append(option);
  }
  select.replaceChildren(fragment);
  if ([...select.options].some((option) => option.value === previous)) {
    select.value = previous;
  }
}

function conditionOperatorOptionsFragment() {
  const fragment = document.createDocumentFragment();
  for (const option of conditionOperator.options) {
    fragment.append(option.cloneNode(true));
  }
  return fragment;
}

function renderConditionExtraRows() {
  conditionExtraRows = normalizeConditionExtraRowsForGroups(conditionExtraRows);
  conditionRows.replaceChildren(...conditionExtraRows.map((row, index) => conditionExtraRowElement(row, index)));
  updateConditionValueState();
}

function normalizeConditionExtraRowsForGroups(rows) {
  const normalized = rows.map((row) => ({
    ...row,
    groupStart: conditionGroupCount(row.groupStart),
    groupEnd: conditionGroupCount(row.groupEnd),
  }));
  rows.forEach((row, index) => {
    if (!row?.groupPrevious || row.groupStart != null || row.groupEnd != null) return;
    if (index === 0) {
      if (conditionPrimaryOpenGroup && conditionGroupCount(conditionPrimaryOpenGroup.value) === 0) {
        conditionPrimaryOpenGroup.value = "1";
      }
    } else {
      normalized[index - 1].groupStart = Math.max(1, normalized[index - 1].groupStart);
    }
    normalized[index].groupEnd = Math.max(1, normalized[index].groupEnd);
  });
  return normalized;
}

function conditionExtraRowElement(row, index) {
  const number = index + 2;
  const fieldset = document.createElement("fieldset");
  fieldset.className = "condition-row condition-row-extra";
  fieldset.dataset.conditionIndex = String(index);

  const legend = document.createElement("legend");
  legend.textContent = tf("condition.legend", { number });
  fieldset.append(legend);

  const joinLabel = document.createElement("label");
  joinLabel.className = "condition-join-control";
  joinLabel.append(conditionLabelText("condition.join"));
  const joinSelect = document.createElement("select");
  joinSelect.className = "condition-join-select";
  joinSelect.dataset.conditionField = "join";
  joinSelect.setAttribute("aria-label", tf("condition.joinAria", { number }));
  for (const value of ["AND", "OR"]) {
    const option = document.createElement("option");
    option.value = value;
    option.textContent = t(value === "OR" ? "condition.join.or" : "condition.join.and");
    joinSelect.append(option);
  }
  joinSelect.value = row.join === "OR" ? "OR" : "AND";
  joinLabel.append(joinSelect);
  fieldset.append(joinLabel);

  const openLabel = document.createElement("label");
  openLabel.className = "condition-group-open-control";
  openLabel.title = t("condition.group.openTitle");
  openLabel.append(conditionLabelText("("));
  const openSelect = conditionGroupCountSelect(
    "groupStart",
    Number(row.groupStart ?? 0),
    "(",
    tf("condition.group.openAria", { number }),
  );
  openLabel.append(openSelect);
  fieldset.append(openLabel);

  const closeLabel = document.createElement("label");
  closeLabel.className = "condition-group-close-control";
  closeLabel.title = t("condition.group.closeTitle");
  closeLabel.append(conditionLabelText(")"));
  const closeSelect = conditionGroupCountSelect(
    "groupEnd",
    Number(row.groupEnd ?? 0),
    ")",
    tf("condition.group.closeAria", { number }),
  );
  closeLabel.append(closeSelect);
  fieldset.append(closeLabel);

  const columnLabel = document.createElement("label");
  columnLabel.append(conditionLabelText("condition.column"));
  const columnSelect = document.createElement("select");
  columnSelect.className = "condition-column";
  columnSelect.dataset.conditionField = "column";
  columnSelect.setAttribute("aria-label", tf("condition.columnAria", { number }));
  populateConditionColumnSelect(columnSelect, row.column || conditionColumn.value || String(sourceColFor(selected.col)));
  columnLabel.append(columnSelect);
  fieldset.append(columnLabel);

  const operatorLabel = document.createElement("label");
  operatorLabel.append(conditionLabelText("condition.operator"));
  const operatorSelect = document.createElement("select");
  operatorSelect.className = "condition-operator";
  operatorSelect.dataset.conditionField = "operator";
  operatorSelect.setAttribute("aria-label", tf("condition.operatorAria", { number }));
  operatorSelect.append(conditionOperatorOptionsFragment());
  operatorSelect.value = row.operator || "contains";
  operatorLabel.append(operatorSelect);
  fieldset.append(operatorLabel);

  const valueLabel = document.createElement("label");
  valueLabel.append(conditionLabelText("condition.value"));
  const valueInput = document.createElement("input");
  valueInput.className = "condition-value";
  valueInput.dataset.conditionField = "value";
  valueInput.type = "text";
  valueInput.autocomplete = "off";
  valueInput.spellcheck = false;
  valueInput.setAttribute("aria-label", tf("condition.valueAria", { number }));
  valueInput.value = row.value || "";
  valueLabel.append(valueInput);
  fieldset.append(valueLabel);

  const actions = document.createElement("div");
  actions.className = "condition-row-actions";

  const moveUpButton = document.createElement("button");
  moveUpButton.type = "button";
  moveUpButton.dataset.conditionAction = "move-up";
  moveUpButton.textContent = t("condition.moveUp");
  moveUpButton.title = tf("condition.moveUpTitle", { number });

  const moveDownButton = document.createElement("button");
  moveDownButton.type = "button";
  moveDownButton.dataset.conditionAction = "move-down";
  moveDownButton.textContent = t("condition.moveDown");
  moveDownButton.title = tf("condition.moveDownTitle", { number });
  moveDownButton.disabled = index >= conditionExtraRows.length - 1;

  const duplicateButton = document.createElement("button");
  duplicateButton.type = "button";
  duplicateButton.dataset.conditionAction = "duplicate";
  duplicateButton.textContent = t("condition.duplicate");
  duplicateButton.title = tf("condition.duplicateTitle", { number });

  const groupButton = document.createElement("button");
  groupButton.type = "button";
  groupButton.dataset.conditionAction = "group-previous";
  groupButton.textContent = t("condition.groupPrevious");
  groupButton.title = tf("condition.groupPreviousTitle", { number });

  const ungroupButton = document.createElement("button");
  ungroupButton.type = "button";
  ungroupButton.dataset.conditionAction = "ungroup-previous";
  ungroupButton.textContent = t("condition.ungroupPrevious");
  ungroupButton.title = tf("condition.ungroupPreviousTitle", { number });

  const removeButton = document.createElement("button");
  removeButton.className = "condition-remove-button";
  removeButton.type = "button";
  removeButton.dataset.conditionAction = "remove";
  removeButton.textContent = t("condition.remove");
  removeButton.title = tf("condition.removeTitle", { number });
  actions.append(moveUpButton, moveDownButton, duplicateButton, groupButton, ungroupButton, removeButton);
  fieldset.append(actions);
  return fieldset;
}

function conditionLabelText(key) {
  const span = document.createElement("span");
  span.textContent = key.startsWith("condition.") ? t(key) : key;
  return span;
}

function conditionGroupCountSelect(field, value, symbol, label) {
  const select = document.createElement("select");
  select.className = `condition-group-${field === "groupStart" ? "open" : "close"}`;
  select.dataset.conditionField = field;
  select.setAttribute("aria-label", label);
  for (let count = 0; count <= 3; count++) {
    const option = document.createElement("option");
    option.value = String(count);
    option.textContent = count === 0 ? "0" : symbol.repeat(count);
    select.append(option);
  }
  select.value = String(conditionGroupCount(value));
  return select;
}

function conditionGroupCount(value) {
  const number = Number(value);
  if (!Number.isFinite(number)) return 0;
  return Math.max(0, Math.min(3, Math.trunc(number)));
}

function updateConditionValueState() {
  updateConditionValueInput(conditionOperator, conditionValue);
  for (const row of conditionRows.querySelectorAll(".condition-row-extra")) {
    updateConditionValueInput(
      row.querySelector(".condition-operator"),
      row.querySelector(".condition-value"),
    );
  }
  const joins = Array.from(conditionRows.querySelectorAll(".condition-join-select")).map((select) => select.value);
  const groupStarts = [conditionPrimaryOpenGroup, ...conditionRows.querySelectorAll(".condition-group-open")]
    .map((select) => conditionGroupCount(select?.value));
  const groupEnds = [conditionPrimaryCloseGroup, ...conditionRows.querySelectorAll(".condition-group-close")]
    .map((select) => conditionGroupCount(select?.value));
  const nestedCount = groupStarts.reduce((sum, count) => sum + count, 0);
  document.body.dataset.conditionGrouped = joins.join(",");
  document.body.dataset.conditionGroupOpen = groupStarts.join(",");
  document.body.dataset.conditionGroupClose = groupEnds.join(",");
  document.body.dataset.conditionNestedGroupCount = String(nestedCount);
  document.body.dataset.conditionRowCount = String(1 + conditionRows.querySelectorAll(".condition-row-extra").length);
}

function updateConditionValueInput(operatorInput, valueInput) {
  if (!operatorInput || !valueInput) return;
  const noValue = !conditionOperatorNeedsValue(operatorInput.value);
  valueInput.disabled = noValue;
  if (noValue) valueInput.value = "";
  valueInput.placeholder = noValue ? "" : conditionValuePlaceholder(operatorInput.value);
}

function conditionOperatorNeedsValue(operator) {
  return !["blank", "notBlank", "booleanTrue", "booleanFalse"].includes(operator);
}

function conditionValuePlaceholder(operator) {
  if (operator === "between" || operator === "notBetween") return t("condition.placeholder.range");
  if (operator === "inList" || operator === "notInList") return t("condition.placeholder.list");
  if (String(operator).startsWith("date")) return t("condition.placeholder.date");
  return "";
}

async function applyConditionBuilder() {
  if (!documentId) return;
  readConditionExtraRowsFromDOM();
  const builderRows = conditionBuilderRowStatesFromDOM();
  const builders = conditionBuilderRowsFromDOM();
  const nodes = [];
  const groupedTerms = [];
  for (const [index, builder] of builders.entries()) {
    const next = conditionExpression(builder.sourceCol, builder.operator, builder.value, builder.valueInput);
    if (!next) return;
    nodes.push({
      join: index === 0 ? "" : builder.join,
      expression: next,
    });
    groupedTerms.push({ ...builder, join: index === 0 ? "" : builder.join, expression: next });
  }
  const hasExplicitGroups = groupedTerms.some((term) => term.groupStart > 0 || term.groupEnd > 0);
  const expression = hasExplicitGroups
    ? conditionGroupedExpression(groupedTerms)
    : conditionNodesExpression(nodes);
  if (!expression) return;
  document.body.dataset.conditionPreview = expression;
  document.body.dataset.conditionNestedGroupCount = String(groupedTerms.reduce((sum, term) => sum + term.groupStart, 0));
  document.body.dataset.conditionRowCount = String(builders.length);
  document.body.dataset.conditionGroupingMode = hasExplicitGroups ? "explicit" : "implicit";
  commandInput.value = `where ${expression}`;
  closeConditionBuilder();
  const visibleCol = visibleColForSource(Number(conditionColumn.value));
  await applyFilterRequest(
    { documentId, mode: "sql", expression },
    t("condition.filterLabel"),
    visibleCol >= 0 ? visibleCol : selected.col,
    {
      matched: "condition.status.matched",
      failed: "condition.status.failed",
    },
    { builderRows },
  );
}

function addConditionRow() {
  readConditionExtraRowsFromDOM();
  conditionExtraRows.push({
    join: "AND",
    column: conditionColumn.value || String(sourceColFor(selected.col)),
    operator: "contains",
    value: "",
    groupStart: 0,
    groupEnd: 0,
  });
  renderConditionExtraRows();
  const valueInput = conditionRows.querySelector(".condition-row-extra:last-child .condition-value");
  valueInput?.focus();
}

function useConditionTemplate() {
  if (!documentId) return;
  const template = conditionTemplateSelect.value;
  if (!template) {
    status.textContent = t("condition.status.chooseQuick");
    return;
  }
  const draft = conditionTemplateDraft(template);
  if (!draft) return;
  conditionColumn.value = String(draft.column);
  conditionOperator.value = draft.operator;
  conditionValue.value = draft.value ?? "";
  updateConditionValueState();
  document.body.dataset.conditionTemplate = `${template}:${columns.find((column) => column.sourceCol === draft.column)?.name || draft.column}:${draft.operator}:${String(draft.value ?? "").slice(0, 80)}`;
  document.body.dataset.conditionPreview = conditionExpression(draft.column, draft.operator, draft.value ?? "", conditionValue);
  status.textContent = tf("condition.status.ready", { label: conditionTemplateLabel(template) });
  conditionValue.focus();
  if (conditionOperatorNeedsValue(draft.operator)) conditionValue.select();
}

function conditionTemplateDraft(template) {
  const activeColumn = sourceColFor(selected.col);
  if (!Number.isInteger(activeColumn) || activeColumn < 0) {
    status.textContent = t("condition.status.selectCell");
    return null;
  }
  const value = String(valueAt(selected.row, selected.col) ?? "");
  switch (template) {
    case "selectedEquals":
      return { column: activeColumn, operator: "equals", value };
    case "selectedContains":
      return { column: activeColumn, operator: "contains", value };
    case "selectedNotContains":
      return { column: activeColumn, operator: "notContains", value };
    case "selectedNotEquals":
      return { column: activeColumn, operator: "notEquals", value };
    case "activeBlank":
      return { column: activeColumn, operator: "blank", value: "" };
    case "activeNotBlank":
      return { column: activeColumn, operator: "notBlank", value: "" };
    case "activeTrue":
      return { column: activeColumn, operator: "booleanTrue", value: "" };
    case "activeFalse":
      return { column: activeColumn, operator: "booleanFalse", value: "" };
    default:
      status.textContent = t("condition.status.unsupportedQuick");
      return null;
  }
}

function conditionTemplateLabel(template) {
  const option = Array.from(conditionTemplateSelect.options).find((item) => item.value === template);
  return option?.textContent || template;
}

function removeConditionRow(index) {
  if (!Number.isInteger(index) || index < 0) return;
  readConditionExtraRowsFromDOM();
  conditionExtraRows.splice(index, 1);
  renderConditionExtraRows();
  conditionValue.focus();
}

function duplicateConditionRow(index) {
  if (!Number.isInteger(index) || index < 0) return;
  readConditionExtraRowsFromDOM();
  const row = conditionExtraRows[index];
  if (!row) return;
  const duplicate = {
    ...row,
    groupStart: 0,
    groupEnd: conditionGroupCount(row.groupEnd),
  };
  conditionExtraRows[index] = {
    ...row,
    groupEnd: 0,
  };
  conditionExtraRows.splice(index + 1, 0, duplicate);
  renderConditionExtraRows();
  const fromNumber = index + 2;
  const toNumber = fromNumber + 1;
  document.body.dataset.conditionDuplicate = `${fromNumber}->${toNumber}`;
  status.textContent = tf("condition.status.duplicated", { number: formatCount(fromNumber) });
  const valueInput = conditionRows.querySelector(`.condition-row-extra[data-condition-index="${index + 1}"] .condition-value`);
  valueInput?.focus();
}

function groupConditionWithPrevious(extraIndex, direction) {
  if (!Number.isInteger(extraIndex) || ![-1, 1].includes(direction)) return;
  const current = extraIndex + 1;
  const previous = current - 1;
  const rows = conditionBuilderRowStatesFromDOM();
  if (current < 1 || current >= rows.length || previous < 0) return;
  const previousStart = conditionGroupCount(rows[previous].groupStart);
  const currentEnd = conditionGroupCount(rows[current].groupEnd);
  const grouping = direction > 0;
  if (grouping) {
    if (previousStart >= 3 || currentEnd >= 3) return;
    rows[previous].groupStart = previousStart + 1;
    rows[current].groupEnd = currentEnd + 1;
  } else {
    if (previousStart <= 0 || currentEnd <= 0) return;
    rows[previous].groupStart = previousStart - 1;
    rows[current].groupEnd = currentEnd - 1;
  }
  setConditionBuilderRows(rows);
  const fromNumber = previous + 1;
  const toNumber = current + 1;
  document.body.dataset.conditionGroupAction = `${grouping ? "group" : "ungroup"}:${fromNumber}-${toNumber}`;
  status.textContent = tf(
    grouping ? "condition.status.groupedPrevious" : "condition.status.ungroupedPrevious",
    {
      from: formatCount(fromNumber),
      to: formatCount(toNumber),
    },
  );
  const field = conditionRows.querySelector(`.condition-row-extra[data-condition-index="${extraIndex}"] .condition-value`);
  field?.focus();
}

function clearConditionGroups() {
  const rows = conditionBuilderRowStatesFromDOM();
  const cleared = rows.reduce(
    (count, row) => count + conditionGroupCount(row.groupStart) + conditionGroupCount(row.groupEnd),
    0,
  );
  const flattened = rows.map((row) => ({ ...row, groupStart: 0, groupEnd: 0 }));
  setConditionBuilderRows(flattened);
  document.body.dataset.conditionGroupAction = "clear";
  document.body.dataset.conditionGroupClear = `${cleared}->0`;
  document.body.dataset.conditionGroupingMode = "implicit";
  status.textContent = tf("condition.status.clearedGroups", { count: formatCount(cleared) });
  conditionValue.focus();
}

function groupAllConditionRows() {
  const rows = conditionBuilderRowStatesFromDOM();
  if (rows.length < 2) {
    status.textContent = t("condition.status.addConditionFirst");
    return;
  }
  const firstStart = conditionGroupCount(rows[0].groupStart);
  const lastIndex = rows.length - 1;
  const lastEnd = conditionGroupCount(rows[lastIndex].groupEnd);
  if (firstStart >= 3 || lastEnd >= 3) {
    status.textContent = t("condition.status.maxGroups");
    return;
  }
  rows[0].groupStart = firstStart + 1;
  rows[lastIndex].groupEnd = lastEnd + 1;
  setConditionBuilderRows(rows);
  document.body.dataset.conditionGroupAction = `group-all:${rows.length}`;
  document.body.dataset.conditionGroupAll = `${rows.length}:${rows[0].groupStart}:${rows[lastIndex].groupEnd}`;
  document.body.dataset.conditionGroupingMode = "explicit";
  status.textContent = tf("condition.status.groupedAll", { count: formatCount(rows.length) });
  conditionValue.focus();
}

function setAllConditionJoins(join) {
  const normalizedJoin = join === "OR" ? "OR" : "AND";
  const rows = conditionBuilderRowStatesFromDOM();
  if (rows.length < 2) {
    status.textContent = t("condition.status.addConditionFirst");
    return;
  }
  for (let index = 1; index < rows.length; index++) {
    rows[index].join = normalizedJoin;
  }
  setConditionBuilderRows(rows);
  document.body.dataset.conditionJoinAction = `all-${normalizedJoin.toLowerCase()}:${rows.length - 1}`;
  document.body.dataset.conditionGrouped = rows.slice(1).map((row) => row.join).join(",");
  status.textContent = t(normalizedJoin === "OR" ? "condition.status.allOr" : "condition.status.allAnd");
  conditionValue.focus();
}

function moveConditionRow(extraIndex, direction) {
  if (!Number.isInteger(extraIndex) || ![-1, 1].includes(direction)) return;
  const from = extraIndex + 1;
  const rows = conditionBuilderRowStatesFromDOM();
  const to = from + direction;
  if (from < 1 || from >= rows.length || to < 0 || to >= rows.length) return;
  const connector = rows[from].join || rows[to].join || "AND";
  [rows[from], rows[to]] = [rows[to], rows[from]];
  rows[0].join = "";
  for (let index = 1; index < rows.length; index++) {
    rows[index].join = rows[index].join || "AND";
  }
  if (to === 0 && rows[1]) {
    rows[1].join = connector;
  }
  setConditionBuilderRows(rows);
  document.body.dataset.conditionReorder = `${from + 1}->${to + 1}`;
  status.textContent = tf("condition.status.reordered", {
    from: formatCount(from + 1),
    to: formatCount(to + 1),
  });
  const movedField = to === 0
    ? conditionValue
    : conditionRows.querySelector(`.condition-row-extra[data-condition-index="${to - 1}"] .condition-value`);
  movedField?.focus();
}

function readConditionExtraRowsFromDOM() {
  conditionExtraRows = Array.from(conditionRows.querySelectorAll(".condition-row-extra")).map((row) => ({
    join: row.querySelector(".condition-join-select")?.value === "OR" ? "OR" : "AND",
    column: row.querySelector(".condition-column")?.value || "",
    operator: row.querySelector(".condition-operator")?.value || "contains",
    value: row.querySelector(".condition-value")?.value || "",
    groupStart: conditionGroupCount(row.querySelector(".condition-group-open")?.value),
    groupEnd: conditionGroupCount(row.querySelector(".condition-group-close")?.value),
  }));
}

function conditionBuilderRowStatesFromDOM() {
  readConditionExtraRowsFromDOM();
  return [
    {
      join: "",
      column: conditionColumn.value,
      operator: conditionOperator.value,
      value: conditionValue.value,
      groupStart: conditionGroupCount(conditionPrimaryOpenGroup?.value),
      groupEnd: conditionGroupCount(conditionPrimaryCloseGroup?.value),
    },
    ...conditionExtraRows.map((row) => ({ ...row })),
  ];
}

function setConditionBuilderRows(rows) {
  const [primary, ...extras] = rows;
  if (!primary) return;
  setSelectExistingValue(conditionColumn, primary.column, String(sourceColFor(selected.col)));
  setSelectExistingValue(conditionOperator, primary.operator, "contains");
  conditionValue.value = primary.value || "";
  if (conditionPrimaryOpenGroup) conditionPrimaryOpenGroup.value = String(conditionGroupCount(primary.groupStart));
  if (conditionPrimaryCloseGroup) conditionPrimaryCloseGroup.value = String(conditionGroupCount(primary.groupEnd));
  conditionExtraRows = extras.map((row) => ({
    join: row.join === "OR" ? "OR" : "AND",
    column: row.column || "",
    operator: row.operator || "contains",
    value: row.value || "",
    groupStart: conditionGroupCount(row.groupStart),
    groupEnd: conditionGroupCount(row.groupEnd),
  }));
  renderConditionExtraRows();
}

function normalizeConditionBuilderRows(rows) {
  if (!Array.isArray(rows)) return [];
  return rows
    .map((row, index) => {
      if (!row || typeof row !== "object") return null;
      const column = String(row.column ?? "");
      const operator = String(row.operator ?? "contains");
      if (!column || !conditionOperatorLabelKeys()[operator]) return null;
      return {
        join: index === 0 ? "" : (row.join === "OR" ? "OR" : "AND"),
        column,
        operator,
        value: String(row.value ?? ""),
        groupStart: conditionGroupCount(row.groupStart),
        groupEnd: conditionGroupCount(row.groupEnd),
      };
    })
    .filter(Boolean);
}

function setSelectExistingValue(select, value, fallback = "") {
  if (!select) return;
  const raw = String(value ?? "");
  if ([...select.options].some((option) => option.value === raw)) {
    select.value = raw;
    return;
  }
  const fallbackValue = String(fallback ?? "");
  if (fallbackValue && [...select.options].some((option) => option.value === fallbackValue)) {
    select.value = fallbackValue;
    return;
  }
  if (select.options.length > 0) select.selectedIndex = 0;
}

function conditionBuilderRowsFromDOM() {
  return [
    {
      join: "",
      sourceCol: Number(conditionColumn.value),
      operator: conditionOperator.value,
      value: conditionValue.value,
      valueInput: conditionValue,
      groupStart: conditionGroupCount(conditionPrimaryOpenGroup?.value),
      groupEnd: conditionGroupCount(conditionPrimaryCloseGroup?.value),
    },
    ...Array.from(conditionRows.querySelectorAll(".condition-row-extra")).map((row) => conditionExtraBuilderFromElement(row)),
  ];
}

function conditionExtraBuilderFromElement(row) {
  return {
    join: row.querySelector(".condition-join-select")?.value === "OR" ? "OR" : "AND",
    sourceCol: Number(row.querySelector(".condition-column")?.value),
    operator: row.querySelector(".condition-operator")?.value,
    value: row.querySelector(".condition-value")?.value ?? "",
    valueInput: row.querySelector(".condition-value"),
    groupStart: conditionGroupCount(row.querySelector(".condition-group-open")?.value),
    groupEnd: conditionGroupCount(row.querySelector(".condition-group-close")?.value),
  };
}

function conditionGroupedExpression(terms) {
  let depth = 0;
  const parts = [];
  for (const [index, term] of terms.entries()) {
    const open = conditionGroupCount(term.groupStart);
    const close = conditionGroupCount(term.groupEnd);
    depth += open;
    if (close > depth) {
      status.textContent = tf("condition.status.tooManyCloses", { number: index + 1 });
      term.valueInput?.focus();
      return "";
    }
    const condition = `${"(".repeat(open)}${term.expression}${")".repeat(close)}`;
    parts.push(index === 0 ? condition : `${term.join} ${condition}`);
    depth -= close;
  }
  if (depth !== 0) {
    status.textContent = t("condition.status.closeGroups");
    const lastInput = terms[terms.length - 1]?.valueInput;
    lastInput?.focus();
    return "";
  }
  return parts.join(" ");
}

function conditionNodesExpression(nodes) {
  if (!nodes.length) return "";
  let expression = nodes[0].expression;
  for (const node of nodes.slice(1)) {
    expression = `(${expression}) ${node.join} (${node.expression})`;
  }
  return expression;
}

function conditionExpression(sourceCol, operator, value, valueInput = conditionValue) {
  if (!Number.isInteger(sourceCol) || sourceCol < 0) {
    status.textContent = t("condition.status.chooseColumn");
    return "";
  }
  const identifier = sqlColumnIdentifier(sourceCol);
  const needsValue = conditionOperatorNeedsValue(operator);
  if (needsValue && value === "") {
    status.textContent = t("condition.status.enterValue");
    valueInput.focus();
    return "";
  }
  switch (operator) {
    case "contains":
      return `${identifier} LIKE ${sqlStringLiteral(`%${value}%`)}`;
    case "notContains":
      return `${identifier} NOT LIKE ${sqlStringLiteral(`%${value}%`)}`;
    case "equals":
      return `${identifier} = ${sqlStringLiteral(value)}`;
    case "notEquals":
      return `${identifier} <> ${sqlStringLiteral(value)}`;
    case "inList":
      return listConditionExpression(identifier, value, valueInput, false);
    case "notInList":
      return listConditionExpression(identifier, value, valueInput, true);
    case "startsWith":
      return `${identifier} LIKE ${sqlStringLiteral(`${value}%`)}`;
    case "endsWith":
      return `${identifier} LIKE ${sqlStringLiteral(`%${value}`)}`;
    case "blank":
      return `${identifier} = ''`;
    case "notBlank":
      return `${identifier} <> ''`;
    case "booleanTrue":
      return booleanConditionExpression(identifier, true);
    case "booleanFalse":
      return booleanConditionExpression(identifier, false);
    case "dateEquals":
      return dateConditionExpression(identifier, "=", value, valueInput);
    case "dateBefore":
      return dateConditionExpression(identifier, "<", value, valueInput);
    case "dateAfter":
      return dateConditionExpression(identifier, ">", value, valueInput);
    case "dateOnOrBefore":
      return dateConditionExpression(identifier, "<=", value, valueInput);
    case "dateOnOrAfter":
      return dateConditionExpression(identifier, ">=", value, valueInput);
    case "gt":
      return `${identifier} > ${sqlComparisonLiteral(value)}`;
    case "gte":
      return `${identifier} >= ${sqlComparisonLiteral(value)}`;
    case "lt":
      return `${identifier} < ${sqlComparisonLiteral(value)}`;
    case "lte":
      return `${identifier} <= ${sqlComparisonLiteral(value)}`;
    case "between":
      return rangeConditionExpression(identifier, value, valueInput, false);
    case "notBetween":
      return rangeConditionExpression(identifier, value, valueInput, true);
    default:
      status.textContent = t("condition.status.chooseOperator");
      return "";
  }
}

function rangeConditionExpression(identifier, value, valueInput, negate) {
  const range = conditionRangeParts(value);
  if (!range) {
    status.textContent = t("condition.status.enterRange");
    valueInput.focus();
    return "";
  }
  const expression = `(${identifier} >= ${sqlComparisonLiteral(range[0])} AND ${identifier} <= ${sqlComparisonLiteral(range[1])})`;
  return negate ? `NOT ${expression}` : expression;
}

function listConditionExpression(identifier, value, valueInput, negate) {
  const values = conditionListValues(value);
  if (!values.length) {
    status.textContent = t("condition.status.enterList");
    valueInput.focus();
    return "";
  }
  const operator = negate ? "NOT IN" : "IN";
  return `${identifier} ${operator} (${values.map(sqlStringLiteral).join(", ")})`;
}

function conditionListValues(value) {
  const values = [];
  let current = "";
  let quoted = false;
  const raw = String(value ?? "");
  for (let index = 0; index < raw.length; index++) {
    const char = raw[index];
    if (quoted) {
      if (char === '"' && raw[index + 1] === '"') {
        current += '"';
        index += 1;
      } else if (char === '"') {
        quoted = false;
      } else {
        current += char;
      }
      continue;
    }
    if (char === '"') {
      quoted = true;
    } else if (char === "," || char === "\n") {
      const next = current.trim();
      if (next) values.push(next);
      current = "";
    } else {
      current += char;
    }
  }
  const last = current.trim();
  if (last) values.push(last);
  return values;
}

function conditionRangeParts(value) {
  const raw = String(value ?? "").trim();
  const match = raw.match(/^(.+?)(?:\s*\.\.\s*|\s+to\s+|\s*,\s*)(.+)$/i);
  if (!match) return null;
  const lower = match[1].trim();
  const upper = match[2].trim();
  return lower && upper ? [lower, upper] : null;
}

function booleanConditionExpression(identifier, expected) {
  const normalized = `LOWER(TRIM(${identifier}))`;
  const trimmed = `TRIM(${identifier})`;
  if (expected) {
    return `(${normalized} = 'true' OR ${trimmed} = '1')`;
  }
  return `(${normalized} = 'false' OR ${trimmed} = '0' OR ${trimmed} = '')`;
}

function dateConditionExpression(identifier, operator, value, valueInput) {
  const date = normalizedDateLiteral(value);
  if (!date) {
    status.textContent = t("condition.status.enterDate");
    valueInput.focus();
    return "";
  }
  return `CAST(${identifier} AS date) ${operator} DATE '${date}'`;
}

function normalizedDateLiteral(value) {
  const raw = String(value ?? "").trim();
  const match = raw.match(/^(\d{4})-(\d{2})-(\d{2})$/);
  if (!match) return "";
  const year = Number(match[1]);
  const month = Number(match[2]);
  const day = Number(match[3]);
  const date = new Date(Date.UTC(year, month - 1, day));
  if (
    date.getUTCFullYear() !== year ||
    date.getUTCMonth() + 1 !== month ||
    date.getUTCDate() !== day
  ) {
    return "";
  }
  return raw;
}

function sqlColumnIdentifier(sourceCol) {
  const name = columnNameForSource(sourceCol);
  if (!metadata?.dialect?.hasHeader || !name || /^Column \d+$/.test(name)) {
    return `c${sourceCol + 1}`;
  }
  return `"${String(name).replaceAll('"', '""')}"`;
}

function sqlStringLiteral(value) {
  return `'${String(value).replaceAll("'", "''")}'`;
}

function sqlComparisonLiteral(value) {
  const trimmed = String(value).trim();
  return /^-?\d+(?:\.\d+)?$/.test(trimmed) ? trimmed : sqlStringLiteral(value);
}

function renderCommandPalette() {
  const query = paletteInput.value.trim().toLowerCase();
  paletteItems = scoredPaletteCommands(query);
  if (paletteIndex >= paletteItems.length) paletteIndex = Math.max(0, paletteItems.length - 1);
  if (!paletteItems.length) {
    paletteList.replaceChildren(emptyPaletteNode("No matching commands"));
    return;
  }
  const fragment = document.createDocumentFragment();
  paletteItems.forEach((command, index) => {
    const button = document.createElement("button");
    button.type = "button";
    button.className = index === paletteIndex ? "palette-item active" : "palette-item";
    button.disabled = isCommandDisabled(command);
    button.setAttribute("role", "option");
    button.setAttribute("aria-selected", String(index === paletteIndex));

    const title = document.createElement("strong");
    title.textContent = localizedCommandTitle(command);
    const description = document.createElement("span");
    description.textContent = localizedCommandDescription(command);
    const shortcut = document.createElement("kbd");
    shortcut.textContent = effectiveShortcut(command);
    button.append(title, description, shortcut);
    button.addEventListener("mouseenter", () => {
      paletteIndex = index;
      renderCommandPalette();
    });
    button.addEventListener("click", () => executePaletteCommand(command.id));
    fragment.append(button);
  });
  paletteList.replaceChildren(fragment);
}

function scoredPaletteCommands(query) {
  const history = commandHistory();
  const recentRank = new Map(history.map((id, index) => [id, history.length - index]));
  return Array.from(commandsById.values())
    .filter((command) => command.implemented !== false)
    .map((command) => ({
      command,
      score: query ? commandFuzzyScore(command, query) : recentRank.get(command.id) ?? 0,
    }))
    .filter((entry) => !query || entry.score > 0)
    .sort((left, right) =>
      right.score - left.score ||
      commandTitle(left.command).localeCompare(commandTitle(right.command)),
    )
    .map((entry) => entry.command);
}

function commandFuzzyScore(command, query) {
  const values = [
    command.id,
    command.title,
    command.description,
    localizedCommandTitle(command),
    localizedCommandDescription(command),
    ...(command.aliases ?? []),
  ].filter(Boolean);
  let best = 0;
  for (const value of values) {
    const score = fuzzyScore(String(value).toLowerCase(), query);
    if (score > best) best = score;
  }
  return best;
}

function fuzzyScore(value, query) {
  if (!query) return 1;
  if (value === query) return 10000;
  if (value.startsWith(query)) return 8000 - value.length;
  const index = value.indexOf(query);
  if (index >= 0) return 6000 - index * 8 - value.length;

  let queryIndex = 0;
  let gap = 0;
  let previousMatch = -1;
  for (let valueIndex = 0; valueIndex < value.length && queryIndex < query.length; valueIndex++) {
    if (value[valueIndex] !== query[queryIndex]) continue;
    if (previousMatch >= 0) gap += valueIndex - previousMatch - 1;
    previousMatch = valueIndex;
    queryIndex++;
  }
  if (queryIndex !== query.length) return 0;
  return Math.max(1, 3000 - gap * 16 - value.length);
}

function commandTitle(command) {
  return localizedCommandTitle(command);
}

function commandHistory() {
  try {
    const parsed = JSON.parse(localStorage.getItem(commandHistoryStorageKey) || "[]");
    if (!Array.isArray(parsed)) return [];
    return parsed
      .filter((id) => typeof id === "string" && commandsById.has(id))
      .slice(0, commandHistoryLimit);
  } catch {
    return [];
  }
}

function rememberCommand(commandID) {
  if (commandID === "command.palette" || !commandsById.has(commandID)) return;
  const next = [
    commandID,
    ...commandHistory().filter((id) => id !== commandID),
  ].slice(0, commandHistoryLimit);
  try {
    localStorage.setItem(commandHistoryStorageKey, JSON.stringify(next));
  } catch {
    status.textContent = "Command history could not be saved.";
  }
}

function emptyPaletteNode(message) {
  const node = document.createElement("div");
  node.className = "profile-empty";
  node.textContent = message;
  return node;
}

async function executePaletteCommand(commandID) {
  const command = commandsById.get(commandID);
  if (!command || isCommandDisabled(command)) {
    status.textContent = command ? `${localizedCommandTitle(command)} is unavailable.` : "Command unavailable.";
    return;
  }
  closeCommandPalette();
  await executeCommand(commandID);
}

async function executeCommand(commandID) {
  rememberCommand(commandID);
  const sqlResultCommand = sqlResultCommandFormats[commandID];
  if (sqlResultCommand?.action === "copy") {
    await copySQLResult(sqlResultCommand.format);
    return;
  }
  if (sqlResultCommand?.action === "copySelection") {
    await copySQLSelection(sqlResultCommand.format);
    return;
  }
  if (sqlResultCommand?.action === "copyFullResult") {
    await copySQLFullResult(sqlResultCommand.format);
    return;
  }
  if (sqlResultCommand?.action === "export") {
    await exportSQLResult(sqlResultCommand.format);
    return;
  }
  if (copyRowFormatCommands.has(commandID)) {
    await copyActiveRow({ format: copyRowFormatCommands.get(commandID) });
    return;
  }
  if (copyColumnFormatCommands.has(commandID)) {
    await copyActiveColumn({ format: copyColumnFormatCommands.get(commandID) });
    return;
  }
  switch (commandID) {
    case "command.palette":
      openCommandPalette();
      break;
    case "settings.open":
      openSettingsPanel();
      break;
    case "settings.keyboardShortcuts":
      openKeyboardShortcuts();
      break;
    case "settings.exportShortcuts":
      exportShortcuts();
      break;
    case "settings.importShortcuts":
      showShortcutImportPanel();
      break;
    case "settings.checkForUpdates":
      await checkForUpdates("manual");
      break;
    case "help.aboutSupport":
      openAboutSupportPanel();
      break;
    case "workbench.toggleTheme":
      cycleTheme();
      break;
    case "workbench.splitEditor":
      splitEditorGroup();
      break;
    case "workbench.closeSplitEditor":
      closeSplitEditorGroup();
      break;
    case "workbench.focusNextGroup":
      await focusNextEditorGroup();
      break;
    case "workbench.moveTabToOtherGroup":
      moveCurrentTabToOtherGroup();
      break;
    case "workbench.toggleSplitOrientation":
      toggleSplitOrientation();
      break;
    case "file.new":
      await createNewFile();
      break;
    case "file.open":
      await openPathInputOrPicker();
      break;
    case "file.closeTab":
      await closeActiveOpenFileTab();
      break;
    case "file.save":
      await saveCurrentDocument();
      break;
    case "file.saveAs":
      openSaveAsDialog();
      break;
    case "file.restoreBackup":
      await openRestoreBackupPanel();
      break;
    case "file.exportView":
      await exportCurrentView();
      break;
    case "file.exportXlsx":
      await exportCurrentView("xlsx");
      break;
    case "file.toggleReadOnly":
      await toggleReadOnly();
      break;
    case "format.savePreset":
      saveCurrentFormatPreset();
      break;
    case "format.applyPreset":
      await applySelectedFormatPreset();
      break;
    case "format.switchWorksheet":
      focusWorkbookSheetSelector();
      break;
    case "format.saveRule":
      saveCurrentFormatRule();
      break;
    case "format.exportPresets":
      exportFormatPresets();
      break;
    case "format.importPresets":
      showFormatPresetImportPanel();
      break;
    case "format.applyDelimiter":
      await applyDelimiter();
      break;
    case "format.applyEncoding":
      await applyEncoding();
      break;
    case "format.applyLineEnding":
      await applyLineEnding();
      break;
    case "format.applyQuotePolicy":
      await applyQuotePolicy();
      break;
    case "format.exportQuoteRules":
      exportQuoteRules();
      break;
    case "format.importQuoteRules":
      showQuoteRulesImportPanel();
      break;
    case "format.toggleHeaderQuote":
      await toggleQuoteRule("header");
      break;
    case "format.toggleBodyQuote":
      await toggleQuoteRule("body");
      break;
    case "format.toggleColumnQuote":
      await toggleQuoteRule("column");
      break;
    case "format.toggleCellQuote":
      await toggleCellQuote();
      break;
    case "format.quoteSelection":
      await setSelectedCellQuoteOverrides("quote");
      break;
    case "format.noQuoteSelection":
      await setSelectedCellQuoteOverrides("no-quote");
      break;
    case "format.clearSelectionQuote":
      await setSelectedCellQuoteOverrides("clear");
      break;
    case "edit.undo":
      await undoCurrentDocument();
      break;
    case "edit.redo":
      await redoCurrentDocument();
      break;
    case "edit.copy":
      await copySelection("tsv");
      break;
    case "edit.copySelectedFormat":
      await copyUsingSelectedFormat();
      break;
    case "edit.copyWithHeaders":
      await copySelection("tsv", { includeHeaders: true });
      break;
    case "edit.copyRow":
      await copyActiveRow({ includeHeaders: false });
      break;
    case "edit.copyRowWithHeaders":
      await copyActiveRow({ includeHeaders: true });
      break;
    case "edit.copyColumn":
      await copyActiveColumn({ includeHeaders: false });
      break;
    case "edit.copyColumnWithHeader":
      await copyActiveColumn({ includeHeaders: true });
      break;
    case "selection.toggleHeaders":
      toggleHeaderInSelection();
      break;
    case "selection.addRange":
      addSmartSelectionRange();
      break;
    case "selection.subtractRange":
      subtractSmartSelectionRange();
      break;
    case "selection.invert":
      invertSmartSelection();
      break;
    case "selection.clearRanges":
      clearSmartSelectionRanges();
      break;
    case "selection.cycleAggregate":
      cycleAggregateMetric();
      break;
    case "selection.copyAggregate":
      await copyActiveAggregate();
      break;
    case "selection.insertAggregateRow":
      insertAggregateRow();
      break;
    case "selection.insertSumRow":
    case "selection.insertAverageRow":
    case "selection.insertMinRow":
    case "selection.insertMaxRow":
    case "selection.insertCountRow":
      insertAggregateRow(aggregateRowCommandMetrics.get(commandID));
      break;
    case "edit.copyCsv":
      await copySelection("csv");
      break;
    case "edit.copyMarkdown":
      await copySelection("markdown");
      break;
    case "edit.copyJson":
      await copySelection("json");
      break;
    case "edit.copyHtml":
      await copySelection("html");
      break;
    case "edit.copySql":
      await copySelection("sql");
      break;
    case "edit.copyLatex":
      await copySelection("latex");
      break;
    case "edit.copyRaw":
      await copySelection("raw");
      break;
    case "edit.paste":
      await pasteClipboard();
      break;
    case "edit.find":
      await focusFindInput({ run: true });
      break;
    case "edit.replacePreview":
      await previewReplace();
      break;
    case "edit.replace":
      await replaceCurrentQuery(false);
      break;
    case "edit.replaceAll":
      await replaceCurrentQuery(true);
      break;
    case "row.insert":
      await insertSelectedRow();
      break;
    case "row.delete":
      await deleteSelectedRow();
      break;
    case "row.cut":
      cutSelectedRow();
      break;
    case "row.paste":
      await pasteRowFromClipboard();
      break;
    case "row.pasteBelow":
      await pasteRowFromClipboard({ placement: "after-selection" });
      break;
    case "row.hide":
      await hideSelectedRow();
      break;
    case "row.showAll":
      await showAllRows();
      break;
    case "row.moveUp":
      await moveSelectedRow(-1);
      break;
    case "row.moveDown":
      await moveSelectedRow(1);
      break;
    case "column.insert":
      await insertSelectedColumn();
      break;
    case "column.delete":
      await deleteSelectedColumn();
      break;
    case "column.cut":
      cutSelectedColumn();
      break;
    case "column.paste":
      await pasteColumnFromClipboard();
      break;
    case "column.pasteRight":
      await pasteColumnFromClipboard({ placement: "after-selection" });
      break;
    case "column.hide":
      hideSelectedColumn();
      break;
    case "column.showAll":
      showAllHiddenColumns();
      break;
    case "column.moveLeft":
      await moveSelectedColumn(-1);
      break;
    case "column.moveRight":
      await moveSelectedColumn(1);
      break;
    case "table.transpose":
      await transposeCurrentDocument();
      break;
    case "table.toggleHeaderRow":
      await toggleHeaderRow();
      break;
    case "table.dedupeRows":
      await dedupeRows();
      break;
    case "table.dedupeColumns":
      await dedupeColumns();
      break;
    case "text.trim":
      await transformSelectedText("trim", localizedCommandLabel(commandID, "Trim Text"));
      break;
    case "text.uppercase":
      await transformSelectedText("upper", localizedCommandLabel(commandID, "Uppercase Text"));
      break;
    case "text.lowercase":
      await transformSelectedText("lower", localizedCommandLabel(commandID, "Lowercase Text"));
      break;
    case "text.snakeCase":
      await transformSelectedText("snake", localizedCommandLabel(commandID, "Snake Case"));
      break;
    case "text.kebabCase":
      await transformSelectedText("kebab", localizedCommandLabel(commandID, "Kebab Case"));
      break;
    case "text.camelCase":
      await transformSelectedText("camel", localizedCommandLabel(commandID, "Camel Case"));
      break;
    case "text.pascalCase":
      await transformSelectedText("pascal", localizedCommandLabel(commandID, "Pascal Case"));
      break;
    case "text.titleCase":
      await transformSelectedText("title", localizedCommandLabel(commandID, "Title Case"));
      break;
    case "text.fullWidth":
      await transformSelectedText("wide", localizedCommandLabel(commandID, "Full Width"));
      break;
    case "text.halfWidth":
      await transformSelectedText("narrow", localizedCommandLabel(commandID, "Half Width"));
      break;
    case "text.normalizeNFC":
      await transformSelectedText("nfc", localizedCommandLabel(commandID, "Normalize NFC"));
      break;
    case "text.normalizeNFD":
      await transformSelectedText("nfd", localizedCommandLabel(commandID, "Normalize NFD"));
      break;
    case "text.normalizeNFKC":
      await transformSelectedText("nfkc", localizedCommandLabel(commandID, "Normalize NFKC"));
      break;
    case "text.normalizeNFKD":
      await transformSelectedText("nfkd", localizedCommandLabel(commandID, "Normalize NFKD"));
      break;
    case "text.column.trim":
    case "text.column.uppercase":
    case "text.column.lowercase":
    case "text.column.snakeCase":
    case "text.column.kebabCase":
    case "text.column.camelCase":
    case "text.column.pascalCase":
    case "text.column.titleCase":
    case "text.column.fullWidth":
    case "text.column.halfWidth":
    case "text.column.normalizeNFC":
    case "text.column.normalizeNFD":
    case "text.column.normalizeNFKC":
    case "text.column.normalizeNFKD":
      await transformColumnTextForCommand(commandID);
      break;
    case "cell.swap":
      await swapSelectedCells();
      break;
    case "cell.swapUp":
      await swapCellWithNeighbor("up");
      break;
    case "cell.swapDown":
      await swapCellWithNeighbor("down");
      break;
    case "cell.swapLeft":
      await swapCellWithNeighbor("left");
      break;
    case "cell.swapRight":
      await swapCellWithNeighbor("right");
      break;
    case "cell.openLink":
      openSelectedCellLink();
      break;
    case "cell.joinHorizontal":
      await joinSelectedCells("horizontal");
      break;
    case "cell.splitHorizontal":
      await splitActiveCell("horizontal");
      break;
    case "cell.joinVertical":
      await joinSelectedCells("vertical");
      break;
    case "cell.splitVertical":
      await splitActiveCell("vertical");
      break;
    case "cell.fillDown":
      autoFillSelection("down");
      break;
    case "cell.fillRight":
      autoFillSelection("right");
      break;
    case "cell.fillUp":
      autoFillSelection("up");
      break;
    case "cell.fillLeft":
      autoFillSelection("left");
      break;
    case "view.filter":
      await filterCurrentQuery();
      break;
    case "view.filterBuilder":
      openConditionBuilder();
      break;
    case "view.saveFilter":
      saveActiveFilter();
      break;
    case "view.saveSort":
      saveActiveSort();
      break;
    case "view.removeLastSortKey":
      await removeLastSortKey();
      break;
    case "view.exportSavedSets":
      exportSavedViews();
      break;
    case "view.importSavedSets":
      showSavedViewImportPanel();
      break;
    case "view.sort":
      await sortSelectedColumn();
      break;
    case "column.sortAscending":
    case "column.sortDescending":
    case "column.sortLengthAscending":
    case "column.sortLengthDescending":
    case "column.addSortAscending":
    case "column.addSortDescending":
      await sortActiveColumnForCommand(commandID);
      break;
    case "column.filterBySelectedValue":
    case "column.filterBlanks":
    case "column.filterNonBlanks":
      await filterActiveColumnForCommand(commandID);
      break;
    case "view.clear":
      await clearView();
      break;
    case "view.sql":
      setViewMode("sql");
      break;
    case "sql.openResult":
      await openSQLResultAsDocument();
      break;
    case "sql.copyCell":
      await copySQLContextCell();
      break;
    case "sql.copyRow":
      await copySQLContextRow(false);
      break;
    case "sql.copyRowWithHeaders":
      await copySQLContextRow(true);
      break;
    case "sql.copyColumn":
      await copySQLContextColumn(false);
      break;
    case "sql.copyColumnWithHeader":
      await copySQLContextColumn(true);
      break;
    case "sql.clearSelection":
      clearSQLResultSelection();
      status.textContent = t("sql.status.selectionCleared");
      break;
    case "view.profile":
      await openProfilePanel();
      break;
    case "profile.run":
      await openProfilePanel();
      await loadProfile(true);
      break;
    case "profile.cancel":
      cancelProfileRequest();
      break;
    case "diff.open":
      openDiffPanel();
      break;
    default:
      status.textContent = `Command ${commandID} has no UI handler.`;
  }
}

function drawHeaders(target, width) {
  const headerNameY = Math.max(14, Math.round(headerHeight * 0.34));
  const headerTypeY = Math.min(headerHeight - 10, Math.round(headerHeight * 0.72));
  target.save();
  target.beginPath();
  target.rect(0, 0, width, headerHeight);
  target.clip();
  target.fillStyle = themeColors.gridHeader;
  target.fillRect(0, 0, width, headerHeight);
  target.fillStyle = themeColors.gridGutter;
  target.fillRect(0, 0, rowNumberWidth, headerHeight);
  target.strokeStyle = themeColors.gridLine;
  target.beginPath();
  target.moveTo(0, headerHeight - 0.5);
  target.lineTo(width, headerHeight - 0.5);
  target.moveTo(rowNumberWidth + 0.5, 0);
  target.lineTo(rowNumberWidth + 0.5, headerHeight);
  target.stroke();

  target.font = "600 12px Inter, system-ui, sans-serif";
  target.textBaseline = "middle";

  if (columns.length === 0) {
    target.fillStyle = themeColors.muted;
    target.fillText("Open a CSV file", rowNumberWidth + 10, headerHeight / 2);
    target.restore();
    return;
  }

  let x = rowNumberWidth - scrollX;
  for (const column of columns) {
    if (x + column.width >= rowNumberWidth && x <= width) {
      target.fillStyle = themeColors.text;
      drawHeaderText(target, column.name, x, column.width, headerNameY, width);
      target.fillStyle = themeColors.muted;
      target.font = "600 11px Inter, system-ui, sans-serif";
      drawHeaderText(target, column.type ?? "string", x, column.width, headerTypeY, width);
      target.font = "600 12px Inter, system-ui, sans-serif";
      drawHeaderMenuGlyph(target, x, column.width, headerNameY);
      target.strokeStyle = themeColors.gridLine;
      target.beginPath();
      target.moveTo(x + column.width + 0.5, 0);
      target.lineTo(x + column.width + 0.5, headerHeight);
      target.stroke();
    }
    x += column.width;
  }
  drawColumnReorderPreview(target, width);
  target.restore();
}

function drawColumnReorderPreview(target, canvasWidth) {
  const drag = pendingColumnReorder;
  if (!drag?.active || !columns[drag.sourceVisibleCol] || !columns[drag.targetVisibleCol]) return;
  const sourceLeft = columnLeft(drag.sourceVisibleCol);
  const sourceWidth = columns[drag.sourceVisibleCol].width;
  const targetLeft = columnLeft(drag.targetVisibleCol);
  const targetWidth = columns[drag.targetVisibleCol].width;
  const insertX = drag.targetVisibleCol > drag.sourceVisibleCol ? targetLeft + targetWidth : targetLeft;
  target.save();
  const sourceX = Math.max(rowNumberWidth, sourceLeft);
  const sourceRight = Math.min(canvasWidth, sourceLeft + sourceWidth);
  if (sourceRight > sourceX) {
    target.globalAlpha = 0.18;
    target.fillStyle = themeColors.accent;
    target.fillRect(sourceX, 0, sourceRight - sourceX, headerHeight);
  }
  target.globalAlpha = 1;
  target.strokeStyle = themeColors.accent;
  target.lineWidth = 3;
  target.beginPath();
  target.moveTo(Math.max(rowNumberWidth, Math.min(canvasWidth, insertX)) + 0.5, 0);
  target.lineTo(Math.max(rowNumberWidth, Math.min(canvasWidth, insertX)) + 0.5, headerHeight);
  target.stroke();
  target.restore();
}

function drawHeaderText(target, text, x, columnWidth, y, canvasWidth) {
  const left = Math.max(x + 8, rowNumberWidth + 8);
  const right = Math.min(x + columnWidth - 28, canvasWidth - 4);
  const width = right - left;
  if (width <= 0) return;
  target.save();
  target.beginPath();
  target.rect(left, 0, width, headerHeight);
  target.clip();
  target.fillText(text, left, y);
  target.restore();
}

function drawHeaderMenuGlyph(target, x, width, centerY) {
  const centerX = x + width - 14;
  if (centerX < rowNumberWidth + 8) return;
  target.save();
  target.fillStyle = themeColors.muted;
  target.beginPath();
  target.moveTo(centerX - 4, centerY - 2);
  target.lineTo(centerX + 4, centerY - 2);
  target.lineTo(centerX, centerY + 3);
  target.closePath();
  target.fill();
  target.restore();
}

function drawBodySeparators(width, height) {
  if (columns.length === 0 || height <= headerHeight) return;

  ctx.save();
  ctx.beginPath();
  ctx.rect(0, headerHeight, width, height - headerHeight);
  ctx.clip();
  ctx.strokeStyle = themeColors.gridLine;
  ctx.beginPath();
  ctx.moveTo(rowNumberWidth + 0.5, headerHeight);
  ctx.lineTo(rowNumberWidth + 0.5, height);
  ctx.stroke();

  ctx.strokeStyle = themeColors.gridLineSoft;
  ctx.beginPath();
  let x = rowNumberWidth - scrollX;
  for (const column of columns) {
    const lineX = x + column.width + 0.5;
    if (lineX >= rowNumberWidth && lineX <= width) {
      ctx.moveTo(lineX, headerHeight);
      ctx.lineTo(lineX, height);
    }
    x += column.width;
  }
  ctx.stroke();
  ctx.restore();
}

function drawRows(width, height) {
  if (columns.length === 0 || totalRows === 0) {
    ctx.fillStyle = themeColors.muted;
    ctx.font = "13px Inter, system-ui, sans-serif";
    ctx.textBaseline = "top";
    ctx.fillText("No file open", rowNumberWidth + 10, headerHeight + 14);
    return;
  }

  const firstRow = Math.max(0, Math.floor(scrollY / rowHeight));
  const yOffset = scrollY % rowHeight;
  const visibleRows = Math.ceil((height - headerHeight) / rowHeight) + 1;
  maybeLoadWindow(firstRow);

  ctx.save();
  ctx.beginPath();
  ctx.rect(0, headerHeight, width, Math.max(0, height - headerHeight));
  ctx.clip();

  ctx.font = gridBodyFont();
  ctx.textBaseline = "middle";

  for (let i = 0; i < visibleRows; i++) {
    const row = firstRow + i;
    if (row >= totalRows) break;
    const y = headerHeight + i * rowHeight - yOffset;
    const aggregateRow = aggregateRowAt(row);

    ctx.fillStyle = aggregateRow
      ? themeColors.panelSubtle
      : (row % 2 === 0 ? themeColors.gridRow : themeColors.gridRowAlt);
    ctx.fillRect(0, y, width, rowHeight);
    ctx.fillStyle = themeColors.gridHeader;
    ctx.fillRect(0, y, rowNumberWidth, rowHeight);
    ctx.fillStyle = themeColors.muted;
    ctx.fillText(aggregateRow ? `agg${row - sourceRowCount + 1}` : String(row + 1), 10, y + rowHeight / 2);

    let x = rowNumberWidth - scrollX;
    for (let col = 0; col < columns.length; col++) {
      const column = columns[col];
      if (x + column.width >= rowNumberWidth && x <= width) {
        const hasCell = hasCellAt(row, col);
        if (!hasCell) {
          ctx.fillStyle = themeColors.gridBg;
          ctx.fillRect(x + 1, y + 1, column.width - 2, rowHeight - 2);
          x += column.width;
          continue;
        }
        const value = valueAt(row, col);
        ctx.fillStyle = themeColors.text;
        ctx.save();
        ctx.beginPath();
        ctx.rect(x + 1, y + 1, column.width - 2, rowHeight - 2);
        ctx.clip();
        drawFindHighlights(value, row, col, x, y, column.width);
        ctx.fillStyle = themeColors.text;
        ctx.fillText(value, x + 8, y + rowHeight / 2);
        if (shouldQuote(value, row, col)) {
          drawQuoteIndicator(x, y, column.width);
        }
        ctx.restore();
      }
      x += column.width;
    }

    ctx.strokeStyle = themeColors.gridLineSoft;
    ctx.beginPath();
    ctx.moveTo(0, y + rowHeight + 0.5);
    ctx.lineTo(width, y + rowHeight + 0.5);
    ctx.stroke();
  }
  ctx.restore();
}

function drawQuoteIndicator(x, y, width) {
  ctx.save();
  ctx.fillStyle = themeColors.quoteMarker;
  ctx.beginPath();
  ctx.moveTo(x + width - 15, y + 1);
  ctx.lineTo(x + width - 2, y + 1);
  ctx.lineTo(x + width - 2, y + 14);
  ctx.closePath();
  ctx.fill();
  ctx.fillStyle = themeColors.quoteMarkerText;
  ctx.font = "700 9px Inter, system-ui, sans-serif";
  ctx.textBaseline = "top";
  ctx.textAlign = "right";
  ctx.fillText('"', x + width - 4, y + 2);
  ctx.restore();
}

function shouldQuote(value, displayRow = selected.row, visibleCol = selected.col) {
  const mode = quoteModeValue();
  if (suppressesQuote(displayRow, visibleCol)) return false;
  if (mode === "always") return true;
  if (mode === "never") return false;
  if (mode === "manual") return preservesQuote(displayRow, visibleCol);
  return requiresQuote(value) || preservesQuote(displayRow, visibleCol);
}

function preservesQuote(displayRow, visibleCol) {
  const dialect = metadata?.dialect;
  if (!dialect || visibleCol < 0) return false;
  const sourceRow = rowSource.get(displayRow) ?? displayRow + dataStartRow;
  const sourceCol = sourceColFor(visibleCol);
  if (dialect.hasHeader && dialect.quoteHeader && sourceRow === 0) return true;
  if (dialect.quoteBody && isBodySourceRow(sourceRow)) return true;
  if (isSourceColumnQuoted(sourceCol) && isBodySourceRow(sourceRow)) return true;
  if (Array.isArray(dialect.quoteCells)) {
    return dialect.quoteCells.some((cell) => cell.row === sourceRow && cell.col === sourceCol);
  }
  return false;
}

function isBodySourceRow(sourceRow) {
  return !metadata?.dialect?.hasHeader || sourceRow > 0;
}

function isSourceColumnQuoted(sourceCol) {
  return Array.isArray(metadata?.dialect?.quoteColumns) && metadata.dialect.quoteColumns.includes(sourceCol);
}

function suppressesQuote(displayRow, visibleCol) {
  const dialect = metadata?.dialect;
  if (!dialect || visibleCol < 0) return false;
  const sourceRow = rowSource.get(displayRow) ?? displayRow + dataStartRow;
  const sourceCol = sourceColFor(visibleCol);
  if (!Array.isArray(dialect.noQuoteCells)) return false;
  return dialect.noQuoteCells.some((cell) => cell.row === sourceRow && cell.col === sourceCol);
}

function requiresQuote(value) {
  const text = String(value ?? "");
  if (text === "") return false;
  const delimiter = delimiterChar();
  const quote = quoteCharLiteral();
  return (
    text.includes(delimiter) ||
    (quote !== "" && text.includes(quote)) ||
    text.includes("\n") ||
    text.includes("\r") ||
    /^\s/u.test(text)
  );
}

function drawFindHighlights(value, displayRow, visibleCol, x, y, width) {
  const ranges = findRangesForCell(value, displayRow, visibleCol);
  if (!ranges.length) return;
  const active = isActiveFindCell(displayRow, visibleCol);
  ctx.save();
  ctx.fillStyle = active ? themeColors.findActiveHighlight : themeColors.findHighlight;
  for (const range of ranges) {
    const before = String(value ?? "").slice(0, range.start);
    const match = String(value ?? "").slice(range.start, range.end);
    const left = x + 8 + ctx.measureText(before).width;
    const highlightWidth = Math.max(2, ctx.measureText(match).width);
    ctx.fillRect(left - 1, y + 4, highlightWidth + 2, rowHeight - 8);
  }
  ctx.restore();
}

function findRangesForCell(value, displayRow, visibleCol) {
  if (!findState.matches.length || !findState.search) return [];
  const sourceRow = rowSource.get(displayRow) ?? displayRow + dataStartRow;
  const sourceCol = sourceColFor(visibleCol);
  const matches = findState.matches.filter((match) => match.row === sourceRow && match.col === sourceCol);
  if (!matches.length) {
    return [];
  }
  const serverRanges = matches
    .flatMap((match) => normalizeServerFindRanges(value, match));
  if (serverRanges.length) return serverRanges;
  return findTextRanges(String(value ?? ""), findState.search);
}

function normalizeServerFindRanges(value, match) {
  if (!match || typeof match !== "object") return [];
  return rawServerMatchRanges(match)
    .map((range) => normalizeServerFindRange(value, range))
    .filter(Boolean);
}

function rawServerMatchRanges(match) {
  if (!match || typeof match !== "object") return [];
  return Array.isArray(match.ranges) && match.ranges.length
    ? match.ranges
    : [match.range];
}

function normalizeServerFindRange(value, range) {
  if (!range || typeof range !== "object") return null;
  const text = String(value ?? "");
  const start = Number(range.start);
  const end = Number(range.end);
  if (!Number.isInteger(start) || !Number.isInteger(end)) return null;
  if (start < 0 || end <= start || end > text.length) return null;
  return { start, end };
}

function formatServerMatchRanges(match) {
  if (!match || typeof match !== "object") return [];
  return rawServerMatchRanges(match)
    .map((range) => {
      const start = Number(range?.start);
      const end = Number(range?.end);
      if (!Number.isInteger(start) || !Number.isInteger(end) || start < 0 || end <= start) return "";
      return `${match.row}:${match.col}:${start}-${end}`;
    })
    .filter(Boolean);
}

function isActiveFindCell(displayRow, visibleCol) {
  const match = findState.matches[findState.index];
  if (!match) return false;
  const sourceRow = rowSource.get(displayRow) ?? displayRow + dataStartRow;
  return match.row === sourceRow && match.col === sourceColFor(visibleCol);
}

function findTextRanges(value, search) {
  if (!value) return [];
  if (search.regex) return regexFindRanges(value, search.query, search.caseSensitive);
  if (search.wholeCell) {
    const haystack = search.caseSensitive ? value : value.toLowerCase();
    const needle = search.caseSensitive ? search.query : search.query.toLowerCase();
    return haystack === needle ? [{ start: 0, end: value.length }] : [];
  }
  return substringFindRanges(value, search.query, search.caseSensitive);
}

function substringFindRanges(value, query, caseSensitive) {
  if (!query) return [];
  const haystack = caseSensitive ? value : value.toLowerCase();
  const needle = caseSensitive ? query : query.toLowerCase();
  const ranges = [];
  let start = 0;
  while (start <= haystack.length) {
    const index = haystack.indexOf(needle, start);
    if (index < 0) break;
    ranges.push({ start: index, end: index + needle.length });
    start = index + Math.max(1, needle.length);
    if (ranges.length >= 20) break;
  }
  return ranges;
}

function regexFindRanges(value, query, caseSensitive) {
  try {
    const flags = caseSensitive ? "g" : "gi";
    const regex = new RegExp(query, flags);
    const ranges = [];
    for (const match of value.matchAll(regex)) {
      if (match.index == null) continue;
      const end = match.index + match[0].length;
      if (end === match.index) continue;
      ranges.push({ start: match.index, end });
      if (ranges.length >= 20) break;
    }
    return ranges;
  } catch {
    return [];
  }
}

function drawSelection() {
  if (columns.length === 0) return;
  ctx.save();
  ctx.beginPath();
  ctx.rect(0, headerHeight, canvas.clientWidth, Math.max(0, canvas.clientHeight - headerHeight));
  ctx.clip();

  ctx.fillStyle = themeColors.selectionBg;
  for (const range of activeSelectionRanges()) {
    const firstVisibleRow = Math.max(range.startRow, Math.floor(scrollY / rowHeight));
    const lastVisibleRow = Math.min(
      range.endRow,
      Math.ceil((scrollY + canvas.clientHeight - headerHeight) / rowHeight),
    );
    const visibleColumns = visibleColumnIndexes(range);
    for (let row = firstVisibleRow; row <= lastVisibleRow; row++) {
      for (const col of visibleColumns) {
        const rect = cellRect(row, col);
        if (!rect) continue;
        ctx.fillRect(rect.x + 1, rect.y + 1, rect.width - 2, rect.height - 2);
      }
    }
  }

  const rect = cellRect(selected.row, selected.col);
  if (rect) {
    ctx.strokeStyle = themeColors.accent;
    ctx.lineWidth = 2;
    ctx.strokeRect(rect.x + 1, rect.y + 1, rect.width - 2, rect.height - 2);
    ctx.lineWidth = 1;
  }
  drawFillHandle(normalizedSelection());
  ctx.restore();
}

function drawRowReorderPreview(width, height) {
  const drag = pendingRowReorder;
  if (!drag?.active || sourceRowCount <= 0) return;
  const sourceTop = rowTop(drag.sourceDisplayRow);
  const targetTop = rowTop(drag.targetDisplayRow);
  const insertY = drag.targetDisplayRow > drag.sourceDisplayRow ? targetTop + rowHeight : targetTop;
  ctx.save();
  ctx.beginPath();
  ctx.rect(0, headerHeight, width, Math.max(0, height - headerHeight));
  ctx.clip();
  const sourceY = Math.max(headerHeight, sourceTop);
  const sourceBottom = Math.min(height, sourceTop + rowHeight);
  if (sourceBottom > sourceY) {
    ctx.globalAlpha = 0.18;
    ctx.fillStyle = themeColors.accent;
    ctx.fillRect(0, sourceY, width, sourceBottom - sourceY);
  }
  ctx.globalAlpha = 1;
  ctx.strokeStyle = themeColors.accent;
  ctx.lineWidth = 3;
  const y = Math.max(headerHeight, Math.min(height, insertY)) + 0.5;
  ctx.beginPath();
  ctx.moveTo(0, y);
  ctx.lineTo(width, y);
  ctx.stroke();
  ctx.restore();
}

function drawFillHandle(range) {
  const handle = fillHandleRect(range);
  if (!handle) return;
  ctx.save();
  ctx.fillStyle = themeColors.accent;
  ctx.strokeStyle = themeColors.gridBg;
  ctx.lineWidth = 2;
  ctx.beginPath();
  ctx.rect(handle.x + 0.5, handle.y + 0.5, handle.width, handle.height);
  ctx.fill();
  ctx.stroke();
  ctx.restore();
}

function fillHandleRect(range = normalizedSelection()) {
  if (smartSelectionRanges.length > 0) return null;
  if (!documentId || readOnly || columns.length === 0 || totalRows === 0 || sourceRowCount <= 0) {
    return null;
  }
  if (range.endRow >= sourceRowCount) return null;
  const rect = cellRect(range.endRow, range.endCol);
  if (!rect) return null;
  const size = 8;
  const x = Math.min(rect.x + rect.width - size - 2, canvas.clientWidth - size - 2);
  const y = Math.min(rect.y + rect.height - size - 2, canvas.clientHeight - size - 2);
  if (x + size < rowNumberWidth || y + size < headerHeight) return null;
  return { x, y, width: size, height: size };
}

function fillHandleAt(x, y) {
  const handle = fillHandleRect();
  if (!handle) return false;
  const slop = 4;
  return (
    x >= handle.x - slop &&
    x <= handle.x + handle.width + slop &&
    y >= handle.y - slop &&
    y <= handle.y + handle.height + slop
  );
}

function visibleColumnIndexes(range = normalizedSelection()) {
  const visible = [];
  let x = rowNumberWidth - scrollX;
  for (let col = 0; col < columns.length; col++) {
    const width = columns[col].width;
    if (
      col >= range.startCol &&
      col <= range.endCol &&
      x + width >= rowNumberWidth &&
      x <= canvas.clientWidth
    ) {
      visible.push(col);
    }
    x += width;
  }
  return visible;
}

function cellRect(row, col) {
  if (!columns[col]) return null;
  const y = headerHeight + row * rowHeight - scrollY;
  let x = rowNumberWidth - scrollX;
  for (let i = 0; i < col; i++) x += columns[i].width;
  const width = columns[col]?.width ?? 120;
  if (y + rowHeight < headerHeight || y > canvas.clientHeight) return null;
  if (x + width < rowNumberWidth || x > canvas.clientWidth) return null;
  return { x, y, width, height: rowHeight };
}

function rowTop(row) {
  return headerHeight + row * rowHeight - scrollY;
}

function columnLeft(col) {
  let x = rowNumberWidth - scrollX;
  for (let index = 0; index < col; index++) x += columns[index]?.width ?? 0;
  return x;
}

function columnAt(x) {
  let current = rowNumberWidth - scrollX;
  for (let col = 0; col < columns.length; col++) {
    const width = columns[col].width;
    if (x >= current && x < current + width) return col;
    current += width;
  }
  return -1;
}

function rowAt(y) {
  if (totalRows <= 0) return -1;
  const row = Math.floor((y - headerHeight + scrollY) / rowHeight);
  return Math.max(0, Math.min(totalRows - 1, row));
}

function dataRowAtOrNearest(y) {
  if (sourceRowCount <= 0) return -1;
  const row = Math.floor((y - headerHeight + scrollY) / rowHeight);
  return Math.max(0, Math.min(sourceRowCount - 1, row));
}

function clampColumn(col) {
  return Math.max(0, Math.min(columns.length - 1, col));
}

function isHeaderMenuTarget(x, col) {
  if (!columns[col]) return false;
  const left = columnLeft(col);
  return x >= left + columns[col].width - 28;
}

function isWholeColumnSelected(col) {
  if (!documentId || columns.length === 0 || sourceRowCount <= 0 || smartSelectionRanges.length > 0) return false;
  const range = normalizedSelection();
  return (
    range.startRow === 0 &&
    range.endRow === Math.max(0, sourceRowCount - 1) &&
    range.startCol === col &&
    range.endCol === col
  );
}

function isWholeRowSelected(row) {
  if (!documentId || columns.length === 0 || sourceRowCount <= 0 || row < 0 || row >= sourceRowCount || smartSelectionRanges.length > 0) {
    return false;
  }
  const range = normalizedSelection();
  return (
    range.startRow === row &&
    range.endRow === row &&
    range.startCol === 0 &&
    range.endCol === columns.length - 1
  );
}

function headerResizeColumnAt(x, y) {
  if (y < 0 || y >= headerHeight || x < rowNumberWidth) return -1;
  let left = rowNumberWidth - scrollX;
  const hitSlop = 5;
  for (let col = 0; col < columns.length; col++) {
    const right = left + columns[col].width;
    if (
      Math.abs(x - right) <= hitSlop &&
      right >= rowNumberWidth + hitSlop &&
      right <= canvas.clientWidth + hitSlop
    ) {
      return col;
    }
    left = right;
  }
  return -1;
}

function columnAtOrNearest(x) {
  if (columns.length === 0) return -1;
  const col = columnAt(x);
  if (col >= 0) return col;
  return x < rowNumberWidth ? 0 : columns.length - 1;
}

function startColumnReorder(col, event) {
  if (!columns[col]) return false;
  pendingColumnReorder = {
    sourceVisibleCol: col,
    sourceCol: sourceColFor(col),
    startClientX: event.clientX,
    startClientY: event.clientY,
    targetVisibleCol: col,
    active: false,
  };
  document.body.dataset.columnDragReorder = `pending:${col}`;
  document.body.dataset.columnDragTarget = String(col);
  clearURLHover();
  clearResizeHover();
  clearFillHandleHover();
  syncCanvasHoverChrome();
  return true;
}

function updateColumnReorder(event) {
  if (!pendingColumnReorder) return;
  const drag = pendingColumnReorder;
  const distance = Math.hypot(event.clientX - drag.startClientX, event.clientY - drag.startClientY);
  if (!drag.active && distance < columnReorderDragThreshold) return;
  if (!drag.active) {
    drag.active = true;
    status.textContent = `Move ${columns[drag.sourceVisibleCol]?.name || `Column ${drag.sourceVisibleCol + 1}`}...`;
  }
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  drag.targetVisibleCol = columnAtOrNearest(x);
  document.body.dataset.columnDragReorder = `${drag.sourceVisibleCol}->${drag.targetVisibleCol}:drag`;
  document.body.dataset.columnDragTarget = String(drag.targetVisibleCol);
  syncCanvasHoverChrome();
  draw();
}

async function finishColumnReorder() {
  if (!pendingColumnReorder) return;
  const drag = pendingColumnReorder;
  pendingColumnReorder = null;
  syncCanvasHoverChrome();
  if (!drag.active) {
    document.body.dataset.columnDragReorder = "";
    document.body.dataset.columnDragTarget = "";
    draw();
    return;
  }
  const targetVisibleCol = Math.max(0, Math.min(columns.length - 1, drag.targetVisibleCol));
  if (targetVisibleCol === drag.sourceVisibleCol) {
    document.body.dataset.columnDragReorder = `${drag.sourceVisibleCol}->${targetVisibleCol}:cancel`;
    document.body.dataset.columnDragTarget = String(targetVisibleCol);
    status.textContent = "Column move canceled.";
    draw();
    return;
  }
  await moveColumnToVisible(drag.sourceVisibleCol, targetVisibleCol, "drag");
}

function startRowReorder(row, event) {
  if (row < 0 || row >= sourceRowCount) return false;
  pendingRowReorder = {
    sourceDisplayRow: row,
    sourceRow: rowSource.get(row) ?? row + dataStartRow,
    startClientX: event.clientX,
    startClientY: event.clientY,
    targetDisplayRow: row,
    active: false,
  };
  document.body.dataset.rowDragReorder = `pending:${row}`;
  document.body.dataset.rowDragTarget = String(row);
  clearURLHover();
  clearResizeHover();
  clearFillHandleHover();
  syncCanvasHoverChrome();
  return true;
}

function updateRowReorder(event) {
  if (!pendingRowReorder) return;
  const drag = pendingRowReorder;
  const distance = Math.hypot(event.clientX - drag.startClientX, event.clientY - drag.startClientY);
  if (!drag.active && distance < rowReorderDragThreshold) return;
  if (!drag.active) {
    drag.active = true;
    status.textContent = `Move row ${drag.sourceDisplayRow + 1}...`;
  }
  const rect = canvas.getBoundingClientRect();
  const y = event.clientY - rect.top;
  drag.targetDisplayRow = dataRowAtOrNearest(y);
  document.body.dataset.rowDragReorder = `${drag.sourceDisplayRow}->${drag.targetDisplayRow}:drag`;
  document.body.dataset.rowDragTarget = String(drag.targetDisplayRow);
  syncCanvasHoverChrome();
  draw();
}

async function finishRowReorder() {
  if (!pendingRowReorder) return;
  const drag = pendingRowReorder;
  pendingRowReorder = null;
  syncCanvasHoverChrome();
  if (!drag.active) {
    document.body.dataset.rowDragReorder = "";
    document.body.dataset.rowDragTarget = "";
    draw();
    return;
  }
  const targetDisplayRow = Math.max(0, Math.min(sourceRowCount - 1, drag.targetDisplayRow));
  if (targetDisplayRow === drag.sourceDisplayRow) {
    document.body.dataset.rowDragReorder = `${drag.sourceDisplayRow}->${targetDisplayRow}:cancel`;
    document.body.dataset.rowDragTarget = String(targetDisplayRow);
    status.textContent = "Row move canceled.";
    draw();
    return;
  }
  await moveRowToDisplay(drag.sourceDisplayRow, targetDisplayRow, "drag");
}

function selectRowRange(row, extend = false) {
  if (totalRows <= 0 || columns.length === 0) return;
  const clampedRow = Math.max(0, Math.min(totalRows - 1, row));
  if (extend) {
    selectionAnchor = { row: selectionAnchor.row, col: columns.length - 1 };
  } else {
    selectionAnchor = { row: clampedRow, col: columns.length - 1 };
  }
  selected = { row: clampedRow, col: 0 };
}

function selectColumnRange(col, extend = false) {
  if (totalRows <= 0 || columns.length === 0) return;
  const clampedCol = clampColumn(col);
  const lastDataRow = Math.max(0, sourceRowCount - 1);
  if (extend) {
    selectionAnchor = { row: lastDataRow, col: selectionAnchor.col };
  } else {
    selectionAnchor = { row: lastDataRow, col: clampedCol };
  }
  selected = { row: 0, col: clampedCol };
}

function selectAllBodyCells() {
  if (totalRows <= 0 || columns.length === 0) return;
  const lastDataRow = Math.max(0, sourceRowCount - 1);
  selectionAnchor = { row: 0, col: 0 };
  selected = { row: lastDataRow, col: columns.length - 1 };
}

function startColumnResize(col, clientX) {
  if (!columns[col]) return;
  resizingColumn = col;
  resizeStartClientX = clientX;
  resizeStartWidth = columns[col].width;
  hoverResizeColumn = col;
  closeColumnMenu();
  document.body.dataset.resizeColumn = String(col);
  document.body.dataset.columnResize = `${col}:${Math.round(resizeStartWidth)}`;
  syncCanvasHoverChrome();
}

function updateColumnResize(clientX) {
  if (resizingColumn < 0 || !columns[resizingColumn]) return;
  const width = clampColumnWidth(resizeStartWidth + clientX - resizeStartClientX);
  setColumnWidth(resizingColumn, width);
  document.body.dataset.columnResize = `${resizingColumn}:${Math.round(width)}`;
  status.textContent = `${columns[resizingColumn].name || `Column ${resizingColumn + 1}`} width ${Math.round(width)}px`;
  draw();
}

function finishColumnResize() {
  if (resizingColumn < 0 || !columns[resizingColumn]) return;
  const col = resizingColumn;
  const width = columns[col].width;
  resizingColumn = -1;
  hoverResizeColumn = -1;
  document.body.dataset.resizeColumn = "";
  document.body.dataset.columnResize = `${col}:${Math.round(width)}`;
  syncCanvasHoverChrome();
  draw();
  status.textContent = `${columns[col].name || `Column ${col + 1}`} width ${Math.round(width)}px`;
}

function autoFitVisibleColumn(col) {
  const column = columns[col];
  if (!column) return;
  const sourceCol = sourceColFor(col);
  const values = Array.from(rows.values(), (record) => record[sourceCol] ?? "");
  const width = autoFitColumnWidth(column, values);
  setColumnWidth(col, width);
  document.body.dataset.columnResize = `${col}:${Math.round(width)}:fit`;
  draw();
  status.textContent = `${column.name || `Column ${col + 1}`} auto-fit to ${Math.round(width)}px`;
}

function setColumnWidth(col, width) {
  const column = columns[col];
  if (!column) return;
  const nextWidth = clampColumnWidth(width);
  column.width = nextWidth;
  const sourceColumn = allColumns.find((candidate) => candidate.sourceCol === column.sourceCol);
  if (sourceColumn) sourceColumn.width = nextWidth;
}

function clampColumnWidth(width) {
  return Math.max(minColumnWidth, Math.min(maxColumnWidth, Math.round(width)));
}

function updateDragSelection(event) {
  if (!draggingSelection) return;
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  if (dragSelectionMode === "row") {
    selected = { row: rowAt(y), col: 0 };
  } else if (dragSelectionMode === "column") {
    selected = { row: 0, col: columnAtOrNearest(x) };
  } else {
    selected = { row: rowAt(y), col: columnAtOrNearest(x) };
  }
  draw();
}

function startFillDrag() {
  const range = normalizedSelection();
  if (range.endRow >= sourceRowCount) return false;
  draggingFillHandle = true;
  hoverFillHandle = true;
  fillDragSourceRange = { ...range };
  fillDragDirection = "";
  document.body.dataset.fillHandle = selectionSummary(range);
  document.body.dataset.fillDrag = "";
  clearURLHover();
  clearResizeHover();
  syncCanvasHoverChrome();
  return true;
}

function updateFillDrag(event) {
  if (!draggingFillHandle || !fillDragSourceRange) return;
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  const row = Math.max(0, Math.min(sourceRowCount - 1, rowAt(y)));
  const col = columnAtOrNearest(x);
  const next = fillDragSelection(fillDragSourceRange, row, col);
  fillDragDirection = next.direction;
  selectionAnchor = next.anchor;
  selected = next.selected;
  document.body.dataset.fillDrag = next.direction ? `${next.direction}:${selectionSummary(next.range)}` : "";
  draw();
}

function finishFillDrag() {
  if (!draggingFillHandle) return;
  const direction = fillDragDirection;
  draggingFillHandle = false;
  hoverFillHandle = false;
  fillDragSourceRange = null;
  fillDragDirection = "";
  document.body.dataset.fillHandle = "";
  document.body.dataset.fillDrag = "";
  syncCanvasHoverChrome();
  if (direction) {
    autoFillSelection(direction);
    return;
  }
  draw();
}

function fillDragSelection(sourceRange, row, col) {
  const clampedRow = Math.max(0, Math.min(sourceRowCount - 1, row));
  const clampedCol = clampColumn(col);
  const down = clampedRow > sourceRange.endRow ? clampedRow - sourceRange.endRow : 0;
  const up = clampedRow < sourceRange.startRow ? sourceRange.startRow - clampedRow : 0;
  const right = clampedCol > sourceRange.endCol ? clampedCol - sourceRange.endCol : 0;
  const left = clampedCol < sourceRange.startCol ? sourceRange.startCol - clampedCol : 0;
  const vertical = Math.max(down, up);
  const horizontal = Math.max(right, left);
  let direction = "";
  if (vertical > 0 || horizontal > 0) {
    direction = vertical >= horizontal ? (down >= up ? "down" : "up") : (right >= left ? "right" : "left");
  }

  if (direction === "down") {
    const range = { ...sourceRange, endRow: clampedRow };
    return {
      direction,
      range,
      anchor: { row: sourceRange.startRow, col: sourceRange.startCol },
      selected: { row: range.endRow, col: sourceRange.endCol },
    };
  }
  if (direction === "up") {
    const range = { ...sourceRange, startRow: clampedRow };
    return {
      direction,
      range,
      anchor: { row: sourceRange.endRow, col: sourceRange.endCol },
      selected: { row: range.startRow, col: sourceRange.startCol },
    };
  }
  if (direction === "right") {
    const range = { ...sourceRange, endCol: clampedCol };
    return {
      direction,
      range,
      anchor: { row: sourceRange.startRow, col: sourceRange.startCol },
      selected: { row: sourceRange.endRow, col: range.endCol },
    };
  }
  if (direction === "left") {
    const range = { ...sourceRange, startCol: clampedCol };
    return {
      direction,
      range,
      anchor: { row: sourceRange.endRow, col: sourceRange.endCol },
      selected: { row: sourceRange.startRow, col: range.startCol },
    };
  }
  return {
    direction: "",
    range: sourceRange,
    anchor: { row: sourceRange.startRow, col: sourceRange.startCol },
    selected: { row: sourceRange.endRow, col: sourceRange.endCol },
  };
}

function openColumnMenu(col, clientX, clientY) {
  if (!documentId || !columns[col]) return;
  columnMenuColumn = col;
  selected = { row: Math.max(0, Math.min(totalRows - 1, selected.row)), col };
  selectionAnchor = { ...selected };
  columnMenuTitle.textContent = columns[col].name || `Column ${col + 1}`;
  const selectedValue = valueAt(selected.row, col);
  columnFilterValue.textContent =
    selectedValue === ""
      ? t("column.filterSelected")
      : tf("column.filterValue", { value: truncateMenuValue(selectedValue) });
  columnFilterValue.disabled = selectedValue === "";
  columnTransformApply.disabled = readOnly;
  syncColumnTransformCommand();
  if (columnStableSortToggle) {
    columnStableSortToggle.checked = sortStable;
  }
  document.body.dataset.columnMenuStableSort = String(sortStable);
  columnShowAll.disabled = hiddenSourceColumns.size === 0;
  renderColumnValuesLoading();
  columnMenu.hidden = false;

  const menuRect = columnMenu.getBoundingClientRect();
  const left = Math.min(Math.max(8, clientX), window.innerWidth - menuRect.width - 8);
  const top = Math.min(Math.max(8, clientY), window.innerHeight - menuRect.height - 8);
  columnMenu.style.left = `${left}px`;
  columnMenu.style.top = `${top}px`;
  draw();
  loadColumnMenuValues(col);
}

function closeColumnMenu() {
  if (columnMenu.hidden) return;
  columnMenu.hidden = true;
  columnMenuColumn = -1;
}

function truncateMenuValue(value) {
  const text = String(value ?? "").replace(/\s+/g, " ").trim();
  return text.length > 24 ? `${text.slice(0, 23)}...` : text || "(empty)";
}

function renderColumnValuesLoading() {
  columnMenuValues.replaceChildren();
  const node = document.createElement("div");
  node.className = "column-values-empty";
  node.textContent = t("column.valuesLoading");
  columnMenuValues.append(node);
  columnApplyValues.disabled = true;
}

async function loadColumnMenuValues(col) {
  if (!documentId || !columns[col]) return;
  const sourceCol = sourceColFor(col);
  try {
    const response = await fetch(`/api/column-values?id=${encodeURIComponent(documentId)}&column=${sourceCol}&limit=12`);
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "column values failed");
    if (columnMenuColumn !== col || columnMenu.hidden) return;
    document.body.dataset.columnValuesStreamed = String(Boolean(payload.streamed));
    renderColumnMenuValues(payload.values ?? [], payload.totalDistinct ?? 0);
  } catch (error) {
    columnMenuValues.replaceChildren();
    const node = document.createElement("div");
    node.className = "column-values-empty";
    node.textContent = tf("column.valuesFailed", { message: error.message });
    columnMenuValues.append(node);
    columnApplyValues.disabled = true;
  }
}

function renderColumnMenuValues(values, totalDistinct) {
  columnMenuValues.replaceChildren();
  if (!values.length) {
    const node = document.createElement("div");
    node.className = "column-values-empty";
    node.textContent = t("column.noValues");
    columnMenuValues.append(node);
    columnApplyValues.disabled = true;
    return;
  }
  const heading = document.createElement("span");
  heading.className = "column-values-heading";
  heading.textContent = totalDistinct > values.length
    ? tf("column.valuesShown", { shown: formatCount(values.length), total: formatCount(totalDistinct) })
    : tf("column.valuesCount", { count: formatCount(totalDistinct) });
  columnMenuValues.append(heading);
  for (const item of values) {
    const label = document.createElement("label");
    label.className = "column-value-option";
    const input = document.createElement("input");
    input.type = "checkbox";
    input.value = item.value ?? "";
    input.dataset.blank = item.blank ? "true" : "false";
    input.addEventListener("change", updateColumnApplyValuesState);
    const text = document.createElement("span");
    text.textContent = `${truncateMenuValue(item.blank ? t("column.blankValue") : item.value)} (${formatCount(item.count)})`;
    label.append(input, text);
    columnMenuValues.append(label);
  }
  updateColumnApplyValuesState();
}

function updateColumnApplyValuesState() {
  columnApplyValues.disabled = columnMenuValues.querySelectorAll("input:checked").length === 0;
}

function commitEditor() {
  if (!editing) return;
  const value = editor.value;
  editing = false;
  editor.style.display = "none";
  setSelectionValue(value);
  canvas.focus();
}

function openEditor() {
  if (!documentId) return;
  if (!ensureEditable("Edit")) return;
  if (aggregateRowAt(selected.row)) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  const rect = cellRect(selected.row, selected.col);
  if (!rect) return;
  editing = true;
  editor.value = valueAt(selected.row, selected.col);
  editor.style.display = "block";
  editor.style.left = `${rect.x}px`;
  editor.style.top = `${rect.y}px`;
  editor.style.width = `${rect.width}px`;
  editor.style.height = `${rect.height}px`;
  editor.focus();
  editor.select();
}

function setSelectionValue(value) {
  if (!ensureEditable("Edit")) return;
  const range = normalizedSelection();
  if (range.endRow >= sourceRowCount) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  windowGeneration++;
  const edits = [];
  for (let row = range.startRow; row <= range.endRow; row++) {
    const record = rows.get(row) ?? [];
    for (let col = range.startCol; col <= range.endCol; col++) {
      const sourceCol = sourceColFor(col);
      while (record.length <= sourceCol) record.push("");
      record[sourceCol] = value;
      edits.push({ row: rowSource.get(row) ?? row + dataStartRow, col: sourceCol, value });
    }
    rows.set(row, record);
  }
  autoFitColumnsForEdits(edits);
  syncCellValueAfterEdit();
  submitCellEdits(edits, range.startRow, `${formatCount(edits.length)} cells edited`);
}

function setActiveCellValue(value) {
  if (!ensureEditable("Edit")) return;
  if (!documentId || columns.length === 0) return;
  if (aggregateRowAt(selected.row)) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  windowGeneration++;
  const record = rows.get(selected.row) ?? [];
  const sourceCol = sourceColFor(selected.col);
  while (record.length <= sourceCol) record.push("");
  record[sourceCol] = value;
  rows.set(selected.row, record);
  selectionAnchor = { ...selected };
  cellValueInputDirty = false;
  const edit = { row: rowSource.get(selected.row) ?? selected.row + dataStartRow, col: sourceCol, value };
  autoFitColumnsForEdits([edit]);
  syncCellValueAfterEdit();
  submitCellEdits(
    [edit],
    selected.row,
    "Cell edited",
  );
}

function ensureEditable(action) {
  if (!readOnly) return true;
  status.textContent = `${action} blocked: document is read-only`;
  return false;
}

function autoFitColumnsForEdits(edits) {
  if (!autoFitColumnWidthOnEdit || !edits.length || columns.length === 0) return;
  const valuesBySourceCol = new Map();
  for (const edit of edits) {
    if (edit.col == null) continue;
    if (!valuesBySourceCol.has(edit.col)) valuesBySourceCol.set(edit.col, []);
    valuesBySourceCol.get(edit.col).push(edit.value);
  }
  const changed = [];
  for (const [sourceCol, values] of valuesBySourceCol.entries()) {
    const visibleCol = columns.findIndex((column) => column.sourceCol === sourceCol);
    if (visibleCol < 0) continue;
    const column = columns[visibleCol];
    const desiredWidth = autoFitColumnWidth(column, values);
    if (desiredWidth <= column.width) continue;
    const previousWidth = column.width;
    column.width = desiredWidth;
    const sourceColumn = allColumns.find((candidate) => candidate.sourceCol === sourceCol);
    if (sourceColumn) sourceColumn.width = desiredWidth;
    changed.push(`${sourceCol}:${Math.round(previousWidth)}>${Math.round(desiredWidth)}`);
  }
  if (changed.length) {
    document.body.dataset.autoFitColumnWidths = changed.join(",");
  }
}

function autoFitColumnWidth(column, values) {
  const headerWidth = measuredTextWidth(column.name, "600 12px Inter, system-ui, sans-serif") + 38;
  const typeWidth = measuredTextWidth(column.type ?? "string", "600 11px Inter, system-ui, sans-serif") + 16;
  const bodyWidth = Math.max(
    0,
    ...values.map(
      (value) =>
        longestLineWidth(
          value,
          gridBodyFont(),
        ) + 24,
    ),
  );
  return Math.max(
    minColumnWidth,
    Math.min(maxColumnWidth, Math.ceil(Math.max(headerWidth, typeWidth, bodyWidth))),
  );
}

function longestLineWidth(value, font) {
  const lines = String(value ?? "").split(/\r\n|\r|\n/);
  return Math.max(0, ...lines.map((line) => measuredTextWidth(line, font)));
}

function measuredTextWidth(value, font) {
  ctx.save();
  ctx.font = font;
  const width = ctx.measureText(String(value ?? "")).width;
  ctx.restore();
  return width;
}

function syncCellValueAfterEdit() {
  if (!documentId) return;
  cellValueInputKey = selectedCellKey();
  cellValueInputDirty = false;
  syncCellValuePanel();
}

function applyCellValue() {
  if (!cellValueInputDirty) return;
  setActiveCellValue(cellValueInput.value);
}

function submitCellEdits(edits, windowRow, message) {
  if (!edits.length) return;
  dirty = true;
  syncActiveFileEntryState();
  profileData = null;
  draw();
  if (!documentId) return;
  fetch(edits.length > 1 ? "/api/edit-batch" : "/api/edit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(
      edits.length > 1
        ? { documentId, edits }
        : { documentId, ...edits[0] },
    ),
  })
    .then(async (response) => {
      const payload = await response.json();
      if (!response.ok) throw new Error(payload.error ?? "edit failed");
      windowGeneration++;
      loadingWindow = false;
      rows = new Map();
      rowSource = new Map();
      windowStart = 0;
      windowEnd = 0;
      await maybeLoadWindow(windowRow);
      applyEditsToLoadedRows(edits);
      draw();
      if (message) status.textContent = message;
      scheduleDirtyDraftSave();
    })
    .catch((error) => {
      status.textContent = `Edit failed: ${error.message}`;
    });
}

function applyEditsToLoadedRows(edits) {
  for (const edit of edits) {
    const displayRow = displayRowForSourceRow(edit.row);
    if (displayRow < 0) continue;
    const record = rows.get(displayRow);
    if (!record) continue;
    while (record.length <= edit.col) record.push("");
    record[edit.col] = edit.value;
    rows.set(displayRow, record);
  }
}

function displayRowForSourceRow(sourceRow) {
  for (const [displayRow, mappedSourceRow] of rowSource.entries()) {
    if (mappedSourceRow === sourceRow) return displayRow;
  }
  if (!viewMapped) return sourceRow - dataStartRow;
  return -1;
}

async function openPath(path, options = {}) {
  captureActiveOpenFileTabState();
  saveCurrentScrollPosition();
  status.textContent = `Opening ${path}...`;
  let preset = options.preset ?? null;
  const formatOverride = options.formatOverride ?? null;
  const urlLocation = options.location ?? null;
  const displayName = options.displayName || fileNameFromPath(path);
  const remember = options.remember !== false;
  let matchedRule = null;
  if (!preset && options.applyRules !== false) {
    const match = await matchingFormatPresetForPath(path);
    if (match) {
      preset = match.preset;
      matchedRule = match.rule;
      document.body.dataset.formatRulePreview = `${formatRuleTarget(matchedRule.pattern)}:${formatRulePattern(matchedRule.pattern)}`;
    }
  }
  const body = { path, previewRows: windowSize, indexRows: -1 };
  if (options.sourceKind) body.sourceKind = options.sourceKind;
  if (options.displayName) body.displayName = options.displayName;
  if (preset) {
    body.delimiter = preset.delimiter;
    body.delimiterName = preset.delimiterName;
    body.hasHeader = Boolean(preset.hasHeader);
    body.headerColumns = normalizeHeaderColumns(preset.headerColumns);
    body.encodingName = preset.encodingName;
    body.lineEnding = normalizeLineEndingName(preset.lineEndingName);
    body.quoteChar = preset.quoteChar || "double";
    body.quoteMode = preset.quoteMode || "minimal";
    body.quoteEscape = normalizeQuoteEscapeValue(preset.quoteEscape);
  }
  applyFormatOverride(body, formatOverride);
  const response = await fetch("/api/open", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const payload = await response.json();
  if (!response.ok) throw new Error(payload.error ?? "open failed");

  documentId = payload.documentId;
  metadata = payload.document.metadata;
  readOnly = false;
  externalFileChanged = false;
  clearStructureClipboard();
  loadDocument(payload.document);
  updateFileList(displayName, options.detail || metadata.path, totalRows, {
    remember,
    documentId,
    sourceKind: payload.sourceKind || options.sourceKind || "file",
  });
  updateFormatRulePlaceholder(path);
  if (options.sourceKind) {
    document.body.dataset.openSourceKind = options.sourceKind;
  } else {
    delete document.body.dataset.openSourceKind;
  }
  document.body.dataset.openDisplayName = displayName;
  delete document.body.dataset.newFile;
  draw();
  startFileStatusChecks();
  await checkFileStatus(false);
  const locationApplied = await applyURLLocation(urlLocation);
  const statusMessages = [];
  if (hasFormatOverride(formatOverride)) statusMessages.push(`Applied URL format: ${documentFormatSummary()}`);
  if (locationApplied) statusMessages.push(`Applied URL location: ${selectionSummary()}`);
  if (statusMessages.length) {
    status.textContent = statusMessages.join(" · ");
  } else if (matchedRule) {
    status.textContent = `Applied format rule: ${matchedRule.pattern} -> ${preset.name || formatPresetName(preset)}`;
  } else if (preset) {
    status.textContent = `Applied format preset: ${preset.name || formatPresetName(preset)}`;
  }
  canvas.focus();
}

function applyFormatOverride(body, override) {
  if (!override) return;
  for (const [key, value] of Object.entries(override)) {
    if (value !== undefined && value !== "") {
      body[key] = value;
    }
  }
}

function hasFormatOverride(override) {
  return Boolean(override && Object.keys(override).length);
}

async function applyURLLocation(location) {
  const range = normalizeURLLocation(location);
  if (!range) return false;
  setSmartSelectionRanges([]);
  selectionAnchor = { row: range.startRow, col: range.startCol };
  selected = { row: range.endRow, col: range.endCol };
  scrollRangeIntoView(range);
  await maybeLoadWindow(selected.row);
  document.body.dataset.urlLocation = selectionSummary();
  draw();
  return true;
}

function normalizeURLLocation(location) {
  if (!location || !documentId || sourceRowCount <= 0 || columns.length === 0) return null;
  const toIndex = (value, fallback) => {
    if (value === "last") return columns.length - 1;
    const parsed = Number.parseInt(String(value ?? ""), 10);
    return Number.isFinite(parsed) && parsed > 0 ? parsed - 1 : fallback;
  };
  const startRow = toIndex(location.startRow, 0);
  const startCol = toIndex(location.startCol, 0);
  const endRow = toIndex(location.endRow, startRow);
  const endCol = toIndex(location.endCol, startCol);
  return normalizeRange({
    startRow: Math.max(0, Math.min(startRow, endRow)),
    endRow: Math.min(sourceRowCount - 1, Math.max(startRow, endRow)),
    startCol: Math.max(0, Math.min(startCol, endCol)),
    endCol: Math.min(columns.length - 1, Math.max(startCol, endCol)),
  });
}

function scrollRangeIntoView(range) {
  scrollY = Math.max(0, Math.min(maxVerticalScroll(), range.startRow * rowHeight - rowHeight * 2));
  const widthBeforeColumn = columns
    .slice(0, range.startCol)
    .reduce((sum, column) => sum + column.width, 0);
  const maxX = maxHorizontalScroll();
  scrollX = Math.max(0, Math.min(maxX, widthBeforeColumn));
}

async function uploadAndOpen(file) {
  captureActiveOpenFileTabState();
  saveCurrentScrollPosition();
  status.textContent = `Opening ${file.name}...`;
  const formData = new FormData();
  formData.append("file", file);
  const response = await fetch("/api/upload-open", {
    method: "POST",
    body: formData,
  });
  const payload = await response.json();
  if (!response.ok) throw new Error(payload.error ?? "upload failed");

  documentId = payload.documentId;
  metadata = payload.document.metadata;
  readOnly = false;
  externalFileChanged = false;
  clearStructureClipboard();
  loadDocument(payload.document);
  pathInput.value = file.name;
  updateFormatRulePlaceholder(file.name);
  updateFileList(file.name, "temporary upload", totalRows, {
    documentId,
    sourceKind: payload.sourceKind || "upload",
  });
  document.body.dataset.openSourceKind = "upload";
  document.body.dataset.openDisplayName = file.name;
  delete document.body.dataset.stdinUrlBoot;
  delete document.body.dataset.newFile;
  draw();
  startFileStatusChecks();
  await checkFileStatus(false);
  canvas.focus();
}

async function createNewFile() {
  captureActiveOpenFileTabState();
  saveCurrentScrollPosition();
  const name = "Untitled.csv";
  status.textContent = t("file.status.creatingNew");
  try {
    const response = await fetch("/api/new", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, rows: 1, columns: 2 }),
    });
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "new file failed");

    documentId = payload.documentId;
    metadata = payload.document.metadata;
    readOnly = false;
    externalFileChanged = false;
    clearStructureClipboard();
    loadDocument(payload.document);
    pathInput.value = name;
    updateFormatRulePlaceholder(name);
    updateFileList(name, "temporary upload", totalRows, {
      documentId,
      sourceKind: payload.sourceKind || "upload",
    });
    document.body.dataset.openSourceKind = "upload";
    document.body.dataset.openDisplayName = name;
    delete document.body.dataset.stdinUrlBoot;
    document.body.dataset.newFile = `${totalRows}:${columns.length}`;
    draw();
    startFileStatusChecks();
    await checkFileStatus(false);
    status.textContent = tf("file.status.newReady", { name });
    canvas.focus();
  } catch (error) {
    status.textContent = tf("file.status.newFailed", { message: error.message });
  }
}

function loadDocument(document) {
  windowGeneration++;
  rows = new Map();
  rowSource = new Map();
  profileData = null;
  lastSQLResult = null;
  profileColumns.replaceChildren(emptyProfileNode(t("profile.status.initialColumns")));
  profileIssues.replaceChildren(emptyProfileNode(t("profile.status.initialIssues")));
  profileSummary.replaceChildren();
  sqlResult.replaceChildren(emptySQLNode(t("sql.empty.prompt")));
  delete window.document.body.dataset.sqlResultPreview;
  delete window.document.body.dataset.sqlResultOpened;
  syncSQLResultActions();
  dataStartRow = document.metadata.dialect.hasHeader ? 1 : 0;
  viewMapped = false;
  sortDescending = false;
  sortKeys = [];
  clearActiveFilter();
  invalidateFindState();
  document.preview.forEach((record, index) => {
    if (index < dataStartRow) return;
    const displayRow = index - dataStartRow;
    rows.set(displayRow, record);
    rowSource.set(displayRow, index);
  });
  windowStart = 0;
  windowEnd = Math.max(0, document.preview.length - dataStartRow);
  sourceRowCount = Math.max(0, (document.rowIndex?.rows || document.preview.length) - dataStartRow);
  aggregateRows = [];
  syncTotalRows();
  allColumns = buildColumns(document);
  hiddenSourceColumns = restoreHiddenColumnsState();
  applyHiddenColumns();
  selected = { row: 0, col: 0 };
  selectionAnchor = { row: 0, col: 0 };
  setSmartSelectionRanges([]);
  scrollX = 0;
  scrollY = 0;
  wheelRemainderY = 0;
  restoreScrollPosition();
  setHeaderInSelection(false);
  clearURLHover();
  delimiterDraftDirty = false;
  encodingDraftDirty = false;
  lineEndingDraftDirty = false;
  dirty = false;
}

function applyMutationPayload(payload, message) {
  const session = payload.session;
  if (!session?.document) {
    status.textContent = message;
    return;
  }
  documentId = session.id ?? documentId;
  readOnly = Boolean(session.readOnly);
  const detail = session.originalPath || session.workingPath || session.document.metadata.path;
  const document = {
    ...session.document,
    metadata: { ...session.document.metadata, path: detail },
  };
  metadata = document.metadata;
  const previousCol = selected.col;
  const previousRow = selected.row;
  const previousAnchor = { ...selectionAnchor };
  loadDocument(document);
  dirty = payload.dirty ?? session.dirty ?? true;
  const maxRow = Math.max(0, totalRows - 1);
  const maxCol = Math.max(0, columns.length - 1);
  selected = {
    row: Math.max(0, Math.min(maxRow, previousRow)),
    col: Math.max(0, Math.min(maxCol, previousCol)),
  };
  selectionAnchor = {
    row: Math.max(0, Math.min(maxRow, previousAnchor.row)),
    col: Math.max(0, Math.min(maxCol, previousAnchor.col)),
  };
  updateFileList(session.displayName || fileNameFromPath(detail), detail, totalRows, {
    documentId,
    sourceKind: session.sourceKind,
  });
  draw();
  status.textContent = message;
  if (dirty) scheduleDirtyDraftSave();
}

function selectedSourceRow() {
  return rowSource.get(selected.row) ?? selected.row + dataStartRow;
}

function selectedRowsForStructure() {
  const displayRows = [];
  const sourceRows = [];
  const seen = new Set();
  for (const range of activeSelectionRanges()) {
    if (range.endRow >= sourceRowCount) {
      status.textContent = "Aggregate rows are view-only.";
      return null;
    }
    for (let row = range.startRow; row <= range.endRow; row++) {
      const sourceRow = rowSource.get(row);
      if (viewMapped && sourceRow == null) {
        status.textContent = "Load the selected filtered rows before using row structure commands.";
        return null;
      }
      const recordRow = sourceRow ?? row + dataStartRow;
      if (seen.has(recordRow)) continue;
      seen.add(recordRow);
      displayRows.push(row);
      sourceRows.push(recordRow);
    }
  }
  return displayRows.length ? { displayRows, sourceRows } : null;
}

function selectedRowPasteAfterTarget() {
  const rowSelection = selectedRowsForStructure();
  if (!rowSelection) return null;
  return Math.max(...rowSelection.sourceRows) + 1;
}

function selectedColumnsForStructure() {
  const sourceCols = [];
  const seen = new Set();
  for (const range of activeSelectionRanges()) {
    for (let col = range.startCol; col <= range.endCol; col++) {
      const sourceCol = sourceColFor(col);
      if (seen.has(sourceCol)) continue;
      seen.add(sourceCol);
      sourceCols.push(sourceCol);
    }
  }
  return sourceCols.length ? sourceCols : null;
}

function selectedColumnPasteRightTarget() {
  const sourceCols = selectedColumnsForStructure();
  if (!sourceCols) return null;
  return Math.min(Math.max(...sourceCols) + 1, Math.max(...allColumns.map((column) => column.sourceCol), -1) + 1);
}

function selectedMatrix() {
  const matrix = [];
  for (const range of activeSelectionRanges()) {
    for (let row = range.startRow; row <= range.endRow; row++) {
      const values = [];
      for (let col = range.startCol; col <= range.endCol; col++) {
        values.push(valueAt(row, col));
      }
      matrix.push(values);
    }
  }
  return matrix;
}

function selectedHeaders() {
  const range = activeSelectionRanges()[0] ?? normalizedSelection();
  const headers = [];
  for (let col = range.startCol; col <= range.endCol; col++) {
    headers.push(columns[col]?.name ?? `Column ${col + 1}`);
  }
  return headers;
}

function selectedCellsForTransform(maxCells = 100_000) {
  const count = selectionCellCount();
  if (count === 0) {
    status.textContent = "Select at least one cell.";
    return null;
  }
  if (count > maxCells) {
    status.textContent = `Selection has ${formatCount(count)} cells; transform limit is ${formatCount(maxCells)}.`;
    return null;
  }
  const cells = [];
  for (const range of activeSelectionRanges()) {
    if (range.endRow >= sourceRowCount) {
      status.textContent = "Aggregate rows are view-only.";
      return null;
    }
    for (let row = range.startRow; row <= range.endRow; row++) {
      const sourceRow = rowSource.get(row);
      if (viewMapped && sourceRow == null) {
        status.textContent = "Load the selected filtered rows before transforming them.";
        return null;
      }
      const recordRow = sourceRow ?? row + dataStartRow;
      for (let col = range.startCol; col <= range.endCol; col++) {
        cells.push({ row: recordRow, col: sourceColFor(col) });
      }
    }
  }
  return cells;
}

function toggleHeaderInSelection() {
  if (!documentId) return;
  setHeaderInSelection(!includeHeadersInSelection);
  draw();
  status.textContent = includeHeadersInSelection
    ? "Header included when copying selection"
    : "Header excluded from selection copy";
}

function setHeaderInSelection(value) {
  includeHeadersInSelection = Boolean(value);
  copyHeadersToggle.checked = includeHeadersInSelection;
  document.body.dataset.headerInSelection = includeHeadersInSelection ? "true" : "false";
}

function setSmartSelectionRanges(ranges) {
  smartSelectionRanges = dedupeRanges(ranges);
  document.body.dataset.smartSelectionRanges = smartSelectionRanges.length ? smartSelectionRanges.map(rangeKey).join(";") : "";
}

function smartSelectionGestureForEvent(event) {
  if (!(event.metaKey || event.ctrlKey)) return "";
  return event.shiftKey ? "subtract" : "add";
}

function startSmartSelectionGesture(action) {
  smartSelectionGesture = action;
  smartSelectionGestureBaseRanges = activeSelectionRanges();
  setSmartSelectionRanges(smartSelectionGestureBaseRanges);
  document.body.dataset.smartSelectionGesture = `${action}:start`;
}

function finishSmartSelectionGesture() {
  if (!smartSelectionGesture) return false;
  const action = smartSelectionGesture;
  const range = normalizedSelection();
  const base = dedupeRanges(smartSelectionGestureBaseRanges);
  smartSelectionGesture = "";
  smartSelectionGestureBaseRanges = [];

  if (action === "subtract") {
    const next = [];
    for (const baseRange of base) next.push(...subtractRange(baseRange, range));
    setSmartSelectionRanges(next);
    focusFirstSmartSelectionCell();
    document.body.dataset.smartSelectionGesture = `subtract:${rangeKey(range)}`;
    draw();
    status.textContent = smartSelectionRanges.length
      ? `Subtracted range by gesture: ${formatCount(selectionCellCount())} cells selected`
      : "Selection range cleared by gesture";
    return true;
  }

  setSmartSelectionRanges([...base, range]);
  document.body.dataset.smartSelectionGesture = `add:${rangeKey(range)}`;
  draw();
  status.textContent = `Added range by gesture: ${formatCount(selectionCellCount())} cells selected`;
  return true;
}

function addSmartSelectionRange() {
  if (!documentId || columns.length === 0 || totalRows === 0) return;
  setSmartSelectionRanges([...smartSelectionRanges, normalizedSelection()]);
  draw();
  status.textContent = `Added range to selection: ${formatCount(selectionCellCount())} cells selected`;
}

function subtractSmartSelectionRange() {
  if (!documentId || columns.length === 0 || totalRows === 0) return;
  const cut = normalizedSelection();
  const base = smartSelectionRanges.length ? normalizedSmartSelectionRanges() : [cut];
  const next = [];
  for (const range of base) next.push(...subtractRange(range, cut));
  setSmartSelectionRanges(next);
  focusFirstSmartSelectionCell();
  draw();
  status.textContent = smartSelectionRanges.length
    ? `Subtracted range: ${formatCount(selectionCellCount())} cells selected`
    : "Selection range cleared";
}

function invertSmartSelection() {
  if (!documentId || columns.length === 0 || totalRows === 0) return;
  let inverted = [{ startRow: 0, endRow: Math.max(0, sourceRowCount - 1), startCol: 0, endCol: columns.length - 1 }];
  for (const cut of activeSelectionRanges()) {
    inverted = inverted.flatMap((range) => subtractRange(range, cut));
  }
  setSmartSelectionRanges(inverted);
  focusFirstSmartSelectionCell();
  draw();
  status.textContent = smartSelectionRanges.length
    ? `Inverted selection: ${formatCount(selectionCellCount())} cells selected`
    : "Inverted selection is empty";
}

function clearSmartSelectionRanges() {
  if (!smartSelectionRanges.length) return false;
  setSmartSelectionRanges([]);
  draw();
  status.textContent = "Extra selection ranges cleared";
  return true;
}

function focusFirstSmartSelectionCell() {
  const first = smartSelectionRanges[0];
  if (!first) {
    selectionAnchor = { ...selected };
    return;
  }
  selected = { row: first.startRow, col: first.startCol };
  selectionAnchor = { ...selected };
}

function subtractRange(range, cut) {
  const base = normalizeRange(range);
  const eraser = normalizeRange(cut);
  if (!base) return [];
  if (!eraser || !rangesOverlap(base, eraser)) return [base];

  const top = { startRow: base.startRow, endRow: eraser.startRow - 1, startCol: base.startCol, endCol: base.endCol };
  const bottom = { startRow: eraser.endRow + 1, endRow: base.endRow, startCol: base.startCol, endCol: base.endCol };
  const middleStart = Math.max(base.startRow, eraser.startRow);
  const middleEnd = Math.min(base.endRow, eraser.endRow);
  const left = { startRow: middleStart, endRow: middleEnd, startCol: base.startCol, endCol: eraser.startCol - 1 };
  const right = { startRow: middleStart, endRow: middleEnd, startCol: eraser.endCol + 1, endCol: base.endCol };
  return [top, bottom, left, right].map(normalizeRange).filter(Boolean);
}

function rangesOverlap(a, b) {
  return a.startRow <= b.endRow && a.endRow >= b.startRow && a.startCol <= b.endCol && a.endCol >= b.startCol;
}

function rangeContainsRange(container, candidate) {
  if (!container || !candidate) return false;
  return (
    container.startRow <= candidate.startRow &&
    container.endRow >= candidate.endRow &&
    container.startCol <= candidate.startCol &&
    container.endCol >= candidate.endCol
  );
}

function rangeContainsCell(range, row, col) {
  return Boolean(range && row >= range.startRow && row <= range.endRow && col >= range.startCol && col <= range.endCol);
}

async function copyUsingSelectedFormat() {
  const format = copyFormatSelect.value || "tsv";
  setHeaderInSelection(copyHeadersToggle.checked);
  document.body.dataset.copyFormat = `${format}:${includeHeadersInSelection ? "headers" : "body"}`;
  await copySelection(format, { includeHeaders: includeHeadersInSelection });
}

async function copySelection(format = "tsv", options = {}) {
  if (!documentId) return;
  const matrix = selectedMatrix();
  const includeHeaders = options.includeHeaders ?? includeHeadersInSelection;
  const copyMatrix = includeHeaders
    ? [selectedHeaders(), ...matrix]
    : matrix;
  const text = formatMatrix(copyMatrix, format);
  if (!text) {
    status.textContent = "No cells selected to copy.";
    return;
  }
  window.__opencsvLastCopiedText = text;
  document.body.dataset.copyPreview = text.slice(0, 500);
  try {
    await writeClipboardText(text);
    const headerLabel = includeHeaders ? " with headers" : "";
    status.textContent = `Copied ${formatCount(selectionCellCount())} cells${headerLabel} as ${format.toUpperCase()}`;
  } catch (error) {
    status.textContent = `Copy failed: ${error.message}`;
  }
}

async function copyActiveRow(options = {}) {
  if (!documentId) return;
  const format = normalizeCopyFormat(options.format || "tsv");
  if (selected.row >= sourceRowCount) {
    await copySelection(format, { includeHeaders: Boolean(options.includeHeaders) });
    return;
  }
  const includeHeaders = Boolean(options.includeHeaders);
  const rowSelection = selectedRowsForStructure();
  if (!rowSelection) return;
  const rowCount = rowSelection.sourceRows.length;
  const copied = await copyExportViewSlice({
    label: includeHeaders
      ? rowCount === 1 ? "row with headers" : "rows with headers"
      : rowCount === 1 ? "row" : "rows",
    body: {
      documentId,
      rows: rowSelection.displayRows,
      columns: columns.map((column) => column.sourceCol),
      format,
      includeHeaders,
    },
    cellCount: columns.length * rowCount,
    format,
    onCopyReady: () => {
      if (!includeHeaders && format === "tsv") setStructureClipboard("row", rowSelection.sourceRows, false);
    },
  });
  if (!copied && !includeHeaders && format === "tsv") clearStructureClipboard();
}

async function copyActiveColumn(options = {}) {
  if (!documentId || !columns[selected.col]) return;
  const format = normalizeCopyFormat(options.format || "tsv");
  const includeHeaders = Boolean(options.includeHeaders);
  const sourceCols = selectedColumnsForStructure();
  if (!sourceCols) return;
  const colCount = sourceCols.length;
  const copied = await copyExportViewSlice({
    label: includeHeaders
      ? colCount === 1 ? "column with header" : "columns with headers"
      : colCount === 1 ? "column" : "columns",
    body: {
      documentId,
      columns: sourceCols,
      format,
      includeHeaders,
    },
    cellCount: sourceRowCount * colCount,
    format,
    onCopyReady: () => {
      if (!includeHeaders && format === "tsv") setStructureClipboard("column", sourceCols, false);
    },
  });
  if (!copied && !includeHeaders && format === "tsv") clearStructureClipboard();
}

function normalizeCopyFormat(format) {
  return ["csv", "tsv", "markdown", "json", "html", "sql", "latex", "raw"].includes(format) ? format : "tsv";
}

function copyFormatLabel(format) {
  return {
    csv: "CSV",
    tsv: "TSV",
    markdown: "Markdown",
    json: "JSON",
    html: "HTML",
    sql: "SQL",
    latex: "LaTeX",
    raw: "Raw",
  }[normalizeCopyFormat(format)];
}

async function copyExportViewSlice({ label, body, cellCount, format = "tsv", onCopyReady = null }) {
  format = normalizeCopyFormat(format);
  const response = await fetch("/api/export-view", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const blob = await response.blob();
  const text = await blob.text();
  if (!response.ok) {
    status.textContent = `Copy ${label} failed: ${errorFromResponseText(text)}`;
    return false;
  }
  if (!text) {
    status.textContent = `No ${label} data to copy.`;
    return false;
  }
  window.__opencsvLastCopiedText = text;
  document.body.dataset.copyPreview = text.slice(0, 500);
  document.body.dataset.copySlice = label;
  if (onCopyReady) onCopyReady(text);
  try {
    await writeClipboardText(text);
    const rows = Number(response.headers.get("X-Opencsv-Rows") || 0);
    const cols = Number(response.headers.get("X-Opencsv-Columns") || 0);
    const shape = rows && cols ? `${formatCount(rows)} rows, ${formatCount(cols)} columns` : `${formatCount(cellCount)} cells`;
    status.textContent = `Copied ${label}: ${shape} as ${copyFormatLabel(format)}`;
  } catch (error) {
    status.textContent = `Copy ${label} failed: ${error.message}`;
  }
  return true;
}

async function exportCurrentView(format = "") {
  if (!documentId) return;
  format = format || exportFormatSelect?.value || "csv";
  document.body.dataset.exportFormat = format;
  const visibleColumns = columns.map((column) => column.sourceCol);
  const response = await fetch("/api/export-view", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, columns: visibleColumns, format }),
  });
  const blob = await response.blob();
  if (!response.ok) {
    status.textContent = `Export failed: ${errorFromResponseText(await blob.text())}`;
    return;
  }

  const filename =
    filenameFromDisposition(response.headers.get("Content-Disposition")) ||
    response.headers.get("X-Opencsv-Filename") ||
    "opencsv.view.csv";
  const contentType = response.headers.get("Content-Type") || "text/csv;charset=utf-8";
  if (format === "xlsx") {
    const bytes = new Uint8Array(await blob.arrayBuffer());
    const signature = Array.from(bytes.slice(0, 2)).map((byte) => String.fromCharCode(byte)).join("");
    window.__opencsvLastExportedBytes = bytes.length;
    document.body.dataset.exportPreview = `${signature}:${bytes.length}`;
    triggerDownloadBlob(new Blob([bytes], { type: contentType }), filename);
  } else {
    const text = await blob.text();
    window.__opencsvLastExportedText = text;
    document.body.dataset.exportPreview = text.slice(0, 500);
    triggerDownloadBlob(blob, filename);
  }
  const rows = Number(response.headers.get("X-Opencsv-Rows") || 0);
  const cols = Number(response.headers.get("X-Opencsv-Columns") || visibleColumns.length);
  status.textContent = `Exported ${formatCount(rows)} rows, ${formatCount(cols)} columns to ${filename}`;
}

async function copyCurrentSQLResult() {
  return copySQLResult(sqlResultFormatSelect?.value || "csv");
}

async function copySQLResult(format = "csv") {
  format = normalizeSQLResultFormat(format);
  const selectedMatrix = sqlSelectedResultMatrix(lastSQLResult);
  const matrix = selectedMatrix ?? sqlResultMatrix(lastSQLResult);
  if (!matrix.length) {
    status.textContent = t("sql.status.runBeforeCopy");
    return;
  }
  return copySQLMatrix(matrix, selectedMatrix ? t("sql.label.selection") : t("sql.label.result"), {
    format,
    bodyRows: selectedMatrix ? matrix.length : Math.max(0, matrix.length - 1),
  });
}

async function copySQLSelection(format = "csv") {
  const matrix = sqlSelectedResultMatrix(lastSQLResult) ?? [];
  return copySQLMatrix(matrix, t("sql.label.selection"), { format });
}

async function copySQLFullResult(format = "csv") {
  const matrix = sqlResultMatrix(lastSQLResult);
  return copySQLMatrix(matrix, t("sql.label.result"), {
    format,
    bodyRows: Math.max(0, matrix.length - 1),
  });
}

async function copySQLContextCell() {
  const target = sqlResultTargetForCommand();
  if (!target) {
    status.textContent = t("sql.status.selectCell");
    return false;
  }
  return copySQLMatrix(sqlResultCellMatrix(target.row, target.col), t("sql.label.cell"));
}

async function copySQLContextRow(includeHeaders = false) {
  const target = sqlResultTargetForCommand();
  if (!target) {
    status.textContent = t("sql.status.selectRow");
    return false;
  }
  return copySQLMatrix(
    sqlResultRowMatrix(target.row, includeHeaders),
    includeHeaders ? t("sql.label.rowWithHeaders") : t("sql.label.row"),
    includeHeaders ? { bodyRows: 1 } : {},
  );
}

async function copySQLContextColumn(includeHeader = false) {
  const target = sqlResultTargetForCommand();
  if (!target) {
    status.textContent = t("sql.status.selectColumn");
    return false;
  }
  return copySQLMatrix(
    sqlResultColumnMatrix(target.col, includeHeader),
    includeHeader ? t("sql.label.columnWithHeader") : t("sql.label.column"),
    includeHeader ? { bodyRows: lastSQLResult?.rows?.length ?? 0 } : {},
  );
}

async function copySQLMatrix(matrix, label, options = {}) {
  if (!matrix.length) {
    status.textContent = tf("sql.status.noCopy", { label });
    return false;
  }
  const format = normalizeSQLResultFormat(options.format || "csv");
  const text = formatMatrix(matrix, format);
  const formatLabel = sqlResultFormatLabel(format);
  window.__opencsvLastSQLResultText = text;
  document.body.dataset.sqlResultPreview = text.slice(0, 500);
  document.body.dataset.sqlResultCopy = `${label}:${format}`;
  document.body.dataset.sqlResultFormat = format;
  try {
    await writeClipboardText(text);
    const rows = options.bodyRows ?? matrix.length;
    status.textContent = tf("sql.status.copied", {
      label,
      rows: formatCount(rows),
      cols: formatCount(matrix[0]?.length ?? 0),
      format: formatLabel,
    });
    return true;
  } catch (error) {
    status.textContent = tf("sql.status.copyFailed", { label, message: error.message });
    return false;
  }
}

async function exportCurrentSQLResult() {
  return exportSQLResult(sqlResultFormatSelect?.value || "csv");
}

async function exportSQLResult(format = "csv") {
  format = normalizeSQLResultFormat(format);
  const matrix = sqlResultMatrix(lastSQLResult);
  if (!matrix.length) {
    status.textContent = t("sql.status.runBeforeExport");
    return;
  }
  const text = formatMatrix(matrix, format);
  const filename = sqlResultFilename(format);
  window.__opencsvLastSQLResultText = text;
  document.body.dataset.sqlResultPreview = text.slice(0, 500);
  document.body.dataset.sqlResultFormat = format;
  triggerDownload(text, filename, sqlResultFormatContentTypes[format] || "text/plain;charset=utf-8");
  status.textContent = tf("sql.status.exported", {
    format: sqlResultFormatLabel(format),
    rows: formatCount(Math.max(0, matrix.length - 1)),
    cols: formatCount(matrix[0]?.length ?? 0),
    filename,
  });
}

async function openSQLResultAsDocument() {
  const result = lastSQLResult;
  const columns = sqlResultColumns(result).map((column, index) => String(column || `Column ${index + 1}`));
  if (!result || !columns.length) {
    status.textContent = t("sql.status.runBeforeOpen");
    return;
  }
  const rows = (result.rows ?? []).map((row) =>
    columns.map((_, index) => {
      const value = row?.[index] ?? "";
      return typeof value === "string" ? value : String(value);
    }),
  );
  const name = sqlResultFilename("csv");
  const bodyRows = rows.length;
  const colCount = columns.length;
  captureActiveOpenFileTabState();
  saveCurrentScrollPosition();
  status.textContent = t("sql.status.openingResult");
  try {
    const response = await fetch("/api/sql-result/open", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, columns, rows }),
    });
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "open SQL result failed");

    documentId = payload.documentId;
    metadata = payload.document.metadata;
    readOnly = false;
    externalFileChanged = false;
    clearStructureClipboard();
    loadDocument(payload.document);
    const displayName = payload.displayName || name;
    const sourceKind = payload.sourceKind || "sql-result";
    pathInput.value = displayName;
    updateFormatRulePlaceholder(displayName);
    updateFileList(displayName, "temporary upload", totalRows, {
      documentId,
      sourceKind,
    });
    document.body.dataset.openSourceKind = sourceKind;
    document.body.dataset.openDisplayName = displayName;
    document.body.dataset.sqlResultOpened = `${bodyRows}:${colCount}:${displayName}`;
    delete document.body.dataset.stdinUrlBoot;
    delete document.body.dataset.newFile;
    setViewMode("grid");
    draw();
    startFileStatusChecks();
    await checkFileStatus(false);
    status.textContent = tf("sql.status.openedResult", {
      name: displayName,
      rows: formatCount(bodyRows),
      cols: formatCount(colCount),
    });
    canvas.focus();
  } catch (error) {
    status.textContent = tf("sql.status.openResultFailed", { message: error.message });
  }
}

function sqlResultMatrix(result) {
  const sqlColumns = sqlResultColumns(result);
  if (!result || (!sqlColumns.length && !result.rows?.length)) return [];
  const rows = (result.rows ?? []).map((row) => {
    const next = Array.from(row ?? []);
    while (next.length < sqlColumns.length) next.push("");
    return next;
  });
  return [sqlColumns, ...rows];
}

function sqlSelectedResultMatrix(result) {
  const range = normalizedSQLResultSelection();
  if (!result || !range) return null;
  const resultRows = result.rows ?? [];
  const matrix = [];
  for (let row = range.startRow; row <= range.endRow; row++) {
    const source = resultRows[row] ?? [];
    const output = [];
    for (let col = range.startCol; col <= range.endCol; col++) {
      output.push(source[col] ?? "");
    }
    matrix.push(output);
  }
  return matrix.length ? matrix : null;
}

function sqlResultCellMatrix(row, col) {
  const rows = lastSQLResult?.rows ?? [];
  if (row < 0 || col < 0 || row >= rows.length) return [];
  return [[rows[row]?.[col] ?? ""]];
}

function sqlResultRowMatrix(row, includeHeaders = false) {
  const resultColumns = sqlResultColumns(lastSQLResult);
  const rows = lastSQLResult?.rows ?? [];
  if (row < 0 || row >= rows.length) return [];
  const output = [];
  if (includeHeaders) output.push(resultColumns);
  output.push(resultColumns.map((_, col) => rows[row]?.[col] ?? ""));
  return output;
}

function sqlResultColumnMatrix(col, includeHeader = false) {
  const resultColumns = sqlResultColumns(lastSQLResult);
  const rows = lastSQLResult?.rows ?? [];
  if (col < 0 || col >= resultColumns.length) return [];
  const output = [];
  if (includeHeader) output.push([resultColumns[col]]);
  for (const row of rows) {
    output.push([row?.[col] ?? ""]);
  }
  return output;
}

function sqlResultColumns(result) {
  if (!result) return [];
  const columns = Array.isArray(result.columns) ? result.columns : [];
  if (columns.length) return columns.map((column, index) => column || `Column ${index + 1}`);
  const width = Math.max(0, ...(result.rows ?? []).map((row) => row?.length ?? 0));
  return Array.from({ length: width }, (_, index) => `Column ${index + 1}`);
}

function normalizeSQLResultFormat(format) {
  return sqlResultFormats.includes(format) ? format : "csv";
}

function sqlResultFormatLabel(format) {
  return sqlResultFormatLabels[normalizeSQLResultFormat(format)] || "CSV";
}

function sqlResultFormatCommandSuffix(format) {
  switch (normalizeSQLResultFormat(format)) {
    case "tsv":
      return "Tsv";
    case "markdown":
      return "Markdown";
    case "json":
      return "Json";
    case "html":
      return "Html";
    case "sql":
      return "Sql";
    case "latex":
      return "Latex";
    case "raw":
      return "Raw";
    default:
      return "Csv";
  }
}

function sqlResultFilename(format = "csv") {
  format = normalizeSQLResultFormat(format);
  const name = activeFileEntry?.name || fileNameFromPath(metadata?.path ?? "") || "opencsv";
  const base = name.replace(/\.[^.]+$/, "") || "opencsv";
  const extension = sqlResultFormatExtensions[format] || "txt";
  return `${base}.sql-result.${extension}`;
}

function syncSQLResultActions() {
  const available = sqlResultAvailable();
  const format = normalizeSQLResultFormat(sqlResultFormatSelect?.value || "csv");
  if (sqlResultFormatSelect) {
    sqlResultFormatSelect.disabled = !available;
    document.body.dataset.sqlResultFormat = format;
  }
  openSqlResultButton.textContent = t("sql.openResult");
  copySqlResultButton.textContent = t("sql.copy");
  exportSqlResultButton.textContent = t("sql.export");
  openSqlResultButton.disabled = !available;
  copySqlResultButton.disabled = !available;
  exportSqlResultButton.disabled = !available;
  updateCommandState();
}

function clearSQLResultSelection() {
  sqlResultSelectionAnchor = null;
  sqlResultSelectionRange = null;
  sqlResultContextTarget = null;
  document.body.dataset.sqlResultSelection = "";
  syncSQLResultSelection();
}

function normalizedSQLResultSelection() {
  if (!sqlResultSelectionRange) return null;
  return {
    startRow: Math.min(sqlResultSelectionRange.startRow, sqlResultSelectionRange.endRow),
    endRow: Math.max(sqlResultSelectionRange.startRow, sqlResultSelectionRange.endRow),
    startCol: Math.min(sqlResultSelectionRange.startCol, sqlResultSelectionRange.endCol),
    endCol: Math.max(sqlResultSelectionRange.startCol, sqlResultSelectionRange.endCol),
  };
}

function syncSQLResultSelection() {
  const range = normalizedSQLResultSelection();
  sqlResult.querySelectorAll("td[data-sql-row]").forEach((cell) => {
    const row = Number(cell.dataset.sqlRow);
    const col = Number(cell.dataset.sqlCol);
    const selectedCell = Boolean(
      range &&
      row >= range.startRow &&
      row <= range.endRow &&
      col >= range.startCol &&
      col <= range.endCol,
    );
    cell.classList.toggle("sql-selected", selectedCell);
    cell.setAttribute("aria-selected", String(selectedCell));
  });
  if (range) {
    document.body.dataset.sqlResultSelection = `${range.startRow}:${range.startCol}-${range.endRow}:${range.endCol}`;
  } else {
    document.body.dataset.sqlResultSelection = "";
  }
}

function selectSQLResultCell(row, col, extend = false) {
  const resultRows = lastSQLResult?.rows ?? [];
  const resultColumns = sqlResultColumns(lastSQLResult);
  if (row < 0 || col < 0 || row >= resultRows.length || col >= resultColumns.length) return;
  if (!extend || !sqlResultSelectionAnchor) {
    sqlResultSelectionAnchor = { row, col };
  }
  sqlResultSelectionRange = {
    startRow: sqlResultSelectionAnchor.row,
    startCol: sqlResultSelectionAnchor.col,
    endRow: row,
    endCol: col,
  };
  sqlResultContextTarget = { row, col };
  syncSQLResultSelection();
  const range = normalizedSQLResultSelection();
  const rowsSelected = range ? range.endRow - range.startRow + 1 : 0;
  const colsSelected = range ? range.endCol - range.startCol + 1 : 0;
  status.textContent = tf("sql.status.selected", {
    rows: formatCount(rowsSelected),
    cols: formatCount(colsSelected),
  });
}

function sqlResultTargetForCommand() {
  const resultRows = lastSQLResult?.rows ?? [];
  const resultColumns = sqlResultColumns(lastSQLResult);
  const range = normalizedSQLResultSelection();
  const target = sqlResultContextTarget ?? (range ? { row: range.startRow, col: range.startCol } : null);
  if (
    !target ||
    target.row < 0 ||
    target.col < 0 ||
    target.row >= resultRows.length ||
    target.col >= resultColumns.length
  ) {
    return null;
  }
  return target;
}

function sqlResultCellFromEvent(event) {
  const cell = event.target?.closest?.("td[data-sql-row]");
  if (!cell || !sqlResult.contains(cell)) return null;
  return {
    row: Number(cell.dataset.sqlRow),
    col: Number(cell.dataset.sqlCol),
  };
}

function sqlSelectionContains(row, col) {
  const range = normalizedSQLResultSelection();
  return Boolean(
    range &&
    row >= range.startRow &&
    row <= range.endRow &&
    col >= range.startCol &&
    col <= range.endCol,
  );
}

function renderSQLResultContextMenu(target, clientX, clientY) {
  if (!lastSQLResult) return;
  sqlResultContextTarget = { row: target.row, col: target.col };
  if (!sqlSelectionContains(target.row, target.col)) {
    selectSQLResultCell(target.row, target.col, false);
    sqlResultContextTarget = { row: target.row, col: target.col };
  }
  const resultColumns = sqlResultColumns(lastSQLResult);
  const columnName = resultColumns[target.col] || `Column ${target.col + 1}`;
  const heading = document.createElement("strong");
  heading.textContent = `SQL ${target.row + 1}, ${columnName}`;
  gridContextMenu.replaceChildren(heading);

  const addSeparator = () => gridContextMenu.append(document.createElement("hr"));
  const addItem = (label, action, options = {}) => {
    const button = document.createElement("button");
    button.type = "button";
    button.setAttribute("role", "menuitem");
    button.textContent = label;
    if (options.commandID) button.dataset.command = options.commandID;
    button.disabled = Boolean(options.disabled);
    button.addEventListener("click", async () => {
      closeGridContextMenu();
      await action();
    });
    gridContextMenu.append(button);
    return button;
  };
  const contextCommandIDs = [];
  const addCommand = (commandID, label = "", options = {}) => {
    const command = commandsById.get(commandID);
    contextCommandIDs.push(commandID);
    addItem(label || commandTitle(command || { id: commandID }), () => executeCommand(commandID), {
      commandID,
      disabled: options.disabled || !command || isCommandDisabled(command),
    });
  };

  const selectedMatrix = sqlSelectedResultMatrix(lastSQLResult);
  addCommand("sql.copyCell", t("sql.context.copyCell"));
  addCommand("sql.copyRow", t("sql.context.copyRow"));
  addCommand("sql.copyRowWithHeaders", t("sql.context.copyRowWithHeaders"));
  addCommand("sql.copyColumn", t("sql.context.copyColumn"));
  addCommand("sql.copyColumnWithHeader", t("sql.context.copyColumnWithHeader"));
  addSeparator();
  for (const format of sqlResultFormats) {
    const formatLabel = sqlResultFormatLabel(format);
    addCommand(`sql.copySelection${sqlResultFormatCommandSuffix(format)}`, tf("sql.context.copySelectionAs", { format: formatLabel }), {
      disabled: !selectedMatrix,
    });
  }
  addSeparator();
  for (const format of sqlResultFormats) {
    const formatLabel = sqlResultFormatLabel(format);
    addCommand(`sql.copyFullResult${sqlResultFormatCommandSuffix(format)}`, tf("sql.context.copyResultAs", { format: formatLabel }));
  }
  addSeparator();
  for (const format of sqlResultFormats) {
    const formatLabel = sqlResultFormatLabel(format);
    addCommand(`sql.exportResult${sqlResultFormatCommandSuffix(format)}`, tf("sql.context.exportResultAs", { format: formatLabel }));
  }
  addSeparator();
  addCommand("sql.openResult", t("sql.context.openResult"));
  addCommand("sql.clearSelection", t("sql.context.clearSelection"), { disabled: !selectedMatrix });

  gridContextMenu.hidden = false;
  const menuRect = gridContextMenu.getBoundingClientRect();
  const left = Math.min(Math.max(8, clientX), window.innerWidth - menuRect.width - 8);
  const top = Math.min(Math.max(8, clientY), window.innerHeight - menuRect.height - 8);
  gridContextMenu.style.left = `${left}px`;
  gridContextMenu.style.top = `${top}px`;
  gridContextMenu.dataset.target = "sql-result";
  document.body.dataset.gridContextTarget = "sql-result";
  document.body.dataset.sqlResultContext = `${target.row}:${target.col}`;
  document.body.dataset.sqlResultContextCommands = contextCommandIDs.join(",");
  document.body.dataset.gridContextCommands = contextCommandIDs.join(",");
}

function errorFromResponseText(text) {
  try {
    return JSON.parse(text).error ?? "unknown error";
  } catch {
    return text.trim() || "unknown error";
  }
}

function filenameFromDisposition(disposition) {
  if (!disposition) return "";
  const match = disposition.match(/filename="([^"]+)"/i);
  return match ? match[1] : "";
}

function triggerDownload(text, filename, contentType) {
  const blob = new Blob([text], { type: contentType });
  triggerDownloadBlob(blob, filename);
}

function triggerDownloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = filename;
  anchor.hidden = true;
  document.body.append(anchor);
  anchor.click();
  anchor.remove();
  setTimeout(() => URL.revokeObjectURL(url), 0);
}

async function writeClipboardText(text) {
  try {
    await navigator.clipboard.writeText(text);
    return;
  } catch (error) {
    if (!fallbackCopyText(text)) throw error;
  }
}

function fallbackCopyText(text) {
  const area = document.createElement("textarea");
  area.value = text;
  area.setAttribute("readonly", "");
  area.style.position = "fixed";
  area.style.left = "-9999px";
  area.style.top = "0";
  document.body.append(area);
  area.focus();
  area.select();
  let ok = false;
  try {
    ok = document.execCommand("copy");
  } finally {
    area.remove();
    canvas.focus();
  }
  return ok;
}

async function pasteClipboard() {
  if (!documentId) return;
  if (!ensureEditable("Paste")) return;
  try {
    if (navigator.clipboard?.read) {
      try {
        const items = await navigator.clipboard.read();
        const html = await clipboardItemsText(items, "text/html");
        if (html && pasteHTMLTable(html)) return;
        const plain = await clipboardItemsText(items, "text/plain");
        if (plain) {
          pasteText(plain);
          return;
        }
      } catch {
        // Fall back to plain text read below when rich clipboard access is unavailable.
      }
    }
    const text = await navigator.clipboard.readText();
    pasteText(text);
  } catch (error) {
    status.textContent = `Paste failed: ${error.message}`;
  }
}

async function clipboardItemsText(items, mimeType) {
  for (const item of items ?? []) {
    if (!Array.isArray(item.types) || !item.types.includes(mimeType)) continue;
    const blob = await item.getType(mimeType);
    return await blob.text();
  }
  return "";
}

function pasteHTMLTable(htmlText) {
  if (!ensureEditable("Paste")) return false;
  const matrix = parseHTMLTable(htmlText);
  if (!matrix.length) return false;
  pasteParsedMatrix(matrix, "html");
  return true;
}

function pasteText(text) {
  if (!ensureEditable("Paste")) return;
  const matrix = parseTSV(text);
  if (!matrix.length) {
    status.textContent = "Clipboard is empty.";
    return;
  }
  pasteParsedMatrix(matrix, "tsv");
}

function pasteParsedMatrix(matrix, source) {
  window.__opencsvLastPasteMatrix = matrix;
  document.body.dataset.pasteSource = source;
  document.body.dataset.pasteMatrix = `${matrix.length}:${Math.max(...matrix.map((row) => row.length), 0)}`;
  document.body.dataset.pastePreview = JSON.stringify(matrix).slice(0, 500);
  pasteMatrix(matrix);
}

function pasteMatrix(matrix) {
  if (!ensureEditable("Paste")) return;
  if (selected.row >= sourceRowCount) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  windowGeneration++;
  const edits = [];
  const startRow = selected.row;
  const startCol = selected.col;
  for (let rowOffset = 0; rowOffset < matrix.length; rowOffset++) {
    const row = startRow + rowOffset;
    if (row >= sourceRowCount) break;
    const record = rows.get(row) ?? [];
    const values = matrix[rowOffset] ?? [];
    for (let colOffset = 0; colOffset < values.length; colOffset++) {
      const col = startCol + colOffset;
      if (col >= columns.length) break;
      const sourceCol = sourceColFor(col);
      while (record.length <= sourceCol) record.push("");
      record[sourceCol] = values[colOffset];
      edits.push({ row: rowSource.get(row) ?? row + dataStartRow, col: sourceCol, value: values[colOffset] });
    }
    rows.set(row, record);
  }
  selected = {
    row: Math.min(totalRows - 1, startRow + matrix.length - 1),
    col: Math.min(columns.length - 1, startCol + Math.max(...matrix.map((row) => row.length), 1) - 1),
  };
  selectionAnchor = { row: startRow, col: startCol };
  autoFitColumnsForEdits(edits);
  syncCellValueAfterEdit();
  submitCellEdits(edits, startRow, `Pasted ${formatCount(edits.length)} cells`);
}

function formatMatrix(matrix, format) {
  if (format === "csv") return matrixToCSV(matrix);
  if (format === "markdown") return matrixToMarkdown(matrix);
  if (format === "json") return JSON.stringify(matrix, null, 2);
  if (format === "html") return matrixToHTML(matrix);
  if (format === "sql") return matrixToSQLValues(matrix);
  if (format === "latex") return matrixToLaTeX(matrix);
  return matrix.map((row) => row.join("\t")).join("\n");
}

function matrixToCSV(matrix) {
  return matrix
    .map((row) => row.map(csvEscape).join(","))
    .join("\n");
}

function csvEscape(value) {
  const text = String(value ?? "");
  if (!/[",\r\n]/.test(text)) return text;
  return `"${text.replaceAll('"', '""')}"`;
}

function matrixToMarkdown(matrix) {
  if (!matrix.length) return "";
  const lines = [];
  lines.push(markdownRow(matrix[0]));
  lines.push(`|${matrix[0].map(() => " --- ").join("|")}|`);
  for (const row of matrix.slice(1)) lines.push(markdownRow(row));
  return lines.join("\n");
}

function markdownRow(row) {
  return `| ${row.map((value) => String(value ?? "").replaceAll("|", "\\|").replaceAll("\n", "<br>")).join(" | ")} |`;
}

function matrixToHTML(matrix) {
  if (!matrix.length) return "";
  const rowsHtml = matrix
    .map((row) => `    <tr>${row.map((value) => `<td>${escapeHTML(value)}</td>`).join("")}</tr>`)
    .join("\n");
  return `<table>\n  <tbody>\n${rowsHtml}\n  </tbody>\n</table>`;
}

function escapeHTML(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function matrixToSQLValues(matrix) {
  if (!matrix.length) return "";
  const rowsSql = matrix
    .map((row) => `  (${row.map((value) => `'${String(value ?? "").replaceAll("'", "''")}'`).join(", ")})`)
    .join(",\n");
  return `VALUES\n${rowsSql};`;
}

function matrixToLaTeX(matrix) {
  if (!matrix.length) return "";
  const cols = Math.max(1, ...matrix.map((row) => row.length));
  const rows = matrix
    .map((row) => {
      const padded = Array.from({ length: cols }, (_, index) => escapeLaTeX(row[index] ?? ""));
      return padded.join(" & ");
    })
    .join(" \\\\\n");
  return `\\begin{tabular}{${"l".repeat(cols)}}\n${rows}\n\\end{tabular}`;
}

function escapeLaTeX(value) {
  return String(value ?? "")
    .replaceAll("\\", "\\textbackslash{}")
    .replaceAll("&", "\\&")
    .replaceAll("%", "\\%")
    .replaceAll("$", "\\$")
    .replaceAll("#", "\\#")
    .replaceAll("_", "\\_")
    .replaceAll("{", "\\{")
    .replaceAll("}", "\\}")
    .replaceAll("~", "\\textasciitilde{}")
    .replaceAll("^", "\\textasciicircum{}")
    .replaceAll("|", "\\textbar{}")
    .replaceAll("\n", "\\\\ ");
}

function parseTSV(text) {
  if (!text) return [];
  if (text.replaceAll("\r", "").replaceAll("\n", "") === "") return [];
  const matrix = [];
  let row = [];
  let field = "";
  let quoted = false;
  let fieldStart = true;
  let rowBreak = false;
  for (let index = 0; index < text.length; index++) {
    const char = text[index];
    if (quoted) {
      if (char === '"') {
        if (text[index + 1] === '"') {
          field += '"';
          index++;
        } else {
          quoted = false;
        }
        rowBreak = false;
        continue;
      }
      if (char === "\r") {
        if (text[index + 1] === "\n") index++;
        field += "\n";
      } else {
        field += char;
      }
      rowBreak = false;
      continue;
    }
    if (char === '"' && fieldStart) {
      quoted = true;
      fieldStart = false;
      rowBreak = false;
    } else if (char === "\t") {
      row.push(field);
      field = "";
      fieldStart = true;
      rowBreak = false;
    } else if (char === "\r" || char === "\n") {
      if (char === "\r" && text[index + 1] === "\n") index++;
      row.push(field);
      field = "";
      matrix.push(row);
      row = [];
      fieldStart = true;
      rowBreak = true;
    } else {
      field += char;
      fieldStart = false;
      rowBreak = false;
    }
  }
  if (!rowBreak || row.length > 0 || field.length > 0) {
    row.push(field);
    matrix.push(row);
  }
  return matrix;
}

function parseHTMLTable(htmlText) {
  if (!htmlText || String(htmlText).trim() === "") return [];
  const parser = new DOMParser();
  const doc = parser.parseFromString(String(htmlText), "text/html");
  const table = doc.querySelector("table");
  if (!table) return [];
  const matrix = [];
  const occupied = new Map();
  const rows = Array.from(table.rows ?? []);
  rows.forEach((tableRow, rowIndex) => {
    const row = matrix[rowIndex] ?? [];
    let colIndex = 0;
    while (occupied.get(`${rowIndex}:${colIndex}`)) {
      row[colIndex] = row[colIndex] ?? "";
      colIndex++;
    }
    for (const cell of Array.from(tableRow.cells ?? [])) {
      while (occupied.get(`${rowIndex}:${colIndex}`)) {
        row[colIndex] = row[colIndex] ?? "";
        colIndex++;
      }
      const colspan = Math.max(1, Number.parseInt(cell.getAttribute("colspan") || "1", 10) || 1);
      const rowspan = Math.max(1, Number.parseInt(cell.getAttribute("rowspan") || "1", 10) || 1);
      row[colIndex] = htmlTableCellText(cell);
      for (let colOffset = 1; colOffset < colspan; colOffset++) {
        row[colIndex + colOffset] = row[colIndex + colOffset] ?? "";
      }
      for (let rowOffset = 1; rowOffset < rowspan; rowOffset++) {
        for (let colOffset = 0; colOffset < colspan; colOffset++) {
          occupied.set(`${rowIndex + rowOffset}:${colIndex + colOffset}`, true);
        }
      }
      colIndex += colspan;
    }
    while (occupied.get(`${rowIndex}:${colIndex}`)) {
      row[colIndex] = row[colIndex] ?? "";
      colIndex++;
    }
    matrix[rowIndex] = row;
  });
  return matrix.filter((row) => row.some((value) => value !== ""));
}

function htmlTableCellText(cell) {
  const clone = cell.cloneNode(true);
  for (const br of clone.querySelectorAll("br")) {
    br.replaceWith("\n");
  }
  return clone.textContent.replaceAll("\u00a0", " ");
}

window.__opencsvParseTSV = parseTSV;
window.__opencsvParseHTMLTable = parseHTMLTable;

function updateFileList(name, detail, rowCount, options = {}) {
  const group = currentOpenTargetGroup(options);
  activeFileEntry = {
    name,
    detail,
    rowCount,
    documentId: options.documentId || documentId,
    sourceKind: options.sourceKind || "",
    group,
    dirty,
    readOnly,
  };
  currentEditorTabGroup = group;
  if (activeFileEntry.documentId) editorGroupActiveTabs[group] = activeFileEntry.documentId;
  if (options.tab !== false) upsertOpenFileTab(activeFileEntry);
  if (options.remember) rememberRecentFile(activeFileEntry);
  renderFileList();
  syncAppTitle();
}

function renderFileList() {
  fileList.replaceChildren();
  if (!activeFileEntry) {
    document.body.dataset.openTabCount = "0";
    document.body.dataset.openTabs = "";
    document.body.dataset.activeOpenTab = "";
    document.body.dataset.activeOpenTabVisible = "";
    renderStartActions();
  } else {
    renderOpenFileTabs();
  }

  const openDetails = new Set(openFileTabs.map((item) => item.detail).filter(Boolean));
  const recent = recentFiles().filter((item) => !openDetails.has(item.detail));
  if (recent.length > 0) {
    const label = document.createElement("div");
    label.className = "file-section-label";
    label.textContent = t("file.recent");
    fileList.append(label);
    for (const item of recent) {
      fileList.append(fileButton(item, { onClick: () => openRecentFile(item) }));
    }
  }
  ensureActiveOpenFileTabVisible();
}

function renderOpenFileTabs() {
  const tabs = openFileTabs.length ? openFileTabs : [activeFileEntry].filter(Boolean);
  if (tabs.length > 1) {
    const label = document.createElement("div");
    label.className = "file-section-label";
    label.textContent = "Open Tabs";
    fileList.append(label);
  }
  for (const tab of tabs) {
    const isActive = tab.documentId && tab.documentId === documentId;
    fileList.append(fileTabItem(tab, isActive));
  }
  document.body.dataset.openTabCount = String(tabs.length);
  document.body.dataset.openTabs = tabs.map((tab) => `${tab.documentId || ""}:${tabGroup(tab)}:${tab.name}`).join("|");
  document.body.dataset.activeOpenTab = activeFileEntry?.documentId || documentId || "";
}

function openFileTabKey(entry) {
  return entry?.documentId || entry?.detail || entry?.name || "";
}

function upsertOpenFileTab(entry) {
  const key = openFileTabKey(entry);
  if (!key) return;
  const existing = openFileTabs.find((item) => (
    openFileTabKey(item) === key || (entry.detail && item.detail === entry.detail)
  ));
  const group = normalizeEditorGroup(entry.group || existing?.group || currentEditorTabGroup || activeEditorGroup);
  const next = {
    ...entry,
    group,
    dirty,
    readOnly,
    viewState: entry.viewState ?? currentOpenFileTabState(),
  };
  if (next.documentId === documentId) {
    currentEditorTabGroup = group;
    editorGroupActiveTabs[group] = next.documentId;
  }
  let updated = false;
  openFileTabs = openFileTabs.map((item) => {
    const matches = openFileTabKey(item) === key || (entry.detail && item.detail === entry.detail);
    if (!matches) return item;
    updated = true;
    return next;
  });
  if (!updated) openFileTabs = [...openFileTabs, next];
}

function currentOpenFileTabState() {
  return {
    selected: { ...selected },
    selectionAnchor: { ...selectionAnchor },
    scrollX: Math.max(0, Math.round(scrollX)),
    scrollY: Math.max(0, Math.round(scrollY)),
    viewMode,
  };
}

function captureActiveOpenFileTabState() {
  if (!activeFileEntry) return;
  const key = openFileTabKey(activeFileEntry);
  if (!key) return;
  openFileTabs = openFileTabs.map((tab) => (
    openFileTabKey(tab) === key
      ? {
        ...tab,
        name: activeFileEntry.name,
        detail: activeFileEntry.detail,
        rowCount: totalRows,
        group: currentEditorTabGroup,
        dirty,
        readOnly,
        viewState: currentOpenFileTabState(),
      }
      : tab
  ));
}

function syncActiveFileEntryState() {
  if (!activeFileEntry) return;
  activeFileEntry = {
    ...activeFileEntry,
    documentId: activeFileEntry.documentId || documentId,
    rowCount: totalRows,
    group: currentEditorTabGroup,
    dirty,
    readOnly,
  };
  upsertOpenFileTab(activeFileEntry);
  renderFileList();
}

function restoreOpenFileTabState(tab) {
  const state = tab?.viewState;
  if (!state) return;
  const maxRow = Math.max(0, sourceRowCount - 1);
  const maxCol = Math.max(0, columns.length - 1);
  selected = {
    row: Math.max(0, Math.min(maxRow, Number(state.selected?.row) || 0)),
    col: Math.max(0, Math.min(maxCol, Number(state.selected?.col) || 0)),
  };
  selectionAnchor = {
    row: Math.max(0, Math.min(maxRow, Number(state.selectionAnchor?.row) || selected.row)),
    col: Math.max(0, Math.min(maxCol, Number(state.selectionAnchor?.col) || selected.col)),
  };
  scrollX = Math.max(0, Math.min(maxHorizontalScroll(), Number(state.scrollX) || 0));
  scrollY = Math.max(0, Math.min(maxVerticalScroll(), Number(state.scrollY) || 0));
  wheelRemainderY = 0;
  if (state.viewMode && state.viewMode !== viewMode) setViewMode(state.viewMode);
}

async function switchOpenFileTab(targetDocumentId, options = {}) {
  const tab = openFileTabs.find((item) => item.documentId === targetDocumentId);
  if (!tab?.documentId) return;
  captureActiveOpenFileTabState();
  const targetGroup = normalizeEditorGroup(options.group || tab.group || activeEditorGroup);
  currentEditorTabGroup = targetGroup;
  activeEditorGroup = targetGroup;
  editorGroupActiveTabs[targetGroup] = tab.documentId;
  setOpenFileTabGroup(tab.documentId, targetGroup);
  status.textContent = `Switching to ${tab.name}...`;
  const response = await fetch(`/api/session?id=${encodeURIComponent(tab.documentId)}`);
  const payload = await response.json();
  if (!response.ok) {
    openFileTabs = openFileTabs.filter((item) => item.documentId !== tab.documentId);
    renderFileList();
    status.textContent = `Open tab unavailable: ${payload.error ?? "unknown error"}`;
    return;
  }
  const session = payload.session;
  documentId = session.id;
  readOnly = Boolean(session.readOnly);
  externalFileChanged = false;
  const detail = session.originalPath || session.workingPath || session.document?.metadata?.path || tab.detail;
  pathInput.value = pathInputValueForOpenSession(session, detail, tab.name);
  updateFormatRulePlaceholder(detail);
  const document = {
    ...session.document,
    metadata: { ...session.document.metadata, path: detail },
  };
  metadata = document.metadata;
  clearStructureClipboard();
  loadDocument(document);
  dirty = Boolean(session.dirty);
  restoreOpenFileTabState(tab);
  rows = new Map();
  rowSource = new Map();
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(selected.row);
  updateFileList(session.displayName || tab.name || fileNameFromPath(detail), detail, totalRows, {
    documentId,
    sourceKind: session.sourceKind,
    remember: false,
    group: targetGroup,
  });
  document.body.dataset.editorGroupSwitch = `${targetGroup}:${activeFileEntry.name}`;
  draw();
  startFileStatusChecks();
  await checkFileStatus(false);
  status.textContent = `Switched to ${activeFileEntry.name}.`;
}

function pathInputValueForOpenSession(session, detail, fallbackName = "") {
  const sourceKind = String(session?.sourceKind || "").toLowerCase();
  if (sourceKind === "stdin") return "stdin";
  if (sourceKind === "upload" || sourceKind === "sql-result") return session?.displayName || fallbackName || detail || "";
  return detail || session?.displayName || fallbackName || "";
}

async function closeActiveOpenFileTab() {
  if (!documentId) return;
  await closeOpenFileTab(documentId);
}

async function closeOpenFileTab(targetDocumentId, options = {}) {
  if (!targetDocumentId) return false;
  const target = openFileTabs.find((item) => item.documentId === targetDocumentId)
    || (activeFileEntry?.documentId === targetDocumentId ? activeFileEntry : null);
  const response = await fetch("/api/close", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId: targetDocumentId,
      force: Boolean(options.force),
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    const message = payload.error ?? "close failed";
    document.body.dataset.closeTabBlocked = `${targetDocumentId}:${message}`;
    if (message === "document has unsaved changes" && !options.force) {
      openDirtyClosePanel(target, message);
    }
    status.textContent = message === "document has unsaved changes"
      ? `Save ${target?.name ?? "the file"} before closing.`
      : `Close tab failed: ${message}`;
    return false;
  }

  const closedName = target?.name || payload.displayName || "tab";
  const closedGroup = tabGroup(target);
  const targetDraftID = target?.detail && target.detail !== "temporary upload"
    ? target.detail
    : [target?.name || "untitled", targetDocumentId].filter(Boolean).join(":");
  clearDirtyDraft(targetDraftID);
  document.body.dataset.closedOpenTab = `${targetDocumentId}:${closedName}`;
  document.body.dataset.closeTabBlocked = "";
  openFileTabs = openFileTabs.filter((item) => item.documentId !== targetDocumentId);
  if (editorGroupActiveTabs[closedGroup] === targetDocumentId) {
    editorGroupActiveTabs[closedGroup] = tabForEditorGroup(closedGroup)?.documentId || "";
  }
  if (targetDocumentId === documentId) {
    const next = tabForEditorGroup(closedGroup)
      || tabForEditorGroup(closedGroup === "primary" ? "secondary" : "primary")
      || openFileTabs[openFileTabs.length - 1];
    if (next?.documentId) {
      await switchOpenFileTab(next.documentId, { group: tabGroup(next) });
      status.textContent = `Closed ${closedName}.`;
      return true;
    }
    clearCurrentDocumentState(closedName);
    return true;
  }

  syncEditorGroupActiveTabs();
  renderEditorGroups();
  renderFileList();
  status.textContent = `Closed ${closedName}.`;
  return true;
}

function openDirtyClosePanel(tab, message = "document has unsaved changes") {
  if (!tab?.documentId || !dirtyClosePanel) return;
  pendingDirtyCloseTab = {
    documentId: tab.documentId,
    name: tab.name || "Unsaved document",
    detail: tab.detail || "",
  };
  dirtyCloseName.textContent = pendingDirtyCloseTab.name;
  dirtyCloseDetail.textContent = pendingDirtyCloseTab.detail
    ? `${pendingDirtyCloseTab.detail} has unsaved edits. Save it before closing, or discard the tab explicitly.`
    : "This tab has unsaved edits. Save it before closing, or discard the tab explicitly.";
  dirtyClosePanel.hidden = false;
  document.body.dataset.closeTabConfirm = `${pendingDirtyCloseTab.documentId}:${pendingDirtyCloseTab.name}`;
  document.body.dataset.closeTabConfirmReason = message;
  discardDirtyCloseButton.focus();
}

function closeDirtyClosePanel() {
  if (!dirtyClosePanel) return;
  dirtyClosePanel.hidden = true;
  pendingDirtyCloseTab = null;
  document.body.dataset.closeTabConfirm = "";
  document.body.dataset.closeTabConfirmReason = "";
}

async function confirmDirtyCloseDiscard() {
  if (!pendingDirtyCloseTab) return;
  const target = { ...pendingDirtyCloseTab };
  closeDirtyClosePanel();
  const closed = await closeOpenFileTab(target.documentId, { force: true });
  if (closed) {
    document.body.dataset.closeTabDiscarded = `${target.documentId}:${target.name}`;
  }
}

function clearCurrentDocumentState(closedName = "") {
  if (fileStatusTimer) clearInterval(fileStatusTimer);
  fileStatusTimer = 0;
  windowGeneration++;
  loadingWindow = false;
  documentId = "";
  activeFileEntry = null;
  currentEditorTabGroup = "primary";
  activeEditorGroup = "primary";
  editorGroupActiveTabs = { primary: "", secondary: "" };
  metadata = null;
  rows = new Map();
  rowSource = new Map();
  windowStart = 0;
  windowEnd = 0;
  sourceRowCount = 0;
  aggregateRows = [];
  syncTotalRows();
  columns = [];
  allColumns = [];
  hiddenSourceColumns = new Set();
  selected = { row: 0, col: 0 };
  selectionAnchor = { row: 0, col: 0 };
  scrollX = 0;
  scrollY = 0;
  wheelRemainderY = 0;
  dirty = false;
  readOnly = false;
  externalFileChanged = false;
  externalChangePanelDismissed = false;
  viewMapped = false;
  sortDescending = false;
  sortKeys = [];
  clearActiveFilter();
  invalidateFindState();
  setSmartSelectionRanges([]);
  clearStructureClipboard();
  renderFileList();
  draw();
  status.textContent = closedName ? `Closed ${closedName}.` : "Closed tab.";
}

function renderStartActions() {
  const newFileButton = document.createElement("button");
  newFileButton.className = "file-button file-action";
  newFileButton.type = "button";
  newFileButton.append(fileText(t("start.newFile"), t("start.newHelp")));
  newFileButton.addEventListener("click", () => {
    void createNewFile();
  });

  const openPathButton = document.createElement("button");
  openPathButton.className = "file-button file-action";
  openPathButton.type = "button";
  openPathButton.append(fileText(t("start.openCsv"), t("start.openHelp")));
  openPathButton.addEventListener("click", () => {
    pathInput.focus();
    pathInput.select();
  });

  const chooseButton = document.createElement("button");
  chooseButton.className = "file-button file-action";
  chooseButton.type = "button";
  chooseButton.append(fileText(t("start.chooseFile"), t("start.chooseHelp")));
  chooseButton.addEventListener("click", () => fileInput.click());

  const settingsButton = document.createElement("button");
  settingsButton.className = "file-button file-action";
  settingsButton.type = "button";
  settingsButton.append(fileText(t("start.settings"), t("start.settingsHelp")));
  settingsButton.addEventListener("click", () => openSettingsPanel());

  const supportButton = document.createElement("button");
  supportButton.className = "file-button file-action";
  supportButton.type = "button";
  supportButton.append(fileText(t("start.support"), t("start.supportHelp")));
  supportButton.addEventListener("click", () => openAboutSupportPanel());

  fileList.append(newFileButton, openPathButton, chooseButton, settingsButton, supportButton);

  const drafts = dirtyDrafts();
  if (drafts.length > 0) {
    const label = document.createElement("div");
    label.className = "file-section-label";
    label.textContent = "Recovery Drafts";
    fileList.append(label);
    for (const draft of drafts) {
      fileList.append(recoveryDraftItem(draft));
    }
  }
}

function recoveryDraftItem(draft) {
  const row = document.createElement("div");
  row.className = "recovery-draft-row";
  row.dataset.draftId = draft.id || "";
  const restoreButton = fileButton({
    name: draft.name || "Recovered draft",
    detail: draft.detail || "unsaved recovery draft",
    rowCount: Number.isFinite(draft.rowCount) ? draft.rowCount : null,
  }, {
    onClick: () => {
      void restoreDirtyDraft(draft.id);
    },
  });
  const discardButton = document.createElement("button");
  discardButton.className = "file-tab-close";
  discardButton.type = "button";
  discardButton.textContent = "x";
  discardButton.title = `Discard ${draft.name || "recovery draft"}`;
  discardButton.setAttribute("aria-label", discardButton.title);
  discardButton.addEventListener("click", (event) => {
    event.stopPropagation();
    discardDirtyDraft(draft.id);
  });
  row.append(restoreButton, discardButton);
  return row;
}

function fileTabItem(tab, isActive) {
  const row = document.createElement("div");
  row.className = isActive ? "file-tab-row active" : "file-tab-row";
  if (tab.documentId) row.dataset.documentId = tab.documentId;
  row.dataset.editorGroup = tabGroup(tab);
  row.append(fileButton(tab, {
    active: isActive,
    onClick: () => {
      if (isActive) {
        void focusEditorGroup(tabGroup(tab), { activateTab: false });
        draw();
        return;
      }
      void switchOpenFileTab(tab.documentId, { group: tabGroup(tab) });
    },
  }));

  const closeButton = document.createElement("button");
  closeButton.className = "file-tab-close";
  closeButton.type = "button";
  closeButton.textContent = "x";
  closeButton.title = `Close ${tab.name}`;
  closeButton.setAttribute("aria-label", `Close ${tab.name}`);
  if (tab.documentId) closeButton.dataset.documentId = tab.documentId;
  closeButton.addEventListener("click", (event) => {
    event.stopPropagation();
    void closeOpenFileTab(tab.documentId);
  });
  row.append(closeButton);
  return row;
}

function ensureActiveOpenFileTabVisible() {
  const activeTab = fileList.querySelector(".file-tab-row.active");
  if (activeTab) {
    scrollOpenFileTabIntoView(activeTab);
  }
  window.scrollTo(0, 0);
  syncOpenFileTabOverflowState();
}

function scrollOpenFileTabIntoView(tabNode) {
  const padding = 4;
  const tabTop = tabNode.offsetTop;
  const tabBottom = tabTop + tabNode.offsetHeight;
  const viewTop = fileList.scrollTop;
  const viewBottom = viewTop + fileList.clientHeight;
  if (tabTop < viewTop) {
    fileList.scrollTop = Math.max(0, tabTop - padding);
  } else if (tabBottom > viewBottom) {
    fileList.scrollTop = Math.max(0, tabBottom - fileList.clientHeight + padding);
  }
}

function syncOpenFileTabOverflowState() {
  const activeTab = fileList.querySelector(".file-tab-row.active");
  const listRect = fileList.getBoundingClientRect();
  const activeRect = activeTab?.getBoundingClientRect();
  const overflows = fileList.scrollHeight > fileList.clientHeight + 1;
  const visible = Boolean(
    activeRect &&
    activeRect.top >= listRect.top - 1 &&
    activeRect.bottom <= listRect.bottom + 1,
  );
  document.body.dataset.openTabOverflow = overflows ? "true" : "false";
  document.body.dataset.openTabScrollTop = String(Math.round(fileList.scrollTop));
  document.body.dataset.activeOpenTabVisible = activeTab ? String(visible) : "";
}

function fileButton(entry, options = {}) {
  const button = document.createElement("button");
  button.className = options.active ? "file-button active" : "file-button";
  button.type = "button";
  button.title = entry.detail;
  if (entry.documentId) button.dataset.documentId = entry.documentId;
  button.append(fileText(entry.name, fileMeta(entry)));
  if (options.onClick) button.addEventListener("click", options.onClick);
  return button;
}

function fileText(name, metaText) {
  const fragment = document.createDocumentFragment();
  const title = document.createElement("span");
  title.className = "file-name";
  title.textContent = name;

  const meta = document.createElement("span");
  meta.className = "file-meta";
  meta.textContent = metaText;

  fragment.append(title, meta);
  return fragment;
}

function fileMeta(entry) {
  const rowsText = Number.isFinite(entry.rowCount)
    ? `${entry.rowCount.toLocaleString()} rows`
    : "recent file";
  const flags = [];
  if (entry.dirty) flags.push("edited");
  if (entry.readOnly) flags.push("read-only");
  const flagText = flags.length ? ` · ${flags.join(", ")}` : "";
  return `${rowsText}${flagText} · ${entry.detail}`;
}

function recentFiles() {
  try {
    const parsed = JSON.parse(localStorage.getItem(recentFilesStorageKey) || "[]");
    if (!Array.isArray(parsed)) return [];
    return parsed
      .filter((item) => item && typeof item.detail === "string" && item.detail)
      .map((item) => ({
        name: String(item.name || fileNameFromPath(item.detail)),
        detail: item.detail,
        rowCount: Number.isFinite(item.rowCount) ? item.rowCount : null,
        lastOpened: Number.isFinite(item.lastOpened) ? item.lastOpened : 0,
      }))
      .sort((left, right) => right.lastOpened - left.lastOpened)
      .slice(0, recentFilesLimit);
  } catch {
    return [];
  }
}

function rememberRecentFile(entry) {
  if (!entry?.detail || entry.detail === "temporary upload") return;
  const next = [
    { ...entry, lastOpened: Date.now() },
    ...recentFiles().filter((item) => item.detail !== entry.detail),
  ].slice(0, recentFilesLimit);
  try {
    localStorage.setItem(recentFilesStorageKey, JSON.stringify(next));
  } catch {
    status.textContent = "Recent files could not be saved.";
  }
}

function removeRecentFile(detail) {
  const next = recentFiles().filter((item) => item.detail !== detail);
  try {
    localStorage.setItem(recentFilesStorageKey, JSON.stringify(next));
  } catch {
    status.textContent = "Recent files could not be updated.";
  }
}

function dirtyDrafts() {
  try {
    const parsed = JSON.parse(localStorage.getItem(dirtyDraftsStorageKey) || "[]");
    if (!Array.isArray(parsed)) return [];
    return parsed
      .filter((item) => item && typeof item.id === "string" && typeof item.csv === "string" && item.csv)
      .map((item) => ({
        id: item.id,
        name: String(item.name || "Recovered draft.csv"),
        detail: String(item.detail || "unsaved recovery draft"),
        csv: item.csv,
        rowCount: Number.isFinite(item.rowCount) ? item.rowCount : null,
        columnCount: Number.isFinite(item.columnCount) ? item.columnCount : null,
        updatedAt: Number.isFinite(item.updatedAt) ? item.updatedAt : 0,
      }))
      .sort((left, right) => right.updatedAt - left.updatedAt)
      .slice(0, dirtyDraftsLimit);
  } catch {
    return [];
  }
}

function dirtyDraftID() {
  const detail = activeFileEntry?.detail || metadata?.path || "";
  if (detail && detail !== "temporary upload") return detail;
  return [activeFileEntry?.name || metadata?.path || "untitled", documentId].filter(Boolean).join(":");
}

function dirtyDraftName() {
  const name = activeFileEntry?.name || fileNameFromPath(metadata?.path || "") || "Recovered draft.csv";
  return name.toLowerCase().endsWith(".csv") ? name : `${name}.csv`;
}

function scheduleDirtyDraftSave() {
  if (!dirty || !documentId) return;
  if (dirtyDraftSaveTimer) window.clearTimeout(dirtyDraftSaveTimer);
  dirtyDraftSaveTimer = window.setTimeout(() => {
    dirtyDraftSaveTimer = 0;
    void saveDirtyDraft();
  }, 300);
}

async function saveDirtyDraft() {
  if (!dirty || !documentId || dirtyDraftSaveInFlight) return;
  const id = dirtyDraftID();
  if (!id) return;
  dirtyDraftSaveInFlight = true;
  try {
    const response = await fetch("/api/export-view", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        documentId,
        columns: (allColumns.length ? allColumns : columns).map((column) => column.sourceCol),
        format: "csv",
        fullDocument: true,
      }),
    });
    const text = await response.text();
    if (!response.ok) throw new Error(errorFromResponseText(text));
    if (!dirty || !documentId) return;
    if (text.length > dirtyDraftMaxBytes) {
      document.body.dataset.dirtyDraft = `too-large:${text.length}`;
      return;
    }
    const draft = {
      id,
      name: dirtyDraftName(),
      detail: activeFileEntry?.detail || metadata?.path || "temporary upload",
      csv: text,
      rowCount: sourceRowCount,
      columnCount: columns.length,
      updatedAt: Date.now(),
    };
    const next = [
      draft,
      ...dirtyDrafts().filter((item) => item.id !== id),
    ].slice(0, dirtyDraftsLimit);
    localStorage.setItem(dirtyDraftsStorageKey, JSON.stringify(next));
    document.body.dataset.dirtyDraft = `${draft.name}:${draft.rowCount}:${draft.columnCount}`;
  } catch (error) {
    document.body.dataset.dirtyDraft = `failed:${error.message}`;
  } finally {
    dirtyDraftSaveInFlight = false;
  }
}

function clearDirtyDraft(id = dirtyDraftID()) {
  if (!id) return;
  try {
    const next = dirtyDrafts().filter((item) => item.id !== id);
    localStorage.setItem(dirtyDraftsStorageKey, JSON.stringify(next));
    document.body.dataset.dirtyDraftCleared = id;
  } catch {
    status.textContent = "Recovery draft could not be cleared.";
  }
}

function discardDirtyDraft(id) {
  clearDirtyDraft(id);
  renderFileList();
  status.textContent = "Recovery draft discarded.";
}

async function restoreDirtyDraft(id) {
  const draft = dirtyDrafts().find((item) => item.id === id);
  if (!draft) {
    status.textContent = "Recovery draft was not found.";
    renderFileList();
    return;
  }
  const file = new File([draft.csv], draft.name || "Recovered draft.csv", { type: "text/csv" });
  await uploadAndOpen(file);
  await markRecoveredDraftDirty();
  clearDirtyDraft(id);
  document.body.dataset.recoveryDraftRestored = `${draft.name}:${draft.rowCount ?? ""}:${draft.columnCount ?? ""}`;
  status.textContent = `Recovered draft ${draft.name}. Save As to keep it.`;
}

async function markRecoveredDraftDirty() {
  if (!documentId || sourceRowCount <= 0 || columns.length === 0) return;
  const row = rowSource.get(0) ?? dataStartRow;
  const col = sourceColFor(0);
  const value = valueAt(0, 0);
  const response = await fetch("/api/edit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, row, col, value }),
  });
  const payload = await response.json();
  if (!response.ok) throw new Error(payload.error ?? "recovery dirty mark failed");
  dirty = true;
  syncActiveFileEntryState();
  scheduleDirtyDraftSave();
}

async function openRecentFile(entry) {
  pathInput.value = entry.detail;
  try {
    await openPath(entry.detail);
  } catch (error) {
    removeRecentFile(entry.detail);
    renderFileList();
    status.textContent = `Removed missing recent file: ${error.message}`;
  }
}

function fileNameFromPath(path) {
  return path.split(/[\\/]/).filter(Boolean).pop() || path;
}

function buildColumns(document) {
  const first = document.preview[0] ?? [];
  const hasHeader = document.metadata.dialect.hasHeader;
  const previewColumnCount = Math.max(...document.preview.map((record) => record.length), 1);
  const columnCount =
    Math.max(document.metadata.dialect.estimatedColumns || 0, previewColumnCount);
  return Array.from({ length: columnCount }, (_, index) => ({
    name: hasHeader && first[index] ? first[index] : `Column ${index + 1}`,
    sourceCol: index,
    width: Math.max(minColumnWidth, Math.min(240, String(first[index] ?? "").length * 9 + 64)),
    type: index === 0 ? "mixed" : "string",
  }));
}

function ensureSourceColumnsForEdits(edits) {
  if (!Array.isArray(edits) || edits.length === 0) return;
  const maxSourceCol = Math.max(-1, ...edits.map((edit) => Number.isInteger(edit.col) ? edit.col : -1));
  if (maxSourceCol < 0) return;
  if (metadata?.dialect) {
    metadata.dialect.estimatedColumns = Math.max(metadata.dialect.estimatedColumns || 0, maxSourceCol + 1);
  }
  let changed = false;
  const existing = new Set(allColumns.map((column) => column.sourceCol));
  for (let sourceCol = 0; sourceCol <= maxSourceCol; sourceCol++) {
    if (existing.has(sourceCol)) continue;
    allColumns.push({
      name: `Column ${sourceCol + 1}`,
      sourceCol,
      width: minColumnWidth,
      type: sourceCol === 0 ? "mixed" : "string",
    });
    existing.add(sourceCol);
    changed = true;
  }
  if (!changed) return;
  allColumns.sort((left, right) => left.sourceCol - right.sourceCol);
  applyHiddenColumns();
  document.body.dataset.extendedColumns = String(maxSourceCol + 1);
}

function ensureSourceRowsForEdits(edits) {
  if (!Array.isArray(edits) || edits.length === 0 || viewMapped) return;
  const maxSourceRow = Math.max(-1, ...edits.map((edit) => Number.isInteger(edit.row) ? edit.row : -1));
  const nextSourceRowCount = Math.max(sourceRowCount, maxSourceRow - dataStartRow + 1);
  if (nextSourceRowCount <= sourceRowCount) return;
  sourceRowCount = nextSourceRowCount;
  syncTotalRows();
  document.body.dataset.extendedRows = String(sourceRowCount);
}

async function maybeLoadWindow(row) {
  if (!documentId || loadingWindow) return;
  if (row >= sourceRowCount) return;
  if (row >= windowStart && row < windowEnd) {
    const nearWindowEnd = row >= windowEnd - 20;
    if (!nearWindowEnd || windowEnd >= sourceRowCount) return;
  }
  loadingWindow = true;
  const generation = windowGeneration;
  const displayStart = Math.max(0, row - 20);
  const sourceStart = viewMapped ? displayStart : displayStart + dataStartRow;
  try {
    const response = await fetch(`/api/window?id=${encodeURIComponent(documentId)}&start=${sourceStart}&count=${windowSize}`);
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "window load failed");
    if (generation !== windowGeneration) return;
    payload.rows.forEach((record, index) => {
      const sourceRow = payload.sourceRows?.[index] ?? payload.start + index;
      const displayRow = viewMapped ? payload.start + index : sourceRow - dataStartRow;
      rows.set(displayRow, record);
      rowSource.set(displayRow, sourceRow);
    });
    windowStart = displayStart;
    windowEnd = displayStart + payload.rows.length;
    sourceRowCount = viewMapped ? payload.totalRows : Math.max(0, payload.totalRows - dataStartRow);
    syncTotalRows();
    draw();
  } catch (error) {
    status.textContent = `Window load failed: ${error.message}`;
  } finally {
    loadingWindow = false;
  }
}

async function findCurrentQuery() {
  await runFind(0, true);
}

function startFindProgress(query) {
  if (activeFindController) activeFindController.abort();
  activeFindController = new AbortController();
  if (findProgressTimer) window.clearInterval(findProgressTimer);
  findProgressStartedAt = performance.now();
  findProgressQuery = query;
  cancelFindButton.hidden = false;
  cancelFindButton.disabled = false;
  findButton.disabled = true;
  findPrevButton.disabled = true;
  findNextButton.disabled = true;
  document.body.dataset.findProgress = "running";
  document.body.dataset.findElapsedMs = "0";
  updateFindProgressStatus();
  findProgressTimer = window.setInterval(updateFindProgressStatus, 250);
  return activeFindController;
}

function updateFindProgressStatus() {
  const elapsed = findElapsedMs();
  document.body.dataset.findElapsedMs = String(Math.round(elapsed));
  status.textContent = tf("find.status.runningElapsed", {
    query: findProgressQuery,
    seconds: formatElapsedSeconds(elapsed),
  });
}

function findElapsedMs() {
  return findProgressStartedAt ? Math.max(0, performance.now() - findProgressStartedAt) : 0;
}

function stopFindProgressTimer() {
  if (findProgressTimer) {
    window.clearInterval(findProgressTimer);
    findProgressTimer = 0;
  }
  const elapsed = findElapsedMs();
  document.body.dataset.findElapsedMs = String(Math.round(elapsed));
  return elapsed;
}

function finishFindProgress(controller) {
  if (activeFindController !== controller) return;
  stopFindProgressTimer();
  activeFindController = null;
  cancelFindButton.hidden = true;
  cancelFindButton.disabled = true;
  findButton.disabled = false;
  findPrevButton.disabled = false;
  findNextButton.disabled = false;
  if (document.body.dataset.findProgress === "running") {
    document.body.dataset.findProgress = "idle";
  }
}

function cancelFindRequest() {
  if (!activeFindController) return;
  activeFindController.abort();
  const elapsed = stopFindProgressTimer();
  activeFindController = null;
  cancelFindButton.hidden = true;
  cancelFindButton.disabled = true;
  findButton.disabled = false;
  findPrevButton.disabled = false;
  findNextButton.disabled = false;
  document.body.dataset.findProgress = "canceled";
  status.textContent = tf("find.status.canceledElapsed", { seconds: formatElapsedSeconds(elapsed) });
}

async function focusFindInput(options = {}) {
  const seeded = applyFindSeed(options.target ?? document.activeElement);
  commandInput.focus();
  commandInput.select();
  if ((options.run || seeded) && commandInput.value.trim()) {
    await findCurrentQuery();
  }
}

function applyFindSeed(target) {
  const seed = findSeedFromSelection(target);
  if (!seed) {
    document.body.dataset.findSeed = "";
    return false;
  }
  commandInput.value = seed.value;
  document.body.dataset.findSeed = `${seed.source}:${String(seed.value).replace(/\s+/g, " ").slice(0, 80)}`;
  invalidateFindState();
  return true;
}

function findSeedFromSelection(target) {
  const selectedText = editableSelectedText(target);
  if (selectedText && target !== commandInput) {
    return {
      source: target === editor
        ? "cell-editor"
        : target === cellValueInput
          ? "cell-panel"
          : target === cellValueDialogInput
            ? "cell-dialog"
            : "text-selection",
      value: selectedText,
    };
  }
  if (commandInput.value.trim()) return null;
  if (selectionCellCount() !== 1 || aggregateRowAt(selected.row)) return null;
  const value = String(valueAt(selected.row, selected.col) ?? "");
  return value ? { source: "cell", value } : null;
}

function editableSelectedText(target) {
  if (!(target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement)) return "";
  if (target.type && ["button", "checkbox", "file", "radio"].includes(target.type)) return "";
  try {
    const start = target.selectionStart;
    const end = target.selectionEnd;
    if (!Number.isInteger(start) || !Number.isInteger(end) || end <= start) return "";
    return target.value.slice(start, end);
  } catch {
    return "";
  }
}

function searchQueryOrSeed(target, emptyMessage) {
  let search = parseSearchQuery(commandInput.value);
  if (search.query) return search;
  if (!applyFindSeed(target ?? document.activeElement)) {
    status.textContent = emptyMessage;
    return null;
  }
  search = parseSearchQuery(commandInput.value);
  return search.query ? search : null;
}

async function findNextMatch(direction) {
  if (!documentId) return;
  await runFind(direction, false);
}

async function runFind(direction, forceRefresh) {
  if (!documentId) return;
  const search = parseSearchQuery(commandInput.value);
  const query = search.query;
  if (!query) {
    status.textContent = t("find.status.enterQuery");
    findState = emptyFindState();
    updateFindMatchCount();
    return;
  }
  const key = findKey(search);
  const criteriaKey = findCriteriaKey(search);
  const sameSelectionScope =
    search.inSelection &&
    !forceRefresh &&
    findState.criteriaKey === criteriaKey &&
    Array.isArray(findState.scopeCells) &&
    findState.scopeCells.length > 0;
  const stateKey = sameSelectionScope ? findState.key : key;
  if (forceRefresh || (!sameSelectionScope && findState.key !== key) || findState.matches.length === 0) {
    const ok = await refreshFindState(search, stateKey, { offset: 0, index: 0 });
    if (!ok) return;
  }
  if (findState.matches.length === 0) {
    status.textContent = tf("find.status.noMatchesTimed", {
      query,
      seconds: formatElapsedSeconds(findState.elapsedMs),
    });
    updateFindMatchCount();
    return;
  }
  if (direction !== 0) {
    const paged = await advanceFindState(search, stateKey, direction);
    if (!paged) return;
  }
  await selectFindMatch(findState.matches[findState.index], search.query);
}

async function advanceFindState(search, key, direction) {
  const nextIndex = findState.index + direction;
  if (nextIndex >= 0 && nextIndex < findState.matches.length) {
    findState.index = nextIndex;
    return true;
  }
  if (direction > 0) {
    const nextOffset = findState.offset + findState.matches.length;
    return refreshFindState(search, key, {
      offset: nextOffset < findState.total ? nextOffset : 0,
      index: 0,
      scopeCells: findState.scopeCells,
    });
  }
  const lastPageOffset = findState.total > 0
    ? Math.floor((findState.total - 1) / findState.limit) * findState.limit
    : 0;
  const previousOffset = findState.offset > 0
    ? Math.max(0, findState.offset - findState.limit)
    : lastPageOffset;
  return refreshFindState(search, key, {
    offset: previousOffset,
    index: -1,
    scopeCells: findState.scopeCells,
  });
}

async function refreshFindState(search, key, options = {}) {
  const offset = Math.max(0, options.offset ?? 0);
  const cells = options.scopeCells ?? (search.inSelection ? selectedCellsForSearch() : []);
  if (cells == null) return false;
  const controller = startFindProgress(search.query);
  let payload = null;
  let elapsed = 0;
  try {
    const response = await fetch("/api/find", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      signal: controller.signal,
      body: JSON.stringify({
        documentId,
        query: search.query,
        regex: search.regex,
        caseSensitive: search.caseSensitive,
        wholeCell: search.wholeCell,
        cells,
        limit: findPageSize,
        offset,
      }),
    });
    payload = await response.json();
    if (!response.ok) {
      stopFindProgressTimer();
      status.textContent = tf("find.status.failed", { message: payload.error ?? "unknown error" });
      return false;
    }
    elapsed = stopFindProgressTimer();
    document.body.dataset.findStreamed = String(Boolean(payload.streamed));
  } catch (error) {
    if (error.name === "AbortError") {
      if (!activeFindController || activeFindController === controller) {
        const canceledElapsed = stopFindProgressTimer();
        document.body.dataset.findProgress = "canceled";
        status.textContent = tf("find.status.canceledElapsed", { seconds: formatElapsedSeconds(canceledElapsed) });
      }
      return false;
    }
    stopFindProgressTimer();
    status.textContent = tf("find.status.failed", { message: error.message });
    return false;
  } finally {
    finishFindProgress(controller);
  }
  findState = {
    key,
    criteriaKey: findCriteriaKey(search),
    search: findCriteriaKeyObject(search),
    matches: payload.matches ?? [],
    total: payload.total ?? payload.matches?.length ?? 0,
    offset: payload.offset ?? offset,
    limit: findPageSize,
    truncated: Boolean(payload.truncated),
    index: 0,
    scopeCells: cells,
    elapsedMs: elapsed,
    streamed: Boolean(payload.streamed),
  };
  if (options.index === -1) {
    findState.index = Math.max(0, findState.matches.length - 1);
  } else {
    findState.index = Math.max(0, Math.min(findState.matches.length - 1, options.index ?? 0));
  }
  updateFindMatchCount();
  updateFindHighlightDataset();
  return true;
}

async function selectFindMatch(match, query) {
  if (viewMapped) await clearView();
  const displayRow = Math.max(0, match.row - dataStartRow);
  selected = { row: displayRow, col: match.col };
  selectionAnchor = { ...selected };
  scrollY = Math.max(0, displayRow * rowHeight - rowHeight * 2);
  wheelRemainderY = 0;
  await maybeLoadWindow(displayRow);
  draw();
  updateFindMatchCount();
  updateFindHighlightDataset();
  const ordinal = findState.offset + findState.index + 1;
  const total = findState.total;
  const pageStart = findState.matches.length ? findState.offset + 1 : 0;
  const pageEnd = findState.offset + findState.matches.length;
  const loadedNote = findState.truncated || findState.offset > 0
    ? tf("find.status.page", { start: formatCount(pageStart), end: formatCount(pageEnd) })
    : "";
  status.textContent = tf("find.status.foundTimed", {
    current: formatCount(ordinal),
    total: formatCount(total),
    query,
    row: formatCount(displayRow + 1),
    column: formatCount(match.col + 1),
    page: loadedNote,
    seconds: formatElapsedSeconds(findState.elapsedMs),
  });
}

async function replaceCurrentQuery(all, options = {}) {
  if (!documentId) return;
  if (!ensureEditable("Replace")) return;
  const search = searchQueryOrSeed(options.target, t("replace.status.enterQuery"));
  if (!search) return;
  const query = search.query;
  const cells = replaceScopeCells(search, all);
  if (cells == null) return;
  const response = await fetch("/api/replace", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      query,
      replacement: replaceInput.value,
      regex: search.regex,
      caseSensitive: search.caseSensitive,
      wholeCell: search.wholeCell,
      cells,
      all,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = tf("replace.status.failed", { message: payload.error ?? "unknown error" });
    return;
  }
  if (!payload.count) {
    status.textContent = tf("replace.status.noReplacements", { query });
    return;
  }
  dirty = payload.dirty ?? true;
  profileData = null;
  const first = payload.changes?.[0];
  document.body.dataset.replaceRanges = (payload.changes ?? [])
    .flatMap((change) => formatServerMatchRanges(change))
    .join(",");
  document.body.dataset.replaceChange = first
    ? replaceChangePreview(first)
    : String(payload.count);
  invalidateFindState();
  rows = new Map();
  rowSource = new Map();
  windowStart = 0;
  windowEnd = 0;
  if (first) {
    const displayRow = Math.max(0, first.row - dataStartRow);
    selected = { row: displayRow, col: first.col };
    selectionAnchor = { ...selected };
    scrollY = Math.max(0, displayRow * rowHeight - rowHeight * 2);
    wheelRemainderY = 0;
    await maybeLoadWindow(displayRow);
  }
  draw();
  status.textContent = all
    ? tf("replace.status.replacedCells", { count: formatCount(payload.count) })
    : document.body.dataset.replaceScope.startsWith("current:")
      ? t("replace.status.replacedCurrent")
      : t("replace.status.replacedOne");
  clearReplacePreview();
}

function replaceChangePreview(change) {
  if (!change || typeof change !== "object") return "";
  const displayRow = Math.max(0, change.row - dataStartRow);
  const visibleCol = visibleColForSource(change.col);
  const before = visibleCol >= 0 ? valueAt(displayRow, visibleCol) : "";
  return `${change.row}:${change.col}:${before}->${change.value ?? ""}`;
}

function replaceScopeCells(search, all) {
  if (all) {
    document.body.dataset.replaceScope = search.inSelection ? "selection" : "all";
    return search.inSelection ? cellsForReplaceScope(search) : [];
  }
  const active = activeFindMatchCell(search);
  if (active) {
    document.body.dataset.replaceScope = `current:${active.row}:${active.col}`;
    return [active];
  }
  document.body.dataset.replaceScope = search.inSelection ? "selection-first" : "first";
  return search.inSelection ? cellsForReplaceScope(search) : [];
}

function activeFindMatchCell(search) {
  if (!findState.matches.length || findState.criteriaKey !== findCriteriaKey(search)) return null;
  const match = findState.matches[findState.index];
  if (!match || !Number.isInteger(match.row) || !Number.isInteger(match.col)) return null;
  return { row: match.row, col: match.col };
}

async function previewReplace(options = {}) {
  if (!documentId) return;
  const search = searchQueryOrSeed(options.target, t("replace.status.enterPreviewQuery"));
  if (!search) return;
  const query = search.query;
  const cells = search.inSelection ? cellsForReplaceScope(search) : [];
  if (cells == null) return;
  const response = await fetch("/api/replace", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      query,
      replacement: replaceInput.value,
      regex: search.regex,
      caseSensitive: search.caseSensitive,
      wholeCell: search.wholeCell,
      cells,
      all: true,
      preview: true,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = tf("replace.status.previewFailed", { message: payload.error ?? "unknown error" });
    return;
  }
  if (!payload.count) {
    clearReplacePreview();
    status.textContent = tf("replace.status.noReplacements", { query });
    return;
  }
  const first = payload.changes?.[0];
  if (first) {
    const displayRow = Math.max(0, first.row - dataStartRow);
    const visibleCol = Math.max(0, visibleColForSource(first.col));
    selected = { row: displayRow, col: visibleCol };
    selectionAnchor = { ...selected };
    scrollY = Math.max(0, displayRow * rowHeight - rowHeight * 2);
    wheelRemainderY = 0;
    await maybeLoadWindow(displayRow);
  }
  draw();
  const firstVisibleCol = first ? visibleColForSource(first.col) : -1;
  const before = first && firstVisibleCol >= 0
    ? valueAt(Math.max(0, first.row - dataStartRow), firstVisibleCol)
    : "";
  const after = first?.value ?? "";
  const rowLabel = first
    ? tf("replace.status.previewFirst", {
        row: formatCount(Math.max(0, first.row - dataStartRow) + 1),
        column: formatCount(first.col + 1),
      })
    : "";
  document.body.dataset.replacePreviewCount = String(payload.count);
  document.body.dataset.replacePreview = first
    ? `${first.row}:${first.col}:${before}->${after}`
    : String(payload.count);
  document.body.dataset.replacePreviewRanges = (payload.changes ?? [])
    .flatMap((change) => formatServerMatchRanges(change))
    .join(",");
  status.textContent = tf("replace.status.preview", {
    count: formatCount(payload.count),
    rowLabel,
    before: truncateMenuValue(before),
    after: truncateMenuValue(after),
  });
}

function clearReplacePreview() {
  document.body.dataset.replacePreview = "";
  document.body.dataset.replacePreviewCount = "";
  document.body.dataset.replacePreviewRanges = "";
}

function parseSearchQuery(raw) {
  const trimmed = raw.trim();
  const match = trimmed.match(/^\/(.+)\/([i]*)$/);
  const query = match ? match[1] : trimmed;
  const slashRegex = Boolean(match);
  const regex = findRegex || slashRegex;
  const caseSensitive = findCaseSensitive || (slashRegex && !match[2].includes("i"));
  return { query, regex, caseSensitive, wholeCell: findWholeCell, inSelection: findInSelection };
}

function findKey(search) {
  return JSON.stringify({
    ...findCriteriaKeyObject(search),
    inSelection: search.inSelection,
    selection: search.inSelection ? selectedSearchScopeKey() : "",
  });
}

function findCriteriaKey(search) {
  return JSON.stringify(findCriteriaKeyObject(search));
}

function findCriteriaKeyObject(search) {
  return {
    documentId,
    query: search.query,
    regex: search.regex,
    caseSensitive: search.caseSensitive,
    wholeCell: search.wholeCell,
  };
}

function cellsForReplaceScope(search) {
  if (
    findState.criteriaKey === findCriteriaKey(search) &&
    Array.isArray(findState.scopeCells) &&
    findState.scopeCells.length > 0
  ) {
    return findState.scopeCells;
  }
  return selectedCellsForSearch();
}

function selectedCellsForSearch(maxCells = 100_000) {
  const count = selectionCellCount();
  if (count <= 1 && !findInSelection) return [];
  if (count === 0) return [];
  if (count > maxCells) {
    status.textContent = tf("find.status.selectionLimit", { count: formatCount(count), limit: formatCount(maxCells) });
    return null;
  }
  const cells = [];
  for (const range of activeSelectionRanges()) {
    for (let row = range.startRow; row <= range.endRow; row++) {
      const sourceRow = rowSource.get(row);
      if (viewMapped && sourceRow == null) {
        status.textContent = t("find.status.loadSelection");
        return null;
      }
      const recordRow = sourceRow ?? row + dataStartRow;
      for (let col = range.startCol; col <= range.endCol; col++) {
        cells.push({ row: recordRow, col: sourceColFor(col) });
      }
    }
  }
  return cells;
}

function selectedSearchScopeKey() {
  return activeSelectionRanges().map((range) => {
    const startSourceRow = rowSource.get(range.startRow) ?? range.startRow + dataStartRow;
    const endSourceRow = rowSource.get(range.endRow) ?? range.endRow + dataStartRow;
    return `${Math.min(startSourceRow, endSourceRow)}:${Math.max(startSourceRow, endSourceRow)}:${sourceColFor(range.startCol)}:${sourceColFor(range.endCol)}`;
  }).join(";");
}

function emptyFindState() {
  return { key: "", criteriaKey: "", search: null, matches: [], total: 0, offset: 0, limit: findPageSize, truncated: false, index: 0, scopeCells: [], elapsedMs: 0, streamed: false };
}

function invalidateFindState() {
  findState = emptyFindState();
  clearReplacePreview();
  updateFindMatchCount();
  updateFindHighlightDataset();
}

function updateFindMatchCount() {
  const current = findState.matches.length ? findState.offset + findState.index + 1 : 0;
  const total = findState.total || 0;
  findMatchCount.textContent = `${current.toLocaleString()}/${total.toLocaleString()}`;
  const enabled = Boolean(documentId && total > 0);
  findPrevButton.disabled = !enabled;
  findNextButton.disabled = !enabled;
}

function updateFindHighlightDataset() {
  const match = findState.matches[findState.index];
  if (!match || !findState.search) {
    document.body.dataset.findHighlight = "";
    return;
  }
  const displayRow = Math.max(0, match.row - dataStartRow);
  const visibleCol = visibleColForSource(match.col);
  const value = visibleCol >= 0 ? valueAt(displayRow, visibleCol) : "";
  const serverRanges = normalizeServerFindRanges(value, match);
  const ranges = serverRanges.length ? serverRanges : findTextRanges(value, findState.search);
  document.body.dataset.findHighlight = ranges
    .map((range) => `${displayRow}:${visibleCol}:${range.start}-${range.end}`)
    .join(",");
}

function toggleFindOption(option) {
  if (option === "case") findCaseSensitive = !findCaseSensitive;
  if (option === "whole") findWholeCell = !findWholeCell;
  if (option === "regex") findRegex = !findRegex;
  if (option === "selection") findInSelection = !findInSelection;
  invalidateFindState();
  updateFindOptionButtons();
}

function updateFindOptionButtons() {
  setPressed(findCaseButton, findCaseSensitive);
  setPressed(findWholeButton, findWholeCell);
  setPressed(findRegexButton, findRegex);
  setPressed(findSelectionButton, findInSelection);
}

function setPressed(button, pressed) {
  button.classList.toggle("active", pressed);
  button.setAttribute("aria-pressed", String(pressed));
}

function setActiveFilter(filter) {
  activeFilter = filter;
  renderViewChips();
}

function clearActiveFilter() {
  activeFilter = null;
  renderViewChips();
}

function localizedFilterDisplay(filter) {
  const display = String(filter?.display || "");
  return display.replace(/^Filter:/i, `${t("view.filter")}:`);
}

function canonicalSortDisplay(keys = sortKeys) {
  return keys.length ? `Sort: ${formatSortKeys(keys)}` : "";
}

function localizedSortDisplay(keys = sortKeys) {
  return keys.length ? tf("view.chip.activeSort", { sort: formatSortKeysLocalized(keys) }) : "";
}

function localizedSavedSortDisplay(savedSort) {
  const keys = Array.isArray(savedSort?.keys) ? savedSort.keys : [];
  if (keys.length) return tf("view.chip.activeSort", { sort: formatStoredSortKeysLocalized(keys) });
  return String(savedSort?.display || "").replace(/^Sort:/i, `${t("view.sort")}:`);
}

function setChipButtonText(button, labelKey, titleKey, values = {}) {
  button.textContent = tf(labelKey, values);
  const title = tf(titleKey, values);
  button.title = title;
  button.setAttribute("aria-label", title);
}

function renderViewChips() {
  const saved = savedFilters();
  const recent = recentFilters();
  const savedSortsList = savedSorts();
  const activeSort = canonicalSortDisplay();
  const activeSortLabel = localizedSortDisplay();
  refreshFilterDropdown(saved, recent);
  viewChips.replaceChildren();
  if (!activeFilter && !activeSort && saved.length === 0 && savedSortsList.length === 0) {
    viewChips.hidden = true;
    document.body.dataset.activeFilter = "";
    document.body.dataset.activeFilterBuilderRows = "0";
    document.body.dataset.activeSort = "";
    document.body.dataset.savedFilterCount = "0";
    document.body.dataset.recentFilterCount = String(recent.length);
    document.body.dataset.savedSortCount = "0";
    return;
  }
  viewChips.hidden = false;
  document.body.dataset.activeFilter = activeFilter?.display ?? "";
  document.body.dataset.activeFilterBuilderRows = String(normalizeConditionBuilderRows(activeFilter?.builderRows).length);
  document.body.dataset.activeSort = activeSort;
  document.body.dataset.savedFilterCount = String(saved.length);
  document.body.dataset.recentFilterCount = String(recent.length);
  document.body.dataset.savedSortCount = String(savedSortsList.length);

  if (activeFilter) {
    const display = localizedFilterDisplay(activeFilter);
    const chip = document.createElement("div");
    chip.className = "view-chip";
    chip.title = display;

    const label = document.createElement("span");
    label.textContent = display;

    const sqlButton = document.createElement("button");
    sqlButton.type = "button";
    setChipButtonText(sqlButton, "view.chip.sql", "view.chip.openSql");
    sqlButton.addEventListener("click", openActiveFilterInSQL);

    const editButton = document.createElement("button");
    editButton.type = "button";
    setChipButtonText(editButton, "view.chip.edit", "view.chip.editBuilder");
    editButton.addEventListener("click", openActiveFilterInBuilder);

    const saveButton = document.createElement("button");
    saveButton.type = "button";
    setChipButtonText(saveButton, "view.chip.save", "view.chip.saveFilter");
    saveButton.addEventListener("click", saveActiveFilter);

    const clearButton = document.createElement("button");
    clearButton.type = "button";
    setChipButtonText(clearButton, "view.chip.clear", "view.chip.clearView");
    clearButton.addEventListener("click", clearView);

    chip.append(label, sqlButton);
    if (normalizeConditionBuilderRows(activeFilter.builderRows).length) {
      chip.append(editButton);
    }
    chip.append(saveButton, clearButton);
    viewChips.append(chip);
  }

  if (activeSort) {
    const chip = document.createElement("div");
    chip.className = "view-chip active-sort-chip";
    chip.title = activeSortLabel;

    const label = document.createElement("span");
    label.textContent = activeSortLabel;

    const saveButton = document.createElement("button");
    saveButton.type = "button";
    setChipButtonText(saveButton, "view.chip.save", "view.chip.saveSort");
    saveButton.addEventListener("click", saveActiveSort);

    const removeButton = document.createElement("button");
    removeButton.type = "button";
    setChipButtonText(removeButton, "view.chip.removeSortKey", "view.chip.removeLastSortKey");
    removeButton.disabled = sortKeys.length <= 1;
    removeButton.addEventListener("click", removeLastSortKey);

    const clearButton = document.createElement("button");
    clearButton.type = "button";
    setChipButtonText(clearButton, "view.chip.clear", "view.chip.clearView");
    clearButton.addEventListener("click", clearView);

    chip.append(label, saveButton, removeButton, clearButton);
    viewChips.append(chip);
  }

  for (const filter of saved) {
    const values = { name: filter.name };
    const chip = document.createElement("div");
    chip.className = "view-chip saved-filter-chip";
    chip.title = localizedFilterDisplay(filter);
    chip.dataset.filterId = filter.id;

    const label = document.createElement("span");
    label.textContent = tf("view.chip.savedFilter", values);

    const applyButton = document.createElement("button");
    applyButton.type = "button";
    setChipButtonText(applyButton, "view.chip.apply", "view.chip.applySavedFilter", values);
    applyButton.addEventListener("click", () => applySavedFilter(filter.id));

    const deleteButton = document.createElement("button");
    deleteButton.type = "button";
    setChipButtonText(deleteButton, "view.chip.delete", "view.chip.deleteSavedFilter", values);
    deleteButton.addEventListener("click", () => deleteSavedFilter(filter.id));

    chip.append(label, applyButton, deleteButton);
    viewChips.append(chip);
  }

  for (const savedSort of savedSortsList) {
    const values = { name: savedSort.name };
    const chip = document.createElement("div");
    chip.className = "view-chip saved-sort-chip";
    chip.title = localizedSavedSortDisplay(savedSort);
    chip.dataset.sortId = savedSort.id;

    const label = document.createElement("span");
    label.textContent = tf("view.chip.savedSort", values);

    const applyButton = document.createElement("button");
    applyButton.type = "button";
    setChipButtonText(applyButton, "view.chip.apply", "view.chip.applySavedSort", values);
    applyButton.addEventListener("click", () => applySavedSort(savedSort.id));

    const deleteButton = document.createElement("button");
    deleteButton.type = "button";
    setChipButtonText(deleteButton, "view.chip.delete", "view.chip.deleteSavedSort", values);
    deleteButton.addEventListener("click", () => deleteSavedSort(savedSort.id));

    chip.append(label, applyButton, deleteButton);
    viewChips.append(chip);
  }
}

function recentFilters() {
  if (filterHistoryLimit === 0) return [];
  try {
    const parsed = JSON.parse(localStorage.getItem(recentFiltersStorageKey) || "[]");
    if (!Array.isArray(parsed)) return [];
    return parsed
      .map(normalizeRecentFilter)
      .filter(Boolean)
      .slice(0, filterHistoryLimit);
  } catch {
    return [];
  }
}

function normalizeRecentFilter(filter) {
  if (!filter || typeof filter !== "object") return null;
  const body = normalizeRecentFilterBody(filter.body);
  const expression = typeof filter.expression === "string" && filter.expression.trim()
    ? filter.expression.trim()
    : (body?.expression || "");
  if (!body || !expression) return null;
  const display = typeof filter.display === "string" && filter.display.trim()
    ? filter.display.trim()
    : `Filter: where ${expression}`;
  return {
    id: typeof filter.id === "string" && filter.id ? filter.id : recentFilterID(body, expression),
    input: typeof filter.input === "string" && filter.input.trim() ? filter.input.trim() : recentFilterInput(body, expression),
    display,
    expression,
    body,
    builderRows: normalizeConditionBuilderRows(filter.builderRows),
    createdAt: typeof filter.createdAt === "string" && filter.createdAt ? filter.createdAt : new Date().toISOString(),
  };
}

function normalizeRecentFilterBody(body) {
  if (!body || typeof body !== "object") return null;
  if (body.mode === "sql" && typeof body.expression === "string" && body.expression.trim()) {
    return { mode: "sql", expression: stripWherePrefix(body.expression) };
  }
  if (typeof body.query === "string" && body.query.trim()) {
    const normalized = { query: body.query.trim() };
    if (Number.isInteger(body.column) && body.column >= 0) normalized.column = body.column;
    return normalized;
  }
  return null;
}

function recentFilterInput(body, expression) {
  if (body?.mode === "sql") return `where ${expression}`;
  return body?.query || `where ${expression}`;
}

function saveRecentFilters(filters) {
  const unique = [];
  const seen = new Set();
  for (const filter of filters) {
    const normalized = normalizeRecentFilter(filter);
    if (!normalized) continue;
    const signature = recentFilterSignature(normalized.body, normalized.expression);
    if (seen.has(signature)) continue;
    seen.add(signature);
    unique.push(normalized);
  }
  const next = filterHistoryLimit === 0 ? [] : unique.slice(0, filterHistoryLimit);
  localStorage.setItem(recentFiltersStorageKey, JSON.stringify(next));
  document.body.dataset.recentFilterCount = String(next.length);
  syncFilterHistoryControls();
  refreshFilterDropdown(savedFilters(), next);
  if (viewMode === "settings") renderSettingsPanel();
}

function rememberRecentFilter(body, descriptor) {
  if (filterHistoryLimit === 0 || !descriptor?.expression) return;
  const normalizedBody = normalizeRecentFilterBody(body);
  if (!normalizedBody) return;
  const expression = descriptor.expression;
  const saved = {
    id: recentFilterID(normalizedBody, expression),
    input: recentFilterInput(normalizedBody, expression),
    display: descriptor.display,
    expression,
    body: normalizedBody,
    builderRows: normalizeConditionBuilderRows(descriptor.builderRows),
    createdAt: new Date().toISOString(),
  };
  saveRecentFilters([saved, ...recentFilters()]);
  document.body.dataset.recentFilterPreview = `${saved.input}|${saved.expression}`;
}

function recentFilterID(body, expression) {
  return `recent-filter-${hashText(recentFilterSignature(body, expression))}`;
}

function recentFilterSignature(body, expression) {
  return `${body?.mode || "text"}:${body?.column ?? ""}:${body?.query || ""}:${expression}`;
}

function refreshFilterDropdown(saved = savedFilters(), recent = recentFilters()) {
  if (!recentFilterSelect) return;
  recentFilterSelect.replaceChildren();
  recentFilterSelect.append(filterDropdownOption("", t("view.filters")));
  if (recent.length) {
    const group = document.createElement("optgroup");
    group.label = t("view.filters.recent");
    for (const filter of recent) {
      group.append(filterDropdownOption(`recent:${filter.id}`, filter.input || filter.display));
    }
    recentFilterSelect.append(group);
  }
  if (saved.length) {
    const group = document.createElement("optgroup");
    group.label = t("view.filters.saved");
    for (const filter of saved) {
      group.append(filterDropdownOption(`saved:${filter.id}`, filter.name));
    }
    recentFilterSelect.append(group);
  }
  recentFilterSelect.value = "";
  recentFilterSelect.disabled = recent.length === 0 && saved.length === 0;
  document.body.dataset.filterDropdownRecentCount = String(recent.length);
  document.body.dataset.filterDropdownSavedCount = String(saved.length);
  document.body.dataset.filterDropdownOptions = [
    ...recent.map((filter) => `recent:${filter.input}`),
    ...saved.map((filter) => `saved:${filter.name}`),
  ].join("|");
}

function filterDropdownOption(value, label) {
  const option = document.createElement("option");
  option.value = value;
  option.textContent = label;
  return option;
}

async function applyFilterDropdownSelection() {
  const value = recentFilterSelect.value;
  recentFilterSelect.value = "";
  if (!value) return;
  const [kind, id] = value.split(":");
  if (kind === "saved") {
    await applySavedFilter(id);
    return;
  }
  if (kind !== "recent" || !documentId) return;
  const filter = recentFilters().find((item) => item.id === id);
  if (!filter) {
    status.textContent = t("view.status.recentFilterNotFound");
    refreshFilterDropdown();
    return;
  }
  commandInput.value = filter.input;
  document.body.dataset.recentFilterApplied = filter.input;
  await applyFilterRequest(
    { documentId, ...filter.body },
    `Recent filter ${savedFilterName(filter.display)}`,
    selected.col,
    { matched: "view.status.recentFilterMatched" },
    { builderRows: filter.builderRows },
  );
}

function savedFilters() {
  try {
    const parsed = JSON.parse(localStorage.getItem(savedFiltersStorageKey) || "[]");
    if (!Array.isArray(parsed)) return [];
    return parsed
      .filter((filter) => filter && typeof filter === "object" && typeof filter.expression === "string")
      .map((filter) => ({
        id: filter.id || savedFilterID(filter.expression),
        name: filter.name || savedFilterName(filter.display || filter.expression),
        display: filter.display || `Filter: where ${filter.expression}`,
        expression: filter.expression,
        builderRows: normalizeConditionBuilderRows(filter.builderRows),
        createdAt: filter.createdAt || new Date().toISOString(),
      }))
      .slice(0, savedFiltersLimit);
  } catch {
    return [];
  }
}

function saveSavedFilters(filters) {
  const unique = [];
  const seen = new Set();
  for (const filter of filters) {
    if (!filter?.expression || seen.has(filter.expression)) continue;
    seen.add(filter.expression);
    unique.push(filter);
  }
  const next = unique.slice(0, savedFiltersLimit);
  localStorage.setItem(savedFiltersStorageKey, JSON.stringify(next));
  document.body.dataset.savedFilterCount = String(next.length);
  refreshFilterDropdown(next, recentFilters());
  if (viewMode === "settings") renderSettingsPanel();
}

function saveActiveFilter() {
  if (!activeFilter?.expression) {
    status.textContent = t("view.status.saveFilterFirst");
    return;
  }
  const saved = {
    id: savedFilterID(activeFilter.expression),
    name: savedFilterName(activeFilter.display),
    display: activeFilter.display,
    expression: activeFilter.expression,
    builderRows: normalizeConditionBuilderRows(activeFilter.builderRows),
    createdAt: new Date().toISOString(),
  };
  const next = [saved, ...savedFilters().filter((filter) => filter.expression !== saved.expression)];
  saveSavedFilters(next);
  renderViewChips();
  document.body.dataset.savedFilterPreview = `${saved.name}|${saved.expression}`;
  status.textContent = tf("view.status.savedFilter", { name: saved.name });
}

async function applySavedFilter(filterID) {
  if (!documentId) return;
  const filter = savedFilters().find((item) => item.id === filterID);
  if (!filter) {
    status.textContent = t("view.status.savedFilterNotFound");
    renderViewChips();
    return;
  }
  commandInput.value = `where ${filter.expression}`;
  document.body.dataset.savedFilterApplied = filter.name;
  await applyFilterRequest(
    { documentId, mode: "sql", expression: filter.expression },
    `Saved filter ${filter.name}`,
    selected.col,
    { matched: "view.status.savedFilterMatched" },
    { builderRows: filter.builderRows },
  );
}

function deleteSavedFilter(filterID) {
  const before = savedFilters();
  const next = before.filter((filter) => filter.id !== filterID);
  saveSavedFilters(next);
  renderViewChips();
  status.textContent = before.length === next.length ? t("view.status.savedFilterNotFound") : t("view.status.deletedSavedFilter");
}

function savedSorts() {
  try {
    const parsed = JSON.parse(localStorage.getItem(savedSortsStorageKey) || "[]");
    if (!Array.isArray(parsed)) return [];
    return parsed
      .map((sort) => {
        const keys = Array.isArray(sort?.keys) ? sort.keys.map(normalizeStoredSortKey).filter(Boolean) : [];
        if (keys.length === 0) return null;
        const display = sort.display || `Sort: ${formatStoredSortKeys(keys)}`;
        return {
          id: sort.id || savedSortID(keys),
          name: sort.name || savedSortName(display),
          display,
          keys,
          stable: sort.stable !== false,
          createdAt: sort.createdAt || new Date().toISOString(),
        };
      })
      .filter(Boolean)
      .slice(0, savedSortsLimit);
  } catch {
    return [];
  }
}

function normalizeStoredSortKey(key) {
  if (!key || typeof key !== "object") return null;
  const column = Number(key.column);
  if (!Number.isInteger(column) || column < 0) return null;
  const mode = ["auto", "number", "date", "length", "text"].includes(key.mode)
    ? key.mode
    : (key.numeric ? "number" : "auto");
  return {
    column,
    descending: Boolean(key.descending),
    numeric: Boolean(key.numeric || mode === "number"),
    mode,
    name: typeof key.name === "string" ? key.name : "",
  };
}

function saveSavedSorts(sorts) {
  const unique = [];
  const seen = new Set();
  for (const sort of sorts) {
    const keys = Array.isArray(sort?.keys) ? sort.keys.map(normalizeStoredSortKey).filter(Boolean) : [];
    if (keys.length === 0) continue;
    const signature = savedSortSignature(keys);
    if (seen.has(signature)) continue;
    seen.add(signature);
    unique.push({ ...sort, keys });
  }
  const next = unique.slice(0, savedSortsLimit);
  localStorage.setItem(savedSortsStorageKey, JSON.stringify(next));
  document.body.dataset.savedSortCount = String(next.length);
  if (viewMode === "settings") renderSettingsPanel();
}

function saveActiveSort() {
  if (!sortKeys.length) {
    status.textContent = t("view.status.saveSortFirst");
    return;
  }
  const keys = sortKeys.map(storedSortKey);
  const display = `Sort: ${formatSortKeys(sortKeys)}`;
  const saved = {
    id: savedSortID(keys),
    name: savedSortName(display),
    display,
    keys,
    stable: sortStable,
    createdAt: new Date().toISOString(),
  };
  const next = [saved, ...savedSorts().filter((sort) => savedSortSignature(sort.keys) !== savedSortSignature(saved.keys))];
  saveSavedSorts(next);
  renderViewChips();
  document.body.dataset.savedSortPreview = `${saved.name}|${savedSortSignature(saved.keys)}`;
  status.textContent = tf("view.status.savedSort", { name: saved.name });
}

async function applySavedSort(sortID) {
  if (!documentId) return;
  const savedSort = savedSorts().find((item) => item.id === sortID);
  if (!savedSort) {
    status.textContent = t("view.status.savedSortNotFound");
    renderViewChips();
    return;
  }
  const resolvedKeys = resolveStoredSortKeys(savedSort.keys);
  if (resolvedKeys.length === 0) {
    status.textContent = t("view.status.savedSortNoColumns");
    renderViewChips();
    return;
  }
  const primaryVisibleCol = Math.max(0, visibleColForSource(resolvedKeys[0].column));
  document.body.dataset.savedSortApplied = savedSort.name;
  await applySortKeys(
    resolvedKeys,
    primaryVisibleCol,
    tf("view.status.savedSortApplied", { name: savedSort.name, sort: formatSortKeysLocalized(resolvedKeys) }),
    { stable: savedSort.stable !== false },
  );
}

function deleteSavedSort(sortID) {
  const before = savedSorts();
  const next = before.filter((sort) => sort.id !== sortID);
  saveSavedSorts(next);
  renderViewChips();
  status.textContent = before.length === next.length ? t("view.status.savedSortNotFound") : t("view.status.deletedSavedSort");
}

function exportSavedViews() {
  const filters = savedFilters();
  const sorts = savedSorts();
  const payload = {
    kind: savedViewsExportKind,
    version: 1,
    exportedAt: new Date().toISOString(),
    filters,
    sorts,
  };
  const text = `${JSON.stringify(payload, null, 2)}\n`;
  window.__opencsvLastSavedViewExportText = text;
  document.body.dataset.savedViewExport = `${filters.length}:${sorts.length}`;
  document.body.dataset.savedViewExportPreview = text.slice(0, 500);
  triggerDownload(text, "opencsv-saved-views.json", "application/json;charset=utf-8");
  status.textContent = tf("view.status.exportedSavedViews", {
    filters: formatCount(filters.length),
    sorts: formatCount(sorts.length),
  });
}

function showSavedViewImportPanel() {
  openSettingsPanel();
  savedViewImportPanel.hidden = false;
  savedViewImportText.value = "";
  document.body.dataset.savedViewImportPanel = "open";
  syncSavedViewImportControls();
  savedViewImportText.focus();
}

function cancelSavedViewImport() {
  savedViewImportPanel.hidden = true;
  savedViewImportText.value = "";
  document.body.dataset.savedViewImportPanel = "closed";
  syncSavedViewImportControls();
}

function syncSavedViewImportControls() {
  if (!savedViewImportText) return;
  applySavedViewImportButton.disabled = savedViewImportText.value.trim() === "";
}

function applySavedViewImport() {
  try {
    const imported = parseSavedViewImport(savedViewImportText.value);
    saveSavedFilters([...imported.filters, ...savedFilters()]);
    saveSavedSorts([...imported.sorts, ...savedSorts()]);
    savedViewImportPanel.hidden = true;
    savedViewImportText.value = "";
    document.body.dataset.savedViewImport = `${imported.filters.length}:${imported.sorts.length}`;
    document.body.dataset.savedViewImportPanel = "closed";
    renderViewChips();
    renderSettingsPanel();
    status.textContent = tf("view.status.importedSavedViews", {
      filters: formatCount(imported.filters.length),
      sorts: formatCount(imported.sorts.length),
    });
  } catch (error) {
    status.textContent = tf("view.status.savedViewImportFailed", { message: error.message });
    document.body.dataset.savedViewImport = "error";
  }
  syncSavedViewImportControls();
}

function parseSavedViewImport(text) {
  let parsed;
  try {
    parsed = JSON.parse(text);
  } catch {
    throw new Error("invalid JSON");
  }
  const rawFilters = Array.isArray(parsed) ? parsed : parsed?.filters;
  const rawSorts = Array.isArray(parsed?.sorts) ? parsed.sorts : [];
  if (!Array.isArray(rawFilters) && !Array.isArray(rawSorts)) {
    throw new Error("missing filters or sorts array");
  }
  return {
    filters: (Array.isArray(rawFilters) ? rawFilters : [])
      .map(normalizeSavedFilterImport)
      .filter(Boolean),
    sorts: (Array.isArray(rawSorts) ? rawSorts : [])
      .map(normalizeSavedSortImport)
      .filter(Boolean),
  };
}

function normalizeSavedFilterImport(filter) {
  if (!filter || typeof filter !== "object" || typeof filter.expression !== "string") return null;
  const expression = filter.expression.trim().replace(/^where\s+/i, "");
  if (!expression) return null;
  const display = typeof filter.display === "string" && filter.display.trim()
    ? filter.display.trim()
    : `Filter: where ${expression}`;
  return {
    id: typeof filter.id === "string" && filter.id ? filter.id : savedFilterID(expression),
    name: typeof filter.name === "string" && filter.name.trim() ? filter.name.trim() : savedFilterName(display),
    display,
    expression,
    builderRows: normalizeConditionBuilderRows(filter.builderRows),
    createdAt: typeof filter.createdAt === "string" && filter.createdAt ? filter.createdAt : new Date().toISOString(),
  };
}

function normalizeSavedSortImport(sort) {
  if (!sort || typeof sort !== "object") return null;
  const keys = Array.isArray(sort.keys) ? sort.keys.map(normalizeStoredSortKey).filter(Boolean) : [];
  if (keys.length === 0) return null;
  const display = typeof sort.display === "string" && sort.display.trim()
    ? sort.display.trim()
    : `Sort: ${formatStoredSortKeys(keys)}`;
  return {
    id: typeof sort.id === "string" && sort.id ? sort.id : savedSortID(keys),
    name: typeof sort.name === "string" && sort.name.trim() ? sort.name.trim() : savedSortName(display),
    display,
    keys,
    stable: sort.stable !== false,
    createdAt: typeof sort.createdAt === "string" && sort.createdAt ? sort.createdAt : new Date().toISOString(),
  };
}

function storedSortKey(key) {
  return {
    column: key.column,
    descending: Boolean(key.descending),
    numeric: Boolean(key.numeric),
    mode: key.mode || (key.numeric ? "number" : "auto"),
    name: columnNameForSource(key.column),
  };
}

function resolveStoredSortKeys(keys) {
  return keys
    .map((key) => {
      const sourceCol = resolveStoredSortColumn(key);
      if (sourceCol < 0 || sourceCol >= allColumns.length) return null;
      return {
        column: sourceCol,
        descending: Boolean(key.descending),
        numeric: Boolean(key.numeric || key.mode === "number"),
        mode: key.mode || (key.numeric ? "number" : "auto"),
      };
    })
    .filter(Boolean);
}

function resolveStoredSortColumn(key) {
  if (key.name) {
    const exact = allColumns.find((column) => column.name === key.name);
    if (exact) return exact.sourceCol;
    const lower = key.name.toLowerCase();
    const folded = allColumns.find((column) => column.name.toLowerCase() === lower);
    if (folded) return folded.sourceCol;
  }
  return Number.isInteger(key.column) ? key.column : -1;
}

function savedSortID(keys) {
  return `sort-${hashText(savedSortSignature(keys))}`;
}

function savedSortSignature(keys) {
  return keys
    .map((key) => `${key.column}:${key.name || ""}:${key.mode || "auto"}:${key.descending ? "desc" : "asc"}`)
    .join("|");
}

function savedSortName(display) {
  const base = String(display || "")
    .replace(/^Sort:\s*/i, "")
    .trim();
  return base.length > 48 ? `${base.slice(0, 45)}...` : base || "Saved sort";
}

function formatStoredSortKeys(keys) {
  return keys
    .map((key) => {
      const label = key.name || columnNameForSource(key.column);
      const modeLabel = key.mode && key.mode !== "auto" ? ` ${key.mode}` : "";
      return `${label}${modeLabel} ${key.descending ? "descending" : "ascending"}`;
    })
    .join(", then ");
}

function formatStoredSortKeysLocalized(keys) {
  return keys
    .map((key) => {
      const label = key.name || columnNameForSource(key.column);
      const mode = sortModeLabel(key.mode);
      const modeLabel = mode ? ` ${mode}` : "";
      const direction = key.descending ? t("sort.direction.descending") : t("sort.direction.ascending");
      return `${label}${modeLabel} ${direction}`;
    })
    .join(t("sort.join.then"));
}

function savedFilterID(expression) {
  return `filter-${hashText(expression)}`;
}

function savedFilterName(display) {
  const base = String(display || "")
    .replace(/^Filter:\s*/i, "")
    .replace(/^where\s+/i, "")
    .trim();
  return base.length > 42 ? `${base.slice(0, 39)}...` : base || "Saved filter";
}

function hashText(value) {
  let hash = 2166136261;
  for (const char of String(value)) {
    hash ^= char.charCodeAt(0);
    hash = Math.imul(hash, 16777619);
  }
  return (hash >>> 0).toString(36);
}

function openActiveFilterInSQL() {
  if (!activeFilter?.sqlQuery) return;
  sqlInput.value = activeFilter.sqlQuery;
  document.body.dataset.sqlFilterPreview = activeFilter.sqlQuery;
  setViewMode("sql");
  status.textContent = t("view.status.filterToSql");
}

function openActiveFilterInBuilder() {
  const builderRows = normalizeConditionBuilderRows(activeFilter?.builderRows);
  if (!builderRows.length) return;
  setViewMode("grid");
  openConditionBuilder();
  setConditionBuilderRows(builderRows);
  document.body.dataset.conditionBuilderReopened = `${builderRows.length}:${activeFilter.expression || ""}`;
  status.textContent = t("view.status.filterToBuilder");
  conditionValue.focus();
  if (conditionOperatorNeedsValue(conditionOperator.value)) conditionValue.select();
}

function filterDescriptorFromRequest(body, label, options = {}) {
  if (body.mode === "sql" && body.expression) {
    return filterDescriptor(`where ${body.expression}`, body.expression, options);
  }
  if (body.query && Number.isInteger(body.column)) {
    const expression = `${sqlColumnIdentifier(body.column)} LIKE ${sqlStringLiteral(`%${body.query}%`)}`;
    return filterDescriptor(`${columnNameForSource(body.column)} contains ${body.query}`, expression, options);
  }
  if (body.query) {
    const expression = containsAnyColumnExpression(body.query);
    return filterDescriptor(`contains ${body.query}`, expression, options);
  }
  return null;
}

function valueFilterDescriptor(sourceCol, values, includeBlank) {
  const parts = [];
  if (values.length === 1) {
    parts.push(`${sqlColumnIdentifier(sourceCol)} = ${sqlStringLiteral(values[0])}`);
  } else if (values.length > 1) {
    parts.push(`${sqlColumnIdentifier(sourceCol)} IN (${values.map(sqlStringLiteral).join(", ")})`);
  }
  if (includeBlank) {
    parts.push(`${sqlColumnIdentifier(sourceCol)} = ''`);
  }
  const expression = parts.join(" OR ");
  const count = values.length + (includeBlank ? 1 : 0);
  return filterDescriptor(`${columnNameForSource(sourceCol)}: ${count} values`, expression);
}

function filterDescriptor(display, expression, options = {}) {
  if (!expression) return null;
  const descriptor = {
    display: `Filter: ${display}`,
    expression,
    sqlQuery: `SELECT * FROM current WHERE ${expression} LIMIT 500`,
  };
  const builderRows = normalizeConditionBuilderRows(options.builderRows);
  if (builderRows.length) descriptor.builderRows = builderRows;
  return descriptor;
}

function containsAnyColumnExpression(query) {
  const sourceColumns = (allColumns.length ? allColumns : columns)
    .map((column) => column.sourceCol)
    .filter((sourceCol, index, values) => values.indexOf(sourceCol) === index);
  const columnsToSearch = sourceColumns.length ? sourceColumns : [0];
  return columnsToSearch
    .map((sourceCol) => `${sqlColumnIdentifier(sourceCol)} LIKE ${sqlStringLiteral(`%${query}%`)}`)
    .join(" OR ");
}

async function filterCurrentQuery() {
  if (!documentId) return;
  const query = commandInput.value.trim();
  if (!query) {
    await clearView();
    return;
  }
  const useSQL = looksLikeSQLFilter(query);
  const expression = stripWherePrefix(query);
  await applyFilterRequest(
    useSQL
      ? { documentId, mode: "sql", expression }
      : { documentId, query },
    useSQL ? "SQL filter" : "Filter",
    0,
    {
      failed: useSQL ? "view.status.sqlFilterFailed" : "view.status.filterFailed",
      matched: useSQL ? "view.status.sqlFilterMatched" : "view.status.filterMatched",
    },
  );
}

async function applyFilterRequest(body, label, preserveCol = selected.col, statusKeys = {}, descriptorOptions = {}) {
  const response = await fetch("/api/filter", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const payload = await response.json();
  if (!response.ok) {
    const message = payload.error ?? "unknown error";
    status.textContent = statusKeys.failed
      ? tf(statusKeys.failed, { message, ...(statusKeys.values || {}) })
      : `${label} failed: ${message}`;
    return false;
  }
  viewMapped = true;
  document.body.dataset.filterStreamed = String(Boolean(payload.streamed));
  sortDescending = false;
  sortKeys = [];
  const descriptor = filterDescriptorFromRequest(body, label, descriptorOptions);
  setActiveFilter(descriptor);
  rememberRecentFilter(body, descriptor);
  rows = new Map();
  rowSource = new Map();
  selected = { row: 0, col: Math.max(0, Math.min(columns.length - 1, preserveCol)) };
  selectionAnchor = { ...selected };
  scrollY = 0;
  wheelRemainderY = 0;
  aggregateRows = [];
  sourceRowCount = payload.totalRows ?? 0;
  syncTotalRows();
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(0);
  const count = sourceRowCount.toLocaleString(localeLanguageTag());
  status.textContent = statusKeys.matched
    ? tf(statusKeys.matched, { count, ...(statusKeys.values || {}) })
    : `${label} matched ${count} rows`;
  return true;
}

function looksLikeSQLFilter(query) {
  const value = query.trim();
  return (
    /^where\s+/i.test(value) ||
    /\b(and|or|not|like|between|is|cast|try_cast|null|date|time|timestamp)\b/i.test(value) ||
    /(?:[<>=!]=|<>|[<>])/.test(value)
  );
}

function stripWherePrefix(query) {
  return query.trim().replace(/^where\s+/i, "");
}

async function sortSelectedColumn() {
  if (!documentId || columns.length === 0) return;
  await sortColumn(selected.col, !sortDescending, false);
}

async function sortActiveColumnForCommand(commandID) {
  const spec = columnSortCommands.get(commandID);
  if (!spec) return;
  const column = activeCommandColumn();
  closeColumnMenu();
  await sortColumn(column, spec.descending, spec.append, spec.mode);
}

async function sortColumn(column, descending, append, modeOverride = "") {
  if (!documentId || columns.length === 0 || !columns[column]) return;
  const nextKeys = nextSortKeys(column, descending, append, modeOverride);
  await applySortKeys(nextKeys, column);
}

async function removeLastSortKey() {
  if (!documentId || sortKeys.length <= 1) {
    status.textContent = t("view.status.removeSortKeyFirst");
    return false;
  }
  const removed = sortKeys[sortKeys.length - 1];
  const nextKeys = sortKeys.slice(0, -1).map((key) => ({ ...key }));
  const snapshot = {
    filter: activeFilter ? { ...activeFilter } : null,
    sortKeys: nextKeys,
    sortStable,
    selectedSourceCol: nextKeys[nextKeys.length - 1]?.column ?? nextKeys[0]?.column ?? selected.col,
  };
  const reset = await resetViewMappingForReplay(snapshot.selectedSourceCol);
  if (!reset) return false;
  const restored = await restoreViewSnapshot(snapshot);
  if (!restored.includes("sort")) return false;
  const removedLabel = formatSortKeysLocalized([removed]);
  const remainingLabel = formatSortKeysLocalized(nextKeys);
  document.body.dataset.sortRemovedKey = `${formatSortKeys([removed])}->${formatSortKeys(nextKeys)}`;
  document.body.dataset.sortKeyCount = String(nextKeys.length);
  status.textContent = tf("view.status.removedSortKey", {
    key: removedLabel,
    sort: remainingLabel,
  });
  renderViewChips();
  return true;
}

async function resetViewMappingForReplay(selectedSourceCol = selected.col) {
  const response = await fetch("/api/view-clear", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = tf("view.status.clearFailed", { message: payload.error ?? "unknown error" });
    return false;
  }
  viewMapped = false;
  sortDescending = false;
  sortKeys = [];
  rows = new Map();
  rowSource = new Map();
  aggregateRows = [];
  sourceRowCount = Math.max(0, (payload.totalRows ?? 0) - dataStartRow);
  syncTotalRows();
  windowStart = 0;
  windowEnd = 0;
  const visibleCol = visibleColForSource(selectedSourceCol);
  selected = { row: 0, col: visibleCol >= 0 ? visibleCol : Math.max(0, Math.min(columns.length - 1, selected.col)) };
  selectionAnchor = { ...selected };
  scrollY = 0;
  wheelRemainderY = 0;
  await maybeLoadWindow(0);
  return true;
}

async function applySortKeys(nextKeys, selectedColumn, successMessage = "", options = {}) {
  if (!documentId || columns.length === 0 || nextKeys.length === 0) return false;
  const primaryKey = nextKeys[0];
  const stable = options.stable ?? sortStable;
  const response = await fetch("/api/sort", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      column: primaryKey.column,
      descending: primaryKey.descending,
      numeric: primaryKey.numeric,
      mode: primaryKey.mode,
      keys: nextKeys,
      hasHeader: metadata?.dialect?.hasHeader ?? false,
      stable,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = tf("view.status.sortFailed", { message: payload.error ?? "unknown error" });
    return false;
  }
  viewMapped = true;
  sortKeys = nextKeys;
  sortStable = payload.stable !== false;
  sortDescending = Boolean(primaryKey.descending);
  rows = new Map();
  rowSource = new Map();
  selected = { row: 0, col: Math.max(0, Math.min(columns.length - 1, selectedColumn)) };
  selectionAnchor = { ...selected };
  scrollY = 0;
  wheelRemainderY = 0;
  aggregateRows = [];
  sourceRowCount = payload.totalRows ?? 0;
  syncTotalRows();
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(0);
  renderViewChips();
  document.body.dataset.sortStable = String(sortStable);
  document.body.dataset.sortKeyCount = String(sortKeys.length);
  status.textContent = successMessage || tf("view.status.sorted", { sort: formatSortKeysLocalized(sortKeys) });
  return true;
}

function nextSortKeys(column, descending, append, modeOverride = "") {
  const key = sortKeyForColumn(column, descending, modeOverride);
  if (!append || sortKeys.length === 0) return [key];
  return [...sortKeys.filter((item) => item.column !== key.column), key];
}

function sortKeyForColumn(column, descending, modeOverride = "") {
  const type = columns[column]?.type ?? "string";
  const mode = modeOverride || sortModeForColumn(column, type);
  return {
    column: sourceColFor(column),
    descending,
    numeric: mode === "number",
    mode,
  };
}

function formatSortKeys(keys) {
  return keys
    .map((key) => {
      const modeLabel = key.mode && key.mode !== "auto" ? ` ${key.mode}` : "";
      return `${columnNameForSource(key.column)}${modeLabel} ${key.descending ? "descending" : "ascending"}`;
    })
    .join(", then ");
}

function sortModeLabel(mode) {
  if (!mode || mode === "auto") return "";
  return translationOrEmpty(`sort.mode.${mode}`) || mode;
}

function formatSortKeysLocalized(keys) {
  return keys
    .map((key) => {
      const mode = sortModeLabel(key.mode);
      const modeLabel = mode ? ` ${mode}` : "";
      const direction = key.descending ? t("sort.direction.descending") : t("sort.direction.ascending");
      return `${columnNameForSource(key.column)}${modeLabel} ${direction}`;
    })
    .join(t("sort.join.then"));
}

function sortModeForColumn(column, type) {
  if (type === "integer" || type === "decimal" || columnLooksNumeric(column)) return "number";
  if (type === "date" || columnLooksDate(column)) return "date";
  return "auto";
}

function columnLooksNumeric(column) {
  let inspected = 0;
  let numeric = 0;
  for (let row = 0; row < Math.min(totalRows, 80); row++) {
    const value = valueAt(row, column);
    if (value == null || String(value).trim() === "") continue;
    inspected++;
    if (Number.isFinite(Number(normalizeNumericForSort(value)))) numeric++;
  }
  return inspected > 0 && numeric === inspected;
}

function columnLooksDate(column) {
  let inspected = 0;
  let dates = 0;
  for (let row = 0; row < Math.min(totalRows, 80); row++) {
    const value = valueAt(row, column);
    const normalized = normalizeNumericForSort(value);
    if (!normalized) continue;
    inspected++;
    if (Number.isFinite(Date.parse(normalized))) dates++;
  }
  return inspected > 0 && dates === inspected;
}

function normalizeNumericForSort(value) {
  return String(value ?? "")
    .trim()
    .replace(/[０-９]/g, (digit) => String.fromCharCode(digit.charCodeAt(0) - 0xff10 + 0x30))
    .replace(/．/g, ".")
    .replace(/－/g, "-")
    .replace(/＋/g, "+")
    .replace(/[，,]/g, "");
}

async function filterActiveColumnForCommand(commandID) {
  const spec = columnFilterCommands.get(commandID);
  if (!spec) return;
  if (spec.kind === "value") {
    await filterColumnBySelectedValue();
    return;
  }
  await filterColumnBlank(spec.kind === "blank");
}

async function filterColumnBySelectedValue(column = activeCommandColumn()) {
  closeColumnMenu();
  if (!documentId || !columns[column]) return;
  const value = valueAt(selected.row, column);
  if (value === "") {
    status.textContent = t("view.status.selectedBlank");
    return;
  }
  const columnName = columns[column].name;
  commandInput.value = value;
  await applyFilterRequest(
    { documentId, query: value, column: sourceColFor(column) },
    `Filtered ${columnName}`,
    column,
    {
      failed: "view.status.filterFailed",
      matched: "view.status.columnFilterMatched",
      values: { column: columnName },
    },
  );
}

async function filterSelectedCellValue() {
  if (!documentId || !columns[selected.col]) return;
  await filterColumnBySelectedValue(selected.col);
}

async function filterColumnBlank(showBlank, column = activeCommandColumn()) {
  closeColumnMenu();
  if (!documentId || !columns[column]) return;
  const expression = `c${sourceColFor(column) + 1} ${showBlank ? "=" : "<>"} ''`;
  const columnName = columns[column].name;
  commandInput.value = `where ${expression}`;
  await applyFilterRequest(
    { documentId, mode: "sql", expression },
    showBlank ? `Blank ${columnName}` : `Non-blank ${columnName}`,
    column,
    {
      failed: "view.status.sqlFilterFailed",
      matched: showBlank ? "view.status.blankFilterMatched" : "view.status.nonblankFilterMatched",
      values: { column: columnName },
    },
  );
}

async function applyColumnValueChecklist() {
  const column = columnMenuColumn;
  if (!documentId || !columns[column]) return;
  const checked = Array.from(columnMenuValues.querySelectorAll("input:checked"));
  if (!checked.length) {
    status.textContent = t("view.status.chooseValueFilter");
    return;
  }
  const values = [];
  let includeBlank = false;
  for (const input of checked) {
    if (input.dataset.blank === "true") {
      includeBlank = true;
    } else {
      values.push(input.value);
    }
  }
  closeColumnMenu();
  commandInput.value = `${columns[column].name} values`;
  const sourceCol = sourceColFor(column);
  const response = await fetch("/api/filter-values", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, column: sourceCol, values, includeBlank }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = tf("view.status.valueFilterFailed", { message: payload.error ?? "unknown error" });
    return;
  }
  viewMapped = true;
  document.body.dataset.valueFilterStreamed = String(Boolean(payload.streamed));
  sortDescending = false;
  sortKeys = [];
  const descriptor = valueFilterDescriptor(sourceCol, values, includeBlank);
  setActiveFilter(descriptor);
  rememberRecentFilter({ mode: "sql", expression: descriptor.expression }, descriptor);
  rows = new Map();
  rowSource = new Map();
  selected = { row: 0, col: column };
  selectionAnchor = { ...selected };
  scrollY = 0;
  wheelRemainderY = 0;
  aggregateRows = [];
  sourceRowCount = payload.totalRows ?? 0;
  syncTotalRows();
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(0);
  status.textContent = tf("view.status.valueFilterMatched", {
    column: columns[column].name,
    count: sourceRowCount.toLocaleString(localeLanguageTag()),
  });
}

async function clearView() {
  if (!documentId) return;
  const response = await fetch("/api/view-clear", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = tf("view.status.clearFailed", { message: payload.error ?? "unknown error" });
    return;
  }
  viewMapped = false;
  sortDescending = false;
  sortKeys = [];
  clearActiveFilter();
  hiddenSourceColumns = new Set();
  applyHiddenColumns();
  saveHiddenColumnsState();
  rows = new Map();
  rowSource = new Map();
  aggregateRows = [];
  sourceRowCount = Math.max(0, (payload.totalRows ?? 0) - dataStartRow);
  syncTotalRows();
  windowStart = 0;
  windowEnd = 0;
  selected = { row: 0, col: Math.max(0, Math.min(columns.length - 1, selected.col)) };
  selectionAnchor = { ...selected };
  scrollY = 0;
  wheelRemainderY = 0;
  await maybeLoadWindow(0);
  status.textContent = t("view.status.cleared");
}

async function undoCurrentDocument() {
  if (!documentId) return;
  if (!ensureEditable("Undo")) return;
  const response = await fetch("/api/undo", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Undo failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  applyMutationPayload(payload, `Undo: ${formatCount(payload.totalRows ?? 0)} rows`);
}

async function redoCurrentDocument() {
  if (!documentId) return;
  if (!ensureEditable("Redo")) return;
  const response = await fetch("/api/redo", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Redo failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  applyMutationPayload(payload, `Redo: ${formatCount(payload.totalRows ?? 0)} rows`);
}

function cutSelectedRow() {
  if (!documentId) return;
  if (!ensureEditable("Cut Row")) return;
  const rowSelection = selectedRowsForStructure();
  if (!rowSelection) return;
  setStructureClipboard("row", rowSelection.sourceRows, true);
  const count = rowSelection.sourceRows.length;
  status.textContent = `Cut ${formatCount(count)} ${count === 1 ? "row" : "rows"} for paste.`;
}

async function pasteRowFromClipboard(options = {}) {
  if (!documentId) return;
  if (!ensureEditable("Paste Row")) return;
  if (aggregateRowAt(selected.row)) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  const clipboard = structureClipboardFor("row");
  if (!clipboard) {
    status.textContent = "Copy or cut a row first.";
    updateCommandState();
    return;
  }
  const placement = options.placement || "before-active";
  const targetRow = placement === "after-selection" ? selectedRowPasteAfterTarget() : selectedSourceRow();
  if (targetRow == null) return;
  const sources = clipboard.sources ?? [clipboard.source];
  const action = clipboard.cut
    ? placement === "after-selection" ? "Move row below" : "Move row"
    : placement === "after-selection" ? "Paste row below" : "Paste row";
  status.textContent = `${action}...`;
  const response = await fetch("/api/row/paste", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      sourceRow: sources[0],
      sourceRows: sources,
      targetRow,
      cut: clipboard.cut,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${action} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const changed = payload.changed ?? 0;
  const count = sources.length;
  const label = `${formatCount(count)} ${count === 1 ? "row" : "rows"}`;
  const message = changed > 0
    ? (clipboard.cut ? `Moved ${label}` : `Pasted ${label}${placement === "after-selection" ? " below" : ""}`)
    : "Row paste did not change the document";
  applyMutationPayload(payload, message);
  document.body.dataset.structurePaste = `row:${clipboard.cut ? "cut" : "copy"}:${placement}:${sources.join(",")}->${targetRow}`;
  if (clipboard.cut) clearStructureClipboard();
}

async function insertSelectedRow() {
  if (aggregateRowAt(selected.row)) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  await mutateStructure("/api/row/insert", { row: selectedSourceRow() }, "Inserted row");
}

async function deleteSelectedRow() {
  if (removeSelectedAggregateRow()) return;
  await mutateStructure("/api/row/delete", { row: selectedSourceRow() }, "Deleted row");
}

async function hideSelectedRow() {
  if (!documentId) return;
  if (!ensureEditable("Hide Row")) return;
  if (aggregateRowAt(selected.row)) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  status.textContent = "Hiding row...";
  const response = await fetch("/api/row/hide", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, row: selectedSourceRow() }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Hide row failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  viewMapped = true;
  rows = new Map();
  rowSource = new Map();
  aggregateRows = [];
  sourceRowCount = payload.totalRows ?? Math.max(0, sourceRowCount - (payload.changed ?? 0));
  syncTotalRows();
  selected = { row: Math.max(0, Math.min(totalRows - 1, selected.row)), col: selected.col };
  selectionAnchor = { ...selected };
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(selected.row);
  status.textContent = payload.changed ? "Row hidden" : "Row was already hidden";
}

async function showAllRows() {
  if (!documentId) return;
  const previous = sourceRowCount;
  const response = await fetch("/api/view-clear", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Show all rows failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  viewMapped = false;
  sortDescending = false;
  sortKeys = [];
  clearActiveFilter();
  rows = new Map();
  rowSource = new Map();
  aggregateRows = [];
  sourceRowCount = Math.max(0, (payload.totalRows ?? 0) - dataStartRow);
  syncTotalRows();
  windowStart = 0;
  windowEnd = 0;
  selected = { row: Math.max(0, Math.min(sourceRowCount - 1, selected.row)), col: selected.col };
  selectionAnchor = { ...selected };
  scrollY = 0;
  wheelRemainderY = 0;
  await maybeLoadWindow(0);
  const restored = Math.max(0, sourceRowCount - previous);
  document.body.dataset.rowShowAll = `${previous}->${sourceRowCount}`;
  document.body.dataset.rowShowAllRestored = String(restored);
  status.textContent = restored > 0
    ? `Restored ${formatCount(restored)} hidden ${restored === 1 ? "row" : "rows"}`
    : "All rows already shown";
}

async function moveSelectedRow(delta) {
  if (!documentId) return;
  if (!ensureEditable("Move Row")) return;
  if (aggregateRowAt(selected.row)) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  const previousRow = selected.row;
  const label = delta < 0 ? "Moved row up" : "Moved row down";
  status.textContent = `${label}...`;
  const response = await fetch("/api/row/move", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, row: selectedSourceRow(), delta }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const changed = payload.changed ?? 0;
  applyMutationPayload(payload, changed > 0 ? label : "Row is already at the edge");
  if (changed > 0) {
    selected = { row: Math.max(0, Math.min(totalRows - 1, previousRow + delta)), col: selected.col };
    selectionAnchor = { ...selected };
    draw();
    status.textContent = label;
  }
}

async function moveRowToDisplay(sourceDisplayRow, targetDisplayRow, source = "command") {
  if (!documentId) return;
  if (!ensureEditable("Move Row")) return;
  if (viewMapped) {
    document.body.dataset.rowDragReorder = `${sourceDisplayRow}->${targetDisplayRow}:blocked`;
    status.textContent = "Clear filter or sort before dragging rows.";
    draw();
    return;
  }
  if (aggregateRowAt(sourceDisplayRow) || aggregateRowAt(targetDisplayRow)) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  const sourceRow = sourceRowForDisplayRow(sourceDisplayRow, "Move Row");
  const targetRow = sourceRowForDisplayRow(targetDisplayRow, "Move Row");
  if (sourceRow == null || targetRow == null) return;
  const delta = targetRow - sourceRow;
  const label = `Moved row ${sourceDisplayRow + 1} to ${targetDisplayRow + 1}`;
  if (delta === 0) {
    document.body.dataset.rowDragReorder = `${sourceDisplayRow}->${targetDisplayRow}:cancel`;
    status.textContent = "Row is already at the target.";
    draw();
    return;
  }
  status.textContent = `${label}...`;
  const response = await fetch("/api/row/move", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, row: sourceRow, delta }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const changed = payload.changed ?? 0;
  applyMutationPayload(payload, changed > 0 ? label : "Row is already at the edge");
  if (changed > 0) {
    const movedRow = Math.max(0, Math.min(sourceRowCount - 1, targetDisplayRow));
    selected = { row: movedRow, col: 0 };
    selectionAnchor = { row: movedRow, col: Math.max(0, columns.length - 1) };
    document.body.dataset.rowDragReorder = `${sourceDisplayRow}->${movedRow}:moved`;
    document.body.dataset.rowDragSource = source;
    document.body.dataset.rowDragPhysical = `${sourceRow}->${targetRow}`;
    draw();
    status.textContent = label;
  }
}

async function insertSelectedColumn() {
  await mutateStructure("/api/column/insert", { col: sourceColFor(selected.col) }, "Inserted column");
}

async function deleteSelectedColumn() {
  await mutateStructure("/api/column/delete", { col: sourceColFor(selected.col) }, "Deleted column");
}

function cutSelectedColumn() {
  if (!documentId) return;
  if (!ensureEditable("Cut Column")) return;
  if (!columns[selected.col]) {
    status.textContent = "Select a column first.";
    return;
  }
  const sourceCols = selectedColumnsForStructure();
  if (!sourceCols) return;
  setStructureClipboard("column", sourceCols, true);
  const count = sourceCols.length;
  status.textContent = `Cut ${formatCount(count)} ${count === 1 ? "column" : "columns"} for paste.`;
}

async function pasteColumnFromClipboard(options = {}) {
  if (!documentId) return;
  if (!ensureEditable("Paste Column")) return;
  if (!columns[selected.col]) {
    status.textContent = "Select a column first.";
    return;
  }
  const clipboard = structureClipboardFor("column");
  if (!clipboard) {
    status.textContent = "Copy or cut a column first.";
    updateCommandState();
    return;
  }
  const placement = options.placement || "before-active";
  const targetCol = placement === "after-selection" ? selectedColumnPasteRightTarget() : sourceColFor(selected.col);
  if (targetCol == null) return;
  const sources = clipboard.sources ?? [clipboard.source];
  const action = clipboard.cut
    ? placement === "after-selection" ? "Move column right" : "Move column"
    : placement === "after-selection" ? "Paste column right" : "Paste column";
  status.textContent = `${action}...`;
  const response = await fetch("/api/column/paste", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      sourceCol: sources[0],
      sourceCols: sources,
      targetCol,
      cut: clipboard.cut,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${action} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const changed = payload.changed ?? 0;
  const count = sources.length;
  const label = `${formatCount(count)} ${count === 1 ? "column" : "columns"}`;
  const message = changed > 0
    ? (clipboard.cut ? `Moved ${label}` : `Pasted ${label}${placement === "after-selection" ? " right" : ""}`)
    : "Column paste did not change the document";
  applyMutationPayload(payload, message);
  document.body.dataset.structurePaste = `column:${clipboard.cut ? "cut" : "copy"}:${placement}:${sources.join(",")}->${targetCol}`;
  if (clipboard.cut) clearStructureClipboard();
}

function hideSelectedColumn() {
  if (!documentId || columns.length <= 1) {
    status.textContent = columns.length <= 1 ? "Cannot hide the last visible column" : "Open a CSV first";
    return;
  }
  if (!ensureEditable("Hide Column")) return;
  const sourceCol = sourceColFor(selected.col);
  hiddenSourceColumns.add(sourceCol);
  applyHiddenColumns();
  saveHiddenColumnsState();
  const nextCol = Math.max(0, Math.min(columns.length - 1, selected.col));
  selected = { row: selected.row, col: nextCol };
  selectionAnchor = { ...selected };
  closeColumnMenu();
  draw();
  updateInspector();
  status.textContent = `Column hidden: ${columnNameForSource(sourceCol)}`;
}

function showAllHiddenColumns() {
  if (!documentId) return;
  const count = hiddenSourceColumns.size;
  if (count === 0) {
    status.textContent = "No hidden columns";
    return;
  }
  const sourceCol = sourceColFor(selected.col);
  hiddenSourceColumns = new Set();
  applyHiddenColumns();
  saveHiddenColumnsState();
  const visibleCol = visibleColForSource(sourceCol);
  selected = {
    row: Math.max(0, Math.min(totalRows - 1, selected.row)),
    col: visibleCol >= 0 ? visibleCol : Math.max(0, Math.min(columns.length - 1, selected.col)),
  };
  selectionAnchor = { ...selected };
  closeColumnMenu();
  draw();
  updateInspector();
  status.textContent = `Restored ${formatCount(count)} hidden ${count === 1 ? "column" : "columns"}`;
}

async function moveSelectedColumn(delta) {
  if (!documentId) return;
  const previousCol = selected.col;
  const label = delta < 0 ? "Moved column left" : "Moved column right";
  await moveColumnToVisible(previousCol, Math.max(0, Math.min(columns.length - 1, previousCol + delta)), "command", label);
}

async function moveColumnToVisible(sourceVisibleCol, targetVisibleCol, source = "command", labelOverride = "") {
  if (!documentId) return;
  if (!ensureEditable("Move Column")) return;
  if (!columns[sourceVisibleCol] || !columns[targetVisibleCol]) {
    status.textContent = "Select a column first.";
    return;
  }
  const sourceCol = sourceColFor(sourceVisibleCol);
  const targetCol = sourceColFor(targetVisibleCol);
  const delta = targetCol - sourceCol;
  const label = labelOverride || `Moved column ${sourceVisibleCol + 1} to ${targetVisibleCol + 1}`;
  if (delta === 0) {
    document.body.dataset.columnDragReorder = `${sourceVisibleCol}->${targetVisibleCol}:cancel`;
    status.textContent = "Column is already at the target.";
    draw();
    return;
  }
  status.textContent = `${label}...`;
  const response = await fetch("/api/column/move", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, col: sourceCol, delta }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const changed = payload.changed ?? 0;
  applyMutationPayload(payload, changed > 0 ? label : "Column is already at the edge");
  if (changed > 0) {
    const movedCol = Math.max(0, Math.min(columns.length - 1, targetVisibleCol));
    selected = { row: 0, col: movedCol };
    selectionAnchor = { row: Math.max(0, sourceRowCount - 1), col: movedCol };
    document.body.dataset.columnDragReorder = `${sourceVisibleCol}->${movedCol}:moved`;
    document.body.dataset.columnDragSource = source;
    document.body.dataset.columnDragPhysical = `${sourceCol}->${targetCol}`;
    draw();
    status.textContent = label;
  }
}

async function transposeCurrentDocument() {
  const ranges = activeSelectionRanges();
  if (ranges.length > 1) {
    if (transposeSelectedRanges(ranges)) return;
    return;
  }
  if (ranges.length === 1 && rangeCellCount(ranges[0]) > 1) {
    if (transposeSelectedRange(ranges[0])) return;
    return;
  }
  const plan = wholeTableTransposePlan();
  if (!plan) return;
  if (plan.cells > transposeMaxCells || plan.cells >= transposeConfirmThresholdCells) {
    openTransposeConfirmPanel(plan);
    return;
  }
  await runWholeTableTranspose(plan);
}

function wholeTableTransposePlan() {
  if (!documentId) return null;
  if (!ensureEditable("Transpose Table")) return null;
  const rows = Math.max(0, sourceRowCount + dataStartRow);
  const columnsCount = Math.max(
    columns.length,
    allColumns.length,
    metadata?.dialect?.estimatedColumns || 0,
  );
  const cells = rows * columnsCount;
  return {
    rows,
    columns: columnsCount,
    cells,
    blocked: cells > transposeMaxCells,
  };
}

function openTransposeConfirmPanel(plan) {
  if (!transposeConfirmPanel || !plan) return;
  pendingWholeTranspose = plan;
  const cells = formatCount(plan.cells);
  const rows = formatCount(plan.rows);
  const columnsCount = formatCount(plan.columns);
  transposeConfirmSummary.textContent = `Transpose ${rows} rows x ${columnsCount} columns`;
  transposeConfirmDetail.textContent = "This will materialize the current table, swap rows and columns, and create one undo step before Save.";
  transposeConfirmWarning.textContent = plan.blocked
    ? `This transpose has ${cells} cells, above the ${formatCount(transposeMaxCells)} cell safety limit.`
    : `This transpose has ${cells} cells and may take noticeable time and memory.`;
  confirmTransposeButton.disabled = plan.blocked;
  transposeConfirmPanel.hidden = false;
  document.body.dataset.transposeConfirm = plan.blocked ? `blocked:${plan.cells}` : `open:${plan.cells}`;
  status.textContent = plan.blocked ? "Transpose table is too large." : "Confirm table transpose.";
  (plan.blocked ? cancelTransposeConfirmButton : confirmTransposeButton).focus();
}

function closeTransposeConfirmPanel() {
  if (!transposeConfirmPanel || transposeConfirmPanel.hidden) return;
  transposeConfirmPanel.hidden = true;
  pendingWholeTranspose = null;
  document.body.dataset.transposeConfirm = "closed";
  canvas.focus();
}

async function confirmWholeTableTranspose() {
  const plan = pendingWholeTranspose;
  closeTransposeConfirmPanel();
  if (!plan || plan.blocked) return;
  await runWholeTableTranspose(plan);
}

async function runWholeTableTranspose(plan = wholeTableTransposePlan()) {
  if (!plan) return;
  document.body.dataset.transposeConfirm = `run:${plan.cells}`;
  await mutateStructure("/api/transpose", { maxCells: transposeMaxCells }, "Transposed table");
}

function transposeSelectedRange(range, maxCells = 100_000) {
  if (!documentId || !range) return false;
  if (!ensureEditable("Transpose Selection")) return true;
  const plan = buildSelectionTransposePlan(range, maxCells);
  if (!plan) return true;
  if (!plan.edits.length) {
    status.textContent = "Transpose Selection: no cells changed.";
    document.body.dataset.transposeSelection = `${plan.rowCount}x${plan.colCount}->${plan.colCount}x${plan.rowCount}:0`;
    return true;
  }
  applySelectionTransposePlans(
    [plan],
    `Transposed ${plan.rowCount}x${plan.colCount} selection`,
  );
  return true;
}

function transposeSelectedRanges(ranges, maxCells = 100_000) {
  if (!documentId || !Array.isArray(ranges) || ranges.length === 0) return false;
  if (!ensureEditable("Transpose Selection")) return true;
  const cellCount = ranges.reduce((count, range) => count + rangeCellCount(range), 0);
  if (cellCount > maxCells) {
    status.textContent = `Selection has ${formatCount(cellCount)} cells; transpose limit is ${formatCount(maxCells)}.`;
    return true;
  }
  const plans = [];
  for (const range of ranges) {
    if (rangeCellCount(range) <= 1) {
      status.textContent = "Each transpose range must contain at least two cells.";
      return true;
    }
    const plan = buildSelectionTransposePlan(range, maxCells);
    if (!plan) return true;
    plans.push(plan);
  }
  for (let i = 0; i < plans.length; i++) {
    for (let j = i + 1; j < plans.length; j++) {
      if (rangesOverlap(plans[i].unionRange, plans[j].unionRange)) {
        status.textContent = "Multi-range transpose targets must not overlap each other.";
        document.body.dataset.transposeSelection = `multi:${plans.length}:overlap`;
        return true;
      }
    }
  }
  const edits = plans.flatMap((plan) => plan.edits);
  if (!edits.length) {
    status.textContent = "Transpose Selection: no cells changed.";
    document.body.dataset.transposeSelection = `multi:${plans.length}:0`;
    return true;
  }
  applySelectionTransposePlans(
    plans,
    `Transposed ${plans.length} ranges (${formatCount(edits.length)} cells changed)`,
  );
  return true;
}

function buildSelectionTransposePlan(range, maxCells = 100_000) {
  if (range.endRow >= sourceRowCount) {
    status.textContent = "Aggregate rows are view-only.";
    return null;
  }
  const rowCount = range.endRow - range.startRow + 1;
  const colCount = range.endCol - range.startCol + 1;
  const cellCount = rowCount * colCount;
  if (cellCount > maxCells) {
    status.textContent = `Selection has ${formatCount(cellCount)} cells; transpose limit is ${formatCount(maxCells)}.`;
    return null;
  }
  const targetEndRow = range.startRow + colCount - 1;
  const targetEndCol = range.startCol + rowCount - 1;
  if (targetEndRow >= sourceRowCount && viewMapped) {
    status.textContent = "Selection transpose target rows must fit while a filtered or sorted view is active.";
    return null;
  }
  const transposeSourceCols = new Map();
  const maxExistingSourceCol = Math.max(-1, ...allColumns.map((column) => column.sourceCol));
  for (let col = range.startCol; col <= Math.max(range.endCol, targetEndCol); col++) {
    const sourceCol = col < columns.length
      ? sourceColFor(col)
      : maxExistingSourceCol + (col - columns.length) + 1;
    transposeSourceCols.set(col, sourceCol);
  }
  const sourceMatrix = [];
  for (let row = range.startRow; row <= range.endRow; row++) {
    if (sourceRowForDisplayRow(row, "Transpose Selection") == null) return null;
    const values = [];
    for (let col = range.startCol; col <= range.endCol; col++) {
      values.push(valueAt(row, col));
    }
    sourceMatrix.push(values);
  }
  const edits = [];
  const unionEndRow = Math.max(range.endRow, targetEndRow);
  const unionEndCol = Math.max(range.endCol, targetEndCol);
  for (let row = range.startRow; row <= unionEndRow; row++) {
    const sourceRow = sourceRowForDisplayRow(row, "Transpose Selection");
    if (sourceRow == null) return null;
    for (let col = range.startCol; col <= unionEndCol; col++) {
      const inSource = row <= range.endRow && col <= range.endCol;
      const inTarget = row <= targetEndRow && col <= targetEndCol;
      if (!inSource && !inTarget) continue;
      const nextValue = inTarget
        ? sourceMatrix[col - range.startCol]?.[row - range.startRow] ?? ""
        : "";
      const sourceCol = transposeSourceCols.get(col) ?? sourceColFor(col);
      if (valueAtSourceColumn(row, sourceCol) === nextValue) continue;
      edits.push({ row: sourceRow, col: sourceCol, value: nextValue });
    }
  }
  return {
    range,
    targetRange: { startRow: range.startRow, endRow: targetEndRow, startCol: range.startCol, endCol: targetEndCol },
    unionRange: { startRow: range.startRow, endRow: unionEndRow, startCol: range.startCol, endCol: unionEndCol },
    rowCount,
    colCount,
    edits,
  };
}

function applySelectionTransposePlans(plans, message) {
  const edits = plans.flatMap((plan) => plan.edits);
  const [firstPlan] = plans;
  if (!firstPlan || !edits.length) return;
  ensureSourceRowsForEdits(edits);
  ensureSourceColumnsForEdits(edits);
  selected = { row: firstPlan.targetRange.startRow, col: firstPlan.targetRange.startCol };
  selectionAnchor = { row: firstPlan.targetRange.endRow, col: firstPlan.targetRange.endCol };
  setSmartSelectionRanges(plans.slice(1).map((plan) => plan.targetRange));
  autoFitColumnsForEdits(edits);
  syncCellValueAfterEdit();
  const summaries = plans.map((plan) => `${plan.rowCount}x${plan.colCount}->${plan.colCount}x${plan.rowCount}`);
  document.body.dataset.transposeSelection = plans.length === 1
    ? `${summaries[0]}:${edits.length}`
    : `multi:${plans.length}:${summaries.join(";")}:${edits.length}`;
  submitCellEdits(
    edits,
    Math.min(...plans.map((plan) => plan.range.startRow)),
    message,
  );
}

function sourceRowForDisplayRow(row, action) {
  const sourceRow = rowSource.get(row);
  if (viewMapped && sourceRow == null) {
    status.textContent = `Load the selected filtered rows before using ${action}.`;
    return null;
  }
  return sourceRow ?? row + dataStartRow;
}

async function toggleHeaderRow() {
  if (!documentId || !metadata) return;
  if (!ensureEditable("Toggle Header Row")) return;
  const hasHeader = !Boolean(metadata.dialect?.hasHeader);
  const label = hasHeader ? "Header row enabled" : "Header row disabled";
  status.textContent = `${label}...`;
  const response = await fetch("/api/header-mode", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, hasHeader }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Toggle Header Row failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  applyMutationPayload(payload, label);
  document.body.dataset.headerMode = hasHeader ? "header" : "no-header";
}

async function dedupeRows() {
  if (!documentId) return;
  if (!ensureEditable("Remove Duplicate Rows")) return;
  status.textContent = "Removing duplicate rows...";
  const response = await fetch("/api/dedupe-rows", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Remove Duplicate Rows failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const removed = payload.changed ?? 0;
  applyMutationPayload(
    payload,
    removed > 0
      ? `Removed ${formatCount(removed)} duplicate rows`
      : "No duplicate rows found",
  );
}

async function dedupeColumns() {
  if (!documentId) return;
  if (!ensureEditable("Remove Duplicate Columns")) return;
  status.textContent = "Removing duplicate columns...";
  const response = await fetch("/api/dedupe-columns", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Remove Duplicate Columns failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const removed = payload.changed ?? 0;
  applyMutationPayload(
    payload,
    removed > 0
      ? `Removed ${formatCount(removed)} duplicate columns`
      : "No duplicate columns found",
  );
}

async function transformSelectedText(mode, label) {
  if (!documentId) return;
  if (!ensureEditable(label)) return;
  const cells = selectedCellsForTransform();
  if (!cells) return;
  status.textContent = localizedTransformStatus("running", label);
  const response = await fetch("/api/text-transform", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, mode, cells }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = localizedTransformStatus("failed", label, { message: payload.error ?? "unknown error" });
    return;
  }
  const changed = payload.changed ?? 0;
  applyMutationPayload(
    payload,
    changed > 0
      ? localizedTransformStatus("changed", label, { count: formatCount(changed) })
      : localizedTransformStatus("noChange", label),
  );
}

async function transformColumnTextForCommand(commandID) {
  const spec = columnTextTransformCommands.get(commandID);
  if (!spec) return;
  await transformColumnText(spec.mode, localizedCommandLabel(commandID, spec.label));
}

async function transformColumnText(mode, label, column = activeTransformColumn()) {
  closeColumnMenu();
  if (!documentId || !columns[column]) return;
  if (!ensureEditable(label)) return;
  const sourceCol = sourceColFor(column);
  status.textContent = localizedTransformStatus("running", label);
  const response = await fetch("/api/text-transform", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, mode, column: sourceCol }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = localizedTransformStatus("failed", label, { message: payload.error ?? "unknown error" });
    return;
  }
  const changed = payload.changed ?? 0;
  document.body.dataset.columnTransform = `${sourceCol}:${mode}:${changed}`;
  document.body.dataset.columnTransformStreamed = String(Boolean(payload.streamed));
  applyMutationPayload(
    payload,
    changed > 0
      ? localizedTransformStatus("changed", label, { count: formatCount(changed) })
      : localizedTransformStatus("noChange", label),
  );
}

function activeTransformColumn() {
  return activeCommandColumn();
}

function activeCommandColumn() {
  return columnMenuColumn >= 0 ? columnMenuColumn : selected.col;
}

function syncColumnTransformCommand() {
  const mode = columnTransformMode?.value || "trim";
  const commandID = columnTextTransformCommandByMode.get(mode) || "text.column.trim";
  if (columnTransformApply) columnTransformApply.dataset.command = commandID;
  return commandID;
}

async function swapSelectedCells() {
  if (!documentId) return;
  if (!ensureEditable("Swap Cells")) return;
  if (selectionCellCount() !== 2) {
    status.textContent = "Select exactly two cells to swap.";
    return;
  }
  const cells = selectedCellsForTransform(2);
  if (!cells) return;
  const range = normalizedSelection();
  await runCellTool("/api/cells/swap", {
    cells,
    activeRow: range.startRow,
    activeCol: range.startCol,
    label: "Swapped cells",
  });
}

async function swapCellWithNeighbor(direction) {
  if (!documentId) return;
  const label = swapNeighborLabel(direction);
  if (!ensureEditable(label)) return;
  const origin = { row: selected.row, col: selected.col };
  if (selected.row >= sourceRowCount) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  const target = swapNeighborTarget(direction);
  if (!target) {
    status.textContent = `${label}: no adjacent cell.`;
    return;
  }
  const activeSourceRow = sourceRowForDisplayRow(selected.row, label);
  if (activeSourceRow == null) return;
  const targetSourceRow = sourceRowForDisplayRow(target.row, label);
  if (targetSourceRow == null) return;
  const activeSourceCol = sourceColFor(selected.col);
  const targetSourceCol = sourceColFor(target.col);
  await runCellTool("/api/cells/swap", {
    cells: [
      { row: activeSourceRow, col: activeSourceCol },
      { row: targetSourceRow, col: targetSourceCol },
    ],
    activeRow: target.row,
    activeCol: target.col,
    label,
  });
  document.body.dataset.cellSwapNeighbor = `${direction}:${origin.row}:${origin.col}->${target.row}:${target.col}`;
}

function swapNeighborTarget(direction) {
  switch (direction) {
    case "up":
      return selected.row > 0 ? { row: selected.row - 1, col: selected.col } : null;
    case "down":
      return selected.row + 1 < sourceRowCount ? { row: selected.row + 1, col: selected.col } : null;
    case "left":
      return selected.col > 0 ? { row: selected.row, col: selected.col - 1 } : null;
    case "right":
      return selected.col + 1 < columns.length ? { row: selected.row, col: selected.col + 1 } : null;
    default:
      return null;
  }
}

function swapNeighborLabel(direction) {
  return {
    up: "Swapped cell up",
    down: "Swapped cell down",
    left: "Swapped cell left",
    right: "Swapped cell right",
  }[direction] ?? "Swapped cell";
}

async function joinSelectedCells(direction = "horizontal") {
  if (!documentId) return;
  const vertical = direction === "vertical";
  const title = vertical ? "Join Cells Vertically" : "Join Cells";
  if (!ensureEditable(title)) return;
  if (selectionCellCount() < 2) {
    status.textContent = `Select at least two cells to ${vertical ? "join vertically" : "join"}.`;
    return;
  }
  const cells = selectedCellsForTransform();
  if (!cells) return;
  const range = normalizedSelection();
  const separator = selectedCellSeparator();
  await runCellTool("/api/cells/join", {
    cells,
    separator,
    direction,
    activeRow: vertical ? range.startRow : selected.row,
    activeCol: vertical ? selected.col : range.startCol,
    label: vertical ? "Joined cells vertically" : "Joined cells",
  });
}

async function splitActiveCell(direction = "horizontal") {
  if (!documentId) return;
  const vertical = direction === "vertical";
  const title = vertical ? "Split Cell Vertically" : "Split Cell";
  if (!ensureEditable(title)) return;
  const separator = selectedCellSeparator();
  const request = {
    cells: [{ row: selectedSourceRow(), col: sourceColFor(selected.col) }],
    separator,
    direction,
    activeRow: selected.row,
    activeCol: selected.col,
    label: vertical ? "Split cell vertically" : "Split cell",
  };
  const result = await runCellTool("/api/cells/split", request);
  if (result?.ok || !result?.error) return;
  const conflict = parseSplitOverwriteConflict(result.error);
  if (conflict) {
    openSplitOverwriteConfirmPanel(
      request,
      {
        sourceRow: selected.row + 1,
        sourceCol: selected.col + 1,
        separator,
        sourceValue: valueAt(selected.row, selected.col),
        direction,
        conflict,
      },
    );
    return;
  }
}

function parseSplitOverwriteConflict(error) {
  const message = String(error ?? "").trim();
  const prefix = "split would overwrite row ";
  if (!message.startsWith(prefix)) return null;
  const body = message.slice(prefix.length);
  const location = body.match(/^(\d+), column (\d+)/);
  if (!location) return null;
  const row = Number.parseInt(location[1], 10);
  const col = Number.parseInt(location[2], 10);
  if (!Number.isFinite(row) || !Number.isFinite(col)) return null;
  const additionalMatch = message.match(/ and (\d+) more;/);
  const additional = additionalMatch ? Number.parseInt(additionalMatch[1], 10) : 0;
  return {
    row,
    col,
    additional: Number.isFinite(additional) ? additional : 0,
  };
}

function splitCellPreviewParts(value, separator) {
  return String(value ?? "").split(separator);
}

function splitWarningMessage(conflictCount, firstRow, firstCol, direction) {
  const directionLabel = direction === "vertical" ? "vertical" : "horizontal";
  const plural = conflictCount === 1 ? "cell" : "cells";
  const directionLabelText = directionLabel === "vertical" ? "vertically" : "horizontally";
  return `This ${directionLabelText} split will overwrite ${formatCount(conflictCount)} existing ${plural}${firstRow && firstCol ? `, starting at row ${formatCount(firstRow)}, column ${formatCount(firstCol)}.` : "."}`;
}

function openSplitOverwriteConfirmPanel(request, preview) {
  if (!splitConfirmPanel || !request) return;
  pendingSplitTool = { ...request };
  const direction = request.direction === "vertical" ? "Vertical" : "Horizontal";
  const sourceRow = preview?.sourceRow || selected.row + 1;
  const sourceCol = preview?.sourceCol || selected.col + 1;
  const separator = printableCellSeparator(request.separator);
  const sourceParts = splitCellPreviewParts(preview?.sourceValue ?? valueAt(selected.row, selected.col), request.separator);
  const conflictInfo = preview?.conflict;
  const conflictCount = conflictInfo ? 1 + (conflictInfo.additional || 0) : 1;
  const firstConflict = conflictInfo ? `row ${formatCount(conflictInfo.row)}, column ${formatCount(conflictInfo.col)}` : null;
  splitConfirmSummary.textContent = `${direction} Split Overwrite`;
  splitConfirmDetail.textContent = `Source cell: row ${formatCount(sourceRow)}, column ${formatCount(sourceCol)} · Separator: ${separator} · ${formatCount(sourceParts.length)} values.`;
  splitConfirmWarning.textContent = conflictInfo
    ? `${splitWarningMessage(conflictCount, conflictInfo.row, conflictInfo.col, request.direction)}${conflictInfo.additional > 0 ? ` ${formatCount(conflictInfo.additional)} more conflict${conflictInfo.additional === 1 ? "" : "s"} beyond ${firstConflict}.` : ""}`
    : `This split would overwrite existing data. Review the first target and retry with overwrite enabled.`;
  splitConfirmPanel.hidden = false;
  status.textContent = conflictInfo
    ? `Split blocked by ${formatCount(conflictCount)} existing cells.`
    : "Split blocked by existing cells. Confirm to overwrite.";
  confirmSplitConfirmButton.focus();
}

function closeSplitOverwriteConfirmPanel() {
  if (!splitConfirmPanel || splitConfirmPanel.hidden) return;
  splitConfirmPanel.hidden = true;
  pendingSplitTool = null;
  canvas.focus();
}

async function confirmSplitOverwrite() {
  const pending = pendingSplitTool;
  closeSplitOverwriteConfirmPanel();
  if (!pending) return;
  await runCellTool("/api/cells/split", {
    ...pending,
    allowOverwrite: true,
  });
}

function autoFillSelection(direction) {
  if (!documentId) return;
  const label = autoFillLabel(direction);
  if (!ensureEditable(label)) return;
  const range = normalizedSelection();
  if (range.endRow >= sourceRowCount) {
    status.textContent = "Aggregate rows are view-only.";
    return;
  }
  const axisLength = direction === "left" || direction === "right"
    ? range.endCol - range.startCol + 1
    : range.endRow - range.startRow + 1;
  if (axisLength < 2 || selectionCellCount() < 2) {
    status.textContent = `${label}: select at least two cells.`;
    return;
  }
  if (selectionCellCount() > 50_000) {
    status.textContent = `${label}: selection limit is ${formatCount(50_000)} cells.`;
    return;
  }
  const edits = autoFillEdits(range, direction);
  if (!edits) return;
  if (!edits.length) {
    document.body.dataset.autoFillPreview = `${direction}:0`;
    status.textContent = `${label}: no cells changed`;
    return;
  }
  document.body.dataset.autoFillPreview = `${direction}:${edits.length}`;
  windowGeneration++;
  applyEditsToLoadedRows(edits);
  autoFitColumnsForEdits(edits);
  syncCellValueAfterEdit();
  submitCellEdits(edits, range.startRow, `${label}: ${formatCount(edits.length)} cells changed`);
}

function autoFillLabel(direction) {
  return {
    down: "Fill Down",
    right: "Fill Right",
    up: "Fill Up",
    left: "Fill Left",
  }[direction] ?? "Auto-fill";
}

function autoFillEdits(range, direction) {
  const edits = [];
  if (direction === "down" || direction === "up") {
    const rowsInOrder = direction === "up"
      ? descendingIndexes(range.startRow, range.endRow)
      : ascendingIndexes(range.startRow, range.endRow);
    if (!ensureRowsLoaded(rowsInOrder, "Auto-fill")) return null;
    for (let col = range.startCol; col <= range.endCol; col++) {
      const sourceCol = sourceColFor(col);
      const values = rowsInOrder.map((row) => valueAt(row, col));
      const filled = autoFillSeries(values);
      rowsInOrder.forEach((row, index) => {
        queueAutoFillEdit(edits, row, sourceCol, filled[index]);
      });
    }
    return edits;
  }

  const colsInOrder = direction === "left"
    ? descendingIndexes(range.startCol, range.endCol)
    : ascendingIndexes(range.startCol, range.endCol);
  const rowsInOrder = ascendingIndexes(range.startRow, range.endRow);
  if (!ensureRowsLoaded(rowsInOrder, "Auto-fill")) return null;
  for (const row of rowsInOrder) {
    const values = colsInOrder.map((col) => valueAt(row, col));
    const filled = autoFillSeries(values);
    colsInOrder.forEach((col, index) => {
      queueAutoFillEdit(edits, row, sourceColFor(col), filled[index]);
    });
  }
  return edits;
}

function ascendingIndexes(start, end) {
  return Array.from({ length: Math.max(0, end - start + 1) }, (_, index) => start + index);
}

function descendingIndexes(start, end) {
  return ascendingIndexes(start, end).reverse();
}

function ensureRowsLoaded(displayRows, action) {
  for (const row of displayRows) {
    if (!rows.has(row)) {
      maybeLoadWindow(row);
      status.textContent = `${action}: load the selected rows before filling.`;
      return false;
    }
    if (viewMapped && rowSource.get(row) == null) {
      status.textContent = `${action}: load the selected filtered rows before filling.`;
      return false;
    }
  }
  return true;
}

function queueAutoFillEdit(edits, displayRow, sourceCol, value) {
  const record = rows.get(displayRow) ?? [];
  const nextValue = String(value ?? "");
  if ((record[sourceCol] ?? "") === nextValue) return;
  while (record.length <= sourceCol) record.push("");
  record[sourceCol] = nextValue;
  rows.set(displayRow, record);
  edits.push({
    row: rowSource.get(displayRow) ?? displayRow + dataStartRow,
    col: sourceCol,
    value: nextValue,
  });
}

function autoFillSeries(values) {
  if (values.length <= 1) return values.map((value) => String(value ?? ""));
  const seed = leadingAutoFillSeed(values);
  if (!seed.length) return values.map((value) => String(value ?? ""));
  const numeric = numericAutoFill(seed, values.length);
  if (numeric) return numeric;
  const dates = dateAutoFill(seed, values.length);
  if (dates) return dates;
  const named = namedSequenceAutoFill(seed, values.length);
  if (named) return named;
  const numberedText = numberedTextAutoFill(seed, values.length);
  if (numberedText) return numberedText;
  return Array.from({ length: values.length }, (_, index) => seed[index % seed.length]);
}

function leadingAutoFillSeed(values) {
  const seed = [];
  for (const value of values) {
    const text = String(value ?? "");
    if (text === "") break;
    seed.push(text);
  }
  return seed.length ? seed : [];
}

function numericAutoFill(seed, length) {
  const numbers = seed.map(parseAutoFillNumber);
  if (numbers.some((value) => value == null)) return null;
  const step = numbers.length >= 2 ? numbers[1] - numbers[0] : 0;
  const decimals = Math.max(...seed.map(decimalPlaces));
  return Array.from({ length }, (_, index) => formatAutoFillNumber(numbers[0] + step * index, decimals));
}

function parseAutoFillNumber(value) {
  const text = String(value ?? "").trim();
  if (!/^[+-]?(?:\d+(?:\.\d*)?|\.\d+)$/.test(text)) return null;
  const number = Number(text);
  return Number.isFinite(number) ? number : null;
}

function decimalPlaces(value) {
  const match = String(value ?? "").trim().match(/\.(\d+)$/);
  return match ? match[1].length : 0;
}

function formatAutoFillNumber(value, decimals) {
  if (decimals > 0) return value.toFixed(decimals);
  return String(Number.isInteger(value) ? value : Number(value.toFixed(12)));
}

function dateAutoFill(seed, length) {
  const dates = seed.map(parseISODate);
  if (dates.some((value) => value == null)) return null;
  const stepDays = dates.length >= 2
    ? Math.round((dates[1] - dates[0]) / 86_400_000)
    : 1;
  return Array.from({ length }, (_, index) => {
    const date = new Date(dates[0] + stepDays * index * 86_400_000);
    return date.toISOString().slice(0, 10);
  });
}

function parseISODate(value) {
  const text = String(value ?? "").trim();
  if (!/^\d{4}-\d{2}-\d{2}$/.test(text)) return null;
  const [year, month, day] = text.split("-").map(Number);
  const time = Date.parse(`${text}T00:00:00Z`);
  if (!Number.isFinite(time)) return null;
  const date = new Date(time);
  if (
    date.getUTCFullYear() !== year ||
    date.getUTCMonth() + 1 !== month ||
    date.getUTCDate() !== day
  ) {
    return null;
  }
  return time;
}

const autoFillNamedSequences = [
  ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"],
  ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
  ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"],
  ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
];

function namedSequenceAutoFill(seed, length) {
  for (const sequence of autoFillNamedSequences) {
    const matches = seed.map((value) => namedSequenceMatch(value, sequence));
    if (matches.some((match) => !match)) continue;
    const step = matches.length >= 2
      ? positiveModulo(matches[1].index - matches[0].index, sequence.length)
      : 1;
    const consistent = matches.every((match, index) =>
      match.index === positiveModulo(matches[0].index + step * index, sequence.length),
    );
    if (!consistent) continue;
    const sample = matches[0].text;
    return Array.from({ length }, (_, index) =>
      applyNamedSequenceCase(sequence[positiveModulo(matches[0].index + step * index, sequence.length)], sample),
    );
  }
  return null;
}

function namedSequenceMatch(value, sequence) {
  const text = String(value ?? "").trim();
  if (!text) return null;
  const index = sequence.findIndex((candidate) => candidate.toLowerCase() === text.toLowerCase());
  return index >= 0 ? { index, text } : null;
}

function positiveModulo(value, size) {
  return ((value % size) + size) % size;
}

function applyNamedSequenceCase(value, sample) {
  if (sample === sample.toUpperCase()) return value.toUpperCase();
  if (sample === sample.toLowerCase()) return value.toLowerCase();
  return value;
}

function numberedTextAutoFill(seed, length) {
  const parsed = seed.map((value) => String(value).match(/^(.*?)(-?\d+)(\D*)$/));
  if (parsed.some((match) => !match)) return null;
  const [firstPrefix, firstSuffix] = [parsed[0][1], parsed[0][3]];
  if (parsed.some((match) => match[1] !== firstPrefix || match[3] !== firstSuffix)) return null;
  const numbers = parsed.map((match) => Number(match[2]));
  if (numbers.some((value) => !Number.isFinite(value))) return null;
  const step = numbers.length >= 2 ? numbers[1] - numbers[0] : 1;
  const width = Math.max(...parsed.map((match) => match[2].replace("-", "").length));
  return Array.from({ length }, (_, index) => {
    const value = numbers[0] + step * index;
    const sign = value < 0 ? "-" : "";
    const digits = String(Math.abs(value)).padStart(width, "0");
    return `${firstPrefix}${sign}${digits}${firstSuffix}`;
  });
}

window.__opencsvAutoFillSeries = autoFillSeries;

async function runCellTool(endpoint, options) {
  status.textContent = `${options.label}...`;
  delete document.body.dataset.cellToolFailed;
  if (options.separator !== undefined) {
    document.body.dataset.cellToolEndpoint = endpoint;
    document.body.dataset.cellToolDirection = options.direction ?? "";
    document.body.dataset.cellToolSeparator = printableCellSeparator(options.separator);
  }
  const response = await fetch(endpoint, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      cells: options.cells,
      separator: options.separator,
      direction: options.direction,
      allowOverwrite: options.allowOverwrite === true,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    const error = payload.error ?? "unknown error";
    document.body.dataset.cellToolFailed = `${endpoint}:${error}`;
    status.textContent = `${options.label} failed: ${error}`;
    return { ok: false, error, status: response.status, endpoint, payload };
  }
  const changed = payload.changed ?? 0;
  applyMutationPayload(payload, changed > 0 ? `${options.label}: ${formatCount(changed)} cells changed` : `${options.label}: no cells changed`);
  selected = {
    row: Math.max(0, Math.min(totalRows - 1, options.activeRow)),
    col: Math.max(0, Math.min(columns.length - 1, options.activeCol)),
  };
  selectionAnchor = { ...selected };
  draw();
  status.textContent = changed > 0 ? `${options.label}: ${formatCount(changed)} cells changed` : `${options.label}: no cells changed`;
  return { ok: true, payload, endpoint };
}

async function mutateStructure(endpoint, body, label) {
  if (!documentId) return;
  if (!ensureEditable(label)) return;
  status.textContent = `${label}...`;
  const response = await fetch(endpoint, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, ...body }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  applyMutationPayload(payload, `${label}: ${formatCount(payload.totalRows ?? 0)} rows`);
}

async function saveCurrentDocument() {
  await saveDocumentToPath("", { saveAs: false });
}

function openSaveAsDialog() {
  if (!documentId) {
    status.textContent = "Open a CSV before using Save As.";
    return;
  }
  if (!ensureEditable("Save As")) return;
  saveAsPath.value = suggestedSaveAsPath();
  saveAsPanel.hidden = false;
  saveAsPath.focus();
  saveAsPath.select();
}

function closeSaveAsDialog() {
  saveAsPanel.hidden = true;
  canvas.focus();
}

async function openRestoreBackupPanel() {
  if (!documentId) {
    status.textContent = "Open a CSV before restoring from backup.";
    return;
  }
  if (!ensureEditable("Restore From Backup")) return;
  restoreBackupPanel.hidden = false;
  renderRestoreBackupMessage("Loading backups...");
  await loadRestoreBackups();
}

function closeRestoreBackupPanel() {
  restoreBackupPanel.hidden = true;
  canvas.focus();
}

function renderRestoreBackupMessage(message) {
  restoreBackupList.replaceChildren(emptyPaletteNode(message));
}

async function loadRestoreBackups() {
  if (!documentId) return;
  const response = await fetch(`/api/backups?id=${encodeURIComponent(documentId)}`);
  const payload = await response.json();
  if (!response.ok) {
    document.body.dataset.backupCount = "0";
    renderRestoreBackupMessage(`Backups failed: ${payload.error ?? "unknown error"}`);
    return;
  }
  const backups = Array.isArray(payload.backups) ? payload.backups : [];
  document.body.dataset.backupCount = String(backups.length);
  renderRestoreBackupList(backups);
}

function renderRestoreBackupList(backups) {
  restoreBackupList.replaceChildren();
  if (!backups.length) {
    renderRestoreBackupMessage("No backups found for this file.");
    return;
  }
  for (const backup of backups) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "restore-backup-item";
    button.dataset.backupPath = backup.path ?? "";
    const detail = document.createElement("div");
    const name = document.createElement("strong");
    name.textContent = backup.name || fileNameFromPath(backup.path || "backup");
    const path = document.createElement("span");
    path.textContent = backup.path ?? "";
    detail.append(name, path);
    const meta = document.createElement("span");
    meta.className = "restore-backup-meta";
    meta.textContent = `${formatBackupSize(backup.sizeBytes)} · ${formatBackupTime(backup.modTimeUnixNano)}`;
    button.append(detail, meta);
    button.addEventListener("click", () => restoreBackup(backup.path));
    restoreBackupList.append(button);
  }
  restoreBackupList.querySelector(".restore-backup-item")?.focus();
}

async function restoreBackup(path) {
  if (!documentId || !path) return;
  if (!ensureEditable("Restore From Backup")) return;
  status.textContent = `Restoring ${fileNameFromPath(path)}...`;
  const response = await fetch("/api/restore-backup", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, path }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Restore From Backup failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  document.body.dataset.backupRestored = path;
  applyMutationPayload(payload, `Restored ${fileNameFromPath(path)}: ${formatCount(payload.totalRows ?? 0)} rows`);
  closeRestoreBackupPanel();
}

function updateExternalChangePanel() {
  const conflict = Boolean(dirty);
  externalChangeTitle.textContent = conflict ? "Unsaved Edits Conflict" : "File Changed on Disk";
  externalChangeMessage.textContent = conflict
    ? "The backing file changed on disk while this document has unsaved edits. Save As keeps both versions; Overwrite Disk replaces the changed file."
    : "The backing file changed on disk and can be refreshed now.";
  compareExternalChangeButton.hidden = !metadata?.path;
  saveAsExternalChangeButton.hidden = !conflict;
  overwriteExternalChangeButton.hidden = !conflict;
  reloadExternalChangeButton.textContent = conflict ? "Reload from Disk" : "Refresh from Disk";
  reloadExternalChangeButton.title = conflict ? "Discard unsaved edits and load the disk file" : "Load the disk file";
}

function openExternalChangePanel() {
  if (!documentId || !externalFileChanged) return;
  updateExternalChangePanel();
  externalChangePanel.hidden = false;
  reloadExternalChangeButton.focus();
}

function closeExternalChangePanel() {
  if (!externalChangePanel.hidden) externalChangePanel.hidden = true;
  externalChangePanelDismissed = true;
  canvas.focus();
}

function saveAsExternalChange() {
  externalChangePanel.hidden = true;
  openSaveAsDialog();
}

async function overwriteExternalChangeOnDisk() {
  if (!documentId || !externalFileChanged) return;
  document.body.dataset.externalOverwrite = "confirm";
  await saveDocumentToPath("", { allowExternalOverwrite: true });
  if (!externalFileChanged) {
    externalChangePanel.hidden = true;
    document.body.dataset.externalOverwrite = "saved";
  }
}

async function compareExternalChangeWithDisk() {
  if (!documentId || !metadata?.path || !externalFileChanged) return;
  const sourceColumns = (allColumns.length ? allColumns : columns).map((column) => column.sourceCol);
  status.textContent = "Preparing disk/local compare...";
  const response = await fetch("/api/export-view", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      columns: sourceColumns,
      format: "csv",
      fullDocument: true,
    }),
  });
  const blob = await response.blob();
  const text = await blob.text();
  if (!response.ok) {
    status.textContent = `Compare failed: ${errorFromResponseText(text)}`;
    return;
  }
  leftDiffPath.value = metadata.path;
  rightDiffPath.value = "";
  diffRightText.value = text;
  diffRightText.dataset.source = "local-document";
  setDiffRightTextMode(true);
  externalChangePanel.hidden = true;
  document.body.dataset.externalChangeCompare = dirty ? "dirty" : "clean";
  document.body.dataset.externalChangeCompareRows = response.headers.get("X-Opencsv-Rows") || "";
  openDiffPanel();
  await runDiffCompare();
}

async function reloadExternalChangeFromDisk() {
  if (!documentId || !externalFileChanged) return;
  const discardDirty = Boolean(dirty);
  status.textContent = discardDirty
    ? "Reloading from disk and discarding local edits..."
    : "Refreshing from disk...";
  const refreshed = await refreshCurrentDocument({ discardDirty, manual: true });
  if (refreshed) {
    externalChangePanel.hidden = true;
    externalChangePanelDismissed = false;
  } else {
    updateExternalChangePanel();
  }
}

function suggestedSaveAsPath() {
  if (activeFileEntry?.detail === "temporary upload") {
    return `/tmp/${activeFileEntry?.name || "opencsv-output.csv"}`;
  }
  const current = metadata?.path || activeFileEntry?.detail || pathInput.value.trim() || "/tmp/output.csv";
  const dot = current.lastIndexOf(".");
  const slash = Math.max(current.lastIndexOf("/"), current.lastIndexOf("\\"));
  if (dot > slash) return `${current.slice(0, dot)} copy${current.slice(dot)}`;
  return `${current} copy`;
}

async function confirmSaveAs() {
  const path = saveAsPath.value.trim();
  if (!path) {
    status.textContent = "Enter a Save As path.";
    saveAsPath.focus();
    return;
  }
  await saveDocumentToPath(path, { saveAs: true });
}

async function saveDocumentToPath(path = "", options = {}) {
  if (!documentId) {
    status.textContent = "Open a CSV through opencsv serve before saving.";
    return;
  }
  if (!ensureEditable(options.saveAs ? "Save As" : "Save")) return;
  const draftIDBeforeSave = dirtyDraftID();
  status.textContent = options.saveAs ? `Saving as ${path}...` : "Saving...";
  const response = await fetch("/api/save", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      path,
      backup: true,
      allowExternalOverwrite: Boolean(options.allowExternalOverwrite),
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    const message = payload.error ?? "unknown error";
    if (String(message).includes("file changed on disk")) {
      externalFileChanged = true;
      externalChangePanelDismissed = false;
      document.body.dataset.externalSaveBlocked = "changed";
      updateExternalChangePanel();
      openExternalChangePanel();
      status.textContent = "Save blocked: file changed on disk.";
      return;
    }
    status.textContent = `${options.saveAs ? "Save As" : "Save"} failed: ${message}`;
    return;
  }
  dirty = false;
  clearDirtyDraft(draftIDBeforeSave);
  externalFileChanged = false;
  if (payload.path) {
    metadata = metadata ? { ...metadata, path: payload.path } : metadata;
    pathInput.value = payload.path;
    updateFileList(fileNameFromPath(payload.path), payload.path, totalRows, { remember: true, documentId });
    updateFormatRulePlaceholder(payload.path);
    delete document.body.dataset.newFile;
  } else {
    updateFileList(activeFileEntry?.name || fileNameFromPath(metadata?.path || ""), activeFileEntry?.detail || metadata?.path || "", totalRows, { documentId });
  }
  if (options.saveAs) closeSaveAsDialog();
  status.textContent = payload.backupPath
    ? `Saved ${payload.path}; backup ${payload.backupPath}`
    : `Saved ${payload.path}`;
  await checkFileStatus(false);
}

async function toggleReadOnly() {
  if (!documentId) return;
  const next = !readOnly;
  const response = await fetch("/api/read-only", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, readOnly: next }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Read-only failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  readOnly = Boolean(payload.readOnly);
  updateFileList(activeFileEntry?.name || fileNameFromPath(metadata?.path || ""), activeFileEntry?.detail || metadata?.path || "", totalRows, { documentId });
  status.textContent = readOnly ? "Read-only mode enabled" : "Read-only mode disabled";
  draw();
}

async function checkFileStatus(autoRefresh = true) {
  if (!documentId) return;
  try {
    const response = await fetch(`/api/file-status?id=${encodeURIComponent(documentId)}`);
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "status failed");
    readOnly = Boolean(payload.readOnly);
    const wasChanged = externalFileChanged;
    externalFileChanged = Boolean(payload.changedExternally);
    if (!externalFileChanged) {
      externalChangePanelDismissed = false;
    }
    if (externalFileChanged && autoRefresh && autoRefreshMode === "clean" && !dirty && canAutoRefreshCleanExternalChange()) {
      await refreshCurrentDocument();
      return;
    }
    if (externalFileChanged && autoRefresh && !wasChanged) {
      externalChangePanelDismissed = false;
      if (dirty) {
        status.textContent = "File changed on disk while local edits are unsaved.";
        if (!externalChangePanelDismissed) openExternalChangePanel();
      } else {
        status.textContent = autoRefreshMode === "notify"
          ? "File changed on disk; auto-refresh is set to notify only."
          : "File changed on disk.";
      }
    }
    updateFormatStatus();
    updateCommandState();
  } catch (error) {
    externalStatus.textContent = "Status failed";
  }
}

function canAutoRefreshCleanExternalChange() {
  return !viewMapped || Boolean(activeFilter?.expression) || sortKeys.length > 0;
}

async function refreshCurrentDocument(options = {}) {
  if (!documentId) return;
  const viewSnapshot = options.discardDirty || options.preserveView === false ? null : captureViewSnapshot();
  const response = await fetch("/api/refresh", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, discardDirty: Boolean(options.discardDirty) }),
  });
  const payload = await response.json();
  if (!response.ok) {
    externalFileChanged = true;
    const message = `Refresh failed: ${payload.error ?? "unknown error"}`;
    updateFormatStatus();
    draw();
    status.textContent = message;
    return false;
  }
  const session = payload.session;
  if (!session?.document) return false;
  metadata = session.document.metadata;
  readOnly = Boolean(session.readOnly);
  externalFileChanged = false;
  externalChangePanelDismissed = false;
  loadDocument(session.document);
  dirty = Boolean(session.dirty);
  updateFileList(session.displayName || fileNameFromPath(metadata.path), metadata.path, totalRows, {
    documentId,
    sourceKind: session.sourceKind,
  });
  const restoredView = await restoreViewSnapshot(viewSnapshot);
  updateFormatStatus();
  updateInspector();
  updateCommandState();
  draw();
  const label = options.manual ? "Refreshed external changes" : "Auto-refreshed external changes";
  status.textContent = restoredView.length
    ? `${label}; restored ${restoredView.join(" and ")}`
    : label;
  return true;
}

function captureViewSnapshot() {
  if (!documentId) return null;
  return {
    filter: activeFilter ? { ...activeFilter } : null,
    sortKeys: sortKeys.map((key) => ({ ...key })),
    sortStable,
    selectedSourceCol: columns[selected.col]?.sourceCol ?? selected.col,
  };
}

async function restoreViewSnapshot(snapshot) {
  if (!documentId || !snapshot) {
    document.body.dataset.refreshViewRestored = "";
    return [];
  }
  const restored = [];
  const selectedSourceCol = Number.isInteger(snapshot.selectedSourceCol)
    ? snapshot.selectedSourceCol
    : selected.col;
  if (snapshot.filter?.expression) {
    const preserveCol = visibleColForSource(selectedSourceCol);
    const filterRestored = await applyFilterRequest(
      { documentId, mode: "sql", expression: snapshot.filter.expression },
      "Restored filter",
      preserveCol >= 0 ? preserveCol : selected.col,
    );
    if (filterRestored) {
      setActiveFilter(snapshot.filter);
      restored.push("filter");
    }
  }
  const resolvedSortKeys = (snapshot.sortKeys ?? [])
    .map((key) => ({
      column: Number(key.column),
      descending: Boolean(key.descending),
      numeric: Boolean(key.numeric || key.mode === "number"),
      mode: key.mode || (key.numeric ? "number" : "auto"),
    }))
    .filter((key) => Number.isInteger(key.column) && key.column >= 0 && key.column < allColumns.length);
  if (resolvedSortKeys.length) {
    const visibleCol = visibleColForSource(resolvedSortKeys[0].column);
    const sortRestored = await applySortKeys(
      resolvedSortKeys,
      visibleCol >= 0 ? visibleCol : Math.max(0, Math.min(columns.length - 1, selected.col)),
      "Restored sort",
      { stable: snapshot.sortStable !== false },
    );
    if (sortRestored) restored.push("sort");
  }
  document.body.dataset.refreshViewRestored = restored.join(",");
  return restored;
}

function startFileStatusChecks() {
  if (fileStatusTimer) clearInterval(fileStatusTimer);
  fileStatusTimer = window.setInterval(() => checkFileStatus(true), 5000);
}

function startSQLProgress() {
  if (activeSQLController) activeSQLController.abort();
  activeSQLController = new AbortController();
  if (sqlProgressTimer) window.clearInterval(sqlProgressTimer);
  sqlProgressStartedAt = performance.now();
  sqlProgressEventCount = 0;
  runSqlButton.disabled = true;
  cancelSqlButton.hidden = false;
  cancelSqlButton.disabled = false;
  document.body.dataset.sqlProgress = "running";
  document.body.dataset.sqlElapsedMs = "0";
  document.body.dataset.sqlProgressEvents = "0";
  document.body.dataset.sqlProgressPhase = "";
  document.body.dataset.sqlProgressRows = "0";
  document.body.dataset.sqlProgressSource = "";
  document.body.dataset.sqlLoadedRows = "0";
  document.body.dataset.sqlProgressSources = "0";
  updateSQLProgressStatus();
  sqlProgressTimer = window.setInterval(updateSQLProgressStatus, 250);
  return activeSQLController;
}

function updateSQLProgressStatus() {
  const elapsed = sqlElapsedMs();
  document.body.dataset.sqlElapsedMs = String(Math.round(elapsed));
  const phase = document.body.dataset.sqlProgressPhase || "";
  const rows = Number(document.body.dataset.sqlProgressRows || 0);
  status.textContent = phase && rows > 0
    ? tf("sql.status.runningRows", {
        phase: sqlProgressPhaseLabel(phase),
        rows: formatCount(rows),
        seconds: formatElapsedSeconds(elapsed),
      })
    : tf("sql.status.runningElapsed", { seconds: formatElapsedSeconds(elapsed) });
}

function sqlElapsedMs() {
  return sqlProgressStartedAt ? Math.max(0, performance.now() - sqlProgressStartedAt) : 0;
}

function stopSQLProgressTimer() {
  if (sqlProgressTimer) {
    window.clearInterval(sqlProgressTimer);
    sqlProgressTimer = 0;
  }
  const elapsed = sqlElapsedMs();
  document.body.dataset.sqlElapsedMs = String(Math.round(elapsed));
  return elapsed;
}

function finishSQLProgress(controller) {
  if (activeSQLController !== controller) return;
  stopSQLProgressTimer();
  activeSQLController = null;
  runSqlButton.disabled = false;
  cancelSqlButton.hidden = true;
  cancelSqlButton.disabled = true;
  if (document.body.dataset.sqlProgress === "running") {
    document.body.dataset.sqlProgress = "idle";
  }
}

function cancelSQLRequest() {
  if (!activeSQLController) return;
  activeSQLController.abort();
  const elapsed = stopSQLProgressTimer();
  activeSQLController = null;
  runSqlButton.disabled = false;
  cancelSqlButton.hidden = true;
  cancelSqlButton.disabled = true;
  document.body.dataset.sqlProgress = "canceled";
  status.textContent = tf("sql.status.canceledElapsed", { seconds: formatElapsedSeconds(elapsed) });
}

function sqlProgressPhaseLabel(phase) {
  return translationOrEmpty(`sql.progress.${phase}`) || phase || t("sql.status.running");
}

function applySQLProgressEvent(event) {
  if (!event || event.type !== "progress") return;
  const phase = String(event.phase || "");
  const rows = Math.max(0, Number(event.rows || 0));
  sqlProgressEventCount++;
  document.body.dataset.sqlProgressEvents = String(sqlProgressEventCount);
  document.body.dataset.sqlProgressPhase = phase;
  document.body.dataset.sqlProgressRows = String(rows);
  document.body.dataset.sqlLoadedRows = String(rows);
  document.body.dataset.sqlProgressSource = String(event.source || "");
  document.body.dataset.sqlProgressDone = event.done ? `${phase}:${rows}` : "";
  updateSQLProgressStatus();
}

function applySQLProgressSummary(progress) {
  if (!progress) return;
  document.body.dataset.sqlLoadedRows = String(progress.loadedRows ?? 0);
  document.body.dataset.sqlProgressSources = String(progress.sources ?? 0);
  document.body.dataset.sqlProgressEvents = String(progress.events ?? sqlProgressEventCount);
}

async function readSQLResponse(response) {
  const contentType = response.headers.get("Content-Type") || "";
  if (!contentType.includes("application/x-ndjson") || !response.body) {
    const payload = await response.json();
    return payload;
  }
  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = "";
  let result = null;
  let errorPayload = null;
  const consumeLine = (line) => {
    const trimmed = line.trim();
    if (!trimmed) return;
    const event = JSON.parse(trimmed);
    if (event.type === "progress") {
      applySQLProgressEvent(event);
    } else if (event.type === "result") {
      result = event.result;
    } else if (event.type === "error") {
      errorPayload = { error: event.error || "SQL failed", sqlError: event.sqlError };
    }
  };
  while (true) {
    const { value, done } = await reader.read();
    buffer += decoder.decode(value || new Uint8Array(), { stream: !done });
    const lines = buffer.split(/\r?\n/);
    buffer = lines.pop() ?? "";
    for (const line of lines) consumeLine(line);
    if (done) break;
  }
  consumeLine(buffer);
  if (errorPayload) return errorPayload;
  if (!result) return { error: "SQL stream ended without a result" };
  return result;
}

async function runSQLQuery() {
  if (!documentId) return;
  const query = sqlInput.value.trim();
  if (!query) {
    status.textContent = t("sql.status.enterQuery");
    return;
  }
  lastSQLResult = null;
  clearSQLErrorState();
  clearSQLResultSelection();
  syncSQLResultActions();
  const controller = startSQLProgress();
  let payload = null;
  try {
    const response = await fetch("/api/sql", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      signal: controller.signal,
      body: JSON.stringify({ documentId, query, limit: 500, streamProgress: true, progressInterval: 1000 }),
    });
    payload = await readSQLResponse(response);
    if (!response.ok || payload.error) {
      renderSQLError(payload);
      clearSQLResultSelection();
      syncSQLResultActions();
      return;
    }
  } catch (error) {
    if (error.name === "AbortError") {
      if (!activeSQLController || activeSQLController === controller) {
        const elapsed = stopSQLProgressTimer();
        document.body.dataset.sqlProgress = "canceled";
        status.textContent = tf("sql.status.canceledElapsed", { seconds: formatElapsedSeconds(elapsed) });
      }
      return;
    }
    renderSQLError({ error: error.message });
    clearSQLResultSelection();
    syncSQLResultActions();
    return;
  } finally {
    finishSQLProgress(controller);
  }
  clearSQLErrorState();
  lastSQLResult = payload;
  document.body.dataset.sqlStreamed = String(Boolean(payload.streamed));
  applySQLProgressSummary(payload.progress);
  renderSQLResult(payload);
  rememberSQLQuery(query);
  const elapsed = stopSQLProgressTimer();
  const statusKey = payload.progress?.loadedRows ? "sql.status.returnedTimedRows" : "sql.status.returnedTimed";
  status.textContent = tf(statusKey, {
    rows: formatCount(payload.rows.length),
    total: formatCount(payload.totalRows),
    loadedRows: formatCount(payload.progress?.loadedRows ?? 0),
    seconds: formatElapsedSeconds(elapsed),
  });
}

function sqlHistory() {
  try {
    const parsed = JSON.parse(localStorage.getItem(sqlHistoryStorageKey) || "[]");
    if (!Array.isArray(parsed)) return [];
    return parsed
      .filter((query) => typeof query === "string" && query.trim())
      .slice(0, sqlHistoryLimit);
  } catch {
    return [];
  }
}

function rememberSQLQuery(query) {
  const normalized = query.trim();
  if (!normalized) return;
  const next = [
    normalized,
    ...sqlHistory().filter((item) => item !== normalized),
  ].slice(0, sqlHistoryLimit);
  try {
    localStorage.setItem(sqlHistoryStorageKey, JSON.stringify(next));
  } catch {
    status.textContent = t("sql.status.historySaveFailed");
  }
  renderSQLHistory();
}

function clearSQLHistory() {
  try {
    localStorage.removeItem(sqlHistoryStorageKey);
  } catch {
    status.textContent = t("sql.status.historyClearFailed");
  }
  renderSQLHistory();
}

function renderSQLHistory() {
  const history = sqlHistory();
  sqlHistoryList.replaceChildren();
  clearSqlHistoryButton.disabled = history.length === 0;
  if (history.length === 0) {
    const empty = document.createElement("span");
    empty.className = "sql-history-empty";
    empty.textContent = t("sql.empty.history");
    sqlHistoryList.append(empty);
    return;
  }
  for (const query of history) {
    const button = document.createElement("button");
    button.type = "button";
    button.className = "sql-history-item";
    button.title = query;
    button.textContent = compactSQLQuery(query);
    button.addEventListener("click", () => {
      sqlInput.value = query;
      sqlInput.focus();
      status.textContent = t("sql.status.historyLoaded");
    });
    sqlHistoryList.append(button);
  }
}

function compactSQLQuery(query) {
  return query.replace(/\s+/g, " ").trim();
}

function clearSQLErrorState() {
  delete document.body.dataset.sqlErrorMessage;
  delete document.body.dataset.sqlErrorLocation;
  delete document.body.dataset.sqlErrorToken;
  delete document.body.dataset.sqlErrorSnippet;
  delete document.body.dataset.sqlErrorCaret;
  delete document.body.dataset.sqlErrorPhase;
}

function renderSQLError(payload = {}) {
  const detail = payload.sqlError && typeof payload.sqlError === "object" ? payload.sqlError : {};
  const message = detail.message || payload.error || "Query failed";
  const line = Number(detail.line);
  const column = Number(detail.column);
  const hasLocation = Number.isInteger(line) && line > 0 && Number.isInteger(column) && column > 0;
  status.textContent = hasLocation
    ? tf("sql.status.failedAt", { line, column, message })
    : tf("sql.status.failed", { message });

  document.body.dataset.sqlErrorMessage = message;
  document.body.dataset.sqlErrorLocation = hasLocation ? `${line}:${column}` : "";
  document.body.dataset.sqlErrorToken = detail.token || "";
  document.body.dataset.sqlErrorSnippet = detail.snippet || "";
  document.body.dataset.sqlErrorCaret = detail.caret || "";
  document.body.dataset.sqlErrorPhase = detail.phase || "";

  const panel = document.createElement("div");
  panel.className = "sql-error";
  panel.setAttribute("role", "alert");

  const title = document.createElement("strong");
  title.textContent = t("sql.error.title");
  const messageNode = document.createElement("p");
  messageNode.textContent = message;
  panel.append(title, messageNode);

  if (hasLocation || detail.phase || detail.token) {
    const meta = document.createElement("span");
    meta.className = "sql-error-meta";
    meta.textContent = [
      hasLocation ? tf("sql.error.location", { line, column }) : "",
      detail.token ? tf("sql.error.token", { token: detail.token }) : "",
      detail.phase ? tf("sql.error.phase", { phase: detail.phase }) : "",
    ].filter(Boolean).join(" · ");
    panel.append(meta);
  }

  if (detail.snippet) {
    const snippet = document.createElement("pre");
    snippet.className = "sql-error-snippet";
    snippet.textContent = detail.caret ? `${detail.snippet}\n${detail.caret}` : detail.snippet;
    panel.append(snippet);
  }

  sqlResult.replaceChildren(panel);
}

async function openProfilePanel() {
  if (!documentId || profileTabButton.disabled) return;
  setViewMode("profile");
  await loadProfile();
}

function parseProfileSchemaText(text) {
  const trimmed = text.trim();
  if (!trimmed) return null;
  const parsed = JSON.parse(trimmed);
  return normalizeProfileSchemaObject(parsed);
}

function activeProfileSchema() {
  return parseProfileSchemaText(profileSchemaText);
}

function normalizeProfileSchemaObject(value) {
  if (!value || typeof value !== "object") {
    throw new Error("Schema must be an object.");
  }
  const source = profileSchemaSourceObject(value);
  const columns = profileSchemaSourceColumns(source);
  if (!columns.length) {
    throw new Error("Schema must include columns or fields.");
  }
  const requiredNames = schemaNameSet(source.required);
  const primaryKeyNames = schemaNameSet(source.primaryKey ?? source.primaryKeys);
  const normalized = columns
    .map((column) => normalizeProfileSchemaColumn(column, requiredNames, primaryKeyNames))
    .filter(Boolean);
  if (!normalized.length) {
    throw new Error("Schema has no valid columns.");
  }
  return { columns: normalized };
}

function profileSchemaSourceObject(value) {
  if (Array.isArray(value.fields) || Array.isArray(value.columns) || isPlainSchemaObject(value.properties)) return value;
  if (value.schema && typeof value.schema === "object") return profileSchemaSourceObject(value.schema);
  if (value.tableSchema && typeof value.tableSchema === "object") return profileSchemaSourceObject(value.tableSchema);
  for (const collectionKey of ["tables", "resources"]) {
    if (!Array.isArray(value[collectionKey])) continue;
    for (const table of value[collectionKey]) {
      if (table?.schema || table?.tableSchema) return profileSchemaSourceObject(table.schema || table.tableSchema);
    }
  }
  return value;
}

function profileSchemaSourceColumns(source) {
  if (Array.isArray(source.columns)) return source.columns;
  if (Array.isArray(source.fields)) return source.fields;
  if (isPlainSchemaObject(source.properties)) {
    return Object.entries(source.properties)
      .sort(([left], [right]) => left.localeCompare(right))
      .map(([name, property]) => {
        if (!isPlainSchemaObject(property)) return null;
        const column = { ...property };
        if (!schemaString(column.name ?? column.title ?? column.label)) column.name = name;
        return column;
      })
      .filter(Boolean);
  }
  return [];
}

function isPlainSchemaObject(value) {
  return Boolean(value && typeof value === "object" && !Array.isArray(value));
}

function schemaNameSet(value) {
  if (typeof value === "string") return new Set([value]);
  if (!Array.isArray(value)) return new Set();
  return new Set(value.filter((item) => typeof item === "string" && item.trim()).map((item) => item.trim()));
}

function normalizeProfileSchemaColumn(column, requiredNames, primaryKeyNames) {
  if (!column || typeof column !== "object") return null;
  const constraints = column.constraints && typeof column.constraints === "object" ? column.constraints : {};
  const datatype = column.datatype ?? column.dataType ?? column.type;
  const datatypeObject = datatype && typeof datatype === "object" ? datatype : {};
  const name = String(column.name ?? column.title ?? column.label ?? "").trim();
  const columnNumber = schemaPositiveInteger(column.column ?? column.index ?? column.number);
  const rule = {};
  if (name) rule.name = name;
  if (columnNumber) rule.column = columnNumber;
  const type = normalizeProfileSchemaType(datatypeObject.base ?? datatypeObject.type ?? datatype ?? constraints.type);
  if (type) rule.type = type;
  const formatType =
    normalizeProfileSchemaFormat(column.format) ||
    normalizeProfileSchemaFormat(datatypeObject.format) ||
    normalizeProfileSchemaFormat(constraints.format);
  if (formatType && (!rule.type || rule.type === "string")) rule.type = formatType;
  if (Boolean(column.required ?? constraints.required) || requiredNames.has(name) || primaryKeyNames.has(name)) rule.required = true;
  if (Boolean(column.unique ?? constraints.unique) || primaryKeyNames.has(name)) rule.unique = true;
  const enumValues = schemaAllowedValues(column, constraints);
  if (enumValues.length) rule.enum = enumValues;
  const pattern = String(column.pattern ?? constraints.pattern ?? constraints.regex ?? "").trim();
  if (pattern) rule.pattern = pattern;
  const rangeSources = [column, constraints, datatypeObject];
  const exclusiveMin = schemaNumberFromSources(rangeSources, ["exclusiveMinimum", "exclusive_minimum", "exclusiveMin"]);
  const exclusiveMax = schemaNumberFromSources(rangeSources, ["exclusiveMaximum", "exclusive_maximum", "exclusiveMax"]);
  const min = exclusiveMin ?? schemaNumberFromSources(rangeSources, ["min", "minimum"]);
  const max = exclusiveMax ?? schemaNumberFromSources(rangeSources, ["max", "maximum"]);
  if (min != null) rule.min = min;
  if (max != null) rule.max = max;
  if (exclusiveMin != null || (min != null && schemaBoolFromSources(rangeSources, ["exclusiveMinimum", "exclusive_minimum", "exclusiveMin"]))) {
    rule.exclusiveMinimum = true;
  }
  if (exclusiveMax != null || (max != null && schemaBoolFromSources(rangeSources, ["exclusiveMaximum", "exclusive_maximum", "exclusiveMax"]))) {
    rule.exclusiveMaximum = true;
  }
  const minLength = schemaNonNegativeInteger(
    column.minLength ??
      column.minimumLength ??
      column.min_length ??
      column.minimum_length ??
      constraints.minLength ??
      constraints.minimumLength ??
      constraints.min_length ??
      constraints.minimum_length ??
      datatypeObject.minLength ??
      datatypeObject.minimumLength ??
      datatypeObject.min_length ??
      datatypeObject.minimum_length,
  );
  const maxLength = schemaNonNegativeInteger(
    column.maxLength ??
      column.maximumLength ??
      column.max_length ??
      column.maximum_length ??
      constraints.maxLength ??
      constraints.maximumLength ??
      constraints.max_length ??
      constraints.maximum_length ??
      datatypeObject.maxLength ??
      datatypeObject.maximumLength ??
      datatypeObject.max_length ??
      datatypeObject.maximum_length,
  );
  if (minLength != null) rule.minLength = minLength;
  if (maxLength != null) rule.maxLength = maxLength;
  const defaultValue = schemaString(column.default ?? column.defaultValue ?? column.default_value ?? constraints.default ?? constraints.defaultValue ?? constraints.default_value);
  if (defaultValue !== "") rule.default = defaultValue;
  return ruleHasValidation(rule) ? rule : null;
}

function schemaPositiveInteger(value) {
  const number = Number.parseInt(value, 10);
  return Number.isInteger(number) && number > 0 ? number : 0;
}

function schemaNumber(value) {
  if (value == null || value === "") return null;
  if (typeof value === "boolean") return null;
  const number = Number(value);
  return Number.isFinite(number) ? number : null;
}

function schemaNumberFromSources(sources, keys) {
  for (const source of sources) {
    if (!source || typeof source !== "object") continue;
    for (const key of keys) {
      const number = schemaNumber(source[key]);
      if (number != null) return number;
    }
  }
  return null;
}

function schemaBoolFromSources(sources, keys) {
  for (const source of sources) {
    if (!source || typeof source !== "object") continue;
    for (const key of keys) {
      if (schemaBool(source[key])) return true;
    }
  }
  return false;
}

function schemaBool(value) {
  if (value === true) return true;
  const text = schemaString(value).toLowerCase();
  return text === "true" || text === "1" || text === "yes";
}

function schemaNonNegativeInteger(value) {
  if (value == null || value === "") return null;
  const number = Number(value);
  return Number.isInteger(number) && number >= 0 ? number : null;
}

function schemaArrayValues(value) {
  if (Array.isArray(value)) {
    return value.map(schemaAllowedValue).filter(Boolean);
  }
  if (typeof value === "string") {
    return value.split(",").map((item) => item.trim()).filter(Boolean);
  }
  return [];
}

function schemaAllowedValues(column, constraints) {
  for (const source of [column, constraints]) {
    const value = schemaAllowedValue(source?.const ?? source?.constant);
    if (value) return [value];
  }
  return schemaArrayValues(column.enum ?? constraints.enum ?? constraints.oneOf ?? constraints.values ?? constraints.in);
}

function schemaAllowedValue(value) {
  if (value && typeof value === "object" && !Array.isArray(value)) {
    return schemaString(value.const ?? value.constant);
  }
  return schemaString(value);
}

function schemaString(value) {
  if (value == null) return "";
  return String(value).trim();
}

function normalizeProfileSchemaType(value) {
  if (Array.isArray(value)) {
    for (const item of value) {
      const type = normalizeProfileSchemaType(item);
      if (type) return type;
    }
    return "";
  }
  const type = normalizeProfileSchemaTypeName(value);
  if (!type || type === "null") return "";
  if (["integer", "int", "long", "short", "byte", "nonnegativeinteger", "nonpositiveinteger", "positiveinteger", "negativeinteger", "unsignedlong", "unsignedint", "unsignedshort", "unsignedbyte"].includes(type)) return "integer";
  if (["number", "numeric", "float", "double"].includes(type)) return "number";
  if (["decimal"].includes(type)) return "decimal";
  if (["boolean", "bool"].includes(type)) return "boolean";
  if (["date"].includes(type)) return "date";
  if (["datetime", "date-time", "timestamp"].includes(type)) return "datetime";
  if (["time"].includes(type)) return "time";
  if (["year", "gyear"].includes(type)) return "year";
  if (["yearmonth", "year-month", "gyearmonth"].includes(type)) return "yearmonth";
  if (["string", "text", "normalizedstring", "token", "uri", "url", "anyuri", "email"].includes(type)) return "string";
  return type;
}

function syncProfileSchemaControls() {
  if (!profileSchemaInput) return;
  if (profileSchemaInput.value !== profileSchemaText) profileSchemaInput.value = profileSchemaText;
  profileSchemaInput.placeholder = defaultProfileSchemaText;
  profileSchemaInput.hidden = !profileSchemaVisible;
  toggleProfileSchemaButton.setAttribute("aria-expanded", profileSchemaVisible ? "true" : "false");
  const text = profileSchemaText.trim();
  resetProfileSchemaButton.disabled = !text && !localStorage.getItem(profileSchemaStorageKey);
  saveProfileSchemaButton.disabled = !text;
  if (!text) {
    profileSchemaStatus.textContent = t("profile.noSchema");
    document.body.dataset.profileSchemaRules = "0";
    renderProfileRuleList([]);
    syncSavedProfileSchemaControls();
    return;
  }
  try {
    const schema = parseProfileSchemaText(text);
    const rules = schema.columns.length;
    profileSchemaStatus.textContent = tf("profile.schemaRules", { count: formatCount(rules) });
    document.body.dataset.profileSchemaRules = String(rules);
    renderProfileRuleList(schema.columns);
  } catch (error) {
    profileSchemaStatus.textContent = t("profile.invalidSchema");
    document.body.dataset.profileSchemaRules = "invalid";
    renderProfileRuleList([]);
  }
  syncSavedProfileSchemaControls();
}

function toggleProfileSchemaEditor() {
  profileSchemaVisible = !profileSchemaVisible;
  syncProfileSchemaControls();
}

function updateProfileSchemaDraft() {
  profileSchemaText = profileSchemaInput.value;
  profileData = null;
  syncProfileSchemaControls();
}

function saveProfileSchema() {
  let schema;
  try {
    schema = parseProfileSchemaText(profileSchemaText);
  } catch (error) {
    status.textContent = tf("profile.schema.statusInvalid", { message: error.message });
    syncProfileSchemaControls();
    return;
  }
  profileSchemaText = JSON.stringify(schema, null, 2);
  localStorage.setItem(profileSchemaStorageKey, profileSchemaText);
  const name = profileSchemaNameInput.value.trim();
  let message = t("profile.schema.saved");
  if (name) {
    const saved = saveNamedProfileSchema(name, profileSchemaText);
    message = tf("profile.schema.savedAs", { name: saved.name });
  }
  profileData = null;
  syncProfileSchemaControls();
  status.textContent = message;
  void loadProfile(true);
}

function resetProfileSchema() {
  localStorage.removeItem(profileSchemaStorageKey);
  profileSchemaText = "";
  profileData = null;
  syncProfileSchemaControls();
  status.textContent = t("profile.schema.reset");
  void loadProfile(true);
}

function loadSavedProfileSchemas() {
  try {
    const parsed = JSON.parse(localStorage.getItem(profileSchemasStorageKey) || "[]");
    const rawSchemas = Array.isArray(parsed) ? parsed : parsed?.schemas;
    if (!Array.isArray(rawSchemas)) return [];
    return uniqueSavedProfileSchemas(rawSchemas.map(normalizeSavedProfileSchema).filter(Boolean));
  } catch {
    return [];
  }
}

function normalizeSavedProfileSchema(item) {
  if (!item || typeof item !== "object") return null;
  const schemaText = typeof item.schemaText === "string"
    ? item.schemaText
    : (item.schema ? JSON.stringify(item.schema, null, 2) : "");
  if (!schemaText.trim()) return null;
  let schema;
  try {
    schema = parseProfileSchemaText(schemaText);
  } catch {
    return null;
  }
  const name = normalizeProfileSchemaName(item.name || item.title || "Validation Schema");
  return {
    id: typeof item.id === "string" && item.id ? item.id : profileSchemaID(name),
    name,
    schemaText: JSON.stringify(schema, null, 2),
    updatedAt: typeof item.updatedAt === "string" && item.updatedAt ? item.updatedAt : new Date().toISOString(),
  };
}

function normalizeProfileSchemaName(name) {
  return String(name || "").replace(/\s+/g, " ").trim().slice(0, 80) || "Validation Schema";
}

function profileSchemaID(name) {
  const slug = normalizeProfileSchemaName(name)
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .slice(0, 32) || "schema";
  return `${slug}-${Date.now().toString(36)}`;
}

function uniqueSavedProfileSchemas(schemas) {
  const unique = [];
  const seen = new Set();
  for (const schema of schemas) {
    const key = schema.name.toLowerCase();
    if (seen.has(key)) continue;
    seen.add(key);
    unique.push(schema);
  }
  return unique.slice(0, profileSchemasLimit);
}

function persistSavedProfileSchemas() {
  const payload = {
    kind: profileSchemasExportKind,
    version: 1,
    schemas: savedProfileSchemasState,
  };
  localStorage.setItem(profileSchemasStorageKey, JSON.stringify(payload));
  document.body.dataset.profileSavedSchemaCount = String(savedProfileSchemasState.length);
}

function saveNamedProfileSchema(name, schemaText) {
  const schema = parseProfileSchemaText(schemaText);
  const normalizedName = normalizeProfileSchemaName(name);
  const existing = savedProfileSchemasState.find((item) => item.name.toLowerCase() === normalizedName.toLowerCase());
  const saved = {
    id: existing?.id || profileSchemaID(normalizedName),
    name: normalizedName,
    schemaText: JSON.stringify(schema, null, 2),
    updatedAt: new Date().toISOString(),
  };
  savedProfileSchemasState = uniqueSavedProfileSchemas([
    saved,
    ...savedProfileSchemasState.filter((item) => item.id !== saved.id && item.name.toLowerCase() !== normalizedName.toLowerCase()),
  ]);
  persistSavedProfileSchemas();
  savedProfileSchemaSelect.value = saved.id;
  document.body.dataset.profileSchemaSaved = saved.name;
  syncSavedProfileSchemaControls(saved.id);
  return saved;
}

function syncSavedProfileSchemaControls(preferredID = "") {
  if (!savedProfileSchemaSelect) return;
  const selectedID = preferredID || savedProfileSchemaSelect.value;
  savedProfileSchemaSelect.replaceChildren(profileSchemaOption("", t("profile.savedSchemas")));
  for (const schema of savedProfileSchemasState) {
    savedProfileSchemaSelect.append(profileSchemaOption(schema.id, schema.name));
  }
  if (savedProfileSchemasState.some((schema) => schema.id === selectedID)) {
    savedProfileSchemaSelect.value = selectedID;
  }
  const hasSelection = savedProfileSchemasState.some((schema) => schema.id === savedProfileSchemaSelect.value);
  applyProfileSchemaButton.disabled = !hasSelection;
  deleteProfileSchemaButton.disabled = !hasSelection;
  exportProfileSchemasButton.disabled = savedProfileSchemasState.length === 0;
  applyProfileSchemaImportButton.disabled = profileSchemaImportText.value.trim() === "";
  document.body.dataset.profileSavedSchemaCount = String(savedProfileSchemasState.length);
}

function profileSchemaOption(value, label) {
  const option = document.createElement("option");
  option.value = value;
  option.textContent = label;
  return option;
}

function applySavedProfileSchema() {
  const schema = savedProfileSchemasState.find((item) => item.id === savedProfileSchemaSelect.value);
  if (!schema) {
    status.textContent = t("profile.schema.notFound");
    syncSavedProfileSchemaControls();
    return;
  }
  profileSchemaText = schema.schemaText;
  profileSchemaVisible = true;
  profileSchemaNameInput.value = schema.name;
  localStorage.setItem(profileSchemaStorageKey, profileSchemaText);
  profileData = null;
  document.body.dataset.profileSchemaApplied = schema.name;
  syncProfileSchemaControls();
  status.textContent = tf("profile.schema.applied", { name: schema.name });
  void loadProfile(true);
}

function deleteSavedProfileSchema() {
  const selectedID = savedProfileSchemaSelect.value;
  const before = savedProfileSchemasState.length;
  savedProfileSchemasState = savedProfileSchemasState.filter((schema) => schema.id !== selectedID);
  persistSavedProfileSchemas();
  savedProfileSchemaSelect.value = "";
  syncSavedProfileSchemaControls();
  status.textContent = before === savedProfileSchemasState.length
    ? t("profile.schema.notFound")
    : t("profile.schema.deleted");
}

function exportProfileSchemas() {
  const payload = {
    kind: profileSchemasExportKind,
    version: 1,
    exportedAt: new Date().toISOString(),
    schemas: savedProfileSchemasState,
  };
  const text = `${JSON.stringify(payload, null, 2)}\n`;
  window.__opencsvLastProfileSchemaExportText = text;
  document.body.dataset.profileSchemaExport = String(savedProfileSchemasState.length);
  document.body.dataset.profileSchemaExportPreview = text.slice(0, 500);
  triggerDownload(text, "opencsv-profile-schemas.json", "application/json;charset=utf-8");
  status.textContent = tf("profile.schema.exported", { count: formatCount(savedProfileSchemasState.length) });
}

function showProfileSchemaImportPanel() {
  profileSchemaImportPanel.hidden = false;
  profileSchemaImportText.value = "";
  document.body.dataset.profileSchemaImportPanel = "open";
  syncSavedProfileSchemaControls();
  profileSchemaImportText.focus();
}

function cancelProfileSchemaImport() {
  profileSchemaImportPanel.hidden = true;
  profileSchemaImportText.value = "";
  document.body.dataset.profileSchemaImportPanel = "closed";
  syncSavedProfileSchemaControls();
}

function applyProfileSchemaImport() {
  try {
    const imported = parseProfileSchemaImport(profileSchemaImportText.value);
    savedProfileSchemasState = uniqueSavedProfileSchemas([...imported, ...savedProfileSchemasState]);
    persistSavedProfileSchemas();
    profileSchemaImportPanel.hidden = true;
    profileSchemaImportText.value = "";
    document.body.dataset.profileSchemaImport = String(imported.length);
    document.body.dataset.profileSchemaImportPanel = "closed";
    syncSavedProfileSchemaControls();
    status.textContent = tf("profile.schema.imported", { count: formatCount(imported.length) });
  } catch (error) {
    status.textContent = tf("profile.schema.importFailed", { message: error.message });
    document.body.dataset.profileSchemaImport = "error";
  }
}

function parseProfileSchemaImport(text) {
  let parsed;
  try {
    parsed = JSON.parse(text);
  } catch {
    throw new Error("invalid JSON");
  }
  const rawSchemas = Array.isArray(parsed) ? parsed : parsed?.schemas;
  let directSchema = [];
  if (!Array.isArray(rawSchemas)) {
    try {
      directSchema = [{ name: parsed.name || parsed.title || "Imported Schema", schema: normalizeProfileSchemaObject(parsed) }];
    } catch {
      directSchema = [];
    }
  }
  const schemas = Array.isArray(rawSchemas) ? rawSchemas : directSchema;
  if (!Array.isArray(schemas)) throw new Error("missing schemas array");
  const normalized = schemas.map(normalizeSavedProfileSchema).filter(Boolean);
  if (normalized.length === 0) throw new Error("no valid schemas");
  return normalized;
}

function schemaDraftFromText() {
  try {
    return parseProfileSchemaText(profileSchemaText) || { columns: [] };
  } catch {
    return { columns: [] };
  }
}

function profileRuleFromBuilder() {
  const rule = {};
  const name = profileRuleNameInput.value.trim();
  const column = Number.parseInt(profileRuleColumnInput.value, 10);
  if (name) rule.name = name;
  if (Number.isInteger(column) && column > 0) rule.column = column;
  if (profileRuleRequiredToggle.checked) rule.required = true;
  if (profileRuleUniqueToggle.checked) rule.unique = true;
  if (profileRuleTypeSelect.value) rule.type = profileRuleTypeSelect.value;
  const enumValues = profileRuleEnumInput.value
    .split(",")
    .map((value) => value.trim())
    .filter(Boolean);
  if (enumValues.length) rule.enum = enumValues;
  const pattern = profileRulePatternInput.value.trim();
  if (pattern) rule.pattern = pattern;
  const defaultValue = profileRuleDefaultInput.value.trim();
  if (defaultValue) rule.default = defaultValue;
  const minValue = Number.parseFloat(profileRuleMinInput.value);
  const maxValue = Number.parseFloat(profileRuleMaxInput.value);
  const minLengthValue = Number.parseInt(profileRuleMinLengthInput.value, 10);
  const maxLengthValue = Number.parseInt(profileRuleMaxLengthInput.value, 10);
  if (profileRuleMinInput.value.trim() !== "" && Number.isFinite(minValue)) rule.min = minValue;
  if (profileRuleMaxInput.value.trim() !== "" && Number.isFinite(maxValue)) rule.max = maxValue;
  if (profileRuleMinLengthInput.value.trim() !== "" && Number.isInteger(minLengthValue) && minLengthValue >= 0) {
    rule.minLength = minLengthValue;
  }
  if (profileRuleMaxLengthInput.value.trim() !== "" && Number.isInteger(maxLengthValue) && maxLengthValue >= 0) {
    rule.maxLength = maxLengthValue;
  }
  if (!rule.name && !rule.column) {
    throw new Error(t("profile.rule.chooseColumn"));
  }
  return rule;
}

function normalizeProfileSchemaFormat(value) {
  const format = normalizeProfileSchemaTypeName(value);
  if (["date-time", "datetime", "timestamp"].includes(format)) return "datetime";
  if (["date", "full-date"].includes(format)) return "date";
  if (["time", "full-time"].includes(format)) return "time";
  return "";
}

function normalizeProfileSchemaTypeName(value) {
  let type = String(value ?? "").toLowerCase().trim();
  if (!type) return "";
  const fragmentIndex = type.lastIndexOf("#");
  const slashIndex = type.lastIndexOf("/");
  if (fragmentIndex >= 0) {
    type = type.slice(fragmentIndex + 1);
  } else if (slashIndex >= 0) {
    type = type.slice(slashIndex + 1);
  }
  return type.replace(/^xsd:/, "").replace(/^xmlschema:/, "");
}

function ruleHasValidation(rule) {
  return Boolean(
    rule.name ||
      rule.column ||
      rule.required ||
      rule.unique ||
      rule.type ||
      rule.pattern ||
      rule.default ||
      rule.min != null ||
      rule.max != null ||
      rule.minLength != null ||
      rule.maxLength != null ||
      (Array.isArray(rule.enum) && rule.enum.length),
  );
}

function saveProfileRuleFromBuilder() {
  let rule;
  try {
    rule = profileRuleFromBuilder();
  } catch (error) {
    status.textContent = tf("profile.rule.invalid", { message: error.message });
    return;
  }
  const schema = schemaDraftFromText();
  const rules = Array.isArray(schema.columns) ? [...schema.columns] : [];
  if (editingProfileRuleIndex >= 0 && editingProfileRuleIndex < rules.length) {
    rules[editingProfileRuleIndex] = rule;
  } else {
    rules.push(rule);
  }
  schema.columns = rules.filter(ruleHasValidation);
  profileSchemaText = JSON.stringify(schema, null, 2);
  localStorage.setItem(profileSchemaStorageKey, profileSchemaText);
  profileData = null;
  document.body.dataset.profileRuleBuilder = `${editingProfileRuleIndex >= 0 ? "update" : "add"}:${schema.columns.length}`;
  const message = editingProfileRuleIndex >= 0 ? t("profile.rule.updated") : t("profile.rule.added");
  clearProfileRuleBuilder(false);
  syncProfileSchemaControls();
  status.textContent = message;
  void loadProfile(true);
}

function clearProfileRuleBuilder(resetStatus = true) {
  editingProfileRuleIndex = -1;
  profileRuleNameInput.value = "";
  profileRuleColumnInput.value = "";
  profileRuleTypeSelect.value = "";
  profileRuleRequiredToggle.checked = false;
  profileRuleUniqueToggle.checked = false;
  profileRuleEnumInput.value = "";
  profileRulePatternInput.value = "";
  profileRuleDefaultInput.value = "";
  profileRuleMinInput.value = "";
  profileRuleMaxInput.value = "";
  profileRuleMinLengthInput.value = "";
  profileRuleMaxLengthInput.value = "";
  saveProfileRuleButton.textContent = t("profile.rule.save");
  if (resetStatus) status.textContent = t("profile.rule.cleared");
}

function editProfileRule(index) {
  const schema = schemaDraftFromText();
  const rule = schema.columns?.[index];
  if (!rule) {
    status.textContent = t("profile.rule.notFound");
    return;
  }
  editingProfileRuleIndex = index;
  profileRuleNameInput.value = rule.name || "";
  profileRuleColumnInput.value = rule.column || "";
  profileRuleTypeSelect.value = rule.type || "";
  profileRuleRequiredToggle.checked = Boolean(rule.required);
  profileRuleUniqueToggle.checked = Boolean(rule.unique);
  profileRuleEnumInput.value = Array.isArray(rule.enum) ? rule.enum.join(", ") : "";
  profileRulePatternInput.value = rule.pattern || "";
  profileRuleDefaultInput.value = rule.default || "";
  profileRuleMinInput.value = rule.min ?? "";
  profileRuleMaxInput.value = rule.max ?? "";
  profileRuleMinLengthInput.value = rule.minLength ?? "";
  profileRuleMaxLengthInput.value = rule.maxLength ?? "";
  saveProfileRuleButton.textContent = t("profile.rule.update");
  document.body.dataset.profileRuleBuilder = `edit:${index}`;
  status.textContent = tf("profile.rule.editing", { index: index + 1 });
}

function removeProfileRule(index) {
  const schema = schemaDraftFromText();
  const rules = Array.isArray(schema.columns) ? [...schema.columns] : [];
  if (index < 0 || index >= rules.length) {
    status.textContent = t("profile.rule.notFound");
    return;
  }
  rules.splice(index, 1);
  schema.columns = rules;
  profileSchemaText = JSON.stringify(schema, null, 2);
  localStorage.setItem(profileSchemaStorageKey, profileSchemaText);
  profileData = null;
  clearProfileRuleBuilder(false);
  document.body.dataset.profileRuleBuilder = `remove:${rules.length}`;
  syncProfileSchemaControls();
  status.textContent = t("profile.rule.removed");
  void loadProfile(true);
}

function renderProfileRuleList(rules) {
  if (!profileRuleList) return;
  profileRuleList.replaceChildren();
  if (!rules.length) {
    const empty = document.createElement("div");
    empty.className = "profile-rule-empty";
    empty.textContent = t("profile.rule.none");
    profileRuleList.append(empty);
    document.body.dataset.profileRuleListCount = "0";
    return;
  }
  for (const [index, rule] of rules.entries()) {
    const item = document.createElement("div");
    item.className = "profile-rule-item";
    item.dataset.ruleIndex = String(index);
    const label = document.createElement("strong");
    label.textContent = profileRuleTitle(rule, index);
    const detail = document.createElement("span");
    detail.textContent = profileRuleSummary(rule);
    const actions = document.createElement("div");
    actions.className = "profile-rule-item-actions";
    const editButton = document.createElement("button");
    editButton.type = "button";
    editButton.textContent = t("profile.rule.edit");
    editButton.addEventListener("click", () => editProfileRule(index));
    const removeButton = document.createElement("button");
    removeButton.type = "button";
    removeButton.textContent = t("profile.rule.remove");
    removeButton.addEventListener("click", () => removeProfileRule(index));
    actions.append(editButton, removeButton);
    item.append(label, detail, actions);
    profileRuleList.append(item);
  }
  document.body.dataset.profileRuleListCount = String(rules.length);
}

function profileRuleTitle(rule, index) {
  if (rule.name) return rule.name;
  if (rule.column) return `${t("profile.columns.column")} ${rule.column}`;
  return tf("profile.rule.rule", { index: index + 1 });
}

function profileRuleSummary(rule) {
  const parts = [];
  if (rule.required) parts.push(t("profile.rule.requiredSummary"));
  if (rule.unique) parts.push(t("profile.rule.uniqueSummary"));
  if (rule.type) parts.push(tf("profile.rule.typeSummary", { type: profileTypeLabel(rule.type) }));
  if (Array.isArray(rule.enum) && rule.enum.length) parts.push(tf("profile.rule.enumSummary", { values: rule.enum.join("|") }));
  if (rule.pattern) parts.push(tf("profile.rule.patternSummary", { pattern: rule.pattern }));
  if (rule.min != null) {
    parts.push(tf(rule.exclusiveMinimum ? "profile.rule.exclusiveMinSummary" : "profile.rule.minSummary", { value: rule.min }));
  }
  if (rule.max != null) {
    parts.push(tf(rule.exclusiveMaximum ? "profile.rule.exclusiveMaxSummary" : "profile.rule.maxSummary", { value: rule.max }));
  }
  if (rule.minLength != null) parts.push(tf("profile.rule.minLengthSummary", { value: rule.minLength }));
  if (rule.maxLength != null) parts.push(tf("profile.rule.maxLengthSummary", { value: rule.maxLength }));
  if (rule.default != null && String(rule.default).trim() !== "") {
    parts.push(tf("profile.rule.defaultSummary", { value: rule.default }));
  }
  return parts.length ? parts.join(", ") : t("profile.rule.locationOnly");
}

function profileTypeLabel(type) {
  return translationOrEmpty(`profile.type.${type}`) || String(type || "");
}

function startProfileProgress() {
  if (activeProfileController) activeProfileController.abort();
  activeProfileController = new AbortController();
  if (profileProgressTimer) window.clearInterval(profileProgressTimer);
  profileProgressStartedAt = performance.now();
  profileProgressEventCount = 0;
  profileLoading = true;
  runProfileButton.disabled = true;
  cancelProfileButton.hidden = false;
  cancelProfileButton.disabled = false;
  document.body.dataset.profileProgress = "running";
  document.body.dataset.profileElapsedMs = "0";
  document.body.dataset.profileProgressEvents = "0";
  document.body.dataset.profileProgressPhase = "";
  document.body.dataset.profileProgressRows = "0";
  document.body.dataset.profileProfileRows = "0";
  document.body.dataset.profileValidationRows = "0";
  updateProfileProgressStatus();
  profileProgressTimer = window.setInterval(updateProfileProgressStatus, 250);
  return activeProfileController;
}

function updateProfileProgressStatus() {
  const elapsed = profileElapsedMs();
  document.body.dataset.profileElapsedMs = String(Math.round(elapsed));
  const phase = document.body.dataset.profileProgressPhase || "";
  const rows = Number(document.body.dataset.profileProgressRows || 0);
  status.textContent = phase && rows > 0
    ? tf("profile.status.runningRows", {
        phase: profileProgressPhaseLabel(phase),
        rows: formatCount(rows),
        seconds: formatElapsedSeconds(elapsed),
      })
    : tf("profile.status.runningElapsed", { seconds: formatElapsedSeconds(elapsed) });
}

function profileElapsedMs() {
  return profileProgressStartedAt ? Math.max(0, performance.now() - profileProgressStartedAt) : 0;
}

function stopProfileProgressTimer() {
  if (profileProgressTimer) {
    window.clearInterval(profileProgressTimer);
    profileProgressTimer = 0;
  }
  const elapsed = profileElapsedMs();
  document.body.dataset.profileElapsedMs = String(Math.round(elapsed));
  return elapsed;
}

function finishProfileProgress(controller) {
  if (activeProfileController !== controller) return;
  stopProfileProgressTimer();
  activeProfileController = null;
  profileLoading = false;
  runProfileButton.disabled = false;
  cancelProfileButton.hidden = true;
  cancelProfileButton.disabled = true;
  if (document.body.dataset.profileProgress === "running") {
    document.body.dataset.profileProgress = "idle";
  }
}

function cancelProfileRequest() {
  if (!activeProfileController) return;
  activeProfileController.abort();
  const elapsed = stopProfileProgressTimer();
  activeProfileController = null;
  profileLoading = false;
  runProfileButton.disabled = false;
  cancelProfileButton.hidden = true;
  cancelProfileButton.disabled = true;
  document.body.dataset.profileProgress = "canceled";
  status.textContent = tf("profile.status.canceledElapsed", { seconds: formatElapsedSeconds(elapsed) });
}

function profileProgressPhaseLabel(phase) {
  return translationOrEmpty(`profile.progress.${phase}`) || phase || t("profile.loading");
}

function applyProfileProgressEvent(event) {
  if (!event || event.type !== "progress") return;
  const phase = String(event.phase || "");
  const rows = Math.max(0, Number(event.rows || 0));
  profileProgressEventCount++;
  document.body.dataset.profileProgressEvents = String(profileProgressEventCount);
  document.body.dataset.profileProgressPhase = phase;
  document.body.dataset.profileProgressRows = String(rows);
  if (phase === "profile") document.body.dataset.profileProfileRows = String(rows);
  if (phase === "validation") document.body.dataset.profileValidationRows = String(rows);
  document.body.dataset.profileProgressDone = event.done ? `${phase}:${rows}` : "";
  updateProfileProgressStatus();
}

function applyProfileProgressSummary(progress) {
  if (!progress) return;
  document.body.dataset.profileProfileRows = String(progress.profileRows ?? 0);
  document.body.dataset.profileValidationRows = String(progress.validationRows ?? 0);
  document.body.dataset.profileProgressTotalRows = String(progress.totalRows ?? 0);
  document.body.dataset.profileProgressPasses = String(progress.passes ?? 0);
  document.body.dataset.profileProgressEvents = String(progress.events ?? profileProgressEventCount);
}

async function readProfileResponse(response) {
  const contentType = response.headers.get("Content-Type") || "";
  if (!contentType.includes("application/x-ndjson") || !response.body) {
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "profile failed");
    return payload;
  }
  if (!response.ok) {
    const payload = await response.json();
    throw new Error(payload.error ?? "profile failed");
  }
  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = "";
  let result = null;
  const consumeLine = (line) => {
    const trimmed = line.trim();
    if (!trimmed) return;
    const event = JSON.parse(trimmed);
    if (event.type === "progress") {
      applyProfileProgressEvent(event);
    } else if (event.type === "result") {
      result = event.result;
    } else if (event.type === "error") {
      throw new Error(event.error || "profile failed");
    }
  };
  while (true) {
    const { value, done } = await reader.read();
    buffer += decoder.decode(value || new Uint8Array(), { stream: !done });
    const lines = buffer.split(/\r?\n/);
    buffer = lines.pop() ?? "";
    for (const line of lines) consumeLine(line);
    if (done) break;
  }
  consumeLine(buffer);
  if (!result) throw new Error("profile stream ended without a result");
  return result;
}

async function loadProfile(force = false) {
  if (!documentId || profileLoading) return;
  if (profileData && !force) {
    renderProfile(profileData);
    return;
  }
  const controller = startProfileProgress();
  profileSummary.replaceChildren(summaryNode(t("profile.summary.status"), t("profile.loading")));
  profileColumns.replaceChildren(emptyProfileNode(t("profile.columns.loading")));
  profileIssues.replaceChildren(emptyProfileNode(t("profile.issues.loading")));
  try {
    const schema = activeProfileSchema();
    syncProfileSchemaControls();
    const response = await fetch("/api/profile", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      signal: controller.signal,
      body: JSON.stringify({ documentId, schema, streamProgress: true, progressInterval: 1000 }),
    });
    const payload = await readProfileResponse(response);
    profileData = payload;
    document.body.dataset.profileStreamed = String(Boolean(payload.streamed));
    applyProfileProgressSummary(payload.progress);
    renderProfile(payload);
    updateInspector();
    const elapsed = stopProfileProgressTimer();
    const progress = payload.progress || {};
    const statusKey = progress.profileRows || progress.validationRows
      ? "profile.status.completeTimedRows"
      : "profile.status.completeTimed";
    status.textContent = tf(statusKey, {
      rows: formatCount(payload.rows ?? 0),
      columns: formatCount(payload.columns?.length ?? 0),
      issues: formatCount(payload.issues?.length ?? 0),
      seconds: formatElapsedSeconds(elapsed),
      profileRows: formatCount(progress.profileRows ?? 0),
      validationRows: formatCount(progress.validationRows ?? 0),
    });
  } catch (error) {
    if (error.name === "AbortError") {
      if (!activeProfileController || activeProfileController === controller) {
        const elapsed = stopProfileProgressTimer();
        document.body.dataset.profileProgress = "canceled";
        status.textContent = tf("profile.status.canceledElapsed", { seconds: formatElapsedSeconds(elapsed) });
      }
      return;
    }
    status.textContent = tf("profile.status.failed", { message: error.message });
    profileColumns.replaceChildren(emptyProfileNode(error.message));
    profileIssues.replaceChildren(emptyProfileNode(t("profile.empty.noValidation")));
  } finally {
    finishProfileProgress(controller);
  }
}

function renderProfile(result) {
  syncProfileSchemaControls();
  const issues = result.issues ?? [];
  const errorCount = issues.filter((issue) => issue.severity === "error").length;
  const warningCount = issues.filter((issue) => issue.severity === "warning").length;
  const summaryNodes = [
    summaryNode(t("profile.summary.rows"), result.rows ?? 0),
    summaryNode(t("profile.summary.columns"), result.columns?.length ?? 0),
    summaryNode(t("profile.summary.errors"), errorCount),
    summaryNode(t("profile.summary.warnings"), warningCount),
  ];
  if (result.schemaRules) summaryNodes.push(summaryNode(t("profile.summary.rules"), result.schemaRules));
  if (result.progress?.passes) {
    summaryNodes.push(summaryNode(t("profile.summary.progress"), tf("profile.summary.progressValue", {
      passes: formatCount(result.progress.passes),
      rows: formatCount((result.progress.profileRows ?? 0) + (result.progress.validationRows ?? 0)),
    })));
  }
  profileSummary.replaceChildren(...summaryNodes);
  document.body.dataset.profileIssueCount = String(issues.length);
  document.body.dataset.profileSchemaIssueCount = String(issues.filter((issue) => String(issue.type || "").startsWith("schema-")).length);
  applyProfileProgressSummary(result.progress);
  renderProfileColumns(result.columns ?? []);
  renderProfileIssues(issues);
}

function renderProfileColumns(columnsProfile) {
  if (!columnsProfile.length) {
    profileColumns.replaceChildren(emptyProfileNode(t("profile.empty.noColumns")));
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const headerRow = document.createElement("tr");
  for (const key of [
    "profile.columns.column",
    "profile.columns.type",
    "profile.columns.rows",
    "profile.columns.nulls",
    "profile.columns.distinct",
    "profile.columns.min",
    "profile.columns.max",
    "profile.columns.average",
    "profile.columns.topValues",
  ]) {
    const th = document.createElement("th");
    th.textContent = t(key);
    headerRow.append(th);
  }
  thead.append(headerRow);
  table.append(thead);

  const tbody = document.createElement("tbody");
  for (const profile of columnsProfile) {
    const tr = document.createElement("tr");
    tr.dataset.columnIndex = String(profile.index);
    if (profile.index === sourceColFor(selected.col)) tr.className = "selected-profile-column";
    appendCell(tr, profile.name || `c${profile.index + 1}`);
    appendCell(tr, profileTypeLabel(profile.type) || t("profile.type.unknown"));
    appendCell(tr, formatCount(profile.rows ?? 0));
    appendCell(tr, `${formatCount(profile.nulls ?? 0)} (${formatPercent(profile.nullRatio ?? 0)})`);
    appendCell(tr, formatCount(profile.distinct ?? 0));
    appendCell(tr, profile.min ?? "");
    appendCell(tr, profile.max ?? "");
    appendCell(tr, profile.average == null ? "" : formatNumber(profile.average));
    appendCell(
      tr,
      (profile.topValues ?? [])
        .map((item) => `${item.value || t("profile.value.empty")} x${formatCount(item.count)}`)
        .join("; "),
    );
    tr.addEventListener("click", () => {
      const visibleCol = visibleColForSource(profile.index);
      selected = { row: selected.row, col: visibleCol >= 0 ? visibleCol : selected.col };
      renderProfileColumns(columnsProfile);
      updateInspector();
    });
    tbody.append(tr);
  }
  table.append(tbody);
  profileColumns.replaceChildren(table);
}

function renderProfileIssues(issues) {
  if (!issues.length) {
    profileIssues.replaceChildren(emptyProfileNode(t("profile.empty.noIssues")));
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const headerRow = document.createElement("tr");
  for (const key of [
    "profile.issues.severity",
    "profile.issues.type",
    "profile.issues.row",
    "profile.issues.column",
    "profile.issues.message",
    "profile.issues.fix",
  ]) {
    const th = document.createElement("th");
    th.textContent = t(key);
    headerRow.append(th);
  }
  thead.append(headerRow);
  table.append(thead);

  const tbody = document.createElement("tbody");
  for (const issue of issues) {
    const tr = document.createElement("tr");
    tr.className = issue.severity === "error" ? "issue-error" : "issue-warning";
    appendCell(tr, profileSeverityLabel(issue.severity));
    appendCell(tr, issue.type ?? "");
    appendCell(tr, issue.row || "");
    appendCell(tr, issue.col || "");
    appendCell(tr, issue.message ?? "");
    const fixCell = document.createElement("td");
    for (const fix of issue.quickFixes ?? []) {
      const button = document.createElement("button");
      button.type = "button";
      button.textContent = profileQuickFixLabel(fix);
      button.addEventListener("click", (event) => {
        event.stopPropagation();
        void applyProfileQuickFix(issue, fix);
      });
      fixCell.append(button);
    }
    tr.append(fixCell);
    tr.addEventListener("click", () => {
      void navigateToProfileIssue(issue);
    });
    tbody.append(tr);
  }
  table.append(tbody);
  profileIssues.replaceChildren(table);
}

function issueDisplayRow(issue) {
  if (!issue?.row || sourceRowCount <= 0) return Math.max(0, Math.min(sourceRowCount - 1, selected.row));
  const sourceRow = Math.max(0, Number(issue.row) - 1);
  const mapped = displayRowForSourceRow(sourceRow);
  if (mapped >= 0) return Math.max(0, Math.min(sourceRowCount - 1, mapped));
  return Math.max(0, Math.min(sourceRowCount - 1, sourceRow - dataStartRow));
}

function issueVisibleCol(issue) {
  if (!issue?.col || columns.length === 0) return Math.max(0, Math.min(columns.length - 1, selected.col));
  const visibleCol = visibleColForSource(Number(issue.col) - 1);
  return Math.max(0, Math.min(columns.length - 1, visibleCol >= 0 ? visibleCol : selected.col));
}

function profileSeverityLabel(severity) {
  return translationOrEmpty(`profile.severity.${severity}`) || String(severity || "warning");
}

function profileQuickFixLabel(fix = {}) {
  if (fix.id) {
    const localized = translationOrEmpty(`profile.quickFix.${fix.id}`);
    if (localized) return localized;
  }
  return fix.label ?? t("profile.issues.defaultFix");
}

async function navigateToProfileIssue(issue) {
  if (!documentId || sourceRowCount <= 0 || columns.length === 0) return;
  const row = issueDisplayRow(issue);
  const col = issueVisibleCol(issue);
  selected = { row, col };
  selectionAnchor = { ...selected };
  clearSmartSelectionRanges();
  scrollRangeIntoView({ startRow: row, endRow: row, startCol: col, endCol: col });
  setViewMode("grid");
  await maybeLoadWindow(row);
  draw();
  document.body.dataset.profileIssueSelection = `${issue.type ?? "issue"}:${issue.row || ""}:${issue.col || ""}`;
  status.textContent = tf("profile.status.selectedIssue", {
    type: issue.type ?? "issue",
    row: issue.row || "-",
    column: issue.col || "-",
  });
}

function profileQuickFixEdits(issue, fix) {
  const fixID = typeof fix === "string" ? fix : fix?.id;
  if (fixID === "header.fill" && issue?.col) {
    const col = Number(issue.col) - 1;
    return [{ row: 0, col, value: `Column ${Number(issue.col)}` }];
  }
  if (fixID === "header.unique" && issue?.col) {
    const col = Number(issue.col) - 1;
    return [{ row: 0, col, value: uniqueHeaderNameForSourceCol(col) }];
  }
  if (fixID === "row.pad" && issue?.row && issue?.expected && issue.actual < issue.expected) {
    const row = Number(issue.row) - 1;
    const edits = [];
    for (let col = Number(issue.actual); col < Number(issue.expected); col++) {
      edits.push({ row, col, value: "" });
    }
    return edits;
  }
  if (
    (fixID === "schema.range.clampMin" ||
      fixID === "schema.range.clampMax" ||
      fixID === "schema.enum.useValue" ||
      fixID === "schema.default.use") &&
    issue?.row &&
    issue?.col &&
    typeof fix?.value === "string"
  ) {
    return [{ row: Number(issue.row) - 1, col: Number(issue.col) - 1, value: fix.value }];
  }
  return [];
}

function uniqueHeaderNameForSourceCol(sourceCol) {
  const current = String(columnNameForSource(sourceCol)).trim() || `Column ${sourceCol + 1}`;
  const taken = new Set(
    allColumns
      .filter((column) => column.sourceCol !== sourceCol)
      .map((column) => String(column.name || "").trim().toLowerCase())
      .filter(Boolean),
  );
  if (!taken.has(current.toLowerCase())) return current;
  for (let index = 2; index < 10000; index++) {
    const candidate = `${current} ${index}`;
    if (!taken.has(candidate.toLowerCase())) return candidate;
  }
  return `${current} ${sourceCol + 1}`;
}

function renameSourceColumn(sourceCol, name) {
  for (const column of allColumns) {
    if (column.sourceCol === sourceCol) column.name = name;
  }
  for (const column of columns) {
    if (column.sourceCol === sourceCol) column.name = name;
  }
}

async function applyProfileQuickFix(issue, fix) {
  if (!documentId || readOnly) {
    status.textContent = readOnly ? t("profile.quickFix.readOnly") : t("profile.quickFix.noFile");
    return;
  }
  const fixID = typeof fix === "string" ? fix : fix?.id;
  if (fixID === "schema.column.add") {
    await applyProfileMissingColumnQuickFix(issue, fix);
    return;
  }
  const edits = profileQuickFixEdits(issue, fix);
  if (!edits.length) {
    status.textContent = t("profile.quickFix.unavailable");
    return;
  }
  try {
    const response = await fetch("/api/edit-batch", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ documentId, edits }),
    });
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "quick fix failed");
    dirty = true;
    profileData = null;
    rows = new Map();
    rowSource = new Map();
    windowStart = 0;
    windowEnd = 0;
    windowGeneration++;
    const row = issueDisplayRow(issue);
    await maybeLoadWindow(row);
    applyEditsToLoadedRows(edits);
    if ((fixID === "header.fill" || fixID === "header.unique") && issue?.col) {
      const sourceCol = Number(issue.col) - 1;
      renameSourceColumn(sourceCol, edits[0]?.value ?? columnNameForSource(sourceCol));
    }
    draw();
    document.body.dataset.profileQuickFix = `${fixID}:${edits.length}`;
    status.textContent = tf("profile.quickFix.applied", { count: formatCount(edits.length) });
    await loadProfile(true);
  } catch (error) {
    status.textContent = tf("profile.quickFix.failed", { message: error.message });
  }
}

function profileMissingColumnEdits(issue, fix, sourceCol) {
  const name = String(issue?.rule || "").trim();
  if (!name || !Number.isInteger(sourceCol) || sourceCol < 0) return [];
  const edits = dataStartRow > 0 ? [{ row: 0, col: sourceCol, value: name }] : [];
  const defaultValue = typeof fix?.value === "string" ? fix.value : "";
  if (defaultValue !== "") {
    for (let sourceRow = dataStartRow; sourceRow < dataStartRow + sourceRowCount; sourceRow++) {
      edits.push({ row: sourceRow, col: sourceCol, value: defaultValue });
    }
  }
  return edits;
}

async function applyProfileMissingColumnQuickFix(issue, fix) {
  const name = String(issue?.rule || "").trim();
  if (!name || dataStartRow <= 0) {
    status.textContent = t("profile.quickFix.unavailable");
    return;
  }
  const sourceCol = allColumns.length;
  try {
    const insertResponse = await fetch("/api/column/insert", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ documentId, col: sourceCol }),
    });
    const insertPayload = await insertResponse.json();
    if (!insertResponse.ok) throw new Error(insertPayload.error ?? "add column failed");
    applyMutationPayload(insertPayload, "");
    const edits = profileMissingColumnEdits(issue, fix, sourceCol);
    if (edits.length) {
      const editResponse = await fetch("/api/edit-batch", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ documentId, edits }),
      });
      const editPayload = await editResponse.json();
      if (!editResponse.ok) throw new Error(editPayload.error ?? "fill column failed");
      dirty = true;
      profileData = null;
      rows = new Map();
      rowSource = new Map();
      windowStart = 0;
      windowEnd = 0;
      windowGeneration++;
      selected = { row: 0, col: Math.max(0, visibleColForSource(sourceCol)) };
      selectionAnchor = { ...selected };
      await maybeLoadWindow(0);
      applyEditsToLoadedRows(edits);
      renameSourceColumn(sourceCol, name);
    }
    draw();
    document.body.dataset.profileQuickFix = `schema.column.add:${edits.length}:col${sourceCol + 1}`;
    status.textContent = tf("profile.quickFix.applied", { count: formatCount(edits.length) });
    await loadProfile(true);
  } catch (error) {
    status.textContent = tf("profile.quickFix.failed", { message: error.message });
  }
}

function openDiffPanel() {
  if (metadata?.path && !leftDiffPath.value) {
    leftDiffPath.value = metadata.path;
  }
  if (!rightDiffPath.value) {
    rightDiffPath.value = commandInput.value.trim();
  }
  setViewMode("diff");
}

async function runDiffCompare() {
  const leftPath = leftDiffPath.value.trim();
  const rightPath = rightDiffPath.value.trim();
  const rightText = diffRightTextMode.checked ? diffRightText.value : "";
  if (!leftPath || (!rightPath && !rightText)) {
    status.textContent = "Choose left CSV and right CSV path or text to compare.";
    return;
  }
  status.textContent = "Comparing files...";
  const response = await fetch("/api/diff", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      leftPath,
      rightPath: diffRightTextMode.checked ? "" : rightPath,
      rightText,
      rightLabel: diffRightTextMode.checked ? diffRightTextLabel() : "",
      maxRows: 1000,
      ...currentDiffOptions(),
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Diff failed: ${payload.error ?? "unknown error"}`;
    diffResult.replaceChildren(emptyDiffNode(payload.error ?? "Diff failed"));
    diffRows = [];
    activeDiffIndex = -1;
    lastDiffResult = null;
    updateDiffNavigationState();
    return;
  }
  lastDiffResult = payload;
  renderDiffResult(payload);
  const filterLabel = diffFilterActive()
    ? tf("diff.status.filteredSuffix", { shown: document.body.dataset.diffFilteredRows || "0" })
    : "";
  status.textContent =
    `${tf("diff.status.summary", {
      changed: payload.summary.changed,
      added: payload.summary.added,
      deleted: payload.summary.deleted,
    })}${filterLabel}`;
}

function currentDiffOptions() {
  return {
    ignoreCase: diffIgnoreCase.checked,
    ignoreWhitespace: diffIgnoreWhitespace.checked,
    ignoredColumns: parseColumnList(diffIgnoredColumns.value),
    ignoredColumnNames: parseColumnNames(diffIgnoredColumns.value),
    matchColumnsByName: diffMatchHeaders.checked,
    keyColumn: diffKeyColumn.value.trim(),
  };
}

function diffRightTextLabel() {
  if (diffRightText.dataset.source === "clipboard") return "clipboard";
  if (diffRightText.dataset.source === "local-document") return "local document";
  return "text";
}

function setDiffRightTextMode(enabled) {
  diffRightTextMode.checked = enabled;
  diffRightText.hidden = !enabled;
  if (!enabled) {
    diffRightText.dataset.source = "";
    document.body.dataset.diffTextSource = "";
  } else {
    document.body.dataset.diffTextSource = diffRightText.dataset.source || "text";
  }
}

async function pasteDiffClipboardText() {
  if (!navigator.clipboard?.readText) {
    status.textContent = "Clipboard read is unavailable; paste CSV text into the right side.";
    setDiffRightTextMode(true);
    diffRightText.focus();
    document.body.dataset.diffClipboard = "unavailable";
    return;
  }
  try {
    const text = await navigator.clipboard.readText();
    if (!text) {
      status.textContent = "Clipboard is empty.";
      document.body.dataset.diffClipboard = "empty";
      return;
    }
    setDiffRightTextMode(true);
    diffRightText.value = text;
    diffRightText.dataset.source = "clipboard";
    document.body.dataset.diffTextSource = "clipboard";
    document.body.dataset.diffClipboard = String(text.length);
    rightDiffPath.value = "";
    diffRightText.focus();
    status.textContent = `Loaded ${text.length.toLocaleString()} clipboard characters for compare.`;
  } catch (error) {
    status.textContent = `Clipboard read failed: ${error.message}`;
    setDiffRightTextMode(true);
    diffRightText.focus();
    document.body.dataset.diffClipboard = "failed";
  }
}

function parseColumnList(value) {
  return parseColumnTokens(value)
    .filter((item) => /^\d+$/.test(item))
    .map((item) => Number.parseInt(item, 10))
    .filter((item) => item > 0);
}

function parseColumnNames(value) {
  return parseColumnTokens(value).filter((item) => !/^\d+$/.test(item));
}

function parseColumnTokens(value) {
  return String(value ?? "")
    .split(",")
    .flatMap((item) => {
      const token = item.trim();
      if (!token) return [];
      if (/^\d+(?:\s+\d+)*$/.test(token)) return token.split(/\s+/);
      return [token];
    });
}

function renderDiffResult(result) {
  const summary = result.summary ?? {};
  const allRows = Array.isArray(result.rows) ? result.rows : [];
  const filteredRows = filteredDiffRows(allRows);
  diffRows = [];
  activeDiffIndex = -1;
  diffSummary.replaceChildren(
    summaryNode(t("diff.summary.unchanged"), summary.unchanged ?? 0),
    summaryNode(t("diff.summary.changed"), summary.changed ?? 0),
    summaryNode(t("diff.summary.added"), summary.added ?? 0),
    summaryNode(t("diff.summary.deleted"), summary.deleted ?? 0),
    summaryNode(t("diff.summary.shown"), filteredRows.length),
  );
  document.body.dataset.diffFilterStatus = diffStatusFilter?.value || "all";
  document.body.dataset.diffFilterQuery = diffSearchInput?.value.trim() || "";
  document.body.dataset.diffSort = diffSortSelect?.value || "source";
  document.body.dataset.diffSortDirection = diffSortDirection();
  document.body.dataset.diffFilteredRows = String(filteredRows.length);
  if (!allRows.length) {
    diffResult.replaceChildren(emptyDiffNode(t("diff.empty.none")));
    updateDiffNavigationState();
    return;
  }
  if (!filteredRows.length) {
    diffResult.replaceChildren(emptyDiffNode(t("diff.empty.filtered")));
    updateDiffNavigationState();
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const tr = document.createElement("tr");
  for (const key of [
    "diff.table.status",
    "diff.table.leftRow",
    "diff.table.left",
    "diff.table.rightRow",
    "diff.table.right",
    "diff.table.changedCells",
    "diff.table.actions",
  ]) {
    const th = document.createElement("th");
    th.textContent = t(key);
    tr.append(th);
  }
  thead.append(tr);
  table.append(thead);
  const tbody = document.createElement("tbody");
  filteredRows.forEach(({ row, index }) => {
    const tr = document.createElement("tr");
    tr.className = row.status;
    if (row.status !== "unchanged") {
      tr.dataset.diffIndex = String(diffRows.length);
      tr.dataset.resultIndex = String(index);
      diffRows.push(tr);
    }
    appendCell(tr, localizedDiffStatus(row.status));
    appendCell(tr, row.leftRow ?? "");
    appendCell(tr, (row.leftRecord ?? []).join(" | "));
    appendCell(tr, row.rightRow ?? "");
    appendCell(tr, (row.rightRecord ?? []).join(" | "));
    const cells = document.createElement("td");
    cells.className = "diff-cells";
    appendChangedCells(cells, row);
    tr.append(cells);
    const actions = document.createElement("td");
    actions.className = "diff-actions";
    appendDiffActions(actions, row, index);
    tr.append(actions);
    tbody.append(tr);
  });
  table.append(tbody);
  diffResult.replaceChildren(table);
  updateDiffNavigationState();
  if (diffRows.length > 0) activateDiffRow(0, false);
}

function filteredDiffRows(rows) {
  const filter = diffStatusFilter?.value || "all";
  const query = (diffSearchInput?.value || "").trim().toLowerCase();
  const filtered = rows
    .map((row, index) => ({ row, index }))
    .filter(({ row }) => diffRowMatchesStatus(row, filter))
    .filter(({ row }) => !query || diffRowSearchText(row).toLowerCase().includes(query));
  return sortedDiffRows(filtered);
}

function diffRowMatchesStatus(row, filter) {
  const statusValue = String(row?.status || "");
  switch (filter) {
    case "differences":
      return statusValue !== "unchanged";
    case "changed":
    case "added":
    case "deleted":
    case "unchanged":
      return statusValue === filter;
    default:
      return true;
  }
}

function diffRowSearchText(row) {
  const changed = (row.changedCells ?? [])
    .map((cell) => [`c${cell.col}`, cell.left, cell.right].join(" "))
    .join(" ");
  return [
    row.status,
    localizedDiffStatus(row.status),
    row.leftRow,
    row.rightRow,
    ...(row.leftRecord ?? []),
    ...(row.rightRecord ?? []),
    changed,
  ].join(" ");
}

function sortedDiffRows(items) {
  const sort = diffSortSelect?.value || "source";
  const direction = diffSortDirection() === "desc" ? -1 : 1;
  const next = [...items];
  next.sort((left, right) => {
    const missing = compareDiffSortMissing(left, right, sort);
    if (missing !== 0) return missing;
    const compared = compareDiffRows(left, right, sort);
    if (compared !== 0) return compared * direction;
    return left.index - right.index;
  });
  return next;
}

function compareDiffRows(left, right, sort) {
  if (sort === "status") {
    return diffStatusRank(left.row?.status) - diffStatusRank(right.row?.status);
  }
  if (sort === "leftRow") {
    return compareNullableNumbers(left.row?.leftRow, right.row?.leftRow);
  }
  if (sort === "rightRow") {
    return compareNullableNumbers(left.row?.rightRow, right.row?.rightRow);
  }
  if (sort === "leftText") {
    return compareDiffText(diffRecordText(left.row?.leftRecord), diffRecordText(right.row?.leftRecord));
  }
  if (sort === "rightText") {
    return compareDiffText(diffRecordText(left.row?.rightRecord), diffRecordText(right.row?.rightRecord));
  }
  return left.index - right.index;
}

function compareDiffSortMissing(left, right, sort) {
  let leftMissing = false;
  let rightMissing = false;
  if (sort === "leftRow") {
    leftMissing = !Number.isFinite(Number(left.row?.leftRow));
    rightMissing = !Number.isFinite(Number(right.row?.leftRow));
  } else if (sort === "rightRow") {
    leftMissing = !Number.isFinite(Number(left.row?.rightRow));
    rightMissing = !Number.isFinite(Number(right.row?.rightRow));
  } else if (sort === "leftText") {
    leftMissing = !diffRecordText(left.row?.leftRecord);
    rightMissing = !diffRecordText(right.row?.leftRecord);
  } else if (sort === "rightText") {
    leftMissing = !diffRecordText(left.row?.rightRecord);
    rightMissing = !diffRecordText(right.row?.rightRecord);
  }
  if (leftMissing && rightMissing) return 0;
  if (leftMissing) return 1;
  if (rightMissing) return -1;
  return 0;
}

function diffStatusRank(statusValue) {
  switch (statusValue) {
    case "changed":
      return 0;
    case "added":
      return 1;
    case "deleted":
      return 2;
    case "unchanged":
      return 3;
    default:
      return 4;
  }
}

function localizedDiffStatus(statusValue) {
  return translationOrEmpty(`diff.status.${statusValue}`) || String(statusValue || "");
}

function compareNullableNumbers(left, right) {
  return Number(left) - Number(right);
}

function diffRecordText(record) {
  return (record ?? []).join(" ");
}

function compareDiffText(left, right) {
  return left.localeCompare(right, undefined, {
    numeric: true,
    sensitivity: "base",
  });
}

function diffSortDirection() {
  return diffSortDirectionButton?.dataset.direction === "desc" ? "desc" : "asc";
}

function setDiffSortDirection(direction) {
  const normalized = direction === "desc" ? "desc" : "asc";
  if (!diffSortDirectionButton) return;
  diffSortDirectionButton.dataset.direction = normalized;
  diffSortDirectionButton.textContent = normalized === "desc" ? t("diff.sort.desc") : t("diff.sort.asc");
  diffSortDirectionButton.setAttribute("aria-pressed", String(normalized === "desc"));
  diffSortDirectionButton.title = normalized === "desc" ? t("diff.sortDescending") : t("diff.sortAscending");
}

function toggleDiffSortDirection() {
  setDiffSortDirection(diffSortDirection() === "desc" ? "asc" : "desc");
  applyDiffResultFilters();
}

function diffFilterActive() {
  return (diffStatusFilter?.value || "all") !== "all" || Boolean(diffSearchInput?.value.trim());
}

function applyDiffResultFilters() {
  if (!lastDiffResult) {
    document.body.dataset.diffFilterStatus = diffStatusFilter?.value || "all";
    document.body.dataset.diffFilterQuery = diffSearchInput?.value.trim() || "";
    return;
  }
  renderDiffResult(lastDiffResult);
  const shown = Number(document.body.dataset.diffFilteredRows || 0);
  const total = Array.isArray(lastDiffResult.rows) ? lastDiffResult.rows.length : 0;
  status.textContent = diffFilterActive()
    ? tf("diff.status.filterShown", { shown: formatCount(shown), total: formatCount(total) })
    : tf("diff.status.filterCleared", { total: formatCount(total) });
}

function clearDiffResultFilters() {
  if (diffStatusFilter) diffStatusFilter.value = "all";
  if (diffSearchInput) diffSearchInput.value = "";
  applyDiffResultFilters();
}

function updateDiffNavigationState() {
  const disabled = diffRows.length === 0;
  diffPrevButton.disabled = disabled;
  diffNextButton.disabled = disabled;
  document.body.dataset.diffCount = String(diffRows.length);
  document.body.dataset.activeDiff = activeDiffIndex >= 0 ? String(activeDiffIndex + 1) : "";
}

function activateDiffRow(index, scroll = true) {
  if (diffRows.length === 0) return;
  const nextIndex = ((index % diffRows.length) + diffRows.length) % diffRows.length;
  if (activeDiffIndex >= 0 && diffRows[activeDiffIndex]) {
    diffRows[activeDiffIndex].classList.remove("active-diff");
  }
  activeDiffIndex = nextIndex;
  const row = diffRows[activeDiffIndex];
  row.classList.add("active-diff");
  if (scroll) {
    row.scrollIntoView({ block: "nearest" });
  }
  updateDiffNavigationState();
  status.textContent = tf("diff.status.active", { current: activeDiffIndex + 1, total: diffRows.length });
}

function applyDiffURLGoto(side, location) {
  const targetSide = side === "right" ? "right" : "left";
  const targetRow = Number.parseInt(String(location?.startRow || ""), 10);
  document.body.dataset.diffGotoSide = targetSide;
  document.body.dataset.diffGotoRow = Number.isFinite(targetRow) && targetRow > 0 ? String(targetRow) : "";
  if (!lastDiffResult || !Array.isArray(lastDiffResult.rows) || !Number.isFinite(targetRow) || targetRow <= 0) {
    document.body.dataset.diffGotoActive = "invalid";
    return false;
  }
  const rowKey = targetSide === "right" ? "rightRow" : "leftRow";
  const resultIndex = lastDiffResult.rows.findIndex((row) => Number(row?.[rowKey]) === targetRow && row?.status !== "unchanged");
  if (resultIndex < 0) {
    document.body.dataset.diffGotoActive = "not-found";
    return false;
  }
  const diffIndex = diffRows.findIndex((row) => Number(row.dataset.resultIndex) === resultIndex);
  if (diffIndex < 0) {
    document.body.dataset.diffGotoActive = "filtered";
    return false;
  }
  activateDiffRow(diffIndex);
  document.body.dataset.diffGotoActive = `${targetSide}:${targetRow}:${resultIndex}`;
  status.textContent = `Diff goto ${targetSide} row ${targetRow}.`;
  return true;
}

function navigateDiff(direction) {
  if (viewMode !== "diff" || diffRows.length === 0) return false;
  activateDiffRow(activeDiffIndex < 0 ? 0 : activeDiffIndex + direction);
  return true;
}

function appendDiffActions(container, row, resultIndex) {
  if (row.status === "unchanged") return;
  appendDiffOpenSourceActions(container, row);
  const rowButton = document.createElement("button");
  rowButton.type = "button";
  rowButton.textContent =
    row.status === "added"
      ? t("diff.action.insertRow")
      : row.status === "deleted"
        ? t("diff.action.deleteLeft")
        : t("diff.action.applyRow");
  rowButton.title = t("diff.action.applyRightRowToLeft");
  rowButton.addEventListener("click", () => mergeDiffChange(resultIndex, "row", 0, "right-to-left"));
  container.append(rowButton);

  if (!diffRightTextMode.checked) {
    const leftButton = document.createElement("button");
    leftButton.type = "button";
    leftButton.textContent =
      row.status === "added"
        ? t("diff.action.deleteRight")
        : row.status === "deleted"
          ? t("diff.action.insertRight")
          : t("diff.action.applyLeftRow");
    leftButton.title = t("diff.action.applyLeftRowToRight");
    leftButton.addEventListener("click", () => mergeDiffChange(resultIndex, "row", 0, "left-to-right"));
    container.append(leftButton);
  }

  if (row.status !== "changed") return;
  for (const cell of row.changedCells ?? []) {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = `c${cell.col}`;
    button.title = tf("diff.action.applyRightCellToLeft", { col: cell.col });
    button.addEventListener("click", () => mergeDiffChange(resultIndex, "cell", cell.col, "right-to-left"));
    container.append(button);
    if (!diffRightTextMode.checked) {
      const leftButton = document.createElement("button");
      leftButton.type = "button";
      leftButton.textContent = `L c${cell.col}`;
      leftButton.title = tf("diff.action.applyLeftCellToRight", { col: cell.col });
      leftButton.addEventListener("click", () => mergeDiffChange(resultIndex, "cell", cell.col, "left-to-right"));
      container.append(leftButton);
    }
  }
}

function appendDiffOpenSourceActions(container, row) {
  const leftRow = Number(row.leftRow);
  if (Number.isFinite(leftRow) && leftRow >= 1) {
    const leftButton = document.createElement("button");
    leftButton.type = "button";
    leftButton.textContent = t("diff.action.openLeftShort");
    leftButton.title = tf("diff.action.openLeftRow", { row: leftRow });
    leftButton.dataset.diffAction = "open-left-row";
    leftButton.addEventListener("click", () => openDiffSourceRow("left", leftRow));
    container.append(leftButton);
  }
  const rightRow = Number(row.rightRow);
  if (!diffRightTextMode.checked && Number.isFinite(rightRow) && rightRow >= 1) {
    const rightButton = document.createElement("button");
    rightButton.type = "button";
    rightButton.textContent = t("diff.action.openRightShort");
    rightButton.title = tf("diff.action.openRightRow", { row: rightRow });
    rightButton.dataset.diffAction = "open-right-row";
    rightButton.addEventListener("click", () => openDiffSourceRow("right", rightRow));
    container.append(rightButton);
  }
}

async function openDiffSourceRow(side, row) {
  const path = side === "right" ? rightDiffPath.value.trim() : leftDiffPath.value.trim();
  if (!path || !Number.isFinite(row) || row < 1) return;
  const label = side === "right" ? "right" : "left";
  try {
    await openPath(path, {
      applyRules: false,
      location: { startRow: row, startCol: 1, endRow: row, endCol: 1 },
    });
    setViewMode("grid");
    document.body.dataset.diffOpenSource = `${label}:${row}`;
    status.textContent = tf("diff.status.openedSource", { side: label, row: formatCount(row) });
  } catch (error) {
    status.textContent = tf("diff.status.openSourceFailed", { message: error.message });
  }
}

async function mergeDiffChange(resultIndex, mode, col, direction = "right-to-left") {
  const leftPath = leftDiffPath.value.trim();
  const rightPath = rightDiffPath.value.trim();
  const rightText = diffRightTextMode.checked ? diffRightText.value : "";
  if (!leftPath || (!rightPath && !rightText)) return;
  status.textContent = "Applying diff change...";
  const response = await fetch("/api/diff-merge", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      leftPath,
      rightPath: diffRightTextMode.checked ? "" : rightPath,
      rightText,
      rightLabel: diffRightTextMode.checked ? diffRightTextLabel() : "",
      resultIndex,
      mode,
      col,
      direction,
      ...currentDiffOptions(),
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Merge failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const backup = payload.save?.backupPath ? ` backup ${payload.save.backupPath}` : "";
  status.textContent = `Merged ${payload.applied}${backup}`;
  await runDiffCompare();
}

function diffSidePath(side) {
  return side === "right" ? rightDiffPath.value.trim() : leftDiffPath.value.trim();
}

function diffSideRow(row, side) {
  return Number(side === "right" ? row.rightRow : row.leftRow);
}

function localizedDiffSide(side) {
  return side === "right" ? t("diff.side.right") : t("diff.side.left");
}

async function editDiffCell(side, row, col, value, previousValue) {
  const path = diffSidePath(side);
  if (!path || !Number.isFinite(row) || row < 1 || !Number.isFinite(col) || col < 1) return;
  if (value === previousValue) {
    status.textContent = t("diff.status.editNoChange");
    return;
  }
  const sideLabel = localizedDiffSide(side);
  status.textContent = tf("diff.status.editingCell", { side: sideLabel, row: formatCount(row), col: formatCount(col) });
  try {
    const response = await fetch("/api/diff-edit-cell", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ path, row, col, value }),
    });
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "diff edit failed");
    document.body.dataset.diffEditCell = `${side}:${row}:${col}`;
    document.body.dataset.diffEditValue = value;
    document.body.dataset.diffEditBackup = payload.save?.backupPath || "";
    delete document.body.dataset.diffEditError;
    await runDiffCompare();
    const backup = payload.save?.backupPath ? tf("diff.status.backupSuffix", { path: payload.save.backupPath }) : "";
    status.textContent = tf("diff.status.editedCell", {
      side: sideLabel,
      row: formatCount(row),
      col: formatCount(col),
      backup,
    });
  } catch (error) {
    document.body.dataset.diffEditCell = "failed";
    document.body.dataset.diffEditError = error.message;
    status.textContent = tf("diff.status.editFailed", { message: error.message });
  }
}

function appendChangedCells(container, row) {
  const changedCells = row.changedCells ?? [];
  if (!changedCells.length) return;
  changedCells.forEach((cell, index) => {
    if (index > 0) container.append(document.createTextNode(" "));
    const item = document.createElement("span");
    item.className = "diff-cell-change";
    item.append(document.createTextNode(`c${cell.col}: `));
    appendDiffSpans(item, cell.leftSpans, cell.left);
    item.append(document.createTextNode(" -> "));
    appendDiffSpans(item, cell.rightSpans, cell.right);
    const editors = document.createElement("span");
    editors.className = "diff-cell-editors";
    appendDiffCellEditor(editors, row, cell, "left");
    appendDiffCellEditor(editors, row, cell, "right");
    if (editors.childElementCount > 0) item.append(editors);
    container.append(item);
  });
}

function appendDiffCellEditor(container, row, cell, side) {
  if (side === "right" && diffRightTextMode.checked) return;
  const sourceRow = diffSideRow(row, side);
  if (!Number.isFinite(sourceRow) || sourceRow < 1 || !diffSidePath(side)) return;
  const col = Number(cell.col);
  if (!Number.isFinite(col) || col < 1) return;
  const value = side === "right" ? (cell.right ?? "") : (cell.left ?? "");
  const wrapper = document.createElement("span");
  wrapper.className = "diff-cell-editor";
  const label = document.createElement("span");
  label.textContent = side === "right" ? t("diff.edit.rightShort") : t("diff.edit.leftShort");
  const input = document.createElement("input");
  input.type = "text";
  input.value = value;
  input.dataset.diffEditSide = side;
  input.dataset.diffEditRow = String(sourceRow);
  input.dataset.diffEditCol = String(col);
  input.setAttribute(
    "aria-label",
    side === "right"
      ? tf("diff.action.editRightCell", { row: sourceRow, col })
      : tf("diff.action.editLeftCell", { row: sourceRow, col }),
  );
  const button = document.createElement("button");
  button.type = "button";
  button.textContent = t("diff.action.editCell");
  button.title =
    side === "right"
      ? tf("diff.action.editRightCell", { row: sourceRow, col })
      : tf("diff.action.editLeftCell", { row: sourceRow, col });
  const submit = () => editDiffCell(side, sourceRow, col, input.value, value);
  button.addEventListener("click", submit);
  input.addEventListener("keydown", (event) => {
    if (event.key !== "Enter") return;
    event.preventDefault();
    submit();
  });
  wrapper.append(label, input, button);
  container.append(wrapper);
}

function appendDiffSpans(container, spans, fallback) {
  const source = spans?.length ? spans : [{ text: fallback ?? "", changed: true }];
  for (const span of source) {
    const node = document.createElement("span");
    node.textContent = span.text;
    if (span.changed) node.className = "diff-inline-change";
    container.append(node);
  }
}

function appendCell(row, value) {
  const cell = document.createElement("td");
  cell.textContent = String(value);
  row.append(cell);
}

function summaryNode(label, value) {
  const node = document.createElement("span");
  const strong = document.createElement("strong");
  strong.textContent = value.toLocaleString();
  node.append(label, " ", strong);
  return node;
}

function emptyDiffNode(message) {
  const node = document.createElement("div");
  node.className = "diff-empty";
  node.textContent = message;
  return node;
}

function emptyProfileNode(message) {
  const node = document.createElement("div");
  node.className = "profile-empty";
  node.textContent = message;
  return node;
}

function formatCount(value) {
  return Number(value || 0).toLocaleString();
}

function formatPercent(value) {
  return `${(Number(value || 0) * 100).toFixed(1)}%`;
}

function formatNumber(value) {
  const number = Number(value);
  if (!Number.isFinite(number)) return "";
  return number.toLocaleString(undefined, { maximumFractionDigits: 3 });
}

function formatElapsedSeconds(milliseconds) {
  const seconds = Math.max(0, Number(milliseconds || 0)) / 1000;
  if (seconds < 10) return seconds.toFixed(1);
  return Math.round(seconds).toLocaleString();
}

function formatBackupSize(value) {
  const bytes = Number(value || 0);
  if (!Number.isFinite(bytes) || bytes <= 0) return "0 B";
  if (bytes < 1024) return `${formatCount(bytes)} B`;
  const units = ["KB", "MB", "GB", "TB"];
  let scaled = bytes;
  let unit = units[0];
  for (const nextUnit of units) {
    scaled /= 1024;
    unit = nextUnit;
    if (scaled < 1024) break;
  }
  return `${scaled.toLocaleString(undefined, { maximumFractionDigits: scaled >= 10 ? 1 : 2 })} ${unit}`;
}

function formatBackupTime(value) {
  const milliseconds = Number(value || 0) / 1_000_000;
  if (!Number.isFinite(milliseconds) || milliseconds <= 0) return "unknown";
  return new Date(milliseconds).toLocaleString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function renderSQLResult(result) {
  const resultColumns = sqlResultColumns(result);
  const resultRows = result.rows ?? [];
  if (!resultColumns.length && !resultRows.length) {
    sqlResult.replaceChildren(emptySQLNode(t("sql.empty.rows")));
    clearSQLResultSelection();
    syncSQLResultActions();
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const headerRow = document.createElement("tr");
  for (const column of resultColumns) {
    const th = document.createElement("th");
    th.textContent = column;
    headerRow.append(th);
  }
  thead.append(headerRow);
  table.append(thead);

  const tbody = document.createElement("tbody");
  resultRows.forEach((row, rowIndex) => {
    const tr = document.createElement("tr");
    for (let colIndex = 0; colIndex < resultColumns.length; colIndex++) {
      const td = document.createElement("td");
      td.textContent = row[colIndex] ?? "";
      td.tabIndex = 0;
      td.dataset.sqlRow = String(rowIndex);
      td.dataset.sqlCol = String(colIndex);
      td.setAttribute("aria-selected", "false");
      tr.append(td);
    }
    tbody.append(tr);
  });
  if (resultRows.length === 0) {
    const tr = document.createElement("tr");
    const td = document.createElement("td");
    td.colSpan = Math.max(1, resultColumns.length);
    td.className = "sql-empty";
    td.textContent = t("sql.empty.rows");
    tr.append(td);
    tbody.append(tr);
  }
  table.append(tbody);
  sqlResult.replaceChildren(table);
  syncSQLResultSelection();
  syncSQLResultActions();
}

function emptySQLNode(message) {
  const node = document.createElement("div");
  node.className = "sql-empty";
  node.textContent = message;
  return node;
}

canvas.addEventListener(
  "wheel",
  (event) => {
    event.preventDefault();
    closeColumnMenu();
    closeGridContextMenu();
    const maxX = maxHorizontalScroll();
    applyVerticalWheelDelta(event.deltaY);
    scrollX = Math.max(0, Math.min(maxX, scrollX + event.deltaX));
    if (editing) commitEditor();
    draw();
  },
  { passive: false },
);

canvas.addEventListener("contextmenu", openGridContextMenu);

canvas.addEventListener("mousedown", (event) => {
  canvas.focus();
  if (event.button !== 0) return;
  closeGridContextMenu();
  if (!documentId || columns.length === 0 || totalRows === 0) return;
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  if (fillHandleAt(x, y)) {
    event.preventDefault();
    if (editing) commitEditor();
    closeColumnMenu();
    startFillDrag();
    return;
  }
  const resizeCol = headerResizeColumnAt(x, y);
  if (resizeCol >= 0) {
    event.preventDefault();
    if (editing) commitEditor();
    startColumnResize(resizeCol, event.clientX);
    return;
  }
  if (event.metaKey || event.ctrlKey) {
    const link = cellLinkAtPoint(x, y);
    if (link.href) {
      event.preventDefault();
      openCellLinkWithPrompt(link);
      return;
    }
  }
  const smartGesture = smartSelectionGestureForEvent(event);
  if (y < headerHeight) {
    const col = columnAt(x);
    if (x < rowNumberWidth) {
      event.preventDefault();
      if (editing) commitEditor();
      closeColumnMenu();
      if (smartGesture) startSmartSelectionGesture(smartGesture);
      else clearSmartSelectionRanges();
      setHeaderInSelection(event.altKey);
      selectAllBodyCells();
      draggingSelection = false;
      dragSelectionMode = "cell";
      if (smartGesture) {
        finishSmartSelectionGesture();
        return;
      }
      draw();
      status.textContent = event.altKey
        ? "Selected all body cells with headers"
        : "Selected all body cells";
      return;
    }
    if (x >= rowNumberWidth && col >= 0) {
      event.preventDefault();
      if (editing) commitEditor();
      if (isHeaderMenuTarget(x, col)) {
        openColumnMenu(col, event.clientX, rect.top + headerHeight + 4);
      } else {
        closeColumnMenu();
        if (!readOnly && !event.shiftKey && !event.altKey && !smartGesture && isWholeColumnSelected(col)) {
          startColumnReorder(col, event);
          return;
        }
        if (smartGesture) startSmartSelectionGesture(smartGesture);
        setHeaderInSelection(event.altKey);
        selectColumnRange(col, smartGesture ? false : event.shiftKey);
        draggingSelection = true;
        dragSelectionMode = "column";
        draw();
      }
    }
    return;
  }
  closeColumnMenu();
  if (x < rowNumberWidth) {
    event.preventDefault();
    if (editing) commitEditor();
    const row = rowAt(y);
    if (!readOnly && !event.shiftKey && !event.altKey && !smartGesture && isWholeRowSelected(row)) {
      startRowReorder(row, event);
      return;
    }
    if (smartGesture) startSmartSelectionGesture(smartGesture);
    setHeaderInSelection(event.altKey);
    selectRowRange(row, smartGesture ? false : event.shiftKey);
    draggingSelection = true;
    dragSelectionMode = "row";
    draw();
    return;
  }
  const row = Math.floor((y - headerHeight + scrollY) / rowHeight);
  const col = columnAt(x);
  if (row >= 0 && row < totalRows && col >= 0) {
    event.preventDefault();
    if (smartGesture) startSmartSelectionGesture(smartGesture);
    selected = { row, col };
    if (!event.shiftKey || smartGesture) {
      selectionAnchor = { ...selected };
    }
    draggingSelection = true;
    dragSelectionMode = "cell";
    if (editing) commitEditor();
    draw();
  }
});

window.addEventListener("mousemove", (event) => {
  if (draggingFillHandle) {
    event.preventDefault();
    updateFillDrag(event);
    return;
  }
  if (resizingColumn >= 0) {
    event.preventDefault();
    updateColumnResize(event.clientX);
    return;
  }
  if (pendingColumnReorder) {
    event.preventDefault();
    updateColumnReorder(event);
    return;
  }
  if (pendingRowReorder) {
    event.preventDefault();
    updateRowReorder(event);
    return;
  }
  if (draggingSelection) {
    event.preventDefault();
    updateDragSelection(event);
    return;
  }
  updateCanvasHover(event);
});

window.addEventListener("mouseup", () => {
  if (draggingFillHandle) {
    finishFillDrag();
    return;
  }
  if (resizingColumn >= 0) {
    finishColumnResize();
    return;
  }
  if (pendingColumnReorder) {
    void finishColumnReorder();
    return;
  }
  if (pendingRowReorder) {
    void finishRowReorder();
    return;
  }
  if (!draggingSelection) return;
  draggingSelection = false;
  dragSelectionMode = "cell";
  finishSmartSelectionGesture();
});

canvas.addEventListener("mouseleave", () => {
  if (resizingColumn >= 0 || draggingFillHandle) return;
  clearURLHover();
  clearResizeHover();
  clearFillHandleHover();
});

canvas.addEventListener("dblclick", (event) => {
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  const resizeCol = headerResizeColumnAt(x, y);
  if (resizeCol >= 0) {
    event.preventDefault();
    autoFitVisibleColumn(resizeCol);
    return;
  }
  if (y < headerHeight || x < rowNumberWidth) return;
  openEditor();
});

window.addEventListener("copy", (event) => {
  if (!documentId || isTypingTarget(event.target) || !defaultShortcutActive("edit.copy")) return;
  event.preventDefault();
  const matrix = includeHeadersInSelection
    ? [selectedHeaders(), ...selectedMatrix()]
    : selectedMatrix();
  const text = formatMatrix(matrix, "tsv");
  event.clipboardData.setData("text/plain", text);
  window.__opencsvLastCopiedText = text;
  document.body.dataset.copyPreview = text.slice(0, 500);
  const headerLabel = includeHeadersInSelection ? " with headers" : "";
  status.textContent = `Copied ${formatCount(selectionCellCount())} cells${headerLabel}`;
});

window.addEventListener("paste", (event) => {
  if (!documentId || isTypingTarget(event.target) || !defaultShortcutActive("edit.paste")) return;
  const html = event.clipboardData.getData("text/html");
  if (html && pasteHTMLTable(html)) {
    event.preventDefault();
    return;
  }
  const text = event.clipboardData.getData("text/plain");
  if (!text) return;
  event.preventDefault();
  pasteText(text);
});

window.addEventListener("keydown", (event) => {
  if (handleShortcutEditorKeydown(event)) return;
  if (shortcutEditorOpen) return;
  if (linkConfirmDialog && !linkConfirmDialog.hidden) {
    if (event.key === "Escape") {
      event.preventDefault();
      cancelCustomLinkConfirmDialog();
    }
    return;
  }
  if (handleGlobalShortcut(event)) return;
  if (conditionBuilderOpen) return;
  if (paletteOpen) return;
  if (editing) return;
  if (isTypingTarget(event.target)) return;
  if (event.key === "Enter") {
    event.preventDefault();
    openEditor();
    return;
  }
  if (columns.length === 0) return;
  let moved = false;
  if (event.key === "ArrowDown") selected.row = Math.min(totalRows - 1, selected.row + 1);
  if (event.key === "ArrowUp") selected.row = Math.max(0, selected.row - 1);
  if (event.key === "ArrowRight") selected.col = Math.min(columns.length - 1, selected.col + 1);
  if (event.key === "ArrowLeft") selected.col = Math.max(0, selected.col - 1);
  moved = ["ArrowDown", "ArrowUp", "ArrowRight", "ArrowLeft"].includes(event.key);
  if (moved) {
    event.preventDefault();
    if (!event.shiftKey) {
      selectionAnchor = { ...selected };
    }
    draw();
  }
});

window.addEventListener("keyup", (event) => {
  if (event.key === "Meta" || event.key === "Control") {
    clearURLHover();
  }
});

function handleGlobalShortcut(event) {
  const primary = event.metaKey || event.ctrlKey;
  const key = event.key.toLowerCase();
  if (runCustomShortcut(event)) return true;
  if (event.altKey && key === "f5") {
    event.preventDefault();
    return navigateDiff(event.shiftKey ? -1 : 1);
  }
  if (event.shiftKey && event.altKey && key === "h" && defaultShortcutActive("selection.toggleHeaders")) {
    event.preventDefault();
    toggleHeaderInSelection();
    return true;
  }
  if (event.shiftKey && event.altKey && !paletteOpen && !isTypingTarget(event.target)) {
    const arrowCommandIDByKey = {
      arrowup: "row.moveUp",
      arrowdown: "row.moveDown",
      arrowleft: "column.moveLeft",
      arrowright: "column.moveRight",
    };
    const arrowCommandID = arrowCommandIDByKey[key];
    if (arrowCommandID && defaultShortcutActive(arrowCommandID)) {
      event.preventDefault();
      executeCommand(arrowCommandID);
      return true;
    }
  }
  if (event.shiftKey && event.altKey && !paletteOpen && !isTypingTarget(event.target)) {
    const commandIDByKey = {
      d: "table.dedupeRows",
      c: "table.dedupeColumns",
      r: "row.hide",
      x: "column.hide",
      w: "cell.swap",
      j: "cell.joinHorizontal",
      s: "cell.splitHorizontal",
      v: "cell.joinVertical",
      b: "cell.splitVertical",
      f: "cell.fillDown",
      g: "cell.fillRight",
      a: "selection.insertAggregateRow",
      m: "selection.addRange",
      "-": "selection.subtractRange",
      "_": "selection.subtractRange",
      minus: "selection.subtractRange",
      i: "selection.invert",
      t: "text.trim",
      u: "text.uppercase",
      l: "text.lowercase",
    };
    const commandID = commandIDByKey[key];
    if (commandID && defaultShortcutActive(commandID)) {
      event.preventDefault();
      executeCommand(commandID);
      return true;
    }
  }
  if (primary && event.shiftKey && event.key.toLowerCase() === "p" && defaultShortcutActive("command.palette")) {
    event.preventDefault();
    openCommandPalette();
    return true;
  }
  if (event.key === "Escape" && paletteOpen) {
    event.preventDefault();
    closeCommandPalette();
    return true;
  }
  if (event.key === "Escape" && !columnMenu.hidden) {
    event.preventDefault();
    closeColumnMenu();
    return true;
  }
  if (event.key === "Escape" && !gridContextMenu.hidden) {
    event.preventDefault();
    closeGridContextMenu();
    return true;
  }
  if (event.key === "Escape" && conditionBuilderOpen) {
    event.preventDefault();
    closeConditionBuilder();
    return true;
  }
  if (event.key === "Escape" && cellValueDialogOpen) {
    event.preventDefault();
    closeCellValueDialog();
    return true;
  }
  if (event.key === "Escape" && aboutSupportPanel && !aboutSupportPanel.hidden) {
    event.preventDefault();
    closeAboutSupportPanel();
    return true;
  }
  if (event.key === "Escape" && !externalChangePanel.hidden) {
    event.preventDefault();
    closeExternalChangePanel();
    return true;
  }
  if (event.key === "Escape" && dirtyClosePanel && !dirtyClosePanel.hidden) {
    event.preventDefault();
    closeDirtyClosePanel();
    return true;
  }
  if (event.key === "Escape" && transposeConfirmPanel && !transposeConfirmPanel.hidden) {
    event.preventDefault();
    closeTransposeConfirmPanel();
    return true;
  }
  if (event.key === "Escape" && splitConfirmPanel && !splitConfirmPanel.hidden) {
    event.preventDefault();
    closeSplitOverwriteConfirmPanel();
    return true;
  }
  if (event.key === "Escape" && smartSelectionRanges.length > 0) {
    event.preventDefault();
    clearSmartSelectionRanges();
    return true;
  }
  if (event.key === "Escape" && selectionCellCount() > 1) {
    event.preventDefault();
    selectionAnchor = { ...selected };
    draw();
    return true;
  }
  if (!primary || paletteOpen) return false;
  if (key === "," && defaultShortcutActive("settings.keyboardShortcuts")) {
    event.preventDefault();
    openKeyboardShortcuts();
    return true;
  }
  if (key === "enter" && event.target === cellValueInput) {
    event.preventDefault();
    applyCellValue();
    return true;
  }
  if (key === "z" && defaultShortcutActive(event.shiftKey ? "edit.redo" : "edit.undo")) {
    event.preventDefault();
    if (event.shiftKey) {
      redoCurrentDocument();
    } else {
      undoCurrentDocument();
    }
    return true;
  }
  if (key === "o" && defaultShortcutActive("file.open")) {
    event.preventDefault();
    fileInput.click();
    return true;
  }
  if (key === "s" && event.shiftKey && defaultShortcutActive("file.saveAs")) {
    event.preventDefault();
    openSaveAsDialog();
    return true;
  }
  if (key === "s" && defaultShortcutActive("file.save")) {
    event.preventDefault();
    saveCurrentDocument();
    return true;
  }
  if (key === "w" && !event.shiftKey && defaultShortcutActive("file.closeTab")) {
    event.preventDefault();
    executeCommand("file.closeTab");
    return true;
  }
  if (event.key === "F4" && event.ctrlKey && !event.metaKey && defaultShortcutActive("file.closeTab")) {
    event.preventDefault();
    executeCommand("file.closeTab");
    return true;
  }
  if (key === "f" && event.shiftKey && defaultShortcutActive("view.filter")) {
    event.preventDefault();
    commandInput.focus();
    commandInput.select();
    return true;
  }
  if (key === "f" && !event.shiftKey && defaultShortcutActive("edit.find")) {
    event.preventDefault();
    focusFindInput({ target: event.target });
    return true;
  }
  return false;
}

function isTypingTarget(target) {
  if (!(target instanceof HTMLElement)) return false;
  const tagName = target.tagName.toLowerCase();
  return tagName === "input" || tagName === "textarea" || target.isContentEditable;
}

editor.addEventListener("keydown", (event) => {
  if (event.key === "Enter") {
    event.preventDefault();
    commitEditor();
  }
  if (event.key === "Escape") {
    event.preventDefault();
    editing = false;
    editor.style.display = "none";
    canvas.focus();
  }
});

openForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  await openPathInputOrPicker();
});

pathInput.addEventListener("keydown", async (event) => {
  if (event.key !== "Enter") return;
  event.preventDefault();
  await openPathInputOrPicker();
});

async function openPathInputOrPicker() {
  const path = pathInput.value.trim();
  if (!path) {
    fileInput.click();
    return;
  }
  try {
    await openPath(path);
  } catch (error) {
    status.textContent = `Open failed: ${error.message}`;
  }
}

async function openInitialURLTarget() {
  const params = new URLSearchParams(window.location.search);
  applyLaunchParams(params);
  const diffLeft = params.get("diffLeft");
  const diffRight = params.get("diffRight");
  const diffText = params.get("diffRightText");
  if (diffLeft || diffRight || diffText) {
    leftDiffPath.value = diffLeft ?? "";
    rightDiffPath.value = diffRight ?? "";
    diffRightText.value = diffText ?? "";
    diffRightText.dataset.source = "";
    setDiffRightTextMode(diffText != null);
    diffIgnoreCase.checked = truthyParam(params.get("diffIgnoreCase"));
    diffIgnoreWhitespace.checked = truthyParam(params.get("diffIgnoreWhitespace"));
    diffMatchHeaders.checked = truthyParam(params.get("diffMatchHeaders"));
    diffKeyColumn.value = params.get("diffKeyColumn") ?? "";
    diffIgnoredColumns.value = params.get("diffIgnoredColumns") ?? "";
    openDiffPanel();
    if (truthyParam(params.get("diffRun")) && diffLeft && (diffRight || diffText)) {
      await runDiffCompare();
      const diffGoto = params.get("diffGoto");
      if (diffGoto) applyDiffURLGoto(params.get("diffGotoSide") || params.get("diffSide"), parseURLLocation(diffGoto));
    }
    return;
  }
  const openPaths = params.getAll("open").filter(Boolean);
  if (!openPaths.length) return;
  const formatOverride = formatOverrideFromParams(params);
  const globalLocation = locationFromParams(params);
  const openSources = params.getAll("openSource");
  const openNames = params.getAll("openName");
  const openGotos = params.getAll("openGoto");
  const openEntries = openPaths.map((path, index) => {
    const source = String(openSources[index] || (openPaths.length === 1 ? params.get("source") : "") || "").toLowerCase();
    const stdinSource = source === "stdin" || (openPaths.length === 1 && truthyParam(params.get("stdin")));
    const goto = openGotos[index] ? parseURLLocation(openGotos[index]) : (index === openPaths.length - 1 ? globalLocation : null);
    return {
      path,
      stdinSource,
      displayName: stdinSource ? (openNames[index] || "stdin") : (openNames[index] || ""),
      location: goto,
    };
  });
  try {
    for (const entry of openEntries) {
      pathInput.value = entry.stdinSource ? "stdin" : entry.path;
      await openPath(entry.path, {
        formatOverride,
        location: entry.location,
        displayName: entry.displayName || undefined,
        remember: !entry.stdinSource,
        sourceKind: entry.stdinSource ? "stdin" : "",
      });
    }
    if (openEntries.length > 1) {
      const active = openEntries[openEntries.length - 1];
      document.body.dataset.multiOpenCount = String(openEntries.length);
      document.body.dataset.multiOpenActive = active.stdinSource ? "stdin" : active.path;
      document.body.dataset.multiOpenFiles = openEntries.map((entry) => entry.stdinSource ? "stdin" : entry.path).join("|");
      status.textContent = `Opened ${openEntries.length} files; active ${active.stdinSource ? "stdin" : fileNameFromPath(active.path)}.`;
    } else {
      delete document.body.dataset.multiOpenCount;
      delete document.body.dataset.multiOpenActive;
      delete document.body.dataset.multiOpenFiles;
    }
    if (openEntries.some((entry) => entry.stdinSource)) {
      document.body.dataset.stdinUrlBoot = "1";
    } else {
      delete document.body.dataset.stdinUrlBoot;
    }
    const sqlQuery = params.get("sql");
    if (sqlQuery != null) {
      sqlInput.value = sqlQuery;
      setViewMode("sql");
      if (truthyParam(params.get("sqlRun"))) {
        await runSQLQuery();
      }
    } else {
      const findQuery = params.get("find");
      if (findQuery != null) {
        commandInput.value = findQuery;
        invalidateFindState();
        document.body.dataset.findUrlBoot = truthyParam(params.get("findRun")) ? "run" : "loaded";
        if (truthyParam(params.get("findRun"))) {
          await findCurrentQuery();
        }
      }
      const filterQuery = params.get("filter");
      if (filterQuery != null) {
        commandInput.value = filterQuery;
        invalidateFindState();
        document.body.dataset.filterUrlBoot = truthyParam(params.get("filterRun")) ? "run" : "loaded";
        if (truthyParam(params.get("filterRun"))) {
          await filterCurrentQuery();
        }
      }
      const profileSchemaParam = params.get("profileSchema");
      if (profileSchemaParam != null) {
        profileSchemaText = profileSchemaParam;
        profileSchemaVisible = true;
        const profileSchemaNameParam = params.get("profileSchemaName");
        if (profileSchemaNameParam != null) profileSchemaNameInput.value = profileSchemaNameParam;
        if (truthyParam(params.get("profileSchemaSave"))) {
          const saved = saveNamedProfileSchema(profileSchemaNameParam || "URL Profile Schema", profileSchemaText);
          localStorage.setItem(profileSchemaStorageKey, saved.schemaText);
          profileSchemaText = saved.schemaText;
          document.body.dataset.profileSchemaUrlSave = saved.name;
        }
        syncProfileSchemaControls();
      }
      if (profileSchemaParam != null || truthyParam(params.get("profileRun"))) {
        await openProfilePanel();
      }
    }
  } catch (error) {
    status.textContent = `Open failed: ${error.message}`;
  }
}

function truthyParam(value) {
  return value === "1" || value === "true" || value === "yes";
}

function applyLaunchParams(params) {
  const newWindow = truthyParam(params.get("newWindow") ?? params.get("new-window"));
  const reuseWindow = truthyParam(params.get("reuseWindow") ?? params.get("reuse-window"));
  const wait = truthyParam(params.get("wait"));
  const flags = [];
  if (newWindow) flags.push("newWindow");
  if (reuseWindow) flags.push("reuseWindow");
  if (wait) flags.push("wait");
  document.body.dataset.launchFlags = flags.join(",");
  document.body.dataset.launchNewWindow = String(newWindow);
  document.body.dataset.launchReuseWindow = String(reuseWindow);
  document.body.dataset.launchWait = String(wait);
}

function formatOverrideFromParams(params) {
  const firstParam = (...names) => {
    for (const name of names) {
      const value = params.get(name);
      if (value) return value;
    }
    for (const [key, value] of params.entries()) {
      if (value && names.some((name) => key.toLowerCase() === name.toLowerCase())) return value;
    }
    return "";
  };
  const override = {};
  const delimiter = firstParam("delimiter", "d");
  if (delimiter) {
    const normalized = normalizeDelimiterOverride(delimiter);
    override.delimiter = normalized.delimiter;
    override.delimiterName = normalized.delimiterName;
  }
  const quoteChar = firstParam("quote", "q");
  if (quoteChar) override.quoteChar = quoteChar;
  const quoteMode = firstParam("quotemode", "quote-mode", "quoteMode", "quote_mode", "qm");
  if (quoteMode) override.quoteMode = quoteMode;
  const quoteEscape = firstParam("quote-escape", "quoteEscape", "quote_escape", "qe");
  if (quoteEscape) override.quoteEscape = normalizeQuoteEscapeValue(quoteEscape);
  const encodingName = firstParam("encoding", "e");
  if (encodingName) override.encodingName = encodingName;
  const lineEnding = firstParam("lineending", "line-ending", "lineEnding", "line_ending", "newline", "eol", "le");
  if (lineEnding) override.lineEnding = normalizeLineEndingName(lineEnding);
  const headerRows = firstParam("headerrows", "header-rows", "headerRows");
  if (headerRows) {
    const rows = Number.parseInt(headerRows, 10);
    if (Number.isFinite(rows) && rows >= 0) override.hasHeader = rows > 0;
  }
  const headerColumns = firstParam("headercolumns", "header-columns", "headerColumns", "header_cols");
  if (headerColumns) {
    const cols = Number.parseInt(headerColumns, 10);
    if (Number.isFinite(cols) && cols >= 0) override.headerColumns = cols;
  }
  const header = firstParam("header");
  if (header) override.hasHeader = truthyParam(header);
  const noHeader = firstParam("noheader", "no-header", "noHeader");
  if (noHeader && truthyParam(noHeader)) override.hasHeader = false;
  const sheetName = firstParam("sheet", "sheetName", "worksheet", "xlsx-sheet");
  if (sheetName) override.sheetName = sheetName;
  const sheetIndex = firstParam("sheet-index", "sheetIndex", "worksheetIndex", "xlsx-sheet-index");
  if (sheetIndex) {
    const index = Number.parseInt(sheetIndex, 10);
    if (Number.isFinite(index) && index > 0) override.sheetIndex = index;
  }
  return override;
}

function locationFromParams(params) {
  const raw = params.get("goto") || params.get("location") || params.get("range");
  if (raw) return parseURLLocation(raw);
  const startRow = params.get("row");
  if (!startRow) return null;
  const endRow = params.get("endRow") || params.get("end-row") || startRow;
  const startCol = params.get("col") || params.get("column") || 1;
  return {
    startRow,
    startCol,
    endRow,
    endCol:
      params.get("endCol") ||
      params.get("end-column") ||
      params.get("endColumn") ||
      (endRow !== startRow && !params.get("col") && !params.get("column") ? "last" : startCol),
  };
}

function parseURLLocation(raw) {
  const value = String(raw).trim();
  let match = value.match(/^([0-9]+):([0-9]+)-([0-9]+):([0-9]+)$/);
  if (match) return orderedURLLocation(match[1], match[2], match[3], match[4]);
  match = value.match(/^([0-9]+)-([0-9]+)$/);
  if (match) return orderedURLLocation(match[1], 1, match[2], "last");
  match = value.match(/^([0-9]+):([0-9]+)$/);
  if (match) return { startRow: match[1], startCol: match[2] };
  match = value.match(/^([0-9]+)$/);
  if (match) return { startRow: match[1], startCol: 1 };
  return null;
}

function orderedURLLocation(startRow, startCol, endRow, endCol) {
  const rowStart = Number.parseInt(String(startRow), 10);
  const rowEnd = Number.parseInt(String(endRow), 10);
  const colStart = Number.parseInt(String(startCol), 10);
  const colEnd = endCol === "last" ? "last" : Number.parseInt(String(endCol), 10);
  const orderedStartRow = Number.isFinite(rowStart) && Number.isFinite(rowEnd) ? Math.min(rowStart, rowEnd) : startRow;
  const orderedEndRow = Number.isFinite(rowStart) && Number.isFinite(rowEnd) ? Math.max(rowStart, rowEnd) : endRow;
  const orderedStartCol = Number.isFinite(colStart) && Number.isFinite(colEnd) ? Math.min(colStart, colEnd) : startCol;
  const orderedEndCol = Number.isFinite(colStart) && Number.isFinite(colEnd) ? Math.max(colStart, colEnd) : endCol;
  return {
    startRow: String(orderedStartRow),
    startCol: String(orderedStartCol),
    endRow: String(orderedEndRow),
    endCol: orderedEndCol === "last" ? "last" : String(orderedEndCol),
  };
}

function normalizeDelimiterOverride(value) {
  const raw = String(value);
  switch (raw.toLowerCase()) {
    case "comma":
      return { delimiter: ",", delimiterName: "comma" };
    case "tab":
    case "\\t":
      return { delimiter: "\t", delimiterName: "tab" };
    case "semicolon":
      return { delimiter: ";", delimiterName: "semicolon" };
    case "pipe":
      return { delimiter: "|", delimiterName: "pipe" };
    case "space":
      return { delimiter: " ", delimiterName: "space" };
    default:
      return { delimiter: raw, delimiterName: raw };
  }
}

openButton.addEventListener("click", openPathInputOrPicker);
fileList.addEventListener("scroll", syncOpenFileTabOverflowState, { passive: true });

fileInput.addEventListener("change", async () => {
  const file = fileInput.files?.[0];
  if (!file) return;
  try {
    await uploadAndOpen(file);
  } catch (error) {
    status.textContent = `Open failed: ${error.message}`;
  } finally {
    fileInput.value = "";
  }
});

saveButton.addEventListener("click", saveCurrentDocument);
saveAsButton.addEventListener("click", openSaveAsDialog);
restoreBackupButton.addEventListener("click", openRestoreBackupPanel);
splitEditorButton.addEventListener("click", () => executeCommand("workbench.splitEditor"));
splitOrientationButton.addEventListener("click", () => executeCommand("workbench.toggleSplitOrientation"));
focusEditorGroupButton.addEventListener("click", () => executeCommand("workbench.focusNextGroup"));
moveEditorTabButton.addEventListener("click", () => executeCommand("workbench.moveTabToOtherGroup"));
closeSplitEditorButton.addEventListener("click", () => executeCommand("workbench.closeSplitEditor"));
primaryEditorGroup.addEventListener("focus", () => {
  void focusEditorGroup("primary");
});
secondaryEditorGroup.addEventListener("focus", () => {
  if (!editorGroupsSplit) return;
  void focusEditorGroup("secondary");
});
focusSecondaryGroupButton.addEventListener("click", () => {
  void focusEditorGroup("secondary");
});
closeSecondaryGroupButton.addEventListener("click", closeSplitEditorGroup);
for (const tab of [primaryEditorTab, secondaryEditorTab]) {
  tab.addEventListener("click", () => {
    void focusEditorGroup(tab === secondaryEditorTab ? "secondary" : "primary");
  });
  tab.addEventListener("mousedown", editorTabPointerDown);
  tab.addEventListener("dragstart", editorTabDragStart);
  tab.addEventListener("dragend", editorTabDragEnd);
  tab.addEventListener("contextmenu", openEditorTabContextMenu);
}
document.addEventListener("mousemove", editorTabPointerMove);
document.addEventListener("mouseup", editorTabPointerUp);
for (const group of [primaryEditorGroup, secondaryEditorGroup]) {
  group.addEventListener("dragover", editorGroupDragOver);
  group.addEventListener("dragleave", editorGroupDragLeave);
  group.addEventListener("drop", editorGroupDrop);
}
exportButton.addEventListener("click", () => exportCurrentView(exportFormatSelect?.value || "csv"));
exportXlsxButton.addEventListener("click", () => exportCurrentView("xlsx"));
readOnlyButton.addEventListener("click", toggleReadOnly);
gridTabButton.addEventListener("click", () => setViewMode("grid"));
sqlTabButton.addEventListener("click", () => {
  if (!sqlTabButton.disabled) setViewMode("sql");
});
profileTabButton.addEventListener("click", openProfilePanel);
runProfileButton.addEventListener("click", () => executeCommand("profile.run"));
cancelProfileButton.addEventListener("click", () => executeCommand("profile.cancel"));
toggleProfileSchemaButton.addEventListener("click", toggleProfileSchemaEditor);
profileSchemaInput.addEventListener("input", updateProfileSchemaDraft);
saveProfileSchemaButton.addEventListener("click", saveProfileSchema);
resetProfileSchemaButton.addEventListener("click", resetProfileSchema);
savedProfileSchemaSelect.addEventListener("change", () => syncSavedProfileSchemaControls());
applyProfileSchemaButton.addEventListener("click", applySavedProfileSchema);
deleteProfileSchemaButton.addEventListener("click", deleteSavedProfileSchema);
exportProfileSchemasButton.addEventListener("click", exportProfileSchemas);
showImportProfileSchemasButton.addEventListener("click", showProfileSchemaImportPanel);
profileSchemaImportText.addEventListener("input", () => syncSavedProfileSchemaControls());
applyProfileSchemaImportButton.addEventListener("click", applyProfileSchemaImport);
cancelProfileSchemaImportButton.addEventListener("click", cancelProfileSchemaImport);
saveProfileRuleButton.addEventListener("click", saveProfileRuleFromBuilder);
clearProfileRuleButton.addEventListener("click", () => clearProfileRuleBuilder());
settingsTabButton.addEventListener("click", openSettingsPanel);
shortcutsTabButton.addEventListener("click", () => executeCommand("settings.keyboardShortcuts"));
closeAboutSupportButton.addEventListener("click", closeAboutSupportPanel);
aboutOpenSettingsButton.addEventListener("click", () => {
  closeAboutSupportPanel();
  openSettingsPanel();
});
aboutKeyboardShortcutsButton.addEventListener("click", () => {
  closeAboutSupportPanel();
  openKeyboardShortcuts();
});
aboutCheckUpdatesButton.addEventListener("click", runAboutSupportUpdateCheck);
paletteButton.addEventListener("click", openCommandPalette);
copyFormatButton.addEventListener("click", () => executeCommand("edit.copySelectedFormat"));
copyFormatSelect.addEventListener("change", () => {
  document.body.dataset.copyFormat = `${copyFormatSelect.value}:${copyHeadersToggle.checked ? "headers" : "body"}`;
});
copyHeadersToggle.addEventListener("change", () => {
  setHeaderInSelection(copyHeadersToggle.checked);
  draw();
});
aggregateStatusButton.addEventListener("click", async (event) => {
  if (event.metaKey || event.ctrlKey) {
    await copyActiveAggregate();
    return;
  }
  cycleAggregateMetric();
});
headerStatus.addEventListener("click", () => executeCommand("table.toggleHeaderRow"));
workbookSheetSelect?.addEventListener("change", () => switchWorkbookSheet());
externalStatus.addEventListener("click", openExternalChangePanel);
saveFormatPresetButton.addEventListener("click", () => executeCommand("format.savePreset"));
applyFormatPresetButton.addEventListener("click", () => executeCommand("format.applyPreset"));
formatPresetSelect.addEventListener("change", () => updateCommandState());
formatRuleInput.addEventListener("keydown", (event) => {
  if (event.key !== "Enter") return;
  event.preventDefault();
  executeCommand("format.saveRule");
});
saveFormatRuleButton.addEventListener("click", () => executeCommand("format.saveRule"));
settingsSearch.addEventListener("input", renderSettingsPanel);
resetFormatPresetsButton.addEventListener("click", resetFormatPresets);
localeSelect.addEventListener("change", updateLocaleFromSettings);
resetLocaleButton.addEventListener("click", resetLocalePreference);
exportSavedViewsButton.addEventListener("click", () => executeCommand("view.exportSavedSets"));
showImportSavedViewsButton.addEventListener("click", () => executeCommand("view.importSavedSets"));
savedViewImportText.addEventListener("input", syncSavedViewImportControls);
applySavedViewImportButton.addEventListener("click", applySavedViewImport);
cancelSavedViewImportButton.addEventListener("click", cancelSavedViewImport);
exportFormatPresetsButton.addEventListener("click", () => executeCommand("format.exportPresets"));
showImportFormatPresetsButton.addEventListener("click", () => executeCommand("format.importPresets"));
formatPresetImportText.addEventListener("input", syncFormatPresetImportControls);
applyFormatPresetImportButton.addEventListener("click", applyFormatPresetImport);
cancelFormatPresetImportButton.addEventListener("click", cancelFormatPresetImport);
themePreferenceSelect.addEventListener("change", updateThemeFromSettings);
resetThemePreferenceButton.addEventListener("click", resetThemePreference);
autoRefreshModeSelect.addEventListener("change", updateAutoRefreshModeFromSettings);
resetAutoRefreshButton.addEventListener("click", resetAutoRefreshMode);
updateCheckModeSelect.addEventListener("change", updateUpdateCheckModeFromSettings);
checkForUpdatesButton.addEventListener("click", () => executeCommand("settings.checkForUpdates"));
resetUpdateChecksButton.addEventListener("click", resetUpdateCheckMode);
filterHistoryLimitSelect.addEventListener("change", updateFilterHistoryLimitFromSettings);
clearRecentFiltersButton.addEventListener("click", clearRecentFilters);
resetSavedViewsButton.addEventListener("click", resetSavedViewsSettings);
applyCustomLinkProtocolsButton.addEventListener("click", updateCustomLinkProtocolsFromSettings);
resetCustomLinkProtocolsButton.addEventListener("click", resetCustomLinkProtocols);
gridFontSelect.addEventListener("change", updateEditorAppearanceFromControls);
gridFontSizeInput.addEventListener("change", updateEditorAppearanceFromControls);
gridFontSizeInput.addEventListener("input", updateEditorAppearanceFromControls);
gridRowHeightSelect.addEventListener("change", updateEditorAppearanceFromControls);
textEditorFontSelect.addEventListener("change", updateEditorAppearanceFromControls);
textEditorFontSizeInput.addEventListener("change", updateEditorAppearanceFromControls);
textEditorFontSizeInput.addEventListener("input", updateEditorAppearanceFromControls);
resetEditorAppearanceButton.addEventListener("click", resetEditorAppearance);
applyDelimiterButton.addEventListener("click", () => executeCommand("format.applyDelimiter"));
delimiterSelect.addEventListener("change", () => {
  delimiterDraftDirty = true;
  if (delimiterSelect.value !== "custom") delimiterCustomInput.value = "";
  syncDelimiterControls();
});
delimiterCustomInput.addEventListener("input", () => {
  delimiterDraftDirty = true;
  delimiterSelect.value = "custom";
  syncDelimiterControls();
});
applyEncodingButton.addEventListener("click", () => executeCommand("format.applyEncoding"));
encodingSelect.addEventListener("change", () => {
  encodingDraftDirty = true;
  document.body.dataset.encodingDraft = encodingSelect.value;
});
applyQuotePolicyButton.addEventListener("click", () => executeCommand("format.applyQuotePolicy"));
applyLineEndingButton.addEventListener("click", () => executeCommand("format.applyLineEnding"));
lineEndingSelect.addEventListener("change", () => {
  lineEndingDraftDirty = true;
  document.body.dataset.lineEndingDraft = normalizeLineEndingName(lineEndingSelect.value);
});
exportQuoteRulesButton.addEventListener("click", () => executeCommand("format.exportQuoteRules"));
showImportQuoteRulesButton.addEventListener("click", () => executeCommand("format.importQuoteRules"));
resetQuotePolicyButton.addEventListener("click", resetQuotePolicySettings);
quoteRulesImportText.addEventListener("input", syncQuoteRulesImportControls);
applyQuoteRulesImportButton.addEventListener("click", applyQuoteRulesImport);
cancelQuoteRulesImportButton.addEventListener("click", cancelQuoteRulesImport);
toggleCellQuoteButton.addEventListener("click", () => executeCommand("format.toggleCellQuote"));
showHiddenColumnsButton.addEventListener("click", () => executeCommand("column.showAll"));
quoteHeaderRuleToggle.addEventListener("change", () => toggleQuoteRule("header", quoteHeaderRuleToggle.checked));
quoteBodyRuleToggle.addEventListener("change", () => toggleQuoteRule("body", quoteBodyRuleToggle.checked));
quoteColumnRuleToggle.addEventListener("change", () => toggleQuoteRule("column", quoteColumnRuleToggle.checked));
addQuoteCellRuleButton.addEventListener("click", () => applyActiveCellQuoteRule("quote"));
addNoQuoteCellRuleButton.addEventListener("click", () => applyActiveCellQuoteRule("no-quote"));
clearCellQuoteRuleButton.addEventListener("click", () => applyActiveCellQuoteRule("clear"));
quoteSelectionRuleButton.addEventListener("click", () => setSelectedCellQuoteOverrides("quote"));
noQuoteSelectionRuleButton.addEventListener("click", () => setSelectedCellQuoteOverrides("no-quote"));
clearSelectionQuoteRuleButton.addEventListener("click", () => setSelectedCellQuoteOverrides("clear"));
applyQuoteColumnRuleButton.addEventListener("click", applyColumnQuoteRule);
applyQuoteColumnPatternRuleButton.addEventListener("click", applyColumnPatternQuoteRule);
applyQuoteValuePatternRuleButton.addEventListener("click", applyValuePatternQuoteRule);
applyQuoteRowRangeRuleButton.addEventListener("click", applyRowRangeQuoteRule);
applyQuoteCoordinateRuleButton.addEventListener("click", applyCoordinateQuoteRule);
quoteModeSelect.addEventListener("change", () => {
  quotePolicyDraftDirty = true;
  settingsQuoteSummary.textContent = selectedQuotePolicySummary();
});
quoteCharSelect.addEventListener("change", () => {
  quotePolicyDraftDirty = true;
  settingsQuoteSummary.textContent = selectedQuotePolicySummary();
});
quoteEscapeSelect.addEventListener("change", () => {
  quotePolicyDraftDirty = true;
  settingsQuoteSummary.textContent = selectedQuotePolicySummary();
});
runSqlButton.addEventListener("click", runSQLQuery);
cancelSqlButton.addEventListener("click", cancelSQLRequest);
sqlResultFormatSelect?.addEventListener("change", syncSQLResultActions);
openSqlResultButton.addEventListener("click", openSQLResultAsDocument);
copySqlResultButton.addEventListener("click", copyCurrentSQLResult);
exportSqlResultButton.addEventListener("click", exportCurrentSQLResult);
clearSqlHistoryButton.addEventListener("click", clearSQLHistory);
sqlResult.addEventListener("click", (event) => {
  const cell = sqlResultCellFromEvent(event);
  if (!cell) return;
  selectSQLResultCell(cell.row, cell.col, event.shiftKey);
});
sqlResult.addEventListener("contextmenu", (event) => {
  const cell = sqlResultCellFromEvent(event);
  if (!cell) return;
  event.preventDefault();
  closeColumnMenu();
  renderSQLResultContextMenu(cell, event.clientX, event.clientY);
});
sqlResult.addEventListener("keydown", async (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "c") {
    event.preventDefault();
    await copyCurrentSQLResult();
    return;
  }
  if (event.key === "Escape") {
    event.preventDefault();
    clearSQLResultSelection();
    status.textContent = t("sql.status.selectionCleared");
    return;
  }
  if (event.key !== "Enter" && event.key !== " ") return;
  const cell = sqlResultCellFromEvent(event);
  if (!cell) return;
  event.preventDefault();
  selectSQLResultCell(cell.row, cell.col, event.shiftKey);
});
diffButton.addEventListener("click", openDiffPanel);
runDiffButton.addEventListener("click", runDiffCompare);
diffPrevButton.addEventListener("click", () => navigateDiff(-1));
diffNextButton.addEventListener("click", () => navigateDiff(1));
diffStatusFilter.addEventListener("change", applyDiffResultFilters);
diffSearchInput.addEventListener("input", applyDiffResultFilters);
diffSortSelect.addEventListener("change", applyDiffResultFilters);
diffSortDirectionButton.addEventListener("click", toggleDiffSortDirection);
diffClearFilterButton.addEventListener("click", clearDiffResultFilters);
diffRightTextMode.addEventListener("change", () => {
  setDiffRightTextMode(diffRightTextMode.checked);
  if (diffRightTextMode.checked) diffRightText.focus();
});
diffPasteClipboardButton.addEventListener("click", pasteDiffClipboardText);
diffRightText.addEventListener("input", () => {
  diffRightText.dataset.source = "";
  document.body.dataset.diffTextSource = diffRightTextMode.checked ? "text" : "";
});
for (const input of [leftDiffPath, rightDiffPath]) {
  input.addEventListener("keydown", async (event) => {
    if (event.key === "Enter") {
      event.preventDefault();
      await runDiffCompare();
    }
  });
}
sqlInput.addEventListener("keydown", async (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
    event.preventDefault();
    await runSQLQuery();
  }
});
paletteInput.addEventListener("input", () => {
  paletteIndex = 0;
  renderCommandPalette();
});
paletteInput.addEventListener("keydown", async (event) => {
  if (event.key === "Escape") {
    event.preventDefault();
    closeCommandPalette();
    return;
  }
  if (event.key === "ArrowDown") {
    event.preventDefault();
    paletteIndex = Math.min(paletteItems.length - 1, paletteIndex + 1);
    renderCommandPalette();
    return;
  }
  if (event.key === "ArrowUp") {
    event.preventDefault();
    paletteIndex = Math.max(0, paletteIndex - 1);
    renderCommandPalette();
    return;
  }
  if (event.key === "Enter") {
    event.preventDefault();
    const command = paletteItems[paletteIndex];
    if (command) await executePaletteCommand(command.id);
  }
});
commandPalette.addEventListener("mousedown", (event) => {
  if (event.target === commandPalette) {
    event.preventDefault();
    closeCommandPalette();
  }
});
shortcutsSearch.addEventListener("input", renderKeyboardShortcuts);
closeShortcutsButton.addEventListener("click", closeKeyboardShortcuts);
applyShortcutPresetButton.addEventListener("click", applyShortcutPreset);
exportShortcutsButton.addEventListener("click", () => executeCommand("settings.exportShortcuts"));
showImportShortcutsButton.addEventListener("click", () => executeCommand("settings.importShortcuts"));
resetAllShortcutsButton.addEventListener("click", resetAllShortcuts);
shortcutImportText.addEventListener("input", syncShortcutImportControls);
applyShortcutImportButton.addEventListener("click", applyShortcutImport);
cancelShortcutImportButton.addEventListener("click", cancelShortcutImport);
shortcutsPanel.addEventListener("mousedown", (event) => {
  if (event.target === shortcutsPanel) {
    event.preventDefault();
    closeKeyboardShortcuts();
  }
});
closeSaveAsButton.addEventListener("click", closeSaveAsDialog);
cancelSaveAsButton.addEventListener("click", closeSaveAsDialog);
confirmSaveAsButton.addEventListener("click", confirmSaveAs);
closeDirtyCloseButton.addEventListener("click", closeDirtyClosePanel);
cancelDirtyCloseButton.addEventListener("click", closeDirtyClosePanel);
discardDirtyCloseButton.addEventListener("click", confirmDirtyCloseDiscard);
closeRestoreBackupButton.addEventListener("click", closeRestoreBackupPanel);
cancelRestoreBackupButton.addEventListener("click", closeRestoreBackupPanel);
refreshRestoreBackupsButton.addEventListener("click", loadRestoreBackups);
closeExternalChangeButton.addEventListener("click", closeExternalChangePanel);
keepExternalChangeButton.addEventListener("click", closeExternalChangePanel);
compareExternalChangeButton.addEventListener("click", compareExternalChangeWithDisk);
saveAsExternalChangeButton.addEventListener("click", saveAsExternalChange);
overwriteExternalChangeButton.addEventListener("click", overwriteExternalChangeOnDisk);
reloadExternalChangeButton.addEventListener("click", reloadExternalChangeFromDisk);
closeTransposeConfirmButton.addEventListener("click", closeTransposeConfirmPanel);
cancelTransposeConfirmButton.addEventListener("click", closeTransposeConfirmPanel);
confirmTransposeButton.addEventListener("click", confirmWholeTableTranspose);
closeSplitConfirmButton.addEventListener("click", closeSplitOverwriteConfirmPanel);
cancelSplitConfirmButton.addEventListener("click", closeSplitOverwriteConfirmPanel);
confirmSplitConfirmButton.addEventListener("click", confirmSplitOverwrite);
closeCellValueDialogButton.addEventListener("click", closeCellValueDialog);
applyCellValueDialogButton.addEventListener("click", () => {
  cellValueInput.value = cellValueDialogInput.value;
  applyCellValue();
});
revertCellValueDialogButton.addEventListener("click", revertCellValueInput);
saveAsPanel.addEventListener("mousedown", (event) => {
  if (event.target === saveAsPanel) {
    event.preventDefault();
    closeSaveAsDialog();
  }
});
restoreBackupPanel.addEventListener("mousedown", (event) => {
  if (event.target === restoreBackupPanel) {
    event.preventDefault();
    closeRestoreBackupPanel();
  }
});
externalChangePanel.addEventListener("mousedown", (event) => {
  if (event.target === externalChangePanel) {
    event.preventDefault();
    closeExternalChangePanel();
  }
});
dirtyClosePanel.addEventListener("mousedown", (event) => {
  if (event.target === dirtyClosePanel) {
    event.preventDefault();
    closeDirtyClosePanel();
  }
});
transposeConfirmPanel.addEventListener("mousedown", (event) => {
  if (event.target === transposeConfirmPanel) {
    event.preventDefault();
    closeTransposeConfirmPanel();
  }
});
splitConfirmPanel.addEventListener("mousedown", (event) => {
  if (event.target === splitConfirmPanel) {
    event.preventDefault();
    closeSplitOverwriteConfirmPanel();
  }
});
cellValueDialog.addEventListener("mousedown", (event) => {
  if (event.target === cellValueDialog) {
    event.preventDefault();
    closeCellValueDialog();
  }
});
saveAsPath.addEventListener("keydown", async (event) => {
  if (event.key === "Enter") {
    event.preventDefault();
    await confirmSaveAs();
    return;
  }
  if (event.key === "Escape") {
    event.preventDefault();
    closeSaveAsDialog();
  }
});
columnMenu.addEventListener("mousedown", (event) => {
  event.stopPropagation();
});
gridContextMenu.addEventListener("mousedown", (event) => {
  event.stopPropagation();
});
window.addEventListener("mousedown", (event) => {
  if (!columnMenu.hidden && !columnMenu.contains(event.target) && event.target !== canvas) {
    closeColumnMenu();
  }
  if (!gridContextMenu.hidden && !gridContextMenu.contains(event.target) && event.target !== canvas) {
    closeGridContextMenu();
  }
});
columnSortAsc.addEventListener("click", () => executeCommand("column.sortAscending"));
columnSortDesc.addEventListener("click", () => executeCommand("column.sortDescending"));
columnSortLengthAsc.addEventListener("click", () => executeCommand("column.sortLengthAscending"));
columnSortLengthDesc.addEventListener("click", () => executeCommand("column.sortLengthDescending"));
columnAddSortAsc.addEventListener("click", () => executeCommand("column.addSortAscending"));
columnAddSortDesc.addEventListener("click", () => executeCommand("column.addSortDescending"));
columnStableSortToggle.addEventListener("change", () => {
  sortStable = columnStableSortToggle.checked;
  document.body.dataset.sortStable = String(sortStable);
  document.body.dataset.columnMenuStableSort = String(sortStable);
});
columnFilterValue.addEventListener("click", () => executeCommand("column.filterBySelectedValue"));
columnFilterBlank.addEventListener("click", () => executeCommand("column.filterBlanks"));
columnFilterFilled.addEventListener("click", () => executeCommand("column.filterNonBlanks"));
columnTransformMode.addEventListener("change", syncColumnTransformCommand);
columnTransformApply.addEventListener("click", () => executeCommand(syncColumnTransformCommand()));
columnApplyValues.addEventListener("click", applyColumnValueChecklist);
columnHide.addEventListener("click", () => {
  const column = columnMenuColumn;
  closeColumnMenu();
  if (column >= 0) {
    selected = { row: selected.row, col: column };
    selectionAnchor = { ...selected };
    hideSelectedColumn();
  }
});
columnShowAll.addEventListener("click", () => {
  closeColumnMenu();
  showAllHiddenColumns();
});
columnClearView.addEventListener("click", async () => {
  closeColumnMenu();
  await clearView();
});
cellValueInput.addEventListener("input", () => {
  syncCellValueDraft(cellValueInput.value);
});
cellValueInput.addEventListener("keydown", (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
    event.preventDefault();
    event.stopPropagation();
    applyCellValue();
    return;
  }
  if (event.key === "Escape" && cellValueInputDirty) {
    event.preventDefault();
    event.stopPropagation();
    revertCellValueInput();
  }
});
applyCellValueButton.addEventListener("click", applyCellValue);
openCellValueDialogButton.addEventListener("click", openCellValueDialog);
toggleCellValueLargeButton.addEventListener("click", toggleCellValueLarge);
revertCellValueButton.addEventListener("click", revertCellValueInput);
cellSeparatorInput.addEventListener("input", updateCellSeparatorFromInput);
cellSeparatorInput.addEventListener("change", () => setCellSeparatorToken(cellSeparatorInput.value));
cellValueDialogInput.addEventListener("input", () => {
  syncCellValueDraft(cellValueDialogInput.value, "dialog");
});
cellValueDialogInput.addEventListener("keydown", (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
    event.preventDefault();
    event.stopPropagation();
    cellValueInput.value = cellValueDialogInput.value;
    applyCellValue();
    return;
  }
  if (event.key === "Escape" && cellValueInputDirty) {
    event.preventDefault();
    event.stopPropagation();
    revertCellValueInput();
    cellValueDialogInput.focus();
    return;
  }
  if (event.key === "Escape") {
    event.preventDefault();
    event.stopPropagation();
    closeCellValueDialog();
  }
});
openCellLink.addEventListener("click", (event) => {
  event.preventDefault();
  if (openCellLink.getAttribute("aria-disabled") === "true") {
    status.textContent = linkUnavailableMessage({
      status: document.body.dataset.cellLinkStatus,
      protocol: document.body.dataset.cellLinkProtocol,
    });
    return;
  }
  executeCommand("cell.openLink");
});
cancelOpenLinkButton.addEventListener("click", cancelCustomLinkConfirmDialog);
confirmOpenLinkButton.addEventListener("click", confirmCustomLinkOpen);
linkConfirmDialog.addEventListener("mousedown", (event) => {
  if (event.target === linkConfirmDialog) {
    event.preventDefault();
    cancelCustomLinkConfirmDialog();
  }
});
findButton.addEventListener("click", (event) => focusFindInput({ run: true, target: event.target }));
cancelFindButton.addEventListener("click", cancelFindRequest);
findPrevButton.addEventListener("click", () => findNextMatch(-1));
findNextButton.addEventListener("click", () => findNextMatch(1));
replacePreviewButton.addEventListener("click", (event) => previewReplace({ target: event.target }));
replaceButton.addEventListener("click", (event) => replaceCurrentQuery(false, { target: event.target }));
replaceAllButton.addEventListener("click", (event) => replaceCurrentQuery(true, { target: event.target }));
findCaseButton.addEventListener("click", () => toggleFindOption("case"));
findWholeButton.addEventListener("click", () => toggleFindOption("whole"));
findRegexButton.addEventListener("click", () => toggleFindOption("regex"));
findSelectionButton.addEventListener("click", () => toggleFindOption("selection"));
filterButton.addEventListener("click", filterCurrentQuery);
recentFilterSelect.addEventListener("change", applyFilterDropdownSelection);
filterBuilderButton.addEventListener("click", openConditionBuilder);
sortButton.addEventListener("click", sortSelectedColumn);
closeConditionButton.addEventListener("click", closeConditionBuilder);
conditionOperator.addEventListener("change", updateConditionValueState);
conditionPrimaryOpenGroup.addEventListener("change", updateConditionValueState);
conditionPrimaryCloseGroup.addEventListener("change", updateConditionValueState);
useConditionTemplateButton.addEventListener("click", useConditionTemplate);
addConditionButton.addEventListener("click", addConditionRow);
clearConditionGroupsButton.addEventListener("click", clearConditionGroups);
groupAllConditionsButton.addEventListener("click", groupAllConditionRows);
conditionAllAndButton.addEventListener("click", () => setAllConditionJoins("AND"));
conditionAllOrButton.addEventListener("click", () => setAllConditionJoins("OR"));
conditionRows.addEventListener("change", updateConditionValueState);
conditionRows.addEventListener("input", () => {
  document.body.dataset.conditionRowCount = String(1 + conditionRows.querySelectorAll(".condition-row-extra").length);
});
conditionRows.addEventListener("click", (event) => {
  const button = event.target.closest("[data-condition-action]");
  if (!button || button.disabled) return;
  const row = button.closest(".condition-row-extra");
  const index = Number(row?.dataset.conditionIndex ?? -1);
  if (button.dataset.conditionAction === "remove") {
    removeConditionRow(index);
  } else if (button.dataset.conditionAction === "duplicate") {
    duplicateConditionRow(index);
  } else if (button.dataset.conditionAction === "group-previous") {
    groupConditionWithPrevious(index, 1);
  } else if (button.dataset.conditionAction === "ungroup-previous") {
    groupConditionWithPrevious(index, -1);
  } else if (button.dataset.conditionAction === "move-up") {
    moveConditionRow(index, -1);
  } else if (button.dataset.conditionAction === "move-down") {
    moveConditionRow(index, 1);
  }
});
applyConditionButton.addEventListener("click", applyConditionBuilder);
clearConditionButton.addEventListener("click", async () => {
  closeConditionBuilder();
  await clearView();
});
conditionPanel.addEventListener("mousedown", (event) => {
  if (event.target === conditionPanel) {
    event.preventDefault();
    closeConditionBuilder();
  }
});
for (const input of [
  conditionColumn,
  conditionOperator,
  conditionValue,
]) {
  input.addEventListener("keydown", async (event) => {
    if (event.key !== "Enter") return;
    event.preventDefault();
    await applyConditionBuilder();
  });
}
conditionRows.addEventListener("keydown", async (event) => {
  if (event.key !== "Enter") return;
  event.preventDefault();
  await applyConditionBuilder();
});
commandInput.addEventListener("keydown", async (event) => {
  if (event.key !== "Enter") return;
  event.preventDefault();
  await findNextMatch(event.shiftKey ? -1 : 1);
});
commandInput.addEventListener("input", invalidateFindState);
replaceInput.addEventListener("input", clearReplacePreview);
window.addEventListener("resize", resize);
window.addEventListener("pagehide", saveWindowGeometryNow);
if ("ResizeObserver" in window) {
  const gridResizeObserver = new ResizeObserver(resize);
  gridResizeObserver.observe(gridFrame);
}
canvas.tabIndex = 0;
profilePanel.tabIndex = -1;
loadCommandRegistry();
resize();
renderFileList();
renderSQLHistory();
renderFormatPresetSelect();
refreshFilterDropdown();
renderSettingsPanel();
syncStatusbarSurface("grid");
syncPlatformDialogButtonOrder();
openInitialURLTarget();
maybeRunStartupUpdateCheck();
