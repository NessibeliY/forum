package notification

import (
	"context"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type NotificationService struct {
	repo models.NotificationRepository
}

func NewNotificationService(repo models.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) CreateNotification(notificationRequest *models.NotificationRequest) (int, error) {
	notification := models.Notification{
		PostID:  notificationRequest.PostID,
		Message: notificationRequest.Message,
	}
	return s.repo.AddNotification(&notification)
}

func (s *NotificationService) GetCountNotifications(userID int) (int, error) {
	return s.repo.GetCountNotifications(userID)
}

func (s *NotificationService) GetCurrentNotifications(userID int) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	current_notification_list, err := s.repo.GetCurrentNotifications(ctx, userID)
	if err != nil {
		return []models.Notification{}, err
	}

	return current_notification_list, nil
}

func (s *NotificationService) GetArchivedNotifications(userID int) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	notification_list, err := s.repo.GetArchivedNotifications(ctx, userID)
	if err != nil {
		return []models.Notification{}, err
	}

	return notification_list, nil
}

func (s *NotificationService) MakeNotificationIsRead(userID, notificationID int) error {
	return s.repo.MakeNotificationIsRead(userID, notificationID)
}

func (s *NotificationService) RemoveNotificationFromPost(postID int) error {
	return s.repo.DeleteNotificationsByPostID(postID)
}

func (s *NotificationService) GetNotificationByID(id int) (*models.Notification, error) {
	return s.repo.GetNotificationByID(id)
}

func (s *NotificationService) GetNotificationsForPost(postID int) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetNotificationsForPost(ctx, postID)
}
