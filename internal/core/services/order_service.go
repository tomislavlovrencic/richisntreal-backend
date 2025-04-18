package services

import (
	"errors"
	"richisntreal-backend/internal/core/domain/models"
)

type OrderService struct {
	orderRepository OrderRepository
	cartRepository  CartRepository
}

func NewOrderService(orderRepository OrderRepository, cartRepository CartRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository, cartRepository: cartRepository}
}

func (s *OrderService) CreateOrder(userID int64) (*models.Order, error) {
	// 1) fetch the cart
	cart, err := s.cartRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// 2) calculate total
	var total float64
	for _, ci := range cart.Items {
		total += float64(ci.Quantity) * ci.UnitPrice
	}

	// 3) insert into orders table
	order := &models.Order{
		UserID: userID,
		Total:  total,
		Status: "pending",
	}
	orderID, err := s.orderRepository.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	order.ID = orderID

	// 4) insert each cart item as an order_item
	for _, ci := range cart.Items {
		oi := &models.OrderItem{
			OrderID:   orderID,
			ProductID: ci.ProductID,
			Quantity:  ci.Quantity,
			UnitPrice: ci.UnitPrice,
		}
		_, err := s.orderRepository.CreateOrderItem(oi)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, *oi)
	}

	// 5) clear the cart
	if err := s.cartRepository.DeleteItemsByCartID(cart.ID); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrdersForUser(userID int64) ([]*models.Order, error) {
	return s.orderRepository.FindOrdersByUser(userID)
}

func (s *OrderService) GetOrderByID(orderID int64) (*models.Order, error) {
	ord, err := s.orderRepository.FindOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	if ord == nil {
		return nil, ErrOrderNotFound
	}
	return ord, nil
}

var ErrOrderNotFound = errors.New("order not found")

type OrderRepository interface {
	CreateOrder(o *models.Order) (int64, error)
	CreateOrderItem(item *models.OrderItem) (int64, error)
	FindOrdersByUser(userID int64) ([]*models.Order, error)
	FindOrderByID(orderID int64) (*models.Order, error)
}
