package session

import (
	"database/sql"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) DeleteSessionByID(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE uuid = $1", sessionID)
	return err
}

func (r *SessionRepository) AddSession(session *models.Session) error {
	return nil
}

func (r *SessionRepository) GetSessionByUserID(userID string) (*models.Session, error) {
	return nil, nil
}

func (r *SessionRepository) UpdateSession(session *models.Session) error {
	return nil
}
