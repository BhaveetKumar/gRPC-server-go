package service

import (
	"context"

	"github.com/BhaveetKumar/gRPC-server-go/internal/domain"
)

type PostService interface {
	CreatePost(ctx context.Context, title, content, author, publicationDate string, tags []string) (*domain.Post, error)
	GetPost(ctx context.Context, id string) (*domain.Post, error)
	UpdatePost(ctx context.Context, id, title, content, author string, tags []string) (*domain.Post, error)
	DeletePost(ctx context.Context, id string) error
}
