package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var tmpl *template.Template

func InitializeHtml() {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Println("Now we're on the home page")
		err := tmpl.ExecuteTemplate(os.Stdout, "index.html", nil)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}
	} else {
		fmt.Println("Should throw error here")
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Now we're on the login page")
	err := tmpl.ExecuteTemplate(os.Stdout, "login.html", nil)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Now we're on the register page")
	err := tmpl.ExecuteTemplate(os.Stdout, "register.html", nil)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func ViewPostPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Now we're on the view posts page")
	err := tmpl.ExecuteTemplate(os.Stdout, "view-post.html", nil)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func PostPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Now we're on the post page")
	err := tmpl.ExecuteTemplate(os.Stdout, "post.html", nil)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
