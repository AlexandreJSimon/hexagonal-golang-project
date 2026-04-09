package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/AlexandreJSimon/hexagonal-golang-project/docs"
	userApp "github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/user"
	env "github.com/AlexandreJSimon/hexagonal-golang-project/internal/config"
	httpx "github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http"
	user_handler "github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/handlers/user"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/http/middleware"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/repositories/user_repository"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/security"
	"github.com/AlexandreJSimon/hexagonal-golang-project/pkg/jwt"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	env := env.Load()

	tokenService, err := jwt.New(jwt.Config{SecretKey: env.JWTSecretKey})
	if err != nil {
		log.Fatalf("failed to initialize JWT: %v", err)
	}

	userRepository := user_repository.NewMemoryRepository()

	userService := userApp.NewUserService(userApp.UserServiceInput{
		UserRepository: userRepository,
		PasswordHasher: security.BcryptHasher{},
	})

	userHandler := user_handler.NewHandler(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/login", userHandler.Login(tokenService))

	httpx.Group(mux, "/api/v1", func(mux *http.ServeMux) {
		mux.HandleFunc("POST /users", userHandler.Create)
		mux.HandleFunc("GET /users", userHandler.List)
	}, middleware.Authentication(tokenService))

	handler := middleware.Chain(middleware.CORS())(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "healthy"}`))
}
