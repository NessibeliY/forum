package comment_reaction

import (
	"context"
	"fmt"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type CommentReactionService struct {
	repo models.CommentReactionRepository
}

func NewCommentReactionService(repo models.CommentReactionRepository) *CommentReactionService {
	return &CommentReactionService{
		repo: repo,
	}
}

func (s *CommentReactionService) GetCommentLikesAndDislikesByID(commentID int) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reactions, err := s.repo.GetReactionsByCommentID(ctx, commentID)
	if err != nil {
		return 0, 0, fmt.Errorf("get reactions by comment id: %v", err)
	}

	var likesCount, dislikesCount int
	for _, reaction := range reactions {
		if reaction.Reaction == "like" {
			likesCount++
		} else {
			dislikesCount++
		}
	}

	return likesCount, dislikesCount, nil
}
