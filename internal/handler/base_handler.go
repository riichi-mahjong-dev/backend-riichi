/*
 * Base Handler to use by other handler, common function used by handler
 *
 * Author: Kristian Ruben
 */
package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/commons"
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
	HasMore     bool  `json:"has_more"`
}

type PaginatedResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    any             `json:"data"`
	Meta    *PaginationMeta `json:"meta"`
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
	hasMore := page < totalPages
	return &PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
		HasMore:     hasMore,
	}
}

func (h *BaseHandler) ParseQueryParams(c *fiber.Ctx, filtersAllowed []string) commons.QueryParams {
	page := max(c.QueryInt("page[number]", 1), 1)
	pageSize := max(c.QueryInt("page[size]", 10), 10)
	search := c.Query("search", "")
	sort := c.Query("sort", "id")
	filters := make(map[string]any)
	order := "asc"
	orderBy := sort

	if len(sort) > 0 && sort[0] == '-' {
		order = "desc"
		orderBy = sort[1:]
	}

	for _, filterAllowed := range filtersAllowed {
		filterValue := c.Query("filter["+filterAllowed+"]", "")
		if filterValue != "" {
			filters[filterAllowed] = strings.Split(filterValue, ",")
		}
	}

	return commons.QueryParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		OrderBy:  orderBy,
		Order:    order,
		Filters:  filters,
	}
}
