package models

import "time"

type Seat struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StudioID  uint      `gorm:"not null;index" json:"studio_id"`
	RowLabel  string    `gorm:"size:2;not null" json:"row_label"`
	ColNumber uint16    `gorm:"not null" json:"col_number"`
	SeatType  string    `gorm:"type:enum('regular','vip');not null;default:regular" json:"seat_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Seat) TableName() string {
	return "seats"
}
