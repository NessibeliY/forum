package category

import (
	"context"
	"database/sql"
	"fmt"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	query := `SELECT * FROM category`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		category := models.Category{}

		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		categories = append(categories, category)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name FROM category WHERE name = $1`, name)
	category := models.Category{}
	err := row.Scan(
		&category.ID,
		&category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("row scan: %w", err)
	}
	return &category, nil
}
