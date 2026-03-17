package handlers

import (
	"github.com/gofiber/fiber/v3"

	dtos "inventory-juanfe/dtos/request"
	"inventory-juanfe/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dtos.LoginRequest

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Email and password are required",
		})
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    resp,
	})
}

func (h *AuthHandler) Me(c fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id": userID,
			"email":   c.Locals("userEmail"),
		},
	})
}
