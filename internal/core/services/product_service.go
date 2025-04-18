package services

import (
	"errors"

	"richisntreal-backend/internal/core/domain/models"
)

var ErrProductNotFound = errors.New("product not found")

// ProductService holds product business logic.
type ProductService struct {
	productRepository ProductRepository
}

func NewProductService(productRepository ProductRepository) *ProductService {
	return &ProductService{productRepository: productRepository}
}

func (s *ProductService) ListProducts() ([]*models.Product, error) {
	return s.productRepository.FindAll()
}

func (s *ProductService) GetProductByID(id int64) (*models.Product, error) {
	return s.productRepository.FindByID(id)
}

func (s *ProductService) CreateProduct(name, description, sku string, price float64) (*models.Product, error) {
	p := &models.Product{
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
	}
	id, err := s.productRepository.Create(p)
	if err != nil {
		return nil, err
	}
	p.ID = id
	return p, nil
}

func (s *ProductService) UpdateProduct(id int64, name, description, sku string, price float64) (*models.Product, error) {
	existing, err := s.productRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrProductNotFound
	}
	existing.Name = name
	existing.Description = description
	existing.SKU = sku
	existing.Price = price
	if err := s.productRepository.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *ProductService) DeleteProduct(id int64) error {
	return s.productRepository.Delete(id)
}

// ProductRepository defines persistence operations for products.
type ProductRepository interface {
	FindAll() ([]*models.Product, error)
	FindByID(id int64) (*models.Product, error)
	Create(p *models.Product) (int64, error)
	Update(p *models.Product) error
	Delete(id int64) error
}
