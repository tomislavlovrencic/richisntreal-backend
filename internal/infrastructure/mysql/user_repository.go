package mysql

import (
	"database/sql"
	"errors"
	"richisntreal-backend/internal/core/domain/models"

	"github.com/jmoiron/sqlx"
)

// UserRepository implements persistence for users using MySQL.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository constructs a new MySQL-backed UserRepository.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ExistsByEmail returns true if a user with the given email already exists.
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(1) FROM users WHERE email = ?`, email)
	return count > 0, err
}

// Create inserts a new user and returns its generated ID.
func (r *UserRepository) Create(user *models.User) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO users (username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FindByEmail retrieves a user by email (including password hash).
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = ?`, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByID retrieves a user (including password) by their ID.
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, `
        SELECT id, username, email, password, created_at, updated_at
          FROM users
         WHERE id = ?`, id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
