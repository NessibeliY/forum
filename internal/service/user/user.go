package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type UserService struct {
	userRepo models.UserRepository
	roleRepo models.RoleRepository
}

func NewUserService(userRepo models.UserRepository, roleRepo models.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
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

	err = s.userRepo.AddUser(user)
	if err != nil {
		return fmt.Errorf("add user: %w", err)
	}

	return nil
}

func (s *UserService) LoginUser(loginRequest *models.LoginRequest) (int, error) {
	user, err := s.userRepo.GetUserByEmail(loginRequest.Email)
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
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *UserService) SendModeratorRequest(userID int) error {
	request := &models.UpdateRoleRequest{
		UserID:    userID,
		Processed: false,
	}
	return s.roleRepo.AddRoleRequest(request)
}

func (s *UserService) GetModeratorRequests() ([]models.ModeratorRequest, error) {
	return s.roleRepo.GetModeratorRequests()
}

func (s *UserService) SetNewRole(request *models.UpdateRoleRequest) error {
	var err error
	switch request.Processed {
	case true:
		err = s.roleRepo.UpdateRoleRequest(request)
	case false:
		err = s.roleRepo.DeleteRoleRequestByUsedID(request.UserID)
	}
	if err != nil {
		return fmt.Errorf("update role request: %w", err)
	}

	if request.Processed {
		err = s.userRepo.UpdateRole(request.UserID, models.ModeratorRole)
		if err != nil {
			return fmt.Errorf("update user role: %w", err)
		}
	}

	return nil
}
