package main

import (
	"forum/dbTools"
	"log"
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
}
