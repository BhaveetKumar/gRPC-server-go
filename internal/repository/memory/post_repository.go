package memory

import (
	"sync"

	"github.com/BhaveetKumar/gRPC-server-go/internal/domain"
	apperrors "github.com/BhaveetKumar/gRPC-server-go/internal/errors"
	"github.com/BhaveetKumar/gRPC-server-go/internal/repository"
)

type PostRepository struct {
	mu    sync.RWMutex
	posts map[string]*domain.Post
}

var _ repository.PostRepository = (*PostRepository)(nil)

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts: make(map[string]*domain.Post),
	}
}

func (r *PostRepository) Create(post *domain.Post) error {
	if post == nil {
		return apperrors.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; exists {
		return apperrors.ErrDuplicatePost
	}

	copy := *post
	r.posts[post.ID] = &copy

	return nil
}

func (r *PostRepository) GetByID(id string) (*domain.Post, error) {
	if id == "" {
		return nil, apperrors.ErrInvalidInput
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return nil, apperrors.ErrPostNotFound
	}

	copy := *post
	return &copy, nil
}

func (r *PostRepository) Update(post *domain.Post) error {
	if post == nil || post.ID == "" {
		return apperrors.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[post.ID]; !ok {
		return apperrors.ErrPostNotFound
	}

	copy := *post
	r.posts[post.ID] = &copy

	return nil
}

func (r *PostRepository) Delete(id string) error {
	if id == "" {
		return apperrors.ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[id]; !ok {
		return apperrors.ErrPostNotFound
	}

	delete(r.posts, id)
	return nil
}

func (r *PostRepository) List() ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Post, 0, len(r.posts))
	for _, post := range r.posts {
		copy := *post
		result = append(result, &copy)
	}

	return result, nil
}
