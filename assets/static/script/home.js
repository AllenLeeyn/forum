 document.getElementById('start-thread-btn').onclick = function () {
    fetch('/start-thread', {
        method: 'GET',
        credentials: 'include', // Ensures cookies are sent with the request
    })
    .then(response => {
        if (response.ok) {
            window.location.href = '/start-thread';
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