package controllers

import (
	"strconv"
	"strings"
	"time"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"

	"github.com/gofiber/fiber/v3"
)

type BazaarController struct{}

func NewBazaarController() *BazaarController {
	return &BazaarController{}
}

// Helper to parse flexible date string
func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	// Try standard YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}
	// Try RFC3339
	return time.Parse(time.RFC3339, dateStr)
}

func toBazaarResponse(b *models.Bazaar) dto.BazaarResponse {
	var outletCount int64
	config.DB.Model(&models.Outlet{}).Where("bazaar_id = ?", b.ID).Count(&outletCount)

	outlets := make([]dto.BazaarOutletBrief, 0, len(b.Outlets))
	for _, o := range b.Outlets {
		outlets = append(outlets, dto.BazaarOutletBrief{
			ID:       o.ID,
			Name:     o.Name,
			Code:     o.Code,
			Location: o.Location,
		})
	}

	return dto.BazaarResponse{
		ID:           b.ID,
		Name:         b.Name,
		Description:  b.Description,
		StartDate:    b.StartDate,
		EndDate:      b.EndDate,
		Status:       b.Status,
		OutletsCount: outletCount,
		Outlets:      outlets,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}

// GetAll returns all bazaars
func (bc *BazaarController) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))

	var bazaars []models.Bazaar
	db := config.DB.Model(&models.Bazaar{}).Preload("Outlets")

	if search != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Order("start_date DESC").Find(&bazaars).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve bazaars: " + err.Error(),
		})
	}

	responses := make([]dto.BazaarResponse, 0, len(bazaars))
	for _, b := range bazaars {
		responses = append(responses, toBazaarResponse(&b))
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Bazaars retrieved successfully",
		Data:    responses,
	})
}

// GetByID returns a single bazaar by ID
func (bc *BazaarController) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid bazaar ID",
		})
	}

	var bazaar models.Bazaar
	if err := config.DB.Preload("Outlets").First(&bazaar, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Bazaar not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Bazaar retrieved successfully",
		Data:    toBazaarResponse(&bazaar),
	})
}

// Create registers a new bazaar event
func (bc *BazaarController) Create(c fiber.Ctx) error {
	var req dto.BazaarRequest
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
			Message: "Bazaar name is required",
		})
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid start_date format (use YYYY-MM-DD)",
		})
	}

	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid end_date format (use YYYY-MM-DD)",
		})
	}

	if endDate.Before(startDate) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "end_date cannot be earlier than start_date",
		})
	}

	status := strings.TrimSpace(strings.ToLower(req.Status))
	if status == "" {
		status = "draft"
	} else if status != "draft" && status != "active" && status != "closed" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Status must be 'draft', 'active', or 'closed'",
		})
	}

	bazaar := models.Bazaar{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      status,
	}

	if err := config.DB.Create(&bazaar).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to create bazaar: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Bazaar created successfully",
		Data:    toBazaarResponse(&bazaar),
	})
}

// Update modifies an existing bazaar
func (bc *BazaarController) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid bazaar ID",
		})
	}

	var req dto.BazaarRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	var bazaar models.Bazaar
	if err := config.DB.First(&bazaar, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Bazaar not found",
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name != "" {
		bazaar.Name = req.Name
	}
	bazaar.Description = req.Description

	if req.StartDate != "" {
		startDate, err := parseDate(req.StartDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Invalid start_date format (use YYYY-MM-DD)",
			})
		}
		bazaar.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := parseDate(req.EndDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Invalid end_date format (use YYYY-MM-DD)",
			})
		}
		bazaar.EndDate = endDate
	}

	if bazaar.EndDate.Before(bazaar.StartDate) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "end_date cannot be earlier than start_date",
		})
	}

	if req.Status != "" {
		status := strings.TrimSpace(strings.ToLower(req.Status))
		if status != "draft" && status != "active" && status != "closed" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Status must be 'draft', 'active', or 'closed'",
			})
		}
		bazaar.Status = status
	}

	if err := config.DB.Save(&bazaar).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to update bazaar: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Bazaar updated successfully",
		Data:    toBazaarResponse(&bazaar),
	})
}

// Delete removes a bazaar
func (bc *BazaarController) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid bazaar ID",
		})
	}

	var bazaar models.Bazaar
	if err := config.DB.First(&bazaar, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Bazaar not found",
		})
	}

	if err := config.DB.Delete(&bazaar).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to delete bazaar: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Bazaar deleted successfully",
	})
}
