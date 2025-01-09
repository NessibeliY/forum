package models

import (
	"time"
)

const (
	ModeratorRole = "moderator"
	AdminRole     = "admin"
	UserRole      = "user"
)

type User struct {
	ID             int
	Username       string
	HashedPassword string
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Role           string
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

type UpdateRoleRequest struct {
	ID        int
	UserID    int
	AdminID   int
	Processed bool
}

type ModeratorRequest struct {
	UserID   int
	Username string
}

type UserService interface {
	SignupUser(user *SignupRequest) error
	LoginUser(user *LoginRequest) (int, error)
	GetUserByID(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	SendModeratorRequest(userID int) error
	CheckModeratorRequestStatus(userID int) (bool, error)
	GetModeratorRequests() ([]ModeratorRequest, error)
	SetNewRole(request *UpdateRoleRequest) error
}

type UserRepository interface {
	AddUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	UpdateRole(userID int, role string) error
}
type RoleRepository interface {
	AddRoleRequest(request *UpdateRoleRequest) error
	ExistsByUserAndRole(userID int) (bool, error)
	GetModeratorRequests() ([]ModeratorRequest, error)
	UpdateRoleRequest(request *UpdateRoleRequest) error
	DeleteRoleRequestByUsedID(usedID int) error
}
