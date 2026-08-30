package dto

import "time"

type TransactionItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type TransactionRequest struct {
	OutletID      *uint                    `json:"outlet_id,omitempty"` // If omitted, uses cashier's assigned outlet
	PaymentMethod string                   `json:"payment_method"`      // "cash", "qris", "transfer", etc.
	Items         []TransactionItemRequest `json:"items"`
}

type TransactionItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

type TransactionResponse struct {
	ID              uint                      `json:"id"`
	OutletID        uint                      `json:"outlet_id"`
	OutletName      string                    `json:"outlet_name,omitempty"`
	CashierID       uint                      `json:"cashier_id"`
	CashierName     string                    `json:"cashier_name,omitempty"`
	TransactionCode string                    `json:"transaction_code"`
	TotalAmount     float64                   `json:"total_amount"`
	PaymentMethod   string                    `json:"payment_method"`
	Items           []TransactionItemResponse `json:"items,omitempty"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

type TransactionSummaryResponse struct {
	TotalTransactions int64   `json:"total_transactions"`
	TotalRevenue      float64 `json:"total_revenue"`
}
