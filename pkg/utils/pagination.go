package utils

type Pagination struct {
	Page   int
	Limit  int
	Offset int
}

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

func GetPagination(page, limit int) Pagination {
	if page < 1 {
		page = DefaultPage
	}

	if limit < 1 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Pagination{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func GetTotalPage(total int64, limit int) int {
	if total == 0 {
		return 0
	}

	return int((total + int64(limit) - 1) / int64(limit))
}
