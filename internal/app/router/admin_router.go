package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/handler"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/middleware"
)

func initiateAdminRouter(api fiber.Router, adminHandler *handler.AdminHandler, authMiddleware *middleware.AuthMiddleware) {
	// Admin routes (super-admin only)
	api.Get("/admins", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.GetAllAdmins)
	api.Get("/admins/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.GetAdminByID)
	api.Post("/admins", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.CreateAdmin)
	api.Put("/admins/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.UpdateAdmin)
	api.Delete("/admins/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.DeleteAdmin)
	api.Post("/admins/change-password", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), adminHandler.ChangePassword)
}
