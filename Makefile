api:
	go run ./cmd/api/main.go

mocks:
	go run go.uber.org/mock/mockgen@latest -destination=mocks/user_service_mock.go -package=mocks ./internal/application/services/user_service UserServiceProvider
	go run go.uber.org/mock/mockgen@latest -destination=mocks/user_repository_mock.go -package=mocks ./internal/application/port UserRepository

test: mocks
	go test -v -race -cover ./...

swagger:
	swag init -g main.go -d ./cmd/api,./internal/infra/api/handlers,./internal/infra/api/dto