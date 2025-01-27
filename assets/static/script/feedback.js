document.addEventListener("DOMContentLoaded", function () {
    // Get all like buttons
    const likeButtons = document.querySelectorAll(".like-button");
  
    // Add click event listeners
    likeButtons.forEach((button) => {
      button.addEventListener("click", async () => {
        const postId = button.getAttribute("data-id"); // Get the post ID
        try {
          // Send a POST request to the server
          const response = await fetch(`/feedback?for=post&id=${postId}`, {
            method: "POST",
          });
  
          if (response.ok) {
            // Add the 'liked' class to the button to make the heart red
            button.classList.add("liked");
  
            // Optionally, update the like count
            const likeCountSpan = button.querySelector("span");
            const currentCount = parseInt(likeCountSpan.textContent, 10);
            likeCountSpan.textContent = currentCount + 1;
          } else {
            // Alert the user if something went wrong
            showMessage("Something went wrong while liking the post. Please ensure you are logged in and try again later.");
          }
        } catch (error) {
          // Handle network or other errors
          showMessage("An error occurred. Please check your connection and try again.");
          console.error("Error:", error);
        }
      });
    });
  });
  