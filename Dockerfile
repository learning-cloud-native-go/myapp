FROM golang:1.27-alpine
WORKDIR /myapp

RUN apk add --no-cache gcc musl-dev tzdata

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app \
    && go build -tags=embed -ldflags '-w -s' -a -o ./bin/migration ./cmd/migration

CMD ["/myapp/bin/app"]
EXPOSE 8080