package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"skymates-api/internal/types"
)

type BaseHandler struct{}

// ResponseJSON writes JSON response with given status code and data
func (h *BaseHandler) ResponseJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	// w.WriteHeader(code) will send the status code and response headers to the client immediately
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(types.Response{
		Code:    code,
		Message: message,
		Data:    data,
	}); err != nil {
		log.Printf("handler.ResponseJSON: failed to encode response: %v", err)
	}
}

func (h *BaseHandler) DecodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	// 禁止未知字段, 不然只有 json 格式不对或者字段类型不匹配才会报错, 但是字段缺失不会报错, 而是直接赋值为零值
	// 比如: v 的类型是 {username string, password string}, 但是传入的 json 是 {username: david, pw: 778899},
	// 那么 pw 字段会被忽略, 不会报错, v 的解析值是 {username: david, password: ""}
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
