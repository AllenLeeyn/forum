package handlers

import (
	"fmt"
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
	// check session from cookie to get feedback data
	cookie, _ := r.Cookie("session-id")
	userID := checkSessionValidity(cookie)
	feedbacks, err := db.SelectFeedbacks("Post", userID)
	if err != nil {
		fmt.Println(err)
		// something went wrong
	}

	// get Queries. No sorting implemented yet
	filterBy := r.URL.Query().Get("filterBy")
	orderBy := r.URL.Query().Get("orderBy")
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id > len(db.Categories) {
		id = -1
	}

	posts, err := db.SelectPosts(filterBy, orderBy, id)
	if err != nil {
		fmt.Println(err)
		// something went wrong
	}
	CustomExecuteTemplate(w, "homepage.html",
		homepageData{
			cookie,
			db.Categories,
			posts,
			feedbacks,
			filterBy,
			orderBy,
			id})
}

func checkSessionValidity(c *http.Cookie) int {
	sessionID := c.Value
	fmt.Println(c)
	fmt.Println(sessionID)
	s, err := db.SelectActiveSessionBy("id", sessionID)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	if s.ExpireTime.Before(time.Now()) {
		s.IsActive = false
		if err := db.UpdateSession(s); err != nil {
			fmt.Println("update failed")
		}
		return -1
	}
	return s.UserID
}
