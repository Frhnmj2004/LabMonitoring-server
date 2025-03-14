package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Alert struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ComputerID uuid.UUID `gorm:"type:uuid;not null" json:"computer_id"`
	Type       string    `gorm:"type:varchar(20);not null" json:"type"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	Timestamp  time.Time `gorm:"not null;index" json:"timestamp"`
	Resolved   bool      `gorm:"not null;default:false" json:"resolved"`
	Computer   Computer  `gorm:"foreignKey:ComputerID" json:"computer,omitempty"`
}

func (a *Alert) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}
	return nil
}
