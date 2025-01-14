package moderation

import (
	"context"
	"database/sql"
	"fmt"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type ModerationRepository struct {
	db *sql.DB
}

func NewModerationRepository(db *sql.DB) *ModerationRepository {
	return &ModerationRepository{
		db: db,
	}
}

func (r *ModerationRepository) AddModerationReport(report *models.ModerationReport) error {
	query := `
		INSERT INTO moderated_post (post_id, moderator_id, moderated) 
		VALUES ($1, $2, $3)
		ON CONFLICT (post_id, moderator_id) DO NOTHING;`
	_, err := r.db.Exec(query, report.PostID, report.ModeratorID, report.IsModerated)
	if err != nil {
		return err
	}
	return nil
}

func (r *ModerationRepository) DeleteModerationReport(report *models.ModerationReport) error {
	query := `DELETE FROM moderated_post WHERE post_id = ? AND moderator_id = ?;`
	result, err := r.db.Exec(query, report.PostID, report.ModeratorID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if count == 0 {
		return nil
	}

	return nil
}

func (r *ModerationRepository) GetAllModeratedPosts(ctx context.Context) ([]models.ModerationReport, error) {
	query := `
	SELECT
		mp.id AS moderation_id,
		mp.post_id,
		p.title AS post_title,
		p.content AS post_content,
		p.author_id AS post_author_id,
		u1.username AS post_author_name,
		p.created_at AS post_created_at,
		p.updated_at AS post_updated_at,
		mp.moderator_id,
		u2.username AS moderator_name,
		mp.admin_answer,
		mp.moderated
	FROM
		moderated_post mp
	JOIN
		post p ON mp.post_id = p.id
	LEFT JOIN
		users u1 ON p.author_id = u1.id
	LEFT JOIN
		users u2 ON mp.moderator_id = u2.id
	ORDER BY mp.id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query execution: %w", err)
	}
	defer rows.Close()

	var reports []models.ModerationReport

	for rows.Next() {
		var report models.ModerationReport
		var post models.Post
		var moderator models.User
		var adminAnswer sql.NullString

		err := rows.Scan(
			&report.ID,
			&report.PostID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CreatedAt,
			&post.UpdatedAt,
			&report.ModeratorID,
			&moderator.Username,
			&adminAnswer,
			&report.IsModerated,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		if adminAnswer.Valid {
			report.AdminAnswer = adminAnswer.String
		} else {
			report.AdminAnswer = ""
		}

		report.Post = &post
		report.Moderator = &moderator

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return reports, nil
}
