package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"inventory-juanfe/handlers"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/services"
)

func SetAssignmentRouter(protected fiber.Router, db *sql.DB) {
	assignRepo := repository.NewAssignmentRepository(db)
	assetRepo := repository.NewAssetRepository(db)

	assignSvc := services.NewAssignmentService(assignRepo, assetRepo)
	assignH := handlers.NewAssignmentHandler(assignSvc)

	assignments := protected.Group("/assignments")
	assignments.Post("/", assignH.Create)
	assignments.Patch("/:id/release", assignH.Release)
}
