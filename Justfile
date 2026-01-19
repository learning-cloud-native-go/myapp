export GOEXPERIMENT := "jsonv2"

server_port := "8080"
db_port := "5432"
db_host := "localhost"

# List available commands
help:
    @just --list --unsorted --list-heading $'🚀MYAPP\n'

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
    gofumpt -l -w . 2>&1 | tee outfile && test -z "$(cat outfile)" && rm outfile
    go vet ./...
    staticcheck ./...
    govulncheck ./...

# Run tests
test:
    go test -v -race ./...

# Generate openapi v3 specification using swag v2
gen-openapi:
    swag init -g cmd/app/main.go -o . -ot yaml --v3.1 --parseDependency && mv swagger.yaml openapi-v3.yml

# Generate gorm repositories using gorm cli
gen-gorm-repos:
    gorm gen -i ./app/book/repository.go -o ./app/book/bookrepo