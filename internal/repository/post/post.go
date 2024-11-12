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
	    posts p
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
				Categories:    []*models.Category{},
				LikesCount:    likesCount,
				DislikesCount: dislikesCount,
				CommentsCount: commentsCount,
			}
		}

		if categoryID != 0 {
			category := &models.Category{
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

func (r *PostRepository) AddPost(post *models.Post) (int, error) {
	createdAt := time.Now()
	updatedAt := createdAt
	query := `
	INSERT INTO post (title, content, author_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`
	err := r.db.QueryRow(query, post.Title, post.Content, post.AuthorID, createdAt, updatedAt).Scan(&post.ID)
	if err != nil {
		return 0, fmt.Errorf("insert post: %w", err)
	}

	err = r.addPostCategories(post.ID, post.Categories)
	if err != nil {
		return 0, fmt.Errorf("add post categories: %w", err)
	}

	return post.ID, nil
}

func (r *PostRepository) addPostCategories(postID int, categories []*models.Category) error {
	for _, category := range categories {
		query := `INSERT INTO post_category (post_id, category_id) VALUES ($1, $2)`
		_, err := r.db.Exec(query, postID, category.ID)
		if err != nil {
			return fmt.Errorf("insert post_category: %w", err)
		}
	}
	return nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	query := `SELECT * FROM post WHERE id = $1`

	post := &models.Post{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	query = `
	SELECT c.name FROM category c2 
	JOIN post_category pc ON pc.post_id = c2.id
	WHERE c2.id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		category := models.Category{}
		err := rows.Scan(&category.Name)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		post.Categories = append(post.Categories, &category)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return post, nil
}
