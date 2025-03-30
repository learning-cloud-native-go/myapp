[![buymeacoffee](https://img.shields.io/badge/Buy%20me%20a%20coffee-dumindu-FFDD00?style=for-the-badge&logo=buymeacoffee&logoColor=ffffff&labelColor=333333)](https://www.buymeacoffee.com/dumindu)


# Learning Cloud Native Go - myapp
> 🌱 Cloud Native Application Development is one way of speeding up the building of web applications using microservices, containers, and orchestration tools.

This repository shows how to build a Dockerized RESTful API application in Go for a simple bookshelf.

## 🔋 Batteries included

- The idiomatic structure based on the resource-oriented design.
- The usage of Docker, Docker compose, Alpine images, and linters on development.
- Healthcheck and CRUD API implementations with OpenAPI specifications.
- The usage of [Goose](https://github.com/pressly/goose) for the database migrations and [GORM](https://gorm.io/) as the database ORM.
- The usage of [Zerolog](https://github.com/rs/zerolog) as the centralized Syslog logger.
- The usage of [Validator.v10](https://github.com/go-playground/validator) as the form validator.
- The usage of GitHub actions to run tests and linters, generate OpenAPI specifications, and build and push production images to the Docker registry.

## 🚀 Endpoints

| Name        | HTTP Method | Route          |
|-------------|-------------|----------------|
| Health      | GET         | /livez         |
| List Books  | GET         | /v1/books      |
| Create Book | POST        | /v1/books      |
| Read Book   | GET         | /v1/books/{id} |
| Update Book | PUT         | /v1/books/{id} |
| Delete Book | DELETE      | /v1/books/{id} |

💡 [swaggo/swag](https://github.com/swaggo/swag) : `swag init -g cmd/api/main.go -o .swagger -ot yaml`

## 🗄️ Database design

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
| deleted_at     | TIMESTAMP |          |             |

## 📦 Container image sizes

- DB: 270MB
- API
    - Development environment: ~800 MB
    - Production environment: ~30 MB ; 💡`docker build -f prod.Dockerfile . -t myapp_app`

## 📁 Project structure

```shell
myapp
├── cmd
│  ├── api
│  │  └── main.go
│  └── migrate
│     └── main.go
│
├── api
│  ├── resource
│  │  ├── book
│  │  │  ├── handler.go
│  │  │  ├── model.go
│  │  │  ├── repository.go
│  │  │  └── repository_test.go
│  │  ├── common
│  │  │  └── err
│  │  │     └── err.go
│  │  └── health
│  │     └── handler.go
│  │
│  └── router
│     ├── middleware
│     │  ├── request_id.go
│     │  ├── request_id_test.go
│     │  ├── requestlog
│     │  │  ├── handler.go
│     │  │  └── log_entry.go
│     │  ├── content_type.go
│     │  └── content_type_test.go
│     └── router.go
│
├── migrations
│  └── 00001_create_books_table.sql
│
├── config
│  └── config.go
│
├── util
│  ├── logger
│  │  └── logger.go
│  └── validator
│     └── validator.go
│
├── .env
│
├── go.mod
├── go.sum
│
├── compose.yml
├── Dockerfile
│
├── prod.Dockerfile
└── k8s
   ├── app-configmap.yaml
   ├── app-secret.yaml
   ├── app-deployment.yaml
   └── app-service.yaml
```

## 📸 Form validations and logs
![Form validation](doc/assets/form_validation.png)

![Logs in app init](doc/assets/logs_app_init.png)
![Logs in crud](doc/assets/logs_crud.png)
