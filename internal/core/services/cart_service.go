package services

import (
	"errors"
	"richisntreal-backend/internal/core/domain/models"
)

type CartService struct {
	cartRepository CartRepository
}

func NewCartService(cartRepository CartRepository) *CartService {
	return &CartService{cartRepository: cartRepository}
}

func (s *CartService) GetCart(userID int64) (*models.Cart, error) {
	cart, err := s.cartRepository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		id, err := s.cartRepository.CreateCart(&models.Cart{UserID: userID})
		if err != nil {
			return nil, err
		}
		cart, _ = s.cartRepository.FindByUserID(userID) // now it exists
		cart.ID = id
	}
	return cart, nil
}

func (s *CartService) AddItem(userID, productID int64, qty int, price float64) (*models.CartItem, error) {
	cart, err := s.GetCart(userID)
	if err != nil {
		return nil, err
	}
	// merge if exists
	if existing, _ := s.cartRepository.FindItemByCartAndProduct(cart.ID, productID); existing != nil {
		existing.Quantity += qty
		if err := s.cartRepository.UpdateItem(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}
	item := &models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  qty,
		UnitPrice: price,
	}
	id, err := s.cartRepository.CreateItem(item)
	if err != nil {
		return nil, err
	}
	item.ID = id
	return item, nil
}

func (s *CartService) UpdateItem(itemID int64, qty int) (*models.CartItem, error) {
	item, err := s.cartRepository.FindItem(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrCartItemNotFound
	}
	item.Quantity = qty
	if err := s.cartRepository.UpdateItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *CartService) RemoveItem(itemID int64) error {
	return s.cartRepository.DeleteItem(itemID)
}

func (s *CartService) ClearCart(userID int64) error {
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}
	return s.cartRepository.DeleteItemsByCartID(cart.ID)
}

var ErrCartItemNotFound = errors.New("cart item not found")

type CartRepository interface {
	FindByUserID(userID int64) (*models.Cart, error)
	CreateCart(cart *models.Cart) (int64, error)
	FindItem(itemID int64) (*models.CartItem, error)
	FindItemByCartAndProduct(cartID, productID int64) (*models.CartItem, error)
	CreateItem(item *models.CartItem) (int64, error)
	UpdateItem(item *models.CartItem) error
	DeleteItem(itemID int64) error
	DeleteItemsByCartID(cartID int64) error
}
