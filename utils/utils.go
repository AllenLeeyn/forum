package utils

import (
	"log"
	"net/http"
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
