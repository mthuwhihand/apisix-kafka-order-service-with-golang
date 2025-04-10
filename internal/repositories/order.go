package repositories

import (
	"errors"
	"hihand/internal/models"

	"gorm.io/gorm"
)

type (
	OrderRepository interface {
		Search(query string, limit int, skip int) ([]*models.Order, error)
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
	if id == "" {
		return errors.New("invalid order ID")
	}

	result := o.db.Where("id = ?", id).Delete(&models.Order{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("order not found or already deleted")
	}

	return nil
}

func (o *orderRepository) Search(query string, limit int, skip int) ([]*models.Order, error) {
	var orders []*models.Order

	tx := o.db.Model(&models.Order{}).Preload("Details")

	if query != "" {
		tx = tx.Where("recipient_name LIKE ? OR email LIKE ? OR contact_phone LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if limit > 0 {
		tx = tx.Limit(limit)
	}
	if skip > 0 {
		tx = tx.Offset(skip)
	}

	err := tx.Find(&orders).Error
	return orders, err
}
