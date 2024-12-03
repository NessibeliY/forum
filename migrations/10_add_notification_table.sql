CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,  -- пользователь, которому нужно отправить уведомление
    message TEXT NOT NULL,     -- текст уведомления
    is_read BOOLEAN DEFAULT 0, -- если уведомление прочитано, будет стоять 1
    created_at DATETIME NOT NULL, -- время создания уведомления
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
);
