package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type AuthMiddleware struct {
	AuthService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
	}
}

func (authMiddleware *AuthMiddleware) CheckAuthorization(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization", "")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Authorization header required",
		})
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token format",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := authMiddleware.AuthService.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token",
			"error":   err.Error(),
		})
	}

	// Create auth user and store in context
	authUser := &models.AuthUser{
		ID:       claims.UserID,
		Username: claims.Username,
		UserType: claims.UserType,
		Role:     claims.Role,
	}

	c.Locals("user", authUser)
	c.Locals("claims", claims)

	return c.Next()
}

func (authMiddleware *AuthMiddleware) CheckRole(requiredRoles []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userData := c.Locals("user").(*models.AuthUser)

		if userData == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "User not found",
			})
		}

		// Check if user type matches required roles
		for _, role := range requiredRoles {
			if role == "any" {
				return c.Next()
			}
			if role == "player" && userData.UserType == models.UserTypePlayer {
				return c.Next()
			}
			if role == "admin" && userData.UserType == models.UserTypeAdmin {
				return c.Next()
			}
			if role == "super-admin" && userData.UserType == models.UserTypeAdmin && userData.Role == "super-admin" {
				return c.Next()
			}
			if role == "staff" && userData.UserType == models.UserTypeAdmin && userData.Role == "staff" {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Access denied",
		})
	}
}
