document.getElementById('submit').onclick = function () {
    event.preventDefault();
    const form = document.getElementById('form');
    const formData = new FormData(form);
  
    const usernameField = document.querySelector("input[name='username']");
    const passwordField = document.querySelector("input[name='password']");
    const confirmPasswordField = document.querySelector("input[name='confirm-password']");
  
    // Validate username
    if (usernameField.value.trim() === "") {
      return showMessage("Username is required.");
    } else if (usernameField.value.length > 25) {
      return showMessage("Username must be less than 25 characters.");
    } else if (!/^[\u0000-\u007F]+$/.test(usernameField.value)) {
      return showMessage("Username must not contain non-ASCII characters.");
    }
  
    // Validate password
    if (passwordField.value.trim() === "") {
      return showMessage("Password is required.");
    } else if (passwordField.value.length <= 6) {
      return showMessage("Password must be longer than 6 characters.");
    } else if (passwordField.value.length >= 20) {
      return showMessage("Password must be less than 20 characters.");
    } else if (passwordField.value.toLowerCase() === "password") {
      return showMessage("Password cannot be 'password'.");
    } else if (passwordField.value !== confirmPasswordField.value){
      return showMessage("Passwords do not match.");
    }
  
    fetch('/signup',{
      method: 'POST',
      credentials: 'include',
      body: formData,
    })
    .then(response => {
      if (response.ok){
        window.location.href = '/';
      } else{
        return response.json().then(errorData => {
          showMessage(errorData.message);
        });
      }
    })
    .catch(error =>{
      console.error("Error:", error);
      showMessage("An error occurred. Please check your connection.");
    });
  
  };