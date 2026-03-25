package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	dtos "inventory-juanfe/dtos/request"
	"inventory-juanfe/services"
	"inventory-juanfe/utils"
)

type AssetHandler struct {
	svc       *services.AssetService
	svcAssign *services.AssignmentService
}

func NewAssetHandler(svc *services.AssetService, svcAssign *services.AssignmentService) *AssetHandler {
	return &AssetHandler{svc: svc, svcAssign: svcAssign}
}

func (h *AssetHandler) List(c fiber.Ctx) error {
	var f dtos.AssetFilter
	if err := c.Bind().Query(&f); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid query params")
	}

	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	userID := utils.GetUserID(c)
	assets, total, err := h.svc.List(f, userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch assets")
	}

	return utils.SuccessPaginated(c, assets, total, f.Page, f.Limit)
}

func (h *AssetHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	userId := utils.GetUserID(c)
	asset, err := h.svc.GetByID(id, userId)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch asset")
	}
	if asset == nil {
		return utils.Error(c, http.StatusNotFound, "asset not found")
	}

	return utils.Success(c, http.StatusOK, asset)
}

func (h *AssetHandler) Create(c fiber.Ctx) error {
	var req dtos.CreateAssetRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.ValidateCreateAsset(req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	userID := utils.GetUserID(c)
	asset, err := h.svc.Create(req, userID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, http.StatusCreated, asset)
}

func (h *AssetHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req dtos.UpdateAssetRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	userId := utils.GetUserID(c)
	asset, err := h.svc.Update(id, req, userId)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if asset == nil {
		return utils.Error(c, http.StatusNotFound, "asset not found")
	}

	return utils.Success(c, http.StatusOK, asset)
}

func (h *AssetHandler) ChangeStatus(c fiber.Ctx) error {
	id := c.Params("id")

	var req dtos.UpdateAssetStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.ValidateUpdateAssetStatus(req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	userID := utils.GetUserID(c)
	asset, err := h.svc.ChangeStatus(id, req, userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if asset == nil {
		return utils.Error(c, http.StatusNotFound, "asset not found")
	}

	return utils.Success(c, http.StatusOK, asset)
}

func (h *AssetHandler) GetHistory(c fiber.Ctx) error {
	assetID := c.Params("id")

	history, err := h.svc.GetHistory(assetID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch history")
	}

	return utils.Success(c, http.StatusOK, history)
}

func (h *AssetHandler) GetByAsset(c fiber.Ctx) error {
	assetID := c.Params("id")

	if assetID == "" {
		return utils.Error(c, http.StatusBadRequest, "asset id is required")
	}

	assignments, err := h.svcAssign.GetByAsset(assetID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch assignments")
	}

	return utils.Success(c, http.StatusOK, assignments)
}
