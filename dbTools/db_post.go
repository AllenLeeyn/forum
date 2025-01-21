package dbTools

import (
	"fmt"
	"strconv"
	"strings"
)

// getWhereQuery() for selectPosts() filterBy query
func getWhereQuery(filterBy string, id int) string {
	switch filterBy {
	case "createdBy":
		return fmt.Sprintf(` WHERE puser_id = %v`, id)
	case "catergory":
		return fmt.Sprintf(` WHERE ',' || category_ids || ',' LIKE '%%,%v,%%'`, id)
	case "likedBy":
		return fmt.Sprintf(` INNER JOIN post_feedback pf ON pf.parent_id = v_posts.id 
							 WHERE pf.user_id = %v AND pf.rating = 1`, id)
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

// splitCategoryIDs into []int to store in post struct
func splitCategoryIDs(catIDs string) ([]int, error) {
	if catIDs == "" {
		return nil, fmt.Errorf("empty string")
	}
	var result []int
	categories := strings.Split(catIDs, ",")
	for _, idStr := range categories {
		if id, err := strconv.Atoi(idStr); err == nil {
			result = append(result, id)
		} else {
			return nil, err
		}
	}
	return result, nil
}

// db.selectPosts() with filter and order options.
// by default, no filter and newest first are applied.
// if invalid options or empty are given, default option is used.
// valid filterBy: createdBy, catergory, likedBy
// valid orderBy: oldest, likeCount, commentCount
func (db *DBContainer) SelectPosts(filterBy, orderBy string, id int) (*[]post, error) {
	qry := `SELECT id, puser_id, user_name, 
			comment_count, like_count, dislike_count,
			title, content, pcreated_at, category_ids
			FROM v_posts` +
		getWhereQuery(filterBy, id) +
		getOrderByQuery(orderBy)
	rows, err := db.conn.Query(qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []post
	for rows.Next() {
		var p post
		var catIDs string
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
			&catIDs)
		if err != nil {
			return nil, err
		}
		p.Categories, err = splitCategoryIDs(catIDs)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return &posts, nil
}

// db.insetPost() into db and record the categories too
func (db *DBContainer) InsertPost(p post) error {
	if err := db.isValidCategories(p.Categories); err != nil {
		return err
	}
	qry := `INSERT INTO posts 
			(user_id, comment_count, like_count, dislike_count, title, content, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := db.conn.Exec(qry,
		p.UserID,
		p.CommentCount,
		p.LikeCount,
		p.DislikeCount,
		p.Title,
		p.Content,
		p.CreatedAt)
	if err != nil {
		return err
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	for _, catID := range p.Categories {
		_, err = db.conn.Exec(`INSERT INTO post_categories (post_id, category_id)
								VALUES (?, ?)`, postID, catID)
	}
	return err
}
