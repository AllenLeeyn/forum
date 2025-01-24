package main

import (
	"errors"
	"forum/dbTools"
	"forum/handlers"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

var idCount int = 0

func signup(w http.ResponseWriter, r *http.Request) {
	name, email, passwd, e := getCredentials(r, true)

	// check that credentials are valid
	if e != nil {
		handlers.ExecuteError(w, e.Error(), 400)
		return
	}
	if user, _ := db.SelectUserByField("email", email); user != nil {
		handlers.ExecuteError(w, "Email has already been taken", 400)
		return
	}
	if user, _ := db.SelectUserByField("name", name); user != nil {
		handlers.ExecuteError(w, "Name has already been taken", 400)
		return
	}

	passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), 0)
	checkErr(err)

	// add the user to the database
	idCount++
	user := &dbTools.User{
		ID:        idCount,
		TypeID:    1,
		Name:      name,
		Email:     email,
		PwHash:    passwdHash,
		RegDate:   time.Now(),
		LastLogin: time.Now(),
	}
	db.InsertUser(user)

	// generate a uuid for the session and set it into a cookie
	id, _ := uuid.NewV4()
	cookie := &http.Cookie{
		Name:  "session-id",
		Value: id.String(),
	}
	http.SetCookie(w, cookie)

	db.InsertSession(&dbTools.Session{
		ID:        id.String(),
		UserID:    user.ID,
		IsActive:  true,
		StartTime: time.Now(),
		// ExpireTime: 0,
		LastAccess: time.Now(),
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	username, _, passwd, e := getCredentials(r, false)

	if e != nil {
		handlers.ExecuteError(w, e.Error(), 400)
		return
	}

	// check that credentials are valid
	user, _ := db.SelectUserByField("name", username)
	if user == nil || bcrypt.CompareHashAndPassword(user.PwHash, []byte(passwd)) != nil {
		handlers.ExecuteError(w, "incorrect username and/or password X", 400)
		return
	}

	// generate a uuid for the session and set it into a cookie
	id, _ := uuid.NewV4()
	cookie := &http.Cookie{
		Name:  "session-id",
		Value: id.String(),
	}
	http.SetCookie(w, cookie)

	db.InsertSession(&dbTools.Session{
		ID:        id.String(),
		UserID:    user.ID,
		IsActive:  true,
		StartTime: time.Now(),
		// ExpireTime: 0,
		LastAccess: time.Now(),
	})
}

func getCredentials(r *http.Request, isSignup bool) (string, string, string, error) {
	username := r.FormValue("name")
	email := r.FormValue("email")
	passwd := r.FormValue("password")

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	usernameRegex := `^[a-zA-Z0-9_-]{3,16}$`

	if r.Method != "POST" {
		return "", "", "", errors.New("invalid method")
	}
	if !validRegex(username, usernameRegex) || !validPsswrd(passwd) {
		return "", "", "", errors.New("invalid username and/or password Y")
	}
	if isSignup && !validRegex(email, emailRegex) {
		return "", "", "", errors.New("invalid email")
	}
	return username, email, passwd, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
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
