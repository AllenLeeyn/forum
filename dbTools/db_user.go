package dbTools

// db.selectUserByEmail(). If no results found, user is not registered/ wrong email.
func (db *DBContainer) SelectUserByEmail(email string) (*user, error) {
	qry := `SELECT * FROM users WHERE email = ?`
	var u user
	err := db.conn.QueryRow(qry, email).Scan(
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

// db.inserUser() insert a user into the database
func (db *DBContainer) InsertUser(u *user) error {
	qry := `INSERT INTO users 
			(type_id, name, email, pw_hash, reg_date, last_login) 
			VALUES ( ?, ?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(qry,
		u.TypeID,
		u.Name,
		u.Email,
		u.PwHash,
		u.RegDate,
		u.LastLogin)
	return err
}

// db.updateUser() info like name, pwHash and lastLogin
func (db *DBContainer) UpdateUser(u *user) error {
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
