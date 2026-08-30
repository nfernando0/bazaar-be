package dto

import "time"

type CategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	ProductsCount int64     `json:"products_count,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
