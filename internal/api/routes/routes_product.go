package routes

import (
	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/api/auth"
	"richisntreal-backend/internal/api/handlers"
	"richisntreal-backend/internal/api/middleware"
)

func RegisterProductRoutes(
	r chi.Router,
	h *handlers.ProductHandler,
	jwtAuth auth.Authenticator,
) {
	// public
	r.Get("/products", h.List)
	r.Get("/products/{id}", h.GetByID)

	// admin
	r.With(middleware.AuthMiddleware(jwtAuth)).
		Post("/products", h.Create)
	r.With(middleware.AuthMiddleware(jwtAuth)).
		Put("/products/{id}", h.Update)
	r.With(middleware.AuthMiddleware(jwtAuth)).
		Delete("/products/{id}", h.Delete)
}
