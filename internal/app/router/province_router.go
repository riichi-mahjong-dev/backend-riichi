package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/handler"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/middleware"
)

func initiateProvinceRouter(api fiber.Router, provinceHandler *handler.ProvinceHandler, authMiddleware *middleware.AuthMiddleware) {
	// Province routes (public view, admin modifications)
	api.Get("/provinces", provinceHandler.GetAllProvinces)     // Public - guests can view
	api.Get("/provinces/:id", provinceHandler.GetProvinceByID) // Public - guests can view
	api.Post("/provinces", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), provinceHandler.CreateProvince)
	api.Put("/provinces/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), provinceHandler.UpdateProvince)
	api.Delete("/provinces/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), provinceHandler.DeleteProvince)
}
