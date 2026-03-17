package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	dtos "inventory-juanfe/dtos/request"
	"inventory-juanfe/services"
	"inventory-juanfe/utils"
)

type AssignmentHandler struct {
	svc *services.AssignmentService
}

func NewAssignmentHandler(svc *services.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{svc: svc}
}

func (h *AssignmentHandler) Create(c fiber.Ctx) error {
	var req dtos.CreateAssignmentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.ValidateCreateAssignment(req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	userID := utils.GetUserID(c)
	assignment, err := h.svc.Create(req, userID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, http.StatusCreated, assignment)
}

func (h *AssignmentHandler) Release(c fiber.Ctx) error {
	id := c.Params("id")

	var req dtos.ReleaseAssignmentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.ValidateReleaseAssignment(req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	if err := h.svc.Release(id, req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessMessage(c, http.StatusOK, "assignment released")
}
