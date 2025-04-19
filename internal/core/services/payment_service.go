package services

import (
	"errors"
	"fmt"
	"strings"

	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/charge"

	"richisntreal-backend/internal/core/domain/models"
)

// ErrPaymentFailed is returned when the gateway reports a failure.
var ErrPaymentFailed = errors.New("payment failed")

// PaymentService handles charging and recording payment transactions.
type PaymentService struct {
	paymentRepository PaymentRepository
	stripeKey         string
}

// NewPaymentService constructs a PaymentService.
// Pass cfg.App.StripeSecretKey from your bootstrap.
func NewPaymentService(paymentRepository PaymentRepository, stripeKey string) *PaymentService {
	return &PaymentService{
		paymentRepository: paymentRepository,
		stripeKey:         stripeKey,
	}
}

// ProcessPayment creates a pending record, sends a Stripe charge,
// updates the record with the result, and returns the final transaction.
func (s *PaymentService) ProcessPayment(
	orderID int64,
	amount float64,
	currency, provider, token string,
) (*models.PaymentTransaction, error) {
	// 1) Create a pending transaction
	tx := &models.PaymentTransaction{
		OrderID:  orderID,
		Amount:   amount,
		Currency: strings.ToLower(currency),
		Provider: provider,
		Token:    token,
		Status:   "pending",
	}
	id, err := s.paymentRepository.Create(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}
	tx.ID = id

	// 2) Configure Stripe and send the charge
	stripe.Key = s.stripeKey

	params := &stripe.ChargeParams{
		Amount:   stripe.Int64(int64(amount * 100)), // convert dollars to cents
		Currency: stripe.String(tx.Currency),
	}
	err = params.SetSource(token)
	if err != nil {
		return nil, err
	} // e.g. "tok_visa" in test mode

	ch, err := charge.New(params)
	if err != nil {
		// update record as failed
		msg := err.Error()
		_ = s.paymentRepository.UpdateStatus(id, "failed", &msg)
		return nil, fmt.Errorf("%w: %s", ErrPaymentFailed, err.Error())
	}

	// 3) On success, update our transaction
	providerTxID := ch.ID
	tx.ProviderTxID = &providerTxID
	tx.Status = string(ch.Status) // e.g. "succeeded"
	if err = s.paymentRepository.UpdateStatus(id, tx.Status, nil); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	return tx, nil
}

// GetPaymentByOrder fetches the transaction associated with an order.
func (s *PaymentService) GetPaymentByOrder(orderID int64) (*models.PaymentTransaction, error) {
	return s.paymentRepository.FindByOrder(orderID)
}

// PaymentRepository required by PaymentService.
type PaymentRepository interface {
	Create(tx *models.PaymentTransaction) (int64, error)
	UpdateStatus(id int64, status string, failureMessage *string) error
	FindByOrder(orderID int64) (*models.PaymentTransaction, error)
}
