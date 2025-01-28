package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Feedback(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID := checkSessionValidity(w, r)
	if userID == -1 { // likely user not login
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}

	userFeedback := struct {
		Tgt      string `json:"tgt"`
		ParentID int    `json:"parentID"`
		Rating   int    `json:"rating"`
	}{}
	body, err := io.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.Unmarshal(body, &userFeedback)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// need to add trigger for dislike count
	fb, err := db.SelectFeedback(userFeedback.Tgt, userID, userFeedback.ParentID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if fb == nil {
		err = db.InsertFeedback(userFeedback.Tgt,
			feedback{
				UserID:   userID,
				ParentID: userFeedback.ParentID,
				Rating:   userFeedback.Rating})
		if err != nil {
			ExecuteError(w, "something went wrong with feedback", http.StatusNotFound)
		}
	} else {
		fb.Rating = userFeedback.Rating
		err = db.UpdateFeedback(userFeedback.Tgt, *fb)
		if err != nil {
			ExecuteError(w, "something went wrong with feedback", http.StatusNotFound)
		}
	}
	extendSession(w, sessionCookie)
}
