package dbTools

// db.SelectUserByEmail(). If no results found, User is not registered/ wrong email.
func (db *DBContainer) SelectUserByEmail(email string) (*User, error) {
	qry := `SELECT * FROM users WHERE email = ?`
	var u User
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

<<<<<<< HEAD
// db.InserUser() insert a User into the database
func (db *DBContainer) InsertUser(u *User) error {
	qry := `INSERT INTO users 
			(type_id, name, email, pw_hash, reg_date, last_login) 
			VALUES ( ?, ?, ?, ?, ?, ?)`
=======
func (db *DBContainer) SelectUserByField(fieldName, fieldValue string) (*User, error) {
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
			(type_id, name, email, pw_hash) 
			VALUES ( ?, ?, ?, ?)`
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
	_, err := db.conn.Exec(qry,
		u.TypeID,
		u.Name,
		u.Email,
<<<<<<< HEAD
		u.PwHash,
		u.RegDate,
		u.LastLogin)
=======
		u.PwHash)
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
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
