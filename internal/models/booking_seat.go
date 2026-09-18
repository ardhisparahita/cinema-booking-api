package models

import "time"

type BookingSeat struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID  uint      `gorm:"not null;index" json:"booking_id"`
	ShowtimeID uint      `gorm:"not null;index" json:"showtime_id"`
	SeatID     uint      `gorm:"not null;index" json:"seat_id"`
	Price      float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	Seat       Seat      `gorm:"foreignKey:SeatID;references:ID" json:"seat,omitempty"`
}

func (bs *BookingSeat) TableName() string {
	return "booking_seats"
}
