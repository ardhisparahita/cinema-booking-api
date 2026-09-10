package request

type CreateAndUpdateStudioRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=50"`
	TotalCols uint16 `json:"total_rows" validate:"required,min=1"`
	TotalRows uint16 `json:"total_cols" validate:"required,min=1"`
}
