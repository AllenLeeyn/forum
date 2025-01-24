package handlers

import (
	"forum/dbTools"
	"net/http"
)

type user = dbTools.User
type session = dbTools.Session
type post = dbTools.Post
type feedback = dbTools.Feedback
type comment = dbTools.Comment

type homepageData struct {
	SessionCookie *http.Cookie
	Categories    []string
	Posts         []post
	Feedbacks     []feedback
	FilterBy      string
	OrderBy       string
	Id            int
}

type postpageData struct {
	SessionCookie *http.Cookie
	Post          post
	Comments      []comment
}

type ErrorData struct {
	ErrorMessage string
	ErrorCode    int
}

type PostError struct {
	ErrorMessage string
}
