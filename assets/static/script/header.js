function showMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('toast-message');
    messageDiv.textContent = message;
    document.body.appendChild(messageDiv);
    setTimeout(() => messageDiv.remove(), 3000); // Remove message after 3 seconds
  }

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
            showMessage("Failed to log out.");
          }
      })
      .catch(error => {
          console.error('Error logging out:', error);
          showMessage("An error occurred during logout.");
      });
};
