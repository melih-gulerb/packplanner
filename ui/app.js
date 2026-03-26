const state = {
  packSizes: [],
  lastPlan: null,
};

const apiBaseUrl =
  (window.PACKPLANNER_CONFIG && window.PACKPLANNER_CONFIG.apiBaseUrl
    ? window.PACKPLANNER_CONFIG.apiBaseUrl
    : "https://packplanner.onrender.com").replace(/\/+$/, "");

const elements = {
  refreshButton: document.querySelector("#refreshButton"),
  packRows: document.querySelector("#packRows"),
  addPackButton: document.querySelector("#addPackButton"),
  savePackButton: document.querySelector("#savePackButton"),
  orderQuantity: document.querySelector("#orderQuantity"),
  calculateButton: document.querySelector("#calculateButton"),
  statusDot: document.querySelector("#statusDot"),
  statusText: document.querySelector("#statusText"),
  packSizeCount: document.querySelector("#packSizeCount"),
  lastOrderValue: document.querySelector("#lastOrderValue"),
  lastShippedValue: document.querySelector("#lastShippedValue"),
  metricTotalItems: document.querySelector("#metricTotalItems"),
  metricExtraItems: document.querySelector("#metricExtraItems"),
  metricTotalPacks: document.querySelector("#metricTotalPacks"),
  resultBadge: document.querySelector("#resultBadge"),
  emptyState: document.querySelector("#emptyState"),
  breakdown: document.querySelector("#breakdown"),
  breakdownRows: document.querySelector("#breakdownRows"),
  toast: document.querySelector("#toast"),
  quickPicks: [...document.querySelectorAll(".chip[data-order]")],
  swaggerLink: document.querySelector("#swaggerLink"),
};

function buildApiUrl(path) {
  return `${apiBaseUrl}${path}`;
}

async function requestJSON(path, options = {}) {
  const response = await fetch(buildApiUrl(path), {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  const payload = await response.json();
  if (!response.ok || payload.success === false) {
    throw new Error(payload.message || "Unexpected request failure");
  }

  return payload;
}

function showToast(message, type = "success") {
  elements.toast.textContent = message;
  elements.toast.className = `toast ${type === "error" ? "error" : ""}`.trim();
  elements.toast.hidden = false;

  window.clearTimeout(showToast.timeoutId);
  showToast.timeoutId = window.setTimeout(() => {
    elements.toast.hidden = true;
  }, 2800);
}

function setButtonBusy(button, isBusy, idleText, busyText) {
  button.disabled = isBusy;
  button.textContent = isBusy ? busyText : idleText;
}

function createPackRow(value = "") {
  const row = document.createElement("div");
  row.className = "pack-row";
  row.innerHTML = `
    <label class="pack-field">
      <span>Pack size</span>
      <input type="number" min="1" inputmode="numeric" value="${value}" />
    </label>
    <button class="icon-button" type="button" aria-label="Remove pack size">×</button>
  `;

  row.querySelector(".icon-button").addEventListener("click", () => {
    row.remove();
    ensureAtLeastOnePackRow();
  });

  return row;
}

function ensureAtLeastOnePackRow() {
  if (elements.packRows.children.length === 0) {
    elements.packRows.appendChild(createPackRow());
  }
}

function renderPackRows(packSizes) {
  elements.packRows.innerHTML = "";
  packSizes.forEach((packSize) => {
    elements.packRows.appendChild(createPackRow(packSize));
  });
  ensureAtLeastOnePackRow();
}

function collectPackSizes() {
  const values = [...elements.packRows.querySelectorAll('input[type="number"]')]
    .map((input) => Number.parseInt(input.value, 10))
    .filter((value) => Number.isFinite(value));

  if (values.length === 0) {
    throw new Error("Please provide at least one pack size");
  }

  return values;
}

function updateStatus(isOnline) {
  elements.statusDot.classList.toggle("online", isOnline);
  elements.statusDot.classList.toggle("offline", !isOnline);
  elements.statusText.textContent = isOnline ? "Connected" : "Unavailable";
}

function renderPlan(plan) {
  state.lastPlan = plan;

  elements.lastOrderValue.textContent = plan.order_quantity;
  elements.lastShippedValue.textContent = plan.total_items;
  elements.metricTotalItems.textContent = plan.total_items;
  elements.metricExtraItems.textContent = plan.total_items - plan.order_quantity;
  elements.metricTotalPacks.textContent = plan.total_packs;
  elements.resultBadge.textContent = "Calculation completed";

  elements.emptyState.hidden = true;
  elements.breakdown.hidden = false;
  elements.breakdownRows.innerHTML = "";

  plan.packs.forEach((pack) => {
    const row = document.createElement("div");
    row.className = "breakdown-row";
    row.innerHTML = `
      <div>
        <strong>${pack.pack_size} items</strong>
        <div class="breakdown-meta">${pack.pack_size * pack.quantity} shipped in this tier</div>
      </div>
      <div class="quantity-pill">${pack.quantity} pack${pack.quantity > 1 ? "s" : ""}</div>
    `;

    elements.breakdownRows.appendChild(row);
  });
}

async function loadPackSizes() {
  const payload = await requestJSON("/api/v1/pack-sizes");
  state.packSizes = payload.data.pack_sizes;
  elements.packSizeCount.textContent = state.packSizes.length;
  renderPackRows(state.packSizes);
}

async function checkHealth() {
  try {
    await requestJSON("/health");
    updateStatus(true);
  } catch {
    updateStatus(false);
  }
}

async function savePackSizes() {
  try {
    setButtonBusy(elements.savePackButton, true, "Save pack sizes", "Saving...");
    const packSizes = collectPackSizes();
    const payload = await requestJSON("/api/v1/pack-sizes", {
      method: "PUT",
      body: JSON.stringify({ pack_sizes: packSizes }),
    });

    state.packSizes = payload.data.pack_sizes;
    elements.packSizeCount.textContent = state.packSizes.length;
    renderPackRows(state.packSizes);
    showToast(payload.message);
  } catch (error) {
    showToast(error.message, "error");
  } finally {
    setButtonBusy(elements.savePackButton, false, "Save pack sizes", "Saving...");
  }
}

async function calculatePlan() {
  try {
    setButtonBusy(elements.calculateButton, true, "Calculate plan", "Calculating...");
    const orderQuantity = Number.parseInt(elements.orderQuantity.value, 10);
    if (!Number.isFinite(orderQuantity) || orderQuantity <= 0) {
      throw new Error("Please enter a valid order quantity");
    }

    const payload = await requestJSON("/api/v1/pack-plans", {
      method: "POST",
      body: JSON.stringify({ order_quantity: orderQuantity }),
    });

    renderPlan(payload.data);
    showToast(payload.message);
  } catch (error) {
    showToast(error.message, "error");
  } finally {
    setButtonBusy(elements.calculateButton, false, "Calculate plan", "Calculating...");
  }
}

function bindEvents() {
  elements.swaggerLink.href = buildApiUrl("/swagger");

  elements.refreshButton.addEventListener("click", async () => {
    try {
      setButtonBusy(elements.refreshButton, true, "Refresh state", "Refreshing...");
      await Promise.all([checkHealth(), loadPackSizes()]);
      showToast("UI state refreshed");
    } catch (error) {
      showToast(error.message, "error");
    } finally {
      setButtonBusy(elements.refreshButton, false, "Refresh state", "Refreshing...");
    }
  });

  elements.addPackButton.addEventListener("click", () => {
    elements.packRows.appendChild(createPackRow());
  });

  elements.savePackButton.addEventListener("click", savePackSizes);
  elements.calculateButton.addEventListener("click", calculatePlan);

  elements.orderQuantity.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
      event.preventDefault();
      calculatePlan();
    }
  });

  elements.quickPicks.forEach((chip) => {
    chip.addEventListener("click", () => {
      elements.orderQuantity.value = chip.dataset.order;
      calculatePlan();
    });
  });
}

async function bootstrap() {
  bindEvents();

  try {
    await Promise.all([checkHealth(), loadPackSizes()]);
  } catch (error) {
    showToast(error.message, "error");
  }
}

bootstrap();
