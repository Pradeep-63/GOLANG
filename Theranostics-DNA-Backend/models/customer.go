// models/customer.go

package models

import (
	"time"

	"gorm.io/gorm"
)

// Customer model
type Customer struct {
	ID            uint           `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	FirstName     string         `gorm:"type:varchar(50);not null" json:"first_name" validate:"required"`
	LastName      string         `gorm:"type:varchar(50);null" json:"last_name"`
	Email         string         `gorm:"type:varchar(100);not null;unique" json:"email" validate:"required,email"`
	PhoneNumber   string         `gorm:"type:varchar(15);not null" json:"phone_number" validate:"required"`
	Country       string         `gorm:"type:varchar(50);not null" json:"country" validate:"required"`
	StreetAddress string         `gorm:"type:varchar(255);not null" json:"street_address" validate:"required"`
	TownCity      string         `gorm:"type:varchar(100);not null" json:"town_city" validate:"required"`
	Region        string         `gorm:"type:varchar(100)" json:"region"`
	Postcode      string         `gorm:"type:varchar(20); not null" json:"postcode"`
	Orders        []Order        `gorm:"foreignKey:CustomerID" json:"customers,omitempty"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	IsDeleted     bool           `gorm:"default:false" json:"is_deleted"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
