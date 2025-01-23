package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//	Place for helper/util functions

/*
	 func PageHandler(w http.ResponseWriter, r *http.Request) {
		switch strings.ToLower(r.URL.Path) {

		case "/":
			IsError(w, r, http.MethodGet, "Incorrect User Method", http.StatusMethodNotAllowed, "method")
			HomePage(w, r)

		case "/login":
			IsError(w, r, http.MethodGet, "Incorrect User Method", http.StatusMethodNotAllowed, "method")
			LoginPage(w, r)

		case "/register":
			IsError(w, r, http.MethodGet, "Incorrect User Method", http.StatusMethodNotAllowed, "method")
			RegisterPage(w, r)

		case "/view-post":
			IsError(w, r, http.MethodGet, "Incorrect User Method", http.StatusMethodNotAllowed, "method")
			ViewPostPage(w, r)

		case "/post":
			IsError(w, r, http.MethodGet, "Incorrect User Method", http.StatusMethodNotAllowed, "method")
			PostPage(w, r)

		default:
			fmt.Println("Now we should throw an error.")

		}
	}
*/

// Custom execute template I wrote to both execute and check for error
func CustomExecuteTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func UpdateAndExecuteHome(w http.ResponseWriter) {
	posts, err := db.SelectPosts("", "", -1)
	if err != nil {
		fmt.Println(err)
		// something went wrong
	}
	PostList := postList{
		Posts: posts,
	}
	CustomExecuteTemplate(w, "homepage.html", PostList)
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
		if len(input) < 8 {
			return false, "*Title too short"
		} else if len(input) > 80 {
			return false, "*Title too long"
		}
	} else if dataType == "postContent" {
		if len(input) < 16 {
			return false, "*Content too short"
		} else if len(input) > 2000 {
			return false, "*Content too long"
		}
	}
	return true, ""
}
