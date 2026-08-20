"use strict";

const $ = (id) => document.getElementById(id);

function setStatus(msg, ok) {
  const el = $("status");
  el.textContent = msg;
  el.className = "status " + (ok ? "ok" : "bad");
}

function loadExample() {
  fetch("/api/example")
    .then((r) => r.json())
    .then((doc) => {
      $("input").value = JSON.stringify(doc, null, 2);
      setStatus("已加载示例，点击计算。", true);
    })
    .catch((e) => setStatus("加载示例失败: " + e, false));
}

function solve() {
  const text = $("input").value;
  let doc;
  try {
    doc = JSON.parse(text);
  } catch (e) {
    setStatus("JSON 解析失败: " + e, false);
    return;
  }
  fetch("/api/solve", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(doc),
  })
    .then((r) => r.json().then((j) => ({ ok: r.ok, j })))
    .then(({ ok, j }) => {
      if (!ok || j.error) {
        setStatus("求解失败: " + (j.error || "未知错误"), false);
        return;
      }
      renderReactions(j.reactions || []);
      renderExtrema(j.extrema || {});
      renderChecks(j.checks || {});
      drawDiagram(j);
      setStatus(
        "求解成功 · 支座 " + j.support + " · 采样 " + (j.samples || []).length + " 点",
        true
      );
    })
    .catch((e) => setStatus("请求失败: " + e, false));
}

function renderReactions(reactions) {
  const tb = $("reactions").querySelector("tbody");
  tb.innerHTML = "";
  reactions.forEach((r) => {
    const tr = document.createElement("tr");
    tr.innerHTML =
      "<td>" + r.at.toFixed(3) + "</td><td>" + r.kind + "</td><td>" +
      r.force.toFixed(4) + "</td>";
    tb.appendChild(tr);
  });
}

function renderExtrema(e) {
  const tb = $("extrema").querySelector("tbody");
  tb.innerHTML = "";
  const rows = [
    ["V max", e.v_max],
    ["V min", e.v_min],
    ["M max", e.m_max],
    ["M min", e.m_min],
    ["y max", e.y_max],
    ["y min", e.y_min],
    ["|y| max", e.y_abs_max],
    ["|y| max @x", e.y_abs_max_x],
  ];
  rows.forEach(([k, v]) => {
    const tr = document.createElement("tr");
    const val = typeof v === "number" ? v.toFixed(4) : v;
    tr.innerHTML = "<td>" + k + "</td><td></td><td>" + val + "</td>";
    tb.appendChild(tr);
  });
}

function renderChecks(c) {
  const box = $("checks");
  box.innerHTML = "";
  (c.items || []).forEach((it) => {
    const div = document.createElement("div");
    div.className = "check-item " + (it.passed ? "pass" : "fail");
    div.textContent = (it.passed ? "✓ " : "✗ ") + it.name + " — " + it.message;
    box.appendChild(div);
  });
}

function drawDiagram(j) {
  const svg = $("diagram");
  const W = 720, H = 360;
  const padL = 50, padR = 30, padT = 30, padB = 60;
  const plotW = W - padL - padR;
  const axisY = padT + (H - padT - padB) / 2;
  const plotH = (H - padT - padB) / 2;
  const L = j.L || 1;
  const samples = j.samples || [];

  let maxAbsM = 0;
  samples.forEach((s) => { maxAbsM = Math.max(maxAbsM, Math.abs(s.M)); });
  if (maxAbsM < 1e-9) maxAbsM = 1;
  const scale = (plotH - 10) / maxAbsM;

  const xMap = (x) => padL + (x / L) * plotW;
  const yMap = (m) => axisY - m * scale;

  let svgParts = "";

  // beam axis
  svgParts += line(padL, axisY, padL + plotW, axisY, "#8b98a9", 2);
  // supports: triangle at x=0 and x=L
  svgParts += support(padL, axisY, "left");
  svgParts += support(padL + plotW, axisY, "right");

  // moment diagram as a filled polyline following M(x) from /api/solve
  if (samples.length > 1) {
    let pts = "M " + xMap(samples[0].X) + " " + yMap(samples[0].M);
    samples.forEach((s) => {
      pts += " L " + xMap(s.X).toFixed(2) + " " + yMap(s.M).toFixed(2);
    });
    pts += " L " + xMap(samples[samples.length - 1].X).toFixed(2) + " " + axisY;
    pts += " L " + xMap(samples[0].X).toFixed(2) + " " + axisY + " Z";
    svgParts += '<path d="' + pts + '" fill="rgba(78,161,255,0.25)" stroke="#4ea1ff" stroke-width="1.5"/>';
  }

  // axis labels
  svgParts += text(padL, H - 18, "x = 0", "#8b98a9", 11);
  svgParts += text(padL + plotW - 40, H - 18, "x = " + L.toFixed(2), "#8b98a9", 11);
  svgParts += text(padL, padT - 12, "M(x) 弯矩 (sagging+)", "#8b98a9", 11);

  svg.innerHTML = svgParts;
}

function line(x1, y1, x2, y2, color, w) {
  return '<line x1="' + x1 + '" y1="' + y1 + '" x2="' + x2 + '" y2="' + y2 +
    '" stroke="' + color + '" stroke-width="' + w + '"/>';
}

function text(x, y, s, color, size) {
  return '<text x="' + x + '" y="' + y + '" fill="' + color + '" font-size="' + size + '">' + s + "</text>";
}

function support(x, y, side) {
  const d = 12;
  if (side === "left") {
    return '<path d="M ' + x + ' ' + y + ' L ' + (x - d) + ' ' + (y + d) + ' L ' +
      (x + d) + ' ' + (y + d) + ' Z" fill="#3fb950"/>';
  }
  return '<path d="M ' + x + ' ' + y + ' L ' + (x - d) + ' ' + (y + d) + ' L ' +
    (x + d) + ' ' + (y + d) + ' Z" fill="#3fb950"/>';
}

window.addEventListener("DOMContentLoaded", () => {
  $("loadExample").addEventListener("click", loadExample);
  $("solve").addEventListener("click", solve);
  loadExample();
});
