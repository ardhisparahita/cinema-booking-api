package models

import "time"

type Studio struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TheaterID uint      `gorm:"not null;index" json:"theater_id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	TotalRows uint16    `gorm:"not null" json:"total_rows"`
	TotalCols uint16    `gorm:"not null" json:"total_cols"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Studio) TableName() string {
	return "studios"
}
