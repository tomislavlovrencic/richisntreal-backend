package routes

import (
	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/api/auth"
	"richisntreal-backend/internal/api/handlers"
	"richisntreal-backend/internal/api/middleware"
)

// RegisterUserRoutes wires up endpoints for users.
func RegisterUserRoutes(
	r chi.Router,
	h *handlers.UserHandler,
	jwtAuth auth.Authenticator,
) {
	// public
	r.Post("/users", h.CreateUser)
	r.Post("/login", h.Login)

	// private
	r.With(middleware.AuthMiddleware(jwtAuth)).
		Get("/users/{id}", h.GetUser)
}
