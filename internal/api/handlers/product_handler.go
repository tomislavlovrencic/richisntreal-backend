package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/core/services"
)

// ProductHandler wires HTTP requests to product services.
type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	r.Get("/products", h.List)
	r.Get("/products/{id}", h.GetByID)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	prods, err := h.productService.ListProducts()
	if err != nil {
		http.Error(w, "could not fetch products", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(prods)
	if err != nil {
		return
	}
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	prod, err := h.productService.GetProductByID(id)
	if err != nil {
		http.Error(w, "could not fetch product", http.StatusInternalServerError)
		return
	}
	if prod == nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(prod)
	if err != nil {
		return
	}
}
