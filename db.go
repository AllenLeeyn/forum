package main

import (
	"database/sql"
	"fmt"

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
	query := "VACUUM"
	_, err := db.conn.Exec(query)
	return err
}
