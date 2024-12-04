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

func (s *NotificationService) GetListNotifications(user_id int) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	notification_list, err := s.repo.GetListNotifications(ctx, user_id)
	if err != nil {
		return []models.Notification{}, err
	}

	return notification_list, nil
}
