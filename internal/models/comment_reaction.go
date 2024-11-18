package models

import "context"

type CommentReaction struct {
	CommentID int
	AuthorID  int
	Reaction  string
}

type CreateCommentReactionRequest struct {
	CommentID int    `json:"comment_id"`
	AuthorID  int    `json:"author_id"`
	Reaction  string `json:"reaction"`
}

type CommentReactionService interface {
	GetCommentLikesAndDislikesByID(commentID int) (int, int, error)
	CreateCommentReaction(createCommentReactionRequest *CreateCommentReactionRequest) error
}

type CommentReactionRepository interface {
	GetReactionsByCommentID(ctx context.Context, commentID int) (reactions []*CommentReaction, err error)
	GetReactionByCommentIDAndAuthorID(ctx context.Context, commentID int, authorID int) (reaction *CommentReaction, err error)
	AddCommentReaction(commentReaction *CommentReaction) error
	UpdateCommentReaction(commentReaction *CommentReaction) error
	DeleteCommentReaction(commentReaction *CommentReaction) error
}
