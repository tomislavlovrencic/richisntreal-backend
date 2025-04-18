package handlers

import (
	"encoding/json"
	"net/http"
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

// RegisterRoutes mounts the pay endpoint.
func (h *PaymentHandler) RegisterRoutes(r chi.Router) {
	// POST /orders/{orderID}/pay
	r.Post("/orders/{orderID}/pay", h.ProcessPayment)
}

// paymentRequest is the JSON body for initiating a payment.
type paymentRequest struct {
	Provider string `json:"provider"` // e.g. "stripe"
	Token    string `json:"token"`    // card or payment method token
}

// ProcessPayment handles POST /orders/{orderID}/pay
func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// 1) parse orderID
	oid, err := strconv.ParseInt(chi.URLParam(r, "orderID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	// 2) decode body
	var req paymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	// 3) fetch order to get amount
	ord, err := h.orderService.GetOrderByID(oid)
	if err != nil {
		http.Error(w, "could not fetch order", http.StatusInternalServerError)
		return
	}
	if ord == nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	// 4) process payment (using USD as default currency)
	tx, err := h.paymentService.ProcessPayment(ord.ID, ord.Total, "USD", req.Provider, req.Token)
	if err != nil {
		http.Error(w, "payment failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 5) return the transaction record
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		return
	}
}
