package dto

import "time"

type BazaarRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"` // Format: YYYY-MM-DD or RFC3339
	EndDate     string `json:"end_date"`   // Format: YYYY-MM-DD or RFC3339
	Status      string `json:"status"`     // "draft", "active", "closed"
}

type BazaarOutletBrief struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Location string `json:"location"`
}

type BazaarResponse struct {
	ID           uint                `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	StartDate    time.Time           `json:"start_date"`
	EndDate      time.Time           `json:"end_date"`
	Status       string              `json:"status"`
	OutletsCount int64               `json:"outlets_count"`
	Outlets      []BazaarOutletBrief `json:"outlets,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}
