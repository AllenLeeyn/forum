INSERT INTO users (id, type_id, name, email, pw_hash)
VALUES (0, 1, 'ADMIN', 'admin@somewhere.earth', 's0m3thingS3cured?');

INSERT INTO posts (user_id, title, content)
VALUES (
    (SELECT id FROM users WHERE email = 'admin@somewhere.earth'),
    'Welcome. Please read and respect each other.',
    'Welcome!

We are so glad to have you here! Whether you are here to share your thoughts, ask questions, or simply learn, this is the perfect place to connect with like-minded individuals.

To make sure everyone has a positive experience, here are a few simple guidelines:

1. Be Respectful: Treat everyone with kindness and respect, even if you disagree. Healthy debate is welcome, but personal attacks are not.
2. Stay on Topic: Keep your posts relevant to the discussion. If you are starting a new topic, make sure it fits in the right category.
3. No Spam: We want to keep the forum focused and valuable for all users. Please avoid unsolicited promotions, ads, or irrelevant links.
4. Help Each Other: If you know the answer to a question, feel free to jump in and share! We are all here to learn from one another.

Feel free to explore the different sections, introduce yourself, and make yourself at home! If you need any help or have any questions, do not hesitate to ask our community or team.

Happy posting!'
);

INSERT INTO post_categories (post_id, category_id)
VALUES(
    (SELECT id FROM posts WHERE user_id = (SELECT id FROM users WHERE email = 'admin@somewhere.earth')),
    (SELECT id FROM categories WHERE name = 'General')
);
