package v1

// Response 通用响应结构
// 若使用 `json:"data,omitempty"`, 则表示
// 当 data 为 nil 时序列化时会忽略这个字段, 得到的 json 无 data 字段
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
