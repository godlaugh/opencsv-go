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
const filterButton = document.querySelector("#filterButton");
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
const diffButton = document.querySelector("#diffButton");
const saveButton = document.querySelector("#saveButton");
const readOnlyButton = document.querySelector("#readOnlyButton");
const commandButtons = Array.from(document.querySelectorAll("[data-command]"));
const gridFrame = document.querySelector(".grid-frame");
const gridTabButton = document.querySelector("#gridTabButton");
const sqlTabButton = document.querySelector("#sqlTabButton");
const profileTabButton = document.querySelector("#profileTabButton");
const paletteButton = document.querySelector("#paletteButton");
const sqlPanel = document.querySelector("#sqlPanel");
const sqlInput = document.querySelector("#sqlInput");
const runSqlButton = document.querySelector("#runSqlButton");
const sqlResult = document.querySelector("#sqlResult");
const diffPanel = document.querySelector("#diffPanel");
const leftDiffPath = document.querySelector("#leftDiffPath");
const rightDiffPath = document.querySelector("#rightDiffPath");
const runDiffButton = document.querySelector("#runDiffButton");
const diffIgnoreCase = document.querySelector("#diffIgnoreCase");
const diffIgnoreWhitespace = document.querySelector("#diffIgnoreWhitespace");
const diffMatchHeaders = document.querySelector("#diffMatchHeaders");
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
const commandPalette = document.querySelector("#commandPalette");
const paletteInput = document.querySelector("#paletteInput");
const paletteList = document.querySelector("#paletteList");
const columnMenu = document.querySelector("#columnMenu");
const columnMenuTitle = document.querySelector("#columnMenuTitle");
const columnSortAsc = document.querySelector("#columnSortAsc");
const columnSortDesc = document.querySelector("#columnSortDesc");
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
const headerHeight = 36;
const rowNumberWidth = 58;
const windowSize = 160;

let totalRows = 0;
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
let profileLoading = false;
let paletteOpen = false;
let paletteItems = [];
let paletteIndex = 0;
let columnMenuColumn = -1;
let readOnly = false;
let externalFileChanged = false;
let fileStatusTimer = 0;
let cellValueInputDirty = false;
let cellValueInputKey = "";
let findCaseSensitive = false;
let findWholeCell = false;
let findRegex = false;
let findState = emptyFindState();
let includeHeadersInSelection = false;
let diffRows = [];
let activeDiffIndex = -1;
const quotePolicy = "minimal";
const themeStorageKey = "opencsv-theme";
const themeMediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
let themePreference = localStorage.getItem(themeStorageKey) || "system";
let themeColors = {};

initializeTheme();

function resize() {
  const rect = gridFrame.getBoundingClientRect();
  const scale = window.devicePixelRatio || 1;
  resizeCanvas(canvas, ctx, rect.width, rect.height, scale);
  draw();
}

function resizeCanvas(targetCanvas, targetCtx, width, height, scale) {
  const pixelWidth = Math.max(1, Math.floor(width * scale));
  const pixelHeight = Math.max(1, Math.floor(height * scale));
  if (targetCanvas.width !== pixelWidth) targetCanvas.width = pixelWidth;
  if (targetCanvas.height !== pixelHeight) targetCanvas.height = pixelHeight;
  targetCtx.setTransform(scale, 0, 0, scale, 0, 0);
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
    quoteMarker: read("--quote-marker"),
    quoteMarkerText: read("--quote-marker-text"),
    selectionBg: read("--selection-bg"),
    text: read("--text"),
  };
}

function valueAt(row, col) {
  if (!documentId) return "";
  const record = rows.get(row);
  if (!record) {
    maybeLoadWindow(row);
    return "";
  }
  return record[sourceColFor(col)] ?? "";
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
  const width = canvas.clientWidth;
  const height = canvas.clientHeight;
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
  const aggregate = selectionAggregate();
  const modeLabel = readOnly ? " - read-only" : "";
  const externalLabel = externalFileChanged ? " - external change" : "";
  document.body.dataset.selectionRange = documentId ? selectionSummary() : "";
  status.textContent =
    documentId
      ? `${totalRows.toLocaleString()} rows, ${columns.length} columns${path}${dirtyLabel}${modeLabel}${externalLabel} - row ${selected.row + 1}, ${selectedColumn}${aggregate}`
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
  sqlPanel.hidden = mode !== "sql";
  diffPanel.hidden = mode !== "diff";
  profilePanel.hidden = mode !== "profile";
  gridTabButton.classList.toggle("active", mode === "grid");
  sqlTabButton.classList.toggle("active", mode === "sql");
  profileTabButton.classList.toggle("active", mode === "profile");
  if (mode === "grid") {
    resize();
    canvas.focus();
  } else if (mode === "sql") {
    sqlInput.focus();
  } else if (mode === "diff") {
    leftDiffPath.focus();
  } else {
    profilePanel.focus();
  }
  updateCommandState();
}

async function loadCommandRegistry() {
  try {
    const response = await fetch("/api/commands");
    const payload = await response.json();
    if (!response.ok) throw new Error(payload.error ?? "commands failed");
    commandsById = new Map((payload.commands ?? []).map((command) => [command.id, command]));
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
    button.title = command.shortcut
      ? `${command.description} (${command.shortcut})`
      : command.description;
    if (button.dataset.command === "edit.replaceAll") {
      button.textContent = "All";
    } else if (button.dataset.command === "command.palette") {
      button.textContent = "Cmd";
    } else if (button.dataset.command === "file.toggleReadOnly") {
      button.textContent = "Read Only";
    } else if (button.dataset.command === "diff.open") {
      button.textContent = "Compare";
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
}

function isCommandDisabled(command) {
  const hasDocument = Boolean(documentId);
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
      command.id === "cell.joinHorizontal" ||
      command.id === "cell.splitHorizontal"
    );
  return (
    command.implemented === false ||
    (command.requiresDocument && !hasDocument) ||
    writeBlocked
  );
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
    readOnlyStatus.textContent = "Editable";
    externalStatus.textContent = "Synced";
    return;
  }
  encodingStatus.textContent = encodingLabel();
  delimiterStatus.textContent = delimiterLabel(metadata.dialect);
  quoteStatus.textContent = "Quote minimal";
  headerStatus.textContent = metadata.dialect?.hasHeader ? "Header row" : "No header";
  readOnlyStatus.textContent = readOnly ? "Read-only" : "Editable";
  externalStatus.textContent = externalFileChanged ? "Changed on disk" : "Synced";
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

function documentFormatSummary() {
  return `${encodingLabel()} · ${delimiterLabel(metadata.dialect)} · ${lineEndingLabel()}`;
}

function quoteLabel(value) {
  return requiresQuote(value) ? "will quote on save" : "no quote needed";
}

function syncCellValuePanel() {
  const key = selectedCellKey();
  const value = documentId ? valueAt(selected.row, selected.col) : "";
  if (key !== cellValueInputKey || !cellValueInputDirty) {
    cellValueInput.value = value;
    cellValueInputKey = key;
    cellValueInputDirty = false;
  }
  cellValueInput.disabled = !documentId || readOnly;
  applyCellValueButton.disabled = !documentId || readOnly || !cellValueInputDirty;
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
  return documentId ? `${documentId}:${selected.row}:${selected.col}` : "";
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

function selectionAggregate() {
  const count = selectionCellCount();
  if (count <= 1) return "";
  const range = normalizedSelection();
  if (count > 5000) {
    const headerNote = includeHeadersInSelection ? ", headers included" : "";
    return ` - selected ${formatCount(count)} cells${headerNote}`;
  }
  let numericCount = 0;
  let sum = 0;
  let minValue = Number.POSITIVE_INFINITY;
  let maxValue = Number.NEGATIVE_INFINITY;
  for (let row = range.startRow; row <= range.endRow; row++) {
    for (let col = range.startCol; col <= range.endCol; col++) {
      const raw = valueAt(row, col).trim();
      if (raw === "") continue;
      const value = Number(raw);
      if (!Number.isFinite(value)) continue;
      numericCount++;
      sum += value;
      minValue = Math.min(minValue, value);
      maxValue = Math.max(maxValue, value);
    }
  }
  if (!numericCount) {
    return ` - selected ${formatCount(count)} cells`;
  }
  const avg = sum / numericCount;
  return ` - selected ${formatCount(count)} cells, sum ${formatNumber(sum)}, avg ${formatNumber(avg)}, min ${formatNumber(minValue)}, max ${formatNumber(maxValue)}`;
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

function renderCommandPalette() {
  const query = paletteInput.value.trim().toLowerCase();
  paletteItems = Array.from(commandsById.values())
    .filter((command) => command.implemented !== false)
    .filter((command) => commandMatchesQuery(command, query));
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
    shortcut.textContent = command.shortcut ?? "";
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

function commandMatchesQuery(command, query) {
  if (!query) return true;
  return [
    command.id,
    command.title,
    command.description,
    ...(command.aliases ?? []),
  ]
    .filter(Boolean)
    .some((value) => value.toLowerCase().includes(query));
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
  switch (commandID) {
    case "command.palette":
      openCommandPalette();
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
    case "file.toggleReadOnly":
      await toggleReadOnly();
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
    case "edit.copyWithHeaders":
      await copySelection("tsv", { includeHeaders: true });
      break;
    case "selection.toggleHeaders":
      toggleHeaderInSelection();
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
    case "cell.joinHorizontal":
      await joinSelectedCells();
      break;
    case "cell.splitHorizontal":
      await splitActiveCell();
      break;
    case "view.filter":
      await filterCurrentQuery();
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
  ctx.save();
  ctx.fillStyle = themeColors.gridHeader;
  ctx.fillRect(0, 0, width, headerHeight);
  ctx.fillStyle = themeColors.gridGutter;
  ctx.fillRect(0, 0, rowNumberWidth, headerHeight);
  ctx.strokeStyle = themeColors.gridLine;
  ctx.beginPath();
  ctx.moveTo(0, headerHeight - 0.5);
  ctx.lineTo(width, headerHeight - 0.5);
  ctx.moveTo(rowNumberWidth + 0.5, 0);
  ctx.lineTo(rowNumberWidth + 0.5, headerHeight);
  ctx.stroke();

  ctx.font = "600 12px Inter, system-ui, sans-serif";
  ctx.textBaseline = "middle";

  if (columns.length === 0) {
    ctx.fillStyle = themeColors.muted;
    ctx.fillText("Open a CSV file", rowNumberWidth + 10, headerHeight / 2);
    ctx.restore();
    return;
  }

  let x = rowNumberWidth - scrollX;
  for (const column of columns) {
    if (x + column.width >= rowNumberWidth && x <= width) {
      ctx.fillStyle = themeColors.text;
      ctx.fillText(column.name, x + 8, 13);
      ctx.fillStyle = themeColors.muted;
      ctx.font = "600 11px Inter, system-ui, sans-serif";
      ctx.fillText(column.type ?? "string", x + 8, 27);
      ctx.font = "600 12px Inter, system-ui, sans-serif";
      drawHeaderMenuGlyph(x, column.width);
      ctx.strokeStyle = themeColors.gridLine;
      ctx.beginPath();
      ctx.moveTo(x + column.width + 0.5, 0);
      ctx.lineTo(x + column.width + 0.5, headerHeight);
      ctx.stroke();
    }
    x += column.width;
  }
  ctx.restore();
}

function drawHeaderMenuGlyph(x, width) {
  const centerX = x + width - 14;
  const centerY = 13;
  if (centerX < rowNumberWidth + 8) return;
  ctx.save();
  ctx.fillStyle = themeColors.muted;
  ctx.beginPath();
  ctx.moveTo(centerX - 4, centerY - 2);
  ctx.lineTo(centerX + 4, centerY - 2);
  ctx.lineTo(centerX, centerY + 3);
  ctx.closePath();
  ctx.fill();
  ctx.restore();
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

    ctx.fillStyle = row % 2 === 0 ? themeColors.gridRow : themeColors.gridRowAlt;
    ctx.fillRect(0, y, width, rowHeight);
    ctx.fillStyle = themeColors.gridHeader;
    ctx.fillRect(0, y, rowNumberWidth, rowHeight);
    ctx.fillStyle = themeColors.muted;
    ctx.fillText(String(row + 1), 10, y + rowHeight / 2);

    let x = rowNumberWidth - scrollX;
    for (let col = 0; col < columns.length; col++) {
      const column = columns[col];
      if (x + column.width >= rowNumberWidth && x <= width) {
        const value = valueAt(row, col);
        ctx.fillStyle = themeColors.text;
        ctx.save();
        ctx.beginPath();
        ctx.rect(x + 1, y + 1, column.width - 2, rowHeight - 2);
        ctx.clip();
        ctx.fillText(value, x + 8, y + rowHeight / 2);
        if (requiresQuote(value)) {
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

function requiresQuote(value) {
  if (quotePolicy !== "minimal") return quotePolicy === "always";
  const text = String(value ?? "");
  if (text === "") return false;
  const delimiter = delimiterChar();
  return (
    text.includes(delimiter) ||
    text.includes('"') ||
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
  if (extend) {
    selectionAnchor = { row: totalRows - 1, col: selectionAnchor.col };
  } else {
    selectionAnchor = { row: totalRows - 1, col: clampedCol };
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
  windowGeneration++;
  const range = normalizedSelection();
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
  submitCellEdits(edits, range.startRow, `${formatCount(edits.length)} cells edited`);
}

function setActiveCellValue(value) {
  if (!ensureEditable("Edit")) return;
  if (!documentId || columns.length === 0) return;
  windowGeneration++;
  const record = rows.get(selected.row) ?? [];
  const sourceCol = sourceColFor(selected.col);
  while (record.length <= sourceCol) record.push("");
  record[sourceCol] = value;
  rows.set(selected.row, record);
  selectionAnchor = { ...selected };
  cellValueInputDirty = false;
  submitCellEdits(
    [{ row: rowSource.get(selected.row) ?? selected.row + dataStartRow, col: sourceCol, value }],
    selected.row,
    "Cell edited",
  );
}

function ensureEditable(action) {
  if (!readOnly) return true;
  status.textContent = `${action} blocked: document is read-only`;
  return false;
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

async function openPath(path) {
  status.textContent = `Opening ${path}...`;
  const response = await fetch("/api/open", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ path, previewRows: windowSize, indexRows: -1 }),
  });
  const payload = await response.json();
  if (!response.ok) throw new Error(payload.error ?? "open failed");

  documentId = payload.documentId;
  metadata = payload.document.metadata;
  readOnly = false;
  externalFileChanged = false;
  loadDocument(payload.document);
  updateFileList(fileNameFromPath(path), metadata.path, totalRows);
  draw();
  startFileStatusChecks();
  await checkFileStatus(false);
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
  profileColumns.replaceChildren(emptyProfileNode("Run Profile to inspect columns."));
  profileIssues.replaceChildren(emptyProfileNode("No validation issues loaded."));
  profileSummary.replaceChildren();
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
  totalRows = Math.max(0, (document.rowIndex?.rows || document.preview.length) - dataStartRow);
  allColumns = buildColumns(document);
  hiddenSourceColumns = new Set();
  applyHiddenColumns();
  selected = { row: 0, col: 0 };
  selectionAnchor = { row: 0, col: 0 };
  scrollX = 0;
  scrollY = 0;
  wheelRemainderY = 0;
  includeHeadersInSelection = false;
  window.document.body.dataset.headerInSelection = "false";
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
  includeHeadersInSelection = !includeHeadersInSelection;
  document.body.dataset.headerInSelection = includeHeadersInSelection ? "true" : "false";
  draw();
  status.textContent = includeHeadersInSelection
    ? "Header included when copying selection"
    : "Header excluded from selection copy";
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
  windowGeneration++;
  const edits = [];
  const startRow = selected.row;
  const startCol = selected.col;
  for (let rowOffset = 0; rowOffset < matrix.length; rowOffset++) {
    const row = startRow + rowOffset;
    if (row >= totalRows) break;
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

function updateFileList(name, detail, rowCount) {
  fileList.replaceChildren();
  const button = document.createElement("button");
  button.className = "file-button active";
  button.type = "button";

  const title = document.createElement("span");
  title.className = "file-name";
  title.textContent = name;

  const meta = document.createElement("span");
  meta.className = "file-meta";
  meta.textContent = `${rowCount.toLocaleString()} rows · ${detail}`;

  button.append(title, meta);
  button.addEventListener("click", () => {
    selected = { row: 0, col: 0 };
    selectionAnchor = { row: 0, col: 0 };
    scrollX = 0;
    scrollY = 0;
    wheelRemainderY = 0;
    canvas.focus();
    draw();
  });
  fileList.append(button);
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
    width: Math.max(92, Math.min(240, String(first[index] ?? "").length * 9 + 64)),
    type: index === 0 ? "mixed" : "string",
  }));
}

async function maybeLoadWindow(row) {
  if (!documentId || loadingWindow) return;
  if (row >= windowStart && row < windowEnd - 20) return;
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
    totalRows = viewMapped ? payload.totalRows : Math.max(0, payload.totalRows - dataStartRow);
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
  const response = await fetch("/api/find", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      query: search.query,
      regex: search.regex,
      caseSensitive: search.caseSensitive,
      wholeCell: search.wholeCell,
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
    matches: payload.matches ?? [],
    total: payload.total ?? payload.matches?.length ?? 0,
    truncated: Boolean(payload.truncated),
    index: 0,
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
  return { query, regex, caseSensitive, wholeCell: findWholeCell };
}

function findKey(search) {
  return JSON.stringify({
    documentId,
    query: search.query,
    regex: search.regex,
    caseSensitive: search.caseSensitive,
    wholeCell: search.wholeCell,
  });
}

function emptyFindState() {
  return { key: "", matches: [], total: 0, truncated: false, index: 0 };
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
  invalidateFindState();
  updateFindOptionButtons();
}

function updateFindOptionButtons() {
  setPressed(findCaseButton, findCaseSensitive);
  setPressed(findWholeButton, findWholeCell);
  setPressed(findRegexButton, findRegex);
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
  totalRows = payload.totalRows ?? 0;
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(0);
  status.textContent = `${label} matched ${totalRows.toLocaleString()} rows`;
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

async function sortColumn(column, descending, append) {
  if (!documentId || columns.length === 0 || !columns[column]) return;
  const nextKeys = nextSortKeys(column, descending, append);
  const primaryKey = nextKeys[0];
  const response = await fetch("/api/sort", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      column: primaryKey.column,
      descending: primaryKey.descending,
      numeric: primaryKey.numeric,
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
  totalRows = payload.totalRows ?? 0;
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(0);
  status.textContent = `Sorted ${formatSortKeys(sortKeys)}`;
}

function nextSortKeys(column, descending, append) {
  const key = sortKeyForColumn(column, descending);
  if (!append || sortKeys.length === 0) return [key];
  return [...sortKeys.filter((item) => item.column !== key.column), key];
}

function sortKeyForColumn(column, descending) {
  const type = columns[column]?.type ?? "string";
  return {
    column: sourceColFor(column),
    descending,
    numeric: type === "integer" || type === "decimal" || columnLooksNumeric(column),
  };
}

function formatSortKeys(keys) {
  return keys
    .map((key) => `${columnNameForSource(key.column)} ${key.descending ? "descending" : "ascending"}`)
    .join(", then ");
}

function columnLooksNumeric(column) {
  let inspected = 0;
  let numeric = 0;
  for (let row = 0; row < Math.min(totalRows, 80); row++) {
    const value = valueAt(row, column);
    if (value == null || String(value).trim() === "") continue;
    inspected++;
    if (Number.isFinite(Number(String(value).trim()))) numeric++;
  }
  return inspected > 0 && numeric === inspected;
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
  totalRows = payload.totalRows ?? 0;
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(0);
  status.textContent = `Value filter ${columns[column].name} matched ${totalRows.toLocaleString()} rows`;
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
  totalRows = Math.max(0, (payload.totalRows ?? 0) - dataStartRow);
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
  await mutateStructure("/api/row/insert", { row: selectedSourceRow() }, "Inserted row");
}

async function deleteSelectedRow() {
  await mutateStructure("/api/row/delete", { row: selectedSourceRow() }, "Deleted row");
}

async function hideSelectedRow() {
  if (!documentId) return;
  if (!ensureEditable("Hide Row")) return;
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
  totalRows = payload.totalRows ?? Math.max(0, totalRows - (payload.changed ?? 0));
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

async function joinSelectedCells() {
  if (!documentId) return;
  if (!ensureEditable("Join Cells")) return;
  if (selectionCellCount() < 2) {
    status.textContent = "Select at least two cells to join.";
    return;
  }
  const cells = selectedCellsForTransform();
  if (!cells) return;
  await runCellTool("/api/cells/join", {
    cells,
    separator: " ",
    activeRow: selected.row,
    activeCol: normalizedSelection().startCol,
    label: "Joined cells",
  });
}

async function splitActiveCell() {
  if (!documentId) return;
  if (!ensureEditable("Split Cell")) return;
  await runCellTool("/api/cells/split", {
    cells: [{ row: selectedSourceRow(), col: sourceColFor(selected.col) }],
    separator: " ",
    activeRow: selected.row,
    activeCol: selected.col,
    label: "Split cell",
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
  status.textContent = "Running SQL...";
  const response = await fetch("/api/sql", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, query: sqlInput.value, limit: 500 }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `SQL failed: ${payload.error ?? "unknown error"}`;
    sqlResult.replaceChildren(emptySQLNode(payload.error ?? "Query failed"));
    return;
  }
  renderSQLResult(payload);
  status.textContent = `SQL returned ${payload.rows.length.toLocaleString()} of ${payload.totalRows.toLocaleString()} rows`;
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
  rowButton.addEventListener("click", () => mergeDiffChange(resultIndex, "row", 0));
  container.append(rowButton);

  if (row.status !== "changed") return;
  for (const cell of row.changedCells ?? []) {
    const button = document.createElement("button");
    button.type = "button";
    button.textContent = `c${cell.col}`;
    button.title = `Apply right cell c${cell.col} to left`;
    button.addEventListener("click", () => mergeDiffChange(resultIndex, "cell", cell.col));
    container.append(button);
  }
}

async function mergeDiffChange(resultIndex, mode, col) {
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
  if (!result.rows?.length) {
    sqlResult.replaceChildren(emptySQLNode("No rows"));
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const headerRow = document.createElement("tr");
  for (const column of result.columns ?? []) {
    const th = document.createElement("th");
    th.textContent = column;
    headerRow.append(th);
  }
  thead.append(headerRow);
  table.append(thead);

  const tbody = document.createElement("tbody");
  for (const row of result.rows) {
    const tr = document.createElement("tr");
    for (const value of row) {
      const td = document.createElement("td");
      td.textContent = value;
      tr.append(td);
    }
    tbody.append(tr);
  }
  table.append(tbody);
  sqlResult.replaceChildren(table);
}

function emptySQLNode(message) {
  const node = document.createElement("div");
  node.className = "sql-empty";
  node.textContent = message;
  return node;
}

canvas.addEventListener("wheel", (event) => {
  event.preventDefault();
  closeColumnMenu();
  const totalWidth = columns.reduce((sum, column) => sum + column.width, 0);
  const maxX = Math.max(0, totalWidth - canvas.clientWidth + rowNumberWidth);
  applyVerticalWheelDelta(event.deltaY);
  scrollX = Math.max(0, Math.min(maxX, scrollX + event.deltaX));
  if (editing) commitEditor();
  draw();
});

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
  if (!documentId || isTypingTarget(event.target)) return;
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
  if (!documentId || isTypingTarget(event.target)) return;
  const text = event.clipboardData.getData("text/plain");
  if (!text) return;
  event.preventDefault();
  pasteText(text);
});

window.addEventListener("keydown", (event) => {
  if (handleGlobalShortcut(event)) return;
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
  if (event.altKey && key === "f5") {
    event.preventDefault();
    return navigateDiff(event.shiftKey ? -1 : 1);
  }
  if (event.shiftKey && event.altKey && key === "h") {
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
    if (arrowCommandID) {
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
      j: "cell.joinHorizontal",
      s: "cell.splitHorizontal",
      t: "text.trim",
      u: "text.uppercase",
      l: "text.lowercase",
    };
    const commandID = commandIDByKey[key];
    if (commandID) {
      event.preventDefault();
      executeCommand(commandID);
      return true;
    }
  }
  if (primary && event.shiftKey && event.key.toLowerCase() === "p") {
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
  if (event.key === "Escape" && selectionCellCount() > 1) {
    event.preventDefault();
    selectionAnchor = { ...selected };
    draw();
    return true;
  }
  if (!primary || paletteOpen) return false;
  if (key === "enter" && event.target === cellValueInput) {
    event.preventDefault();
    applyCellValue();
    return true;
  }
  if (key === "z") {
    event.preventDefault();
    if (event.shiftKey) {
      redoCurrentDocument();
    } else {
      undoCurrentDocument();
    }
    return true;
  }
  if (key === "o") {
    event.preventDefault();
    fileInput.click();
    return true;
  }
  if (key === "s") {
    event.preventDefault();
    saveCurrentDocument();
    return true;
  }
  if (key === "f") {
    event.preventDefault();
    commandInput.focus();
    commandInput.select();
    return true;
  }
  if (key === "l") {
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
readOnlyButton.addEventListener("click", toggleReadOnly);
gridTabButton.addEventListener("click", () => setViewMode("grid"));
sqlTabButton.addEventListener("click", () => {
  if (!sqlTabButton.disabled) setViewMode("sql");
});
profileTabButton.addEventListener("click", openProfilePanel);
paletteButton.addEventListener("click", openCommandPalette);
runSqlButton.addEventListener("click", runSQLQuery);
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
filterButton.addEventListener("click", filterCurrentQuery);
sortButton.addEventListener("click", sortSelectedColumn);
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
openInitialURLTarget();
