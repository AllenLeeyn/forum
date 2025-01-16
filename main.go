package main

import (
	"fmt"
	"forum/utils"
	"net/http"
)

func main() {
	fmt.Println("Starting Forum on http://localhost:8080/...")
	utils.InitializeHtml()
	http.HandleFunc("/", utils.HomePage)
	http.HandleFunc("/login", utils.LoginPage)
	http.HandleFunc("/register", utils.RegisterPage)
	//http.HandleFunc("/post", utils.PostPage)
	//http.HandleFunc("/view-post", utils.ViewPostPage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
