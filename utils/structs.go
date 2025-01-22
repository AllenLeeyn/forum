package utils

import "time"

// Place structs and any global variables/data in here
type RegisterData struct {
	username string
	password string
	email    string
}

type LoginData struct {
	username string
	password string
}

type User struct {
	typeID    int
	name      string
	email     string
	pwHash    string
	regDate   time.Time
	lastLogin time.Time
}

type Session struct {
	id         string
	userID     int
	isActive   bool
	startTime  time.Time
	expireTime time.Time
	lastAccess time.Time
}

type Post struct {
	ID           int
	UserID       int
	UserName     string
	CommentCount int
	LikeCount    int
	DislikeCount int
	Title        string
	Content      string
	CreatedAt    time.Time
}

type Posts struct {
	index []Post
}

type ErrorData struct {
	ErrorMessage string
	ErrorCode    int
}
