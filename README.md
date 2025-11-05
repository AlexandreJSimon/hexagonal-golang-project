# Go Hexagonal Architecture Base Project

This project provides a **base structure** for quickly setting up a new **Go (Golang)** application following **Hexagonal Architecture** and **Domain-Driven Design (DDD)** principles.

It includes ready-to-use configurations for running the API, generating mocks, executing tests, and initializing Swagger documentation.

---

## 🧱 Project Overview

The goal of this repository is to serve as a **scalable and modular starting point** for Go applications, making it easier to maintain clean boundaries between:
- **Domain logic**
- **Application services**
- **Infrastructure adapters**

This base setup promotes testability, flexibility, and a clear separation of concerns.

---

## 🚀 Getting Started

### Prerequisites
Make sure you have the following installed:
- **Go 1.24+**
- **Swag CLI** (`go install github.com/swaggo/swag/cmd/swag@latest`)
- **Mockgen** (`go install go.uber.org/mock/mockgen@latest`)

---

## 🧩 Project Structure

```
cmd/
  └── api/
      └── main.go           # API entry point
internal/
  ├── domain/               # Entities and domain interfaces
  ├── application/          # Business logic (services, use cases)
  └── infra/                # Infrastructure (repositories, API handlers)
mocks/                      # Generated mocks for testing
```

---

## ⚙️ Makefile Commands

### Run the API
Starts the main Go application:
```bash
make api
```
Equivalent to:
```bash
go run ./cmd/api/main.go
```

---

### Generate Mocks
Creates mock implementations for interfaces using **mockgen**:
```bash
make mocks
```
This will generate:
- `mocks/user_service_mock.go`
- `mocks/user_repository_mock.go`

Manual command example:
```bash
go run go.uber.org/mock/mockgen@latest -destination=mocks/user_service_mock.go -package=mocks ./internal/application/services/user_service UserServiceProvider
go run go.uber.org/mock/mockgen@latest -destination=mocks/user_repository_mock.go -package=mocks ./internal/domain/user UserRepository
```

---

### Run Tests
Executes all unit tests with race detection and coverage:
```bash
make test
```
Equivalent to:
```bash
go test -v -race -cover ./...
```

---

### Generate Swagger Docs
Generates API documentation using **swag**:
```bash
make swagger
```
Equivalent to:
```bash
swag init -g main.go -d ./cmd/api,./internal/infra/api/handlers,./internal/infra/api/dto
```

---

## 📘 Notes

- All interfaces and services should follow **dependency inversion** principles.
- Mocks are automatically regenerated before running tests.
- Swagger configuration assumes handlers and DTOs are under `internal/infra/api/`.

---

## 🧑‍💻 Author
**Alexandre Simon**
