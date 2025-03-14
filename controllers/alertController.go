package controllers

import (
	"time"

	"github.com/Frhnmj2004/LabMonitoring-server/config"
	"github.com/Frhnmj2004/LabMonitoring-server/models"
	"github.com/Frhnmj2004/LabMonitoring-server/utils"
	"github.com/Frhnmj2004/LabMonitoring-server/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAlerts returns all alerts with optional filtering
func GetAlerts(c *fiber.Ctx) error {
	var alerts []models.Alert
	query := config.DB.Order("timestamp desc")

	// Apply filters
	if computerID := c.Query("computer_id"); computerID != "" {
		query = query.Where("computer_id = ?", computerID)
	}

	if alertType := c.Query("type"); alertType != "" {
		query = query.Where("type = ?", alertType)
	}

	if resolved := c.Query("resolved"); resolved != "" {
		query = query.Where("resolved = ?", resolved == "true")
	}

	// Add pagination
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	var total int64
	if err := query.Model(&models.Alert{}).Count(&total).Error; err != nil {
		utils.LogError("Failed to count alerts: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch alerts",
		})
	}

	if err := query.Limit(limit).Offset(offset).Find(&alerts).Error; err != nil {
		utils.LogError("Failed to fetch alerts: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch alerts",
		})
	}

	return c.JSON(fiber.Map{
		"data": alerts,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
			"total_items": total,
			"per_page":    limit,
		},
	})
}

// GetActiveAlerts returns all unresolved alerts
func GetActiveAlerts(c *fiber.Ctx) error {
	var alerts []models.Alert
	query := config.DB.Where("resolved = ?", false).Order("timestamp desc")

	// Apply computer filter if provided
	if computerID := c.Query("computer_id"); computerID != "" {
		query = query.Where("computer_id = ?", computerID)
	}

	if err := query.Find(&alerts).Error; err != nil {
		utils.LogError("Failed to fetch active alerts: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch active alerts",
		})
	}

	return c.JSON(fiber.Map{
		"data": alerts,
	})
}

// GetAlertHistory returns alert history with pagination
func GetAlertHistory(c *fiber.Ctx) error {
	var alerts []models.Alert
	query := config.DB.Order("timestamp desc")

	// Apply filters
	if computerID := c.Query("computer_id"); computerID != "" {
		query = query.Where("computer_id = ?", computerID)
	}

	if alertType := c.Query("type"); alertType != "" {
		query = query.Where("type = ?", alertType)
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
	if err := query.Model(&models.Alert{}).Count(&total).Error; err != nil {
		utils.LogError("Failed to count alert history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch alert history",
		})
	}

	if err := query.Limit(limit).Offset(offset).Find(&alerts).Error; err != nil {
		utils.LogError("Failed to fetch alert history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch alert history",
		})
	}

	return c.JSON(fiber.Map{
		"data": alerts,
		"pagination": fiber.Map{
			"current_page": page,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
			"total_items": total,
			"per_page":    limit,
		},
	})
}

// ResolveAlert marks an alert as resolved
func ResolveAlert(c *fiber.Ctx) error {
	alertID := c.Params("id")
	if alertID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Alert ID is required",
		})
	}

	// Parse UUID
	id, err := uuid.Parse(alertID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid alert ID format",
		})
	}

	var alert models.Alert
	if err := config.DB.First(&alert, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Alert not found",
		})
	}

	alert.Resolved = true
	if err := config.DB.Save(&alert).Error; err != nil {
		utils.LogError("Failed to resolve alert: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to resolve alert",
		})
	}

	// Broadcast alert resolution
	websocket.BroadcastResourceUpdate(fiber.Map{
		"type": "alert_resolved",
		"data": alert,
	})

	return c.JSON(fiber.Map{
		"message": "Alert resolved successfully",
		"data":    alert,
	})
}

// GetAlertStats returns alert statistics
func GetAlertStats(c *fiber.Ctx) error {
	var stats struct {
		TotalAlerts      int64 `json:"total_alerts"`
		ActiveAlerts     int64 `json:"active_alerts"`
		ResolvedAlerts   int64 `json:"resolved_alerts"`
		HighCPUAlerts    int64 `json:"high_cpu_alerts"`
		HighMemoryAlerts int64 `json:"high_memory_alerts"`
	}

	db := config.DB.Model(&models.Alert{})

	// Get total alerts
	db.Count(&stats.TotalAlerts)

	// Get active alerts
	db.Where("resolved = ?", false).Count(&stats.ActiveAlerts)

	// Get resolved alerts
	db.Where("resolved = ?", true).Count(&stats.ResolvedAlerts)

	// Get alerts by type
	db.Where("type = ?", "HIGH_CPU").Count(&stats.HighCPUAlerts)
	db.Where("type = ?", "HIGH_MEMORY").Count(&stats.HighMemoryAlerts)

	// Add time-based statistics
	var last24Hours struct {
		Count int64
	}
	db.Where("timestamp >= ?", time.Now().Add(-24*time.Hour)).Count(&last24Hours.Count)

	return c.JSON(fiber.Map{
		"total_stats": stats,
		"last_24h": fiber.Map{
			"total_alerts": last24Hours.Count,
		},
	})
}
