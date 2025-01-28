package handlers

import (
	"net/http"
	"strconv"
)

// Home page
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Error 404, Page not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}

	// check session from cookie to get feedback data
	sessionCookie, userID := checkSessionValidity(w, r)
	userName := "Guest"
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

	id, err := strconv.Atoi(idStr)
	if err != nil || id > len(db.Categories) {
		id = -1
	}
	posts, err := db.SelectPosts(filterBy, orderBy, id)
	if err != nil {
		ExecuteError(w, "Problem occur when getting data", http.StatusInternalServerError)
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
