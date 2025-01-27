document.addEventListener("DOMContentLoaded", function () {
    const startThreadBtn = document.querySelector('.start-thread-btn');

    if (startThreadBtn) {
        startThreadBtn.addEventListener('click', function (event) {
            try{
                const response = fetch(`/start-thread`, {method: "GET",});

                if (response.ok){

                } else {
                    showMessage(" Please ensure you are logged in and try again later.");
                }
            }catch(error){
                showMessage("An error occurred. Please check your connection and try again.");
                console.error("Error:", error);}
        });
    }
  });
