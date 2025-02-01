package handlers

import (
	"net/http"
)

// Profile page
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in using session id
	sessionCookie, userID := checkSessionValidity(w, r)
	if userID == -1 {
		ExecuteError(w, "json", "Invalid session", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		posts, err := db.SelectPosts("createdBy", "", userID, userID)
		if err != nil {
			ExecuteError(w, "Tmpl", "Error getting user posts", http.StatusInternalServerError)
		}
		user, err := db.SelectUserByField("id", userID)
		if err != nil {
			ExecuteError(w, "Tmpl", "Error getting user details", http.StatusInternalServerError)
			return
		}
		data := profilepageData{
			Name:          user.Name,
			Email:         user.Email,
			Posts:         posts,
			SessionCookie: sessionCookie,
			Categories:    db.Categories,
		}
		extendSession(w, sessionCookie)
		ExecuteTmpl(w, "profile.html", data)
	} else {
		ExecuteError(w, "Tmpl", "Invalid User Method", http.StatusMethodNotAllowed)
	}
}
