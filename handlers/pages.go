package handlers

import (
	"fmt"
	"forum/dbTools"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

//	Place for functions that executes different pages

var tmpl *template.Template
var db *dbTools.DBContainer

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
		http.Error(w, "Error 404, Page not found", http.StatusNotFound)
		return
	}
	// check session from cookie

	// get Queries. No sorting implemented yet
	filterBy := r.URL.Query().Get("filterBy")
	// orderBy := r.URL.Query().Get("orderBy")

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id > len(db.Categories) {
		id = -1
	}

	posts, err := db.SelectPosts(filterBy, "", id)
	if err != nil {
		fmt.Println(err)
		// something went wrong
	}
	CustomExecuteTemplate(w, "homepage.html", homepageData{posts, db.Categories})
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
func SignupPage(w http.ResponseWriter, r *http.Request) {
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

// Page for making posts
func PostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("something went wrong with post_id")
		}

		post, err := db.SelectPost(id)
		if err != nil || post == nil {
			fmt.Println("something went wrong with getting post or nothing found")
		}

		comments, err := db.SelectComments(id, "")
		if err != nil {
			fmt.Println("something went wrong with grabbing comments")
		}
		CustomExecuteTemplate(w, "post.html", postpageData{*post, comments})

	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			//	Error parsing data
		}
		postData := post{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			//UserID:    However we get the users ID,
			//UserName: However we get the users Name,+
		}
		titleIsValid, titleInvalidReason := CheckValidity(postData.Title, "postTitle")
		contentIsValid, contentInvalidReason := CheckValidity(postData.Content, "postContent")
		if titleIsValid && contentIsValid {
			fmt.Println("Valid user, check if duplicate")
		} else {
			//	Placeholder:
			reason := titleInvalidReason + "\n" + contentInvalidReason
			http.Error(w, reason, http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}
}
