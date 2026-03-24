package routers

import (
	"database/sql"
	"inventory-juanfe/handlers"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/services"

	"github.com/gofiber/fiber/v3"
)

func SetExporterRouter(db *sql.DB, protected fiber.Router) {
	exportRepo := repository.NewExporterRepository(db)
	exportService := services.NewExportAssetsToXlsx(exportRepo)
	exportHandler := handlers.NewExportHandlerXlsx(exportService)

	protected.Get("/export/xlsx", exportHandler.ExportXlsx)
}
