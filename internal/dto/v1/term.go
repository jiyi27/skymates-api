package v1

import (
	"github.com/google/uuid"
)

// ListTermsRequest 列出术语请求
type ListTermsRequest struct {
	LastID *uuid.UUID `json:"last_id"`
	Limit  int        `json:"limit"`
}

// ListTermsResponse 列出术语响应
type ListTermsResponse struct {
	Terms   []TermResponse `json:"terms"`    // 术语列表
	LastID  *uuid.UUID     `json:"last_id"`  // 用于客户端下次查询
	HasMore bool           `json:"has_more"` // 是否还有更多数据，用于前端滚动加载
}

// TermResponse 术语响应
type TermResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
