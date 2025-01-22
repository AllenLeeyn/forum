"use strict";

const usernameField = document.querySelector("input[name='username']");
const passwordField = document.querySelector("input[name='password']");
const form = document.getElementById("form");

// Create a div dynamically to display error messages
let errorElement = document.createElement("div");
errorElement.setAttribute("id", "error");
errorElement.style.color = "white";
errorElement.style.marginTop = "10px";
errorElement.style.marginBottom = "10px";
form.appendChild(errorElement);

form.addEventListener("submit", (e) => {
  let messages = [];

  // Validate username
  if (usernameField.value.trim() === "") {
    messages.push("Username is required.");
  } else if (usernameField.value.length > 25) {
    messages.push("Username must be less than 25 characters.");
  } else if (!/^[\u0000-\u007F]+$/.test(usernameField.value)) {
    messages.push("Username must not contain non-ASCII characters.");
  }

  // Validate password
  if (passwordField.value.trim() === "") {
    messages.push("Password is required.");
  } else if (passwordField.value.length <= 6) {
    messages.push("Password must be longer than 6 characters.");
  } else if (passwordField.value.length >= 20) {
    messages.push("Password must be less than 20 characters.");
  } else if (passwordField.value.toLowerCase() === "password") {
    messages.push("Password cannot be 'password'.");
  }

  // Show error messages if validation fails
  if (messages.length > 0) {
    e.preventDefault(); // Prevent form submission
    errorElement.innerText = messages.join("\n"); // Display error messages
  }
});
