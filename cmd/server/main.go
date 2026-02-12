package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/BhaveetKumar/gRPC-server-go/internal/config"
	"github.com/BhaveetKumar/gRPC-server-go/internal/handler"
	"github.com/BhaveetKumar/gRPC-server-go/internal/logger"
	"github.com/BhaveetKumar/gRPC-server-go/internal/repository/memory"
	"github.com/BhaveetKumar/gRPC-server-go/internal/service"
	blogv1 "github.com/BhaveetKumar/gRPC-server-go/proto/blog/v1"
	"google.golang.org/grpc"
)

func main() {
	log.Println("starting gRPC blog server")

	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	baseLogger := logger.NewWithConfig(cfg.Log.EnableRequestID)
	repo := memory.NewPostRepository()
	postService := service.NewPostService(repo)
	blogHandler := handler.NewBlogHandler(postService, baseLogger)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.UnaryServerInterceptor(baseLogger)),
	)

	blogv1.RegisterBlogServiceServer(grpcServer, blogHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}

	go func() {
		log.Printf("gRPC server listening on %s", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("shutting down gRPC server")
	grpcServer.GracefulStop()
}
