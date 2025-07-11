package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type AuthHandler struct {
	BaseHandler
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) LoginPlayer(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	response, err := h.AuthService.LoginPlayer(req.Username, req.Password)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials", err)
	}

	return h.SuccessResponse(c, "Login successful", response)
}

func (h *AuthHandler) LoginAdmin(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	response, err := h.AuthService.LoginAdmin(req.Username, req.Password)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials", err)
	}

	return h.SuccessResponse(c, "Login successful", response)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	response, err := h.AuthService.RefreshToken(req.RefreshToken)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token", err)
	}

	return h.SuccessResponse(c, "Token refreshed successfully", response)
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.AuthUser)
	return h.SuccessResponse(c, "Profile retrieved successfully", user)
}
