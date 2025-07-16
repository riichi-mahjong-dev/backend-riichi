package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type MatchHandler struct {
	BaseHandler
	MatchService *services.MatchService
}

func NewMatchHandler(matchService *services.MatchService) *MatchHandler {
	return &MatchHandler{
		MatchService: matchService,
	}
}

func (h *MatchHandler) CreateMatch(c *fiber.Ctx) error {
	var req models.MatchRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	match, err := h.MatchService.CreateMatch(&req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create match", err)
	}

	return h.SuccessResponse(c, "Match created successfully", match)
}

func (h *MatchHandler) GetMatchByID(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	match, err := h.MatchService.GetMatchByID(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusNotFound, "Match not found", err)
	}

	return h.SuccessResponse(c, "Match retrieved successfully", match)
}

func (h *MatchHandler) GetAllMatches(c *fiber.Ctx) error {
	queryPaginate := h.GetPaginationParams(c)

	matches, err := h.MatchService.GetAllMatches(queryPaginate)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve matches", err)
	}

	// Count total matches for pagination
	total, err := h.MatchService.Count(&models.Match{})
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count matches", err)
	}

	meta := h.CalculatePaginationMeta(int(c.QueryInt("page", 1)), queryPaginate.Limit, total)
	return h.PaginatedSuccessResponse(c, "Matches retrieved successfully", matches, meta)
}

func (h *MatchHandler) UpdateMatch(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req models.MatchRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	match, err := h.MatchService.UpdateMatch(id, &req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update match", err)
	}

	return h.SuccessResponse(c, "Match updated successfully", match)
}

func (h *MatchHandler) DeleteMatch(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	err = h.MatchService.DeleteMatch(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete match", err)
	}

	return h.SuccessResponse(c, "Match deleted successfully", nil)
}

func (h *MatchHandler) ApproveMatch(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	// Get approver ID from request or context
	var approverReq struct {
		ApprovedBy uint64 `json:"approved_by" validate:"required"`
	}
	if err := c.BodyParser(&approverReq); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	match, err := h.MatchService.ApproveMatch(id, approverReq.ApprovedBy)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to approve match", err)
	}

	return h.SuccessResponse(c, "Match approved successfully", match)
}

func (h *MatchHandler) PointMatch(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req models.PointMatchRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	match, err := h.MatchService.PointMatch(id, &req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save point match", err)
	}

	return h.SuccessResponse(c, "Point match saved", match)
}
