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
	cookie, _ := r.Cookie("session-id")
	CustomExecuteTemplate(w, "homepage.html", homepageData{posts, db.Categories, cookie})
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
