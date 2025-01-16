package utils

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

//	Place for functions that executes different pages

var tmpl *template.Template

// Initializes all html files in templates folder
func InitializeHtml() {
	var err error
	tmpl, err = template.ParseGlob("template/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

// Home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Println("Now we're on the home page")
		CustomExecuteTemplate(w, "index.html", nil)
	} else {
		fmt.Println("Error 404")
	}
}

// Login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CustomExecuteTemplate(w, "login.html", nil)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			//	Error parsing data
		}
		loginData := LoginData{
			username: r.FormValue("name"),
			password: r.FormValue("password"),
		}
		nameIsValid, _ := CheckValidity(loginData.username, "username")
		passIsValid, _ := CheckValidity(loginData.password, "password")
		if nameIsValid && passIsValid {
			fmt.Println("Valid user, check if duplicate")
		} else {
			//	Placeholder:
			fmt.Println("Invalid name, or pass.")
			http.Error(w, "Invalid name, or pass.", http.StatusBadRequest)
		}
	} else {
		//	Error wrong method
	}
}

// Register page
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CustomExecuteTemplate(w, "register.html", nil)
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
		nameIsValid, _ := CheckValidity(registerData.username, "username")
		passIsValid, _ := CheckValidity(registerData.password, "password")
		emailIsValid, _ := CheckValidity(registerData.email, "email")
		if nameIsValid && passIsValid && emailIsValid {
			fmt.Println("Valid user, check if duplicate")
		} else {
			//	Placeholder:
			fmt.Println("Invalid name, pass, or email.")
			http.Error(w, "Invalid name, pass, or email.", http.StatusBadRequest)
		}
	} else {
		//	Error wrong method
	}
}

// Page for viewing posts
/* func ViewPostPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Now we're on the view posts page")
	err := tmpl.ExecuteTemplate(w, "view-post.html", nil)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
} */

// Page for making posts
/* func PostPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Now we're on the post page")
	err := tmpl.ExecuteTemplate(w, "post.html", nil)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
} */
