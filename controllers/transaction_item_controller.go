package controllers

import (
	"strconv"
	"strings"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"

	"github.com/gofiber/fiber/v3"
)

type TransactionItemController struct{}

func NewTransactionItemController() *TransactionItemController {
	return &TransactionItemController{}
}

func toTransactionItemDetailResponse(item *models.TransactionItem) dto.TransactionItemDetailResponse {
	trxCode := ""
	var outletID uint
	outletName := ""

	if item.Transaction != nil {
		trxCode = item.Transaction.TransactionCode
		outletID = item.Transaction.OutletID
		if item.Transaction.Outlet != nil {
			outletName = item.Transaction.Outlet.Name
		}
	}

	var vendorID uint
	vendorName := ""
	if item.Product != nil && item.Product.Vendor != nil {
		vendorID = item.Product.VendorID
		vendorName = item.Product.Vendor.Name
	}

	return dto.TransactionItemDetailResponse{
		ID:              item.ID,
		TransactionID:   item.TransactionID,
		TransactionCode: trxCode,
		OutletID:        outletID,
		OutletName:      outletName,
		ProductID:       item.ProductID,
		ProductName:     item.ProductName,
		VendorID:        vendorID,
		VendorName:      vendorName,
		Quantity:        item.Quantity,
		UnitPrice:       item.UnitPrice,
		Subtotal:        item.Subtotal,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

// GetAll returns transaction items with optional filtering by transaction_id, product_id, outlet_id, or vendor_id
func (tic *TransactionItemController) GetAll(c fiber.Ctx) error {
	trxIDStr := strings.TrimSpace(c.Query("transaction_id"))
	productIDStr := strings.TrimSpace(c.Query("product_id"))
	outletIDStr := strings.TrimSpace(c.Query("outlet_id"))
	vendorIDStr := strings.TrimSpace(c.Query("vendor_id"))

	var items []models.TransactionItem
	db := config.DB.Model(&models.TransactionItem{}).
		Preload("Transaction.Outlet").
		Preload("Product.Vendor")

	if trxIDStr != "" {
		if trxID, err := strconv.Atoi(trxIDStr); err == nil {
			db = db.Where("transaction_id = ?", trxID)
		}
	}

	if productIDStr != "" {
		if productID, err := strconv.Atoi(productIDStr); err == nil {
			db = db.Where("product_id = ?", productID)
		}
	}

	if outletIDStr != "" {
		if outletID, err := strconv.Atoi(outletIDStr); err == nil {
			db = db.Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
				Where("transactions.outlet_id = ?", outletID)
		}
	}

	if vendorIDStr != "" {
		if vendorID, err := strconv.Atoi(vendorIDStr); err == nil {
			db = db.Joins("JOIN products ON products.id = transaction_items.product_id").
				Where("products.vendor_id = ?", vendorID)
		}
	}

	if err := db.Order("created_at DESC").Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve transaction items: " + err.Error(),
		})
	}

	responses := make([]dto.TransactionItemDetailResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toTransactionItemDetailResponse(&item))
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Transaction items retrieved successfully",
		Data:    responses,
	})
}

// GetByID returns a single transaction item detail by ID
func (tic *TransactionItemController) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid transaction item ID",
		})
	}

	var item models.TransactionItem
	if err := config.DB.
		Preload("Transaction.Outlet").
		Preload("Product.Vendor").
		First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Transaction item not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Transaction item retrieved successfully",
		Data:    toTransactionItemDetailResponse(&item),
	})
}

// GetTopSelling returns top selling products ranked by sales quantity
func (tic *TransactionItemController) GetTopSelling(c fiber.Ctx) error {
	limitStr := strings.TrimSpace(c.Query("limit"))
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	type AggregateResult struct {
		ProductID         uint    `gorm:"column:product_id"`
		ProductName       string  `gorm:"column:product_name"`
		TotalQuantitySold int64   `gorm:"column:total_quantity_sold"`
		TotalRevenue      float64 `gorm:"column:total_revenue"`
	}

	var results []AggregateResult
	err := config.DB.Model(&models.TransactionItem{}).
		Select("product_id, product_name, SUM(quantity) as total_quantity_sold, SUM(subtotal) as total_revenue").
		Group("product_id, product_name").
		Order("total_quantity_sold DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve top selling products: " + err.Error(),
		})
	}

	responses := make([]dto.TopSellingProductResponse, 0, len(results))
	for _, r := range results {
		responses = append(responses, dto.TopSellingProductResponse{
			ProductID:         r.ProductID,
			ProductName:       r.ProductName,
			TotalQuantitySold: r.TotalQuantitySold,
			TotalRevenue:      r.TotalRevenue,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Top selling products retrieved successfully",
		Data:    responses,
	})
}
