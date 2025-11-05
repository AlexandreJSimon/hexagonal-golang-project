package main

import (
	"log"
	"net/http"

	_ "github.com/AlexandreJSimon/hexagonal-golang-project/docs"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/services/user_service"
	api_handlers "github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/api/handlers"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/repositories/user_repository"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Hexagonal Go API
// @version 1.0
// @description This is a sample server for a hexagonal architecture in Go.
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	userRepository := user_repository.NewMemoryRepository()

	userService := user_service.NewUserService(user_service.UserServiceInput{
		UserRepository: userRepository,
	})

	handlers := api_handlers.NewHandler(api_handlers.HandlerInput{
		UserService: userService,
	})

	mux := http.NewServeMux()

	apiV1Handler := http.StripPrefix("/api/v1", func() http.Handler {
		mux := http.NewServeMux()

		mux.HandleFunc("POST /users", handlers.CreateUser)
		mux.HandleFunc("GET /users", handlers.ListUsers)

		return mux
	}())

	mux.Handle("/api/v1/", apiV1Handler)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
