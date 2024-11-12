package post_reaction

import (
	"context"
	"fmt"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type PostReactionService struct {
	repo models.PostReactionRepository
}

func NewPostReactionService(repo models.PostReactionRepository) *PostReactionService {
	return &PostReactionService{
		repo: repo,
	}
}

func (s *PostReactionService) GetPostLikesAndDislikesByID(postID int) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	reactions, err := s.repo.GetReactionsByPostID(ctx, postID)
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
