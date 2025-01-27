package handlers

import (
	"forum/dbTools"
	"net/http"
	"strconv"
	"strings"
)

// Home page
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Error 404, Page not found", http.StatusNotFound)
		return
	}
	UpdateAndExecuteHome(w, r)
}

func UpdateAndExecuteHome(w http.ResponseWriter, r *http.Request) {
	var posts []dbTools.Post
	var err error
	id := -1
	// check session from cookie to get feedback data
	session, sessionCookie := GetSessionAndCookie(w, r)
	userID, userName := checkSessionValidity(w, sessionCookie), "Guest"
	u, err := db.SelectUserByField("id", userID)
	if err == nil && u != nil {
		userName = u.Name
	}
	feedbacks, err := db.SelectFeedbacks("post", userID)
	if err != nil {
		ExecuteError(w, "Something went wrong fetching comments", http.StatusInternalServerError)
		return
	}

	// get Queries. No sorting implemented yet
	filterBy := r.URL.Query().Get("filterBy")
	orderBy := r.URL.Query().Get("orderBy")
	idStr := r.URL.Query().Get("id")

	if strings.Contains(filterBy, "category") {
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
	extendSession(w, sessionCookie)
	ExecuteTemp(w, "home.html",
		homepageData{
			sessionCookie,
			db.Categories,
			posts,
			userName,
			feedbacks,
			filterBy,
			orderBy,
			id})
}
