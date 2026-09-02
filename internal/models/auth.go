package models

import "time"

type RefreshToken struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Token     string     `gorm:"size:255;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (r *RefreshToken) TableName() string {
	return "refresh_tokens"
}
