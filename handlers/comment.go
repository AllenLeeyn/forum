package handlers

import (
	"encoding/json"
	"fmt"
	"forum/dbTools"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie("session-id")
	userID := checkSessionValidity(w, sessionCookie)
	if userID == -1 { // likely user not login
		w.WriteHeader(http.StatusNotFound)
		return
	}

	s := struct {
		Comment string
	}{}
	b, e := io.ReadAll(r.Body)
	checkErr(e)
	e = json.Unmarshal(b, &s)
	checkErr(e)

	postId := r.URL.Query().Get("postId")
	pId, err := strconv.Atoi(postId)
	checkErr(err)

	post, err := db.SelectPost(pId)
	if err != nil || post == nil {
		fmt.Println("something went wrong with getting post or nothing found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := db.SelectUserByField("id", userID)
	if err != nil || user == nil {
		fmt.Println("something went wrong with getting user or nothing found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db.InsertComment(dbTools.Comment{
		UserID:    userID,
		PostID:    pId,
		Content:   s.Comment,
		UserName:  user.Name,
		LikeCount: 7,
		CreatedAt: time.Now(),
	})

}
