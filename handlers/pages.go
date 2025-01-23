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
	UpdateAndExecuteHome(w, r)
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
		cookie, _ := r.Cookie("session-id")
		CustomExecuteTemplate(w, "post.html", postpageData{*post, comments, cookie})

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
			http.Error(w, "Session not found", http.StatusInternalServerError)
			return
		}
		session, err := db.SelectActiveSessionBy("id", cookie.Value)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
		fmt.Println("No error getting user by ID1")
		//	Placeholder:
		categories := []int{1, 2, 3}
		fmt.Println("No error getting user by ID2")
		post := dbTools.Post{
			UserID:       session.UserID,
			CommentCount: 0,
			LikeCount:    0,
			DislikeCount: 0,
			Title:        title,
			Content:      content,
			Categories:   categories,
		}
		fmt.Println("No error getting user by ID3")
		fmt.Println(post)
		fmt.Println("No error getting user by ID4")
		err = db.InsertPost(post)
		fmt.Println("No error getting user by ID5")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("No error getting user by ID")
		UpdateAndExecuteHome(w, r)
	} else {
		ExecuteError(w, http.StatusMethodNotAllowed, "Invalid User Method")
	}
}
