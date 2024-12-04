package models

import (
	"context"
	"time"
)

type Notification struct {
	ID        int
	PostID    int
	Message   string
	IsRead    bool
	CreatedAt time.Time
}

type NotificationRequest struct {
	PostID  int
	Message string
}

type NotificationService interface {
	CreateNotification(notification *NotificationRequest) (int, error)
	GetCountNotifications(user_id int) (int, error)
	GetListNotifications(user_id int) ([]Notification, error)
}

type NotificationRepository interface {
	AddNotification(notification *Notification) (int, error)
	GetCountNotifications(user_id int) (int, error)
	GetListNotifications(ctx context.Context, user_id int) ([]Notification, error)
}
