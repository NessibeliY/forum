package comment_reaction

import (
	"01.alem.school/git/nyeltay/forum/internal/models"
)

type CommentReactionService struct {
	repo models.CommentReactionRepository
}

func NewCommentReactionService(repo models.CommentReactionRepository) *CommentReactionService {
	return &CommentReactionService{
		repo: repo,
	}
}
