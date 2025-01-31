package handlers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := checkSessionValidity(w, r)
	if sessionCookie != nil { // logged in user trying to login again?
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		ExecuteTmpl(w, "login.html", nil) // Going to the login page
		return
	} else if r.Method != http.MethodPost {
		ExecuteError(w, "Tmpl", "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, _, passwd, e := getCredentials(r, false)
	if e != nil {
		ExecuteError(w, "json", e.Error(), http.StatusBadRequest)
		return
	}

	// check that credentials are valid
	user, _ := db.SelectUserByField("name", username)
	if user == nil || bcrypt.CompareHashAndPassword(user.PwHash, []byte(passwd)) != nil {
		ExecuteError(w, "json", "incorrect username and/or password", http.StatusBadRequest)
		return
	}

	s, e := db.SelectActiveSessionBy("user_id", user.ID)
	if e == nil {
		expireSession(w, s.ID)
	}
	createSession(w, user)
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie("session-id")
	if sessionCookie == nil {
		ExecuteError(w, "json", "user not logged in", http.StatusBadRequest)
	} else {
		expireSession(w, sessionCookie.Value)
	}
	http.Redirect(w, r, "./", http.StatusSeeOther)
}
