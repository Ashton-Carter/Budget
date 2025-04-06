window.handleCredentialResponse = function(response) {
    console.log("Encoded JWT ID token: " + response.credential);
  
    const user = parseJwt(response.credential);
    console.log("Sending to backend:", {
      google_id: user.sub,
      email: user.email,
      username: user.name,
      picture_url: user.picture
    });
  
    fetch("http://localhost:8080/auth/google", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        google_id: user.sub,
        email: user.email,
        username: user.name,
        picture_url: user.picture
      })
    })
      .then(res => res.json())
      .then(data => {
        localStorage.setItem("user", JSON.stringify(data)); // store user
        window.location.href = "dashboard.html"; // go to next page
      })
      .catch(err => console.error("Backend error:", err));
  };
  
  function parseJwt(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64).split('').map(c =>
        '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
      ).join('')
    );
    return JSON.parse(jsonPayload);
  }
  