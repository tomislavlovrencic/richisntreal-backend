package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/core/services"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	// Create and list per-user
	r.Route("/users/{userID}/orders", func(r chi.Router) {
		r.Post("/", h.CreateOrder)
		r.Get("/", h.ListOrders)
	})
	// Retrieve any order by ID
	r.Get("/orders/{orderID}", h.GetOrder)
}

type createOrderResponse struct {
	ID     int64   `json:"id"`
	UserID int64   `json:"user_id"`
	Total  float64 `json:"total"`
}

// CreateOrder converts a cart into an order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
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

// ListOrders fetches all orders for a user
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
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

// GetOrder fetches any single order by its ID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID, _ := strconv.ParseInt(chi.URLParam(r, "orderID"), 10, 64)
	ord, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "could not fetch order", http.StatusInternalServerError)
		return
	}
	if ord == nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(ord)
	if err != nil {
		return
	}
}
