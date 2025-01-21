package main

import "fmt"

// db.selectFeedback() select a list of feedback that user has given.
// can be use to identify if user like a post/ comment and call insert/uodate accordingly.
// valid tgt: "post", "comment"
func (db *dataBase) selectFeedback(tgt string, userID int) (*[]feedback, error) {
	if tgt != "post" && tgt != "comment" {
		return nil, fmt.Errorf("invalid target")
	}
	qry := `SELECT * FROM ` + tgt + `_feedback WHERE user_id = ?`
	rows, err := db.conn.Query(qry, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []feedback
	for rows.Next() {
		var fb feedback
		err := rows.Scan(
			&fb.userID,
			&fb.parentID,
			&fb.rating,
			&fb.createdAt)
		if err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, fb)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &feedbacks, nil
}

// db.insertFeedback() inserts feedback into tgt table.
// valid tgt: "post", "comment"
func (db *dataBase) insertFeedback(tgt string, fb feedback) error {
	if tgt != "post" && tgt != "comment" {
		return fmt.Errorf("invalid target")
	}
	qry := `INSERT INTO ` + tgt + `_feedback 
			(user_id, parent_id, rating, created_at) 
			VALUES ( ?, ?, ?, ?)`
	_, err := db.conn.Exec(qry,
		fb.userID,
		fb.parentID,
		fb.rating,
		fb.createdAt)
	if err != nil {
		return err
	}
	return db.updateFeedbackCount(tgt, fb.parentID)
}

// db.updateFeedback() updates feedback in tgt table. for when user unlike.
// valid tgt: "post", "comment"
func (db *dataBase) updateFeedback(tgt string, fb feedback) error {
	if tgt != "post" && tgt != "comment" {
		return fmt.Errorf("invalid target")
	}
	qry := `UPDATE ` + tgt + `_feedback
			SET rating = ?, created_at = ? 
			WHERE user_id = ? AND parent_id = ?`
	_, err := db.conn.Exec(qry,
		fb.rating,
		fb.createdAt,
		fb.userID,
		fb.parentID)
	if err != nil {
		return err
	}
	return db.updateFeedbackCount(tgt, fb.parentID)
}

// db.updateFeedback() updates like_count and dislike_count
func (db *dataBase) updateFeedbackCount(tgt string, id int) error {
	qry := `UPDATE ` + tgt + `s
    	   SET like_count = (
        	   SELECT COUNT(*)
        	   FROM ` + tgt + `_feedback
        	   WHERE parent_id = ? AND rating = 1)
		   WHERE id = ?`
	_, err := db.conn.Exec(qry, id, id)
	return err
}
