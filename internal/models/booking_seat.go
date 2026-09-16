package models

type BookingSeat struct {
	ID         uint `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID  uint `gorm:"not null;index" json:"booking_id"`
	ShowtimeID uint `gorm:"not null;index" json:"showtime_id"`
	SeatID     uint `gorm:"not null;index" json:"seat_id"`
	Price      uint `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt  uint `json:"created_at"`
}

func (bs *BookingSeat) TableName() string {
	return "booking_seats"
}
