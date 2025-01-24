package handlers

import (
	"forum/dbTools"
	"net/http"
)

// type user = dbTools.User
type session = dbTools.Session
type post = dbTools.Post

// type feedback = dbTools.Feedback
type comment = dbTools.Comment

type homepageData struct {
	Posts         []post
	Categories    []string
	SessionCookie *http.Cookie
}

type postpageData struct {
	Post          post
	Comments      []comment
	SessionCookie *http.Cookie
}

type ErrorData struct {
	ErrorMessage string
	ErrorCode    int
}

type PostError struct {
	ErrorMessage string
}
