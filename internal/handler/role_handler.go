package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type RoleHandler struct {
	BaseHandler
	RoleService *services.RoleService
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		RoleService: roleService,
	}
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req models.RoleRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	role, err := h.RoleService.CreateRole(&req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create role", err)
	}

	return h.SuccessResponse(c, "Role created successfully", role)
}

func (h *RoleHandler) GetRoleByID(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	role, err := h.RoleService.GetRoleByID(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusNotFound, "Role not found", err)
	}

	return h.SuccessResponse(c, "Role retrieved successfully", role)
}

func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	queryPaginate := h.GetPaginationParams(c)

	roles, err := h.RoleService.GetAllRoles(queryPaginate.Limit, queryPaginate.Offset)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve roles", err)
	}

	// Count total roles for pagination
	total, err := h.RoleService.Count(&models.Role{})
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to count roles", err)
	}

	meta := h.CalculatePaginationMeta(int(c.QueryInt("page", 1)), queryPaginate.Limit, total)
	return h.PaginatedSuccessResponse(c, "Roles retrieved successfully", roles, meta)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req models.RoleRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	role, err := h.RoleService.UpdateRole(id, &req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update role", err)
	}

	return h.SuccessResponse(c, "Role updated successfully", role)
}

func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	err = h.RoleService.DeleteRole(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete role", err)
	}

	return h.SuccessResponse(c, "Role deleted successfully", nil)
}
