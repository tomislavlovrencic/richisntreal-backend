package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/core/services"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users/{userID}/cart", func(r chi.Router) {
		r.Get("/", h.GetCart)
		r.Post("/items", h.AddItem)
		r.Put("/items/{itemID}", h.UpdateItem)
		r.Delete("/items/{itemID}", h.RemoveItem)
		r.Delete("/", h.ClearCart)
	})
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		http.Error(w, "could not fetch cart", http.StatusInternalServerError)
		return
	}
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
	userID, _ := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	var req addItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	itemID, _ := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
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
	itemID, _ := strconv.ParseInt(chi.URLParam(r, "itemID"), 10, 64)
	if err := h.cartService.RemoveItem(itemID); err != nil {
		http.Error(w, "could not remove item", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err := h.cartService.ClearCart(userID); err != nil {
		http.Error(w, "could not clear cart", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
