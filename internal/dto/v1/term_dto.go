package v1

import "time"

// SearchTermsResponse 搜索术语的响应 DTO
type SearchTermsResponse struct {
	Terms []TermSummary `json:"terms"`
}

// TermSummary 术语概要 DTO
type TermSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TermDetailResponse 术语详情的响应 DTO
type TermDetailResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Explanation string    `json:"explanation"`
	SourceURL   string    `json:"source_url"`
	CategoryIDs []int64   `json:"category_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListTermsByCategoryResponse 列出分类下术语的响应 DTO
type ListTermsByCategoryResponse struct {
	Terms   []TermSummary `json:"terms"`
	HasMore bool          `json:"has_more"`
}

// CreateTermRequest 创建术语的请求 DTO
type CreateTermRequest struct {
	Name        string  `json:"name" validate:"required"`
	Explanation string  `json:"explanation" validate:"required"`
	SourceURL   string  `json:"source_url"`
	CategoryIDs []int64 `json:"category_ids"`
}

// UpdateTermRequest 更新术语的请求 DTO
type UpdateTermRequest struct {
	Name        string  `json:"name" validate:"required"`
	Explanation string  `json:"explanation" validate:"required"`
	SourceURL   string  `json:"source_url"`
	CategoryIDs []int64 `json:"category_ids"`
}
