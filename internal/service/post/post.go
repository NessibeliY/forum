package post

import (
	"context"
	"fmt"
	"time"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type PostService struct {
	repo models.PostRepository
}

func NewPostService(repo models.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	posts, err := s.repo.GetAllPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all posts: %w", err)
	}

	return posts, nil
}

func (s *PostService) CreatePost(createPostRequest *models.CreatePostRequest) (int, error) {
	post := &models.Post{
		Title:      createPostRequest.Title,
		Content:    createPostRequest.Content,
		AuthorID:   createPostRequest.AuthorID,
		Categories: createPostRequest.Categories,
	}
	return s.repo.AddPost(post)
}

func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetPostByID(ctx, id)
}

func (s *PostService) GetPostsByAuthorID(authorID int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetPostsByAuthorID(ctx, authorID)
}

func (s *PostService) GetLikedPosts(userID int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetLikedPosts(ctx, userID)
}

func (s *PostService) GetPostsByCategories(categories []string) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.repo.GetPostsByCategories(ctx, categories)
}

func (s *PostService) DeletePost(request *models.DeletePostRequest) error {
	return s.repo.DeletePost(request.ID)
}
