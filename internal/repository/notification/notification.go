package notification

import (
	"context"
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

func (r *NotificationRepository) GetCountNotifications(userID int) (int, error) {
	var count int
	query := `
        SELECT COUNT(n.id)
        FROM notifications n
        JOIN post p ON n.post_id = p.id
        JOIN users u ON p.author_id = u.id
        WHERE u.id = $1 AND n.is_read = false
		ORDER BY n.created_at DESC;
    `

	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error get count notifications: %w", err)
	}
	return count, nil
}

func (r *NotificationRepository) GetCurrentNotifications(ctx context.Context, userID int) ([]models.Notification, error) {
	query := `
	SELECT 
		n.id,n.post_id,n.message,n.is_read,n.created_at 
	FROM notifications n
	JOIN post p ON p.id = n.post_id
	JOIN users u ON u.id = p.author_id
	WHERE u.id = $1 AND n.is_read = 0
	ORDER BY n.created_at DESC;
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var id, postID int
		var message string
		var isRead bool
		var createdAt time.Time

		err := rows.Scan(
			&id,
			&postID,
			&message,
			&isRead,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		notification := models.Notification{
			ID:        id,
			PostID:    postID,
			Message:   message,
			IsRead:    isRead,
			CreatedAt: createdAt,
		}
		notifications = append(notifications, notification)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return notifications, nil
}

func (r *NotificationRepository) GetArchivedNotifications(ctx context.Context, userID int) ([]models.Notification, error) {
	query := `
	SELECT 
		n.id,n.post_id,n.message,n.is_read,n.created_at 
	FROM notifications n
	JOIN post p ON p.id = n.post_id
	JOIN users u ON u.id = p.author_id
	WHERE u.id = $1 AND  n.is_read = 1
	ORDER BY n.created_at DESC;
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var id, postID int
		var message string
		var isRead bool
		var createdAt time.Time

		err := rows.Scan(
			&id,
			&postID,
			&message,
			&isRead,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		notification := models.Notification{
			ID:        id,
			PostID:    postID,
			Message:   message,
			IsRead:    isRead,
			CreatedAt: createdAt,
		}
		notifications = append(notifications, notification)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return notifications, nil
}

func (r *NotificationRepository) MakeNotificationIsRead(userID, notificationID int) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE
		 WHERE id IN (
            SELECT n.id
            FROM notifications n
            JOIN post p ON p.id = n.post_id
            WHERE p.author_id = $1 AND n.id = $2
        );
	`

	if _, err := r.db.Exec(query, userID, notificationID); err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}

	return nil
}

func (r *NotificationRepository) DeleteNotificationsByPostID(postID int) error {
	query := `DELETE FROM notifications WHERE post_id = $1`

	result, err := r.db.Exec(query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete notification for post %d: %w", postID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for post %d: %w", postID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no notifications found for post %d to delete", postID)
	}

	return nil
}

func (r *NotificationRepository) GetNotificationByID(id int) (*models.Notification, error) {
	query := `
		SELECT id, post_id, message, is_read, created_at
		FROM notifications
		WHERE id = $1
	`

	var notification models.Notification
	err := r.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.PostID,
		&notification.Message,
		&notification.IsRead,
		&notification.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notification by ID: %w", err)
	}

	return &notification, nil
}

func (r *NotificationRepository) GetNotificationsForPost(ctx context.Context, postID int) ([]models.Notification, error) {
	query := `
	SELECT 
		id,post_id,message,is_read,created_at 
	FROM notifications 
	WHERE  post_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var id, postID int
		var message string
		var isRead bool
		var createdAt time.Time

		err := rows.Scan(
			&id,
			&postID,
			&message,
			&isRead,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		notification := models.Notification{
			ID:        id,
			PostID:    postID,
			Message:   message,
			IsRead:    isRead,
			CreatedAt: createdAt,
		}
		notifications = append(notifications, notification)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return notifications, nil
}
