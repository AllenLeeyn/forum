package handlers

import (
	"forum/dbTools"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie("session-id")
	if userID := checkSessionValidity(w, sessionCookie); userID != -1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		//	Going to the login page
		ExecuteTemp(w, "signup.html", nil)
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name, email, passwd, e := getCredentials(r, true)

	// check that credentials are valid
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	if user, _ := db.SelectUserByField("email", email); user != nil {
		http.Error(w, "email is already used", 400)
		return
	}
	if user, _ := db.SelectUserByField("name", name); user != nil {
		http.Error(w, "name is already used", 400)
		return
	}

	passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), 0)
	checkErr(err)

	// add the user to the database
	user := &dbTools.User{
		TypeID: 1,
		Name:   name,
		Email:  email,
		PwHash: passwdHash,
	}
	user.ID, err = db.InsertUser(user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	createSession(w, user)
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

// Terms page
func Terms(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ExecuteTemp(w, "terms.html", nil)
	} else {
		ExecuteError(w, http.StatusMethodNotAllowed, "Invalid User Method")
	}
}
