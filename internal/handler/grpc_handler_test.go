package handler

import (
	"context"
	"testing"

	"github.com/BhaveetKumar/gRPC-server-go/internal/logger"
	"github.com/BhaveetKumar/gRPC-server-go/internal/repository/memory"
	"github.com/BhaveetKumar/gRPC-server-go/internal/service"
	blogv1 "github.com/BhaveetKumar/gRPC-server-go/proto/blog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setupHandler() *BlogHandler {
	repo := memory.NewPostRepository()
	svc := service.NewPostService(repo)
	log := logger.New()
	return NewBlogHandler(svc, log)
}

func TestBlogHandler_CreatePost(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()
	req := &blogv1.CreatePostRequest{
		Title:   "Test Title",
		Content: "Test Content",
		Author:  "Test Author",
		Tags:    []string{"golang", "grpc"},
	}

	resp, err := handler.CreatePost(ctx, req)
	if err != nil {
		t.Fatalf("create post failed: %v", err)
	}
	if resp.GetPost() == nil {
		t.Fatal("expected post in response")
	}
	if resp.GetPost().GetPostId() == "" {
		t.Fatal("expected post id to be set")
	}
}

func TestBlogHandler_CreatePostInvalid(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()
	req := &blogv1.CreatePostRequest{Title: "", Content: "Content", Author: "Author"}

	_, err := handler.CreatePost(ctx, req)
	if err == nil {
		t.Fatal("expected error for empty title")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", err)
	}
}

func TestBlogHandler_GetPost(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()
	createReq := &blogv1.CreatePostRequest{Title: "Test", Content: "Content", Author: "Author"}
	createResp, _ := handler.CreatePost(ctx, createReq)
	postID := createResp.GetPost().GetPostId()

	getReq := &blogv1.GetPostRequest{PostId: postID}
	getResp, err := handler.GetPost(ctx, getReq)
	if err != nil {
		t.Fatalf("get post failed: %v", err)
	}
	if getResp.GetPost().GetPostId() != postID {
		t.Fatal("post id mismatch")
	}
}

func TestBlogHandler_GetPostNotFound(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()
	req := &blogv1.GetPostRequest{PostId: "nonexistent"}

	_, err := handler.GetPost(ctx, req)
	if err == nil {
		t.Fatal("expected error for nonexistent post")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", err)
	}
}

func TestBlogHandler_UpdatePost(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()
	createReq := &blogv1.CreatePostRequest{Title: "Original", Content: "Content", Author: "Author"}
	createResp, _ := handler.CreatePost(ctx, createReq)
	postID := createResp.GetPost().GetPostId()

	updateReq := &blogv1.UpdatePostRequest{PostId: postID, Title: "Updated", Content: "Updated Content", Author: "Updated Author"}
	updateResp, err := handler.UpdatePost(ctx, updateReq)
	if err != nil {
		t.Fatalf("update post failed: %v", err)
	}
	if updateResp.GetPost().GetTitle() != "Updated" {
		t.Fatal("title not updated")
	}
}

func TestBlogHandler_UpdatePostNotFound(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()
	req := &blogv1.UpdatePostRequest{PostId: "nonexistent", Title: "Title", Content: "Content", Author: "Author"}

	_, err := handler.UpdatePost(ctx, req)
	if err == nil {
		t.Fatal("expected error for nonexistent post")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", err)
	}
}

func TestBlogHandler_DeletePost(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()
	createReq := &blogv1.CreatePostRequest{Title: "Test", Content: "Content", Author: "Author"}
	createResp, _ := handler.CreatePost(ctx, createReq)
	postID := createResp.GetPost().GetPostId()

	deleteReq := &blogv1.DeletePostRequest{PostId: postID}
	deleteResp, err := handler.DeletePost(ctx, deleteReq)
	if err != nil {
		t.Fatalf("delete post failed: %v", err)
	}
	if !deleteResp.GetSuccess() {
		t.Fatal("expected success")
	}

	_, err = handler.GetPost(ctx, &blogv1.GetPostRequest{PostId: postID})
	if err == nil {
		t.Fatal("expected error for deleted post")
	}
}

func TestBlogHandler_DeletePostNotFound(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()
	req := &blogv1.DeletePostRequest{PostId: "nonexistent"}

	_, err := handler.DeletePost(ctx, req)
	if err == nil {
		t.Fatal("expected error for nonexistent post")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.NotFound {
		t.Fatalf("expected NotFound, got %v", err)
	}
}

func TestBlogHandler_FullCRUDFlow(t *testing.T) {
	handler := setupHandler()
	ctx := context.Background()

	createReq := &blogv1.CreatePostRequest{Title: "Flow Test", Content: "Testing CRUD", Author: "Tester", Tags: []string{"test"}}
	createResp, err := handler.CreatePost(ctx, createReq)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	postID := createResp.GetPost().GetPostId()

	updateReq := &blogv1.UpdatePostRequest{PostId: postID, Title: "Updated", Content: "Updated content", Author: "Updated"}
	_, err = handler.UpdatePost(ctx, updateReq)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	deleteReq := &blogv1.DeletePostRequest{PostId: postID}
	_, err = handler.DeletePost(ctx, deleteReq)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = handler.GetPost(ctx, &blogv1.GetPostRequest{PostId: postID})
	if err == nil {
		t.Fatal("get after delete should fail")
	}
}
