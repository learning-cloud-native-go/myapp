export GOEXPERIMENT := "jsonv2"

server_port := "8080"
db_port := "5432"
db_host := "localhost"

# List available commands
help:
    @just --list --unsorted --list-heading $'🚀MYAPP\n'

# Install development tools
install:
    go install tool

# Run a specific cmd (defaults to app)
go-run cmd="app":
    @export $(grep -v '^#' .env | xargs) && \
    DB_PORT={{db_port}} DB_HOST={{db_host}} SERVER_PORT={{server_port}} \
    go run ./cmd/{{cmd}}

# Run database migrations (defaults to up)
go-run-migrate cmd="up":
    @export $(grep -v '^#' .env | xargs) && \
    DB_PORT={{db_port}} DB_HOST={{db_host}} SERVER_PORT={{server_port}} \
    go run ./cmd/migrate {{cmd}}

# Run docker compose build
build:
    @docker compose build

# Run docker compose up
up cmd="":
    @docker compose up {{cmd}}

# Run docker compose down
down:
    @docker compose down

# Run lints using gofumpt, go vet, staticcheck and govulncheck
lint:
    go tool gofumpt -d -e .
    go vet ./...
    go tool staticcheck ./...
    go tool govulncheck ./...

# Run tests
test:
    go test -v -race ./...

# Run go generate for all packages
gen:
    go generate ./...

# Generate openapi v3 specification using swag v2
gen-openapi:
    go tool swag init -g cmd/app/main.go -o . -ot yaml --v3.1 --parseDependency && mv swagger.yaml openapi-v3.yaml

# Generate gorm repositories using gorm cli
gen-gorm-repos:
    go tool gorm gen -i ./app/book/repository.go -o ./app/book/bookrepo