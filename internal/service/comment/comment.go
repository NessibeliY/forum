package comment

import (
	"01.alem.school/git/nyeltay/forum/internal/models"
)

type CommentService struct {
	repo models.CommentRepository
}

func NewCommentService(repo models.CommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (s *CommentService) CreateComment(createCommentRequest *models.CreateCommentRequest) error {
	return nil
}
