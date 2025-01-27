document.addEventListener("DOMContentLoaded", function () {
    // Get all like buttons
    const likeButtons = document.querySelectorAll(".like-button");
    const dislikeButtons = document.querySelectorAll(".dislike-button");
  
    // Add click event listeners
    likeButtons.forEach((button) => {
      button.addEventListener("click", async () => {
        const postId = button.getAttribute("data-id"); // Get the post ID
        const forType = button.getAttribute("data-for"); // Get the post ID
        const isLiked = button.classList.contains("liked"); // Check if already liked

        const dislikeButton = document.querySelector(`.dislike-button[data-id="${postId}"]`);
        const isDisliked = dislikeButton.classList.contains("disliked");

        try {
          // Send a POST request to the server
          const response = await fetch(`/feedback?for=${forType}&id=${postId}&action=${isLiked ? 'unlike' : 'like'}`, {
            method: "POST",
          });
          const likeCountSpan = button.querySelector("span");
          const currentCount = parseInt(likeCountSpan.textContent, 10);
  
          if (response.ok) {
            if (button.classList.contains("liked")) {
              button.classList.remove("liked");
              likeCountSpan.textContent = currentCount - 1;
            } else{
              button.classList.add("liked");
              likeCountSpan.textContent = currentCount + 1;
              if(isDisliked){
                dislikeButton.classList.remove("disliked");
                const dislikeCountSpan = dislikeButton.querySelector("span");
                const currentCount = parseInt(dislikeCountSpan.textContent, 10);
                dislikeCountSpan.textContent = currentCount - 1;
              }
            }
          } else {
            // Alert the user if something went wrong
            showMessage("Please ensure you are logged in and try again later.");
          }
        } catch (error) {
          // Handle network or other errors
          showMessage("An error occurred. Please check your connection and try again.");
          console.error("Error:", error);
        }
      });
    });

    // Add click event listeners
    dislikeButtons.forEach((button) => {
      button.addEventListener("click", async () => {
        const postId = button.getAttribute("data-id"); // Get the post ID
        const forType = button.getAttribute("data-for"); // Get the post ID
        const isDisliked = button.classList.contains("disliked"); // Check if already disliked

        const likeButton = document.querySelector(`.like-button[data-id="${postId}"]`);
        const isLiked = likeButton.classList.contains("liked");

        try {
          // Send a POST request to the server
          const response = await fetch(`/feedback?for=${forType}&id=${postId}&action=${isDisliked ? 'undislike' : 'dislike'}`, {
            method: "POST",
          });
          const dislikeCountSpan = button.querySelector("span");
          const currentCount = parseInt(dislikeCountSpan.textContent, 10);
  
          if (response.ok) {
            if (button.classList.contains("disliked")) {
              button.classList.remove("disliked");
              dislikeCountSpan.textContent = currentCount - 1;
            } else{
              button.classList.add("disliked");
              dislikeCountSpan.textContent = currentCount + 1;
              if (isLiked){
                likeButton.classList.remove("liked");
                const likeCountSpan = likeButton.querySelector("span");
                const currentCount = parseInt(likeCountSpan.textContent, 10);
                likeCountSpan.textContent = currentCount - 1;
              }
            }
          } else {
            // Alert the user if something went wrong
            showMessage("Please ensure you are logged in and try again later.");
          }
        } catch (error) {
          // Handle network or other errors
          showMessage("An error occurred. Please check your connection and try again.");
          console.error("Error:", error);
        }
      });
    });
  });
  