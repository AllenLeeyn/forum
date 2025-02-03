package dbTools

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// getWhereQuery() for selectPosts() filterBy query
func getWhereQuery(filterBy string, id int) string {
	switch filterBy {
	case "createdBy":
		return fmt.Sprintf(` WHERE puser_id = %v`, id)
	case "category":
		return fmt.Sprintf(` WHERE ',' || category_ids || ',' LIKE '%%,%v,%%'`, id)
	case "likedBy":
		return fmt.Sprintf(` WHERE pf.user_id = %v AND pf.rating = 1`, id)
	}
	return ``
}

// getOrderByQuery() for selectPosts() orderBy query
func getOrderByQuery(orderBy string) string {
	switch orderBy {
	case "oldest":
		return ` ORDER BY pcreated_at ASC`
	case "likeCount":
		return ` ORDER BY like_count DESC`
	case "commentCount":
		return ` ORDER BY comment_count DESC`
	}
	return ` ORDER BY pcreated_at DESC`
}

// splitCategoryIDs() into []int to store in Post struct
func (db *DBContainer) splitCategoryIDs(catIDs string) ([]int, string, error) {
	if catIDs == "" {
		return nil, "", fmt.Errorf("empty string")
	}
	var result []int
	var resultNames string
	categories := strings.Split(catIDs, ",")
	for _, idStr := range categories {
		if id, err := strconv.Atoi(idStr); err == nil {
			result = append(result, id)
			resultNames = resultNames + db.Categories[id] + ", "
		} else {
			return nil, "", err
		}
	}
	return result, resultNames[:len(resultNames)-2], nil
}

/*
	db.SelectPosts() returns a list of posts.

By default, no filter and newest first are applied.
If invalid options or empty are given, default option is used.

Valid filterBy: createdBy, category, likedBy.
Valid orderBy: oldest, likeCount, commentCount.
*/
func (db *DBContainer) SelectPosts(filterBy, orderBy string, id, userID int) ([]Post, error) {
	qry := `SELECT v_posts.id, puser_id, user_name, 
				comment_count, like_count, dislike_count,
				title, content, pcreated_at, category_ids, pf.rating
			FROM v_posts
			LEFT JOIN post_feedback pf ON pf.parent_id = v_posts.id AND pf.user_id = ?` +
		getWhereQuery(filterBy, id) +
		getOrderByQuery(orderBy)
	rows, err := db.conn.Query(qry, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var catIDs string
		var rating sql.NullInt64
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.UserName,
			&p.CommentCount,
			&p.LikeCount,
			&p.DislikeCount,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&catIDs,
			&rating)
		if err != nil {
			return nil, err
		}
		p.Categories, p.CatNames, err = db.splitCategoryIDs(catIDs)
		if err != nil {
			return nil, err
		}
		if len(p.Content) > 50 {
			p.Content = p.Content[:50] + "..."
		}
		p.Rating = int(rating.Int64)
		p.TimeAgo = getTimeAgo(p.CreatedAt)
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return posts, nil
}

func (db *DBContainer) SelectPost(id, userID int) (*Post, error) {
	qry := `SELECT v_posts.id, puser_id, user_name, 
			comment_count, like_count, dislike_count,
			title, content, pcreated_at, category_ids, pf.rating
			FROM v_posts
			LEFT JOIN post_feedback pf ON pf.parent_id = v_posts.id AND pf.user_id = ?
			WHERE id = ?`
	var p Post
	var catIDs string
	var rating sql.NullInt64
	err := db.conn.QueryRow(qry, userID, id).Scan(
		&p.ID,
		&p.UserID,
		&p.UserName,
		&p.CommentCount,
		&p.LikeCount,
		&p.DislikeCount,
		&p.Title,
		&p.Content,
		&p.CreatedAt,
		&catIDs,
		&rating)

	if err != nil {
		return nil, checkErrNoRows(err)
	}
	p.Categories, p.CatNames, err = db.splitCategoryIDs(catIDs)
	if err != nil {
		return nil, err
	}
	p.Rating = int(rating.Int64)
	p.TimeAgo = getTimeAgo(p.CreatedAt)
	return &p, err
}

// db.InsetPost() into db and record the categories too
// include created at for testing for now
func (db *DBContainer) InsertPost(p Post) (int, error) {
	if err := db.isValidCategories(p.Categories); err != nil {
		return -1, err
	}
	qry := `INSERT INTO posts 
			(user_id, title, content)
			VALUES (?, ?, ?)`

	res, err := db.conn.Exec(qry,
		p.UserID,
		p.Title,
		p.Content,
		p.CreatedAt)
	if err != nil {
		return -1, err
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	for _, catID := range p.Categories {
		_, err = db.conn.Exec(`INSERT INTO post_categories (post_id, category_id)
							   VALUES (?, ?)`, postID, catID)
	}
	return int(postID), err
}

// db.UpdatePost() for updating post when user make changes
func (db *DBContainer) UpdatePost(p Post) error {
	qry := `UPDATE posts SET title = ?, content = ?	WHERE id = ?`
	_, err := db.conn.Exec(qry,
		p.Title,
		p.Content, p.ID)
	if err != nil {
		return err
	}
	qry = `DELETE FROM post_categories WHERE post_id = ?`
	_, err = db.conn.Exec(qry, p.ID)
	if err != nil {
		return err
	}
	for _, catID := range p.Categories {
		_, err = db.conn.Exec(`INSERT INTO post_categories (post_id, category_id)
							   VALUES (?, ?)`, p.ID, catID)
	}
	return err
}
