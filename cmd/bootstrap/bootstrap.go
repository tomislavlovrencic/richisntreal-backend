package bootstrap

import (
	"errors"
	"fmt"
	"github.com/go-chi/cors"
	"log"
	"richisntreal-backend/internal/api/auth"
	"richisntreal-backend/internal/api/routes"

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
	userSvc := services.NewUserService(userRepo, cfg.JWT.Secret)
	userHandler := handlers.NewUserHandler(userSvc)

	cartRepo := mysql.NewCartRepository(mysqlClient.DB)
	cartService := services.NewCartService(cartRepo)
	cartHandler := handlers.NewCartHandler(cartService)

	prodRepo := mysql.NewProductRepository(mysqlClient.DB)
	prodService := services.NewProductService(prodRepo)
	prodHandler := handlers.NewProductHandler(prodService)

	orderRepo := mysql.NewOrderRepository(mysqlClient.DB)
	orderService := services.NewOrderService(orderRepo, cartRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	payRepo := mysql.NewPaymentRepository(mysqlClient.DB)
	paySvc := services.NewPaymentService(payRepo, cfg.Stripe.SecretKey)
	payHandler := handlers.NewPaymentHandler(paySvc, orderService)

	jwtAuth := auth.NewJWTAuthenticator(cfg.JWT.Secret)

	// 5) Mount routes
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// <-- in dev you’ll want to allow your front‑end origin
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // if you ever use cookies or credentialed requests
		MaxAge:           300,  // how long browser can cache the preflight response
	}))

	routes.RegisterUserRoutes(r, userHandler, jwtAuth)
	routes.RegisterProductRoutes(r, prodHandler, jwtAuth)
	routes.RegisterCartRoutes(r, cartHandler, jwtAuth)
	routes.RegisterOrderRoutes(r, orderHandler, jwtAuth)
	routes.RegisterPaymentRoutes(r, payHandler, jwtAuth)
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
