package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OutletID     *uint          `gorm:"index" json:"outlet_id"`
	Outlet       *Outlet        `gorm:"foreignKey:OutletID" json:"outlet,omitempty"`
	Name         string         `gorm:"size:150;not null" json:"name"`
	Email        string         `gorm:"size:150;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Role         string         `gorm:"size:30;not null;default:cashier" json:"role"` // admin, cashier
	Transactions []Transaction  `gorm:"foreignKey:CashierID" json:"-"`
	Tokens       []UserToken    `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
