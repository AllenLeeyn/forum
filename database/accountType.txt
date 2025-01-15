INSERT INTO account_type (id, name, can_create_post, can_comment, can_feedback, can_moderate, can_ban_user) 
VALUES (1, 'user', 1, 1, 1, 0, 0),
       (2, 'moderator', 1, 1, 1, 1, 1),
       (3, 'administrator', 1, 1, 1, 1, 1);