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

func (r *CommentRepository) GetUserCommentedPosts(ctx context.Context, authorID int) ([]models.Post, error) {
	query := `
		SELECT 
			p.id AS post_id, 
			p.title, 
			p.content, 
			p.author_id AS post_author_id, 
			p.created_at AS post_created_at, 
			p.updated_at AS post_updated_at, 
			u.username AS post_author_name,
			c.id AS comment_id, 
			c.content AS comment_content, 
			c.author_id AS comment_author_id
		FROM post p
		JOIN comment c ON c.post_id = p.id
		JOIN users u ON p.author_id = u.id
		WHERE c.author_id = $1
		ORDER BY p.id DESC, c.created_at ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	var currentPost *models.Post

	for rows.Next() {
		var (
			postID, postAuthorID, commentID, commentAuthorID       int
			postTitle, postContent, postAuthorName, commentContent string
			postCreatedAt, postUpdatedAt                           time.Time
		)

		err := rows.Scan(
			&postID,
			&postTitle,
			&postContent,
			&postAuthorID,
			&postCreatedAt,
			&postUpdatedAt,
			&postAuthorName,
			&commentID,
			&commentContent,
			&commentAuthorID,
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
				Title:      postTitle,
				Content:    postContent,
				AuthorID:   postAuthorID,
				AuthorName: postAuthorName,
				CreatedAt:  postCreatedAt,
				UpdatedAt:  postUpdatedAt,
				Comments:   []models.Comment{}, // Инициализируем список комментариев
			}
		}

		// Добавляем комментарий к текущему посту
		currentPost.Comments = append(currentPost.Comments, models.Comment{
			ID:       commentID,
			Content:  commentContent,
			AuthorID: commentAuthorID,
		})
	}

	// Добавляем последний пост в список
	if currentPost != nil {
		posts = append(posts, *currentPost)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return posts, nil
}

func (r *CommentRepository) UpdateComment(updateCommentRequest *models.UpdateCommentRequest) error {
	query := `
	UPDATE comment
	SET content = $1, created_at = $2
	WHERE id = $3
	`

	result, err := r.db.Exec(query, updateCommentRequest.Content, time.Now(), updateCommentRequest.ID)
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
