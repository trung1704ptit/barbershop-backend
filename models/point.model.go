package models

import (
	"time"

	"github.com/google/uuid"
)

type Point struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	Points      int64     `gorm:"type:int;default:0;not null" json:"points,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreatePointRequest struct {
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type UpdatePointRequest struct {
	UserID      uuid.UUID `json:"user_id,omitempty"`
	Points      int64     `json:"points,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
