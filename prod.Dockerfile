# Build environment
# -----------------
FROM --platform=$BUILDPLATFORM golang:1.26rc2-alpine as build-env
WORKDIR /myapp

ENV GOEXPERIMENT=jsonv2

RUN apk add --no-cache tzdata ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

ARG TARGETOS TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags '-w -s' -o ./bin/app ./cmd/app \
    && go build -ldflags '-w -s' -o ./bin/migrate ./cmd/migrate


# Deployment environment
# ----------------------
FROM gcr.io/distroless/static-debian12

ENV TZ=Asia/Singapore

COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build-env /myapp/bin/app /myapp/
COPY --from=build-env /myapp/bin/migrate /myapp/

USER 65532:65532

CMD ["/myapp/app"]