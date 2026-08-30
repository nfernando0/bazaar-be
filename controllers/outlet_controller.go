package controllers

import (
	"strconv"
	"strings"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"

	"github.com/gofiber/fiber/v3"
)

type OutletController struct{}

func NewOutletController() *OutletController {
	return &OutletController{}
}

func toOutletResponse(o *models.Outlet) dto.OutletResponse {
	var vendorsCount int64
	var usersCount int64
	config.DB.Model(&models.VendorOutlet{}).Where("outlet_id = ?", o.ID).Count(&vendorsCount)
	config.DB.Model(&models.User{}).Where("outlet_id = ?", o.ID).Count(&usersCount)

	bazaarName := ""
	if o.Bazaar != nil {
		bazaarName = o.Bazaar.Name
	}

	return dto.OutletResponse{
		ID:           o.ID,
		BazaarID:     o.BazaarID,
		BazaarName:   bazaarName,
		Name:         o.Name,
		Code:         o.Code,
		Location:     o.Location,
		VendorsCount: vendorsCount,
		UsersCount:   usersCount,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}

// GetAll returns all outlets with optional bazaar_id or search filtering
func (oc *OutletController) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	bazaarIDStr := strings.TrimSpace(c.Query("bazaar_id"))

	var outlets []models.Outlet
	db := config.DB.Model(&models.Outlet{}).Preload("Bazaar")

	if bazaarIDStr != "" {
		if bazaarID, err := strconv.Atoi(bazaarIDStr); err == nil {
			db = db.Where("bazaar_id = ?", bazaarID)
		}
	}

	if search != "" {
		db = db.Where("name LIKE ? OR code LIKE ? OR location LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := db.Order("name ASC").Find(&outlets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve outlets: " + err.Error(),
		})
	}

	responses := make([]dto.OutletResponse, 0, len(outlets))
	for _, o := range outlets {
		responses = append(responses, toOutletResponse(&o))
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Outlets retrieved successfully",
		Data:    responses,
	})
}

// GetByID returns a single outlet by ID
func (oc *OutletController) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid outlet ID",
		})
	}

	var outlet models.Outlet
	if err := config.DB.Preload("Bazaar").First(&outlet, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Outlet not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Outlet retrieved successfully",
		Data:    toOutletResponse(&outlet),
	})
}

// Create creates a new outlet under a bazaar
func (oc *OutletController) Create(c fiber.Ctx) error {
	var req dto.OutletRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	req.Location = strings.TrimSpace(req.Location)

	if req.BazaarID == 0 || req.Name == "" || req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "bazaar_id, name, and code are required",
		})
	}

	// Verify bazaar exists
	var bazaar models.Bazaar
	if err := config.DB.First(&bazaar, req.BazaarID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Specified bazaar_id not found",
		})
	}

	// Check unique code
	var existing models.Outlet
	if err := config.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(dto.APIResponse{
			Success: false,
			Message: "Outlet code already exists",
		})
	}

	outlet := models.Outlet{
		BazaarID: req.BazaarID,
		Name:     req.Name,
		Code:     req.Code,
		Location: req.Location,
	}

	if err := config.DB.Create(&outlet).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to create outlet: " + err.Error(),
		})
	}

	config.DB.Preload("Bazaar").First(&outlet, outlet.ID)

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Outlet created successfully",
		Data:    toOutletResponse(&outlet),
	})
}

// Update modifies an existing outlet
func (oc *OutletController) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid outlet ID",
		})
	}

	var req dto.OutletRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	var outlet models.Outlet
	if err := config.DB.First(&outlet, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Outlet not found",
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	req.Location = strings.TrimSpace(req.Location)

	if req.Name != "" {
		outlet.Name = req.Name
	}
	if req.Location != "" {
		outlet.Location = req.Location
	}

	if req.BazaarID != 0 && req.BazaarID != outlet.BazaarID {
		var bazaar models.Bazaar
		if err := config.DB.First(&bazaar, req.BazaarID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Specified bazaar_id not found",
			})
		}
		outlet.BazaarID = req.BazaarID
	}

	if req.Code != "" && req.Code != outlet.Code {
		var duplicate models.Outlet
		if err := config.DB.Where("code = ? AND id != ?", req.Code, id).First(&duplicate).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(dto.APIResponse{
				Success: false,
				Message: "Outlet code already exists",
			})
		}
		outlet.Code = req.Code
	}

	if err := config.DB.Save(&outlet).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to update outlet: " + err.Error(),
		})
	}

	config.DB.Preload("Bazaar").First(&outlet, outlet.ID)

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Outlet updated successfully",
		Data:    toOutletResponse(&outlet),
	})
}

// Delete removes an outlet
func (oc *OutletController) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid outlet ID",
		})
	}

	var outlet models.Outlet
	if err := config.DB.First(&outlet, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Outlet not found",
		})
	}

	// Check if transactions exist for this outlet
	var trxCount int64
	config.DB.Model(&models.Transaction{}).Where("outlet_id = ?", id).Count(&trxCount)
	if trxCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Cannot delete outlet: existing transactions are recorded for this outlet",
		})
	}

	if err := config.DB.Delete(&outlet).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to delete outlet: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Outlet deleted successfully",
	})
}
