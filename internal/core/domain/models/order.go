package models

import "time"

// Order represents a user's purchase.
type Order struct {
	ID        int64       `db:"id" json:"id"`
	UserID    int64       `db:"user_id" json:"user_id"`
	Total     float64     `db:"total" json:"total"`
	Status    string      `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
	Items     []OrderItem `json:"items"`
}

// OrderItem is a single line item in an order.
type OrderItem struct {
	ID        int64     `db:"id" json:"id"`
	OrderID   int64     `db:"order_id" json:"order_id"`
	ProductID int64     `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	UnitPrice float64   `db:"unit_price" json:"unit_price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
