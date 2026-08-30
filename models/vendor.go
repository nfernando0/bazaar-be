package models

import (
	"time"

	"gorm.io/gorm"
)

type Vendor struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:150;not null" json:"name"`
	ContactPerson string         `gorm:"size:150" json:"contact_person"`
	Phone         string         `gorm:"size:30" json:"phone"`
	Products      []Product      `gorm:"foreignKey:VendorID" json:"products,omitempty"`
	VendorOutlets []VendorOutlet `gorm:"foreignKey:VendorID" json:"vendor_outlets,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
