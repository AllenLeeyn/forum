document.getElementsByClassName("close")[0].onclick = function() {
document.getElementById("newPostModal").style.display = "none";
}

document.getElementById('new-post').onclick = function () {
  fetch('/new-post', {
      method: 'GET',
      credentials: 'include', // Ensures cookies are sent with the request
  })
  .then(response => {
      if (response.ok) {
        document.getElementById("newPostModal").style.display = "block";
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

document.getElementById('create-new-post').onclick = function (event) {
event.preventDefault();
const form = document.getElementById('post-thread-form');
const formData = new FormData(form);

fetch('/new-post', {
    method: 'POST',
    credentials: 'include', // Ensures cookies are sent with the request
    body: formData,
})
.then(response => {
    if (response.redirected) {
      window.location.href = response.url;
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
