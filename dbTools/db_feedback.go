package dbTools

import "fmt"

/*
	db.SelectFeedback() selects a list of Feedback that User has given.

Use to identify if User like a Post/ Comment and call insert/uodate accordingly.
Valid tgt: "Post", "Comment".
*/
func (db *DBContainer) SelectFeedbacks(tgt string, userID int) ([]Feedback, error) {
	if tgt != "post" && tgt != "comment" {
		return nil, fmt.Errorf("invalid target")
	}
	qry := `SELECT * FROM ` + tgt + `_feedback WHERE user_id = ?`
	rows, err := db.conn.Query(qry, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var fb Feedback
		err := rows.Scan(
			&fb.UserID,
			&fb.ParentID,
			&fb.Rating,
			&fb.CreatedAt)
		if err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, fb)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return feedbacks, nil
}

// db.insertFeedback() inserts Feedback into tgt table.
// valid tgt: "Post", "Comment"
func (db *DBContainer) InsertFeedback(tgt string, fb Feedback) error {
	if tgt != "post" && tgt != "comment" {
		return fmt.Errorf("invalid target")
	}
	qry := `INSERT INTO ` + tgt + `_feedback 
			(user_id, parent_id, rating) 
			VALUES ( ?, ?, ?)`
	_, err := db.conn.Exec(qry,
		fb.UserID,
		fb.ParentID,
		fb.Rating)
	return err
}

// db.updateFeedback() updates Feedback in tgt table. for when User unlike.
// valid tgt: "Post", "Comment"
func (db *DBContainer) UpdateFeedback(tgt string, fb Feedback) error {
	if tgt != "post" && tgt != "comment" {
		return fmt.Errorf("invalid target")
	}
	qry := `UPDATE ` + tgt + `_feedback
			SET rating = ?, created_at = ? 
			WHERE user_id = ? AND parent_id = ?`
	_, err := db.conn.Exec(qry,
		fb.Rating,
		fb.CreatedAt,
		fb.UserID,
		fb.ParentID)
	return err
}
