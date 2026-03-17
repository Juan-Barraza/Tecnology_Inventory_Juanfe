package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"inventory-juanfe/handlers"
	repository "inventory-juanfe/repositories"
	"inventory-juanfe/services"
)

func SetCatalogRouter(protected fiber.Router, db *sql.DB) {
	cityRepo := repository.NewCityRepository(db)
	areaRepo := repository.NewAreaRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	acctGroupRepo := repository.NewAccountingGroupRepository(db)

	catalogSvc := services.NewCatalogService(cityRepo, areaRepo, categoryRepo, acctGroupRepo)
	catalogH := handlers.NewCatalogHandler(catalogSvc)

	catalogs := protected.Group("/catalogs")
	catalogs.Get("/cities", catalogH.ListCities)
	catalogs.Get("/areas", catalogH.ListAreas)
	catalogs.Get("/categories", catalogH.ListCategories)
	catalogs.Get("/accounting-groups", catalogH.ListAccountingGroups)
	catalogs.Put("/accounting-groups/:id", catalogH.UpdateAccountingGroup)
}
