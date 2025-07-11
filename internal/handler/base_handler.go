package handler

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type BaseHandler struct{}

// Common response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data"`
	Meta    *PaginationMeta `json:"meta"`
}

// Helper functions
func (h *BaseHandler) SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(200).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (h *BaseHandler) ErrorResponse(c *fiber.Ctx, statusCode int, message string, err error) error {
	response := Response{
		Success: false,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	return c.Status(statusCode).JSON(response)
}

func (h *BaseHandler) PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, meta *PaginationMeta) error {
	return c.Status(200).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func (h *BaseHandler) GetPaginationParams(c *fiber.Ctx) (int, int) {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	
	offset := (page - 1) * limit
	return limit, offset
}

func (h *BaseHandler) GetIDParam(c *fiber.Ctx) (uint64, error) {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *BaseHandler) CalculatePaginationMeta(page, limit int, total int64) *PaginationMeta {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}
}
