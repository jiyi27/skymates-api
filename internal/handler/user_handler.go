package handler

import (
	"log"
	"net/http"
	servererrors "skymates-api/errors"
	"skymates-api/internal/dto/v1"
	"skymates-api/internal/service"
	"skymates-api/internal/validator"
)

// UserHandler 处理用户相关的 HTTP 请求
type UserHandler struct {
	BaseHandler
	userService service.UserService
}

// NewUserHandler 创建 UserHandler 实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: BaseHandler{},
		userService: userService,
	}
}

// Register 处理用户注册请求
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req v1.RegisterDto
	if err := h.ReadJSON(r, &req); err != nil {
		h.SendJSON(w, http.StatusBadRequest, "无效的请求格式，可能传递了未知字段", nil)
		return
	}

	// 验证结构体
	if errMsg, err := validator.ValidateRequest(req); err != nil {
		h.SendJSON(w, http.StatusBadRequest, errMsg, nil)
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		// 统一用 HTTPStatus 映射
		status := servererrors.HTTPStatus(err)
		// 5xx 记日志
		if status >= 500 {
			log.Printf("UserHandler.Register: %v", err)
		}
		h.SendJSON(w, status, err.Error(), nil)
		return
	}

	h.SendJSON(w, http.StatusCreated, "用户创建成功", user)
}

// Login 处理用户登录请求
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req v1.LoginDto
	if err := h.ReadJSON(r, &req); err != nil {
		h.SendJSON(w, http.StatusBadRequest, "无效的请求格式", nil)
		return
	}

	// 验证结构体
	if errMsg, err := validator.ValidateRequest(req); err != nil {
		h.SendJSON(w, http.StatusBadRequest, errMsg, nil)
		return
	}

	user, token, err := h.userService.Login(req)
	if err != nil {
		status := servererrors.HTTPStatus(err)
		if status >= 500 {
			log.Printf("UserHandler.Login: %v", err)
		}
		h.SendJSON(w, status, err.Error(), nil)
		return
	}

	//// 调用 http.SetCookie 函数，只是设置响应头部中的 Set-Cookie 字段，还没有真的把响应发送出去
	//// 当你把 Cookie 的 SameSite 设置为 None（表示可以跨站发送）时，必须把 Secure 设置为 true
	//// 否则，现代浏览器（例如 Chrome）会拒绝这个 Cookie，不会保存
	//http.SetCookie(w, &http.Cookie{
	//	Name:     "token",
	//	Value:    token,
	//	Expires:  time.Now().Add(24 * time.Hour),
	//	HttpOnly: true,                  // 禁止客户端通过 js 访问 cookie
	//	Secure:   true,                  // 仅在 https 下发送 cookie
	//	SameSite: http.SameSiteNoneMode, // 允许跨域发送 cookie
	//})

	data := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	h.SendJSON(w, http.StatusOK, "登录成功", data)
}
