package middleware

import (
	"strings"

	"bazar-be/config"
	"bazar-be/dto"
	"bazar-be/models"
	"bazar-be/utils"

	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware checks the validity of the JWT token passed in Authorization header
// and verifies that the token is active (not null/logged out) in the database.
func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Success: false,
				Message: "Authorization header is required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Success: false,
				Message: "Invalid authorization format. Expected 'Bearer <token>'",
			})
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
		}

		// Check if token exists and is active in database (not null / logged out)
		var userToken models.UserToken
		if err := config.DB.Where("user_id = ? AND token = ?", claims.UserID, tokenString).First(&userToken).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
				Success: false,
				Message: "Session has expired or you have logged out",
			})
		}

		// Store user info and current token in context locals
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("outlet_id", claims.OutletID)
		c.Locals("token", tokenString)

		return c.Next()
	}
}

// RequireRoles restricts route access to specific roles (e.g. "admin", "cashier")
func RequireRoles(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole == "" {
			return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
				Success: false,
				Message: "Access forbidden: role not found",
			})
		}

		for _, role := range roles {
			if strings.EqualFold(userRole, role) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(dto.APIResponse{
			Success: false,
			Message: "Access forbidden: insufficient permissions",
		})
	}
}
