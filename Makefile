PROTO_FILES=proto/blog/v1/blog.proto

proto:
	PATH="$$(go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(PROTO_FILES)

build:
	go build ./cmd/server
	go build ./cmd/client

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client

test:
	go test ./...

clean:
	rm -f ./cmd/server/server ./cmd/client/client
