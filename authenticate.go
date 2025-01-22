package main

import (
	"forum/dbTools"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

var idCount int

func signup(w http.ResponseWriter, r *http.Request) {
	name, email, passwdHash, _ := getCredentials(r)

	// check that credentials are valid
	if _, e := db.SelectUserByField("email:" + email); e == nil {
		http.Error(w, "email is already used", 400)
		return
	}
	if _, e := db.SelectUserByField("name:" + name); e == nil {
		http.Error(w, "name is already used", 400)
		return
	}

	// add the user to the database
	user := &dbTools.User{
		ID:        idCount + 1,
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
	_, email, passwdHash, _ := getCredentials(r)

	// check that credentials are valid
	user, e := db.SelectUserByField("email:" + email)
	if e != nil || string(user.PwHash) != string(passwdHash) {
		http.Error(w, "incorrect email and/or password", 400)
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

func getCredentials(r *http.Request) (string, string, []byte, error) {
	username := r.FormValue("name")
	email := r.FormValue("email")
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 0)
	checkErr(err)

	//

	return username, email, passwdHash, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
