package routers

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"inventory-juanfe/middleware"
)

// SetupRoutes configures all application routes.
func SetupRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Global middleware
	// CORS
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{allowedOrigins},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	// Health check
	v1.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Inventory Juanfe API is running ",
		})
	})
	// public router
	SetAuthRouter(v1, db)

	// Protected group
	protected := v1.Group("/")
	protected.Use(middleware.JWTAuth())

	// Register routers
	SetAssetRouter(protected, db)
	SetAssignmentRouter(protected, db)
	SetInventoryRouter(protected, db)
	SetCatalogRouter(protected, db)
}
