package routes

import (
	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/api/auth"
	"richisntreal-backend/internal/api/handlers"
	"richisntreal-backend/internal/api/middleware"
)

func RegisterOrderRoutes(
	r chi.Router,
	h *handlers.OrderHandler,
	jwtAuth auth.Authenticator,
) {
	// user’s orders
	r.Route("/users/{userID}/orders", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtAuth))
		r.Post("/", h.CreateOrder)
		r.Get("/", h.ListOrders)
	})

	// fetch any single order
	r.With(middleware.AuthMiddleware(jwtAuth)).
		Get("/orders/{orderID}", h.GetOrder)
}
