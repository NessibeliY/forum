package post_reaction

import "database/sql"

type PostReactionRepository struct {
	db *sql.DB
}

func NewPostReactionRepository(db *sql.DB) *PostReactionRepository {
	return &PostReactionRepository{
		db: db,
	}
}
