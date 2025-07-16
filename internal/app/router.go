package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/handler"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/middleware"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/services"
)

func InitializeRoute(app *fiber.App, appConfig *commons.AppConfig) {
	db := appConfig.Db
	env := appConfig.Env

	// Initialize services
	playerService := services.NewPlayerService(db.Conn)
	roleService := services.NewRoleService(db.Conn)
	adminService := services.NewAdminService(db.Conn)
	parlourService := services.NewParlourService(db.Conn)
	matchService := services.NewMatchService(db.Conn)
	provinceService := services.NewProvinceService(db.Conn)
	postService := services.NewPostService(db.Conn)

	// Initialize auth service
	jwtConfig := env.LoadJwtConfig()
	authService := services.NewAuthService(db.Conn, playerService, adminService, jwtConfig.SecretKey)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize handlers
	playerHandler := handler.NewPlayerHandler(playerService)
	roleHandler := handler.NewRoleHandler(roleService)
	adminHandler := handler.NewAdminHandler(adminService)
	parlourHandler := handler.NewParlourHandler(parlourService)
	matchHandler := handler.NewMatchHandler(matchService)
	provinceHandler := handler.NewProvinceHandler(provinceService)
	postHandler := handler.NewPostHandler(postService)
	authHandler := handler.NewAuthHandler(authService)

	// Authentication routes (public)
	auth := app.Group("/auth")
	auth.Post("/login/player", authHandler.LoginPlayer)
	auth.Post("/login/admin", authHandler.LoginAdmin)
	auth.Post("/refresh", authHandler.RefreshToken)

	// API routes with authentication
	api := app.Group("/api")

	// Profile route (requires authentication)
	api.Get("/profile", authMiddleware.CheckAuthorization, authHandler.GetProfile)

	// Player routes (guests can view, admins can manage)
	api.Get("/players", playerHandler.GetAllPlayers)     // Public - guests can view
	api.Get("/players/:id", playerHandler.GetPlayerByID) // Public - guests can view
	api.Post("/players", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}) playerHandler.CreatePlayer)     // Public registration
	api.Put("/players/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), playerHandler.UpdatePlayer)
	api.Delete("/players/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), playerHandler.DeletePlayer)

	// Role routes (admin only)
	api.Get("/roles", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), roleHandler.GetAllRoles)
	api.Get("/roles/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), roleHandler.GetRoleByID)
	api.Post("/roles", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), roleHandler.CreateRole)
	api.Put("/roles/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), roleHandler.UpdateRole)
	api.Delete("/roles/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), roleHandler.DeleteRole)

	// Admin routes (super-admin only)
	api.Get("/admins", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.GetAllAdmins)
	api.Get("/admins/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.GetAdminByID)
	api.Post("/admins", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.CreateAdmin)
	api.Put("/admins/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.UpdateAdmin)
	api.Delete("/admins/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), adminHandler.DeleteAdmin)

	// Province routes (public view, admin modifications)
	api.Get("/provinces", provinceHandler.GetAllProvinces)     // Public - guests can view
	api.Get("/provinces/:id", provinceHandler.GetProvinceByID) // Public - guests can view
	api.Post("/provinces", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), provinceHandler.CreateProvince)
	api.Put("/provinces/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), provinceHandler.UpdateProvince)
	api.Delete("/provinces/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), provinceHandler.DeleteProvince)

	// Parlour routes (public view, admin modifications)
	api.Get("/parlours", parlourHandler.GetAllParlours)     // Public - guests can view
	api.Get("/parlours/:id", parlourHandler.GetParlourByID) // Public - guests can view
	api.Post("/parlours", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), parlourHandler.CreateParlour)
	api.Put("/parlours/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), parlourHandler.UpdateParlour)
	api.Delete("/parlours/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"super-admin"}), parlourHandler.DeleteParlour)

	// Match routes (guests can view, players/admins can manage)
	api.Get("/matches", matchHandler.GetAllMatches)    // Public - guests can view
	api.Get("/matches/:id", matchHandler.GetMatchByID) // Public - guests can view
	api.Post("/matches", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"player", "admin"}), matchHandler.CreateMatch)
	api.Put("/matches/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"player", "admin"}), matchHandler.UpdateMatch)
	api.Delete("/matches/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), matchHandler.DeleteMatch)
	api.Post("/matches/:id/approve", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), matchHandler.ApproveMatch)

	// Post routes (public view, admin modifications)
	api.Get("/posts", postHandler.GetAllPosts)     // Public - guests can view
	api.Get("/posts/:id", postHandler.GetPostByID) // Public - guests can view
	api.Post("/posts", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), postHandler.CreatePost)
	api.Put("/posts/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), postHandler.UpdatePost)
	api.Delete("/posts/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin"}), postHandler.DeletePost)

	// Health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
}
