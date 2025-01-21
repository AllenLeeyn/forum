package main

import (
	"log"
)

var db *dataBase

func init() {
	var err error
	db, err = openDB("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	db.deleteAllUsers()

	db.categories, _ = db.selectFieldFromTable("name", "categories")
}

func main() {
}
