let user = JSON.parse(localStorage.getItem("user"));
const categories = ["Food", "Gas", "Entertainment", "Shopping", "Subscriptions", "Transfers", "Alcohol", "Other"];
const categoryMap = {
    Food: 1,
    Gas: 2,
    Entertainment: 3,
    Shopping: 4,
    Subscriptions: 5,
    Transfers: 6,
    Alcohol: 7,
    Other: 9,
    total: 0
  };

function initialLoad(){
    user = JSON.parse(localStorage.getItem("user"));
    if (!user) location.href = "/index.html";
    document.getElementById("user-info").innerText = `Logged in as ${user.username} (${user.email})`;
    console.log("User object in localStorage:", user);
}

async function refreshBudgets(){
    try {
        const res = await fetch(`http://localhost:8080/budgets/${user.google_id}`);
        return await res.json();
      } catch (err) {
        console.error("Fetch failed:", err);
        return null;
      }
}

function showBudgetModal() {
  document.getElementById("budget-modal").style.display = "flex";
}

function hideBudgetModal() {
  document.getElementById("budget-modal").style.display = "none";
}

function setupBudgetModal() {
  const typeRadios = document.getElementsByName("budget-type");
  const totalInput = document.getElementById("total-budget-input");
  const categoryInputs = document.getElementById("category-budget-inputs");
  const categoryList = document.getElementById("category-input-list");

  // Populate category fields
  categoryList.innerHTML = "";
  categories.forEach(cat => {
    const wrapper = document.createElement("div");
    wrapper.innerHTML = `
      <label>${cat}:</label>
      <input type="number" min="0" name="category-${cat}" />
    `;
    categoryList.appendChild(wrapper);
  });

  // Toggle input visibility
  typeRadios.forEach(radio => {
    radio.addEventListener("change", () => {
      if (radio.value === "total" && radio.checked) {
        totalInput.style.display = "block";
        categoryInputs.style.display = "none";
      } else if (radio.value === "category" && radio.checked) {
        totalInput.style.display = "none";
        categoryInputs.style.display = "block";
      }
    });
  });

  // Hook up buttons
  document.getElementById("cancel-budget").onclick = hideBudgetModal;

  document.getElementById("save-budget").onclick = () => {
    const name = document.getElementById("budget-name").value.trim();
    if (!name) {
        alert("Please enter a budget name.");
        return;
    }
    const type = [...typeRadios].find(r => r.checked).value;
  
    if (type === "total") {
      const total = parseFloat(document.getElementById("total-budget").value);
      if (isNaN(total) || total <= 0) {
        alert("Please enter a valid total.");
        return;
      }
  
      // Call submitBudget with category_id = 0 for "By Total"
      submitBudget(name, 0, total);
      hideBudgetModal();
  
    } else if (type === "category") {
        let submittedAny = false;

        categories.forEach(cat => {
          const input = document.querySelector(`input[name='category-${cat}']`);
          const value = parseFloat(input.value);
    
          if (!isNaN(value) && value > 0) {
            const categoryId = categoryMap[cat];
            submitBudget(name, categoryId, value);
            submittedAny = true;
          }
        });
    
        if (!submittedAny) {
          alert("Please enter at least one positive amount for a category.");
          return;
        }
    
        hideBudgetModal();
    }
  };
}


async function submitBudget(name, categoryId, amount) {
    const payload = {
      user_id: user.google_id,      
      name: name,
      category_id: categoryId,
      amount: amount,
    };
  
    try {
      const res = await fetch("http://localhost:8080/budgets", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
      });
  
      if (!res.ok) {
        const errMsg = await res.text();
        throw new Error(errMsg);
      }
  
      const data = await res.json();
      console.log("Budget created:", data);
    } catch (err) {
      console.error("Failed to create budget:", err.message);
      alert("Could not create budget.");
    }
  }



function EventListeners(){
    document.getElementById("create-budget-btn").addEventListener("click", () => {
        setupBudgetModal();
        showBudgetModal();
      });
      
}
function ApplyMonthBudgets(budgets){
    document.getElementById("apply-month").addEventListener("click", async () => {
        displayBudgets(budgets);
      });
}
function getMonthDateRange(monthValue) {
    // monthValue is in format "YYYY-MM"
    const [year, month] = monthValue.split("-").map(Number);
  
    const startDate = new Date(year, month - 1, 1); // 1st of the month
    const endDate = new Date(year, month, 0); // last day of the month
  
    // Format as YYYY-MM-DD
    const format = (d) => d.toISOString().split("T")[0];
  
    return {
      start: format(startDate),
      end: format(endDate),
    };
  }

async function displayBudgets(budgets){
    const monthValue = document.getElementById("month-picker").value;
    if (!monthValue) {
      alert("Please select a month.");
      return;
    }
  
    const { start, end } = getMonthDateRange(monthValue);
  
    try {
      const res = await fetch(`http://localhost:8080/transactions/${user.google_id}?start_date=${start}&end_date=${end}`);
      const transactions = await res.json();
      populateBudgets(budgets, getTransactionAmounts(transactions));
    } catch (err) {
      console.error("Failed to fetch budgets:", err);
      alert("Could not load budgets for selected month.");
    }
}

function getTransactionAmounts(transactions){
    let transactionMap = new Map();
    transactions.forEach((transaction) => {
        const catID = categoryMap[transaction.T_type];
        const currValue = transactionMap.get(catID);
        if(currValue){
            transactionMap.set(catID, currValue + Math.abs(transaction.Transaction.Amount));
        } else {
            transactionMap.set(catID, Math.abs(transaction.Transaction.Amount));
        }
        if(transactionMap.has(0)){
            transactionMap.set(0, transactionMap.get(0) + Math.abs(transaction.Transaction.Amount));
        } else {
            transactionMap.set(0, Math.abs(transaction.Transaction.Amount));
        }
    });
    console.log(transactionMap);
    Object.keys(categoryMap).forEach((key) => {
        if (!transactionMap.has(categoryMap[key])) {
          transactionMap.set(categoryMap[key], 0);
        }
      });
    return transactionMap;
}

function populateBudgets(budgets, transactionMap) {
    const container = document.getElementById("budgets-list");
    container.innerHTML = ""; // clear existing content
  
    const categoryDesc = new Map();
  
    budgets.forEach((budget) => {
      const key = `${budget.budget_id}-${budget.name}`;
      if (categoryDesc.has(key)) {
        categoryDesc.get(key).push([budget.category_id, budget.amount]);
      } else {
        categoryDesc.set(key, [[budget.category_id, budget.amount]]);
      }
    });
  
    // Map of category IDs to readable names
    const reverseCategoryMap = Object.entries(categoryMap).reduce((acc, [name, id]) => {
      acc[id] = name;
      return acc;
    }, {});
  
    categoryDesc.forEach((categoryItems, key) => {
      const [budgetId, name] = key.split("-");
  
      const budgetItem = document.createElement("li");
      const title = document.createElement("strong");
      title.textContent = name;
      budgetItem.appendChild(title);
  
      const innerList = document.createElement("ul");
      categoryItems.forEach(([catID, amt]) => {
        const categoryName = reverseCategoryMap[parseInt(catID)] || "Unknown";
        const catLi = document.createElement("li");
        catLi.textContent = `${categoryName}: $${transactionMap.get(catID).toFixed(2)}/$${amt.toFixed(2)}`;
        innerList.appendChild(catLi);
      });
  
      budgetItem.appendChild(innerList);
      container.appendChild(budgetItem);
    });
  }

async function populateGoals(){
    const res = await fetch(`http://localhost:8080/goals/${user.google_id}`);
    const goals = await res.json();
    console.log(goals);
}


async function main(){
    initialLoad();
    
    let budgets = await refreshBudgets();
    if(budgets){
        ApplyMonthBudgets(budgets);
        console.log("Displaying...");
    }
    populateGoals();
    EventListeners();
}

main();