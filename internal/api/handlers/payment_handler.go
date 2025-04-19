package handlers

import (
	"encoding/json"
	"net/http"
	"richisntreal-backend/internal/api/middleware"
	"strconv"

	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/core/services"
)

// PaymentHandler wires payment endpoints.
type PaymentHandler struct {
	paymentService *services.PaymentService
	orderService   *services.OrderService
}

// NewPaymentHandler constructs.
func NewPaymentHandler(paymentService *services.PaymentService, orderService *services.OrderService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService, orderService: orderService}
}

// paymentRequest is the JSON body for initiating a payment.
type paymentRequest struct {
	Provider string `json:"provider"` // e.g. "stripe"
	Token    string `json:"token"`    // card or payment method token
}

func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// 0) who’s calling?
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 1) parse orderID
	oid, err := strconv.ParseInt(chi.URLParam(r, "orderID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	// 2) fetch order to get amount & owner
	ord, err := h.orderService.GetOrderByID(oid)
	if err != nil {
		http.Error(w, "could not fetch order", http.StatusInternalServerError)
		return
	}
	if ord == nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	// 3) enforce ownership
	if ord.UserID != caller {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 4) decode body
	var req paymentRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	// 5) process payment
	tx, err := h.paymentService.ProcessPayment(ord.ID, ord.Total, "USD", req.Provider, req.Token)
	if err != nil {
		http.Error(w, "payment failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 6) return
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		return
	}
}
