package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/dbTools"
	"log"
	"net/http"
	"regexp"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
)

//	Place for functions that executes different pages

var tmpl *template.Template
var db *dbTools.DBContainer

type user = dbTools.User
type session = dbTools.Session
type post = dbTools.Post
type feedback = dbTools.Feedback
type comment = dbTools.Comment

type homepageData struct {
	SessionCookie *http.Cookie
	Categories    []string
	Posts         []post
	UserName      string
	FilterBy      string
	OrderBy       string
	Id            int
}

type postpageData struct {
	SessionCookie *http.Cookie
	Post          post
	Comments      []comment
}

type profilepageData struct {
	Name          string
	Email         string
	Posts         []post
	SessionCookie *http.Cookie
}

type startThreadData struct {
	SessionCookie *http.Cookie
	Categories    []string
}

type ErrorData struct {
	Type    string
	Message string
	Code    int
}

type PostError struct {
	ErrorMessage string
}

// Initializes all html files in templates folder
func Init(dbMain *dbTools.DBContainer) {
	var err error
	tmpl, err = template.ParseGlob("template/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
	db = dbMain
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

/*----------- Execute func -----------*/

// Custom execute template I wrote to both execute and check for error
func ExecuteTmpl(w http.ResponseWriter, name string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

// Execute error page if possible, otherwise use inbuilt http error
func ExecuteError(w http.ResponseWriter, errtype, msg string, code int) {
	errorData := ErrorData{errtype, msg, code}
	if errorData.Type == "Tmpl" {
		err := tmpl.ExecuteTemplate(w, "error.html", errorData)
		if err != nil {
			message := fmt.Sprintf("Error %d\n%s", errorData.Code, errorData.Message)
			http.Error(w, message, http.StatusNotFound)
		}
		return
	}
	w.WriteHeader(errorData.Code)
	type errorJson struct {
		Message string `json:"message"`
	}
	errJson := errorJson{errorData.Message}
	json.NewEncoder(w).Encode(errJson)
}

/*----------- aunthenticate func -----------*/

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
		return "", "", "", errors.New("invalid username and/or password Y")
	}
	if !validPsswrd(passwd) {
		return "", "", "", errors.New("invalid username and/or password Y")
	}
	if isSignup && !validRegex(email, emailRegex) {
		return "", "", "", errors.New("invalid email")
	}
	return username, email, passwd, nil
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

/*----------- session func -----------*/

func checkSessionValidity(w http.ResponseWriter, r *http.Request) (*http.Cookie, int) {
	sessionCookie, err := r.Cookie("session-id")
	if err != nil || sessionCookie == nil {
		return nil, -1
	}
	sessionID := sessionCookie.Value
	s, err := db.SelectActiveSessionBy("id", sessionID)
	if err != nil {
		fmt.Println(err)
		return nil, -1
	}
	if s.ExpireTime.Before(time.Now()) {
		expireSession(w, sessionCookie.Value)
		return nil, -1
	}
	return sessionCookie, s.UserID
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

func extendSession(w http.ResponseWriter, sessionCookie *http.Cookie) {
	if sessionCookie == nil {
		return
	}
	// generate a uuid for the session and set it into a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    sessionCookie.Value,
		MaxAge:   7200,
		HttpOnly: true,
	})
	db.UpdateSession(&session{
		IsActive:   true,
		ExpireTime: time.Now().Add(2 * time.Hour),
		LastAccess: time.Now(),
		ID:         sessionCookie.Value,
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
