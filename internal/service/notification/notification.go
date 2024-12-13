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

func (s *NotificationService) GetCountNotifications(user_id int) (int, error) {
	return s.repo.GetCountNotifications(user_id)
}

func (s *NotificationService) GetCurrentNotifications(user_id int) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	current_notification_list, err := s.repo.GetCurrentNotifications(ctx, user_id)
	if err != nil {
		return []models.Notification{}, err
	}

	return current_notification_list, nil
}

func (s *NotificationService) GetArchivedNotifications(user_id int) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	notification_list, err := s.repo.GetArchivedNotifications(ctx, user_id)
	if err != nil {
		return []models.Notification{}, err
	}

	return notification_list, nil
}

func (s *NotificationService) MakeNotificationIsRead(user_id, notification_id int) error {
	return s.repo.MakeNotificationIsRead(user_id, notification_id)
}

func (s *NotificationService) RemoveNotificationFromPost(post_id int) error {
	return s.repo.RemoveNotificationFromPost(post_id)
}

func (s *NotificationService) GetNotificationByID(id int) (*models.Notification, error) {
	return s.repo.GetNotificationByID(id)
}
