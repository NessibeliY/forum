package models

import "context"

type Category struct {
	ID   int
	Name string
}

type CategoryService interface {
	GetAllCategories() ([]Category, error)
	CreateCategory(category *Category) (string, error)
}

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]Category, error)
	AddCategory(category *Category) (string, error)
}
