package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID := checkSessionValidity(w, r)
	if userID == -1 {
		ExecuteError(w, "json", "Please login and try again", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		ExecuteError(w, "Tmpl", "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s := struct {
		Comment string
	}{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ExecuteError(w, "json", "Error reading body", http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(body, &s); err != nil {
		ExecuteError(w, "json", "Error reading json", http.StatusInternalServerError)
		return
	}

	postId := r.URL.Query().Get("postId")
	pId, err := strconv.Atoi(postId)
	if err != nil {
		ExecuteError(w, "Tmpl", "Erroneous post ID", http.StatusBadRequest)
		return
	}

	post, err := db.SelectPost(pId, userID)
	if err != nil || post == nil {
		ExecuteError(w, "json", "Post not found", http.StatusNotFound)
		return
	}

	user, err := db.SelectUserByField("id", userID)
	if err != nil || user == nil {
		ExecuteError(w, "json", "User not found", http.StatusNotFound)
		return
	}
	c := comment{
		UserID:   userID,
		PostID:   pId,
		Content:  s.Comment,
		UserName: user.Name}

	if err := db.InsertComment(c); err != nil {
		ExecuteError(w, "json", "Error creating comment", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	extendSession(w, sessionCookie)
}
