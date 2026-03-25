package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	dtos "inventory-juanfe/dtos/request"
	"inventory-juanfe/services"
	"inventory-juanfe/utils"
)

type InventoryHandler struct {
	svc *services.InventoryService
}

func NewInventoryHandler(svc *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) ListPeriods(c fiber.Ctx) error {
	userId := utils.GetUserID(c)
	periods, err := h.svc.ListPeriods(userId)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch periods")
	}
	return utils.Success(c, http.StatusOK, periods)
}

func (h *InventoryHandler) CreatePeriod(c fiber.Ctx) error {
	var req dtos.CreatePeriodRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.ValidateCreatePeriod(req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	userID := utils.GetUserID(c)
	period, err := h.svc.CreatePeriod(req.PeriodYear, req.PeriodMonth, req.PeriodDay, userID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, http.StatusCreated, period)
}

func (h *InventoryHandler) ClosePeriod(c fiber.Ctx) error {
	id := c.Params("id")
	userID := utils.GetUserID(c)

	if err := h.svc.ClosePeriod(id, userID); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.SuccessMessage(c, http.StatusOK, "period closed")
}

func (h *InventoryHandler) GetRecords(c fiber.Ctx) error {
	periodID := c.Params("id")
	userId := utils.GetUserID(c)
	records, err := h.svc.GetRecords(periodID, userId)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch records")
	}
	return utils.Success(c, http.StatusOK, records)
}

func (h *InventoryHandler) RecordAsset(c fiber.Ctx) error {
	var req dtos.RecordAssetRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.ValidateRecordAsset(req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	userID := utils.GetUserID(c)
	if err := h.svc.RecordAsset(req, userID); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.SuccessMessage(c, http.StatusOK, "asset recorded")
}

func (h *InventoryHandler) GetProgress(c fiber.Ctx) error {
	periodID := c.Params("id")
	userId := utils.GetUserID(c)
	progress, err := h.svc.GetProgress(periodID, userId)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch progress")
	}
	return utils.Success(c, http.StatusOK, progress)
}

func (h *InventoryHandler) GetPeriodAssets(c fiber.Ctx) error {
	periodID := c.Params("id")
	if periodID == "" {
		return utils.Error(c, http.StatusBadRequest, "period id is required")
	}
	userId := utils.GetUserID(c)
	assets, err := h.svc.GetPeriodAssets(periodID, userId)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, http.StatusOK, assets)
}
