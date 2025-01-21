package main

import "time"

type user struct {
	id        int
	typeID    int
	name      string
	email     string
	pwHash    string
	regDate   time.Time
	lastLogin time.Time
}

type session struct {
	id         string
	userID     int
	isActive   bool
	startTime  time.Time
	expireTime time.Time
	lastAccess time.Time
}

type post struct {
	ID           int
	UserID       int
	UserName     string
	CommentCount int
	LikeCount    int
	DislikeCount int
	Title        string
	Content      string
	CreatedAt    time.Time
	categories   []int
}

type feedback struct {
	userID    int
	parentID  int
	rating    int
	createdAt time.Time
}
