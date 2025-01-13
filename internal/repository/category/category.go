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

func (r *CategoryRepository) DeleteCategory(categoryID int) error {
	query := `DELETE FROM category WHERE id = $1`
	result, err := r.db.Exec(query, categoryID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil
	}

	return nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, categoryID int) (*models.Category, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name FROM category WHERE id = $1`, categoryID)
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

func (r *CategoryRepository) AddCategory(category *models.Category) (int, error) {
	query := `INSERT INTO category (name) VALUES ($1) RETURNING id;`

	result, err := r.db.Exec(query, category.Name)
	if err != nil {
		return 0, fmt.Errorf("insert category: %w", err)
	}

	// Получаем последний вставленный id
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	fmt.Println("lastInsertID", lastInsertID)

	// Возвращаем id новой категории
	category.ID = int(lastInsertID)
	return category.ID, nil
}
