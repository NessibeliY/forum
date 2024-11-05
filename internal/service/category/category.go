package category

import (
	"01.alem.school/git/nyeltay/forum/internal/models"
)

type CategoryService struct {
	repo models.CategoryRepository
}

func NewCategoryService(repo models.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}
