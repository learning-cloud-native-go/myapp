FROM golang:1.26rc2-alpine
WORKDIR /myapp

ENV GOEXPERIMENT=jsonv2

RUN apk add --no-cache gcc musl-dev tzdata

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app ./cmd/app \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate

CMD ["/myapp/bin/app"]
EXPOSE 8080