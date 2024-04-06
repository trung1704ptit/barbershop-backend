package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Image       string    `gorm:"type:varchar(255);not null" json:"image,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
