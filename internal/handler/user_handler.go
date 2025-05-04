package handler

import (
	"log"
	"net/http"
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
		h.SendJSON(w, http.StatusBadRequest, "无效的请求格式, 可能传递了未知字段", nil)
		return
	}

	// 验证结构体
	if errMsg, err := validator.ValidateRequest(req); err != nil {
		h.SendJSON(w, http.StatusBadRequest, errMsg, nil)
		return
	}

	return

	// 调用 service 处理业务逻辑
	user, err := h.userService.Register(req)
	if err != nil {
		// 根据错误类型返回不同的状态码
		switch err.Error() {
		case "username already exists", "email already exists":
			h.SendJSON(w, http.StatusConflict, err.Error(), nil)
		default:
			h.SendJSON(w, http.StatusInternalServerError, "服务器内部错误", nil)
			log.Printf("UserHandler.Register: %v", err)
		}
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

	// 调用 service 处理业务逻辑
	user, token, err := h.userService.Login(req)
	if err != nil {
		// 根据错误类型返回不同的状态码
		switch err.Error() {
		case "user not found":
			h.SendJSON(w, http.StatusNotFound, "用户不存在", nil)
		case "invalid credentials":
			h.SendJSON(w, http.StatusUnauthorized, "无效的凭证", nil)
		default:
			h.SendJSON(w, http.StatusInternalServerError, "服务器内部错误", nil)
			log.Printf("UserHandler.Login: %v", err)
		}
		return
	}

	data := map[string]interface{}{"token": token, "user": user}
	h.SendJSON(w, http.StatusOK, "登录成功", data)

	//// 添加响应头 Set-Cookie 字段, 还未发送, 只是设置了响应头
	//// If you are setting SameSite=None it must always be Secure = True.
	//// If you do not set Secure, the cookie will be rejected by the browser.
	//http.SetCookie(w, &http.Cookie{
	//	Name:     "token",
	//	Value:    jwtToken,
	//	Expires:  time.Now().Add(24 * time.Hour),
	//	HttpOnly: true,                  // 禁止客户端通过 js 访问 cookie
	//	Secure:   true,                  // 仅在 https 下发送 cookie
	//	SameSite: http.SameSiteNoneMode, // 允许跨域发送 cookie
	//})
}
