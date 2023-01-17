# Learning Cloud Native Go - myapp
Cloud Native Application Development is a one way of speeding up building web applications, using micro-services, containers and orchestration tools.

As the first step, this repository shows **How to build a Dockerized RESTful API application using Go**. 

## Points to Highlight
- Usage of Docker and Docker Compose.
- Usage of Golang and MySQL Alpine images.
- Usage of Docker Multistage builds.
- [Health API for K8s liveness & readiness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/).
- Usage of [Goose](https://github.com/pressly/goose) for Migrations.
- Usage of [GORM](https://gorm.io/) as the ORM.
- Usage of [Chi](https://github.com/go-chi/chi) as the Router.
- Usage of [Zerolog](https://github.com/rs/zerolog) as the Logger.
- Usage of [Validator.v10](https://github.com/go-playground/validator) as the Form Validator.

### Endpoints
![endpoints](doc/assets/endpoints.png)

### Docker Image Sizes
- DB: 230MB
- App
    - Development environment: 667MB
    - Production environment: 21MB

> 💡 Building Docker image for production
> `docker build -f prod.Dockerfile . -t myapp_app`

## Design Decisions & Project Folder Structure
- Store executable packages inside the `cmd` folder.
- Store database migrations inside the `migrations` folder.
- Store main application code inside the `app` folder.
- Store reusable packages like configs, utils in separate folders. This will be helpful if you are adding more executable applications to support web front-ends, [publish/subscribe systems](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern), [document stores](https://en.wikipedia.org/wiki/Document-oriented_database) and etc.

```
.
├── docker-compose.yml
├── Dockerfile
├── prod.Dockerfile
│
├── cmd
│  ├── api
│  │  └── main.go
│  └── migrate
│     └── main.go
│
├── migrations
│  └── 20190805170000_create_books_table.sql
│
├── api
│  ├── resource
│  │  ├── health
│  │  │  └── handler.go
│  │  ├── book
│  │  │  ├── app.go
│  │  │  ├── handler.go
│  │  │  ├── model.go
│  │  │  └── repository.go
│  │  └── error
│  │     └── handler.go
│  │
│  ├── router
│  │  ├── router.go
│  │  └── middleware
│  │     ├── content_type_json.go
│  │     └── content_type_json_test.go
│  │
│  └── requestlog
│     ├── handler.go
│     └── log_entry.go
│
├── config
│  └── config.go
│
├── util
│  ├── logger
│  │  ├── logger.go
│  │  └── logger_test.go
│  └── validator
│     └── validator.go
│     └── validator_test.go
│
├── k8s
│  ├── app-configmap.yaml
│  ├── app-secret.yaml
│  ├── app-deployment.yaml
│  └── app-service.yaml
│
├── go.mod
└── go.sum
```

### Form Validation
![Form validation](doc/assets/form_validation.png)

### Logs
![Logs in app init](doc/assets/logs_app_init.png)
![Logs in crud](doc/assets/logs_crud.png)
