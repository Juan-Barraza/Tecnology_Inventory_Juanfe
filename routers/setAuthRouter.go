package routers

import (
	"database/sql"
	"inventory-juanfe/handlers"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/services"

	"github.com/gofiber/fiber/v3"
)

func SetAuthRouter(v1 fiber.Router, db *sql.DB) {
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	auth := v1.Group("/auth")
	auth.Post("/login", authHandler.Login)
}
