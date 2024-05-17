package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	BarberID    uuid.UUID `gorm:"type:uuid;not null" json:"barber_id"`
	Barber      User      `gorm:"foreignKey:BarberID;references:ID" json:"barber"`
	GuestID     uuid.UUID `gorm:"type:uuid;not null" json:"guest_id"`
	Guest       User      `gorm:"foreignKey:GuestID;references:ID" json:"guest"`
	Status      string    `gorm:"type:varchar(100);not null;default:'open'" json:"status"`
	Services    []Service `gorm:"many2many:booking_services;" json:"services"`
	BookingTime time.Time `gorm:"not null" json:"booking_time,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateBookingRequest struct {
	BarberID    uuid.UUID   `json:"barber_id" binding:"required"`
	GuestID     uuid.UUID   `json:"guest_id" binding:"required"`
	BookingTime time.Time   `json:"booking_time" binding:"required"`
	ServiceIDs  []uuid.UUID `json:"service_ids" binding:"required"`
}

type BookingService struct {
	BookingID uuid.UUID `gorm:"type:uuid" json:"booking_id,omitempty"`
	ServiceID uuid.UUID `gorm:"type:uuid" json:"service_id,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
