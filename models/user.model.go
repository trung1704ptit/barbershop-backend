package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"uniqueIndex"`
	Password  string    `gorm:"not null"`
	Phone     string    `gorm:"uniqueIndex;type:varchar(10);not null"`
	Birthday  time.Time `gorm:"not null"`
	Role      string    `gorm:"type:varchar(255);not null"`
	Provider  string    `gorm:"not null"`
	Services  []Service `gorm:"many2many:user_services;"`
	Points    []Point   `json:"points,omitempty"`
	Photo     string    `gorm:"not null"`
	Verified  bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type SignUpInput struct {
	Name            string    `json:"name" binding:"required"`
	Email           string    `json:"email"`
	Password        string    `json:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"passwordConfirm" binding:"required"`
	Phone           string    `json:"phone" binding:"required"`
	Birthday        time.Time `json:"birthday"`
	Photo           string    `json:"photo"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Birthday  time.Time `json:"birthday,omitempty"`
	Points    []Point   `json:"points,omitempty"`
	Services  []Service `json:"services,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	Name     string    `gorm:"type:varchar(255);not null"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Phone    string    `gorm:"uniqueIndex;not null"`
	Birthday time.Time `gorm:"not null"`
}
