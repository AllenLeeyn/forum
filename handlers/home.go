package handlers

import (
	"fmt"
	"net/http"
	"strconv"
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
	// check session from cookie to get feedback data
	sessionCookie, _ := r.Cookie("session-id")
	userID, userName := checkSessionValidity(w, sessionCookie), "Guest"
	u, err := db.SelectUserByField("id", userID)
	if err == nil && u != nil {
		userName = u.Name
	}
	feedbacks, err := db.SelectFeedbacks("post", userID)
	if err != nil {
		fmt.Println(err)
		// something went wrong
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
		fmt.Println(err)
		// something went wrong
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
