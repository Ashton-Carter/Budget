let user = JSON.parse(localStorage.getItem("user"));
const categories = ["Food", "Gas", "Entertainment", "Shopping", "Subscriptions", "Transfers", "Alcohol", "Other"];

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
    const type = [...typeRadios].find(r => r.checked).value;
  
    if (type === "total") {
      const total = parseFloat(document.getElementById("total-budget").value);
      if (isNaN(total) || total <= 0) {
        alert("Please enter a valid total.");
        return;
      }
  
      // Call submitBudget with category_id = 0 for "By Total"
      submitBudget("Weekly Total Budget", 0, total);
      hideBudgetModal();
  
    } else if (type === "category") {
      const categoryValues = {};
      let valid = true;
  
      categories.forEach(cat => {
        const val = parseFloat(document.querySelector(`input[name='category-${cat}']`).value || 0);
        if (val < 0) valid = false;
        categoryValues[cat] = val;
      });
  
      if (!valid) {
        alert("Please enter valid amounts for all categories.");
        return;
      }
  
      // Map category names to IDs
      const categoryMap = {
        Food: 1,
        Gas: 2,
        Entertainment: 3,
        Shopping: 4,
        Subscriptions: 5,
        Transfers: 6,
        Alcohol: 7,
        Other: 8,
      };
  
      // Submit one budget per category
      Object.entries(categoryValues).forEach(([cat, amount]) => {
        if (amount > 0) {
          const categoryId = categoryMap[cat];
          submitBudget(`Weekly - ${cat}`, categoryId, amount);
        }
      });
  
      hideBudgetModal();
    }
  };
  

async function submitBudget(name, categoryId, amount) {
    const payload = {
      user_id: user.google_id,         // assuming your backend maps this
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

function displayBudgets(budgets) {
    budgets.forEach(bud => {
        console.log(bud);
    });
}

function EventListeners(){
    document.getElementById("create-budget-btn").addEventListener("click", () => {
        setupBudgetModal();
        showBudgetModal();
      });
}

async function main(){
    initialLoad();
    EventListeners();
    let budgets = await refreshBudgets();
    if(budgets){displayBudgets(budgets)};
}

main();