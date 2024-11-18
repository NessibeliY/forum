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

func (s *CommentReactionService) CreateCommentReaction(request *models.CreateCommentReactionRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	commentReaction := &models.CommentReaction{
		AuthorID:  request.AuthorID,
		CommentID: request.CommentID,
		Reaction:  request.Reaction,
	}

	currentReaction, err := s.repo.GetReactionByCommentIDAndAuthorID(ctx, request.CommentID, request.AuthorID)
	if err != nil {
		return fmt.Errorf("get reaction by comment id: %v", err)
	}
	if currentReaction != nil {
		if currentReaction.Reaction == request.Reaction {
			err = s.repo.DeleteCommentReaction(commentReaction)
			if err != nil {
				return fmt.Errorf("delete comment reaction: %v", err)
			}
			return nil
		}

		err = s.repo.UpdateCommentReaction(commentReaction)
		if err != nil {
			return fmt.Errorf("update comment reaction: %v", err)
		}
		return nil
	}

	err = s.repo.AddCommentReaction(commentReaction)
	if err != nil {
		return fmt.Errorf("add comment reaction: %v", err)
	}

	return nil
}
