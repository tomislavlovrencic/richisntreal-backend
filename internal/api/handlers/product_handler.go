package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"richisntreal-backend/internal/core/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type productRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price"`
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
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
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

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	prod, err := h.productService.CreateProduct(req.Name, req.Description, req.SKU, req.Price)
	if err != nil {
		http.Error(w, "could not create product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(prod)
	if err != nil {
		return
	}
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}
	var req productRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	prod, err := h.productService.UpdateProduct(id, req.Name, req.Description, req.SKU, req.Price)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
		} else {
			http.Error(w, "could not update product", http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(prod)
	if err != nil {
		return
	}
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}
	if err = h.productService.DeleteProduct(id); err != nil {
		http.Error(w, "could not delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
