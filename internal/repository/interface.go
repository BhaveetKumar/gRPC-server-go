package repository

import "github.com/BhaveetKumar/gRPC-server-go/internal/domain"

type PostRepository interface {
	Create(post *domain.Post) error
	GetByID(id string) (*domain.Post, error)
	Update(post *domain.Post) error
	Delete(id string) error
	List() ([]*domain.Post, error)
}
