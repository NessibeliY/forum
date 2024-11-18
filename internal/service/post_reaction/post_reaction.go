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

func (s *PostReactionService) CreatePostReaction(request *models.CreatePostReactionRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	postReaction := &models.PostReaction{
		AuthorID: request.AuthorID,
		PostID:   request.PostID,
		Reaction: request.Reaction,
	}

	currentReaction, err := s.repo.GetReactionByPostIDAndAuthorID(ctx, request.PostID, request.AuthorID)
	if err != nil {
		return fmt.Errorf("get reaction by post id: %v", err)
	}
	if currentReaction != nil {
		if currentReaction.Reaction == request.Reaction {
			err = s.repo.DeletePostReaction(request.PostID, request.AuthorID)
			if err != nil {
				return fmt.Errorf("delete post reaction: %v", err)
			}
			return nil
		}

		err = s.repo.UpdatePostReaction(postReaction)
		if err != nil {
			return fmt.Errorf("update post reaction: %v", err)
		}
	}

	err := s.repo.AddPostReaction(postReaction)
	if err != nil {
		return fmt.Errorf("add post reaction: %v", err)
	}

	return nil
}
