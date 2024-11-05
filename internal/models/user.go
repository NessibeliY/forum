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
	SignupUser(user *Signup) error
	LoginUser(user *Login) (int, error)
	UpdateUser(user *UpdateUser) error
}

type UserRepository interface {
	AddUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *UpdateUser) error
}
