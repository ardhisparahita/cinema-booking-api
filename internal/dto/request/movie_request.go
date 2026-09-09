package request

type CreateMovieRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Synopsis    string `json:"synopsis"`
	DurationMin uint16 `json:"duration_min" validate:"required,min=1"`
	Rating      string `json:"rating" validate:"required,oneof=SU 13+ 17+ 21+"`
	PosterURL   string `json:"poster_url"`
	ReleaseDate string `json:"release_date"`
	GenreIDs    []uint `json:"genre_ids" validate:"required,min=1"`
}

type UpdateMovieRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Synopsis    string `json:"synopsis"`
	DurationMin uint16 `json:"duration_min" validate:"required,min=1"`
	Rating      string `json:"rating" validate:"required,oneof=SU 13+ 17+ 21+"`
	PosterURL   string `json:"poster_url"`
	ReleaseDate string `json:"release_date"`
	IsActive    bool   `json:"is_active"`
	GenreIDs    []uint `json:"genre_ids" validate:"required,min=1"`
}
