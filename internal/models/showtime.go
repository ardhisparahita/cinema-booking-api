package models

import "time"

type Showtime struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	MovieID   uint      `gorm:"not null;index" json:"movie_id"`
	StudioID  uint      `gorm:"not null;index" json:"studio_id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Showtime) TableName() string {
	return "showtimes"
}
