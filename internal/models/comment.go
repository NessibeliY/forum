package models

import "time"

type Comment struct {
	ID            int
	Content       string
	AuthorID      int
	PostID        int
	CreatedAt     time.Time
	LikesCount    int
	DislikesCount int
}

type CreateCommentRequest struct {
	Content  string `json:"content"`
	AuthorID int    `json:"author_id"`
	PostID   int    `json:"post_id"`
}

type CommentService interface {
	CreateComment(createCommentRequest *CreateCommentRequest) error
}

type CommentRepository interface {
	AddComment(comment *Comment) error
}
