package request

type PaginationRequest struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}
