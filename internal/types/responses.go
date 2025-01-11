package types

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	// omitempty tag tells the JSON encoder Data field is optional
	Data interface{} `json:"data,omitempty"`
}

// ListTermsResponse 用户查询某一分类下的术语列表响应
type ListTermsResponse struct {
	Terms   []Term `json:"terms"`    // 术语列表
	HasMore bool   `json:"has_more"` // 是否还有更多数据, 用于前端滚动加载
}
