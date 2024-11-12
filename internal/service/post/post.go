package post

import (
	"context"
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

	return s.repo.GetAllPosts(ctx)
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
