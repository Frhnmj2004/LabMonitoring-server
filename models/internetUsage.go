package models

import (
	"github.com/google/uuid"
	"time"
)

type InternetUsage struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    ComputerID uuid.UUID `gorm:"type:uuid;not null" json:"computer_id"` // Foreign key to Computer
    Computer   Computer  `gorm:"foreignKey:ComputerID" json:"computer"` // Relationship with Computer
    Domain     string    `gorm:"type:varchar(255);not null" json:"domain"` // e.g., "google.com"
    Timestamp  time.Time `gorm:"not null" json:"timestamp"` // When the domain was accessed
    CreatedAt  time.Time `json:"created_at"`
}