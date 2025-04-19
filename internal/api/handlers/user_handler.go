package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"richisntreal-backend/internal/api/middleware"
	"richisntreal-backend/internal/core/services"
	"strconv"

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

type loginResponse struct {
	Token string      `json:"token"`
	User  userProfile `json:"user"`
}

type userProfile struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Login handles user authentication.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// 2a) Authenticate & get token
	token, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 2b) Fetch the user record (so we can include first/last name):
	user, err := h.userService.GetByEmail(req.Email)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// 3) Marshal combined response
	resp := loginResponse{
		Token: token,
		User: userProfile{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

// GetUser handles GET /users/{id}; only the logged‑in user may fetch their own record.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// 1) Extract the caller’s userID from the context (AuthMiddleware must have run).
	caller := middleware.FromContext(r.Context())
	if caller == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2) Parse the URL param
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// 3) Enforce ownership
	if caller != id {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 4) Fetch & return
	user, err := h.userService.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 5) Write JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return
	}
}
