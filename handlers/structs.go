package handlers

import (
	"forum/dbTools"
	"net/http"
)

type post = dbTools.Post

type homepageData struct {
	Posts         []post
	Categories    []string
	SessionCookie *http.Cookie
}

type ErrorData struct {
	ErrorMessage string
	ErrorCode    int
}

type PostError struct {
	ErrorMessage string
}
