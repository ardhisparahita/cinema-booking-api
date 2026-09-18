package models

import "time"

type Booking struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingCode string     `gorm:"size:20;not null;uniqueIndex" json:"booking_code"`
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	ShowtimeID  uint       `gorm:"not null;index" json:"showtime_id"`
	TotalPrice  float64    `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status      string     `gorm:"type:enum('pending','confirmed','cancelled','expired');not null;default:pending" json:"status"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	BookingSeats []BookingSeat `gorm:"foreignKey:BookingID;references:ID" json:"booking_seats,omitempty"`
}

func (b *Booking) TableName() string {
	return "bookings"
}
