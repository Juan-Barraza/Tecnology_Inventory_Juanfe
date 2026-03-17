package utils

import (
	"math"
	"net/http"

	"github.com/gofiber/fiber/v3"

	response "inventory-juanfe/dtos/response"
)

// Success returns a standard JSON success response with data.
func Success(c fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(response.APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessMessage returns a standard JSON success response with a message only.
func SuccessMessage(c fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(response.APIResponse{
		Success: true,
		Message: msg,
	})
}

// SuccessPaginated returns a paginated JSON response using PaginatedResponse.
func SuccessPaginated(c fiber.Ctx, items interface{}, total, page, limit int) error {
	return c.Status(http.StatusOK).JSON(response.APIResponse{
		Success: true,
		Data: response.PaginatedResponse{
			Items:      items,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

// Error returns a standard JSON error response.
func Error(c fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(response.APIResponse{
		Success: false,
		Error:   msg,
	})
}
