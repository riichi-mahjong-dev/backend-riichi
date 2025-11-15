package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/handler"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/middleware"
)

func initiatePostRouter(api fiber.Router, postHandler *handler.PostHandler, authMiddleware *middleware.AuthMiddleware) {
	// Post routes (public view, admin modifications)
	api.Get("/posts", postHandler.GetAllPosts)     // Public - guests can view
	api.Get("/posts/:id", postHandler.GetPostByID) // Public - guests can view
	api.Post("/posts", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), postHandler.CreatePost)
	api.Put("/posts/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), postHandler.UpdatePost)
	api.Delete("/posts/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), postHandler.DeletePost)
}
