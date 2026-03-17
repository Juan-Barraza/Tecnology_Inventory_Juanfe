package routers

import (
	"database/sql"
	"inventory-juanfe/handlers"
	"inventory-juanfe/middleware"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// SetupRoutes configures all application routes.
// Add route groups here as features are implemented.
func SetupRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Global middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
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

	SetAuthRouter(v1, db)
	protected := v1.Group("/")
	protected.Use(middleware.JWTAuth())

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	protected.Get("/auth/me", authHandler.Me)

}
