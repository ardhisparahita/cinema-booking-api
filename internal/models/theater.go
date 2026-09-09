package models

import (
	"time"
)

type Theater struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:150;not null" json:"name"`
	City      string    `gorm:"size:100;not null" json:"city"`
	Address   string    `gorm:"size:255;not null" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Theater) TableName() string {
	return "theaters"
}
