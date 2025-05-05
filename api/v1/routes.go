package v1

import (
	"net/http"
	"skymates-api/internal/service"
)

// RegisterRoutes 注册V1版本的所有API路由
func RegisterRoutes(mux *http.ServeMux, services *service.Services) {
	registerUserRoutes(mux, services.UserService)
}
