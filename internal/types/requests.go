package types

import "github.com/google/uuid"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ListTermsRequest struct {
	LastID *uuid.UUID `json:"last_id,omitempty"` // 上次加载的最后一个ID
	Limit  int        `json:"limit"`             // 每次加载数量
}

type CreatePostRequest struct {
	Content string `json:"content"`
}
