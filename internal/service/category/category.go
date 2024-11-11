package category

import (
	"context"
	"time"

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

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetAllCategories(ctx)
}
