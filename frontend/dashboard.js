// Ensure user is authenticated
const user = JSON.parse(localStorage.getItem("user"));
if (!user) location.href = "/index.html";

let chart = null;

// ========== FETCH ========== //
async function fetchTransactions(start, end) {
  try {
    const res = await fetch(`http://localhost:8080/transactions/${user.google_id}?start_date=${start}&end_date=${end}`);
    return await res.json();
  } catch (err) {
    console.error("Fetch failed:", err);
    return null;
  }
}

async function fetchTotals(start, end) {
  try {
    const res = await fetch(`http://localhost:8080/categorytotals/${user.google_id}?start_date=${start}&end_date=${end}`);
    return await res.json();
  } catch (err) {
    console.error("Fetch failed:", err);
    return null;
  }
}

// ========== DATA TRANSFORMATION ========== //
function buildRunningTotals(transactions) {
  if (!transactions) return null;

  const categories = ["Food", "Gas", "Entertainment", "Shopping", "Subscriptions", "Transfers", "Alcohol", "Other"];
  const categorySums = {};
  const dateSet = new Set();

  categories.forEach(cat => categorySums[cat] = {});

  transactions.forEach(tx => {
    const { Posting_date, Amount } = tx.Transaction;
    const category = tx.T_type;
    if (Amount > 0) return;

    dateSet.add(Posting_date);
    categorySums[category] ||= {};
    categorySums[category][Posting_date] = (categorySums[category][Posting_date] || 0) + Amount;
  });

  const sortedDates = Array.from(dateSet).sort();

  const datasets = categories.map(category => {
    let runningTotal = 0;
    const data = sortedDates.map(date => {
      runningTotal += categorySums[category][date] || 0;
      return Math.abs(runningTotal).toFixed(2);
    });
    return { label: category, data, borderWidth: 2, fill: false };
  });

  return { labels: sortedDates, datasets };
}

function formatToMonthName(dateStr) {
  return new Date(dateStr).toLocaleString('default', { month: 'long', year: 'numeric' });
}

function buildMonthlyTotals(transactions) {
  const monthlyTotals = {};
  transactions.forEach(tx => {
    const { Posting_date, Amount } = tx.Transaction;
    if (Amount > 0) return;
    const key = formatToMonthName(Posting_date);
    monthlyTotals[key] = (monthlyTotals[key] || 0) + Amount;
  });
  return monthlyTotals;
}

// ========== RENDER UI ========== //
function renderMonthlyTotals(totals) {
  const list = document.getElementById("monthly-totals");
  list.innerHTML = "";
  const sorted = Object.keys(totals).sort((a, b) => new Date(a) - new Date(b));

  sorted.forEach(month => {
    const li = document.createElement("li");
    li.textContent = `${month}: $${Math.abs(totals[month]).toFixed(2)}`;
    list.appendChild(li);
  });
}

function createChart(transactions) {
  if (chart) chart.destroy();

  const { labels, datasets } = buildRunningTotals(transactions);
  const ctx = document.getElementById('spending-chart').getContext('2d');

  chart = new Chart(ctx, {
    type: 'line',
    data: { labels, datasets },
    options: {
      responsive: true,
      scales: { y: { beginAtZero: true } }
    }
  });
}

function renderCategoryTotals(cat_totals) {
  const items = ["Food", "Gas", "Entertainment", "Shopping", "Subscriptions", "Transfers", "Alcohol", "Other"];
  let itemsPrice = {};
  let totalSpending = 0;
  const list = document.createElement("ul");

  cat_totals.forEach(tx => {
    const Amount = tx.total;
    const category = tx.name;

    totalSpending += Amount;
    itemsPrice[category] = Amount;
  });

  document.getElementById("total-spending").textContent = Math.abs(totalSpending.toFixed(2));

  items.forEach(item => {
    const li = document.createElement("li");
    const amount = Math.abs((itemsPrice[item] || 0).toFixed(2));
    li.textContent = `${item}: $${amount}`;
    list.appendChild(li);
  });

  const container = document.getElementById("list-container");
  container.innerHTML = "";
  container.appendChild(list);
}

// ========== EVENT HANDLERS ========== //
function setupUploadListener() {
  document.getElementById("csv-form").addEventListener("submit", function (e) {
    e.preventDefault();

    const fileInput = document.getElementById("csv-file");
    const file = fileInput.files[0];
    if (!file) return alert("Please choose a file.");

    const formData = new FormData();
    formData.append("file", file);
    formData.append("google_id", user.google_id);

    fetch("http://localhost:8080/upload", {
      method: "POST",
      body: formData
    })
      .then(res => {
        if (res.ok) return res.text();
        throw new Error("Upload failed");
      })
      .then(msg => {
        alert("CSV uploaded successfully!");
        console.log(msg);
        refreshUI(); // refresh the dashboard with new data
      })
      .catch(err => {
        console.error("Upload error:", err);
        alert("Failed to upload CSV.");
      });
  });
}

function setDefaultDateRange() {
  const endDate = new Date();
  const startDate = new Date();
  startDate.setMonth(startDate.getMonth() - 3);

  // Format as YYYY-MM-DD
  const formatDate = (date) => date.toISOString().split('T')[0];

  document.getElementById("start-date").value = formatDate(startDate);
  document.getElementById("end-date").value = formatDate(endDate);
}

function setupLogoutListener() {
  document.getElementById("logout").addEventListener("click", () => {
    localStorage.removeItem("user");
    window.location.href = "index.html";
  });
}
function setupDateListener(){
  document.getElementById("apply-date-range").addEventListener("click", async () => {
    refreshUI();
  });
}

// ========== MAIN UI LOAD ========== //
async function refreshUI() {
  const transactions = await fetchTransactions(document.getElementById("start-date").value, document.getElementById("end-date").value);
  const cat_totals = await fetchTotals(document.getElementById("start-date").value, document.getElementById("end-date").value);

  if (!transactions || transactions.length == 0) {
    alert("No transactions found. Please upload a CSV file to continue or pick different date range");
    document.getElementById("csv-form").style.display = "block";
    return;
  }

  document.getElementById("csv-form").style.display = "none";
  createChart(transactions);
  renderMonthlyTotals(buildMonthlyTotals(transactions));
  
  renderCategoryTotals(cat_totals);
}

async function setup() {
  document.getElementById("user-info").innerText = `Logged in as ${user.username} (${user.email})`;
  console.log("User object in localStorage:", user);
  setDefaultDateRange();
  setupDateListener();
  setupUploadListener();
  setupLogoutListener();
  await refreshUI();
}

// Init
setup();
