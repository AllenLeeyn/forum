package main

import "fmt"

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
