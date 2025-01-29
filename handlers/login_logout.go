package handlers

import (
	"forum/dbTools"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := checkSessionValidity(w, r)
	if sessionCookie != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		//	Going to the login page
		ExecuteTemp(w, "login.html", nil)
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, _, passwd, e := getCredentials(r, false)

	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	// check that credentials are valid
	user, _ := db.SelectUserByField("name", username)
	if user == nil || bcrypt.CompareHashAndPassword(user.PwHash, []byte(passwd)) != nil {
		http.Error(w, "incorrect username and/or password X", 400)
		return
	}

	session, e := db.SelectActiveSessionBy("user_id", user.ID)
	if e == nil {
		db.UpdateSession(&dbTools.Session{
			IsActive:   false,
			ExpireTime: time.Now(),
			LastAccess: time.Now(),
			ID:         session.ID,
		})
	}
	createSession(w, user)
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie("session-id")
	if sessionCookie == nil {
		log.Println("no session cookie??")
	} else {
		expireSession(w, sessionCookie.Value)
	}
	http.Redirect(w, r, "./login", http.StatusSeeOther)
}
