package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name            string           `gorm:"type:varchar(255);not null" json:"name"`
	Email           string           `gorm:"uniqueIndex" json:"email"`
	Position        string           `gorm:"type:varchar(255)" json:"position,omitempty"`
	Intro           string           `gorm:"type:varchar(255)" json:"intro,omitempty" `
	Password        string           `gorm:"not null" json:"-"`
	Phone           string           `gorm:"uniqueIndex;type:varchar(10);not null" json:"phone"`
	Birthday        time.Time        `json:"birthday,omitempty"`
	Roles           pq.StringArray   `gorm:"type:text[];not null" json:"roles,omitempty"`
	Provider        string           `gorm:"not null" json:"provider"`
	Services        []Service        `gorm:"many2many:user_services;" json:"services,omitempty"`
	ServicesHistory []ServiceHistory `gorm:"foreignKey:UserID;" json:"services_history,omitempty"`
	Points          []Point          `gorm:"foreignKey:UserID;" json:"points,omitempty"`
	Photo           string           `gorm:"force" json:"photo,omitempty"`
	Verified        bool             `gorm:"not null" json:"verified"`
	CreatedAt       time.Time        `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"not null" json:"updated_at"`
}

type SignUpInput struct {
	Name            string         `json:"name" binding:"required"`
	Email           string         `json:"email"`
	Password        string         `json:"password" binding:"required,min=8"`
	PasswordConfirm string         `json:"passwordConfirm" binding:"required"`
	Phone           string         `json:"phone" binding:"required"`
	Birthday        time.Time      `json:"birthday,omitempty"`
	Photo           string         `gorm:"force" json:"photo,omitempty"`
	Roles           pq.StringArray `gorm:"type:text[];not null" json:"roles" binding:"required"`
	Position        string         `json:"position"`
	Intro           string         `json:"intro"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UpdateUserRequest struct {
	Name     string         `gorm:"type:varchar(255);not null"`
	Email    string         `gorm:"uniqueIndex;not null"`
	Phone    string         `gorm:"uniqueIndex;not null"`
	Birthday time.Time      `json:"birthday,omitempty"`
	Position string         `json:"position,omitempty"`
	Intro    string         `json:"intro,omitempty"`
	Roles    pq.StringArray `json:"roles,omitempty"`
	Photo    string         `gorm:"force" json:"photo,omitempty"`
}
