package repository

import (
	"database/sql"

	"01.alem.school/git/nyeltay/forum/internal/models"
	"01.alem.school/git/nyeltay/forum/internal/repository/category"
	"01.alem.school/git/nyeltay/forum/internal/repository/comment"
	"01.alem.school/git/nyeltay/forum/internal/repository/comment_reaction"
	"01.alem.school/git/nyeltay/forum/internal/repository/moderation"
	"01.alem.school/git/nyeltay/forum/internal/repository/notification"
	"01.alem.school/git/nyeltay/forum/internal/repository/post"
	"01.alem.school/git/nyeltay/forum/internal/repository/post_reaction"
	"01.alem.school/git/nyeltay/forum/internal/repository/role"
	"01.alem.school/git/nyeltay/forum/internal/repository/session"
	"01.alem.school/git/nyeltay/forum/internal/repository/user"
)

type Repository struct {
	CategoryRepo        models.CategoryRepository
	CommentRepo         models.CommentRepository
	CommentReactionRepo models.CommentReactionRepository
	PostRepo            models.PostRepository
	PostReactionRepo    models.PostReactionRepository
	SessionRepo         models.SessionRepository
	UserRepo            models.UserRepository
	NotificationRepo    models.NotificationRepository
	ModerationRepo      models.ModerationRepository
	RoleRepo            models.RoleRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		CategoryRepo:        category.NewCategoryRepository(db),
		CommentRepo:         comment.NewCommentRepository(db),
		CommentReactionRepo: comment_reaction.NewCommentReactionRepository(db),
		PostRepo:            post.NewPostRepository(db),
		PostReactionRepo:    post_reaction.NewPostReactionRepository(db),
		SessionRepo:         session.NewSessionRepository(db),
		UserRepo:            user.NewUserRepository(db),
		NotificationRepo:    notification.NewNotificationRepository(db),
		ModerationRepo:      moderation.NewModerationRepository(db),
		RoleRepo:            role.NewRoleRepository(db),
	}
}
