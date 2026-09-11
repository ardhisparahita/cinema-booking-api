package request

type CreateAndUpdateSeatRequest struct {
	RowLabel  string `json:"row_label" validate:"required,max=2"`
	ColNumber uint16 `json:"col_number" validate:"required,min=1"`
	SeatType  string `json:"seat_type" validate:"required,oneof=regular vip"`
}
