package services

import (
	"fmt"

	"github.com/google/uuid"
	"richisntreal-backend/internal/core/domain/models"
)

type PaymentService struct {
	paymentRepository PaymentRepository
}

func NewPaymentService(paymentRepository PaymentRepository) *PaymentService {
	return &PaymentService{paymentRepository: paymentRepository}
}

func (s *PaymentService) ProcessPayment(
	orderID int64,
	amount float64,
	currency, provider, token string,
) (*models.PaymentTransaction, error) {
	// 1) Create a pending transaction
	tx := &models.PaymentTransaction{
		OrderID:  orderID,
		Amount:   amount,
		Currency: currency,
		Provider: provider,
		Token:    token,
		Status:   "pending",
	}
	id, err := s.paymentRepository.Create(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}
	tx.ID = id

	// 2) Simulate charging—here we always succeed with a random provider ID
	providerTxID := uuid.New().String()
	tx.ProviderTxID = &providerTxID
	tx.Status = "succeeded"

	// 3) Update the record
	if err := s.paymentRepository.UpdateStatus(id, tx.Status, nil); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	return tx, nil
}

func (s *PaymentService) GetPaymentByOrder(orderID int64) (*models.PaymentTransaction, error) {
	return s.paymentRepository.FindByOrder(orderID)
}

// PaymentRepository required by PaymentService.
type PaymentRepository interface {
	Create(tx *models.PaymentTransaction) (int64, error)
	UpdateStatus(id int64, status string, failureMessage *string) error
	FindByOrder(orderID int64) (*models.PaymentTransaction, error)
}
