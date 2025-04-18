package mysql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"richisntreal-backend/internal/core/domain/models"
)

// ProductRepository implements persistence for products.
type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll() ([]*models.Product, error) {
	var prods []*models.Product
	err := r.db.Select(&prods, `
        SELECT id, name, description, price, sku, created_at, updated_at
          FROM products
    `)
	return prods, err
}

func (r *ProductRepository) FindByID(id int64) (*models.Product, error) {
	var p models.Product
	err := r.db.Get(&p, `
        SELECT id, name, description, price, sku, created_at, updated_at
          FROM products
         WHERE id = ?
    `, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
