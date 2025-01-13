package models

import "context"

type Category struct {
	ID   int
	Name string
}

type CategoryService interface {
	GetAllCategories() ([]Category, error)
	GetCategoryByName(name string) (*Category, error)
	GetCategoryByID(categoryID int) (*Category, error)
	CreateCategory(category *Category) (int, error)
	DeleteCategory(categoryID int) error
}

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]Category, error)
	GetCategoryByName(ctx context.Context, name string) (*Category, error)
	GetCategoryByID(ctx context.Context, categoryID int) (*Category, error)
	AddCategory(category *Category) (int, error)
	DeleteCategory(categoryID int) error
}
