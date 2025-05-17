package models

type ListResponse[T any] struct {
	TotalRecords int64 `json:"total_records"`
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	Data         []T   `json:"data"`
}
