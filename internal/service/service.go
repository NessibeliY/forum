package service

import (
	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/internal/repository"
	"01.alem.school/git/nyeltay/forum/internal/service/category"
	"01.alem.school/git/nyeltay/forum/internal/service/comment"
	"01.alem.school/git/nyeltay/forum/internal/service/comment_reaction"
	"01.alem.school/git/nyeltay/forum/internal/service/post"
	"01.alem.school/git/nyeltay/forum/internal/service/post_reaction"
	"01.alem.school/git/nyeltay/forum/internal/service/session"
	"01.alem.school/git/nyeltay/forum/internal/service/user"
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
		CategoryService:        category.NewCategoryService(repo.CategoryRepo),
		CommentService:         comment.NewCommentService(repo.CommentRepo),
		CommentReactionService: comment_reaction.NewCommentReactionService(repo.CommentReactionRepo),
		PostService:            post.NewPostService(repo.PostRepo),
		PostReactionService:    post_reaction.NewPostReactionService(repo.PostReactionRepo),
		SessionService:         session.NewSessionService(repo.SessionRepo),
		UserService:            user.NewUserService(repo.UserRepo, repo.RoleRepo),
	}
}
