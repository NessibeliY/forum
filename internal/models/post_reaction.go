package models

type PostReaction struct {
	AuthorID     int
	PostID       int
	ReactionType string
}

type PostReactionRequest struct {
	AuthorID     int    `json:"author_id"`
	PostID       int    `json:"post_id"`
	ReactionType string `json:"reaction_type"`
}

type PostReactionService interface {
	CreatePostReaction(request *PostReactionRequest) error
	UpdatePostReaction(request *PostReactionRequest) error
	DeletePostReaction(request *PostReactionRequest) error
}

type PostReactionRepository interface {
	AddPostReaction(postReaction *PostReaction) error
	UpdatePostReaction(postReaction *PostReaction) error
	DeletePostReaction(postID, authorID int) error
}
