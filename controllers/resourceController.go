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

type Alert struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ComputerID  uuid.UUID `json:"computer_id"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
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
		alert := Alert{
			ComputerID: data.ComputerID,
			Type:      "HIGH_CPU",
			Message:   "CPU usage exceeds 90%",
			Timestamp: time.Now(),
		}
		config.DB.Create(&alert)
		websocket.BroadcastResourceUpdate(fiber.Map{
			"type": "alert",
			"data": alert,
		})
	}

	if data.Memory > 90 {
		alert := Alert{
			ComputerID: data.ComputerID,
			Type:      "HIGH_MEMORY",
			Message:   "Memory usage exceeds 90%",
			Timestamp: time.Now(),
		}
		config.DB.Create(&alert)
		websocket.BroadcastResourceUpdate(fiber.Map{
			"type": "alert",
			"data": alert,
		})
	}

	// Broadcast update
	websocket.BroadcastResourceUpdate(fiber.Map{
		"type": "resource_update",
		"data": log,
	})

	return c.Status(fiber.StatusOK).JSON(log)
}

func GetHistory(c *fiber.Ctx) error {
	var logs []models.ResourceLog
	computerID := c.Query("computer_id")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	query := config.DB.Order("timestamp desc")

	if computerID != "" {
		query = query.Where("computer_id = ?", computerID)
	}

	if startTime != "" {
		query = query.Where("timestamp >= ?", startTime)
	}

	if endTime != "" {
		query = query.Where("timestamp <= ?", endTime)
	}

	if err := query.Find(&logs).Error; err != nil {
		utils.LogError("Failed to fetch resource history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch resource history",
		})
	}

	return c.JSON(logs)
}

func GetAlerts(c *fiber.Ctx) error {
	var alerts []Alert
	computerID := c.Query("computer_id")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	query := config.DB.Order("timestamp desc")

	if computerID != "" {
		query = query.Where("computer_id = ?", computerID)
	}

	if startTime != "" {
		query = query.Where("timestamp >= ?", startTime)
	}

	if endTime != "" {
		query = query.Where("timestamp <= ?", endTime)
	}

	if err := query.Find(&alerts).Error; err != nil {
		utils.LogError("Failed to fetch alerts: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch alerts",
		})
	}

	return c.JSON(alerts)
}
