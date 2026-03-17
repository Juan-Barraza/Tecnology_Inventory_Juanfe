package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"inventory-juanfe/handlers"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/services"
)

func SetDashboardRouter(protected fiber.Router, db *sql.DB) {
	repo := repository.NewDashboardRepository(db)
	svc := services.NewDashboardService(repo)
	h := handlers.NewDashboardHandler(svc)

	protected.Get("/dashboard", h.GetDashboard)
}
