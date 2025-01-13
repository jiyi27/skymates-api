package types

import "github.com/google/uuid"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ListTermsResponse 用户查询某一分类下的术语列表响应
type ListTermsResponse struct {
	Terms   []Term     `json:"terms"`    // 术语列表
	LastID  *uuid.UUID `json:"last_id"`  // 用于客户端下次查询
	HasMore bool       `json:"has_more"` // 是否还有更多数据, 用于前端滚动加载
}
