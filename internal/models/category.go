package models

import "context"

type Category struct {
	ID   int
	Name string
}

type CategoryService interface {
	GetAllCategories() ([]Category, error)
	GetCategoryByName(name string) (*Category, error)
	//CreateCategory(category *Category) (string, error)
}

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]Category, error)
	GetCategoryByName(ctx context.Context, name string) (*Category, error)
	//AddCategory(category *Category) (string, error)
}
