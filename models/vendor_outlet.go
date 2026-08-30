package models

import "time"

type VendorOutlet struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VendorID    uint      `gorm:"not null;index:idx_vendor_outlet,unique" json:"vendor_id"`
	Vendor      *Vendor   `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	OutletID    uint      `gorm:"not null;index:idx_vendor_outlet,unique" json:"outlet_id"`
	Outlet      *Outlet   `gorm:"foreignKey:OutletID" json:"outlet,omitempty"`
	BoothNumber string    `gorm:"size:20" json:"booth_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
