package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name            string           `gorm:"type:varchar(255);not null"`
	Email           string           `gorm:"uniqueIndex"`
	Position        string           `gorm:"type:varchar(255)" json:"position,omitempty"`
	Intro           string           `gorm:"type:varchar(255)" json:"intro,omitempty"`
	Password        string           `gorm:"not null"`
	Phone           string           `gorm:"uniqueIndex;type:varchar(10);not null"`
	Birthday        time.Time        `gorm:"not null"`
	Roles           pq.StringArray   `gorm:"type:text[];not null" json:"roles,omitempty"`
	Provider        string           `gorm:"not null"`
	Services        []Service        `gorm:"many2many:user_services;" json:"services,omitempty"`
	ServicesHistory []ServiceHistory `gorm:"foreignKey:UserID;" json:"services_history,omitempty"`
	Points          []Point          `gorm:"foreignKey:UserID;" json:"points,omitempty"`
	Photo           string           `gorm:"force" json:"photo,omitempty"`
	Verified        bool             `gorm:"not null"`
	CreatedAt       time.Time        `gorm:"not null"`
	UpdatedAt       time.Time        `gorm:"not null"`
}

type SignUpInput struct {
	Name            string         `json:"name" binding:"required"`
	Email           string         `json:"email"`
	Password        string         `json:"password" binding:"required,min=8"`
	PasswordConfirm string         `json:"passwordConfirm" binding:"required"`
	Phone           string         `json:"phone" binding:"required"`
	Birthday        time.Time      `json:"birthday"`
	Photo           string         `gorm:"force" json:"photo,omitempty"`
	Roles           pq.StringArray `gorm:"type:text[];not null" json:"roles" binding:"required"`
	Position        string         `json:"position"`
	Intro           string         `json:"intro"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID              uuid.UUID        `json:"id,omitempty"`
	Name            string           `json:"name,omitempty"`
	Email           string           `json:"email,omitempty"`
	Phone           string           `json:"phone,omitempty"`
	Position        string           `json:"position,omitempty"`
	Intro           string           `json:"intro,omitempty"`
	Birthday        time.Time        `json:"birthday,omitempty"`
	Points          []Point          `json:"points,omitempty"`
	Services        []Service        `json:"services,omitempty"`
	ServicesHistory []ServiceHistory `json:"services_history"`
	Roles           pq.StringArray   `json:"roles,omitempty"`
	Photo           string           `gorm:"force" json:"photo,omitempty"`
	Provider        string           `json:"provider"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type UpdateUserRequest struct {
	Name     string         `gorm:"type:varchar(255);not null"`
	Email    string         `gorm:"uniqueIndex;not null"`
	Phone    string         `gorm:"uniqueIndex;not null"`
	Birthday time.Time      `gorm:"not null"`
	Position string         `json:"position,omitempty"`
	Intro    string         `json:"intro,omitempty"`
	Roles    pq.StringArray `json:"roles,omitempty"`
	Photo    string         `gorm:"force" json:"photo,omitempty"`
}
