package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type ParlourHandler struct {
	BaseHandler
	ParlourService *services.ParlourService
}

func NewParlourHandler(parlourService *services.ParlourService) *ParlourHandler {
	return &ParlourHandler{
		ParlourService: parlourService,
	}
}

func (h *ParlourHandler) CreateParlour(c *fiber.Ctx) error {
	var req models.ParlourRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	parlour, err := h.ParlourService.CreateParlour(&req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create parlour", err)
	}

	return h.SuccessResponse(c, "Parlour created successfully", parlour)
}

func (h *ParlourHandler) GetParlourByID(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	parlour, err := h.ParlourService.GetParlourByID(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusNotFound, "Parlour not found", err)
	}

	return h.SuccessResponse(c, "Parlour retrieved successfully", parlour)
}

func (h *ParlourHandler) GetAllParlours(c *fiber.Ctx) error {
	queryPaginate := h.GetPaginationParams(c)

	parlours, err := h.ParlourService.GetAllParlours(queryPaginate)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve parlours", err)
	}

	// Count total parlours for pagination
	total, err := h.ParlourService.Count(&models.Parlour{})
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count parlours", err)
	}

	meta := h.CalculatePaginationMeta(int(c.QueryInt("page", 1)), queryPaginate.Limit, total)
	return h.PaginatedSuccessResponse(c, "Parlours retrieved successfully", parlours, meta)
}

func (h *ParlourHandler) UpdateParlour(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req models.ParlourRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	parlour, err := h.ParlourService.UpdateParlour(id, &req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update parlour", err)
	}

	return h.SuccessResponse(c, "Parlour updated successfully", parlour)
}

func (h *ParlourHandler) DeleteParlour(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	err = h.ParlourService.DeleteParlour(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete parlour", err)
	}

	return h.SuccessResponse(c, "Parlour deleted successfully", nil)
}
