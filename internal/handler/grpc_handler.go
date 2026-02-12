package handler

import (
	"context"

	"github.com/BhaveetKumar/gRPC-server-go/internal/domain"
	"github.com/BhaveetKumar/gRPC-server-go/internal/errors"
	"github.com/BhaveetKumar/gRPC-server-go/internal/logger"
	"github.com/BhaveetKumar/gRPC-server-go/internal/service"
	blogv1 "github.com/BhaveetKumar/gRPC-server-go/proto/blog/v1"
)

type BlogHandler struct {
	blogv1.UnimplementedBlogServiceServer

	service service.PostService
	logger  *logger.Logger
}

func NewBlogHandler(s service.PostService, l *logger.Logger) *BlogHandler {
	return &BlogHandler{
		service: s,
		logger:  l,
	}
}

func (h *BlogHandler) CreatePost(ctx context.Context, req *blogv1.CreatePostRequest) (*blogv1.CreatePostResponse, error) {
	post, err := h.service.CreatePost(ctx, req.GetTitle(), req.GetContent(), req.GetAuthor(), req.GetPublicationDate(), req.GetTags())
	if err != nil {
		return nil, errors.ToStatus(err, h.logger)
	}

	return &blogv1.CreatePostResponse{Post: toProtoPost(post)}, nil
}

func (h *BlogHandler) GetPost(ctx context.Context, req *blogv1.GetPostRequest) (*blogv1.GetPostResponse, error) {
	post, err := h.service.GetPost(ctx, req.GetPostId())
	if err != nil {
		return nil, errors.ToStatus(err, h.logger)
	}

	return &blogv1.GetPostResponse{Post: toProtoPost(post)}, nil
}

func (h *BlogHandler) UpdatePost(ctx context.Context, req *blogv1.UpdatePostRequest) (*blogv1.UpdatePostResponse, error) {
	post, err := h.service.UpdatePost(ctx, req.GetPostId(), req.GetTitle(), req.GetContent(), req.GetAuthor(), req.GetTags())
	if err != nil {
		return nil, errors.ToStatus(err, h.logger)
	}

	return &blogv1.UpdatePostResponse{Post: toProtoPost(post)}, nil
}

func (h *BlogHandler) DeletePost(ctx context.Context, req *blogv1.DeletePostRequest) (*blogv1.DeletePostResponse, error) {
	if err := h.service.DeletePost(ctx, req.GetPostId()); err != nil {
		return nil, errors.ToStatus(err, h.logger)
	}

	return &blogv1.DeletePostResponse{Success: true}, nil
}

func toProtoPost(p *domain.Post) *blogv1.Post {
	if p == nil {
		return nil
	}

	return &blogv1.Post{
		PostId:          p.ID,
		Title:           p.Title,
		Content:         p.Content,
		Author:          p.Author,
		PublicationDate: p.PublicationDate,
		Tags:            p.Tags,
	}
}
