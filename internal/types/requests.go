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
	// omitempty 只对序列化有效, 不会对反序列化产生影响
	LastID *uuid.UUID `json:"last_id"`
	Limit  int        `json:"limit"`
}

type CreatePostRequest struct {
	Content string `json:"content"`
}
