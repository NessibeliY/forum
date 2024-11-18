package comment

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (r *CommentRepository) GetAllCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error) {
	query := `SELECT id, content, post_id, author_id, created_at FROM comment WHERE post_id=? ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	comments := make([]*models.Comment, 0)
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.PostID,
			&comment.AuthorID,
			&comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return comments, nil
}

func (r *CommentRepository) AddComment(comment *models.Comment) error {
	createdAt := time.Now()
	query := `INSERT INTO comment (content, post_id, author_id, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, comment.Content, comment.PostID, comment.AuthorID, createdAt)
	return err
}
