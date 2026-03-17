package routers

import (
	"database/sql"
	"inventory-juanfe/middleware"

	"github.com/gofiber/fiber/v3"
)

// SetupRoutes configures all application routes.
// Add route groups here as features are implemented.
func SetupRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Global middleware
	app.Use(middleware.CORS())

	// Health check
	v1.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Inventory Juanfe API is running ",
		})
	})

}
