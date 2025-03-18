package helper

import (
	"github.com/Frhnmj2004/LabMonitoring-server/models"
	//"gorm.io/gorm"
	"github.com/Frhnmj2004/LabMonitoring-server/config"
	"github.com/google/uuid"
)

func SaveInternetUsage(iu models.InternetUsage) error {
	return config.DB.Create(&iu).Error
}

// GetInternetUsageByComputer retrieves internet usage logs for a specific computer.
func GetInternetUsageByComputer(computerID uuid.UUID, limit int) ([]models.InternetUsage, error) {
	var logs []models.InternetUsage
	err := config.DB.Where("computer_id = ?", computerID).Order("timestamp DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
