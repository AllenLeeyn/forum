package main

import (
	"fmt"
	"forum/dbTools"
	"forum/handlers"
	"log"
	"net/http"
)

var db *dbTools.DBContainer

type user = dbTools.User
type session = dbTools.Session
type post = dbTools.Post
type feedback = dbTools.Feedback
type comment = dbTools.Comment

func init() {
	var err error
	db, err = dbTools.OpenDB("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	db.Categories, _ = db.SelectFieldFromTable("name", "categories")
	handlers.Init(db)
}

func main() {
	http.Handle("/static/", http.FileServer(http.Dir("assets/")))

	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/signup", handlers.SignupPage)
	http.HandleFunc("/logout", handlers.LogOut)
	http.HandleFunc("/terms", handlers.TermsPage)
	http.HandleFunc("/start-thread", handlers.StartThread)
	http.HandleFunc("/post-thread", handlers.PostThread)
	http.HandleFunc("/post", handlers.PostPage)
	http.HandleFunc("/view-post", handlers.ViewPostPage)

	fmt.Println("Starting Forum on http://localhost:8080/...")
	http.ListenAndServe(":8080", nil)
}
