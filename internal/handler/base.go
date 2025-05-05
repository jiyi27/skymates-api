package handler

import (
	"encoding/json"
	"log"
	"net/http"
	v1 "skymates-api/internal/dto/v1"
)

// BaseHandler 基础处理器
type BaseHandler struct{}

// ResponseJSON 写入JSON响应
func (h *BaseHandler) ResponseJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(v1.Response{
		Message: message,
		Data:    data,
	}); err != nil {
		log.Printf("handler.SendJSON: failed to encode response: %v", err)
	}
}

// DecodeJSON 解码JSON请求
func (h *BaseHandler) DecodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	// 禁止未知字段，不然只有json格式不对或者字段类型不匹配才会报错，但是字段缺失不会报错，而是直接赋值为零值
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// ValidateEmail 验证邮箱格式（从原代码中提取）
func (h *BaseHandler) ValidateEmail(email string) bool {
	// 实现邮箱验证逻辑
	return true // 简化实现，应该添加真正的验证逻辑
}
