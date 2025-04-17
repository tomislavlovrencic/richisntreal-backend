package bootstrap

import (
	"richisntreal-backend/cmd/config"
	"richisntreal-backend/internal/core/services"
	"richisntreal-backend/internal/infrastructure/database"
	"richisntreal-backend/internal/infrastructure/handlers"
	"richisntreal-backend/internal/infrastructure/middleware"
	"richisntreal-backend/internal/infrastructure/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
	Router *gin.Engine
}

func NewApp() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, cfg)

	// Initialize router
	router := gin.Default()

	// Basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes example
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.GetUint("userID")
			user, err := userService.GetUserByID(userID)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, user)
		})
	}

	return &App{
		Config: cfg,
		DB:     db,
		Router: router,
	}, nil
}

func (app *App) Run() error {
	port := app.Config.Port
	if port == "" {
		port = "8080"
	}

	return app.Router.Run(":" + port)
}
