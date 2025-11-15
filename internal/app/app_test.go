package app

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/riichi-mahjong-dev/backend-riichi/commons"
	"github.com/riichi-mahjong-dev/backend-riichi/configs"
	"github.com/riichi-mahjong-dev/backend-riichi/database"
	"github.com/riichi-mahjong-dev/backend-riichi/internal/app/router"
	"github.com/riichi-mahjong-dev/backend-riichi/utils"
)

func TestCreateApp(t *testing.T) {
	env := &configs.EnvConfig{}
	emailConfig := env.LoadEmailConfig()
	db := database.MockConnectDatabase()

	mailer, err := utils.InitializeEmailer(emailConfig)
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	app.Static("/images", "./images")
	router.InitializeRoute(app, &commons.AppConfig{
		Db:     db,
		Mailer: mailer,
		Env:    env,
	})

	// Make a test request to verify the route was created
	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	req = httptest.NewRequest("GET", "/api/test", nil)
	resp, err = app.Test(req)

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}

	req = httptest.NewRequest("GET", "/api/players", nil)
	resp, err = app.Test(req)

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}
