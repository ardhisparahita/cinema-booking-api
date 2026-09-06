package request

type CreateGenreRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UpdateGenreRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}
