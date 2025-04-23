let user = JSON.parse(localStorage.getItem("user"));
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

async function main(){
    initialLoad();
    refreshBudgets();
}

main();