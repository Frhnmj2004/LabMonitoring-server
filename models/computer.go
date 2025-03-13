package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Computer struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	IPAddress string    `gorm:"uniqueIndex;not null" json:"ip_address"`
	Logs      []ResourceLog `gorm:"foreignKey:ComputerID" json:"logs,omitempty"`
}

func (c *Computer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
