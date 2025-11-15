package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
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
	adminService := services.NewAdminService(db.Conn)
	parlourService := services.NewParlourService(db.Conn)
	matchService := services.NewMatchService(db.Conn)
	provinceService := services.NewProvinceService(db.Conn)
	postService := services.NewPostService(db.Conn)
	logService := services.NewLogService(db.Conn)

	// Initialize auth service
	jwtConfig := env.LoadJwtConfig()
	authService := services.NewAuthService(db.Conn, playerService, adminService, jwtConfig.SecretKey)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize handlers
	playerHandler := handler.NewPlayerHandler(playerService)
	adminHandler := handler.NewAdminHandler(adminService)
	parlourHandler := handler.NewParlourHandler(parlourService)
	matchHandler := handler.NewMatchHandler(matchService)
	provinceHandler := handler.NewProvinceHandler(provinceService)
	postHandler := handler.NewPostHandler(postService)
	authHandler := handler.NewAuthHandler(authService)
	logHandler := handler.NewLogHandler(logService)

	// Authentication routes (public)
	auth := app.Group("/auth")
	auth.Post("/login/player", authHandler.LoginPlayer)
	auth.Post("/login/admin", authHandler.LoginAdmin)
	auth.Post("/refresh", authHandler.RefreshToken)

	// API routes with authentication
	api := app.Group("/api")

	api.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	initiateAdminRouter(api, adminHandler, authMiddleware)
	initiateMatchRouter(api, matchHandler, authMiddleware)
	initiateParlourRouter(api, parlourHandler, authMiddleware)
	initiatePlayerRouter(api, playerHandler, authHandler, authMiddleware)
	initiatePostRouter(api, postHandler, authMiddleware)
	initiateProvinceRouter(api, provinceHandler, authMiddleware)

	// Log routes (public view)
	api.Get("/logs", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), logHandler.GetAllLogs)
	api.Get("/logs/:id", authMiddleware.CheckAuthorization, authMiddleware.CheckRole([]string{"admin", "super-admin"}), logHandler.GetLogByID)

	// Health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running well",
		})
	})
}
