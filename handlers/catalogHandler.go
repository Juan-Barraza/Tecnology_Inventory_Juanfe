package handlers

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"

	dtos "inventory-juanfe/dtos/request"
	"inventory-juanfe/services"
	"inventory-juanfe/utils"
)

type CatalogHandler struct {
	svc *services.CatalogService
}

func NewCatalogHandler(svc *services.CatalogService) *CatalogHandler {
	return &CatalogHandler{svc: svc}
}

func (h *CatalogHandler) ListCities(c fiber.Ctx) error {
	cities, err := h.svc.ListCities()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch cities")
	}
	return utils.Success(c, http.StatusOK, cities)
}

func (h *CatalogHandler) ListAreas(c fiber.Ctx) error {
	areas, err := h.svc.ListAreas()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch areas")
	}
	return utils.Success(c, http.StatusOK, areas)
}

func (h *CatalogHandler) ListCategories(c fiber.Ctx) error {
	cats, err := h.svc.ListCategories()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch categories")
	}
	return utils.Success(c, http.StatusOK, cats)
}

func (h *CatalogHandler) ListAccountingGroups(c fiber.Ctx) error {
	groups, err := h.svc.ListAccountingGroups()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch accounting groups")
	}
	return utils.Success(c, http.StatusOK, groups)
}

func (h *CatalogHandler) UpdateAccountingGroup(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id == 0 {
		return utils.Error(c, http.StatusBadRequest, "invalid id")
	}

	var req dtos.UpdateAccountingGroupRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.ValidateUpdateAccountingGroup(req.Name); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := h.svc.UpdateAccountingGroup(id, req.Name); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not update accounting group")
	}

	return utils.SuccessMessage(c, http.StatusOK, "accounting group updated")
}
