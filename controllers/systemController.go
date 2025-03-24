package controllers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Frhnmj2004/LabMonitoring-server/config"
	"github.com/Frhnmj2004/LabMonitoring-server/models"
	"github.com/Frhnmj2004/LabMonitoring-server/utils"
	"github.com/gofiber/fiber/v2"
)

type SystemSignupRequest struct {
	College string `json:"college"`
	LabName string `json:"lab_name"`
}

// SystemSignup handles the registration of a new lab computer
func SystemSignup(c *fiber.Ctx) error {
	var req SystemSignupRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.College == "" || req.LabName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "College and Lab Name are required",
		})
	}

	// Create new computer record
	computer := &models.Computer{
		College: req.College,
		LabName: req.LabName,
	}

	// Register computer in database
	if err := computer.RegisterComputer(config.DB); err != nil {
		utils.LogError("Failed to register computer: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register computer",
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Computer registered successfully",
		"system_id": computer.ComputerID,
	})
}

// DownloadCollector serves the collector executable file
func DownloadCollector(c *fiber.Ctx) error {
	computerID := c.Query("ComputerID")
	if computerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Computer ID is required",
		})
	}

	// Verify computer ID exists
	computer, err := models.GetComputerBySystemID(config.DB, computerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invalid Computer ID",
		})
	}

	// Path to collector executable
	collectorPath := filepath.Join("downloads", "collector.exe")

	// Check if file exists
	if _, err := os.Stat(collectorPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Collector executable not found",
		})
	}

	// Set filename for download
	filename := fmt.Sprintf("collector_%s.exe", computer.ComputerID)
	
	// Serve the file
	return c.Download(collectorPath, filename)
}

// GetAllComputers returns a list of all registered computers
func GetAllComputers(c *fiber.Ctx) error {
	computers, err := models.GetAllComputers(config.DB)
	if err != nil {
		utils.LogError("Failed to fetch computers: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch computers",
		})
	}

	// Add online status to each computer
	var response []fiber.Map
	for _, computer := range computers {
		response = append(response, fiber.Map{
			"id":         computer.ID,
			"system_id":  computer.ComputerID,
			"college":    computer.College,
			"lab_name":   computer.LabName,
			"last_seen":  computer.LastSeen,
			"is_online":  computer.IsOnline(),
			"created_at": computer.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}
