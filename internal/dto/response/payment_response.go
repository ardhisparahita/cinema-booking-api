package response

import "time"

type PaymentResponse struct {
	ID          uint       `json:"id"`
	BookingID   uint       `json:"booking_id"`
	Method      string     `json:"method"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"`
	ProviderRef string     `json:"provider_ref,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
