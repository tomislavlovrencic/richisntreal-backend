package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/api/middleware"
	"richisntreal-backend/internal/core/services"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type createOrderResponse struct {
	ID     int64   `json:"id"`
	UserID int64   `json:"user_id"`
	Total  float64 `json:"total"`
}

// CreateOrder converts a cart into a new order, only for the logged‑in user.
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// 0) who’s calling?
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 1) which user’s cart?
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid userID", http.StatusBadRequest)
		return
	}
	// 2) enforce ownership
	if caller != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 3) create the order
	ord, err := h.orderService.CreateOrder(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createOrderResponse{
		ID:     ord.ID,
		UserID: ord.UserID,
		Total:  ord.Total,
	})
	if err != nil {
		return
	}
}

// ListOrders fetches all orders for the logged‑in user.
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid userID", http.StatusBadRequest)
		return
	}
	if caller != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orders, err := h.orderService.GetOrdersForUser(userID)
	if err != nil {
		http.Error(w, "could not fetch orders", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		return
	}
}

// GetOrder fetches a single order by its ID, only if the logged‑in user owns it.
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orderID, err := strconv.ParseInt(chi.URLParam(r, "orderID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid orderID", http.StatusBadRequest)
		return
	}

	ord, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "could not fetch order", http.StatusInternalServerError)
		return
	}
	if ord == nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	// enforce ownership
	if ord.UserID != caller {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err = json.NewEncoder(w).Encode(ord)
	if err != nil {
		return
	}
}
