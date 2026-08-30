package dto

import "time"

type VendorOutletRequest struct {
	VendorID    uint   `json:"vendor_id"`
	OutletID    uint   `json:"outlet_id"`
	BoothNumber string `json:"booth_number"`
}

type VendorOutletResponse struct {
	ID          uint      `json:"id"`
	VendorID    uint      `json:"vendor_id"`
	VendorName  string    `json:"vendor_name"`
	OutletID    uint      `json:"outlet_id"`
	OutletName  string    `json:"outlet_name"`
	OutletCode  string    `json:"outlet_code"`
	BoothNumber string    `json:"booth_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
