package main

import (
	"fmt"
	"forum/dbTools"
<<<<<<< HEAD
	"forum/utils"
=======
	"forum/handlers"
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
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
<<<<<<< HEAD
	utils.Init(db)
=======
	handlers.Init(db)
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
}

func main() {
	http.Handle("/static/", http.FileServer(http.Dir("assets/")))

<<<<<<< HEAD
	http.HandleFunc("/", utils.HomePage)
	http.HandleFunc("/login", utils.LoginPage)
	http.HandleFunc("/signup", utils.RegisterPage)
	http.HandleFunc("/terms", utils.TermsPage)
	http.HandleFunc("/start-thread", utils.StartThread)
	http.HandleFunc("/post-thread", utils.PostThread)
	// http.HandleFunc("/view-post", utils.ViewPostPage)
=======
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/signup", handlers.SignupPage)
	http.HandleFunc("/logout", handlers.LogOut)
	http.HandleFunc("/terms", handlers.TermsPage)
	http.HandleFunc("/start-thread", handlers.StartThread)
	http.HandleFunc("/post/", handlers.ViewPostPage)
	http.HandleFunc("/profile-page", handlers.ProfilePage)
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c

	fmt.Println("Starting Forum on http://localhost:8080/...")
	http.ListenAndServe(":8080", nil)
}
