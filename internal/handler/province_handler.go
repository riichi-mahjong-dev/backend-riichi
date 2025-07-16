package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type ProvinceHandler struct {
	BaseHandler
	ProvinceService *services.ProvinceService
}

func NewProvinceHandler(provinceService *services.ProvinceService) *ProvinceHandler {
	return &ProvinceHandler{
		ProvinceService: provinceService,
	}
}

func (h *ProvinceHandler) CreateProvince(c *fiber.Ctx) error {
	var req models.ProvinceRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	province, err := h.ProvinceService.CreateProvince(&req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create province", err)
	}

	return h.SuccessResponse(c, "Province created successfully", province)
}

func (h *ProvinceHandler) GetProvinceByID(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	province, err := h.ProvinceService.GetProvinceByID(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusNotFound, "Province not found", err)
	}

	return h.SuccessResponse(c, "Province retrieved successfully", province)
}

func (h *ProvinceHandler) GetAllProvinces(c *fiber.Ctx) error {
	queryPaginate := h.GetPaginationParams(c)

	provinces, err := h.ProvinceService.GetAllProvinces(queryPaginate.Limit, queryPaginate.Offset)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve provinces", err)
	}

	// Count total provinces for pagination
	total, err := h.ProvinceService.Count(&models.Province{})
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count provinces", err)
	}

	meta := h.CalculatePaginationMeta(int(c.QueryInt("page", 1)), queryPaginate.Limit, total)
	return h.PaginatedSuccessResponse(c, "Provinces retrieved successfully", provinces, meta)
}

func (h *ProvinceHandler) UpdateProvince(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req models.ProvinceRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	province, err := h.ProvinceService.UpdateProvince(id, &req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update province", err)
	}

	return h.SuccessResponse(c, "Province updated successfully", province)
}

func (h *ProvinceHandler) DeleteProvince(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	err = h.ProvinceService.DeleteProvince(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete province", err)
	}

	return h.SuccessResponse(c, "Province deleted successfully", nil)
}
