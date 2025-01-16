package main

import (
	"fmt"
	"forum/utils"
	"net/http"
)

func main() {

	utils.InitializeHtml()
	http.HandleFunc("/", utils.HomePage)
	http.HandleFunc("/login", utils.LoginPage)
	http.HandleFunc("/signup", utils.SignUpPage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Starting Forum on http://localhost:8080/...")
	http.ListenAndServe(":8080", nil)
}
