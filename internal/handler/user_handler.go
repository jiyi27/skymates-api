package handler

import (
	"errors"
	"log"
	"net/http"
	serverErrors "skymates-api/errors"
	v1 "skymates-api/internal/dto/v1"
	"skymates-api/internal/service"
	"skymates-api/internal/validator"
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
	var registerDto v1.RegisterDto
	if err := h.DecodeJSON(r, &registerDto); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid request format", nil)
		return
	}

	// 基本验证
	msg, err := validator.ValidateRequest(registerDto)
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, msg, nil)
		return
	}

	// 调用服务层注册用户
	user, err := h.userService.Register(registerDto)
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
	var loginDto v1.LoginDto
	if err := h.DecodeJSON(r, &loginDto); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid request format", nil)
		return
	}

	// 基本验证
	msg, err := validator.ValidateRequest(loginDto)
	if err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, msg, nil)
		return
	}

	user, token, err := h.userService.Login(loginDto)
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
