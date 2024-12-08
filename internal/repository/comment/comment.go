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
	query := `
	SELECT c.id, c.content, c.post_id, c.author_id, u.username AS author_name, c.created_at 
	FROM comment c
	JOIN users u ON c.author_id = u.id
	WHERE c.post_id=? 
	ORDER BY c.id DESC`
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
			&comment.AuthorName,
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

func (r *CommentRepository) DeleteComment(id int) error {
	query := `DELETE FROM comment WHERE id = $1`
	result, err := r.db.Exec(query, id)
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

func (r *CommentRepository) GetCommentByID(ctx context.Context, id int) (*models.Comment, error) {
	query := `SELECT * FROM comment WHERE id = $1 ORDER BY id DESC`

	comment := &models.Comment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.PostID,
		&comment.AuthorID,
		&comment.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	return comment, nil
}

func (r *CommentRepository) GetUserCommentedPosts(сtx context.Context, author_id int) ([]models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at, u.username
		FROM post p
		JOIN comment c ON c.post_id = p.id
		JOIN users u ON p.author_id = u.id 
		WHERE c.author_id = $1
		GROUP BY p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at, u.username
		ORDER BY p.id DESC;
	`

	rows, err := r.db.QueryContext(сtx, query, author_id)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	var posts []models.Post
	var currentPost *models.Post

	for rows.Next() {
		var postID, authorID int
		var title, content, authorName string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&postID,
			&title,
			&content,
			&authorID,
			&createdAt,
			&updatedAt,
			&authorName,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		if currentPost == nil || currentPost.ID != postID {
			if currentPost != nil {
				posts = append(posts, *currentPost)
			}
			currentPost = &models.Post{
				ID:         postID,
				Title:      title,
				Content:    content,
				AuthorID:   authorID,
				AuthorName: authorName,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			}
		}

	}

	if currentPost != nil {
		posts = append(posts, *currentPost)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return posts, nil
}
