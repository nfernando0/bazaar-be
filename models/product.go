package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	VendorID   uint              `gorm:"not null;index" json:"vendor_id"`
	Vendor     *Vendor           `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	CategoryID *uint             `gorm:"index" json:"category_id"`
	Category   *Category         `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name       string            `gorm:"size:150;not null" json:"name"`
	Price      float64           `gorm:"type:decimal(12,2);not null" json:"price"`
	Stock      int               `gorm:"not null;default:0" json:"stock"`
	Items      []TransactionItem `gorm:"foreignKey:ProductID" json:"-"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  gorm.DeletedAt    `gorm:"index" json:"-"`
}
