package dbTools

func (db *DBContainer) SelectUserByField(fieldName string, fieldValue interface{}) (*User, error) {
	qry := `SELECT * FROM users WHERE ` + fieldName + ` = ?`
	var u User
	err := db.conn.QueryRow(qry, fieldValue).Scan(
		&u.ID,
		&u.TypeID,
		&u.Name,
		&u.Email,
		&u.PwHash,
		&u.RegDate,
		&u.LastLogin)
	if err != nil {
		return nil, checkErrNoRows(err)
	}
	return &u, nil
}

// db.InserUser() insert a User into the database
func (db *DBContainer) InsertUser(u *User) error {
	qry := `INSERT INTO users 
			(id, type_id, name, email, pw_hash) 
			VALUES (?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(qry,
		u.ID,
		u.TypeID,
		u.Name,
		u.Email,
		u.PwHash)
	return err
}

// db.UpdateUser() info like name, pwHash and lastLogin
func (db *DBContainer) UpdateUser(u *User) error {
	qry := `UPDATE users
			SET name = ?, pw_hash = ?, last_login = ?
			WHERE id = ?`
	_, err := db.conn.Exec(qry,
		u.Name,
		u.PwHash,
		u.LastLogin,
		u.ID)
	return err
}
