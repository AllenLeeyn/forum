function showMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('toast-message');
    messageDiv.textContent = message;
    document.body.appendChild(messageDiv);
    setTimeout(() => messageDiv.remove(), 3000); // Remove message after 3 seconds
  }

document.getElementById('new-post').onclick = function () {
    fetch('/new-post', {
        method: 'GET',
        credentials: 'include', // Ensures cookies are sent with the request
    })
    .then(response => {
      console.log("Response status:", response.status);
      console.log("Response OK:", response.ok);
      
        if (response.ok) {
            window.location.href = '/new-post';
        } else {
          return response.json().then(errorData => {
            showMessage(errorData.message);
          });
        }
    })
    .catch(error => {
      showMessage("An error occurred. Please check your connection and try again.");
      console.error("Error:", error);
    });
  };
  
document.getElementById('logout-btn').onclick = function () {
      fetch('/logout', {
          method: 'POST',
          credentials: 'include', // Ensures cookies are sent with the request
      })
      .then(response => {
          if (response.ok) {
            showMessage("Logout successful!");
            setTimeout(() => {
              window.location.href = '/';
            }, 2000);
          } else {
            return response.json().then(errorData => {
              showMessage(errorData.message);
            });
          }
      })
      .catch(error => {
          console.error('Error logging out:', error);
          showMessage("An error occurred during logout.");
      });
};
