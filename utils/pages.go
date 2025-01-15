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
		CustomExecuteTemplate(w, "index.html", nil)
	} else {
		fmt.Println("Error 404")
	}
}

// Login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	CustomExecuteTemplate(w, "login.html", nil)
}

// Register page
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	CustomExecuteTemplate(w, "register.html", nil)
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
