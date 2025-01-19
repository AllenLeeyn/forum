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

// checkErrNoRows() checks if no result from sql query.
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
	qry := `SELECT * FROM users WHERE email = ?`
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

// db.updateUser() info like name, pwHash and lastLogin
func (db *dataBase) updateUser(u *user) error {
	qry := `UPDATE users
			SET name = ?, pw_hash = ?, last_login = ?
			WHERE id = ?`
	_, err := db.conn.Exec(qry,
		u.name,
		u.pwHash,
		u.lastLogin,
		u.id)
	return err
}

// db.selectActiveSessionBy() id or user_id
func (db *dataBase) selectActiveSessionBy(field string, id interface{}) (*session, error) {
	if field != "id" && field != "user_id" {
		return nil, fmt.Errorf("invalid field")
	}
	var s session
	qry := `SELECT * FROM sessions WHERE ` + field + ` = ? AND is_active = 1`
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

// db.updateSession() for when session is expired, logout or refreshed
func (db *dataBase) updateSession(s *session) error {
	qry := `UPDATE sessions
			SET is_active = ?, expire_time = ?, last_access= ?
			WHERE id = ?`
	_, err := db.conn.Exec(qry,
		s.isActive,
		s.expireTime,
		s.lastAccess,
		s.id)
	return err
}

// db.isValidCategories*() check if given categories are valid with categories in db
func (db *dataBase) isValidCategories(categories []int) error {
	if len(categories) == 0 {
		return fmt.Errorf("no categories")
	}
	for _, catID := range categories {
		if catID == 0 || catID > len(db.categories)+1 {
			return fmt.Errorf("invalid category")
		}
	}
	return nil
}

// db.insetPost() into db and record the categories too
func (db *dataBase) insertPost(p post) error {
	if err := db.isValidCategories(p.categories); err != nil {
		return err
	}
	qry := `INSERT INTO posts 
			(user_id, comment_count, like_count, dislike_count, title, content, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := db.conn.Exec(qry,
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
	postID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	for _, catID := range p.categories {
		_, err = db.conn.Exec(`INSERT INTO post_categories (post_id, category_id)
								VALUES (?, ?)`, postID, catID)
	}
	return err
}

// getWhereQuery() for selectPosts() filterBy query
func getWhereQuery(filterBy string, id int) string {
	switch filterBy {
	case "createdBy":
		return fmt.Sprintf(` WHERE user_id = %v`, id)
	case "catergory":
		return fmt.Sprintf(` WHERE ',' || category_ids || ',' LIKE '%%,%v,%%'`, id)
	case "likedBy":
		return fmt.Sprintf(` INNER JOIN post_feedback pf ON pf.post_id = v_posts.id 
							 WHERE pf.user_id = %v AND pf.rating = 1`, id)
	}
	return ``
}

// getOrderByQuery() for selectPosts() orderBy query
func getOrderByQuery(orderBy string) string {
	switch orderBy {
	case "oldest":
		return ` ORDER BY created_at ASC`
	case "likeCount":
		return ` ORDER BY like_count DESC`
	case "commentCount":
		return ` ORDER BY comment_count DESC`
	}
	return ` ORDER BY created_at DESC`
}

// db.selectPosts() with filter and order options.
// by default, no filter and newest first are applied.
// if invalid options or empty are given, default option is used.
// valid filterBy: createdBy, catergory, likedBy
// valid orderBy: oldest, likeCount, commentCount
func (db *dataBase) selectPosts(filterBy, orderBy string, id int) (*posts, error) {
	qry := `SELECT * FROM v_posts` +
		getWhereQuery(filterBy, id) +
		getOrderByQuery(orderBy)
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

// splitCategoryIDs into []int to store in post struct
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

// db.deleteAllusers() for testing purposes
func (db *dataBase) deleteAllUsers() error {
	query := "DELETE FROM users"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

// db.deleteAllSessions() for testing purposes
func (db *dataBase) deleteAllSessions() error {
	query := "DELETE FROM sessions"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

// db.deleteAllPosts() for testing purposes
func (db *dataBase) deleteAllPosts() error {
	query := "DELETE FROM posts"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

// db.vacuumDB) for testing purposes
func (db *dataBase) vacuumDB() error {
	query := "VACUUm"
	_, err := db.conn.Exec(query)
	return err
}
