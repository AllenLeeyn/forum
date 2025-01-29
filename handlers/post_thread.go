package handlers

import (
	"fmt"
	"forum/dbTools"
	"net/http"
	"strconv"
)

// Page for viewing individual post
func Post(w http.ResponseWriter, r *http.Request) {
	// check session from cookie to get feedback data
	sessionCookie, userID := checkSessionValidity(w, r)

	if r.Method == http.MethodGet {

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ExecuteError(w, "Tmpl", "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		post, err := db.SelectPost(id, userID)
		if err != nil || post == nil {
			ExecuteError(w, "Tmpl", "Error getting post: nothing found", http.StatusInternalServerError)
			return
		}
		comments, err := db.SelectComments(id, userID, "oldest")
		if err != nil {
			ExecuteError(w, "Tmpl", "Error getting comments: "+err.Error(), http.StatusInternalServerError)
			return
		}
		extendSession(w, sessionCookie)
		ExecuteTmpl(w, "post.html", postpageData{sessionCookie, *post, comments})
	} else {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Page for user to draft their post (if method == get), otherwise post it (if method == post)
func StartThread(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in using session id
	sessionCookie, userID := checkSessionValidity(w, r)
	if userID == -1 {
		ExecuteError(w, "json", "Please login and try again", http.StatusNotFound)
		return
	}

	// only allows page to display or accept request if user is logged in
	if r.Method == http.MethodGet {
		ExecuteTmpl(w, "start-thread.html", startThreadData{sessionCookie, db.Categories})
	} else if r.Method == http.MethodPost {
		title, content, categoriesInt, err := GetData(w, r)
		if err != nil {
			ExecuteError(w, "json", "Error reading form:"+err.Error(), 400)
			return
		}
		titleIsValid, titleError := CheckPostValidity(title, "postTitle")
		contentIsValid, contentError := CheckPostValidity(content, "postContent")
		if !(titleIsValid && contentIsValid) {
			ExecuteError(w, "json", "Error :"+titleError+contentError, 400)
			return
		}
		post := dbTools.Post{
			UserID:     userID,
			Title:      title,
			Content:    content,
			Categories: categoriesInt,
		}
		postNum, err := db.InsertPost(post)
		if err != nil {
			ExecuteError(w, "json", "Error creating post: "+err.Error(), http.StatusInternalServerError)
		}
		extendSession(w, sessionCookie)
		postID := "/post?id=" + strconv.Itoa(postNum)
		http.Redirect(w, r, postID, http.StatusSeeOther)

	} else {
		ExecuteError(w, "Tmpl", "Method not allowed", http.StatusMethodNotAllowed)
	}
}

/*----------- helper functions for posts -----------*/

// Checks the validity of posts (like length requirements)
func CheckPostValidity(input string, dataType string) (bool, string) {
	//	Add to these conditions later, including special character check
	if dataType == "postTitle" {
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

// Gets and parses data from front end post
func GetData(w http.ResponseWriter, r *http.Request) (title string, content string, categoriesInt []int, err error) {
	err = r.ParseForm()
	if err != nil {
		return "", "", nil, fmt.Errorf("error parsing form data")
	}
	title = r.FormValue("threadTitle")
	content = r.FormValue("threadContent")
	categoriesStr := r.Form["category"]
	for _, value := range categoriesStr {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return "", "", nil, fmt.Errorf("error parsing categories")
		}
		categoriesInt = append(categoriesInt, intVal)
	}
	return title, content, categoriesInt, nil
}
