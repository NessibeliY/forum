package post

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	query := `
	SELECT
	    p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at,
	    c.id AS category_id, c.name AS category_name,
	    COUNT(CASE WHEN pr.reaction = 'like' THEN 1 END) AS likes_count,
	    COUNT(CASE WHEN pr.reaction = 'dislike' THEN 1 END) AS dislikes_count,
	    COUNT(co.id) AS comments_count
	FROM
	    post p
	LEFT JOIN
	        post_category pc ON p.id = pc.post_id
	LEFT JOIN
	        category c ON pc.category_id = c.id
	LEFT JOIN
	        post_reaction pr ON p.id = pr.post_id
	LEFT JOIN
	        comment co ON p.id = co.post_id
	GROUP BY
	    p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at, c.id, c.name
	ORDER BY p.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	var currentPost *models.Post

	for rows.Next() {
		var postID, categoryID, likesCount, dislikesCount, commentsCount int
		var title, content, categoryName string
		var authorID int
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&postID,
			&title,
			&content,
			&authorID,
			&createdAt,
			&updatedAt,
			&categoryID,
			&categoryName,
			&likesCount,
			&dislikesCount,
			&commentsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		if currentPost == nil || currentPost.ID != postID {
			if currentPost != nil {
				posts = append(posts, *currentPost)
			}
			currentPost = &models.Post{
				ID:            postID,
				Title:         title,
				Content:       content,
				AuthorID:      authorID,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
				Categories:    []models.Category{},
				LikesCount:    likesCount,
				DislikesCount: dislikesCount,
				CommentsCount: commentsCount,
			}
		}

		if categoryID != 0 {
			category := models.Category{
				ID:   categoryID,
				Name: categoryName,
			}
			currentPost.Categories = append(currentPost.Categories, category)
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
