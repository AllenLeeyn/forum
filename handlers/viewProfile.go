package handlers

import (
	"net/http"
	"strconv"
)

// Profile page
func ViewProfile(w http.ResponseWriter, r *http.Request) {
	var data profilepageData
	if r.Method == http.MethodGet {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ExecuteError(w, "Tmpl", "Erroneous profile ID", http.StatusBadRequest)
			return
		}
		user, err := db.SelectUserByField("id", id)
		if err != nil || user == nil {
			ExecuteError(w, "Tmpl", "Error getting user details", http.StatusInternalServerError)
			return
		}
		sessionCookie, userID := checkSessionValidity(w, r)
		if sessionCookie != nil {
			extendSession(w, sessionCookie)
		}
		posts, err := db.SelectPosts("createdBy", "", id, userID)
		if err != nil {
			ExecuteError(w, "Tmpl", "Error getting user posts", http.StatusInternalServerError)
			return
		}
		data = profilepageData{
			ProfileID:     id,
			ViewerID:      userID,
			Name:          user.Name,
			Email:         user.Email,
			Posts:         posts,
			SessionCookie: sessionCookie,
			Categories:    db.Categories,
		}
		if userID == id {
			data.Email = user.Email
			data.SessionCookie = sessionCookie
		}
		ExecuteTmpl(w, "profile.html", data)
	} else {
		ExecuteError(w, "Tmpl", "Invalid User Method", http.StatusMethodNotAllowed)
		return
	}
}
