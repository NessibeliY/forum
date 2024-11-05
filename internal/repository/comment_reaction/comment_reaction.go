package comment_reaction

import "database/sql"

type CommentReactionRepository struct {
	db *sql.DB
}

func NewCommentReactionRepository(db *sql.DB) *CommentReactionRepository {
	return &CommentReactionRepository{
		db: db,
	}
}
