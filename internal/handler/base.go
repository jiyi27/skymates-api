package handler

import (
	"encoding/json"
	"net/http"
	"skymates-api/internal/types"
)

type BaseHandler struct{}

func (h *BaseHandler) ResponseJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(types.Response{
		Code: code,
		Data: data,
	})
}

func (h *BaseHandler) DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
