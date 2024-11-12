package models

import (
	"context"
	"time"
)

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
	GetAllCommentsByPostID(postID int) ([]*Comment, error)
	//CreateComment(createCommentRequest *CreateCommentRequest) error
}

type CommentRepository interface {
	GetAllCommentsByPostID(ctx context.Context, postID int) ([]*Comment, error)
	//AddComment(comment *Comment) error
}
