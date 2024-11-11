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
