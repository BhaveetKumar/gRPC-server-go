#!/usr/bin/env bash
set -euo pipefail

dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$dir"

echo "=== Starting gRPC Server ==="
go run ./cmd/server &
SERVER_PID=$!

sleep 2

echo ""
echo "=== Testing CRUD Operations ==="
echo ""

echo "1. CREATE Post"
CREATE_OUTPUT=$(go run ./cmd/client create -title "My First Blog Post" -content "This is the content of my first post" -author "John Doe" -tags "golang,grpc,testing" 2>&1)
echo "$CREATE_OUTPUT"

POST_ID=$(echo "$CREATE_OUTPUT" | grep -o 'post_id:"[^"]*"' | sed 's/post_id:"\(.*\)"/\1/')
echo "Created Post ID: $POST_ID"
echo ""

sleep 1

echo "2. GET Post by ID"
go run ./cmd/client get -id "$POST_ID"
echo ""

sleep 1

echo "3. UPDATE Post"
go run ./cmd/client update -id "$POST_ID" -title "Updated Blog Post" -content "This content has been updated" -author "Jane Smith" -tags "updated,modified"
echo ""

sleep 1

echo "4. GET Updated Post"
go run ./cmd/client get -id "$POST_ID"
echo ""

sleep 1

echo "5. DELETE Post"
go run ./cmd/client delete -id "$POST_ID"
echo ""

sleep 1

echo "6. Verify Deletion (should fail)"
go run ./cmd/client get -id "$POST_ID" 2>&1 || echo "Post successfully deleted (not found)"
echo ""

echo "=== All CRUD Operations Completed ==="
echo ""
echo "=== Interactive Client Mode ==="
echo "Server is running on localhost:50051"
echo ""
echo "Available commands:"
echo "  1 - Create Post"
echo "  2 - Get Post"
echo "  3 - Update Post"
echo "  4 - Delete Post"
echo "  q - Quit and stop server"
echo ""

while true; do
    echo -n "Enter command (1-4 or q): "
    read -r cmd
    
    case $cmd in
        1)
            echo -n "Title: "
            read -r title
            echo -n "Content: "
            read -r content
            echo -n "Author: "
            read -r author
            echo -n "Tags (comma-separated, optional): "
            read -r tags
            
            if [ -n "$tags" ]; then
                go run ./cmd/client create -title "$title" -content "$content" -author "$author" -tags "$tags"
            else
                go run ./cmd/client create -title "$title" -content "$content" -author "$author"
            fi
            echo ""
            ;;
        2)
            echo -n "Post ID: "
            read -r post_id
            go run ./cmd/client get -id "$post_id"
            echo ""
            ;;
        3)
            echo -n "Post ID: "
            read -r post_id
            echo -n "New Title: "
            read -r title
            echo -n "New Content: "
            read -r content
            echo -n "New Author: "
            read -r author
            echo -n "New Tags (comma-separated, optional): "
            read -r tags
            
            if [ -n "$tags" ]; then
                go run ./cmd/client update -id "$post_id" -title "$title" -content "$content" -author "$author" -tags "$tags"
            else
                go run ./cmd/client update -id "$post_id" -title "$title" -content "$content" -author "$author"
            fi
            echo ""
            ;;
        4)
            echo -n "Post ID: "
            read -r post_id
            go run ./cmd/client delete -id "$post_id"
            echo ""
            ;;
        q|Q)
            echo "Stopping server..."
            kill $SERVER_PID 2>/dev/null || true
            echo "Server stopped. Goodbye!"
            exit 0
            ;;
        *)
            echo "Invalid command. Please enter 1-4 or q"
            echo ""
            ;;
    esac
done
