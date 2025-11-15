/*
 * Admin handler to handle data admin, select, insert, delete, update
 *
 * Author: Kristian Ruben
 */
package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/models"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

type AdminHandler struct {
	BaseHandler
	AdminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		AdminService: adminService,
	}
}

func (h *AdminHandler) CreateAdmin(c *fiber.Ctx) error {
	var req models.AdminRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	admin, err := h.AdminService.CreateAdmin(&req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create admin", err)
	}

	return h.SuccessResponse(c, "Admin created successfully", admin)
}

func (h *AdminHandler) GetAdminByID(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	admin, err := h.AdminService.GetAdminByID(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusNotFound, "Admin not found", err)
	}

	return h.SuccessResponse(c, "Admin retrieved successfully", admin)
}

func (h *AdminHandler) GetAllAdmins(c *fiber.Ctx) error {
	queryPaginate := h.ParseQueryParams(c, nil)

	admins, total, err := h.AdminService.GetAllAdmins(queryPaginate)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve admins", err)
	}

	meta := h.CalculatePaginationMeta(queryPaginate.Page, queryPaginate.PageSize, total)
	return h.PaginatedSuccessResponse(c, "Admins retrieved successfully", admins, meta)
}

func (h *AdminHandler) UpdateAdmin(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	var req models.AdminRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	admin, err := h.AdminService.UpdateAdmin(id, &req)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update admin", err)
	}

	return h.SuccessResponse(c, "Admin updated successfully", admin)
}

func (h *AdminHandler) DeleteAdmin(c *fiber.Ctx) error {
	id, err := h.GetIDParam(c)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID", err)
	}

	err = h.AdminService.DeleteAdmin(id)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete admin", err)
	}

	return h.SuccessResponse(c, "Admin deleted successfully", nil)
}

func (h *AdminHandler) ChangePassword(c *fiber.Ctx) error {
	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return h.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	userData := c.Locals("user").(*models.AuthUser)

	err := h.AdminService.UpdatePassword(userData.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		return h.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to change password", err)
	}

	return h.SuccessResponse(c, "Password changed", nil)
}
