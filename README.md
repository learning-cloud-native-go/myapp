[![buymeacoffee](https://img.shields.io/badge/Buy%20me%20a%20coffee-dumindu-FFDD00?style=for-the-badge&logo=buymeacoffee&logoColor=ffffff&labelColor=333333)](https://www.buymeacoffee.com/dumindu)

[![learning-cloud-native-go/myapp](https://img.shields.io/github/stars/learning-cloud-native-go/myapp?style=for-the-badge&logo=go&logoColor=ffffff&label=learning-cloud-native-go%2Fmyapp&labelColor=333333&color=00ADD8)](https://github.com/learning-cloud-native-go/myapp)

[![learning-rust.github.io](https://img.shields.io/github/stars/learning-rust/learning-rust.github.io?style=for-the-badge&logo=rust&label=learning-rust.github.io&labelColor=333333&color=F46623)](https://learning-rust.github.io)
[![dumindu/axum](https://img.shields.io/github/stars/dumindu/axum?style=for-the-badge&logo=rust&label=dumindu%2Faxum&labelColor=333333&color=F46623)](https://github.com/dumindu/axum)
[![E25DX](https://img.shields.io/github/stars/dumindu/E25DX?style=for-the-badge&logo=hugo&logoColor=ffffff&label=E25DX&labelColor=333333&color=FF4088)](https://themes.gohugo.io/themes/e25dx/)

# Learning Cloud Native Go - myapp

## In this series,

We build a production-ready containerized RESTful API server application using following packages and tools.
- Go Standard Library `net/http`: The most idiomatic way to write web API applications in Go.
- [Chi](https://github.com/go-chi/chi): The most idiomatic router with middleware and route groups support.
- [Zerolog](https://github.com/rs/zerolog): Zero Allocation JSON Logger and Faster than `slog`.
- [Goose](https://github.com/pressly/goose): The database migration CLI builder library with lease dependency.
- [Gorm CLI](https://gorm.io/cli/): Generate Go generics-based, type safe, repository functions with no runtime wrappers.
- [Swag](https://github.com/swaggo/swag) and [Validator v10](https://github.com/go-playground/validator): The most prominent OpenAPI 3.1 specification generator and validation library in Go.
- Use of GitHub Actions to run linters and tests, and to build and push production images to the registry.
- Use of GitOps with ArgoCD to automate declarative environment orchestration and application lifecycle management.

## Containerization Environment

| Environment    | Go Image Type                      | Go Image Size | Postgres Image Type | Postgres Image Size |
|----------------|------------------------------------|---------------|---------------------|---------------------|
| Development    | golang:1.27-alpine                 | ~ 800 MB      | postgres:18-alpine  | ~ 300MB             |
| Production     | distroless/static-debian13:nonroot | ~ 30 MB       |                     |                     |

## Endpoints

| Name        | HTTP Method | Route          |
|-------------|-------------|----------------|
| List Books  | GET         | /v1/books      |
| Create Book | POST        | /v1/books      |
| Read Book   | GET         | /v1/books/{id} |
| Update Book | PUT         | /v1/books/{id} |
| Delete Book | DELETE      | /v1/books/{id} |
| Health      | GET         | /livez         |

### Request (`POST`/`PUT`)
```json
{
  "title": "Harry Potter and the Deathly Hallows",
  "description": "It is the seventh and final novel in the Harry Potter series",
  "image_url": "https://upload.wikimedia.org/wikipedia/en/a/a9/Harry_Potter_and_the_Deathly_Hallows.jpg",
  "published_date": "2007-07-21",
  "status": "verified"
}
```

### Response (`GET`/`POST`/`PUT`)
```json
{
  "id": "01bbbbbb-bbbb-7bbb-8bbb-bbbbbbbbbbbb",
  "created_at": "2027-01-01T00:00:00.123456Z",
  "updated_at": "2027-01-01T00:00:00.123456Z",
  "published_date": "2007-07-21",
  "title": "Harry Potter and the Deathly Hallows",
  "description": "It is the seventh and final novel in the Harry Potter series",
  "image_url": "https://upload.wikimedia.org/wikipedia/en/a/a9/Harry_Potter_and_the_Deathly_Hallows.jpg"
  "status": "verified",
}
```
> [!note]
> The list endpoint returns an array of above response JSON.

## Form Validation

```json
{
  "errors": {
    "image_url": "Must be a valid URL",
    "status": "Must be one of pending, verified",
    "title": "This field is required",
    "published_date": "Must be a valid date"
  }
}
```

## Database Design

To keep this simple, we use only a single database table named `books`.

| Column Name    | Datatype    | Not Null | Primary Key |
|----------------|-------------|----------|-------------|
| created_at     | TIMESTAMPTZ | ✅       |             |
| updated_at     | TIMESTAMPTZ | ✅       |             |
| id             | UUID        | ✅       | ✅          |
| published_date | DATE        | ✅       |             |
| status         | SMALLINT    | ✅       |             |
| title          | TEXT        | ✅       |             |
| description    | TEXT        |          |             |
| image_url      | TEXT        |          |             |

> [!important]
> - For high-traffic systems with very large PostgreSQL tables that containing millions/billions of rows, arranging fixed-width columns by decreasing alignment requirements can reduce tuple alignment padding; potentially minimize row/ storage size. This technique is called "**Column Tetris**".
> - For this optimization, order fixed-width table columns by decreasing alignment requirements.
>   - 8-byte alignment types: `bigint`, `bigserial`, `double precision`/ `float8`, `timestamp`, `timestamptz`, `time`, `interval`
>   - 4-byte alignment types: `integer`, `serial`, `real`/ `float4`, `uuid`, `date`
>   - 2-byte alignment types: `smallint`, `smallserial`
>   - 1-byte alignment types: `boolean`
>   - Variable-width types (at last): `numeric`, `text`, `character varying`/ `varchar`, `bytea`
> - However, it's ok to follow a more readable column format, when your table schema changes frequently.

## Just commands

```just
MYAPP
    help             # List available commands
    install          # Install development tools
    app              # Run server app
    migrate cmd="up" # Run DB migration CLI (defaults to up)
    build            # Run docker compose build
    up cmd=""        # Run docker compose up
    down             # Run docker compose down
    lint             # Run lints using gofumpt, go vet, staticcheck and govulncheck
    test             # Run tests
    gen              # Run go generate for all packages
    apidoc           # Generate openapi.yaml
    repos            # Generate gorm repositories using gorm cli
```

## Sample Request Logs

```json lines
db-1  | 2018-01-10 01:00:00.000 +08 [1] LOG:  database system is ready to accept connections
Container myapp-db-1 Healthy
app-1  | 2018/01/10 01:00:00 OK   00001_create_books_table.sql (2.21ms)
app-1  | 2018/01/10 01:00:00 goose: successfully migrated database to version: 1
app-1  |
app-1  | {"level":"info","time":"2018-01-10T02:00:00+08:00","message":"Starting server :8080"}
app-1  |
app-1  | [7.218ms] [rows:1] INSERT INTO books (id, created_at, updated_at, title, published_date, image_url, description, status) VALUES ('38ba23d1-9565-40ed-b781-aacd2f84018d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Death Note', 'Light Yagami', '2006-10-04 00:00:00', 'https://static.wikia.nocookie.net/deathnote/images/9/94/A_Death_Note.jpg', 'A supernatural volume dropped into the human world by the Shinigami Ryuk', 0) RETURNING *
app-1  | {"level":"info","request_id":"d5mq7oi6hkls7397s43g","received_time":"2018-01-10T03:00:00+08:00","method":"POST","url":"/v1/books","header_size":135,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":201,"resp_header_size":47,"resp_body_size":296,"latency":8.307,"time":"2018-01-10T03:00:00+08:00"}
app-1  | {"level":"info","request_id":"d5mq7oi6hkls7397s43g","id":"38ba23d1-9565-40ed-b781-aacd2f84018d","time":"2018-01-10T03:00:00+08:00","message":"new book created"}
app-1  |
app-1  | [2.541ms] [rows:1] SELECT * FROM books WHERE id = '38ba23d1-9565-40ed-b781-aacd2f84018d'
app-1  | {"level":"info","request_id":"d5mqa6a6hkls7397s44g","received_time":"2018-01-10T04:00:00+08:00","method":"GET","url":"/v1/books/38ba23d1-9565-40ed-b781-aacd2f84018d","header_size":82,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":200,"resp_header_size":47,"resp_body_size":296,"latency":2.674625,"time":"2018-01-10T04:00:00+08:00"}
app-1  |
app-1  | [3.744ms] [rows:1] UPDATE books SET updated_at=CURRENT_TIMESTAMP, title='Death Note', published_date='2004-11-04', image_url='https://static.wikia.nocookie.net/deathnote/images/9/94/A_Death_Note.jpg', description='Light Yagami''s buried notebook', status=1 WHERE id = '38ba23d1-9565-40ed-b781-aacd2f84018d' RETURNING *
app-1  | {"level":"info","request_id":"d5mqesa6hkls7397s45g","id":"38ba23d1-9565-40ed-b781-aacd2f84018d","time":"2018-01-10T05:00:00+08:00","message":"book updated"}
app-1  | {"level":"info","request_id":"d5mqesa6hkls7397s45g","received_time":"2018-01-10T05:00:00+08:00","method":"PUT","url":"/v1/books/38ba23d1-9565-40ed-b781-aacd2f84018d","header_size":135,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":200,"resp_header_size":47,"resp_body_size":252,"latency":4.018875,"time":"2018-01-10T05:00:00+08:00"}
app-1  |
app-1  | [3.035ms] [rows:1] DELETE FROM books WHERE id = '38ba23d1-9565-40ed-b781-aacd2f84018d' RETURNING true
app-1  | {"level":"info","request_id":"d5mqfgi6hkls7397s460","received_time":"2018-01-10T06:00:00+08:00","method":"DELETE","url":"/v1/books/38ba23d1-9565-40ed-b781-aacd2f84018d","header_size":82,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":200,"resp_header_size":47,"resp_body_size":0,"latency":3.265,"time":"2018-01-10T06:00:00+08:00"}
app-1  | {"level":"info","request_id":"d5mqfgi6hkls7397s460","id":"38ba23d1-9565-40ed-b781-aacd2f84018d","time":"2018-01-10T06:00:00+08:00","message":"book deleted"}
app-1  |
app-1  | [2.573ms] [rows:1] SELECT * FROM books LIMIT 10 OFFSET 0
app-1  | {"level":"info","request_id":"d5mq9gi6hkls7397s440","received_time":"2018-01-10T07:00:00+08:00","method":"GET","url":"/v1/books","header_size":82,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":200,"resp_header_size":47,"resp_body_size":298,"latency":2.926916,"time":"2018-01-10T07:00:00+08:00"}
app-1  |
app-1  |
app-1  | [1.661ms] [rows:0] DELETE FROM books WHERE id = '38ba23d1-9565-40ed-b781-aacd2f84018d' RETURNING true
app-1  | {"level":"info","request_id":"d5mrj8ppsdvs73dkfct0","received_time":"2018-02-01T01:00:00+08:00","method":"DELETE","url":"/v1/books/38ba23d1-9565-40ed-b781-aacd2f84018d","header_size":82,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":404,"resp_header_size":47,"resp_body_size":0,"latency":1.80125,"time":"2018-02-01T01:00:00+08:00"}
app-1  | [1.384ms] [rows:0] UPDATE books SET updated_at=CURRENT_TIMESTAMP, title='Death Note', published_date='2004-11-04 00:00:00', image_url='https://static.wikia.nocookie.net/deathnote/images/9/94/A_Death_Note.jpg', description='Light Yagami''s buried notebook', status=1 WHERE id = '38ba23d1-9565-40ed-b781-aacd2f84018d' RETURNING *
app-1  | {"level":"info","request_id":"d5mqjmhqvtmc73foh3dg","received_time":"2018-02-02T08:00:00:00","method":"PUT","url":"/v1/books/38ba23d1-9565-40ed-b781-aacd2f84018d","header_size":135,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":404,"resp_header_size":47,"resp_body_size":0,"latency":1.576,"time":"2018-02-02T08:00:00+08:00"}

// 💯 Real logs collected locally but with few rearrangements to make it easier to read.
```

## Folder Structure

### API Server

```shell
├── cmd
│   ├── app
│   │   └── main.go
│   └── migrate
│       ├── main.go
│       └── migrations
│           └── 00001_create_books_table.sql
│
├── form  # 💡Form validation middleware rely on this and pkg folder only
│   └── book.go
│
├── app
│   ├── book
│   │   ├── bookrepo  # 💡generated with gorm-cli via the interface in book/repository.go
│   │   │   └── repository.go
│   │   ├── form_util.go
│   │   ├── handler.go
│   │   └── repository.go
│   └── router
│       └── router.go
│
├── model
│   ├── book.go
│   ├── book_status.go
│   └── date.go
│
├── config
│   └── config.go
│
├── pkg (middleware, validator, ctxutil, paramsutil, errors)
│
├── openapi.yaml
├── compose.yml
├── Dockerfile
└── prod.Dockerfile
```

### ArgoCD and Kustomize

```shell
└── k8s
    │
    ├── bootstrap
    │   ├── argocd
    │   └── argocd-config
    │       ├── clusters
    │       ├── projects
    │       └── applications
    │
    ├── platform
    │   ├── metrics-server
    │   ├── gateway-api
    │   ├── istio-ambient
    │   └── cloudnative-pg
    │
    ├── components
    │   └── myapp-db
    ├── services
    │   ├── base
    │   │   └── myapp
    │   └── overlays
    │       ├── dev
    │       ├── prod
    │       └── stage
    │
    └── gateways
```

> [!tip]
> Sample Kind Dev Cluster
> ```shell
> kind create cluster --name dev
> kubectl apply -k k8s/bootstrap/argocd
> kubectl apply -k k8s/bootstrap/argocd-config
> kubectl apply -k k8s/platform/istio-ambient
> kubectl apply -k k8s/platform/gateway-api
> kubectl apply -k k8s/gateways
> kubectl port-forward svc/shared-gateway-dev-istio -n istio-ingress 8081:8081 # 💡 Shared Dev Gateway
> curl -X GET 'localhost:8081/myapp/v1/books' --header 'Accept: application/json'
> 
> kubectl port-forward svc/argocd-server -n argocd 8080:443 # 💡 ArgoCD Dashboard(admin/password)
> ```