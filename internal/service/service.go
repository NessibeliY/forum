package service

import (
	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/internal/repository"
)

type Service struct {
	CategoryService        models.CategoryService
	CommentService         models.CommentService
	CommentReactionService models.CommentReactionService
	PostService            models.PostService
	PostReactionService    models.PostReactionService
	SessionService         models.SessionService
	UserService            models.UserService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		CategoryService:        category.NewService(repo),
		CommentService:         comment.NewService(repo),
		CommentReactionService: comment_reaction.NewService(repo),
		PostService:            post.NewService(repo),
		PostReactionService:    post_reaction.NewService(repo),
		SessionService:         session.NewService(repo),
		UserService:            user.NewService(repo),
	}
}
