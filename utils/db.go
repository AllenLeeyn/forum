package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type dataBase struct {
	conn *sql.DB
}

func OpenDB(driver, dataSource string) (*dataBase, error) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}
	return &dataBase{conn: db}, nil
}

func (db *dataBase) SelectUserEmails() ([]string, error) {
	rows, err := db.conn.Query("SELECT email FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return emails, nil
}

func (db *dataBase) InsertNewUser(u *User) error {
	if (u.typeID < 1 || u.typeID > 3) ||
		u.name == "" ||
		u.email == "" ||
		u.pwHash == "" ||
		u.regDate.IsZero() ||
		u.lastLogin.IsZero() {
		return fmt.Errorf("ERROR: invalid data")
	}
	query := `INSERT INTO Users (type_id, name, email, pw_hash, reg_date, last_login) 
		VALUES ( ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query,
		u.typeID,
		u.name,
		u.email,
		u.pwHash,
		u.regDate,
		u.lastLogin)
	if err != nil {
		return err
	}
	return nil
}

func (db *dataBase) SelectUserByEmail(email string) (*User, error) {
	query := `SELECT type_id, name, email, pw_hash, reg_date, last_login 
		FROM Users 
		WHERE email = ?`

	var u User
	err := db.conn.QueryRow(query, email).Scan(
		&u.typeID,
		&u.name,
		&u.email,
		&u.pwHash,
		&u.regDate,
		&u.lastLogin)
	if err != nil {
		return nil, fmt.Errorf("failed to query User by email: %v", err)
	}
	return &u, nil
}

func (db *dataBase) SelectActiveSessionByUserID(UserID int) (*Session, error) {
	var s Session
	query := `SELECT id, User_id, is_active, start_time, expire_time, last_access 
		FROM sessions 
		WHERE User_id = ? AND is_active = 1`
	err := db.conn.QueryRow(query, UserID).Scan(
		&s.id,
		&s.userID,
		&s.isActive,
		&s.startTime,
		&s.expireTime,
		&s.lastAccess)
	if err != nil {
		return nil, err
	}
	return &s, err
}

func (db *dataBase) SelectPosts() (*Posts, error) {
	query := `SELECT p.id, p.User_id, u.name AS User_name, 
			p.comment_count, p.like_count, p.dislike_count, 
			p.title, p.content,p.created_at 
		FROM Posts p
		INNER JOIN Users u ON p.User_id = u.id`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts Posts

	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.UserName,
			&p.CommentCount,
			&p.LikeCount,
			&p.DislikeCount,
			&p.Title,
			&p.Content,
			&p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts.index = append(posts.index, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &posts, nil
}

func (db *dataBase) InsertPost(p Post) error {
	query := `
		INSERT INTO Posts (User_id, comment_count, like_count, dislike_count, title, content, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.Exec(
		query,
		p.UserID,
		p.CommentCount,
		p.LikeCount,
		p.DislikeCount,
		p.Title,
		p.Content,
		p.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (db *dataBase) DeleteAllUsers() error {
	query := "DELETE FROM Users"
	_, err := db.conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
