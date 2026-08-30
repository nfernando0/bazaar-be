package controllers

import (
	"strconv"
	"strings"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"

	"github.com/gofiber/fiber/v3"
)

type VendorOutletController struct{}

func NewVendorOutletController() *VendorOutletController {
	return &VendorOutletController{}
}

// Helper to format VendorOutlet model to VendorOutletResponse DTO
func toVendorOutletResponse(vo *models.VendorOutlet) dto.VendorOutletResponse {
	vendorName := ""
	if vo.Vendor != nil {
		vendorName = vo.Vendor.Name
	}

	outletName := ""
	outletCode := ""
	if vo.Outlet != nil {
		outletName = vo.Outlet.Name
		outletCode = vo.Outlet.Code
	}

	return dto.VendorOutletResponse{
		ID:          vo.ID,
		VendorID:    vo.VendorID,
		VendorName:  vendorName,
		OutletID:    vo.OutletID,
		OutletName:  outletName,
		OutletCode:  outletCode,
		BoothNumber: vo.BoothNumber,
		CreatedAt:   vo.CreatedAt,
		UpdatedAt:   vo.UpdatedAt,
	}
}

// GetAll returns all vendor outlet assignments with optional filtering by vendor_id or outlet_id
func (voc *VendorOutletController) GetAll(c fiber.Ctx) error {
	var vendorOutlets []models.VendorOutlet
	db := config.DB.Model(&models.VendorOutlet{}).Preload("Vendor").Preload("Outlet")

	if vendorIDStr := c.Query("vendor_id"); vendorIDStr != "" {
		if vendorID, err := strconv.Atoi(vendorIDStr); err == nil {
			db = db.Where("vendor_id = ?", vendorID)
		}
	}

	if outletIDStr := c.Query("outlet_id"); outletIDStr != "" {
		if outletID, err := strconv.Atoi(outletIDStr); err == nil {
			db = db.Where("outlet_id = ?", outletID)
		}
	}

	if err := db.Order("id ASC").Find(&vendorOutlets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve vendor outlets: " + err.Error(),
		})
	}

	responses := make([]dto.VendorOutletResponse, 0, len(vendorOutlets))
	for _, vo := range vendorOutlets {
		responses = append(responses, toVendorOutletResponse(&vo))
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor outlet assignments retrieved successfully",
		Data:    responses,
	})
}

// GetByID returns a single vendor outlet assignment by ID
func (voc *VendorOutletController) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid assignment ID",
		})
	}

	var vo models.VendorOutlet
	if err := config.DB.Preload("Vendor").Preload("Outlet").First(&vo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor outlet assignment not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor outlet assignment retrieved successfully",
		Data:    toVendorOutletResponse(&vo),
	})
}

// Assign creates a new vendor to outlet assignment
func (voc *VendorOutletController) Assign(c fiber.Ctx) error {
	var req dto.VendorOutletRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.BoothNumber = strings.TrimSpace(req.BoothNumber)

	if req.VendorID == 0 || req.OutletID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Both vendor_id and outlet_id are required",
		})
	}

	// Verify vendor exists
	var vendor models.Vendor
	if err := config.DB.First(&vendor, req.VendorID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor not found",
		})
	}

	// Verify outlet exists
	var outlet models.Outlet
	if err := config.DB.First(&outlet, req.OutletID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Outlet not found",
		})
	}

	// Check if already assigned
	var existing models.VendorOutlet
	if err := config.DB.Where("vendor_id = ? AND outlet_id = ?", req.VendorID, req.OutletID).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor is already assigned to this outlet",
		})
	}

	vo := models.VendorOutlet{
		VendorID:    req.VendorID,
		OutletID:    req.OutletID,
		BoothNumber: req.BoothNumber,
	}

	if err := config.DB.Create(&vo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to assign vendor to outlet: " + err.Error(),
		})
	}

	// Preload for response
	config.DB.Preload("Vendor").Preload("Outlet").First(&vo, vo.ID)

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor assigned to outlet successfully",
		Data:    toVendorOutletResponse(&vo),
	})
}

// Update modifies booth number or assignment details
func (voc *VendorOutletController) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid assignment ID",
		})
	}

	var req dto.VendorOutletRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.BoothNumber = strings.TrimSpace(req.BoothNumber)

	var vo models.VendorOutlet
	if err := config.DB.First(&vo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor outlet assignment not found",
		})
	}

	// If changing vendor or outlet
	if req.VendorID != 0 && req.VendorID != vo.VendorID {
		var vendor models.Vendor
		if err := config.DB.First(&vendor, req.VendorID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Specified vendor not found",
			})
		}
		vo.VendorID = req.VendorID
	}

	if req.OutletID != 0 && req.OutletID != vo.OutletID {
		var outlet models.Outlet
		if err := config.DB.First(&outlet, req.OutletID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Specified outlet not found",
			})
		}
		vo.OutletID = req.OutletID
	}

	// Check unique combination if changed
	var duplicate models.VendorOutlet
	if err := config.DB.Where("vendor_id = ? AND outlet_id = ? AND id != ?", vo.VendorID, vo.OutletID, id).First(&duplicate).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor is already assigned to this outlet",
		})
	}

	vo.BoothNumber = req.BoothNumber

	if err := config.DB.Save(&vo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to update assignment: " + err.Error(),
		})
	}

	// Preload for response
	config.DB.Preload("Vendor").Preload("Outlet").First(&vo, vo.ID)

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor outlet assignment updated successfully",
		Data:    toVendorOutletResponse(&vo),
	})
}

// Delete removes vendor assignment from outlet
func (voc *VendorOutletController) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid assignment ID",
		})
	}

	var vo models.VendorOutlet
	if err := config.DB.First(&vo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor outlet assignment not found",
		})
	}

	if err := config.DB.Delete(&vo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to delete assignment: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor outlet assignment deleted successfully",
	})
}
