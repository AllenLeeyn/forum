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
		ExecuteError(w, "Error 404, Page not found", http.StatusNotFound)
		return
	}
	UpdateAndExecuteHome(w, r)
}

// Page for viewing posts
func ViewPostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			ExecuteError(w, "Post ID is missing", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ExecuteError(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}
		post, err := db.SelectPost(id)
		if err != nil || post == nil {
			ExecuteError(w, "Post not found", http.StatusNotFound)
			return
		}
		comments, err := db.SelectComments(id, "")
		if err != nil {
			ExecuteError(w, "Failed to load comments", http.StatusInternalServerError)
			return
		}
		cookie, _ := r.Cookie("session-id")
		CustomExecuteTemplate(w, "post.html", postpageData{*post, comments, cookie})
	} else {
		ExecuteError(w, "Invalid Request Method", http.StatusMethodNotAllowed)
	}
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
			postNum, err := db.InsertPost(post)
			if err != nil {
				ExecuteError(w, "Failed to insert post into database", http.StatusInternalServerError)
			}
			postID := "/post?id=" + strconv.Itoa(postNum)
			http.Redirect(w, r, postID, http.StatusSeeOther)
		}
	} else {
		ExecuteError(w, "Invalid User Method", http.StatusMethodNotAllowed)
	}
}
