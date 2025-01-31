document.getElementById('submit').onclick = function (event) {
  event.preventDefault();
  const form = document.getElementById('form');
  const formData = new FormData(form);

  const usernameField = document.querySelector("input[name='username']");
  const passwordField = document.querySelector("input[name='password']");

  // Validate username
    if (usernameField.value.trim() === "") {
      return showMessage("Username is required.");
    } else if (usernameField.value.length < 3) {
      return showMessage("Invalid username");
    } else if (usernameField.value.length > 16) {
      return showMessage("Invalid username");
    } else if (!/^[\u0000-\u007F]+$/.test(usernameField.value)) {
      return showMessage("Invalid username");
    }

  // Validate password
    if (passwordField.value.trim() === "") {
      return showMessage("Password is required.");
    } else if (passwordField.value.length <= 7) {
      return showMessage("Password must be at least 8 characters long.");
    }

  fetch('/login',{
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