package handler

import (
	"errors"
	"log"
	"net/http"
	serverErrors "skymates-api/errors"
	v1 "skymates-api/internal/dto/v1"
	"skymates-api/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	BaseHandler
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 处理用户注册
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ResponseJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	var req v1.RegisterDto
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid request format", nil)
		return
	}

	// 基本验证
	if len(req.Username) < 3 {
		h.ResponseJSON(w, http.StatusBadRequest, "username must be at least 3 characters", nil)
		return
	}
	if len(req.Password) < 6 {
		h.ResponseJSON(w, http.StatusBadRequest, "password must be at least 6 characters", nil)
		return
	}
	if !h.ValidateEmail(req.Email) {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid email", nil)
		return
	}

	// 调用服务层注册用户
	user, err := h.userService.Register(req)
	if err != nil {
		var serverErr *serverErrors.ServerError
		if errors.As(err, &serverErr) {
			switch serverErr.Kind {
			case serverErrors.KindValidation:
				h.ResponseJSON(w, http.StatusConflict, serverErr.Message, nil)
			default:
				h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
				log.Printf("UserHandler.Register: %v", err)
			}
			return
		}

		h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
		log.Printf("UserHandler.Register: %v", err)
		return
	}

	h.ResponseJSON(w, http.StatusCreated, "user created successfully", user)
}

// Login 处理用户登录
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ResponseJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	var req v1.LoginDto
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid request format", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		h.ResponseJSON(w, http.StatusBadRequest, "email and password are required", nil)
		return
	}

	user, token, err := h.userService.Login(req)
	if err != nil {
		var serverErr *serverErrors.ServerError
		if errors.As(err, &serverErr) {
			switch serverErr.Kind {
			case serverErrors.KindNotFound:
				h.ResponseJSON(w, http.StatusNotFound, "user not found", nil)
			case serverErrors.KindUnauthorized:
				h.ResponseJSON(w, http.StatusUnauthorized, "invalid credentials", nil)
			default:
				h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
				log.Printf("UserHandler.Login: %v", err)
			}
			return
		}

		h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
		log.Printf("UserHandler.Login: %v", err)
		return
	}

	data := map[string]interface{}{"token": token, "user": user}
	h.ResponseJSON(w, http.StatusOK, "login successful", data)
}
