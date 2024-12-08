package post_reaction

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (r *PostReactionRepository) DeletePostReaction(postReaction *models.PostReaction) error {
	query := `DELETE FROM post_reaction WHERE post_id=$1 AND author_id=$2`
	result, err := r.db.Exec(query, postReaction.PostID, postReaction.AuthorID)
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

func (r *PostReactionRepository) GetUserReactionPosts(ctx context.Context, author_id int) ([]models.UserReactionPost, error) {
	query := `
		SELECT
    p.id, p.title, p.content, p.author_id, u.username AS author_name, p.created_at, p.updated_at,
    c.id AS category_id, c.name AS category_name,
    (SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'like') AS likes_count,
    (SELECT COUNT(*) FROM post_reaction pr WHERE pr.post_id = p.id AND pr.reaction = 'dislike') AS dislikes_count,
    (SELECT COUNT(*) FROM comment co WHERE co.post_id = p.id) AS comments_count,
    (SELECT reaction FROM post_reaction pr WHERE pr.post_id = p.id AND pr.author_id = $1 LIMIT 1) AS user_reaction
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
        AND pr2.reaction IN ('like', 'dislike')
    )
GROUP BY
    p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at, c.id, c.name
ORDER BY p.id DESC

	`

	rows, err := r.db.QueryContext(ctx, query, author_id)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var posts []models.UserReactionPost
	var currentPost *models.UserReactionPost

	for rows.Next() {
		var postID, categoryID, likesCount, dislikesCount, commentsCount int
		var title, content, categoryName, authorName, userReaction string
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
			&userReaction,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		if currentPost == nil || currentPost.ID != postID {
			if currentPost != nil {
				posts = append(posts, *currentPost)
			}
			currentPost = &models.UserReactionPost{
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
				UserReaction:  userReaction,
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
