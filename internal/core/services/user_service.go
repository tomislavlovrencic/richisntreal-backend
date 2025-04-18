package services

import (
	"errors"
	"richisntreal-backend/internal/core/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// UserService is the default implementation of UserService.
type UserService struct {
	userRepository UserRepository
	jwtSecret      string
}

// NewUserService constructs a new UserService.
func NewUserService(userRepository UserRepository, jwtSecret string) *UserService {
	return &UserService{userRepository: userRepository, jwtSecret: jwtSecret}
}

// CreateUser registers a new user, hashing their password.
func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
	exists, err := s.userRepository.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  string(hash),
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	user.Password = "" // clear password before returning
	return user, nil
}

// Authenticate verifies credentials and returns a JWT on success.
func (s *UserService) Authenticate(email, password string) (string, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ErrUserExists is returned when trying to register with an email that's already taken.
var ErrUserExists = errors.New("user already exists")

// ErrInvalidCredentials is returned when login credentials are incorrect.
var ErrInvalidCredentials = errors.New("invalid credentials")

// UserRepository defines persistence operations for users.
type UserRepository interface {
	ExistsByEmail(email string) (bool, error)
	Create(user *models.User) (int64, error)
	FindByEmail(email string) (*models.User, error)
}
