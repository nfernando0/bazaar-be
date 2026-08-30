package models

import (
	"time"

	"gorm.io/gorm"
)

type Outlet struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	BazaarID      uint           `gorm:"not null;index" json:"bazaar_id"`
	Bazaar        *Bazaar        `gorm:"foreignKey:BazaarID" json:"bazaar,omitempty"`
	Name          string         `gorm:"size:150;not null" json:"name"`
	Code          string         `gorm:"size:30;uniqueIndex;not null" json:"code"`
	Location      string         `gorm:"size:255" json:"location"`
	VendorOutlets []VendorOutlet `gorm:"foreignKey:OutletID" json:"vendor_outlets,omitempty"`
	Users         []User         `gorm:"foreignKey:OutletID" json:"users,omitempty"`
	Transactions  []Transaction  `gorm:"foreignKey:OutletID" json:"transactions,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
