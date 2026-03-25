package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"inventory-juanfe/services"
	"inventory-juanfe/utils"
)

type DashboardHandler struct {
	svc *services.DashboardService
}

func NewDashboardHandler(svc *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetDashboard(c fiber.Ctx) error {
	userID := utils.GetUserID(c)
	data, err := h.svc.GetDashboard(userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "could not fetch dashboard data")
	}
	return utils.Success(c, http.StatusOK, data)
}
