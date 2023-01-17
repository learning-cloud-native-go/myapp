# Build environment
# -----------------
FROM golang:1.19-alpine as build-env
WORKDIR /myapp

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate


# Deployment environment
# ----------------------
FROM alpine
RUN apk update

COPY --from=build-env /myapp/bin/api /myapp/
COPY --from=build-env /myapp/bin/migrate /myapp/
COPY --from=build-env /myapp/migrations /myapp/migrations

EXPOSE 8080
CMD ["/myapp/app"]