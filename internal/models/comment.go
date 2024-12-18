package models

import (
	"context"
	"time"
)

type Comment struct {
	ID            int
	Content       string
	AuthorID      int
	AuthorName    string
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

type DeleteCommentRequest struct {
	ID int `json:"id"`
}

type UpdateCommentRequest struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CommentService interface {
	GetAllCommentsByPostID(postID int) ([]*Comment, error)
	CreateComment(createCommentRequest *CreateCommentRequest) error
	DeleteComment(deleteCommentRequest *DeleteCommentRequest) error
	GetCommentByID(id int) (*Comment, error)
	GetUserCommentedPosts(authorID int) ([]Post, error)
	UpdateComment(updateCommentRequest *UpdateCommentRequest) error
}

type CommentRepository interface {
	GetAllCommentsByPostID(ctx context.Context, postID int) ([]*Comment, error)
	AddComment(comment *Comment) error
	DeleteComment(id int) error
	GetCommentByID(ctx context.Context, id int) (*Comment, error)
	GetUserCommentedPosts(сtx context.Context, authorID int) ([]Post, error)
	UpdateComment(updateCommentRequest *UpdateCommentRequest) error
}
