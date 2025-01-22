package main

import (
	"fmt"
	"forum/dbTools"
	"forum/utils"
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
	utils.Init(db)
}

func main() {
	http.Handle("/static/", http.FileServer(http.Dir("assets/")))

	http.HandleFunc("/", utils.HomePage)
	http.HandleFunc("/login", utils.LoginPage)
	http.HandleFunc("/signup", utils.RegisterPage)
	http.HandleFunc("/terms", utils.TermsPage)
	http.HandleFunc("/start-thread", utils.StartThread)
	// http.HandleFunc("/post", utils.PostPage)
	// http.HandleFunc("/view-post", utils.ViewPostPage)

	fmt.Println("Starting Forum on http://localhost:8080/...")
	http.ListenAndServe(":8080", nil)
}
