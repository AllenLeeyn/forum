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
	http.HandleFunc("/", handlers.Home)

	http.HandleFunc("/signup", handlers.Signup)
	http.HandleFunc("/terms", handlers.Terms)
	http.HandleFunc("/about", handlers.About)

	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.LogOut)
	http.HandleFunc("/profile-page", handlers.Profile)

	http.HandleFunc("/post", handlers.ViewPost)
	http.HandleFunc("/profile", handlers.ViewProfile)
	http.HandleFunc("/new-post", handlers.NewPost)
	http.HandleFunc("/add-comment", handlers.Comment)

	http.HandleFunc("/feedback", handlers.Feedback)

	fmt.Println("Starting Forum on http://localhost:8080/...")
	log.Fatal(http.ListenAndServe(":8080", nil))
	db.Close()
}
