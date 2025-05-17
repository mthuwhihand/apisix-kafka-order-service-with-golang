package models

import (
	order_statuses "hihand/enums/order"
	"hihand/pkgs/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID            uuid.UUID      `gorm:"type:uuid;index;primaryKey" json:"id" form:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;index;primaryKey" json:"user_id" form:"user_id"`
	RecipientName string         `gorm:"index;not null" json:"recipient_name" validate:"required" form:"recipient_name"`
	ContactPhone  string         `gorm:"index;not null" json:"contact_phone" validate:"required" form:"contact_phone"`
	Email         string         `gorm:"index;not null" json:"email" validate:"required,email" form:"email"`
	Address       string         `gorm:"not null" json:"address" validate:"required" form:"address"`
	Status        string         `gorm:"not null" json:"status" validate:"required" form:"status"`
	Total         float64        `gorm:"not null" json:"total" validate:"required,gte=0" form:"total"`
	Note          string         `json:"note" form:"note"`
	Details       []*OrderDetail `gorm:"foreignKey:OrderID,UserID;references:ID,UserID;constraint:OnDelete:CASCADE" json:"details" validate:"dive" form:"details"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderDetail struct {
	ID        uuid.UUID `gorm:"type:uuid;index;primaryKey" json:"id" form:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;primaryKey" json:"user_id" form:"user_id"` // shard key
	OrderID   uuid.UUID `gorm:"type:uuid;index" json:"order_id" form:"order_id"`
	ProductID string    `gorm:"index;not null" json:"product_id" validate:"required" form:"product_id"`
	Name      string    `gorm:"index;not null" json:"name" validate:"required" form:"name"`
	Price     float64   `gorm:"not null" json:"price" validate:"required,gte=0" form:"price"`
	Quantity  int       `gorm:"not null" json:"quantity" validate:"required,gte=0" form:"quantity"`
	Total     float64   `gorm:"not null" json:"total" validate:"required,gte=0" form:"total"`
}

// BeforeCreate hooks tạo UUID mới cho Order
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.Status == "" {
		o.Status = order_statuses.Created
	}
	return nil
}

// BeforeCreate hooks tạo UUID mới cho OrderDetail
func (od *OrderDetail) BeforeCreate(tx *gorm.DB) (err error) {
	if od.ID == uuid.Nil {
		od.ID = uuid.New()
	}
	return nil
}

// Hàm chuyển đổi OrderDetail sang map với UUID thành string
func (od *OrderDetail) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         od.ID.String(),
		"order_id":   od.OrderID.String(),
		"user_id":    od.UserID.String(),
		"product_id": od.ProductID,
		"name":       od.Name,
		"price":      od.Price,
		"quantity":   od.Quantity,
		"total":      od.Total,
	}
}

// Hàm FromMap chuyển string sang uuid.UUID và kiểu dữ liệu khác
func (od *OrderDetail) FromMap(data map[string]interface{}) {
	od.ID, _ = uuid.Parse(utils.GetString(data, "id"))
	od.OrderID, _ = uuid.Parse(utils.GetString(data, "order_id"))
	od.UserID, _ = uuid.Parse(utils.GetString(data, "user_id"))
	od.ProductID = utils.GetString(data, "product_id")
	od.Name = utils.GetString(data, "name")
	od.Price = utils.GetFloat64(data, "price")
	od.Quantity = utils.GetInt(data, "quantity")
	od.Total = utils.GetFloat64(data, "total")
}

// ToMap Order -> map
func (o *Order) ToMap() (map[string]interface{}, error) {
	details := []map[string]interface{}{}
	for _, detail := range o.Details {
		details = append(details, detail.ToMap())
	}

	return map[string]interface{}{
		"id":             o.ID.String(),
		"user_id":        o.UserID.String(),
		"recipient_name": o.RecipientName,
		"contact_phone":  o.ContactPhone,
		"address":        o.Address,
		"email":          o.Email,
		"status":         o.Status,
		"total":          o.Total,
		"details":        details,
		"note":           o.Note,
		"created_at":     o.CreatedAt,
		"updated_at":     o.UpdatedAt,
	}, nil
}

// FromMap map -> Order
func (o *Order) FromMap(data map[string]interface{}) error {
	o.ID, _ = uuid.Parse(utils.GetString(data, "id"))
	o.UserID, _ = uuid.Parse(utils.GetString(data, "user_id"))
	o.RecipientName = utils.GetString(data, "recipient_name")
	o.ContactPhone = utils.GetString(data, "contact_phone")
	o.Address = utils.GetString(data, "address")
	o.Email = utils.GetString(data, "email")
	o.Status = utils.GetString(data, "status")
	o.Note = utils.GetString(data, "note")
	o.CreatedAt = utils.GetTime(data, "created_at")
	o.UpdatedAt = utils.GetTime(data, "updated_at")
	o.Total = utils.GetFloat64(data, "total")

	if details, ok := data["details"].([]interface{}); ok {
		for _, d := range details {
			if detailMap, ok := d.(map[string]interface{}); ok {
				detail := &OrderDetail{}
				detail.FromMap(detailMap)
				o.Details = append(o.Details, detail)
			}
		}
	}
	return nil
}
