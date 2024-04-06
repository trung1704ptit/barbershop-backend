package models

import (
	"time"

	"github.com/google/uuid"
)

type Point struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	User        uuid.UUID `gorm:"not null" json:"user,omitempty"`
	Points      int64     `gorm:"type:int;not null" json:"points,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreatePointRequest struct {
	User        uuid.UUID `json:"user,omitempty"`
	Points      int64     `json:"points,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type UpdatePointRequest struct {
	User        uuid.UUID `json:"user,omitempty"`
	Points      int64     `json:"points,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
