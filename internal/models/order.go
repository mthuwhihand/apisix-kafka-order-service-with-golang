package models

import (
	order_statuses "hihand/enums/order"
	"hihand/pkgs/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID            string         `bson:"id" json:"id" form:"id" gorm:"primaryKey"`
	UserID        string         `bson:"user_id" json:"user_id" form:"user_id"`
	RecipientName string         `bson:"recipient_name" json:"recipient_name" validate:"required" form:"recipient_name" gorm:"index"`
	ContactPhone  string         `bson:"contact_phone" json:"contact_phone" validate:"required" form:"contact_phone" gorm:"index"`
	Email         string         `bson:"email" json:"email" validate:"required" form:"email" gorm:"index"`
	Address       string         `bson:"address" json:"address" validate:"required" form:"address"`
	OrderDate     time.Time      `bson:"order_date" json:"order_date" validate:"required" form:"order_date"`
	Status        string         `bson:"status" json:"status" validate:"required" form:"status"`
	Total         float64        `bson:"total" json:"total" validate:"required" form:"total"`
	Note          string         `bson:"note" json:"note" form:"note"`
	Details       []*OrderDetail `gorm:"foreignKey:OrderID" bson:"details" json:"details" validate:"dive" form:"details"`
}

type OrderDetail struct {
	ID        string  `bson:"id" json:"id" form:"id" gorm:"primaryKey"`
	OrderID   string  `bson:"order_id" json:"order_id" validate:"required" form:"order_id" gorm:"index;onDelete:CASCADE"`
	ProductID string  `bson:"product_id" json:"product_id" validate:"required" form:"product_id" gorm:"index"`
	Name      string  `bson:"name" json:"name" validate:"required" form:"name" gorm:"index"`
	Price     float64 `bson:"price" json:"price" validate:"required,gte=0" form:"price"`
	Quantity  int     `bson:"quantity" json:"quantity" validate:"required,gte=0" form:"quantity"`
	Total     float64 `bson:"total" json:"total" validate:"required,gte=0" form:"total"`
}

// ToMap converts an OrderDetail struct to a map
func (od *OrderDetail) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         od.ID,
		"order_id":   od.OrderID,
		"product_id": od.ProductID,
		"name":       od.Name,
		"price":      od.Price,
		"quantity":   od.Quantity,
		"total":      od.Total,
	}
}

// FromMap populates an OrderDetail struct from a map
func (od *OrderDetail) FromMap(data map[string]interface{}) {
	od.ID = utils.GetString(data, "id")
	od.OrderID = utils.GetString(data, "order_id")
	od.ProductID = utils.GetString(data, "product_id")
	od.Name = utils.GetString(data, "name")
	od.Price = utils.GetFloat64(data, "price")
	od.Quantity = utils.GetInt(data, "quantity")
	od.Total = utils.GetFloat64(data, "total")
}

// ToMap converts an Order struct to a map
func (o *Order) ToMap() (map[string]interface{}, error) {
	details := []map[string]interface{}{}
	for _, detail := range o.Details {
		detailMap := detail.ToMap()
		details = append(details, detailMap)
	}

	return map[string]interface{}{
		"id":             o.ID,
		"user_id":        o.UserID,
		"recipient_name": o.RecipientName,
		"contact_phone":  o.ContactPhone,
		"address":        o.Address,
		"email":          o.Email,
		"order_date":     o.OrderDate,
		"status":         o.Status,
		"total":          o.Total,
		"details":        details,
		"note":           o.Note,
		"created_at":     o.CreatedAt,
		"updated_at":     o.UpdatedAt,
	}, nil
}

// FromMap populates an Order struct from a map
func (o *Order) FromMap(data map[string]interface{}) error {
	o.ID = utils.GetString(data, "id")
	o.UserID = utils.GetString(data, "user_id")
	o.RecipientName = utils.GetString(data, "recipient_name")
	o.ContactPhone = utils.GetString(data, "contact_phone")
	o.Address = utils.GetString(data, "address")
	o.Email = utils.GetString(data, "email")
	o.OrderDate = utils.GetTime(data, "order_date")
	o.Status = utils.GetString(data, "status")
	o.Note = utils.GetString(data, "note")
	o.CreatedAt = utils.GetTime(data, "created_at")
	o.UpdatedAt = utils.GetTime(data, "updated_at")
	o.Total = utils.GetFloat64(data, "total")

	if details, ok := data["details"].([]interface{}); ok {
		for _, detailData := range details {
			if detailMap, ok := detailData.(map[string]interface{}); ok {
				detail := &OrderDetail{}
				detail.FromMap(detailMap)
				o.Details = append(o.Details, detail)
			}
		}
	}

	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New().String()
	o.Status = order_statuses.Created
	return
}

func (od *OrderDetail) BeforeCreate(tx *gorm.DB) (err error) {
	if od.ID == "" {
		od.ID = uuid.New().String()
	}
	return nil
}
