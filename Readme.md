# gRPC Blog Post Service

A simple gRPC-based API in Go for managing blog posts with CRUD operations.

## Features

- **gRPC API** with 4 operations: Create, Get, Update, Delete
- **In-memory storage** for blog posts
- **Input/output logging** for all requests
- **Unit tests** for all layers

## Quick Start

Run automated tests and interactive client:
```bash
./scripts/test_crud.sh
```

The script will:
1. Start the gRPC server on `localhost:50051`
2. Run automated CRUD tests
3. Enter interactive mode where you can manually test operations

## API Operations

- `CreatePost` - Create a new blog post
- `GetPost` - Retrieve a post by ID
- `UpdatePost` - Update an existing post
- `DeletePost` - Remove a post

See `proto/blog/v1/blog.proto` for the complete API definition.

## Project Structure

```
cmd/server/        - gRPC server
cmd/client/        - CLI client
internal/handler/  - gRPC request handlers
internal/service/  - Business logic
internal/repository/ - In-memory data storage
proto/blog/v1/     - Protocol buffer definitions
test/              - Unit and integration tests
```

## Configuration

Edit `.env` file to configure:
- Server host and port
- Client timeout
- Request ID logging (disabled by default)

## Testing

Run all tests:
```bash
go test ./...
```

39 tests covering repository, service, handler, and integration layers.