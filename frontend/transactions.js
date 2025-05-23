let user = JSON.parse(localStorage.getItem("user"));
async function getMap(){
    let catMap = new Map();
    let data;
    try {
        const res = await fetch(`http://localhost:8080/categories`);
        data = await res.json();
    } catch (error) {
        console.log(error);
        alert("Couldn't get map");
    }
    data.forEach(cat =>{
        if(cat.name != "total"){
            catMap.set(cat.name, cat.category_id);
        }
    });
    return catMap;

}

async function deleteTransaction(transaction) {
    try {
        const res = await fetch(`http://localhost:8080/transactions/${transaction.Transaction.Id}?google_id=${user.google_id}`, {method: "DELETE"});
        if (!res.ok) {
            console.log(await res.text());
            alert("Failed to delete transaction");
        }
    } catch (error) {
        alert("failed to dellete transaction");
        console.log(error);
    }
    main();
    
}

async function populateCategoryDropdown() {
    const select = document.getElementById("tx-category");
    select.innerHTML = ""; // Clear existing options
  
    const categoryMap = await getMap();
  
    categoryMap.forEach((id, name) => {
      const option = document.createElement("option");
      option.value = id;
      option.textContent = name;
      select.appendChild(option);
    });
  }
  

async function getTransactions() {
    document.getElementById("add-transaction-btn").addEventListener("click", () => {
        showAddModal()
      });
    document.getElementById("save-to-file").addEventListener("click", () => {
        saveToFile()
      });
      
    try {
        const res = await fetch(`http://localhost:8080/transactions/${user.google_id}`);
        const data = await res.json();
        return data;
    } catch (error) {
        console.log(error);
        alert("Couldn't fetch transactions");
    }
}

async function showEditModal(transaction) {
    // Set form title
    document.getElementById("form-title").textContent = "Edit Transaction";
  
    // Pre-fill form fields
    document.getElementById("tx-date").value = transaction.Transaction.Posting_date;
    document.getElementById("tx-description").value = transaction.Transaction.Description;
    document.getElementById("tx-amount").value = transaction.Transaction.Amount;
    const translationMap = await getMap();
    document.getElementById("tx-category").value = translationMap.get(transaction.T_type);
  
    // Show the form
    document.getElementById("transaction-form").style.display = "block";
  
    // Save button updates transaction
    const saveButton = document.getElementById("save-transaction");
    saveButton.onclick = async () => {
      const updatedTx = {
        google_id: user.google_id,
        date: document.getElementById("tx-date").value,
        description: document.getElementById("tx-description").value,
        amount: parseFloat(document.getElementById("tx-amount").value),
        category_id: parseInt(document.getElementById("tx-category").value),
      };
      console.log(updatedTx);
  
      try {
        const res = await fetch(`http://localhost:8080/transactions/${transaction.Transaction.Id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(updatedTx),
        });
  
        if (!res.ok) throw new Error(await res.text());
  
        alert("Transaction updated.");
        document.getElementById("transaction-form").style.display = "none";
        main(); 
      } catch (err) {
        console.error("Update failed:", err);
        alert("Could not update transaction.");
      }
    };
  
    // Cancel button
    document.getElementById("cancel-transaction").onclick = () => {
      document.getElementById("transaction-form").style.display = "none";
    };
  }
  

function displayTransactions(transactions) {
    const container = document.getElementById("transaction-table");
    container.innerHTML = "";

    if (user) {
    document.getElementById("user-info").textContent = `Logged in as ${user.username} (${user.email})`;
    }

  
    // Header
    const header = document.createElement("div");
    header.className = "transaction-row transaction-row-header";
    header.innerHTML = `
      <div>Date</div>
      <div>Description</div>
      <div>Category</div>
      <div></div> <!-- for edit button -->
    `;
    container.appendChild(header);
  
    transactions.forEach(tx => {
      const row = document.createElement("div");
      row.className = "transaction-row";
  
      const date = tx.Transaction.Posting_date || "";
      const description = tx.Transaction.Description || "";
      const category = tx.T_type;
  
      row.innerHTML = `
        <div>${date}</div>
        <div>${description}</div>
        <div>${category}</div>
        <div><button class="edit-tx-btn" data-id="${tx.Transaction.id}">Edit</button></div>
        <div></div> <div></div> <div></div>
        <div><button class="delete-tx-btn" data-id="${tx.Transaction.id}">Delete</button></div>
      `;

      const editBtn = row.querySelector(".edit-tx-btn");
      const deleteBtn = row.querySelector(".delete-tx-btn");
      editBtn.addEventListener("click", () => {
        showEditModal(tx);
      });
      deleteBtn.addEventListener("click", () => {
        deleteTransaction(tx);
      });
  
      container.appendChild(row);
    });
  }


async function showAddModal(){
    document.getElementById("form-title").textContent = "Add Transaction";
        document.getElementById("transaction-form").style.display = "block";
      
        // Clear form values
        document.getElementById("tx-date").value = "";
        document.getElementById("tx-description").value = "Description";
        document.getElementById("tx-amount").value = 0.00;
        document.getElementById("tx-category").value = 9;
        
        const saveButton = document.getElementById("save-transaction");
        saveButton.onclick = async () => {
          const updatedTx = {
            google_id: user.google_id,
            date: document.getElementById("tx-date").value,
            description: document.getElementById("tx-description").value,
            amount: parseFloat(document.getElementById("tx-amount").value),
            category_id: parseInt(document.getElementById("tx-category").value),
          };
          console.log(updatedTx);
      
          try {
            const res = await fetch(`http://localhost:8080/transactions`, {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify(updatedTx),
            });
      
            if (!res.ok) throw new Error(await res.text());
      
            alert("Transaction Added.");
            document.getElementById("transaction-form").style.display = "none";
            main(); 
          } catch (err) {
            console.error("Add failed:", err);
            alert("Could not add transaction.");
          }
        }
}
  
async function saveToFile(){
    const link = document.createElement("a");
    link.href = `http://localhost:8080/transactions/download?google_id=${user.google_id}`;
    link.download = "transactions.csv";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}


async function main(){
    await populateCategoryDropdown();
    const data = await getTransactions();
    displayTransactions(data);
}


main();
