package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func GetInternetUsage(c *fiber.Ctx) error {
	computerID := c.Query("computer_id")
	if computerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing computer_id parameter",
		})
	}

	return c.JSON(fiber.Map{
		"error": "Not implemented",
	})
}

func PostInternetUsage(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": "Not implemented",
	})
}
