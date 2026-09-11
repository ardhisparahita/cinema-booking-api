package response

import "time"

type ShowtimeResponse struct {
	ID        uint      `json:"id"`
	MovieID   uint      `json:"movie_id"`
	StudioID  uint      `json:"studio_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
