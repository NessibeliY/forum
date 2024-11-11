package models

import "time"

type Session struct {
	UUID      string
	UserID    int
	ExpiresAt time.Time
}

type SessionService interface {
	SetSession(userID int) (*Session, error)
	GetSession(userID int) (*Session, error)
	UpdateSession(session *Session) error
	DeleteSession(cookieValue string) error
}

type SessionRepository interface {
	AddSession(session *Session) error
	GetSessionByUserID(userID string) (*Session, error)
	UpdateSession(session *Session) error
	DeleteSessionByID(sessionID string) error
}
