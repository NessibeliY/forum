package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

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

func (s *UserService) SignupUser(signupRequest *models.SignupRequest) error {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate hashed password: %w", err)
	}

	user := &models.User{
		Username:       signupRequest.Username,
		HashedPassword: string(hashedPW),
		Email:          signupRequest.Email,
	}

	err = s.repo.AddUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) LoginUser(loginRequest *models.LoginRequest) (int, error) {
	user, err := s.repo.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return 0, fmt.Errorf("get user by email: %w", err)
	}
	if user == nil {
		return 0, models.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginRequest.Password))
	if err != nil {
		return 0, models.ErrInvalidCredentials
	}

	return user.ID, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetUserByID(id)
}
