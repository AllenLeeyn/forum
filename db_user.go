package main

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
