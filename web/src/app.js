const canvas = document.querySelector("#grid");
const ctx = canvas.getContext("2d");
const editor = document.querySelector("#editor");
const status = document.querySelector("#status");
const columnName = document.querySelector("#columnName");

const rowHeight = 28;
const headerHeight = 32;
const rowNumberWidth = 58;
const totalRows = 120000;
const columns = [
  { name: "id", width: 92, type: "integer" },
  { name: "segment", width: 140, type: "string" },
  { name: "amount", width: 120, type: "decimal" },
  { name: "created_at", width: 140, type: "date" },
  { name: "status", width: 132, type: "string" },
  { name: "notes", width: 300, type: "string" },
  { name: "region", width: 140, type: "string" },
  { name: "score", width: 110, type: "decimal" },
];

let scrollX = 0;
let scrollY = 0;
let selected = { row: 0, col: 0 };
let editing = false;

function resize() {
  const rect = canvas.getBoundingClientRect();
  const scale = window.devicePixelRatio || 1;
  canvas.width = Math.floor(rect.width * scale);
  canvas.height = Math.floor(rect.height * scale);
  ctx.setTransform(scale, 0, 0, scale, 0, 0);
  draw();
}

function valueAt(row, col) {
  switch (columns[col].name) {
    case "id":
      return String(row + 1);
    case "segment":
      return `segment_${row % 12}`;
    case "amount":
      return ((row + 1) * 19.73).toFixed(2);
    case "created_at":
      return `2026-05-${String((row % 28) + 1).padStart(2, "0")}`;
    case "status":
      return ["new", "active", "review", "archived"][row % 4];
    case "notes":
      return `editable generated row ${row + 1}`;
    case "region":
      return ["APAC", "NA", "EU", "LATAM"][row % 4];
    case "score":
      return ((row * 7) % 1000 / 10).toFixed(1);
    default:
      return "";
  }
}

function draw() {
  const width = canvas.clientWidth;
  const height = canvas.clientHeight;
  ctx.clearRect(0, 0, width, height);

  ctx.fillStyle = "#fbfcfc";
  ctx.fillRect(0, 0, width, height);

  drawHeaders(width);
  drawRows(width, height);
  drawSelection();

  status.textContent =
    `${totalRows.toLocaleString()} rows, ${columns.length} columns - ` +
    `row ${selected.row + 1}, ${columns[selected.col].name}`;
  columnName.textContent = columns[selected.col].name;
}

function drawHeaders(width) {
  ctx.fillStyle = "#eef2f3";
  ctx.fillRect(0, 0, width, headerHeight);
  ctx.fillStyle = "#f5f7f7";
  ctx.fillRect(0, 0, rowNumberWidth, headerHeight);
  ctx.strokeStyle = "#d8e0e2";
  ctx.beginPath();
  ctx.moveTo(0, headerHeight + 0.5);
  ctx.lineTo(width, headerHeight + 0.5);
  ctx.stroke();

  ctx.font = "600 12px Inter, system-ui, sans-serif";
  ctx.textBaseline = "middle";

  let x = rowNumberWidth - scrollX;
  for (const column of columns) {
    if (x + column.width >= rowNumberWidth && x <= width) {
      ctx.fillStyle = "#162022";
      ctx.fillText(column.name, x + 8, headerHeight / 2);
      ctx.fillStyle = "#697579";
      ctx.fillText(column.type, x + column.width - 56, headerHeight / 2);
      ctx.strokeStyle = "#d8e0e2";
      ctx.beginPath();
      ctx.moveTo(x + column.width + 0.5, 0);
      ctx.lineTo(x + column.width + 0.5, canvas.clientHeight);
      ctx.stroke();
    }
    x += column.width;
  }
}

function drawRows(width, height) {
  const firstRow = Math.max(0, Math.floor(scrollY / rowHeight));
  const yOffset = scrollY % rowHeight;
  const visibleRows = Math.ceil((height - headerHeight) / rowHeight) + 1;

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
        ctx.fillStyle = "#162022";
        const value = valueAt(row, col);
        ctx.save();
        ctx.beginPath();
        ctx.rect(x + 1, y + 1, column.width - 2, rowHeight - 2);
        ctx.clip();
        ctx.fillText(value, x + 8, y + rowHeight / 2);
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
}

function drawSelection() {
  const rect = cellRect(selected.row, selected.col);
  if (!rect) return;
  ctx.strokeStyle = "#0f7f78";
  ctx.lineWidth = 2;
  ctx.strokeRect(rect.x + 1, rect.y + 1, rect.width - 2, rect.height - 2);
  ctx.lineWidth = 1;
}

function cellRect(row, col) {
  const y = headerHeight + row * rowHeight - scrollY;
  let x = rowNumberWidth - scrollX;
  for (let i = 0; i < col; i++) x += columns[i].width;
  const width = columns[col].width;
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

function commitEditor() {
  editing = false;
  editor.style.display = "none";
  canvas.focus();
}

function openEditor() {
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

canvas.addEventListener("wheel", (event) => {
  event.preventDefault();
  const maxY = Math.max(0, totalRows * rowHeight - canvas.clientHeight);
  const totalWidth = columns.reduce((sum, column) => sum + column.width, 0);
  const maxX = Math.max(0, totalWidth - canvas.clientWidth + rowNumberWidth);
  scrollY = Math.max(0, Math.min(maxY, scrollY + event.deltaY));
  scrollX = Math.max(0, Math.min(maxX, scrollX + event.deltaX));
  if (editing) commitEditor();
  draw();
});

canvas.addEventListener("mousedown", (event) => {
  const rect = canvas.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  if (y < headerHeight || x < rowNumberWidth) return;
  const row = Math.floor((y - headerHeight + scrollY) / rowHeight);
  const col = columnAt(x);
  if (row >= 0 && row < totalRows && col >= 0) {
    selected = { row, col };
    if (editing) commitEditor();
    draw();
  }
});

canvas.addEventListener("dblclick", openEditor);

window.addEventListener("keydown", (event) => {
  if (editing) return;
  if (event.key === "Enter") {
    event.preventDefault();
    openEditor();
    return;
  }
  if (event.key === "ArrowDown") selected.row = Math.min(totalRows - 1, selected.row + 1);
  if (event.key === "ArrowUp") selected.row = Math.max(0, selected.row - 1);
  if (event.key === "ArrowRight") selected.col = Math.min(columns.length - 1, selected.col + 1);
  if (event.key === "ArrowLeft") selected.col = Math.max(0, selected.col - 1);
  draw();
});

editor.addEventListener("keydown", (event) => {
  if (event.key === "Enter") {
    event.preventDefault();
    commitEditor();
  }
  if (event.key === "Escape") {
    event.preventDefault();
    commitEditor();
  }
});

window.addEventListener("resize", resize);
canvas.tabIndex = 0;
resize();
