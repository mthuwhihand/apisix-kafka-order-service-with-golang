package services

import (
	"hihand/internal/models"
	"hihand/internal/repositories"
)

type (
	OrderService interface {
		Search(query string, limit int, skip int) ([]*models.Order, error)
		Create(order *models.Order) error
		Update(id string, updates map[string]interface{}) error
		Delete(orderID string) error
	}

	orderService struct {
		repo repositories.OrderRepository
	}
)

func NewOrderService(repo repositories.OrderRepository) OrderService {
	return &orderService{
		repo: repo,
	}
}

func (s *orderService) Create(order *models.Order) error {
	return s.repo.Create(order)
}

func (s *orderService) Update(id string, updates map[string]interface{}) error {
	return s.repo.Update(id, updates)
}
func (s *orderService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *orderService) Search(query string, limit int, skip int) ([]*models.Order, error) {
	return s.repo.Search(query, limit, skip)
}
