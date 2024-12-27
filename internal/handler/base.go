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
	return json.NewDecoder(r.Body).Decode(v)
}
