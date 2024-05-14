package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	BarberID  uuid.UUID `gorm:"type:uuid;not null" json:"barber_id"`
	Barber    User      `gorm:"foreignKey:BarberID;references:ID" json:"barber"`
	GuestID   uuid.UUID `gorm:"type:uuid;not null" json:"guest_id"`
	Guest     User      `gorm:"foreignKey:GuestID;references:ID" json:"guest"`
	Status    string    `gorm:"type:varchar(100);not null" json:"status"`
	StartTime time.Time `gorm:"not null" json:"start_time,omitempty"`
	EndTime   time.Time `gorm:"not null" json:"end_time,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateBookingRequest struct {
	BarberID  uuid.UUID `gorm:"type:uuid;not null" json:"barber_id"`
	GuestID   uuid.UUID `gorm:"type:uuid;not null" json:"guest_id"`
	StartTime time.Time `gorm:"not null" json:"start_time,omitempty"`
	EndTime   time.Time `gorm:"not null" json:"end_time,omitempty"`
}
