package handlers

import (
	"forum/dbTools"
	"net/http"
	"strconv"
)

// Page for viewing individual post
func Post(w http.ResponseWriter, r *http.Request) {
	// check session from cookie to get feedback data
	sessionCookie, _ := checkSessionValidity(w, r)

	if r.Method == http.MethodGet {

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ExecuteError(w, "something went wrong with post_id", http.StatusInternalServerError)
			return
		}
		post, err := db.SelectPost(id)
		if err != nil || post == nil {
			ExecuteError(w, "something went wrong with getting post or nothing found", http.StatusInternalServerError)
			return
		}
		comments, err := db.SelectComments(id, "oldest")
		if err != nil {
			ExecuteError(w, "something went wrong with grabbing comments", http.StatusInternalServerError)
			return
		}
		extendSession(w, sessionCookie)
		ExecuteTemp(w, "post.html", postpageData{sessionCookie, *post, comments})
	} else {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Page for user to draft their post (if method == get), otherwise post it (if method == post)
func StartThread(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in using session id
	sessionCookie, userID := checkSessionValidity(w, r)
	if userID == -1 {
		// write header for toast message and do nothing
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// only allows page to display or accept request if user is logged in
	if r.Method == http.MethodGet {
		ExecuteTemp(w, "start-thread.html", startThreadData{sessionCookie, db.Categories})
	} else if r.Method == http.MethodPost {
		title, content, categoriesInt := GetData(w, r)
		if title == "" && content == "" && categoriesInt == nil {
			return
		}
		titleIsValid, titleInvalidReason := CheckPostValidity(title, "postTitle")
		contentIsValid, contentInvalidReason := CheckPostValidity(content, "postContent")
		if !(titleIsValid && contentIsValid) {
			errData := ErrorData{
				ErrorMessage: titleInvalidReason + "\n" + contentInvalidReason,
			}
			ExecuteTemp(w, "start-thread.html", errData)
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
			ExecuteError(w, "Failed to insert post into database", http.StatusInternalServerError)
		}
		extendSession(w, sessionCookie)
		postID := "/post?id=" + strconv.Itoa(postNum)
		http.Redirect(w, r, postID, http.StatusSeeOther)

	} else {
		ExecuteError(w, "Invalid User Method", http.StatusMethodNotAllowed)
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
