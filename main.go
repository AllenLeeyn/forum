package main

import (
	"fmt"
	"log"
	"time"
)

var db *dataBase

func init() {
	var err error
	db, err = openDB("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	db.conn.Exec("PRAGMA foreign_keys = ON;")
	db.deleteAllUsers()
}

func main() {
	err := db.insertNewUser(&user{
		typeID:    3,
		name:      "allen",
		email:     "leeyn.shun@gmail.com",
		pwHash:    "password123",
		regDate:   time.Now(),
		lastLogin: time.Now()})
	if err != nil {
		fmt.Println(err)
	}

	err = db.insertNewUser(&user{
		typeID:    3,
		name:      "bruceWayne",
		email:     "bat.man@bat.cave",
		pwHash:    "password123",
		regDate:   time.Now(),
		lastLogin: time.Now()})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(db.selectUserEmails())
	user, err := db.selectUserByEmail("bat.man@bat.cave")
	fmt.Println(user)
	fmt.Println(user.lastLogin)
}
