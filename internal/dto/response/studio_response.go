package response

import "time"

type StudioResponse struct {
	ID        uint      `json:"id"`
	TheaterID uint      `json:"theater_id"`
	Name      string    `json:"name"`
	TotalCols uint16    `json:"total_cols"`
	TotalRows uint16    `json:"total_rows"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
