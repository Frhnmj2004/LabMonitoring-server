package controllers

import (
	"time"

	"github.com/Frhnmj2004/LabMonitoring-server/config"
	"github.com/Frhnmj2004/LabMonitoring-server/models"
	"github.com/Frhnmj2004/LabMonitoring-server/storage"
	"github.com/Frhnmj2004/LabMonitoring-server/utils"
	"github.com/Frhnmj2004/LabMonitoring-server/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ResourceData struct {
	ComputerID  uuid.UUID `json:"computer_id"`
	CPU         float64   `json:"cpu"`
	Memory      float64   `json:"memory"`
	NetworkIn   float64   `json:"network_in"`
	NetworkOut  float64   `json:"network_out"`
}

func PostResource(c *fiber.Ctx) error {
	var data ResourceData
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate data
	if data.CPU < 0 || data.CPU > 100 || data.Memory < 0 || data.Memory > 100 || data.NetworkIn < 0 || data.NetworkOut < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid resource values",
		})
	}

	// Create resource log
	log := &models.ResourceLog{
		ComputerID:  data.ComputerID,
		CPU:         data.CPU,
		Memory:      data.Memory,
		NetworkIn:   data.NetworkIn,
		NetworkOut:  data.NetworkOut,
		Timestamp:   time.Now(),
	}

	// Try to save to database
	if err := config.DB.Create(log).Error; err != nil {
		// If database save fails, store in buffer
		utils.LogWarning("Failed to save to database, writing to buffer: %v", err)
		if bufferErr := storage.WriteBuffer(log); bufferErr != nil {
			utils.LogError("Failed to write to buffer: %v", bufferErr)
		}
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Service temporarily unavailable, data buffered",
		})
	}

	// Check for alerts
	if data.CPU > 90 {
		alert := &models.Alert{
			ComputerID: data.ComputerID,
			Type:      "HIGH_CPU",
			Message:   "CPU usage exceeds 90%",
			Timestamp: time.Now(),
			Resolved:  false,
		}
		if err := config.DB.Create(alert).Error; err != nil {
			utils.LogError("Failed to create CPU alert: %v", err)
		} else {
			websocket.BroadcastResourceUpdate(fiber.Map{
				"type": "alert",
				"data": alert,
			})
		}
	}

	if data.Memory > 90 {
		alert := &models.Alert{
			ComputerID: data.ComputerID,
			Type:      "HIGH_MEMORY",
			Message:   "Memory usage exceeds 90%",
			Timestamp: time.Now(),
			Resolved:  false,
		}
		if err := config.DB.Create(alert).Error; err != nil {
			utils.LogError("Failed to create memory alert: %v", err)
		} else {
			websocket.BroadcastResourceUpdate(fiber.Map{
				"type": "alert",
				"data": alert,
			})
		}
	}

	// Broadcast resource update
	websocket.BroadcastResourceUpdate(fiber.Map{
		"type": "resource_update",
		"data": log,
	})

	return c.Status(fiber.StatusOK).JSON(log)
}

func GetHistory(c *fiber.Ctx) error {
	var logs []models.ResourceLog
	query := config.DB.Order("timestamp desc")

	// Apply filters
	if computerID := c.Query("computer_id"); computerID != "" {
		query = query.Where("computer_id = ?", computerID)
	}

	if startTime := c.Query("start_time"); startTime != "" {
		query = query.Where("timestamp >= ?", startTime)
	}

	if endTime := c.Query("end_time"); endTime != "" {
		query = query.Where("timestamp <= ?", endTime)
	}

	// Add pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	var total int64
	if err := query.Model(&models.ResourceLog{}).Count(&total).Error; err != nil {
		utils.LogError("Failed to count resource logs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch resource history",
		})
	}

	if err := query.Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		utils.LogError("Failed to fetch resource logs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch resource history",
		})
	}

	return c.JSON(fiber.Map{
		"data": logs,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
			"total_items": total,
			"per_page":    limit,
		},
	})
}
