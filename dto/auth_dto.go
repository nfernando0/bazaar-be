package dto

import (
	"time"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`      // "admin" or "cashier" (default: "cashier")
	OutletID *uint  `json:"outlet_id,omitempty"` // optional for cashier / admin
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OutletInfo struct {
	ID       uint   `json:"id"`
	BazaarID uint   `json:"bazaar_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Location string `json:"location"`
}

type UserResponse struct {
	ID        uint        `json:"id"`
	OutletID  *uint       `json:"outlet_id,omitempty"`
	Outlet    *OutletInfo `json:"outlet,omitempty"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      string      `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
