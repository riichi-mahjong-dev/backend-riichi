package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/handler"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/middleware"
)

func initiateMatchRouter(api fiber.Router, matchHandler *handler.MatchHandler, authMiddleware *middleware.AuthMiddleware) {
	api.Get("/matches", matchHandler.GetAllMatches)    // Public - guests can view
	api.Get("/matches/:id", matchHandler.GetMatchByID) // Public - guests can view
	api.Post("/matches", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"player", "admin", "super-admin"}), matchHandler.CreateMatch)
	api.Put("/matches/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"player", "admin", "super-admin"}), matchHandler.UpdateMatch)
	api.Delete("/matches/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), matchHandler.DeleteMatch)
	api.Post("/matches/:id/approve", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), matchHandler.ApproveMatch)
	api.Post("/matches/:id/point", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), matchHandler.PointMatch)
}
