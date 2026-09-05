package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Email       string         `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	PhoneNumber string         `gorm:"size:20" json:"phone_number"`
	Role        string         `gorm:"type:enum('user','admin');not null;default:user" json:"role"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) TableName() string {
	return "users"
}
