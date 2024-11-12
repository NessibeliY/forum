package comment_reaction

import (
	"database/sql"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type CommentReactionRepository struct {
	db *sql.DB
}

func NewCommentReactionRepository(db *sql.DB) *CommentReactionRepository {
	return &CommentReactionRepository{
		db: db,
	}
}

func (r *CommentReactionRepository) AddCommentReaction(commentReaction *models.CommentReaction) error {
	return nil
}

func (r *CommentReactionRepository) UpdateCommentReaction(commentReaction *models.CommentReaction) error {
	return nil
}

func (r *CommentReactionRepository) DeleteCommentReaction(commentID, authorID int) error {
	return nil
}
