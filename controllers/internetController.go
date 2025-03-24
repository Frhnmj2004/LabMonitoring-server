package controllers

import (
	"github.com/Frhnmj2004/LabMonitoring-server/helper"
	"github.com/Frhnmj2004/LabMonitoring-server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetInternetUsage retrieves internet usage logs for a specific computer
func GetInternetUsage(c *fiber.Ctx) error {
	computerID := c.Query("computer_id")
	if computerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing computer_id parameter",
		})
	}

	compUUID, err := uuid.Parse(computerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid computer_id format",
		})
	}

	logs, err := helper.GetInternetUsageByComputer(compUUID, 100) // Limit to 100 logs
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve internet usage logs",
		})
	}

	return c.JSON(fiber.Map{
		"data": logs,
	})
}

// PostInternetUsage saves new internet usage data
func PostInternetUsage(c *fiber.Ctx) error {
	var usage models.InternetUsage
	if err := c.BodyParser(&usage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := helper.SaveInternetUsage(usage); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save internet usage",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Internet usage saved successfully",
	})
}
