package models

import "time"

type User struct {
	ID             int
	Username       string
	HashedPassword string
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService interface {
	SignupUser(user *SignupRequest) error
	LoginUser(user *LoginRequest) (int, error)
	GetUserByID(id int) (*User, error)
	// UpdateUser(user *UpdateUserRequest) error
}

type UserRepository interface {
	AddUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	// UpdateUser(user *UpdateUser) error
}
