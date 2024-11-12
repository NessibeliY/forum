package post_reaction

import (
	"database/sql"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type PostReactionRepository struct {
	db *sql.DB
}

func NewPostReactionRepository(db *sql.DB) *PostReactionRepository {
	return &PostReactionRepository{
		db: db,
	}
}

func (r *PostReactionRepository) AddPostReaction(postReaction *models.PostReaction) error {
	return nil
}

func (r *PostReactionRepository) UpdatePostReaction(postReaction *models.PostReaction) error {
	return nil
}

func (r *PostReactionRepository) DeletePostReaction(postID, authorID int) error {
	return nil
}
