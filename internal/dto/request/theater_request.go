package request

type CreateAndUpdateTheaterRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=150"`
	City    string `json:"city" validate:"required,min=2,max=100"`
	Address string `json:"address" validate:"required,min=5,max=255"`
}
