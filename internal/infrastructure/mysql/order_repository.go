package mysql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"richisntreal-backend/internal/core/domain/models"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(o *models.Order) (int64, error) {
	res, err := r.db.Exec(`
        INSERT INTO orders (user_id, total, status, created_at, updated_at)
        VALUES (?, ?, ?, NOW(), NOW())
    `, o.UserID, o.Total, o.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *OrderRepository) CreateOrderItem(item *models.OrderItem) (int64, error) {
	res, err := r.db.Exec(`
        INSERT INTO order_items (order_id, product_id, quantity, unit_price, created_at, updated_at)
        VALUES (?, ?, ?, ?, NOW(), NOW())
    `, item.OrderID, item.ProductID, item.Quantity, item.UnitPrice)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *OrderRepository) FindOrdersByUser(userID int64) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Select(&orders, `
        SELECT id, user_id, total, status, created_at, updated_at
        FROM orders WHERE user_id = ?
    `, userID); err != nil {
		return nil, err
	}
	// load items per order
	for _, ord := range orders {
		var items []models.OrderItem
		r.db.Select(&items, `
            SELECT id, order_id, product_id, quantity, unit_price, created_at, updated_at
            FROM order_items WHERE order_id = ?
        `, ord.ID)
		ord.Items = items
	}
	return orders, nil
}

func (r *OrderRepository) FindOrderByID(orderID int64) (*models.Order, error) {
	var ord models.Order
	if err := r.db.Get(&ord, `
        SELECT id, user_id, total, status, created_at, updated_at
        FROM orders WHERE id = ?
    `, orderID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	var items []models.OrderItem
	r.db.Select(&items, `
        SELECT id, order_id, product_id, quantity, unit_price, created_at, updated_at
        FROM order_items WHERE order_id = ?
    `, ord.ID)
	ord.Items = items
	return &ord, nil
}
