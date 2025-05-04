package v1

import (
	"errors"
	"log"
	"net/http"
	serverErrors "skymates-api/errors"
	v1 "skymates-api/internal/dto/v1"
	"skymates-api/internal/service"
	"strings"
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
	user, err := h.userService.Register(req.Username, req.Password, req.Email)
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

	user, token, err := h.userService.Login(req.Email, req.Password)
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

// GetUser 获取用户信息
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ResponseJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	// 从URL路径中获取用户ID
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	userID := path[3]
	if userID == "" {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid user id", nil)
		return
	}

	// 从上下文获取认证用户名，确保用户有权限访问该资源
	username := r.Context().Value("username")
	if username == nil {
		h.ResponseJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// 调用服务获取用户
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		var serverErr *serverErrors.ServerError
		if errors.As(err, &serverErr) && serverErr.Kind == serverErrors.KindNotFound {
			h.ResponseJSON(w, http.StatusNotFound, "user not found", nil)
			return
		}

		h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
		log.Printf("UserHandler.GetUser: %v", err)
		return
	}

	h.ResponseJSON(w, http.StatusOK, "user found", user)
}
