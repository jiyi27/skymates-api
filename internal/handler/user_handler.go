package handler

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"skymates-api/internal/auth"
	servererrors "skymates-api/internal/errors"
	"skymates-api/internal/middleware"
	"skymates-api/internal/repositories"
	"skymates-api/internal/types"
)

type UserHandler struct {
	BaseHandler
	userRepo repositories.UserRepository
}

func RegisterUserRoutes(us repositories.UserRepository, mux *http.ServeMux) {
	h := &UserHandler{userRepo: us}
	mux.HandleFunc("/api/auth/register", middleware.Use(h.Register, middleware.Logger, middleware.CORS(nil)))
	mux.HandleFunc("/api/auth/login", middleware.Use(h.Login, middleware.Logger, middleware.CORS(nil)))
	mux.HandleFunc("/api/users/{id}", middleware.Use(h.GetUser, middleware.Logger, middleware.CORS(nil), middleware.Auth))
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid request format", nil)
		return
	}

	// Basic validation
	if len(req.Username) < 3 {
		h.ResponseJSON(w, http.StatusBadRequest, "username must be at least 3 characters", nil)
		return
	}
	if len(req.Password) < 6 {
		h.ResponseJSON(w, http.StatusBadRequest, "password must be at least 6 characters", nil)
		return
	}
	if !isValidEmail(req.Email) {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid email", nil)
		return
	}

	// Check if username exists
	exists, err := h.userRepo.CheckUsernameExists(req.Username)
	if err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
		log.Printf("handler.CreateUser: failed to check username exists: %v", err)
		return
	}
	if exists {
		h.ResponseJSON(w, http.StatusConflict, "username already exists", nil)
		return
	}

	// Check if email exists
	exists, err = h.userRepo.CheckEmailExists(req.Email)
	if err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
		log.Printf("handler.CreateUser: failed to check email exists: %v", err)
		return
	}
	if exists {
		h.ResponseJSON(w, http.StatusConflict, "email already exists", nil)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
		log.Printf("handler.CreateUser: failed to hash password: %v", err)
		return
	}

	// Create user
	user := &types.User{
		Username:       req.Username,
		HashedPassword: string(hashedPassword),
		Email:          req.Email,
	}

	if err := h.userRepo.Create(user); err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "internal server error", nil)
		log.Printf("handler.CreateUser: failed to create user: %v", err)
		return
	}

	h.ResponseJSON(w, http.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := h.DecodeJSON(r, &req); err != nil {
		h.ResponseJSON(w, http.StatusBadRequest, "invalid request format", nil)
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		var serverErr *servererrors.ServerError
		if errors.As(err, &serverErr) {
			fmt.Printf("handler.Login: failed to get user by username: %v", err)
			switch serverErr.Kind {
			case servererrors.KindNotFound:
				h.ResponseJSON(w, http.StatusNotFound, "User not found", nil)
			default:
				h.ResponseJSON(w, http.StatusInternalServerError, "Internal server error", nil)
				log.Printf("handler.Login: failed to get user by username: %v", err)
			}
		}

		h.ResponseJSON(w, http.StatusInternalServerError, "Internal server error", nil)
		log.Printf("handler.Login: failed to get user by username: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		h.ResponseJSON(w, http.StatusUnauthorized, "invalid credentials", nil)
		return
	}

	jwtToken, err := auth.GenerateJwtToken(user)
	if err != nil {
		h.ResponseJSON(w, http.StatusInternalServerError, "Internal server error", nil)
		log.Printf("handler.Login: failed to generate jwt token: %v", err)
		return
	}

	h.ResponseJSON(w, http.StatusOK, "Login successful", map[string]string{"token": jwtToken})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// 实现获取用户逻辑
}
