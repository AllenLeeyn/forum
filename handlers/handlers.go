package handlers

import (
	"fmt"
	"forum/dbTools"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//	Place for helper/util functions

// Custom execute template I wrote to both execute and check for error
func CustomExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func CheckValidity(input string, dataType string) (bool, string) {
	//	Add to these conditions later, including special character check
	if dataType == "username" {
		if len(input) > 25 {
			return false, "*Name too long"
		}
	} else if dataType == "password" {
		if len(input) < 6 {
			return false, "*Password too short"
		}
	} else if dataType == "email" {
		if !strings.Contains(input, "@") {
			return false, "*Invalid email"
		}
	} else if dataType == "postTitle" {
		if len(input) < 10 {
			return false, "*Title too short"
		} else if len(input) > 200 {
			return false, "*Title too long"
		}
	} else if dataType == "postContent" {
		if len(input) < 10 {
			return false, "*Content too short"
		} else if len(input) > 2000 {
			return false, "*Content too long"
		}
	}
	return true, ""
}

func GetSessionAndCookie(w http.ResponseWriter, r *http.Request) (*dbTools.Session, *http.Cookie) {
	cookie, err := r.Cookie("session-id")
	if err != nil {
		ExecuteError(w, "Session not found", http.StatusInternalServerError)
		return nil, nil
	}
	session, err := db.SelectActiveSessionBy("id", cookie.Value)
	if err != nil {
		ExecuteError(w, "Invalid session", http.StatusUnauthorized)
		return nil, cookie
	}
	return session, cookie
}

// Put together data into a post and return it
func CreatePost(w http.ResponseWriter, r *http.Request, title string, content string, categoriesInt []int) dbTools.Post {
	session, _ := GetSessionAndCookie(w, r)
	if session == nil {
		post := dbTools.Post{
			LikeCount: -1,
		}
		return post
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
	return post
}

func GetData(w http.ResponseWriter, r *http.Request) (title string, content string, categoriesInt []int) {
	err := r.ParseForm()
	if err != nil {
		ExecuteError(w, "Error parsing form data", http.StatusBadRequest)
		return "", "", nil
	}
	title = r.FormValue("threadTitle")
	content = r.FormValue("threadContent")
	categoriesStr := r.Form["category"]
	for _, value := range categoriesStr {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			ExecuteError(w, "Error parsing categories", http.StatusBadRequest)
			return "", "", nil
		}
		categoriesInt = append(categoriesInt, intVal)
	}
	return title, content, categoriesInt
}

// Execute error page if possible, otherwise use inbuilt http error
func ExecuteError(w http.ResponseWriter, errorMessage string, errorStatus int) {
	errorData := ErrorData{
		ErrorCode:    errorStatus,
		ErrorMessage: errorMessage,
	}
	err := tmpl.ExecuteTemplate(w, "error.html", errorData)
	if err != nil {
		message := fmt.Sprintf("Error %d\n%s", errorStatus, errorMessage)
		http.Error(w, message, http.StatusNotFound)
	}
}

func UpdateAndExecuteHome(w http.ResponseWriter, r *http.Request) {
	var posts []dbTools.Post
	var err error
	// check session from cookie
	session, cookie := GetSessionAndCookie(w, r)

	// get Queries. No sorting implemented yet
	filterBy := r.URL.Query().Get("filterBy")
	// orderBy := r.URL.Query().Get("orderBy")

	if strings.Contains(filterBy, "category") {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id > len(db.Categories) {
			id = -1
		}
		posts, err = db.SelectPosts(filterBy, "", id)
		if err != nil {
			ExecuteError(w, "Filtering posts went wrong", http.StatusInternalServerError)
			return
		}
	} else if strings.Contains(filterBy, "createdBy") {
		posts, err = db.SelectPosts("createdBy", "", session.UserID)
		if err != nil {
			ExecuteError(w, "Error getting user posts", http.StatusInternalServerError)
		}
	} else if strings.Contains(filterBy, "likedBy") {
		posts, err = db.SelectPosts("likedBy", "", session.UserID)
		if err != nil {
			ExecuteError(w, "Error getting user posts", http.StatusInternalServerError)
		}
	} else {
		posts, err = db.SelectPosts(filterBy, "", -1)
		if err != nil {
			ExecuteError(w, "Getting all posts failed", http.StatusInternalServerError)
			return
		}
	}

	CustomExecuteTemplate(w, "homepage.html", homepageData{posts, db.Categories, cookie})
}

/* SendUserToPost() {
	http.Redirect(w, r, "post", 200)
	CustomExecuteTemplate()
} */
