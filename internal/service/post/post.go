package post

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"

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

func (s *PostService) CreatePostWithImage(request *models.CreatePostRequest) (int, error) {
	if request.ImageFile == nil {
		return s.CreatePost(request)
	}

	data, err := io.ReadAll(request.ImageFile)
	if err != nil {
		return 0, fmt.Errorf("read image file: %w", err)
	}

	fileName, err := uuid.NewV4()
	if err != nil {
		return 0, fmt.Errorf("generate uuid: %w", err)
	}

	filePath := "ui/static/img/" + fileName.String()

	post := &models.Post{
		Title:      request.Title,
		Content:    request.Content,
		AuthorID:   request.AuthorID,
		Categories: request.Categories,
		ImagePath:  filePath,
	}

	id, err := s.repo.AddPostWithImage(post)
	if err != nil {
		return 0, fmt.Errorf("add post with image: %w", err)
	}

	err = os.WriteFile(filePath, data, 0o666)
	if err != nil {
		rollbackErr := s.repo.DeletePostWithImage(id)
		if rollbackErr != nil {
			return 0, fmt.Errorf("rollback delete post with image: %w", rollbackErr)
		}
		return 0, fmt.Errorf("write file: %w", err)
	}

	return id, nil
}

func (s *PostService) GetPostByID(id int) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	if post.ImagePath == "" {
		return post, nil
	}

	post.ImagePath = ".." + strings.TrimPrefix(post.ImagePath, "ui")

	return post, nil
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
