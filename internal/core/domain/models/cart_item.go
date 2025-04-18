package models

import "time"

type CartItem struct {
	ID        int64     `db:"id" json:"id"`
	CartID    int64     `db:"cart_id" json:"cart_id"`
	ProductID int64     `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	UnitPrice float64   `db:"unit_price" json:"unit_price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
