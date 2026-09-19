package request

type CreatePaymentRequest struct {
	BookingID uint   `json:"booking_Id" validate:"required"`
	Method    string `json:"method" validate:"required"`
}
