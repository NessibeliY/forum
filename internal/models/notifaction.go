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
	GetCountNotifications(userID int) (int, error)
	GetCurrentNotifications(userID int) ([]Notification, error)
	MakeNotificationIsRead(userID, notificationID int) error
	GetArchivedNotifications(userID int) ([]Notification, error)
	RemoveNotificationFromPost(postID int) error
	GetNotificationByID(id int) (*Notification, error)
	GetNotificationsForPost(postID int) ([]Notification, error)
}

type NotificationRepository interface {
	AddNotification(notification *Notification) (int, error)
	GetCountNotifications(user_id int) (int, error)
	GetCurrentNotifications(ctx context.Context, user_id int) ([]Notification, error)
	MakeNotificationIsRead(userID, notificationID int) error
	GetArchivedNotifications(ctx context.Context, userID int) ([]Notification, error)
	DeleteNotificationsByPostID(postID int) error
	GetNotificationByID(id int) (*Notification, error)
	GetNotificationsForPost(ctx context.Context, postID int) ([]Notification, error)
}
