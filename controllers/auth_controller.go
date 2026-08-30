package controllers

import (
	"strings"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"
	"bazar-be/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// Helper to format User model into UserResponse DTO
func toUserResponse(user *models.User) dto.UserResponse {
	res := dto.UserResponse{
		ID:        user.ID,
		OutletID:  user.OutletID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.Outlet != nil {
		res.Outlet = &dto.OutletInfo{
			ID:       user.Outlet.ID,
			BazaarID: user.Outlet.BazaarID,
			Name:     user.Outlet.Name,
			Code:     user.Outlet.Code,
			Location: user.Outlet.Location,
		}
	}

	return res
}

// Register handles new user registration
func (ac *AuthController) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	// Basic validation
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)
	req.Role = strings.TrimSpace(strings.ToLower(req.Role))

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Name, email, and password are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Password must be at least 6 characters long",
		})
	}

	// Default role to "cashier" if empty
	if req.Role == "" {
		req.Role = "cashier"
	} else if req.Role != "admin" && req.Role != "cashier" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid role. Allowed roles are 'admin' and 'cashier'",
		})
	}

	// Verify outlet if outlet_id is provided
	if req.OutletID != nil {
		var outlet models.Outlet
		if err := config.DB.First(&outlet, *req.OutletID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
				Success: false,
				Message: "Specified outlet_id not found",
			})
		}
	}

	// Check if email already registered
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(dto.APIResponse{
			Success: false,
			Message: "Email is already registered",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to hash password",
		})
	}

	// Save user to database
	user := models.User{
		OutletID:     req.OutletID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to create user: " + err.Error(),
		})
	}

	// Preload outlet relation if assigned
	if user.OutletID != nil {
		config.DB.Preload("Outlet").First(&user, user.ID)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    toUserResponse(&user),
	})
}

// Login handles user authentication, JWT token issuance, and storing bearer token to database
func (ac *AuthController) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Email and password are required",
		})
	}

	// Find user by email including Outlet relationship
	var user models.User
	if err := config.DB.Preload("Outlet").Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid email or password",
		})
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid email or password",
		})
	}

	// Generate JWT Token with role and outlet_id
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, user.OutletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to generate token",
		})
	}

	// Save bearer token in UserToken model
	userToken := models.UserToken{
		UserID: user.ID,
		Token:  &token,
	}
	if err := config.DB.Create(&userToken).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to save session token: " + err.Error(),
		})
	}

	authResponse := dto.AuthResponse{
		Token: token,
		User:  toUserResponse(&user),
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    authResponse,
	})
}

// Logout invalidates the active bearer token by setting it to NULL in database
func (ac *AuthController) Logout(c fiber.Ctx) error {
	tokenVal := c.Locals("token")
	tokenStr, ok := tokenVal.(string)
	if !ok || tokenStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Active token not found in session",
		})
	}

	// Update the token field to NULL for the current session token
	if err := config.DB.Model(&models.UserToken{}).
		Where("token = ?", tokenStr).
		Update("token", nil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: "Failed to invalidate token: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "Logout successful. Bearer token has been invalidated.",
	})
}

// GetProfile handles fetching the authenticated user's profile
func (ac *AuthController) GetProfile(c fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	var user models.User
	if err := config.DB.Preload("Outlet").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Success: false,
			Message: "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data:    toUserResponse(&user),
	})
}
