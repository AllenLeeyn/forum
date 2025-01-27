package handlers

import (
	"fmt"
	"forum/dbTools"
	"net/http"
	"strconv"
)

// Page for making posts
func Post(w http.ResponseWriter, r *http.Request) {

	// check session from cookie to get feedback data
	sessionCookie, _ := r.Cookie("session-id")
	/* 	userID := checkSessionValidity(w, session)
	   	feedbacks, err := db.SelectFeedbacks("comment", userID)
	   	if err != nil {
	   		fmt.Println(err)
	   		// something went wrong
	   	} */

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

		comments, err := db.SelectComments(id, "oldest")
		if err != nil {
			fmt.Println("something went wrong with grabbing comments")
		}
		cookie, _ := r.Cookie("session-id")
		ExecuteTemp(w, "post.html", postpageData{cookie, *post, comments})

	} else if r.Method == http.MethodPost {
		// for create post maybe
	} else {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}
	extendSession(w, sessionCookie)
}

// Page for user to draft their post (if method == get), otherwise post it (if method == post)
func StartThread(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ExecuteTemp(w, "start-thread.html", nil)
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
		post := CreatePost(w, r, title, content, categoriesInt)
		//	LikeCount in the negatives indicates an error
		if post.LikeCount == -1 {
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
