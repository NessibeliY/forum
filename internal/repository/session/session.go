package session

import (
	"database/sql"
	"fmt"

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
	_, err := r.db.Exec("DELETE FROM session WHERE uuid = $1", sessionID)
	return err
}

func (r *SessionRepository) GetSessionBySessionID(sessionID string) (session *models.Session, err error) {
	session = &models.Session{}
	query := `SELECT * FROM session WHERE uuid = $1`
	err = r.db.QueryRow(query, sessionID).Scan(&session.UUID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query row: %v", err)
	}

	return session, nil
}

func (r *SessionRepository) GetSessionByUserID(userID int) (session *models.Session, err error) {
	session = &models.Session{}
	query := `SELECT * FROM session WHERE user_id = ?`
	err = r.db.QueryRow(query, userID).Scan(&session.UUID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query row: %v", err)
	}

	return session, nil
}

func (r *SessionRepository) AddSession(session *models.Session) error {
	query := `INSERT INTO session (uuid, user_id, expire_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, session.UUID, session.UserID, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("insert row: %v", err)
	}
	return nil
}
