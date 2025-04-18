package mysql

import (
	"github.com/jmoiron/sqlx"
	"richisntreal-backend/internal/core/domain/models"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(tx *models.PaymentTransaction) (int64, error) {
	res, err := r.db.Exec(`
        INSERT INTO payment_transactions
            (order_id, amount, currency, provider, token, status, created_at, updated_at)
        VALUES
            (?, ?, ?, ?, ?, ?, NOW(), NOW())
    `,
		tx.OrderID, tx.Amount, tx.Currency, tx.Provider, tx.Token, tx.Status,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PaymentRepository) UpdateStatus(id int64, status string, failureMessage *string) error {
	_, err := r.db.Exec(`
        UPDATE payment_transactions
           SET status = ?, failure_message = ?, updated_at = NOW()
         WHERE id = ?
    `,
		status, failureMessage, id,
	)
	return err
}

func (r *PaymentRepository) FindByOrder(orderID int64) (*models.PaymentTransaction, error) {
	var tx models.PaymentTransaction
	err := r.db.Get(&tx, `
        SELECT id, order_id, amount, currency, provider, provider_tx_id, token, status, failure_message, created_at, updated_at
          FROM payment_transactions
         WHERE order_id = ?
    `, orderID)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}
