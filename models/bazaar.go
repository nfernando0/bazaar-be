package models

import (
	"time"

	"gorm.io/gorm"
)

type Bazaar struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `gorm:"not null" json:"end_date"`
	Status      string         `gorm:"size:20;not null;default:draft" json:"status"` // draft, active, closed
	Outlets     []Outlet       `gorm:"foreignKey:BazaarID" json:"outlets,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
