package handlers

import (
	"net/http"
	"strconv"
)

// Profile page
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	// check if user is logged in using session id
	sessionCookie, userID := checkSessionValidity(w, r)
	if userID == -1 {
		ExecuteError(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		posts, err := db.SelectPosts("createdBy", "", userID)
		if err != nil {
			ExecuteError(w, "Error getting user posts", http.StatusInternalServerError)
		}
		user, err := db.SelectUserByField("id", strconv.Itoa(userID))
		if err != nil {
			ExecuteError(w, "Error getting user details", http.StatusInternalServerError)
			return
		}
		//	Placeholder data
		data := profilepageData{
			Name:          user.Name,
			Email:         user.Email,
			Posts:         posts,
			SessionCookie: sessionCookie,
		}
		extendSession(w, sessionCookie)
		ExecuteTemp(w, "profile.html", data)
	} else {
		ExecuteError(w, "Invalid User Method", http.StatusMethodNotAllowed)
	}
}
