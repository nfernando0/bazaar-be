package dto

import "time"

type TransactionItemDetailResponse struct {
	ID              uint      `json:"id"`
	TransactionID   uint      `json:"transaction_id"`
	TransactionCode string    `json:"transaction_code,omitempty"`
	OutletID        uint      `json:"outlet_id,omitempty"`
	OutletName      string    `json:"outlet_name,omitempty"`
	ProductID       uint      `json:"product_id"`
	ProductName     string    `json:"product_name"`
	VendorID        uint      `json:"vendor_id,omitempty"`
	VendorName      string    `json:"vendor_name,omitempty"`
	Quantity        int       `json:"quantity"`
	UnitPrice       float64   `json:"unit_price"`
	Subtotal        float64   `json:"subtotal"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TopSellingProductResponse struct {
	ProductID         uint    `json:"product_id"`
	ProductName       string  `json:"product_name"`
	TotalQuantitySold int64   `json:"total_quantity_sold"`
	TotalRevenue      float64 `json:"total_revenue"`
}
