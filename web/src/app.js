const canvas = document.querySelector("#grid");
const ctx = canvas.getContext("2d");
const editor = document.querySelector("#editor");
const status = document.querySelector("#status");
const encodingStatus = document.querySelector("#encodingStatus");
const delimiterStatus = document.querySelector("#delimiterStatus");
const quoteStatus = document.querySelector("#quoteStatus");
const headerStatus = document.querySelector("#headerStatus");
const readOnlyStatus = document.querySelector("#readOnlyStatus");
const externalStatus = document.querySelector("#externalStatus");
const aggregateStatusButton = document.querySelector("#aggregateStatusButton");
const formatPresetSelect = document.querySelector("#formatPresetSelect");
const saveFormatPresetButton = document.querySelector("#saveFormatPresetButton");
const applyFormatPresetButton = document.querySelector("#applyFormatPresetButton");
const formatRuleInput = document.querySelector("#formatRuleInput");
const saveFormatRuleButton = document.querySelector("#saveFormatRuleButton");
const columnName = document.querySelector("#columnName");
const columnType = document.querySelector("#columnType");
const columnNulls = document.querySelector("#columnNulls");
const columnDistinct = document.querySelector("#columnDistinct");
const cellQuote = document.querySelector("#cellQuote");
const formatSummary = document.querySelector("#formatSummary");
const cellValueInput = document.querySelector("#cellValueInput");
const applyCellValueButton = document.querySelector("#applyCellValueButton");
const openCellLink = document.querySelector("#openCellLink");
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
const filterBuilderButton = document.querySelector("#filterBuilderButton");
const sortButton = document.querySelector("#sortButton");
const findButton = document.querySelector("#findButton");
const findPrevButton = document.querySelector("#findPrevButton");
const findNextButton = document.querySelector("#findNextButton");
const findMatchCount = document.querySelector("#findMatchCount");
const replaceButton = document.querySelector("#replaceButton");
const replaceAllButton = document.querySelector("#replaceAllButton");
const findCaseButton = document.querySelector("#findCaseButton");
const findWholeButton = document.querySelector("#findWholeButton");
const findRegexButton = document.querySelector("#findRegexButton");
const findSelectionButton = document.querySelector("#findSelectionButton");
const diffButton = document.querySelector("#diffButton");
const exportButton = document.querySelector("#exportButton");
const exportXlsxButton = document.querySelector("#exportXlsxButton");
const saveButton = document.querySelector("#saveButton");
const readOnlyButton = document.querySelector("#readOnlyButton");
const commandButtons = Array.from(document.querySelectorAll("[data-command]"));
const gridFrame = document.querySelector(".grid-frame");
const gridTabButton = document.querySelector("#gridTabButton");
const sqlTabButton = document.querySelector("#sqlTabButton");
const profileTabButton = document.querySelector("#profileTabButton");
const settingsTabButton = document.querySelector("#settingsTabButton");
const paletteButton = document.querySelector("#paletteButton");
const sqlPanel = document.querySelector("#sqlPanel");
const sqlInput = document.querySelector("#sqlInput");
const runSqlButton = document.querySelector("#runSqlButton");
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
const diffIgnoreWhitespace = document.querySelector("#diffIgnoreWhitespace");
const diffMatchHeaders = document.querySelector("#diffMatchHeaders");
const diffKeyColumn = document.querySelector("#diffKeyColumn");
const diffIgnoredColumns = document.querySelector("#diffIgnoredColumns");
const diffPrevButton = document.querySelector("#diffPrevButton");
const diffNextButton = document.querySelector("#diffNextButton");
const diffRightTextMode = document.querySelector("#diffRightTextMode");
const diffRightText = document.querySelector("#diffRightText");
const diffSummary = document.querySelector("#diffSummary");
const diffResult = document.querySelector("#diffResult");
const profilePanel = document.querySelector("#profilePanel");
const profileSummary = document.querySelector("#profileSummary");
const profileColumns = document.querySelector("#profileColumns");
const profileIssues = document.querySelector("#profileIssues");
const settingsPanel = document.querySelector("#settingsPanel");
const settingsSearch = document.querySelector("#settingsSearch");
const settingsModifiedBadge = document.querySelector("#settingsModifiedBadge");
const settingsQuoteSummary = document.querySelector("#settingsQuoteSummary");
const quoteModeSelect = document.querySelector("#quoteModeSelect");
const quoteCharSelect = document.querySelector("#quoteCharSelect");
const applyQuotePolicyButton = document.querySelector("#applyQuotePolicyButton");
const settingsFormatSummary = document.querySelector("#settingsFormatSummary");
const settingsPresetList = document.querySelector("#settingsPresetList");
const resetFormatPresetsButton = document.querySelector("#resetFormatPresetsButton");
const commandPalette = document.querySelector("#commandPalette");
const paletteInput = document.querySelector("#paletteInput");
const paletteList = document.querySelector("#paletteList");
const shortcutsPanel = document.querySelector("#shortcutsPanel");
const shortcutsSearch = document.querySelector("#shortcutsSearch");
const shortcutsList = document.querySelector("#shortcutsList");
const closeShortcutsButton = document.querySelector("#closeShortcutsButton");
const conditionPanel = document.querySelector("#conditionPanel");
const conditionColumn = document.querySelector("#conditionColumn");
const conditionOperator = document.querySelector("#conditionOperator");
const conditionValue = document.querySelector("#conditionValue");
const closeConditionButton = document.querySelector("#closeConditionButton");
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
const columnFilterValue = document.querySelector("#columnFilterValue");
const columnFilterBlank = document.querySelector("#columnFilterBlank");
const columnFilterFilled = document.querySelector("#columnFilterFilled");
const columnMenuValues = document.querySelector("#columnMenuValues");
const columnApplyValues = document.querySelector("#columnApplyValues");
const columnHide = document.querySelector("#columnHide");
const columnClearView = document.querySelector("#columnClearView");

const rowHeight = 28;
const headerHeight = cssPixelValue("--grid-header-height", 48);
const rowNumberWidth = 58;
const windowSize = 160;
const minColumnWidth = 92;
const maxColumnWidth = 420;
const autoFitColumnWidthOnEdit = true;
const recentFilesStorageKey = "opencsv-recent-files";
const recentFilesLimit = 8;
const commandHistoryStorageKey = "opencsv-command-history";
const commandHistoryLimit = 8;
const customShortcutsStorageKey = "opencsv-custom-shortcuts";
const sqlHistoryStorageKey = "opencsv-sql-history";
const sqlHistoryLimit = 8;
const formatPresetsStorageKey = "opencsv-format-presets";
const formatPresetsLimit = 12;
const formatRuleSampleBytes = 65536;

let totalRows = 0;
let sourceRowCount = 0;
let aggregateRows = [];
let columns = [];
let allColumns = [];
let hiddenSourceColumns = new Set();

let scrollX = 0;
let scrollY = 0;
let wheelRemainderY = 0;
let selected = { row: 0, col: 0 };
let selectionAnchor = { row: 0, col: 0 };
let draggingSelection = false;
let dragSelectionMode = "cell";
let hoverURL = "";
let editing = false;
let documentId = "";
let activeFileEntry = null;
let metadata = null;
let rows = new Map();
let rowSource = new Map();
let windowStart = 0;
let windowEnd = 0;
let loadingWindow = false;
let windowGeneration = 0;
let dirty = false;
let dataStartRow = 0;
let viewMapped = false;
let sortDescending = false;
let sortKeys = [];
let commandsById = new Map();
let viewMode = "grid";
let profileData = null;
let lastSQLResult = null;
let sqlResultSelectionAnchor = null;
let sqlResultSelectionRange = null;
let profileLoading = false;
let paletteOpen = false;
let paletteItems = [];
let paletteIndex = 0;
let customShortcuts = new Map();
let shortcutEditorOpen = false;
let recordingShortcutID = "";
let conditionBuilderOpen = false;
let columnMenuColumn = -1;
let readOnly = false;
let externalFileChanged = false;
let fileStatusTimer = 0;
let cellValueInputDirty = false;
let cellValueInputKey = "";
let findCaseSensitive = false;
let findWholeCell = false;
let findRegex = false;
let findInSelection = false;
let findState = emptyFindState();
let includeHeadersInSelection = false;
let activeAggregateMetric = "sum";
let diffRows = [];
let activeDiffIndex = -1;
let quotePolicyDraftDirty = false;
const themeStorageKey = "opencsv-theme";
const themeMediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
let themePreference = localStorage.getItem(themeStorageKey) || "system";
let themeColors = {};

initializeTheme();

function resize() {
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

function applyTheme(preference, announce = true) {
  if (!["light", "dark", "system"].includes(preference)) preference = "system";
  themePreference = preference;
  localStorage.setItem(themeStorageKey, preference);
  const resolved = preference === "system"
    ? (themeMediaQuery.matches ? "dark" : "light")
    : preference;
  document.body.dataset.theme = resolved;
  document.body.dataset.themePreference = preference;
  refreshThemeColors();
  draw();
  if (announce && status) {
    status.textContent = `Theme: ${preference === "system" ? `System (${resolved})` : preference}`;
  }
}

function cycleTheme() {
  const order = ["light", "dark", "system"];
  const current = order.includes(themePreference) ? order.indexOf(themePreference) : 2;
  applyTheme(order[(current + 1) % order.length]);
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

function applyHiddenColumns() {
  columns = allColumns.filter((column) => !hiddenSourceColumns.has(column.sourceCol));
}

function draw() {
  if (viewMode !== "grid") {
    updateCommandState();
    return;
  }
  const { width, height } = syncCanvasSize();
  ctx.clearRect(0, 0, width, height);

  ctx.fillStyle = themeColors.gridBg;
  ctx.fillRect(0, 0, width, height);

  drawRows(width, height);
  drawBodySeparators(width, height);
  drawSelection();
  drawHeaders(width);

  const selectedColumn = columns[selected.col]?.name ?? "";
  const path = metadata?.path ? ` - ${metadata.path}` : "";
  const dirtyLabel = dirty ? " - edited" : "";
  const aggregateStats = selectionAggregateStats();
  const aggregate = selectionAggregateSummary(aggregateStats);
  const modeLabel = readOnly ? " - read-only" : "";
  const externalLabel = externalFileChanged ? " - external change" : "";
  const aggregateLabel = aggregateRows.length ? ` + ${aggregateRows.length} aggregate` : "";
  document.body.dataset.selectionRange = documentId ? selectionSummary() : "";
  document.body.dataset.columnWidths = columns.map((column) => Math.round(column.width)).join(",");
  document.body.dataset.gridHeaderHeight = String(headerHeight);
  copyHeadersToggle.checked = includeHeadersInSelection;
  updateAggregateStatus(aggregateStats);
  status.textContent =
    documentId
      ? `${sourceRowCount.toLocaleString()} rows${aggregateLabel}, ${columns.length} columns${path}${dirtyLabel}${modeLabel}${externalLabel} - row ${selected.row + 1}, ${selectedColumn}${aggregate}`
      : "Ready";
  updateFormatStatus();
  updateInspector();
  updateCommandState();
}

function setViewMode(mode) {
  closeColumnMenu();
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
      ? `${command.description} (${shortcut})`
      : command.description;
    if (button.dataset.command === "edit.replaceAll") {
      button.textContent = "All";
    } else if (button.dataset.command === "command.palette") {
      button.textContent = "Cmd";
    } else if (button.dataset.command === "file.toggleReadOnly") {
      button.textContent = "Read Only";
    } else if (button.dataset.command === "diff.open") {
      button.textContent = "Compare";
    } else if (button.dataset.command === "file.exportView") {
      button.textContent = "Export";
    } else if (button.dataset.command === "file.exportXlsx") {
      button.textContent = "XLSX";
    } else if (button.dataset.command === "edit.copySelectedFormat") {
      button.textContent = "Copy";
    } else if (button.dataset.command === "view.filterBuilder") {
      button.textContent = "Builder";
    } else if (button.dataset.command === "format.savePreset") {
      button.textContent = "Save Preset";
    } else if (button.dataset.command === "format.applyPreset") {
      button.textContent = "Apply";
    } else if (button.dataset.command === "format.saveRule") {
      button.textContent = "Save Rule";
    } else if (button.dataset.command === "format.applyQuotePolicy") {
      button.textContent = "Apply Quote";
    } else if (button.dataset.command === "sql.copyResultCsv") {
      button.textContent = "Copy CSV";
    } else if (button.dataset.command === "sql.exportResultCsv") {
      button.textContent = "Export CSV";
    } else if (button.dataset.command === "selection.cycleAggregate") {
      updateAggregateStatus(selectionAggregateStats());
    } else if (button.dataset.command !== "view.sql" && button.dataset.command !== "view.profile") {
      button.textContent = command.title.replace("Replace All", "All");
    }
  }
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
  applyFormatPresetButton.disabled = !hasDocument || !formatPresetSelect.value;
  formatRuleInput.disabled = !hasDocument;
  saveFormatRuleButton.disabled = !hasDocument;
  quoteModeSelect.disabled = !hasDocument;
  quoteCharSelect.disabled = !hasDocument;
  applyQuotePolicyButton.disabled = !hasDocument || readOnly;
}

function isCommandDisabled(command) {
  const hasDocument = Boolean(documentId);
  const needsSQLResult = command.id === "sql.copyResultCsv" || command.id === "sql.exportResultCsv";
  const needsAggregate =
    command.id === "selection.copyAggregate" ||
    command.id === "selection.cycleAggregate" ||
    command.id === "selection.insertAggregateRow";
  const writeBlocked =
    readOnly &&
    (
      command.id === "file.save" ||
      command.id === "edit.undo" ||
      command.id === "edit.redo" ||
      command.id === "edit.paste" ||
      command.id === "edit.replace" ||
      command.id === "edit.replaceAll" ||
      command.id === "row.insert" ||
      command.id === "row.delete" ||
      command.id === "row.hide" ||
      command.id === "row.moveUp" ||
      command.id === "row.moveDown" ||
      command.id === "column.insert" ||
      command.id === "column.delete" ||
      command.id === "column.hide" ||
      command.id === "column.moveLeft" ||
      command.id === "column.moveRight" ||
      command.id === "table.transpose" ||
      command.id === "table.dedupeRows" ||
      command.id === "table.dedupeColumns" ||
      command.id === "text.trim" ||
      command.id === "text.uppercase" ||
      command.id === "text.lowercase" ||
      command.id === "cell.swap" ||
      command.id === "cell.joinHorizontal" ||
      command.id === "cell.splitHorizontal" ||
      command.id === "cell.joinVertical" ||
      command.id === "cell.splitVertical" ||
      command.id === "cell.fillDown" ||
      command.id === "cell.fillRight" ||
      command.id === "cell.fillUp" ||
      command.id === "cell.fillLeft" ||
      command.id === "format.applyQuotePolicy"
    );
  return (
    command.implemented === false ||
    (command.requiresDocument && !hasDocument) ||
    (needsSQLResult && !sqlResultAvailable()) ||
    (needsAggregate && selectionAggregateStats().count <= 1) ||
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
  cellQuote.textContent = documentId ? quoteLabel(valueAt(selected.row, selected.col)) : "unknown";
  formatSummary.textContent = metadata ? documentFormatSummary() : "unknown";
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
    readOnlyStatus.textContent = "Editable";
    externalStatus.textContent = "Synced";
    return;
  }
  encodingStatus.textContent = encodingLabel();
  delimiterStatus.textContent = delimiterLabel(metadata.dialect);
  quoteStatus.textContent = quotePolicyLabel(metadata.dialect);
  headerStatus.textContent = metadata.dialect?.hasHeader ? "Header row" : "No header";
  headerStatus.disabled = !documentId;
  headerStatus.setAttribute("aria-pressed", String(Boolean(metadata.dialect?.hasHeader)));
  readOnlyStatus.textContent = readOnly ? "Read-only" : "Editable";
  externalStatus.textContent = externalFileChanged ? "Changed on disk" : "Synced";
  syncQuotePolicyControls();
}

function encodingLabel(info = metadata?.encoding) {
  const name = info?.name ?? "UTF-8";
  if (name.toLowerCase() === "utf-8-bom") return "UTF-8 BOM";
  if (name.toLowerCase() === "utf-8") return "UTF-8";
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

function lineEndingLabel(info = metadata?.lineEnding) {
  return (info?.name ?? "LF").toUpperCase();
}

function quoteModeValue(dialect = metadata?.dialect) {
  const raw = String(dialect?.quoteMode || "minimal").toLowerCase();
  return ["minimal", "always", "never"].includes(raw) ? raw : "minimal";
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

function quotePolicyLabel(dialect = metadata?.dialect) {
  return `Quote ${quoteModeValue(dialect)} ${quoteCharLabel(dialect)}`;
}

function syncQuotePolicyControls() {
  if (!quoteModeSelect || !quoteCharSelect) return;
  if (quotePolicyDraftDirty) {
    settingsQuoteSummary.textContent = `${quoteModeSelect.value} ${quoteCharSelect.value}`;
    return;
  }
  quoteModeSelect.value = quoteModeValue();
  quoteCharSelect.value = quoteCharValue();
  settingsQuoteSummary.textContent = metadata ? `${quoteModeValue()} ${quoteCharLabel()}` : "minimal double";
}

function documentFormatSummary() {
  return `${encodingLabel()} · ${delimiterLabel(metadata.dialect)} · ${lineEndingLabel()} · ${quotePolicyLabel(metadata.dialect)}`;
}

function currentFormatPreset() {
  if (!metadata?.dialect) return null;
  return {
    id: "",
    name: documentFormatSummary(),
    delimiter: delimiterChar(metadata.dialect),
    delimiterName: metadata.dialect.delimiterName || delimiterLabel(metadata.dialect).toLowerCase(),
    hasHeader: Boolean(metadata.dialect.hasHeader),
    encodingName: metadata.encoding?.name ?? "utf-8",
    lineEndingName: metadata.lineEnding?.name ?? "LF",
    quoteChar: quoteCharValue(metadata.dialect),
    quoteMode: quoteModeValue(metadata.dialect),
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
        quoteChar: preset.quoteChar || "double",
        quoteMode: preset.quoteMode || "minimal",
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
  applyFormatPresetButton.disabled = !documentId || !formatPresetSelect.value;
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

function openSettingsPanel() {
  setViewMode("settings");
}

function renderSettingsPanel() {
  if (!settingsPresetList) return;
  syncQuotePolicyControls();
  const presets = formatPresets();
  const query = settingsSearch.value.trim().toLowerCase();
  const visiblePresets = query
    ? presets.filter((preset) => formatPresetSearchText(preset).includes(query))
    : presets;
  settingsModifiedBadge.textContent = presets.length ? "Modified" : "Defaults";
  settingsFormatSummary.textContent = `${presets.length} presets, ${presets.reduce((sum, preset) => sum + (preset.rules?.length ?? 0), 0)} rules`;
  resetFormatPresetsButton.disabled = presets.length === 0;
  settingsPresetList.replaceChildren();
  document.body.dataset.settingsPresetCount = String(presets.length);
  document.body.dataset.settingsVisiblePresetCount = String(visiblePresets.length);
  if (visiblePresets.length === 0) {
    settingsPresetList.append(emptySettingsNode(presets.length ? "No matching presets" : "No format presets saved"));
    return;
  }
  for (const preset of visiblePresets) {
    settingsPresetList.append(formatPresetSettingsRow(preset));
  }
}

function formatPresetSearchText(preset) {
  return [
    preset.name,
    preset.id,
    preset.delimiterName,
    preset.encodingName,
    preset.quoteMode,
    preset.quoteChar,
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
  const meta = document.createElement("span");
  meta.textContent = `${formatPresetPreview(preset)} · ${preset.rules?.length ?? 0} rules`;
  const actions = document.createElement("div");
  actions.className = "settings-preset-actions";
  const applyButton = settingsActionButton("Apply");
  applyButton.disabled = !documentId;
  applyButton.addEventListener("click", async () => {
    formatPresetSelect.value = preset.id;
    await applyFormatPresetToCurrentPath(preset);
  });
  const deleteButton = settingsActionButton("Delete");
  deleteButton.addEventListener("click", () => deleteFormatPreset(preset.id));
  actions.append(applyButton, deleteButton);
  header.append(title, meta, actions);

  const rules = document.createElement("div");
  rules.className = "settings-rule-list";
  if ((preset.rules ?? []).length === 0) {
    const empty = document.createElement("span");
    empty.className = "settings-rule-empty";
    empty.textContent = "No automatic rules";
    rules.append(empty);
  } else {
    for (const rule of preset.rules) {
      const item = document.createElement("div");
      item.className = "settings-rule-row";
      const label = document.createElement("code");
      label.textContent = rule.pattern;
      const remove = settingsActionButton("Remove");
      remove.addEventListener("click", () => deleteFormatRule(preset.id, rule.pattern));
      item.append(label, remove);
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

  row.append(header, rules, addRule);
  return row;
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
  status.textContent = "Deleted format preset.";
  updateCommandState();
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
  status.textContent = `Removed format rule: ${pattern}`;
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
  status.textContent = `Added format rule: ${normalized.pattern}`;
}

function resetFormatPresets() {
  localStorage.removeItem(formatPresetsStorageKey);
  formatPresetSelect.value = "";
  document.body.dataset.formatPresets = "0";
  document.body.dataset.formatRulePreview = "reset";
  document.body.dataset.formatPresetPreview = "";
  renderFormatPresetSelect();
  renderSettingsPanel();
  updateCommandState();
  status.textContent = "Reset format presets.";
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
  const response = await fetch("/api/quote-policy", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, quoteChar, quoteMode }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Quote policy failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  quotePolicyDraftDirty = false;
  applyMutationPayload(payload, `Quote policy set to ${quoteMode} ${quoteChar}.`);
  if (metadata?.dialect) {
    metadata.dialect.quoteMode = quoteMode;
    metadata.dialect.quoteChar = quoteChar === "single" ? "'" : quoteChar === "none" ? "" : '"';
  }
  syncQuotePolicyControls();
  updateFormatStatus();
  updateInspector();
  status.textContent = `Quote policy set to ${quoteMode} ${quoteChar}.`;
  document.body.dataset.quotePolicy = `${quoteMode}:${quoteChar}`;
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
  return [
    preset.delimiter === "\t" ? "tab" : preset.delimiter,
    preset.hasHeader ? "header" : "no-header",
    (preset.encodingName || "utf-8").toLowerCase(),
    preset.quoteChar || "double",
    preset.quoteMode || "minimal",
  ].join(":");
}

function formatPresetName(preset) {
  const delimiter = preset.delimiterName || (preset.delimiter === "\t" ? "tab" : preset.delimiter);
  return `${delimiterLabel({ delimiter: preset.delimiter, delimiterName: delimiter })} · ${preset.hasHeader ? "Header" : "No header"} · ${(preset.encodingName || "UTF-8").toUpperCase()} · ${quoteModeValue(preset)} ${quoteCharLabel(preset)}`;
}

function formatPresetPreview(preset) {
  return `${preset.delimiter === "\t" ? "\\t" : preset.delimiter}|${preset.hasHeader ? "header" : "no-header"}|${preset.encodingName || "utf-8"}|${preset.quoteMode || "minimal"}|${preset.quoteChar || "double"}`;
}

function quoteLabel(value) {
  return shouldQuote(value) ? "will quote on save" : "no quote needed";
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
  cellValueInput.disabled = !documentId || readOnly || aggregate;
  applyCellValueButton.disabled = !documentId || readOnly || aggregate || !cellValueInputDirty;
  const url = normalizeURL(value);
  if (url) {
    openCellLink.href = url;
    openCellLink.setAttribute("aria-disabled", "false");
    openCellLink.classList.remove("disabled");
  } else {
    openCellLink.removeAttribute("href");
    openCellLink.setAttribute("aria-disabled", "true");
    openCellLink.classList.add("disabled");
  }
}

function selectedCellKey() {
  return documentId ? `${documentId}:${selected.row}:${selected.col}:${aggregateRowAt(selected.row) ? "aggregate" : "data"}` : "";
}

function normalizeURL(value) {
  const text = String(value ?? "").trim();
  if (!text) return "";
  if (/^https?:\/\/[^\s]+$/i.test(text)) return text;
  if (/^mailto:[^\s@]+@[^\s@]+\.[^\s@]+$/i.test(text)) return text;
  if (/^file:\/\/[^\s]+$/i.test(text)) return text;
  if (/^[a-z0-9.-]+\.[a-z]{2,}(?:\/[^\s]*)?$/i.test(text)) return `https://${text}`;
  return "";
}

function cellAtPoint(x, y) {
  if (x < rowNumberWidth || y < headerHeight) return null;
  const row = rowAt(y);
  const col = columnAt(x);
  if (row < 0 || col < 0) return null;
  return { row, col };
}

function urlAtPoint(x, y) {
  const cell = cellAtPoint(x, y);
  if (!cell) return "";
  return normalizeURL(valueAt(cell.row, cell.col));
}

function clearURLHover() {
  if (!hoverURL && canvas.style.cursor !== "pointer" && !canvas.title) return;
  hoverURL = "";
  canvas.style.cursor = "";
  canvas.removeAttribute("title");
  document.body.dataset.hoverUrl = "";
}

function updateURLHover(event) {
  if (draggingSelection || editing || !documentId) {
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
  canvas.style.cursor = "pointer";
  canvas.title = url;
  document.body.dataset.hoverUrl = url;
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

function selectionCellCount() {
  if (!documentId || columns.length === 0 || totalRows === 0) return 0;
  const range = normalizedSelection();
  return (range.endRow - range.startRow + 1) * (range.endCol - range.startCol + 1);
}

function selectionSummary(range = normalizedSelection()) {
  if (!documentId || columns.length === 0 || totalRows === 0) return "";
  const base = `${range.startRow}:${range.startCol}-${range.endRow}:${range.endCol}`;
  return includeHeadersInSelection ? `${base}:headers` : base;
}

function selectionAggregateStats() {
  const count = selectionCellCount();
  const stats = {
    count,
    numericCount: 0,
    limited: false,
    sum: 0,
    avg: 0,
    min: 0,
    max: 0,
  };
  if (count <= 1) return stats;
  const range = normalizedSelection();
  if (range.endRow >= sourceRowCount) {
    stats.limited = true;
    return stats;
  }
  if (count > 5000) {
    stats.limited = true;
    return stats;
  }
  let minValue = Number.POSITIVE_INFINITY;
  let maxValue = Number.NEGATIVE_INFINITY;
  for (let row = range.startRow; row <= range.endRow; row++) {
    for (let col = range.startCol; col <= range.endCol; col++) {
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
  if (stats.numericCount) {
    stats.avg = stats.sum / stats.numericCount;
    stats.min = minValue;
    stats.max = maxValue;
  }
  return stats;
}

function selectionAggregateSummary(stats = selectionAggregateStats()) {
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

function aggregateMetricOptions(stats = selectionAggregateStats()) {
  if (stats.count <= 1) return [];
  if (!stats.numericCount) return ["count"];
  return ["sum", "avg", "min", "max", "count"];
}

function normalizeAggregateMetric(stats = selectionAggregateStats()) {
  const options = aggregateMetricOptions(stats);
  if (!options.length) return "";
  if (!options.includes(activeAggregateMetric)) activeAggregateMetric = options[0];
  return activeAggregateMetric;
}

function activeAggregate(stats = selectionAggregateStats()) {
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

function updateAggregateStatus(stats = selectionAggregateStats()) {
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

function insertAggregateRow() {
  if (!documentId) return;
  const range = normalizedSelection();
  if (range.endRow >= sourceRowCount) {
    status.textContent = "Select data rows before adding an aggregate row.";
    return;
  }
  const stats = selectionAggregateStats();
  const metric = normalizeAggregateMetric(stats);
  if (!metric) {
    status.textContent = "Select numeric cells to add an aggregate row.";
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
  shortcutsSearch.value = "";
  shortcutsPanel.hidden = false;
  renderKeyboardShortcuts();
  shortcutsSearch.focus();
}

function closeKeyboardShortcuts() {
  shortcutEditorOpen = false;
  recordingShortcutID = "";
  shortcutsPanel.hidden = true;
  shortcutsList.replaceChildren();
  canvas.focus();
}

function renderKeyboardShortcuts() {
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
    title.textContent = command.title;
    const description = document.createElement("span");
    description.textContent = command.description;
    meta.append(title, description);

    const shortcut = document.createElement("kbd");
    shortcut.textContent = command.id === recordingShortcutID
      ? "Press shortcut"
      : effectiveShortcut(command);

    const actions = document.createElement("div");
    actions.className = "shortcut-actions";
    const recordButton = document.createElement("button");
    recordButton.type = "button";
    recordButton.textContent = command.id === recordingShortcutID ? "Cancel" : "Record";
    recordButton.addEventListener("click", () => {
      recordingShortcutID = command.id === recordingShortcutID ? "" : command.id;
      renderKeyboardShortcuts();
    });
    const resetButton = document.createElement("button");
    resetButton.type = "button";
    resetButton.textContent = "Reset";
    resetButton.disabled = !customShortcuts.has(command.id);
    resetButton.addEventListener("click", () => resetCustomShortcut(command.id));
    actions.append(recordButton, resetButton);

    row.append(meta, shortcut, actions);
    fragment.append(row);
  }
  shortcutsList.replaceChildren(fragment);
}

function effectiveShortcut(command) {
  return customShortcuts.get(command.id) || command.shortcut || "";
}

function defaultShortcutActive(commandID) {
  return !customShortcuts.has(commandID);
}

function loadCustomShortcuts() {
  try {
    const parsed = JSON.parse(localStorage.getItem(customShortcutsStorageKey) || "{}");
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) return new Map();
    return new Map(
      Object.entries(parsed)
        .filter(([id, shortcut]) => typeof id === "string" && typeof shortcut === "string" && shortcut.trim())
        .map(([id, shortcut]) => [id, shortcut]),
    );
  } catch {
    return new Map();
  }
}

function saveCustomShortcuts() {
  const value = Object.fromEntries(customShortcuts);
  localStorage.setItem(customShortcutsStorageKey, JSON.stringify(value));
  document.body.dataset.shortcutPreview = JSON.stringify(value);
  applyCommandMetadata();
  renderCommandPalette();
  if (shortcutEditorOpen) renderKeyboardShortcuts();
}

function resetCustomShortcut(commandID) {
  const command = commandsById.get(commandID);
  if (!command) return;
  customShortcuts.delete(commandID);
  saveCustomShortcuts();
  status.textContent = `Reset shortcut for ${command.title}`;
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
  const conflict = commandForShortcut(shortcut, command.id);
  if (conflict) {
    status.textContent = `${shortcut} is already assigned to ${conflict.title}`;
    return true;
  }
  if (shortcut === command.shortcut) {
    customShortcuts.delete(command.id);
  } else {
    customShortcuts.set(command.id, shortcut);
  }
  recordingShortcutID = "";
  saveCustomShortcuts();
  status.textContent = `Shortcut for ${command.title}: ${shortcut}`;
  return true;
}

function shortcutAllowed(shortcut) {
  return shortcut.includes("Cmd+") || shortcut.includes("Alt+") || /^F\d+$/.test(shortcut);
}

function commandForShortcut(shortcut, exceptID = "") {
  for (const command of commandsById.values()) {
    if (command.id === exceptID) continue;
    if (effectiveShortcut(command) === shortcut) return command;
  }
  return null;
}

function customCommandForShortcut(shortcut) {
  for (const [commandID, customShortcut] of customShortcuts.entries()) {
    if (customShortcut === shortcut) return commandsById.get(commandID) ?? null;
  }
  return null;
}

function runCustomShortcut(event) {
  if (paletteOpen || shortcutEditorOpen) return false;
  const shortcut = shortcutFromEvent(event);
  if (!shortcut) return false;
  const command = customCommandForShortcut(shortcut);
  if (!command) return false;
  event.preventDefault();
  if (isCommandDisabled(command)) {
    status.textContent = `${command.title} is unavailable.`;
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
  const fragment = document.createDocumentFragment();
  for (const column of columns) {
    const option = document.createElement("option");
    option.value = String(column.sourceCol);
    option.textContent = column.name;
    fragment.append(option);
  }
  conditionColumn.replaceChildren(fragment);
  if ([...conditionColumn.options].some((option) => option.value === previous)) {
    conditionColumn.value = previous;
  }
}

function updateConditionValueState() {
  const noValue = conditionOperator.value === "blank" || conditionOperator.value === "notBlank";
  conditionValue.disabled = noValue;
  if (noValue) conditionValue.value = "";
}

async function applyConditionBuilder() {
  if (!documentId) return;
  const sourceCol = Number(conditionColumn.value);
  const expression = conditionExpression(sourceCol, conditionOperator.value, conditionValue.value);
  if (!expression) return;
  document.body.dataset.conditionPreview = expression;
  commandInput.value = `where ${expression}`;
  closeConditionBuilder();
  const visibleCol = visibleColForSource(sourceCol);
  await applyFilterRequest(
    { documentId, mode: "sql", expression },
    "Condition filter",
    visibleCol >= 0 ? visibleCol : selected.col,
  );
}

function conditionExpression(sourceCol, operator, value) {
  if (!Number.isInteger(sourceCol) || sourceCol < 0) {
    status.textContent = "Choose a column to filter.";
    return "";
  }
  const identifier = sqlColumnIdentifier(sourceCol);
  const needsValue = !["blank", "notBlank"].includes(operator);
  if (needsValue && value === "") {
    status.textContent = "Enter a filter value.";
    conditionValue.focus();
    return "";
  }
  switch (operator) {
    case "contains":
      return `${identifier} LIKE ${sqlStringLiteral(`%${value}%`)}`;
    case "equals":
      return `${identifier} = ${sqlStringLiteral(value)}`;
    case "notEquals":
      return `${identifier} <> ${sqlStringLiteral(value)}`;
    case "startsWith":
      return `${identifier} LIKE ${sqlStringLiteral(`${value}%`)}`;
    case "endsWith":
      return `${identifier} LIKE ${sqlStringLiteral(`%${value}`)}`;
    case "blank":
      return `${identifier} = ''`;
    case "notBlank":
      return `${identifier} <> ''`;
    case "gt":
      return `${identifier} > ${sqlComparisonLiteral(value)}`;
    case "gte":
      return `${identifier} >= ${sqlComparisonLiteral(value)}`;
    case "lt":
      return `${identifier} < ${sqlComparisonLiteral(value)}`;
    case "lte":
      return `${identifier} <= ${sqlComparisonLiteral(value)}`;
    default:
      status.textContent = "Choose a filter operator.";
      return "";
  }
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
    title.textContent = command.title;
    const description = document.createElement("span");
    description.textContent = command.description;
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
  return command.title || command.id || "";
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
    status.textContent = command ? `${command.title} is unavailable.` : "Command unavailable.";
    return;
  }
  closeCommandPalette();
  await executeCommand(commandID);
}

async function executeCommand(commandID) {
  rememberCommand(commandID);
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
    case "workbench.toggleTheme":
      cycleTheme();
      break;
    case "file.open":
      await openPathInputOrPicker();
      break;
    case "file.save":
      await saveCurrentDocument();
      break;
    case "file.exportView":
      await exportCurrentView();
      break;
    case "file.exportXlsx":
      await exportCurrentView("xlsx");
      break;
    case "sql.copyResultCsv":
      await copySQLResultCSV();
      break;
    case "sql.exportResultCsv":
      await exportSQLResultCSV();
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
    case "format.saveRule":
      saveCurrentFormatRule();
      break;
    case "format.applyQuotePolicy":
      await applyQuotePolicy();
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
    case "selection.toggleHeaders":
      toggleHeaderInSelection();
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
      commandInput.focus();
      if (commandInput.value.trim()) await findCurrentQuery();
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
    case "row.hide":
      await hideSelectedRow();
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
    case "column.hide":
      hideSelectedColumn();
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
      await transformSelectedText("trim", "Trim Text");
      break;
    case "text.uppercase":
      await transformSelectedText("upper", "Uppercase Text");
      break;
    case "text.lowercase":
      await transformSelectedText("lower", "Lowercase Text");
      break;
    case "cell.swap":
      await swapSelectedCells();
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
    case "view.sort":
      await sortSelectedColumn();
      break;
    case "view.clear":
      await clearView();
      break;
    case "view.sql":
      setViewMode("sql");
      break;
    case "view.profile":
      await openProfilePanel();
      break;
    case "diff.open":
      openDiffPanel();
      break;
    default:
      status.textContent = `Command ${commandID} has no UI handler.`;
  }
}

function drawHeaders(width) {
  const headerNameY = Math.max(14, Math.round(headerHeight * 0.34));
  const headerTypeY = Math.min(headerHeight - 10, Math.round(headerHeight * 0.72));
  const target = ctx;
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

  ctx.font =
    "13px ui-monospace, SFMono-Regular, Menlo, Consolas, Liberation Mono, monospace";
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
        ctx.fillText(value, x + 8, y + rowHeight / 2);
        if (shouldQuote(value)) {
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

function shouldQuote(value) {
  const mode = quoteModeValue();
  if (mode !== "minimal") return mode === "always";
  return requiresQuote(value);
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

function drawSelection() {
  if (columns.length === 0) return;
  const range = normalizedSelection();
  ctx.save();
  ctx.beginPath();
  ctx.rect(0, headerHeight, canvas.clientWidth, Math.max(0, canvas.clientHeight - headerHeight));
  ctx.clip();

  const firstVisibleRow = Math.max(range.startRow, Math.floor(scrollY / rowHeight));
  const lastVisibleRow = Math.min(
    range.endRow,
    Math.ceil((scrollY + canvas.clientHeight - headerHeight) / rowHeight),
  );
  const visibleColumns = visibleColumnIndexes(range);

  ctx.fillStyle = themeColors.selectionBg;
  for (let row = firstVisibleRow; row <= lastVisibleRow; row++) {
    for (const col of visibleColumns) {
      const rect = cellRect(row, col);
      if (!rect) continue;
      ctx.fillRect(rect.x + 1, rect.y + 1, rect.width - 2, rect.height - 2);
    }
  }

  const rect = cellRect(selected.row, selected.col);
  if (!rect) {
    ctx.restore();
    return;
  }
  ctx.strokeStyle = themeColors.accent;
  ctx.lineWidth = 2;
  ctx.strokeRect(rect.x + 1, rect.y + 1, rect.width - 2, rect.height - 2);
  ctx.lineWidth = 1;
  ctx.restore();
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

function clampColumn(col) {
  return Math.max(0, Math.min(columns.length - 1, col));
}

function isHeaderMenuTarget(x, col) {
  if (!columns[col]) return false;
  const left = columnLeft(col);
  return x >= left + columns[col].width - 28;
}

function columnAtOrNearest(x) {
  if (columns.length === 0) return -1;
  const col = columnAt(x);
  if (col >= 0) return col;
  return x < rowNumberWidth ? 0 : columns.length - 1;
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

function openColumnMenu(col, clientX, clientY) {
  if (!documentId || !columns[col]) return;
  columnMenuColumn = col;
  selected = { row: Math.max(0, Math.min(totalRows - 1, selected.row)), col };
  selectionAnchor = { ...selected };
  columnMenuTitle.textContent = columns[col].name || `Column ${col + 1}`;
  const selectedValue = valueAt(selected.row, col);
  columnFilterValue.textContent =
    selectedValue === ""
      ? "Filter by selected value"
      : `Filter: ${truncateMenuValue(selectedValue)}`;
  columnFilterValue.disabled = selectedValue === "";
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
  node.textContent = "Loading values...";
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
    renderColumnMenuValues(payload.values ?? [], payload.totalDistinct ?? 0);
  } catch (error) {
    columnMenuValues.replaceChildren();
    const node = document.createElement("div");
    node.className = "column-values-empty";
    node.textContent = `Values failed: ${error.message}`;
    columnMenuValues.append(node);
    columnApplyValues.disabled = true;
  }
}

function renderColumnMenuValues(values, totalDistinct) {
  columnMenuValues.replaceChildren();
  if (!values.length) {
    const node = document.createElement("div");
    node.className = "column-values-empty";
    node.textContent = "No values";
    columnMenuValues.append(node);
    columnApplyValues.disabled = true;
    return;
  }
  const heading = document.createElement("span");
  heading.className = "column-values-heading";
  heading.textContent = totalDistinct > values.length
    ? `Values (${values.length} of ${formatCount(totalDistinct)})`
    : `Values (${formatCount(totalDistinct)})`;
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
    text.textContent = `${truncateMenuValue(item.blank ? "(blank)" : item.value)} (${formatCount(item.count)})`;
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
          "13px ui-monospace, SFMono-Regular, Menlo, Consolas, Liberation Mono, monospace",
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
  cellValueInput.value = valueAt(selected.row, selected.col);
  cellValueInputKey = selectedCellKey();
  cellValueInputDirty = false;
  applyCellValueButton.disabled = true;
}

function applyCellValue() {
  if (!cellValueInputDirty) return;
  setActiveCellValue(cellValueInput.value);
}

function submitCellEdits(edits, windowRow, message) {
  if (!edits.length) return;
  dirty = true;
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
      rows = new Map();
      rowSource = new Map();
      windowStart = 0;
      windowEnd = 0;
      await maybeLoadWindow(windowRow);
      applyEditsToLoadedRows(edits);
      draw();
      if (message) status.textContent = message;
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
  status.textContent = `Opening ${path}...`;
  let preset = options.preset ?? null;
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
  if (preset) {
    body.delimiter = preset.delimiter;
    body.delimiterName = preset.delimiterName;
    body.hasHeader = Boolean(preset.hasHeader);
    body.encodingName = preset.encodingName;
    body.quoteChar = preset.quoteChar || "double";
    body.quoteMode = preset.quoteMode || "minimal";
  }
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
  loadDocument(payload.document);
  updateFileList(fileNameFromPath(path), metadata.path, totalRows, { remember: true });
  updateFormatRulePlaceholder(path);
  draw();
  startFileStatusChecks();
  await checkFileStatus(false);
  if (matchedRule) {
    status.textContent = `Applied format rule: ${matchedRule.pattern} -> ${preset.name || formatPresetName(preset)}`;
  } else if (preset) {
    status.textContent = `Applied format preset: ${preset.name || formatPresetName(preset)}`;
  }
  canvas.focus();
}

async function uploadAndOpen(file) {
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
  loadDocument(payload.document);
  pathInput.value = file.name;
  updateFormatRulePlaceholder(file.name);
  updateFileList(file.name, "temporary upload", totalRows);
  draw();
  startFileStatusChecks();
  await checkFileStatus(false);
  canvas.focus();
}

function loadDocument(document) {
  windowGeneration++;
  rows = new Map();
  rowSource = new Map();
  profileData = null;
  lastSQLResult = null;
  profileColumns.replaceChildren(emptyProfileNode("Run Profile to inspect columns."));
  profileIssues.replaceChildren(emptyProfileNode("No validation issues loaded."));
  profileSummary.replaceChildren();
  sqlResult.replaceChildren(emptySQLNode("Run SQL to inspect results."));
  delete window.document.body.dataset.sqlResultPreview;
  syncSQLResultActions();
  dataStartRow = document.metadata.dialect.hasHeader ? 1 : 0;
  viewMapped = false;
  sortDescending = false;
  sortKeys = [];
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
  hiddenSourceColumns = new Set();
  applyHiddenColumns();
  selected = { row: 0, col: 0 };
  selectionAnchor = { row: 0, col: 0 };
  scrollX = 0;
  scrollY = 0;
  wheelRemainderY = 0;
  setHeaderInSelection(false);
  clearURLHover();
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
  updateFileList(session.displayName || fileNameFromPath(detail), detail, totalRows);
  draw();
  status.textContent = message;
}

function selectedSourceRow() {
  return rowSource.get(selected.row) ?? selected.row + dataStartRow;
}

function selectedMatrix() {
  const range = normalizedSelection();
  const matrix = [];
  for (let row = range.startRow; row <= range.endRow; row++) {
    const values = [];
    for (let col = range.startCol; col <= range.endCol; col++) {
      values.push(valueAt(row, col));
    }
    matrix.push(values);
  }
  return matrix;
}

function selectedHeaders() {
  const range = normalizedSelection();
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
  const range = normalizedSelection();
  if (range.endRow >= sourceRowCount) {
    status.textContent = "Aggregate rows are view-only.";
    return null;
  }
  const cells = [];
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

async function exportCurrentView(format = "csv") {
  if (!documentId) return;
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

async function copySQLResultCSV() {
  const selectedMatrix = sqlSelectedResultMatrix(lastSQLResult);
  const matrix = selectedMatrix ?? sqlResultMatrix(lastSQLResult);
  if (!matrix.length) {
    status.textContent = "Run a SQL query before copying results.";
    return;
  }
  const text = formatMatrix(matrix, "csv");
  window.__opencsvLastSQLResultText = text;
  document.body.dataset.sqlResultPreview = text.slice(0, 500);
  try {
    await writeClipboardText(text);
    if (selectedMatrix) {
      status.textContent = `Copied SQL selection: ${formatCount(matrix.length)} rows, ${formatCount(matrix[0]?.length ?? 0)} columns`;
    } else {
      status.textContent = `Copied SQL result: ${formatCount(Math.max(0, matrix.length - 1))} rows, ${formatCount(matrix[0]?.length ?? 0)} columns`;
    }
  } catch (error) {
    status.textContent = `SQL result copy failed: ${error.message}`;
  }
}

async function exportSQLResultCSV() {
  const matrix = sqlResultMatrix(lastSQLResult);
  if (!matrix.length) {
    status.textContent = "Run a SQL query before exporting results.";
    return;
  }
  const text = formatMatrix(matrix, "csv");
  const filename = sqlResultFilename();
  window.__opencsvLastSQLResultText = text;
  document.body.dataset.sqlResultPreview = text.slice(0, 500);
  triggerDownload(text, filename, "text/csv;charset=utf-8");
  status.textContent = `Exported SQL result: ${formatCount(Math.max(0, matrix.length - 1))} rows, ${formatCount(matrix[0]?.length ?? 0)} columns to ${filename}`;
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

function sqlResultColumns(result) {
  if (!result) return [];
  const columns = Array.isArray(result.columns) ? result.columns : [];
  if (columns.length) return columns.map((column, index) => column || `Column ${index + 1}`);
  const width = Math.max(0, ...(result.rows ?? []).map((row) => row?.length ?? 0));
  return Array.from({ length: width }, (_, index) => `Column ${index + 1}`);
}

function sqlResultFilename() {
  const name = activeFileEntry?.name || fileNameFromPath(metadata?.path ?? "") || "opencsv";
  const base = name.replace(/\.[^.]+$/, "") || "opencsv";
  return `${base}.sql-result.csv`;
}

function syncSQLResultActions() {
  const available = sqlResultAvailable();
  copySqlResultButton.disabled = !available;
  exportSqlResultButton.disabled = !available;
  updateCommandState();
}

function clearSQLResultSelection() {
  sqlResultSelectionAnchor = null;
  sqlResultSelectionRange = null;
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
  syncSQLResultSelection();
  const range = normalizedSQLResultSelection();
  const rowsSelected = range ? range.endRow - range.startRow + 1 : 0;
  const colsSelected = range ? range.endCol - range.startCol + 1 : 0;
  status.textContent = `Selected SQL result: ${formatCount(rowsSelected)} rows, ${formatCount(colsSelected)} columns`;
}

function sqlResultCellFromEvent(event) {
  const cell = event.target?.closest?.("td[data-sql-row]");
  if (!cell || !sqlResult.contains(cell)) return null;
  return {
    row: Number(cell.dataset.sqlRow),
    col: Number(cell.dataset.sqlCol),
  };
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
    const text = await navigator.clipboard.readText();
    pasteText(text);
  } catch (error) {
    status.textContent = `Paste failed: ${error.message}`;
  }
}

function pasteText(text) {
  if (!ensureEditable("Paste")) return;
  const matrix = parseTSV(text);
  if (!matrix.length) {
    status.textContent = "Clipboard is empty.";
    return;
  }
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
  const normalized = text.replaceAll("\r\n", "\n").replaceAll("\r", "\n").replace(/\n$/, "");
  if (!normalized) return [];
  return normalized.split("\n").map((line) => line.split("\t"));
}

function updateFileList(name, detail, rowCount, options = {}) {
  activeFileEntry = { name, detail, rowCount };
  if (options.remember) rememberRecentFile(activeFileEntry);
  renderFileList();
}

function renderFileList() {
  fileList.replaceChildren();
  if (!activeFileEntry) {
    renderStartActions();
  } else {
    fileList.append(fileButton(activeFileEntry, {
      active: true,
      onClick: () => {
        selected = { row: 0, col: 0 };
        selectionAnchor = { row: 0, col: 0 };
        scrollX = 0;
        scrollY = 0;
        wheelRemainderY = 0;
        canvas.focus();
        draw();
      },
    }));
  }

  const recent = recentFiles().filter((item) => item.detail !== activeFileEntry?.detail);
  if (recent.length > 0) {
    const label = document.createElement("div");
    label.className = "file-section-label";
    label.textContent = "Recent";
    fileList.append(label);
    for (const item of recent) {
      fileList.append(fileButton(item, { onClick: () => openRecentFile(item) }));
    }
  }
}

function renderStartActions() {
  const openPathButton = document.createElement("button");
  openPathButton.className = "file-button file-action";
  openPathButton.type = "button";
  openPathButton.append(fileText("Open CSV/TSV", "Type a path in the toolbar"));
  openPathButton.addEventListener("click", () => {
    pathInput.focus();
    pathInput.select();
  });

  const chooseButton = document.createElement("button");
  chooseButton.className = "file-button file-action";
  chooseButton.type = "button";
  chooseButton.append(fileText("Choose File", "Open with browser picker"));
  chooseButton.addEventListener("click", () => fileInput.click());

  fileList.append(openPathButton, chooseButton);
}

function fileButton(entry, options = {}) {
  const button = document.createElement("button");
  button.className = options.active ? "file-button active" : "file-button";
  button.type = "button";
  button.title = entry.detail;
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
  return `${rowsText} · ${entry.detail}`;
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
  const columnCount =
    document.metadata.dialect.estimatedColumns ||
    Math.max(...document.preview.map((record) => record.length), 1);
  return Array.from({ length: columnCount }, (_, index) => ({
    name: hasHeader && first[index] ? first[index] : `Column ${index + 1}`,
    sourceCol: index,
    width: Math.max(minColumnWidth, Math.min(240, String(first[index] ?? "").length * 9 + 64)),
    type: index === 0 ? "mixed" : "string",
  }));
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

async function findNextMatch(direction) {
  if (!documentId) return;
  await runFind(direction, false);
}

async function runFind(direction, forceRefresh) {
  if (!documentId) return;
  const search = parseSearchQuery(commandInput.value);
  const query = search.query;
  if (!query) {
    status.textContent = "Type text to find.";
    findState = emptyFindState();
    updateFindMatchCount();
    return;
  }
  const key = findKey(search);
  if (forceRefresh || findState.key !== key || findState.matches.length === 0) {
    const ok = await refreshFindState(search, key);
    if (!ok) return;
  }
  if (findState.matches.length === 0) {
    status.textContent = `No matches for ${query}`;
    updateFindMatchCount();
    return;
  }
  if (direction !== 0) {
    findState.index = (findState.index + direction + findState.matches.length) % findState.matches.length;
  }
  await selectFindMatch(findState.matches[findState.index], search.query);
}

async function refreshFindState(search, key) {
  const cells = search.inSelection ? selectedCellsForSearch() : [];
  if (cells == null) return false;
  const response = await fetch("/api/find", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      query: search.query,
      regex: search.regex,
      caseSensitive: search.caseSensitive,
      wholeCell: search.wholeCell,
      cells,
      limit: 5000,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Find failed: ${payload.error ?? "unknown error"}`;
    return false;
  }
  findState = {
    key,
    criteriaKey: findCriteriaKey(search),
    matches: payload.matches ?? [],
    total: payload.total ?? payload.matches?.length ?? 0,
    truncated: Boolean(payload.truncated),
    index: 0,
    scopeCells: cells,
  };
  updateFindMatchCount();
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
  const ordinal = findState.index + 1;
  const total = findState.total;
  const loadedNote = findState.truncated ? ` (showing first ${findState.matches.length.toLocaleString()})` : "";
  status.textContent = `Found ${ordinal.toLocaleString()} of ${total.toLocaleString()} for ${query} at row ${displayRow + 1}, column ${match.col + 1}${loadedNote}`;
}

async function replaceCurrentQuery(all) {
  if (!documentId) return;
  if (!ensureEditable("Replace")) return;
  const search = parseSearchQuery(commandInput.value);
  const query = search.query;
  if (!query) {
    status.textContent = "Type text to replace.";
    return;
  }
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
      all,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Replace failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  if (!payload.count) {
    status.textContent = `No replacements for ${query}`;
    return;
  }
  dirty = payload.dirty ?? true;
  profileData = null;
  invalidateFindState();
  const first = payload.changes?.[0];
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
    ? `Replaced ${payload.count.toLocaleString()} cells`
    : "Replaced 1 cell";
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
    status.textContent = `Selection has ${formatCount(count)} cells; search-in-selection limit is ${formatCount(maxCells)}.`;
    return null;
  }
  const range = normalizedSelection();
  const cells = [];
  for (let row = range.startRow; row <= range.endRow; row++) {
    const sourceRow = rowSource.get(row);
    if (viewMapped && sourceRow == null) {
      status.textContent = "Load the selected filtered rows before searching within selection.";
      return null;
    }
    const recordRow = sourceRow ?? row + dataStartRow;
    for (let col = range.startCol; col <= range.endCol; col++) {
      cells.push({ row: recordRow, col: sourceColFor(col) });
    }
  }
  return cells;
}

function selectedSearchScopeKey() {
  const range = normalizedSelection();
  const startSourceRow = rowSource.get(range.startRow) ?? range.startRow + dataStartRow;
  const endSourceRow = rowSource.get(range.endRow) ?? range.endRow + dataStartRow;
  return `${Math.min(startSourceRow, endSourceRow)}:${Math.max(startSourceRow, endSourceRow)}:${sourceColFor(range.startCol)}:${sourceColFor(range.endCol)}`;
}

function emptyFindState() {
  return { key: "", criteriaKey: "", matches: [], total: 0, truncated: false, index: 0, scopeCells: [] };
}

function invalidateFindState() {
  findState = emptyFindState();
  updateFindMatchCount();
}

function updateFindMatchCount() {
  const current = findState.matches.length ? findState.index + 1 : 0;
  const total = findState.total || 0;
  findMatchCount.textContent = `${current.toLocaleString()}/${total.toLocaleString()}`;
  const enabled = Boolean(documentId && total > 0);
  findPrevButton.disabled = !enabled;
  findNextButton.disabled = !enabled;
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
  );
}

async function applyFilterRequest(body, label, preserveCol = selected.col) {
  const response = await fetch("/api/filter", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  viewMapped = true;
  sortDescending = false;
  sortKeys = [];
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
  status.textContent = `${label} matched ${sourceRowCount.toLocaleString()} rows`;
}

function looksLikeSQLFilter(query) {
  const value = query.trim();
  return (
    /^where\s+/i.test(value) ||
    /\b(and|or|not|like|is|cast|try_cast|null|date|time|timestamp)\b/i.test(value) ||
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

async function sortColumn(column, descending, append, modeOverride = "") {
  if (!documentId || columns.length === 0 || !columns[column]) return;
  const nextKeys = nextSortKeys(column, descending, append, modeOverride);
  const primaryKey = nextKeys[0];
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
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Sort failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  viewMapped = true;
  sortKeys = nextKeys;
  sortDescending = descending;
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
  status.textContent = `Sorted ${formatSortKeys(sortKeys)}`;
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

async function filterColumnBySelectedValue() {
  const column = columnMenuColumn;
  closeColumnMenu();
  if (!documentId || !columns[column]) return;
  const value = valueAt(selected.row, column);
  if (value === "") return;
  commandInput.value = value;
  await applyFilterRequest(
    { documentId, query: value, column: sourceColFor(column) },
    `Filtered ${columns[column].name}`,
    column,
  );
}

async function filterColumnBlank(showBlank) {
  const column = columnMenuColumn;
  closeColumnMenu();
  if (!documentId || !columns[column]) return;
  const expression = `c${sourceColFor(column) + 1} ${showBlank ? "=" : "<>"} ''`;
  commandInput.value = `where ${expression}`;
  await applyFilterRequest(
    { documentId, mode: "sql", expression },
    showBlank ? `Blank ${columns[column].name}` : `Non-blank ${columns[column].name}`,
    column,
  );
}

async function applyColumnValueChecklist() {
  const column = columnMenuColumn;
  if (!documentId || !columns[column]) return;
  const checked = Array.from(columnMenuValues.querySelectorAll("input:checked"));
  if (!checked.length) {
    status.textContent = "Choose at least one value to filter.";
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
    status.textContent = `Value filter failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  viewMapped = true;
  sortDescending = false;
  sortKeys = [];
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
  status.textContent = `Value filter ${columns[column].name} matched ${sourceRowCount.toLocaleString()} rows`;
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
    status.textContent = `Clear failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  viewMapped = false;
  sortDescending = false;
  sortKeys = [];
  hiddenSourceColumns = new Set();
  applyHiddenColumns();
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
  status.textContent = "Filter/sort cleared";
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

async function insertSelectedColumn() {
  await mutateStructure("/api/column/insert", { col: sourceColFor(selected.col) }, "Inserted column");
}

async function deleteSelectedColumn() {
  await mutateStructure("/api/column/delete", { col: sourceColFor(selected.col) }, "Deleted column");
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
  const nextCol = Math.max(0, Math.min(columns.length - 1, selected.col));
  selected = { row: selected.row, col: nextCol };
  selectionAnchor = { ...selected };
  closeColumnMenu();
  draw();
  updateInspector();
  status.textContent = `Column hidden: ${columnNameForSource(sourceCol)}`;
}

async function moveSelectedColumn(delta) {
  if (!documentId) return;
  if (!ensureEditable("Move Column")) return;
  const previousCol = selected.col;
  const label = delta < 0 ? "Moved column left" : "Moved column right";
  status.textContent = `${label}...`;
  const response = await fetch("/api/column/move", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, col: sourceColFor(selected.col), delta }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const changed = payload.changed ?? 0;
  applyMutationPayload(payload, changed > 0 ? label : "Column is already at the edge");
  if (changed > 0) {
    selected = { row: selected.row, col: Math.max(0, Math.min(columns.length - 1, previousCol + delta)) };
    selectionAnchor = { ...selected };
    draw();
    status.textContent = label;
  }
}

async function transposeCurrentDocument() {
  await mutateStructure("/api/transpose", { maxCells: 1_000_000 }, "Transposed table");
}

async function toggleHeaderRow() {
  if (!documentId || !metadata) return;
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
  status.textContent = `${label}...`;
  const response = await fetch("/api/text-transform", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, mode, cells }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${label} failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const changed = payload.changed ?? 0;
  applyMutationPayload(
    payload,
    changed > 0
      ? `${label}: ${formatCount(changed)} cells changed`
      : `${label}: no cells changed`,
  );
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
  await runCellTool("/api/cells/join", {
    cells,
    separator: " ",
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
  await runCellTool("/api/cells/split", {
    cells: [{ row: selectedSourceRow(), col: sourceColFor(selected.col) }],
    separator: " ",
    direction,
    activeRow: selected.row,
    activeCol: selected.col,
    label: vertical ? "Split cell vertically" : "Split cell",
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

async function runCellTool(endpoint, options) {
  status.textContent = `${options.label}...`;
  const response = await fetch(endpoint, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      cells: options.cells,
      separator: options.separator,
      direction: options.direction,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `${options.label} failed: ${payload.error ?? "unknown error"}`;
    return;
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
  if (!documentId) {
    status.textContent = "Open a CSV through opencsv serve before saving.";
    return;
  }
  if (!ensureEditable("Save")) return;
  status.textContent = "Saving...";
  const response = await fetch("/api/save", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, backup: true }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Save failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  dirty = false;
  status.textContent = payload.backupPath
    ? `Saved ${payload.path}; backup ${payload.backupPath}`
    : `Saved ${payload.path}`;
  externalFileChanged = false;
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
    externalFileChanged = Boolean(payload.changedExternally);
    if (externalFileChanged && autoRefresh && !dirty && !readOnly && !viewMapped) {
      await refreshCurrentDocument();
      return;
    }
    updateFormatStatus();
    updateCommandState();
  } catch (error) {
    externalStatus.textContent = "Status failed";
  }
}

async function refreshCurrentDocument() {
  if (!documentId) return;
  const response = await fetch("/api/refresh", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId }),
  });
  const payload = await response.json();
  if (!response.ok) {
    externalFileChanged = true;
    status.textContent = `Refresh failed: ${payload.error ?? "unknown error"}`;
    draw();
    return;
  }
  const session = payload.session;
  if (!session?.document) return;
  metadata = session.document.metadata;
  readOnly = Boolean(session.readOnly);
  externalFileChanged = false;
  loadDocument(session.document);
  updateFileList(session.displayName || fileNameFromPath(metadata.path), metadata.path, totalRows);
  draw();
  status.textContent = "Auto-refreshed external changes";
}

function startFileStatusChecks() {
  if (fileStatusTimer) clearInterval(fileStatusTimer);
  fileStatusTimer = window.setInterval(() => checkFileStatus(true), 5000);
}

async function runSQLQuery() {
  if (!documentId) return;
  const query = sqlInput.value.trim();
  if (!query) {
    status.textContent = "Enter a SQL query.";
    return;
  }
  status.textContent = "Running SQL...";
  lastSQLResult = null;
  clearSQLResultSelection();
  syncSQLResultActions();
  const response = await fetch("/api/sql", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, query, limit: 500 }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `SQL failed: ${payload.error ?? "unknown error"}`;
    sqlResult.replaceChildren(emptySQLNode(payload.error ?? "Query failed"));
    clearSQLResultSelection();
    syncSQLResultActions();
    return;
  }
  lastSQLResult = payload;
  renderSQLResult(payload);
  rememberSQLQuery(query);
  status.textContent = `SQL returned ${payload.rows.length.toLocaleString()} of ${payload.totalRows.toLocaleString()} rows`;
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
    status.textContent = "SQL query history could not be saved.";
  }
  renderSQLHistory();
}

function clearSQLHistory() {
  try {
    localStorage.removeItem(sqlHistoryStorageKey);
  } catch {
    status.textContent = "SQL query history could not be cleared.";
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
    empty.textContent = "No query history";
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
      status.textContent = "Loaded SQL query from history.";
    });
    sqlHistoryList.append(button);
  }
}

function compactSQLQuery(query) {
  return query.replace(/\s+/g, " ").trim();
}

async function openProfilePanel() {
  if (!documentId || profileTabButton.disabled) return;
  setViewMode("profile");
  await loadProfile();
}

async function loadProfile(force = false) {
  if (!documentId || profileLoading) return;
  if (profileData && !force) {
    renderProfile(profileData);
    return;
  }
  profileLoading = true;
  profileSummary.replaceChildren(summaryNode("Status", "loading"));
  profileColumns.replaceChildren(emptyProfileNode("Profiling columns..."));
  profileIssues.replaceChildren(emptyProfileNode("Checking data quality..."));
  status.textContent = "Profiling data...";
  try {
    const response = await fetch("/api/profile", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ documentId }),
    });
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "profile failed");
    profileData = payload;
    renderProfile(payload);
    updateInspector();
    status.textContent =
      `Profile: ${formatCount(payload.rows ?? 0)} rows, ${formatCount(payload.columns?.length ?? 0)} columns, ${formatCount(payload.issues?.length ?? 0)} issues`;
  } catch (error) {
    status.textContent = `Profile failed: ${error.message}`;
    profileColumns.replaceChildren(emptyProfileNode(error.message));
    profileIssues.replaceChildren(emptyProfileNode("No validation result."));
  } finally {
    profileLoading = false;
  }
}

function renderProfile(result) {
  const issues = result.issues ?? [];
  const errorCount = issues.filter((issue) => issue.severity === "error").length;
  const warningCount = issues.filter((issue) => issue.severity === "warning").length;
  profileSummary.replaceChildren(
    summaryNode("Rows", result.rows ?? 0),
    summaryNode("Columns", result.columns?.length ?? 0),
    summaryNode("Errors", errorCount),
    summaryNode("Warnings", warningCount),
  );
  renderProfileColumns(result.columns ?? []);
  renderProfileIssues(issues);
}

function renderProfileColumns(columnsProfile) {
  if (!columnsProfile.length) {
    profileColumns.replaceChildren(emptyProfileNode("No columns found."));
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const headerRow = document.createElement("tr");
  for (const label of ["Column", "Type", "Rows", "Nulls", "Distinct", "Min", "Max", "Average", "Top Values"]) {
    const th = document.createElement("th");
    th.textContent = label;
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
    appendCell(tr, profile.type ?? "unknown");
    appendCell(tr, formatCount(profile.rows ?? 0));
    appendCell(tr, `${formatCount(profile.nulls ?? 0)} (${formatPercent(profile.nullRatio ?? 0)})`);
    appendCell(tr, formatCount(profile.distinct ?? 0));
    appendCell(tr, profile.min ?? "");
    appendCell(tr, profile.max ?? "");
    appendCell(tr, profile.average == null ? "" : formatNumber(profile.average));
    appendCell(
      tr,
      (profile.topValues ?? [])
        .map((item) => `${item.value || "(empty)"} x${formatCount(item.count)}`)
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
    profileIssues.replaceChildren(emptyProfileNode("No validation issues."));
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const headerRow = document.createElement("tr");
  for (const label of ["Severity", "Type", "Row", "Column", "Message"]) {
    const th = document.createElement("th");
    th.textContent = label;
    headerRow.append(th);
  }
  thead.append(headerRow);
  table.append(thead);

  const tbody = document.createElement("tbody");
  for (const issue of issues) {
    const tr = document.createElement("tr");
    tr.className = issue.severity === "error" ? "issue-error" : "issue-warning";
    appendCell(tr, issue.severity ?? "warning");
    appendCell(tr, issue.type ?? "");
    appendCell(tr, issue.row || "");
    appendCell(tr, issue.col || "");
    appendCell(tr, issue.message ?? "");
    tbody.append(tr);
  }
  table.append(tbody);
  profileIssues.replaceChildren(table);
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
      rightLabel: diffRightTextMode.checked ? "text" : "",
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
    updateDiffNavigationState();
    return;
  }
  renderDiffResult(payload);
  status.textContent =
    `Diff: ${payload.summary.changed} changed, ${payload.summary.added} added, ${payload.summary.deleted} deleted`;
}

function currentDiffOptions() {
  return {
    ignoreCase: diffIgnoreCase.checked,
    ignoreWhitespace: diffIgnoreWhitespace.checked,
    ignoredColumns: parseColumnList(diffIgnoredColumns.value),
    matchColumnsByName: diffMatchHeaders.checked,
    keyColumn: diffKeyColumn.value.trim(),
  };
}

function parseColumnList(value) {
  return String(value ?? "")
    .split(/[,\s]+/)
    .map((item) => Number.parseInt(item, 10))
    .filter((item) => Number.isInteger(item) && item > 0);
}

function renderDiffResult(result) {
  const summary = result.summary ?? {};
  diffRows = [];
  activeDiffIndex = -1;
  diffSummary.replaceChildren(
    summaryNode("Unchanged", summary.unchanged ?? 0),
    summaryNode("Changed", summary.changed ?? 0),
    summaryNode("Added", summary.added ?? 0),
    summaryNode("Deleted", summary.deleted ?? 0),
  );
  if (!result.rows?.length) {
    diffResult.replaceChildren(emptyDiffNode("No diff rows"));
    updateDiffNavigationState();
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const tr = document.createElement("tr");
  for (const label of ["Status", "Left Row", "Left", "Right Row", "Right", "Changed Cells", "Actions"]) {
    const th = document.createElement("th");
    th.textContent = label;
    tr.append(th);
  }
  thead.append(tr);
  table.append(thead);
  const tbody = document.createElement("tbody");
  result.rows.forEach((row, index) => {
    const tr = document.createElement("tr");
    tr.className = row.status;
    if (row.status !== "unchanged") {
      tr.dataset.diffIndex = String(diffRows.length);
      tr.dataset.resultIndex = String(index);
      diffRows.push(tr);
    }
    appendCell(tr, row.status);
    appendCell(tr, row.leftRow ?? "");
    appendCell(tr, (row.leftRecord ?? []).join(" | "));
    appendCell(tr, row.rightRow ?? "");
    appendCell(tr, (row.rightRecord ?? []).join(" | "));
    const cells = document.createElement("td");
    cells.className = "diff-cells";
    appendChangedCells(cells, row.changedCells ?? []);
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
  status.textContent = `Diff ${activeDiffIndex + 1} of ${diffRows.length}`;
}

function navigateDiff(direction) {
  if (viewMode !== "diff" || diffRows.length === 0) return false;
  activateDiffRow(activeDiffIndex < 0 ? 0 : activeDiffIndex + direction);
  return true;
}

function appendDiffActions(container, row, resultIndex) {
  if (row.status === "unchanged") return;
  const rowButton = document.createElement("button");
  rowButton.type = "button";
  rowButton.textContent =
    row.status === "added"
      ? "Insert row"
      : row.status === "deleted"
        ? "Delete left"
        : "Apply row";
  rowButton.title = "Apply right row to left file";
  rowButton.addEventListener("click", () => mergeDiffChange(resultIndex, "row", 0, "right-to-left"));
  container.append(rowButton);

  if (!diffRightTextMode.checked) {
    const leftButton = document.createElement("button");
    leftButton.type = "button";
    leftButton.textContent =
      row.status === "added"
        ? "Delete right"
        : row.status === "deleted"
          ? "Insert right"
          : "Apply left row";
    leftButton.title = "Apply left row to right file";
    leftButton.addEventListener("click", () => mergeDiffChange(resultIndex, "row", 0, "left-to-right"));
    container.append(leftButton);
  }

  if (row.status !== "changed") return;
  for (const cell of row.changedCells ?? []) {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = `c${cell.col}`;
    button.title = `Apply right cell c${cell.col} to left`;
    button.addEventListener("click", () => mergeDiffChange(resultIndex, "cell", cell.col, "right-to-left"));
    container.append(button);
    if (!diffRightTextMode.checked) {
      const leftButton = document.createElement("button");
      leftButton.type = "button";
      leftButton.textContent = `L c${cell.col}`;
      leftButton.title = `Apply left cell c${cell.col} to right`;
      leftButton.addEventListener("click", () => mergeDiffChange(resultIndex, "cell", cell.col, "left-to-right"));
      container.append(leftButton);
    }
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
      rightLabel: diffRightTextMode.checked ? "text" : "",
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

function appendChangedCells(container, changedCells) {
  if (!changedCells.length) return;
  changedCells.forEach((cell, index) => {
    if (index > 0) container.append(document.createTextNode("; "));
    const item = document.createElement("span");
    item.className = "diff-cell-change";
    item.append(document.createTextNode(`c${cell.col}: `));
    appendDiffSpans(item, cell.leftSpans, cell.left);
    item.append(document.createTextNode(" -> "));
    appendDiffSpans(item, cell.rightSpans, cell.right);
    container.append(item);
  });
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

function renderSQLResult(result) {
  const resultColumns = sqlResultColumns(result);
  const resultRows = result.rows ?? [];
  if (!resultColumns.length && !resultRows.length) {
    sqlResult.replaceChildren(emptySQLNode("No rows"));
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
    td.textContent = "No rows";
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
    const totalWidth = columns.reduce((sum, column) => sum + column.width, 0);
    const maxX = Math.max(0, totalWidth - canvas.clientWidth + rowNumberWidth);
    applyVerticalWheelDelta(event.deltaY);
    scrollX = Math.max(0, Math.min(maxX, scrollX + event.deltaX));
    if (editing) commitEditor();
    draw();
  },
  { passive: false },
);

canvas.addEventListener("mousedown", (event) => {
  canvas.focus();
  if (!documentId || columns.length === 0 || totalRows === 0) return;
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  if (event.metaKey || event.ctrlKey) {
    const url = urlAtPoint(x, y);
    if (url) {
      event.preventDefault();
      window.open(url, "_blank", "noopener");
      status.textContent = `Opened ${url}`;
      return;
    }
  }
  if (y < headerHeight) {
    const col = columnAt(x);
    if (x >= rowNumberWidth && col >= 0) {
      event.preventDefault();
      if (editing) commitEditor();
      if (isHeaderMenuTarget(x, col)) {
        openColumnMenu(col, event.clientX, rect.top + headerHeight + 4);
      } else {
        closeColumnMenu();
        if (event.altKey) {
          includeHeadersInSelection = true;
          document.body.dataset.headerInSelection = "true";
        }
        selectColumnRange(col, event.shiftKey);
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
    selectRowRange(rowAt(y), event.shiftKey);
    draggingSelection = true;
    dragSelectionMode = "row";
    draw();
    return;
  }
  const row = Math.floor((y - headerHeight + scrollY) / rowHeight);
  const col = columnAt(x);
  if (row >= 0 && row < totalRows && col >= 0) {
    event.preventDefault();
    selected = { row, col };
    if (!event.shiftKey) {
      selectionAnchor = { ...selected };
    }
    draggingSelection = true;
    dragSelectionMode = "cell";
    if (editing) commitEditor();
    draw();
  }
});

window.addEventListener("mousemove", (event) => {
  if (draggingSelection) {
    event.preventDefault();
    updateDragSelection(event);
    return;
  }
  updateURLHover(event);
});

window.addEventListener("mouseup", () => {
  if (!draggingSelection) return;
  draggingSelection = false;
  dragSelectionMode = "cell";
});

canvas.addEventListener("mouseleave", clearURLHover);

canvas.addEventListener("dblclick", (event) => {
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
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
  const text = event.clipboardData.getData("text/plain");
  if (!text) return;
  event.preventDefault();
  pasteText(text);
});

window.addEventListener("keydown", (event) => {
  if (handleShortcutEditorKeydown(event)) return;
  if (shortcutEditorOpen) return;
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
  if (event.key === "Escape" && conditionBuilderOpen) {
    event.preventDefault();
    closeConditionBuilder();
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
  if (key === "s" && defaultShortcutActive("file.save")) {
    event.preventDefault();
    saveCurrentDocument();
    return true;
  }
  if (key === "f" && defaultShortcutActive("edit.find")) {
    event.preventDefault();
    commandInput.focus();
    commandInput.select();
    return true;
  }
  if (key === "l" && defaultShortcutActive("view.filter")) {
    event.preventDefault();
    commandInput.focus();
    commandInput.select();
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
  const diffLeft = params.get("diffLeft");
  const diffRight = params.get("diffRight");
  const diffText = params.get("diffRightText");
  if (diffLeft || diffRight || diffText) {
    leftDiffPath.value = diffLeft ?? "";
    rightDiffPath.value = diffRight ?? "";
    diffRightTextMode.checked = diffText != null;
    diffRightText.hidden = !diffRightTextMode.checked;
    diffRightText.value = diffText ?? "";
    diffIgnoreCase.checked = truthyParam(params.get("diffIgnoreCase"));
    diffIgnoreWhitespace.checked = truthyParam(params.get("diffIgnoreWhitespace"));
    diffMatchHeaders.checked = truthyParam(params.get("diffMatchHeaders"));
    diffKeyColumn.value = params.get("diffKeyColumn") ?? "";
    diffIgnoredColumns.value = params.get("diffIgnoredColumns") ?? "";
    openDiffPanel();
    if (truthyParam(params.get("diffRun")) && diffLeft && (diffRight || diffText)) {
      await runDiffCompare();
    }
    return;
  }
  const path = params.get("open");
  if (!path) return;
  pathInput.value = path;
  try {
    await openPath(path);
  } catch (error) {
    status.textContent = `Open failed: ${error.message}`;
  }
}

function truthyParam(value) {
  return value === "1" || value === "true" || value === "yes";
}

openButton.addEventListener("click", openPathInputOrPicker);

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
exportButton.addEventListener("click", exportCurrentView);
exportXlsxButton.addEventListener("click", () => exportCurrentView("xlsx"));
readOnlyButton.addEventListener("click", toggleReadOnly);
gridTabButton.addEventListener("click", () => setViewMode("grid"));
sqlTabButton.addEventListener("click", () => {
  if (!sqlTabButton.disabled) setViewMode("sql");
});
profileTabButton.addEventListener("click", openProfilePanel);
settingsTabButton.addEventListener("click", openSettingsPanel);
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
applyQuotePolicyButton.addEventListener("click", () => executeCommand("format.applyQuotePolicy"));
quoteModeSelect.addEventListener("change", () => {
  quotePolicyDraftDirty = true;
  settingsQuoteSummary.textContent = `${quoteModeSelect.value} ${quoteCharSelect.value}`;
});
quoteCharSelect.addEventListener("change", () => {
  quotePolicyDraftDirty = true;
  settingsQuoteSummary.textContent = `${quoteModeSelect.value} ${quoteCharSelect.value}`;
});
runSqlButton.addEventListener("click", runSQLQuery);
copySqlResultButton.addEventListener("click", () => executeCommand("sql.copyResultCsv"));
exportSqlResultButton.addEventListener("click", () => executeCommand("sql.exportResultCsv"));
clearSqlHistoryButton.addEventListener("click", clearSQLHistory);
sqlResult.addEventListener("click", (event) => {
  const cell = sqlResultCellFromEvent(event);
  if (!cell) return;
  selectSQLResultCell(cell.row, cell.col, event.shiftKey);
});
sqlResult.addEventListener("keydown", async (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "c") {
    event.preventDefault();
    await copySQLResultCSV();
    return;
  }
  if (event.key === "Escape") {
    event.preventDefault();
    clearSQLResultSelection();
    status.textContent = "Cleared SQL result selection.";
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
diffRightTextMode.addEventListener("change", () => {
  diffRightText.hidden = !diffRightTextMode.checked;
  if (diffRightTextMode.checked) diffRightText.focus();
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
shortcutsPanel.addEventListener("mousedown", (event) => {
  if (event.target === shortcutsPanel) {
    event.preventDefault();
    closeKeyboardShortcuts();
  }
});
columnMenu.addEventListener("mousedown", (event) => {
  event.stopPropagation();
});
window.addEventListener("mousedown", (event) => {
  if (!columnMenu.hidden && !columnMenu.contains(event.target) && event.target !== canvas) {
    closeColumnMenu();
  }
});
columnSortAsc.addEventListener("click", async () => {
  const column = columnMenuColumn;
  closeColumnMenu();
  await sortColumn(column, false, false);
});
columnSortDesc.addEventListener("click", async () => {
  const column = columnMenuColumn;
  closeColumnMenu();
  await sortColumn(column, true, false);
});
columnSortLengthAsc.addEventListener("click", async () => {
  const column = columnMenuColumn;
  closeColumnMenu();
  await sortColumn(column, false, false, "length");
});
columnSortLengthDesc.addEventListener("click", async () => {
  const column = columnMenuColumn;
  closeColumnMenu();
  await sortColumn(column, true, false, "length");
});
columnAddSortAsc.addEventListener("click", async () => {
  const column = columnMenuColumn;
  closeColumnMenu();
  await sortColumn(column, false, true);
});
columnAddSortDesc.addEventListener("click", async () => {
  const column = columnMenuColumn;
  closeColumnMenu();
  await sortColumn(column, true, true);
});
columnFilterValue.addEventListener("click", filterColumnBySelectedValue);
columnFilterBlank.addEventListener("click", () => filterColumnBlank(true));
columnFilterFilled.addEventListener("click", () => filterColumnBlank(false));
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
columnClearView.addEventListener("click", async () => {
  closeColumnMenu();
  await clearView();
});
cellValueInput.addEventListener("input", () => {
  cellValueInputDirty = true;
  applyCellValueButton.disabled = readOnly || !documentId;
});
cellValueInput.addEventListener("keydown", (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
    event.preventDefault();
    applyCellValue();
  }
});
applyCellValueButton.addEventListener("click", applyCellValue);
openCellLink.addEventListener("click", (event) => {
  if (openCellLink.getAttribute("aria-disabled") === "true") {
    event.preventDefault();
    status.textContent = "Selected cell is not a link";
  }
});
findButton.addEventListener("click", findCurrentQuery);
findPrevButton.addEventListener("click", () => findNextMatch(-1));
findNextButton.addEventListener("click", () => findNextMatch(1));
replaceButton.addEventListener("click", () => replaceCurrentQuery(false));
replaceAllButton.addEventListener("click", () => replaceCurrentQuery(true));
findCaseButton.addEventListener("click", () => toggleFindOption("case"));
findWholeButton.addEventListener("click", () => toggleFindOption("whole"));
findRegexButton.addEventListener("click", () => toggleFindOption("regex"));
findSelectionButton.addEventListener("click", () => toggleFindOption("selection"));
filterButton.addEventListener("click", filterCurrentQuery);
filterBuilderButton.addEventListener("click", openConditionBuilder);
sortButton.addEventListener("click", sortSelectedColumn);
closeConditionButton.addEventListener("click", closeConditionBuilder);
conditionOperator.addEventListener("change", updateConditionValueState);
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
for (const input of [conditionColumn, conditionOperator, conditionValue]) {
  input.addEventListener("keydown", async (event) => {
    if (event.key !== "Enter") return;
    event.preventDefault();
    await applyConditionBuilder();
  });
}
commandInput.addEventListener("keydown", async (event) => {
  if (event.key !== "Enter") return;
  event.preventDefault();
  await findNextMatch(event.shiftKey ? -1 : 1);
});
commandInput.addEventListener("input", invalidateFindState);
window.addEventListener("resize", resize);
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
renderSettingsPanel();
openInitialURLTarget();
