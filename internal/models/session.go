package models

import "time"

type Session struct {
	UUID      string
	UserID    int
	ExpiresAt time.Time
}

type SessionService interface {
	SetSession(userID int) error
	GetSession(userID int) (*Session, error)
	UpdateSession(session *Session) error
	DeleteSession(uuid int) error
}

type SessionRepository interface {
	AddSession(session *Session) error
	GetSessionByUserID(userID string) (*Session, error)
	UpdateSession(session *Session) error
	DeleteSessionByUUID(uuid int) error
}
