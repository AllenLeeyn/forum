document.getElementById('submit-comment').onclick = function () {
    const commentInput = document.getElementById('comment-input')
    const commentText = commentInput.value.trim()

    if (!commentText) {
        showMessage('Comment cannot be empty!')
        return;
    }

    const url = new URL(window.location.href)
    const postId = url.searchParams.get('id')

    if (!postId) {
        showMessage('Post ID not found in URL.')
        return;
    }

    // Send the comment to the server
    fetch(`/add-comment?postId=${postId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ Comment: commentText }),
    })
        .then(response => {
            if (!response.ok) {
                console.error()
                return response.json().then(errorData => {
                  showMessage(errorData.message);
                });
            } else {
              location.reload()
            }
        })

    commentInput.value = '';
};
