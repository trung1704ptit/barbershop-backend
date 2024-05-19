package models

import (
	"time"

	"github.com/google/uuid"
)

type Gallery struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	ImageURL  string    `gorm:"type:text;not null" json:"image_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
type GalleryRequest struct {
	ImageURL string `gorm:"not null" json:"image_url,omitempty"`
}
