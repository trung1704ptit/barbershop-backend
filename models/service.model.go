package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name,omitempty"`
	Image       string    `gorm:"type:varchar(255);not null" json:"image,omitempty"`
	Price       float64   `json:"price,omitempty"`
	PriceText   string    `json:"price_text,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateServiceRequest struct {
	Name        string    `json:"name" binding:"required"`
	Image       string    `json:"image" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	PriceText   string    `json:"price_text,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type UpdateService struct {
	Name        string    `json:"name,omitempty"`
	Image       string    `json:"image,omitempty" `
	Price       float64   `json:"price,omitempty" `
	PriceText   string    `json:"price_text,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type UserService struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID          uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	ServiceID       uuid.UUID `gorm:"type:uuid" json:"service_id,omitempty"`
	CreatedAt       time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt       time.Time `gorm:"not null" json:"updated_at,omitempty"`
	UserIDServiceID string    `gorm:"uniqueIndex:user_id_service_id" json:"user_id_service_id,omitempty"`
}

type UserWithServicesRequest struct {
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	ServiceID uuid.UUID `gorm:"type:uuid" json:"service_id,omitempty"`
}
