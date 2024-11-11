package user

import (
	"database/sql"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) AddUser(user *models.User) error {
	createdAt := time.Now()
	updatedAt := createdAt

	_, err := r.db.Exec("INSERT INTO user (username, hashed_pw, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Username,
		user.HashedPassword,
		user.Email,
		createdAt,
		updatedAt,
	)
	if err != nil {
		switch err.Error() {
		case "UNIQUE constraint failed: user.email":
			return models.ErrDuplicateEmail
		case "UNIQUE constraint failed: user.username":
			return models.ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT * FROM user WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
