package model

import "time"

// TermSummary 术语概要模型，仅包含 ID 和名称
type TermSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Term 术语模型，对应 terms 表
type Term struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Explanation string    `json:"explanation"`
	SourceURL   string    `json:"source_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TermDetail 术语详情模型，包含分类 ID 列表
type TermDetail struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Explanation string    `json:"explanation"`
	SourceURL   string    `json:"source_url"`
	CategoryIDs []int64   `json:"category_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
