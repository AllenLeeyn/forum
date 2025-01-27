package handlers

import "net/http"

func Feedback(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie("session-id")
	userID := checkSessionValidity(w, sessionCookie)
	if userID == -1 { // likely user not login
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
