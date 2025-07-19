package handler

import (
	"fmt"
	"net/url"
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
}

type PaginatedResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    any             `json:"data"`
	Meta    *PaginationMeta `json:"meta"`
}

type QueryParams struct {
	Page     int               `form:"page" query:"page"`
	PageSize int               `form:"pageSize" query:"pageSize"`
	Search   string            `form:"search" query:"search"`
	Filters  map[string]string // Custom filters, e.g., age, status
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

func (h *BaseHandler) GetPaginationParams(c *fiber.Ctx) commons.QueryPagination {
	var queryPaginate commons.QueryPagination

	if err := c.QueryParser(&queryPaginate); err != nil {
		fmt.Errorf("Invalid query param")
	}

	rawQuery := c.Context().QueryArgs().String()
	parsed, _ := url.ParseQuery(rawQuery)

	queryPaginate.Filters = make(map[string]string)

	for key, value := range parsed {
		if strings.Contains(key, "[") && strings.HasSuffix(key, "]") {
			innerKey := key[len("filters[") : len(key)-1]
			if len(value) > 0 {
				queryPaginate.Filters[innerKey] = value[len(value)-1]
			}
		}
	}

	if queryPaginate.Page < 1 {
		queryPaginate.Page = 1
	}

	if queryPaginate.Limit < 1 {
		queryPaginate.Limit = 10
	}

	if queryPaginate.Sort == "" {
		queryPaginate.Sort = "ASC"
	}

	if queryPaginate.SortBy == "" {
		queryPaginate.SortBy = "id"
	}

	queryPaginate.Offset = (queryPaginate.Page - 1) * queryPaginate.Limit
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

func ParseQueryParams(c *fiber.Ctx, filtersAllowed []string) QueryParams {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	search := c.Query("search", "")
	filters := make(map[string]string)

	for _, filterAllowed := range filtersAllowed {
		filterValue := c.Query("filter["+filterAllowed+"]", "")
		if filterValue != "" {
			filters[filterAllowed] = filterValue
		}
	}

	return QueryParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		Filters:  filters,
	}
}
