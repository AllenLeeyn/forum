package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func Feedback(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID := checkSessionValidity(w, r)
	if userID == -1 { // likely user not login
		ExecuteError(w, "json", "Please login and try again", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		ExecuteError(w, "Tmpl", "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userFeedback := struct {
		Tgt      string `json:"tgt"`
		ParentID int    `json:"parentID"`
		Rating   int    `json:"rating"`
	}{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ExecuteError(w, "json", "Error reading body", http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(body, &userFeedback); err != nil {
		ExecuteError(w, "json", "Error reading json", http.StatusInternalServerError)
		return
	}

	fb, err := db.SelectFeedback(userFeedback.Tgt, userID, userFeedback.ParentID)
	if err != nil {
		ExecuteError(w, "json", "Error getting feedback", http.StatusInternalServerError)
		return
	} else if fb == nil {
		fb = &feedback{
			UserID:   userID,
			ParentID: userFeedback.ParentID,
			Rating:   userFeedback.Rating}
		if err = db.InsertFeedback(userFeedback.Tgt, *fb); err != nil {
			ExecuteError(w, "json", "Error giving feedback", http.StatusInternalServerError)
			return
		}
	} else {
		fb.Rating = userFeedback.Rating
		if err = db.UpdateFeedback(userFeedback.Tgt, *fb); err != nil {
			ExecuteError(w, "json", "Error giving feedback", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	extendSession(w, sessionCookie)
}
