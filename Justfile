db_host := "localhost"

# List available commands
help:
    @just --list --unsorted --list-heading $'MYAPP\n'

# Install development tools
install:
    go install tool

# Run server app
app:
    @export $(grep -v '^#' .env | xargs) && \
    DB_HOST={{ db_host }} \
    go run ./cmd/app

# Run DB migration CLI (defaults to up)
migration *cmd="up":
    @export $(grep -v '^#' .env | xargs) && \
    DB_HOST={{ db_host }} \
    go run ./cmd/migration -dir={{ justfile_directory() }}/cmd/migration/migrations {{ cmd }}

# Run docker compose build
build:
    @docker compose build

# Run docker compose up
up cmd="":
    @docker compose up {{ cmd }}

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

# Generate openapi.yaml
apidoc:
    go tool swag init -g cmd/app/main.go -o . -ot yaml --v3.1 --parseDependency && mv swagger.yaml openapi.yaml

# Generate gorm repositories using gorm cli
repos:
    go tool gorm gen -i ./app/book/repository.go -o ./app/book/bookrepo

# Build production distroless image
build-for-prod:
    docker build -f prod.Dockerfile . -t myapp-app