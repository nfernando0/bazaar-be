package dto

import "time"

type VendorRequest struct {
	Name          string `json:"name"`
	ContactPerson string `json:"contact_person"`
	Phone         string `json:"phone"`
}

type VendorOutletBrief struct {
	ID          uint   `json:"id"`
	OutletID    uint   `json:"outlet_id"`
	OutletName  string `json:"outlet_name"`
	OutletCode  string `json:"outlet_code"`
	BoothNumber string `json:"booth_number"`
}

type VendorResponse struct {
	ID            uint                `json:"id"`
	Name          string              `json:"name"`
	ContactPerson string              `json:"contact_person"`
	Phone         string              `json:"phone"`
	Outlets       []VendorOutletBrief `json:"outlets,omitempty"`
	ProductsCount int64               `json:"products_count"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}
