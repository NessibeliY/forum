package session

import (
	"01.alem.school/git/nyeltay/forum/internal/models"
)

type SessionService struct {
	repo models.SessionRepository
}

func NewSessionService(repo models.SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

func (s *SessionService) DeleteSession(cookieValue string) error {
	return s.repo.DeleteSessionByID(cookieValue)
}
