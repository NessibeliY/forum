package post_reaction

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *PostReactionRepository) GetReactionsByPostID(ctx context.Context, postID int) (reactions []*models.PostReaction, err error) {
	query := `SELECT * FROM post_reaction WHERE post_id=?`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reaction models.PostReaction
		err = rows.Scan(
			&reaction.AuthorID,
			&reaction.PostID,
			&reaction.Reaction)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		reactions = append(reactions, &reaction)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return reactions, nil
}
