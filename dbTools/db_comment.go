package dbTools

// db.SelectComments() select all comments made in a post.
func (db *DBContainer) SelectComments(id int, orderBy string) ([]Comment, error) {
	qry := `SELECT c.id, u.id, u.name, c.post_id, c.parent_id, c.content, 
				   c.like_count, c.dislike_count, c.created_at
			FROM comments c
			INNER JOIN users u ON c.user_id = u.id
			WHERE post_id = ?`
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
	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.UserName,
			&c.PostID,
			&c.ParentID,
			&c.Content,
			&c.LikeCount,
			&c.DislikeCount,
			&c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return comments, nil
}

// db.InsertComment() inserts a comment for a post.
func (db *DBContainer) InsertComment(c Comment) error {
	qry := `INSERT INTO comments
			(user_id, user_name, post_id, parent_id, content, like_count)
			VALUES (?, ?, ?, ?, ?, ?)`

	var parentID interface{}
	if c.ParentID.Valid {
		parentID = c.ParentID.Int64
	} else {
		parentID = nil
	}
	_, err := db.conn.Exec(qry,
		c.UserID,
		c.UserName,
		c.PostID,
		parentID,
		c.Content,
		c.LikeCount)
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
