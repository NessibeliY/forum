CREATE TABLE IF NOT EXISTS post_reaction (
    author_id INTEGER NOT NULL,
    post_id INTEGER,
    reaction TEXT NOT NULL,
    PRIMARY KEY (author_id, post_id),
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
);