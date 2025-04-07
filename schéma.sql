/*Table s'occupant des informations users*/
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
/*Table s'occupant des posts des utilisateurs*/
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
/*Table s'occupant des commentaires utilisateurs*/
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_id INTEGER,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES posts(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);

INSERT INTO users (username, email) VALUES ('Alice', 'alice@example.com');
INSERT INTO users (username, email) VALUES ('Bob', 'bob@example.com');

INSERT INTO posts (title, content, user_id) VALUES ('Bienvenue sur le forum', 'Test premier message.', 1);
INSERT INTO posts (title, content, user_id) VALUES ('Deuxième post', 'Salut tout le monde!', 2);

INSERT INTO comments (post_id, user_id, content) VALUES (1, 2, 'Super, merci Alice !');
INSERT INTO comments (post_id, user_id, content) VALUES (2, 1, 'Bienvenue Bob !');

INSERT INTO likes (post_id, user_id, type) VALUES (1, 2, 'like');
INSERT INTO likes (post_id, user_id, type) VALUES (2, 1, 'like');

SELECT posts.title, users.username FROM posts
JOIN users ON posts.user_id = users.id;
