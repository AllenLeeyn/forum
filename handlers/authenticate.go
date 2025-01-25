package handlers

import (
	"errors"
	"fmt"
	"forum/dbTools"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

var idCount int = 1

func getCredentials(r *http.Request, isSignup bool) (string, string, string, error) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	passwd := r.FormValue("password")

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	usernameRegex := `^[a-zA-Z0-9_-]{3,16}$`

	if r.Method != "POST" {
		return "", "", "", errors.New("invalid method")
	}
	if !validRegex(username, usernameRegex) {
		fmt.Println("username")
		return "", "", "", errors.New("invalid username and/or password Y")
	}
	if !validPsswrd(passwd) {
		fmt.Println("passwd")
		return "", "", "", errors.New("invalid username and/or password Y")
	}
	if isSignup && !validRegex(email, emailRegex) {
		return "", "", "", errors.New("invalid email")
	}
	return username, email, passwd, nil
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func validRegex(input, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}

func validPsswrd(password string) bool {
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString
	hasDigit := regexp.MustCompile(`\d`).MatchString
	hasSpecial := regexp.MustCompile(`[@$!%*?&]`).MatchString
	isValidLength := len(password) >= 8

	return hasLowercase(password) &&
		hasUppercase(password) &&
		hasDigit(password) &&
		hasSpecial(password) &&
		isValidLength
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	curCookie, _ := r.Cookie("session-id")
	if userID := checkSessionValidity(w, curCookie); userID != -1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		//	Going to the login page
		CustomExecuteTemplate(w, "signup.html", nil)
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
		ID:        idCount,
		TypeID: 1,
		Name:   name,
		Email:  email,
		PwHash: passwdHash,
	}
	db.InsertUser(user)

	createSession(w, user)
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	curCookie, _ := r.Cookie("session-id")
	if userID := checkSessionValidity(w, curCookie); userID != -1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		//	Going to the login page
		CustomExecuteTemplate(w, "login.html", nil)
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

	createSession(w, user)
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session-id")
	if c == nil {
		log.Println("no session cookie??")
		return
	}
	expireSession(w, c.Value)
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

func checkSessionValidity(w http.ResponseWriter, c *http.Cookie) int {
	if c == nil {
		return -1
	}
	sessionID := c.Value
	s, err := db.SelectActiveSessionBy("id", sessionID)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	if s.ExpireTime.Before(time.Now()) {
		expireSession(w, c.Value)
		return -1
	}
	return s.UserID
}

func createSession(w http.ResponseWriter, user *dbTools.User) {
	// generate a uuid for the session and set it into a cookie
	id, _ := uuid.NewV4()
	cookie := &http.Cookie{
		Name:     "session-id",
		Value:    id.String(),
		MaxAge:   7200,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	db.InsertSession(&dbTools.Session{
		ID:         id.String(),
		UserID:     user.ID,
		IsActive:   true,
		ExpireTime: time.Now().Add(2 * time.Hour),
	})
}

func extendSession(w http.ResponseWriter, cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	// generate a uuid for the session and set it into a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    cookie.Value,
		MaxAge:   7200,
		HttpOnly: true,
	})
	db.UpdateSession(&session{
		IsActive:   true,
		ExpireTime: time.Now().Add(2 * time.Hour),
		LastAccess: time.Now(),
		ID:         cookie.Value,
	})
}

func expireSession(w http.ResponseWriter, sessionId string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    "", // Empty the cookie's value
		MaxAge:   -1, // Invalidate the cookie immediately
		HttpOnly: true,
	})
	db.UpdateSession(&session{
		IsActive:   false,
		ExpireTime: time.Now(),
		LastAccess: time.Now(),
		ID:         sessionId,
	})
}
