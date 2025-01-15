package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type dataBase struct {
	conn *sql.DB
}

func openDB(driver, dataSource string) (*dataBase, error) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}
	return &dataBase{conn: db}, nil
}

func (db *dataBase) selectUserEmails() ([]string, error) {
	rows, err := db.conn.Query("SELECT email FROM users")
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

func (db *dataBase) insertNewUser(u *user) error {
	if (u.typeID < 1 || u.typeID > 3) ||
		u.name == "" ||
		u.email == "" ||
		u.pwHash == "" ||
		u.regDate.IsZero() ||
		u.lastLogin.IsZero() {
		return fmt.Errorf("ERROR: invalid data")
	}
	query := `INSERT INTO users (type_id, name, email, pw_hash, reg_date, last_login) 
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

func (db *dataBase) selectUserByEmail(email string) (*user, error) {
	query := `SELECT type_id, name, email, pw_hash, reg_date, last_login 
		FROM users 
		WHERE email = ?`

	var u user
	err := db.conn.QueryRow(query, email).Scan(
		&u.typeID,
		&u.name,
		&u.email,
		&u.pwHash,
		&u.regDate,
		&u.lastLogin)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %v", err)
	}
	return &u, nil
}

func (db *dataBase) selectActiveSessionByUserID(userID int) (*session, error) {
	var s session
	query := `SELECT id, user_id, is_active, start_time, expire_time, last_access 
		FROM sessions 
		WHERE user_id = ? AND is_active = 1`
	err := db.conn.QueryRow(query, userID).Scan(
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

func (db *dataBase) selectPosts() (*posts, error) {
	query := `SELECT p.id, p.user_id, u.name AS user_name, 
			p.comment_count, p.like_count, p.dislike_count, 
			p.title, p.content,p.created_at 
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts posts

	for rows.Next() {
		var p post
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

func (db *dataBase) insertPost(p post) error {
	query := `
		INSERT INTO posts (user_id, comment_count, like_count, dislike_count, title, content, created_at)
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

func (db *dataBase) deleteAllUsers() error {
	query := "DELETE FROM users"
	_, err := db.conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
