[![buymeacoffee](https://img.shields.io/badge/Buy%20me%20a%20coffee-dumindu-FFDD00?style=for-the-badge&logo=buymeacoffee&logoColor=ffffff&labelColor=333333)](https://www.buymeacoffee.com/dumindu)

# Learning Cloud Native Go - myapp

## 🔋 Batteries Included

- Use of Go linters, Docker, Docker Compose, Alpine development images, and Distroless production images.
- Use of [GORM CLI](https://gorm.io/cli/) to generate Go generics-based database repositories.
- Use of [Goose](https://github.com/pressly/goose) to build a DB migration CLI with embedded migrations.
- Use of [Validator.v10](https://github.com/go-playground/validator) to validate forms via Go generics-based, fail-fast validation middleware.
- Use of [Zerolog](https://github.com/rs/zerolog) to generate request logs and centralized Syslog logging.
- Use of [Swag.v2](https://github.com/swaggo/swag) to generate OpenAPI v3 specifications.
- Use of GitHub Actions to run linters and tests, and to build and push production images to the registry.
- Use of GitOps with ArgoCD to automate declarative environment orchestration and application lifecycle management.

| Environment    | Go 1.26 Image Size | Postgres v18 Image Size |
|----------------|--------------------|-------------------------|
| Development    | 800 MB             | 300MB                   |
| Production     | 30 MB              |                         |

## 📟 Available Commands

```just
$ just
🚀MYAPP
    help                    # List available commands
    install                 # Install development tools
    go-run cmd="app"        # Run a specific cmd (defaults to app)
    go-run-migrate cmd="up" # Run database migrations (defaults to up)
    build                   # Run docker compose build
    up cmd=""               # Run docker compose up
    down                    # Run docker compose down
    lint                    # Run lints using gofumpt, go vet, staticcheck and govulncheck
    test                    # Run tests
    gen                     # Run go generate for all packages
    gen-openapi             # Generate openapi v3 specification using swag v2
    gen-gorm-repos          # Generate gorm repositories using gorm cli
```

## 🛬 Endpoints

| Name        | HTTP Method | Route          |
|-------------|-------------|----------------|
| Health      | GET         | /livez         |
| List Books  | GET         | /v1/books      |
| Create Book | POST        | /v1/books      |
| Read Book   | GET         | /v1/books/{id} |
| Update Book | PUT         | /v1/books/{id} |
| Delete Book | DELETE      | /v1/books/{id} |

## 🗄️ Database Design

| Column Name    | Datatype  | Not Null | Primary Key |
|----------------|-----------|----------|-------------|
| id             | UUID      | ✅        | ✅           |
| title          | TEXT      | ✅        |             |
| author         | TEXT      | ✅        |             |
| published_date | DATE      | ✅        |             |
| image_url      | TEXT      |          |             |
| description    | TEXT      |          |             |
| created_at     | TIMESTAMP | ✅        |             |
| updated_at     | TIMESTAMP | ✅        |             |

## ⛔️ Form Validation

```json
{
  "errors": {
    "title": "This is a required field",
    "author": "This can only contain alphabetic and space characters",
    "published_date": "This must be a valid date",
    "image_url": "This must be a valid URL"
  }
}
```

## 📝 Request Logs and Centralized Syslog Logging

```json lines
db-1  | 2018-01-10 01:00:00.000 +08 [1] LOG:  database system is ready to accept connections
Container myapp-db-1 Healthy
app-1  | 2018/01/10 01:00:00 OK   00001_create_books_table.sql (2.21ms)
app-1  | 2018/01/10 01:00:00 goose: successfully migrated database to version: 1
app-1  |
app-1  | {"level":"info","time":"2018-01-10T02:00:00+08:00","message":"Starting server :8080"}
app-1  |
app-1  | [7.218ms] [rows:1] INSERT INTO books (id, created_at, updated_at, title, author, published_date, image_url, description) VALUES ('38ba23d1-9565-40ed-b781-aacd2f84018d', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Death Note', 'Light Yagami', '2006-10-04 00:00:00', 'https://static.wikia.nocookie.net/deathnote/images/9/94/A_Death_Note.jpg', 'A supernatural volume dropped into the human world by the Shinigami Ryuk') RETURNING *
app-1  | {"level":"info","request_id":"d5mq7oi6hkls7397s43g","received_time":"2018-01-10T03:00:00+08:00","method":"POST","url":"/v1/books","header_size":135,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":201,"resp_header_size":47,"resp_body_size":296,"latency":8.307,"time":"2018-01-10T03:00:00+08:00"}
app-1  | {"level":"info","request_id":"d5mq7oi6hkls7397s43g","id":"38ba23d1-9565-40ed-b781-aacd2f84018d","time":"2018-01-10T03:00:00+08:00","message":"new book created"}
app-1  |
app-1  | [2.541ms] [rows:1] SELECT * FROM books WHERE id = '38ba23d1-9565-40ed-b781-aacd2f84018d'
app-1  | {"level":"info","request_id":"d5mqa6a6hkls7397s44g","received_time":"2018-01-10T04:00:00+08:00","method":"GET","url":"/v1/books/38ba23d1-9565-40ed-b781-aacd2f84018d","header_size":82,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":200,"resp_header_size":47,"resp_body_size":296,"latency":2.674625,"time":"2018-01-10T04:00:00+08:00"}
app-1  |
app-1  | [3.744ms] [rows:1] UPDATE books SET updated_at=CURRENT_TIMESTAMP, title='Death Note', author='Misa Amane', published_date='2004-11-04 00:00:00', image_url='https://static.wikia.nocookie.net/deathnote/images/9/94/A_Death_Note.jpg', description='Light Yagami''s buried notebook' WHERE id = '38ba23d1-9565-40ed-b781-aacd2f84018d' RETURNING *
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
app-1  | [1.384ms] [rows:0] UPDATE books SET updated_at=CURRENT_TIMESTAMP, title='Death Note', author='Misa Amane', published_date='2004-11-04 00:00:00', image_url='https://static.wikia.nocookie.net/deathnote/images/9/94/A_Death_Note.jpg', description='Light Yagami''s buried notebook' WHERE id = '38ba23d1-9565-40ed-b781-aacd2f84018d' RETURNING *
app-1  | {"level":"info","request_id":"d5mqjmhqvtmc73foh3dg","received_time":"2018-02-02T08:00:00:00","method":"PUT","url":"/v1/books/38ba23d1-9565-40ed-b781-aacd2f84018d","header_size":135,"body_size":0,"agent":"yaak","referer":"","proto":"HTTP/1.1","remote_ip":"192.168.65.1","server_ip":"172.19.0.3","status":404,"resp_header_size":47,"resp_body_size":0,"latency":1.576,"time":"2018-02-02T08:00:00+08:00"}

// 💯 Real logs collected locally but with few rearrangements to make it easier to read.
```

## 🗂️ Project Folder Structure

```shell
├── compose.yml
├── Dockerfile
│
├── openapi-v3.yaml
│
├── app
│   ├── book
│   │   ├── bookrepo      # 💡generated with gorm-cli via the interface in book/repository.go
│   │   │   └── repository.go
│   │   ├── form_util.go
│   │   ├── handler.go
│   │   └── repository.go
│   └── router
│       └── router.go
├── form    # 💡Form validation middleware rely on this and pkg folder only
│   └── book.go
├── model
│   └── book.go
│
├── config
│   └── config.go
│
├── cmd   # 💡Entrypoint for app and migrate executables
│   ├── app
│   │   └── main.go
│   └── migrate
│       ├── main.go
│       └── migrations
│           └── 00001_create_books_table.sql
│
└── pkg (middleware, logger, validator, ctxutil, paramsutil, errors)
```

## 🏗️ ArgoCD and Kustomize

ArgoCD and `Kustomize` based cloud native IaC & GitOps setup.

> 💡 Consider moving to a Hub-and-Spoke architecture for Argo CD, combined with a separate repository strategy.

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

> 💡 Sample Kind Dev Cluster
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