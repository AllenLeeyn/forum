CREATE VIEW v_posts AS
SELECT 
    p.id,
    p.user_id AS puser_id,
    u.name AS user_name,
    p.comment_count,
    p.like_count,
    p.dislike_count,
    p.title,
    p.content,
    p.created_at AS pcreated_at,
    GROUP_CONCAT(pc.category_id) AS category_ids
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN post_categories pc ON p.id = pc.post_id
GROUP BY 
    p.id, p.user_id, u.name,
    p.comment_count, p.like_count, p.dislike_count,
    p.title, p.content, p.created_at;