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

let currentGoalId = null;

function showAddToGoalModal(goalId) {
  currentGoalId = goalId;
  document.getElementById("goal-add-amount").value = "";
  document.getElementById("goal-modal").style.display = "flex";
}

function hideAddToGoalModal() {
  currentGoalId = null;
  document.getElementById("goal-modal").style.display = "none";
}

document.getElementById("confirm-add-to-goal").addEventListener("click", async () => {
  const amount = parseFloat(document.getElementById("goal-add-amount").value);
  if (isNaN(amount) || amount <= 0) {
    alert("Please enter a valid amount.");
    return;
  }

  try {
    const res = await fetch("http://localhost:8080/goals/add", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        goal_id: parseInt(currentGoalId),
        amount: amount,
      }),
    });

    if (!res.ok) throw new Error(await res.text());

    hideAddToGoalModal();
    const updatedGoals = await fetchGoals();
    populateGoals(updatedGoals);
  } catch (err) {
    console.error("Add to goal failed:", err);
    alert("Could not add to goal.");
  }
});

// Modal controls
function showCreateGoalModal() {
    document.getElementById("goal-name").value = "";
    document.getElementById("goal-amount").value = "";
    document.getElementById("create-goal-modal").style.display = "flex";
  }
  
  function hideCreateGoalModal() {
    document.getElementById("create-goal-modal").style.display = "none";
  }
  
  // POST request to create a new goal
  async function submitGoal(name, amount) {
    const payload = {
      user_id: user.google_id,
      name: name,
      amount: amount
    };
  
    try {
      const res = await fetch("http://localhost:8080/goals", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
  
      if (!res.ok) throw new Error(await res.text());
  
      const newGoals = await fetchGoals();
      populateGoals(newGoals);
      hideCreateGoalModal();
    } catch (err) {
      console.error("Failed to create goal:", err);
      alert("Could not create goal.");
    }
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

async function setupBudgetModal() {
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
    console.log(JSON.stringify(payload));
  
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

      const updatedBudgets = await refreshBudgets();
        await displayBudgets(updatedBudgets);
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
    document.getElementById("create-goal-btn").addEventListener("click", () => {
        showCreateGoalModal();
    });
    
    document.getElementById("save-goal").addEventListener("click", () => {
        const name = document.getElementById("goal-name").value.trim();
        const amount = parseFloat(document.getElementById("goal-amount").value);
        
        if (!name || isNaN(amount) || amount <= 0) {
            alert("Please enter a valid name and amount.");
            return;
        }
        
        submitGoal(name, amount);
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

  async function fetchGoals() {
    try {
      const res = await fetch(`http://localhost:8080/goals/${user.google_id}`);
      if (!res.ok) throw new Error("Failed to fetch goals");
      return await res.json();
    } catch (err) {
      console.error("Fetch goals error:", err);
      return [];
    }
  }
  
  function populateGoals(goals) {
    const container = document.getElementById("goals-list");
    container.innerHTML = "";
  
    if (goals.length === 0) {
      const li = document.createElement("li");
      li.textContent = "No goals found.";
      container.appendChild(li);
      return;
    }
  
    goals.forEach(goal => {
      const li = document.createElement("li");
      li.innerHTML = `
        <strong>${goal.name}</strong> — $${goal.current_amount.toFixed(2)} / $${goal.amount.toFixed(2)}
        <button class="add-to-goal-btn" data-goal-id="${goal.goal_id}">Add</button>
        <button class="delete-goal-btn" data-goal-id="${goal.goal_id}">Delete</button>

      `;
      container.appendChild(li);
    });
  
    // Add event listeners for "Add" buttons
    document.querySelectorAll(".add-to-goal-btn").forEach(btn => {
      btn.addEventListener("click", () => {
        const goalId = btn.getAttribute("data-goal-id");
        showAddToGoalModal(goalId);
      });
    });

    document.querySelectorAll(".delete-goal-btn").forEach(btn => {
        btn.addEventListener("click", async () => {
          const goalId = btn.getAttribute("data-goal-id");
          const confirmed = confirm("Are you sure you want to delete this goal?");
          if (confirmed) {
            await deleteGoal(goalId);
          }
        });
      });
  }


  async function deleteGoal(goalId) {
    try {
      const res = await fetch(`http://localhost:8080/goals/${goalId}`, {
        method: "DELETE"
      });
  
      if (!res.ok) {
        const errMsg = await res.text();
        throw new Error(errMsg);
      }
  
      const updatedGoals = await fetchGoals();
      populateGoals(updatedGoals);
    } catch (err) {
      console.error("Failed to delete goal:", err);
      alert("Could not delete goal.");
    }
  }
  

async function main(){
    initialLoad();
    
    let budgets = await refreshBudgets();
    let goals = await fetchGoals();
    if(budgets){
        ApplyMonthBudgets(budgets);
        console.log("Displaying...");
    }
    if(goals){
        populateGoals(goals);
    }
    EventListeners();
}

main();