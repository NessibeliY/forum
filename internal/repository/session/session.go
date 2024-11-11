package session

import "database/sql"

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
