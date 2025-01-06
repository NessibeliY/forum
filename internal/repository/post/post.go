package post

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
	    p.id, p.title, p.content, p.author_id, u.username AS author_name, p.created_at, p.updated_at,
	    c.id AS category_id, c.name AS category_name,
	    (SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'like') AS likes_count,
    	(SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'dislike') AS dislikes_count,
	    (SELECT COUNT(*) FROM comment co WHERE co.post_id = p.id) AS comments_count
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
	LEFT JOIN
	        users u ON p.author_id = u.id
	GROUP BY
	    p.id, p.title, p.content, p.author_id, u.username, p.created_at, p.updated_at, c.id, c.name
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
		var title, content, categoryName, authorName string
		var authorID int
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&postID,
			&title,
			&content,
			&authorID,
			&authorName,
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
				AuthorName:    authorName,
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

func (r *PostRepository) DeletePostWithImage(id int) error {
	query := `DELETE FROM post WHERE id = $1`
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

	query = `DELETE FROM image WHERE post_id = $1`
	result, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("exec image: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected image: %w", err)
	}

	if rowsAffected == 0 {
		return nil
	}

	return nil
}

func (r *PostRepository) AddPostWithImage(post *models.Post) (int, error) {
	createdAt := time.Now()
	updatedAt := createdAt
	query := `
	INSERT INTO post (title, content, author_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;`

	err := r.db.QueryRow(query, post.Title, post.Content, post.AuthorID, createdAt, updatedAt).Scan(&post.ID)
	if err != nil {
		return 0, fmt.Errorf("insert post: %w", err)
	}

	err = r.addPostCategories(post.ID, post.Categories)
	if err != nil {
		return 0, fmt.Errorf("add post categories: %w", err)
	}

	query = `INSERT INTO image (post_id, image_path) VALUES ($1, $2);`
	_, err = r.db.Exec(query, post.ID, post.ImagePath)
	if err != nil {
		return 0, fmt.Errorf("insert post image: %w", err)
	}

	return post.ID, nil
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

func (r *PostRepository) UpdatePost(post *models.Post) (int, error) {
	query := `
		UPDATE post
		SET title = $1, content = $2, updated_at = $3
		WHERE author_id = $4
		RETURNING id;
	`

	err := r.db.QueryRow(query, post.Title, post.Content, time.Now(), post.AuthorID).Scan(&post.ID)
	if err != nil {
		return 0, fmt.Errorf("update post: %w", err)
	}

	if post.ID == 0 {
		return 0, fmt.Errorf("no rows updated, post with author_id %d not found", post.AuthorID)
	}

	err = r.updateCategories(post.ID, post.Categories)
	if err != nil {
		return 0, fmt.Errorf("update post categories: %w", err)
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

func (r *PostRepository) updateCategories(postID int, categories []*models.Category) error {
	deleteQuery := `DELETE FROM post_category WHERE post_id = $1`
	_, err := r.db.Exec(deleteQuery, postID)
	if err != nil {
		return fmt.Errorf("delete old post categories: %w", err)
	}

	for _, category := range categories {
		insertQuery := `
			INSERT INTO post_category (post_id, category_id) 
			VALUES ($1, $2)
			ON CONFLICT (post_id, category_id) 
			DO NOTHING
		`

		_, err := r.db.Exec(insertQuery, postID, category.ID)
		if err != nil {
			return fmt.Errorf("insert post_category: %w", err)
		}
	}
	return nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	query := `SELECT * FROM post WHERE id = $1 ORDER BY id DESC`

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
	SELECT c2.name 
	FROM category c2 
	JOIN post_category pc ON pc.category_id = c2.id
	WHERE pc.post_id = $1`

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

	query = `SELECT image_path FROM image WHERE post_id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	err = row.Scan(&post.ImagePath)
	if err != nil {
		if err == sql.ErrNoRows {
			return post, nil
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	return post, nil
}

func (r *PostRepository) GetPostsByAuthorID(ctx context.Context, authorID int) ([]models.Post, error) {
	query := `
	SELECT
	    p.id, p.title, p.content, p.author_id, u.username AS author_name, p.created_at, p.updated_at,
	    c.id AS category_id, c.name AS category_name,
	    (SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'like') AS likes_count,
    	(SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'dislike') AS dislikes_count,
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
	LEFT JOIN
	        users u ON p.author_id = u.id
	WHERE
	    p.author_id = $1
	GROUP BY
	    p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at, c.id, c.name
	ORDER BY p.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	var currentPost *models.Post

	for rows.Next() {
		var postID, categoryID, likesCount, dislikesCount, commentsCount int
		var title, content, categoryName, authorName string
		var authorID int
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&postID,
			&title,
			&content,
			&authorID,
			&authorName,
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
				AuthorName:    authorName,
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

func (r *PostRepository) GetLikedPosts(ctx context.Context, userID int) ([]models.Post, error) {
	query := `
	SELECT
	    p.id, p.title, p.content, p.author_id, u.username AS author_name, p.created_at, p.updated_at,
	    c.id AS category_id, c.name AS category_name,
	    (SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'like') AS likes_count,
    	(SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'dislike') AS dislikes_count,
	    (SELECT COUNT(*) FROM comment co WHERE co.post_id = p.id) AS comments_count
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
	LEFT JOIN
	    users u ON p.author_id = u.id
	WHERE
	    EXISTS (
	        SELECT 1
	        FROM post_reaction pr2
	        WHERE pr2.post_id = p.id
	          AND pr2.author_id = $1
	          AND pr2.reaction = 'like'
	    )
	GROUP BY
	    p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at, c.id, c.name
	ORDER BY p.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	var currentPost *models.Post

	for rows.Next() {
		var postID, categoryID, likesCount, dislikesCount, commentsCount int
		var title, content, categoryName, authorName string
		var authorID int
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&postID,
			&title,
			&content,
			&authorID,
			&authorName,
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
				AuthorName:    authorName,
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

func (r *PostRepository) GetPostsByCategories(ctx context.Context, categories []string) ([]models.Post, error) {
	query := `
	WITH FilteredPosts AS (
		SELECT
			p.id
		FROM
			post p
		LEFT JOIN
			post_category pc ON p.id = pc.post_id
		LEFT JOIN
			category c ON pc.category_id = c.id
		WHERE
			c.name IN (` + placeholders(len(categories)) + `)
		GROUP BY
			p.id
	)
	SELECT
		p.id, p.title, p.content, p.author_id, u.username AS author_name, p.created_at, p.updated_at,
		c.id AS category_id, c.name AS category_name,
		(SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'like') AS likes_count,
		(SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'dislike') AS dislikes_count,
		(SELECT COUNT(*) FROM comment co WHERE co.post_id = p.id) AS comments_count
	FROM
		post p
	LEFT JOIN
		post_category pc ON p.id = pc.post_id
	LEFT JOIN
		category c ON pc.category_id = c.id
	LEFT JOIN
		users u ON p.author_id = u.id
	WHERE
		p.id IN (SELECT id FROM FilteredPosts)
	ORDER BY
		p.id DESC	
	`

	args := make([]interface{}, len(categories))
	for i, category := range categories {
		args[i] = category
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	var currentPost *models.Post

	for rows.Next() {
		var postID, categoryID, likesCount, dislikesCount, commentsCount int
		var title, content, categoryName, authorName string
		var authorID int
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&postID,
			&title,
			&content,
			&authorID,
			&authorName,
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
				AuthorName:    authorName,
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

func placeholders(count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat("?, ", count-1) + "?"
}

func (r *PostRepository) DeletePost(id int) error {
	query := `DELETE FROM post WHERE id = $1`
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
