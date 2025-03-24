package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Computer struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ComputerID string    `json:"computer_id" gorm:"uniqueIndex;not null;column:computer_id"`
	College    string    `json:"college" gorm:"not null"`
	LabName    string    `json:"lab_name" gorm:"not null"`
	LastSeen   time.Time `json:"last_seen"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BeforeCreate is a GORM hook that runs before creating a new computer
func (c *Computer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// RegisterComputer creates a new computer record in the database
func (c *Computer) RegisterComputer(db *gorm.DB) error {
	return db.Create(c).Error
}

// GetAllComputers returns all registered computers
func GetAllComputers(db *gorm.DB) ([]Computer, error) {
	var computers []Computer
	err := db.Order("created_at desc").Find(&computers).Error
	return computers, err
}

// GetComputerBySystemID finds a computer by its SystemID
func GetComputerBySystemID(db *gorm.DB, computerID string) (*Computer, error) {
	var computer Computer
	err := db.Where("computer_id = ?", computerID).First(&computer).Error
	return &computer, err
}

// UpdateLastSeen updates the LastSeen timestamp for a computer
func (c *Computer) UpdateLastSeen(db *gorm.DB) error {
	c.LastSeen = time.Now()
	return db.Model(c).Update("last_seen", c.LastSeen).Error
}

// IsOnline returns true if the computer has been seen in the last 5 minutes
func (c *Computer) IsOnline() bool {
	return time.Since(c.LastSeen) <= 5*time.Minute
}
