package structs

import (
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	Posts          int64     `json:"posts"`
	Comments       int64     `json:"comments"`
	Likes          int64     `json:"likes"`
	Dislikes       int64     `json:"dislikes"`
	RecentActivity *Post     `json:"recent_activity"`
}

type Session struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	Status   string `json:"status"`
}

type Post struct {
	ID            int64    `json:"id"`
	Title         string   `json:"title"`
	Content       string   ` json:"content"`
	UserID        int64    ` json:"user_id"`
	CreatedAt     string   ` json:"created_at"`
	Author        string   `json:"author"`
	TotalLikes    int64    `json:"total_likes"`
	TotalDislikes int64    `json:"total_dislikes"`
	TotalComments int64    `json:"total_comments"`
	Categories    []string `json:"categories"`
}

var PostsShowing []Post

type Comment struct {
	ID            int64  `json:"id"`
	Content       string `json:"content"`
	UserID        int64  `json:"user_id"`
	PostID        int64  `json:"post_id"`
	CreatedAt     string `json:"created_at"`
	Author        string `json:"author"`
	TotalLikes    int64  `json:"total_likes"`
	TotalDislikes int64  `json:"total_dislikes"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Page    string `json:"page"`
	Path    string `json:"path"`
}