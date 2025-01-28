package handlers

import "net/http"

func Feedback(w http.ResponseWriter, r *http.Request) {
	_, userID := checkSessionValidity(w, r)
	if userID == -1 { // likely user not login
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
