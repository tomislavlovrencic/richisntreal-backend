package bootstrap

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"richisntreal-backend/cmd/config"
	"richisntreal-backend/internal/api/handlers"
	"richisntreal-backend/internal/core/services"
	mysql "richisntreal-backend/internal/infrastructure/mysql"
)

func NewRouter() *chi.Mux {
	// 1) Load config
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := config.Get()

	// 2) Run migrations
	runMigrations(cfg.MySQL)

	// 3) Init MySQL client
	mysqlClient, err := mysql.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}

	// 4) Wire services & handlers
	userRepo := mysql.NewUserRepository(mysqlClient.DB)
	userSvc := services.NewUserService(userRepo, cfg.App.JWTSecret)
	userHandler := handlers.NewUserHandler(userSvc)

	cartRepo := mysql.NewCartRepository(mysqlClient.DB)
	cartService := services.NewCartService(cartRepo)
	cartHandler := handlers.NewCartHandler(cartService)

	prodRepo := mysql.NewProductRepository(mysqlClient.DB)
	prodService := services.NewProductService(prodRepo)
	prodHandler := handlers.NewProductHandler(prodService)

	// 5) Mount routes
	r := chi.NewRouter()
	userHandler.RegisterRoutes(r)
	cartHandler.RegisterRoutes(r)
	prodHandler.RegisterRoutes(r)
	return r
}

// runMigrations applies all “.up.sql” scripts in migrations/ against your DB.
func runMigrations(mysqlCfg config.MySQL) {
	// source://directory and database://dsn
	sourceURL := "file://internal/infrastructure/mysql/migrations"
	dbURL := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		mysqlCfg.Username,
		mysqlCfg.Password,
		mysqlCfg.Host,
		mysqlCfg.Port,
		mysqlCfg.Database,
	)

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatalf("migrations: failed to initialize: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrations: failed to run up: %v", err)
	}
	log.Println("migrations: applied all available migrations")
}
