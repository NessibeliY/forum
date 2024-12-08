package models

import (
	"context"
)

type PostReaction struct {
	AuthorID int
	PostID   int
	Reaction string
}

type CreatePostReactionRequest struct {
	AuthorID int    `json:"author_id"`
	PostID   int    `json:"post_id"`
	Reaction string `json:"reaction"`
}

type PostReactionService interface {
	GetPostLikesAndDislikesByID(postID int) (int, int, error)
	CreatePostReaction(request *CreatePostReactionRequest) error
	GetUserReactionPosts(author_id int) ([]UserReactionPost, error)
}

type PostReactionRepository interface {
	GetReactionsByPostID(ctx context.Context, postID int) (reactions []*PostReaction, err error)
	GetReactionByPostIDAndAuthorID(ctx context.Context, postID, authorID int) (reaction *PostReaction, err error)
	AddPostReaction(postReaction *PostReaction) error
	UpdatePostReaction(postReaction *PostReaction) error
	DeletePostReaction(postReaction *PostReaction) error
	GetUserReactionPosts(ctx context.Context, author_id int) ([]UserReactionPost, error)
}
