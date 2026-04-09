include .env
export
MOCKGEN = mockgen

api:
	go run ./cmd/api/main.go

debug:
	go run ./cmd/debug/main.go

mocks:
	$(MOCKGEN) -destination=mocks/user_service_mock.go -package=mocks github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/handlers/user UserServiceProvider
	$(MOCKGEN) -destination=mocks/user_repository_mock.go -package=mocks github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain/user UserRepository

test: mocks
	go test -v -race -cover ./...

swagger:
	swag init -g main.go -d ./cmd/api,./internal/infra/http/handlers,./internal/infra/http/dto