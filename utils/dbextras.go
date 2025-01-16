package utils

import (
	"fmt"
	"log"
	"time"
)

var db *dataBase

func init() {
	var err error
	db, err = OpenDB("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	db.conn.Exec("PRAGMA foreign_keys = ON;")
	db.DeleteAllUsers()
}

func DbTest() {
	err := db.InsertNewUser(&User{
		typeID:    3,
		name:      "allen",
		email:     "leeyn.shun@gmail.com",
		pwHash:    "password123",
		regDate:   time.Now(),
		lastLogin: time.Now()})
	if err != nil {
		fmt.Println(err)
	}

	err = db.InsertNewUser(&User{
		typeID:    3,
		name:      "bruceWayne",
		email:     "bat.man@bat.cave",
		pwHash:    "password123",
		regDate:   time.Now(),
		lastLogin: time.Now()})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(db.SelectUserEmails())
	user, err := db.SelectUserByEmail("bat.man@bat.cave")
	if err != nil {
		fmt.Println("Error getting batman:", err)
	}
	fmt.Println(user)
	fmt.Println(user.lastLogin)
}
