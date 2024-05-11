package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title        string    `gorm:"uniqueIndex;not null" json:"title,omitempty"`
	Slug         string    `gorm:"uniqueIndex;not null" json:"slug,omitempty"`
	Description  string    `gorm:"not null" json:"description,omitempty"`
	Price        float64   `gorm:"uniqueIndex;not null" json:"price,omitempty"`
	PriceText    string    `gorm:"uniqueIndex;not null" json:"price_text,omitempty"`
	PreviewImage string    `gorm:"not null" json:"preview_image,omitempty"`
	Images       []string  `gorm:"uniqueIndex" json:"images"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateProductRequest struct {
	Title        string    `gorm:"uniqueIndex;not null" json:"title,omitempty"`
	Slug         string    `gorm:"uniqueIndex;not null" json:"slug,omitempty"`
	Description  string    `gorm:"not null" json:"description,omitempty"`
	Price        float64   `gorm:"uniqueIndex;not null" json:"price,omitempty"`
	PriceText    string    `gorm:"uniqueIndex;not null" json:"price_text,omitempty"`
	PreviewImage string    `gorm:"not null" json:"preview_image,omitempty"`
	Images       []string  `gorm:"uniqueIndex" json:"images"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type UpdateProduct struct {
	Title        string    `gorm:"uniqueIndex;not null" json:"title,omitempty"`
	Slug         string    `gorm:"uniqueIndex;not null" json:"slug,omitempty"`
	Description  string    `gorm:"not null" json:"description,omitempty"`
	Price        float64   `gorm:"uniqueIndex;not null" json:"price,omitempty"`
	PriceText    string    `gorm:"uniqueIndex;not null" json:"price_text,omitempty"`
	PreviewImage string    `gorm:"not null" json:"preview_image,omitempty"`
	Images       []string  `gorm:"uniqueIndex" json:"images"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
