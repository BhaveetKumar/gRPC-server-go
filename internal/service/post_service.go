package service

import (
	"context"

	"github.com/BhaveetKumar/gRPC-server-go/internal/domain"
	apperrors "github.com/BhaveetKumar/gRPC-server-go/internal/errors"
	"github.com/BhaveetKumar/gRPC-server-go/internal/repository"
	"github.com/google/uuid"
)

type postService struct {
	repo repository.PostRepository
}

var _ PostService = (*postService)(nil)

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) CreatePost(ctx context.Context, title, content, author, publicationDate string, tags []string) (*domain.Post, error) {
	post := &domain.Post{
		ID:              uuid.NewString(),
		Title:           title,
		Content:         content,
		Author:          author,
		PublicationDate: publicationDate,
		Tags:            tags,
	}

	if err := post.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	if id == "" {
		return nil, apperrors.ErrInvalidInput
	}

	return s.repo.GetByID(id)
}

func (s *postService) UpdatePost(ctx context.Context, id, title, content, author string, tags []string) (*domain.Post, error) {
	if id == "" {
		return nil, apperrors.ErrInvalidInput
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	existing.Title = title
	existing.Content = content
	existing.Author = author
	existing.Tags = tags

	if err := existing.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *postService) DeletePost(ctx context.Context, id string) error {
	if id == "" {
		return apperrors.ErrInvalidInput
	}

	return s.repo.Delete(id)
}
