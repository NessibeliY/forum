package notification

import (
	"database/sql"
	"fmt"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) AddNotification(notification *models.Notification) (int, error) {
	createdAt := time.Now()
	query := `INSERT INTO notifications (post_id,message,is_read,created_at) VALUES($1,$2,$3,$4) RETURNING id;`
	err := r.db.QueryRow(query, notification.PostID, notification.Message, notification.IsRead, createdAt).Scan(&notification.ID)
	if err != nil {
		return 0, fmt.Errorf("insert notification: %w", err)
	}
	return notification.ID, nil
}

func (r *NotificationRepository) GetCountNotifications(user_id int) (int, error) {
	var count int
	query := `
        SELECT COUNT(n.id)
        FROM notifications n
        JOIN post p ON n.post_id = p.id
        JOIN users u ON p.author_id = u.id
        WHERE u.id = ?
    `

	err := r.db.QueryRow(query, user_id).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("ошибка при получении количества уведомлений: %w", err)
	}
	return count, nil
}
