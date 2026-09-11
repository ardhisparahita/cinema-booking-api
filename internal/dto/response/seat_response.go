package response

import "time"

type SeatResponse struct {
	ID        uint      `json:"id"`
	StudioID  uint      `json:"studio_id"`
	RowLabel  string    `json:"row_label"`
	ColNumber uint16    `json:"col_number"`
	SeatType  string    `json:"seat_type"`
	CreatedAt time.Time `json:"created_at"`
}
