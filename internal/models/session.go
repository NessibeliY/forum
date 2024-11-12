package models

import "time"

type Session struct {
	UUID      string
	UserID    int
	ExpiresAt time.Time
}

type SessionService interface {
	SetSession(userID int) (*Session, error)
	GetSession(uuid string) (*Session, error)
	DeleteSession(cookieValue string) error
}

type SessionRepository interface {
	AddSession(session *Session) error
	GetSessionBySessionID(sessionID string) (*Session, error)
	GetSessionByUserID(userID int) (session *Session, err error)
	DeleteSessionByID(sessionID string) error
}
