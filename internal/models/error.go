package models

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRows             = sql.ErrNoRows
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateUsername  = errors.New("duplicate username")
	ErrPostNotFound       = errors.New("post not found")
)
