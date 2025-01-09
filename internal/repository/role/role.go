package role

import (
	"database/sql"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) AddRoleRequest(request *models.UpdateRoleRequest) error {
	query := `INSERT INTO new_role_request (user_id, processed) VALUES (?, ?)`
	_, err := r.db.Exec(query, request.UserID, request.Processed)
	return err
}

func (r *RoleRepository) ExistsByUserAndRole(userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM new_role_request WHERE user_id = $1 AND processed = 0`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
