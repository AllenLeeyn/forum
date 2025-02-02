function showMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('toast-message');
    messageDiv.textContent = message;
    document.body.appendChild(messageDiv);
    setTimeout(() => messageDiv.remove(), 3000); // Remove message after 3 seconds
  }
  