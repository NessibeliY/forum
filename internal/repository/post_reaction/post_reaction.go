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

func (r *PostReactionRepository) GetReactionByPostIDAndAuthorID(ctx context.Context, postID, authorID int) (reaction *models.PostReaction, err error) {
	reaction = &models.PostReaction{}
	query := `SELECT * FROM post_reaction WHERE post_id=$1 AND author_id=$2`
	err = r.db.QueryRowContext(ctx, query, postID, authorID).Scan(
		&reaction.AuthorID,
		&reaction.PostID,
		&reaction.Reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return reaction, nil
}

func (r *PostReactionRepository) DeletePostReaction(postID, authorID int) error {
	query := `DELETE FROM post_reaction WHERE post_id=$1 AND author_id=$2`
	result, err := r.db.Exec(query, postID, authorID)
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

func (r *PostReactionRepository) UpdatePostReaction(postReaction *models.PostReaction) error {
	query := `UPDATE post_reaction SET reaction=$1 WHERE post_id=$2 AND author_id=$3`
	result, err := r.db.Exec(query, postReaction.Reaction, postReaction.PostID, postReaction.AuthorID)
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

func (r *PostReactionRepository) AddPostReaction(postReaction *models.PostReaction) error {
	query := `INSERT INTO post_reaction(post_id, author_id, reaction) VALUES($1, $2, $3)`
	_, err := r.db.Exec(query, postReaction.PostID, postReaction.AuthorID, postReaction.Reaction)
	return err
}
