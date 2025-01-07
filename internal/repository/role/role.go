package role

import (
	"database/sql"
	"fmt"

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

func (r *RoleRepository) UpdateRoleRequest(request *models.UpdateRoleRequest) error {
	query := `UPDATE new_role_request SET admin_id = ?, processed = ? WHERE user_id = ?`
	_, err := r.db.Exec(query, request.AdminID, request.Processed, request.UserID)
	return err
}

func (r *RoleRepository) DeleteRoleRequestByUsedID(usedID int) error {
	query := `DELETE FROM new_role_request WHERE user_id = ?`
	_, err := r.db.Exec(query, usedID)
	return err
}

func (r *RoleRepository) GetModeratorRequests() ([]models.ModeratorRequest, error) {
	query := `
SELECT u.username, n.user_id
FROM users u
JOIN new_role_request n ON n.user_id = u.id
WHERE n.processed = ?`
	rows, err := r.db.Query(query, "false")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var moderatorRequests []models.ModeratorRequest
	for rows.Next() {
		var moderatorRequest models.ModeratorRequest

		err := rows.Scan(&moderatorRequest.Username, &moderatorRequest.UserID)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		moderatorRequests = append(moderatorRequests, moderatorRequest)
	}

	return moderatorRequests, nil
}
