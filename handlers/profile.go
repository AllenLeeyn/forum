package handlers

import (
	"net/http"
	"strconv"
)

// Profile page
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, cookie := GetSessionAndCookie(w, r)
		if (cookie == nil) || (session == nil) {
			return
		}
		posts, err := db.SelectPosts("createdBy", "", session.UserID)
		if err != nil {
			ExecuteError(w, "Error getting user posts", http.StatusInternalServerError)
		}
		user, err := db.SelectUserByField("id", strconv.Itoa(session.UserID))
		if err != nil {
			ExecuteError(w, "Error getting user details", http.StatusInternalServerError)
			return
		}
		//	Placeholder data
		data := profilepageData{
			Name:          user.Name,
			Email:         user.Email,
			Posts:         posts,
			SessionCookie: cookie,
		}
		ExecuteTemp(w, "profile.html", data)
	} else {
		ExecuteError(w, "Invalid User Method", http.StatusMethodNotAllowed)
	}
}
