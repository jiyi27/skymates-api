package handler

import (
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
		h.ResponseJSON(w, http.StatusBadRequest, nil)
		return
	}

	// 实现注册逻辑
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, nil)
		return
	}
	// 实现登录逻辑
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// 实现获取用户逻辑
}
