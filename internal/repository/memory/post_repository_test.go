package memory

import (
	"fmt"
	"sync"
	"testing"

	"github.com/BhaveetKumar/gRPC-server-go/internal/domain"
	apperrors "github.com/BhaveetKumar/gRPC-server-go/internal/errors"
)

func TestPostRepository_CreateAndGet(t *testing.T) {
	repo := NewPostRepository()

	post := &domain.Post{ID: "id1", Title: "title", Content: "content", Author: "author"}

	if err := repo.Create(post); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	loaded, err := repo.GetByID("id1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if loaded.Title != post.Title {
		t.Fatalf("unexpected title: %s", loaded.Title)
	}
	if loaded.Content != post.Content {
		t.Fatalf("unexpected content: %s", loaded.Content)
	}
	if loaded.Author != post.Author {
		t.Fatalf("unexpected author: %s", loaded.Author)
	}
}

func TestPostRepository_CreateNilPost(t *testing.T) {
	repo := NewPostRepository()

	if err := repo.Create(nil); err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input for nil post, got %v", err)
	}
}

func TestPostRepository_Duplicate(t *testing.T) {
	repo := NewPostRepository()
	post := &domain.Post{ID: "id1", Title: "title", Content: "content", Author: "author"}

	_ = repo.Create(post)
	if err := repo.Create(post); err != apperrors.ErrDuplicatePost {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestPostRepository_GetByIDEmptyString(t *testing.T) {
	repo := NewPostRepository()

	_, err := repo.GetByID("")
	if err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input for empty id, got %v", err)
	}
}

func TestPostRepository_GetByIDNotFound(t *testing.T) {
	repo := NewPostRepository()

	_, err := repo.GetByID("nonexistent")
	if err != apperrors.ErrPostNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestPostRepository_Update(t *testing.T) {
	repo := NewPostRepository()

	original := &domain.Post{ID: "id1", Title: "original", Content: "original content", Author: "author1"}
	_ = repo.Create(original)

	updated := &domain.Post{ID: "id1", Title: "updated", Content: "updated content", Author: "author2", Tags: []string{"tag1", "tag2"}}
	if err := repo.Update(updated); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	loaded, _ := repo.GetByID("id1")
	if loaded.Title != "updated" {
		t.Fatalf("title not updated: %s", loaded.Title)
	}
	if loaded.Content != "updated content" {
		t.Fatalf("content not updated: %s", loaded.Content)
	}
	if loaded.Author != "author2" {
		t.Fatalf("author not updated: %s", loaded.Author)
	}
	if len(loaded.Tags) != 2 {
		t.Fatalf("tags not updated: %v", loaded.Tags)
	}
}

func TestPostRepository_UpdateNilPost(t *testing.T) {
	repo := NewPostRepository()

	if err := repo.Update(nil); err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input for nil post, got %v", err)
	}
}

func TestPostRepository_UpdateEmptyID(t *testing.T) {
	repo := NewPostRepository()

	post := &domain.Post{ID: "", Title: "title", Content: "content", Author: "author"}
	if err := repo.Update(post); err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input for empty id, got %v", err)
	}
}

func TestPostRepository_UpdateNotFound(t *testing.T) {
	repo := NewPostRepository()

	post := &domain.Post{ID: "nonexistent", Title: "title", Content: "content", Author: "author"}
	if err := repo.Update(post); err != apperrors.ErrPostNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestPostRepository_Delete(t *testing.T) {
	repo := NewPostRepository()

	post := &domain.Post{ID: "id1", Title: "title", Content: "content", Author: "author"}
	_ = repo.Create(post)

	if err := repo.Delete("id1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := repo.GetByID("id1")
	if err != apperrors.ErrPostNotFound {
		t.Fatalf("expected post to be deleted, got %v", err)
	}
}

func TestPostRepository_DeleteEmptyID(t *testing.T) {
	repo := NewPostRepository()

	if err := repo.Delete(""); err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input for empty id, got %v", err)
	}
}

func TestPostRepository_DeleteNotFound(t *testing.T) {
	repo := NewPostRepository()

	if err := repo.Delete("missing"); err != apperrors.ErrPostNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestPostRepository_List(t *testing.T) {
	repo := NewPostRepository()

	post1 := &domain.Post{ID: "id1", Title: "title1", Content: "content1", Author: "author1"}
	post2 := &domain.Post{ID: "id2", Title: "title2", Content: "content2", Author: "author2"}
	post3 := &domain.Post{ID: "id3", Title: "title3", Content: "content3", Author: "author3"}

	_ = repo.Create(post1)
	_ = repo.Create(post2)
	_ = repo.Create(post3)

	posts, err := repo.List()
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}
}

func TestPostRepository_ListEmpty(t *testing.T) {
	repo := NewPostRepository()

	posts, err := repo.List()
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	if len(posts) != 0 {
		t.Fatalf("expected empty list, got %d posts", len(posts))
	}
}

func TestPostRepository_ConcurrentAccess(t *testing.T) {
	repo := NewPostRepository()

	const workers = 1000
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			id := fmt.Sprintf("id-%d", i)
			post := &domain.Post{ID: id, Title: "title", Content: "content", Author: "author"}

			_ = repo.Create(post)
			_, _ = repo.GetByID(id)
		}(i)
	}

	wg.Wait()

	posts, err := repo.List()
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(posts) == 0 {
		t.Fatalf("expected some posts after concurrent access")
	}
}
