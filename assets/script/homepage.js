// homepage.js

function openPopup() {
    document.getElementById("popupOverlay").style.display = "block";
    document.getElementById("popup").style.display = "block";
}

function closePopup() {
    document.getElementById("popupOverlay").style.display = "none";
    document.getElementById("popup").style.display = "none";
}

// Redirect to sign-in and sign-up pages (assuming these pages are already created)
document.getElementById("signIn").addEventListener("click", function() {
    window.location.href = "signin.html";
});

document.getElementById("signUp").addEventListener("click", function() {
    window.location.href = "signup.html";
});
