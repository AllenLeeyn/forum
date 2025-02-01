package handlers

import (
	"net/http"
	"strconv"
)

// Profile page
func ViewProfilePage(w http.ResponseWriter, r *http.Request) {
	var viewer int
	var data profilepageData
	if r.Method == http.MethodGet {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ExecuteError(w, "Tmpl", "Erroneous profile ID", http.StatusBadRequest)
			return
		}
		user, err := db.SelectUserByField("id", id)
		if err != nil {
			ExecuteError(w, "Tmpl", "Error getting user details", http.StatusInternalServerError)
			return
		}
		sessionCookie, userID := checkSessionValidity(w, r)
		if sessionCookie != nil {
			extendSession(w, sessionCookie)
		}
		if userID != -1 {
			viewer = userID
		} else {
			viewer = -1
		}
		posts, err := db.SelectPosts("createdBy", "", id, viewer)
		if err != nil {
			ExecuteError(w, "Tmpl", "Error getting user posts", http.StatusInternalServerError)
		}
		if viewer == id {
			data = profilepageData{
				ProfileID:     id,
				ViewerID:      viewer,
				Name:          user.Name,
				Email:         user.Email,
				Posts:         posts,
				SessionCookie: sessionCookie,
				Categories:    db.Categories,
			}
		} else {
			data = profilepageData{
				ProfileID:  id,
				ViewerID:   viewer,
				Name:       user.Name,
				Posts:      posts,
				Categories: db.Categories,
			}
		}
		ExecuteTmpl(w, "profile.html", data)
	} else {
		ExecuteError(w, "Tmpl", "Invalid User Method", http.StatusMethodNotAllowed)
	}
}
