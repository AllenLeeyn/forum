package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

// Page for making posts
func Post(w http.ResponseWriter, r *http.Request) {

	// check session from cookie to get feedback data
	sessionCookie, _ := r.Cookie("session-id")
	/* 	userID := checkSessionValidity(w, session)
	   	feedbacks, err := db.SelectFeedbacks("comment", userID)
	   	if err != nil {
	   		fmt.Println(err)
	   		// something went wrong
	   	} */

	if r.Method == http.MethodGet {

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("something went wrong with post_id")
		}

		post, err := db.SelectPost(id)
		if err != nil || post == nil {
			fmt.Println("something went wrong with getting post or nothing found")
		}

		comments, err := db.SelectComments(id, "oldest")
		if err != nil {
			fmt.Println("something went wrong with grabbing comments")
		}
		cookie, _ := r.Cookie("session-id")
		ExecuteTemp(w, "post.html", postpageData{cookie, *post, comments})

	} else if r.Method == http.MethodPost {
		// for create post maybe
	} else {
		http.Error(w, "Error 405, Method not allowed", http.StatusMethodNotAllowed)
	}
	extendSession(w, sessionCookie)
}

func StartThread(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie("session-id")
	userID := checkSessionValidity(w, sessionCookie)
	if userID == -1 { // likely user not login
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {

		ExecuteTemp(w, "start-thread.html", startThreadData{Categories: db.Categories})
		/* 	} else if r.Method == http.MethodPost {
		 */
	} else {
		ExecuteError(w, http.StatusMethodNotAllowed, "Invalid User Method")
	}
}
