package notification

import "01.alem.school/git/nyeltay/forum/internal/models"

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
