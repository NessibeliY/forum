package user

import (
	"context"
	"database/sql"
	"fmt"
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

	_, err := r.db.Exec("INSERT INTO users (username, hashed_pw, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Username,
		user.HashedPassword,
		user.Email,
		createdAt,
		updatedAt,
	)
	if err != nil {
		switch err.Error() {
		case "UNIQUE constraint failed: users.email":
			return models.ErrDuplicateEmail
		case "UNIQUE constraint failed: users.username":
			return models.ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.HashedPassword, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(id int) (user *models.User, err error) {
	user = &models.User{}
	query := `SELECT * FROM users WHERE id = $1`
	err = r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateRole(userID int, role string) error {
	query := `UPDATE users SET role = ? WHERE id = ?`
	_, err := r.db.Exec(query, role, userID)
	return err
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT * FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := models.User{}

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.HashedPassword,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return users, nil
}
