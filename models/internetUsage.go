package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InternetUsage struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ComputerID string    `gorm:"not null" json:"computer_id"`              // Foreign key to Computer
	Computer   Computer  `gorm:"foreignKey:ComputerID" json:"computer"`    // Relationship with Computer
	Domain     string    `gorm:"type:varchar(255);not null" json:"domain"` // e.g., "google.com"
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`                // When the domain was accessed
	CreatedAt  time.Time `json:"created_at"`
}

func (i *InternetUsage) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	if i.Timestamp.IsZero() {
		i.Timestamp = time.Now()
	}
	return nil
}
