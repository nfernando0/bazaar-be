package controllers

import (
	"strconv"
	"strings"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"

	"github.com/gofiber/fiber/v3"
)

type VendorController struct{}

func NewVendorController() *VendorController {
	return &VendorController{}
}

// Helper to format Vendor model to VendorResponse DTO
func toVendorResponse(vendor *models.Vendor) dto.VendorResponse {
	var productCount int64
	config.DB.Model(&models.Product{}).Where("vendor_id = ?", vendor.ID).Count(&productCount)

	outlets := make([]dto.VendorOutletBrief, 0, len(vendor.VendorOutlets))
	for _, vo := range vendor.VendorOutlets {
		outletName := ""
		outletCode := ""
		if vo.Outlet != nil {
			outletName = vo.Outlet.Name
			outletCode = vo.Outlet.Code
		}
		outlets = append(outlets, dto.VendorOutletBrief{
			ID:          vo.ID,
			OutletID:    vo.OutletID,
			OutletName:  outletName,
			OutletCode:  outletCode,
			BoothNumber: vo.BoothNumber,
		})
	}

	return dto.VendorResponse{
		ID:            vendor.ID,
		Name:          vendor.Name,
		ContactPerson: vendor.ContactPerson,
		Phone:         vendor.Phone,
		Outlets:       outlets,
		ProductsCount: productCount,
		CreatedAt:     vendor.CreatedAt,
		UpdatedAt:     vendor.UpdatedAt,
	}
}

// GetAll returns all vendors with optional search query
func (vc *VendorController) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	var vendors []models.Vendor
	db := config.DB.Model(&models.Vendor{}).Preload("VendorOutlets.Outlet")

	if search != "" {
		db = db.Where("name LIKE ? OR contact_person LIKE ? OR phone LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := db.Order("name ASC").Find(&vendors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve vendors: " + err.Error(),
		})
	}

	responses := make([]dto.VendorResponse, 0, len(vendors))
	for _, vendor := range vendors {
		responses = append(responses, toVendorResponse(&vendor))
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendors retrieved successfully",
		Data:    responses,
	})
}

// GetByID returns a single vendor with their outlet assignments and products
func (vc *VendorController) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid vendor ID",
		})
	}

	var vendor models.Vendor
	if err := config.DB.Preload("VendorOutlets.Outlet").First(&vendor, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor retrieved successfully",
		Data:    toVendorResponse(&vendor),
	})
}

// Create registers a new vendor
func (vc *VendorController) Create(c fiber.Ctx) error {
	var req dto.VendorRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.ContactPerson = strings.TrimSpace(req.ContactPerson)
	req.Phone = strings.TrimSpace(req.Phone)

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor name is required",
		})
	}

	vendor := models.Vendor{
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
	}

	if err := config.DB.Create(&vendor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to create vendor: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor created successfully",
		Data:    toVendorResponse(&vendor),
	})
}

// Update modifies an existing vendor's information
func (vc *VendorController) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid vendor ID",
		})
	}

	var req dto.VendorRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.ContactPerson = strings.TrimSpace(req.ContactPerson)
	req.Phone = strings.TrimSpace(req.Phone)

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor name is required",
		})
	}

	var vendor models.Vendor
	if err := config.DB.Preload("VendorOutlets.Outlet").First(&vendor, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor not found",
		})
	}

	vendor.Name = req.Name
	vendor.ContactPerson = req.ContactPerson
	vendor.Phone = req.Phone

	if err := config.DB.Save(&vendor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to update vendor: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor updated successfully",
		Data:    toVendorResponse(&vendor),
	})
}

// Delete removes a vendor (soft delete)
func (vc *VendorController) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid vendor ID",
		})
	}

	var vendor models.Vendor
	if err := config.DB.First(&vendor, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "Vendor not found",
		})
	}

	// Check if vendor has products attached
	var productCount int64
	config.DB.Model(&models.Product{}).Where("vendor_id = ?", id).Count(&productCount)
	if productCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Cannot delete vendor: products are still assigned to this vendor",
		})
	}

	if err := config.DB.Delete(&vendor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to delete vendor: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Vendor deleted successfully",
	})
}
