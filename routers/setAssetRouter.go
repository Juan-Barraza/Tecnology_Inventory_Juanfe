package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"inventory-juanfe/handlers"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/services"
)

func SetAssetRouter(protected fiber.Router, db *sql.DB) {
	assetRepo := repository.NewAssetRepository(db)
	historyRepo := repository.NewStatusHistoryRepository(db)
	assignRepo := repository.NewAssignmentRepository(db)

	assetSvc := services.NewAssetService(assetRepo, historyRepo, assignRepo)
	assignSvc := services.NewAssignmentService(assignRepo, assetRepo)
	assetH := handlers.NewAssetHandler(assetSvc, assignSvc)

	assets := protected.Group("/assets")
	assets.Get("/", assetH.List)
	assets.Get("/:id", assetH.GetByID)
	assets.Post("/", assetH.Create)
	assets.Put("/:id", assetH.Update)
	assets.Patch("/:id/status", assetH.ChangeStatus)
	assets.Get("/:id/history", assetH.GetHistory)
	assets.Get("/:id/assignments", assetH.GetByAsset)

}
