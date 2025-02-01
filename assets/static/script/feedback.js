document.addEventListener("DOMContentLoaded", function () {
  // Get all like and dislike buttons
  const likeButtons = document.querySelectorAll(".like-button");
  const dislikeButtons = document.querySelectorAll(".dislike-button");

  // Unified function to handle both like and dislike actions
  function handlePostFeedback(forType, parentID, action) {
      const postElement = document.querySelector(`.post-actions[data-id="${parentID}"][data-for="${forType}"]`);
      const currentState = parseInt(postElement.getAttribute('data-state'), 10);
      let newState = currentState;
      let newAction = 0; // Default action (neutral)

      // Get like and dislike button elements for the post
      const likeButton = document.querySelector(`.like-button[data-id="${parentID}"][data-for="${forType}"]`);
      const likeCountSpan = likeButton.querySelector("span");
      let newLikeCount = parseInt(likeCountSpan.textContent);

      const dislikeButton = document.querySelector(`.dislike-button[data-id="${parentID}"][data-for="${forType}"]`);
      const dislikeCountSpan = dislikeButton.querySelector("span");
      let newDislikeCount = parseInt(dislikeCountSpan.textContent);

      // Handle the "like" action
      if (action === "like") {
          if (currentState === 1) { // unLike
              newState = 0
              newLikeCount = newLikeCount - 1;
              newAction = 0;
          } else {
              if (currentState === -1) { // remove dislike
                newDislikeCount = newDislikeCount - 1;
              }
              newState = 1
              newLikeCount = newLikeCount + 1;
              newAction = 1;
          }
      } 
      // Handle the "dislike" action
      else if (action === "dislike") {
          if (currentState === -1) { // unDislike
            newState = 0
            newDislikeCount = newDislikeCount - 1;
            newAction = 0;
          } else {
              if (currentState === 1) { // remove like
                newLikeCount = newLikeCount - 1;
              }
              newState = -1
              newDislikeCount = newDislikeCount + 1;
              newAction = -1; // Dislike action
          }
      }
      
      // Create the data to send in the POST request
      const feedback = {
        tgt: forType,
        parentID: parentID,
        rating: newAction
      };

      // You could handle fetch success/failure like this:
      fetch('/feedback', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(feedback)
      })
      .then(response => {
        if (!response.ok) {
          return response.json().then(errorData => {
            showMessage(errorData.message);
          });
        } else {
          postElement.setAttribute('data-state', newState);
          likeCountSpan.textContent = newLikeCount;
          dislikeCountSpan.textContent = newDislikeCount;
        }
      })
      .catch(error => {
        console.error("Error:", error);
        showMessage("An error occurred. Please check your connection.");
      });
  }

  // Attach event listeners to both like and dislike buttons
  likeButtons.forEach((button) => {
      button.addEventListener("click", () => {
          const forType = button.getAttribute("data-for");
          const parentID = parseInt(button.getAttribute("data-id"), 10);
          handlePostFeedback(forType, parentID, "like");
      });
  });

  dislikeButtons.forEach((button) => {
      button.addEventListener("click", () => {
          const forType = button.getAttribute("data-for");
          const parentID = parseInt(button.getAttribute("data-id"), 10);
          handlePostFeedback(forType, parentID, "dislike");
      });
  });
});
