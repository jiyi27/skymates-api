package handler

import (
	"encoding/json"
	"log"
	"net/http"
	v1 "skymates-api/internal/dto/v1"
)

type BaseHandler struct{}

// SendJSON 发送固定格式的 JSON 响应
func (h *BaseHandler) SendJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	// w.WriteHeader(code) will send the status code and response headers to the client immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := v1.Response{
		Code:    code,
		Message: message,
		Data:    data,
	}

	// 如果编码失败, 尝试发送一个简单的错误响应
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("handler.ResponseJSON: failed to encode response: %v", err)
		w.Write([]byte(`{"code":500,"message":"Internal server error"}`))
	}
}

// ReadJSON 从请求中读取 JSON 数据并解码到指定的结构体中
// 禁止未知字段, 当 JSON 中出现结构体里没有定义的字段时, 会导致 Decode() 报错
func (h *BaseHandler) ReadJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
