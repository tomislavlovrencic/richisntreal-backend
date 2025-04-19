package routes

import (
	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/api/auth"
	"richisntreal-backend/internal/api/handlers"
	"richisntreal-backend/internal/api/middleware"
)

func RegisterPaymentRoutes(
	r chi.Router,
	h *handlers.PaymentHandler,
	jwtAuth auth.Authenticator,
) {
	r.Route("/orders/{orderID}/pay", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtAuth))
		r.Post("/", h.ProcessPayment)
	})
}
