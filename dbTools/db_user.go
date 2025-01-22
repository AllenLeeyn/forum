package dbTools

import "strings"

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

func GetSubstringBeforeChar(input, char string) string {
	index := strings.Index(input, char)
	if index == -1 {
		return input
	}
	return input[:index]
}

// field example: "email:abc@def.gh"
// field example: "name:Bob"
func (db *DBContainer) SelectUserByField(field string) (*User, error) {
	fieldName := GetSubstringBeforeChar(field, ":")
	fieldValue := field[len(fieldName)+1:]
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
