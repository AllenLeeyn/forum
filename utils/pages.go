package utils

import (
	"fmt"
	"forum/dbTools"
	"log"
	"net/http"
	"text/template"
)

//	Place for functions that executes different pages

var tmpl *template.Template
var db *dbTools.DBContainer

type user = dbTools.User
type session = dbTools.Session
type post = dbTools.Post
type feedback = dbTools.Feedback
type comment = dbTools.Comment

type postList struct {
	Posts     []post
	sessionID int
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

// Home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ExecuteError(w, http.StatusNotFound, "Page Not Found")
		return
	}
	UpdateAndExecuteHome(w)
}

// Login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		//	Going to the login page
		CustomExecuteTemplate(w, "login.html", nil)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			//	Error parsing data
		}
		loginData := LoginData{
			username: r.FormValue("name"),
			password: r.FormValue("password"),
		}
		nameIsValid, nameInvalidReason := CheckValidity(loginData.username, "username")
		passIsValid, passInvalidReason := CheckValidity(loginData.password, "password")
		if nameIsValid && passIsValid {
			fmt.Println("Valid user, check if duplicate")
			http.Error(w, "placeholder successful login (check for if user exists in db first)", http.StatusMethodNotAllowed)
		} else {
			//	Placeholder:
			reason := nameInvalidReason + "\n" + passInvalidReason
			http.Error(w, reason, http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Register page
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CustomExecuteTemplate(w, "signup.html", nil)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			//	Error parsing data
		}
		if r.FormValue("password") != r.FormValue("passwordConf") {
			//	Placeholder:
			fmt.Println("The two passwords given are different")
			http.Error(w, "The two password are different.", http.StatusBadRequest)
		}
		registerData := RegisterData{
			username: r.FormValue("name"),
			password: r.FormValue("password"),
			email:    r.FormValue("email"),
		}
		nameIsValid, nameInvalidReason := CheckValidity(registerData.username, "username")
		passIsValid, passInvalidReason := CheckValidity(registerData.password, "password")
		emailIsValid, emailInvalidReason := CheckValidity(registerData.email, "email")
		if nameIsValid && passIsValid && emailIsValid {
			fmt.Println("Valid user, check if duplicate")
			http.Error(w, "placeholder successful registration (check for if user exists in db first)", http.StatusMethodNotAllowed)
		} else {
			//	Placeholder:
			reason := nameInvalidReason + "\n" + passInvalidReason + "\n" + emailInvalidReason
			http.Error(w, reason, http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Page for viewing posts
func ViewPostPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Now we're on the view posts page")
	CustomExecuteTemplate(w, "view-post.html", nil)
}

// Terms page
func TermsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CustomExecuteTemplate(w, "terms.html", nil)
	} else {
		ExecuteError(w, http.StatusMethodNotAllowed, "Invalid User Method")
	}
}

func StartThread(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CustomExecuteTemplate(w, "start-thread.html", nil)
		/* 	} else if r.Method == http.MethodPost {
		 */
	} else {
		ExecuteError(w, http.StatusMethodNotAllowed, "Invalid User Method")
	}
}

// Page for making posts
func PostThread(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		query := r.URL.Query()
		title := query.Get("threadTitle")
		content := query.Get("threadContent")
		reason := ""
		titleIsValid, titleInvalidReason := CheckValidity(title, "postTitle")
		contentIsValid, contentInvalidReason := CheckValidity(content, "postContent")
		if titleIsValid && contentIsValid {
			fmt.Println("Valid user, check if duplicate")
		} else {
			//	Placeholder:
			reason = titleInvalidReason + "\n" + contentInvalidReason
		}
		if !titleIsValid && !contentIsValid {
			errData := ErrorData{
				ErrorMessage: reason,
			}
			CustomExecuteTemplate(w, "start-thread.html", errData)
		} else {
			categories := []int{1, 2, 3}
			post := dbTools.Post{
				ID:           0,
				UserID:       0,
				UserName:     "Me",
				CommentCount: 0,
				LikeCount:    0,
				DislikeCount: 0,
				Title:        title,
				Content:      content,
				//	Placeholder:
				Categories: categories,
			}
			fmt.Println(post)
			err := db.InsertPost(post)
			if err != nil {
				fmt.Println(err)
			}
			UpdateAndExecuteHome(w)
		}
	} else {
		ExecuteError(w, http.StatusMethodNotAllowed, "Invalid User Method")
	}
}
