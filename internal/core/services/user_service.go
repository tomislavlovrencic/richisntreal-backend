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

func (s *UserService) CreateUser(
	username, email, password, firstName, lastName, country string,
	dateOfBirth *time.Time,
) (*models.User, error) {
	// 1) check for duplicate email
	exists, err := s.userRepository.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// 2) hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &models.User{
		Username:    username,
		Email:       email,
		Password:    string(hash),
		FirstName:   firstName,
		LastName:    lastName,
		Country:     country,
		DateOfBirth: dateOfBirth,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 3) persist
	id, err := s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}
	user.ID = id

	// 4) clear password hash before returning
	user.Password = ""

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

// GetByID looks up a user by ID (stripping out their password).
func (s *UserService) GetByID(id int64) (*models.User, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.Password = ""
	return user, nil
}

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

// UserRepository defines persistence operations for users.
type UserRepository interface {
	ExistsByEmail(email string) (bool, error)
	Create(user *models.User) (int64, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
}
