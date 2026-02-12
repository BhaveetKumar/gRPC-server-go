package integration

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/BhaveetKumar/gRPC-server-go/internal/handler"
	"github.com/BhaveetKumar/gRPC-server-go/internal/logger"
	"github.com/BhaveetKumar/gRPC-server-go/internal/repository/memory"
	"github.com/BhaveetKumar/gRPC-server-go/internal/service"
	blogv1 "github.com/BhaveetKumar/gRPC-server-go/proto/blog/v1"
	"google.golang.org/grpc"
)

func startTestServer(t *testing.T) (blogv1.BlogServiceClient, func()) {
	t.Helper()

	baseLogger := logger.New()
	repo := memory.NewPostRepository()
	postService := service.NewPostService(repo)
	blogHandler := handler.NewBlogHandler(postService, baseLogger)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.UnaryServerInterceptor(baseLogger)),
	)
	blogv1.RegisterBlogServiceServer(grpcServer, blogHandler)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		grpcServer.Stop()
		lis.Close()
		t.Fatalf("dial: %v", err)
	}

	client := blogv1.NewBlogServiceClient(conn)

	cleanup := func() {
		conn.Close()
		grpcServer.Stop()
		lis.Close()
	}

	return client, cleanup
}

func TestBlogService_EndToEnd(t *testing.T) {
	client, cleanup := startTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := client.CreatePost(ctx, &blogv1.CreatePostRequest{
		Title:   "title",
		Content: "content",
		Author:  "author",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	id := created.GetPost().GetPostId()

	got, err := client.GetPost(ctx, &blogv1.GetPostRequest{PostId: id})
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if got.GetPost().GetTitle() != "title" {
		t.Fatalf("unexpected title: %s", got.GetPost().GetTitle())
	}

	_, err = client.UpdatePost(ctx, &blogv1.UpdatePostRequest{
		PostId:  id,
		Title:   "new",
		Content: "content",
		Author:  "author",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}

	_, err = client.DeletePost(ctx, &blogv1.DeletePostRequest{PostId: id})
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
}
