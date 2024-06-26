package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Service struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name,omitempty"`
	Image       string         `gorm:"type:varchar(255);not null" json:"image,omitempty"`
	Price       float64        `json:"price,omitempty"`
	PriceText   string         `json:"price_text,omitempty"`
	Todos       pq.StringArray `gorm:"type:text[]" json:"todos,omitempty"`
	Category    string         `gorm:"type:text" json:"category,omitempty"`
	ServiceType string         `gorm:"type:text;default:'one_time'" json:"service_type,omitempty"`
	Limit       int64          `json:"limit,omitempty"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateServiceRequest struct {
	Name        string         `json:"name" binding:"required"`
	Image       string         `json:"image" binding:"required"`
	Price       float64        `json:"price" binding:"required"`
	PriceText   string         `json:"price_text,omitempty"`
	Description string         `json:"description,omitempty"`
	Todos       pq.StringArray `json:"todos,omitempty"`
	Category    string         `gorm:"type:text" json:"category,omitempty"`
	ServiceType string         `json:"service_type,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
}

type UpdateService struct {
	Name        string         `json:"name,omitempty"`
	Image       string         `json:"image,omitempty" `
	Price       float64        `json:"price,omitempty" `
	PriceText   string         `json:"price_text,omitempty"`
	Todos       pq.StringArray `json:"todos,omitempty"`
	Category    string         `gorm:"type:text" json:"category,omitempty"`
	Description string         `json:"description,omitempty"`
	ServiceType string         `json:"service_type,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
}

type UserService struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID          uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	ServiceID       uuid.UUID `gorm:"type:uuid" json:"service_id,omitempty"`
	CreatedAt       time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt       time.Time `gorm:"not null" json:"updated_at,omitempty"`
	UserIDServiceID string    `gorm:"uniqueIndex:user_id_service_id" json:"user_id_service_id,omitempty"`
}

type UserServiceRequest struct {
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	ServiceID uuid.UUID `gorm:"type:uuid" json:"service_id,omitempty"`
}

type ServiceHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	ServiceID uuid.UUID `gorm:"type:uuid" json:"service_id,omitempty"`
	Count     int       `gorm:"not null" json:"count,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
	User      User      `gorm:"foreignkey:UserID" json:"-"`
	Service   Service   `gorm:"foreignkey:ServiceID" json:"-"`
}
