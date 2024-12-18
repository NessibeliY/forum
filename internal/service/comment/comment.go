package comment

import (
	"context"
	"time"

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

func (s *CommentService) GetAllCommentsByPostID(postID int) ([]*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetAllCommentsByPostID(ctx, postID)
}

func (s *CommentService) CreateComment(createCommentRequest *models.CreateCommentRequest) error {
	comment := &models.Comment{
		Content:  createCommentRequest.Content,
		AuthorID: createCommentRequest.AuthorID,
		PostID:   createCommentRequest.PostID,
	}

	return s.repo.AddComment(comment)
}

func (s *CommentService) DeleteComment(request *models.DeleteCommentRequest) error {
	return s.repo.DeleteComment(request.ID)
}

func (s *CommentService) GetCommentByID(id int) (*models.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetCommentByID(ctx, id)
}

func (s *CommentService) GetUserCommentedPosts(authorID int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetUserCommentedPosts(ctx, authorID)
}

func (s *CommentService) UpdateComment(updateCommentRequest *models.UpdateCommentRequest) error {
	return s.repo.UpdateComment(updateCommentRequest)
}
