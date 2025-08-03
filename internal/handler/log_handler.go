package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type LogHandler struct {
	BaseHandler
	LogService *services.LogService
}

func NewLogHandler(logService *services.LogService) *LogHandler {
	return &LogHandler{
		LogService: logService,
	}
}

func (h *LogHandler) GetLogByID(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	log, err := h.LogService.GetLogById(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusNotFound, "Log not found", err)
	}

	return h.SuccessResponse(c, "Log retrieved successfully", log)
}

func (h *LogHandler) GetAllLogs(c *fiber.Ctx) error {
	queryPaginate := h.ParseQueryParams(c, []string{"created_between"})

	matches, total, err := h.LogService.GetAllLogs(queryPaginate)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve logs", err)
	}

	meta := h.CalculatePaginationMeta(queryPaginate.Page, queryPaginate.PageSize, total)
	return h.PaginatedSuccessResponse(c, "Logs retrieved successfully", matches, meta)
}
