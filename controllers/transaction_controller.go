package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"

	"github.com/gofiber/fiber/v3"
)

type TransactionController struct{}

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

// Generate unique transaction code: TRX-YYYYMMDD-HHMMSS-XXXX
func generateTransactionCode() string {
	now := time.Now().Format("20060102-150405")
	n, _ := rand.Int(rand.Reader, big.NewInt(9000))
	randPart := n.Int64() + 1000
	return fmt.Sprintf("TRX-%s-%d", now, randPart)
}

func toTransactionResponse(t *models.Transaction) dto.TransactionResponse {
	outletName := ""
	if t.Outlet != nil {
		outletName = t.Outlet.Name
	}

	cashierName := ""
	if t.Cashier != nil {
		cashierName = t.Cashier.Name
	}

	items := make([]dto.TransactionItemResponse, 0, len(t.Items))
	for _, item := range t.Items {
		items = append(items, dto.TransactionItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
		})
	}

	return dto.TransactionResponse{
		ID:              t.ID,
		OutletID:        t.OutletID,
		OutletName:      outletName,
		CashierID:       t.CashierID,
		CashierName:     cashierName,
		TransactionCode: t.TransactionCode,
		TotalAmount:     t.TotalAmount,
		PaymentMethod:   t.PaymentMethod,
		Items:           items,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

// Create processes a new cashier checkout transaction
func (tc *TransactionController) Create(c fiber.Ctx) error {
	cashierIDVal := c.Locals("user_id")
	cashierID, ok := cashierIDVal.(uint)
	if !ok || cashierID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Success: false,
			Message: "Unauthorized cashier session",
		})
	}

	var req dto.TransactionRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Transaction must contain at least one item",
		})
	}

	req.PaymentMethod = strings.TrimSpace(strings.ToLower(req.PaymentMethod))
	if req.PaymentMethod == "" {
		req.PaymentMethod = "cash"
	}

	// Determine Outlet ID
	var targetOutletID uint
	if req.OutletID != nil && *req.OutletID != 0 {
		targetOutletID = *req.OutletID
	} else if outletIDLocVal := c.Locals("outlet_id"); outletIDLocVal != nil {
		if ptr, ok := outletIDLocVal.(*uint); ok && ptr != nil {
			targetOutletID = *ptr
		}
	}

	// Fallback to query cashier's assigned outlet
	if targetOutletID == 0 {
		var cashier models.User
		if err := config.DB.First(&cashier, cashierID).Error; err == nil && cashier.OutletID != nil {
			targetOutletID = *cashier.OutletID
		}
	}

	if targetOutletID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Outlet ID could not be determined. Please specify outlet_id or assign cashier to an outlet.",
		})
	}

	// Verify outlet exists
	var outlet models.Outlet
	if err := config.DB.First(&outlet, targetOutletID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Target outlet not found",
		})
	}

	// Begin DB Transaction for atomic inventory deduction and checkout
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var totalAmount float64
	var transactionItems []models.TransactionItem

	for _, itemReq := range req.Items {
		if itemReq.ProductID == 0 || itemReq.Quantity <= 0 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Invalid product_id or quantity",
			})
		}

		var product models.Product
		if err := tx.First(&product, itemReq.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: fmt.Sprintf("Product ID %d not found", itemReq.ProductID),
			})
		}

		// Check stock
		if product.Stock < itemReq.Quantity {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: fmt.Sprintf("Insufficient stock for product '%s'. Available: %d, Requested: %d",
					product.Name, product.Stock, itemReq.Quantity),
			})
		}

		// Deduct stock
		product.Stock -= itemReq.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
				Success: false,
				Message: "Failed to update product stock: " + err.Error(),
			})
		}

		subtotal := float64(itemReq.Quantity) * product.Price
		totalAmount += subtotal

		transactionItems = append(transactionItems, models.TransactionItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    itemReq.Quantity,
			UnitPrice:   product.Price,
			Subtotal:    subtotal,
		})
	}

	transaction := models.Transaction{
		OutletID:        targetOutletID,
		CashierID:       cashierID,
		TransactionCode: generateTransactionCode(),
		TotalAmount:     totalAmount,
		PaymentMethod:   req.PaymentMethod,
		Items:           transactionItems,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to record transaction: " + err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to commit transaction: " + err.Error(),
		})
	}

	// Preload relationships for response
	config.DB.Preload("Outlet").Preload("Cashier").Preload("Items").First(&transaction, transaction.ID)

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Transaction completed successfully",
		Data:    toTransactionResponse(&transaction),
	})
}

// GetAll returns transactions with filter support
func (tc *TransactionController) GetAll(c fiber.Ctx) error {
	outletIDStr := strings.TrimSpace(c.Query("outlet_id"))
	cashierIDStr := strings.TrimSpace(c.Query("cashier_id"))
	paymentMethod := strings.TrimSpace(c.Query("payment_method"))

	var transactions []models.Transaction
	db := config.DB.Model(&models.Transaction{}).
		Preload("Outlet").
		Preload("Cashier").
		Preload("Items")

	if outletIDStr != "" {
		if outletID, err := strconv.Atoi(outletIDStr); err == nil {
			db = db.Where("outlet_id = ?", outletID)
		}
	}

	if cashierIDStr != "" {
		if cashierID, err := strconv.Atoi(cashierIDStr); err == nil {
			db = db.Where("cashier_id = ?", cashierID)
		}
	}

	if paymentMethod != "" {
		db = db.Where("payment_method = ?", strings.ToLower(paymentMethod))
	}

	if err := db.Order("created_at DESC").Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve transactions: " + err.Error(),
		})
	}

	responses := make([]dto.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		responses = append(responses, toTransactionResponse(&t))
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Transactions retrieved successfully",
		Data:    responses,
	})
}

// GetByID returns a single transaction receipt by ID
func (tc *TransactionController) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid transaction ID",
		})
	}

	var transaction models.Transaction
	if err := config.DB.Preload("Outlet").Preload("Cashier").Preload("Items").First(&transaction, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Transaction not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Transaction retrieved successfully",
		Data:    toTransactionResponse(&transaction),
	})
}

// GetSummary returns total transaction count and gross revenue
func (tc *TransactionController) GetSummary(c fiber.Ctx) error {
	outletIDStr := strings.TrimSpace(c.Query("outlet_id"))

	db := config.DB.Model(&models.Transaction{})
	if outletIDStr != "" {
		if outletID, err := strconv.Atoi(outletIDStr); err == nil {
			db = db.Where("outlet_id = ?", outletID)
		}
	}

	var count int64
	db.Count(&count)

	var totalRevenue float64
	db.Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Transaction summary retrieved successfully",
		Data: dto.TransactionSummaryResponse{
			TotalTransactions: count,
			TotalRevenue:      totalRevenue,
		},
	})
}
