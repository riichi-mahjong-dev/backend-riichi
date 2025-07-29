package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type PlayerHandler struct {
	BaseHandler
	PlayerService *services.PlayerService
}

func NewPlayerHandler(playerService *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		PlayerService: playerService,
	}
}

func (h *PlayerHandler) CreatePlayer(c *fiber.Ctx) error {
	var req models.PlayerRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	player, err := h.PlayerService.CreatePlayer(&req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create player", err)
	}

	return h.SuccessResponse(c, "Player created successfully", player)
}

func (h *PlayerHandler) GetPlayerByID(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	player, err := h.PlayerService.GetPlayerByID(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusNotFound, "Player not found", err)
	}

	return h.SuccessResponse(c, "Player retrieved successfully", player)
}

func (h *PlayerHandler) GetAllPlayers(c *fiber.Ctx) error {
	queryPaginate := h.ParseQueryParams(c, []string{"mmr_between"})

	players, total, err := h.PlayerService.GetAllPlayers(queryPaginate)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve players", err)
	}

	meta := h.CalculatePaginationMeta(queryPaginate.Page, queryPaginate.PageSize, total)
	return h.PaginatedSuccessResponse(c, "Players retrieved successfully", players, meta)
}

func (h *PlayerHandler) UpdatePlayer(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req models.PlayerRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	player, err := h.PlayerService.UpdatePlayer(id, &req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update player", err)
	}

	return h.SuccessResponse(c, "Player updated successfully", player)
}

func (h *PlayerHandler) DeletePlayer(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	err = h.PlayerService.DeletePlayer(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete player", err)
	}

	return h.SuccessResponse(c, "Player deleted successfully", nil)
}

func (h *PlayerHandler) ChangePassword(c *fiber.Ctx) error {
	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	userData := c.Locals("user").(*models.AuthUser)

	err := h.PlayerService.UpdatePassword(userData.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to change password", err)
	}

	return h.SuccessResponse(c, "Password changed", nil)
}
