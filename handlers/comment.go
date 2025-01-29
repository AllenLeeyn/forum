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
		ExecuteError(w, "json", "Please login and try again", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		ExecuteError(w, "Tmpl", "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s := struct {
		Comment string
	}{}
	b, e := io.ReadAll(r.Body)
	checkErr(e) // need to handle error here?
	e = json.Unmarshal(b, &s)
	checkErr(e) // need to handle error here?

	postId := r.URL.Query().Get("postId")
	pId, err := strconv.Atoi(postId)
	checkErr(err) // need to handle error here?

	post, err := db.SelectPost(pId, userID)
	if err != nil || post == nil {
		ExecuteError(w, "json", "Post not found: "+err.Error(), http.StatusNotFound)
		return
	}

	user, err := db.SelectUserByField("id", userID)
	if err != nil || user == nil {
		ExecuteError(w, "json", "User not found: "+err.Error(), http.StatusNotFound)
		return
	}
	c := comment{
		UserID:   userID,
		PostID:   pId,
		Content:  s.Comment,
		UserName: user.Name}

	if err := db.InsertComment(c); err != nil {
		ExecuteError(w, "json", "Error creating comment: "+err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	extendSession(w, sessionCookie)
}
