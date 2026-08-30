package dto

import "time"

type ProductRequest struct {
	VendorID   uint    `json:"vendor_id"`
	CategoryID *uint   `json:"category_id,omitempty"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
}

type ProductResponse struct {
	ID           uint      `json:"id"`
	VendorID     uint      `json:"vendor_id"`
	VendorName   string    `json:"vendor_name,omitempty"`
	CategoryID   *uint     `json:"category_id,omitempty"`
	CategoryName string    `json:"category_name,omitempty"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Stock        int       `json:"stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
