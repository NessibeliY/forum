package comment

import (
	"database/sql"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) AddComment(comment *models.Comment) error {
	return nil
}
