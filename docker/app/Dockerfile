# Build environment
# -----------------
FROM golang:1.13-alpine as build-env
WORKDIR /myapp

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate


# Deployment environment
# ----------------------
FROM alpine
RUN apk update && apk add --no-cache bash mysql-client

COPY --from=build-env /myapp/bin/app /myapp/
COPY --from=build-env /myapp/bin/migrate /myapp/
COPY --from=build-env /myapp/migrations /myapp/migrations

COPY --from=build-env /myapp/docker/app/bin /usr/local/bin/myapp/
RUN chmod +x /usr/local/bin/myapp/*

EXPOSE 8080
CMD ["/myapp/app"]