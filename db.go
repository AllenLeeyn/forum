package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// note: fields need to be validated before calling insert functions.
// The variables in the struct are initialized with default value,
// meaning they are not null when inserting to db.
// This means the variables/fields will not be recognise as empty/null by sql.

// dataBase struct comes with a set of functions.
// This should be easier to reference the database and call its functions.
type dataBase struct {
	conn       *sql.DB
	categories []string
}

// openDB() opens a sql database with the driver and dataSource given.
func openDB(driver, dataSource string) (*dataBase, error) {
	conn, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}
	conn.Exec("PRAGMA foreign_keys = ON;")
	return &dataBase{conn: conn}, nil
}

// check for no results
func checkErrNoRows(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// db.selectFieldFromTable() is a generic function to grab a column of data from a table.
func (db *dataBase) selectFieldFromTable(field, table string) ([]string, error) {
	rows, err := db.conn.Query("SELECT " + field + " FROM " + table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return values, nil
}

// db.selectUserByEmail(). If no results found, user is not registered/ wrong email.
func (db *dataBase) selectUserByEmail(email string) (*user, error) {
	qry := `SELECT id, type_id, name, email, pw_hash, reg_date, last_login 
			FROM users 
			WHERE email = ?`
	var u user
	err := db.conn.QueryRow(qry, email).Scan(
		&u.id,
		&u.typeID,
		&u.name,
		&u.email,
		&u.pwHash,
		&u.regDate,
		&u.lastLogin)
	if err != nil {
		return nil, checkErrNoRows(err)
	}
	return &u, nil
}

// db.inserUser() insert a user into the database
func (db *dataBase) insertUser(u *user) error {
	qry := `INSERT INTO users 
			(type_id, name, email, pw_hash, reg_date, last_login) 
			VALUES ( ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(qry,
		u.typeID,
		u.name,
		u.email,
		u.pwHash,
		u.regDate,
		u.lastLogin)
	return err
}

func (db *dataBase) updateUser(u *user) error {
	qry := `UPDATE users
			name = ?, pw_hash = ?, last_login = ?
			WHERE id = ?`
	_, err := db.conn.Exec(qry,
		u.name,
		u.pwHash,
		u.lastLogin)
	return err
}

// db.selectActiveSessionBy() id or user_id
func (db *dataBase) selectActiveSessionBy(field string, id interface{}) (*session, error) {
	if field != "id" && field != "user_id" {
		return nil, fmt.Errorf("invalid field")
	}
	var s session
	qry := `SELECT id, user_id, is_active, start_time, expire_time, last_access 
			FROM sessions 
			WHERE ` + field + ` = ? AND is_active = 1`
	err := db.conn.QueryRow(qry, id).Scan(
		&s.id,
		&s.userID,
		&s.isActive,
		&s.startTime,
		&s.expireTime,
		&s.lastAccess)
	if err != nil {
		return nil, checkErrNoRows(err)
	}
	return &s, nil
}

// db.insertSession() when user login is successful
func (db *dataBase) insertSession(s *session) error {
	qry := `INSERT INTO sessions
			(id, user_id, is_active, start_time, expire_time, last_access)
			VALUES ( ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(qry,
		s.id,
		s.userID,
		s.isActive,
		s.startTime,
		s.expireTime,
		s.lastAccess)
	return err
}

// db.updateSession() for when session is expired or refreshed
func (db *dataBase) updateSession(s *session) error {
	qry := `UPDATE session
			SET is_active = ?, expire_time = ?, last_access= ?
			WHERE id = ?`
	_, err := db.conn.Exec(qry,
		s.isActive,
		s.expireTime,
		s.lastAccess,
		s.id)
	return err
}

func getPostsQuery(filterBy, sortBy string, catID int) string {
	qry := `SELECT * FROM v_posts`
	return qry
}

// filter by categories, created posts and liked posts
// order by newest/oldest, comment count, like count
func (db *dataBase) selectPosts(filterBy, sortBy string, catID int) (*posts, error) {
	qry := getPostsQuery(filterBy, sortBy, catID)
	rows, err := db.conn.Query(qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts posts
	for rows.Next() {
		var p post
		var catIDs string
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.UserName,
			&p.CommentCount,
			&p.LikeCount,
			&p.DislikeCount,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&catIDs)
		if err != nil {
			return nil, err
		}
		p.categories, err = splitCategoryIDs(catIDs)
		if err != nil {
			return nil, err
		}
		posts.index = append(posts.index, p)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &posts, nil
}

func splitCategoryIDs(catIDs string) ([]int, error) {
	if catIDs == "" {
		return nil, fmt.Errorf("empty string")
	}
	var result []int
	categories := strings.Split(catIDs, ",")
	for _, idStr := range categories {
		if id, err := strconv.Atoi(idStr); err == nil {
			result = append(result, id)
		} else {
			return nil, err
		}
	}
	return result, nil
}

// needs categories
func (db *dataBase) insertPost(p post) error {
	qry := `INSERT INTO posts 
			(user_id, comment_count, like_count, dislike_count, title, content, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(qry,
		p.UserID,
		p.CommentCount,
		p.LikeCount,
		p.DislikeCount,
		p.Title,
		p.Content,
		p.CreatedAt)
	return err
}

func (db *dataBase) deleteAllUsers() error {
	query := "DELETE FROM users"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

func (db *dataBase) deleteAllSessions() error {
	query := "DELETE FROM sessions"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

func (db *dataBase) vacuumDB() error {
	query := "VACUUm"
	_, err := db.conn.Exec(query)
	return err
}
