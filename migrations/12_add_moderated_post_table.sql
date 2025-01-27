CREATE TABLE IF NOT EXISTS moderated_post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    moderator_id INTEGER NOT NULL,
    admin_answer TEXT,
    moderated BOOLEAN,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (moderator_id) REFERENCES users(id),
    UNIQUE (post_id, moderator_id)
);