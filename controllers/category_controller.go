package controllers

import (
	"strconv"
	"strings"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"

	"github.com/gofiber/fiber/v3"
)

type CategoryController struct{}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

// GetAll returns all categories with optional name filtering
func (cc *CategoryController) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	var categories []models.Category
	db := config.DB.Model(&models.Category{})

	if search != "" {
		db = db.Where("name LIKE ?", "%"+search+"%")
	}

	if err := db.Order("name ASC").Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve categories: " + err.Error(),
		})
	}

	responses := make([]dto.CategoryResponse, 0, len(categories))
	for _, cat := range categories {
		var count int64
		config.DB.Model(&models.Product{}).Where("category_id = ?", cat.ID).Count(&count)

		responses = append(responses, dto.CategoryResponse{
			ID:            cat.ID,
			Name:          cat.Name,
			ProductsCount: count,
			CreatedAt:     cat.CreatedAt,
			UpdatedAt:     cat.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Categories retrieved successfully",
		Data:    responses,
	})
}

// GetByID returns a single category by ID
func (cc *CategoryController) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid category ID",
		})
	}

	var category models.Category
	if err := config.DB.Preload("Products").First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Category not found",
		})
	}

	var count int64
	config.DB.Model(&models.Product{}).Where("category_id = ?", category.ID).Count(&count)

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Category retrieved successfully",
		Data: dto.CategoryResponse{
			ID:            category.ID,
			Name:          category.Name,
			ProductsCount: count,
			CreatedAt:     category.CreatedAt,
			UpdatedAt:     category.UpdatedAt,
		},
	})
}

// Create creates a new category
func (cc *CategoryController) Create(c fiber.Ctx) error {
	var req dto.CategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Category name is required",
		})
	}

	// Check if category name already exists
	var existing models.Category
	if err := config.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(dto.APIResponse{
			Success: false,
			Message: "Category name already exists",
		})
	}

	category := models.Category{
		Name: req.Name,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to create category: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Category created successfully",
		Data: dto.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		},
	})
}

// Update updates an existing category
func (cc *CategoryController) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid category ID",
		})
	}

	var req dto.CategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Category name is required",
		})
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Category not found",
		})
	}

	// Check unique name on other records
	var duplicate models.Category
	if err := config.DB.Where("name = ? AND id != ?", req.Name, id).First(&duplicate).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(dto.APIResponse{
			Success: false,
			Message: "Category name already exists",
		})
	}

	category.Name = req.Name
	if err := config.DB.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to update category: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Category updated successfully",
		Data: dto.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		},
	})
}

// Delete removes a category (soft delete)
func (cc *CategoryController) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid category ID",
		})
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Category not found",
		})
	}

	// Check if any product is using this category
	var productCount int64
	config.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
	if productCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Cannot delete category: products are currently assigned to this category",
		})
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to delete category: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}
