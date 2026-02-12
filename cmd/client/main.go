package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BhaveetKumar/gRPC-server-go/internal/config"
	blogv1 "github.com/BhaveetKumar/gRPC-server-go/proto/blog/v1"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: client <command> [flags]")
		log.Println("commands: create, get, update, delete")
		os.Exit(1)
	}

	command := os.Args[1]

	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Client.TimeoutSeconds)*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.Client.ServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := blogv1.NewBlogServiceClient(conn)

	switch command {
	case "create":
		runCreate(ctx, client, os.Args[2:])
	case "get":
		runGet(ctx, client, os.Args[2:])
	case "update":
		runUpdate(ctx, client, os.Args[2:])
	case "delete":
		runDelete(ctx, client, os.Args[2:])
	default:
		log.Fatalf("unknown command: %s", command)
	}
}

func runCreate(ctx context.Context, client blogv1.BlogServiceClient, args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	title := fs.String("title", "", "post title")
	content := fs.String("content", "", "post content")
	author := fs.String("author", "", "post author")
	date := fs.String("date", "", "publication date")
	tags := fs.String("tags", "", "comma separated tags")
	_ = fs.Parse(args)

	req := &blogv1.CreatePostRequest{
		Title:           *title,
		Content:         *content,
		Author:          *author,
		PublicationDate: *date,
		Tags:            splitTags(*tags),
	}

	resp, err := client.CreatePost(ctx, req)
	if err != nil {
		log.Fatalf("create failed: %v", err)
	}

	fmt.Printf("created post: %+v\n", resp.GetPost())
}

func runGet(ctx context.Context, client blogv1.BlogServiceClient, args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	id := fs.String("id", "", "post id")
	_ = fs.Parse(args)

	req := &blogv1.GetPostRequest{PostId: *id}
	resp, err := client.GetPost(ctx, req)
	if err != nil {
		log.Fatalf("get failed: %v", err)
	}

	fmt.Printf("post: %+v\n", resp.GetPost())
}

func runUpdate(ctx context.Context, client blogv1.BlogServiceClient, args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	id := fs.String("id", "", "post id")
	title := fs.String("title", "", "post title")
	content := fs.String("content", "", "post content")
	author := fs.String("author", "", "post author")
	tags := fs.String("tags", "", "comma separated tags")
	_ = fs.Parse(args)

	req := &blogv1.UpdatePostRequest{
		PostId:  *id,
		Title:   *title,
		Content: *content,
		Author:  *author,
		Tags:    splitTags(*tags),
	}

	resp, err := client.UpdatePost(ctx, req)
	if err != nil {
		log.Fatalf("update failed: %v", err)
	}

	fmt.Printf("updated post: %+v\n", resp.GetPost())
}

func runDelete(ctx context.Context, client blogv1.BlogServiceClient, args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.String("id", "", "post id")
	_ = fs.Parse(args)

	req := &blogv1.DeletePostRequest{PostId: *id}
	resp, err := client.DeletePost(ctx, req)
	if err != nil {
		log.Fatalf("delete failed: %v", err)
	}

	fmt.Printf("delete success: %v\n", resp.GetSuccess())
}

func splitTags(raw string) []string {
	if raw == "" {
		return nil
	}

	var result []string
	current := ""
	for i := 0; i < len(raw); i++ {
		if raw[i] == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
			continue
		}
		current += string(raw[i])
	}
	if current != "" {
		result = append(result, current)
	}

	return result
}
