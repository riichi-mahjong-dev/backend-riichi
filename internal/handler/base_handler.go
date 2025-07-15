package handler

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/unicode/rangetable"
)

type BaseHandler struct{}

// Common response structure
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
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
	Data    any             `json:"data"`
	Meta    *PaginationMeta `json:"meta"`
}

type QueryPagination struct {
	Search string `json:"q"`
	SortBy string `json:"sortBy"`
	Sort   string `json:"sort"`
	Page string `json:"page"`
	Limit  int `json:"limit"`
	Offset int
	Filters map[string]string `json:"filters"`
}

// Helper functions
func (h *BaseHandler) SuccessResponse(c *fiber.Ctx, message string, data any) error {
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

func (h *BaseHandler) PaginatedSuccessResponse(c *fiber.Ctx, message string, data any, meta *PaginationMeta) error {
	return c.Status(200).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func (h *BaseHandler) GetPaginationParams(c *fiber.Ctx) QueryPagination {
	var queryPaginate QueryPagination

	if err := c.QueryParser(&queryPaginate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query param"
		})
	}

	rawQuery := c.Context().QueryArgs().String()
	parsed, _ := url.ParseQuery(rawQuery)

	queryPaginate.Filters = make(map[string]string)

	for key, value := range parsed {
		if strings.Contains(key, "[") && strings.HasSuffix(key, "]") {
			innerKey := key[len("filters[") : len(key)-1]
			if len(values) > 0 {
				queryPaginate.Filters[innerKey] = values[len(values)-1]
			}
		}
	}

	if queryPaginate.Page < 1 {
		queryPaginate.Page = 1
	}

	if queryPaginate.Limit < 1 {
		limit = 10
	}

	if queryPaginate.Sort == "" {
		queryPaginate.Sort = "ASC"
	}

	if queryPaginate.SortBy == "" {
		queryPaginate.SortBy = "id"
	}

	queryPaginate.Offset := (page - 1) * limit
	return queryPaginate
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
