CREATE TRIGGER update_comment_count_after_insert
AFTER INSERT ON comments
FOR EACH ROW
BEGIN
    UPDATE posts
    SET comment_count = (
        SELECT COUNT(*) 
        FROM comments 
        WHERE post_id = NEW.post_id)
    WHERE id = NEW.post_id;
END;

CREATE TRIGGER update_like_count_after_insert_in_comment_feedback
AFTER INSERT ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET like_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = 1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_like_count_after_update_in_comment_feeback
AFTER UPDATE ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET like_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = 1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_like_count_after_insert_in_post_feedback
AFTER INSERT ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET like_count = (
        SELECT COUNT(*) 
        FROM post_feedback
        WHERE parent_id = NEW.parent_id AND rating = 1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_like_count_after_update_in_post_feedback
AFTER UPDATE ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET like_count = (
        SELECT COUNT(*) 
        FROM post_feedback
        WHERE parent_id = NEW.parent_id AND rating = 1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_dislike_count_after_insert_in_comment_feedback
AFTER INSERT ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = -1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_dislike_count_after_update_in_comment_feeback
AFTER UPDATE ON comment_feedback
FOR EACH ROW
BEGIN
    UPDATE comments
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM comment_feedback 
        WHERE parent_id = NEW.parent_id AND rating = -1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_dislike_count_after_insert_in_post_feedback
AFTER INSERT ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM post_feedback
        WHERE parent_id = NEW.parent_id AND rating = -1)
    WHERE id = NEW.parent_id;
END;

CREATE TRIGGER update_dislike_count_after_update_in_post_feedback
AFTER UPDATE ON post_feedback
FOR EACH ROW
BEGIN
    UPDATE posts
    SET dislike_count = (
        SELECT COUNT(*) 
        FROM post_feedback
        WHERE parent_id = NEW.parent_id AND rating = -1)
    WHERE id = NEW.parent_id;
END;
