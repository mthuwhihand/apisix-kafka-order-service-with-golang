package repositories

import (
	"hihand/internal/models"

	"gorm.io/gorm"
)

type (
	OrderRepository interface {
		Create(order *models.Order) error
		Update(id string, updates map[string]interface{}) error
		Delete(orderID string) error
	}

	orderRepository struct {
		db *gorm.DB
	}
)

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (o *orderRepository) Create(order *models.Order) error {
	return o.db.Create(order).Error
}

func (o *orderRepository) Update(id string, updates map[string]interface{}) error {
	return o.db.Model(&models.Order{}).Where("id = ?", id).Updates(updates).Error
}

func (o *orderRepository) Delete(id string) error {
	return o.db.Delete(&models.Order{}, id).Error
}
