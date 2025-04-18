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

func (r *ProductRepository) Create(p *models.Product) (int64, error) {
	res, err := r.db.Exec(`
        INSERT INTO products (name, description, price, sku, created_at, updated_at)
        VALUES (?, ?, ?, ?, NOW(), NOW())
    `, p.Name, p.Description, p.Price, p.SKU)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ProductRepository) Update(p *models.Product) error {
	_, err := r.db.Exec(`
        UPDATE products
           SET name = ?, description = ?, price = ?, sku = ?, updated_at = NOW()
         WHERE id = ?
    `, p.Name, p.Description, p.Price, p.SKU, p.ID)
	return err
}

func (r *ProductRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM products WHERE id = ?`, id)
	return err
}
