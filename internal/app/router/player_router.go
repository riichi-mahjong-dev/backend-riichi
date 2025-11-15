package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/handler"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/middleware"
)

func initiatePlayerRouter(api fiber.Router, playerHandler *handler.PlayerHandler, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	// Profile route (requires authentication)
	api.Get("/profile", authMiddleware.CheckAuthorization, authHandler.GetProfile)

	// Player routes (guests can view, admins can manage)
	api.Get("/players", playerHandler.GetAllPlayers)                                                                                                // Public - guests can view
	api.Get("/players/:id", playerHandler.GetPlayerByID)                                                                                            // Public - guests can view
	api.Post("/players", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), playerHandler.CreatePlayer) // Public registration
	api.Put("/players/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), playerHandler.UpdatePlayer)
	api.Delete("/players/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), playerHandler.DeletePlayer)
	api.Post("/players/change-password", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"player"}), playerHandler.ChangePassword)
}
