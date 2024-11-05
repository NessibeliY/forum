package models

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
	CreateCommentReaction(commentReactionRequest *CommentReactionRequest) error
	UpdateCommentReaction(commentReactionRequest *CommentReactionRequest) error
	DeleteCommentReaction(commentID, authorID int) error
}

type CommentReactionRepository interface {
	AddCommentReaction(commentReaction *CommentReaction) error
	UpdateCommentReaction(commentReaction *CommentReaction) error
	DeleteCommentReaction(commentID, authorID int) error
}
