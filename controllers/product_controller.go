package controllers

import (
	"strconv"
	"strings"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"

	"github.com/gofiber/fiber/v3"
)

type ProductController struct{}

func NewProductController() *ProductController {
	return &ProductController{}
}

func toProductResponse(p *models.Product) dto.ProductResponse {
	vendorName := ""
	if p.Vendor != nil {
		vendorName = p.Vendor.Name
	}

	var categoryName string
	if p.Category != nil {
		categoryName = p.Category.Name
	}

	return dto.ProductResponse{
		ID:           p.ID,
		VendorID:     p.VendorID,
		VendorName:   vendorName,
		CategoryID:   p.CategoryID,
		CategoryName: categoryName,
		Name:         p.Name,
		Price:        p.Price,
		Stock:        p.Stock,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

// GetAll returns all products with filtering options
func (pc *ProductController) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	vendorIDStr := strings.TrimSpace(c.Query("vendor_id"))
	categoryIDStr := strings.TrimSpace(c.Query("category_id"))

	var products []models.Product
	db := config.DB.Model(&models.Product{}).Preload("Vendor").Preload("Category")

	if search != "" {
		db = db.Where("name LIKE ?", "%"+search+"%")
	}

	if vendorIDStr != "" {
		if vendorID, err := strconv.Atoi(vendorIDStr); err == nil {
			db = db.Where("vendor_id = ?", vendorID)
		}
	}

	if categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			db = db.Where("category_id = ?", categoryID)
		}
	}

	if err := db.Order("name ASC").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve products: " + err.Error(),
		})
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		responses = append(responses, toProductResponse(&p))
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    responses,
	})
}

// GetByID returns a single product by ID
func (pc *ProductController) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	var product models.Product
	if err := config.DB.Preload("Vendor").Preload("Category").First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Product not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    toProductResponse(&product),
	})
}

// Create adds a new product
func (pc *ProductController) Create(c fiber.Ctx) error {
	var req dto.ProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.VendorID == 0 || req.Name == "" || req.Price < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "vendor_id, name, and valid non-negative price are required",
		})
	}

	// Verify vendor exists
	var vendor models.Vendor
	if err := config.DB.First(&vendor, req.VendorID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Specified vendor_id not found",
		})
	}

	// Verify category if provided
	if req.CategoryID != nil && *req.CategoryID != 0 {
		var category models.Category
		if err := config.DB.First(&category, *req.CategoryID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Specified category_id not found",
			})
		}
	} else {
		req.CategoryID = nil
	}

	product := models.Product{
		VendorID:   req.VendorID,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to create product: " + err.Error(),
		})
	}

	config.DB.Preload("Vendor").Preload("Category").First(&product, product.ID)

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    toProductResponse(&product),
	})
}

// Update modifies an existing product
func (pc *ProductController) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	var req dto.ProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Product not found",
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name != "" {
		product.Name = req.Name
	}

	if req.Price >= 0 {
		product.Price = req.Price
	}

	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if req.VendorID != 0 && req.VendorID != product.VendorID {
		var vendor models.Vendor
		if err := config.DB.First(&vendor, req.VendorID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Specified vendor_id not found",
			})
		}
		product.VendorID = req.VendorID
	}

	if req.CategoryID != nil {
		if *req.CategoryID != 0 {
			var category models.Category
			if err := config.DB.First(&category, *req.CategoryID).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
					Success: false,
					Message: "Specified category_id not found",
				})
			}
			product.CategoryID = req.CategoryID
		} else {
			product.CategoryID = nil
		}
	}

	if err := config.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to update product: " + err.Error(),
		})
	}

	config.DB.Preload("Vendor").Preload("Category").First(&product, product.ID)

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Product updated successfully",
		Data:    toProductResponse(&product),
	})
}

// Delete removes a product (soft delete)
func (pc *ProductController) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Product not found",
		})
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to delete product: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}
