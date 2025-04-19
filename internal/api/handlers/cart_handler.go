package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/api/middleware"
	"richisntreal-backend/internal/core/services"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	// 1) Authenticate & authorize
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if caller != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 2) Fetch cart
	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		http.Error(w, "could not fetch cart", http.StatusInternalServerError)
		return
	}

	// 3) Respond
	err = json.NewEncoder(w).Encode(cart)
	if err != nil {
		return
	}
}

type addItemReq struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if caller != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req addItemReq
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	item, err := h.cartService.AddItem(userID, req.ProductID, req.Quantity, req.UnitPrice)
	if err != nil {
		http.Error(w, "could not add item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		return
	}
}

type updateItemReq struct {
	Quantity int `json:"quantity"`
}

func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if caller != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	var req updateItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	item, err := h.cartService.UpdateItem(itemID, req.Quantity)
	if err != nil {
		http.Error(w, "could not update item", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		return
	}
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if caller != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	itemID, err := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	if err = h.cartService.RemoveItem(itemID); err != nil {
		http.Error(w, "could not remove item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	if caller != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err = h.cartService.ClearCart(userID); err != nil {
		http.Error(w, "could not clear cart", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
