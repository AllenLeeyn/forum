package dbTools

import "fmt"

// db.selectActiveSessionBy() id or user_id
func (db *DBContainer) SelectActiveSessionBy(field string, id interface{}) (*session, error) {
	if field != "id" && field != "user_id" {
		return nil, fmt.Errorf("invalid field")
	}
	var s session
	qry := `SELECT * FROM sessions WHERE ` + field + ` = ? AND is_active = 1`
	err := db.conn.QueryRow(qry, id).Scan(
		&s.ID,
		&s.UserID,
		&s.IsActive,
		&s.StartTime,
		&s.ExpireTime,
		&s.LastAccess)
	if err != nil {
		return nil, checkErrNoRows(err)
	}
	return &s, nil
}

// db.insertSession() when user login is successful
func (db *DBContainer) InsertSession(s *session) error {
	qry := `INSERT INTO sessions
			(id, user_id, is_active, start_time, expire_time, last_access)
			VALUES ( ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(qry,
		s.ID,
		s.UserID,
		s.IsActive,
		s.StartTime,
		s.ExpireTime,
		s.LastAccess)
	return err
}

// db.updateSession() for when session is expired, logout or refreshed
func (db *DBContainer) UpdateSession(s *session) error {
	qry := `UPDATE sessions
			SET is_active = ?, expire_time = ?, last_access= ?
			WHERE id = ?`
	_, err := db.conn.Exec(qry,
		s.IsActive,
		s.ExpireTime,
		s.LastAccess,
		s.ID)
	return err
}
