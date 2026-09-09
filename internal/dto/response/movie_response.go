package response

import (
	"time"
)

type MovieResponse struct {
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	Synopsis    string          `json:"synopsis"`
	DurationMin uint16          `json:"duration_min"`
	Rating      string          `json:"rating"`
	PosterURL   string          `json:"poster_url"`
	ReleaseDate string          `json:"release_date"`
	IsActive    bool            `json:"is_active"`
	Genres      []GenreResponse `json:"genres"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
