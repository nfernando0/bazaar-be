package controllers

import (
	"bazar-be/config"
	"bazar-be/dto"

	"github.com/gofiber/fiber/v3"
)

type DevController struct{}

func NewDevController() *DevController {
	return &DevController{}
}

type SeedRequest struct {
	Fresh bool `json:"fresh"`
}

// Seed populates the database with comprehensive dummy data
func (dc *DevController) Seed(c fiber.Ctx) error {
	var req SeedRequest
	_ = c.Bind().Body(&req)

	// Also check query param ?fresh=true
	if c.Query("fresh") == "true" || c.Query("fresh") == "1" {
		req.Fresh = true
	}

	if err := config.SeedAll(config.DB, req.Fresh); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to seed database: " + err.Error(),
		})
	}

	mode := "Safe update mode"
	if req.Fresh {
		mode = "Fresh reset mode (all tables truncated and re-populated)"
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Database successfully seeded with comprehensive dummy data! (" + mode + ")",
	})
}
