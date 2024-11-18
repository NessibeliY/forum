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

func (r *CommentReactionRepository) GetReactionByCommentIDAndAuthorID(ctx context.Context, commentID int, authorID int) (reaction *models.CommentReaction, err error) {
	reaction = &models.CommentReaction{}
	query := `SELECT * FROM comment_reaction WHERE comment_id=$1 AND author_id=$2`
	err = r.db.QueryRowContext(ctx, query, commentID, authorID).Scan(
		&reaction.AuthorID,
		&reaction.CommentID,
		&reaction.Reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return reaction, nil
}

func (r *CommentReactionRepository) DeleteCommentReaction(commentReaction *models.CommentReaction) error {
	query := `DELETE FROM comment_reaction WHERE comment_id=$1 AND author_id=$2`
	result, err := r.db.Exec(query, commentReaction.CommentID, commentReaction.AuthorID)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil
	}
	return nil
}

func (r CommentReactionRepository) UpdateCommentReaction(commentReaction *models.CommentReaction) error {
	query := `UPDATE comment_reaction SET reaction=$1 WHERE comment_id=$2 AND author_id=$3`
	result, err := r.db.Exec(query, commentReaction.Reaction, commentReaction.CommentID, commentReaction.AuthorID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil
	}
	return nil
}

func (r *CommentReactionRepository) AddCommentReaction(commentReaction *models.CommentReaction) error {
	query := `INSERT INTO comment_reaction(comment_id, author_id, reaction) VALUES($1, $2, $3)`
	_, err := r.db.Exec(query, commentReaction.CommentID, commentReaction.AuthorID, commentReaction.Reaction)
	return err
}
