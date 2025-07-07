package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	// "github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	// "github.com/riichi-mahjong-dev/backend-riichi/internal/services"
	"github.com/riichi-mahjong-dev/backend-riichi/utils"
)

type AuthMiddleware struct {
	userService   *services.UserService
	authGenerator *utils.AuthToken
}

func InitializeAuthMiddleware(userService *services.UserService, authGenerator *utils.AuthToken) *AuthMiddleware {
	return &AuthMiddleware{
		userService:   userService,
		authGenerator: authGenerator,
	}
}

func (authMiddleware *AuthMiddleware) CheckAuthorization(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization", "")

	if authHeader == "" {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", fmt.Errorf("authorization not found"))(c)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", fmt.Errorf("invalid token format"))(c)
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	exp, userId, err := authMiddleware.authGenerator.ValidateToken(tokenString)

	if err != nil {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", err)(c)
	}

	now := time.Now().Unix()

	if int64(exp) < now {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", fmt.Errorf("token expired"))(c)
	}

	user, err := authMiddleware.userService.GetUser(uint(userId))

	if err != nil {
		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", fmt.Errorf("user not found"))(c)
	}

	c.Locals("user", user)

	return c.Next()
}

func (authMiddleware *AuthMiddleware) CheckRole(requiredRole []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userData := c.Locals("user").(*models.User)

		if userData == nil {
			return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", fmt.Errorf("user not found"))(c)
		}

		for _, role := range requiredRole {
			if role == userData.Role {
				return c.Next()
			}
		}

		return utils.ResponseError(fiber.StatusUnauthorized, "Unauthorized", fmt.Errorf("user not authorized to this route"))(c)
	}
}
