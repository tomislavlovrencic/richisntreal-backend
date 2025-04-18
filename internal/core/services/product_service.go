package services

import "richisntreal-backend/internal/core/domain/models"

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

// ProductRepository defines persistence operations for products.
type ProductRepository interface {
	FindAll() ([]*models.Product, error)
	FindByID(id int64) (*models.Product, error)
}
