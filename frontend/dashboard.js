const user = JSON.parse(localStorage.getItem("user"));
if (!user) {
  location.href = "/index.html";
}

document.getElementById("user-info").innerText = `Logged in as ${user.username} (${user.email})`;

// Handle CSV upload
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
