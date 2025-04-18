package models

import "time"

// PaymentTransaction records an attempt to pay for an order.
type PaymentTransaction struct {
	ID             int64     `db:"id" json:"id"`
	OrderID        int64     `db:"order_id" json:"order_id"`
	Amount         float64   `db:"amount" json:"amount"`
	Currency       string    `db:"currency" json:"currency"`
	Provider       string    `db:"provider" json:"provider"`
	ProviderTxID   *string   `db:"provider_tx_id" json:"provider_tx_id,omitempty"`
	Token          string    `db:"token" json:"token"`
	Status         string    `db:"status" json:"status"`
	FailureMessage *string   `db:"failure_message" json:"failure_message,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
