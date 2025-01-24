package dbTools

// db.SelectComments() select all comments made in a post.
<<<<<<< HEAD
func (db *DBContainer) SelectComments(id int, orderBy string) ([]*Comment, error) {
	qry := `SELECT * FROM comments WHERE post_id = ?`
=======
func (db *DBContainer) SelectComments(id int, orderBy string) ([]Comment, error) {
	qry := `SELECT c.id, u.id, u.name, c.post_id, c.parent_id, c.content, 
				   c.like_count, c.dislike_count, c.created_at
			FROM comments c
			INNER JOIN users u ON c.user_id = u.id
			WHERE post_id = ?`
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
	orderByQry := ` ORDER BY created_at DESC`
	switch orderBy {
	case "oldest":
		orderByQry = ` ORDER BY created_at ASC`
	case "likeCount":
		orderByQry = ` ORDER BY like_count DESC`
	}
	qry += orderByQry

	rows, err := db.conn.Query(qry, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
<<<<<<< HEAD
	var comments []*Comment
=======
	var comments []Comment
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.ID,
			&c.UserID,
<<<<<<< HEAD
=======
			&c.UserName,
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
			&c.PostID,
			&c.ParentID,
			&c.Content,
			&c.LikeCount,
			&c.DislikeCount,
			&c.CreatedAt)
		if err != nil {
			return nil, err
		}
<<<<<<< HEAD
		comments = append(comments, &c)
=======
		comments = append(comments, c)
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return comments, nil
}

// db.InsertComment() inserts a comment for a post.
func (db *DBContainer) InsertComment(c Comment) error {
	qry := `INSERT INTO comments
<<<<<<< HEAD
			(user_id, post_id, parent_id, content, like_count, dislike_count, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`
=======
			(user_id, post_id, parent_id, content)
			VALUES (?, ?, ?, ?)`
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c

	var parentID interface{}
	if c.ParentID.Valid {
		parentID = c.ParentID.Int64
	} else {
		parentID = nil
	}
	_, err := db.conn.Exec(qry,
		c.UserID,
		c.PostID,
		parentID,
<<<<<<< HEAD
		c.Content,
		c.LikeCount,
		c.DislikeCount,
		c.CreatedAt)
=======
		c.Content)
>>>>>>> 8be9eaccc3a15efa4a6390295bf7249a6f455b5c
	return err
}

// db.UpdateComment() based on changes in comment
func (db *DBContainer) UpdateComment(c Comment) error {
	qry := `UPDATE comments	SET content = ?	WHERE id = ?`
	_, err := db.conn.Exec(qry,
		c.Content,
		c.ID)
	return err
}
