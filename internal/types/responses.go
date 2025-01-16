package types

import "github.com/google/uuid"

// Response 通用响应结构
// 当 data 为 nil 时, 序列化时会忽略这个字段, 得到的 json 是 {"message": "xxx"}, 无 data 字段
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ListTermsResponse 用户查询某一分类下的术语列表响应
type ListTermsResponse struct {
	Terms   []Term     `json:"terms"`    // 术语列表
	LastID  *uuid.UUID `json:"last_id"`  // 用于客户端下次查询
	HasMore bool       `json:"has_more"` // 是否还有更多数据, 用于前端滚动加载
}
