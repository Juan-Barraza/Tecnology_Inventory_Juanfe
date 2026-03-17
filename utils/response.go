package utils

import "github.com/gofiber/fiber/v3"

func GetUserID(c fiber.Ctx) string {
	val, _ := c.Locals("userID").(string)
	return val
}

func GetUserEmail(c fiber.Ctx) string {
	val, _ := c.Locals("userEmail").(string)
	return val
}
