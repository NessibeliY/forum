package models

import "context"

type CommentReaction struct {
	CommentID int
	AuthorID  int
	Reaction  string
}

type CommentReactionRequest struct {
	CommentID int    `json:"comment_id"`
	AuthorID  int    `json:"author_id"`
	Reaction  string `json:"reaction"`
}

type CommentReactionService interface {
	GetCommentLikesAndDislikesByID(commentID int) (int, int, error)
	//CreateCommentReaction(commentReactionRequest *CommentReactionRequest) error
	//UpdateCommentReaction(commentReactionRequest *CommentReactionRequest) error
	//DeleteCommentReaction(commentID, authorID int) error
}

type CommentReactionRepository interface {
	GetReactionsByCommentID(ctx context.Context, commentID int) (reactions []*CommentReaction, err error)
	//AddCommentReaction(commentReaction *CommentReaction) error
	//UpdateCommentReaction(commentReaction *CommentReaction) error
	//DeleteCommentReaction(commentID, authorID int) error
}
