FROM golang:1.19-alpine
WORKDIR /myapp

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate

CMD ["/myapp/bin/api"]
EXPOSE 8080