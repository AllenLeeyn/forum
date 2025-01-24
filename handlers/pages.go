package handlers

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
		ExecuteError(w, "Error 404, Page not found", http.StatusNotFound)
		return
	}
	UpdateAndExecuteHome(w, r)
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
		ExecuteError(w, "Invalid User Method", http.StatusMethodNotAllowed)
	}
}

// Page for user to draft their post (if method == get), otherwise post it (if method == post)
func StartThread(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CustomExecuteTemplate(w, "start-thread.html", nil)
	} else if r.Method == http.MethodPost {
		title, content, categoriesInt := GetData(w, r)
		if title == "" && content == "" && categoriesInt == nil {
			return
		}
		fmt.Println(categoriesInt)
		titleIsValid, titleInvalidReason := CheckValidity(title, "postTitle")
		contentIsValid, contentInvalidReason := CheckValidity(content, "postContent")
		if !(titleIsValid && contentIsValid) {
			errData := ErrorData{
				ErrorMessage: titleInvalidReason + "\n" + contentInvalidReason,
			}
			CustomExecuteTemplate(w, "start-thread.html", errData)
			return
		}
		post := CreatePost(w, r, title, content, categoriesInt)
		//	DislikeCount in the negatives indicates an error
		if post.DislikeCount == -1 {
			return
		} else {
			fmt.Println(post)
			err := db.InsertPost(post)
			if err != nil {
				fmt.Println(err)
			}
			UpdateAndExecuteHome(w, r)
		}
	} else {
		ExecuteError(w, "Invalid User Method", http.StatusMethodNotAllowed)
	}
}
