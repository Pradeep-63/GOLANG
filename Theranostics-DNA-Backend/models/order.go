package models

import (
	"time"

	"gorm.io/gorm"
)

// Order model
type Order struct {
	ID                 uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	CustomerID         uint           `gorm:"not null" json:"customer_id" validate:"required"`
	Customer           Customer       `gorm:"foreignKey:CustomerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"customer,omitempty"`
	ProductName        string         `gorm:"type:varchar(100);not null" json:"product_name" validate:"required"`
	ProductDescription string         `gorm:"type:text" json:"product_description"`
	ProductImage       string         `gorm:"type:text" json:"product_image"`
	ProductPrice       float64        `gorm:"type:decimal(10,2);not null" json:"product_price" validate:"required,gt=0"`
	Quantity           int            `gorm:"type:int;not null" json:"quantity" validate:"required,min=1"`
	TotalPrice         float64        `gorm:"type:decimal(10,2);not null" json:"total_price" validate:"required,gt=0"`
	PaymentStatus      string         `gorm:"type:varchar(50);not null" json:"dna_payment_status" validate:"required,oneof=Pending Completed Failed"`
	OrderStatus        string         `gorm:"type:varchar(50);not null" json:"dna_order_status" validate:"required,oneof=Pending Shipped Delivered Cancelled"`
	Payments           []Payment      `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
	CreatedAt          time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsDeleted          bool           `gorm:"default:false" json:"is_deleted"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
