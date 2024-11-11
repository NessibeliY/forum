CREATE TABLE IF NOT EXISTS comment_reaction (
    author_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    reaction TEXT NOT NULL,
    PRIMARY KEY (author_id, comment_id),
    FOREIGN KEY (author_id) REFERENCES user(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comment(id) ON DELETE CASCADE
);