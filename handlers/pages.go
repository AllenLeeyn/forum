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

// Page for user to draft their post
func StartThread(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CustomExecuteTemplate(w, "start-thread.html", nil)
	} else {
		ExecuteError(w, "Invalid User Method", http.StatusMethodNotAllowed)
	}
}

// Page that handles and uploads post made by user
func PostThread(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		query := r.URL.Query()
		title := query.Get("threadTitle")
		content := query.Get("threadContent")
		categoriesStr := query["category"]
		var categoriesInt []int
		for _, value := range categoriesStr {
			intVal, err := strconv.Atoi(value)
			if err != nil {
				ExecuteError(w, "Error parsing categories", http.StatusBadRequest)
			}
			categoriesInt = append(categoriesInt, intVal)
		}
		fmt.Println(categoriesInt)
		titleIsValid, titleInvalidReason := CheckValidity(title, "postTitle")
		contentIsValid, contentInvalidReason := CheckValidity(content, "postContent")
		if titleIsValid && contentIsValid {
			fmt.Println("Valid post, check if duplicate")
		} else {
			//	Placeholder:
			errData := ErrorData{
				ErrorMessage: titleInvalidReason + "\n" + contentInvalidReason,
			}
			CustomExecuteTemplate(w, "start-thread.html", errData)
			return
		}
		cookie, err := r.Cookie("session-id")
		if err != nil {
			ExecuteError(w, "Session not found", http.StatusInternalServerError)
			return
		}
		session, err := db.SelectActiveSessionBy("id", cookie.Value)
		if err != nil {
			ExecuteError(w, "Invalid session", http.StatusUnauthorized)
			return
		}
		//	Placeholder:
		post := dbTools.Post{
			UserID:       session.UserID,
			CommentCount: 0,
			LikeCount:    0,
			DislikeCount: 0,
			Title:        title,
			Content:      content,
			Categories:   categoriesInt,
		}
		fmt.Println(post)
		err = db.InsertPost(post)
		if err != nil {
			fmt.Println(err)
		}
		UpdateAndExecuteHome(w, r)
	} else {
		ExecuteError(w, "Invalid User Method", http.StatusMethodNotAllowed)
	}
}
