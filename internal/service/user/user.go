package user

import (
	"01.alem.school/git/nyeltay/forum/internal/models"
)

type UserService struct {
	repo models.UserRepository
}

func NewUserService(repo models.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
