package models

import (
	"context"
	"mime/multipart"
	"time"
)

type Post struct {
	ID            int
	Title         string
	Content       string
	AuthorID      int
	AuthorName    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Categories    []*Category
	Comments      []Comment
	LikesCount    int
	DislikesCount int
	CommentsCount int
	ImagePath     string
}

type UserReactionPost struct {
	ID            int
	Title         string
	Content       string
	AuthorID      int
	AuthorName    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Categories    []*Category
	Comments      []Comment
	LikesCount    int
	DislikesCount int
	CommentsCount int
	UserReaction  string
}

type CreatePostRequest struct {
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	AuthorID   int            `json:"author_id"`
	Categories []*Category    `json:"categories"`
	ImageFile  multipart.File `json:"image_file"`
}

type UpdatePostRequest struct {
	PostID     int
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	AuthorID   int            `json:"author_id"`
	Categories []*Category    `json:"categories"`
	ImageFile  multipart.File `json:"image_file"`
}

type DeletePostRequest struct {
	ID int `json:"id"`
}

type SendReportRequest struct {
	PostID      int
	IsModerated bool
	ModeratorID int
	Post        *Post
	Moderator   *User
}

type ModerationReport struct {
	ID          int
	PostID      int
	ModeratorID int
	AdminAnswer string
	IsModerated bool
	Post        *Post
	Moderator   *User
}

type PostService interface {
	GetAllPosts() ([]Post, error)
	CreatePost(request *CreatePostRequest) (int, error)
	CreatePostWithImage(request *CreatePostRequest) (int, error)
	GetPostByID(id int) (*Post, error)
	GetPostsByAuthorID(authorID int) ([]Post, error)
	GetLikedPosts(userID int) ([]Post, error)
	GetPostsByCategories(categories []string) ([]Post, error)
	DeletePost(request *DeletePostRequest) error
	UpdatePost(request *UpdatePostRequest) (int, error)
	UpdatePostWithImage(request *UpdatePostRequest) (int, error)
	SendReport(request *SendReportRequest) error
	GetAllModeratedPosts() ([]ModerationReport, error)
	UpdateModerationReport(report *ModerationReport) error
}

type PostRepository interface {
	GetAllPosts(ctx context.Context) ([]Post, error)
	AddPost(post *Post) (int, error)
	AddPostWithImage(post *Post) (int, error)
	GetPostByID(ctx context.Context, id int) (*Post, error)
	GetPostsByAuthorID(ctx context.Context, authorID int) ([]Post, error)
	GetLikedPosts(ctx context.Context, userID int) ([]Post, error)
	GetPostsByCategories(ctx context.Context, categories []string) ([]Post, error)
	DeletePost(id int) error
	UpdatePost(post *Post) (int, error)
	UpdatePostWithImage(post *Post) (int, error)
	DeletePostWithImage(id int) error
}

type ModerationRepository interface {
	AddModerationReport(report *ModerationReport) error
	DeleteModerationReport(report *ModerationReport) error
	GetAllModeratedPosts(ctx context.Context) ([]ModerationReport, error)
}
