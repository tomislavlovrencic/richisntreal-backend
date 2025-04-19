package routes

import (
	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/api/auth"
	"richisntreal-backend/internal/api/handlers"
	"richisntreal-backend/internal/api/middleware"
)

func RegisterCartRoutes(
	r chi.Router,
	h *handlers.CartHandler,
	jwtAuth auth.Authenticator,
) {
	r.Route("/users/{userID}/cart", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtAuth))
		r.Get("/", h.GetCart)
		r.Post("/items", h.AddItem)
		r.Put("/items/{itemID}", h.UpdateItem)
		r.Delete("/items/{itemID}", h.RemoveItem)
		r.Delete("/", h.ClearCart)
	})
}
