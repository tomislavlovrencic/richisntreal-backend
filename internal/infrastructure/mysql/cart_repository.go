package mysql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"richisntreal-backend/internal/core/domain/models"
)

type CartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) FindByUserID(userID int64) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Get(&cart, `
        SELECT id, user_id, created_at, updated_at
          FROM carts
         WHERE user_id = ?`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	err = r.db.Select(&cart.Items, `
        SELECT id, cart_id, product_id, quantity, unit_price, created_at, updated_at
          FROM cart_items
         WHERE cart_id = ?`, cart.ID)
	return &cart, err
}

func (r *CartRepository) CreateCart(cart *models.Cart) (int64, error) {
	res, err := r.db.Exec(`
        INSERT INTO carts (user_id, created_at, updated_at)
             VALUES (?, NOW(), NOW())`, cart.UserID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CartRepository) FindItem(itemID int64) (*models.CartItem, error) {
	var it models.CartItem
	err := r.db.Get(&it, `
        SELECT id, cart_id, product_id, quantity, unit_price, created_at, updated_at
          FROM cart_items
         WHERE id = ?`, itemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &it, nil
}

func (r *CartRepository) FindItemByCartAndProduct(cartID, productID int64) (*models.CartItem, error) {
	var it models.CartItem
	err := r.db.Get(&it, `
        SELECT id, cart_id, product_id, quantity, unit_price, created_at, updated_at
          FROM cart_items
         WHERE cart_id = ? AND product_id = ?`, cartID, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &it, nil
}

func (r *CartRepository) CreateItem(item *models.CartItem) (int64, error) {
	res, err := r.db.Exec(`
        INSERT INTO cart_items (cart_id, product_id, quantity, unit_price, created_at, updated_at)
             VALUES (?, ?, ?, ?, NOW(), NOW())`,
		item.CartID, item.ProductID, item.Quantity, item.UnitPrice,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CartRepository) UpdateItem(item *models.CartItem) error {
	_, err := r.db.Exec(`
        UPDATE cart_items
           SET quantity = ?, updated_at = NOW()
         WHERE id = ?`,
		item.Quantity, item.ID,
	)
	return err
}

func (r *CartRepository) DeleteItem(itemID int64) error {
	_, err := r.db.Exec(`DELETE FROM cart_items WHERE id = ?`, itemID)
	return err
}

func (r *CartRepository) DeleteItemsByCartID(cartID int64) error {
	_, err := r.db.Exec(`DELETE FROM cart_items WHERE cart_id = ?`, cartID)
	return err
}
