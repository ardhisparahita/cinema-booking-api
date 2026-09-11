package request

import "time"

type CreateAndUpdateShowtimeRequest struct {
	MovieID   uint      `json:"movie_id" validate:"required"`
	StudioID  uint      `json:"studio_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
	Price     float64   `json:"price" validate:"required,gt=0"`
}
