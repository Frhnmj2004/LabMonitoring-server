package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ComputerID uuid.UUID `gorm:"type:uuid;not null" json:"computer_id"`
	Timestamp  time.Time `gorm:"not null;index" json:"timestamp"`
	CPU        float64   `gorm:"not null;check:cpu >= 0 AND cpu <= 100" json:"cpu"`
	Memory     float64   `gorm:"not null;check:memory >= 0 AND memory <= 100" json:"memory"`
	NetworkIn  float64   `gorm:"not null;check:network_in >= 0" json:"network_in"`
	NetworkOut float64   `gorm:"not null;check:network_out >= 0" json:"network_out"`
	Computer   Computer  `gorm:"foreignKey:ComputerID" json:"computer,omitempty"`
}

func (r *ResourceLog) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.Timestamp.IsZero() {
		r.Timestamp = time.Now()
	}
	return nil
}
