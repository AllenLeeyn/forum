const usernameField = document.querySelector("input[name='username']");
const passwordField = document.querySelector("input[name='password']");
const form = document.getElementById("form");

let errorElement = document.createElement("div");
errorElement.setAttribute("id", "error");
errorElement.style.color = "#ff847c";
errorElement.style.marginTop = "10px";
errorElement.style.marginBottom = "10px";
form.appendChild(errorElement);

form.addEventListener("submit", (e) => {
  let messages = [];

  // Validate username
  if (usernameField.value.trim() === "") {
    showMessage("Username is required.");
  } else if (usernameField.value.length > 25) {
    showMessage("Username must be less than 25 characters.");
  } else if (!/^[\u0000-\u007F]+$/.test(usernameField.value)) {
    showMessage("Username must not contain non-ASCII characters.");
  }

  // Validate password
  if (passwordField.value.trim() === "") {
    showMessage("Password is required.");
  } else if (passwordField.value.length <= 6) {
    showMessage("Password must be longer than 6 characters.");
  } else if (passwordField.value.length >= 20) {
    showMessage("Password must be less than 20 characters.");
  } else if (passwordField.value.toLowerCase() === "password") {
    showMessage("Password cannot be 'password'.");
  }

  if (messages.length > 0) {
    e.preventDefault();
    errorElement.innerText = messages.join("\n");
  }
});
