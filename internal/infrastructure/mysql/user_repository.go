// File: internal/infrastructure/mysql/user_repository.go
package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"richisntreal-backend/internal/core/domain/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// UserRepository implements persistence for users using MySQL.
type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(1) FROM users WHERE email = ?`, email)
	return count > 0, err
}

func (r *UserRepository) Create(user *models.User) (int64, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
    INSERT INTO users
        (username, email, password, first_name, last_name, country, date_of_birth, created_at, updated_at)
    VALUES
        (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Country,
		user.DateOfBirth,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("UserRepository.Create: %w", err)
	}
	return res.LastInsertId()
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	query := `
    SELECT id, username, email, password,
           first_name, last_name, country, date_of_birth,
           created_at, updated_at
      FROM users
     WHERE email = ?
     LIMIT 1`
	err := r.db.Get(&u, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepository.FindByEmail: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var u models.User
	query := `
    SELECT id, username, email, password,
           first_name, last_name, country, date_of_birth,
           created_at, updated_at
      FROM users
     WHERE id = ?`
	err := r.db.Get(&u, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepository.FindByID: %w", err)
	}
	return &u, nil
}
