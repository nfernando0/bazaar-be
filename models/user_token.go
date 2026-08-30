package models

import (
	"time"

	"gorm.io/gorm"
)

type UserToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Token     *string        `gorm:"type:text" json:"token"` // Becomes NULL when user logs out
	ExpiresAt *time.Time     `json:"expires_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
