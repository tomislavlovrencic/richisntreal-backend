package bootstrap

import (
	"log"

	"github.com/go-chi/chi/v5"
	"richisntreal-backend/cmd/config"
	"richisntreal-backend/internal/api/handlers"
	"richisntreal-backend/internal/core/services"
	"richisntreal-backend/internal/infrastructure/mysql"
)

// NewRouter wires up config, DB, services, handlers, and returns your mux.
func NewRouter() *chi.Mux {
	// 1) load config
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := config.Get()

	// 2) init MySQL
	mysqlClient, err := mysql.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}

	// 3) init repo / service / handler
	userRepo := mysql.NewUserRepository(mysqlClient.DB)
	userSvc := services.NewUserService(userRepo, cfg.App.JWTSecret)
	userHandler := handlers.NewUserHandler(userSvc)

	// 4) mount routes
	r := chi.NewRouter()
	userHandler.RegisterRoutes(r)

	return r
}
