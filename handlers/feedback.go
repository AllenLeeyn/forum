package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func Feedback(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID := checkSessionValidity(w, r)
	if userID == -1 { // likely user not login
		ExecuteError(w, "json", "Please login and try again", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		ExecuteError(w, "Tmpl", "Method not allowed", http.StatusMethodNotAllowed)
	}

	userFeedback := struct {
		Tgt      string `json:"tgt"`
		ParentID int    `json:"parentID"`
		Rating   int    `json:"rating"`
	}{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ExecuteError(w, "json", "Error reading body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(body, &userFeedback); err != nil {
		ExecuteError(w, "json", "Error reading json: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fb, err := db.SelectFeedback(userFeedback.Tgt, userID, userFeedback.ParentID)
	if err != nil {
		ExecuteError(w, "json", "Error getting feedback: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if fb == nil {
		fb = &feedback{
			UserID:   userID,
			ParentID: userFeedback.ParentID,
			Rating:   userFeedback.Rating}
		if err = db.InsertFeedback(userFeedback.Tgt, *fb); err != nil {
			ExecuteError(w, "json", "Error giving feedback: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		fb.Rating = userFeedback.Rating
		if err = db.UpdateFeedback(userFeedback.Tgt, *fb); err != nil {
			ExecuteError(w, "json", "Error giving feedback: "+err.Error(), http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusOK)
	extendSession(w, sessionCookie)
}
