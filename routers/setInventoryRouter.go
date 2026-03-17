package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"inventory-juanfe/handlers"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/services"
)

func SetInventoryRouter(protected fiber.Router, db *sql.DB) {
	inventoryRepo := repository.NewInventoryRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	historyRepo := repository.NewStatusHistoryRepository(db)

	inventorySvc := services.NewInventoryService(inventoryRepo, assetRepo, historyRepo)
	inventoryH := handlers.NewInventoryHandler(inventorySvc)

	inventory := protected.Group("/inventory")
	inventory.Get("/periods", inventoryH.ListPeriods)
	inventory.Post("/periods", inventoryH.CreatePeriod)
	inventory.Patch("/periods/:id/close", inventoryH.ClosePeriod)
	inventory.Get("/periods/:id/records", inventoryH.GetRecords)
	inventory.Post("/records", inventoryH.RecordAsset)
	inventory.Get("/periods/:id/progress", inventoryH.GetProgress)
	inventory.Get("/periods/:id/assets", inventoryH.GetPeriodAssets)
}
