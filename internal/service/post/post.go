package post

import (
	"context"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type PostService struct {
	repo models.PostRepository
}

func NewPostService(repo models.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetAllPosts(ctx)
}
