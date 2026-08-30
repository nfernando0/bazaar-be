package dto

import "time"

type OutletRequest struct {
	BazaarID uint   `json:"bazaar_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Location string `json:"location"`
}

type OutletResponse struct {
	ID           uint      `json:"id"`
	BazaarID     uint      `json:"bazaar_id"`
	BazaarName   string    `json:"bazaar_name,omitempty"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Location     string    `json:"location"`
	VendorsCount int64     `json:"vendors_count"`
	UsersCount   int64     `json:"users_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
