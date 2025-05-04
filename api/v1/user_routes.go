package v1

import (
	"net/http"
	"skymates-api/internal/handler"
	"skymates-api/internal/service"
)

// RegisterRoutes 注册V1版本的所有 User API 路由
func registerUserRoutes(mux *http.ServeMux, userService service.UserService) {
	userHandler := handler.NewUserHandler(userService)

	// 公开路由
	mux.HandleFunc("POST /api/v1/users/login", userHandler.Login)
	mux.HandleFunc("POST /api/v1/users/register", userHandler.Register)
}
