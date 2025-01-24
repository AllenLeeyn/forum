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
<<<<<<< HEAD
	PwHash    string
=======
	PwHash    []byte
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
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
<<<<<<< HEAD
=======
	CatNames     string
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
}

type Comment struct {
	ID           int
	UserID       int
<<<<<<< HEAD
=======
	UserName     string
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
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
