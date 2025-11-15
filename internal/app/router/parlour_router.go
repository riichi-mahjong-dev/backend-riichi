package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/handler"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/middleware"
)

func initiateParlourRouter(api fiber.Router, parlourHandler *handler.ParlourHandler, authMiddleware *middleware.AuthMiddleware) {
	// Parlour routes (public view, admin modifications)
	api.Get("/parlours", parlourHandler.GetAllParlours)     // Public - guests can view
	api.Get("/parlours/:id", parlourHandler.GetParlourByID) // Public - guests can view
	api.Post("/parlours", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), parlourHandler.CreateParlour)
	api.Put("/parlours/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), parlourHandler.UpdateParlour)
	api.Delete("/parlours/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), parlourHandler.DeleteParlour)
}
