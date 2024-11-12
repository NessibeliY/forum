package comment_reaction

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *CommentReactionRepository) GetReactionsByCommentID(ctx context.Context, commentID int) (reactions []*models.CommentReaction, err error) {
	query := `SELECT * FROM comment_reaction WHERE comment_id=?`
	rows, err := r.db.QueryContext(ctx, query, commentID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reaction models.CommentReaction
		err = rows.Scan(
			&reaction.AuthorID,
			&reaction.CommentID,
			&reaction.Reaction)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		reactions = append(reactions, &reaction)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return reactions, nil
}
