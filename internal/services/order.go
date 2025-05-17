package services

import (
	"errors"
	"hihand/internal/models"

	"gorm.io/gorm"
)

type (
	OrderService interface {
		Search(query string, limit int, skip int) (*models.ListResponse[models.Order], error)
		Create(order *models.Order) error
		Update(id string, updates map[string]interface{}) error
		Delete(orderID string) error
	}

	orderService struct {
		db *gorm.DB
	}
)

func NewOrderService(db *gorm.DB) OrderService {
	return &orderService{
		db: db,
	}
}

func (o *orderService) Create(order *models.Order) error {
	return o.db.Create(order).Error
}

func (o *orderService) Update(id string, updates map[string]interface{}) error {
	return o.db.Model(&models.Order{}).Where("id = ?", id).Updates(updates).Error
}

func (o *orderService) Delete(id string) error {
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

func (o *orderService) Search(query string, limit int, skip int) (*models.ListResponse[models.Order], error) {
	var orders []models.Order

	tx := o.db.Model(&models.Order{}).Preload("Details")

	if query != "" {
		tx = tx.Where("recipient_name LIKE ? OR email LIKE ? OR contact_phone LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}
	// Đếm tổng số bản ghi phù hợp query
	var totalRecords int64
	if err := tx.Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	// Áp dụng limit, skip (offset)
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	if skip > 0 {
		tx = tx.Offset(skip)
	}

	// Lấy dữ liệu
	err := tx.Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	// Tính các thông số phân trang
	currentPage := 1
	if limit > 0 {
		currentPage = (skip / limit) + 1
	}

	totalPages := 1
	if limit > 0 {
		totalPages = int((totalRecords + int64(limit) - 1) / int64(limit)) // làm tròn lên
	}
	hasNext := currentPage < totalPages

	// Chuẩn bị kết quả trả về
	listRes := models.ListResponse[models.Order]{
		TotalRecords: totalRecords,
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		HasNext:      hasNext,
		Data:         orders,
	}

	return &listRes, nil
}
