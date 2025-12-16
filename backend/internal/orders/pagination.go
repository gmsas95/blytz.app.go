package orders

import "math"

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize  int `json:"page_size"`
}

// GetOffset calculates database offset
func (p PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// PaginatedResponse represents paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	PageSize    int   `json:"page_size"`
	Total       int64  `json:"total"`
	TotalPages  int    `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total int64, page, pageSize int) *PaginatedResponse {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginatedResponse{
		Data: data,
		Pagination: Pagination{
			Page:      page,
			PageSize:   pageSize,
			Total:     total,
			TotalPages: totalPages,
		},
	}
}