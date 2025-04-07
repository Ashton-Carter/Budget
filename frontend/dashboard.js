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
  



async function main(){
    let transactions;
    document.getElementById("user-info").innerText = `Logged in as ${user.username} (${user.email})`;
    console.log("User object in localStorage:", user);

    transactions = await fetchTransactions();
    console.log("User transactions:", transactions);
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

main();