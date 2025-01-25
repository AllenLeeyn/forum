package handlers

import (
	"encoding/json"
	"fmt"
	"forum/dbTools"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//	Place for helper/util functions

// Custom execute template I wrote to both execute and check for error
func CustomExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func CheckValidity(input string, dataType string) (bool, string) {
	//	Add to these conditions later, including special character check
	if dataType == "username" {
		if len(input) > 25 {
			return false, "*Name too long"
		}
	} else if dataType == "password" {
		if len(input) < 6 {
			return false, "*Password too short"
		}
	} else if dataType == "email" {
		if !strings.Contains(input, "@") {
			return false, "*Invalid email"
		}
	} else if dataType == "postTitle" {
		if len(input) < 10 {
			return false, "*Title too short"
		} else if len(input) > 200 {
			return false, "*Title too long"
		}
	} else if dataType == "postContent" {
		if len(input) < 10 {
			return false, "*Content too short"
		} else if len(input) > 2000 {
			return false, "*Content too long"
		}
	}
	return true, ""
}

// Execute error page if possible, otherwise use inbuilt http error
func ExecuteError(w http.ResponseWriter, errorStatus int, errorMessage string) {
	errorData := ErrorData{
		ErrorCode:    errorStatus,
		ErrorMessage: errorMessage,
	}
	err := tmpl.ExecuteTemplate(w, "error.html", errorData)
	if err != nil {
		message := fmt.Sprintf("Error %d\n%s", errorStatus, errorMessage)
		http.Error(w, message, http.StatusNotFound)
	}
}

func UpdateAndExecuteHome(w http.ResponseWriter, r *http.Request) {
	// check session from cookie

	// get Queries. No sorting implemented yet
	filterBy := r.URL.Query().Get("filterBy")
	// orderBy := r.URL.Query().Get("orderBy")

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id > len(db.Categories) {
		id = -1
	}

	posts, err := db.SelectPosts(filterBy, "", id)
	if err != nil {
		fmt.Println(err)
		// something went wrong
	}
	cookie, _ := r.Cookie("session-id")
	CustomExecuteTemplate(w, "homepage.html", homepageData{posts, db.Categories, cookie})
}

var commCount int = 1

func Comment(w http.ResponseWriter, r *http.Request) {
	sessionCookie, e := r.Cookie("session-id")
	checkErr(e)

	s := struct {
		Comment string
	}{}
	b, e := io.ReadAll(r.Body)
	checkErr(e)
	e = json.Unmarshal(b, &s)
	checkErr(e)

	query := r.URL.Query()
	postId := query.Get("postId")
	pId, err := strconv.Atoi(postId)
	checkErr(err)

	session, e := db.SelectActiveSessionBy("id", sessionCookie.Value)
	checkErr(e)

	user, e := db.SelectUserByField("id", strconv.Itoa(session.UserID))
	checkErr(e)

	db.InsertComment(dbTools.Comment{
		ID:        commCount,
		UserID:    session.UserID,
		PostID:    pId,
		Content:   s.Comment,
		UserName:  user.Name,
		LikeCount: 7,
		CreatedAt: time.Now(),
	})
	commCount++

	post, err := db.SelectPost(pId)
	if err != nil || post == nil {
		fmt.Println("something went wrong with getting post or nothing found")
	}
}
