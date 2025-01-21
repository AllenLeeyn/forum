package dbTools

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int
	TypeID    int
	Name      string
	Email     string
	PwHash    string
	RegDate   time.Time
	LastLogin time.Time
}

type Session struct {
	ID         string
	UserID     int
	IsActive   bool
	StartTime  time.Time
	ExpireTime time.Time
	LastAccess time.Time
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
	Categories   []int
}

type Comment struct {
	ID           int
	UserID       int
	ParentID     sql.NullInt64
	PostID       int
	Content      string
	LikeCount    int
	DislikeCount int
	CreatedAt    time.Time
}

type Feedback struct {
	UserID    int
	ParentID  int
	Rating    int
	CreatedAt time.Time
}
