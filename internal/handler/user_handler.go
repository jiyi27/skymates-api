package handler

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"skymates-api/internal/repositories"
	"skymates-api/internal/types"
)

type UserHandler struct {
	BaseHandler
	userRepo repositories.UserRepository
}

func NewUserHandler(us repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: us}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/register", h.Register)
	mux.HandleFunc("/api/auth/login", h.Login)
	mux.HandleFunc("/api/users/{id}", h.GetUser)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "请求格式错误", nil)
		return
	}

	// Basic validation
	if len(req.Username) < 3 {
		h.ResponseJSON(w, http.StatusBadRequest, "用户名至少三个字母", nil)
		return
	}
	if len(req.Password) < 6 {
		h.ResponseJSON(w, http.StatusBadRequest, "密码长度至少为6", nil)
		return
	}
	if !isValidEmail(req.Email) {
		h.ResponseJSON(w, http.StatusBadRequest, "邮箱格式错误", nil)
		return
	}

	// Check if username exists
	exists, err := h.userRepo.CheckUsernameExists(req.Username)
	if err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "服务器内部发生错误", nil)
		return
	}
	if exists {
		h.ResponseJSON(w, http.StatusConflict, "用户名已存在", nil)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "服务器内部发生错误", nil)
		return
	}

	// Create user
	user := &types.User{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		Email:          req.Email,
	}

	if err := h.userRepo.Create(user); err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "服务器内部发生错误", nil)
		return
	}

	h.ResponseJSON(w, http.StatusCreated, "用户创建成功", user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "请求格式错误", nil)
		return
	}
	// 实现登录逻辑
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// 实现获取用户逻辑
}
