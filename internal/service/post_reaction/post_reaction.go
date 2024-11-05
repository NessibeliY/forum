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
