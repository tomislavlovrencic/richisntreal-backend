package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"richisntreal-backend/internal/core/services"

	"github.com/go-chi/chi/v5"
)

// UserHandler wires HTTP requests to user-related services.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler constructs a new UserHandler.
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRoutes binds user routes onto the router.
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/users", h.CreateUser)
	r.Post("/login", h.Login)
}

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// CreateUser handles user registration.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createUserResponse{ID: user.ID, Username: user.Username, Email: user.Email})
	if err != nil {
		return
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login handles user authentication.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = json.NewEncoder(w).Encode(loginResponse{Token: token})
	if err != nil {
		return
	}
}
