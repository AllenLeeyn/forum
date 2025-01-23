package main

import (
	"fmt"
	"forum/dbTools"
	"log"
	"net/http"
)

var db *dbTools.DBContainer

func init() {
	var err error
	db, err = dbTools.OpenDB("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	db.DeleteAllUsers()

	db.Categories, _ = db.SelectFieldFromTable("name", "categories")
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)

	fmt.Println("server started")
	http.ListenAndServe(":8080", nil)
}
