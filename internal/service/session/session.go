package session

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

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

func (s *SessionService) GetSession(cookieValue string) (*models.Session, error) {
	session, err := s.repo.GetSessionBySessionID(cookieValue)
	if err != nil {
		return nil, fmt.Errorf("get session: %v", err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, models.ErrSessionExpired
	}

	return session, nil
}

func (s *SessionService) SetSession(userID int) (*models.Session, error) {
	oldSession, err := s.repo.GetSessionByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get old session: %v", err)
	}
	if oldSession != nil {
		err := s.repo.DeleteSessionByID(oldSession.UUID)
		if err != nil {
			return nil, fmt.Errorf("delete old session: %v", err)
		}
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("new uuid: %v", err)
	}

	session := &models.Session{
		UUID:      uuid.String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err = s.repo.AddSession(session)
	if err != nil {
		return nil, fmt.Errorf("add session: %v", err)
	}

	return session, nil
}
