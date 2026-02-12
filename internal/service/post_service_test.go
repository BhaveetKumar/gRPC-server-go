package service

import (
	"context"
	"testing"

	apperrors "github.com/BhaveetKumar/gRPC-server-go/internal/errors"
	"github.com/BhaveetKumar/gRPC-server-go/internal/repository/memory"
)

func TestPostService_CreateValidate(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	post, err := service.CreatePost(ctx, "title", "content", "author", "", nil)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if post.ID == "" {
		t.Fatalf("expected id to be set")
	}
	if post.Title != "title" {
		t.Fatalf("expected title to be set")
	}
}

func TestPostService_CreateWithTags(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	tags := []string{"golang", "grpc", "testing"}
	post, err := service.CreatePost(ctx, "title", "content", "author", "2026-02-12", tags)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if len(post.Tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(post.Tags))
	}
	if post.PublicationDate != "2026-02-12" {
		t.Fatalf("unexpected publication date: %s", post.PublicationDate)
	}
}

func TestPostService_CreateInvalidEmptyTitle(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	_, err := service.CreatePost(ctx, "", "content", "author", "", nil)
	if err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestPostService_CreateInvalidEmptyContent(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	_, err := service.CreatePost(ctx, "title", "", "author", "", nil)
	if err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestPostService_CreateInvalidEmptyAuthor(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	_, err := service.CreatePost(ctx, "title", "content", "", "", nil)
	if err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestPostService_GetPost(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	created, _ := service.CreatePost(ctx, "title", "content", "author", "", nil)

	retrieved, err := service.GetPost(ctx, created.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Fatalf("id mismatch: got %s, want %s", retrieved.ID, created.ID)
	}
	if retrieved.Title != "title" {
		t.Fatalf("title mismatch: %s", retrieved.Title)
	}
}

func TestPostService_GetPostNotFound(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	_, err := service.GetPost(ctx, "nonexistent")
	if err != apperrors.ErrPostNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestPostService_GetPostEmptyID(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	_, err := service.GetPost(ctx, "")
	if err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	created, _ := service.CreatePost(ctx, "original", "original content", "author1", "", nil)

	updated, err := service.UpdatePost(ctx, created.ID, "updated", "updated content", "author2", []string{"tag1"})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if updated.Title != "updated" {
		t.Fatalf("title not updated: %s", updated.Title)
	}
	if updated.Content != "updated content" {
		t.Fatalf("content not updated: %s", updated.Content)
	}
	if updated.Author != "author2" {
		t.Fatalf("author not updated: %s", updated.Author)
	}
	if len(updated.Tags) != 1 {
		t.Fatalf("tags not updated: %v", updated.Tags)
	}
}

func TestPostService_UpdateNotFound(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	_, err := service.UpdatePost(ctx, "missing", "title", "content", "author", nil)
	if err != apperrors.ErrPostNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestPostService_UpdateInvalidTitle(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	created, _ := service.CreatePost(ctx, "original", "original content", "author1", "", nil)

	_, err := service.UpdatePost(ctx, created.ID, "", "content", "author", nil)
	if err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestPostService_DeletePost(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	created, _ := service.CreatePost(ctx, "title", "content", "author", "", nil)

	err := service.DeletePost(ctx, created.ID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = service.GetPost(ctx, created.ID)
	if err != apperrors.ErrPostNotFound {
		t.Fatalf("expected post to be deleted")
	}
}

func TestPostService_DeleteInvalidEmptyID(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	if err := service.DeletePost(ctx, ""); err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestPostService_DeleteNotFound(t *testing.T) {
	repo := memory.NewPostRepository()
	service := NewPostService(repo)

	ctx := context.Background()
	if err := service.DeletePost(ctx, "nonexistent"); err != apperrors.ErrPostNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
