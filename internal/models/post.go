package models

import (
	"context"
	"time"
)

type Post struct {
	ID            int
	Title         string
	Content       string
	AuthorID      int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Categories    []*Category
	Comments      []Comment
	LikesCount    int
	DislikesCount int
	CommentsCount int
}

type CreatePostRequest struct {
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	AuthorID   int         `json:"author_id"`
	Categories []*Category `json:"categories"`
}

type UpdatePostRequest struct {
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Categories []*Category `json:"categories"`
}

type DeletePostRequest struct {
	ID int `json:"id"`
}

type PostService interface {
	GetAllPosts() ([]Post, error)
	CreatePost(request *CreatePostRequest) (int, error)
	//UpdatePost(request *UpdatePostRequest) error
	//DeletePost(request *DeletePostRequest) error
}

type PostRepository interface {
	GetAllPosts(ctx context.Context) ([]Post, error)
	AddPost(post *Post) (int, error)
	//UpdatePost(post *Post) error
	//DeletePost(id int) error
}
