const canvas = document.querySelector("#grid");
const ctx = canvas.getContext("2d");
const frozenCanvas = document.querySelector("#frozenGrid");
const frozenCtx = frozenCanvas.getContext("2d");
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
const replaceButton = document.querySelector("#replaceButton");
const replaceAllButton = document.querySelector("#replaceAllButton");
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
const columnClearView = document.querySelector("#columnClearView");

const rowHeight = 28;
const headerHeight = 36;
const rowNumberWidth = 58;
const windowSize = 160;

let totalRows = 0;
let columns = [];

let scrollX = 0;
let scrollY = 0;
let selected = { row: 0, col: 0 };
let selectionAnchor = { row: 0, col: 0 };
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
const quotePolicy = "minimal";

function resize() {
  const rect = gridFrame.getBoundingClientRect();
  const scale = window.devicePixelRatio || 1;
  resizeCanvas(canvas, ctx, rect.width, rect.height, scale);
  resizeCanvas(frozenCanvas, frozenCtx, rect.width, headerHeight, scale);
  draw();
}

function resizeCanvas(targetCanvas, targetCtx, width, height, scale) {
  const pixelWidth = Math.max(1, Math.floor(width * scale));
  const pixelHeight = Math.max(1, Math.floor(height * scale));
  if (targetCanvas.width !== pixelWidth) targetCanvas.width = pixelWidth;
  if (targetCanvas.height !== pixelHeight) targetCanvas.height = pixelHeight;
  targetCtx.setTransform(scale, 0, 0, scale, 0, 0);
}

function valueAt(row, col) {
  if (!documentId) return "";
  const record = rows.get(row);
  if (!record) {
    maybeLoadWindow(row);
    return "";
  }
  return record[col] ?? "";
}

function draw() {
  if (viewMode !== "grid") {
    updateCommandState();
    return;
  }
  const width = canvas.clientWidth;
  const height = canvas.clientHeight;
  ctx.clearRect(0, 0, width, height);

  ctx.fillStyle = "#fbfcfc";
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
      command.id === "column.insert" ||
      command.id === "column.delete" ||
      command.id === "table.transpose"
    );
  return (
    command.implemented === false ||
    (command.requiresDocument && !hasDocument) ||
    writeBlocked
  );
}

function updateInspector() {
  const selectedColumn = columns[selected.col]?.name ?? "";
  const profile = profileData?.columns?.[selected.col];
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

function selectionAggregate() {
  const count = selectionCellCount();
  if (count <= 1) return "";
  const range = normalizedSelection();
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
    case "file.open":
      fileInput.click();
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
    case "column.insert":
      await insertSelectedColumn();
      break;
    case "column.delete":
      await deleteSelectedColumn();
      break;
    case "table.transpose":
      await transposeCurrentDocument();
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
  frozenCtx.clearRect(0, 0, frozenCanvas.clientWidth, headerHeight);
  frozenCtx.fillStyle = "#eef2f3";
  frozenCtx.fillRect(0, 0, width, headerHeight);
  frozenCtx.fillStyle = "#f5f7f7";
  frozenCtx.fillRect(0, 0, rowNumberWidth, headerHeight);
  frozenCtx.strokeStyle = "#d8e0e2";
  frozenCtx.beginPath();
  frozenCtx.moveTo(0, headerHeight - 0.5);
  frozenCtx.lineTo(width, headerHeight - 0.5);
  frozenCtx.moveTo(rowNumberWidth + 0.5, 0);
  frozenCtx.lineTo(rowNumberWidth + 0.5, headerHeight);
  frozenCtx.stroke();

  frozenCtx.font = "600 12px Inter, system-ui, sans-serif";
  frozenCtx.textBaseline = "middle";

  if (columns.length === 0) {
    frozenCtx.fillStyle = "#697579";
    frozenCtx.fillText("Open a CSV file", rowNumberWidth + 10, headerHeight / 2);
    return;
  }

  let x = rowNumberWidth - scrollX;
  for (const column of columns) {
    if (x + column.width >= rowNumberWidth && x <= width) {
      frozenCtx.fillStyle = "#162022";
      frozenCtx.fillText(column.name, x + 8, 13);
      frozenCtx.fillStyle = "#697579";
      frozenCtx.font = "600 11px Inter, system-ui, sans-serif";
      frozenCtx.fillText(column.type ?? "string", x + 8, 27);
      frozenCtx.font = "600 12px Inter, system-ui, sans-serif";
      drawHeaderMenuGlyph(x, column.width);
      frozenCtx.strokeStyle = "#d8e0e2";
      frozenCtx.beginPath();
      frozenCtx.moveTo(x + column.width + 0.5, 0);
      frozenCtx.lineTo(x + column.width + 0.5, headerHeight);
      frozenCtx.stroke();
    }
    x += column.width;
  }
}

function drawHeaderMenuGlyph(x, width) {
  const centerX = x + width - 14;
  const centerY = 13;
  if (centerX < rowNumberWidth + 8) return;
  frozenCtx.save();
  frozenCtx.fillStyle = "#697579";
  frozenCtx.beginPath();
  frozenCtx.moveTo(centerX - 4, centerY - 2);
  frozenCtx.lineTo(centerX + 4, centerY - 2);
  frozenCtx.lineTo(centerX, centerY + 3);
  frozenCtx.closePath();
  frozenCtx.fill();
  frozenCtx.restore();
}

function drawBodySeparators(width, height) {
  if (columns.length === 0 || height <= headerHeight) return;

  ctx.save();
  ctx.beginPath();
  ctx.rect(0, headerHeight, width, height - headerHeight);
  ctx.clip();
  ctx.strokeStyle = "#d8e0e2";
  ctx.beginPath();
  ctx.moveTo(rowNumberWidth + 0.5, headerHeight);
  ctx.lineTo(rowNumberWidth + 0.5, height);
  ctx.stroke();

  ctx.strokeStyle = "#e4eaec";
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
    ctx.fillStyle = "#697579";
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

    ctx.fillStyle = row % 2 === 0 ? "#ffffff" : "#f8faf9";
    ctx.fillRect(0, y, width, rowHeight);
    ctx.fillStyle = "#eef2f3";
    ctx.fillRect(0, y, rowNumberWidth, rowHeight);
    ctx.fillStyle = "#697579";
    ctx.fillText(String(row + 1), 10, y + rowHeight / 2);

    let x = rowNumberWidth - scrollX;
    for (let col = 0; col < columns.length; col++) {
      const column = columns[col];
      if (x + column.width >= rowNumberWidth && x <= width) {
        const value = valueAt(row, col);
        ctx.fillStyle = "#162022";
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

    ctx.strokeStyle = "#edf1f2";
    ctx.beginPath();
    ctx.moveTo(0, y + rowHeight + 0.5);
    ctx.lineTo(width, y + rowHeight + 0.5);
    ctx.stroke();
  }
  ctx.restore();
}

function drawQuoteIndicator(x, y, width) {
  ctx.save();
  ctx.fillStyle = "#e3b341";
  ctx.beginPath();
  ctx.moveTo(x + width - 15, y + 1);
  ctx.lineTo(x + width - 2, y + 1);
  ctx.lineTo(x + width - 2, y + 14);
  ctx.closePath();
  ctx.fill();
  ctx.fillStyle = "#5b4300";
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

  ctx.fillStyle = "rgba(15, 127, 120, 0.12)";
  for (let row = range.startRow; row <= range.endRow; row++) {
    for (let col = range.startCol; col <= range.endCol; col++) {
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
  ctx.strokeStyle = "#0f7f78";
  ctx.lineWidth = 2;
  ctx.strokeRect(rect.x + 1, rect.y + 1, rect.width - 2, rect.height - 2);
  ctx.lineWidth = 1;
  ctx.restore();
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

function columnAt(x) {
  let current = rowNumberWidth - scrollX;
  for (let col = 0; col < columns.length; col++) {
    const width = columns[col].width;
    if (x >= current && x < current + width) return col;
    current += width;
  }
  return -1;
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
  try {
    const response = await fetch(`/api/column-values?id=${encodeURIComponent(documentId)}&column=${col}&limit=12`);
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
    while (record.length <= range.endCol) record.push("");
    for (let col = range.startCol; col <= range.endCol; col++) {
      record[col] = value;
      edits.push({ row: rowSource.get(row) ?? row + dataStartRow, col, value });
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
  while (record.length <= selected.col) record.push("");
  record[selected.col] = value;
  rows.set(selected.row, record);
  selectionAnchor = { ...selected };
  cellValueInputDirty = false;
  submitCellEdits(
    [{ row: rowSource.get(selected.row) ?? selected.row + dataStartRow, col: selected.col, value }],
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
  document.preview.forEach((record, index) => {
    if (index < dataStartRow) return;
    const displayRow = index - dataStartRow;
    rows.set(displayRow, record);
    rowSource.set(displayRow, index);
  });
  windowStart = 0;
  windowEnd = Math.max(0, document.preview.length - dataStartRow);
  totalRows = Math.max(0, (document.rowIndex?.rows || document.preview.length) - dataStartRow);
  columns = buildColumns(document);
  selected = { row: 0, col: 0 };
  selectionAnchor = { row: 0, col: 0 };
  scrollX = 0;
  scrollY = 0;
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
  loadDocument(document);
  dirty = payload.dirty ?? session.dirty ?? true;
  selected = {
    row: Math.max(0, Math.min(totalRows - 1, previousRow)),
    col: Math.max(0, Math.min(columns.length - 1, previousCol)),
  };
  selectionAnchor = { ...selected };
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

async function copySelection(format = "tsv") {
  if (!documentId) return;
  const matrix = selectedMatrix();
  const text = formatMatrix(matrix, format);
  if (!text) {
    status.textContent = "No cells selected to copy.";
    return;
  }
  try {
    await navigator.clipboard.writeText(text);
    status.textContent = `Copied ${formatCount(selectionCellCount())} cells as ${format.toUpperCase()}`;
  } catch (error) {
    status.textContent = `Copy failed: ${error.message}`;
  }
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
      while (record.length <= col) record.push("");
      record[col] = values[colOffset];
      edits.push({ row: rowSource.get(row) ?? row + dataStartRow, col, value: values[colOffset] });
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
  if (!documentId) return;
  const search = parseSearchQuery(commandInput.value);
  const query = search.query;
  if (!query) {
    status.textContent = "Type text to find.";
    return;
  }
  const response = await fetch("/api/find", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      documentId,
      query,
      regex: search.regex,
      caseSensitive: search.caseSensitive,
      limit: 1,
    }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Find failed: ${payload.error ?? "unknown error"}`;
    return;
  }
  const match = payload.matches?.[0];
  if (!match) {
    status.textContent = `No matches for ${query}`;
    return;
  }
  if (viewMapped) await clearView();
  const displayRow = Math.max(0, match.row - dataStartRow);
  selected = { row: displayRow, col: match.col };
  selectionAnchor = { ...selected };
  scrollY = Math.max(0, displayRow * rowHeight - rowHeight * 2);
  await maybeLoadWindow(displayRow);
  draw();
  status.textContent = `Found ${query} at row ${displayRow + 1}, column ${match.col + 1}`;
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
  return match
    ? { query: match[1], regex: true, caseSensitive: !match[2].includes("i") }
    : { query: trimmed, regex: false, caseSensitive: false };
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
  totalRows = payload.totalRows ?? 0;
  windowStart = 0;
  windowEnd = 0;
  await maybeLoadWindow(0);
  status.textContent = `Sorted ${formatSortKeys(sortKeys)}`;
}

function nextSortKeys(column, descending, append) {
  const key = sortKeyForColumn(column, descending);
  if (!append || sortKeys.length === 0) return [key];
  return [...sortKeys.filter((item) => item.column !== column), key];
}

function sortKeyForColumn(column, descending) {
  const type = columns[column]?.type ?? "string";
  return {
    column,
    descending,
    numeric: type === "integer" || type === "decimal" || columnLooksNumeric(column),
  };
}

function formatSortKeys(keys) {
  return keys
    .map((key) => `${columns[key.column]?.name ?? `column ${key.column + 1}`} ${key.descending ? "descending" : "ascending"}`)
    .join(", then ");
}

function columnLooksNumeric(column) {
  let inspected = 0;
  let numeric = 0;
  for (let row = 0; row < Math.min(totalRows, 80); row++) {
    const value = rows.get(row)?.[column];
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
    { documentId, query: value, column },
    `Filtered ${columns[column].name}`,
    column,
  );
}

async function filterColumnBlank(showBlank) {
  const column = columnMenuColumn;
  closeColumnMenu();
  if (!documentId || !columns[column]) return;
  const expression = `c${column + 1} ${showBlank ? "=" : "<>"} ''`;
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
  const response = await fetch("/api/filter-values", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ documentId, column, values, includeBlank }),
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
  rows = new Map();
  rowSource = new Map();
  totalRows = Math.max(0, (payload.totalRows ?? 0) - dataStartRow);
  windowStart = 0;
  windowEnd = 0;
  selected = { row: 0, col: selected.col };
  selectionAnchor = { ...selected };
  scrollY = 0;
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

async function insertSelectedColumn() {
  await mutateStructure("/api/column/insert", { col: selected.col }, "Inserted column");
}

async function deleteSelectedColumn() {
  await mutateStructure("/api/column/delete", { col: selected.col }, "Deleted column");
}

async function transposeCurrentDocument() {
  await mutateStructure("/api/transpose", { maxCells: 1_000_000 }, "Transposed table");
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
    if (profile.index === selected.col) tr.className = "selected-profile-column";
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
      selected = { row: selected.row, col: profile.index };
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
  if (!leftPath || !rightPath) {
    status.textContent = "Choose left and right CSV paths to compare.";
    return;
  }
  status.textContent = "Comparing files...";
  const response = await fetch("/api/diff", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ leftPath, rightPath, maxRows: 1000 }),
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = `Diff failed: ${payload.error ?? "unknown error"}`;
    diffResult.replaceChildren(emptyDiffNode(payload.error ?? "Diff failed"));
    return;
  }
  renderDiffResult(payload);
  status.textContent =
    `Diff: ${payload.summary.changed} changed, ${payload.summary.added} added, ${payload.summary.deleted} deleted`;
}

function renderDiffResult(result) {
  const summary = result.summary ?? {};
  diffSummary.replaceChildren(
    summaryNode("Unchanged", summary.unchanged ?? 0),
    summaryNode("Changed", summary.changed ?? 0),
    summaryNode("Added", summary.added ?? 0),
    summaryNode("Deleted", summary.deleted ?? 0),
  );
  if (!result.rows?.length) {
    diffResult.replaceChildren(emptyDiffNode("No diff rows"));
    return;
  }
  const table = document.createElement("table");
  const thead = document.createElement("thead");
  const tr = document.createElement("tr");
  for (const label of ["Status", "Left Row", "Left", "Right Row", "Right", "Changed Cells"]) {
    const th = document.createElement("th");
    th.textContent = label;
    tr.append(th);
  }
  thead.append(tr);
  table.append(thead);
  const tbody = document.createElement("tbody");
  for (const row of result.rows) {
    const tr = document.createElement("tr");
    tr.className = row.status;
    appendCell(tr, row.status);
    appendCell(tr, row.leftRow ?? "");
    appendCell(tr, (row.leftRecord ?? []).join(" | "));
    appendCell(tr, row.rightRow ?? "");
    appendCell(tr, (row.rightRecord ?? []).join(" | "));
    const cells = document.createElement("td");
    cells.className = "diff-cells";
    cells.textContent = (row.changedCells ?? [])
      .map((cell) => `c${cell.col}: ${cell.left} -> ${cell.right}`)
      .join("; ");
    tr.append(cells);
    tbody.append(tr);
  }
  table.append(tbody);
  diffResult.replaceChildren(table);
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
  const viewportRowsHeight = Math.max(0, canvas.clientHeight - headerHeight);
  const maxY = Math.max(0, totalRows * rowHeight - viewportRowsHeight);
  const totalWidth = columns.reduce((sum, column) => sum + column.width, 0);
  const maxX = Math.max(0, totalWidth - canvas.clientWidth + rowNumberWidth);
  scrollY = Math.max(0, Math.min(maxY, scrollY + event.deltaY));
  scrollX = Math.max(0, Math.min(maxX, scrollX + event.deltaX));
  if (editing) commitEditor();
  draw();
});

canvas.addEventListener("mousedown", (event) => {
  canvas.focus();
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  if (y < headerHeight) {
    const col = columnAt(x);
    if (x >= rowNumberWidth && col >= 0) {
      event.preventDefault();
      if (editing) commitEditor();
      openColumnMenu(col, event.clientX, rect.top + headerHeight + 4);
    }
    return;
  }
  closeColumnMenu();
  if (x < rowNumberWidth) return;
  const row = Math.floor((y - headerHeight + scrollY) / rowHeight);
  const col = columnAt(x);
  if (row >= 0 && row < totalRows && col >= 0) {
    selected = { row, col };
    if (!event.shiftKey) {
      selectionAnchor = { ...selected };
    }
    if (editing) commitEditor();
    draw();
  }
});

canvas.addEventListener("dblclick", openEditor);

window.addEventListener("copy", (event) => {
  if (!documentId || isTypingTarget(event.target)) return;
  event.preventDefault();
  event.clipboardData.setData("text/plain", formatMatrix(selectedMatrix(), "tsv"));
  status.textContent = `Copied ${formatCount(selectionCellCount())} cells`;
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

function handleGlobalShortcut(event) {
  const primary = event.metaKey || event.ctrlKey;
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
  const key = event.key.toLowerCase();
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
replaceButton.addEventListener("click", () => replaceCurrentQuery(false));
replaceAllButton.addEventListener("click", () => replaceCurrentQuery(true));
filterButton.addEventListener("click", filterCurrentQuery);
sortButton.addEventListener("click", sortSelectedColumn);
commandInput.addEventListener("keydown", async (event) => {
  if (event.key !== "Enter") return;
  event.preventDefault();
  await findCurrentQuery();
});
window.addEventListener("resize", resize);
if ("ResizeObserver" in window) {
  const gridResizeObserver = new ResizeObserver(resize);
  gridResizeObserver.observe(gridFrame);
}
canvas.tabIndex = 0;
profilePanel.tabIndex = -1;
loadCommandRegistry();
resize();
