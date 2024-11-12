package post_reaction

import (
	"01.alem.school/git/nyeltay/forum/internal/models"
)

type PostReactionService struct {
	repo models.PostReactionRepository
}

func NewPostReactionService(repo models.PostReactionRepository) *PostReactionService {
	return &PostReactionService{
		repo: repo,
	}
}

func (s *PostReactionService) CreatePostReaction(request *models.PostReactionRequest) error {
	return nil
}

func (s *PostReactionService) UpdatePostReaction(request *models.PostReactionRequest) error {
	return nil
}

func (s *PostReactionService) DeletePostReaction(request *models.PostReactionRequest) error {
	return nil
}
