package models

import "time"

type Transaction struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	OutletID        uint              `gorm:"not null;index" json:"outlet_id"`
	Outlet          *Outlet           `gorm:"foreignKey:OutletID" json:"outlet,omitempty"`
	CashierID       uint              `gorm:"not null;index" json:"cashier_id"`
	Cashier         *User             `gorm:"foreignKey:CashierID" json:"cashier,omitempty"`
	TransactionCode string            `gorm:"size:40;uniqueIndex;not null" json:"transaction_code"`
	TotalAmount     float64           `gorm:"type:decimal(14,2);not null" json:"total_amount"`
	PaymentMethod   string            `gorm:"size:30;not null" json:"payment_method"` // cash, qris, transfer, etc.
	Items           []TransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type TransactionItem struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	TransactionID uint         `gorm:"not null;index" json:"transaction_id"`
	Transaction   *Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
	ProductID     uint         `gorm:"not null;index" json:"product_id"`
	Product       *Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductName   string       `gorm:"size:255;not null" json:"product_name"`
	Quantity      int          `gorm:"not null" json:"quantity"`
	UnitPrice     float64      `gorm:"type:decimal(14,2);not null" json:"unit_price"`
	Subtotal      float64      `gorm:"type:decimal(14,2);not null" json:"subtotal"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
