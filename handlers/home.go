package handlers

import (
	"net/http"
	"strconv"
)

// Home page
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ExecuteError(w, "Tmpl", "Page not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		ExecuteError(w, "Tmpl", "Method not allowed", http.StatusMethodNotAllowed)
	}

	// check session from cookie to get feedback data
	sessionCookie, userID := checkSessionValidity(w, r)
	userName := "Guest"
	u, err := db.SelectUserByField("id", userID)
	if err == nil && u != nil {
		userName = u.Name
	}

	// get Queries. No sorting implemented yet
	filterBy := r.URL.Query().Get("filterBy")
	orderBy := r.URL.Query().Get("orderBy")
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id > len(db.Categories) {
		id = -1
	}
	if filterBy != "category" {
		id = userID
	}
	posts, err := db.SelectPosts(filterBy, orderBy, id, userID)
	if err != nil {
		ExecuteError(w, "Tmpl", "Error getting posts: "+err.Error(), http.StatusInternalServerError)
	}
	extendSession(w, sessionCookie)
	ExecuteTmpl(w, "home.html",
		homepageData{
			sessionCookie,
			db.Categories,
			posts,
			userName,
			filterBy,
			orderBy,
			id})
}
