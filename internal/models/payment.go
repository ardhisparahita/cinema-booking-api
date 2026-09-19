package models

import "time"

type Payment struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID   uint       `gorm:"not null;index" json:"booking_id"`
	Method      string     `gorm:"size:30;not null" json:"method"`
	Amount      float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status      string     `gorm:"type:enum('pending','paid','failed','refunded');not null;default:pending" json:"status"`
	ProviderRef string     `gorm:"size:100" json:"provider_ref,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Booking Booking `gorm:"foreignKey:BookingID;references:ID" json:"booking,omitempty"`
}

func (c *Payment) TableName() string {
	return "payments"
}
