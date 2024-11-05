package models

type Category struct {
	Name string
}

type CategoryService interface {
	CreateCategory(category *Category) (string, error)
}

type CategoryRepository interface {
	AddCategory(category *Category) (string, error)
}
