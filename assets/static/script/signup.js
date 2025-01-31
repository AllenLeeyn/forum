document.getElementById('submit').onclick = function (event) {
    event.preventDefault();
    const form = document.getElementById('form');
    const formData = new FormData(form);
  
    const usernameField = document.querySelector("input[name='username']");
    const emailField = document.querySelector("input[name='email']");
    const passwordField = document.querySelector("input[name='password']");
    const confirmPasswordField = document.querySelector("input[name='confirm-password']");
  
    // Validate username
    if (usernameField.value.trim() === "") {
      return showMessage("Username is required.");
    } else if (usernameField.value.length < 3) {
      return showMessage("Username must be between 3 to 16 alphanumeric characters, '_' or '-'");
    } else if (usernameField.value.length > 16) {
      return showMessage("Username must be between 3 to 16 alphanumeric characters, '_' or '-'");
    } else if (!/^[\u0000-\u007F]+$/.test(usernameField.value)) {
      return showMessage("Username must be between 3 to 16 alphanumeric characters, '_' or '-'");
    }
  
    // Validate email
    if (emailField.value.trim() === "") {
      return showMessage("Email is required.");
    } else if (!/^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(emailField.value)) {
      return showMessage("Please enter a valid email address.");
    }

    // Validate password
    if (passwordField.value.trim() === "") {
      return showMessage("Password is required.");
    } else if (passwordField.value.length <= 7) {
      return showMessage("Password must be at least 8 characters long.");
    } else if (passwordField.value.toLowerCase() === "password") {
      return showMessage("Password cannot be 'password'.");
    }
    if (confirmPasswordField.value.trim() === ""){
      return showMessage("Passwords do not match.");
    } else if(passwordField.value !== confirmPasswordField.value){
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