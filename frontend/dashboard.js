const user = JSON.parse(localStorage.getItem("user"));
    if (!user) {
        location.href = "/index.html";
    }
async function fetchTransactions() {
    try {
      const res = await fetch("http://localhost:8080/transactions/" + user.google_id);
      const data = await res.json();
      return await data;
    } catch (err) {
      console.error("Fetch failed:", err);
    }
  }
  
  function buildRunningTotals(transactions) {
    const categories = ["Food", "Gas", "Entertainment", "Shopping", "Subscriptions", "Transfers", "Alcohol", "Other"];
    const categorySums = {};
    const dateSet = new Set();
  
    // Initialize data structures
    categories.forEach(cat => {
      categorySums[cat] = {};
    });
  
    // Group amounts by date + category
    transactions.forEach(tx => {
      const date = tx.Transaction.Posting_date;
      const category = tx.T_type;
      const amount = tx.Transaction.Amount;
  
      if (amount > 0) return; // only spending
  
      dateSet.add(date);
  
      if (!categorySums[category]) categorySums[category] = {};
      categorySums[category][date] = (categorySums[category][date] || 0) + amount;
    });
  
    const sortedDates = Array.from(dateSet).sort(); // ensure dates are in order
  
    const datasets = categories.map(category => {
      let runningTotal = 0;
      const data = sortedDates.map(date => {
        const val = categorySums[category][date] || 0;
        runningTotal += val;
        return Math.abs(runningTotal).toFixed(2); // Use abs to show positive
      });
  
      return {
        label: category,
        data,
        borderWidth: 2,
        fill: false
      };
    });
  
    return { labels: sortedDates, datasets };
  }

  function formatToMonthName(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleString('default', { month: 'long', year: 'numeric' });
  }
  
  function buildMonthlyTotals(transactions) {
    const monthlyTotals = {};
  
    transactions.forEach(tx => {
      const amount = tx.Transaction.Amount;
      if (amount > 0) return;
  
      const key = formatToMonthName(tx.Transaction.Posting_date); // e.g., "March 2025"
      monthlyTotals[key] = (monthlyTotals[key] || 0) + amount;
    });
  
    return monthlyTotals;
  }
  
  function renderMonthlyTotals(totals) {
    const list = document.getElementById("monthly-totals");
    list.innerHTML = ""; // clear old list
  
    // Sort months by date
    const sorted = Object.keys(totals).sort((a, b) => new Date(a) - new Date(b));
  
    sorted.forEach(month => {
      const li = document.createElement("li");
      li.textContent = `${month}: $${Math.abs(totals[month]).toFixed(2)}`;
      list.appendChild(li);
    });
  }
  


async function main(){
    let transactions;
    document.getElementById("user-info").innerText = `Logged in as ${user.username} (${user.email})`;
    console.log("User object in localStorage:", user);

    transactions = await fetchTransactions();

    const { labels, datasets } = buildRunningTotals(transactions);

    const ctx = document.getElementById('spending-chart').getContext('2d');
    const chart = new Chart(ctx, {
      type: 'line',
      data: {
        labels,
        datasets
      },
      options: {
        responsive: true,
        scales: {
          y: {
            beginAtZero: true
          }
        }
      }
    });

    console.log("User transactions:", transactions);
    const items = ["Food", "Gas", "Entertainment", "Shopping", "Subscriptions", "Transfers", "Alcohol", "Other"];
    let itemsPrice = {}
    let total_transactions = 0;
    const list = document.createElement("ul");
    transactions.forEach(transaction => {
        if (transaction.Transaction.Amount > 0){
            return;
        }
        total_transactions += transaction.Transaction.Amount;
        if (transaction.T_type in itemsPrice) {
            itemsPrice[transaction.T_type] += transaction.Transaction.Amount;
        } else {
            itemsPrice[transaction.T_type] = transaction.Transaction.Amount;
        }
    });
    document.getElementById("total-spending").textContent = Math.abs(total_transactions.toFixed(2));
    const monthlyTotals = buildMonthlyTotals(transactions);
    renderMonthlyTotals(monthlyTotals);
    
    items.forEach(item => {
    const li = document.createElement("li");
    let amount;
    if(item in itemsPrice){
        amount = Math.abs(itemsPrice[item].toFixed(2));
    } else {
        amount = 0.00;
    }
    li.textContent = item + ": $" + amount;
    list.appendChild(li);
    });

    document.getElementById("list-container").appendChild(list);



    document.getElementById("csv-form").addEventListener("submit", function (e) {
        e.preventDefault();
        const fileInput = document.getElementById("csv-file");
        const file = fileInput.files[0];
        if (!file) {
          alert("Please choose a file.");
          return;
        }
      
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
          })
          .catch(err => {
            console.error("Upload error:", err);
            alert("Failed to upload CSV.");
          });
      });
      // Logout
    document.getElementById("logout").addEventListener("click", () => {
        localStorage.removeItem("user");
        window.location.href = "index.html";
    });
}
function updateItems(){
    const items = ["Food", "Gas", "Entertainment", "Shopping", "Subscriptions", "Transfers", "Alcohol", "Income", "Other"];

    const list = document.createElement("ul");

    items.forEach(item => {
    const li = document.createElement("li");
    li.textContent = item;
    list.appendChild(li);
    });

    document.getElementById("list-container").appendChild(list);
}


main();